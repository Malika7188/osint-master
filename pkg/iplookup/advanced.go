package iplookup

import (
	"fmt"
	"strings"
	"time"
)

// AdvancedLookupIP performs enhanced IP address analysis
func AdvancedLookupIP(ip string) (string, error) {
	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced IP address analysis\n")
	// result.WriteString("This mode performs additional checks and takes longer\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard lookup first
	startTime := time.Now()
	standardResult, err := LookupIP(ip)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")
