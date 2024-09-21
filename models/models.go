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

type LeadNote struct {
	LeadNoteID    int    `json:"lead_note_id" form:"lead_note_id" schema:"lead_note_id"`
	Note          string `json:"note" form:"note" schema:"note"`
	LeadID        int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	DateAdded     int64  `json:"date_added" form:"date_added" schema:"date_added"`
	AddedByUserID int    `json:"added_by_user_id" form:"added_by_user_id" schema:"added_by_user_id"`
}

type LeadImage struct {
	LeadImageID   int    `json:"lead_image_id" form:"lead_image_id" schema:"lead_image_id"`
	Src           string `json:"src" form:"src" schema:"src"`
	LeadID        int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	DateAdded     int64  `json:"date_added" form:"date_added" schema:"date_added"`
	AddedByUserID int    `json:"added_by_user_id" form:"added_by_user_id" schema:"added_by_user_id"`
}

type LeadMarketing struct {
	LeadMarketingID  int64  `json:"lead_marketing_id"`
	LeadID           int64  `json:"lead_id"`
	Source           string `json:"source"`
	Medium           string `json:"medium"`
	Channel          string `json:"channel"`
	LandingPage      string `json:"landing_page"`
	Longitude        string `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude         string `json:"latitude" form:"latitude" schema:"latitude"`
	Keyword          string `json:"keyword"`
	Referrer         string `json:"referrer"`
	GCLID            string `json:"gclid"`
	CampaignID       int64  `json:"campaign_id"`
	AdCampaign       string `json:"ad_campaign"`
	AdGroupID        int64  `json:"ad_group_id"`
	AdGroupName      string `json:"ad_group_name"`
	AdSetID          int64  `json:"ad_set_id"`
	AdSetName        string `json:"ad_set_name"`
	AdID             int64  `json:"ad_id"`
	AdHeadline       int64  `json:"ad_headline"`
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
	SessionID   int    `json:"session_id" form:"session_id" schema:"session_id"`
	UserID      int    `json:"user_id" form:"user_id" schema:"user_id"`
	CSRFSecret  string `json:"csrf_secret" form:"csrf_secret" schema:"csrf_secret"`
	ExternalID  string `json:"external_id" form:"external_id" schema:"external_id"`
	DateCreated int64  `json:"date_created" form:"date_created" schema:"date_created"`
	DateExpires int64  `json:"date_expires" form:"date_expires" schema:"date_expires"`
}

type Location struct {
	LocationID          int    `json:"location_id" form:"location_id" schema:"location_id"`
	VendingLocationID   int    `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	BusinessID          int    `json:"business_id" form:"business_id" schema:"business_id"`
	CityID              int    `json:"city_id" form:"city_id" schema:"city_id"`
	DateStarted         int64  `json:"date_started" form:"date_started" schema:"date_started"`
	Name                string `json:"name" form:"name" schema:"name"`
	Longitude           string `json:"longitude" form:"longitude" schema:"longitude"`
	Latitude            string `json:"latitude" form:"latitude" schema:"latitude"`
	StreetAdressLineOne string `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAdressLineTwo string `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	ZipCode             string `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State               string `json:"state" form:"state" schema:"state"`
	Opening             string `json:"opening" form:"opening" schema:"opening"`
	Closing             string `json:"closing" form:"closing" schema:"closing"`
}

type Business struct {
	BusinessID            int    `json:"business_id" form:"business_id" schema:"business_id"`
	Name                  string `json:"name" form:"name" schema:"name"`
	IsActive              bool   `json:"is_active" form:"is_active" schema:"is_active"`
	DateCreated           int64  `json:"date_created" form:"date_created" schema:"date_created"`
	Website               string `json:"website" form:"website" schema:"website"`
	Industry              string `json:"industry" form:"industry" schema:"industry"`
	GoogleBusinessProfile string `json:"google_business_profile" form:"google_business_profile" schema:"google_business_profile"`
}

type BusinessContact struct {
	BusinessContactID      int    `json:"business_contact_id" form:"business_contact_id" schema:"business_contact_id"`
	FirstName              string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName               string `json:"last_name" form:"last_name" schema:"last_name"`
	Phone                  string `json:"phone" form:"phone" schema:"phone"`
	Email                  string `json:"email" form:"email" schema:"email"`
	PreferredContactMethod string `json:"preferred_contact_method" form:"preferred_contact_method" schema:"preferred_contact_method"`
	PreferredContactTime   string `json:"preferred_contact_time" form:"preferred_contact_time" schema:"preferred_contact_time"`
	BusinessID             int    `json:"business_id" form:"business_id" schema:"business_id"`
	BusinessPosition       string `json:"business_position" form:"business_position" schema:"business_position"`
	IsPrimaryContact       bool   `json:"is_primary_contact" form:"is_primary_contact" schema:"is_primary_contact"`
}

type MachineStatus struct {
	MachineStatusID int    `json:"machine_status_id" form:"machine_status_id" schema:"machine_status_id"`
	Status          string `json:"status" form:"status" schema:"status"`
}

type Machine struct {
	MachineID              int     `json:"machine_id" form:"machine_id" schema:"machine_id"`
	VendingTypeID          int     `json:"vending_type_id" form:"vending_type_id" schema:"vending_type_id"`
	MachineStatusID        int     `json:"machine_status_id" form:"machine_status_id" schema:"machine_status_id"`
	LocationID             int     `json:"location_id" form:"location_id" schema:"location_id"`
	VendorID               int     `json:"vendor_id" form:"vendor_id" schema:"vendor_id"`
	Year                   int     `json:"year" form:"year" schema:"year"`
	Make                   string  `json:"make" form:"make" schema:"make"`
	Model                  string  `json:"model" form:"model" schema:"model"`
	PurchasePrice          float64 `json:"purchase_price" form:"purchase_price" schema:"purchase_price"`
	CardReaderSerialNumber string  `json:"card_reader_serial_number" form:"card_reader_serial_number" schema:"card_reader_serial_number"`
	ColumnsQty             int     `json:"columns_qty" form:"columns_qty" schema:"columns_qty"`
	RowsQty                int     `json:"rows_qty" form:"rows_qty" schema:"rows_qty"`
	TotalSlots             int     `json:"total_slots" form:"total_slots" schema:"total_slots"`
	PurchaseDate           int64   `json:"purchase_date" form:"purchase_date" schema:"purchase_date"`
}

type Vendor struct {
	VendorID               int    `json:"vendor_id" form:"vendor_id" schema:"vendor_id"`
	Name                   string `json:"name" form:"name" schema:"name"`
	FirstName              string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName               string `json:"last_name" form:"last_name" schema:"last_name"`
	Phone                  string `json:"phone" form:"phone" schema:"phone"`
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

type Supplier struct {
	SupplierID            int     `json:"supplier_id" form:"supplier_id" schema:"supplier_id"`
	Name                  string  `json:"name" form:"name" schema:"name"`
	MembershipID          string  `json:"membership_id" form:"membership_id" schema:"membership_id"`
	MembershipCost        float64 `json:"membership_cost" form:"membership_cost" schema:"membership_cost"`          // Use float64 for MONEY
	MembershipRenewal     int64   `json:"membership_renewal" form:"membership_renewal" schema:"membership_renewal"` // Represented as Unix timestamp
	StreetAddressLineOne  string  `json:"street_address_line_one" form:"street_address_line_one" schema:"street_address_line_one"`
	StreetAddressLineTwo  string  `json:"street_address_line_two" form:"street_address_line_two" schema:"street_address_line_two"`
	CityID                int     `json:"city_id" form:"city_id" schema:"city_id"`
	ZipCode               string  `json:"zip_code" form:"zip_code" schema:"zip_code"`
	State                 string  `json:"state" form:"state" schema:"state"`
	GoogleBusinessProfile string  `json:"google_business_profile,omitempty" form:"google_business_profile" schema:"google_business_profile"`
}

type TicketType struct {
	TicketTypeID int    `json:"ticket_type_id" form:"ticket_type_id" schema:"ticket_type_id"`
	Type         string `json:"type" form:"type" schema:"type"`
	UrgencyLevel int    `json:"urgency_level" form:"urgency_level" schema:"urgency_level"`
	Description  string `json:"description" form:"description" schema:"description"`
}

type TicketStatus struct {
	TicketStatusID int    `json:"ticket_status_id" form:"ticket_status_id" schema:"ticket_status_id"`
	Status         string `json:"status" form:"status" schema:"status"`
	Description    string `json:"description" form:"description" schema:"description"`
}

type Ticket struct {
	TicketID       int    `json:"ticket_id" form:"ticket_id" schema:"ticket_id"`
	MachineID      int    `json:"machine_id" form:"machine_id" schema:"machine_id"`
	TicketTypeID   int    `json:"ticket_type_id" form:"ticket_type_id" schema:"ticket_type_id"`
	Content        string `json:"content" form:"content" schema:"content"`
	CreatedAt      int64  `json:"created_at" form:"created_at" schema:"created_at"`
	UpdatedAt      int64  `json:"updated_at" form:"updated_at" schema:"updated_at"`
	TicketStatusID int    `json:"ticket_status_id" form:"ticket_status_id" schema:"ticket_status_id"`
	AssignedTo     int    `json:"assigned_to" form:"assigned_to" schema:"assigned_to"`
	Priority       int    `json:"priority" form:"priority" schema:"priority"`
	Summary        string `json:"summary" form:"summary" schema:"summary"`
}

type TicketImage struct {
	TicketImageID int    `json:"ticket_image_id" form:"ticket_image_id" schema:"ticket_image_id"`
	TicketID      int    `json:"ticket_id" form:"ticket_id" schema:"ticket_id"`
	URL           string `json:"url" form:"url" schema:"url"`
	Caption       string `json:"caption" form:"caption" schema:"caption"`
}
