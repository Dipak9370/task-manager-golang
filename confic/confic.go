package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                string
	DBUrl               string
	JWTSecret           string
	AutoCompleteMinutes int
}

func LoadEnv() *Config {
	// load .env
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, using system env")
	}

	cfg := &Config{
		Port:                getEnv("PORT", "8080"),
		DBUrl:               getEnv("DB_URL", ""),
		JWTSecret:           getEnv("JWT_SECRET", "secret"),
		AutoCompleteMinutes: getEnvAsInt("AUTO_COMPLETE_MINUTES", 5),
	}

	if cfg.DBUrl == "" {
		log.Fatal("DB_URL is required in .env")
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}
	return val
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := os.Getenv(key)
	if valStr == "" {
		return defaultVal
	}

	var val int
	_, err := fmt.Sscanf(valStr, "%d", &val)
	if err != nil {
		return defaultVal
	}
	return val
}
