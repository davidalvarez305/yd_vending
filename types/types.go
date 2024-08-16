package types

import (
	"time"

	"github.com/davidalvarez305/yd_vending/models"
)

type QuoteForm struct {
	FirstName        *string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName         *string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber      *string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Rent             *string `json:"rent" form:"rent" schema:"rent"`
	City             *int    `json:"city" form:"city" schema:"city"`
	LocationType     *int    `json:"location_type" form:"location_type" schema:"location_type"`
	MachineType      *int    `json:"machine_type" form:"machine_type" schema:"machine_type"`
	FootTraffic      *string `json:"foot_traffic" form:"foot_traffic" schema:"foot_traffic"`
	FootTrafficType  *string `json:"foot_traffic_type" form:"foot_traffic_type" schema:"foot_traffic_type"`
	Message          *string `json:"message" form:"message" schema:"message"`
	Source           *string `json:"source" form:"source" schema:"source"`
	Medium           *string `json:"medium" form:"medium" schema:"medium"`
	Channel          *string `json:"channel" form:"channel" schema:"channel"`
	LandingPage      *string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	Keyword          *string `json:"keyword" form:"keyword" schema:"keyword"`
	Referrer         *string `json:"referrer" form:"referrer" schema:"referrer"`
	ClickID          *string `json:"click_id" form:"click_id" schema:"click_id"`
	CampaignID       *int    `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	AdCampaign       *string `json:"ad_campaign" form:"ad_campaign" schema:"ad_campaign"`
	AdGroupID        *int    `json:"ad_group_id" form:"ad_group_id" schema:"ad_group_id"`
	AdGroupName      *string `json:"ad_group_name" form:"ad_group_name" schema:"ad_group_name"`
	AdSetID          *int    `json:"ad_set_id" form:"ad_set_id" schema:"ad_set_id"`
	AdSetName        *string `json:"ad_set_name" form:"ad_set_name" schema:"ad_set_name"`
	AdID             *int    `json:"ad_id" form:"ad_id" schema:"ad_id"`
	AdHeadline       *int    `json:"ad_headline" form:"ad_headline" schema:"ad_headline"`
	Language         *string `json:"language" form:"language" schema:"language"`
	Longitude        *string `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude         *string `json:"latitude" form:"latitude" schema:"latitude"`
	UserAgent        *string `json:"user_agent" form:"user_agent" schema:"user_agent"`
	ButtonClicked    *string `json:"button_clicked" form:"button_clicked" schema:"button_clicked"`
	IP               *string `json:"ip" form:"ip" schema:"ip"`
	CSRFToken        *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	ExternalID       *string `json:"external_id" form:"external_id" schema:"external_id"`
	GoogleClientID   *string `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	FacebookClickID  *string `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID *string `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	CityString       *string `json:"city_string" form:"city_string" schema:"city_string"`
	CSRFSecret       *[]byte `json:"csrf_secret"`
}

type ContactForm struct {
	CSRFToken string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	FirstName string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName  string `json:"last_name" form:"last_name" schema:"last_name"`
	Email     string `json:"email" form:"email" schema:"email"`
	Message   string `json:"message" form:"message" schema:"message"`
}

type OutboundMessageForm struct {
	To        string `json:"to" form:"to" schema:"to"`
	Body      string `json:"body" form:"body" schema:"body"`
	From      string `json:"from" form:"from" schema:"from"`
	CSRFToken string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
}

type LeadDetails struct {
	LeadID          int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FirstName       string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName        string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber     string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	VendingType     string `json:"vending_type" form:"vending_type" schema:"vending_type"`
	VendingLocation string `json:"vending_location" form:"vending_location" schema:"vending_location"`
	CampaignName    string `json:"campaign_name" form:"campaign_name" schema:"campaign_name"`
	Medium          string `json:"medium" form:"medium" schema:"medium"`
	Source          string `json:"source" form:"source" schema:"source"`
	Referrer        string `json:"referrer" form:"referrer" schema:"referrer"`
	LandingPage     string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	IP              string `json:"ip" form:"ip" schema:"ip"`
	Keyword         string `json:"keyword" form:"keyword" schema:"keyword"`
	Channel         string `json:"channel" form:"channel" schema:"channel"`
	Language        string `json:"language" form:"language" schema:"language"`
	City            string `json:"city" form:"city" schema:"city"`
}

type LeadList struct {
	LeadID            int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FirstName         string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName          string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber       string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	CreatedAt         int64  `json:"created_at" form:"created_at" schema:"created_at"`
	Rent              string `json:"rent" form:"rent" schema:"rent"`
	FootTraffic       string `json:"foot_traffic" form:"foot_traffic" schema:"foot_traffic"`
	FootTrafficType   string `json:"foot_traffic_type" form:"foot_traffic_type" schema:"foot_traffic_type"`
	MachineType       string `json:"machine_type" form:"machine_type" schema:"machine_type"`
	LocationType      string `json:"location_type" form:"location_type" schema:"location_type"`
	City              string `json:"city" form:"city" schema:"city"`
	Language          string `json:"language" form:"language" schema:"language"`
	CityID            int    `json:"city_id" form:"city_id" schema:"city_id"`
	VendingTypeID     int    `json:"vending_type_id" form:"vending_type_id" schema:"vending_type_id"`
	VendingLocationID int    `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	TotalRows         int    `json:"total_rows" form:"total_rows" schema:"total_rows"`
}

