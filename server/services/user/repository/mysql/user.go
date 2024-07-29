package mysql

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/entity"
)

type UserRepository interface {
	GetByUserId(userId string) (*entity.UserEntity, error)
	CreateUser(ctx context.Context, userEntity entity.UserEntity) (*entity.UserEntity, error)
	GetByEmail(email string) (*entity.UserEntity, error)
}
