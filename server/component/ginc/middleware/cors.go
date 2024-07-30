package middleware

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"time"
)

func CORSMiddleware(serviceCtx core.ServiceContext) gin.HandlerFunc {
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"}, // Change to your allowed origins
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
	return cors.New(config)
}
