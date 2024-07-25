package business

import (
	"context"
	"encoding/base64"
	"encoding/json"
	userGrpcRepository "github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/grpc"
	"time"

	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	authFactory "github.com/tranTriDev61/GoDownloadEngine/services/auth/business/factory"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/entity"
	authRepo "github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/mysql"
	redisRepo "github.com/tranTriDev61/GoDownloadEngine/services/auth/repository/redis"
)

type authBusiness struct {
	Sctx                 core.ServiceContext
	logger               core.Logger
	Repository           authRepo.AuthRepository
	UserRepository       userGrpcRepository.UserGrpcRepository
	RedisRepositoryToken redisRepo.TokenRepository
}

func NewAuthBusiness(
	Sctx core.ServiceContext,
	repo authRepo.AuthRepository,
	RedisRepositoryToken redisRepo.TokenRepository,
	UserRepository userGrpcRepository.UserGrpcRepository,
) *authBusiness {
	return &authBusiness{
		Sctx:                 Sctx,
		logger:               Sctx.Logger("auth-business"),
		Repository:           repo,
		RedisRepositoryToken: RedisRepositoryToken,
		UserRepository:       UserRepository,
	}
}
func (biz *authBusiness) RegisterBiz(ctx context.Context, data entity.AuthRegister, ip string) (*entity.RegisterResponse, error) {
	//New registration
	registration := authFactory.NewRegistration(&data, biz.Repository, biz.UserRepository)
	if registration == nil {
		return nil, entity.ErrAuthTypeIsNotValid
	}
	newAuth, err := registration.Register(ctx, ip)
	if err != nil {
		return nil, err
	}
	return &entity.RegisterResponse{UserId: newAuth.UserId, Email: newAuth.Email, AuthType: data.AuthType}, nil
}

func (biz *authBusiness) genDataNewUserAccess(authInfo *entity.AuthEntity) (*entity.UserAccessResponse, error) {
	//Update privateKey and publicKey to auth
	jwtC := biz.Sctx.MustGet(common.KeyCompJWT).(common.JWTProvider)
	accessTokenSubJsonData, err := json.Marshal(common.AccessTokenSubject{UserId: authInfo.UserId, Role: 1})
	if err != nil {
		biz.logger.Error("Error marshalling subject access token JSON:", err)
		return nil, err
	}
	_, publicKey, accessToken, err := jwtC.GenAccessToken(authInfo.UserId, base64.StdEncoding.EncodeToString(accessTokenSubJsonData))
	if err != nil {
		return nil, err
	}
	refreshTokenSubJsonData, err := json.Marshal(common.RefreshTokenSubject{UserId: authInfo.UserId})
	if err != nil {
		biz.logger.Error("Error marshalling subject refresh token JSON:", err)
		return nil, err
	}
	refreshToken, expirationTime, err := jwtC.GenRefreshToken(authInfo.UserId, base64.StdEncoding.EncodeToString(refreshTokenSubJsonData))
	if err != nil {
		return nil, err
	}
	// save refresh token to cache db
	err = biz.RedisRepositoryToken.SetRefreshTokenByUserId(authInfo.UserId, refreshToken, expirationTime)
	if err != nil {
		return nil, err
	}
	// save private key and public key to db
	err = biz.Repository.UpdatePubPriKeyByUserId(authInfo.UserId, "", publicKey)
	if err != nil {
		return nil, err
	}
	// save publicKey to cache db
	expirationTimePublicKey := time.Now().Add(60 * time.Minute)
	err = biz.RedisRepositoryToken.SetUserPublicKeyByUserId(authInfo.UserId, publicKey, &expirationTimePublicKey)
	if err != nil {
		return nil, err
	}
	res := entity.UserAccessResponse{
		Auth: entity.AuthRes{
			UserId: authInfo.UserId,
			Email:  authInfo.Email,
		},
		Token: entity.TokenResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
	return &res, nil
}

func (biz *authBusiness) LoginBiz(data entity.AuthLogging) (*entity.UserAccessResponse, error) {
	//New logging
	logging := authFactory.NewLogging(&data, biz.Repository)
	if logging == nil {
		return nil, entity.ErrAuthTypeValidation
	}
	user, err := logging.Logging()
	if err != nil {
		return nil, err
	}
	res, err := biz.genDataNewUserAccess(user)
	if err != nil {
		return nil, err
	}
	return res, nil
}
