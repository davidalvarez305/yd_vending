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
	CampaignID       *int64  `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	AdCampaign       *string `json:"ad_campaign" form:"ad_campaign" schema:"ad_campaign"`
	AdGroupID        *int64  `json:"ad_group_id" form:"ad_group_id" schema:"ad_group_id"`
	AdGroupName      *string `json:"ad_group_name" form:"ad_group_name" schema:"ad_group_name"`
	AdSetID          *int64  `json:"ad_set_id" form:"ad_set_id" schema:"ad_set_id"`
	AdSetName        *string `json:"ad_set_name" form:"ad_set_name" schema:"ad_set_name"`
	AdID             *int64  `json:"ad_id" form:"ad_id" schema:"ad_id"`
	AdHeadline       *int64  `json:"ad_headline" form:"ad_headline" schema:"ad_headline"`
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
	CSRFSecret       *string `json:"csrf_secret" form:"csrf_secret"`
	Email            *string `json:"email" form:"email" schema:"email"`
	LeadTypeID       *int    `json:"lead_type_id" form:"lead_type_id" schema:"lead_type_id"`
}

type LeadFormWebhook struct {
	LeadID         string           `json:"lead_id" form:"lead_id" schema:"lead_id"`
	APIVersion     string           `json:"api_version" form:"api_version" schema:"api_version"`
	FormID         int64            `json:"form_id" form:"form_id" schema:"form_id"`
	CampaignID     int64            `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	AdGroupID      *int64           `json:"adgroup_id" form:"adgroup_id" schema:"adgroup_id"`
	CreativeID     *int64           `json:"creative_id" form:"creative_id" schema:"creative_id"`
	GCLID          string           `json:"gcl_id" form:"gcl_id" schema:"gcl_id"`
	GoogleKey      string           `json:"google_key" form:"google_key" schema:"google_key"`
	IsTest         *bool            `json:"is_test" form:"is_test" schema:"is_test"`
	UserColumnData []UserColumnData `json:"user_column_data" form:"user_column_data" schema:"user_column_data"`
}

type UserColumnData struct {
	ColumnID    string `json:"column_id" form:"column_id" schema:"column_id"`
	StringValue string `json:"string_value" form:"string_value" schema:"string_value"`
	ColumnName  string `json:"column_name" form:"column_name" schema:"column_name"`
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
	LeadID           int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	LeadTypeID       int    `json:"lead_type_id" form:"lead_type_id" schema:"lead_type_id"`
	FirstName        string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName         string `json:"last_name" form:"last_name" schema:"last_name"`
	Email            string `json:"email" form:"email" schema:"email"`
	PhoneNumber      string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	VendingType      string `json:"vending_type" form:"vending_type" schema:"vending_type"`
	VendingLocation  string `json:"vending_location" form:"vending_location" schema:"vending_location"`
	CampaignName     string `json:"campaign_name" form:"campaign_name" schema:"campaign_name"`
	Medium           string `json:"medium" form:"medium" schema:"medium"`
	Source           string `json:"source" form:"source" schema:"source"`
	Referrer         string `json:"referrer" form:"referrer" schema:"referrer"`
	LandingPage      string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	IP               string `json:"ip" form:"ip" schema:"ip"`
	Keyword          string `json:"keyword" form:"keyword" schema:"keyword"`
	Channel          string `json:"channel" form:"channel" schema:"channel"`
	Language         string `json:"language" form:"language" schema:"language"`
	Message          string `json:"message" form:"message" schema:"message"`
	FacebookClickID  string `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID string `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	UserAgent        string `json:"user_agent" form:"user_agent" schema:"facebook_click_id"`
	ExternalID       string `json:"external_id" form:"external_id" schema:"external_id"`
	ClickID          string `json:"click_id" form:"click_id" schema:"click_id"`
	GoogleClientID   string `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
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
	Language          string `json:"language" form:"language" schema:"language"`
	VendingTypeID     int    `json:"vending_type_id" form:"vending_type_id" schema:"vending_type_id"`
	VendingLocationID int    `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	TotalRows         int    `json:"total_rows" form:"total_rows" schema:"total_rows"`
}

type GetLeadsParams struct {
	VendingType  *string `json:"vending_type" form:"vending_type" schema:"vending_type"`
	LocationType *string `json:"location_type" form:"location_type" schema:"location_type"`
	PageNum      *string `json:"page_num" form:"page_num" schema:"page_num"`
	LeadTypeID   *int    `json:"lead_type_id" form:"lead_type_id" schema:"lead_type_id"`
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
	DateCreated string `json:"date_created"`
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
	Called           string `json:"called" form:"called" schema:"called"`
	ToState          string `json:"to_state" form:"to_state" schema:"to_state"`
	DialCallStatus   string `json:"dial_call_status" form:"dial_call_status" schema:"dial_call_status"`
	CallerCountry    string `json:"caller_country" form:"caller_country" schema:"caller_country"`
	Direction        string `json:"direction" form:"direction" schema:"direction"`
	CallerState      string `json:"caller_state" form:"caller_state" schema:"caller_state"`
	ToZip            string `json:"to_zip" form:"to_zip" schema:"to_zip"`
	DialCallSid      string `json:"dial_call_sid" form:"dial_call_sid" schema:"dial_call_sid"`
	CallSid          string `json:"call_sid" form:"call_sid" schema:"call_sid"`
	To               string `json:"to" form:"to" schema:"to"`
	CallerZip        string `json:"caller_zip" form:"caller_zip" schema:"caller_zip"`
	ToCountry        string `json:"to_country" form:"to_country" schema:"to_country"`
	CalledZip        string `json:"called_zip" form:"called_zip" schema:"called_zip"`
	ApiVersion       string `json:"api_version" form:"api_version" schema:"api_version"`
	CalledCity       string `json:"called_city" form:"called_city" schema:"called_city"`
	CallStatus       string `json:"call_status" form:"call_status" schema:"call_status"`
	From             string `json:"from" form:"from" schema:"from"`
	DialBridged      bool   `json:"dial_bridged" form:"dial_bridged" schema:"dial_bridged"`
	AccountSid       string `json:"account_sid" form:"account_sid" schema:"account_sid"`
	DialCallDuration int    `json:"dial_call_duration" form:"dial_call_duration" schema:"dial_call_duration"`
	CalledCountry    string `json:"called_country" form:"called_country" schema:"called_country"`
	CallerCity       string `json:"caller_city" form:"caller_city" schema:"caller_city"`
	ToCity           string `json:"to_city" form:"to_city" schema:"to_city"`
	FromCountry      string `json:"from_country" form:"from_country" schema:"from_country"`
	Caller           string `json:"caller" form:"caller" schema:"caller"`
	FromCity         string `json:"from_city" form:"from_city" schema:"from_city"`
	CalledState      string `json:"called_state" form:"called_state" schema:"called_state"`
	FromZip          string `json:"from_zip" form:"from_zip" schema:"from_zip"`
	FromState        string `json:"from_state" form:"from_state" schema:"from_state"`
	RecordingURL     string `json:"recording_url" form:"recording_url" schema:"recording_url"`
}

