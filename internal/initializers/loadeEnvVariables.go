package initializers

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	err := godotenv.Load()
	
	if err != nil {
		log.Fatalf("failed to load environment variables: %s/n", err.Error())
	}
}