package factory

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/entity"
	userGrpcRepository "github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/grpc"
	authRepo "github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/mysql"
)

type Registration interface {
	Register(ctx context.Context, ip string) (*entity.AuthEntity, error)
}

type GmailPassRegistration struct {
	Data           *entity.AuthRegister
	AuthRepository authRepo.AuthRepository
	UserRepository userGrpcRepository.UserGrpcRepository
}

func (r *GmailPassRegistration) Register(ctx context.Context, ip string) (*entity.AuthEntity, error) {
	exit, _ := r.AuthRepository.GetByEmail(r.Data.Email)
	if exit != nil {
		return nil, entity.ErrEmailHasExisted
	}
	//grpc create new user
	userId, err := r.UserRepository.CreateUserByEmailAndIp(ctx, r.Data.Email, ip, r.Data.FirstName, r.Data.LastName)
	if err != nil {
		return nil, err
	}
	// has pass
	salt, err := common.RandomStr(5)
	if err != nil {
		return nil, err
	}
	newPass, err := common.HashPassword(salt, r.Data.Password)
	if err != nil {
		return nil, err
	}
	// create auth
	newAuth := entity.NewAuthEntity(*userId, r.Data.Email, salt, newPass)
	_, err = r.AuthRepository.CreateNewAuth(newAuth)
	if err != nil {
		return nil, err
	}
	return newAuth, nil

}

func NewRegistration(requestData *entity.AuthRegister, authRepo authRepo.AuthRepository, UserRepository userGrpcRepository.UserGrpcRepository) Registration {

	switch requestData.AuthType {
	case common.AuthTypeEmailPassword:
		return &GmailPassRegistration{
			Data:           requestData,
			AuthRepository: authRepo,
			UserRepository: UserRepository,
		}
	default:
		return nil
	}
}
