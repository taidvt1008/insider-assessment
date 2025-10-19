package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost       string
	DBPort       string
	DBUser       string
	DBPassword   string
	DBName       string
	RedisHost    string
	WebhookURL   string
	SendInterval time.Duration
	ServerPort   string
}

func Load() *Config {
	_ = godotenv.Load()

	interval, err := time.ParseDuration(getEnv("SEND_INTERVAL", false, "2m"))
	if err != nil {
		log.Fatalf("Invalid SEND_INTERVAL: %v", err)
	}

	return &Config{
		DBHost:       getEnv("DB_HOST", true, ""),
		DBPort:       getEnv("DB_PORT", false, "5432"),
		DBUser:       getEnv("DB_USER", true, ""),
		DBPassword:   getEnv("DB_PASSWORD", true, ""),
		DBName:       getEnv("DB_NAME", true, ""),
		RedisHost:    getEnv("REDIS_ADDR", true, ""),
		WebhookURL:   getEnv("WEBHOOK_URL", true, ""),
		SendInterval: interval,
		ServerPort:   getEnv("SERVER_PORT", false, "8080"),
	}
}

func getEnv(key string, required bool, fallback string) string {
	v := os.Getenv(key)
	if v != "" {
		return v
	}
	if !required {
		return fallback
	}
	panic(fmt.Sprintf("%s is required", key))
}
