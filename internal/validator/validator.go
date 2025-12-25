package validator

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// ValidateIP validates if a string is a valid IP address
func ValidateIP(ip string) error {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return fmt.Errorf("invalid IP address format: %s", ip)
	}
	return nil
}

// ValidateDomain validates if a string is a valid domain name
func ValidateDomain(domain string) error {
	// Remove protocol if present
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimSuffix(domain, "/")

	// Basic domain validation regex
	domainRegex := regexp.MustCompile(`^([a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?\.)+[a-zA-Z]{2,}$`)

	if !domainRegex.MatchString(domain) {
		return fmt.Errorf("invalid domain format: %s", domain)
	}

	return nil
}

// ValidateUsername validates if a username is reasonable
func ValidateUsername(username string) error {
	// Remove @ symbol if present
	username = strings.TrimPrefix(username, "@")

	if len(username) == 0 {
		return fmt.Errorf("username cannot be empty")
	}

	