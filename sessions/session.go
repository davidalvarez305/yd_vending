package sessions

import (
	"encoding/hex"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitializeSessions() error {
	authKey, err := hex.DecodeString(constants.AuthSecretKey)

	if err != nil {
		return err
	}

	encKey, err := hex.DecodeString(constants.EncSecretKey)

	if err != nil {
		return err
	}

	var store = sessions.NewCookieStore(authKey, encKey)

	Store = store

	return nil
}
