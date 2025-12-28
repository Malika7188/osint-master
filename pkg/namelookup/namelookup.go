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
