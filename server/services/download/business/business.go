package business

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/entity"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/repository/mysql"
	"time"
)

const (
	downloadTaskMetadataFieldNameFileName = "file-name"
)

type DownloadTaskBusiness struct {
	Sctx       core.ServiceContext
	logger     core.Logger
	Repository mysql.DownloadTaskRepository
	producer   common.ProducerComponent
	fileClient common.FileClient
}

func NewDownloadTaskBusiness(
	sctx core.ServiceContext,
	repo mysql.DownloadTaskRepository,
	producer common.ProducerComponent,
	fileClient common.FileClient) *DownloadTaskBusiness {
	return &DownloadTaskBusiness{
		Sctx:       sctx,
		logger:     sctx.Logger("DownloadTaskBusiness"),
		Repository: repo,
		producer:   producer,
		fileClient: fileClient,
	}
}

func (b *DownloadTaskBusiness) CreateDownloadTask(ctx context.Context, userId string, task entity.Create) (*entity.DownloadTask, error) {
	newDownloadTask := entity.NewDownloadTask(userId, task.Name, task.DownloadType, task.URL, common.DownloadTaskPending, &task.Description)
	result, err := b.Repository.Create(ctx, newDownloadTask)
	if err != nil {
		b.logger.Error("can't create new download task", err)
		return nil, entity.ErrCannotCreateDownloadTask
	}
	//send download task to mp
	valueBytes, err := json.Marshal(entity.DownloadTaskMessageMp{DownloadTaskId: newDownloadTask.DownloadID})
	if err != nil {
		b.logger.Error("can't marshal download task message", err)
		return nil, entity.ErrCannotCreateDownloadTask
	}
	err = b.producer.Produce(context.Background(), common.DownloadTaskTopic, valueBytes)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (b *DownloadTaskBusiness) ExecuteDownloadTask(ctx context.Context, downloadTaskId string) (*entity.DownloadTask, error) {
	//update download task to
	b.logger.Debugf("start execute download task %s, time now %s", downloadTaskId, time.Now())
	update, downloadTask, err := b.Repository.UpdateStatusDownloadTaskPendingTLoading(ctx, downloadTaskId)
	if err != nil {
		return nil, err
	}
	if !update {
		return nil, nil
	}
	var downloader core.Downloader
	switch downloadTask.DownloadType {
	case common.DownloadTaskTypeHTTP:
		logger := b.Sctx.Logger("Http downloader")
		downloader = core.NewHTTPDownloader(downloadTask.URL, logger)
	default:
		b.logger.Errorf("unsupported download type: %s", downloadTask.DownloadType)
		b.updateDownloadTaskToFailed(ctx, downloadTaskId)
		return nil, err
	}
	fileName := fmt.Sprintf("%s_%s", common.DownloadFileNamePrefix, downloadTask.DownloadID)
	fileWriteCloser, err := b.fileClient.Write(ctx, fileName)
	if err != nil {
		b.logger.Error("failed to get download file writer")
		b.updateDownloadTaskToFailed(ctx, downloadTaskId)
		return nil, err
	}
	defer fileWriteCloser.Close()

	metadata, err := downloader.Download(ctx, fileWriteCloser)
	if err != nil {
		b.logger.Error("failed to download")
		b.updateDownloadTaskToFailed(ctx, downloadTaskId)
		return nil, err
	}
	metadata[downloadTaskMetadataFieldNameFileName] = fileName
	downloadTask.DownloadStatus = common.DownloadTaskSuccess
	downloadTask.Metadata = entity.JSON{
		Data: metadata,
	}
	err = b.Repository.UpdateDownloadTask(ctx, *downloadTask)
	if err != nil {
		b.logger.Error("failed to update download task status to success")
		return nil, err
	}
	b.logger.Debugf("end execute download task %s", downloadTaskId, time.Now())
	return nil, nil
}

func (b *DownloadTaskBusiness) updateDownloadTaskToFailed(ctx context.Context, downloadTaskId string) {
	err := b.Repository.UpdateStatus(downloadTaskId, common.DownloadTaskFailed)
	if err != nil {
		b.logger.Errorf("can't update download task status to failed: %s", err)
	}
}
