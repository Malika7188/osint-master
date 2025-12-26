package emaillookup

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// EmailInfo holds information about an email address
type EmailInfo struct {
	Email           string
	IsValid         bool
	Domain          string
	IsDisposable    bool
	BreachCount     int
	Breaches        []string
	GravatarExists  bool
	GravatarURL     string
	Reputation      string
	Suspicious      bool
	References      int
}

// LookupEmail performs comprehensive email address lookup
func LookupEmail(email string) (string, error) {
	return LookupEmailWithConfig(email, "")
}

// LookupEmailWithConfig performs email lookup with API key support
func LookupEmailWithConfig(email, hibpAPIKey string) (string, error) {
	// Validate email format
	if !isValidEmail(email) {
		return "", fmt.Errorf("invalid email format: %s", email)
	}

	email = strings.ToLower(strings.TrimSpace(email))

	info := &EmailInfo{
		Email:    email,
		IsValid:  true,
		Domain:   extractDomain(email),
		Breaches: make([]string, 0),
	}
