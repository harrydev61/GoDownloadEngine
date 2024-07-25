package mysql_impl

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
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
	err := repo.db.Table(data.GetTableName()).Where("user_id=?", userId).Where("is_deleted=?", 0).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}
	if &data == nil {
		return nil, core.ErrRecordNotFound
	}
	return &data, nil

}

func (repo *UserRepositoryImpl) CreateByEmailAndIp(email, ip string) (*entity.UserEntity, error) {
	exit, err := repo.GetByEmail(email)
	if exit != nil {
		return nil, entity.ErrEmailHasExisted
	}
	timeNow := time.Now().UTC()
	userEntity := entity.NewUserEntity(email, email, email, "", "", 4, "unknown", &timeNow, -1, ip)
	err = repo.db.Table(userEntity.GetTableName()).Create(&userEntity).Error
	if err != nil {
		return nil, err
	}

	return userEntity, nil
}

func (repo *UserRepositoryImpl) GetByEmail(email string) (*entity.UserEntity, error) {
	var data entity.UserEntity
	err := repo.db.Table(data.GetTableName()).Where("email=?", email).Where("is_deleted=?", 0).First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}
	if &data == nil {
		return nil, core.ErrRecordNotFound
	}
	return &data, nil

}
func (repo *UserRepositoryImpl) UpdateStatusByUserId(userId string, status int) (*entity.UserEntity, error) {
	var data entity.UserEntity
	exit, err := repo.GetByUserId(userId)
	fmt.Println(exit, err)
	if err != nil {
		return nil, err
	}
	if exit == nil {
		return nil, errors.New(common.GetCodeText(common.EntityNotExists))
	}
	exit.Status = status
	if err := repo.db.Table(data.GetTableName()).Where("user_id=?", userId).Updates(exit).Error; err != nil {
		return nil, err
	}

	return exit, nil
}
