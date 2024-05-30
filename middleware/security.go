package middleware

import (
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
		cspDirective := "default-src 'self' " + os.Getenv("AWS_STORAGE_BUCKET")
		w.Header().Set("Content-Security-Policy", cspDirective)

		next.ServeHTTP(w, r)
	})
}

func CSRFProtectMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete {

			csrfToken := r.Header.Get("X-CSRF-Token")
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
