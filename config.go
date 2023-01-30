package utils

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	DEFAULT_QUERY_SIZE int
	ENV_MODE           string
	ENV_PROD           bool
)

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	DEFAULT_QUERY_SIZE = 1000
	ENV_MODE = os.Getenv("ENV_MODE")
	ENV_PROD = ENV_MODE == "production"
}
