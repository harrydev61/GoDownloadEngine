package grpc

import (
	"context"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/proto/pb"
)

type AuthServiceGrpcBiz interface {
	CreateUserByEmailAndIp(email, ip string) (*string, error)
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

func (s *userService) CreateUserByEmailAndIp(ctx context.Context, req *pb.CreateUserByEmailAndIpReq) (*pb.CreateUserByEmailAndIpResp, error) {
	result, err := s.business.CreateUserByEmailAndIp(req.Email, req.Ip)
	if err != nil {
		return nil, err
	}
	resp := &pb.CreateUserByEmailAndIpResp{
		UserId: *result,
	}
	return resp, nil
}
