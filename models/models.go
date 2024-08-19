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

type User struct {
	UserID      int    `json:"user_id"`
	Username    string `json:"username"`
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
	IsAdmin     bool   `json:"is_admin"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
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
	LeadStatusID      int    `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
	Message           string `json:"message"`
}

type LeadStatus struct {
	LeadStatusID int    `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
	Status       string `json:"status" form:"status" schema:"status"`
}

type LeadStatusLog struct {
	LeadStatusLogID int   `json:"lead_status_log_id" form:"lead_status_log_id" schema:"lead_status_log_id"`
	LeadID          int   `json:"lead_id" form:"lead_id" schema:"lead_id"`
	LeadStatusID    int   `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
	DateAdded       int64 `json:"date_added" form:"date_added" schema:"date_added"`
	ChangedByUserID int   `json:"changed_by_user_id" form:"changed_by_user_id" schema:"changed_by_user_id"`
}

type LeadNote struct {
	LeadNoteID    int    `json:"lead_note_id" form:"lead_note_id" schema:"lead_note_id"`
	Note          string `json:"note" form:"note" schema:"note"`
	LeadID        int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	DateAdded     int64  `json:"date_added" form:"date_added" schema:"date_added"`
	AddedByUserID int    `json:"added_by_user_id" form:"added_by_user_id" schema:"added_by_user_id"`
}

type LeadMarketing struct {
	LeadMarketingID  int    `json:"lead_marketing_id"`
	LeadID           int    `json:"lead_id"`
	Source           string `json:"source"`
	Medium           string `json:"medium"`
	Channel          string `json:"channel"`
	LandingPage      string `json:"landing_page"`
	Longitude        string `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude         string `json:"latitude" form:"latitude" schema:"latitude"`
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
	ExternalID       string `json:"external_id"`
	GoogleClientID   string `json:"google_client_id"`
	FacebookClickID  string `json:"facebook_click_id"`
	FacebookClientID string `json:"facebook_client_id"`
	CSRFSecret       string `json:"csrf_secret"`
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
	Status      string `json:"status" form:"status" schema:"status"`
}

type PhoneCall struct {
	PhoneCallID  int    `json:"phone_call_id" form:"phone_call_id" schema:"phone_call_id"`
	ExternalID   string `json:"external_id" form:"external_id" schema:"external_id"`
	UserID       int    `json:"user_id" form:"user_id" schema:"user_id"`
	LeadID       int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	CallDuration int    `json:"call_duration" form:"call_duration" schema:"call_duration"`
	DateCreated  int64  `json:"date_created" form:"date_created" schema:"date_created"`
	CallFrom     string `json:"call_from" form:"call_from" schema:"call_from"`
	CallTo       string `json:"call_to" form:"call_to" schema:"call_to"`
	IsInbound    bool   `json:"is_inbound" form:"is_inbound" schema:"is_inbound"`
	RecordingURL string `json:"recording_url" form:"recording_url" schema:"recording_url"`
	Status       string `json:"status" form:"status" schema:"status"`
}

type Session struct {
	SessionID        int    `json:"session_id" form:"session_id" schema:"session_id"`
	UserID           int    `json:"user_id" form:"user_id" schema:"user_id"`
	CSRFSecret       string `json:"csrf_secret" form:"csrf_secret" schema:"csrf_secret"`
	ExternalID       string `json:"external_id" form:"external_id" schema:"external_id"`
	GoogleClientID   string `json:"google_client_id" form:"google_client_id" schema:"google_client_id"`
	FacebookClickID  string `json:"facebook_click_id" form:"facebook_click_id" schema:"facebook_click_id"`
	FacebookClientID string `json:"facebook_client_id" form:"facebook_client_id" schema:"facebook_client_id"`
	DateCreated      int64  `json:"date_created" form:"date_created" schema:"date_created"`
	DateExpires      int64  `json:"date_expires" form:"date_expires" schema:"date_expires"`
}
