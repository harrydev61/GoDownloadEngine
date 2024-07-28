package jobs

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/core"
)

type downloadBusiness interface {
	ExecuteAllPendingDownloadTask(ctx context.Context) error
	UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error
}
type downloader struct {
	Sctx     core.ServiceContext
	logger   core.Logger
	business downloadBusiness
}

func NewCronJobsDownloader(s core.ServiceContext, business downloadBusiness) *downloader {
	return &downloader{Sctx: s, logger: s.Logger("download.transport.cronjob"), business: business}
}

func (d *downloader) ExecuteAllPendingDownloadTask(ctx context.Context) error {
	return d.business.ExecuteAllPendingDownloadTask(ctx)
}

func (d *downloader) UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error {
	return d.business.UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx)
}
