package username

import (
	"fmt"
	"math/rand"
	"strings"
	"sync"
	"time"

	"github.com/playwright-community/playwright-go"
)

// ⚠️ WARNING: This file contains advanced OSINT techniques
// Use ONLY for:
// - Educational purposes
// - Authorized penetration testing
// - Bug bounty programs (within scope)
// - Testing on YOUR OWN systems
//
// DO NOT use for unauthorized access or ToS violations

// AdvancedSearchUsername uses browser automation to bypass basic bot detection
func AdvancedSearchUsername(username string) (string, error) {
	// Remove @ symbol if present
	username = strings.TrimPrefix(username, "@")

	// Validate username format
	if username == "" {
		return "", fmt.Errorf("username cannot be empty")
	}

	if strings.Contains(username, " ") {
		return "", fmt.Errorf("invalid username: usernames cannot contain spaces")
	}

	if !isValidUsername(username) {
		return "", fmt.Errorf("invalid username: only letters, numbers, underscores, and hyphens allowed")
	}

	// fmt.Println("⚠️  Advanced Mode: Using browser automation")
	// fmt.Println("⚠️  This mode is slower but more accurate")
	// fmt.Println("⚠️  Use only for authorized testing")

	// Define social networks to check
	networks := []SocialNetwork{
		{Name: "GitHub", URL: fmt.Sprintf("https://github.com/%s", username)},
		{Name: "Reddit", URL: fmt.Sprintf("https://www.reddit.com/user/%s", username)},
		{Name: "Twitter", URL: fmt.Sprintf("https://twitter.com/%s", username)},
		{Name: "Medium", URL: fmt.Sprintf("https://medium.com/@%s", username)},
	}

	// Use browser automation for checking
	var wg sync.WaitGroup
	results := make([]UsernameResult, len(networks))

	// Rate limiting: check one at a time to be polite
	for i, network := range networks {
		wg.Add(1)
		go func(index int, net SocialNetwork) {
			defer wg.Done()

			fmt.Printf("Checking %s... ", net.Name)
			found := checkWithBrowser(net.URL, net.Name)
			results[index] = UsernameResult{
				Network: net.Name,
				Found:   found,
			}

			if found {
				fmt.Println("✓ Found")
			} else {
				fmt.Println("✗ Not Found")
			}

			// Polite delay between checks
			time.Sleep(2 * time.Second)
		}(i, network)

		wg.Wait() // Wait for each to complete (sequential, not parallel)
	}

	// Also check platforms that block bots (with warnings)
	blockedPlatforms := []SocialNetwork{
		{Name: "LinkedIn", URL: fmt.Sprintf("https://linkedin.com/in/%s", username)},
		{Name: "Instagram", URL: fmt.Sprintf("https://www.instagram.com/%s/", username)},
		{Name: "Facebook", URL: fmt.Sprintf("https://www.facebook.com/%s", username)},
	}

	for _, platform := range blockedPlatforms {
		results = append(results, UsernameResult{
			Network: platform.Name,
			Found:   false, // Don't check, add warning instead
		})
	}
