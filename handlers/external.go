package handlers

import (
	"fmt"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/database"
	"github.com/davidalvarez305/yd_vending/helpers"
	"github.com/davidalvarez305/yd_vending/utils"
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

	var location string
	parts := strings.Split(r.URL.Path, "/")

	if len(parts) > 2 {
		location = parts[2]
	} else {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Location ID not found in URL.", http.StatusInternalServerError)
		return
	}

	locationId, err := database.GetLocationIDFromURL(location)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting location id from URL.", http.StatusInternalServerError)
		return
	}

	report, err := database.GetCommissionReport(locationId, start, end)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error getting commission report.", http.StatusInternalServerError)
		return
	}

	dates, err := database.GetAvailableReportDates(locationId)
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
	data["PageTitle"] = "Commission Report — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CommissionReport"] = report
	data["Revenue"] = revenue
	data["Costs"] = costs
	data["CommissionDue"] = commissionDue
	data["GrossProfit"] = grossProfit
	data["Dates"] = dates

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}
