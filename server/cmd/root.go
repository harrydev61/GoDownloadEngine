package cmd

import (
	"context"
	"fmt"
	"github.com/gin-contrib/cors"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/tranTriDev61/GoDownloadEngine/component/filec"
	"github.com/tranTriDev61/GoDownloadEngine/component/gormc"
	"github.com/tranTriDev61/GoDownloadEngine/component/mq"
	"github.com/tranTriDev61/GoDownloadEngine/component/mq/consumer"
	"github.com/tranTriDev61/GoDownloadEngine/component/mq/producer"
	"github.com/tranTriDev61/GoDownloadEngine/component/redisc"
	grpCcomposer "github.com/tranTriDev61/GoDownloadEngine/composer/grpc"
	mqCcomposer "github.com/tranTriDev61/GoDownloadEngine/composer/mq"
	_ "github.com/tranTriDev61/GoDownloadEngine/docs"
	"github.com/tranTriDev61/GoDownloadEngine/proto/pb"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/component/ginc"
	gincMiddleware "github.com/tranTriDev61/GoDownloadEngine/component/ginc/middleware"
	"github.com/tranTriDev61/GoDownloadEngine/component/grpcc"
	"github.com/tranTriDev61/GoDownloadEngine/component/jwtc"
	"github.com/tranTriDev61/GoDownloadEngine/core"
)

func newServiceCtx() core.ServiceContext {
	return core.NewServiceContext(
		core.WithComponent(gormc.NewGormDB(common.KeyCompMySQL, common.KeyCompMySQL)),
		core.WithComponent(redisc.NewRedisC(common.KeyCompRedis)),
		core.WithComponent(jwtc.NewJWT(common.KeyCompJWT)),
		core.WithComponent(ginc.NewGin(common.KeyCompGIN)),
		core.WithComponent(mq.NewConfig(common.KeyCompMqConfig)),
		core.WithComponent(grpcc.NewGrpcConfig(common.KeyCompGRPC)),
		core.WithComponent(filec.NewFileC(common.KeyCompFileClient)),
	)
}

var rootCmd = &cobra.Command{
	Use:   "goDownloadEngine",
	Short: "Start services goDownloadEngine",
	Run: func(cmd *cobra.Command, args []string) {
		serviceCtx := newServiceCtx()
		logger := core.GlobalLogger().GetLogger("services")
		// Make some delay for DB ready (migration)
		// remove it if you already had your own DB
		time.Sleep(serviceCtx.GetTimeSleep())
		if err := serviceCtx.Load(); err != nil {
			logger.Fatal(err)
		}

		go setupGrpcService(serviceCtx)
		setupProducerMq(serviceCtx)
		go setupConsumerHandler(serviceCtx)
		server := setupGinServer(serviceCtx)

		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logger.Warnf("Shutdown Server ...")
		ctx, cancel := context.WithTimeout(context.Background(), serviceCtx.GetTimeSleep())
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Fatal("Server Shutdown:", err)
		}
		select {
		case <-ctx.Done():
			logger.Infoln(fmt.Sprintf("timeout of %s seconds.", serviceCtx.GetTimeSleep()))
		}
		_ = serviceCtx.Stop()
		logger.Info("Server exited")

	},
}

func setupGinServer(serviceCtx core.ServiceContext) *http.Server {
	logger := core.GlobalLogger().GetLogger("gin.core.services")
	ginComp := serviceCtx.MustGet(common.KeyCompGIN).(common.GINComponent)
	router := ginComp.GetRouter()
	// Set up CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // Change to your allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	router.Use(gin.Recovery(), gin.Logger(), gincMiddleware.RecoveryMiddleware(serviceCtx))
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	api := router.Group(ginComp.GetPrefix())
	SetupHttpRouter(serviceCtx, api)
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

func setupGrpcService(serviceCtx core.ServiceContext) {
	grpcComponent := serviceCtx.MustGet(common.KeyCompGRPC).(common.GrpcComponent)
	logger := serviceCtx.Logger(common.KeyCompGRPC)

	lis, err := net.Listen("tcp", grpcComponent.GetGRPCServerAddress())
	if err != nil {
		log.Fatal(err)
	}
	s := grpcComponent.GetServer()
	pb.RegisterUserServiceServer(s, grpCcomposer.UserServiceServer(serviceCtx))
	logger.Infof("GRPC Server is listening on %d ...\n", grpcComponent.GetGRPCServerAddress())
	if err := s.Serve(lis); err != nil {
		log.Fatalln(err)
	}
}

func setupConsumerHandler(sctx core.ServiceContext) {
	ctx := context.Background()
	logger := sctx.Logger("Setup consumer mq")
	mqConfigComponent := sctx.MustGet(common.KeyCompMqConfig).(common.MqComponent)
	composerDownloadTask := mqCcomposer.DownloadTaskMq(sctx)
	consumerComponent, err := consumer.NewConsumerC(common.KeyCompConsumer, sctx, mqConfigComponent.GetMqConfig())
	if err != nil {
		logger.Panicf(err.Error())
	}
	consumerComponent.RegisterHandler(common.DownloadTaskTopic, composerDownloadTask.ProcessingDownloadTask())
	consumerComponent.Start(ctx)
}

func setupProducerMq(sctx core.ServiceContext) {
	logger := sctx.Logger("Setup producer mq")
	mqConfigComponent := sctx.MustGet(common.KeyCompMqConfig).(common.MqComponent)
	producerComponent, err := producer.NewProducerC(common.KeyCompProducer, sctx, mqConfigComponent.GetMqConfig())
	if err != nil {
		logger.Panicf(err.Error())
	}
	if err := sctx.AddComponent(producerComponent); err != nil {
		logger.Panicf(err.Error())
	}
	logger.Infof("Started producer mq with key: %s", common.KeyCompProducer)

}
func Execute() {
	rootCmd.AddCommand(outEnvCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
