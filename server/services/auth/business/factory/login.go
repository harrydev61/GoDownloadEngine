package factory

import (
	"github.com/IBM/sarama"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/entity"
	authRepo "github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/mysql"
)

type Logging interface {
	Logging() (*entity.AuthEntity, error)
}

type EmailPasswordLogging struct {
	Data           *entity.AuthLogging
	AuthRepository authRepo.AuthRepository
	producer       *sarama.SyncProducer
}

func (r *EmailPasswordLogging) Logging() (*entity.AuthEntity, error) {
	//Check auth user info
	authUserInfo, err := r.AuthRepository.GetByEmail(r.Data.Email)
	if err != nil {
		return nil, entity.ErrLoginFailed
	}
	if !common.CompareHashPassword(authUserInfo.Password, authUserInfo.Salt, r.Data.Password) {
		return nil, entity.ErrLoginFailed
	}

	return authUserInfo, nil
}

func NewLogging(requestData *entity.AuthLogging, authRepo authRepo.AuthRepository) Logging {
	switch requestData.AuthType {
	case common.AuthTypeEmailPassword:
		return &EmailPasswordLogging{
			Data:           requestData,
			AuthRepository: authRepo,
		}
	default:
		return nil
	}
}
