package mq

import (
	"context"
	"encoding/json"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/entity"
)

type downloadBusiness interface {
	ExecuteDownloadTask(ctx context.Context, downloadTaskId string) (*entity.DownloadTask, error)
}

type downloader struct {
	Sctx     core.ServiceContext
	logger   core.Logger
	business downloadBusiness
}

func NewDownloader(s core.ServiceContext, business downloadBusiness) *downloader {
	return &downloader{Sctx: s, logger: s.Logger("download.transport.mq"), business: business}
}

func (d *downloader) ProcessingDownloadTask() common.MqHandlerFunc {
	return func(ctx context.Context, queueName string, payload []byte) error {

		var event entity.DownloadTaskMessageMp
		if err := json.Unmarshal(payload, &event); err != nil {
			return err
		}
		d.logger.Infof("\"download task created event received\" %s, %v", queueName, event)
		_, err := d.business.ExecuteDownloadTask(ctx, event.DownloadTaskId)
		return err
	}

}
