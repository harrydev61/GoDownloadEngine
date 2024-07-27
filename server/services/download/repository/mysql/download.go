package mysql

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/entity"
)

type DownloadTaskRepository interface {
	GetByName(name string) (*entity.DownloadTask, error)
	Create(ctx context.Context, downloadTask *entity.DownloadTask) (*entity.DownloadTask, error)
	UpdateStatus(downloadId string, status int) error
	GetByUserIdAndDownloadId(ctx context.Context, userId string, downloadId string) (*entity.DownloadTask, error)
	UpdateStatusDownloadTaskPendingTLoading(ctx context.Context, id string) (bool, *entity.DownloadTask, error)
	GetDownloadTaskByDownloadIdWithXLock(ctx context.Context, id string) (*entity.DownloadTask, error)
	UpdateDownloadTask(ctx context.Context, task entity.DownloadTask) error
	GetListByUserId(ctx context.Context, userId string, page, limit, sortTime int) ([]entity.DownloadTask, error)
	GetCountByUserId(ctx context.Context, userId string) (int64, error)
	GetPendingDownloadTaskIDList(ctx context.Context, status int) ([]entity.DownloadTask, error)
	UpdateDownloadingAndFailedDownloadTaskStatusToPending(ctx context.Context) error
	TenderlyDeleteDownloadTask(ctx context.Context, userId, downloadTaskId string) (*string, error)
}
