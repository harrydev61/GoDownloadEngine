package mysql_impl

import (
	"context"
	"errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"

	"github.com/tranTriDev61/GoDownloadEngine/services/user/entity"
	"gorm.io/gorm"
)

type UserRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepositoryImpl(db *gorm.DB) *UserRepositoryImpl {
	return &UserRepositoryImpl{
		db: db,
	}
}

func (repo *UserRepositoryImpl) GetByUserId(userId string) (*entity.UserEntity, error) {
	var data entity.UserEntity
	err := repo.db.Table(data.GetTableName()).Where("user_id = ? AND is_deleted = ?", userId, 0).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}
	return &data, nil
}

func getByEmailTx(tx *gorm.DB, email string) (*entity.UserEntity, error) {
	var user entity.UserEntity
	if err := tx.Table(user.GetTableName()).Where("email = ? AND is_deleted = ?", email, common.RecordNotDeleted).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *UserRepositoryImpl) CreateUser(ctx context.Context, userEntity entity.UserEntity) (*entity.UserEntity, error) {
	txErr := repo.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		exist, err := getByEmailTx(tx, userEntity.Email)
		if err != nil {
			return err
		}
		if exist != nil {
			return entity.ErrEmailHasExisted
		}

		if err := tx.Table(userEntity.GetTableName()).Create(&userEntity).Error; err != nil {
			return err
		}

		if err := tx.Table(userEntity.GetTableName()).
			Where("user_id = ? AND is_deleted = ?", userEntity.UserId, common.RecordNotDeleted).
			First(&userEntity).Error; err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return nil, txErr
	}
	return &userEntity, nil
}

func (repo *UserRepositoryImpl) GetByEmail(email string) (*entity.UserEntity, error) {
	var data entity.UserEntity
	err := repo.db.Table(data.GetTableName()).Where("email = ? AND is_deleted = ?", email, common.RecordNotDeleted).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}
	return &data, nil
}
