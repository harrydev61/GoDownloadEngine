package mysql_impl

import (
	"github.com/pkg/errors"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/entity"
	"gorm.io/gorm"
)

type authRepositoryImpl struct {
	db *gorm.DB
}

func NewAuthRepositoryImpl(db *gorm.DB) *authRepositoryImpl {
	return &authRepositoryImpl{
		db: db,
	}
}

func (repo *authRepositoryImpl) CreateNewAuth(data *entity.AuthEntity) (*entity.AuthEntity, error) {
	var result entity.AuthEntity
	exit, _ := repo.GetByUserIdT(data.UserId)
	if exit != nil {
		return nil, errors.New(common.GetCodeText(common.EntityNotExists))
	}
	if err := repo.db.Table(result.GetTableName()).Create(&data).Error; err != nil {
		return nil, err
	}
	return data, nil
}

func (repo *authRepositoryImpl) GetByUserIdT(userId string) (*entity.AuthEntity, error) {
	var data entity.AuthEntity

	if err := repo.db.Table(data.GetTableName()).Where("user_id=?", userId).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	if &data == nil {
		return nil, core.ErrRecordNotFound
	}
	return &data, nil

}

func (repo *authRepositoryImpl) GetByEmail(email string) (*entity.AuthEntity, error) {
	var data entity.AuthEntity
	if err := repo.db.Table(data.GetTableName()).Where("email=?", email).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, core.ErrRecordNotFound
		}
		return nil, err
	}

	if &data == nil {
		return nil, core.ErrRecordNotFound
	}
	return &data, nil

}

func (repo *authRepositoryImpl) UpdatePubPriKeyByUserId(userId, privateKey, publicKey string) error {
	var data entity.AuthEntity

	// Update the private_key and public_key columns for the specified user ID
	if err := repo.db.Table(data.GetTableName()).Where("user_id = ?", userId).Updates(map[string]interface{}{
		"private_key": privateKey,
		"public_key":  publicKey,
	}).Error; err != nil {
		return err
	}

	return nil
}

func (repo *authRepositoryImpl) GetKeysByUserId(userId string) (privateKey, publicKey string, err error) {
	var data entity.AuthEntity

	// Select only privateKey and publicKey columns from the database
	if err := repo.db.Table(data.GetTableName()).Select("private_key, public_key").Where("user_id=?", userId).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", "", core.ErrRecordNotFound
		}
		return "", "", err
	}

	// Check if no data found
	if &data == nil {
		return "", "", core.ErrRecordNotFound
	}

	// Return the values of privateKey and publicKey
	return data.PrivateKey, data.PublicKey, nil
}
