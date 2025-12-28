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
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")

	// Parse name
	firstName, lastName := parseName(fullName)
	fullNameEncoded := strings.ReplaceAll(fullName, " ", "%20")

	result.WriteString("Professional Networks:\n")
	result.WriteString(fmt.Sprintf("  - LinkedIn: https://www.linkedin.com/search/results/all/?keywords=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Indeed Resume: https://www.indeed.com/resumes?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - AngelList: https://angel.co/search?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Crunchbase: https://www.crunchbase.com/discover/people?q=%s\n", fullNameEncoded))

	result.WriteString("\nSocial Media Deep Search:\n")
	result.WriteString(fmt.Sprintf("  - Facebook People: https://www.facebook.com/search/people/?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Twitter Advanced: https://twitter.com/search?q=%s&f=user\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Instagram: https://www.instagram.com/explore/tags/%s/\n", strings.ToLower(strings.ReplaceAll(fullName, " ", ""))))
	result.WriteString(fmt.Sprintf("  - TikTok: https://www.tiktok.com/search/user?q=%s\n", fullNameEncoded))
	result.WriteString(fmt.Sprintf("  - Reddit: https://www.reddit.com/search/?q=%s\n", fullNameEncoded))
