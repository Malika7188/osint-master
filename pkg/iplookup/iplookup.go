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

// lookupIPAPI queries ip-api.com for IP information
func lookupIPAPI(ip string) (*IPInfo, error) {
	url := fmt.Sprintf("http://ip-api.com/json/%s?fields=status,message,country,countryCode,region,city,lat,lon,timezone,isp,org,as,query", ip)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ip-api.com returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Check if lookup was successful
	if status, ok := result["status"].(string); ok && status != "success" {
		if msg, ok := result["message"].(string); ok {
			return nil, fmt.Errorf("ip-api error: %s", msg)
		}
		return nil, fmt.Errorf("ip-api lookup failed")
	}

	// Parse response into IPInfo
	info := &IPInfo{
		IP: ip,
	}

	if city, ok := result["city"].(string); ok {
		info.City = city
	}
	if region, ok := result["region"].(string); ok {
		info.Region = region
	}
	if country, ok := result["country"].(string); ok {
		info.Country = country
	}
	if countryCode, ok := result["countryCode"].(string); ok {
		info.CountryCode = countryCode
	}
	if isp, ok := result["isp"].(string); ok {
		info.ISP = isp
	}
	if asn, ok := result["as"].(string); ok {
		info.ASN = asn
	}
	if timezone, ok := result["timezone"].(string); ok {
		info.Timezone = timezone
	}
	if lat, ok := result["lat"].(float64); ok {
		info.Latitude = lat
	}
	if lon, ok := result["lon"].(float64); ok {
		info.Longitude = lon
	}

	return info, nil
}

// lookupIPInfo queries ipinfo.io for IP information
func lookupIPInfo(ip string) (*IPInfo, error) {
	url := fmt.Sprintf("https://ipinfo.io/%s/json", ip)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ipinfo.io returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	info := &IPInfo{
		IP: ip,
	}

	if city, ok := result["city"].(string); ok {
		info.City = city
	}
	if region, ok := result["region"].(string); ok {
		info.Region = region
	}
	if country, ok := result["country"].(string); ok {
		info.CountryCode = country
		// Convert country code to full name
		info.Country = getCountryName(country)
	}
	if org, ok := result["org"].(string); ok {
		info.ISP = org
		// Extract ASN if present in org field (format: "AS15169 Google LLC")
		if strings.Contains(org, "AS") {
			parts := strings.Fields(org)
			if len(parts) > 0 && strings.HasPrefix(parts[0], "AS") {
				info.ASN = parts[0]
			}
		}
	}
	if timezone, ok := result["timezone"].(string); ok {
		info.Timezone = timezone
	}
	if loc, ok := result["loc"].(string); ok {
		// loc format: "latitude,longitude"
		parts := strings.Split(loc, ",")
		if len(parts) == 2 {
			fmt.Sscanf(parts[0], "%f", &info.Latitude)
			fmt.Sscanf(parts[1], "%f", &info.Longitude)
		}
	}

	return info, nil
}

// lookupIPApiCo queries ipapi.co for IP information
func lookupIPApiCo(ip string) (*IPInfo, error) {
	url := fmt.Sprintf("https://ipapi.co/%s/json/", ip)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// ipapi.co requires User-Agent header
	req.Header.Set("User-Agent", "OSINT-Master-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ipapi.co returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	// Check for error response
	if errMsg, ok := result["error"].(bool); ok && errMsg {
		if reason, ok := result["reason"].(string); ok {
			return nil, fmt.Errorf("ipapi.co error: %s", reason)
		}
		return nil, fmt.Errorf("ipapi.co lookup failed")
	}

	info := &IPInfo{
		IP: ip,
	}

	if city, ok := result["city"].(string); ok {
		info.City = city
	}
	if region, ok := result["region"].(string); ok {
		info.Region = region
	}
	if country, ok := result["country_name"].(string); ok {
		info.Country = country
	}
	if countryCode, ok := result["country_code"].(string); ok {
		info.CountryCode = countryCode
	}
	if org, ok := result["org"].(string); ok {
		info.ISP = org
	}
	if asn, ok := result["asn"].(string); ok {
		info.ASN = asn
	}
	if timezone, ok := result["timezone"].(string); ok {
		info.Timezone = timezone
	}
	if lat, ok := result["latitude"].(float64); ok {
		info.Latitude = lat
	}
	if lon, ok := result["longitude"].(float64); ok {
		info.Longitude = lon
	}

	return info, nil
}
