package models

import (
	"time"
)

type VendingType struct {
	VendingTypeID int    `json:"vending_type_id" form:"vending_type_id" schema:"vending_type_id"`
	MachineType   string `json:"machine_type" form:"machine_type" schema:"machine_type"`
}

type VendingLocation struct {
	VendingLocationID int    `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	LocationType      string `json:"location_type" form:"location_type" schema:"location_type"`
}

type City struct {
	CityID int    `json:"city_id" form:"city_id" schema:"city_id"`
	Name   string `json:"name" form:"name" schema:"name"`
}

type User struct {
	UserID      int    `json:"user_id" form:"user_id" schema:"user_id"`
	Username    string `json:"username" form:"username" schema:"username"`
	PhoneNumber string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Password    string `json:"password" form:"password" schema:"password"`
	UserRoleID  int    `json:"user_role_id" form:"user_role_id" schema:"user_role_id"`
	FirstName   string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName    string `json:"last_name" form:"last_name" schema:"last_name"`
}

type LeadType struct {
	LeadTypeID int    `json:"lead_type_id" form:"lead_type_id" schema:"lead_type_id"`
	LeadType   string `json:"lead_type" form:"lead_type" schema:"lead_type"`
}

type Lead struct {
	LeadID             int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	LeadTypeID         int    `json:"lead_type_id" form:"lead_type_id" schema:"lead_type_id"`
	FirstName          string `json:"first_name" form:"first_name" schema:"first_name"`
	LastName           string `json:"last_name" form:"last_name" schema:"last_name"`
	PhoneNumber        string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email              string `json:"email" form:"email" schema:"email"`
	CreatedAt          int64  `json:"created_at" form:"created_at" schema:"created_at"`
	VendingTypeID      int    `json:"vending_type_id" form:"vending_type_id" schema:"vending_type_id"`
	VendingLocationID  int    `json:"vending_location_id" form:"vending_location_id" schema:"vending_location_id"`
	LeadStatusID       int    `json:"lead_status_id" form:"lead_status_id" schema:"lead_status_id"`
	Message            string `json:"message" form:"message" schema:"message"`
	OptInTextMessaging bool   `json:"opt_in_text_messaging" form:"opt_in_text_messaging" schema:"opt_in_text_messaging"`
}

type LeadApplication struct {
	LeadApplicationID int    `json:"lead_application_id" form:"lead_application_id" schema:"lead_application_id"`
	LeadID            int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	Website           string `json:"website" form:"website" schema:"website"`
	CompanyName       string `json:"company_name" form:"company_name" schema:"company_name"`
	YearsInBusiness   int    `json:"years_in_business" form:"years_in_business" schema:"years_in_business"`
	NumLocations      int    `json:"num_locations" form:"num_locations" schema:"num_locations"`
	City              string `json:"city" form:"city" schema:"city"`
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
	ClickID          string `json:"click_id"`
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

type LocationStatus struct {
	LocationStatusID int    `json:"location_status_id" form:"location_status_id" schema:"location_status_id"`
	Status           string `json:"status" form:"status" schema:"status"`
}

type Business struct {
	BusinessID            int    `json:"business_id" form:"business_id" schema:"business_id"`
	Name                  string `json:"name" form:"name" schema:"name"`
	IsActive              bool   `json:"is_active" form:"is_active" schema:"is_active"`
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
	BusinessPosition       string `json:"business_position" form:"business_position" schema:"business_position"`
}

type BusinessLocationContact struct {
	BusinessContactID int  `json:"business_contact_id" form:"business_contact_id" schema:"business_contact_id"`
	LocationID        int  `json:"location_id" form:"location_id" schema:"location_id"`
	BusinessID        int  `json:"business_id" form:"business_id" schema:"business_id"`
	IsPrimaryContact  bool `json:"is_primary_contact" form:"is_primary_contact" schema:"is_primary_contact"`
}

type MachineStatus struct {
	MachineStatusID int    `json:"machine_status_id" form:"machine_status_id" schema:"machine_status_id"`
	Status          string `json:"status" form:"status" schema:"status"`
}

type Machine struct {
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

type Vendor struct {
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

type Supplier struct {
	SupplierID            int     `json:"supplier_id" form:"supplier_id" schema:"supplier_id"`
	Name                  string  `json:"name" form:"name" schema:"name"`
	MembershipID          string  `json:"membership_id" form:"membership_id" schema:"membership_id"`
	MembershipCost        float64 `json:"membership_cost" form:"membership_cost" schema:"membership_cost"` // Use float64 for MONEY
	MembershipRenewal     string  `json:"membership_renewal" form:"membership_renewal" schema:"membership_renewal"`
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

type Image struct {
	ImageID       int    `json:"image_id" form:"image_id" schema:"image_id"`
	Src           string `json:"src" form:"src" schema:"src"`
	DateAdded     int64  `json:"date_added" form:"date_added" schema:"date_added"`
	AddedByUserID int    `json:"added_by_user_id" form:"added_by_user_id" schema:"added_by_user_id"`
}

type SeedLiveTransaction struct {
	SeedLiveTransactionID  int       `json:"seed_live_transaction_id" form:"seed_live_transaction_id" schema:"seed_live_transaction_id"`
	TerminalNumber         string    `json:"terminal_number" form:"terminal_number" schema:"terminal_number"`
	TransactionRefNumber   string    `json:"transaction_ref_number" form:"transaction_ref_number" schema:"transaction_ref_number"`
	TransactionType        string    `json:"transaction_type" form:"transaction_type" schema:"transaction_type"`
	CardNumber             string    `json:"card_number" form:"card_number" schema:"card_number"`
	TotalAmount            float64   `json:"total_amount" form:"total_amount" schema:"total_amount"`
	VendedColumns          int       `json:"vended_columns" form:"vended_columns" schema:"vended_columns"`
	Price                  float64   `json:"price" form:"price" schema:"price"`
	MDBNumber              int       `json:"mdb_number" form:"mdb_number" schema:"mdb_number"`
	NumberOfProductsVended int       `json:"number_of_products_vended" form:"number_of_products_vended" schema:"number_of_products_vended"`
	Timestamp              time.Time `json:"timestamp" form:"timestamp" schema:"timestamp"`
	CardId                 string    `json:"card_id" form:"card_id" schema:"card_id"`
}

type Product struct {
	ProductID         int     `json:"product_id" form:"product_id" schema:"product_id"`
	Name              string  `json:"name" form:"name" schema:"name"`
	ProductCategoryID int     `json:"product_category_id" form:"product_category_id" schema:"product_category_id"`
	Size              float64 `json:"size" form:"size" schema:"size"`
	SizeType          string  `json:"size_type" form:"size_type" schema:"size_type"`
	UPC               string  `json:"upc" form:"upc" schema:"upc"`
}

type ProductCategory struct {
	ProductCategoryID int    `json:"product_category_id" form:"product_category_id"`
	Name              string `json:"name" form:"name" schema:"name"`
}

type Slot struct {
	SlotID      int     `json:"slot_id" form:"slot_id" schema:"slot_id"`
	Nickname    string  `json:"nickname" form:"nickname" schema:"nickname"`
	Slot        string  `json:"slot" form:"slot" schema:"slot"`
	MachineCode string  `json:"machine_code" form:"machine_code" schema:"machine_code"`
	MachineID   int     `json:"machine_id" form:"machine_id" schema:"machine_id"`
	Price       float64 `json:"price" form:"price" schema:"price"`
	Capacity    int     `json:"capacity" form:"capacity" schema:"capacity"`
}

type ProductSlotAssignment struct {
	ProductSlotAssignmentID int     `json:"product_slot_assignment_id" form:"product_slot_assignment_id" schema:"product_slot_assignment_id"`
	SlotID                  int     `json:"slot_id" form:"slot_id" schema:"slot_id"`
	ProductID               int     `json:"product_id" form:"product_id" schema:"product_id"`
	SupplierID              int     `json:"supplier_id" form:"supplier_id" schema:"supplier_id"`
	ExpirationDate          int64   `json:"expiration_date" form:"expiration_date" schema:"expiration_date"`
	UnitCost                float64 `json:"unit_cost" form:"unit_cost" schema:"unit_cost"`
	Quantity                int     `json:"quantity" form:"quantity" schema:"quantity"`
	DateAssigned            int64   `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
}

