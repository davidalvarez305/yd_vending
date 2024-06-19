package handlers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/types"
	"github.com/gorilla/schema"
)

var crmBaseFilePath = constants.CRM_TEMPLATES_DIR + "base.html"
var crmFooterFilePath = constants.CRM_TEMPLATES_DIR + "footer.html"

func createCrmContext() map[string]any {
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

func CRMHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createCrmContext()

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/crm/dashboard":
			GetDashboard(w, r, ctx)
		case "/crm/leads":
			GetLeads(w, r, ctx)
		case "/crm/machines":
			GetMachines(w, r, ctx)
		case "/crm/tickets":
			GetTickets(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetLeads(w http.ResponseWriter, r *http.Request, ctx map[string]interface{}) {
	fileName := "quotes_table.html"
	baseFile := constants.CRM_TEMPLATES_DIR + "leads.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, constants.PARTIAL_TEMPLATES_DIR + fileName}

	// Retrieve nonce from request context
	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	// Retrieve CSRF token from request context
	csrfToken, ok := r.Context().Value("csrf_token").(string)
	if !ok {
		http.Error(w, "Error retrieving CSRF token.", http.StatusInternalServerError)
		return
	}

	// Parse query parameters into GetQuotesParams struct
	var params types.GetLeadsParams
	err := schema.NewDecoder().Decode(&params, r.URL.Query())
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error parsing query parameters.", http.StatusBadRequest)
		return
	}

	quotes, totalRows, err := database.GetLeadList(params)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting quotes from DB.", http.StatusInternalServerError)
		return
	}

	if len(r.URL.RawQuery) > 0 {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "quotes_table",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "quotes_table.html",
			Data: map[string]any{
				"Quotes":      quotes,
				"TotalRows":   totalRows,
				"CurrentPage": params.PageNum,
			},
		}
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pages := helpers.GenerateSequence(1, totalRows)

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Quotes"] = quotes
	data["Pages"] = pages

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetMachines(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "machines.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName}
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetTickets(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "tickets.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName}
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetDashboard(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "dashboard.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName}
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

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetLeadDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "lead_detail.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName}
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

	leadId, err := helpers.GetLeadIDFromURLPath(r.URL.Path)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Invalid URL path.", http.StatusInternalServerError)
		return
	}

	leadDetails, err := database.GetLeadDetails(leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead details from DB.", http.StatusInternalServerError)
		return
	}

	messages, err := database.GetMessagesByLeadID(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting quotes from DB.", http.StatusInternalServerError)
		return
	}

	userId, err := helpers.GetUserIDFromSession(r)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user ID from session.", http.StatusInternalServerError)
		return
	}

	phoneNumber, err := database.GetPhoneNumberFromUserID(userId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting phone number from user ID.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PagePath"] = "http://localhost" + r.URL.Path
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Lead"] = leadDetails
	data["Messages"] = messages
	data["CRMUserPhoneNumber"] = phoneNumber

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}
