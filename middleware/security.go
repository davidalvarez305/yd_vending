package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davidalvarez305/budgeting/database"
	"github.com/davidalvarez305/budgeting/helpers"
	"github.com/davidalvarez305/budgeting/models"
)

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// CORS Settings
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ROOT_DOMAIN"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// CSP Settings
		cspDirective := fmt.Sprintf(`
		default-src 'self';
		script-src 'self' https://www.googletagmanager.com %s;
		font-src 'self' https://fonts.bunny.net;
		script-src-elem 'self' https://jspm.dev https://www.googletagmanager.com;
		style-src-attr 'self';
		img-src 'self' https://www.google-analytics.com data:;
		connect-src 'self' https://www.google-analytics.com;
		style-src-elem 'self' https://fonts.bunny.net;
		`, os.Getenv("AWS_STORAGE_BUCKET"))

		w.Header().Set("Content-Security-Policy", cspDirective)

		// Generate a random nonce
		nonce := make([]byte, 16)
		if _, err := rand.Read(nonce); err != nil {
			http.Error(w, "Error creating nonce.", http.StatusInternalServerError)
			return
		}
		nonceBase64 := base64.StdEncoding.EncodeToString(nonce)

		// Pass the nonce to the next handler
		r = r.WithContext(context.WithValue(r.Context(), "nonce", nonceBase64))

		next.ServeHTTP(w, r)
	})
}

func CSRFProtectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/call/") || strings.Contains(r.URL.Path, "/sms/") {
			if err := validateTwilioWebhook(r); err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error validating Twilio webhook.", http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodGet {
			token, err := GetTokenFromSession(r)

			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error getting user token from session.", http.StatusBadRequest)
				return
			}

			var unixTime = time.Now().Unix() + 300 // 5 minutes

			encryptedToken, err := helpers.EncryptToken(unixTime, token)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error encrypting CSRF token.", http.StatusInternalServerError)
				return
			}

			csrfToken := models.CSRFToken{
				ExpiryTime: unixTime,
				Token:      encryptedToken,
				IsUsed:     false,
			}

			err = database.InsertCSRFToken(csrfToken)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error inserting CSRF token.", http.StatusBadRequest)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "csrf_token", encryptedToken))

			next.ServeHTTP(w, r)

		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {

			csrfToken := r.FormValue("csrf_token")
			if csrfToken == "" {
				http.Error(w, "CSRF token is missing", http.StatusForbidden)
				return
			}

			token, err := GetTokenFromSession(r)

			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error getting user token from session.", http.StatusBadRequest)
				return
			}

			dbToken, err := helpers.ValidateCSRFToken(csrfToken, token)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error validating token.", http.StatusBadRequest)
				return
			}

			err = database.MarkCSRFTokenAsUsed(dbToken.Token)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error marking token as used.", http.StatusBadRequest)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
