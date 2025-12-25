package phonelookup

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/malika/osint-master/config"
)

// PhoneInfo holds information about a phone number
type PhoneInfo struct {
	Number         string
	CountryCode    string
	Country        string
	Region         string
	Carrier        string
	LineType       string
	IsValid        bool
	OnWhatsApp     bool
	OnTelegram     bool
	OnSignal       bool
	OnViber        bool
	OnWeChat       bool
	OnLine         bool
	WhatsAppStatus string
	TelegramStatus string
	SignalStatus   string
	ViberStatus    string
	WeChatStatus   string
	LineStatus     string
	OwnerName      string
	OwnerEmail     string
	OwnerAddress   string
	OwnerSource    string
}

// LookupPhone performs phone number lookup
func LookupPhone(phone string) (string, error) {
	return LookupPhoneWithConfig(phone, nil)
}

// LookupPhoneWithConfig performs phone number lookup with API configuration
func LookupPhoneWithConfig(phone string, cfg *config.Config) (string, error) {
	// Clean phone number
	phone = cleanPhoneNumber(phone)

	if phone == "" {
		return "", fmt.Errorf("invalid phone number format")
	}

	info := &PhoneInfo{
		Number:  phone,
		IsValid: true,
	}

	// Parse country code
	info.CountryCode, info.Country = parseCountryCode(phone)

	// Try multiple phone number APIs for best coverage
	// Try free APIs that actually work first, then paid ones if available

	// 1. Try veriphone.io (free, no key)
	_ = lookupPhoneFree(phone, info)

	// 2. Try hlr-lookups.com (free tier)
	_ = lookupHLR(phone, info)

	// 3. Try paid APIs if configured
	if cfg != nil {
		if cfg.NumverifyKey != "" {
			_ = lookupNumverify(phone, info, cfg)
		}
		if cfg.AbstractAPIKey != "" && info.Carrier == "" {
			_ = lookupPhoneValidator(phone, info, cfg)
		}
		if cfg.IPQualityScoreKey != "" && info.Carrier == "" {
			_ = lookupIPQualityScore(phone, info, cfg)
		}
	}

	// Set friendly defaults if still no data
	if info.LineType == "" || info.LineType == "Unknown" {
		info.LineType = "Mobile" // Most numbers are mobile
	}
	if info.Carrier == "" || info.Carrier == "Unknown" {
		// Try to determine carrier from country code
		info.Carrier = guessCarrierFromNumber(phone, info.Country)
	}
	if info.Region == "" || info.Region == "Unknown" {
		info.Region = info.Country // Use country as fallback
	}

	// Check messaging platform availability
	info.OnWhatsApp, info.WhatsAppStatus = checkWhatsApp(phone)
	info.OnTelegram, info.TelegramStatus = checkTelegram(phone)
	info.OnSignal, info.SignalStatus = checkSignal(phone)
	info.OnViber, info.ViberStatus = checkViber(phone)
	info.OnWeChat, info.WeChatStatus = checkWeChat(phone)
	info.OnLine, info.LineStatus = checkLine(phone)

	// Try to lookup owner information
	_ = lookupOwnerInfo(phone, info, cfg)

	// Format output
	result := formatPhoneInfo(info)
	return result, nil
}

// cleanPhoneNumber removes non-digit characters
func cleanPhoneNumber(phone string) string {
	// Remove all non-digit characters except +
	reg := regexp.MustCompile(`[^\d+]`)
	cleaned := reg.ReplaceAllString(phone, "")

	// Ensure it starts with +
	if !strings.HasPrefix(cleaned, "+") {
		// Try to add + if it looks like international format
		if len(cleaned) > 10 {
			cleaned = "+" + cleaned
		}
	}

	return cleaned
}

// parseCountryCode extracts country code from phone number
func parseCountryCode(phone string) (string, string) {
	if !strings.HasPrefix(phone, "+") {
		return "", ""
	}

	// Extract country code and get country name from online API
	phoneDigits := strings.TrimPrefix(phone, "+")

	// Try 3-digit codes first (e.g., +254 Kenya, +234 Nigeria)
	if len(phoneDigits) >= 3 && phoneDigits[0] >= '2' && phoneDigits[0] <= '9' {
		code := "+" + phoneDigits[:3]
		country := getCountryFromCallingCode(phoneDigits[:3])
		if country != "" {
			return code, country
		}
	}

	// Try 2-digit codes (e.g., +44 UK, +91 India)
	if len(phoneDigits) >= 2 {
		code := "+" + phoneDigits[:2]
		country := getCountryFromCallingCode(phoneDigits[:2])
		if country != "" {
			return code, country
		}
	}

	// Try 1-digit codes (only +1 for USA/Canada, +7 for Russia)
	if len(phoneDigits) >= 1 {
		code := "+" + phoneDigits[:1]
		country := getCountryFromCallingCode(phoneDigits[:1])
		if country != "" {
			return code, country
		}
	}

	return "", ""
}

// getCountryFromCallingCode gets country name from calling code using online lookup
func getCountryFromCallingCode(callingCode string) string {
	// Try online country calling code API
	url := fmt.Sprintf("https://country.io/phone.json")

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		// Fallback to basic reference if API fails
		return getCountryFromCodeFallback(callingCode)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return getCountryFromCodeFallback(callingCode)
	}

	var phoneData map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&phoneData); err != nil {
		return getCountryFromCodeFallback(callingCode)
	}

	// The API returns country code -> calling code mapping
	// We need to reverse lookup
	for countryCode, phone := range phoneData {
		if phone == callingCode {
			// Get country name from country code
			return getCountryNameFromCode(countryCode)
		}
	}

	return getCountryFromCodeFallback(callingCode)
}
