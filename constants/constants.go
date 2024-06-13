package constants

import "os"

var (
	AccessToken      string
	DatasetID        string
	MeasurementID    string
	APISecret        string
	PostgresHost     string
	PostgresPort     string
	PostgresUser     string
	PostgresPassword string
	PostgresDBName   string
	DavidPhoneNumber string
)

func init() {
	AccessToken = os.Getenv("FACEBOOK_ACCESS_TOKEN")
	DatasetID = os.Getenv("FACEBOOK_DATASET_ID")
	MeasurementID = os.Getenv("GOOGLE_ANALYTICS_ID")
	APISecret = os.Getenv("GOOGLE_ANALYTICS_API_KEY")
	PostgresHost = os.Getenv("POSTGRES_HOST")
	PostgresPort = os.Getenv("POSTGRES_PORT")
	PostgresUser = os.Getenv("PGUSER")
	PostgresPassword = os.Getenv("POSTGRES_PASSWORD")
	PostgresDBName = os.Getenv("POSTGRES_DB")
	DavidPhoneNumber = os.Getenv("DAVID_TWILIO_PHONE_NUMBER")
}

var TEMPLATES_DIR = "./templates/"
var WEBSITE_TEMPLATES_DIR = TEMPLATES_DIR + "website/"

var PUBLIC_DIR = "./public/"
var WEBSITE_PUBLIC_DIR = PUBLIC_DIR + "website/"

var PARTIAL_TEMPLATES_DIR = TEMPLATES_DIR + "partials/"
