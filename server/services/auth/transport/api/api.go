package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/auth/entity"
	"net/http"
)

type AuthBusiness interface {
	RegisterBiz(ctx context.Context, data entity.AuthRegister, ip string) (*entity.RegisterResponse, error)
	LoginBiz(data entity.AuthLogging) (*entity.UserAccessResponse, error)
}

type authApi struct {
	ServiceCtx core.ServiceContext
	Business   AuthBusiness
}

func NewAuthApi(serviceCtx core.ServiceContext, biz AuthBusiness) *authApi {
	return &authApi{ServiceCtx: serviceCtx, Business: biz}
}

// LoginHandler godoc
//
//	@Summary		User login
//	@Description	User login
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			login 	body		entity.AuthLogging		true	"User login"
//	@Success		200		{object}	common.BaseSuccessResponse
//	@Failure		404		{object}	common.BaseErrorResponse
//	@Failure		422		{object}	common.BaseErrorResponse
//	@Failure		500		{object}	common.BaseErrorResponse
//	@Router			/auth/login [post]
func (h *authApi) LoginHdl() func(*gin.Context) {
	return func(c *gin.Context) {
		var data entity.AuthLogging
		if err := c.ShouldBind(&data); err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		//custom validate
		if err := data.Validate(); err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		result, err := h.Business.LoginBiz(data)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		common.WriteBaseSuccessResponse(c, common.NewBaseSuccessResponse(http.StatusText(http.StatusOK), http.StatusOK, result))
	}
}

// registerHandler godoc
//
//	@Summary		User register
//	@Description	User register
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			register 	body		entity.AuthRegister		true	"User register"
//	@Success		200		{object}	common.BaseSuccessResponse
//	@Failure		404		{object}	common.BaseErrorResponse
//	@Failure		422		{object}	common.BaseErrorResponse
//	@Failure		500		{object}	common.BaseErrorResponse
//	@Router			/auth/register [post]
func (authApi *authApi) RegisterHdl() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var data entity.AuthRegister

		if err := c.ShouldBind(&data); err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		//custom validate
		if err := data.Validate(); err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		ip := c.ClientIP()
		result, err := authApi.Business.RegisterBiz(c.Request.Context(), data, ip)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, entity.ErrCannotRegister)
			return
		}
		common.WriteBaseSuccessResponse(c, common.NewBaseSuccessResponse(http.StatusText(http.StatusCreated), http.StatusCreated, result))
		return
	}
}