type WebsiteContext struct {
	PageTitle                    string                   `json:"page_title" form:"page_title"`
	MetaDescription              string                   `json:"meta_description" form:"meta_description"`
	SiteName                     string                   `json:"site_name" form:"site_name"`
	StaticPath                   string                   `json:"static_path" form:"static_path"`
	MediaPath                    string                   `json:"media_path" form:"media_path"`
	PhoneNumber                  string                   `json:"phone_number" form:"phone_number"`
	CurrentYear                  int                      `json:"current_year" form:"current_year"`
	LeadTypeID                   int                      `json:"lead_type_id" form:"lead_type_id"`
	GoogleAnalyticsID            string                   `json:"google_analytics_id" form:"google_analytics_id"`
	FacebookDataSetID            string                   `json:"facebook_data_set_id" form:"facebook_data_set_id"`
	CompanyName                  string                   `json:"company_name" form:"company_name"`
	PagePath                     string                   `json:"page_path" form:"page_path"`
	Nonce                        string                   `json:"nonce" form:"nonce"`
	Features                     []string                 `json:"features" form:"features"`
	CSRFToken                    string                   `json:"csrf_token" form:"csrf_token"`
	VendingTypes                 []models.VendingType     `json:"vending_types" form:"vending_types"`
	VendingLocations             []models.VendingLocation `json:"vending_locations" form:"vending_location"`
	ExternalID                   string                   `json:"external_id" form:"external_id"`
	GoogleAdsID                  string                   `json:"google_ads_id"`
	GoogleAdsCallConversionLabel string                   `json:"google_ads_call_conversion_label"`
	MarketingImages              []string                 `json:"marketing_images"`
}

type FacebookUserData struct {
	Phone           string `json:"ph"`
	FirstName       string `json:"fn"`
	LastName        string `json:"ln"`
	Email           string `json:"em"`
	ClientIPAddress string `json:"client_ip_address"`
	ClientUserAgent string `json:"client_user_agent"`
	FBC             string `json:"fbc"`
	FBP             string `json:"fbp"`
	State           string `json:"st"`
	ExternalID      string `json:"external_id"`
}

type FacebookEventData struct {
	EventName      string           `json:"event_name"`
	EventTime      int64            `json:"event_time"`
	ActionSource   string           `json:"action_source"`
	EventSourceURL string           `json:"event_source_url"`
	UserData       FacebookUserData `json:"user_data"`
}

type FacebookPayload struct {
	Data []FacebookEventData `json:"data"`
}

type GoogleEventParamsLead struct {
	GCLID string `json:"gclid" form:"gclid" schema:"gclid"`
}

type GoogleEventLead struct {
	Name   string                `json:"name" form:"name" schema:"name"`
	Params GoogleEventParamsLead `json:"params" form:"params" schema:"params"`
}

type GooglePayload struct {
	ClientID string            `json:"client_id" form:"client_id" schema:"client_id"`
	UserId   string            `json:"userId" form:"userId" schema:"userId"`
	Events   []GoogleEventLead `json:"events" form:"events" schema:"events"`
	UserData GoogleUserData    `json:"user_data" form:"user_data" schema:"user_data"`
}

type GoogleUserData struct {
	Sha256EmailAddress []string            `json:"sha256_email_address,omitempty" form:"sha256_email_address,omitempty" schema:"sha256_email_address,omitempty"`
	Sha256PhoneNumber  []string            `json:"sha256_phone_number,omitempty" form:"sha256_phone_number,omitempty" schema:"sha256_phone_number,omitempty"`
	Address            []GoogleUserAddress `json:"address,omitempty" form:"address,omitempty" schema:"address,omitempty"`
}

type GoogleUserAddress struct {
	Sha256FirstName string `json:"sha256_first_name,omitempty" form:"sha256_first_name,omitempty" schema:"sha256_first_name,omitempty"`
	Sha256LastName  string `json:"sha256_last_name,omitempty" form:"sha256_last_name,omitempty" schema:"sha256_last_name,omitempty"`
	Sha256Street    string `json:"sha256_street,omitempty" form:"sha256_street,omitempty" schema:"sha256_street,omitempty"`
	City            string `json:"city,omitempty" form:"city,omitempty" schema:"city,omitempty"`
	Region          string `json:"region,omitempty" form:"region,omitempty" schema:"region,omitempty"`
	PostalCode      string `json:"postal_code,omitempty" form:"postal_code,omitempty" schema:"postal_code,omitempty"`
	Country         string `json:"country,omitempty" form:"country,omitempty" schema:"country,omitempty"`
}

