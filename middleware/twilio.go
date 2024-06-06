package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"os"
	"strings"
)

func validateTwilioWebhook(r *http.Request) error {
	authToken := os.Getenv("TWILLIO_AUTH_TOKEN")
	twilioSignature := r.Header.Get("X-Twilio-Signature")

	url := "http://" + r.Host + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	data := url
	if r.Method == "POST" {
		r.ParseForm()
		data += strings.Join(r.Form["Body"], "")
	}

	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(data))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if twilioSignature != expectedSignature {
		return errors.New("invalid Twilio signature")
	}

	return nil
}
