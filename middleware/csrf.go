package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"math/big"
)

const (
	CSRF_SECRET_LENGTH = 32
	CSRF_TOKEN_LENGTH  = 2 * CSRF_SECRET_LENGTH
	CSRF_ALLOWED_CHARS = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func generateCSRFSecret() string {
	secret := make([]byte, CSRF_SECRET_LENGTH)

	_, err := rand.Read(secret)
	if err != nil {
		panic(err)
	}

	encodedSecret := base64.URLEncoding.EncodeToString(secret)

	return encodedSecret
}

func GenerateCSRFToken(secret string) string {
	token := make([]byte, CSRF_TOKEN_LENGTH)

	secretBytes := []byte(secret)

	for i := 0; i < CSRF_TOKEN_LENGTH; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(secretBytes))))
		if err != nil {
			panic(err)
		}
		token[i] = secretBytes[randomIndex.Int64()]
	}

	return string(token)
}
