package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
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
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {

			csrfToken := r.FormValue("csrf_token")
			if csrfToken == "" {
				http.Error(w, "CSRF token is missing", http.StatusForbidden)
				return
			}

			valid := validateCSRFToken(csrfToken)
			if !valid {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

// Placeholder
func validateCSRFToken(token string) bool {
	return token != ""
}
