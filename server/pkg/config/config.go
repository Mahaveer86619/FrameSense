package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	PORT string

	PROFILE string

	DB_HOST     string
	DB_PORT     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string

	ID_SALT string

	JWT_SECRET string

	STORAGE_DRIVER string
	STORAGE_PATH   string

	S3_ENDPOINT   string
	S3_PUBLIC_URL string
	S3_BUCKET     string
	S3_REGION     string
	S3_ACCESS_KEY string
	S3_SECRET_KEY string
}

var AppConfig Config

func LoadConfig() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Printf("Error loading .env file: %v", err)
	}

	AppConfig = Config{
		PORT:    getEnv("PORT", "7000"),
		PROFILE: getEnv("PROFILE", "dev"),

		DB_HOST:     getEnv("DB_HOST", "localhost"),
		DB_PORT:     getEnv("DB_PORT", "5432"),
		DB_USER:     getEnv("DB_USER", "admin"),
		DB_PASSWORD: getEnv("DB_PASSWORD", "password123"),
		DB_NAME:     getEnv("DB_NAME", "framesense_db"),

		ID_SALT:    getEnv("ID_SALT", "framesense-secret-salt-change-me"),
		JWT_SECRET: getEnv("JWT_SECRET", "your_secret_key"),

		STORAGE_DRIVER: getEnv("STORAGE_DRIVER", "local"),
		STORAGE_PATH:   getEnv("STORAGE_PATH", "./uploads"),

		S3_ENDPOINT:   getEnv("S3_ENDPOINT", ""),
		S3_PUBLIC_URL: getEnv("S3_PUBLIC_URL", ""),
		S3_BUCKET:     getEnv("S3_BUCKET", ""),
		S3_REGION:     getEnv("S3_REGION", "us-east-1"),
		S3_ACCESS_KEY: getEnv("S3_ACCESS_KEY", ""),
		S3_SECRET_KEY: getEnv("S3_SECRET_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := lookupEnv(key); exists {
		return value
	}
	return defaultValue
}

var lookupEnv = func(key string) (string, bool) {
	return os.LookupEnv(key)
}
