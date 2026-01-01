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

	// Format results
	output := formatAdvancedResults(username, results)
	return output, nil
}

// checkWithBrowser uses Playwright to check if username exists
func checkWithBrowser(url, platform string) bool {
	// Initialize Playwright
	pw, err := playwright.Run()
	if err != nil {
		return false
	}
	defer pw.Stop()

	// Launch browser with stealth settings
	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(true), // Run without visible window
		Args: []string{
			"--disable-blink-features=AutomationControlled", // Hide automation
			"--disable-dev-shm-usage",
			"--no-sandbox",
			"--disable-setuid-sandbox",
		},
	})
	if err != nil {
		return false
	}
	defer browser.Close()

	// Create context with realistic settings
	userAgent := getRandomUserAgent()
	context, err := browser.NewContext(playwright.BrowserNewContextOptions{
		UserAgent: &userAgent,
		Viewport: &playwright.Size{
			Width:  1920,
			Height: 1080,
		},
		Locale:     playwright.String("en-US"),
		TimezoneId: playwright.String("America/New_York"),
	})
	if err != nil {
		return false
	}
	defer context.Close()

	page, err := context.NewPage()
	if err != nil {
		return false
	}

	// Inject JavaScript to hide automation markers
	err = page.AddInitScript(playwright.Script{
		Content: playwright.String(`
			// Hide webdriver property
			Object.defineProperty(navigator, 'webdriver', {
				get: () => false
			});

			// Add chrome object
			window.chrome = {
				runtime: {}
			};

			// Mock plugins
			Object.defineProperty(navigator, 'plugins', {
				get: () => [1, 2, 3, 4, 5]
			});

			// Mock languages
			Object.defineProperty(navigator, 'languages', {
				get: () => ['en-US', 'en']
			});
		`),
	})

	// Navigate to URL
	_, err = page.Goto(url, playwright.PageGotoOptions{
		WaitUntil: playwright.WaitUntilStateNetworkidle,
		Timeout:   playwright.Float(30000), // 30 second timeout
	})
	if err != nil {
		return false
	}

	// Wait a bit to simulate human behavior
	page.WaitForTimeout(1000 + float64(rand.Intn(2000))) // Random 1-3 second delay

	// Get page content
	content, err := page.Content()
	if err != nil {
		return false
	}

	// Platform-specific detection
	return detectProfile(content, platform)
}

// detectProfile checks if a profile exists based on page content
func detectProfile(content, platform string) bool {
	content = strings.ToLower(content)

	// Common "not found" indicators
	notFoundIndicators := []string{
		"page not found",
		"this page isn't available",
		"sorry, this page isn't available",
		"the page you requested was not found",
	}

	for _, indicator := range notFoundIndicators {
		if strings.Contains(content, indicator) {
			return false
		}
	}

	// Platform-specific checks
	switch platform {
	case "GitHub":
		// GitHub shows 404 for non-existent users
		if strings.Contains(content, "404") && strings.Contains(content, "not found") {
			return false
		}
		// Look for profile indicators
		return strings.Contains(content, "profile") || strings.Contains(content, "repositories")

	case "Reddit":
		if strings.Contains(content, "sorry, nobody on reddit goes by that name") {
			return false
		}
		return strings.Contains(content, "karma") || strings.Contains(content, "posts")

	case "Twitter":
		if strings.Contains(content, "this account doesn't exist") {
			return false
		}
		return strings.Contains(content, "tweets") || strings.Contains(content, "followers")

	case "Medium":
		return strings.Contains(content, "stories") || strings.Contains(content, "followers")

	default:
		// Generic check
		return !strings.Contains(content, "not found")
	}
}
