package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Config holds application configuration
type Config struct {
	// API Keys (load from environment variables or .env file)
	HIBPAPIKey         string // Have I Been Pwned
	IPAPIKey           string
	AbuseIPDBKey       string
	PiplAPIKey         string
	SecurityTrailsKey  string
	NumverifyKey       string // Phone number verification
	GoogleAPIKey       string
	TwitterAPIKey      string
	TwitterAPISecret   string
}

// LoadConfig loads configuration from environment variables and .env file
func LoadConfig() *Config {
	// Try to load from .env file first
	loadEnvFile()

	config := &Config{
		HIBPAPIKey:        os.Getenv("HIBP_API_KEY"),
		IPAPIKey:          os.Getenv("IPAPI_KEY"),
		AbuseIPDBKey:      os.Getenv("ABUSEIPDB_KEY"),
		PiplAPIKey:        os.Getenv("PIPL_API_KEY"),
		SecurityTrailsKey: os.Getenv("SECURITYTRAILS_KEY"),
		NumverifyKey:      os.Getenv("NUMVERIFY_KEY"),
		GoogleAPIKey:      os.Getenv("GOOGLE_API_KEY"),
		TwitterAPIKey:     os.Getenv("TWITTER_API_KEY"),
		TwitterAPISecret:  os.Getenv("TWITTER_API_SECRET"),
	}

	return config
}

// loadEnvFile loads environment variables from .env file
func loadEnvFile() {
	// Get user's home directory
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	// Try ~/.osintmaster/.env first
	envPath := filepath.Join(home, ".osintmaster", ".env")
	if _, err := os.Stat(envPath); os.IsNotExist(err) {
		// Try current directory .env
		envPath = ".env"
		if _, err := os.Stat(envPath); os.IsNotExist(err) {
			return // No .env file found
		}
	}

	// Read .env file
	file, err := os.Open(envPath)
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip comments and empty lines
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse KEY=VALUE
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		value = strings.Trim(value, "\"'")

		// Set environment variable (don't override existing)
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}
}

// GetConfigPath returns the path to the config directory
func GetConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	configDir := filepath.Join(home, ".osintmaster")
	return configDir, nil
}

// EnsureConfigDir creates the config directory if it doesn't exist
func EnsureConfigDir() error {
	configDir, err := GetConfigPath()
	if err != nil {
		return err
	}

	return os.MkdirAll(configDir, 0700)
}

// CreateSampleEnvFile creates a sample .env file with instructions
func CreateSampleEnvFile() error {
	if err := EnsureConfigDir(); err != nil {
		return err
	}

	configDir, err := GetConfigPath()
	if err != nil {
		return err
	}

	envPath := filepath.Join(configDir, ".env")

	// Check if file already exists
	if _, err := os.Stat(envPath); err == nil {
		return fmt.Errorf("config file already exists at: %s", envPath)
	}

	