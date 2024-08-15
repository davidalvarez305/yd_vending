package middleware

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/davidalvarez305/yd_vending/constants"
)

func validateTwilioWebhook(r *http.Request) error {
	authToken := constants.TwilioAuthToken
	twilioSignature := r.Header.Get("X-Twilio-Signature")

	if r.Method == "POST" {
		return nil
	}

	url := "https://" + r.Host + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	data := url
	if r.Method == "POST" {
		// Copy the body so it can be read multiple times
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			return errors.New("failed to read request body")
		}
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		r.ParseForm()
		data += strings.Join(r.Form["Body"], "")
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
