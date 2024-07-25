package mysql_impl

import (
	"context"
	"errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/entity"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DownloadTaskRepositoryImpl struct {
	db *gorm.DB
}

func NewDownloadTaskRepositoryImpl(db *gorm.DB) *DownloadTaskRepositoryImpl {
	return &DownloadTaskRepositoryImpl{
		db: db,
	}
}
func (d *DownloadTaskRepositoryImpl) GetByName(name string) (*entity.DownloadTask, error) {
	var task entity.DownloadTask
	if err := d.db.Table(task.GetTableName()).Where("name=? AND is_deleted = ?", name, common.RecordNotDeleted).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	if &task == nil {
		return nil, core.ErrRecordNotFound
	}
	return &task, nil
}

func (d *DownloadTaskRepositoryImpl) getByNameTx(tx *gorm.DB, name string) (*entity.DownloadTask, error) {
	var task entity.DownloadTask
	if err := tx.Table(task.GetTableName()).Where("name = ?", name).First(&task).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (d *DownloadTaskRepositoryImpl) Create(ctx context.Context, downloadTask *entity.DownloadTask) (*entity.DownloadTask, error) {
	var newTask *entity.DownloadTask

	txErr := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// check exist
		exist, err := d.getByNameTx(tx, downloadTask.Name)
		if err != nil {
			return err
		}
		if exist != nil {
			return errors.New("conflict download task")
		}

		// create new download task
		if err := tx.Table(downloadTask.GetTableName()).Create(&downloadTask).Error; err != nil {
			return err
		}

		// get new task to response
		newTask, err = d.getByNameTx(tx, downloadTask.Name)
		if err != nil {
			return err
		}
		return nil
	})

	if txErr != nil {
		return nil, txErr
	}

	return newTask, nil
}

func (d *DownloadTaskRepositoryImpl) UpdateStatus(downloadId string, status int) error {
	var task entity.DownloadTask
	if err := d.db.Table(task.GetTableName()).
		Where("download_id = ? AND is_deleted = ?", downloadId, common.RecordNotDeleted).
		Updates(map[string]interface{}{"download_status": status}).
		Error; err != nil {
		return err
	}
	return nil
}

func (d *DownloadTaskRepositoryImpl) GetByUserIdAndDownloadId(userId int, downloadId int) (*entity.DownloadTask, error) {
	var downloadTask entity.DownloadTask
	if err := d.db.Table(downloadTask.GetTableName()).
		Where("user_id = ?,download_id = ? AND is_deleted = ?", userId, downloadId, common.RecordNotDeleted).
		First(&downloadTask).Error; err != nil {
		return nil, err
	}
	return &downloadTask, nil
}

func (d *DownloadTaskRepositoryImpl) GetDownloadTaskByDownloadIdWithXLock(ctx context.Context, id string) (*entity.DownloadTask, error) {
	var downloadTask entity.DownloadTask
	err := d.db.WithContext(ctx).
		Table(entity.DownloadTask{}.GetTableName()).
		Where("download_id = ? AND is_deleted = ?", id, common.RecordNotDeleted).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&downloadTask).Error
	if err != nil {
		return nil, err
	}

	return &downloadTask, nil
}

func (d *DownloadTaskRepositoryImpl) UpdateStatusDownloadTaskPendingTLoading(ctx context.Context, id string) (bool, *entity.DownloadTask, error) {
	var updated bool
	var downloadTask entity.DownloadTask
	txErr := d.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		err := tx.Table(entity.DownloadTask{}.GetTableName()).
			Where("download_id = ? AND is_deleted = ?", id, common.RecordNotDeleted).
			Clauses(clause.Locking{Strength: "UPDATE"}).
			First(&downloadTask).Error

		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil
			}
			return err
		}
		if downloadTask.DownloadStatus != common.DownloadTaskPending {
			updated = false
			return nil
		}
		downloadTask.DownloadStatus = common.DownloadTaskPending
		err = tx.Save(&downloadTask).Error
		if err != nil {
			return err
		}
		updated = true
		return nil
	})

	if txErr != nil {
		return false, nil, txErr
	}
	return updated, &downloadTask, nil
}

func (d *DownloadTaskRepositoryImpl) UpdateDownloadTask(ctx context.Context, task entity.DownloadTask) error {
	if err := d.db.WithContext(ctx).
		Table(entity.DownloadTask{}.GetTableName()).
		Where("download_id = ?", task.DownloadID).
		Updates(&task).Error; err != nil {
		return status.Error(codes.Internal, "failed to update download task")
	}
	return nil
}

func (d *DownloadTaskRepositoryImpl) GetListByUserId(ctx context.Context, userId string, page, limit, sortTime int) ([]entity.DownloadTask, error) {

	var tasks []entity.DownloadTask

	// Calculate the offset for pagination
	offset := (page - 1) * limit

	// Build the query
	query := d.db.WithContext(ctx).Table(entity.DownloadTask{}.GetTableName()).Where("user_id = ? AND is_deleted", userId, common.RecordNotDeleted)

	// Apply sorting based on sortTime parameter
	if sortTime == 1 {
		// Sort by CreatedAt ascending
		query = query.Order("created_at asc")
	} else if sortTime == -1 {
		// Sort by CreatedAt descending
		query = query.Order("created_at desc")
	}

	// Fetch the results with pagination
	if err := query.Limit(limit).Offset(offset).Find(&tasks).Error; err != nil {
		return nil, err
	}

	return tasks, nil
}

func (d *DownloadTaskRepositoryImpl) GetCountByUserId(ctx context.Context, userId string) (int64, error) {
	var count int64

	// Use the context for the query
	if err := d.db.WithContext(ctx).Table(entity.DownloadTask{}.GetTableName()).
		Where("user_id = ? AND is_deleted", userId, common.RecordNotDeleted). // Filter by user ID
		Count(&count).Error;                                                  // Count the records
	err != nil {
		return 0, err
	}

	return count, nil
}
