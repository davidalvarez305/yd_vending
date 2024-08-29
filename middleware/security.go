package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/csrf"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/models"
	"github.com/davidalvarez305/yd_vending/sessions"
)

var allowedOrigins = []string{
	constants.RootDomain,
	"https://www.googleadservices.com",
}

func isOriginAllowed(origin string) bool {
	for _, allowedOrigin := range allowedOrigins {
		if origin == allowedOrigin {
			return true
		}
	}
	return false
}

func SecurityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS Settings
		origin := r.Header.Get("Origin")
		if isOriginAllowed(origin) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		} else {
			w.Header().Set("Access-Control-Allow-Origin", "none")
		}
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Generate a random nonce
		nonce := make([]byte, 16)
		if _, err := rand.Read(nonce); err != nil {
			fmt.Printf("ERROR CREATIONG NONCE SECURITY MIDDLEWARE: %+v\n", err)
			http.Error(w, "Error creating nonce.", http.StatusInternalServerError)
			return
		}
		nonceBase64 := base64.StdEncoding.EncodeToString(nonce)

		// CSP Settings
		cspDirective := fmt.Sprintf(`default-src 'self';
		script-src 'self' https://www.googletagmanager.com %s 'nonce-%s';
		font-src 'self' https://fonts.bunny.net;
		frame-src 'self' https://td.doubleclick.net https://www.facebook.com;
		script-src-elem 'self' https://jspm.dev https://www.googletagmanager.com 'nonce-%s' https://connect.facebook.net http://cdn.jsdelivr.net https://cdn.jsdelivr.net https://code.jquery.com;
		style-src 'self' %s;
		img-src 'self' https://www.google-analytics.com data: https://cdn.tailkit.com https://www.facebook.com https://www.google.com https://adservice.google.com https://www.googletagmanager.com %s %s;
		connect-src 'self' https://www.google-analytics.com https://www.googleadservices.com https://www.google.com https://adservice.google.com https://www.facebook.com;
		style-src-elem 'self' https://fonts.bunny.net %s https://cdn.jsdelivr.net;
		style-src-attr 'self' 'unsafe-inline';`, constants.AWSStorageBucket, nonceBase64, nonceBase64, constants.AWSStorageBucket, constants.AWSS3LiveImagesPath, constants.AWSStorageBucket, constants.AWSStorageBucket)

		w.Header().Set("Content-Security-Policy", cspDirective)

		// Pass the nonce to the next handler
		r = r.WithContext(context.WithValue(r.Context(), "nonce", nonceBase64))

		next.ServeHTTP(w, r)
	})
}

func CSRFProtectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if csrf.UrlsListHasCurrentPath([]string{"/static/", "/partials/", "/webhooks/"}, path) {
			next.ServeHTTP(w, r)
			return
		}

		if strings.HasPrefix(path, "/crm/lead/") && strings.Contains(path, "/messages") {
			next.ServeHTTP(w, r)
			return
		}

		if strings.Contains(path, "/call/inbound") || strings.Contains(path, "/sms/inbound") {
			if err := validateTwilioWebhook(r); err != nil {
				fmt.Printf("ERROR VALIDATING TWILIO SECURITY MIDDLEWARE: %+v\n", err)
				http.Error(w, "Error validating Twilio webhook.", http.StatusInternalServerError)
				return
			}
			next.ServeHTTP(w, r)
			return
		}

		var csrfURLs = []string{"/contact", "/quote", "/login", "/crm", "/atm-services", "ice-vending-services"}

		if r.Method == http.MethodGet && (csrf.UrlsListHasCurrentPath(csrfURLs, path) || path == "/") {
			csrfSecret, ok := r.Context().Value("csrf_secret").(string)
			if !ok {
				fmt.Printf("ERROR GETTING CSRF SECRET IN SECURITY MIDDLEWARE")
				http.Error(w, "Error retrieving user secret token in middleware.", http.StatusInternalServerError)
				return
			}

			decodedSecret, err := hex.DecodeString(csrfSecret)
			if err != nil {
				fmt.Printf("ERROR GETTING DECODED STRING FROM CSRF SECRET SECURITY MIDDLEWARE")
				http.Error(w, "Error decoding user secret token in middleware.", http.StatusInternalServerError)
				return
			}

			var unixTime = time.Now().Unix() + constants.TokenExpiryLength

			encryptedToken, err := csrf.EncryptToken(unixTime, decodedSecret)
			if err != nil {
				fmt.Printf("ERROR ENCRYPTING TOKEN SECURITY MIDDLEWARE: %+v\n", err)
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
				fmt.Printf("ERROR INSERTING TOKEN SECURITY MIDDLEWARE: %+v\n", err)
				http.Error(w, "Error inserting CSRF token.", http.StatusBadRequest)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "csrf_token", encryptedToken))

			next.ServeHTTP(w, r)
			return
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {
			csrfToken := r.FormValue("csrf_token")
			if csrfToken == "" {
				// If CSRF token is not in form values, check the request headers
				csrfToken = r.Header.Get("X-Csrf-Token")
				if csrfToken == "" {
					fmt.Printf("MISSING CSRF TOKEN")
					http.Error(w, "CSRF token is missing.", http.StatusForbidden)
					return
				}
			}

			token, err := helpers.GetTokenFromSession(r)
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

			err = csrf.ValidateCSRFToken(isUsed, csrfToken, token)
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

		next.ServeHTTP(w, r)
	})
}

func AuthRequired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		values, err := sessions.Get(r)
		if err != nil || values.UserID == 0 {
			fmt.Printf("SESSION NOT FOUND, REDIRECTING TO LOGIN PAGE: %+v\n", err)
			http.Redirect(w, r, "/login?redirect="+r.URL.Path, http.StatusSeeOther)
			return
		}

		user, err := database.GetUserById(values.UserID)
		if err != nil {
			fmt.Printf("CANNOT GET USER PERMISSION DENIED: %+v\n", err)
			http.Error(w, "Permission denied", http.StatusUnauthorized)
			return
		}

		if !user.IsAdmin {
			fmt.Printf("IS NOT ADMIN PERMISSION DENIED: %+v\n", err)
			http.Error(w, "Admins only.", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}
