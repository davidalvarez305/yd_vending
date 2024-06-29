package helpers

import (
	"encoding/hex"
	"net/http"

	"github.com/davidalvarez305/yd_vending/sessions"
	"golang.org/x/crypto/bcrypt"
)

const ()

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
	session, err := sessions.Get(r)

	if err != nil {
		return nil, err
	}

	decodedSecret, err := hex.DecodeString(session.CSRFSecret)
	if err != nil {
		return nil, err
	}

	return decodedSecret, nil
}
