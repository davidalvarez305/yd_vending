package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/sessions"
)

func GetUserIDFromSession(r *http.Request) (int, error) {
	session, err := sessions.Store.Get(r, constants.SessionName)

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

func SaveUserIDInSession(w http.ResponseWriter, r *http.Request, userID int) error {
	session, err := sessions.Store.Get(r, constants.SessionName)
	if err != nil {
		return err
	}

	session.Values["user_id"] = userID

	err = session.Save(r, w)
	if err != nil {
		return err
	}

	return nil
}

func GetSessionValueByKey(r *http.Request, key string) string {
	if len(key) == 0 {
		return ""
	}

	session, err := sessions.Store.Get(r, constants.SessionName)
	if err != nil {
		return ""
	}

	if clientID, ok := session.Values[key]; ok {
		if strClientID, ok := clientID.(string); ok {
			return strClientID
		}
		return ""
	}

	return ""
}

func GetGoogleClientIDFromRequest(r *http.Request) (string, error) {
	gaCookie, err := r.Cookie("_ga")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _ga cookie found")
		}
		return "", err
	}

	parts := strings.Split(gaCookie.Value, ".")
	if len(parts) != 4 {
		return "", fmt.Errorf("unexpected _ga cookie format")
	}

	return parts[2] + "." + parts[3], nil
}

func GetFacebookClickIDFromRequest(r *http.Request) (string, error) {
	fbcCookie, err := r.Cookie("_fbc")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _fbc cookie found")
		}
		return "", err
	}

	return fbcCookie.Value, nil
}

func GetFacebookClientIDFromRequest(r *http.Request) (string, error) {
	fbpCookie, err := r.Cookie("_fbp")
	if err != nil {
		if err == http.ErrNoCookie {
			return "", fmt.Errorf("no _fbc cookie found")
		}
		return "", err
	}

	return fbpCookie.Value, nil
}

func HashString(input string) string {
	hasher := sha256.New()
	hasher.Write([]byte(input))
	hashBytes := hasher.Sum(nil)
	hashedString := hex.EncodeToString(hashBytes)
	return hashedString
}
