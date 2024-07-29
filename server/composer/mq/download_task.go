package mq

import (
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/business"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/repository/mysql/mysql_impl"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/transport/mq"
)

type DownloadTaskMqHandler interface {
	ProcessingDownloadTask() common.MqHandlerFunc
}

func DownloadTaskMq(sctx core.ServiceContext) DownloadTaskMqHandler {
	db := sctx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	repo := mysql_impl.NewDownloadTaskRepositoryImpl(db.GetDB())
	dlBusiness := business.NewDownloadTaskBusiness(sctx, repo)
	downloader := mq.NewDownloader(sctx, dlBusiness)
	return downloader
}
