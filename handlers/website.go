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
	"github.com/davidalvarez305/yd_vending/types"
	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

var websiteBaseFilePath = constants.WEBSITE_TEMPLATES_DIR + "base.html"
var websiteFooterFilePath = constants.WEBSITE_TEMPLATES_DIR + "footer.html"

func createWebsiteContext() map[string]any {
	return map[string]any{
		"PageTitle":         "Request Quote",
		"MetaDescription":   "Get a quote for vending machine services.",
		"SiteName":          constants.SiteName,
		"PagePath":          "http://localhost/quote",
		"StaticPath":        "/static",
		"PhoneNumber":       constants.DavidPhoneNumber,
		"CurrentYear":       time.Now().Year(),
		"GoogleAnalyticsID": constants.GoogleAnalyticsID,
		"FacebookPixelID":   constants.FacebookPixelID,
		"CompanyName":       constants.CompanyName,
	}
}

func WebsiteHandler(w http.ResponseWriter, r *http.Request) {

	ctx := createWebsiteContext()

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
		case "/lp":
			GetLP(w, r, ctx)
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

func GetHome(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "home.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	googleUserId := helpers.GetSessionValueByKey(r, "google_user_id")

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["GoogleUserID"] = googleUserId
	data["Features"] = []string{
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

func GetAbout(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "about.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	googleUserId := helpers.GetSessionValueByKey(r, "google_user_id")

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["GoogleUserID"] = googleUserId

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetPrivacyPolicy(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "privacy.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	googleUserId := helpers.GetSessionValueByKey(r, "google_user_id")

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["GoogleUserID"] = googleUserId

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetTermsAndConditions(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "terms.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	googleUserId := helpers.GetSessionValueByKey(r, "google_user_id")

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["GoogleUserID"] = googleUserId

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetLP(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "lp.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	googleUserId := helpers.GetSessionValueByKey(r, "google_user_id")

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["GoogleUserID"] = googleUserId

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetQuoteForm(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
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

	googleUserId := helpers.GetSessionValueByKey(r, "google_user_id")

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["VendingTypes"] = vendingTypes
	data["VendingLocations"] = vendingLocations
	data["Cities"] = cities
	data["GoogleUserID"] = googleUserId

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = helpers.ServeContent(w, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostQuote(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error parsing form data.", http.StatusBadRequest)
		return
	}

	var form types.QuoteForm
	err = decoder.Decode(&form, r.PostForm)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error decoding form data.", http.StatusBadRequest)
		return
	}

	csrfSecret, err := helpers.GetTokenFromSession(r)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user token.", http.StatusBadRequest)
		return
	}

	googleUserID := helpers.GetSessionValueByKey(r, "google_user_id")
	googleClientID := helpers.GetSessionValueByKey(r, "google_client_id")
	fbClickID := helpers.GetSessionValueByKey(r, "facebook_click_id")
	fbClientID := helpers.GetSessionValueByKey(r, "facebook_client_id")

	// User Marketing Variables
	form.UserAgent = r.Header.Get("User-Agent")
	form.IP = helpers.GetUserIPFromRequest(r)
	form.FacebookClickID = fbClickID
	form.FacebookClientID = fbClientID
	form.GoogleClientID = googleClientID
	form.GoogleUserID = googleUserID
	form.CSRFSecret = csrfSecret

	err = database.CreateLeadAndMarketing(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error creating lead and marketing data.", http.StatusInternalServerError)
		return
	}

	fbEvent := conversions.FacebookEventData{
		EventName:      "quote",
		EventTime:      time.Now().Unix(),
		ActionSource:   "Web",
		EventSourceURL: r.URL.String(),
		UserData: conversions.FacebookUserData{
			FirstName:       helpers.HashString(form.FirstName),
			LastName:        helpers.HashString(form.LastName),
			Phone:           helpers.HashString(form.PhoneNumber),
			FBC:             fbClickID,
			FBP:             fbClientID,
			ClientIPAddress: form.IP,
			ClientUserAgent: form.UserAgent,
		},
	}

	metaPayload := conversions.FacebookPayload{
		Data: []conversions.FacebookEventData{fbEvent},
	}

	payload := conversions.GooglePayload{
		ClientID: googleClientID,
		UserId:   googleUserID,
		Events: []conversions.GoogleEventLead{
			{
				Name: "quote",
				Params: conversions.GoogleEventParamsLead{
					GCLID: form.GCLID,
				},
			},
		},
	}

	// Send conversion events
	go conversions.SendGoogleConversion(payload)
	go conversions.SendFacebookConversion(metaPayload)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, constants.PARTIAL_TEMPLATES_DIR+"modal.html")
}

func GetContactForm(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "contact_form.html"
	files := []string{websiteBaseFilePath, websiteFooterFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}

	googleUserID := helpers.GetSessionValueByKey(r, "google_user_id")

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
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["GoogleUserID"] = googleUserID

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := helpers.ServeContent(w, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostContactForm(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	var form types.ContactForm
	err := decoder.Decode(&form, r.PostForm)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error decoding form data.", http.StatusInternalServerError)
		return
	}

	subject := "Contact Form: YD Vending"
	recipient := constants.GmailEmail
	templateFile := constants.PARTIAL_TEMPLATES_DIR + "contact_form_email.html"

	template, err := helpers.BuildStringFromTemplate(templateFile, "email", form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error building e-mail template.", http.StatusInternalServerError)
		return
	}

	body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)
	if err := services.SendGmail(recipient, subject, body); err != nil {
		fmt.Printf("Error sending email: %v\n", err)
		http.Error(w, "Failed to send message.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PARTIAL_TEMPLATES_DIR+"modal.html")
}

func GetLogin(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "login.html"
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

	googleUserID := helpers.GetSessionValueByKey(r, "google_user_id")

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["GoogleUserID"] = googleUserID

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := helpers.ServeContent(w, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := database.GetUserByEmail(email)
	if err != nil {
		http.Error(w, "Email not found.", http.StatusBadRequest)
		return
	}

	isValid := helpers.ValidatePassword(password, user.Password)
	if !isValid {
		http.Error(w, "Invalid password.", http.StatusBadRequest)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     constants.CookieName,
		Value:    email,
		Path:     "/",
		Domain:   constants.RootDomain,
		Expires:  time.Now().Add(24 * time.Hour), // Expires in 24 hours
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	http.Redirect(w, r, "/crm", http.StatusFound)
}

func PostLogout(w http.ResponseWriter, r *http.Request, ctx map[string]any) {

	http.SetCookie(w, &http.Cookie{
		Name:     constants.CookieName,
		Value:    "",
		Path:     "/",
		Domain:   constants.RootDomain,
		Expires:  time.Now().Add(-1 * time.Hour), // Set expiration time to a past date
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
