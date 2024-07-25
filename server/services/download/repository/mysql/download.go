package mysql

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/entity"
)

type DownloadTaskRepository interface {
	GetByName(name string) (*entity.DownloadTask, error)
	Create(ctx context.Context, downloadTask *entity.DownloadTask) (*entity.DownloadTask, error)
	UpdateStatus(downloadId string, status int) error
	GetByUserIdAndDownloadId(userId int, downloadId int) (*entity.DownloadTask, error)
	UpdateStatusDownloadTaskPendingTLoading(ctx context.Context, id string) (bool, *entity.DownloadTask, error)
	GetDownloadTaskByDownloadIdWithXLock(ctx context.Context, id string) (*entity.DownloadTask, error)
	UpdateDownloadTask(ctx context.Context, task entity.DownloadTask) error
	GetListByUserId(ctx context.Context, userId string, page, limit, sortTime int) ([]entity.DownloadTask, error)
	GetCountByUserId(ctx context.Context, userId string) (int64, error)
}
