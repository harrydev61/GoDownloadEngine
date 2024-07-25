package composer

import (
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/proto/pb"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/business"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/repository/mysql/mysql_impl"
	authTspGrpc "github.com/tranTriDev61/GoDownloadEngine/services/user/transport/grpc"
)

func UserServiceServer(serviceCtx core.ServiceContext) pb.UserServiceServer {
	db := serviceCtx.MustGet(common.KeyCompMySQL).(common.GormComponent)
	authMysqlRepository := mysql_impl.NewUserRepositoryImpl(db.GetDB())
	authGrpcBiz := business.NewUserBusiness(serviceCtx, authMysqlRepository)
	return authTspGrpc.NewUserServiceGrpc(serviceCtx, authGrpcBiz)
}
