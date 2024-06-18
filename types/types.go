package types

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

type SuccessModal struct {
	TemplateName string
	AlertHeader  string
	AlertMessage string
}