type GetLeadsParams struct {
	VendingType  *string `json:"vending_type" form:"vending_type" schema:"vending_type"`
	LocationType *string `json:"location_type" form:"location_type" schema:"location_type"`
	City         *string `json:"city" form:"city" schema:"city"`
	PageNum      *string `json:"page_num" form:"page_num" schema:"page_num"`
}

type DynamicPartialTemplate struct {
	TemplateName string
	TemplatePath string
	Data         map[string]any
}

type TwilioMessage struct {
	MessageSid          string    `json:"MessageSid"`
	AccountSid          string    `json:"AccountSid"`
	MessagingServiceSid string    `json:"MessagingServiceSid"`
	From                string    `json:"From"`
	To                  string    `json:"To"`
	Body                string    `json:"Body"`
	NumMedia            string    `json:"NumMedia"`
	NumSegments         string    `json:"NumSegments"`
	SmsStatus           string    `json:"SmsStatus"`
	ApiVersion          string    `json:"ApiVersion"`
	DateCreated         time.Time `json:"DateCreated"`
}

type FrontendMessage struct {
	ClientName  string `json:"client_name"`
	UserName    string `json:"user_name"`
	DateCreated int64  `json:"date_created"`
	Message     string `json:"message"`
	IsInbound   bool   `json:"is_inbound"`
}

type UpdateLeadForm struct {
	Method          *string `json:"_method" form:"_method" schema:"_method"`
	CSRFToken       *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	LeadID          *string `json:"lead_id" form:"lead_id" schema:"lead_id"`
	FirstName       *string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName        *string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber     *string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	City            *int    `json:"city" form:"city" schema:"city"`
	VendingType     *int    `json:"vending_type" form:"vending_type" schema:"vending_type"`
	VendingLocation *int    `json:"vending_location" form:"vending_location" schema:"vending_location"`
}

type UpdateLeadMarketingForm struct {
	Method       *string `json:"_method" form:"_method" schema:"_method"`
	CSRFToken    *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	LeadID       *string `json:"lead_id" form:"lead_id" schema:"lead_id"`
	CampaignName *string `json:"campaign_name" form:"campaign_name" schema:"campaign_name"`
	Medium       *string `json:"medium" form:"medium" schema:"medium"`
	Source       *string `json:"source" form:"source" schema:"source"`
	Referrer     *string `json:"referrer" form:"referrer" schema:"referrer"`
	LandingPage  *string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	IP           *string `json:"ip" form:"ip" schema:"ip"`
	Keyword      *string `json:"keyword" form:"keyword" schema:"keyword"`
	Channel      *string `json:"channel" form:"channel" schema:"channel"`
	Language     *string `json:"language" form:"language" schema:"language"`
}

type TwilioSMSResponse struct {
	Sid                 string            `json:"sid"`
	DateCreated         string            `json:"date_created"`
	DateUpdated         string            `json:"date_updated"`
	DateSent            string            `json:"date_sent"`
	AccountSid          string            `json:"account_sid"`
	To                  string            `json:"to"`
	From                string            `json:"from"`
	MessagingServiceSid string            `json:"messaging_service_sid"`
	Body                string            `json:"body"`
	Status              string            `json:"status"`
	NumSegments         string            `json:"num_segments"`
	NumMedia            string            `json:"num_media"`
	Direction           string            `json:"direction"`
	ApiVersion          string            `json:"api_version"`
	Price               string            `json:"price"`
	PriceUnit           string            `json:"price_unit"`
	ErrorCode           string            `json:"error_code"`
	ErrorMessage        string            `json:"error_message"`
	Uri                 string            `json:"uri"`
	SubresourceUris     map[string]string `json:"subresource_uris"`
}

