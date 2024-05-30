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
