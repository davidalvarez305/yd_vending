package services

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/types"
)

func SendOutboundMessage(form types.OutboundMessageForm) (string, error) {
	var response string

	accountSID := constants.TwilioAccountSID
	authToken := constants.TwilioAuthToken

	twilioURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)

	formData := url.Values{}
	formData.Set("To", form.To)
	formData.Set("From", form.From)
	formData.Set("Body", form.Body)

	req, err := http.NewRequest("POST", twilioURL, bytes.NewBufferString(formData.Encode()))
	if err != nil {
		log.Printf("Error creating request: %s", err)
		return response, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	auth := base64.StdEncoding.EncodeToString([]byte(accountSID + ":" + authToken))
	req.Header.Set("Authorization", "Basic "+auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("Error sending request: %s", err)
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var twilioResp struct {
			SID string `json:"sid"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&twilioResp); err != nil {
			log.Printf("Error parsing Twilio response: %s", err)
			return response, err
		}
		return twilioResp.SID, nil // Return the Twilio SID if successful
	} else {
		log.Printf("Error sending request: status %d", resp.StatusCode)
		return response, errors.New("failed to send message")
	}
}
