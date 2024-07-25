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

func (u *UserGrpcRepositoryImpl) CreateUserByEmailAndIp(ctx context.Context, email, ip string) (*string, error) {
	result, err := u.Client.CreateUserByEmailAndIp(ctx, &pb.CreateUserByEmailAndIpReq{Email: email, Ip: ip})
	if err != nil {
		return nil, err
	}
	userId := result.GetUserId()
	return &userId, nil
}
