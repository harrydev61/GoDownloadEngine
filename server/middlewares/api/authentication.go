package middleware

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/tranTriDev61/GoDownloadEngine/common"
	composerApi "github.com/tranTriDev61/GoDownloadEngine/composer/api"
	"github.com/tranTriDev61/GoDownloadEngine/core"
	"strings"
)

func AuthenticationMiddleware(serviceCtx core.ServiceContext) func(*gin.Context) {
	return func(c *gin.Context) {
		logger := serviceCtx.Logger("AuthenticationMiddleware")
		token, err := extractTokenFromHeaderString(c.GetHeader("Authorization"))

		if err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			c.Abort()
			return
		}
		userId := c.GetHeader("X-USER-ID")
		if userId == "" {
			common.WriteErrorResponseBadRequest(c, errors.New("missing user Id"))
			c.Abort()
			return
		}
		publicKey, err := composerApi.GetSetPublicKey(serviceCtx, userId)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			return
		}
		jwtComp := serviceCtx.MustGet(common.KeyCompJWT).(common.JWTProvider)
		claims, err := jwtComp.VerifyAccessToken(token, *publicKey)
		if err != nil {
			common.WriteErrorResponseUnauthorized(c)
			c.Abort()
			return
		}
		decodedData, err := base64.StdEncoding.DecodeString(claims.Subject)
		if err != nil {
			common.WriteErrorResponseBadRequest(c, err)
			c.Abort()
			return
		}
		// Unmarshal JSON data into struct
		var accessTokenSub common.AccessTokenSubject
		err = json.Unmarshal(decodedData, &accessTokenSub)
		if err != nil {
			logger.Error("decode data claims error", err)
			common.WriteErrorResponseBadRequest(c, err)
			c.Abort()
			return
		}
		c.Set("user", accessTokenSub)
		c.Next()
	}
}

func extractTokenFromHeaderString(s string) (string, error) {
	parts := strings.Split(s, " ")
	//"Authorization" : "Bearer {token}"
	if parts[0] != "Bearer" || len(parts) < 2 || strings.TrimSpace(parts[1]) == "" {
		return "", errors.New("missing access token")
	}
	return parts[1], nil
}
