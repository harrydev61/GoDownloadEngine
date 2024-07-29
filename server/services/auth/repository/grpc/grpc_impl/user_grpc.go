package userGrpcRepositoryImpl

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/proto/pb"
)

type UserGrpcRepositoryImpl struct {
	Client pb.UserServiceClient
}

func NewUserGrpcRepository(client pb.UserServiceClient) *UserGrpcRepositoryImpl {
	return &UserGrpcRepositoryImpl{
		Client: client,
	}
}

func (u *UserGrpcRepositoryImpl) CreateUserByEmailAndIp(ctx context.Context, email, ip, firstName, lastname string) (*string, error) {
	result, err := u.Client.CreateUser(ctx, &pb.CreateUserReq{Email: email, Ip: ip, Firstname: firstName, Lastname: lastname})
	if err != nil {
		return nil, err
	}
	userId := result.GetUserId()
	return &userId, nil
}
