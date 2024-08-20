package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/types"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/webhooks/lead-form":
			handleGoogleLeadFormWebhook(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func handleGoogleLeadFormWebhook(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var leadForm types.LeadFormWebhook
	if err := json.NewDecoder(r.Body).Decode(&leadForm); err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	if leadForm.GoogleKey != constants.GoogleWebhookKey {
		http.Error(w, "Permission denied. Not valid key.", http.StatusForbidden)
		return
	}

	fmt.Printf("LEAD FORM: %+v\n", leadForm)
}
