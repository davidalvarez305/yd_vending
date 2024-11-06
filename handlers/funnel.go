package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/conversions"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/services"
	"github.com/davidalvarez305/yd_vending/sessions"
	"github.com/davidalvarez305/yd_vending/types"
)

var funnelBaseFilePath = constants.FUNNEL_TEMPLATES_DIR + "base.html"

func createFunnelContext() types.WebsiteContext {
	return types.WebsiteContext{
		PageTitle:                    constants.CompanyName,
		MetaDescription:              "Get a quote for vending machine services.",
		SiteName:                     constants.SiteName,
		StaticPath:                   constants.StaticPath,
		MediaPath:                    constants.MediaPath,
		PhoneNumber:                  helpers.FormatPhoneNumber(constants.CompanyPhoneNumber),
		CurrentYear:                  time.Now().Year(),
		GoogleAnalyticsID:            constants.GoogleAnalyticsID,
		GoogleAdsID:                  constants.GoogleAdsID,
		GoogleAdsCallConversionLabel: constants.GoogleAdsCallConversionLabel,
		FacebookDataSetID:            constants.FacebookDatasetID,
		CompanyName:                  constants.CompanyName,
	}
}

func FunnelHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createFunnelContext()
	ctx.PagePath = constants.RootDomain + r.URL.Path

	externalId, ok := r.Context().Value("external_id").(string)
	if !ok {
		http.Error(w, "Error retrieving external id in context.", http.StatusInternalServerError)
		return
	}

	ctx.ExternalID = externalId

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/funnel/90-day-challenge":
			Get90DayVendingChallengeOptIn(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/funnel/90-day-challenge":
			Post90DayVendingChallengeOptIn(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Get90DayVendingChallengeOptIn(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	isMobile := helpers.IsMobileRequest(r)
	heroImagePath := "90_day_hero_image_desktop.html"
	if isMobile {
		heroImagePath = "90_day_hero_image_mobile.html"
	}

	fileName := "90_day_challenge.html"
	optInForm := constants.FUNNEL_TEMPLATES_DIR + "90_day_challenge_opt_in_form.html"
	files := []string{funnelBaseFilePath, constants.FUNNEL_TEMPLATES_DIR + heroImagePath, optInForm, constants.FUNNEL_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Get 5 Locations in 90 Days Challenge — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func Post90DayVendingChallengeOptIn(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid request.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.QuoteForm
	form.FirstName = helpers.GetStringPointerFromForm(r, "first_name")
	form.LastName = helpers.GetStringPointerFromForm(r, "last_name")
	form.PhoneNumber = helpers.GetStringPointerFromForm(r, "phone_number")
	form.Rent = helpers.GetStringPointerFromForm(r, "rent")
	form.LocationType = helpers.GetIntPointerFromForm(r, "location_type")
	form.MachineType = helpers.GetIntPointerFromForm(r, "machine_type")
	form.FootTraffic = helpers.GetStringPointerFromForm(r, "foot_traffic")
	form.FootTrafficType = helpers.GetStringPointerFromForm(r, "foot_traffic_type")
	form.Message = helpers.GetStringPointerFromForm(r, "message")
	form.Source = helpers.GetStringPointerFromForm(r, "source")
	form.Medium = helpers.GetStringPointerFromForm(r, "medium")
	form.Channel = helpers.GetStringPointerFromForm(r, "channel")
	form.LandingPage = helpers.GetStringPointerFromForm(r, "landing_page")
	form.Keyword = helpers.GetStringPointerFromForm(r, "keyword")
	form.Referrer = helpers.GetStringPointerFromForm(r, "referrer")
	form.ClickID = helpers.GetStringPointerFromForm(r, "click_id")
	form.CampaignID = helpers.GetInt64PointerFromForm(r, "campaign_id")
	form.AdCampaign = helpers.GetStringPointerFromForm(r, "ad_campaign")
	form.AdGroupID = helpers.GetInt64PointerFromForm(r, "ad_group_id")
	form.AdGroupName = helpers.GetStringPointerFromForm(r, "ad_group_name")
	form.AdSetID = helpers.GetInt64PointerFromForm(r, "ad_set_id")
	form.AdSetName = helpers.GetStringPointerFromForm(r, "ad_set_name")
	form.AdID = helpers.GetInt64PointerFromForm(r, "ad_id")
	form.AdHeadline = helpers.GetInt64PointerFromForm(r, "ad_headline")
	form.Language = helpers.GetStringPointerFromForm(r, "language")
	form.Longitude = helpers.GetStringPointerFromForm(r, "longitude")
	form.Latitude = helpers.GetStringPointerFromForm(r, "latitude")
	form.UserAgent = helpers.GetStringPointerFromForm(r, "user_agent")
	form.ButtonClicked = helpers.GetStringPointerFromForm(r, "button_clicked")
	form.IP = helpers.GetStringPointerFromForm(r, "ip")

	facebookClientId, _ := helpers.GetFacebookClientIDFromRequest(r)
	facebookClickId, _ := helpers.GetFacebookClickIDFromRequest(r)
	googleClientId, _ := helpers.GetGoogleClientIDFromRequest(r)

	if facebookClickId != "" {
		form.FacebookClickID = &facebookClickId
	}

	if facebookClientId != "" {
		form.FacebookClientID = &facebookClientId
	}

	if googleClientId != "" {
		form.GoogleClientID = &googleClientId
	}

	session, err := sessions.Get(r)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to retrieve session.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	// User Marketing Variables
	var userIP = helpers.GetUserIPFromRequest(r)
	var userAgent = r.Header.Get("User-Agent")

	if userIP != "" {
		form.IP = &userIP
	}

	if userAgent != "" {
		form.UserAgent = &userAgent
	}

	if session.ExternalID != "" {
		form.ExternalID = &session.ExternalID
	}

	if session.CSRFSecret != "" {
		form.CSRFSecret = &session.CSRFSecret
	}

	leadID, err := database.CreateLeadAndMarketing(form)
	if err != nil {
		fmt.Printf("Error creating lead: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating quote request.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	// HTML successful lead creation
	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Awesome!",
			"AlertMessage": "We received your request and will be right with you.",
		},
	}

	lead, err := database.GetConversionLeadInfo(leadID)

	if err != nil {
		fmt.Printf("ERROR GETTING NEW LEAD FROM DB: %+v\n", err)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	fbEvent := types.FacebookEventData{
		EventName:      "Lead",
		EventTime:      time.Now().Unix(),
		ActionSource:   "website",
		EventSourceURL: helpers.SafeString(form.LandingPage),
		UserData: types.FacebookUserData{
			FirstName:       helpers.HashString(helpers.SafeString(form.FirstName)),
			LastName:        helpers.HashString(helpers.SafeString(form.LastName)),
			Phone:           helpers.HashString(helpers.SafeString(form.PhoneNumber)),
			FBC:             helpers.SafeString(form.FacebookClickID),
			FBP:             helpers.SafeString(form.FacebookClientID),
			State:           helpers.HashString("Florida"),
			ExternalID:      helpers.HashString(helpers.SafeString(form.ExternalID)),
			ClientIPAddress: helpers.SafeString(form.IP),
			ClientUserAgent: helpers.SafeString(form.UserAgent),
		},
	}

	metaPayload := types.FacebookPayload{
		Data: []types.FacebookEventData{fbEvent},
	}

	payload := types.GooglePayload{
		ClientID: helpers.SafeString(form.GoogleClientID),
		UserId:   helpers.SafeString(form.ExternalID),
		Events: []types.GoogleEventLead{
			{
				Name: "quote",
				Params: types.GoogleEventParamsLead{
					GCLID: helpers.SafeString(form.ClickID),
				},
			},
		},
	}

	// Send conversion events
	err = conversions.SendGoogleConversion(payload)

	if err != nil {
		fmt.Printf("Error sending Google conversion: %+v\n", err)
	}

	err = conversions.SendFacebookConversion(metaPayload)

	if err != nil {
		fmt.Printf("Error sending Facebook conversion: %+v\n", err)
	}

	// New lead notification
	subject := "YD Vending: New Lead"
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
		fmt.Printf("ERROR BUILDING QUOTE NOTIFICATION TEMPLATE: %+v\n", err)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)
	err = services.SendGmail(recipients, subject, constants.CompanyEmail, body)
	if err != nil {
		fmt.Printf("ERROR SENDING QUOTE NOTIFICATION EMAIL: %+v\n", err)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