type TwilioIncomingCallBody struct {
	CallSid       string  `json:"CallSid" form:"CallSid" schema:"CallSid"`
	AccountSid    string  `json:"AccountSid" form:"AccountSid" schema:"AccountSid"`
	From          string  `json:"From" form:"From" schema:"From"`
	To            string  `json:"To" form:"To" schema:"To"`
	CallStatus    string  `json:"CallStatus" form:"CallStatus" schema:"CallStatus"`
	ApiVersion    string  `json:"ApiVersion" form:"ApiVersion" schema:"ApiVersion"`
	Direction     string  `json:"Direction" form:"Direction" schema:"Direction"`
	ForwardedFrom string  `json:"ForwardedFrom" form:"ForwardedFrom" schema:"ForwardedFrom"`
	CallerName    string  `json:"CallerName" form:"CallerName" schema:"CallerName"`
	FromCity      string  `json:"FromCity" form:"FromCity" schema:"FromCity"`
	FromState     string  `json:"FromState" form:"FromState" schema:"FromState"`
	FromZip       string  `json:"FromZip" form:"FromZip" schema:"FromZip"`
	FromCountry   string  `json:"FromCountry" form:"FromCountry" schema:"FromCountry"`
	ToCity        string  `json:"ToCity" form:"ToCity" schema:"ToCity"`
	ToState       string  `json:"ToState" form:"ToState" schema:"ToState"`
	ToZip         string  `json:"ToZip" form:"ToZip" schema:"ToZip"`
	ToCountry     string  `json:"ToCountry" form:"ToCountry" schema:"ToCountry"`
	Caller        string  `json:"Caller" form:"Caller" schema:"Caller"`
	Digits        string  `json:"Digits" form:"Digits" schema:"Digits"`
	SpeechResult  string  `json:"SpeechResult" form:"SpeechResult" schema:"SpeechResult"`
	Confidence    float64 `json:"Confidence" form:"Confidence" schema:"Confidence"`
}

type IncomingPhoneCallForwarding struct {
	FirstName          string `json:"first_name"`
	UserID             int    `json:"user_id"`
	LeadID             int    `json:"lead_id"`
	ForwardPhoneNumber string `json:"forward_phone_number"`
}

type IncomingPhoneCallDialStatus struct {
	DialCallStatus   string `json:"dial_call_status" form:"dial_call_status" schema:"dial_call_status"`
	DialCallSid      string `json:"dial_call_sid" form:"dial_call_sid" schema:"dial_call_sid"`
	DialCallDuration int    `json:"dial_call_duration" form:"dial_call_duration" schema:"dial_call_duration"`
	DialBridged      bool   `json:"dial_bridged" form:"dial_bridged" schema:"dial_bridged"`
	RecordingURL     string `json:"recording_url" form:"recording_url" schema:"recording_url"`
}

type WebsiteContext struct {
	PageTitle         string                   `json:"page_title" form:"page_title"`
	MetaDescription   string                   `json:"meta_description" form:"meta_description"`
	SiteName          string                   `json:"site_name" form:"site_name"`
	StaticPath        string                   `json:"static_path" form:"static_path"`
	MediaPath         string                   `json:"media_path" form:"media_path"`
	PhoneNumber       string                   `json:"phone_number" form:"phone_number"`
	CurrentYear       int                      `json:"current_year" form:"current_year"`
	GoogleAnalyticsID string                   `json:"google_analytics_id" form:"google_analytics_id"`
	FacebookDataSetID string                   `json:"facebook_data_set_id" form:"facebook_data_set_id"`
	CompanyName       string                   `json:"company_name" form:"company_name"`
	PagePath          string                   `json:"page_path" form:"page_path"`
	Nonce             string                   `json:"nonce" form:"nonce"`
	Features          []string                 `json:"features" form:"features"`
	CSRFToken         string                   `json:"csrf_token" form:"csrf_token"`
	VendingTypes      []models.VendingType     `json:"vending_types" form:"vending_types"`
	VendingLocations  []models.VendingLocation `json:"vending_locations" form:"vending_location"`
	Cities            []models.City            `json:"cities" form:"cities"`
	ExternalID        string                   `json:"external_id" form:"external_id"`
}
