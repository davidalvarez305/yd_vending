package conversions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davidalvarez305/yd_vending/constants"
)

type FacebookUserData struct {
	Phone           string `json:"ph"`
	FirstName       string `json:"fn"`
	LastName        string `json:"ln"`
	ClientIPAddress string `json:"client_ip_address"`
	ClientUserAgent string `json:"client_user_agent"`
	FBC             string `json:"fbc"`
	FBP             string `json:"fbp"`
	City            string `json:"ct"`
	State           string `json:"st"`
	ExternalID      string `json:"external_id"`
}

type FacebookEventData struct {
	EventName      string           `json:"event_name"`
	EventTime      int64            `json:"event_time"`
	ActionSource   string           `json:"action_source"`
	EventSourceURL string           `json:"event_source_url"`
	UserData       FacebookUserData `json:"user_data"`
}

type FacebookPayload struct {
	Data []FacebookEventData `json:"data"`
}

func SendFacebookConversion(payload FacebookPayload) error {
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/events?access_token=%s", constants.FacebookDatasetID, constants.FacebookAppAccessToken)
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Printf("Error marshaling meta payload: %+v\n", err)
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Printf("Error sending meta request: %+v\n", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Meta conversions request error: %+v\n", err)
		return fmt.Errorf("facebook API returned non-200 status code: %d", resp.StatusCode)
	}

	return nil
}
