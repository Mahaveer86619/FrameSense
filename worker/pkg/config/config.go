package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	// AWS/SQS Configuration
	AWSRegion    string
	AWSAccessKey string
	AWSSecretKey string
	SQSEndpoint  string
	SQSQueueURL  string

	// Worker Configuration
	WorkerCount int
	UseGPU      bool
	GPUType     string

	// FFmpeg Configuration
	FFmpegPreset   string
	HLSSegmentTime int
	HLSListSize    int

	// Logging
	LogLevel string
}

var AppConfig *Config

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	// Load .env file if it exists (ignore error if it doesn't)
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	config := &Config{
		// AWS/SQS Configuration
		AWSRegion:    getEnv("AWS_REGION", "us-east-1"),
		AWSAccessKey: getEnv("AWS_ACCESS_KEY_ID", "test"),
		AWSSecretKey: getEnv("AWS_SECRET_ACCESS_KEY", "test"),
		SQSEndpoint:  getEnv("SQS_ENDPOINT", ""),
		SQSQueueURL:  getEnv("SQS_QUEUE_URL", ""),

		// Worker Configuration
		WorkerCount: getEnvAsInt("WORKER_COUNT", 5),
		UseGPU:      getEnvAsBool("USE_GPU", false),
		GPUType:     getEnv("GPU_TYPE", "nvidia"),

		// FFmpeg Configuration
		FFmpegPreset:   getEnv("FFMPEG_PRESET", "fast"),
		HLSSegmentTime: getEnvAsInt("HLS_SEGMENT_TIME", 10),
		HLSListSize:    getEnvAsInt("HLS_LIST_SIZE", 0),

		// Logging
		LogLevel: getEnv("LOG_LEVEL", "info"),
	}

	// Validate required fields
	if config.SQSQueueURL == "" {
		log.Fatal("SQS_QUEUE_URL is required")
	}

	AppConfig = config
	return config
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid integer value for %s, using default: %d", key, defaultValue)
		return defaultValue
	}
	return value
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Printf("Warning: Invalid boolean value for %s, using default: %t", key, defaultValue)
		return defaultValue
	}
	return value
}

// PrintConfig logs the current configuration (excluding sensitive data)
func (c *Config) PrintConfig() {
	log.Println("=== Worker Configuration ===")
	log.Printf("AWS Region: %s", c.AWSRegion)
	log.Printf("SQS Endpoint: %s", c.SQSEndpoint)
	log.Printf("SQS Queue URL: %s", c.SQSQueueURL)
	log.Printf("Worker Count: %d", c.WorkerCount)
	log.Printf("Use GPU: %t", c.UseGPU)
	if c.UseGPU {
		log.Printf("GPU Type: %s", c.GPUType)
	}
	log.Printf("FFmpeg Preset: %s", c.FFmpegPreset)
	log.Printf("HLS Segment Time: %d", c.HLSSegmentTime)
	log.Printf("Log Level: %s", c.LogLevel)
	log.Println("===========================")
}
