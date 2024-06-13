package conversions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davidalvarez305/yd_vending/constants"
)

type GoogleEventParamsLead struct {
	GCLID string `json:"gclid"`
}

type GoogleEventLead struct {
	Name   string                `json:"name"`
	Params GoogleEventParamsLead `json:"params"`
}

type GooglePayload struct {
	ClientID string            `json:"client_id"`
	UserId   string            `json:"userId"`
	Events   []GoogleEventLead `json:"events"`
}

func SendGoogleConversion(payload GooglePayload) error {
	endpoint := "https://www.google-analytics.com/mp/collect"

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling Google payload: %+v\n", err)
		return err
	}

	url := endpoint + "?measurement_id=" + constants.GoogleAnalyticsID + "&api_secret=" + constants.GoogleAnalyticsAPISecretKey

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error sending Google request: %+v\n", err)
		return err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Google conversions request error: %+v\n", err)
		return fmt.Errorf("facebook API returned non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
