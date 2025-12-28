package namelookup

import (
	"fmt"
	"strings"
)

// NameInfo holds information about a person
type NameInfo struct {
	FirstName    string
	LastName     string
	PhoneNumber  string
	Address      string
	LinkedInURL  string
	FacebookURL  string
	TwitterURL   string
	InstagramURL string
}

// SearchByName searches for information based on a full name
func SearchByName(fullName string) (string, error) {
	if fullName == "" {
		return "", fmt.Errorf("name cannot be empty")
	}

	// Parse the full name into first and last name
	firstName, lastName := parseName(fullName)

	// TODO: Implement actual API calls to people search services
	// Note: Most people search APIs are paid or have strict rate limits
	// Examples: Pipl, Whitepages, etc.
	// For now, this is a placeholder implementation

	info := &NameInfo{
		FirstName:    firstName,
		LastName:     lastName,
		PhoneNumber:  "Not Available (API Required)",
		Address:      "Not Available (API Required)",
		LinkedInURL:  fmt.Sprintf("https://www.linkedin.com/search/results/all/?keywords=%s+%s", firstName, lastName),
		FacebookURL:  fmt.Sprintf("https://www.facebook.com/search/top/?q=%s+%s", firstName, lastName),
		TwitterURL:   fmt.Sprintf("https://twitter.com/search?q=%s+%s", firstName, lastName),
		InstagramURL: fmt.Sprintf("https://www.instagram.com/explore/tags/%s%s/", strings.ToLower(firstName), strings.ToLower(lastName)),
	}

	result := formatNameInfo(info)
	return result, nil
}

// parseName splits a full name into first and last name
func parseName(fullName string) (string, string) {
	parts := strings.Fields(fullName)

	if len(parts) == 0 {
		return "", ""
	} else if len(parts) == 1 {
		return parts[0], ""
	}

	// First word is first name, rest is last name
	firstName := parts[0]
	lastName := strings.Join(parts[1:], " ")

	return firstName, lastName
}
