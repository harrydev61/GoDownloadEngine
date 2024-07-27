package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/download/entity"
	"net/http"
)

type DownloadBusiness interface {
	CreateDownloadTask(ctx context.Context, userId string, task entity.Create) (*entity.DownloadTask, error)
	GetListDownloadTask(ctx context.Context, userId string, page, limit int) ([]entity.DownloadTask, error)
	TenderlyDeleteDownloadTask(ctx context.Context, userId, downloadTaskId string) (*string, error)
	GetDetailDownloadTask(ctx context.Context, userId, downloadTaskId string) (*entity.DownloadTask, error)
}
type downloadTransport struct {
	Sctx     core.ServiceContext
	business DownloadBusiness
	logger   core.Logger
}

func NewDownloadTransport(sctx core.ServiceContext, business DownloadBusiness) *downloadTransport {
	return &downloadTransport{Sctx: sctx, business: business, logger: sctx.Logger("download-transport")}
}

// createDownloadTaskHandler godoc
//
//	@Summary		Create download task
//	@Description	Create download task
//	@Tags			Download task
//	@Accept			json
//	@Produce		json
//	@Param			download-task 	body		entity.Create		true	"Create download task"
//	@Success		200		{object}	common.BaseSuccessResponse
//	@Failure		404		{object}	common.BaseErrorResponse
//	@Failure		422		{object}	common.BaseErrorResponse
//	@Failure		500		{object}	common.BaseErrorResponse
//	@Router			/download-task/create [post]
func (t *downloadTransport) CreateDownloadTask() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var body entity.Create
		if err := c.ShouldBind(&body); err != nil {
			t.logger.Warnf("bind body err: %v", err)
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		if err := body.Validate(); err != nil {
			t.logger.Warnf("validate body err: %v", err)
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		user := c.MustGet("user").(common.AccessTokenSubject)
		result, err := t.business.CreateDownloadTask(c, user.UserId, body)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		common.WriteBaseSuccessResponse(c, common.NewBaseSuccessResponse(http.StatusText(http.StatusCreated), http.StatusCreated, result))

	}
}

// getDetailDownloadTaskHandler godoc
//
//	@Summary		Get detail download task
//	@Description	Get detail download task
//	@Tags			Download task
//	@Accept			json
//	@Produce		json
//	@Param			downloadID 	path		uuid		true	"Get detail download task"
//	@Success		200		{object}	common.BaseSuccessResponse
//	@Failure		404		{object}	common.BaseErrorResponse
//	@Failure		422		{object}	common.BaseErrorResponse
//	@Failure		500		{object}	common.BaseErrorResponse
//	@Router			/download-task/:downloadID [get]
func (t *downloadTransport) GetDetailDownloadTask() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var params entity.GetDetail
		if err := c.ShouldBindUri(&params); err != nil {
			t.logger.Warnf("bind param err: %v", err)
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		if err := params.Validate(); err != nil {
			t.logger.Warnf("validate body err: %v", err)
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		user := c.MustGet("user").(common.AccessTokenSubject)
		result, err := t.business.GetDetailDownloadTask(c, user.UserId, params.DownloadID)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		common.WriteBaseSuccessResponse(c, common.NewBaseSuccessResponse(http.StatusText(http.StatusOK), http.StatusOK, result))

	}
}

// getListDownloadTaskHandler godoc
//
//	@Summary		Get list download task
//	@Description	Get list download task
//	@Tags			Download task
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	common.BaseSuccessResponse
//	@Failure		404		{object}	common.BaseErrorResponse
//	@Failure		422		{object}	common.BaseErrorResponse
//	@Failure		500		{object}	common.BaseErrorResponse
//	@Router			/download-task/list [get]
func (t *downloadTransport) GetListDownloadTask() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var query entity.GetLists
		if err := c.ShouldBindQuery(&query); err != nil {
			t.logger.Warnf("bind query err: %v", err)
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		if err := query.Validate(); err != nil {
			t.logger.Warnf("validate query err: %v", err)
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		if query.Page <= 0 || query.Limit <= 0 {
			query = entity.GetLists{
				Page:  1,
				Limit: 10,
			}
		}
		user := c.MustGet("user").(common.AccessTokenSubject)
		result, err := t.business.GetListDownloadTask(c, user.UserId, query.Page, query.Limit)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		common.WriteBaseSuccessResponse(c, common.NewBaseSuccessResponse(http.StatusText(http.StatusOK), http.StatusOK, result))

	}
}

// tenderlyDeleteDownloadTaskHandler godoc
//
//	@Summary		Tenderly delete download task
//	@Description	Tenderly delete download task
//	@Tags			Download task
//	@Accept			json
//	@Produce		json
//	@Param			downloadID 	path		uuid		true	"Tenderly delete download task"
//	@Success		200		{object}	common.BaseSuccessResponse
//	@Failure		404		{object}	common.BaseErrorResponse
//	@Failure		422		{object}	common.BaseErrorResponse
//	@Failure		500		{object}	common.BaseErrorResponse
//	@Router			/download-task/delete/:downloadID [delete]
func (t *downloadTransport) TenderlyDeleteDownloadTask() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		var params entity.Delete
		if err := c.ShouldBindUri(&params); err != nil {
			t.logger.Warnf("bind param err: %v", err)
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		if err := params.Validate(); err != nil {
			t.logger.Warnf("validate body err: %v", err)
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		user := c.MustGet("user").(common.AccessTokenSubject)
		result, err := t.business.TenderlyDeleteDownloadTask(c, user.UserId, params.DownloadID)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		common.WriteBaseSuccessResponse(c, common.NewBaseSuccessResponse(http.StatusText(http.StatusOK), http.StatusOK, result))
	}
}
