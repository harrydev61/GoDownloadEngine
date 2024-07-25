package redisc

import (
	"flag"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/tranTriDev61/GoDownloadEngine/core"
)

const (
	redisHostDefault    = "127.0.0.1"
	redisPortDefault    = 6379
	redisPassDefault    = ""
	redisDbCacheDefault = 0
)

type redisC struct {
	id      string
	logger  core.Logger
	client  *redis.Client
	secret  string
	host    string
	port    int
	pass    string
	dbCache int
}

func NewRedisC(id string) *redisC {
	return &redisC{id: id}
}

func (r *redisC) ID() string {
	return r.id
}

func (r *redisC) InitFlags() {
	flag.StringVar(
		&r.host,
		"redis_cache_host",
		redisHostDefault,
		"Host for auth connect redis",
	)
	flag.IntVar(
		&r.port,
		"redis_cache_port",
		redisPortDefault,
		"Port for connect redis",
	)
	flag.StringVar(
		&r.pass,
		"redis_cache_auth",
		redisPassDefault,
		"Password for auth connect redis",
	)
	flag.IntVar(
		&r.dbCache,
		"redis_cache_db",
		redisDbCacheDefault,
		"Database name for storage redis",
	)

}

func (r *redisC) Activate(_ core.ServiceContext) error {
	r.logger = core.GlobalLogger().GetLogger(r.id)
	connectOption := redis.Options{
		Addr:     fmt.Sprintf("%s:%d", r.host, r.port), // Redis server address
		Password: r.pass,                               // no password set
		DB:       r.dbCache,                            // use default DB
	}
	client := redis.NewClient(&connectOption)

	// Ping Redis to check if the connection is successful
	_, err := client.Ping().Result()
	if err != nil {
		// Try connecting without password
		connectOption.Password = ""
		client = redis.NewClient(&connectOption)
		_, err = client.Ping().Result()
		if err != nil {
			r.logger.Fatalf("Error connecting to Redis without password: %v", err)
			return err
		}
	}
	r.client = client
	r.logger.Info("Connected to Redis:", fmt.Sprintf("%s:%d", r.host, r.port))
	return nil
}

func (r *redisC) Stop() error {
	err := r.client.Close()
	if err != nil {
		fmt.Println("Error closing connection:", err)
		return err
	}
	return nil
}

func (r *redisC) GetClient() *redis.Client {
	return r.client
}
