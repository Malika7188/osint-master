package domain

import (
	"fmt"
	"strings"
	"time"
)

// AdvancedEnumerateDomain performs enhanced domain enumeration with additional analysis
func AdvancedEnumerateDomain(domain string) (string, error) {
	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced domain analysis\n")
	// result.WriteString("This mode performs additional checks and takes longer\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard enumeration first
	startTime := time.Now()
	standardResult, err := EnumerateDomain(domain)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")

	// Clean domain
	cleanDomain := strings.TrimPrefix(domain, "http://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "https://")
	cleanDomain = strings.TrimSuffix(cleanDomain, "/")
