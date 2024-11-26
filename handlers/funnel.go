package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/conversions"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/models"
	"github.com/davidalvarez305/yd_vending/services"
	"github.com/davidalvarez305/yd_vending/types"
)

var funnelBaseFilePath = constants.FUNNEL_TEMPLATES_DIR + "base.html"

func createFunnelContext() map[string]any {
	return map[string]any{
		"PageTitle":                    constants.CompanyName,
		"MetaDescription":              "Get a quote for vending machine services.",
		"SiteName":                     constants.SiteName,
		"StaticPath":                   constants.StaticPath,
		"MediaPath":                    constants.MediaPath,
		"PhoneNumber":                  helpers.FormatPhoneNumber(constants.CompanyPhoneNumber),
		"CurrentYear":                  time.Now().Year(),
		"GoogleAnalyticsID":            constants.GoogleAnalyticsID,
		"GoogleAdsID":                  constants.GoogleAdsID,
		"GoogleAdsCallConversionLabel": constants.GoogleAdsCallConversionLabel,
		"FacebookDataSetID":            constants.FacebookDatasetID,
		"CompanyName":                  constants.CompanyName,
	}
}

func FunnelHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createFunnelContext()
	ctx["PagePath"] = constants.RootDomain + r.URL.Path

	externalId, ok := r.Context().Value("external_id").(string)
	if !ok {
		http.Error(w, "Error retrieving external id in context.", http.StatusInternalServerError)
		return
	}

	ctx["ExternalID"] = externalId

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/funnel/90-day-challenge":
			Get90DayVendingChallenge(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/funnel/90-day-challenge":
			PostLeadApplication(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Get90DayVendingChallenge(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	isMobile := helpers.IsMobileRequest(r)
	heroImagePath := "90_day_hero_image_desktop.html"
	if isMobile {
		heroImagePath = "90_day_hero_image_mobile.html"
	}

	fileName := "90_day_challenge.html"
	applicationForm := constants.FUNNEL_TEMPLATES_DIR + "90_day_challenge_application_form.html"
	files := []string{funnelBaseFilePath, constants.FUNNEL_TEMPLATES_DIR + heroImagePath, applicationForm, constants.FUNNEL_TEMPLATES_DIR + fileName}
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
	data["PageTitle"] = "Get 5 Locations in 90 Days Challenge — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["LeadTypeID"] = constants.LeadApplicationLeadTypeID
	data["LeadApplicationEventName"] = constants.LeadApplicationEventName

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostLeadApplication(w http.ResponseWriter, r *http.Request) {
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

	var form types.LeadApplicationForm

	form.FirstName = helpers.GetStringPointerFromForm(r, "first_name")
	form.LastName = helpers.GetStringPointerFromForm(r, "last_name")
	form.PhoneNumber = helpers.GetStringPointerFromForm(r, "phone_number")
	form.Email = helpers.GetStringPointerFromForm(r, "email")
	form.Website = helpers.GetStringPointerFromForm(r, "website")
	form.CompanyName = helpers.GetStringPointerFromForm(r, "company_name")
	form.YearsInBusiness = helpers.GetIntPointerFromForm(r, "years_in_business")
	form.NumLocations = helpers.GetIntPointerFromForm(r, "num_locations")
	form.City = helpers.GetStringPointerFromForm(r, "city")
	form.OptInTextMessaging = helpers.GetBoolPointerFromForm(r, "opt_in_text_messaging")

	form.LeadTypeID = helpers.GetIntPointerFromForm(r, "lead_type_id")
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
	form.CSRFToken = helpers.GetStringPointerFromForm(r, "csrf_token")
	form.ExternalID = helpers.GetStringPointerFromForm(r, "external_id")

	// Cookies
	form.FacebookClickID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_fbc", "facebook_click_id")
	form.FacebookClientID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_fbp", "facebook_client_id")
	form.GoogleClientID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_ga", "google_client_id")
	form.CSRFSecret = helpers.GetMarketingCookiesFromRequestOrForm(r, constants.CookieName, "csrf_secret")

	// User Marketing Variables
	var userIP = helpers.GetUserIPFromRequest(r)
	var userAgent = r.Header.Get("User-Agent")

	if userIP != "" {
		form.IP = &userIP
	}

	if userAgent != "" {
		form.UserAgent = &userAgent
	}

	leadId, err := database.CreateLeadApplication(form)
	if err != nil {
		fmt.Printf("Error creating lead application: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating lead application.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Your application has been received.",
		},
	}

	fbEvent := types.FacebookEventData{
		EventName:      constants.LeadApplicationEventName,
		EventTime:      time.Now().UTC().Unix(),
		ActionSource:   "website",
		EventSourceURL: helpers.SafeString(form.LandingPage),
		UserData: types.FacebookUserData{
			Email:           helpers.HashString(helpers.SafeString(form.Email)),
			FirstName:       helpers.HashString(helpers.SafeString(form.FirstName)),
			LastName:        helpers.HashString(helpers.SafeString(form.LastName)),
			Phone:           helpers.HashString(helpers.SafeString(form.PhoneNumber)),
			FBC:             helpers.SafeString(form.FacebookClickID),
			FBP:             helpers.SafeString(form.FacebookClientID),
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
				Name: constants.LeadApplicationEventName,
				Params: types.GoogleEventParamsLead{
					GCLID: helpers.SafeString(form.ClickID),
				},
			},
		},
		UserData: types.GoogleUserData{
			Sha256EmailAddress: []string{helpers.HashString(helpers.SafeString(form.Email))},
			Sha256PhoneNumber:  []string{helpers.HashString(helpers.SafeString(form.PhoneNumber))},

			Address: []types.GoogleUserAddress{
				{
					Sha256FirstName: helpers.HashString(helpers.SafeString(form.FirstName)),
					Sha256LastName:  helpers.HashString(helpers.SafeString(form.LastName)),
				},
			},
		},
	}

	go conversions.SendGoogleConversion(payload)
	go conversions.SendFacebookConversion(metaPayload)

	go func() {
		phoneNumber := helpers.SafeString(form.PhoneNumber)
		textBody := "Thanks for signing up!"

		text, err := services.SendTextMessage(phoneNumber, constants.CompanyPhoneNumber, textBody)
		if err != nil {
			fmt.Printf("ERROR SENDING TEXT MESSAGE: %+v\n", err)
			return
		}

		externalId := helpers.SafeString(text.Sid)
		status := helpers.SafeString(text.Status)

		msg := models.Message{
			ExternalID:  externalId,
			UserID:      constants.DavidUserID,
			LeadID:      leadId,
			Text:        textBody,
			DateCreated: time.Now().Unix(),
			TextFrom:    constants.CompanyPhoneNumber,
			TextTo:      phoneNumber,
			IsInbound:   false,
			Status:      status,
		}

		err = database.SaveSMS(msg)
		if err != nil {
			fmt.Printf("ERROR SAVING TEXT MESSAGE: %+v\n", err)
			return
		}
	}()

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
