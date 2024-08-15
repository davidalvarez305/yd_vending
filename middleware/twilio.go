package middleware

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/davidalvarez305/yd_vending/constants"
)

func validateTwilioWebhook(r *http.Request) error {
	authToken := constants.TwilioAuthToken
	twilioSignature := r.Header.Get("X-Twilio-Signature")

	scheme := "https"

	url := scheme + "://" + r.Host + r.URL.Path
	if r.URL.RawQuery != "" {
		url += "?" + r.URL.RawQuery
	}

	var data strings.Builder
	data.WriteString(url)

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return errors.New("error parsing form data")
		}

		var sortedKeys []string
		for key := range r.Form {
			sortedKeys = append(sortedKeys, key)
		}
		sort.Strings(sortedKeys)

		for _, key := range sortedKeys {
			values := r.Form[key]
			for _, value := range values {
				data.WriteString(key)
				data.WriteString(value)
			}
		}
	}

	fmt.Println(data.String())

	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(data.String()))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(twilioSignature), []byte(expectedSignature)) {
		return errors.New("invalid Twilio signature")
	}

	return nil
}
