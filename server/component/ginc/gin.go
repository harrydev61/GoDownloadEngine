package ginc

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"time"
)

const (
	defaultPort       = 8537
	defaultIp         = "127.0.0.1"
	defaultPrefixPath = "api/v1"
	defaultSsl        = false
	defaultMode       = "debug"
	releaseMode       = "release"
	stagingMode       = "staging"
)

type Config struct {
	port     int
	ip       string
	prefix   string
	ssl      bool
	protocol string
	ginMode  string
	timeExit int
}

type ginEngine struct {
	*Config
	name   string
	id     string
	logger core.Logger
	router *gin.Engine
}

func NewGin(id string) *ginEngine {
	return &ginEngine{
		Config: new(Config),
		id:     id,
	}
}

func (gs *ginEngine) ID() string {
	return gs.id
}

func (gs *ginEngine) Activate(sv core.ServiceContext) error {
	gs.logger = sv.Logger(gs.id)
	gs.name = sv.GetName()
	gs.timeExit = 0
	if gs.ginMode == releaseMode || gs.ginMode == stagingMode {
		gin.SetMode(gin.ReleaseMode)
		gs.timeExit = 5
	}

	gs.logger.Info("init engine...")
	gs.router = gin.New()
	return nil
}

func (gs *ginEngine) Stop() error {
	return nil
}

func (gs *ginEngine) InitFlags() {
	flag.IntVar(&gs.Config.port, "gin_port", defaultPort, "gin server port. Default "+string(rune(defaultPort)))
	flag.StringVar(&gs.Config.ip, "gin_ip", defaultIp, "gin server port. Default 127.0.0.1")
	flag.StringVar(&gs.Config.prefix, "gin_prefix_path", defaultPrefixPath, "gin server port. Default "+defaultPrefixPath)
	flag.BoolVar(&gs.Config.ssl, "gin_ssl", defaultSsl, "gin server protocol. Default http://")
	flag.StringVar(&gs.Config.ginMode, "gin_mode", defaultMode, "gin mode (debug | release). Default debug")
}

func (gs *ginEngine) GetPort() int {
	return gs.port
}
func (gs *ginEngine) GetIP() string {
	return gs.ip
}

func (gs *ginEngine) GetAddr() string {
	return fmt.Sprintf("%s:%d", gs.GetIP(), gs.GetPort())
}
func (gs *ginEngine) GetPrefix() string {
	return gs.prefix

}

func (gs *ginEngine) GetTimeExit() time.Duration {
	return time.Duration(gs.timeExit) * time.Second

}
func (gs *ginEngine) GetRouter() *gin.Engine {
	return gs.router
}

func (gs *ginEngine) GetProtocol() string {
	var protocol = "http"
	if gs.ssl {
		protocol = "https"
	}
	return protocol
}
