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
