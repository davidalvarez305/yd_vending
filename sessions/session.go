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
	fmt.Println("=========================")
	fmt.Println("RequestURI", r.RequestURI)
	fmt.Println("Method", r.Method)
	fmt.Println("URL", r.URL)
	fmt.Println("r.URL.User", r.URL.User)
	fmt.Println("Cookies", r.Cookies())
	fmt.Println("Method", r.UserAgent())
	fmt.Println("cookie", cookie)
	fmt.Println("=========================")
	if err != nil {
		return "", err
	}

	return cookie.Value, err
}

func IsNew(r *http.Request) (bool, error) {
	_, err := r.Cookie(constants.CookieName)
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
		ExternalID:       uuid.New().String(),
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

	SetCookie(w, expirationTime, session.CSRFSecret)

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

	SetCookie(w, expirationTime, secret)

	return nil
}

func SetCookie(w http.ResponseWriter, expires time.Time, value string) {
	http.SetCookie(w, &http.Cookie{
		Name:     constants.CookieName,
		Value:    value,
		Path:     "/",
		Domain:   constants.DomainHost,
		Expires:  expires,
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})
}
