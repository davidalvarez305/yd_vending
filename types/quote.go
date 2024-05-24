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
	Source          string `json:"source"`
	Medium          string `json:"medium"`
	Channel         string `json:"channel"`
	LandingPage     string `json:"landing_page"`
	Keyword         string `json:"keyword"`
	Referrer        string `json:"referrer"`
	Gclid           string `json:"gclid"`
	CampaignID      int    `json:"campaign_id"`
	AdCampaign      string `json:"ad_campaign"`
	AdGroupID       int    `json:"ad_group_id"`
	AdGroupName     string `json:"ad_group_name"`
	AdSetID         int    `json:"ad_set_id"`
	AdSetName       string `json:"ad_set_name"`
	AdID            int    `json:"ad_id"`
	AdHeadline      int    `json:"ad_headline"`
	Language        string `json:"language"`
	OS              string `json:"os"`
	UserAgent       string `json:"user_agent"`
	ButtonClicked   string `json:"button_clicked"`
	DeviceType      string `json:"device_type"`
	IP              string `json:"ip"`
}
