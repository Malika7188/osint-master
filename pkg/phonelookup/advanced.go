package phonelookup

import (
	"fmt"
	"strings"
	"time"

	"github.com/malika/osint-master/config"
)

// AdvancedLookupPhone performs enhanced phone number lookup
func AdvancedLookupPhone(phone string) (string, error) {
	return AdvancedLookupPhoneWithConfig(phone, nil)
}

// AdvancedLookupPhoneWithConfig performs enhanced phone number lookup with API configuration
func AdvancedLookupPhoneWithConfig(phone string, cfg *config.Config) (string, error) {
	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced phone number analysis\n")
	// result.WriteString("This mode performs additional checks and takes longer\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard lookup first
	startTime := time.Now()
	standardResult, err := LookupPhoneWithConfig(phone, cfg)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")

	// Clean phone for display
	cleanedPhone := cleanPhoneNumber(phone)

	result.WriteString("Enhanced Verification:\n")
	result.WriteString(fmt.Sprintf("  Cleaned Number: %s\n\n", cleanedPhone))

	// Perform automated messaging platform checks
	result.WriteString("🔍 Automated Messaging Platform Detection:\n")
	result.WriteString("  Checking WhatsApp registration...\n")
	onWhatsApp, whatsAppStatus := checkWhatsApp(cleanedPhone)
	if onWhatsApp {
		result.WriteString(fmt.Sprintf("  ✓ WhatsApp: %s\n", whatsAppStatus))
	} else {
		result.WriteString(fmt.Sprintf("  ✗ WhatsApp: %s\n", whatsAppStatus))
	}

	result.WriteString("  Checking Telegram registration...\n")
	onTelegram, telegramStatus := checkTelegram(cleanedPhone)
	if onTelegram {
		result.WriteString(fmt.Sprintf("  ✓ Telegram: %s\n", telegramStatus))
	} else {
		result.WriteString(fmt.Sprintf("  ✗ Telegram: %s\n", telegramStatus))
	}
