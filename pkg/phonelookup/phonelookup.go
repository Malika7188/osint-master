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

	