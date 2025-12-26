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

	// Check if disposable email
	info.IsDisposable = isDisposableEmail(info.Domain)

	// Check Gravatar
	info.GravatarExists, info.GravatarURL = checkGravatar(email)

	// Check email reputation (FREE - no API key needed)
	checkEmailReputation(email, info)

	// Check Have I Been Pwned (HIBP)
	breaches, err := checkHIBPWithKey(email, hibpAPIKey)
	if err == nil {
		info.Breaches = breaches
		info.BreachCount = len(breaches)
	}

	// Automatically check social media accounts
	fmt.Println("\nChecking social media accounts...")
	socialAccounts := checkSocialMediaAccounts(email)

	// Format output
	result := formatEmailInfo(info)
	result += "\n" + formatSocialAccounts(socialAccounts)
	return result, nil
}

// isValidEmail validates email format using regex
func isValidEmail(email string) bool {
	// RFC 5322 compliant email regex (simplified)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// extractDomain extracts domain from email
func extractDomain(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) == 2 {
		return parts[1]
	}
	return ""
}

// isDisposableEmail checks if email is from disposable email service
func isDisposableEmail(domain string) bool {
	// List of common disposable email domains
	disposableDomains := []string{
		"tempmail.com", "10minutemail.com", "guerrillamail.com",
		"mailinator.com", "throwaway.email", "temp-mail.org",
		"yopmail.com", "maildrop.cc", "trashmail.com",
	}

	for _, disposable := range disposableDomains {
		if strings.Contains(domain, disposable) {
			return true
		}
	}
	return false
}
