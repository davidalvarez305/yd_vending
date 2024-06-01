package sessions

import (
	"github.com/gorilla/sessions"
	"os"
)

var Store *Session

func InitializeSessions() *Session {
	var store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_SECRET_KEY")))

	Store = store

	return store
}
