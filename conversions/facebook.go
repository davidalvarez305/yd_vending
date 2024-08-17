package conversions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/types"
)

func SendFacebookConversion(payload types.FacebookPayload) error {
	url := fmt.Sprintf("https://graph.facebook.com/v20.0/%s/events?access_token=%s", constants.FacebookDatasetID, constants.FacebookAccessToken)
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
