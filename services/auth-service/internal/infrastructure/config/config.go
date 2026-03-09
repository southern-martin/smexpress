package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port             string
	DBHost           string
	DBPort           int
	DBUser           string
	DBPassword       string
	DBName           string
	DBSchema         string
	JWTSecret        string
	JWTAccessExpiry  time.Duration
	JWTRefreshExpiry time.Duration
	NATSUrl          string
}

func Load() Config {
	return Config{
		Port:             envOrDefault("PORT", "8081"),
		DBHost:           envOrDefault("DB_HOST", "localhost"),
		DBPort:           envIntOrDefault("DB_PORT", 5432),
		DBUser:           envOrDefault("DB_USER", "postgres"),
		DBPassword:       envOrDefault("DB_PASSWORD", "postgres"),
		DBName:           envOrDefault("DB_NAME", "smexpress"),
		DBSchema:         envOrDefault("DB_SCHEMA", "imcs_auth"),
		JWTSecret:        envOrDefault("JWT_SECRET", "dev-secret-change-in-production"),
		JWTAccessExpiry:  envDurationOrDefault("JWT_ACCESS_EXPIRY", 15*time.Minute),
		JWTRefreshExpiry: envDurationOrDefault("JWT_REFRESH_EXPIRY", 24*time.Hour),
		NATSUrl:          envOrDefault("NATS_URL", "nats://localhost:4222"),
	}
}

func envOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func envIntOrDefault(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

func envDurationOrDefault(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	d, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return d
}
