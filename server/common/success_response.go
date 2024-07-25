package common

import (
	"github.com/gin-gonic/gin"
)

type SuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data"`
	Paging  interface{} `json:"paging,omitempty"`
	Extra   interface{} `json:"extra,omitempty"`
}

type BaseSuccessResponse struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data"`
}

func NewSuccessResponse(message string, code int, data interface{}) *SuccessResponse {
	return &SuccessResponse{Status: 1, Message: message, Code: code, Data: data}
}

func NewSuccessFullResponse(status int, message string, code int, data, paging, extra interface{}) *SuccessResponse {
	return &SuccessResponse{Status: status, Message: message, Code: code, Data: data, Paging: paging, Extra: extra}
}

func WriteSuccessResponse(c *gin.Context, response *SuccessResponse) {
	c.JSON(response.Code, response)
}

func WriteBaseSuccessResponse(c *gin.Context, response *BaseSuccessResponse) {
	c.JSON(response.Code, response)
}

func NewBaseSuccessResponse(message string, code int, data interface{}) *BaseSuccessResponse {
	return &BaseSuccessResponse{Status: 1, Message: message, Code: code, Data: data}
}
