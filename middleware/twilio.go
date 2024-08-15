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
	fmt.Println(twilioSignature)

	scheme := "https"
	if r.TLS == nil {
		scheme = "http"
	}

	baseURL := scheme + "://" + r.Host + r.URL.Path
	fmt.Println(baseURL)

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			return errors.New("error parsing form data")
		}
	}

	var sortedParams []string
	for key, values := range r.Form {
		for _, value := range values {
			sortedParams = append(sortedParams, key+value)
		}
	}
	sort.Strings(sortedParams)
	fmt.Println(sortedParams)

	data := baseURL + strings.Join(sortedParams, "")

	mac := hmac.New(sha1.New, []byte(authToken))
	mac.Write([]byte(data))
	expectedSignature := base64.StdEncoding.EncodeToString(mac.Sum(nil))
	fmt.Println(expectedSignature)
	fmt.Println(twilioSignature)

	if twilioSignature != expectedSignature {
		return errors.New("invalid Twilio signature")
	}

	return nil
}
