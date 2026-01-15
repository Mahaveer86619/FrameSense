package config

import (
	"fmt"
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

	// Queue Configuration
	QUEUE_DRIVER string

	// SQS Configuration
	SQS_REGION     string
	SQS_ACCESS_KEY string
	SQS_SECRET_KEY string
	SQS_ENDPOINT   string
	SQS_QUEUE_NAME string

	// Kafka Configuration
	KAFKA_BROKERS string
	KAFKA_TOPIC   string
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

		QUEUE_DRIVER: getEnv("QUEUE_DRIVER", "local"),

		SQS_REGION:     getEnv("SQS_REGION", "us-east-1"),
		SQS_ACCESS_KEY: getEnv("SQS_ACCESS_KEY", ""),
		SQS_SECRET_KEY: getEnv("SQS_SECRET_KEY", ""),
		SQS_ENDPOINT:   getEnv("SQS_ENDPOINT", ""),
		SQS_QUEUE_NAME: getEnv("SQS_QUEUE_NAME", "framesense-queue"),

		KAFKA_BROKERS: getEnv("KAFKA_BROKERS", "localhost:9092"),
		KAFKA_TOPIC:   getEnv("KAFKA_TOPIC", "framesense-topic"),
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

func PrettyPrintConfig() {
	fmt.Printf("PORT: %s\n", AppConfig.PORT)
	fmt.Printf("PROFILE: %s\n", AppConfig.PROFILE)
	fmt.Printf("DB_HOST: %s\n", AppConfig.DB_HOST)
	fmt.Printf("DB_PORT: %s\n", AppConfig.DB_PORT)
	fmt.Printf("DB_USER: %s\n", AppConfig.DB_USER)
	fmt.Printf("DB_PASSWORD: %s\n", AppConfig.DB_PASSWORD)
	fmt.Printf("DB_NAME: %s\n", AppConfig.DB_NAME)
	fmt.Printf("ID_SALT: %s\n", AppConfig.ID_SALT)
	fmt.Printf("JWT_SECRET: %s\n", AppConfig.JWT_SECRET)
	fmt.Printf("STORAGE_DRIVER: %s\n", AppConfig.STORAGE_DRIVER)
	fmt.Printf("STORAGE_PATH: %s\n", AppConfig.STORAGE_PATH)
}
