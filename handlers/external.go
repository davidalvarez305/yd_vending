package handlers

import (
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/types"
	"github.com/davidalvarez305/yd_vending/utils"

	"github.com/xuri/excelize/v2"
)

var externalReportsBasePath = constants.EXTERNAL_REPORTS_TEMPLATES_DIR + "base.html"

func createExternalReportsHandler() map[string]any {
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

func ExternalPagesHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	ctx := createExternalReportsHandler()
	ctx["PagePath"] = constants.RootDomain + path

	switch r.Method {
	case http.MethodGet:

		if strings.Contains(r.URL.Path, "/external/commission-report") && strings.Contains(r.URL.Path, "download") {
			GetExternalReportDownload(w, r, ctx)
			return
		}

		if strings.Contains(r.URL.Path, "/external/commission-report") {
			GetExternalReportHandler(w, r, ctx)
			return
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func GetExternalReportHandler(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	baseFile := constants.EXTERNAL_REPORTS_TEMPLATES_DIR + "commission_report.html"
	files := []string{externalReportsBasePath, baseFile}

	nonce, ok := r.Context().Value("nonce").(string)
	if !ok {
		http.Error(w, "Error retrieving nonce.", http.StatusInternalServerError)
		return
	}

	monthYear := r.URL.Query().Get("monthYear")

	if !r.URL.Query().Has("monthYear") || len(monthYear) == 0 {
		currentTime := time.Now()
		monthYear = currentTime.Format("January, 2006")
	}

	start, end, err := utils.GetStartAndEndDatesFromMonthYear(monthYear)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting start and end dates for commission report.", http.StatusInternalServerError)
		return
	}

	var business string
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) > 3 {
		decodedBusinessName, err := url.PathUnescape(parts[3])
		if err != nil {
			http.Error(w, "Failed to decode business name.", http.StatusInternalServerError)
			return
		}

		business = decodedBusinessName
	} else {
		http.Error(w, "Business name not found in URL.", http.StatusInternalServerError)
		return
	}

	businessId, err := database.GetBusinessIDFromURL(business)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting business id from URL.", http.StatusInternalServerError)
		return
	}

	report, err := database.GetCommissionReport(fmt.Sprint(businessId), start, end)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting commission report.", http.StatusInternalServerError)
		return
	}

	dates, err := database.GetAvailableReportDatesByBusiness(fmt.Sprint(businessId))
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting available dates for report.", http.StatusInternalServerError)
		return
	}

	var revenue, costs, grossProfit, commissionDue float64

	for _, line := range report {
		revenue += line.Revenue
		costs += line.Cost + line.CreditCardFee
		grossProfit += line.GrossProfit
		commissionDue += line.CommissionDue
	}

	// Round to 2 decimal points
	revenue = math.Round(revenue*100) / 100
	costs = math.Round(costs*100) / 100
	grossProfit = math.Round(grossProfit*100) / 100
	commissionDue = math.Round(commissionDue*100) / 100

	data := ctx
	data["PageTitle"] = fmt.Sprintf("%s Commission Report — %s", business, constants.CompanyName)
	data["Nonce"] = nonce
	data["CommissionReport"] = report
	data["Revenue"] = revenue
	data["Costs"] = costs
	data["CommissionDue"] = commissionDue
	data["GrossProfit"] = grossProfit
	data["Dates"] = dates
	data["BusinessName"] = business

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}

func GetExternalReportDownload(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	monthYear := r.URL.Query().Get("monthYear")

	if !r.URL.Query().Has("monthYear") || len(monthYear) == 0 {
		currentTime := time.Now()
		monthYear = currentTime.Format("January, 2006")
	}

	start, end, err := utils.GetStartAndEndDatesFromMonthYear(monthYear)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error parsing dates for commission report.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	var business string
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) > 3 {
		decodedBusinessName, err := url.PathUnescape(parts[3])
		if err != nil {
			fmt.Printf("%+v\n", err)
			tmplCtx := types.DynamicPartialTemplate{
				TemplateName: "error",
				TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
				Data: map[string]any{
					"Message": "Failed to decode business name.",
				},
			}
			w.WriteHeader(http.StatusBadRequest)
			helpers.ServeDynamicPartialTemplate(w, tmplCtx)
			return
		}

		business = decodedBusinessName
	} else {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Business not found in URL.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	businessId, err := database.GetBusinessIDFromURL(business)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Failed to decode business id from URL.",
			},
		}
		w.WriteHeader(http.StatusBadRequest)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	report, err := database.GetCommissionReport(fmt.Sprint(businessId), start, end)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error getting commission report.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}

	// Create a new Excel file
	f := excelize.NewFile()
	sheetName := "Commission Report"
	index := f.NewSheet(sheetName)

	headers := []string{"Product", "Amount Sold", "Revenue", "Cost", "Credit Card Fee", "Gross Profit", "Commission Due"}
	for i, header := range headers {
		cell := fmt.Sprintf("%s%d", string('A'+i), 1) // A1, B1, C1, etc.
		f.SetCellValue(sheetName, cell, header)
	}

	for i, row := range report {
		f.SetCellValue(sheetName, fmt.Sprintf("A%d", i+2), row.Product)
		f.SetCellValue(sheetName, fmt.Sprintf("B%d", i+2), row.AmountSold)
		f.SetCellValue(sheetName, fmt.Sprintf("C%d", i+2), row.Revenue)
		f.SetCellValue(sheetName, fmt.Sprintf("D%d", i+2), row.Cost)
		f.SetCellValue(sheetName, fmt.Sprintf("E%d", i+2), row.CreditCardFee)
		f.SetCellValue(sheetName, fmt.Sprintf("F%d", i+2), row.GrossProfit)
		f.SetCellValue(sheetName, fmt.Sprintf("G%d", i+2), row.CommissionDue)
	}

	f.SetActiveSheet(index)

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	w.Header().Set("Content-Disposition", "attachment; filename=commission_report.xlsx")

	err = f.Write(w)
	if err != nil {
		fmt.Printf("%+v\n", err)
		tmplCtx := types.DynamicPartialTemplate{
			TemplateName: "error",
			TemplatePath: constants.PARTIAL_TEMPLATES_DIR + "error_banner.html",
			Data: map[string]any{
				"Message": "Error writing Excel file.",
			},
		}
		w.WriteHeader(http.StatusInternalServerError)
		helpers.ServeDynamicPartialTemplate(w, tmplCtx)
		return
	}
}
