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

// formatNameInfo formats the name information into a readable string
func formatNameInfo(info *NameInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("First Name: %s\n", info.FirstName))
	sb.WriteString(fmt.Sprintf("Last Name: %s\n", info.LastName))
	sb.WriteString(fmt.Sprintf("Phone Number: %s\n", info.PhoneNumber))
	sb.WriteString(fmt.Sprintf("Address: %s\n", info.Address))
	sb.WriteString("\nSocial Media Search URLs:\n")
	sb.WriteString(fmt.Sprintf("LinkedIn: %s\n", info.LinkedInURL))
	sb.WriteString(fmt.Sprintf("Facebook: %s\n", info.FacebookURL))
	sb.WriteString(fmt.Sprintf("Twitter: %s\n", info.TwitterURL))
	sb.WriteString(fmt.Sprintf("Instagram: %s\n", info.InstagramURL))
	sb.WriteString("\nNote: Phone and address require paid API access (e.g., Pipl, Whitepages)\n")
	sb.WriteString("Visit the URLs above to manually search on each platform.\n")

	return sb.String()
}
