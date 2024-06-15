package handlers

import (
	"fmt"
	"log"
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

var baseFilePath = constants.WEBSITE_TEMPLATES_DIR + "base.html"
var footerFilePath = constants.WEBSITE_TEMPLATES_DIR + "footer.html"
var decoder = schema.NewDecoder()

var websiteContext = map[string]any{
	"PageTitle":         "Request Quote",
	"MetaDescription":   "Get a quote for vending machine services.",
	"SiteName":          "YD Vending",
	"PagePath":          "http://localhost/quote",
	"StaticPath":        "/static",
	"PhoneNumber":       constants.DavidPhoneNumber,
	"CurrentYear":       time.Now().Year(),
	"GoogleAnalyticsID": constants.GoogleAnalyticsID,
	"FacebookPixelID":   constants.FacebookPixelID,
	"CompanyName":       "YD Vending, LLC",
}

func WebsiteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/quote":
			GetQuoteForm(w, r)
		case "/contact":
			GetContactForm(w, r)
		case "/login":
			GetLogin(w, r)
		case "/lp":
			GetLP(w, r)
		case "/":
			GetHome(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		case "/quote":
			PostQuote(w, r)
		case "/contact":
			PostContactForm(w, r)
		case "/login":
			PostLogin(w, r)
		case "/logout":
			PostLogout(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	fileName := "home.html"
	files := []string{baseFilePath, footerFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	googleUserID := helpers.GetSessionValueByKey(r, "google_user_id")

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["GoogleUserID"] = googleUserID

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := helpers.ServeContent(w, fileName, files, nil)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func GetLP(w http.ResponseWriter, r *http.Request) {
	fileName := "lp.html"
	files := []string{baseFilePath, footerFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	googleUserId := helpers.GetSessionValueByKey(r, "google_user_id")

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["GoogleUserID"] = googleUserId

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, fileName, files, data)
}

func GetQuoteForm(w http.ResponseWriter, r *http.Request) {
	fileName := "quote.html"
	files := []string{baseFilePath, footerFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}

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

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["VendingTypes"] = vendingTypes
	data["VendingLocations"] = vendingLocations
	data["Cities"] = cities
	data["GoogleUserID"] = googleUserId

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err = helpers.ServeContent(w, fileName, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostQuote(w http.ResponseWriter, r *http.Request) {
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

func GetContactForm(w http.ResponseWriter, r *http.Request) {
	fileName := "contact_form.html"
	files := []string{baseFilePath, footerFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}

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

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["GoogleUserID"] = googleUserID

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := helpers.ServeContent(w, fileName, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostContactForm(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	var form types.ContactForm
	err := decoder.Decode(&form, r.PostForm)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error decoding form data.", http.StatusBadRequest)
		return
	}

	subject := "Contact Form: YD Vending"
	senderEmail := form.Email
	recipient := constants.GmailEmail
	templateName := "contact_form_email.html"

	// Send email
	if err := services.SendSMTPEmail(subject, recipient, senderEmail, form, templateName); err != nil {
		log.Printf("Error sending email: %s", err)
		http.Error(w, "Failed to send message.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PARTIAL_TEMPLATES_DIR+"modal.html")
}

func GetLogin(w http.ResponseWriter, r *http.Request) {
	fileName := "login.html"
	files := []string{baseFilePath, footerFilePath, constants.WEBSITE_TEMPLATES_DIR + fileName}

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

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["GoogleUserID"] = googleUserID

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	err := helpers.ServeContent(w, fileName, files, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
	}
}

func PostLogin(w http.ResponseWriter, r *http.Request) {
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

func PostLogout(w http.ResponseWriter, r *http.Request) {

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
