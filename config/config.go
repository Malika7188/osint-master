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
	