package handlers

import (
	"net/http"
	"time"

	"github.com/davidalvarez305/yd_vending/constants"
	"github.com/davidalvarez305/yd_vending/helpers"
)

const (
	OptInEventName       string = "opt_in"
	ApplicationEventName string = "lead_application"
	AppointmentEventName string = "lead_appointment"
)

var funnelBaseFilePath = constants.FUNNEL_TEMPLATES_DIR + "base.html"

func createFunnelContext() map[string]any {
	return map[string]any{
		"PageTitle":                    constants.CompanyName,
		"MetaDescription":              "Get a quote for vending machine services.",
		"SiteName":                     constants.SiteName,
		"StaticPath":                   constants.StaticPath,
		"MediaPath":                    constants.MediaPath,
		"PhoneNumber":                  helpers.FormatPhoneNumber(constants.CompanyPhoneNumber),
		"CurrentYear":                  time.Now().Year(),
		"GoogleAnalyticsID":            constants.GoogleAnalyticsID,
		"GoogleAdsID":                  constants.GoogleAdsID,
		"GoogleAdsCallConversionLabel": constants.GoogleAdsCallConversionLabel,
		"FacebookDataSetID":            constants.FacebookDatasetID,
		"CompanyName":                  constants.CompanyName,
	}
}

func FunnelHandler(w http.ResponseWriter, r *http.Request) {
	ctx := createFunnelContext()
	ctx["PagePath"] = constants.RootDomain + r.URL.Path

	externalId, ok := r.Context().Value("external_id").(string)
	if !ok {
		http.Error(w, "Error retrieving external id in context.", http.StatusInternalServerError)
		return
	}

	ctx["ExternalID"] = externalId

	switch r.Method {
	case http.MethodGet:
		switch r.URL.Path {
		case "/funnel/90-day-challenge":
			Get90DayVendingChallenge(w, r, ctx)
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	case http.MethodPost:
		switch r.URL.Path {
		default:
			http.Error(w, "Not Found", http.StatusNotFound)
		}
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}

func Get90DayVendingChallenge(w http.ResponseWriter, r *http.Request, ctx map[string]any) {
	isMobile := helpers.IsMobileRequest(r)
	heroImagePath := "90_day_hero_image_desktop.html"
	if isMobile {
		heroImagePath = "90_day_hero_image_mobile.html"
	}

	fileName := "90_day_challenge.html"
	applicationForm := constants.FUNNEL_TEMPLATES_DIR + "90_day_challenge_application_form.html"
	files := []string{funnelBaseFilePath, constants.FUNNEL_TEMPLATES_DIR + heroImagePath, applicationForm, constants.FUNNEL_TEMPLATES_DIR + fileName}
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
	data["PageTitle"] = "Get 5 Locations in 90 Days Challenge — " + constants.CompanyName
	data["Nonce"] = nonce
	data["CSRFToken"] = csrfToken
	data["LeadTypeID"] = constants.LeadApplicationLeadTypeID

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	helpers.ServeContent(w, files, data)
}
