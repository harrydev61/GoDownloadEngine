package composer

import (
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	grpcComposer "github.com/tranTriDev61/GoDownloadEngine/composer/grpc"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/business"
	userGrpcRepositoryImpl "github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/grpc/grpc_impl"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/mysql/mysql_impl"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/redis/redis_impl"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/transport/api"
	"time"
)

type AuthHandler interface {
	LoginHdl() func(*gin.Context)
	RegisterHdl() func(ctx *gin.Context)
}

func AuthApi(serviceCtx core.ServiceContext) AuthHandler {
	db := serviceCtx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	redisC := serviceCtx.MustGet(common.KeyCompRedis).(common.RedisComponent)
	authMysqlRepository := mysql_impl.NewAuthRepositoryImpl(db.GetDB())
	newRedisTokenRepository := redis_impl.NewTokenRepositoryImpl(redisC.GetClient())
	client, err := grpcComposer.NewUserServiceGrpcClient(grpcComposer.GetAddressDf(serviceCtx))
	if err != nil {
		panic(err)
	}
	userGrpcRepo := userGrpcRepositoryImpl.NewUserGrpcRepository(client)
	authBusiness := business.NewAuthBusiness(serviceCtx, authMysqlRepository, newRedisTokenRepository, userGrpcRepo)
	authApi := api.NewAuthApi(serviceCtx, authBusiness)
	return authApi
}

func GetSetPublicKey(serviceCtx core.ServiceContext, userId string) (*string, error) {
	logger := serviceCtx.Logger("composer api.GetSetPublicKey")
	db := serviceCtx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	redisC := serviceCtx.MustGet(common.KeyCompRedis).(common.RedisComponent)
	authMysqlRepo := mysql_impl.NewAuthRepositoryImpl(db.GetDB())
	authRedisRepo := redis_impl.NewTokenRepositoryImpl(redisC.GetClient())
	publicKey, err := authRedisRepo.GetUserPublicKeyByUserId(userId)
	if err != nil {
		logger.Error("get public key from cache fail", err)
	}
	if publicKey == "" {
		_, publicKey, err := authMysqlRepo.GetKeysByUserId(userId)
		if err != nil {
			logger.Error("get public key from db fail", err)
			return nil, err
		}
		expirationTimePublicKey := time.Now().Add(60 * time.Minute)
		err = authRedisRepo.SetUserPublicKeyByUserId(userId, publicKey, &expirationTimePublicKey)
		if err != nil {
			logger.Error("set public key to cache fail", err)
			return nil, err
		}
		return &publicKey, nil
	}
	return &publicKey, nil
}
