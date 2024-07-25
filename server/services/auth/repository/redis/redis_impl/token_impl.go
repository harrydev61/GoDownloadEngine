package redis_impl

import (
	"strings"
	"time"

	"github.com/go-redis/redis"
)

const (
	userRefreshTokenPrefix  = "user-auth-refresh-token-{{USERID}}"
	userAuthPublicKeyPrefix = "user-auth-public-key-{{USERID}}"
)

type tokenRepositoryImpl struct {
	redisClient *redis.Client
}

func genUserRefreshTokenPrefix(userID string) string {
	return strings.Replace(userRefreshTokenPrefix, "{{USERID}}", userID, -1)
}
func genUserAuthPublicKeyPrefix(userID string) string {
	return strings.Replace(userAuthPublicKeyPrefix, "{{USERID}}", userID, -1)
}
func NewTokenRepositoryImpl(redisClient *redis.Client) *tokenRepositoryImpl {
	return &tokenRepositoryImpl{
		redisClient: redisClient,
	}
}

func (r *tokenRepositoryImpl) SetRefreshTokenByUserId(userId, refreshToken string, expirationTime *time.Time) error {
	//gen ttl time form expiration time default
	expirationDuration := expirationTime.Sub(time.Now())
	expirationMilliseconds := expirationDuration.Milliseconds()
	// gen key
	key := genUserRefreshTokenPrefix(userId)
	//save data refreshToken to redis cache
	statusCmd := r.redisClient.Set(key, refreshToken, time.Duration(expirationMilliseconds)*time.Millisecond)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}

func (r *tokenRepositoryImpl) GetRefreshTokenByUserId(userId string) (string, error) {
	key := genUserRefreshTokenPrefix(userId)
	// Find refresh token form cache
	result, err := r.redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}

func (r *tokenRepositoryImpl) DeleteRefreshTokenByUserId(userId string) (bool, error) {
	key := genUserRefreshTokenPrefix(userId)
	// Delete refresh token form cache
	statusCmd := r.redisClient.Del(key)
	if statusCmd.Err() != nil {
		return false, statusCmd.Err()
	}
	return true, nil
}
func (r *tokenRepositoryImpl) DeleteUserPublicKeyByUserId(userId string) (bool, error) {
	key := genUserAuthPublicKeyPrefix(userId)
	// Delete refresh token form cache
	statusCmd := r.redisClient.Del(key)
	if statusCmd.Err() != nil {
		return false, statusCmd.Err()
	}
	return true, nil
}

func (r *tokenRepositoryImpl) SetUserPublicKeyByUserId(userId, userPublicKey string, expirationTime *time.Time) error {
	//gen ttl time form expiration time default
	expirationDuration := expirationTime.Sub(time.Now())
	expirationMilliseconds := expirationDuration.Milliseconds()
	// gen key
	key := genUserAuthPublicKeyPrefix(userId)
	//save data refreshToken to redis cache
	statusCmd := r.redisClient.Set(key, userPublicKey, time.Duration(expirationMilliseconds)*time.Millisecond)
	if statusCmd.Err() != nil {
		return statusCmd.Err()
	}
	return nil
}
func (r *tokenRepositoryImpl) GetUserPublicKeyByUserId(userId string) (string, error) {
	key := genUserAuthPublicKeyPrefix(userId)
	// Find refresh token form cache
	result, err := r.redisClient.Get(key).Result()
	if err != nil {
		return "", err
	}
	return result, nil
}
