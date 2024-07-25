package middleware

import (
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/core"
)

type CanGetStatusCode interface {
	StatusCode() int
}

func RecoveryMiddleware(serviceCtx core.ServiceContext) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				c.Header("Content-Type", "application/json")

				if appErr, ok := err.(CanGetStatusCode); ok {
					c.AbortWithStatusJSON(appErr.StatusCode(), appErr)
				} else {
					// General panic cases
					c.AbortWithStatusJSON(http.StatusInternalServerError, common.ErrInternalServerError)
				}

				serviceCtx.Logger("services").Errorf("%+v \n", err)

				// Must go with gin recovery
				if gin.IsDebugging() {
					panic(err)
				}
			}
		}()
		c.Next()
	}
}
