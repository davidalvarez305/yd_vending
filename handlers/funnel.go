package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/conversions"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/sessions"
	"github.com/davidalvarez305/yd_vending/types"
)

const (
	OptInEventName       string = "opt_in"
	ApplicationEventName string = "90_day_application"
	AppointmentEventName string = "90_day_appointment"
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
		case "/funnel/90-day-challenge-video":
			Get90DayVendingChallengeOptInVideo(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/funnel/90-day-challenge":
			Post90DayVendingChallengeOptIn(w, r)
		case "/funnel/90-day-challenge-application":
			Post90DayVendingChallengeApplication(w, r)
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

func Get90DayVendingChallengeOptInVideo(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "90_day_challenge_video.html"
	files := []string{funnelBaseFilePath, constants.FUNNEL_TEMPLATES_DIR + fileName}
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
	data.PageTitle = "Get 5 Locations in 90 Days Challenge Video — " + constants.CompanyName
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

	var form types.OptIn90DayChallengeForm
	form.Email = helpers.GetStringPointerFromForm(r, "email")
	form.LandingPageID = helpers.GetIntPointerFromForm(r, "landing_page_id")
	form.PercentScrolled = helpers.GetFloat64PointerFromForm(r, "percent_scrolled")
	form.TimeSpentOnPage = helpers.GetInt64PointerFromForm(r, "time_spent_on_page")
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

	form.FacebookClickID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_fbc", "facebook_click_id")
	form.FacebookClientID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_fbp", "facebook_client_id")
	form.GoogleClientID = helpers.GetMarketingCookiesFromRequestOrForm(r, "_ga", "google_client_id")

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

	err = database.Create90DayChallengeOptIn(form)
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

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Awesome!",
			"AlertMessage": "Your guide will be sent to your e-mail.",
		},
	}

	fbEvent := types.FacebookEventData{
		EventName:      OptInEventName,
		EventTime:      time.Now().UTC().Unix(),
		ActionSource:   "website",
		EventSourceURL: helpers.SafeString(form.LandingPage),
		UserData: types.FacebookUserData{
			Email:           helpers.HashString(helpers.SafeString(form.Email)),
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
				Name: OptInEventName,
				Params: types.GoogleEventParamsLead{
					GCLID: helpers.SafeString(form.ClickID),
				},
			},
		},
		UserData: types.GoogleUserData{
			Sha256EmailAddress: []string{helpers.HashString(helpers.SafeString(form.Email))},
		},
	}

	go conversions.SendGoogleConversion(payload)
	go conversions.SendFacebookConversion(metaPayload)

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func Post90DayVendingChallengeApplication(w http.ResponseWriter, r *http.Request) {
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

	var form types.Application90DayChallengeForm

	form.FirstName = helpers.GetStringPointerFromForm(r, "first_name")
	form.LastName = helpers.GetStringPointerFromForm(r, "last_name")
	form.PhoneNumber = helpers.GetStringPointerFromForm(r, "phone_number")
	form.Email = helpers.GetStringPointerFromForm(r, "email")
	form.Website = helpers.GetStringPointerFromForm(r, "website")
	form.CompanyName = helpers.GetStringPointerFromForm(r, "company_name")
	form.YearsInBusiness = helpers.GetIntPointerFromForm(r, "years_in_business")
	form.NumLocations = helpers.GetIntPointerFromForm(r, "num_locations")
	form.City = helpers.GetStringPointerFromForm(r, "city")

	form.LandingPageID = helpers.GetIntPointerFromForm(r, "landing_page_id")
	form.LeadTypeID = helpers.GetIntPointerFromForm(r, "lead_type_id")
	form.PercentScrolled = helpers.GetFloat64PointerFromForm(r, "percent_scrolled")
	form.TimeSpentOnPage = helpers.GetInt64PointerFromForm(r, "time_spent_on_page")
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

	err = database.Create90DayChallengeApplication(form)
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

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Your application has been received.",
		},
	}

	fbEvent := types.FacebookEventData{
		EventName:      ApplicationEventName,
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
				Name: ApplicationEventName,
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

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
