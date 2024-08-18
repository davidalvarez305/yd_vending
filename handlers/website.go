package handlers

import (
	"encoding/hex"
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
		PageTitle:         constants.CompanyName,
		MetaDescription:   "Get a quote for vending machine services.",
		SiteName:          constants.SiteName,
		StaticPath:        constants.StaticPath,
		MediaPath:         constants.MediaPath,
		PhoneNumber:       constants.CompanyPhoneNumber,
		CurrentYear:       time.Now().Year(),
		GoogleAnalyticsID: constants.GoogleAnalyticsID,
		FacebookDataSetID: constants.FacebookDatasetID,
		CompanyName:       constants.CompanyName,
	}
}

func WebsiteHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createWebsiteContext()
	ctx.PagePath = constants.RootDomain + r.URL.Path

	externalId, ok := r.Context().Value("external_id").(string)
	if !ok {
		http.Error(w, "Error retrieving session in context.", http.StatusInternalServerError)
		return
	}

	ctx.ExternalID = externalId

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/quote":
			GetQuoteForm(w, r, ctx)
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

func GetHome(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "home.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data.PageTitle = "Miami Vending Services — " + constants.CompanyName
	data.Nonce = nonce
	data.Features = []string{
		"Innovative Payment Options",
		"24/7 Customer Support",
		"Regular Product Rotation",
		"Health-Conscious Choices",
		"Advanced Security Features",
		"Energy Efficiency",
		"Detailed Reporting and Analytics",
		"Local Sourcing Partnerships",
		"Flexible Contract Terms",
	}

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

func GetQuoteForm(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "quote.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}

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

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
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
	data.PageTitle = "Request A Quote — " + constants.CompanyName
	data.Nonce = nonce
	data.CSRFToken = csrfToken
	data.VendingTypes = vendingTypes
	data.VendingLocations = vendingLocations
	data.Cities = cities

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = helpers.ServeContent(w, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
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
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var form types.QuoteForm
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
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
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

	decodedSecret, err := hex.DecodeString(session.CSRFSecret)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to decode csrf secret.",
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

	if session.FacebookClickID != "" {
		form.FacebookClickID = &session.FacebookClickID
	}

	if session.FacebookClientID != "" {
		form.FacebookClientID = &session.FacebookClientID
	}

	if session.GoogleClientID != "" {
		form.GoogleClientID = &session.GoogleClientID
	}

	if session.ExternalID != "" {
		form.ExternalID = &session.ExternalID
	}

	if len(decodedSecret) > 0 {
		form.CSRFSecret = &decodedSecret
	}

	err = database.CreateLeadAndMarketing(form)
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
			City:            helpers.HashString(helpers.SafeString(form.CityString)),
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

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Awesome!",
			"AlertMessage": "We received your request and will be right with you.",
		},
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
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	subject := "Contact Form: YD Vending"
	recipient := constants.DavidEmail
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
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)
	err = services.SendGmail(recipient, subject, form.Email, body)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to send message.",
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
			"AlertHeader":  "Sent!",
			"AlertMessage": "We've received your message and will be quick to respond.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLogin(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {
	fileName := "login.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}

	session, err := sessions.Get(r)
	if err != nil {
		http.Error(w, "Error trying to retrieve session.", http.StatusInternalServerError)
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

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// Error handling
	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "error",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
		Data:         map[string]any{},
	}

	user, err := database.GetUserByEmail(email)
	if err != nil {
		tmplCtx.Data["Message"] = "Invalid e-mail."
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	isValid := helpers.ValidatePassword(password, user.Password)
	if !isValid {
		tmplCtx.Data["Message"] = "Invalid password."
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	session, err := sessions.Get(r)
	if err != nil {
		tmplCtx.Data["Message"] = "Could not retrieve session."
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	session.UserID = user.UserID
	err = sessions.Update(session)
	if err != nil {
		tmplCtx.Data["Message"] = "Could not update session."
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	sessions.SetCookie(w, time.Now().Add(24*time.Hour), session.CSRFSecret)

	w.WriteHeader(http.StatusOK)
}

func PostLogout(w http.ResponseWriter, r *http.Request, ctx types.WebsiteContext) {

	sessions.SetCookie(w, time.Now().Add(-1*time.Hour), "")

	w.WriteHeader(http.StatusOK)
}
