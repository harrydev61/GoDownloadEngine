package mysql

import "github.com/tranTriDev61/GoDownloadEngine/services/user/entity"

type UserRepository interface {
	GetByUserId(userId string) (*entity.UserEntity, error)
	CreateByEmailAndIp(email, ip string) (*entity.UserEntity, error)
	GetByEmail(email string) (*entity.UserEntity, error)
	UpdateStatusByUserId(userId string, status int) (*entity.UserEntity, error)
}
