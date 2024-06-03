package sessions

import (
	"encoding/hex"
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitializeSessions() error {
	authKey, err := hex.DecodeString(os.Getenv("AUTH_SECRET_KEY"))

	if err != nil {
		return err
	}

	encKey, err := hex.DecodeString(os.Getenv("ENC_SECRET_KEY"))

	if err != nil {
		return err
	}

	var store = sessions.NewCookieStore(authKey, encKey)

	Store = store

	return nil
}
