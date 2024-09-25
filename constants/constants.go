package constants

import (
	"os"
)

var (
	FacebookAccessToken          string
	FacebookDatasetID            string
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
	SessionLength                int64
	StaticPath                   string
	MediaPath                    string
	MaxOpenConnections           string
	MaxIdleConnections           string
	MaxConnectionLifetime        string
	CompanyEmail                 string
	GoogleWebhookKey             string
	TokenExpiryLength            int64
)

func Init() {
	FacebookAccessToken = os.Getenv("FACEBOOK_ACCESS_TOKEN")
	FacebookDatasetID = os.Getenv("FACEBOOK_DATASET_ID")
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
	SessionLength = 604800
	TokenExpiryLength = 604800
	StaticPath = os.Getenv("STATIC_PATH")
	MediaPath = os.Getenv("MEDIA_PATH")
	MaxOpenConnections = os.Getenv("MAX_OPEN_CONNECTIONS")
	MaxIdleConnections = os.Getenv("MAX_IDLE_CONNECTIONS")
	MaxConnectionLifetime = os.Getenv("MAX_CONN_LIFETIME")
	DomainHost = os.Getenv("DOMAIN_HOST")
	CompanyEmail = os.Getenv("COMPANY_EMAIL")
	GoogleAdsID = os.Getenv("GOOGLE_ADS_ID")
	GoogleAdsCallConversionLabel = os.Getenv("GOOGLE_ADS_CALL_CONVERSION_LABEL")
}

var TEMPLATES_DIR = "./templates/"
var WEBSITE_TEMPLATES_DIR = TEMPLATES_DIR + "website/"
var CRM_TEMPLATES_DIR = TEMPLATES_DIR + "crm/"
var PARTIAL_TEMPLATES_DIR = TEMPLATES_DIR + "partials/"
