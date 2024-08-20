package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/csrf"
	"github.com/davidalvarez305/yd_vending/sessions"
)

var urlsToSkip = []string{"/static/", "/partials/", "/sms/", "/call/", "/webhooks/"}

func UserTracking(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if csrf.UrlsListHasCurrentPath(urlsToSkip, r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		isNew, err := sessions.IsNew(r)
		if err != nil {
			http.Error(w, "Unable to check if session is new.", http.StatusInternalServerError)
			return
		}

		var externalId, csrfSecret string

		if isNew {
			session, err := sessions.Create(r, w)
			if err != nil {
				http.Error(w, "Failed to create session.", http.StatusInternalServerError)
				return
			}

			expirationTime := time.Unix(session.DateExpires, 0).UTC()

			sessions.SetCookie(w, expirationTime, session.CSRFSecret)

			externalId = session.ExternalID
			csrfSecret = session.CSRFSecret
		}

		if !isNew {
			session, err := sessions.Get(r)
			if err != nil {
				http.Error(w, "Failed to retrieve session in user middleware.", http.StatusInternalServerError)
				return
			}

			externalId = session.CSRFSecret
			csrfSecret = session.CSRFSecret
		}

		r = r.WithContext(context.WithValue(r.Context(), "external_id", externalId))
		r = r.WithContext(context.WithValue(r.Context(), "csrf_secret", csrfSecret))
		next.ServeHTTP(w, r)
	})
}
