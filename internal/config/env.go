package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	DatabaseURL string
	SessionKey  string
	ServerPort  string
	BaseURL     string
	TimeZone    string
}

func LoadEnv() *EnvConfig {
	// Find the .env file
	envPath, err := findEnvFile()
	if err != nil {
		log.Printf("Warning: .env file not found, using environment variables: %v", err)
	} else {
		// Load the .env file
		err = godotenv.Load(envPath)
		if err != nil {
			log.Printf("Warning: Error loading .env file: %v", err)
		}
	}

	return &EnvConfig{
		DatabaseURL: getEnvWithDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/url_shortener?sslmode=disable"),
		SessionKey:  getEnvWithDefault("SESSION_KEY", "your-secret-key-here"),
		ServerPort:  getEnvWithDefault("SERVER_PORT", "8080"),
		BaseURL:     getEnvWithDefault("BASE_URL", "http://localhost:8080"),
		TimeZone:    getEnvWithDefault("TZ", "UTC"),
	}
}

func getEnvWithDefault(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// findEnvFile looks for .env file in current and parent directories
func findEnvFile() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return envPath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
