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
	