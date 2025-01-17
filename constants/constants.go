package constants

import (
	"os"
)

const (
	UserAdminRoleID        int = 1
	CommissionReportRoleID int = 2

	DavidUserID int = 1

	EmailMIMEBoundary          string = "my-boundary-12345"
	TimeZone                   string = "America/New_York"
	CommissionReportFilename   string = "commission_report.xlsx"
	StructSpreadsheetHeaderTag string = "spreadsheet_header"

	LeadApplicationLeadTypeID int = 2
	VendingLeadTypeID         int = 1

	LeadApplicationEventName   string = "lead_application"
	LeadAppointmentEventName   string = "lead_appointment"
	LeadOfferEventName         string = "lead_offer"
	LeadOfferAcceptedEventName string = "lead_offer_accepted"
	OptInEventName             string = "opt_in"
	QuoteEventName             string = "quote"
	LeadEventName              string = "Lead"

	MiniSiteBuildCommand               string = "next build"
	MiniSiteDevCommand                 string = "npm run dev"
	MiniSiteBranchName                 string = "main"
	MiniSiteFramework                  string = "nextjs"
	MiniSiteOutputDirectory            string = ".next"
	MiniSiteEnvironmentVariablesTag    string = "env"
	MiniSiteEnvironmentVariablesPrefix string = "NEXT_PUBLIC_"

	LeadOfferSentID     int = 1
	LeadOfferAcceptedID int = 2
	LeadOfferRejectedID int = 3
	LeadOfferExpiredID  int = 4

	VercelProjectEnvinronmentVariableTarget string = "production"
	VercelProjectEnvinronmentVariableType   string = "plain"
	VercelUsername                          string = "ydvending"

	LeadOfferAcceptedSuccessPath string = "/funnel/90-day-channel-offer-accepted"
	LeadOfferCanceledPath        string = "/funnel/90-day-challenge-offer-canceled"

	LinkedInEventID   string = "19410058"
	MicrosoftAdsTagID string = "187165719"
)

var (
	FacebookAccessToken          string
	FacebookDatasetID            string
	LinkedInTagID                string
	GoogleAnalyticsID            string
	GoogleAdsID                  string
	GoogleAdsCallConversionLabel string
	GoogleAnalyticsAPISecretKey  string
	GoogleRefreshToken           string
	GoogleJSONPath               string
	PostgresHost                 string
	PostgresPort                 string
	PostgresUser                 string
	PostgresPassword             string
	PostgresDBName               string
	DavidPhoneNumber             string
	DavidEmail                   string
	YovaEmail                    string
	YovaPhoneNumber              string
	ServerPort                   string
	RootDomain                   string
	AWSStorageBucket             string
	AWSS3BucketName              string
	AWSS3LiveImagesPath          string
	AWSS3MarketingImagesPath     string
	AWSRegion                    string
	CookieName                   string
	DomainHost                   string
	SecretAESKey                 string
	AuthSecretKey                string
	EncSecretKey                 string
	TwilioAccountSID             string
	TwilioAuthToken              string
	TwilioCallbackWebhook        string
	CompanyName                  string
	SiteName                     string
	SessionName                  string
	LeadsPerPage                 int
	CompanyPhoneNumber           string
	SessionLength                int
	CSRFTokenLength              int
	StaticPath                   string
	MediaPath                    string
	MaxOpenConnections           string
	MaxIdleConnections           string
	MaxConnectionLifetime        string
	CompanyEmail                 string
	GoogleWebhookKey             string
	VercelAccessToken            string
	MiniSiteGithubRepo           string
	MiniSiteVercelTeamID         string
	MiniSiteRepoID               string
	GithubOwnerName              string
	GithubAccessToken            string
	StrikeAPIKey                 string
)