type ConversionLeadInfo struct {
	LeadID       int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	MachineType  string `json:"machine_type" form:"machine_type" schema:"machine_type"`
	LocationType string `json:"location_type" form:"location_type" schema:"location_type"`
	CreatedAt    int64  `json:"created_at" form:"created_at" schema:"created_at"`
}

type FrontendNote struct {
	UserName  string `json:"user_name"`
	DateAdded string `json:"date_added"`
	Note      string `json:"note"`
}

type BusinessForm struct {
	CSRFToken             *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	Name                  *string `json:"name" form:"name" schema:"name"`
	Website               *string `json:"website" form:"website" schema:"website"`
	Industry              *string `json:"industry" form:"industry" schema:"industry"`
	IsActive              *bool   `json:"is_active" form:"is_active" schema:"is_active"`
	GoogleBusinessProfile *string `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
}

type BusinessContactForm struct {
	CSRFToken              *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	FirstName              *string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName               *string `json:"last_name" form:"last_name" schema:"last_name"`
	Phone                  *string `json:"phone" form:"phone" schema:"phone"`
	Email                  *string `json:"email" form:"email" schema:"email"`
	PreferredContactMethod *string `json:"preferred_contact_method" form:"preferred_contact_method" schema:"preferred_contact_method"`
	PreferredContactTime   *string `json:"preferred_contact_time" form:"preferred_contact_time" schema:"preferred_contact_time"`
	BusinessPosition       *string `json:"business_position" form:"business_position" schema:"business_position"`
}

type LocationForm struct {
	CSRFToken            *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	VendingLocationID    *int     `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	LocationStatusID     *int     `json:"location_status_id" form:"location_status_id" schema:"location_status_id"`
	BusinessID           *int     `json:"business_id" form:"business_id" schema:"business_id"`
	LocationContact      *int     `json:"location_contact_id" form:"location_contact_id" schema:"location_contact_id"`
	DateStarted          int64    `json:"date_started" form:"date_started" schema:"date_started"`
	Name                 *string  `json:"name" form:"name" schema:"name"`
	Longitude            *float64 `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude             *float64 `json:"latitude" form:"latitude" schema:"latitude"`
	StreetAddressLineOne *string  `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo *string  `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	CityID               *int     `json:"city_id" form:"city_id" schema:"city_id"`
	ZipCode              *string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                *string  `json:"state" form:"state" schema:"state"`
	Opening              *string  `json:"opening" form:"opening" schema:"opening"`
	Closing              *string  `json:"closing" form:"closing" schema:"closing"`
}

type MachineForm struct {
	CSRFToken       *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	Year            *int     `json:"year" form:"year" schema:"year"`
	Make            *string  `json:"make" form:"make" schema:"make"`
	Model           *string  `json:"model" form:"model" schema:"model"`
	PurchasePrice   *float64 `json:"purchase_price" form:"purchase_price" schema:"purchase_price"`
	PurchaseDate    *int64   `json:"purchase_date" form:"purchase_date" schema:"purchase_date"`
	VendingTypeID   *int     `json:"vending_type_id" form:"vending_type_id" schema:"vending_type_id"`
	MachineStatusID *int     `json:"machine_status_id" form:"machine_status_id" schema:"machine_status_id"`
	VendorID        *int     `json:"vendor_id" form:"vendor_id" schema:"vendor_id"`
}

type MachineList struct {
	MachineID              int    `json:"machine_id" form:"machine_id" schema:"machine_id"`
	MachineName            string `json:"machine_name" form:"machine_name" schema:"machine_name"`
	CardReaderSerialNumber string `json:"card_reader_serial_number" form:"card_reader_serial_number" schema:"card_reader_serial_number"`
	MachineStatus          string `json:"machine_status" form:"machine_status" schema:"machine_status"`
	Location               string `json:"location" form:"location" schema:"location"`
	PurchaseDate           string `json:"purchase_date" form:"purchase_date" schema:"purchase_date"`
}

type VendorList struct {
	VendorID               int    `json:"vendor_id" form:"vendor_id" schema:"vendor_id"`
	Name                   string `json:"name" form:"name" schema:"name"`
	FirstName              string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName               string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber            string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email                  string `json:"email" form:"email" schema:"email"`
	PreferredContactMethod string `json:"preferred_contact_method" form:"preferred_contact_method" schema:"preferred_contact_method"`
	PreferredContactTime   string `json:"preferred_contact_time" form:"preferred_contact_time" schema:"preferred_contact_time"`
	StreetAddressLineOne   string `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo   string `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	CityID                 int    `json:"city_id" form:"city_id" schema:"city_id"`
	CityName               string `json:"city_name" form:"city_name" schema:"city_name"`
	ZipCode                string `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                  string `json:"state" form:"state" schema:"state"`
	GoogleBusinessProfile  string `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
}

type VendorForm struct {
	CSRFToken              *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	Name                   *string `json:"name" form:"name" schema:"name"`
	FirstName              *string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName               *string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber            *string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email                  *string `json:"email" form:"email" schema:"email"`
	PreferredContactMethod *string `json:"preferred_contact_method" form:"preferred_contact_method" schema:"preferred_contact_method"`
	PreferredContactTime   *string `json:"preferred_contact_time" form:"preferred_contact_time" schema:"preferred_contact_time"`
	StreetAddressLineOne   *string `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo   *string `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	CityID                 *int    `json:"city_id" form:"city_id" schema:"city_id"`
	ZipCode                *string `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                  *string `json:"state" form:"state" schema:"state"`
	GoogleBusinessProfile  *string `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
}

type SupplierForm struct {
	CSRFToken             *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	Name                  *string  `json:"name" form:"name" schema:"name"`
	MembershipID          *string  `json:"membership_id" form:"membership_id" schema:"membership_id"`
	MembershipCost        *float64 `json:"membership_cost" form:"membership_cost" schema:"membership_cost"` // Use float64 for MONEY
	MembershipRenewal     *string  `json:"membership_renewal" form:"membership_renewal" schema:"membership_renewal"`
	StreetAddressLineOne  *string  `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo  *string  `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	CityID                *int     `json:"city_id" form:"city_id" schema:"city_id"`
	ZipCode               *string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                 *string  `json:"state" form:"state" schema:"state"`
	GoogleBusinessProfile *string  `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
}

