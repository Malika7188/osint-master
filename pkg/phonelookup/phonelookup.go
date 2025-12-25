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

	