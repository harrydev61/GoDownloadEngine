package composer

import (
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/business"
	userRepositoryImpl "github.com/tranTriDev61/GoDownloadEngine/services/user/repository/mysql/mysql_impl"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/transport/api"
)

type UserHandler interface {
	GetDetailUser() func(*gin.Context)
}

func UserApi(serviceCtx core.ServiceContext) UserHandler {
	db := serviceCtx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	userRepo := userRepositoryImpl.NewUserRepositoryImpl(db.GetDB())
	authBusiness := business.NewUserBusiness(serviceCtx, userRepo)
	return api.NewUserTransport(serviceCtx, authBusiness)
}
