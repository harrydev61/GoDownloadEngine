package jobs

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/business"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/repository/mysql/mysql_impl"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/transport/jobs"
)

type DownloadTaskCronJobHdl interface {
	ExecuteAllPendingDownloadTask(ctx context.Context) error
	UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error
}

func TakeDownloadTaskCronJobHdl(sctx core.ServiceContext) DownloadTaskCronJobHdl {
	db := sctx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	repo := mysql_impl.NewDownloadTaskRepositoryImpl(db.GetDB())
	dtBusiness := business.NewDownloadTaskBusiness(sctx, repo)
	cronJobTransport := jobs.NewCronJobsDownloader(sctx, dtBusiness)
	return cronJobTransport
}
