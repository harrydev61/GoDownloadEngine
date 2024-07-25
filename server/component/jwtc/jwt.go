package jwtc

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"flag"
	"fmt"
	"time"

	"github.com/tranTriDev61/GoDownloadEngine/core"

	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
)

const (
	defaultSecret                      = "very-important-please-change-it!" // in 32 bytes
	defaultExpireTokenInSeconds        = 60 * 60 * 24 * 7                   // 7d
	defaultExpireRefreshTokenInSeconds = 60 * 60 * 24 * 10                  // 10d
)

var (
	ErrSecretKeyNotValid            = errors.New("secret key must be in 32 bytes")
	ErrSecretRefreshKeyNotValid     = errors.New("Secret Refresh key must be in 32 bytes")
	ErrTokenLifeTimeTooShort        = errors.New("token life time too short")
	ErrRefreshTokenLifeTimeTooShort = errors.New("Refresh token life time too short")
)

type jwtx struct {
	id                          string
	secret                      string
	refreshSecret               string
	expireTokenInSeconds        int
	expireRefreshTokenInSeconds int
}

func NewJWT(id string) *jwtx {
	return &jwtx{id: id}
}

func (j *jwtx) ID() string {
	return j.id
}

func (j *jwtx) InitFlags() {
	flag.IntVar(
		&j.expireTokenInSeconds,
		"jwt_access_token_ttl",
		defaultExpireTokenInSeconds,
		"Number of seconds token will expired",
	)
	flag.StringVar(
		&j.refreshSecret,
		"jwt_refresh_token_secret",
		"",
		"Refresh secret key to sign JWT",
	)
	flag.IntVar(
		&j.expireRefreshTokenInSeconds,
		"jwt_refresh_token_ttl",
		defaultExpireRefreshTokenInSeconds,
		"Number of seconds refresh token will expired",
	)
}

func (j *jwtx) Activate(_ core.ServiceContext) error {
	if len(j.refreshSecret) < 32 {
		return errors.WithStack(ErrSecretRefreshKeyNotValid)
	}

	if j.expireTokenInSeconds <= 60 {
		return errors.WithStack(ErrTokenLifeTimeTooShort)
	}
	if j.expireRefreshTokenInSeconds < 60*60*24 {
		return errors.WithStack(ErrRefreshTokenLifeTimeTooShort)
	}

	return nil
}

func (j *jwtx) Stop() error {
	return nil
}

func (j *jwtx) VerifyAccessToken(accessToken string, publicKey string) (claims *jwt.RegisteredClaims, err error) {
	publicKeyBytes, err := decodePublicKeyFromPEM(publicKey)
	if err != nil {
		return nil, err
	}
	var rc jwt.RegisteredClaims

	// Parse JWT token with RSA public key
	token, err := jwt.ParseWithClaims(accessToken, &rc, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return publicKeyBytes, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract claims
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
func (j *jwtx) VerifyRefreshToken(refreshToken string) (claims *jwt.RegisteredClaims, err error) {

	var rc jwt.RegisteredClaims

	token, err := jwt.ParseWithClaims(refreshToken, &rc, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(j.refreshSecret), nil
	})

	if !token.Valid {
		return nil, errors.WithStack(err)
	}

	return &rc, nil
}

func (j *jwtx) GenAccessToken(userId string, sub string) (string, string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", "", err
	}

	publicKey := &privateKey.PublicKey

	privateKeyPEM := encodePrivateKeyToPEM(privateKey)
	publicKeyPEM := encodePublicKeyToPEM(publicKey)

	now := time.Now().UTC()
	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Second * time.Duration(j.expireTokenInSeconds))),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        userId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	accessToken, err := token.SignedString(privateKey)
	if err != nil {
		return "", "", "", err
	}

	// Trả về chuỗi PEM của cả hai khóa và token
	return privateKeyPEM, publicKeyPEM, accessToken, nil
}
func (j *jwtx) GenRefreshToken(userId, sub string) (string, *time.Time, error) {
	// Create the JWT token
	now := time.Now().UTC()
	expirationTime := now.Add(time.Second * time.Duration(j.expireRefreshTokenInSeconds))
	claims := jwt.RegisteredClaims{
		Subject:   sub,
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		NotBefore: jwt.NewNumericDate(now),
		IssuedAt:  jwt.NewNumericDate(now),
		ID:        userId,
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenSignedStr, err := t.SignedString([]byte(j.refreshSecret))

	if err != nil {
		return "", nil, err
	}
	return tokenSignedStr, &expirationTime, nil
}

func encodePrivateKeyToPEM(privateKey *rsa.PrivateKey) string {
	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	return string(privateKeyPEM)
}

func encodePublicKeyToPEM(publicKey *rsa.PublicKey) string {
	publicKeyBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	})
	return string(publicKeyPEM)
}

func decodePublicKeyFromPEM(publicKeyPEM string) (*rsa.PublicKey, error) {
	// Decode PEM-encoded public key
	block, _ := pem.Decode([]byte(publicKeyPEM))
	if block == nil {
		return nil, errors.New("failed to decode PEM block containing public key")
	}

	// Parse the key
	pubKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	rsaPubKey, ok := pubKey.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("invalid RSA public key")
	}
	return rsaPubKey, nil
}
