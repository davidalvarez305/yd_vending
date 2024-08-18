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
	twilio "github.com/twilio/twilio-go"
	openapi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendOutboundMessageREST(form types.OutboundMessageForm) (types.TwilioSMSResponse, error) {
	var response types.TwilioSMSResponse

	accountSID := constants.TwilioAccountSID
	authToken := constants.TwilioAuthToken

	twilioURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", accountSID)

	formData := url.Values{}
	formData.Set("To", "+1"+form.To)
	formData.Set("From", "+1"+form.From)
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
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			log.Printf("Error parsing Twilio response: %s", err)
			return response, err
		}
		return response, nil
	} else {
		log.Printf("Error sending request: status %d", resp.StatusCode)
		return response, errors.New("failed to send message")
	}
}

func SendOutboundMessage(form types.OutboundMessageForm) (*openapi.ApiV2010Message, error) {
	client := twilio.NewRestClient()

	params := &openapi.CreateMessageParams{}
	params.SetTo("+1" + form.To)
	params.SetFrom("+1" + form.From)
	params.SetBody(form.Body)

	return client.Api.CreateMessage(params)
}
