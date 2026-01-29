package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerHost     string
	ServerPort     string
	DatabaseURL    string
	RedisURL       string
	LogLevel       string
	LogFormat      string
	WADebug        bool
	GlobalAPIToken string

	// WhatsApp Cloud API (Meta)
	CloudAPIPhoneNumberID string
	CloudAPIAccessToken   string
}

func Load() *Config {
	_ = godotenv.Load()

	return &Config{
		ServerHost:     getEnv("SERVER_HOST", "0.0.0.0"),
		ServerPort:     getEnv("SERVER_PORT", "8080"),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://fiozap:fiozap123@localhost:5432/fiozap?sslmode=disable"),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		LogLevel:       getEnv("LOG_LEVEL", "debug"),
		LogFormat:      getEnv("LOG_FORMAT", "console"),
		WADebug:        getEnv("WA_DEBUG", "false") == "true",
		GlobalAPIToken: getEnv("GLOBAL_API_TOKEN", ""),

		CloudAPIPhoneNumberID: getEnv("CLOUD_API_PHONE_NUMBER_ID", ""),
		CloudAPIAccessToken:   getEnv("CLOUD_API_ACCESS_TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
