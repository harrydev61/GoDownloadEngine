package app

import (
	"github.com/gin-gonic/gin"
	composer "github.com/tranTriDev61/GoDownloadEngine/composer/api"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	middleware "github.com/tranTriDev61/GoDownloadEngine/middlewares/api"
)

func SetupHttpRouter(serviceCtx core.ServiceContext, router *gin.RouterGroup) {
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"ping": "pong"})
	})
	composerAuthApi := composer.AuthApi(serviceCtx)
	sourceGroup := router.Group("/auth")
	{
		sourceGroup.POST("/login", composerAuthApi.LoginHdl())
		sourceGroup.POST("/register", composerAuthApi.RegisterHdl())
	}
	authMiddleware := middleware.AuthenticationMiddleware(serviceCtx)
	composerUser := composer.UserApi(serviceCtx)
	userGroup := router.Group("/user", authMiddleware)
	{
		userGroup.GET("/detail", composerUser.GetDetailUser())

	}
	composerDownloadTask := composer.DownloadTaskApi(serviceCtx)
	downloadGroup := router.Group("/download-task", authMiddleware)
	{
		downloadGroup.POST("/create", composerDownloadTask.CreateDownloadTask())
		downloadGroup.GET("/:downloadID", composerDownloadTask.GetDetailDownloadTask())
		downloadGroup.GET("/list", composerDownloadTask.GetListDownloadTask())
		downloadGroup.DELETE("/delete/:downloadID", composerDownloadTask.TenderlyDeleteDownloadTask())
	}

}
