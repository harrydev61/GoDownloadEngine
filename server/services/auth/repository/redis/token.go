package redis

import (
	"time"
)

type TokenRepository interface {
	SetRefreshTokenByUserId(userId, refreshToken string, expirationTime *time.Time) error
	GetRefreshTokenByUserId(userId string) (string, error)
	DeleteRefreshTokenByUserId(userId string) (bool, error)
	SetUserPublicKeyByUserId(userId, userPublicKey string, expirationTime *time.Time) error
	GetUserPublicKeyByUserId(userId string) (string, error)
	DeleteUserPublicKeyByUserId(userId string) (bool, error)
}
