package configs

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config represents all configuration variables needed to run this application
type Config struct {
	IsDev            bool
	PublicHost       string
	PublicPort       string
	OPENAI_KEY       string
	CLOUDMERSIVE_KEY string
	DBHost           string
	DBPort           int64
	DBUser           string
	DBPassword       string
	DBName           string
}

// Envs represents the access point for using all configuration variables
var Envs = initConfig()

// initConfig() loads all config variables into memory for use
func initConfig() Config {
	godotenv.Load()

	return Config{
		IsDev:            getEnvBool("IsDev", true),
		PublicHost:       getEnv("PublicHost", "PublicHost"),
		PublicPort:       getEnv("PublicPort", "PublicPort"),
		OPENAI_KEY:       getEnv("OPENAI_KEY", "OPENAI_KEY"),
		CLOUDMERSIVE_KEY: getEnv("CLOUDMERSIVE_KEY", "CLOUDMERSIVE_KEY"),
		DBHost:           getEnv("DBHost", "DBHost"),
		DBPort:           getEnvAsint("DBPort", -1),
		DBUser:           getEnv("DBUser", "DBUser"),
		DBPassword:       getEnv("DBPassword", "DBPassword"),
		DBName:           getEnv("DBName", "DBName"),
	}
}

// getEnv() returns the value of the provided key if found, else returns fallback string
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// getEnvAsint() returns the value of the provided key if found as an int64 data type
func getEnvAsint(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
		return fallback
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
