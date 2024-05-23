package handlers

import (
	"net/http"
	"time"

	"github.com/davidalvarez305/budgeting/constants"
	"github.com/davidalvarez305/budgeting/helpers"
	"github.com/davidalvarez305/budgeting/models"
)

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
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetHome(w http.ResponseWriter, r *http.Request) {
	fileName := "home.html"

	err := helpers.BuildFile(fileName, constants.PUBLIC_DIR+fileName, constants.TEMPLATES_DIR+fileName, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PUBLIC_DIR+fileName)
}

func GetQuoteForm(w http.ResponseWriter, r *http.Request) {
	fileName := "quote.html"

	// @TODO: Build core data that will be in all "website handler" pages.
	data := struct {
		PageTitle         string
		MetaDescription   string
		SiteName          string
		PagePath          string
		StaticPath        string
		PhoneNumber       string
		CurrentYear       int
		GoogleAnalyticsID string
	}{
		PageTitle:         "Request Quote",
		MetaDescription:   "Get a quote for vending machine services.",
		SiteName:          "YD Vending",
		PagePath:          "http://localhost" + r.URL.Path,
		StaticPath:        "../static",
		PhoneNumber:       "(123) - 456 7890",
		CurrentYear:       time.Now().Year(),
		GoogleAnalyticsID: "G-1231412312",
	}

	err := helpers.BuildFile(fileName, constants.WEBSITE_PUBLIC_DIR+fileName, constants.WEBSITE_TEMPLATES_DIR+fileName, data)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.WEBSITE_PUBLIC_DIR+fileName)
}

func PostQuote(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		http.Error(w, "Form cannot be parsed.", http.StatusBadRequest)
		return
	}

	/* transaction, err := helpers.ParseTransaction(r.Form)

	if err != nil {
		http.Error(w, "Error parsing transaction.", http.StatusInternalServerError)
		return
	} */

	var lead models.Lead

	fileName := "form.html"

	err = helpers.BuildFile(fileName, constants.PUBLIC_DIR+fileName, constants.TEMPLATES_DIR+fileName, lead)

	if err != nil {
		http.Error(w, "Error building HTML file.", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, constants.PUBLIC_DIR+fileName)
}
