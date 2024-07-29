package common

import (
	"context"
	"crypto/rsa"
	"github.com/go-redis/redis"
	"google.golang.org/grpc"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

type GINComponent interface {
	GetPort() int
	GetRouter() *gin.Engine
	GetAddr() string
	GetIP() string
	GetPrefix() string
	GetTimeExit() time.Duration
	GetProtocol() string
}

type TokenComponent interface {
	GenPriPubKey() (*rsa.PrivateKey, *rsa.PublicKey, error)
}

type GormComponent interface {
	GetDB() *gorm.DB
}

type RedisComponent interface {
	GetClient() *redis.Client
}

type JWTProvider interface {
	VerifyAccessToken(accessToken string, publicKey string) (claims *jwt.RegisteredClaims, err error)
	VerifyRefreshToken(tokenString string) (claims *jwt.RegisteredClaims, err error)
	GenAccessToken(userId string, sub string) (string, string, string, error)
	GenRefreshToken(userId, sub string) (string, *time.Time, error)
}

type GrpcComponent interface {
	GetServer() *grpc.Server
	GetGRPCPort() int
	GetGRPCServerAddress() string
	GetGRPCPathHealthCheck() string
	GetGRPCServerHost() string
}

type MqComponent interface {
	GetMqConfig() MqConfig
}

type ProducerComponent interface {
	ID() string
	Produce(ctx context.Context, queueName string, payload []byte) error
	Stop() error
}

type FileClientComponent interface {
	GetClient() FileClient
}
