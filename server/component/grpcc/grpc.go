package grpcc

import (
	"flag"
	"fmt"
	"time"

	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/component/grpcc/interceptors"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"google.golang.org/grpc"
)

const (
	GrpcServerHostDefault = "127.0.0.1"
	GrpcServerPortDefault = 3300
	debugMode             = "debug"
)

type grpcConfig struct {
	id         string
	server     *grpc.Server
	serverHost string
	serverPort int
	serverMode string
}

func NewGrpcConfig(id string) *grpcConfig {
	return &grpcConfig{id: id}
}

func (c *grpcConfig) ID() string {
	return c.id
}

func (c *grpcConfig) InitFlags() {
	flag.StringVar(
		&c.serverHost,
		"grpc_server_host",
		GrpcServerHostDefault,
		"gRPC server host default: 127.0.0.1",
	)
	flag.IntVar(
		&c.serverPort,
		"grpc_server_port",
		GrpcServerPortDefault,
		"gRPC server port default: 3300",
	)
	flag.StringVar(
		&c.serverMode,
		"grpc_server_mode",
		debugMode,
		"gRPC server mode default: debug",
	)
}

func (c *grpcConfig) Activate(serviceCtx core.ServiceContext) error {
	logger := serviceCtx.Logger(common.KeyCompGRPC)
	var serverOptions []grpc.ServerOption
	serverOptions = append(serverOptions, grpc.UnaryInterceptor(interceptors.RecoverInterceptor(logger, c.serverMode == debugMode, 120*time.Second)))
	s := grpc.NewServer(serverOptions...)
	c.server = s
	return nil
}

func (c *grpcConfig) Stop() error {
	c.server.GracefulStop()
	return nil
}

func (c *grpcConfig) GetServer() *grpc.Server {
	return c.server
}

func (c *grpcConfig) GetGRPCPort() int {
	return c.serverPort
}

func (c *grpcConfig) GetGRPCServerAddress() string {
	return fmt.Sprintf("%s:%d", c.serverHost, c.serverPort)
}

func (c *grpcConfig) GetGRPCServerHost() string {
	return c.serverHost
}

func (c *grpcConfig) GetGRPCPathHealthCheck() string {
	return fmt.Sprintf("%s:%s", c.GetGRPCServerAddress(), "/pb.AuthCoreService/HealthCheck")
}
