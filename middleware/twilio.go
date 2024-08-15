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
	fmt.Println(authToken)

	scheme := "https"
	baseURL := scheme + "://" + r.Host + r.URL.Path
	fmt.Println(baseURL)

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return errors.New("error parsing form data")
		}
	}

	var sortedKeys []string

	sortedKeys = append(sortedKeys, r.Form["Body"]...)
	sort.Strings(sortedKeys)

	var sortedParams strings.Builder
	for _, key := range sortedKeys {
		value := r.Form[key][0]

		sortedParams.WriteString(key)
		sortedParams.WriteString(value)
	}

	data := baseURL + sortedParams.String()
	fmt.Println(data)

	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(data))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	if twilioSignature != expectedSignature {
		return errors.New("invalid Twilio signature")
	}

	return nil
}
