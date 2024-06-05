package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/davidalvarez305/budgeting/constants"
	"github.com/davidalvarez305/budgeting/helpers"
	"github.com/davidalvarez305/budgeting/middleware"
	"github.com/davidalvarez305/budgeting/models"
	"github.com/davidalvarez305/budgeting/services"
	"github.com/davidalvarez305/budgeting/types"
)

var baseFilePath = constants.WEBSITE_TEMPLATES_DIR + "base.html"
var footerFilePath = constants.WEBSITE_TEMPLATES_DIR + "footer.html"

var websiteContext = map[string]any{
	"PageTitle":         "Request Quote",
	"MetaDescription":   "Get a quote for vending machine services.",
	"SiteName":          "YD Vending",
	"PagePath":          "http://localhost/quote",
	"StaticPath":        "/static",
	"PhoneNumber":       "(123) - 456 7890",
	"CurrentYear":       time.Now().Year(),
	"GoogleAnalyticsID": "G-1231412312",
}

func WebsiteHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/quote":
			GetQuoteForm(w, r)
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
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	fileName := "home.html"

	err := helpers.BuildFile(fileName, baseFilePath, footerFilePath, constants.PUBLIC_DIR+fileName, constants.TEMPLATES_DIR+fileName, nil)

	if err != nil {
		http.Error(w, "Error building home page.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PUBLIC_DIR+fileName)
}

func GetQuoteForm(w http.ResponseWriter, r *http.Request) {
	fileName := "quote.html"

	token, err := middleware.GetTokenFromSession(r)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user token from session.", http.StatusBadRequest)
		return
	}

	// Only do this on first user session
	csrfToken, err := middleware.Encrypt(time.Now().Unix(), token)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error generating CSRF token.", http.StatusInternalServerError)
		return
	}

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = r.Context().Value("nonce").(string)
	data["CSRFToken"] = csrfToken

	err = helpers.BuildFile(fileName, baseFilePath, footerFilePath, constants.WEBSITE_PUBLIC_DIR+fileName, constants.WEBSITE_TEMPLATES_DIR+fileName, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error building quote form.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.WEBSITE_PUBLIC_DIR+fileName)
}

func PostQuote(w http.ResponseWriter, r *http.Request) {
	var form types.QuoteForm

	err := json.NewDecoder(r.Body).Decode(&form)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error decoding JSON.", http.StatusBadRequest)
		return
	}

	token, err := middleware.GetTokenFromSession(r)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user token from session.", http.StatusBadRequest)
		return
	}

	err = middleware.ValidateCSRFToken(form.CSRFToken, token)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error validating token.", http.StatusBadRequest)
		return
	}

	lead := &models.Lead{
		FirstName:         form.FirstName,
		LastName:          form.LastName,
		PhoneNumber:       form.PhoneNumber,
		CreatedAt:         time.Now().Unix(),
		Rent:              form.Rent,
		VendingTypeID:     form.MachineType,
		CityID:            form.City,
		VendingLocationID: form.LocationType,
		LeadMarketing: &models.LeadMarketing{
			Source:        form.Source,
			Medium:        form.Medium,
			Channel:       form.Channel,
			LandingPage:   form.LandingPage,
			Keyword:       form.Keyword,
			Referrer:      form.Referrer,
			Gclid:         form.Gclid,
			CampaignID:    form.CampaignID,
			AdCampaign:    form.AdCampaign,
			AdGroupID:     form.AdGroupID,
			AdGroupName:   form.AdGroupName,
			AdSetID:       form.AdSetID,
			AdSetName:     form.AdSetName,
			AdID:          form.AdID,
			AdHeadline:    form.AdHeadline,
			Language:      form.Language,
			OS:            form.OS,
			UserAgent:     form.UserAgent,
			ButtonClicked: form.ButtonClicked,
			DeviceType:    form.DeviceType,
			IP:            form.IP,
		},
	}

	fmt.Printf("%+v\n", lead)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PARTIAL_TEMPLATES_DIR+"modal.html")
}

func PostContactForm(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	// Extract form values
	firstName := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	email := r.FormValue("email")
	message := r.FormValue("message")

	// Compose email message
	subject := "Contact Form: YD Vending"
	body := fmt.Sprintf("Name: %s %s\nEmail: %s\nMessage:\n%s", firstName, lastName, email, message)

	// Send email
	if err := services.SendSMTPEmail(subject, body, email); err != nil {
		log.Printf("Error sending email: %s", err)
		http.Error(w, "Failed to send message.", http.StatusInternalServerError)
		return
	}

	// Respond to client
	fmt.Fprintf(w, "Message sent successfully!")
}
