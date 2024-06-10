package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/models"
)

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// CORS Settings
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("ROOT_DOMAIN"))
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Generate a random nonce
		nonce := make([]byte, 16)
		if _, err := rand.Read(nonce); err != nil {
			http.Error(w, "Error creating nonce.", http.StatusInternalServerError)
			return
		}
		nonceBase64 := base64.StdEncoding.EncodeToString(nonce)

		// CSP Settings
		cspDirective := fmt.Sprintf(`default-src 'self';
		script-src 'self' https://www.googletagmanager.com %s 'nonce-%s';
		font-src 'self' https://fonts.bunny.net;
		script-src-elem 'self' https://jspm.dev https://www.googletagmanager.com 'nonce-%s';
		style-src 'self';
		img-src 'self' https://www.google-analytics.com data:;
		connect-src 'self' https://www.google-analytics.com;
		style-src-elem 'self' https://fonts.bunny.net;
		style-src-attr 'self' 'unsafe-inline';`, os.Getenv("AWS_STORAGE_BUCKET"), nonceBase64, nonceBase64)

		w.Header().Set("Content-Security-Policy", cspDirective)

		// Pass the nonce to the next handler
		r = r.WithContext(context.WithValue(r.Context(), "nonce", nonceBase64))

		next.ServeHTTP(w, r)
	})
}

func CSRFProtectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if strings.Contains(r.URL.Path, "/static/") {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(r.URL.Path, "/call/inbound") || strings.Contains(r.URL.Path, "/sms/inbound") {
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
			return
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			csrfToken := r.FormValue("csrf_token")

			// If CSRF token is not in form values, check the JSON body
			if csrfToken == "" {
				var jsonBody map[string]string
				err := json.NewDecoder(r.Body).Decode(&jsonBody)
				if err != nil {
					http.Error(w, "Invalid JSON body.", http.StatusBadRequest)
					return
				}

				csrfToken = jsonBody["csrf_token"]
				if csrfToken == "" {
					http.Error(w, "CSRF token is missing.", http.StatusForbidden)
					return
				}
			}

			token, err := GetTokenFromSession(r)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error getting user token from session.", http.StatusBadRequest)
				return
			}

			// Check if string exists in DB
			isUsed, err := database.CheckIsTokenUsed(csrfToken)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Token doesn't exist in DB.", http.StatusBadRequest)
				return
			}

			err = helpers.ValidateCSRFToken(isUsed, csrfToken, token)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error validating token.", http.StatusBadRequest)
				return
			}

			err = database.MarkCSRFTokenAsUsed(csrfToken)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error marking token as used.", http.StatusBadRequest)
				return
			}

			next.ServeHTTP(w, r)
			return
		}
	})
}

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie(os.Getenv("COOKIE_NAME"))
		if err != nil || cookie == nil {
			http.Error(w, "Permission denied", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
