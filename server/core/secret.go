package core

import (
	"crypto/rand"
	"encoding/base32"
)

type secret struct {
	key string
}

func NewSecret(key string) *secret {
	return &secret{
		key: key,
	}
}

func (s *secret) GenNewSecret32() (string, error) {
	secretLength := 32

	// Generate random bytes
	randomBytes := make([]byte, secretLength)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Combine random bytes and key
	tokenWithKey := append(randomBytes, []byte(s.key)...)

	// Encode the combined bytes to base32
	secret := base32.StdEncoding.EncodeToString(tokenWithKey)

	// Trim the base32 encoded string to 32 characters
	secret = secret[:secretLength]

	return secret, nil
}
