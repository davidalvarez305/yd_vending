package sessions

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/csrf"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/models"
	"github.com/google/uuid"
)

func getSessionFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie(constants.CookieName)
	if err != nil {
		return "", err
	}

	return cookie.Value, err
}

func IsNew(r *http.Request) (bool, error) {
	cookie, err := r.Cookie(constants.CookieName)
	fmt.Printf("cookie: %+v\n", cookie)
	fmt.Printf("err: %+v\n", err)
	if err == http.ErrNoCookie {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return false, nil
}

func Get(r *http.Request) (models.Session, error) {
	var sessions models.Session

	userSecret, err := getSessionFromRequest(r)
	if err != nil {
		return sessions, err
	}

	sessions, err = database.GetSession(userSecret)
	if err != nil {
		return sessions, err
	}

	return sessions, nil
}

func Create(r *http.Request, w http.ResponseWriter) error {
	secret, err := csrf.GenerateCSRFSecret()
	if err != nil {
		return err
	}

	googleClientID, err := GetGoogleClientIDFromRequest(r)

	if err != nil {
		fmt.Printf("%+v\n", err)
		fmt.Println("Couldn't extract client ID from GA.")
	}

	fbClickID, err := GetFacebookClickIDFromRequest(r)

	if err != nil {
		fmt.Printf("%+v\n", err)
		fmt.Println("Couldn't extract FB ClickID.")
	}

	fbClientID, err := GetFacebookClientIDFromRequest(r)

	if err != nil {
		fmt.Printf("%+v\n", err)
		fmt.Println("Couldn't extract FB ClientID.")
	}

	expirationTime := time.Now().Add(time.Duration(constants.SessionLength) * time.Second)

	session := models.Session{
		CSRFSecret:       secret,
		GoogleUserID:     uuid.New().String(),
		GoogleClientID:   googleClientID,
		FacebookClickID:  fbClickID,
		FacebookClientID: fbClientID,
		DateCreated:      time.Now().Unix(),
		DateExpires:      expirationTime.Unix(),
	}

	err = database.CreateSession(session)
	if err != nil {
		fmt.Printf("FAILED TO CREATE SESSION: %+v\n", err)
		return err
	}

	cookie := &http.Cookie{
		Name:     constants.SessionName,
		Value:    session.CSRFSecret,
		Path:     "/",
		Domain:   constants.RootDomain,
		HttpOnly: true,
		Secure:   true,
		Expires:  expirationTime,
		SameSite: http.SameSiteStrictMode,
	}
	http.SetCookie(w, cookie)

	return nil
}

func Update(values models.Session) error {
	err := database.UpdateSession(values)
	if err != nil {
		return err
	}

	return nil
}

func Destroy(r *http.Request, w http.ResponseWriter) error {
	secret, err := getSessionFromRequest(r)
	if err != nil {
		return err
	}

	err = database.DeleteSession(secret)
	if err != nil {
		return err
	}

	expirationTime := time.Now().Add(-24 * time.Hour)

	cookie := &http.Cookie{
		Name:     constants.SessionName,
		Value:    secret,
		Domain:   constants.RootDomain,
		HttpOnly: true,
		Secure:   true,
		Expires:  expirationTime,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	return nil
}