type SupplierList struct {
	SupplierID            int     `json:"supplier_id" form:"supplier_id" schema:"supplier_id"`
	Name                  string  `json:"name" form:"name" schema:"name"`
	MembershipID          string  `json:"membership_id" form:"membership_id" schema:"membership_id"`
	MembershipCost        float64 `json:"membership_cost" form:"membership_cost" schema:"membership_cost"` // Use float64 for MONEY
	MembershipRenewal     string  `json:"membership_renewal" form:"membership_renewal" schema:"membership_renewal"`
	StreetAddressLineOne  string  `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo  string  `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	ZipCode               string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                 string  `json:"state" form:"state" schema:"state"`
	GoogleBusinessProfile string  `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
	City                  string  `json:"city" form:"city" schema:"city"`
}

type VendorDetails struct {
	VendorID               int    `json:"vendor_id" form:"vendor_id" schema:"vendor_id"`
	Name                   string `json:"name" form:"name" schema:"name"`
	FirstName              string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName               string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber            string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email                  string `json:"email" form:"email" schema:"email"`
	PreferredContactMethod string `json:"preferred_contact_method" form:"preferred_contact_method" schema:"preferred_contact_method"`
	PreferredContactTime   string `json:"preferred_contact_time" form:"preferred_contact_time" schema:"preferred_contact_time"`
	StreetAddressLineOne   string `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo   string `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	CityID                 int    `json:"city_id" form:"city_id" schema:"city_id"`
	ZipCode                string `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                  string `json:"state" form:"state" schema:"state"`
	GoogleBusinessProfile  string `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
}

type SupplierDetails struct {
	SupplierID            int     `json:"supplier_id" form:"supplier_id" schema:"supplier_id"`
	Name                  string  `json:"name" form:"name" schema:"name"`
	MembershipID          string  `json:"membership_id" form:"membership_id" schema:"membership_id"`
	MembershipCost        float64 `json:"membership_cost" form:"membership_cost" schema:"membership_cost"`
	MembershipRenewal     string  `json:"membership_renewal" form:"membership_renewal" schema:"membership_renewal"`
	StreetAddressLineOne  string  `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo  string  `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	CityID                int     `json:"city_id" form:"city_id" schema:"city_id"`
	ZipCode               string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                 string  `json:"state" form:"state" schema:"state"`
	GoogleBusinessProfile string  `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
}

type BusinessDetails struct {
	BusinessID            int    `json:"business_id" form:"business_id" schema:"business_id"`
	Name                  string `json:"name" form:"name" schema:"name"`
	IsActive              bool   `json:"is_active" form:"is_active" schema:"is_active"`
	Website               string `json:"website" form:"website" schema:"website"`
	Industry              string `json:"industry" form:"industry" schema:"industry"`
	GoogleBusinessProfile string `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
}

type LocationDetails struct {
	LocationID           int     `json:"location_id" form:"location_id" schema:"location_id"`
	BusinessID           int     `json:"business_id" form:"business_id" schema:"business_id"`
	VendingLocationID    int     `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	CityID               int     `json:"city_id" form:"city_id" schema:"city_id"`
	LocationStatusID     int     `json:"location_status_id" form:"location_status_id" schema:"location_status_id"`
	DateStarted          int64   `json:"date_started" form:"date_started" schema:"date_started"`
	Name                 string  `json:"name" form:"name" schema:"name"`
	Longitude            float64 `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude             float64 `json:"latitude" form:"latitude" schema:"latitude"`
	StreetAddressLineOne string  `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo string  `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	ZipCode              string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                string  `json:"state" form:"state" schema:"state"`
	Opening              string  `json:"opening" form:"opening" schema:"opening"`
	Closing              string  `json:"closing" form:"closing" schema:"closing"`
}

type LocationList struct {
	LocationID           string  `json:"location_id" form:"location_id" schema:"location_id"`
	BusinessID           string  `json:"business_id" form:"business_id" schema:"business_id"`
	VendingLocationType  string  `json:"vending_location_type" form:"vending_location_type" schema:"vending_location_type"`
	Name                 string  `json:"name" form:"name" schema:"name"`
	Longitude            float64 `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude             float64 `json:"latitude" form:"latitude" schema:"latitude"`
	StreetAddressLineOne string  `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo string  `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	City                 string  `json:"city" form:"city" schema:"city"`
	ZipCode              string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                string  `json:"state" form:"state" schema:"state"`
	Opening              string  `json:"opening" form:"opening" schema:"opening"`
	Closing              string  `json:"closing" form:"closing" schema:"closing"`
	LocationStatus       string  `json:"location_status" form:"location_status" schema:"location_status"`
}

