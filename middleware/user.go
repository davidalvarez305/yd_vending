package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/sessions"
)

func UserTracking(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.URL.Path, "/static/") || strings.Contains(r.URL.Path, "/partials/") {
			next.ServeHTTP(w, r)
			return
		}

		if helpers.UserAgentIsBot(r.Header.Get("User-Agent")) {
			next.ServeHTTP(w, r)
			return
		}

		isNew, err := sessions.IsNew(r)
		if err != nil {
			http.Error(w, "Unable to check if session is new.", http.StatusForbidden)
			return
		}

		if isNew {
			err = sessions.Create(r, w)
			if err != nil {
				http.Error(w, "Failed to create session.", http.StatusForbidden)
				return
			}
		}

		session, err := sessions.Get(r)
		if err != nil {
			http.Error(w, "Error getting session in context.", http.StatusInternalServerError)
			return
		}
		r = r.WithContext(context.WithValue(r.Context(), "external_id", session.ExternalID))

		next.ServeHTTP(w, r)
	})
}