func Init() {
	FacebookAccessToken = os.Getenv("FACEBOOK_ACCESS_TOKEN")
	FacebookDatasetID = os.Getenv("FACEBOOK_DATASET_ID")
	LinkedInTagID = os.Getenv("LINKEDIN_TAG_ID")
	GoogleAnalyticsID = os.Getenv("GOOGLE_ANALYTICS_ID")
	GoogleAnalyticsAPISecretKey = os.Getenv("GOOGLE_ANALYTICS_API_KEY")
	GoogleWebhookKey = os.Getenv("GOOGLE_WEBHOOK_KEY")
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresUser = os.Getenv("PGUSER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresDBName = os.Getenv("POSTGRES_DB")
	DavidPhoneNumber = os.Getenv("DAVID_PHONE_NUMBER")
	YovaPhoneNumber = os.Getenv("YOVA_PHONE_NUMBER")
	ServerPort = os.Getenv("SERVER_PORT")
	RootDomain = os.Getenv("ROOT_DOMAIN")
	AWSStorageBucket = os.Getenv("AWS_STORAGE_BUCKET")
	AWSS3BucketName = os.Getenv("AWS_S3_BUCKET_NAME")
	AWSS3LiveImagesPath = os.Getenv("AWS_S3_IMAGES_PATH")
	AWSS3MarketingImagesPath = os.Getenv("AWS_S3_MARKETING_IMAGES_PATH")
	AWSRegion = os.Getenv("AWS_REGION")
	CookieName = os.Getenv("COOKIE_NAME")
	SecretAESKey = os.Getenv("SECRET_AES_KEY")
	AuthSecretKey = os.Getenv("AUTH_SECRET_KEY")
	EncSecretKey = os.Getenv("ENC_SECRET_KEY")
	TwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	CompanyName = os.Getenv("COMPANY_NAME")
	SiteName = os.Getenv("SITE_NAME")
	GoogleRefreshToken = os.Getenv("GOOGLE_REFRESH_TOKEN")
	GoogleJSONPath = "./google.json"
	DavidEmail = os.Getenv("DAVID_EMAIL")
	YovaEmail = os.Getenv("YOVA_EMAIL")
	SessionName = "yd_vending_sessions"
	LeadsPerPage = 10
	TwilioCallbackWebhook = "/call/inbound/end"
	CompanyPhoneNumber = os.Getenv("COMPANY_PHONE_NUMBER")
	SessionLength = 7
	CSRFTokenLength = 1
	StaticPath = os.Getenv("STATIC_PATH")
	MediaPath = os.Getenv("MEDIA_PATH")
	MaxOpenConnections = os.Getenv("MAX_OPEN_CONNECTIONS")
	MaxIdleConnections = os.Getenv("MAX_IDLE_CONNECTIONS")
	MaxConnectionLifetime = os.Getenv("MAX_CONN_LIFETIME")
	DomainHost = os.Getenv("DOMAIN_HOST")
	CompanyEmail = os.Getenv("COMPANY_EMAIL")
	GoogleAdsID = os.Getenv("GOOGLE_ADS_ID")
	GoogleAdsCallConversionLabel = os.Getenv("GOOGLE_ADS_CALL_CONVERSION_LABEL")
	VercelAccessToken = os.Getenv("VERCEL_ACCESS_TOKEN")
	MiniSiteGithubRepo = os.Getenv("MINI_SITE_GITHUB_REPO")
	MiniSiteVercelTeamID = os.Getenv("MINI_SITE_GITHUB_TEAM_ID")
	MiniSiteRepoID = os.Getenv("MINI_SITE_REPO_ID")
	GithubOwnerName = os.Getenv("GITHUB_OWNER_NAME")
	GithubAccessToken = os.Getenv("GITHUB_ACCESS_TOKEN")
	StrikeAPIKey = os.Getenv("STRIPE_API_KEY")
}

var TEMPLATES_DIR = "./templates/"
var LOCAL_FILES_DIR = "./local_files/"
var WEBSITE_TEMPLATES_DIR = TEMPLATES_DIR + "website/"
var CRM_TEMPLATES_DIR = TEMPLATES_DIR + "crm/"
var INVENTORY_TEMPLATES_DIR = TEMPLATES_DIR + "inventory/"
var PARTIAL_TEMPLATES_DIR = TEMPLATES_DIR + "partials/"
var EXTERNAL_REPORTS_TEMPLATES_DIR = TEMPLATES_DIR + "external/"
var FUNNEL_TEMPLATES_DIR = TEMPLATES_DIR + "funnel/"
var EMAIL_ATTACHMENTS_S3_BUCKET = "email-attachmennts/"
var SQL_FILES_S3_BUCKET = "sql/"
