package helpers

import (
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/davidalvarez305/yd_vending/sessions"
	"golang.org/x/crypto/bcrypt"
)

const (
	SessionName = "yd_vending_sessions"
)

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func ValidatePassword(formPassword, userPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(formPassword))
	return err == nil
}

func GetTokenFromSession(r *http.Request) ([]byte, error) {
	session, err := sessions.Store.Get(r, SessionName)

	if err != nil {
		return nil, err
	}

	if csrfSecret, ok := session.Values["csrf_secret"].(string); ok {
		decodedSecret, err := hex.DecodeString(csrfSecret)
		if err != nil {
			return nil, err
		}
		return decodedSecret, nil
	}
	return nil, errors.New("csrf_secret not found in session values or is not a string")
}
