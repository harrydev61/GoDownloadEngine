package downloader

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/tranTriDev61/GoDownloadEngine/app"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/component/filec"
	"github.com/tranTriDev61/GoDownloadEngine/component/ginc"
	"github.com/tranTriDev61/GoDownloadEngine/component/gormc"
	"github.com/tranTriDev61/GoDownloadEngine/component/grpcc"
	"github.com/tranTriDev61/GoDownloadEngine/component/jwtc"
	"github.com/tranTriDev61/GoDownloadEngine/component/mq"
	"github.com/tranTriDev61/GoDownloadEngine/component/redisc"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	_ "github.com/tranTriDev61/GoDownloadEngine/docs"
	"log"
)

var (
	version string
)

func newServiceCtx() core.ServiceContext {
	return core.NewServiceContext(
		core.WithComponent(gormc.NewGormDB(common.KeyCompMySQL, common.KeyCompMySQL)),
		core.WithComponent(redisc.NewRedisC(common.KeyCompRedis)),
		core.WithComponent(jwtc.NewJWT(common.KeyCompJWT)),
		core.WithComponent(ginc.NewGin(common.KeyCompGIN)),
		core.WithComponent(mq.NewConfig(common.KeyCompMqConfig)),
		core.WithComponent(grpcc.NewGrpcConfig(common.KeyCompGRPC)),
		core.WithComponent(filec.NewFileC(common.KeyCompFileClient)),
	)
}

func server() *cobra.Command {
	command := &cobra.Command{
		Use:   "standalone-server",
		Short: "Start services goDownloadEngine - grpc, http server + kafka consumer + cronjob",
		Run: func(cmd *cobra.Command, args []string) {
			newSer := newServiceCtx()
			err := app.NewStandaloneServer(&newSer).Start()
			if err != nil {
				panic(err)
			}
		},
	}
	return command
}

func Execute() {
	rootCommand := &cobra.Command{
		Version: fmt.Sprintf("%s", version),
	}

	rootCommand.AddCommand(outEnvCmd)
	rootCommand.AddCommand(
		server(),
	)

	if err := rootCommand.Execute(); err != nil {
		log.Panic(err)
	}
}
