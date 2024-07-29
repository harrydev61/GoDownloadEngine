package business

import (
	"context"
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

func (b *UserBusiness) CreateUser(ctx context.Context, data entity.CreateUserRequest) (*string, error) {
	newUser := entity.NewUserEntity(
		data.Firstname,
		data.Lastname,
		data.Email,
		"",
		"",
		3,
		"unknown",
		nil,
		1,
		data.Ip)
	newUser, err := b.repository.CreateUser(ctx, *newUser)
	if err != nil {
		return nil, err
	}
	return &newUser.UserId, nil
}

func (b *UserBusiness) GetDetailUser(ctx context.Context, userId string) (*entity.UserDataResponse, error) {
	user, err := b.repository.GetByUserId(userId)
	if err != nil {
		return nil, err
	}
	var userData entity.UserDataResponse
	if user != nil {
		userData.UserId = user.UserId
		userData.Email = user.Email
		userData.FirstName = user.FirstName
		userData.LastName = user.LastName
		userData.Phone = user.Phone
		userData.Dob = user.Dob
		userData.Avatar = user.Avatar
		userData.Role = user.Role
		userData.Gender = user.Gender
		userData.CreatedAt = user.CreatedAt
		userData.UpdatedAt = user.UpdatedAt
	}
	return &userData, nil

}
