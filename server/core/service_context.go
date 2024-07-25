package core

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	DevEnv = "development"
	StgEnv = "staging"
	PrdEnv = "production"
)

type Component interface {
	ID() string
	InitFlags()
	Activate(ServiceContext) error
	Stop() error
}

type ServiceContext interface {
	Load() error
	MustGet(id string) interface{}
	AddComponent(c Component) error
	Get(id string) (interface{}, bool)
	Logger(prefix string) Logger
	EnvName() string
	GetName() string
	GetIP() string
	GetID() string
	Stop() error
	OutEnv()
	GetTimeSleep() time.Duration
}

type serviceCtx struct {
	name       string
	env        string
	ip         string
	id         string
	components []Component
	store      map[string]Component
	cmdLine    *AppFlagSet
	logger     Logger
	timeSleep  int
}

func NewServiceContext(opts ...Option) ServiceContext {
	sv := &serviceCtx{
		store: make(map[string]Component),
	}
	sv.components = []Component{defaultLogger}

	for _, opt := range opts {
		opt(sv)
	}
	sv.initFlags()

	sv.cmdLine = newFlagSet(sv.name, flag.CommandLine)
	sv.parseFlags()

	sv.logger = defaultLogger.GetLogger(sv.name)

	return sv
}

func (s *serviceCtx) initFlags() {
	flag.StringVar(&s.name, "service_name", "", "Service name")
	flag.StringVar(&s.id, "service_id", "", "Service id")
	flag.StringVar(&s.env, "service_env", DevEnv, "Env for services. Ex: dev | stg | prd")
	flag.StringVar(&s.ip, "service_ip", "127.0.0.1", "ip for services. Ex: 127.0.0.1")
	flag.IntVar(&s.timeSleep, "service_time_sleep", 5, "Time sleep for start server")

	for _, c := range s.components {
		c.InitFlags()
	}
}

func (s *serviceCtx) Get(id string) (interface{}, bool) {
	c, ok := s.store[id]

	if !ok {
		return nil, false
	}

	return c, true
}

func (s *serviceCtx) MustGet(id string) interface{} {
	c, ok := s.Get(id)

	if !ok {
		panic(fmt.Sprintf("can not get %s\n", id))
	}

	return c
}

func (s *serviceCtx) Load() error {
	if len(s.name) < 10 {
		s.logger.Panicf("Service name is required and length must be greater than 10")
	}
	if len(s.id) < 5 {
		s.logger.Panicf("Service id is required and length must be greater than 5")
	}
	s.logger.Infoln("Service context " + s.name + " is loading...")
	for _, c := range s.components {
		if err := c.Activate(s); err != nil {
			return err
		}
	}

	return nil
}

func (s *serviceCtx) Logger(prefix string) Logger {
	return defaultLogger.GetLogger(prefix)
}

func (s *serviceCtx) Stop() error {
	s.logger.Infoln("Stopping services context")
	for i := range s.components {
		if err := s.components[i].Stop(); err != nil {
			return err
		}
	}

	s.logger.Infoln("Service context stopped")

	return nil
}

func (s *serviceCtx) AddComponent(c Component) error {
	if _, ok := s.store[c.ID()]; !ok {
		s.components = append(s.components, c)
		s.store[c.ID()] = c
	}
	return nil
}

func (s *serviceCtx) GetName() string { return s.name }
func (s *serviceCtx) GetID() string   { return s.id }
func (s *serviceCtx) GetIP() string   { return s.ip }
func (s *serviceCtx) EnvName() string { return s.env }
func (s *serviceCtx) OutEnv()         { s.cmdLine.GetSampleEnvs() }

type Option func(*serviceCtx)

func WithComponent(c Component) Option {
	return func(s *serviceCtx) {
		if _, ok := s.store[c.ID()]; !ok {
			s.components = append(s.components, c)
			s.store[c.ID()] = c
		}
	}
}

func (s *serviceCtx) parseFlags() {
	envFile := os.Getenv("ENV_FILE")
	if envFile == "" {
		envFile = ".env"
	}

	_, err := os.Stat(envFile)
	if err == nil {
		err := godotenv.Load(envFile)
		if err != nil {
			log.Fatalf("Loading env(%s): %s", envFile, err.Error())
		}
	} else if envFile != ".env" {
		log.Fatalf("Loading env(%s): %s", envFile, err.Error())
	}

	s.cmdLine.Parse([]string{})
}

func (s *serviceCtx) GetTimeSleep() time.Duration {
	return time.Duration(s.timeSleep) * time.Second
}
