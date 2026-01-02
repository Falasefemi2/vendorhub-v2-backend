package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

func Load() {
	err := godotenv.Load()
	if err != nil {
		err = godotenv.Load("../../.env")
		if err != nil {
			fmt.Println("Warning: .env file not found, using system environment variables")
		}
	}
}

func GetDBURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		panic("DATABASE_URL environment variable is not set")
	}
	return dbURL
}