type MachineDetails struct {
	MachineID       int     `json:"machine_id" form:"machine_id" schema:"machine_id"`
	VendingTypeID   int     `json:"vending_type_id" form:"vending_type_id" schema:"vending_type_id"`
	MachineStatusID int     `json:"machine_status_id" form:"machine_status_id" schema:"machine_status_id"`
	VendorID        int     `json:"vendor_id" form:"vendor_id" schema:"vendor_id"`
	Year            int     `json:"year" form:"year" schema:"year"`
	Make            string  `json:"make" form:"make" schema:"make"`
	Model           string  `json:"model" form:"model" schema:"model"`
	PurchasePrice   float64 `json:"purchase_price" form:"purchase_price" schema:"purchase_price"`
	PurchaseDate    int64   `json:"purchase_date" form:"purchase_date" schema:"purchase_date"`
}

type SeedLiveTransaction struct {
	Month        string `json:"Month"`
	Day          string `json:"Day"`
	HourOfDay    string `json:"Hour of Day"`
	Device       string `json:"Device"`
	Location     string `json:"Location"`
	Item         string `json:"Item"`
	TransType    string `json:"Trans Type"`
	CardID       string `json:"Card Id"`
	CardNumber   string `json:"Card Number"`
	NumOfTrans   string `json:"# of Trans"`
	ItemQuantity string `json:"Item Quantity"`
}

type DashboardStats struct {
	Leads      int `json:"leads" form:"leads" schema:"leads"`
	Businesses int `json:"businesses" form:"businesses" schema:"businesses"`
	Vendors    int `json:"vendors" form:"vendors" schema:"vendors"`
	Suppliers  int `json:"suppliers" form:"suppliers" schema:"suppliers"`
	Machines   int `json:"machines" form:"machines" schema:"machines"`
}

type ProductForm struct {
	CSRFToken         *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	Name              *string  `json:"name" form:"name" schema:"name"`
	ProductCategoryID *int     `json:"product_category_id" form:"product_category_id" schema:"product_category_id"`
	Size              *float64 `json:"size" form:"size" schema:"size"`
	SizeType          *string  `json:"size_type" form:"size_type" schema:"size_type"`
	UPC               *string  `json:"upc" form:"upc" schema:"upc"`
}

type ProductList struct {
	ProductID int     `json:"product_id" form:"product_id" schema:"product_id"`
	Name      string  `json:"name" form:"name" schema:"name"`
	Category  string  `json:"category" form:"category" schema:"category"`
	Size      float64 `json:"size" form:"size" schema:"size"`
	SizeType  string  `json:"size_type" form:"size_type" schema:"size_type"`
	UPC       string  `json:"upc" form:"upc" schema:"upc"`
}

type ProductDetails struct {
	ProductID         int     `json:"product_id" form:"product_id" schema:"product_id"`
	Name              string  `json:"name" form:"name" schema:"name"`
	ProductCategoryID int     `json:"product_category_id" form:"product_category_id" schema:"product_category_id"`
	Size              float64 `json:"size" form:"size" schema:"size"`
	SizeType          string  `json:"size_type" form:"size_type" schema:"size_type"`
	UPC               string  `json:"upc" form:"upc" schema:"upc"`
}

type SlotList struct {
	SlotID       int     `json:"slot_id" form:"slot_id" schema:"slot_id"`
	Slot         string  `json:"slot" form:"slot" schema:"slot"`
	MachineID    int     `json:"machine_id" form:"machine_id" schema:"machine_id"`
	MachineCode  string  `json:"machine_code" form:"machine_code" schema:"machine_code"`
	Price        float64 `json:"price" form:"price" schema:"price"`
	Capacity     int     `json:"capacity" form:"capacity" schema:"capacity"`
	LastRefill   string  `json:"last_refill" form:"last_refill" schema:"last_refill"`
	LastRefillID int     `json:"last_refill_id" form:"last_refill_id" schema:"last_refill_id"`
}

type SlotDetails struct {
	SlotID        int    `json:"slot_id" form:"slot_id" schema:"slot_id"`
	Nickname      string `json:"nickname" form:"nickname" schema:"nickname"`
	Slot          string `json:"slot" form:"slot" schema:"slot"`
	MachineID     int    `json:"machine_id" form:"machine_id" schema:"machine_id"`
	MachineCode   string `json:"machine_code" form:"machine_code" schema:"machine_code"`
	Capacity      int    `json:"capacity" form:"capacity" schema:"capacity"`
	HasCommission bool   `json:"has_commission" form:"has_commission" schema:"has_commission"`
}

type SlotForm struct {
	CSRFToken               *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	Nickname                *string `json:"nickname" form:"nickname" schema:"nickname"`
	Slot                    *string `json:"slot" form:"slot" schema:"slot"`
	MachineCode             *string `json:"machine_code" form:"machine_code" schema:"machine_code"`
	MachineID               *int    `json:"machine_id" form:"machine_id" schema:"machine_id"`
	Capacity                *int    `json:"capacity" form:"capacity" schema:"capacity"`
	ProductSlotAssignmentID int64   `json:"product_slot_assignment_id" form:"product_slot_assignment_id" schema:"product_slot_assignment_id"`
}

