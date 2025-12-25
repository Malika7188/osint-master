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
	IPQualityScoreKey  string // Phone/IP quality validation
	AbstractAPIKey     string // Phone validation (AbstractAPI)
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
		IPQualityScoreKey: os.Getenv("IPQUALITYSCORE_KEY"),
		AbstractAPIKey:    os.Getenv("ABSTRACTAPI_KEY"),
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

	content := `# OSINT Master API Configuration
# Copy this file to ~/.osintmaster/.env and add your API keys
# Get API keys from the URLs provided in comments

# Have I Been Pwned API Key
# Get key at: https://haveibeenpwned.com/API/Key ($3.50/month)
# Enables: Full data breach checking for email addresses
HIBP_API_KEY=your_hibp_api_key_here

# Numverify Phone Number Validation API
# Get key at: https://numverify.com/product (Free: 100 requests/month)
# Enables: Phone carrier, line type, and location lookup
NUMVERIFY_KEY=your_numverify_key_here

# IPQualityScore Phone/IP Validation API
# Get key at: https://www.ipqualityscore.com/create-account (Free: 5k requests/month)
# Enables: Enhanced phone carrier, fraud detection, and IP validation
IPQUALITYSCORE_KEY=your_ipqualityscore_key_here

# AbstractAPI Phone Validation
# Get key at: https://app.abstractapi.com/api/phone-validation/pricing (Free: 100 requests/month)
# Enables: Phone number validation and carrier lookup
ABSTRACTAPI_KEY=your_abstractapi_key_here

# IPapi.co API Key (Optional - has free tier)
# Get key at: https://ipapi.co/api/ (Free: 30k requests/month)
# Enables: IP geolocation lookups
IPAPI_KEY=your_ipapi_key_here

# AbuseIPDB API Key (Optional)
# Get key at: https://www.abuseipdb.com/api (Free: 1k requests/day)
# Enables: IP abuse/blacklist checking
ABUSEIPDB_KEY=your_abuseipdb_key_here

# SecurityTrails API Key (Optional)
# Get key at: https://securitytrails.com/app/account/credentials
# Enables: Enhanced domain/subdomain enumeration
SECURITYTRAILS_KEY=your_securitytrails_key_here

# Twitter API Keys (Optional - requires paid plan)
# Get keys at: https://developer.twitter.com/
# Enables: Twitter username and email association lookup
TWITTER_API_KEY=your_twitter_api_key_here
TWITTER_API_SECRET=your_twitter_api_secret_here

# Google API Key (Optional)
# Get key at: https://console.cloud.google.com/apis/credentials
# Enables: Enhanced Google account lookup
GOOGLE_API_KEY=your_google_api_key_here

# Instructions:
# 1. Copy this file to ~/.osintmaster/.env
# 2. Replace "your_*_key_here" with actual API keys
# 3. Remove keys you don't have (tool works without them)
# 4. Never commit this file to git!
# 5. Keep your API keys secret!
`

	if err := os.WriteFile(envPath, []byte(content), 0600); err != nil {
		return err
	}

	fmt.Printf("Sample config file created at: %s\n", envPath)
	fmt.Println("Edit this file and add your API keys.")

	return nil
}
