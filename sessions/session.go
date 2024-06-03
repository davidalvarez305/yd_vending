package sessions

import (
	"os"

	"github.com/gorilla/sessions"
)

var Store *sessions.CookieStore

func InitializeSessions() *sessions.CookieStore {
	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET_KEY")))

	Store = store

	return store
}
