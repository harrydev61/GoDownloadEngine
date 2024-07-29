package grpc

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/proto/pb"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/entity"
)

type AuthServiceGrpcBiz interface {
	CreateUser(ctx context.Context, request entity.CreateUserRequest) (*string, error)
}

type userService struct {
	serviceCtx core.ServiceContext
	business   AuthServiceGrpcBiz
}

func NewUserServiceGrpc(serviceCtx core.ServiceContext, biz AuthServiceGrpcBiz) *userService {
	return &userService{
		serviceCtx: serviceCtx,
		business:   biz,
	}
}

func (s *userService) CreateUser(ctx context.Context, req *pb.CreateUserReq) (*pb.CreateUserResp, error) {
	newUser := entity.CreateUserRequest{
		Email:     req.Email,
		Ip:        req.Ip,
		Firstname: req.Firstname,
		Lastname:  req.Lastname,
	}

	userId, err := s.business.CreateUser(ctx, newUser)
	if err != nil {
		return nil, err
	}
	resp := &pb.CreateUserResp{
		UserId: *userId,
	}
	return resp, nil
}