type ProductSlotAssignment struct {
	ProductSlotAssignmentID int64   `json:"product_slot_assignment_id" form:"product_slot_assignment_id" schema:"product_slot_assignment_id"`
	Slot                    string  `json:"slot" form:"slot" schema:"slot"`
	DateAssigned            string  `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
	Product                 string  `json:"product" form:"product" schema:"product"`
	Supplier                string  `json:"supplier" form:"supplier" schema:"supplier"`
	UnitCost                float64 `json:"unit_cost" form:"unit_cost" schema:"unit_cost"`
	Quantity                int     `json:"quantity" form:"quantity" schema:"quantity"`
	ExpirationDate          string  `json:"expiration_date" form:"expiration_date" schema:"expiration_date"`
}

type ProductSlotAssignmentDetails struct {
	ProductSlotAssignmentID int64   `json:"product_slot_assignment_id" form:"product_slot_assignment_id" schema:"product_slot_assignment_id"`
	SlotID                  int     `json:"slot_id" form:"slot_id" schema:"slot_id"`
	DateAssigned            int64   `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
	ProductID               int     `json:"product_id" form:"product_id" schema:"product_id"`
	SupplierID              int     `json:"supplier_id" form:"supplier_id" schema:"supplier_id"`
	UnitCost                float64 `json:"unit_cost" form:"unit_cost" schema:"unit_cost"`
	Quantity                int     `json:"quantity" form:"quantity" schema:"quantity"`
	ExpirationDate          int64   `json:"expiration_date" form:"expiration_date" schema:"expiration_date"`
}

type ProductSlotAssignmentForm struct {
	CSRFToken               *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	ProductSlotAssignmentID *int     `json:"product_slot_assignment_id" form:"product_slot_assignment_id" schema:"product_slot_assignment_id"`
	SlotID                  *int     `json:"slot_id" form:"slot_id" schema:"slot_id"`
	ProductID               *int     `json:"product_id" form:"product_id" schema:"product_id"`
	DateAssigned            *int64   `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
	SupplierID              *int     `json:"supplier_id" form:"supplier_id" schema:"supplier_id"`
	UnitCost                *float64 `json:"unit_cost" form:"unit_cost" schema:"unit_cost"`
	Quantity                *int     `json:"quantity" form:"quantity" schema:"quantity"`
	ExpirationDate          *int64   `json:"expiration_date" form:"expiration_date" schema:"expiration_date"`
}

type RefillForm struct {
	CSRFToken    *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	SlotID       *int    `json:"slot_id" form:"slot_id" schema:"slot_id"`
	DateRefilled *int64  `json:"date_refilled" form:"date_refilled" schema:"date_refilled"`
}

type GetTransactionsParams struct {
	TransactionType *string `json:"transaction_type" schema:"transaction_type" form:"transaction_type"`
	Machine         *string `json:"machine" schema:"machine" form:"machine"`
	Location        *string `json:"location" schema:"location" form:"location"`
	Product         *string `json:"product" schema:"product" form:"product"`
	PageNum         *string `json:"page_num" schema:"page_num" form:"page_num"`
	DateFrom        *int64  `json:"date_from" form:"date_from" schema:"date_from"`
	DateTo          *int64  `json:"date_to" form:"date_to" schema:"date_to"`
}

type TransactionList struct {
	TransactionLogID        int64   `json:"transaction_log_id" db:"transaction_log_id" form:"transaction_log_id" schema:"transaction_log_id"`
	Machine                 string  `json:"machine" db:"machine" form:"machine" schema:"machine"`
	Location                string  `json:"location" db:"location" form:"location" schema:"location"`
	MachineSelection        string  `json:"machine_selection" db:"machine_selection" form:"machine_selection" schema:"machine_selection"`
	Product                 string  `json:"product" db:"product" form:"product" schema:"product"`
	TransactionTimestamp    string  `json:"transaction_timestamp" db:"transaction_timestamp" form:"transaction_timestamp" schema:"transaction_timestamp"`
	TransactionType         string  `json:"transaction_type" db:"transaction_type" form:"transaction_type" schema:"transaction_type"`
	CardNumber              string  `json:"card_number" db:"card_number" form:"card_number" schema:"card_number"`
	Revenue                 float64 `json:"revenue" db:"revenue" form:"revenue" schema:"revenue"`
	Items                   int     `json:"items" db:"items" form:"items" schema:"items"`
	IsValidated             bool    `json:"is_validated" db:"is_validated" form:"is_validated" schema:"is_validated"`
	TransactionValidationID int     `json:"transaction_validation_id" db:"transaction_validation_id" form:"transaction_validation_id" schema:"transaction_validation_id"`
}

type SlotPriceLogList struct {
	SlotPriceLogID int     `json:"slot_price_log_id" form:"slot_price_log_id" schema:"slot_price_log_id"`
	MachineID      int     `json:"machine_id" form:"machine_id" schema:"machine_id"`
	SlotID         int     `json:"slot_id" form:"slot_id" schema:"slot_id"`
	Price          float64 `json:"price" form:"price" schema:"price"`
	DateAssigned   int64   `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
}

