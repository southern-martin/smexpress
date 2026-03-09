package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port       string
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSchema   string
}

func Load() Config {
	return Config{
		Port:       envOrDefault("PORT", "8086"),
		DBHost:     envOrDefault("DB_HOST", "localhost"),
		DBPort:     envIntOrDefault("DB_PORT", 5432),
		DBUser:     envOrDefault("DB_USER", "postgres"),
		DBPassword: envOrDefault("DB_PASSWORD", "postgres"),
		DBName:     envOrDefault("DB_NAME", "smexpress"),
		DBSchema:   envOrDefault("DB_SCHEMA", "imcs_addresses"),
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
