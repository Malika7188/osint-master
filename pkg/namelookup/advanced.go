package namelookup

import (
	"fmt"
	"strings"
	"time"
)

// AdvancedSearchByName performs enhanced people search with additional sources
func AdvancedSearchByName(fullName string) (string, error) {
	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced people search\n")
	// result.WriteString("This mode provides additional search sources\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard search first
	startTime := time.Now()
	standardResult, err := SearchByName(fullName)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
