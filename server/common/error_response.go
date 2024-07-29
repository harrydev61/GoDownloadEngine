package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/core"
)

type BaseErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message,omitempty"`
	Code    int    `json:"code,omitempty"`
}

func NewBaseErrorResponse(message string, code int) *BaseErrorResponse {
	return &BaseErrorResponse{
		Status:  0,
		Message: message,
		Code:    code,
	}
}

func WriteErrorResponseBadRequest(c *gin.Context, err error) {
	if errSt, ok := err.(core.StatusCodeCarrier); ok {
		c.JSON(errSt.StatusCode(), errSt)
		return
	}
	c.JSON(http.StatusBadRequest, NewBaseErrorResponse(err.Error(), http.StatusBadRequest))
}

func WriteErrorResponseUnauthorized(c *gin.Context) {
	c.JSON(http.StatusUnauthorized, NewBaseErrorResponse(http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized))
}
