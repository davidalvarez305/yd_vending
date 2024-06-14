package middleware

import (
	"fmt"
	"net/http"

	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/sessions"
	"github.com/google/uuid"
)

func UserTracking(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := sessions.Store.Get(r, "yd_vending_sessions")
		if err != nil {
			http.Error(w, "Unable to retrieve session.", http.StatusForbidden)
			return
		}

		if session.IsNew {
			secret, err := helpers.GenerateCSRFSecret()
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error generating secret.", http.StatusForbidden)
				return
			}

			googleClientID, err := helpers.GetGoogleClientIDFromRequest(r)

			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Println("Couldn't extract client ID from GA.")
			}

			fbClickID, err := helpers.GetFacebookClickIDFromRequest(r)

			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Println("Couldn't extract FB ClickID.")
			}

			fbClientID, err := helpers.GetFacebookClientIDFromRequest(r)

			if err != nil {
				fmt.Printf("%+v\n", err)
				fmt.Println("Couldn't extract FB ClientID.")
			}

			googleUserId := uuid.New().String()
			session.Values["csrf_secret"] = secret
			session.Values["google_user_id"] = googleUserId
			session.Values["google_client_id"] = googleClientID
			session.Values["facebook_click_id"] = fbClickID
			session.Values["facebook_client_id"] = fbClientID

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
