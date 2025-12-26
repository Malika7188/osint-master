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

// checkGoogle checks if email is a Gmail/Google account
func checkGoogle(email, username string) (bool, string) {
	// If it's a gmail.com email, account definitely exists
	if strings.HasSuffix(email, "@gmail.com") {
		return true, fmt.Sprintf("https://mail.google.com/mail/u/%s", email)
	}
	return false, ""
}

// checkGitHub checks if email is associated with GitHub
func checkGitHub(email, username string) (bool, string) {
	// Try to check GitHub API for user by username (from email)
	url := fmt.Sprintf("https://api.github.com/users/%s", username)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false, ""
	}

	req.Header.Set("User-Agent", "OSINT-Master-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, fmt.Sprintf("https://github.com/%s", username)
	}

	return false, ""
}

// checkTwitterByEmail checks Twitter
func checkTwitterByEmail(email, username string) (bool, string) {
	// Twitter doesn't allow email-based lookup without API key
	// Return search URL for manual verification
	return false, fmt.Sprintf("https://twitter.com/search?q=%s", email)
}

// checkFacebook checks Facebook
func checkFacebook(email, username string) (bool, string) {
	// Facebook requires login for email-based search
	// Return search URL for manual verification
	return false, fmt.Sprintf("https://www.facebook.com/search/people/?q=%s", email)
}

// checkLinkedIn checks LinkedIn
func checkLinkedIn(email, username string) (bool, string) {
	// LinkedIn requires login for email-based search
	// Return search URL for manual verification
	return false, fmt.Sprintf("https://www.linkedin.com/search/results/people/?keywords=%s", email)
}

// checkInstagram checks Instagram
func checkInstagram(email, username string) (bool, string) {
	// Instagram doesn't allow email-based lookup
	// Try username instead
	url := fmt.Sprintf("https://www.instagram.com/%s/", username)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return false, ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return true, url
	}

	return false, ""
}

// formatSocialAccounts formats social media account results
func formatSocialAccounts(accounts []SocialAccount) string {
	var sb strings.Builder

	sb.WriteString("Social Media Account Detection:\n")
	sb.WriteString("================================\n\n")

	foundCount := 0
	for _, account := range accounts {
		if account.Found {
			foundCount++
			sb.WriteString(fmt.Sprintf("✓ %s: FOUND\n", account.Platform))
			sb.WriteString(fmt.Sprintf("  URL: %s\n", account.URL))
		} else {
			sb.WriteString(fmt.Sprintf("✗ %s: Not found (or requires login to verify)\n", account.Platform))
			if account.URL != "" {
				sb.WriteString(fmt.Sprintf("  Search: %s\n", account.URL))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Summary: %d accounts found automatically\n", foundCount))
	sb.WriteString("\nNote: Some platforms require login to search by email.\n")
	sb.WriteString("      Links provided for manual verification where needed.\n")

	return sb.String()
}

// formatEmailInfo formats email information into readable string
func formatEmailInfo(info *EmailInfo) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Email Address: %s\n", info.Email))
	sb.WriteString(fmt.Sprintf("Domain: %s\n", info.Domain))
	sb.WriteString(fmt.Sprintf("Valid Format: %v\n", info.IsValid))
	sb.WriteString(fmt.Sprintf("Disposable Email: %v\n", info.IsDisposable))

	if info.IsDisposable {
		sb.WriteString("⚠️  Warning: This is a temporary/disposable email service\n")
	}

	// Email Reputation (FREE API - EmailRep.io)
	sb.WriteString("\nEmail Reputation (via EmailRep.io - FREE):\n")
	if info.Reputation != "" {
		sb.WriteString(fmt.Sprintf("  Reputation: %s\n", info.Reputation))
		sb.WriteString(fmt.Sprintf("  Suspicious: %v\n", info.Suspicious))
		if info.References > 0 {
			sb.WriteString(fmt.Sprintf("  References: %d (times seen in data)\n", info.References))
		}
		if info.Suspicious {
			sb.WriteString("  ⚠️  Warning: Email marked as suspicious by reputation database\n")
		}
	} else {
		sb.WriteString("  Status: No reputation data available\n")
		sb.WriteString("  Note: This is a free service that may not have all emails\n")
	}

	sb.WriteString("\nGravatar:\n")
	if info.GravatarExists {
		sb.WriteString(fmt.Sprintf("  Found: Yes\n"))
		sb.WriteString(fmt.Sprintf("  URL: %s\n", info.GravatarURL))
	} else {
		sb.WriteString("  Found: No\n")
	}

	sb.WriteString("\nData Breach Check (Have I Been Pwned):\n")
	if info.BreachCount == 0 {
		sb.WriteString("  Status: No breaches found (or API key required)\n")
		sb.WriteString(fmt.Sprintf("  🔗 Check directly: https://haveibeenpwned.com/account/%s\n", info.Email))
		sb.WriteString("  Note: Visit the link above to see detailed breach information\n")
	} else {
		sb.WriteString(fmt.Sprintf("  Breaches Found: %d\n", info.BreachCount))
		sb.WriteString("  Breach Names:\n")
		for _, breach := range info.Breaches {
			sb.WriteString(fmt.Sprintf("    - %s\n", breach))
		}
		sb.WriteString("\n⚠️  WARNING: This email has been found in data breaches!\n")
		sb.WriteString("  Recommendation: Change passwords immediately\n")
		sb.WriteString(fmt.Sprintf("  🔗 View details: https://haveibeenpwned.com/account/%s\n", info.Email))
	}

	return sb.String()
}
