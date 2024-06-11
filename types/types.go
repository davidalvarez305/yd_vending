package types

type QuoteForm struct {
	FirstName       string `json:"first_name" form:"first_name"`
	LastName        string `json:"last_name" form:"last_name"`
	PhoneNumber     string `json:"phone_number" form:"phone_number"`
	Rent            string `json:"rent" form:"rent"`
	City            int    `json:"city" form:"city"`
	LocationType    int    `json:"location_type" form:"location_type"`
	MachineType     int    `json:"machine_type" form:"machine_type"`
	FootTraffic     string `json:"foot_traffic" form:"foot_traffic"`
	FootTrafficType string `json:"foot_traffic_type" form:"foot_traffic_type"`
	Message         string `json:"message" form:"message"`
	Source          string `json:"source" form:"source"`
	Medium          string `json:"medium" form:"medium"`
	Channel         string `json:"channel" form:"channel"`
	LandingPage     string `json:"landing_page" form:"landing_page"`
	Keyword         string `json:"keyword" form:"keyword"`
	Referrer        string `json:"referrer" form:"referrer"`
	Gclid           string `json:"gclid" form:"gclid"`
	CampaignID      int    `json:"campaign_id" form:"campaign_id"`
	AdCampaign      string `json:"ad_campaign" form:"ad_campaign"`
	AdGroupID       int    `json:"ad_group_id" form:"ad_group_id"`
	AdGroupName     string `json:"ad_group_name" form:"ad_group_name"`
	AdSetID         int    `json:"ad_set_id" form:"ad_set_id"`
	AdSetName       string `json:"ad_set_name" form:"ad_set_name"`
	AdID            int    `json:"ad_id" form:"ad_id"`
	AdHeadline      int    `json:"ad_headline" form:"ad_headline"`
	Language        string `json:"language" form:"language"`
	OS              string `json:"os" form:"os"`
	UserAgent       string `json:"user_agent" form:"user_agent"`
	ButtonClicked   string `json:"button_clicked" form:"button_clicked"`
	DeviceType      string `json:"device_type" form:"device_type"`
	IP              string `json:"ip" form:"ip"`
	CSRFToken       string `json:"csrf_token" form:"csrf_token"`
}

type ContactForm struct {
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"last_name" form:"last_name"`
	Email     string `json:"email" form:"email"`
	Message   string `json:"message" form:"message"`
}

type OutboundMessageForm struct {
	To   string `json:"to" form:"to"`
	Body string `json:"body" form:"body"`
	From string `json:"from" form:"from"`
}
