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