type SlotPriceLogForm struct {
	CSRFToken      *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	SlotPriceLogID *int     `json:"slot_price_log_id" form:"slot_price_log_id" schema:"slot_price_log_id"`
	MachineID      *int     `json:"machine_id" form:"machine_id" schema:"machine_id"`
	SlotID         *int     `json:"slot_id" form:"slot_id" schema:"slot_id"`
	Price          *float64 `json:"price" form:"price" schema:"price"`
	DateAssigned   *int64   `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
}

type PrepReport struct {
	Machine    string `json:"machine" form:"machine" schema:"machine"`
	Location   string `json:"location" form:"location" schema:"location"`
	Product    string `json:"product" form:"product" schema:"product"`
	AmountSold string `json:"amount_sold" form:"amount_sold" schema:"amount_sold"`
}

type CommissionReport struct {
	Product       string  `json:"product" form:"product" schema:"product" spreadsheet_header:"Product"`
	AmountSold    float64 `json:"amount_sold" form:"amount_sold" schema:"amount_sold" spreadsheet_header:"Amount Sold"`
	Revenue       float64 `json:"revenue" form:"revenue" schema:"revenue" spreadsheet_header:"Revenue"`
	Cost          float64 `json:"cost" form:"cost" schema:"cost" spreadsheet_header:"Cost"`
	CreditCardFee float64 `json:"credit_card_fee" form:"credit_card_fee" schema:"credit_card_fee" spreadsheet_header:"Credit Card Fee"`
	GrossProfit   float64 `json:"gross_profit" form:"gross_profit" schema:"gross_profit" spreadsheet_header:"Gross Profit"`
	CommissionDue float64 `json:"commission_due" form:"commission_due" schema:"commission_due" spreadsheet_header:"Commission Due"`
}

type EmailScheduleForm struct {
	CSRFToken       *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	EmailName       *string `json:"email_name" form:"email_name" schema:"email_name"`
	IntervalSeconds *int64  `json:"interval_seconds" form:"interval_seconds" schema:"interval_seconds"`
	Recipients      *string `json:"recipients" form:"recipients" schema:"recipients"`
	Subject         *string `json:"subject" form:"subject" schema:"subject"`
	Body            *string `json:"body" form:"body" schema:"body"`
	Sender          *string `json:"sender" form:"sender" schema:"sender"`
	SQLFile         *string `json:"sql_file" form:"sql_file" schema:"sql_file"`
	LastSent        *int64  `json:"last_sent" form:"last_sent" schema:"last_sent"`
	IsActive        *bool   `json:"is_active" form:"is_active" schema:"is_active"`
}

type EmailScheduleList struct {
	EmailScheduleID int    `json:"email_schedule_id" form:"email_schedule_id" schema:"email_schedule_id"`
	EmailName       string `json:"email_name" form:"email_name" schema:"email_name"`
	IntervalSeconds int64  `json:"interval_seconds" form:"interval_seconds" schema:"interval_seconds"`
	Recipients      string `json:"recipients" form:"recipients" schema:"recipients"`
	Subject         string `json:"subject" form:"subject" schema:"subject"`
	Body            string `json:"body" form:"body" schema:"body"`
	Sender          string `json:"sender" form:"sender" schema:"sender"`
	SQLFile         string `json:"sql_file" form:"sql_file" schema:"sql_file"`
	LastSent        string `json:"last_sent" form:"last_sent" schema:"last_sent"`
	IsActive        bool   `json:"is_active" form:"is_active" schema:"is_active"`
}

type SentEmail struct {
	CSRFToken       *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	EmailScheduleID *int    `json:"email_schedule_id" form:"email_schedule_id" schema:"email_schedule_id"`
	DeliveryStatus  *string `json:"delivery_status" form:"delivery_status" schema:"delivery_status"`
	DateSent        *int64  `json:"date_sent" form:"date_sent" schema:"date_sent"`
	ErrorMessage    *string `json:"error_message" form:"error_message" schema:"error_message"`
}

type MachineLocationAssignmentForm struct {
	CSRFToken            *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	LocationID           *int    `json:"location_id" form:"location_id" schema:"location_id"`
	MachineID            *int    `json:"machine_id" form:"machine_id" schema:"machine_id"`
	LocationDateAssigned *int64  `json:"location_date_assigned" form:"location_date_assigned" schema:"location_date_assigned"`
	IsLocationActive     *bool   `json:"is_location_active" form:"is_location_active" schema:"is_location_active"`
}

type MachineCardReaderAssignmentForm struct {
	CSRFToken                     *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	CardReaderSerialNumber        *string `json:"card_reader_serial_number" form:"card_reader_serial_number" schema:"card_reader_serial_number"`
	MachineID                     *int    `json:"machine_id" form:"machine_id" schema:"machine_id"`
	MachineCardReaderDateAssigned *int64  `json:"machine_card_reader_date_assigned" form:"machine_card_reader_date_assigned" schema:"machine_card_reader_date_assigned"`
	IsCardReaderActive            *bool   `json:"is_card_reader_active" form:"is_card_reader_active" schema:"is_card_reader_active"`
}

type MachineLocationAssignment struct {
	MachineLocationAssignmentID int   `json:"machine_location_assignment_id" form:"machine_location_assignment_id" schema:"machine_location_assignment_id"`
	LocationID                  int   `json:"location_id" form:"location_id" schema:"location_id"`
	MachineID                   int   `json:"machine_id" form:"machine_id" schema:"machine_id"`
	LocationDateAssigned        int64 `json:"location_date_assigned" form:"location_date_assigned" schema:"location_date_assigned"`
	IsLocationActive            bool  `json:"is_location_active" form:"is_location_active" schema:"is_location_active"`
}

type MachineCardReaderAssignment struct {
	MachineCardReaderID           int    `json:"machine_card_reader_assignment_id" form:"machine_card_reader_assignment_id" schema:"machine_card_reader_assignment_id"`
	CardReaderSerialNumber        string `json:"card_reader_serial_number" form:"card_reader_serial_number" schema:"card_reader_serial_number"`
	MachineID                     int    `json:"machine_id" form:"machine_id" schema:"machine_id"`
	MachineCardReaderDateAssigned int64  `json:"machine_card_reader_date_assigned" form:"machine_card_reader_date_assigned" schema:"machine_card_reader_date_assigned"`
	IsCardReaderActive            bool   `json:"is_card_reader_active" form:"is_card_reader_active" schema:"is_card_reader_active"`
}

type OptIn90DayChallengeForm struct {
	Email            *string  `json:"email" form:"email" schema:"email"`
	LandingPageID    *int     `json:"landing_page_id" form:"landing_page_id" schema:"landing_page_id"`
	PercentScrolled  *float64 `json:"percent_scrolled" form:"percent_scrolled" schema:"percent_scrolled"`
	TimeSpentOnPage  *int64   `json:"time_spent_on_page" form:"time_spent_on_page" schema:"time_spent_on_page"`
	Source           *string  `json:"source" form:"source" schema:"source"`
	Medium           *string  `json:"medium" form:"medium" schema:"medium"`
	Channel          *string  `json:"channel" form:"channel" schema:"channel"`
	LandingPage      *string  `json:"landing_page" form:"landing_page" schema:"landing_page"`
	Keyword          *string  `json:"keyword" form:"keyword" schema:"keyword"`
	Referrer         *string  `json:"referrer" form:"referrer" schema:"referrer"`
	ClickID          *string  `json:"click_id" form:"click_id" schema:"click_id"`
	CampaignID       *int64   `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	AdCampaign       *string  `json:"ad_campaign" form:"ad_campaign" schema:"ad_campaign"`
	AdGroupID        *int64   `json:"ad_group_id" form:"ad_group_id" schema:"ad_group_id"`
	AdGroupName      *string  `json:"ad_group_name" form:"ad_group_name" schema:"ad_group_name"`
	AdSetID          *int64   `json:"ad_set_id" form:"ad_set_id" schema:"ad_set_id"`
	AdSetName        *string  `json:"ad_set_name" form:"ad_set_name" schema:"ad_set_name"`
	AdID             *int64   `json:"ad_id" form:"ad_id" schema:"ad_id"`
	AdHeadline       *int64   `json:"ad_headline" form:"ad_headline" schema:"ad_headline"`
	Language         *string  `json:"language" form:"language" schema:"language"`
	Longitude        *string  `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude         *string  `json:"latitude" form:"latitude" schema:"latitude"`
	UserAgent        *string  `json:"user_agent" form:"user_agent" schema:"user_agent"`
	ButtonClicked    *string  `json:"button_clicked" form:"button_clicked" schema:"button_clicked"`
	IP               *string  `json:"ip" form:"ip" schema:"ip"`
	CSRFToken        *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	ExternalID       *string  `json:"external_id" form:"external_id" schema:"external_id"`
	GoogleClientID   *string  `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	FacebookClickID  *string  `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID *string  `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	CSRFSecret       *string  `json:"csrf_secret" form:"csrf_secret"`
}