type Refill struct {
	RefillID     int   `json:"refill_id" form:"refill_id" schema:"refill_id"`
	SlotID       int   `json:"slot_id" form:"slot_id" schema:"slot_id"`
	DateRefilled int64 `json:"date_refilled" form:"date_refilled" schema:"date_refilled"`
}

type MachineLocationAssignment struct {
	MachineLocationAssignmentID int   `json:"machine_location_assignment_id" form:"machine_location_assignment_id" schema:"machine_location_assignment_id"`
	LocationID                  int   `json:"location_id" form:"location_id" schema:"location_id"`
	MachineID                   int   `json:"machine_id" form:"machine_id" schema:"machine_id"`
	DateAssigned                int64 `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
	IsActive                    bool  `json:"is_active" form:"is_active" schema:"is_active"`
}

type MachineCardReaderAssignment struct {
	MachineCardReaderID    int    `json:"machine_card_reader_assignment_id" form:"machine_card_reader_assignment_id" schema:"machine_card_reader_assignment_id"`
	CardReaderSerialNumber string `json:"card_reader_serial_number" form:"card_reader_serial_number" schema:"card_reader_serial_number"`
	MachineID              int    `json:"machine_id" form:"machine_id" schema:"machine_id"`
	DateAssigned           int64  `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
	IsActive               bool   `json:"is_active" form:"is_active" schema:"is_active"`
}

