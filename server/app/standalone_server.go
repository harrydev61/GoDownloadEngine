package app

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-co-op/gocron/v2"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	gincMiddleware "github.com/tranTriDev61/GoDownloadEngine/component/ginc/middleware"
	grpCcomposer "github.com/tranTriDev61/GoDownloadEngine/composer/grpc"
	jobsCcomposer "github.com/tranTriDev61/GoDownloadEngine/composer/jobs"
	mqCcomposer "github.com/tranTriDev61/GoDownloadEngine/composer/mq"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/dataaccess/mq/consumer"
	"github.com/tranTriDev61/GoDownloadEngine/dataaccess/mq/producer"
	"github.com/tranTriDev61/GoDownloadEngine/proto/pb"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type StandaloneServer struct {
	serviceCtx core.ServiceContext
}

func NewStandaloneServer(serviceCtx *core.ServiceContext) *StandaloneServer {
	return &StandaloneServer{
		serviceCtx: *serviceCtx,
	}
}

func (s *StandaloneServer) Start() error {
	logger := core.GlobalLogger().GetLogger("standalone server")
	// Make some delay for DB ready (migration)
	// remove it if you already had your own DB
	time.Sleep(s.serviceCtx.GetTimeSleep())
	if err := s.serviceCtx.Load(); err != nil {
		logger.Fatal(err)
	}
	s.setupProducerSession()
	s.setupGrpcService()
	s.setupConsumerHandler()
	jobServer, err := s.setupJobsHandler()
	if err != nil {
		_ = jobServer.Shutdown()
		_ = s.serviceCtx.Stop()
		panic(err)
		return err
	}
	server := s.setupGinHttpServer()
	BlockUntilSignal(syscall.SIGINT, syscall.SIGTERM)
	logger.Warnf("Shutdown Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), s.serviceCtx.GetTimeSleep())
	defer cancel()
	if err := jobServer.Shutdown(); err != nil {
		logger.Fatal("Cron job server Shutdown error:", err)
	}
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown　error:", err)
	}
	select {
	case <-ctx.Done():
		logger.Infoln(fmt.Sprintf("timeout of %s seconds.", s.serviceCtx.GetTimeSleep()))
	}
	_ = s.serviceCtx.Stop()
	logger.Info("Server exited")
	return nil
}

func BlockUntilSignal(signals ...os.Signal) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, signals...)
	<-quit
}

func (s *StandaloneServer) setupGinHttpServer() *http.Server {
	logger := core.GlobalLogger().GetLogger("gin.core.services")
	ginComp := s.serviceCtx.MustGet(common.KeyCompGIN).(common.GINComponent)
	router := ginComp.GetRouter()
	// Set up CORS middleware
	router.Use(gincMiddleware.CORSMiddleware(s.serviceCtx))
	router.Use(gin.Recovery(), gin.Logger(), gincMiddleware.RecoveryMiddleware(s.serviceCtx))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group(ginComp.GetPrefix())
	SetupHttpRouter(s.serviceCtx, api)
	server := &http.Server{
		Addr:         ginComp.GetAddr(),
		Handler:      router,
		ReadTimeout:  120 * time.Second,
		WriteTimeout: 120 * time.Second,
	}
	logger.Info(fmt.Sprintf(" Listening and serving HTTP on %s", server.Addr))
	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			logger.Fatalf("listen: %s\n", err)
		}
	}()
	return server

}

func (s *StandaloneServer) setupGrpcService() {
	grpcComponent := s.serviceCtx.MustGet(common.KeyCompGRPC).(common.GrpcComponent)
	logger := s.serviceCtx.Logger(common.KeyCompGRPC)

	lis, err := net.Listen("tcp", grpcComponent.GetGRPCServerAddress())
	if err != nil {
		log.Fatal(err)
	}
	grpcSer := grpcComponent.GetServer()
	pb.RegisterUserServiceServer(grpcSer, grpCcomposer.UserServiceServer(s.serviceCtx))
	logger.Infof("GRPC Server is listening on %d ...\n", grpcComponent.GetGRPCServerAddress())
	go func() {
		if err := grpcSer.Serve(lis); err != nil {
			logger.Fatal(err)
		}
	}()
}

func (s *StandaloneServer) setupProducerSession() {
	mqConfigComponent := s.serviceCtx.MustGet(common.KeyCompMqConfig).(common.MqComponent)
	logger := s.serviceCtx.Logger(common.KeyCompProducer)
	newProducer, err := producer.NewProducerClient(common.KeyCompProducer, logger, mqConfigComponent.GetMqConfig())
	if err != nil {
		log.Fatal(err)
	}
	s.serviceCtx.SetService(common.KeyCompProducer, newProducer)
}

func (s *StandaloneServer) setupConsumerHandler() {
	logger := s.serviceCtx.Logger("Setup consumer mq")
	mqConfigComponent := s.serviceCtx.MustGet(common.KeyCompMqConfig).(common.MqComponent)
	composerDownloadTask := mqCcomposer.DownloadTaskMq(s.serviceCtx)
	newConsumer, err := consumer.NewConsumer(common.KeyCompConsumer, s.serviceCtx, mqConfigComponent.GetMqConfig())
	if err != nil {
		logger.Panicf(err.Error())
	}
	newConsumer.RegisterHandler(common.DownloadTaskTopic, composerDownloadTask.ProcessingDownloadTask())

	go func() {
		consumerStartErr := newConsumer.Start(context.Background())
		logger.Warnf("message queue consumer stopped: %v", consumerStartErr)
	}()
}

func (sdl_s *StandaloneServer) setupJobsHandler() (gocron.Scheduler, error) {
	logger := sdl_s.serviceCtx.Logger("Setup jobs cron job handler")
	// create a scheduler
	jobsServer, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		return nil, err
	}

	// add a job to the scheduler
	cronJobComposer := jobsCcomposer.TakeDownloadTaskCronJobHdl(sdl_s.serviceCtx)
	j, err := jobsServer.NewJob(
		gocron.CronJob(common.ExecuteAllPendingDownloadTaskSchedule, true),
		gocron.NewTask(func() {
			logger.Info("run execute all pending download task job action ")
			if err := cronJobComposer.ExecuteAllPendingDownloadTask(context.Background()); err != nil {
				logger.Errorf("failed to run execute all pending download task job , error: %v", err)
			}
		}),
	)
	if err != nil {
		logger.Errorf("failed to run job cron job , error: %v", err)
		return nil, err
	}

	b, err := jobsServer.NewJob(
		gocron.CronJob(common.UpdateDownloadingAndFailedDownloadTaskStatusToPendingSchedule, true),
		gocron.NewTask(func() {
			if err := cronJobComposer.UpdateDownloadingAndFailedDownloadTaskStatusToPending(context.Background()); err != nil {
				logger.Errorf("failed to run update download task job , error: %v", err)
			}
		}),
	)
	if err != nil {
		logger.Errorf("failed to run job cron job , error: %v", err)
		return nil, err
	}
	jobsServer.Start()
	logger.Infof("Add cron job id:%s", j.ID())
	logger.Infof("Add cron job id:%s", b.ID())
	return jobsServer, nil

}
