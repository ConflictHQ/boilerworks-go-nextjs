package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	Port        int
	DatabaseURL string
	RedisURL    string
	SessionKey  string
	Environment string
	FrontendURL string
}

func Load() (*Config, error) {
	port, err := strconv.Atoi(getEnv("PORT", "8088"))
	if err != nil {
		return nil, fmt.Errorf("invalid PORT: %w", err)
	}

	dbURL := getEnv("DATABASE_URL", "postgres://boilerworks:boilerworks@localhost:5447/boilerworks?sslmode=disable")
	redisURL := getEnv("REDIS_URL", "redis://localhost:6390/0")
	sessionKey := getEnv("SESSION_KEY", "change-me-in-production-32-chars!")
	env := getEnv("ENVIRONMENT", "development")
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3004")

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
		RedisURL:    redisURL,
		SessionKey:  sessionKey,
		Environment: env,
		FrontendURL: frontendURL,
	}, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
