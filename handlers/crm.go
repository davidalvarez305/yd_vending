package handlers

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/conversions"
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
		if strings.HasPrefix(path, "/crm/email-schedule/") {
			if len(path) > len("/crm/email-schedule/") && helpers.IsNumeric(path[len("/crm/email-schedule/"):]) && strings.Contains(path, "test") {
				GetEmailScheduleTest(w, r)
				return
			}
			if len(path) > len("/crm/email-schedule/") && helpers.IsNumeric(path[len("/crm/email-schedule/"):]) {
				GetEmailScheduleDetail(w, r, ctx)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/vendor/") {
			if len(path) > len("/crm/vendor/") && helpers.IsNumeric(path[len("/crm/vendor/"):]) {
				GetVendorDetail(w, r, ctx)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/mini-site/") {
			if len(path) > len("/crm/mini-site/") && helpers.IsNumeric(path[len("/crm/mini-site/"):]) {
				GetMiniSiteDetail(w, r, ctx)
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
		case "/crm/email-schedule":
			GetEmailSchedules(w, r, ctx)
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
		case "/crm/mini-site":
			GetMiniSites(w, r, ctx)
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
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/location-assignment") {
			PutLocationAssignment(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/card-reader-assignment") {
			PutCardReaderAssignment(w, r)
			return
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
		if strings.HasPrefix(path, "/crm/email-schedule/") {
			if len(path) > len("/crm/email-schedule/") && helpers.IsNumeric(path[len("/crm/email-schedule/"):]) {
				PutEmailSchedule(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/vendor/") {
			if len(path) > len("/crm/vendor/") && helpers.IsNumeric(path[len("/crm/vendor/"):]) {
				PutVendor(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/mini-site/") {
			if len(path) > len("/crm/mini-site/") && helpers.IsNumeric(path[len("/crm/mini-site/"):]) {
				PutMiniSite(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/mini-site/") {
			if len(path) > len("/crm/mini-site/") && strings.Contains(path, "/vercel") && strings.Contains(path, "/env") {
				PutVercelProjectEnvironmentVariables(w, r)
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
		if strings.HasPrefix(path, "/crm/lead/") && strings.Contains(path, "/offer") {
			PostLeadOffer(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/lead/") && strings.Contains(path, "/appointment") {
			PostLeadAppointment(w, r)
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
		if strings.HasPrefix(path, "/crm/mini-site/") && strings.Contains(path, "/vercel") {
			PostVercelProject(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/mini-site/") && strings.Contains(path, "/deploy") {
			PostVercelDeployProject(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/location-assignment") {
			PostLocationAssignment(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/card-reader-assignment") {
			PostCardReaderAssignment(w, r)
			return
		}

		switch path {
		case "/crm/business":
			PostBusiness(w, r)
		case "/crm/mini-site":
			PostMiniSite(w, r)
		case "/crm/machine":
			PostMachine(w, r)
		case "/crm/email-schedule":
			PostEmailSchedule(w, r)
		case "/crm/vendor":
			PostVendor(w, r)
		case "/crm/supplier":
			PostSupplier(w, r)
		case "/crm/upload-images":
			PostImagesUpload(w, r)
		case "/crm/product-slot-assignment":
			PostProductSlotAssignment(w, r)
		case "/crm/slot-price-log":
			PostPriceSlotLog(w, r)
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
		if strings.HasPrefix(path, "/crm/lead/") {
			if len(path) > len("/crm/lead/") && helpers.IsNumeric(path[len("/crm/lead/"):]) {
				DeleteLead(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/machine/") {
			if len(path) > len("/crm/machine/") && helpers.IsNumeric(path[len("/crm/machine/"):]) {
				DeleteMachine(w, r)
				return
			}
			if len(parts) >= 6 && parts[4] == "slot" && helpers.IsNumeric(parts[3]) && helpers.IsNumeric(parts[5]) && strings.Contains(path, "/refill") {
				DeleteRefill(w, r)
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
		if strings.HasPrefix(path, "/crm/mini-site/") {
			if len(path) > len("/crm/mini-site/") && helpers.IsNumeric(path[len("/crm/mini-site/"):]) {
				DeleteMiniSite(w, r)
				return
			}
		}
		if strings.HasPrefix(path, "/crm/email-schedule/") {
			if len(path) > len("/crm/email-schedule/") && helpers.IsNumeric(path[len("/crm/email-schedule/"):]) {
				DeleteEmailSchedule(w, r)
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
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/location-assignment") {
			DeleteLocationAssignment(w, r)
			return
		}
		if strings.HasPrefix(path, "/crm/machine/") && strings.Contains(path, "/card-reader-assignment") {
			DeleteCardReaderAssignment(w, r)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetLeads(w http.ResponseWriter, r *http.Request, ctx map[string]interface{}) {
	baseFile := constants.CRM_TEMPLATES_DIR + "leads.html"
	leadsTable := constants.PARTIAL_TEMPLATES_DIR + "leads_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, leadsTable, baseFile}

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
	params.LeadTypeID = helpers.SafeStringToIntPointer(r.URL.Query().Get("lead_type"))
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

	leadTypes, err := database.GetLeadTypes()
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting lead types.", http.StatusInternalServerError)
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
	data["LeadTypes"] = leadTypes

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
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, constants.PARTIAL_TEMPLATES_DIR + "messages.html", constants.PARTIAL_TEMPLATES_DIR + "notes.html", constants.PARTIAL_TEMPLATES_DIR + "lead_images.html", constants.CRM_TEMPLATES_DIR + "create_lead_appointment_form.html", constants.CRM_TEMPLATES_DIR + "create_lead_offer_form.html"}
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
	data["LeadAppointmentEventName"] = constants.LeadAppointmentEventName
	data["LeadOfferEventName"] = constants.LeadOfferEventName

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

	err = database.CreateMachine(form)
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
	cardReaderAssignmentsTable := "card_reader_assignments_table.html"
	createCardReaderAssignmentsForm := "create_card_reader_assignment_form.html"
	locationsAssignmentsTable := "location_assignments_table.html"
	createLocationAssignmentsForm := "create_location_assignment_form.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, constants.PARTIAL_TEMPLATES_DIR + slotsTable, constants.CRM_TEMPLATES_DIR + createSlotForm, constants.PARTIAL_TEMPLATES_DIR + cardReaderAssignmentsTable, constants.CRM_TEMPLATES_DIR + createCardReaderAssignmentsForm, constants.PARTIAL_TEMPLATES_DIR + locationsAssignmentsTable, constants.CRM_TEMPLATES_DIR + createLocationAssignmentsForm}
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

	cardReaderAssignments, err := database.GetMachineCardReaderAssignments(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machine card reader assignments.", http.StatusInternalServerError)
		return
	}

	locationAssignments, err := database.GetMachineLocationAssignments(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting machine location assignments.", http.StatusInternalServerError)
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
	data["MachineCardReaderAssignments"] = cardReaderAssignments
	data["MachineLocationAssignments"] = locationAssignments

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

	err = database.CreateSlot(form)
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
	createPriceSlotLogForm := "create_price_slot_log_form.html"
	productPriceSlotLogsTable := "price_slot_logs_table.html"
	files := []string{crmBaseFilePath, crmFooterFilePath, constants.CRM_TEMPLATES_DIR + fileName, constants.PARTIAL_TEMPLATES_DIR + productSlotAssignmentTables, constants.CRM_TEMPLATES_DIR + createProductSlotAssignmentForm, constants.CRM_TEMPLATES_DIR + createPriceSlotLogForm, constants.PARTIAL_TEMPLATES_DIR + productPriceSlotLogsTable}
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

	slotId, err := helpers.GetSecondIDFromPath(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		fmt.Printf("Error deleting refill: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "No machine id in URL.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	refillId, err := helpers.GetThirdIDFromPath(r, "/crm/machine/")
	if err != nil {
		fmt.Printf("Error deleting refill: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "No refill id in URL.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
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

	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	dateRefilled := helpers.GetInt64PointerFromForm(r, "date_refilled")
	if dateRefilled == nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.CreateRefillAll(machineId, helpers.SafeInt64(dateRefilled))
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

func PostPriceSlotLog(w http.ResponseWriter, r *http.Request) {
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

	var form types.SlotPriceLogForm
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

	err = database.CreateSlotPriceLog(form)
	if err != nil {
		fmt.Printf("Error creating: %+v\n", err)
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

	var form types.SlotPriceLogForm
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
	baseFile := constants.CRM_TEMPLATES_DIR + "email_schedules.html"
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
	err := r.ParseMultipartForm(50 << 20)
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

	file, fileHeader, err := r.FormFile("sql_file")
	if err != nil {
		fmt.Printf("Error opening file: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to open sql file.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}
	defer file.Close()

	sqlFileS3Key := constants.SQL_FILES_S3_BUCKET + fileHeader.Filename

	err = services.UploadFileToS3(file, fileHeader.Size, sqlFileS3Key)
	if err != nil {
		fmt.Printf("Error creating email schedule: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to upload SQL File.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
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

	form.SQLFile = &fileHeader.Filename

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

func GetEmailScheduleTest(w http.ResponseWriter, r *http.Request) {
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

	emailScheduleId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/email-schedule/")
	if err != nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Bad request.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	email, err := database.GetEmailScheduleDetails(fmt.Sprint(emailScheduleId))
	if err != nil {
		fmt.Printf("Error updating: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to ge scheduled email.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	now := time.Now()

	recipients := strings.Split(email.Recipients, ", ")
	subject := email.Subject
	sender := email.Sender

	fileName := fmt.Sprintf("%s_%s_%d.xls", email.EmailName, now.Local().Month().String(), now.Local().Year())
	uploadReportS3Key := constants.EMAIL_ATTACHMENTS_S3_BUCKET + fileName
	localFilePath := constants.LOCAL_FILES_DIR + fileName

	sqlFileS3Key := constants.SQL_FILES_S3_BUCKET + email.SQLFile
	sqlFileLocalPath := constants.SQL_FILES_S3_BUCKET + email.SQLFile

	sqlFile, err := services.DownloadFileFromS3(sqlFileS3Key, sqlFileLocalPath)
	if err != nil {
		fmt.Printf("Failed to download SQL file from S3: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to download SQL file from S3.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	sqlQuery, err := os.ReadFile(sqlFile)
	if err != nil {
		fmt.Printf("Failed to read SQL query from file: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to read SQL query from file.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	data, err := database.ExecuteQueryFromSQLFile(string(sqlQuery))
	if err != nil {
		fmt.Printf("Failed to execute SQL query from file: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to execute SQL query from file.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	template, err := helpers.InsertHTMLIntoEmailTemplate(services.EmailTemplateFilePath, "content.html", email.Body, data)
	if err != nil {
		fmt.Printf("Failed to insert HTML into e-mail template: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to insert HTML into e-mail template.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	body := fmt.Sprintf("Content-Type: text/html; charset=UTF-8\r\n%s", template)

	excelFilePath, err := helpers.GenerateExcelFile(data, "data", localFilePath)
	if err != nil {
		fmt.Printf("Failed to generate XLSX file: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to generate XLSX file.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	fileInfo, err := os.Open(excelFilePath)
	if err != nil {
		fmt.Printf("Failed to generate XLSX file: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to generate XLSX file.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}
	defer fileInfo.Close()

	info, err := fileInfo.Stat()
	if err != nil {
		fmt.Printf("Failed to get file info: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get file info.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = services.UploadFileToS3(fileInfo, info.Size(), uploadReportS3Key)
	if err != nil {
		fmt.Printf("Failed to upload xlsx file to S3: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to upload xlsx file to S3.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = services.SendGmailWithAttachment(recipients, subject, sender, body, excelFilePath)
	if err != nil {
		fmt.Printf("Failed to send e-mail with attachment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to send e-mail with attachment.",
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
			"AlertMessage": "Email has been sent successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteCardReaderAssignment(w http.ResponseWriter, r *http.Request) {
	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	cardReaderAssignmentId, err := helpers.GetSecondIDFromPath(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteMachineCardReaderAssignment(cardReaderAssignmentId)
	if err != nil {
		fmt.Printf("Error deleting card reader assignment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete card reader assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	cardReaderAssignments, err := database.GetMachineCardReaderAssignments(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting card reader assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "card_reader_assignments_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "card_reader_assignments_table.html",
		Data: map[string]any{
			"MachineCardReaderAssignments": cardReaderAssignments,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteLocationAssignment(w http.ResponseWriter, r *http.Request) {
	machineId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	locationAssignmentId, err := helpers.GetSecondIDFromPath(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteMachineLocationAssignment(locationAssignmentId)
	if err != nil {
		fmt.Printf("Error deleting location assignment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete location assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	locationAssignments, err := database.GetMachineLocationAssignments(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting location assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	locations, err := database.GetLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting location from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "location_assignments_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "location_assignments_table.html",
		Data: map[string]any{
			"Locations":                  locations,
			"MachineLocationAssignments": locationAssignments,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostCardReaderAssignment(w http.ResponseWriter, r *http.Request) {
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

	var form types.MachineCardReaderAssignmentForm
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

	err = database.CreateMachineCardReaderAssignment(form)
	if err != nil {
		fmt.Printf("Error creating slot price log: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create card reader assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	cardReaderAssignments, err := database.GetMachineCardReaderAssignments(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting card reader assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "card_reader_assignments_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "card_reader_assignments_table.html",
		Data: map[string]any{
			"MachineCardReaderAssignments": cardReaderAssignments,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostLocationAssignment(w http.ResponseWriter, r *http.Request) {
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

	var form types.MachineLocationAssignmentForm
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

	err = database.CreateMachineLocationAssignment(form)
	if err != nil {
		fmt.Printf("Error creating location assignment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create location assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	locationAssignments, err := database.GetMachineLocationAssignments(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting location assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	locations, err := database.GetLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting location from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "location_assignments_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "location_assignments_table.html",
		Data: map[string]any{
			"Locations":                  locations,
			"MachineLocationAssignments": locationAssignments,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutCardReaderAssignment(w http.ResponseWriter, r *http.Request) {
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

	cardReaderId, err := helpers.GetSecondIDFromPath(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var form types.MachineCardReaderAssignmentForm
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

	err = database.UpdateMachineCardReaderAssignment(cardReaderId, form)
	if err != nil {
		fmt.Printf("Error updating card reader assignment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update card reader assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	cardReaderAssignments, err := database.GetMachineCardReaderAssignments(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting card reader assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "card_reader_assignments_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "card_reader_assignments_table.html",
		Data: map[string]any{
			"MachineCardReaderAssignments": cardReaderAssignments,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutLocationAssignment(w http.ResponseWriter, r *http.Request) {
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

	locationAssignmentId, err := helpers.GetSecondIDFromPath(r, "/crm/machine/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var form types.MachineLocationAssignmentForm
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

	err = database.UpdateMachineLocationAssignment(locationAssignmentId, form)
	if err != nil {
		fmt.Printf("Error updating location assignment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update location assignment.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	locationAssignments, err := database.GetMachineLocationAssignments(machineId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting location assignments from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	locations, err := database.GetLocations()
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting location from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "location_assignments_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "location_assignments_table.html",
		Data: map[string]any{
			"Locations":                  locations,
			"MachineLocationAssignments": locationAssignments,
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostLeadAppointment(w http.ResponseWriter, r *http.Request) {
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

	var form types.LeadAppointmentForm

	form.LeadID = helpers.GetIntPointerFromForm(r, "lead_id")
	form.AppointmentTime = helpers.GetInt64PointerFromForm(r, "appointment_time")
	form.Attendee = helpers.GetStringPointerFromForm(r, "attendee")

	if form.LeadID == nil || form.Attendee == nil {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Lead ID or Attendee cannot be nill.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	lead, err := database.GetLeadDetails(fmt.Sprint(*form.LeadID))
	if err != nil {
		fmt.Printf("Error retrieving lead details: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get lead details from DB.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	eventTitle := "YD Vending, LLC - Call With " + lead.FirstName + " " + lead.LastName
	description := "YD Vending demonstration for 90 day vending challenge."
	location := "https://ydvending.com/"

	bookedTime := utils.CreateNullInt64(form.AppointmentTime)
	if !bookedTime.Valid {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Booked time cannot be nil.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	startTime, err := utils.ConvertTimestampToESTDateTime(bookedTime.Int64)
	if err != nil {
		fmt.Printf("Error converting time: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to convert booked time to EST date time.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	endTime := startTime.Add(30 * time.Minute)
	attendees := strings.Split(*form.Attendee, ", ")

	link, err := services.ScheduleGoogleCalendarEvent(eventTitle, description, location, startTime, endTime, attendees)
	if err != nil {
		fmt.Printf("Error creating event: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create google calendar event.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	form.Link = &link
	err = database.CreateLeadAppointment(form)
	if err != nil {
		fmt.Printf("Error creating appointment: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating appointment.",
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
			"AlertMessage": "Appointment has been booked.",
		},
	}

	fbEvent := types.FacebookEventData{
		EventName:      constants.LeadAppointmentEventName,
		EventTime:      time.Now().UTC().Unix(),
		ActionSource:   "website",
		EventSourceURL: lead.LandingPage,
		UserData: types.FacebookUserData{
			Email:           helpers.HashString(lead.Email),
			FirstName:       helpers.HashString(lead.FirstName),
			LastName:        helpers.HashString(lead.LastName),
			Phone:           helpers.HashString(lead.PhoneNumber),
			FBC:             lead.FacebookClickID,
			FBP:             lead.FacebookClientID,
			ExternalID:      helpers.HashString(lead.ExternalID),
			ClientIPAddress: lead.IP,
			ClientUserAgent: lead.UserAgent,
		},
	}

	metaPayload := types.FacebookPayload{
		Data: []types.FacebookEventData{fbEvent},
	}

	payload := types.GooglePayload{
		ClientID: lead.GoogleClientID,
		UserId:   lead.ExternalID,
		Events: []types.GoogleEventLead{
			{
				Name: constants.LeadAppointmentEventName,
				Params: types.GoogleEventParamsLead{
					GCLID: lead.ClickID,
				},
			},
		},
		UserData: types.GoogleUserData{
			Sha256EmailAddress: []string{helpers.HashString(lead.Email)},
			Sha256PhoneNumber:  []string{helpers.HashString(lead.PhoneNumber)},

			Address: []types.GoogleUserAddress{
				{
					Sha256FirstName: helpers.HashString(lead.FirstName),
					Sha256LastName:  helpers.HashString(lead.LastName),
				},
			},
		},
	}

	go conversions.SendGoogleConversion(payload)
	go conversions.SendFacebookConversion(metaPayload)

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutLeadApplication(w http.ResponseWriter, r *http.Request) {
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

	var form types.UpdateLeadApplicationForm
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

	err = database.UpdateLeadApplication(form)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error updating lead application.",
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
			"AlertMessage": "Lead application has been successfully updated.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostMiniSite(w http.ResponseWriter, r *http.Request) {
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

	var form types.MiniSiteForm
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

	err = database.CreateMiniSite(form)
	if err != nil {
		fmt.Printf("Error creating mini site: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create mini site.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1 // Always default to one after new entity is created
	miniSites, totalRows, err := database.GetMiniSiteList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting mini sites from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "mini_sites_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "mini_sites_table.html",
		Data: map[string]any{
			"MiniSites":   miniSites,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetMiniSites(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.CRM_TEMPLATES_DIR + "mini_sites.html"
	createVendorForm := constants.CRM_TEMPLATES_DIR + "create_mini_site_form.html"
	table := constants.PARTIAL_TEMPLATES_DIR + "mini_sites_table.html"
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

	miniSites, totalRows, err := database.GetMiniSiteList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting mini sites from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Mini Sites — " + constants.CompanyName

	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["MiniSites"] = miniSites
	data["MaxPages"] = helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage)
	data["CurrentPage"] = pageNum

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func DeleteMiniSite(w http.ResponseWriter, r *http.Request) {
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

	miniSiteId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/mini-site/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteMiniSite(miniSiteId)
	if err != nil {
		fmt.Printf("Error deleting mini site: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete mini site.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := 1
	miniSites, totalRows, err := database.GetMiniSiteList(pageNum)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting mini sites from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "mini_sites_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "mini_sites_table.html",
		Data: map[string]any{
			"MiniSites":   miniSites,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func GetMiniSiteDetail(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	fileName := "mini_site_detail.html"
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

	miniSiteId := strings.TrimPrefix(r.URL.Path, "/crm/mini-site/")

	miniSiteDetails, err := database.GetMiniSiteDetails(miniSiteId)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting mini site details from DB.", http.StatusInternalServerError)
		return
	}

	data := ctx
	data["PageTitle"] = "Mini Site Detail — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["MiniSite"] = miniSiteDetails
	data["VercelUsername"] = constants.VercelUsername

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func PutMiniSite(w http.ResponseWriter, r *http.Request) {
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

	var form types.MiniSiteForm
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

	miniSiteId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/mini-site/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.UpdateMiniSite(miniSiteId, form)
	if err != nil {
		fmt.Printf("Error updating mini site: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update mini site.",
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
			"AlertMessage": "Mini site updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostVercelProject(w http.ResponseWriter, r *http.Request) {
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

	var form types.VercelProjectForm
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

	miniSiteId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/mini-site")
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Couldn't get mini site ID from request.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	slug := helpers.SafeString(form.Slug)
	projectName := helpers.SafeString(form.ProjectName)

	envVars := []types.EnvironmentVariable{}
	v := reflect.ValueOf(form)

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		value := v.Field(i)

		envTag := field.Tag.Get(constants.MiniSiteEnvironmentVariablesTag)
		if !strings.Contains(envTag, constants.MiniSiteEnvironmentVariablesPrefix) || value.IsNil() {
			continue
		}

		var val string
		if ptr, ok := value.Interface().(*string); ok && ptr != nil {
			val = *ptr
		}

		envVars = append(envVars, types.EnvironmentVariable{
			Key:    envTag,
			Target: "production",
			Type:   "plain",
			Value:  val,
		})
	}

	project := types.VercelProjectRequestBody{
		Name:                              projectName,
		BuildCommand:                      constants.MiniSiteBuildCommand,
		DevCommand:                        constants.MiniSiteDevCommand,
		EnableAffectedProjectsDeployments: true,
		EnvironmentVariables:              envVars,
		Framework:                         constants.MiniSiteFramework,
		GitRepository: types.GitRepository{
			Repo: constants.MiniSiteGithubRepo,
			Type: "github",
		},
		OIDCTokenConfig: types.OIDCTokenConfig{
			Enabled:    true,
			IssuerMode: "team",
		},
		OutputDirectory: constants.MiniSiteOutputDirectory,
		PublicSource:    true,
	}

	resp, err := services.CreateVercelProject(slug, constants.MiniSiteVercelTeamID, constants.VercelAccessToken, project)
	if err != nil {
		fmt.Printf("Error creating vercel project: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create vercel project.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.CreateMiniSiteEnvironmentVariables(miniSiteId, resp.Env)
	if err != nil {
		fmt.Printf("Error creating vercel project environment variables: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create vercel project environment variables.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.UpdateMiniSiteProjectID(miniSiteId, resp.ID)
	if err != nil {
		fmt.Printf("Error updating vercel project: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update vercel project.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	git, err := services.GetLatestGithubCommit(constants.MiniSiteGithubRepo, constants.MiniSiteBranchName)
	if err != nil {
		fmt.Printf("Error getting latest git commit: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get latest git commit.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	deploymentBody := types.DeployVercelProjectBody{
		Name:   helpers.SafeString(form.ProjectName),
		Target: constants.VercelProjectEnvinronmentVariableTarget,
		GitSource: types.GitSource{
			Type:   "github",
			RepoId: constants.MiniSiteRepoID,
			Ref:    constants.MiniSiteBranchName,
			Sha:    git.GetSHA(),
		},
	}

	err = services.CreateVercelProjectDeployment(slug, constants.MiniSiteVercelTeamID, constants.VercelAccessToken, deploymentBody)
	if err != nil {
		fmt.Printf("Error deploying project: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to deploy project.",
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
			"AlertMessage": "Vercel project launched successfully.",
		},
	}
	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PutVercelProjectEnvironmentVariables(w http.ResponseWriter, r *http.Request) {
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

	var form types.VercelProjectForm
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

	miniSiteId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/mini-site/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	miniSite, err := database.GetMiniSiteDetails(fmt.Sprint(miniSiteId))
	if err != nil {
		fmt.Printf("Error querying mini site: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get mini site from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	slug := helpers.SafeString(form.Slug)

	environmentVariables, err := database.GetMiniSiteEnvironmentVariablesByProject(miniSiteId)
	if err != nil {
		fmt.Printf("Error querying mini site environment variables: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get mini site environment variables from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var body []types.UpdateVercelEnvironmentVariablesBody

	formValue := reflect.ValueOf(form)
	formType := formValue.Type()

	for i := 0; i < formValue.NumField(); i++ {
		field := formValue.Field(i)
		fieldType := formType.Field(i)

		envTag := fieldType.Tag.Get(constants.MiniSiteEnvironmentVariablesTag)

		if envTag == "" || !strings.Contains(envTag, constants.MiniSiteEnvironmentVariablesPrefix) {
			continue
		}

		for _, environmentVariable := range environmentVariables {
			if environmentVariable.Key == envTag {
				body = append(body, types.UpdateVercelEnvironmentVariablesBody{
					Key:    environmentVariable.Key,
					Value:  helpers.GetStringValueFromField(field),
					ID:     environmentVariable.EnvironmentVariableUniqueID,
					Target: []string{constants.VercelProjectEnvinronmentVariableTarget},
					Type:   constants.VercelProjectEnvinronmentVariableType,
				})
			}
		}
	}

	err = database.UpdateMiniSiteEnvironmentVariables(miniSiteId, body)
	if err != nil {
		fmt.Printf("Error updating mini site: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update mini site.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = services.UpdateVercelEnvironmentVariables(constants.VercelAccessToken, miniSite.VercelProjectID, constants.MiniSiteBranchName, slug, constants.MiniSiteVercelTeamID, body)
	if err != nil {
		fmt.Printf("Error updating environment variables: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to update environment variables.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	git, err := services.GetLatestGithubCommit(constants.MiniSiteGithubRepo, constants.MiniSiteBranchName)
	if err != nil {
		fmt.Printf("Error getting latest git commit: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get latest git commit.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	deploymentBody := types.DeployVercelProjectBody{
		Name:   helpers.SafeString(form.ProjectName),
		Target: constants.VercelProjectEnvinronmentVariableTarget,
		GitSource: types.GitSource{
			Type:   "github",
			RepoId: constants.MiniSiteRepoID,
			Ref:    constants.MiniSiteBranchName,
			Sha:    git.GetSHA(),
		},
	}

	err = services.CreateVercelProjectDeployment(slug, constants.MiniSiteVercelTeamID, constants.VercelAccessToken, deploymentBody)
	if err != nil {
		fmt.Printf("Error deploying project: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to deploy project.",
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
			"AlertMessage": "Mini site updated successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func DeleteLead(w http.ResponseWriter, r *http.Request) {
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

	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = database.DeleteLead(leadId)
	if err != nil {
		fmt.Printf("Error deleting lead: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to delete lead.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var params types.GetLeadsParams
	params.LocationType = helpers.SafeStringToPointer(r.URL.Query().Get("location_type"))
	params.VendingType = helpers.SafeStringToPointer(r.URL.Query().Get("vending_type"))
	params.LeadTypeID = helpers.SafeStringToIntPointer(r.URL.Query().Get("lead_type"))
	params.PageNum = helpers.SafeStringToPointer(r.URL.Query().Get("page_num"))

	leads, totalRows, err := database.GetLeadList(params)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting leads from DB.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	pageNum := "1"
	safePageNum := helpers.SafeString(params.PageNum)
	if safePageNum != "" {
		pageNum = safePageNum
	}

	tmplCtx := types.DynamicPartialTemplate{
		TemplateName: "leads_table.html",
		TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "leads_table.html",
		Data: map[string]any{
			"Leads":       leads,
			"CurrentPage": pageNum,
			"MaxPages":    helpers.CalculateMaxPages(totalRows, constants.LeadsPerPage),
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostVercelDeployProject(w http.ResponseWriter, r *http.Request) {
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

	slug := r.URL.Query().Get("slug")
	projectName := r.URL.Query().Get("projectName")

	if slug == "" || projectName == "" {
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Missing slug from query.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	git, err := services.GetLatestGithubCommit(constants.MiniSiteGithubRepo, constants.MiniSiteBranchName)
	if err != nil {
		fmt.Printf("Error getting latest git commit: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get latest git commit.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	deploymentBody := types.DeployVercelProjectBody{
		Name:   projectName,
		Target: constants.VercelProjectEnvinronmentVariableTarget,
		GitSource: types.GitSource{
			Type:   "github",
			RepoId: constants.MiniSiteRepoID,
			Ref:    constants.MiniSiteBranchName,
			Sha:    git.GetSHA(),
		},
	}

	err = services.CreateVercelProjectDeployment(slug, constants.MiniSiteVercelTeamID, constants.VercelAccessToken, deploymentBody)
	if err != nil {
		fmt.Printf("Error deploying project: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to deploy project.",
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
			"AlertMessage": "Mini site deployed successfully.",
		},
	}

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}

func PostLeadOffer(w http.ResponseWriter, r *http.Request) {
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

	var form types.CreateLeadOfferForm
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

	leadId, err := helpers.GetFirstIDAfterPrefix(r, "/crm/lead/")
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get lead ID from request.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	lead, err := database.GetLeadDetails(fmt.Sprint(leadId))
	if err != nil {
		fmt.Printf("Error retrieving lead details: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to get lead details from DB.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var quantity int
	if form.Quantity != nil {
		quantity = *form.Quantity
	}

	leadQs := "?lead=" + fmt.Sprint(lead.LeadID)
	successUrl := constants.RootDomain + constants.LeadOfferAcceptedSuccessPath + leadQs
	cancelUrl := constants.RootDomain + constants.LeadOfferCanceledPath + leadQs
	clientReferenceId := fmt.Sprint(lead.LeadID)

	link, err := services.CreateStripeCheckout(helpers.SafeString(form.Price), int64(quantity), successUrl, cancelUrl, clientReferenceId)
	if err != nil {
		fmt.Printf("Error creating stripe checkout: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to create stripe checkout.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	recipients := []string{lead.Email}
	subject := constants.CompanyName + " - 90 Day Challenge Offer"
	sender := constants.CompanyEmail
	body := link

	err = services.SendGmail(recipients, subject, sender, body)
	if err != nil {
		fmt.Printf("Error offer e-mail: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while offer e-mail.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	leadOffer := models.LeadOffer{
		LeadID:            leadId,
		Offer:             link,
		LeadOfferStatusID: constants.LeadOfferSentID,
	}

	leadOfferId, err := database.CreateLeadOffer(leadOffer)
	if err != nil {
		fmt.Printf("Error creating lead offer: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating lead offer.",
			},
		}

		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	err = database.CreateLeadOfferStatusLog(leadOfferId, constants.LeadOfferSentID)
	if err != nil {
		fmt.Printf("Error creating lead offer log: %+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Server error while creating lead offer log.",
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
			"AlertMessage": "Appointment has been booked.",
		},
	}

	fbEvent := types.FacebookEventData{
		EventName:      constants.LeadOfferEventName,
		EventTime:      time.Now().UTC().Unix(),
		ActionSource:   "website",
		EventSourceURL: lead.LandingPage,
		UserData: types.FacebookUserData{
			Email:           helpers.HashString(lead.Email),
			FirstName:       helpers.HashString(lead.FirstName),
			LastName:        helpers.HashString(lead.LastName),
			Phone:           helpers.HashString(lead.PhoneNumber),
			FBC:             lead.FacebookClickID,
			FBP:             lead.FacebookClientID,
			ExternalID:      helpers.HashString(lead.ExternalID),
			ClientIPAddress: lead.IP,
			ClientUserAgent: lead.UserAgent,
		},
	}

	metaPayload := types.FacebookPayload{
		Data: []types.FacebookEventData{fbEvent},
	}

	payload := types.GooglePayload{
		ClientID: lead.GoogleClientID,
		UserId:   lead.ExternalID,
		Events: []types.GoogleEventLead{
			{
				Name: constants.LeadOfferEventName,
				Params: types.GoogleEventParamsLead{
					GCLID: lead.ClickID,
				},
			},
		},
		UserData: types.GoogleUserData{
			Sha256EmailAddress: []string{helpers.HashString(lead.Email)},
			Sha256PhoneNumber:  []string{helpers.HashString(lead.PhoneNumber)},

			Address: []types.GoogleUserAddress{
				{
					Sha256FirstName: helpers.HashString(lead.FirstName),
					Sha256LastName:  helpers.HashString(lead.LastName),
				},
			},
		},
	}

	go conversions.SendGoogleConversion(payload)
	go conversions.SendFacebookConversion(metaPayload)

	helpers.ServeDynamicPartialTemplate(w, tmplCtx)
}
