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

	