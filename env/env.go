package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var BASE_URL string = goDotEnvVariable("BASE_URL")
var SECRET_KEY string = goDotEnvVariable("SECRET_KEY")
var CONFIG_SMTP_HOST string = goDotEnvVariable("CONFIG_SMTP_HOST")
var CONFIG_EMAIL_USERNAME string = goDotEnvVariable("CONFIG_EMAIL_USERNAME")
var CONFIG_EMAIL_PASSWORD string = goDotEnvVariable("CONFIG_EMAIL_PASSWORD")
var MAX_FILE_SIZE string = goDotEnvVariable("MAX_FILE_SIZE")
var MAX_FILE_PER_REQUEST string = goDotEnvVariable("MAX_FILE_PER_REQUEST")


func goDotEnvVariable(key string) string {
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
