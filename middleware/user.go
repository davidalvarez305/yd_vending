package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

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
			http.Error(w, "Unable to check if session is new.", http.StatusInternalServerError)
			return
		}

		var externalId string

		if isNew {
			session, err := sessions.Create(r, w)
			if err != nil {
				http.Error(w, "Failed to create session.", http.StatusInternalServerError)
				return
			}

			expirationTime := time.Unix(session.DateExpires, 0).UTC()

			sessions.SetCookie(w, expirationTime, session.CSRFSecret)

			externalId = session.CSRFSecret
		}

		if !isNew {
			session, err := sessions.Get(r)
			if err != nil {
				http.Error(w, "Failed to retrieve session in user middleware.", http.StatusInternalServerError)
				return
			}

			externalId = session.CSRFSecret
		}

		r = r.WithContext(context.WithValue(r.Context(), "external_id", externalId))
		next.ServeHTTP(w, r)
	})
}
