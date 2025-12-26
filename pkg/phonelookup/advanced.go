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

	result.WriteString("  Checking Signal registration...\n")
	onSignal, signalStatus := checkSignal(cleanedPhone)
	if onSignal {
		result.WriteString(fmt.Sprintf("  ✓ Signal: %s\n", signalStatus))
	} else {
		result.WriteString(fmt.Sprintf("  ✗ Signal: %s\n", signalStatus))
	}

	result.WriteString("  Checking Viber registration...\n")
	onViber, viberStatus := checkViber(cleanedPhone)
	if onViber {
		result.WriteString(fmt.Sprintf("  ✓ Viber: %s\n", viberStatus))
	} else {
		result.WriteString(fmt.Sprintf("  ✗ Viber: %s\n", viberStatus))
	}

	result.WriteString("  Checking WeChat registration...\n")
	onWeChat, wechatStatus := checkWeChat(cleanedPhone)
	if onWeChat {
		result.WriteString(fmt.Sprintf("  ✓ WeChat: %s\n", wechatStatus))
	} else {
		result.WriteString(fmt.Sprintf("  ✗ WeChat: %s\n", wechatStatus))
	}

	result.WriteString("  Checking LINE registration...\n")
	onLine, lineStatus := checkLine(cleanedPhone)
	if onLine {
		result.WriteString(fmt.Sprintf("  ✓ LINE: %s\n", lineStatus))
	} else {
		result.WriteString(fmt.Sprintf("  ✗ LINE: %s\n", lineStatus))
	}

	// Additional lookup services
	result.WriteString("\nAdditional Lookup Services:\n")
	result.WriteString(fmt.Sprintf("  - TrueCaller: https://www.truecaller.com/search/us/%s\n", cleanedPhone))
	result.WriteString(fmt.Sprintf("  - WhitePages: https://www.whitepages.com/phone/%s\n", strings.ReplaceAll(cleanedPhone, "+", "")))
	(fmt.Println("  - Spy Dialer: https://www.spydialer.com/"))
	fmt.Println("  - NumLookup: https://www.numlookup.com/")

	result.WriteString("\nSocial Media Registration:\n")
	result.WriteString("  Some platforms allow phone number registration\n")
	result.WriteString("  - Facebook: May be linked to account\n")
	result.WriteString("  - Twitter: May be linked to account\n")
	result.WriteString("  - Instagram: May be linked to account\n")
	result.WriteString("  - LinkedIn: May be linked to account\n")

	result.WriteString("\nCarrier & CNAM Lookup:\n")
	(fmt.Println("  - FreeCarrierLookup: https://freecarrierlookup.com/"))
	(fmt.Println("  - Carrier Lookup: https://www.carrierlookup.com/"))
