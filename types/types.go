package types

import "time"

type QuoteForm struct {
	FirstName        string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName         string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber      string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Rent             string `json:"rent" form:"rent" schema:"rent"`
	City             int    `json:"city" form:"city" schema:"city"`
	LocationType     int    `json:"location_type" form:"location_type" schema:"location_type"`
	MachineType      int    `json:"machine_type" form:"machine_type" schema:"machine_type"`
	FootTraffic      string `json:"foot_traffic" form:"foot_traffic" schema:"foot_traffic"`
	FootTrafficType  string `json:"foot_traffic_type" form:"foot_traffic_type" schema:"foot_traffic_type"`
	Message          string `json:"message" form:"message" schema:"message"`
	Source           string `json:"source" form:"source" schema:"source"`
	Medium           string `json:"medium" form:"medium" schema:"medium"`
	Channel          string `json:"channel" form:"channel" schema:"channel"`
	LandingPage      string `json:"landing_page" form:"landing_page" schema:"landing_page"`
	Keyword          string `json:"keyword" form:"keyword" schema:"keyword"`
	Referrer         string `json:"referrer" form:"referrer" schema:"referrer"`
	GCLID            string `json:"gclid" form:"gclid" schema:"gclid"`
	CampaignID       int    `json:"campaign_id" form:"campaign_id" schema:"campaign_id"`
	AdCampaign       string `json:"ad_campaign" form:"ad_campaign" schema:"ad_campaign"`
	AdGroupID        int    `json:"ad_group_id" form:"ad_group_id" schema:"ad_group_id"`
	AdGroupName      string `json:"ad_group_name" form:"ad_group_name" schema:"ad_group_name"`
	AdSetID          int    `json:"ad_set_id" form:"ad_set_id" schema:"ad_set_id"`
	AdSetName        string `json:"ad_set_name" form:"ad_set_name" schema:"ad_set_name"`
	AdID             int    `json:"ad_id" form:"ad_id" schema:"ad_id"`
	AdHeadline       int    `json:"ad_headline" form:"ad_headline" schema:"ad_headline"`
	Language         string `json:"language" form:"language" schema:"language"`
	Longitude        string `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude         string `json:"latitude" form:"latitude" schema:"latitude"`
	UserAgent        string `json:"user_agent" form:"user_agent" schema:"user_agent"`
	ButtonClicked    string `json:"button_clicked" form:"button_clicked" schema:"button_clicked"`
	IP               string `json:"ip" form:"ip" schema:"ip"`
	CSRFToken        string `json:"csrf_token" form:"csrf_token" schema:"csrf_token"`
	GoogleUserID     string `json:"google_user_id" form:"google_user_id" schema:"google_user_id"`
	GoogleClientID   string `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	FacebookClickID  string `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID string `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	CSRFSecret       []byte `json:"csrf_secret"`
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
	FullName        string `json:"full_name" form:"full_name" schema:"full_name"`
	PhoneNumber     string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email           string `json:"email" form:"email" schema:"email"`
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
	VendingType     string `json:"vending_type" form:"vending_type"`
	LocationType    string `json:"location_type" form:"location_type"`
	City            string `json:"city" form:"city"`
	SearchFieldType string `json:"search_field_type" form:"search_field_type"`
	PageNum         int    `json:"page_num" form:"page_num" schema:"page_num"`
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
