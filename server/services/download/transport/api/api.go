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
}
type downloadTransport struct {
	Sctx     core.ServiceContext
	business DownloadBusiness
	logger   core.Logger
}

func NewDownloadTransport(sctx core.ServiceContext, business DownloadBusiness) *downloadTransport {
	return &downloadTransport{Sctx: sctx, business: business, logger: sctx.Logger("download-transport")}
}

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
