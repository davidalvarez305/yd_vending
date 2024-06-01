package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"io"
	"math/big"
	"strconv"
	"strings"

	"github.com/davidalvarez305/budgeting/database"
	"github.com/davidalvarez305/budgeting/models"
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

func generateCSRFToken() string {
	token := make([]byte, CSRF_TOKEN_LENGTH)
	secret := generateCSRFSecret()

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

func Encrypt(unixTime int64, key []byte, db *sql.DB) (string, error) {
	var encryptedString string
	var token = generateCSRFToken()

	csrfToken := models.CSRFToken{
		ExpiryTime: unixTime,
		Token:      token,
	}

	err := database.InsertCSRFToken(csrfToken)
	if err != nil {
		return encryptedString, err
	}

	tokenBytes := []byte(token)
	unixTimeBytes := []byte(strconv.FormatInt(unixTime, 10))

	joinedData := strings.Join([]string{string(tokenBytes), string(unixTimeBytes)}, ":")

	paddedData := pad([]byte(joinedData), aes.BlockSize)

	block, err := aes.NewCipher(key)
	if err != nil {
		return encryptedString, err
	}

	cipherText := make([]byte, aes.BlockSize+len(paddedData))
	iv := cipherText[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return encryptedString, err
	}
	mode := cipher.NewCBCEncrypter(block, iv)

	mode.CryptBlocks(cipherText[aes.BlockSize:], paddedData)

	encryptedString = base64.StdEncoding.EncodeToString(cipherText)

	return encryptedString, nil
}

func Decrypt(encryptedStr string) (string, int64, error) {
	key := []byte(os.Getenv("SECRET_AES_KEY"))

	encryptedData, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", 0, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", 0, err
	}

	if len(encryptedData) < aes.BlockSize {
		return "", 0, errors.New("encrypted data is too short")
	}

	iv := encryptedData[:aes.BlockSize]
	cipherText := encryptedData[aes.BlockSize:]

	mode := cipher.NewCBCDecrypter(block, iv)

	mode.CryptBlocks(cipherText, cipherText)

	unpaddedData := unpad(cipherText)

	parts := strings.Split(string(unpaddedData), ":")
	if len(parts) != 2 {
		return "", 0, errors.New("invalid decrypted data format")
	}

	unixTime, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", 0, err
	}

	return parts[0], unixTime, nil
}

func pad(data []byte, blockSize int) []byte {
	padding := blockSize - (len(data) % blockSize)
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func unpad(data []byte) []byte {
	length := len(data)
	unpadding := int(data[length-1])
	return data[:(length - unpadding)]
}

/* func generateAESKey(length int) ([]byte, error) {
	key := make([]byte, length)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
} */

func ValidateCSRFToken(token string) error {
	decryptedStr, decryptedUnixTime, err := Decrypt(token)
	if err != nil {
		return err
	}

	// Check if string exists in DB
	csrfToken, err := database.DB.GetCSRFToken(decryptedStr)
	if err != nil {
		return err
	}

	// Unix time validation
	if decryptedUnixTime > time.Now().Unix() || csrfToken.UnixTime != decryptedUnixTime {
		return errors.New("invalid token UNIX time.")
	}

	// Check if used
	if csrfToken.IsUsed == false {
		return errors.New("token already used.")
	}

	return nil
}
