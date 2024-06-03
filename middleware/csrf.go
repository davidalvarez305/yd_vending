package middleware

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/davidalvarez305/budgeting/database"
	"github.com/davidalvarez305/budgeting/models"
)

const (
	CSRF_SECRET_LENGTH = 32
)

func generateCSRFSecret() (string, error) {
	secret := make([]byte, CSRF_SECRET_LENGTH)

	_, err := rand.Read(secret)
	if err != nil {
		return "", err
	}

	encodedSecret := hex.EncodeToString(secret)

	return encodedSecret, nil
}

func Encrypt(unixTime int64, key []byte) (string, error) {
	var encryptedString string
	var token, err = generateCSRFSecret()
	if err != nil {
		return encryptedString, err
	}

	csrfToken := models.CSRFToken{
		ExpiryTime: unixTime,
		Token:      token,
		IsUsed:     false,
	}

	err = database.InsertCSRFToken(csrfToken)
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

func Decrypt(encryptedStr string, userToken []byte) (string, int64, error) {
	encryptedData, err := base64.StdEncoding.DecodeString(encryptedStr)
	if err != nil {
		return "", 0, err
	}

	block, err := aes.NewCipher(userToken)
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

func ValidateCSRFToken(token string, userToken []byte) error {
	decryptedStr, decryptedUnixTime, err := Decrypt(token, userToken)
	if err != nil {
		return err
	}

	// Check if string exists in DB
	csrfToken, err := database.GetCSRFToken(decryptedStr)
	if err != nil {
		return err
	}

	// Unix time validation
	if decryptedUnixTime > time.Now().Unix() || csrfToken.ExpiryTime != decryptedUnixTime {
		return errors.New("invalid token UNIX time")
	}

	// Check if used
	if !csrfToken.IsUsed {
		return errors.New("token already used")
	}

	return nil
}
