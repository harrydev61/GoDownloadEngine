package composer

import (
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/proto/pb"
	"google.golang.org/grpc"
)

func GetAddressDf(sctx core.ServiceContext) string {
	grpcComponent := sctx.MustGet(common.KeyCompGRPC).(common.GrpcComponent)
	return grpcComponent.GetGRPCServerAddress()
}

func NewUserServiceGrpcClient(address string) (pb.UserServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return pb.NewUserServiceClient(conn), nil
}
