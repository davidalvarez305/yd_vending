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
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

var websiteBaseFilePath = constants.WEBSITE_TEMPLATES_DIR + "base.html"
var websiteFooterFilePath = constants.WEBSITE_TEMPLATES_DIR + "footer.html"

func createWebsiteContext() types.WebsiteContext {
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

func WebsiteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createWebsiteContext()
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
		case "/contact":
			GetContactForm(w, r, ctx)
		case "/login":
			GetLogin(w, r, ctx)
		case "/about":
			GetAbout(w, r, ctx)
		case "/privacy-policy":
			GetPrivacyPolicy(w, r, ctx)
		case "/terms-and-conditions":
			GetTermsAndConditions(w, r, ctx)
		case "/robots.txt":
			GetRobots(w, r, ctx)
		case "/ice-vending-services":
			GetIceLandingPage(w, r, ctx)
		case "/atm-services":
			GetATMLandingPage(w, r, ctx)
		case "/servicio-de-maquinas-expendedoras":
			GetVendingESP(w, r, ctx)
		case "/cdd-proposal":
			GetCDDProposal(w, r, ctx)
		case "/":
			GetHome(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/quote":
			PostQuote(w, r, ctx)
		case "/contact":
			PostContactForm(w, r, ctx)
		case "/login":
			PostLogin(w, r, ctx)
		case "/logout":
			PostLogout(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetCDDProposal(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "proposal.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Coral Bay Community Dev District — Proposal by " + constants.CompanyName
	data.Nonce = nonce

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetHome(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "home.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}
	vendingTypes, err := database.GetVendingTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending types.", http.StatusInternalServerError)
		return
	}

	vendingLocations, err := database.GetVendingLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Miami Vending Services — " + constants.CompanyName
	data.Nonce = nonce
	data.Features = []string{
		"Regular Product Rotation",
		"Modern & Smart Machines",
		"Health-Conscious Choices",
		"Flexible Working Hours",
		"Forward-Thinking Product Selection",
		"Attentive Customer Care",
		"Responsive Repair & Servicing Agents",
		"High Upkeep & Maintanance",
	}
	data.CSRFToken = csrfToken
	data.VendingTypes = vendingTypes
	data.VendingLocations = vendingLocations

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetVendingESP(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "vending_esp.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}
	vendingTypes, err := database.GetVendingTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending types.", http.StatusInternalServerError)
		return
	}

	vendingLocations, err := database.GetVendingLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Servicios de Máquinas Expendedoras en Miami — " + constants.CompanyName
	data.Nonce = nonce
	data.Features = []string{
		"Rotación Regular de Productos", "Máquinas Modernas e Inteligentes", "Opciones Saludables", "Horarios de Trabajo Flexibles", "Selección de Productos Innovadora", "Atención al Cliente Cuidadosa", "Agentes de Reparación y Servicio Receptivos", "Mantenimiento y Cuidado de Alta Calidad",
	}
	data.CSRFToken = csrfToken
	data.VendingTypes = vendingTypes
	data.VendingLocations = vendingLocations

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetIceLandingPage(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "ice.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}
	vendingTypes, err := database.GetVendingTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending types.", http.StatusInternalServerError)
		return
	}

	vendingLocations, err := database.GetVendingLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Miami Ice Vending Services — " + constants.CompanyName
	data.Nonce = nonce
	data.Features = []string{
		"Free Installation & Inspection",
		"Regular Upkeep & Maintance",
		"Impeccable Hygiene & Safety Standards",
		"Modern & Smart Machines",
		"Secure Payment Solutions Implemented",
		"Energy Efficiency",
		"Attentive & Efficient Service",
		"Commercial Ice Makers Available",
		"Flexible Working Partnerships",
	}
	data.CSRFToken = csrfToken
	data.VendingTypes = vendingTypes
	data.VendingLocations = vendingLocations

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetATMLandingPage(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "atm.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}
	vendingTypes, err := database.GetVendingTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending types.", http.StatusInternalServerError)
		return
	}

	vendingLocations, err := database.GetVendingLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Miami ATM Rental Services — " + constants.CompanyName
	data.Nonce = nonce
	data.Features = []string{
		"Complimentary Installation & Initial Inspection",
		"Ongoing Maintenance & Regular Upkeep",
		"High Standards for Hygiene & Security",
		"Advanced & User-Friendly Machines",
		"Robust Payment Security Measures",
		"Energy-Efficient Technology",
		"Responsive & Professional Customer Support",
		"Commercial-Grade ATM Solutions",
		"Customizable Partnership Options",
	}
	data.CSRFToken = csrfToken
	data.VendingTypes = vendingTypes
	data.VendingLocations = vendingLocations

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetRobots(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	robotsTxtContent := `
	# robots.txt for https://ydvending.com/

	# Allow all robots complete access
	User-agent: *
	Disallow:
	`

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")

	_, err := w.Write([]byte(robotsTxtContent))
	if err != nil {
		http.Error(w, "Error writing robots.txt content.", http.StatusInternalServerError)
	}
}

func GetAbout(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "about.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "About Us — " + constants.CompanyName
	data.Nonce = nonce

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetPrivacyPolicy(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "privacy.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Privacy Policy — " + constants.CompanyName
	data.Nonce = nonce

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetTermsAndConditions(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "terms.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Terms & Conditions — " + constants.CompanyName
	data.Nonce = nonce

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostQuote(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
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

		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
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
	facebookClickId, _ := helpers.GetFacebookClientIDFromRequest(r)
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

		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
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

		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
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
		EventName:      "quote",
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

func GetContactForm(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "contact_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}

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
	data.PageTitle = "Contact Us — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := helpers.ServeContent(w, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostContactForm(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	err := r.ParseForm()

	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to parse form data.",
			},
		}

		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.ContactForm
	err = decoder.Decode(&form, r.PostForm)

	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error decoding form data.",
			},
		}
		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	subject := "Contact Form: YD Vending"
	recipients := []string{constants.DavidEmail, constants.YovaEmail}
	templateFile := constants.PARTIAL_TEMPLATES_DIR + "contact_form_email.html"

	template, err := helpers.BuildStringFromTemplate(templateFile, "email", form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error building e-mail template.",
			},
		}
		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)
	err = services.SendGmail(recipients, subject, form.Email, body)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to send message.",
			},
		}
		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Sent!",
			"AlertMessage": "We've received your message and will be quick to respond.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLogin(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "login.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}

	csrfSecret, ok := r.Context().Value("csrf_secret").(string)
	if !ok {
		http.Error(w, "Error retrieving external id in login page.", http.StatusInternalServerError)
		return
	}

	session, err := database.GetSession(csrfSecret)
	if err != nil {
		http.Error(w, "Error trying to get session in login page.", http.StatusInternalServerError)
		return
	}

	if session.UserID > 0 {
		user, err := database.GetUserById(session.UserID)
		if err != nil {
			http.Error(w, "Error trying to get existing user from DB.", http.StatusInternalServerError)
			return
		}

		if user.IsAdmin {
			http.Redirect(w, r, "/crm/dashboard", http.StatusSeeOther)
			return
		}
	}

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
	data.PageTitle = "Login — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = helpers.ServeContent(w, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	username := r.Form.Get("username")
	password := r.Form.Get("password")

	// Error handling
	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "error",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
		Data:         map[string]any{},
	}

	user, err := database.GetUserByUsername(username)
	if err != nil {
		tmplCtx.Data["Message"] = "Invalid username."

		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
		}
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	isValid := helpers.ValidatePassword(password, user.Password)
	if !isValid {
		tmplCtx.Data["Message"] = "Invalid password."

		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
		}
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	session, err := sessions.Get(r)
	if err != nil {
		tmplCtx.Data["Message"] = "Could not retrieve session."

		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
		}
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	session.UserID = user.UserID
	err = sessions.Update(session)
	if err != nil {
		tmplCtx.Data["Message"] = "Could not update session."

		token, err := helpers.GenerateTokenInHeader(w, r)
		if err == nil {
			w.Header().Set("X-Csrf-Token", token)
		}
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	sessions.SetCookie(w, time.Now().Add(24*time.Hour), session.CSRFSecret)

	w.WriteHeader(http.StatusOK)
}

func PostLogout(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {

	sessions.SetCookie(w, time.Now().Add(-1*time.Hour), "")

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
