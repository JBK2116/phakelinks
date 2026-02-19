package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents all configuration variables needed to run this application
type Config struct {
	IsDev      bool
	PublicHost string
	PublicPort string
	DBHost     string
	DBPort     int64
	DBUser     string
	DBPassword string
	DBName     string
}

// Envs represents the access point for using all configuration variables
var Envs = initConfig()

// initConfig() loads all config variables into memory for use
func initConfig() Config {
	godotenv.Load()

	return Config{
		IsDev:      getEnvBool("IsDev", true),
		PublicHost: getEnv("PublicHost", "http://localhost"),
		PublicPort: getEnv("PublicPort", "8080"),
		DBHost:     getEnv("DBHost", "localhost"),
		DBPort:     getEnvAsInt("DBPort", 5432),
		DBUser:     getEnv("DBUser", "postgres"),
		DBPassword: getEnv("DBPassword", "random"),
		DBName:     getEnv("DBName", "phakelinks"),
	}
}

// getEnv() returns the value of the provided key if found, else returns fallback string
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return fallback
}

// getEnvAsInt() returns the value of the provided key if found, else returns the fallback int64
func getEnvAsInt(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			return fallback
		}

		return i
	}

	return fallback
}

// getEnvBool() returns the value of the provided key if found, else returns the fallback bool
func getEnvBool(key string, fallback bool) bool {
	if value, ok := os.LookupEnv(key); ok {
		i, err := strconv.ParseBool(value)
		if err != nil {
			return fallback
		}
		return i
	}
	return fallback
}
