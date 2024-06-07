package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/middleware"
	"github.com/davidalvarez305/yd_vending/services"
	"github.com/davidalvarez305/yd_vending/types"
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
		case "/contact":
			GetContactForm(w, r)
		case "/login":
			GetLogin(w, r)
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

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = r.Context().Value("nonce").(string)
	data["CSRFToken"] = r.Context().Value("csrf_token").(string)
	data["VendingTypes"] = vendingTypes
	data["VendingLocations"] = vendingLocations
	data["Cities"] = cities

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
		http.Error(w, "Error getting user token.", http.StatusBadRequest)
		return
	}

	err = database.CreateLeadAndMarketing(form, token)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user token.", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PARTIAL_TEMPLATES_DIR+"modal.html")
}

func GetContactForm(w http.ResponseWriter, r *http.Request) {
	fileName := "contact_form.html"

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = r.Context().Value("nonce").(string)
	data["CSRFToken"] = r.Context().Value("csrf_token").(string)

	err := helpers.BuildFile(fileName, baseFilePath, footerFilePath, constants.WEBSITE_PUBLIC_DIR+fileName, constants.WEBSITE_TEMPLATES_DIR+fileName, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error building quote form.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.WEBSITE_PUBLIC_DIR+fileName)
}

func PostContactForm(w http.ResponseWriter, r *http.Request) {
	// Parse form data
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form data.", http.StatusBadRequest)
		return
	}

	var form types.ContactForm

	err := json.NewDecoder(r.Body).Decode(&form)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error decoding JSON.", http.StatusBadRequest)
		return
	}

	subject := "Contact Form: YD Vending"
	senderEmail := os.Getenv("GMAIL_EMAIL")
	recipient := form.Email
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

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = r.Context().Value("nonce").(string)
	data["CSRFToken"] = r.Context().Value("csrf_token").(string)

	err := helpers.BuildFile(fileName, baseFilePath, footerFilePath, constants.WEBSITE_PUBLIC_DIR+fileName, constants.WEBSITE_TEMPLATES_DIR+fileName, data)

	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error building login form.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.WEBSITE_PUBLIC_DIR+fileName)
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
		Name:     os.Getenv("COOKIE_NAME"),
		Value:    email,
		Path:     "/",
		Domain:   os.Getenv("ROOT_DOMAIN"),
		Expires:  time.Now().Add(24 * time.Hour), // Expires in 24 hours
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	http.Redirect(w, r, "/crm", http.StatusFound)
}

func PostLogout(w http.ResponseWriter, r *http.Request) {

	http.SetCookie(w, &http.Cookie{
		Name:     os.Getenv("COOKIE_NAME"),
		Value:    "",
		Path:     "/",
		Domain:   os.Getenv("ROOT_DOMAIN"),
		Expires:  time.Now().Add(-1 * time.Hour), // Set expiration time to a past date
		HttpOnly: false,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}
