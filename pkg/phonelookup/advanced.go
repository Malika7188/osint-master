package phonelookup

import (
	"fmt"
	"strings"
	"time"

	"github.com/malika/osint-master/config"
)

// AdvancedLookupPhone performs enhanced phone number lookup
func AdvancedLookupPhone(phone string) (string, error) {
	return AdvancedLookupPhoneWithConfig(phone, nil)
}

// AdvancedLookupPhoneWithConfig performs enhanced phone number lookup with API configuration
func AdvancedLookupPhoneWithConfig(phone string, cfg *config.Config) (string, error) {
	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced phone number analysis\n")
	// result.WriteString("This mode performs additional checks and takes longer\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard lookup first
	startTime := time.Now()
	standardResult, err := LookupPhoneWithConfig(phone, cfg)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")