type SlotPriceLog struct {
	SlotPriceLogID int     `json:"slot_price_log_id" form:"slot_price_log_id" schema:"slot_price_log_id"`
	SlotID         int     `json:"slot_id" form:"slot_id" schema:"slot_id"`
	Price          float64 `json:"price" form:"price" schema:"price"`
	DateAssigned   int64   `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
}

type SeedTransaction struct {
	TransactionLogID     int64  `json:"transaction_log_id" db:"transaction_log_id" form:"transaction_log_id" schema:"transaction_log_id"`
	TransactionTimestamp int64  `json:"transaction_timestamp" db:"transaction_timestamp" form:"transaction_timestamp" schema:"transaction_timestamp"`
	Device               string `json:"device" db:"device" form:"device" schema:"device"`
	Item                 string `json:"item" db:"item" form:"item" schema:"item"`
	TransactionType      string `json:"transaction_type" db:"transaction_type" form:"transaction_type" schema:"transaction_type"`
	CardID               string `json:"card_id" db:"card_id" form:"card_id" schema:"card_id"`
	CardNumber           string `json:"card_number" db:"card_number" form:"card_number" schema:"card_number"`
	NumTransactions      int    `json:"num_transactions" db:"num_transactions" form:"num_transactions" schema:"num_transactions"`
	Items                int    `json:"items" db:"items" form:"items" schema:"items"`
}

type TransactionValidation struct {
	TransactionValidationID int64 `json:"transaction_validation_id" db:"transaction_validation_id"`
	TransactionID           int64 `json:"transaction_id" db:"transaction_id" form:"transaction_id" schema:"transaction_id"`
	IsValidated             bool  `json:"is_validated" db:"is_validated" form:"is_validated" schema:"is_validated"`
}

type LocationCommission struct {
	LocationID   int     `json:"location_id" form:"location_id" schema:"location_id"`
	Commission   float64 `json:"commission" form:"commission" schema:"commission"`
	DateAssigned int64   `json:"date_assigned" form:"date_assigned" schema:"date_assigned"`
}

type UserRole struct {
	RoleID int    `json:"role_id" form:"role_id" schema:"role_id"`
	Role   string `json:"role" form:"role" schema:"role"`
}

type UserExternalReportsRole struct {
	UserID     string `json:"user_id" form:"user_id" schema:"user_id"`
	BusinessID string `json:"business_id" form:"business_id" schema:"business_id"`
}

type EmailSchedule struct {
	EmailScheduleID int    `json:"email_schedule_id" form:"email_schedule_id" schema:"email_schedule_id"`
	EmailName       string `json:"email_name" form:"email_name" schema:"email_name"`
	IntervalSeconds int64  `json:"interval_seconds" form:"interval_seconds" schema:"interval_seconds"`
	Recipients      string `json:"recipients" form:"recipients" schema:"recipients"`
	Subject         string `json:"subject" form:"subject" schema:"subject"`
	Body            string `json:"body" form:"body" schema:"body"`
	Sender          string `json:"sender" form:"sender" schema:"sender"`
	SQLFile         string `json:"sql_file" form:"sql_file" schema:"sql_file"`
	LastSent        int64  `json:"last_sent" form:"last_sent" schema:"last_sent"`
	IsActive        bool   `json:"is_active" form:"is_active" schema:"is_active"`
}

type SentEmail struct {
	SentEmailID     int    `json:"sent_email_id" form:"sent_email_id" schema:"sent_email_id"`
	EmailScheduleID int    `json:"email_schedule_id" form:"email_schedule_id" schema:"email_schedule_id"`
	DeliveryStatus  string `json:"delivery_status" form:"delivery_status" schema:"delivery_status"`
	DateSent        int64  `json:"date_sent" form:"date_sent" schema:"date_sent"`
	ErrorMessage    string `json:"error_message" form:"error_message" schema:"error_message"`
}

type LeadAppointment struct {
	LeadAppointmentID int    `json:"lead_appointment_id" form:"lead_appointment_id" schema:"lead_appointment_id"`
	LeadID            int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	AppointmentTime   int64  `json:"appointment_time" form:"appointment_time" schema:"appointment_time"`
	DateCreated       int64  `json:"date_created" form:"date_created" schema:"date_created"`
	Link              string `json:"link" form:"link" schema:"link"`
	Attendee          string `json:"attendee" form:"attendee" schema:"attendee"`
}

type MiniSite struct {
	MiniSiteID      int    `json:"mini_site_id" form:"mini_site_id" schema:"mini_site_id"`
	LeadID          int    `json:"lead_id" form:"lead_id" schema:"lead_id"`
	Website         string `json:"website" form:"website" schema:"website"`
	DateCreated     int64  `json:"date_created" form:"date_created" schema:"date_created"`
	PhoneNumber     string `json:"phone_number" form:"phone_number" schema:"phone_number"`
	Email           string `json:"email" form:"email" schema:"email"`
	VercelProjectID string `json:"vercel_project_id" form:"vercel_project_id" schema:"vercel_project_id"`
}

type MiniSiteEnvironmentVariable struct {
	MiniSiteEnvironmentVariableID int    `json:"mini_site_environment_variable_id" form:"mini_site_environment_variable_id" schema:"mini_site_environment_variable_id"`
	MiniSiteID                    int    `json:"mini_site_id" form:"mini_site_id" schema:"mini_site_id"`
	EnvironmentVariableUniqueID   string `json:"environment_variable_unique_id" form:"environment_variable_unique_id" schema:"environment_variable_unique_id"`
	Key                           string `json:"key" form:"key" schema:"key"`
	Value                         string `json:"value" form:"value" schema:"value"`
}
