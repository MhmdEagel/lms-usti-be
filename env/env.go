package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

var DEFAULT_PORT string = goDotEnvVariable("DEFAULT_PORT")
var SECRET_KEY string = goDotEnvVariable("SECRET_KEY")
var CONFIG_SMTP_HOST string = goDotEnvVariable("CONFIG_SMTP_HOST")
var CONFIG_EMAIL_USERNAME string = goDotEnvVariable("CONFIG_EMAIL_USERNAME")
var CONFIG_EMAIL_PASSWORD string = goDotEnvVariable("CONFIG_EMAIL_PASSWORD")


func goDotEnvVariable(key string) string {
	err := godotenv.Load("./.env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}
