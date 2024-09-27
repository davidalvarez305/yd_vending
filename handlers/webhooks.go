package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/conversions"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/services"
	"github.com/davidalvarez305/yd_vending/types"
)

func WebhookHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		switch r.URL.Path {
		case "/webhooks/lead-form":
			handleGoogleLeadFormWebhook(w, r)
		case "/webhooks/seed-live-hourly":
			handleSeedLiveHourly(w, r)
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

	var form types.QuoteForm
	var state string

	for _, item := range leadForm.UserColumnData {
		switch item.ColumnID {
		case "STATE":
			state = item.StringValue
		case "FIRST_NAME":
			form.FirstName = helpers.SafeStringToPointer(item.StringValue)
		case "LAST_NAME":
			form.LastName = helpers.SafeStringToPointer(item.StringValue)
		case "PHONE_NUMBER":
			form.PhoneNumber = helpers.SafeStringToPointer(item.StringValue)
		case "RENT":
			form.Rent = helpers.SafeStringToPointer(item.StringValue)
		case "LOCATION_TYPE":
			form.LocationType = helpers.SafeStringToIntPointer(item.StringValue)
		case "MACHINE_TYPE":
			form.MachineType = helpers.SafeStringToIntPointer(item.StringValue)
		case "FOOT_TRAFFIC":
			form.FootTraffic = helpers.SafeStringToPointer(item.StringValue)
		case "FOOT_TRAFFIC_TYPE":
			form.FootTrafficType = helpers.SafeStringToPointer(item.StringValue)
		case "MESSAGE":
			form.Message = helpers.SafeStringToPointer(item.StringValue)
		case "SOURCE":
			form.Source = helpers.SafeStringToPointer(item.StringValue)
		case "MEDIUM":
			form.Medium = helpers.SafeStringToPointer(item.StringValue)
		case "CHANNEL":
			form.Channel = helpers.SafeStringToPointer(item.StringValue)
		case "LANDING_PAGE":
			form.LandingPage = helpers.SafeStringToPointer(item.StringValue)
		case "KEYWORD":
			form.Keyword = helpers.SafeStringToPointer(item.StringValue)
		case "REFERRER":
			form.Referrer = helpers.SafeStringToPointer(item.StringValue)
		case "CLICK_ID":
			form.ClickID = helpers.SafeStringToPointer(item.StringValue)
		case "CAMPAIGN_ID":
			form.CampaignID = helpers.SafeStringToInt64Pointer(item.StringValue)
		case "AD_CAMPAIGN":
			form.AdCampaign = helpers.SafeStringToPointer(item.StringValue)
		case "AD_GROUP_ID":
			form.AdGroupID = helpers.SafeStringToInt64Pointer(item.StringValue)
		case "AD_GROUP_NAME":
			form.AdGroupName = helpers.SafeStringToPointer(item.StringValue)
		case "AD_SET_ID":
			form.AdSetID = helpers.SafeStringToInt64Pointer(item.StringValue)
		case "AD_SET_NAME":
			form.AdSetName = helpers.SafeStringToPointer(item.StringValue)
		case "AD_ID":
			form.AdID = helpers.SafeStringToInt64Pointer(item.StringValue)
		case "AD_HEADLINE":
			form.AdHeadline = helpers.SafeStringToInt64Pointer(item.StringValue)
		case "LANGUAGE":
			form.Language = helpers.SafeStringToPointer(item.StringValue)
		case "LONGITUDE":
			form.Longitude = helpers.SafeStringToPointer(item.StringValue)
		case "LATITUDE":
			form.Latitude = helpers.SafeStringToPointer(item.StringValue)
		case "USER_AGENT":
			form.UserAgent = helpers.SafeStringToPointer(item.StringValue)
		case "BUTTON_CLICKED":
			form.ButtonClicked = helpers.SafeStringToPointer(item.StringValue)
		case "IP":
			form.IP = helpers.SafeStringToPointer(item.StringValue)
		case "CSRF_TOKEN":
			form.CSRFToken = helpers.SafeStringToPointer(item.StringValue)
		case "EXTERNAL_ID":
			form.ExternalID = helpers.SafeStringToPointer(item.StringValue)
		case "GOOGLE_CLIENT_ID":
			form.GoogleClientID = helpers.SafeStringToPointer(item.StringValue)
		case "FACEBOOK_CLICK_ID":
			form.FacebookClickID = helpers.SafeStringToPointer(item.StringValue)
		case "FACEBOOK_CLIENT_ID":
			form.FacebookClientID = helpers.SafeStringToPointer(item.StringValue)
		case "CSRF_SECRET":
			form.CSRFSecret = helpers.SafeStringToPointer(item.StringValue)
		}
	}

	leadID, err := database.CreateLeadAndMarketing(form)
	if err != nil {
		fmt.Printf("Error creating lead in lead form webhook: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	lead, err := database.GetConversionLeadInfo(leadID)

	if err != nil {
		fmt.Printf("Error getting lead conversion info in lead form webhook: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fbEvent := types.FacebookEventData{
		EventName:    "quote",
		EventTime:    time.Now().Unix(),
		ActionSource: "other",
		UserData: types.FacebookUserData{
			FirstName: helpers.HashString(helpers.SafeString(form.FirstName)),
			LastName:  helpers.HashString(helpers.SafeString(form.LastName)),
			Phone:     helpers.HashString(helpers.SafeString(form.PhoneNumber)),
			State:     state,
		},
	}

	metaPayload := types.FacebookPayload{
		Data: []types.FacebookEventData{fbEvent},
	}

	err = conversions.SendFacebookConversion(metaPayload)

	if err != nil {
		fmt.Printf("Error sending Facebook conversion: %+v\n", err)
	}

	// New lead notification
	subject := "YD Vending: New Lead From Lead Form Webhook"
	recipients := []string{constants.DavidEmail, constants.YovaEmail}
	templateFile := constants.PARTIAL_TEMPLATES_DIR + "new_lead_notification_email.html"

	var notificationTemplateData = map[string]any{
		"Name":           helpers.SafeString(form.FirstName) + " " + helpers.SafeString(form.LastName),
		"PhoneNumber":    helpers.SafeString(form.PhoneNumber),
		"DateCreated":    time.Unix(lead.CreatedAt, 0).Format("01/02/2006 3 PM"),
		"MachineType":    lead.MachineType,
		"LocationType":   lead.LocationType,
		"Message":        helpers.SafeString(form.Message),
		"LeadDetailsURL": fmt.Sprintf("%s/crm/lead/%d", constants.RootDomain, leadID),
		"Location":       "",
	}

	if helpers.SafeString(form.Longitude) != "0.0" && len(helpers.SafeString(form.Longitude)) > 0 || helpers.SafeString(form.Latitude) != "0.0" && len(helpers.SafeString(form.Latitude)) > 0 {
		notificationTemplateData["Location"] = fmt.Sprintf("https://www.google.com/maps?q=%s,%s", helpers.SafeString(form.Latitude), helpers.SafeString(form.Longitude))
	}

	template, err := helpers.BuildStringFromTemplate(templateFile, "email", notificationTemplateData)

	if err != nil {
		fmt.Printf("ERROR BUILDING QUOTE NOTIFICATION TEMPLATE IN LEAD FORM WEBHOOK: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)
	err = services.SendGmail(recipients, subject, constants.CompanyEmail, body)
	if err != nil {
		fmt.Printf("ERROR SENDING QUOTE NOTIFICATION EMAIL IN LEAD FORM WEBHOOK: %+v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func handleSeedLiveHourly(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/json" {
		http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
		return
	}

	var payload []types.SeedLiveTransaction
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	fmt.Printf("%+v\n", payload)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Received successfully"))
}
