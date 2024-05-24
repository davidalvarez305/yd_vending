package types

type QuoteForm struct {
	FirstName     string `json:"first_name" form:"first_name"`
	LastName      string `json:"last_name" form:"last_name"`
	PhoneNumber   string `json:"phone_number" form:"phone_number"`
	Rent          string `json:"rent" form:"rent"`
	Service       string `json:"service" form:"service"`
	City          string `json:"city" form:"city"`
	Message       string `json:"message" form:"message"`
	Source        string `json:"source"`
	Medium        string `json:"medium"`
	Channel       string `json:"channel"`
	LandingPage   string `json:"landing_page"`
	Keyword       string `json:"keyword"`
	Referrer      string `json:"referrer"`
	Gclid         string `json:"gclid"`
	CampaignID    string `json:"campaign_id"`
	AdCampaign    string `json:"ad_campaign"`
	AdGroupID     string `json:"ad_group_id"`
	AdGroupName   string `json:"ad_group_name"`
	AdSetID       string `json:"ad_set_id"`
	AdSetName     string `json:"ad_set_name"`
	AdID          string `json:"ad_id"`
	AdHeadline    string `json:"ad_headline"`
	Language      string `json:"language"`
	OS            string `json:"os"`
	UserAgent     string `json:"user_agent"`
	ButtonClicked string `json:"button_clicked"`
	DeviceType    string `json:"device_type"`
	IP            string `json:"ip"`
}
