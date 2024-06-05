package middleware

import (
	"encoding/hex"
	"errors"
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

			secret, err := generateCSRFSecret()

			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error generating secret.", http.StatusForbidden)
				return
			}

			session.Values["csrf_secret"] = secret

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

	if csrfSecret, ok := session.Values["csrf_secret"]; ok {
		return hex.DecodeString(csrfSecret.(string))
	}

	return nil, errors.New("no secret in session")
}

func GetUserIDFromSession(r *http.Request) (int, error) {
	session, err := sessions.Store.Get(r, "yd_vending_sessions")

	if err != nil {
		return 0, err
	}

	if userID, ok := session.Values["user_id"]; ok {
		if intUserID, ok := userID.(int); ok {
			return intUserID, nil
		}
		return 0, errors.New("user_id is not of type int")
	}

	return 0, errors.New("no user_id in session")
}
