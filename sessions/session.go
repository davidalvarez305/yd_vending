package sessions

import (
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitializeSessions() {
	authKey := []byte(os.Getenv("AUTH_SECRET_KEY"))
	encKey := []byte(os.Getenv("ENC_SECRET_KEY"))

	var store = sessions.NewCookieStore(authKey, encKey)

	Store = store
}
