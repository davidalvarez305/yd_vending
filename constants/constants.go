package constants

import (
	"os"
)

var (
	FacebookAccessToken         string
	FacebookDatasetID           string
	FacebookPixelID             string
	GoogleAnalyticsID           string
	GoogleAnalyticsAPISecretKey string
	PostgresHost                string
	PostgresPort                string
	PostgresUser                string
	PostgresPassword            string
	PostgresDBName              string
	DavidPhoneNumber            string
	ServerPort                  string
	RootDomain                  string
	AWSStorageBucket            string
	CookieName                  string
	SecretAESKey                string
	AuthSecretKey               string
	EncSecretKey                string
	TwilioAccountSID            string
	TwilioAuthToken             string
	GmailUsername               string
	GmailPassword               string
	GmailEmail                  string
	CompanyName                 string
)

func Init() {
	FacebookAccessToken = os.Getenv("FACEBOOK_ACCESS_TOKEN")
	FacebookDatasetID = os.Getenv("FACEBOOK_DATASET_ID")
	FacebookPixelID = os.Getenv("FACEBOOK_PIXEL_ID")
	GoogleAnalyticsID = os.Getenv("GOOGLE_ANALYTICS_ID")
	GoogleAnalyticsAPISecretKey = os.Getenv("GOOGLE_ANALYTICS_API_KEY")
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresUser = os.Getenv("PGUSER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresDBName = os.Getenv("POSTGRES_DB")
	DavidPhoneNumber = os.Getenv("DAVID_TWILIO_PHONE_NUMBER")
	ServerPort = os.Getenv("SERVER_PORT")
	RootDomain = os.Getenv("ROOT_DOMAIN")
	AWSStorageBucket = os.Getenv("AWS_STORAGE_BUCKET")
	CookieName = os.Getenv("COOKIE_NAME")
	SecretAESKey = os.Getenv("SECRET_AES_KEY")
	AuthSecretKey = os.Getenv("AUTH_SECRET_KEY")
	EncSecretKey = os.Getenv("ENC_SECRET_KEY")
	TwilioAccountSID = os.Getenv("TWILIO_ACCOUNT_SID")
	TwilioAuthToken = os.Getenv("TWILIO_AUTH_TOKEN")
	GmailUsername = os.Getenv("GMAIL_USERNAME")
	GmailPassword = os.Getenv("GMAIL_PASSWORD")
	GmailEmail = os.Getenv("GMAIL_EMAIL")
	CompanyName = os.Getenv("COMPANY_NAME")
}

var TEMPLATES_DIR = "./templates/"
var WEBSITE_TEMPLATES_DIR = TEMPLATES_DIR + "website/"

var PUBLIC_DIR = "./public/"
var WEBSITE_PUBLIC_DIR = PUBLIC_DIR + "website/"

var PARTIAL_TEMPLATES_DIR = TEMPLATES_DIR + "partials/"
