package handlers

import (
	"fmt"
	"net/http"

	"github.com/davidalvarez305/yd_vending/helpers"
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
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data", http.StatusBadRequest)
		return
	}

	var leadForm types.LeadFormWebhook

	leadForm.LeadID = r.FormValue("lead_id")
	leadForm.APIVersion = r.FormValue("api_version")
	leadForm.FormID = helpers.ParseInt64(r.FormValue("form_id"))
	leadForm.CampaignID = helpers.ParseInt64(r.FormValue("campaign_id"))
	leadForm.GCLID = r.FormValue("gcl_id")
	leadForm.GoogleKey = r.FormValue("google_key")

	if adGroupIDStr := r.FormValue("adgroup_id"); adGroupIDStr != "" {
		adGroupID := helpers.ParseInt64(adGroupIDStr)
		leadForm.AdGroupID = &adGroupID
	}

	if creativeIDStr := r.FormValue("creative_id"); creativeIDStr != "" {
		creativeID := helpers.ParseInt64(creativeIDStr)
		leadForm.CreativeID = &creativeID
	}

	if isTestStr := r.FormValue("is_test"); isTestStr != "" {
		isTest := isTestStr == "true"
		leadForm.IsTest = &isTest
	}

	userColumnData := []types.UserColumnData{}
	for key, values := range r.Form {
		if key == "user_column_data" {
			for _, value := range values {
				userColumnData = append(userColumnData, types.UserColumnData{
					ColumnID:    value,
					ColumnValue: value,
				})
			}
		}
	}
	leadForm.UserColumnData = userColumnData

	fmt.Printf("LEAD FORM: %+v\n", leadForm)
}
