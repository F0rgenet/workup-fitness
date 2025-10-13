package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Port string
	JwtSecret string
)

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file, please check your .env file")
	}

	Port = getEnv("PORT", "8080")
	JwtSecret = getEnv("JWT_SECRET", "")
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}