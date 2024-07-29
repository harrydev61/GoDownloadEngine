package business

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gammazero/workerpool"
	"github.com/pkg/errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/entity"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/repository/mysql"
	"time"
)

type DownloadTaskBusiness struct {
	Sctx       core.ServiceContext
	logger     core.Logger
	Repository mysql.DownloadTaskRepository
}

func NewDownloadTaskBusiness(
	sctx core.ServiceContext,
	repo mysql.DownloadTaskRepository,
) *DownloadTaskBusiness {
	return &DownloadTaskBusiness{
		Sctx:       sctx,
		logger:     sctx.Logger("DownloadTaskBusiness"),
		Repository: repo,
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
	producerService := b.Sctx.GetService(common.KeyCompProducer).(common.ProducerComponent)
	err = producerService.Produce(context.Background(), common.DownloadTaskTopic, valueBytes)
	if err != nil {
		return nil, err
	}

	return result, nil
}
func (b *DownloadTaskBusiness) GetListDownloadTask(ctx context.Context, userId string, page, limit int) ([]entity.DownloadTask, error) {
	downloadTask, err := b.Repository.GetListByUserId(ctx, userId, page, limit, 1)
	if err != nil {
		b.logger.Errorf("can't get list download task by userId: %s, %d,%d, err: %v", userId, page, limit, err)
		return make([]entity.DownloadTask, 0), err
	}
	return downloadTask, nil
}
func (b *DownloadTaskBusiness) TenderlyDeleteDownloadTask(ctx context.Context, userId, downloadTaskId string) (*string, error) {
	deleteId, err := b.Repository.TenderlyDeleteDownloadTask(ctx, userId, downloadTaskId)
	if err != nil {
		b.logger.Errorf("can't delete download task by userId: %s, %s, err: %v", userId, downloadTaskId, err)
		return nil, errors.New("can't delete download task")
	}
	return deleteId, nil
}
func (b *DownloadTaskBusiness) GetDetailDownloadTask(ctx context.Context, userId, downloadTaskId string) (*entity.DownloadTask, error) {
	downloadTask, err := b.Repository.GetByUserIdAndDownloadId(ctx, userId, downloadTaskId)
	if err != nil {
		b.logger.Errorf("can't get detail download task by id %s, err: %v", downloadTaskId, err)
		return nil, err
	}
	return downloadTask, nil

}

func (b *DownloadTaskBusiness) ExecuteDownloadTask(ctx context.Context, downloadTaskId string) (*entity.DownloadTask, error) {
	//update download task to
	b.logger.Infof("start execute download task %s, time now %s", downloadTaskId, time.Now())
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
	fileComponent := b.Sctx.MustGet(common.KeyCompFileClient).(common.FileClientComponent)
	fileWriteCloser, err := fileComponent.GetClient().Write(ctx, fileName)
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
	metadata[common.DownloadTaskMetadataFieldNameFileName] = fileName
	downloadTask.DownloadStatus = common.DownloadTaskSuccess
	downloadTask.Metadata = entity.JSON{
		Data: metadata,
	}
	err = b.Repository.UpdateDownloadTask(ctx, *downloadTask)
	if err != nil {
		b.logger.Error("failed to update download task status to success")
		return nil, err
	}
	b.logger.Infof("end execute download task %s", downloadTaskId, time.Now())
	return nil, nil
}
func (b *DownloadTaskBusiness) updateDownloadTaskToFailed(ctx context.Context, downloadTaskId string) {
	err := b.Repository.UpdateStatus(downloadTaskId, common.DownloadTaskFailed)
	if err != nil {
		b.logger.Errorf("can't update download task status to failed: %s", err)
	}
}
func (d *DownloadTaskBusiness) ExecuteAllPendingDownloadTask(ctx context.Context) error {
	pendingDownloadTaskIDList, err := d.Repository.GetPendingDownloadTaskIDList(ctx, common.DownloadTaskPending)
	if err != nil {
		return err
	}
	if len(pendingDownloadTaskIDList) == 0 {
		d.logger.Info("no pending download task found")
		return nil
	}

	d.logger.Infof("start execute all pending download task tasks %s", len(pendingDownloadTaskIDList))

	workerPool := workerpool.New(common.ExecuteAllPendingDownloadTaskConcurrencyLimit)
	jobLogger := d.Sctx.Logger("ExecuteAllPendingDownloadTask")
	for _, task := range pendingDownloadTaskIDList {
		workerPool.Submit(func() {
			if _, exDownloadTaskErr := d.ExecuteDownloadTask(ctx, task.DownloadID); exDownloadTaskErr != nil {
				jobLogger.Errorf("failed to execute download task %s: %s", task.DownloadID, exDownloadTaskErr)
			}
		})
	}

	workerPool.StopWait()
	return nil
}
func (d *DownloadTaskBusiness) UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error {
	return d.Repository.UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx)
}
