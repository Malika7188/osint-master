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
