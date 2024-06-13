package middleware

import (
	"context"
	"encoding/hex"
	"errors"
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

			marketingUserId := uuid.New().String()
			session.Values["csrf_secret"] = secret
			session.Values["marketing_user_id"] = marketingUserId

			err = session.Save(r, w)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error saving session.", http.StatusForbidden)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "marketing_user_id", marketingUserId))
		} else {
			userId, err := GetMarketingUserIDFromSession(r)
			if err != nil {
				fmt.Printf("%+v\n", err)
				http.Error(w, "Error getting marketing user ID from session.", http.StatusForbidden)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "marketing_user_id", userId))
		}

		next.ServeHTTP(w, r)
	})
}

func GetTokenFromSession(r *http.Request) ([]byte, error) {
	session, err := sessions.Store.Get(r, "yd_vending_sessions")

	if err != nil {
		return nil, err
	}

	if csrfSecret, ok := session.Values["csrf_secret"].(string); ok {
		decodedSecret, err := hex.DecodeString(csrfSecret)
		if err != nil {
			return nil, err
		}
		return decodedSecret, nil
	}
	return nil, errors.New("csrf_secret not found in session values or is not a string")
}

func GetUserIDFromSession(r *http.Request) (int, error) {
	session, err := sessions.Store.Get(r, "yd_vending_sessions")

	if err != nil {
		return 0, err
	}

	if userID, ok := session.Values["user_id"]; ok {
		if intUserID, ok := userID.(int); ok {
			return intUserID, nil
		}
		return 0, errors.New("user_id is not of type int")
	}

	return 0, errors.New("no user_id in session")
}

func GetMarketingUserIDFromSession(r *http.Request) (string, error) {
	session, err := sessions.Store.Get(r, "yd_vending_sessions")
	if err != nil {
		return "", err
	}

	if marketingUserID, ok := session.Values["marketing_user_id"]; ok {
		if strUserID, ok := marketingUserID.(string); ok {
			return strUserID, nil
		}
		return "", errors.New("marketing_user_id is not of type string")
	}

	return "", errors.New("no marketing_user_id in session")
}
