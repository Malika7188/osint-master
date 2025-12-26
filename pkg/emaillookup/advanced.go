package emaillookup

import (
	"fmt"
	"strings"
	"time"
)

// AdvancedLookupEmail performs enhanced email lookup with additional checks
func AdvancedLookupEmail(email, hibpAPIKey string) (string, error) {
	// Validate email format
	if !isValidEmail(email) {
		return "", fmt.Errorf("invalid email format: %s", email)
	}

	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced email analysis\n")
	// result.WriteString("This mode performs additional checks and takes longer\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard lookup first
	startTime := time.Now()
	standardResult, err := LookupEmailWithConfig(email, hibpAPIKey)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")

	// Additional advanced checks
	result.WriteString("Enhanced Social Media Discovery:\n")

	// Check more platforms in advanced mode
	username := strings.Split(email, "@")[0]
	domain := extractDomain(email)

	result.WriteString(fmt.Sprintf("  Username pattern: %s\n", username))
	result.WriteString(fmt.Sprintf("  Domain: %s\n", domain))
	result.WriteString("  Checking extended platforms...\n\n")

	// Additional platforms to check in advanced mode
	advancedPlatforms := []struct {
		name string
		url  string
	}{
		{"Reddit", fmt.Sprintf("https://www.reddit.com/user/%s", username)},
		{"Medium", fmt.Sprintf("https://medium.com/@%s", username)},
		{"Dev.to", fmt.Sprintf("https://dev.to/%s", username)},
		{"Behance", fmt.Sprintf("https://www.behance.net/%s", username)},
