package models

type VendingType struct {
	VendingTypeID int    `json:"vending_type_id"`
	MachineType   string `json:"machine_type"`
}

type VendingLocation struct {
	VendingLocationID int    `json:"vending_location_id"`
	LocationType      string `json:"location_type"`
}

type City struct {
	CityID int    `json:"city_id"`
	Name   string `json:"name"`
}

type Lead struct {
	LeadID            int    `json:"lead_id"`
	FirstName         string `json:"first_name"`
	LastName          string `json:"last_name"`
	PhoneNumber       string `json:"phone_number"`
	CreatedAt         int64  `json:"created_at"`
	Rent              string `json:"rent"`
	FootTraffic       int    `json:"foot_traffic"`
	FootTrafficType   string `json:"foot_traffic_type"`
	VendingTypeID     int    `json:"vending_type_id"`
	VendingLocationID int    `json:"vending_location_id"`
	CityID            int    `json:"city_id"`
	UserKey           string `json:"user_key"`
}

type LeadMarketing struct {
	LeadMarketingID  int    `json:"lead_marketing_id"`
	LeadID           int    `json:"lead_id"`
	Source           string `json:"source"`
	Medium           string `json:"medium"`
	Channel          string `json:"channel"`
	LandingPage      string `json:"landing_page"`
	Keyword          string `json:"keyword"`
	Referrer         string `json:"referrer"`
	GCLID            string `json:"gclid"`
	CampaignID       int    `json:"campaign_id"`
	AdCampaign       string `json:"ad_campaign"`
	AdGroupID        int    `json:"ad_group_id"`
	AdGroupName      string `json:"ad_group_name"`
	AdSetID          int    `json:"ad_set_id"`
	AdSetName        string `json:"ad_set_name"`
	AdID             int    `json:"ad_id"`
	AdHeadline       int    `json:"ad_headline"`
	Language         string `json:"language"`
	OS               string `json:"os"`
	UserAgent        string `json:"user_agent"`
	ButtonClicked    string `json:"button_clicked"`
	DeviceType       string `json:"device_type"`
	IP               string `json:"ip"`
	GoogleUserID     string `json:"google_user_id"`
	GoogleClientID   string `json:"google_client_id"`
	FacebookClickID  string `json:"facebook_click_id"`
	FacebookClientID string `json:"facebook_client_id"`
	CSRFSecret       []byte `json:"csrf_secret"`
}

type CSRFToken struct {
	CSRFTokenID int    `json:"csrf_token_id"`
	ExpiryTime  int64  `json:"expiry_time"`
	Token       string `json:"token"`
	IsUsed      bool   `json:"is_used"`
}

type Message struct {
	MessageID   int    `json:"message_id"`
	ExternalID  string `json:"external_id"`
	UserID      int    `json:"user_id"`
	LeadID      int    `json:"lead_id"`
	Text        string `json:"text"`
	DateCreated int64  `json:"date_created"`
	TextFrom    string `json:"text_from"`
	TextTo      string `json:"text_to"`
	IsInbound   bool   `json:"is_inbound"`
}

type User struct {
	UserID      int    `json:"user_id"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	IsAdmin     bool   `json:"is_admin"`
}
