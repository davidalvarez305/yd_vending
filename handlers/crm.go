package handlers

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/models"
	"github.com/davidalvarez305/yd_vending/services"
	"github.com/davidalvarez305/yd_vending/sessions"
	"github.com/davidalvarez305/yd_vending/types"
	"github.com/davidalvarez305/yd_vending/utils"
	"github.com/google/uuid"
)

var crmBaseFilePath = constants.CRM_TEMPLATES_DIR + "base.html"
var crmFooterFilePath = constants.CRM_TEMPLATES_DIR + "footer.html"

func createCrmContext() map[string]any {
	return map[string]any{
		"PageTitle":       constants.CompanyName,
		"MetaDescription": "Get a quote for vending machine services.",
		"SiteName":        constants.SiteName,
		"StaticPath":      constants.StaticPath,
		"MediaPath":       constants.MediaPath,
		"PhoneNumber":     constants.DavidPhoneNumber,
		"CurrentYear":     time.Now().Year(),
		"CompanyName":     constants.CompanyName,
	}
}

func CRMHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createCrmContext()
	ctx["PagePath"] = constants.RootDomain + r.URL.Path
	path := r.URL.Path

	stats, err := database.GetDashboardStats()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting dashboard stats from DB.", http.StatusInternalServerError)
		return
	}
	ctx["DashboardStats"] = stats

	switch r.Method {
	case http.MethodGet:
		parts := strings.Split(path, "/")
		if strings.HasPrefix(path, "/crm/business/") {
			if len(parts) >= 6 && parts[4] == "location" && helpers.IsNumeric(parts[3]) && helpers.IsNumeric(parts[5]) {
				GetLocationDetail(w, r, ctx)
				return
			}
		}

		if strings.HasPrefix(path, "/crm/lead/") && strings.Contains(path, "/messages") {
			GetLeadMessagesPartial(w, r)
			return
		}

		if strings.HasPrefix(path, "/crm/lead/") && strings.Contains(path, "/images") {
			GetLeadImagesPartial(w, r)
			return
		}

		if strings.HasPrefix(path, "/crm/lead/") {
			if len(path) > len("/crm/lead/") && helpers.IsNumeric(path[len("/crm/lead/"):]) {
				GetLeadDetail(w, r, ctx)
				return
			}
			return
		}

		if strings.HasPrefix(path, "/crm/business/") {
			if len(path) > len("/crm/business/") && helpers.IsNumeric(path[len("/crm/business/"):]) {
				GetBusinessDetail(w, r, ctx)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/machine/") {
			if len(path) > len("/crm/machine/") && helpers.IsNumeric(path[len("/crm/machine/"):]) {
				GetMachineDetail(w, r, ctx)
				return
			}
			if len(parts) >= 6 && parts[4] == "slot" && helpers.IsNumeric(parts[3]) && helpers.IsNumeric(parts[5]) {
				GetSlotDetail(w, r, ctx)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/supplier/") {
			if len(path) > len("/crm/supplier/") && helpers.IsNumeric(path[len("/crm/supplier/"):]) {
				GetSupplierDetail(w, r, ctx)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/vendor/") {
			if len(path) > len("/crm/vendor/") && helpers.IsNumeric(path[len("/crm/vendor/"):]) {
				GetVendorDetail(w, r, ctx)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/product-slot-assignment/") {
			if len(path) > len("/crm/product-slot-assignment/") && helpers.IsNumeric(path[len("/crm/product-slot-assignment/"):]) {
				GetProductSlotAssignmentDetail(w, r, ctx)
				return
			}
		}

		switch path {
		case "/crm/dashboard":
			GetDashboard(w, r, ctx)
		case "/crm/lead":
			GetLeads(w, r, ctx)
		case "/crm/machine":
			GetMachines(w, r, ctx)
		case "/crm/business":
			GetBusinesses(w, r, ctx)
		case "/crm/vendor":
			GetVendors(w, r, ctx)
		case "/crm/supplier":
			GetSuppliers(w, r, ctx)
		case "/crm/location":
			GetLocation(w, r, ctx)
		case "/crm/ticket":
			GetTickets(w, r, ctx)
		case "/crm/upload-images":
			GetImagesUpload(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPut:
		parts := strings.Split(path, "/")
		if strings.HasPrefix(path, "/crm/lead/") {
			if len(parts) >= 5 && parts[4] == "marketing" && helpers.IsNumeric(parts[3]) {
				PutLeadMarketing(w, r)
				return
			}
			if len(path) > len("/crm/lead/") && helpers.IsNumeric(path[len("/crm/lead/"):]) {
				PutLead(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/business/") {
			parts := strings.Split(path, "/")
			if len(parts) >= 6 && parts[4] == "location" && helpers.IsNumeric(parts[3]) && helpers.IsNumeric(parts[5]) {
				PutLocation(w, r)
				return
			}
			return
		}

		if len(path) > len("/crm/product-slot-assignment/") && helpers.IsNumeric(path[len("/crm/product-slot-assignment/"):]) {
			PutProductSlotAssignment(w, r)
			return
		}

		if strings.HasPrefix(path, "/crm/business/") {
			if len(path) > len("/crm/business/") && helpers.IsNumeric(path[len("/crm/business/"):]) {
				PutBusiness(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/machine/") {
			if len(path) > len("/crm/machine/") && helpers.IsNumeric(path[len("/crm/machine/"):]) {
				PutMachine(w, r)
				return
			}
			if len(parts) >= 6 && parts[4] == "slot" && helpers.IsNumeric(parts[3]) && helpers.IsNumeric(parts[5]) {
				PutSlot(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/supplier/") {
			if len(path) > len("/crm/supplier/") && helpers.IsNumeric(path[len("/crm/supplier/"):]) {
				PutSupplier(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/vendor/") {
			if len(path) > len("/crm/vendor/") && helpers.IsNumeric(path[len("/crm/vendor/"):]) {
				PutVendor(w, r)
				return
			}
		}
		if len(path) > len("/crm/slot-price-log/") && helpers.IsNumeric(path[len("/crm/slot-price-log/"):]) {
			PutSlotPriceLog(w, r)
			return
		}

		switch path {
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		if strings.HasPrefix(path, "/crm/lead/") && strings.Contains(path, "/images") {
			PostLeadImages(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/lead/") && strings.Contains(path, "/notes") {
			PostLeadNotes(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/slot") && strings.Contains(path, "/refill") {
			PostRefill(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/refill") {
			PostRefillAll(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/business/") && strings.Contains(path, "/location") {
			PostLocation(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/slot") {
			PostSlot(w, r)
			return
		}

		switch path {
		case "/crm/business":
			PostBusiness(w, r)
		case "/crm/machine":
			PostMachine(w, r)
		case "/crm/vendor":
			PostVendor(w, r)
		case "/crm/supplier":
			PostSupplier(w, r)
		case "/crm/upload-images":
			PostImagesUpload(w, r)
		case "/crm/product-slot-assignment":
			PostProductSlotAssignment(w, r)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodDelete:
		parts := strings.Split(path, "/")
		if strings.HasPrefix(path, "/crm/business/") {
			if len(parts) >= 6 && parts[4] == "location" && helpers.IsNumeric(parts[3]) && helpers.IsNumeric(parts[5]) {
				DeleteLocation(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/business/") {
			if len(path) > len("/crm/business/") && helpers.IsNumeric(path[len("/crm/business/"):]) {
				DeleteBusiness(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/machine/") {
			if len(path) > len("/crm/machine/") && helpers.IsNumeric(path[len("/crm/machine/"):]) {
				DeleteMachine(w, r)
				return
			}
			if len(parts) >= 6 && parts[4] == "slot" && helpers.IsNumeric(parts[3]) && helpers.IsNumeric(parts[5]) {
				DeleteSlot(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/supplier/") {
			if len(path) > len("/crm/supplier/") && helpers.IsNumeric(path[len("/crm/supplier/"):]) {
				DeleteSupplier(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/vendor/") {
			if len(path) > len("/crm/vendor/") && helpers.IsNumeric(path[len("/crm/vendor/"):]) {
				DeleteVendor(w, r)
				return
			}
		}
		if len(path) > len("/crm/product-slot-assignment/") && helpers.IsNumeric(path[len("/crm/product-slot-assignment/"):]) {
			DeleteProductSlotAssigment(w, r)
			return
		}
		if len(path) > len("/crm/slot-price-log/") && helpers.IsNumeric(path[len("/crm/slot-price-log/"):]) {
			DeleteSlotPriceLog(w, r)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetLeads(w http.ResponseWriter, r *http.Request, ctx map[string]interface{}) {
	baseFile := constants.CRM_TEMPLATES_DIR + "leads.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile}

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

	var params types.GetLeadsParams
	params.LocationType = helpers.SafeStringToPointer(r.URL.Query().Get("location_type"))
	params.VendingType = helpers.SafeStringToPointer(r.URL.Query().Get("vending_type"))
	params.PageNum = helpers.SafeStringToPointer(r.URL.Query().Get("page_num"))

	leads, totalRows, err := database.GetLeadList(params)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting leads from DB.", http.StatusInternalServerError)
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

	data := ctx
	data["PageTitle"] = "Leads — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Leads"] = leads
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["VendingTypes"] = vendingTypes
	data["VendingLocations"] = vendingLocations

	data["CurrentPage"] = 1
	if params.PageNum != nil {
		data["CurrentPage"] = *params.PageNum
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetMachines(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "machines.html"
	createMachineForm := constants.CRM_TEMPLATES_DIR + "create_machine_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "machines_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createMachineForm}

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

	pageNum := 1
	hasPageNum := r.URL.Query().Has("page_num")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("page_num"))
		if err == nil && num > 1 {
			pageNum = num
		}
	}

	machines, totalRows, err := database.GetMachineList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machines from DB.", http.StatusInternalServerError)
		return
	}

	locations, err := database.GetLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting locations from DB.", http.StatusInternalServerError)
		return
	}

	vendingTypes, err := database.GetVendingTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending types from DB.", http.StatusInternalServerError)
		return
	}

	machineStatuses, err := database.GetMachineStatuses()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machine statuses from DB.", http.StatusInternalServerError)
		return
	}

	vendors, err := database.GetVendors()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vendors from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Machines — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Machines"] = machines
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum
	data["Locations"] = locations
	data["VendingTypes"] = vendingTypes
	data["MachineStatuses"] = machineStatuses
	data["Vendors"] = vendors

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetBusinesses(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "businesses.html"
	createBusinessForm := constants.CRM_TEMPLATES_DIR + "create_business_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "businesses_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createBusinessForm}

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

	pageNum := 1
	hasPageNum := r.URL.Query().Has("page_num")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("page_num"))
		if err == nil && num > 1 {
			pageNum = num
		}
	}

	businesses, totalRows, err := database.GetBusinessList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting businesses from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Businesses — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Businesses"] = businesses
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetLocation(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
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
	data["PageTitle"] = "Machines — " + constants.CompanyName

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
	data["PageTitle"] = "Dashboard — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetLeadDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "lead_detail.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, constants.PARTIAL_TEMPLATES_DIR + "messages.html", constants.PARTIAL_TEMPLATES_DIR + "notes.html", constants.PARTIAL_TEMPLATES_DIR + "lead_images.html"}
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

	leadId := strings.TrimPrefix(r.URL.Path, "/crm/lead/")

	leadDetails, err := database.GetLeadDetails(leadId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead details from DB.", http.StatusInternalServerError)
		return
	}

	messages, err := database.GetMessagesByLeadID(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead messages from DB.", http.StatusInternalServerError)
		return
	}

	leadNotes, err := database.GetLeadNotesByLeadID(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead notes from DB.", http.StatusInternalServerError)
		return
	}

	leadImages, err := database.GetLeadImagesByLeadID(leadDetails.LeadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead images from DB.", http.StatusInternalServerError)
		return
	}

	values, err := sessions.Get(r)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting user ID from session.", http.StatusInternalServerError)
		return
	}

	phoneNumber, err := database.GetPhoneNumberFromUserID(values.UserID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting phone number from user ID.", http.StatusInternalServerError)
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

	data := ctx
	data["PageTitle"] = "Lead Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Lead"] = leadDetails
	data["Messages"] = messages
	data["CRMUserPhoneNumber"] = phoneNumber
	data["VendingTypes"] = vendingTypes
	data["VendingLocations"] = vendingLocations
	data["LeadNotes"] = leadNotes
	data["LeadImagesCount"] = len(leadImages)
	data["LeadImages"] = leadImages

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PutLead(w http.ResponseWriter, r *http.Request) {
	token, err := helpers.GenerateTokenInHeader(w, r)
	if err != nil {
		fmt.Printf("Error generating token: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error generating new token. Reload page.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	w.Header().Set("X-Csrf-Token", token)

	err = r.ParseForm()
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

	var form types.UpdateLeadForm
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

	err = database.UpdateLead(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating lead.",
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
			"AlertMessage": "Lead has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutLeadMarketing(w http.ResponseWriter, r *http.Request) {
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

	var form types.UpdateLeadMarketingForm
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

	err = database.UpdateLeadMarketing(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating lead marketing.",
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
			"AlertMessage": "Lead marketing has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLeadMessagesPartial(w http.ResponseWriter, r *http.Request) {
	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := database.GetMessagesByLeadID(leadId)
	if err != nil {
		fmt.Printf("Error getting messages: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get new messages.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "messages.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "messages.html",
		Data: map[string]any{
			"Messages": messages,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostLeadImages(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(50 << 20) // 50 MB limit
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

	leadIDStr := r.FormValue("lead_id")
	leadID, err := strconv.Atoi(leadIDStr)
	if err != nil {
		fmt.Printf("Error converting lead_id to int: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid lead ID.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	values, err := sessions.Get(r)
	if err != nil || values.UserID == 0 {
		fmt.Printf("COULD NOT GET SESSION WHILE UPLOADING IMAGES: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Could not get session while uploading images.",
			},
		}
		w.WriteHeader(http.StatusForbidden)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	files := r.MultipartForm.File["upload_images"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Printf("%+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to open uploaded file.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
		defer file.Close()

		fileExtension := filepath.Ext(fileHeader.Filename)
		src := uuid.New().String() + fileExtension
		filePath := "images/" + src

		err = services.UploadFileToS3(file, fileHeader.Size, filePath)
		if err != nil {
			fmt.Printf("Failed to upload image to S3: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to upload image to S3.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		dateAdded, err := utils.GetCurrentTimeInEST()
		if err != nil {
			fmt.Printf("Error getting time as EST: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error getting time as EST.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		form := models.LeadImage{
			Src:           src,
			LeadID:        leadID,
			DateAdded:     dateAdded,
			AddedByUserID: values.UserID,
		}

		err = database.CreateLeadImage(form)
		if err != nil {
			fmt.Printf("%+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error saving image metadata.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Images have been successfully uploaded.",
		},
	}
	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostLeadNotes(w http.ResponseWriter, r *http.Request) {
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

	leadIdForm := r.FormValue("lead_id")
	note := r.FormValue("note")

	leadID, err := strconv.Atoi(leadIdForm)
	if err != nil {
		fmt.Printf("Error converting lead_id to int: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid lead ID.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	values, err := sessions.Get(r)
	if err != nil || values.UserID == 0 {
		fmt.Printf("COULD NOT GET SESSION WHILE UPLOADING IMAGES: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Could not get session while uploading images.",
			},
		}
		w.WriteHeader(http.StatusForbidden)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	dateAdded, err := utils.GetCurrentTimeInEST()
	if err != nil {
		fmt.Printf("Error getting time as EST: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting time as EST.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadNote := models.LeadNote{
		LeadID:        leadID,
		Note:          note,
		DateAdded:     dateAdded,
		AddedByUserID: values.UserID,
	}

	err = database.CreateLeadNote(leadNote)
	if err != nil {
		fmt.Printf("Error creating note: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to save note.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	notes, err := database.GetLeadNotesByLeadID(leadID)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get new notes.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "notes.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "notes.html",
		Data: map[string]any{
			"LeadNotes": notes,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLeadImagesPartial(w http.ResponseWriter, r *http.Request) {
	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	images, err := database.GetLeadImagesByLeadID(leadId)
	if err != nil {
		fmt.Printf("Error getting images: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get images.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "lead_images.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "lead_images.html",
		Data: map[string]any{
			"LeadImages": images,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostBusiness(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.BusinessForm
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

	err = database.CreateBusiness(form)
	if err != nil {
		fmt.Printf("Error creating business: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create business.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1 // Always default to one after new business is created
	businesses, totalRows, err := database.GetBusinessList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting businesses from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "businesses_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "businesses_table.html",
		Data: map[string]any{
			"Businesses":  businesses,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostLocation(w http.ResponseWriter, r *http.Request) {
	token, err := helpers.GenerateTokenInHeader(w, r)
	if err != nil {
		fmt.Printf("Error generating token: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error generating new token. Reload page.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	w.Header().Set("X-Csrf-Token", token)

	err = r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.LocationForm
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

	businessId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateLocation(businessId, form)
	if err != nil {
		fmt.Printf("Error creating location: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create location.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "success.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Location created successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostMachine(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.MachineForm
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

	machineId, err := database.CreateMachine(form)
	if err != nil {
		fmt.Printf("Error creating machine: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create machine.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	cardReaderSerialNumber := utils.CreateNullString(form.CardReaderSerialNumber)
	locationId := utils.CreateNullInt(form.LocationID)
	locationDateAssigned := utils.CreateNullInt64(form.LocationDateAssigned)
	machineCardReaderDateAssigned := utils.CreateNullInt64(form.DateAssigned)
	isLocationActive := utils.CreateNullBool(form.IsLocationActive)
	isCardReaderActive := utils.CreateNullBool(form.IsCardReaderActive)

	// Assign machine to location
	if locationId.Valid && locationDateAssigned.Valid && isLocationActive.Valid {
		assignment := models.MachineLocationAssignment{
			LocationID:   int(locationId.Int64),
			MachineID:    machineId,
			DateAssigned: locationDateAssigned.Int64,
			IsActive:     isLocationActive.Bool,
		}
		err = database.CreateMachineLocationAssignment(assignment)
		if err != nil {
			fmt.Printf("Error creating machine location assignemnt: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to create machine location assignemnt.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	// Assign card reader to machine
	if cardReaderSerialNumber.Valid && machineCardReaderDateAssigned.Valid && isCardReaderActive.Valid {
		cardReaderAssignment := models.MachineCardReaderAssignment{
			CardReaderSerialNumber: cardReaderSerialNumber.String,
			MachineID:              machineId,
			DateAssigned:           machineCardReaderDateAssigned.Int64,
			IsActive:               isCardReaderActive.Bool,
		}
		err = database.CreateMachineCardReaderAssignment(cardReaderAssignment)
		if err != nil {
			fmt.Printf("Error creating machine card reader assignemnt: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to create machine card reader assignemnt.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	pageNum := 1 // Always default to one after new business is created
	machines, totalRows, err := database.GetMachineList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting machines from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "machines_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "machines_table.html",
		Data: map[string]any{
			"Machines":    machines,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutBusiness(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.BusinessForm
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

	businessId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateBusiness(businessId, form)
	if err != nil {
		fmt.Printf("Error updating business: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update business.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "success.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Business updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutLocation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.LocationForm
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

	businessId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	locationId, err := helpers.GetSecondIDFromPath(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateLocation(businessId, locationId, form)
	if err != nil {
		fmt.Printf("Error updating location: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update location.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "success.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Location updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutBusinessContact(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.BusinessContactForm
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

	businessId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	businessContactId, err := helpers.GetSecondIDFromPath(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateBusinessContact(businessId, businessContactId, form)
	if err != nil {
		fmt.Printf("Error updating business contact: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update business contact.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "success.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Business contact updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutMachine(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.MachineForm
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

	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateMachine(machineId, form)
	if err != nil {
		fmt.Printf("Error updating machine: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update machine.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	dateAssigned := utils.CreateNullInt64(form.DateAssigned)
	locationDateAssigned := utils.CreateNullInt64(form.LocationDateAssigned)
	cardReaderSerialNumber := utils.CreateNullString(form.CardReaderSerialNumber)
	locationId := utils.CreateNullInt(form.LocationID)
	isCardReaderActive := utils.CreateNullBool(form.IsCardReaderActive)
	isLocationActive := utils.CreateNullBool(form.IsLocationActive)

	machine, err := database.GetMachineDetails(machineId)
	if err != nil {
		fmt.Printf("Error updating machine: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting machine details.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	// Assign machine to location if not equal to current location
	if locationId.Valid && locationId.Int64 != int64(machine.LocationID) && locationDateAssigned.Valid && locationDateAssigned.Int64 != machine.LocationDateAssigned {
		assignment := models.MachineLocationAssignment{
			LocationID:   int(locationId.Int64),
			MachineID:    machineId,
			DateAssigned: locationDateAssigned.Int64,
			IsActive:     isLocationActive.Bool,
		}
		err = database.CreateMachineLocationAssignment(assignment)
		if err != nil {
			fmt.Printf("Error creating machine: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to create machine location assignment.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	// Assign card reader to machine if not equal to current card reader
	if cardReaderSerialNumber.Valid && cardReaderSerialNumber.String != machine.CardReaderSerialNumber && dateAssigned.Valid && dateAssigned.Int64 != machine.DateAssigned {
		cardReaderAssignment := models.MachineCardReaderAssignment{
			CardReaderSerialNumber: cardReaderSerialNumber.String,
			MachineID:              machineId,
			DateAssigned:           dateAssigned.Int64,
			IsActive:               isCardReaderActive.Bool,
		}
		err = database.CreateMachineCardReaderAssignment(cardReaderAssignment)
		if err != nil {
			fmt.Printf("Error creating machine: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to create machine card reader assignment.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "success.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Machine updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetVendors(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "vendors.html"
	createVendorForm := constants.CRM_TEMPLATES_DIR + "create_vendor_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "vendors_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createVendorForm}

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

	pageNum := 1
	hasPageNum := r.URL.Query().Has("page_num")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("page_num"))
		if err == nil && num > 1 {
			pageNum = num
		}
	}

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities from DB.", http.StatusInternalServerError)
		return
	}

	vendors, totalRows, err := database.GetVendorList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vendors from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Vendors — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Vendors"] = vendors
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum
	data["Cities"] = cities

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetCreateVendorForm(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
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

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("Error getting locations: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting cities.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "create_vendor_form.html",
		TemplatePath: constants.CRM_TEMPLATES_DIR + "create_vendor_form.html",
		Data: map[string]any{
			"CSRFToken": csrfToken,
			"Nonce":     nonce,
			"Cities":    cities,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostVendor(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.VendorForm
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

	err = database.CreateVendor(form)
	if err != nil {
		fmt.Printf("Error creating vendor: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create vendor.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1 // Always default to one after new entity is created
	vendors, totalRows, err := database.GetVendorList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting vendors from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "vendors_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "vendors_table.html",
		Data: map[string]any{
			"Vendors":     vendors,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutVendor(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.VendorForm
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

	vendorId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/vendor/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateVendor(vendorId, form)
	if err != nil {
		fmt.Printf("Error updating vendor: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update vendor.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "success.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Vendor updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetSuppliers(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "suppliers.html"
	createSupplierForm := constants.CRM_TEMPLATES_DIR + "create_supplier_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "suppliers_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createSupplierForm}

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

	pageNum := 1
	hasPageNum := r.URL.Query().Has("page_num")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("page_num"))
		if err == nil && num > 1 {
			pageNum = num
		}
	}

	suppliers, totalRows, err := database.GetSupplierList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting suppliers from DB.", http.StatusInternalServerError)
		return
	}

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Suppliers — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Suppliers"] = suppliers
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum
	data["Cities"] = cities

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostSupplier(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.SupplierForm
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

	err = database.CreateSupplier(form)
	if err != nil {
		fmt.Printf("Error creating supplier: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create supplier.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1 // Always default to one after new entity is created
	suppliers, totalRows, err := database.GetSupplierList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting suppliers from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "suppliers_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "suppliers_table.html",
		Data: map[string]any{
			"Suppliers":   suppliers,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutSupplier(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.SupplierForm
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

	supplierId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/supplier/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateSupplier(supplierId, form)
	if err != nil {
		fmt.Printf("Error updating supplier: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update supplier.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "success.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Supplier updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteMachine(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteMachine(machineId)
	if err != nil {
		fmt.Printf("Error deleting machine: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete machine.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	machines, totalRows, err := database.GetMachineList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting machines from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "machines_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "machines_table.html",
		Data: map[string]any{
			"Machines":    machines,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteBusiness(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	businessId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteBusiness(businessId)
	if err != nil {
		fmt.Printf("Error deleting business: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete business.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	businesses, totalRows, err := database.GetBusinessList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting businesses from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "businesses_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "businesses_table.html",
		Data: map[string]any{
			"Businesses":  businesses,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteVendor(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	vendorId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/vendor/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteVendor(vendorId)
	if err != nil {
		fmt.Printf("Error deleting vendor: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete vendor.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	vendors, totalRows, err := database.GetVendorList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting vendors from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "vendors_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "vendors_table.html",
		Data: map[string]any{
			"Vendors":     vendors,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteSupplier(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	supplierId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/supplier/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteSupplier(supplierId)
	if err != nil {
		fmt.Printf("Error deleting supplier: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete supplier.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	suppliers, totalRows, err := database.GetSupplierList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting suppliers from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "suppliers_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "suppliers_table.html",
		Data: map[string]any{
			"Suppliers":   suppliers,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetVendorDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "vendor_detail.html"
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

	vendorId := strings.TrimPrefix(r.URL.Path, "/crm/vendor/")

	vendorDetails, err := database.GetVendorDetails(vendorId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vendor details from DB.", http.StatusInternalServerError)
		return
	}

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Vendor Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Vendor"] = vendorDetails
	data["Cities"] = cities

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetSupplierDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "supplier_detail.html"
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

	supplierId := strings.TrimPrefix(r.URL.Path, "/crm/supplier/")

	supplierDetails, err := database.GetSupplierDetails(supplierId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting supplier details from DB.", http.StatusInternalServerError)
		return
	}

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Supplier Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Supplier"] = supplierDetails
	data["Cities"] = cities

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetBusinessDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "business_detail.html"
	locationsTable := "locations_table.html"
	createLocationForm := "create_location_form.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, constants.PARTIAL_TEMPLATES_DIR + locationsTable, constants.CRM_TEMPLATES_DIR + createLocationForm}
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

	businessId := strings.TrimPrefix(r.URL.Path, "/crm/business/")

	businessDetails, err := database.GetBusinessDetails(businessId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting business details from DB.", http.StatusInternalServerError)
		return
	}

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
	}

	businessLocations, err := database.GetLocationsByBusiness(businessId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting business locations.", http.StatusInternalServerError)
		return
	}

	businesses, err := database.GetBusinesses()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting businesses.", http.StatusInternalServerError)
		return
	}

	vendingLocations, err := database.GetVendingLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	locationStatuses, err := database.GetLocationStatuses()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Business Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Business"] = businessDetails
	data["Cities"] = cities
	data["BusinessLocations"] = businessLocations
	data["Businesses"] = businesses
	data["VendingLocations"] = vendingLocations
	data["LocationStatuses"] = locationStatuses

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func DeleteLocation(w http.ResponseWriter, r *http.Request) {
	businessId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	locationId, err := helpers.GetSecondIDFromPath(r, "/crm/business/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteLocation(locationId)
	if err != nil {
		fmt.Printf("Error deleting location: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete location.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	locations, err := database.GetLocationsByBusiness(fmt.Sprint(businessId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting locations from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "locations_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "locations_table.html",
		Data: map[string]any{
			"Locations": locations,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetLocationDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "location_detail.html"
	machineTables := "machines_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, constants.PARTIAL_TEMPLATES_DIR + machineTables}
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

	businessId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/business/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting business ID from URL.", http.StatusInternalServerError)
		return
	}

	locationId, err := helpers.GetSecondIDFromPath(r, "/crm/business/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting location ID from URL.", http.StatusInternalServerError)
		return
	}

	locationDetails, err := database.GetLocationDetails(businessId, locationId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting location details from DB.", http.StatusInternalServerError)
		return
	}

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
	}

	locationMachines, err := database.GetMachinesByLocation(locationId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting location locations.", http.StatusInternalServerError)
		return
	}

	vendingLocations, err := database.GetVendingLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vending locations.", http.StatusInternalServerError)
		return
	}

	businesses, err := database.GetBusinesses()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting businesses.", http.StatusInternalServerError)
		return
	}

	locationStatuses, err := database.GetLocationStatuses()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting location stauses.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Location Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Location"] = locationDetails
	data["Cities"] = cities
	data["Machines"] = locationMachines
	data["VendingLocations"] = vendingLocations
	data["Businesses"] = businesses
	data["LocationStatuses"] = locationStatuses

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetImagesUpload(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "images_upload.html"
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
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	helpers.ServeContent(w, files, data)
}

func PostImagesUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(50 << 20) // 50 MB limit
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

	values, err := sessions.Get(r)
	if err != nil || values.UserID == 0 {
		fmt.Printf("COULD NOT GET SESSION WHILE UPLOADING IMAGES: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Could not get session while uploading images.",
			},
		}
		w.WriteHeader(http.StatusForbidden)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	files := r.MultipartForm.File["upload_images"]

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			fmt.Printf("%+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to open uploaded file.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
		defer file.Close()

		fileExtension := filepath.Ext(fileHeader.Filename)
		src := uuid.New().String() + fileExtension
		filePath := "marketing-images/" + src

		err = services.UploadFileToS3(file, fileHeader.Size, filePath)
		if err != nil {
			fmt.Printf("Failed to upload image to S3: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to upload image to S3.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		dateAdded, err := utils.GetCurrentTimeInEST()
		if err != nil {
			fmt.Printf("Error getting time as EST: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error getting time as EST.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		form := models.Image{
			Src:           src,
			DateAdded:     dateAdded,
			AddedByUserID: values.UserID,
		}

		err = database.CreateMarketingImage(form)
		if err != nil {
			fmt.Printf("%+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Error saving image metadata.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Images have been successfully uploaded.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetMachineDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "machine_detail.html"
	slotsTable := "slots_table.html"
	createSlotForm := "create_slot_form.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, constants.PARTIAL_TEMPLATES_DIR + slotsTable, constants.CRM_TEMPLATES_DIR + createSlotForm}
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

	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machine ID from URL.", http.StatusInternalServerError)
		return
	}

	machineDetails, err := database.GetMachineDetails(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting location details from DB.", http.StatusInternalServerError)
		return
	}

	cities, err := database.GetCities()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
	}

	vendingTypes, err := database.GetVendingTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
	}

	vendors, err := database.GetVendors()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
	}

	machineStatuses, err := database.GetMachineStatuses()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting cities.", http.StatusInternalServerError)
		return
	}

	locations, err := database.GetLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting locations.", http.StatusInternalServerError)
		return
	}

	slots, err := database.GetMachineSlotsByMachineID(fmt.Sprint(machineId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machine slots.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Machine Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Machine"] = machineDetails
	data["Cities"] = cities
	data["VendingTypes"] = vendingTypes
	data["Vendors"] = vendors
	data["MachineStatuses"] = machineStatuses
	data["Locations"] = locations
	data["Slots"] = slots

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostSlot(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.SlotForm
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

	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slotId, err := database.CreateSlot(form)
	if err != nil {
		fmt.Printf("Error creating location: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create slot.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	price := utils.CreateNullFloat64(form.Price)
	dateAssigned, err := utils.GetCurrentTimeInEST()
	if err != nil {
		fmt.Printf("Error getting time as EST: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting time as EST.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	if price.Valid {
		slotPriceLog := models.SlotPriceLog{
			SlotID:       slotId,
			Price:        price.Float64,
			DateAssigned: dateAssigned,
		}

		err = database.CreateSlotPriceLog(slotPriceLog)
		if err != nil {
			fmt.Printf("Error creating slot price log: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to create slot price log.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	slots, err := database.GetMachineSlotsByMachineID(fmt.Sprint(machineId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting slots from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "slots_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "slots_table.html",
		Data: map[string]any{
			"Slots": slots,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteSlot(w http.ResponseWriter, r *http.Request) {
	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	slotId, err := helpers.GetSecondIDFromPath(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteSlot(slotId)
	if err != nil {
		fmt.Printf("Error deleting slot: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete slot.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	slots, err := database.GetMachineSlotsByMachineID(fmt.Sprint(machineId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting slots from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "slots_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "slots_table.html",
		Data: map[string]any{
			"Slots": slots,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetSlotDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "slot_detail.html"
	productSlotAssignmentTables := "product_slot_assignments_table.html"
	createProductSlotAssignmentForm := "create_product_slot_assignments_form.html"
	productPriceSlotLogsTable := "price_slot_logs_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, constants.PARTIAL_TEMPLATES_DIR + productSlotAssignmentTables, constants.CRM_TEMPLATES_DIR + createProductSlotAssignmentForm, constants.PARTIAL_TEMPLATES_DIR + productPriceSlotLogsTable}
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

	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machine ID from URL.", http.StatusInternalServerError)
		return
	}

	slotId, err := helpers.GetSecondIDFromPath(r, "/crm/machine/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting slot ID from URL.", http.StatusInternalServerError)
		return
	}

	slotDetails, err := database.GetSlotDetails(fmt.Sprint(machineId), fmt.Sprint(slotId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting slot details from DB.", http.StatusInternalServerError)
		return
	}

	products, err := database.GetProducts()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting products from DB.", http.StatusInternalServerError)
		return
	}

	suppliers, err := database.GetSuppliers()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting suppliers from DB.", http.StatusInternalServerError)
		return
	}

	productSlotAssignments, err := database.GetProductSlotAssignments(fmt.Sprint(slotId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting product slot assignments.", http.StatusInternalServerError)
		return
	}

	priceSlotLogs, err := database.GetSlotPriceLogs(fmt.Sprint(slotId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting price change logs.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Slot Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["Slot"] = slotDetails
	data["ProductSlotAssignments"] = productSlotAssignments
	data["Suppliers"] = suppliers
	data["Products"] = products
	data["PriceSlotLogs"] = priceSlotLogs

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PutSlot(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.SlotForm
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

	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slotId, err := helpers.GetSecondIDFromPath(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	slot, err := database.GetSlotDetails(fmt.Sprint(machineId), fmt.Sprint(slotId))
	if err != nil {
		fmt.Printf("Error updating: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get slot.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.UpdateSlot(slotId, form)
	if err != nil {
		fmt.Printf("Error updating: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update slot.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	price := utils.CreateNullFloat64(form.Price)
	dateAssigned, err := utils.GetCurrentTimeInEST()
	if err != nil {
		fmt.Printf("Error getting time as EST: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting time as EST.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}
	if price.Valid && price.Float64 != slot.Price {
		slotPriceLog := models.SlotPriceLog{
			SlotID:       slotId,
			Price:        price.Float64,
			DateAssigned: dateAssigned,
		}

		err = database.CreateSlotPriceLog(slotPriceLog)
		if err != nil {
			fmt.Printf("Error creating slot price log: %+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to create slot price log.",
				},
			}
			w.WriteHeader(http.StatusInternalServerError)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Slot details have been updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostProductSlotAssignment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.ProductSlotAssignmentForm
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

	err = database.CreateProductSlotAssignment(form)
	if err != nil {
		fmt.Printf("Error creating product slot assignment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create product slot assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	slotId := utils.CreateNullInt(form.SlotID)
	if !slotId.Valid {
		fmt.Printf("Invalid slot id: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Invalid slot id.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	productSlotAssignments, err := database.GetProductSlotAssignments(fmt.Sprint(slotId.Int64))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting product slot assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "product_slot_assignments_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "product_slot_assignments_table.html",
		Data: map[string]any{
			"ProductSlotAssignments": productSlotAssignments,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutProductSlotAssignment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.ProductSlotAssignmentForm
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

	err = database.UpdateProductSlotAssignment(form)
	if err != nil {
		fmt.Printf("Error updating: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update product slot assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Product slot assigment details have been updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteProductSlotAssigment(w http.ResponseWriter, r *http.Request) {
	slotId := r.URL.Query().Get("slotId")

	if len(slotId) == 0 {
		http.Error(w, "Missing slotId querystring.", http.StatusBadRequest)
		return
	}

	productSlotAssignmentId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/product-slot-assignment/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteProductSlotAssignment(productSlotAssignmentId)
	if err != nil {
		fmt.Printf("Error deleting product slot assignment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete product slot assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	productSlotAssignments, err := database.GetProductSlotAssignments(slotId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting product slot assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "product_slot_assignments_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "product_slot_assignments_table.html",
		Data: map[string]any{
			"ProductSlotAssignments": productSlotAssignments,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostRefill(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.RefillForm
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

	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateRefill(form)
	if err != nil {
		fmt.Printf("Error creating refill: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create refill.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	slots, err := database.GetMachineSlotsByMachineID(fmt.Sprint(machineId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting slots from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "slots_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "slots_table.html",
		Data: map[string]any{
			"Slots": slots,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteRefill(w http.ResponseWriter, r *http.Request) {
	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	refillId, err := helpers.GetThirdIDFromPath(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteRefill(refillId)
	if err != nil {
		fmt.Printf("Error deleting refill: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete refill.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	slots, err := database.GetMachineSlotsByMachineID(fmt.Sprint(machineId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting slots from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "slots_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "slots_table.html",
		Data: map[string]any{
			"Slots": slots,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostRefillAll(w http.ResponseWriter, r *http.Request) {
	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateRefillAll(machineId)
	if err != nil {
		fmt.Printf("Error creating bulk refill: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create bulk refill.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	slots, err := database.GetMachineSlotsByMachineID(fmt.Sprint(machineId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting slots from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "slots_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "slots_table.html",
		Data: map[string]any{
			"Slots": slots,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetProductSlotAssignmentDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "product_slot_assignment_detail.html"
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

	productSlotAssignmentId := strings.TrimPrefix(r.URL.Path, "/crm/product-slot-assignment/")

	productSlotAssignmentDetails, err := database.GetProductSlotAssignmentDetails(productSlotAssignmentId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vendor details from DB.", http.StatusInternalServerError)
		return
	}

	suppliers, err := database.GetSuppliers()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting suppliers.", http.StatusInternalServerError)
		return
	}

	products, err := database.GetProducts()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting products.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Product Slot Assignment Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["ProductSlotAssignment"] = productSlotAssignmentDetails
	data["Products"] = products
	data["Suppliers"] = suppliers

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func DeleteSlotPriceLog(w http.ResponseWriter, r *http.Request) {
	slotId := r.URL.Query().Get("slotId")

	if len(slotId) == 0 {
		http.Error(w, "Missing slotId querystring.", http.StatusBadRequest)
		return
	}

	priceSlotId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/slot-price-log/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeletePriceSlotLog(fmt.Sprint(priceSlotId))
	if err != nil {
		fmt.Printf("Error deleting slot price log: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete slot price log.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	priceSlotLogs, err := database.GetSlotPriceLogs(slotId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting slot price logs from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "price_slot_logs_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "price_slot_logs_table.html",
		Data: map[string]any{
			"PriceSlotLogs": priceSlotLogs,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutSlotPriceLog(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.SlotPriceLog
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

	err = database.UpdateSlotPriceLog(form)
	if err != nil {
		fmt.Printf("Error updating: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update slot price log.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	slotId := utils.CreateNullInt(form.SlotID)
	if !slotId.Valid {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error parsing slot ID.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	priceSlotLogs, err := database.GetSlotPriceLogs(fmt.Sprint(slotId.Int64))
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting slot price logs from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "price_slot_logs_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "price_slot_logs_table.html",
		Data: map[string]any{
			"PriceSlotLogs": priceSlotLogs,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetEmailSchedules(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "email_schedule.html"
	createEmailScheduleForm := constants.CRM_TEMPLATES_DIR + "create_email_schedule_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "email_schedule_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, baseFile, table, createEmailScheduleForm}

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

	pageNum := 1
	hasPageNum := r.URL.Query().Has("page_num")

	if hasPageNum {
		num, err := strconv.Atoi(r.URL.Query().Get("page_num"))
		if err == nil && num > 1 {
			pageNum = num
		}
	}

	emailSchedules, totalRows, err := database.GetEmailSchedules(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machines from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Machines — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["EmailSchedules"] = emailSchedules
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetEmailScheduleDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "email_schedule_detail.html"
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

	emailScheduleId := strings.TrimPrefix(r.URL.Path, "/crm/email-schedule/")

	emailScheduleDetails, err := database.GetEmailScheduleDetails(emailScheduleId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting vendor details from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Product Slot Assignment Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["EmailSchedule"] = emailScheduleDetails

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PostEmailSchedule(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.EmailScheduleForm
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

	err = database.CreateEmailSchedule(form)
	if err != nil {
		fmt.Printf("Error creating email schedule: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create email schedule.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	emailSchedules, totalRows, err := database.GetEmailSchedules(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting email schedules from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "email_schedule_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "email_schedule_table.html",
		Data: map[string]any{
			"EmailSchedules": emailSchedules,
			"CurrentPage":    pageNum,
			"MaxPages":       helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutEmailSchedule(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Printf("Error parsing form: %+v\n", err)
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

	var form types.EmailScheduleForm
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

	emailScheduleId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/email-schedule/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateEmailSchedule(emailScheduleId, form)
	if err != nil {
		fmt.Printf("Error updating: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update email schedule.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "modal",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "modal.html",
		Data: map[string]any{
			"AlertHeader":  "Success!",
			"AlertMessage": "Email schedule details have been updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteEmailSchedule(w http.ResponseWriter, r *http.Request) {
	emailScheduleId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/email-schedule/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteEmailSchedule(emailScheduleId)
	if err != nil {
		fmt.Printf("Error deleting product slot assignment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete product slot assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	emailSchedules, totalRows, err := database.GetEmailSchedules(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting email schedules from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "email_schedule_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "email_schedule_table.html",
		Data: map[string]any{
			"EmailSchedules": emailSchedules,
			"CurrentPage":    pageNum,
			"MaxPages":       helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
