package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"net/http"
	"sort"
	"strings"

	"github.com/davidalvarez305/yd_vending/constants"
)

func validateTwilioWebhook(r *http.Request) error {
	authToken := constants.TwilioAuthToken
	twilioSignature := r.Header.Get("X-Twilio-Signature")

	scheme := "https"
	baseURL := scheme + "://" + r.Host + r.URL.Path

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return errors.New("error parsing form data")
		}
	}

	var sortedKeys []string
	for key := range r.Form {
		sortedKeys = append(sortedKeys, key)
	}

	sort.Strings(sortedKeys)

	var sortedParams []string
	for _, key := range sortedKeys {
		values := r.Form[key]
		for _, value := range values {
			sortedParams = append(sortedParams, key+value)
		}
	}

	data := baseURL + strings.Join(sortedParams, "")

	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(data))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if twilioSignature != expectedSignature {
		return errors.New("invalid Twilio signature")
	}

	return nil
}
