package api

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"github.com/tranTriDev61/GoDownloadEngine/services/user/entity"
	"net/http"
)

type UserBusiness interface {
	GetDetailUser(ctx context.Context, userId string) (*entity.UserDataResponse, error)
}
type downloadTransport struct {
	Sctx     core.ServiceContext
	business UserBusiness
	logger   core.Logger
}

func NewUserTransport(sctx core.ServiceContext, business UserBusiness) *downloadTransport {
	return &downloadTransport{Sctx: sctx, business: business, logger: sctx.Logger("user-transport")}
}

// getDetailUserHandler godoc
//
//	@Summary		Get detail user
//	@Description	Get detail user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	common.BaseSuccessResponse
//	@Failure		404		{object}	common.BaseErrorResponse
//	@Failure		422		{object}	common.BaseErrorResponse
//	@Failure		500		{object}	common.BaseErrorResponse
//	@Router			/user/detail [get]
func (t *downloadTransport) GetDetailUser() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		user := c.MustGet("user").(common.AccessTokenSubject)
		result, err := t.business.GetDetailUser(c, user.UserId)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, entity.ErrUserHasExisted)
			return
		}
		common.WriteBaseSuccessResponse(c, common.NewBaseSuccessResponse(http.StatusText(http.StatusOK), http.StatusOK, result))
	}
}
