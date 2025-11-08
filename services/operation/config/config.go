package config

import (
	"fmt"
	"os"
)

// Config
// Includes some of the server configs.
type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	GRPCPort   string
}

// LoadConfig
// Gets the required configs for the service.
//
// Returns:
//   - *Config
//
// Panic:
//   - If any of the values are missing.
func LoadConfig() (*Config, error) {
	cfg := &Config{
		DBHost:     getEnvRequired("DB_HOST"),
		DBPort:     getEnvRequired("DB_PORT"),
		DBUser:     getEnvRequired("DB_USER"),
		DBPassword: getEnvRequired("DB_PASSWORD"),
		DBName:     getEnvRequired("DB_NAME"),
		GRPCPort:   getEnvRequired("GRPC_PORT"),
	}
	return cfg, nil
}

// getEnvRequired
// Gets the Env Variable, or fails if missing.
// Useful to make sure all required variables are available before starting the service.
func getEnvRequired(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintf("Required environment variable %s is missing", key))
	}
	return value
}
