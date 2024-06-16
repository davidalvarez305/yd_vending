package handlers

import (
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/helpers"
)

var crmBaseFilePath = constants.CRM_TEMPLATES_DIR + "base.html"
var crmFooterFilePath = constants.CRM_TEMPLATES_DIR + "footer.html"

var crmContext = map[string]any{
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

func CRMHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/crm/leads":
			GetLeads(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetLeads(w http.ResponseWriter, r *http.Request) {
	fileName := "leads.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName}
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	data := websiteContext
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, fileName, files, data)
}
