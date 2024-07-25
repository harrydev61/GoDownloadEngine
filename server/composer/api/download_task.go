package composer

import (
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/business"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/repository/mysql/mysql_impl"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/transport/api"
)

type DownloadTaskHdl interface {
	CreateDownloadTask() func(ctx *gin.Context)
}

func DownloadTaskApi(sctx core.ServiceContext) DownloadTaskHdl {
	db := sctx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	producer := sctx.MustGet(common.KeyCompProducer).(common.ProducerComponent)
	fileComponent := sctx.MustGet(common.KeyCompFileClient).(common.FileClientComponent)
	repo := mysql_impl.NewDownloadTaskRepositoryImpl(db.GetDB())
	dtBusiness := business.NewDownloadTaskBusiness(sctx, repo, producer, fileComponent.GetClient())
	apiTransport := api.NewDownloadTransport(sctx, dtBusiness)
	return apiTransport
}
