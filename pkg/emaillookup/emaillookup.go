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

// checkGravatar checks if email has associated Gravatar
func checkGravatar(email string) (bool, string) {
	// Generate MD5 hash of email
	hash := md5.Sum([]byte(strings.ToLower(strings.TrimSpace(email))))
	hashStr := fmt.Sprintf("%x", hash)

	gravatarURL := fmt.Sprintf("https://www.gravatar.com/avatar/%s?d=404", hashStr)

	// Check if Gravatar exists
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(gravatarURL)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	// If status is 200, Gravatar exists
	if resp.StatusCode == http.StatusOK {
		return true, fmt.Sprintf("https://www.gravatar.com/avatar/%s", hashStr)
	}

	return false, ""
}

// checkEmailReputation checks email reputation using EmailRep.io (FREE - no API key needed)
func checkEmailReputation(email string, info *EmailInfo) {
	url := fmt.Sprintf("https://emailrep.io/%s", email)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	// EmailRep.io requires User-Agent
	req.Header.Set("User-Agent", "OSINT-Master-Educational-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return
	}

	// Parse reputation data
	if reputation, ok := result["reputation"].(string); ok {
		info.Reputation = reputation
	}

	if suspicious, ok := result["suspicious"].(bool); ok {
		info.Suspicious = suspicious
	}

	if references, ok := result["references"].(float64); ok {
		info.References = int(references)
	}
}

// checkHIBP checks Have I Been Pwned API for data breaches (without API key)
func checkHIBP(email string) ([]string, error) {
	return checkHIBPWithKey(email, "")
}

// checkHIBPWithKey checks Have I Been Pwned API with optional API key
func checkHIBPWithKey(email, apiKey string) ([]string, error) {
	// HIBP API v3 requires API key for email search
	// For educational purposes, we'll use the public breach list
	// In production, get API key from: https://haveibeenpwned.com/API/Key

	url := fmt.Sprintf("https://haveibeenpwned.com/api/v3/breachedaccount/%s?truncateResponse=false", email)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add required headers
	req.Header.Set("User-Agent", "OSINT-Master-Educational-Tool")

	// Add API key if provided
	if apiKey != "" {
		req.Header.Set("hibp-api-key", apiKey)
		fmt.Println("Using HIBP API key for breach check...")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 404 means no breaches found
	if resp.StatusCode == http.StatusNotFound {
		return []string{}, nil
	}

	// 401/403 means API key required or invalid
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		if apiKey == "" {
			return []string{"API key required - Get yours at: https://haveibeenpwned.com/API/Key"}, nil
		} else {
			return []string{"API key invalid - Check your HIBP_API_KEY"}, nil
		}
	}

	// 200 means breaches found
	if resp.StatusCode == http.StatusOK {
		var breaches []map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&breaches); err != nil {
			return nil, err
		}

		breachNames := make([]string, 0, len(breaches))
		for _, breach := range breaches {
			if name, ok := breach["Name"].(string); ok {
				breachNames = append(breachNames, name)
			}
		}
		return breachNames, nil
	}

	return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
}

// SocialAccount represents a social media account
type SocialAccount struct {
	Platform string
	Found    bool
	URL      string
	Method   string // How it was detected
}

// checkSocialMediaAccounts automatically checks for social media accounts
func checkSocialMediaAccounts(email string) []SocialAccount {
	accounts := []SocialAccount{}

	// Extract username from email for some platforms
	username := strings.Split(email, "@")[0]

	// Check various platforms
	platforms := []struct {
		name      string
		checkFunc func(string, string) (bool, string)
	}{
		{"Google/Gmail", checkGoogle},
		{"GitHub", checkGitHub},
		{"Twitter", checkTwitterByEmail},
		{"Facebook", checkFacebook},
		{"LinkedIn", checkLinkedIn},
		{"Instagram", checkInstagram},
	}

	for _, platform := range platforms {
		found, url := platform.checkFunc(email, username)
		accounts = append(accounts, SocialAccount{
			Platform: platform.name,
			Found:    found,
			URL:      url,
			Method:   "Automatic check",
		})
	}

	return accounts
}