type LeadApplicationForm struct {
	FirstName   *string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName    *string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber *string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email       *string `json:"email" form:"email" schema:"email"`

	Website         *string `json:"website" form:"website" schema:"website"`
	CompanyName     *string `json:"company_name" form:"company_name" schema:"company_name"`
	YearsInBusiness *int    `json:"years_in_business" form:"years_in_business" schema:"years_in_business"`
	NumLocations    *int    `json:"num_locations" form:"num_locations" schema:"num_locations"`
	City            *string `json:"city" form:"city" schema:"city"`

	LandingPageID    *int     `json:"landing_page_id" form:"landing_page_id" schema:"landing_page_id"`
	LeadTypeID       *int     `json:"lead_type_id" form:"lead_type_id" schema:"lead_type_id"`
	PercentScrolled  *float64 `json:"percent_scrolled" form:"percent_scrolled" schema:"percent_scrolled"`
	TimeSpentOnPage  *int64   `json:"time_spent_on_page" form:"time_spent_on_page" schema:"time_spent_on_page"`
	Source           *string  `json:"source" form:"source" schema:"source"`
	Medium           *string  `json:"medium" form:"medium" schema:"medium"`
	Channel          *string  `json:"channel" form:"channel" schema:"channel"`
	LandingPage      *string  `json:"landing_page" form:"landing_page" schema:"landing_page"`
	Keyword          *string  `json:"keyword" form:"keyword" schema:"keyword"`
	Referrer         *string  `json:"referrer" form:"referrer" schema:"referrer"`
	ClickID          *string  `json:"click_id" form:"click_id" schema:"click_id"`
	CampaignID       *int64   `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	AdCampaign       *string  `json:"ad_campaign" form:"ad_campaign" schema:"ad_campaign"`
	AdGroupID        *int64   `json:"ad_group_id" form:"ad_group_id" schema:"ad_group_id"`
	AdGroupName      *string  `json:"ad_group_name" form:"ad_group_name" schema:"ad_group_name"`
	AdSetID          *int64   `json:"ad_set_id" form:"ad_set_id" schema:"ad_set_id"`
	AdSetName        *string  `json:"ad_set_name" form:"ad_set_name" schema:"ad_set_name"`
	AdID             *int64   `json:"ad_id" form:"ad_id" schema:"ad_id"`
	AdHeadline       *int64   `json:"ad_headline" form:"ad_headline" schema:"ad_headline"`
	Language         *string  `json:"language" form:"language" schema:"language"`
	Longitude        *string  `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude         *string  `json:"latitude" form:"latitude" schema:"latitude"`
	UserAgent        *string  `json:"user_agent" form:"user_agent" schema:"user_agent"`
	ButtonClicked    *string  `json:"button_clicked" form:"button_clicked" schema:"button_clicked"`
	IP               *string  `json:"ip" form:"ip" schema:"ip"`
	CSRFToken        *string  `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	ExternalID       *string  `json:"external_id" form:"external_id" schema:"external_id"`
	GoogleClientID   *string  `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	FacebookClickID  *string  `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID *string  `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	CSRFSecret       *string  `json:"csrf_secret" form:"csrf_secret"`
}

type LeadAppointmentForm struct {
	LeadID     *int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	BookedTime *int64  `json:"booked_time" form:"booked_time" schema:"booked_time"`
	Attendee   *string `json:"attendee" form:"attendee" schema:"attendee"`
	CSRFToken  *string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
}
