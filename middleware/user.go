package middleware

import (
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/davidalvarez305/budgeting/sessions"
)

func UserTracking(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Store.Get(r, "yd_vending_sessions")

		if err != nil {
			http.Error(w, "Unable to retrieve session.", http.StatusForbidden)
			return
		}

		if session.IsNew {
			// Cookie Settings
			/* http.SetCookie(w, &http.Cookie{
				Name:     os.Getenv("COOKIE_NAME"),
				Value:    "cookie_value",
				Path:     "/",
				Domain:   os.Getenv("ROOT_DOMAIN"),
				Expires:  time.Now().Add(24 * time.Hour), // Expires in 24 hours
				HttpOnly: false,
				SameSite: http.SameSiteStrictMode,
				Secure:   true,
			}) */

			session.Values["csrfToken"] = generateCSRFToken()

			err = session.Save(r, w)

			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error saving session.", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func GetTokenFromSession(r *http.Request) ([]byte, error) {
	session, err := sessions.Store.Get(r, "yd_vending_sessions")

	if err != nil {
		return nil, err
	}

	return hex.DecodeString(session.Values["csrfToken"].(string))
}
