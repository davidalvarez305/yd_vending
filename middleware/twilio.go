package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"

	"github.com/davidalvarez305/yd_vending/constants"
)

func validateTwilioWebhook(r *http.Request) error {
	authToken := constants.TwilioAuthToken
	twilioSignature := r.Header.Get("X-Twilio-Signature")

	url := "https://" + r.Host + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	data := url
	if r.Method == "POST" {
		r.ParseForm()
		data += strings.Join(r.Form["Body"], "")
		return nil
	}

	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(data))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(twilioSignature), []byte(expectedSignature)) {
		// return errors.New("invalid Twilio signature")
		fmt.Println("invalid Twilio signature")
	}

	return nil
}
