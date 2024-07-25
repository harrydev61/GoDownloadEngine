package business

import (
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/entity"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/repository/mysql"
)

type UserBusiness struct {
	ctx        core.ServiceContext
	repository mysql.UserRepository
	logger     core.Logger
}

func NewUserBusiness(ctx core.ServiceContext, repo mysql.UserRepository) *UserBusiness {
	return &UserBusiness{
		ctx:        ctx,
		repository: repo,
		logger:     ctx.Logger("user-business"),
	}
}

func (b *UserBusiness) CreateUserByEmailAndIp(email, ip string) (*string, error) {
	//check user exit
	exit, err := b.repository.GetByEmail(email)
	if exit != nil {
		return nil, err
	}
	if exit != nil {
		return nil, entity.ErrEmailHasExisted
	}
	user, err := b.repository.CreateByEmailAndIp(email, ip)
	if err != nil {
		return nil, err
	}
	return &user.UserId, nil
}
