package conversions

import (
	"bytes"
	"encoding/json"
	"net/http"
)

const (
	measurementID = "G-XXXXXXXXXX"
	apiSecret     = "<secret_value>"
)

type EventParamsLead struct {
	Gclid string `json:"gclid"`
}

type EventLead struct {
	Name   string          `json:"name"`
	Params EventParamsLead `json:"params"`
}

type PayloadLead struct {
	ClientID string      `json:"client_id"`
	UserId   string      `json:"userId"`
	Events   []EventLead `json:"events"`
}

func SendLeadEvent(payload PayloadLead) error {
	endpoint := "https://www.google-analytics.com/mp/collect"

	// Convert payload to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Construct the URL with measurement_id and secret_key as URL parameters
	url := endpoint + "?measurement_id=" + measurementID + "&api_secret=" + apiSecret

	// Send HTTP POST request
	_, err = http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	return nil
}
