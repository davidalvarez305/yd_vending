package models

import "time"

// Vending machine, ice machine, coffee machine, snack machine, healthy snack machine, etc...
type VendingType struct {
	VendingTypeID int    `gorm:"primaryKey;column:vending_type_id" json:"vending_type_id"`
	MachineType   string `gorm:"unique;not null;column:machine_type" json:"machine_type"`
}

// Barbershop, office building, gym, etc...
type VendingLocation struct {
	VendingLocationID int    `gorm:"primaryKey;column:vending_location_id" json:"vending_location_id"`
	LocationType      string `gorm:"unique;not null;column:location_type" json:"location_type"`
}

type City struct {
	CityID int    `gorm:"primaryKey;column:city_id" json:"city_id"`
	Name   string `gorm:"unique;not null;column:name" json:"name"`
}

type Lead struct {
	LeadID            int              `gorm:"primaryKey;column:lead_id" json:"lead_id"`
	FirstName         string           `gorm:"not null;column:first_name" json:"first_name"`
	LastName          string           `gorm:"not null;column:last_name" json:"last_name"`
	PhoneNumber       string           `gorm:"not null;column:phone_number" json:"phone_number"`
	CreatedAt         int64            `gorm:"not null;column:created_at" json:"created_at"`
	Rent              string           `gorm:"column:rent" json:"rent"`
	FootTraffic       int              `gorm:"column:foot_traffic" json:"foot_traffic"`
	FootTrafficType   string           `gorm:"column:foot_traffic_type" json:"foot_traffic_type"`
	VendingTypeID     int              `gorm:"column:vending_type_id" json:"vending_type_id"`
	VendingType       *VendingType     `gorm:"not null;column:vending_type_id;foreignKey:VendingTypeID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"vending_type"`
	VendingLocationID int              `gorm:"column:vending_location_id" json:"vending_location_id"`
	VendingLocation   *VendingLocation `gorm:"not null;column:vending_location_id;foreignKey:VendingLocationID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"vending_location"`
	CityID            int              `gorm:"column:city_id" json:"city_id"`
	City              *City            `gorm:"not null;column:city_id;foreignKey:CityID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"city"`
	LeadMarketing     *LeadMarketing   `json:"lead_marketing"`
}

type LeadMarketing struct {
	LeadMarketingID int    `gorm:"primaryKey;column:lead_marketing_id" json:"lead_marketing_id"`
	LeadID          int    `gorm:"not null;column:lead_id;foreignKey:LeadID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE" json:"lead_id"`
	Source          string `gorm:"column:source" json:"source"`
	Medium          string `gorm:"column:medium" json:"medium"`
	Channel         string `gorm:"column:channel" json:"channel"`
	LandingPage     string `gorm:"column:landing_page" json:"landing_page"`
	Keyword         string `gorm:"column:keyword" json:"keyword"`
	Referrer        string `gorm:"column:referrer" json:"referrer"`
	Gclid           string `gorm:"column:gclid" json:"gclid"`
	CampaignID      int    `gorm:"column:campaign_id" json:"campaign_id"`
	AdCampaign      string `gorm:"column:ad_campaign" json:"ad_campaign"`
	AdGroupID       int    `gorm:"column:ad_group_id" json:"ad_group_id"`
	AdGroupName     string `gorm:"column:ad_group_name" json:"ad_group_name"`
	AdSetID         int    `gorm:"column:ad_set_id" json:"ad_set_id"`
	AdSetName       string `gorm:"column:ad_set_name" json:"ad_set_name"`
	AdID            int    `gorm:"column:ad_id" json:"ad_id"`
	AdHeadline      int    `gorm:"column:ad_headline" json:"ad_headline"`
	Language        string `gorm:"column:language" json:"language"`
	OS              string `gorm:"column:os" json:"os"`
	UserAgent       string `gorm:"column:user_agent" json:"user_agent"`
	ButtonClicked   string `gorm:"column:button_clicked" json:"button_clicked"`
	DeviceType      string `gorm:"column:device_type" json:"device_type"`
	IP              string `gorm:"column:ip" json:"ip"`
}

type CSRFToken struct {
	CSRFTokenID int    `json:"csrf_token_id"`
	ExpiryTime  int64  `json:"expiry_time"`
	Token       string `json:"token"`
	IsUsed      bool   `json:"is_used"`
}

type TextMessage struct {
	TextMessageID int       `json:"text_message_id"`
	MessageSID    string    `json:"message_sid"`
	UserID        int       `json:"user_id"`
	FromNumber    string    `json:"from_number"`
	ToNumber      string    `json:"to_number"`
	Body          string    `json:"body"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
	IsInbound     bool      `json:"is_inbound"`
}
