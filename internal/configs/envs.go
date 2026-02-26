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
}

// Envs represents the access point for using all configuration variables
var Envs = initConfig()

// initConfig() loads all config variables into memory for use
func initConfig() Config {
	godotenv.Load()

	return Config{
		IsDev:            getEnvBool("IsDev", true),
		PublicHost:       getEnv("PublicHost", "http://localhost"),
		PublicPort:       getEnv("PublicPort", "8080"),
		OPENAI_KEY:       getEnv("OPENAI_KEY", "OPENAI_KEY"),
		CLOUDMERSIVE_KEY: getEnv("CLOUDMERSIVE_KEY", "CLOUDMERSIVE_KEY"),
	}
}

// getEnv() returns the value of the provided key if found, else returns fallback string
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
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
