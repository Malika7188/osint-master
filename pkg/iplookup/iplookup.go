package iplookup

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// IPInfo holds geolocation information about an IP address
type IPInfo struct {
	IP          string  `json:"ip"`
	City        string  `json:"city"`
	Region      string  `json:"region"`
	Country     string  `json:"country"`
	CountryCode string  `json:"country_code"`
	ISP         string  `json:"org"`
	ASN         string  `json:"asn"`
	Timezone    string  `json:"timezone"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

// LookupIP performs IP geolocation lookup using multiple API providers
func LookupIP(ip string) (string, error) {
	// Validate IP address
	if ip == "" {
		return "", fmt.Errorf("IP address cannot be empty")
	}

	ip = strings.TrimSpace(ip)

	// Try multiple APIs for redundancy
	var info *IPInfo
	var err error

	// Try ip-api.com first (free, no key required, 45 requests/minute)
	info, err = lookupIPAPI(ip)
	if err == nil && info != nil {
		return formatIPInfo(info), nil
	}

	// Try ipinfo.io as fallback (free tier available)
	info, err = lookupIPInfo(ip)
	if err == nil && info != nil {
		return formatIPInfo(info), nil
	}

	// Try ipapi.co as second fallback
	info, err = lookupIPApiCo(ip)
	if err == nil && info != nil {
		return formatIPInfo(info), nil
	}

	// Try ipwhois.app as last resort
	info, err = lookupIPWhois(ip)
	if err == nil && info != nil {
		return formatIPInfo(info), nil
	}

	return "", fmt.Errorf("failed to lookup IP address from all providers")
}
