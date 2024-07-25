package mysql

import (
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/entity"
)

type AuthRepository interface {
	CreateNewAuth(data *entity.AuthEntity) (*entity.AuthEntity, error)
	UpdatePubPriKeyByUserId(userId, privateKey, publicKey string) error
	GetByUserIdT(userId string) (*entity.AuthEntity, error)
	GetKeysByUserId(userId string) (privateKey, publicKey string, err error)
	GetByEmail(email string) (*entity.AuthEntity, error)
}
