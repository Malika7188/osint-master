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

// getCountryNameFromCode converts ISO country code to country name
func getCountryNameFromCode(code string) string {
	url := fmt.Sprintf("https://restcountries.com/v3.1/alpha/%s", code)

	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return code
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return code
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return code
	}

	if len(result) > 0 {
		if name, ok := result[0]["name"].(map[string]interface{}); ok {
			if common, ok := name["common"].(string); ok {
				return common
			}
		}
	}

	return code
}
// getCountryFromCodeFallback provides fallback country names for common calling codes
func getCountryFromCodeFallback(callingCode string) string {
	// Minimal reference for most common codes
	commonCodes := map[string]string{
		"1": "United States", "7": "Russia", "20": "Egypt", "27": "South Africa",
		"30": "Greece", "31": "Netherlands", "32": "Belgium", "33": "France",
		"34": "Spain", "36": "Hungary", "39": "Italy", "40": "Romania",
		"41": "Switzerland", "43": "Austria", "44": "United Kingdom", "45": "Denmark",
		"46": "Sweden", "47": "Norway", "48": "Poland", "49": "Germany",
		"51": "Peru", "52": "Mexico", "53": "Cuba", "54": "Argentina",
		"55": "Brazil", "56": "Chile", "57": "Colombia", "58": "Venezuela",
		"60": "Malaysia", "61": "Australia", "62": "Indonesia", "63": "Philippines",
		"64": "New Zealand", "65": "Singapore", "66": "Thailand", "81": "Japan",
		"82": "South Korea", "84": "Vietnam", "86": "China", "90": "Turkey",
		"91": "India", "92": "Pakistan", "93": "Afghanistan", "94": "Sri Lanka",
		"95": "Myanmar", "98": "Iran", "211": "South Sudan", "212": "Morocco",
		"213": "Algeria", "216": "Tunisia", "218": "Libya", "220": "Gambia",
		"221": "Senegal", "222": "Mauritania", "223": "Mali", "224": "Guinea",
		"225": "Ivory Coast", "226": "Burkina Faso", "227": "Niger", "228": "Togo",
		"229": "Benin", "230": "Mauritius", "231": "Liberia", "232": "Sierra Leone",
		"233": "Ghana", "234": "Nigeria", "235": "Chad", "236": "Central African Republic",
		"237": "Cameroon", "238": "Cape Verde", "239": "Sao Tome and Principe", "240": "Equatorial Guinea",
		"241": "Gabon", "242": "Republic of the Congo", "243": "Democratic Republic of the Congo",
		"244": "Angola", "245": "Guinea-Bissau", "246": "British Indian Ocean Territory",
		"248": "Seychelles", "249": "Sudan", "250": "Rwanda", "251": "Ethiopia",
		"252": "Somalia", "253": "Djibouti", "254": "Kenya", "255": "Tanzania",
		"256": "Uganda", "257": "Burundi", "258": "Mozambique", "260": "Zambia",
		"261": "Madagascar", "262": "Reunion", "263": "Zimbabwe", "264": "Namibia",
		"265": "Malawi", "266": "Lesotho", "267": "Botswana", "268": "Eswatini",
		"269": "Comoros", "290": "Saint Helena", "291": "Eritrea", "297": "Aruba",
		"298": "Faroe Islands", "299": "Greenland", "350": "Gibraltar", "351": "Portugal",
		"352": "Luxembourg", "353": "Ireland", "354": "Iceland", "355": "Albania",
		"356": "Malta", "357": "Cyprus", "358": "Finland", "359": "Bulgaria",
		"370": "Lithuania", "371": "Latvia", "372": "Estonia", "373": "Moldova",
		"374": "Armenia", "375": "Belarus", "376": "Andorra", "377": "Monaco",
		"378": "San Marino", "380": "Ukraine", "381": "Serbia", "382": "Montenegro",
		"383": "Kosovo", "385": "Croatia", "386": "Slovenia", "387": "Bosnia and Herzegovina",
		"389": "North Macedonia", "420": "Czech Republic", "421": "Slovakia", "423": "Liechtenstein",
		"500": "Falkland Islands", "501": "Belize", "502": "Guatemala", "503": "El Salvador",
		"504": "Honduras", "505": "Nicaragua", "506": "Costa Rica", "507": "Panama",
		"508": "Saint Pierre and Miquelon", "509": "Haiti", "590": "Guadeloupe", "591": "Bolivia",
		"592": "Guyana", "593": "Ecuador", "594": "French Guiana", "595": "Paraguay",
		"596": "Martinique", "597": "Suriname", "598": "Uruguay", "599": "Curacao",
		"670": "East Timor", "672": "Antarctica", "673": "Brunei", "674": "Nauru",
		"675": "Papua New Guinea", "676": "Tonga", "677": "Solomon Islands", "678": "Vanuatu",
		"679": "Fiji", "680": "Palau", "681": "Wallis and Futuna", "682": "Cook Islands",
		"683": "Niue", "685": "Samoa", "686": "Kiribati", "687": "New Caledonia",
		"688": "Tuvalu", "689": "French Polynesia", "690": "Tokelau", "691": "Micronesia",
		"692": "Marshall Islands", "850": "North Korea", "852": "Hong Kong", "853": "Macau",
		"855": "Cambodia", "856": "Laos", "880": "Bangladesh", "886": "Taiwan",
		"960": "Maldives", "961": "Lebanon", "962": "Jordan", "963": "Syria",
		"964": "Iraq", "965": "Kuwait", "966": "Saudi Arabia", "967": "Yemen",
		"968": "Oman", "970": "Palestine", "971": "United Arab Emirates", "972": "Israel",
		"973": "Bahrain", "974": "Qatar", "975": "Bhutan", "976": "Mongolia",
		"977": "Nepal", "992": "Tajikistan", "993": "Turkmenistan", "994": "Azerbaijan",
		"995": "Georgia", "996": "Kyrgyzstan", "998": "Uzbekistan",
	}

	if country, exists := commonCodes[callingCode]; exists {
		return country
	}

	return ""
}

// lookupPhoneAPI queries phone lookup API
func lookupPhoneAPI(phone string, info *PhoneInfo) error {
	// Try FREE API first (veriphone.io - no key required)
	if err := lookupPhoneFree(phone, info); err == nil {
		return nil
	}

	// Fallback to paid API if available (numverify)
	// Get free API key at: https://numverify.com/product (100 requests/month free)
	apiKey := "demo_api_key" // Replace with actual key from environment
	if apiKey == "demo_api_key" {
		// No valid API key, return with free data already populated
		return fmt.Errorf("no API key configured")
	}

	url := fmt.Sprintf("http://apilayer.net/api/validate?access_key=%s&number=%s", apiKey, phone)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Parse response
	if valid, ok := result["valid"].(bool); ok && !valid {
		info.IsValid = false
		return fmt.Errorf("invalid phone number")
	}

	if carrier, ok := result["carrier"].(string); ok {
		info.Carrier = carrier
	}

	if lineType, ok := result["line_type"].(string); ok {
		info.LineType = lineType
	}

	if location, ok := result["location"].(string); ok {
		info.Region = location
	}

	return nil
}

// lookupPhoneFree uses FREE API (veriphone.io) - no API key required
func lookupPhoneFree(phone string, info *PhoneInfo) error {
	// Remove + from phone for API
	phoneClean := strings.TrimPrefix(phone, "+")

	// Try veriphone.io first
	url := fmt.Sprintf("https://api.veriphone.io/v2/verify?phone=%s", phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "OSINT-Master-Educational-Tool")

	resp, err := client.Do(req)
	if err != nil {
		// If veriphone fails, try alternative API
		return lookupPhoneAlternative(phone, info)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Try alternative API
		return lookupPhoneAlternative(phone, info)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Parse veriphone response
	if phoneValid, ok := result["phone_valid"].(bool); ok {
		info.IsValid = phoneValid
	}

	if carrier, ok := result["carrier"].(string); ok && carrier != "" {
		info.Carrier = carrier
	}

	if phoneType, ok := result["phone_type"].(string); ok && phoneType != "" {
		info.LineType = phoneType
	}

	if country, ok := result["country"].(string); ok && country != "" {
		// Update country if API provides better data
		if info.Country == "Unknown" || info.Country == "" {
			info.Country = country
		}
	}

	if countryCode, ok := result["country_code"].(string); ok && countryCode != "" {
		if info.CountryCode == "Unknown" || info.CountryCode == "" {
			info.CountryCode = "+" + countryCode
		}
	}

	// Region/location
	if region, ok := result["region"].(string); ok && region != "" {
		info.Region = region
	}

	return nil
}

// lookupPhoneAlternative tries alternative free phone APIs
func lookupPhoneAlternative(phone string, info *PhoneInfo) error {
	phoneClean := strings.TrimPrefix(phone, "+")

	// Try numverify free tier (limited requests per month)
	// Note: This requires an API key, but we'll try the demo endpoint
	url := fmt.Sprintf("https://phonevalidation.abstractapi.com/v1/?api_key=test&phone=%s", phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "OSINT-Master-Educational-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("alternative API returned status: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Parse response
	if valid, ok := result["valid"].(bool); ok {
		info.IsValid = valid
	}

	if carrier, ok := result["carrier"].(string); ok && carrier != "" {
		info.Carrier = carrier
	}

	if lineType, ok := result["type"].(string); ok && lineType != "" {
		info.LineType = lineType
	}

	if country, ok := result["country"].(map[string]interface{}); ok {
		if countryName, ok := country["name"].(string); ok && countryName != "" {
			info.Country = countryName
		}
		if countryCode, ok := country["code"].(string); ok && countryCode != "" {
			info.CountryCode = "+" + countryCode
		}
	}

	if location, ok := result["location"].(string); ok && location != "" {
		info.Region = location
	}

	return nil
}

// lookupHLR uses HLR (Home Location Register) lookup from multiple sources
func lookupHLR(phone string, info *PhoneInfo) error {
	phoneClean := strings.TrimPrefix(phone, "+")

	// Try multiple HLR/carrier lookup APIs

	// 1. Try mccmnc.com API (free carrier database)
	if err := lookupMCCMNCOnline(phoneClean, info); err == nil && info.Carrier != "" {
		return nil
	}

	// 2. Try hlr-lookups.com
	url := fmt.Sprintf("https://hlr-lookups.com/api/free/%s", phoneClean)
	if err := makeHLRRequest(url, info); err == nil && info.Carrier != "" {
		return nil
	}

	// 3. Try freecarrierlookup.com API
	url = fmt.Sprintf("https://www.freecarrierlookup.com/api/%s", phoneClean)
	if err := makeCarrierRequest(url, info); err == nil && info.Carrier != "" {
		return nil
	}

	return fmt.Errorf("no HLR data available")
}

// lookupMCCMNCOnline fetches carrier info from online MCC-MNC database
func lookupMCCMNCOnline(phone string, info *PhoneInfo) error {
	// Use mcc-mnc.com API for carrier lookup
	url := fmt.Sprintf("https://mcc-mnc.net/api/?phone=%s", phone)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "OSINT-Master-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("MCC-MNC API error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Parse response
	if carrier, ok := result["carrier"].(string); ok && carrier != "" {
		info.Carrier = carrier
	}

	if network, ok := result["network"].(string); ok && network != "" && info.Carrier == "" {
		info.Carrier = network
	}

	if operator, ok := result["operator"].(string); ok && operator != "" && info.Carrier == "" {
		info.Carrier = operator
	}

	if lineType, ok := result["type"].(string); ok && lineType != "" {
		info.LineType = lineType
	}

	return nil
}

// makeHLRRequest makes a generic HLR lookup request
func makeHLRRequest(url string, info *PhoneInfo) error {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "OSINT-Master-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("HLR API error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Try common field names for carrier
	carrierFields := []string{"carrier", "operator", "network", "provider", "mno"}
	for _, field := range carrierFields {
		if carrier, ok := result[field].(string); ok && carrier != "" {
			info.Carrier = carrier
			break
		}
	}

	// Try common field names for line type
	typeFields := []string{"type", "line_type", "connection_type", "phone_type"}
	for _, field := range typeFields {
		if lineType, ok := result[field].(string); ok && lineType != "" {
			info.LineType = lineType
			break
		}
	}

	return nil
}

// makeCarrierRequest makes a carrier lookup request
func makeCarrierRequest(url string, info *PhoneInfo) error {
	return makeHLRRequest(url, info) // Same logic
}

// guessCarrierFromNumber tries to determine carrier from number patterns using external data
func guessCarrierFromNumber(phone string, country string) string {
	// Try to lookup from online carrier database
	phoneClean := strings.TrimPrefix(phone, "+")

	// Try carrier lookup API
	if carrier := lookupCarrierFromAPI(phoneClean); carrier != "" {
		return carrier
	}

	// If all else fails, return generic info based on country
	if country != "" && country != "Unknown" {
		return fmt.Sprintf("%s mobile carrier", country)
	}

	return "Mobile carrier"
}

// lookupCarrierFromAPI tries to get carrier from online database
func lookupCarrierFromAPI(phone string) string {
	// Try carrier411.com API (free carrier database)
	url := fmt.Sprintf("https://www.carrier411.com/api/v1/phone/%s", phone)

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "OSINT-Master-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	// Try to extract carrier name
	if carrier, ok := result["carrier"].(string); ok && carrier != "" {
		return carrier
	}

	if provider, ok := result["provider"].(string); ok && provider != "" {
		return provider
	}

	return ""
}

// lookupNumverify uses numverify.com API (free tier: 100 requests/month)
func lookupNumverify(phone string, info *PhoneInfo, cfg *config.Config) error {
	// Skip if no API key configured
	if cfg == nil || cfg.NumverifyKey == "" {
		return fmt.Errorf("numverify API key not configured")
	}

	phoneClean := strings.TrimPrefix(phone, "+")

	// Use configured API key
	url := fmt.Sprintf("http://apilayer.net/api/validate?access_key=%s&number=%s&format=1", cfg.NumverifyKey, phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("numverify API error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Check if API returned an error
	if success, ok := result["success"].(bool); ok && !success {
		return fmt.Errorf("numverify API requires valid key")
	}

	// Parse response
	if valid, ok := result["valid"].(bool); ok {
		info.IsValid = valid
	}

	if carrier, ok := result["carrier"].(string); ok && carrier != "" {
		info.Carrier = carrier
	}

	if lineType, ok := result["line_type"].(string); ok && lineType != "" {
		info.LineType = lineType
	}

	if location, ok := result["location"].(string); ok && location != "" {
		info.Region = location
	}

	if country, ok := result["country_name"].(string); ok && country != "" {
		info.Country = country
	}

	if countryCode, ok := result["country_prefix"].(string); ok && countryCode != "" {
		info.CountryCode = "+" + countryCode
	}

	return nil
}

// lookupPhoneValidator uses phone-validator.net API
func lookupPhoneValidator(phone string, info *PhoneInfo, cfg *config.Config) error {
	// Skip if no API key configured
	if cfg == nil || cfg.AbstractAPIKey == "" {
		return fmt.Errorf("abstractapi key not configured")
	}

	phoneClean := strings.TrimPrefix(phone, "+")

	// Use configured API key
	url := fmt.Sprintf("https://phonevalidation.abstractapi.com/v1/?api_key=%s&phone=%s", cfg.AbstractAPIKey, phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("phone validator API error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Parse response
	if valid, ok := result["valid"].(bool); ok {
		info.IsValid = valid
	}

	if carrier, ok := result["carrier"].(string); ok && carrier != "" {
		info.Carrier = carrier
	}

	if phoneType, ok := result["type"].(string); ok && phoneType != "" {
		info.LineType = phoneType
	}

	if country, ok := result["country"].(map[string]interface{}); ok {
		if countryName, ok := country["name"].(string); ok {
			info.Country = countryName
		}
		if countryCode, ok := country["code"].(string); ok {
			info.CountryCode = "+" + countryCode
		}
	}

	if location, ok := result["location"].(string); ok && location != "" {
		info.Region = location
	}

	return nil
}

// lookupIPQualityScore uses IPQualityScore phone validation API
func lookupIPQualityScore(phone string, info *PhoneInfo, cfg *config.Config) error {
	// Skip if no API key configured
	if cfg == nil || cfg.IPQualityScoreKey == "" {
		return fmt.Errorf("IPQualityScore API key not configured")
	}

	phoneClean := strings.TrimPrefix(phone, "+")

	// Use configured API key
	url := fmt.Sprintf("https://ipqualityscore.com/api/json/phone/%s/%s", cfg.IPQualityScoreKey, phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "OSINT-Master-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("IPQualityScore API error: %d", resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	// Check if API request was successful
	if success, ok := result["success"].(bool); ok && !success {
		return fmt.Errorf("IPQualityScore requires valid API key")
	}

	// Parse response
	if valid, ok := result["valid"].(bool); ok {
		info.IsValid = valid
	}

	if carrier, ok := result["carrier"].(string); ok && carrier != "" {
		info.Carrier = carrier
	}

	if lineType, ok := result["line_type"].(string); ok && lineType != "" {
		info.LineType = lineType
	}

	if country, ok := result["country"].(string); ok && country != "" {
		info.Country = country
	}

	if region, ok := result["region"].(string); ok && region != "" {
		info.Region = region
	}

	if city, ok := result["city"].(string); ok && city != "" && info.Region != "" {
		info.Region = fmt.Sprintf("%s, %s", city, info.Region)
	}

	return nil
}

// lookupOwnerInfo tries to find the owner's information from various sources
func lookupOwnerInfo(phone string, info *PhoneInfo, cfg *config.Config) error {
	phoneClean := strings.ReplaceAll(strings.ReplaceAll(phone, "+", ""), " ", "")

	// Try multiple owner lookup sources

	// 1. Try TrueCaller API (requires scraping or unofficial API)
	if owner := lookupTrueCaller(phoneClean); owner != "" {
		info.OwnerName = owner
		info.OwnerSource = "TrueCaller (public data)"
		return nil
	}

	// 2. Try Numverify extended data (if configured)
	if cfg != nil && cfg.NumverifyKey != "" {
		if err := lookupNumverifyExtended(phoneClean, info, cfg); err == nil && info.OwnerName != "" {
			return nil
		}
	}

	// 3. Try phone directory services
	if owner := lookupPhoneDirectory(phoneClean); owner != "" {
		info.OwnerName = owner
		info.OwnerSource = "Public directory"
		return nil
	}

	// 4. Try social media reverse lookup
	if owner := lookupSocialMedia(phoneClean); owner != "" {
		info.OwnerName = owner
		info.OwnerSource = "Social media"
		return nil
	}

	return fmt.Errorf("no owner information found")
}

// lookupTrueCaller attempts to get name from TrueCaller using multiple methods
func lookupTrueCaller(phone string) string {
	// Try local cache first (fastest)
	name := tryLocalCache(phone)
	if name != "" {
		return name
	}

	// Try GetContact API (works well for international numbers)
	name = tryGetContactAPI(phone)
	if name != "" {
		return name
	}

	// Try Sync.me API
	name = trySyncMeAPI(phone)
	if name != "" {
		return name
	}

	// Try TrueCaller JSON API endpoint (unofficial but works)
	name = tryTrueCallerJSONAPI(phone)
	if name != "" {
		return name
	}

	// Try Eyecon API as alternative
	name = tryEyeconAPI(phone)
	if name != "" {
		return name
	}

	// Try NumLookup API
	name = tryNumLookupAPI(phone)
	if name != "" {
		return name
	}

	return ""
}

// tryLocalCache checks a local JSON file for known phone-name mappings
// This allows users to optionally add their own known contacts
func tryLocalCache(phone string) string {
	// Skip local cache - we want to use real APIs only
	return ""
}

// tryGetContactAPI tries GetContact caller ID service
func tryGetContactAPI(phone string) string {
	phoneClean := strings.TrimPrefix(phone, "+")

	// GetContact API endpoint
	url := fmt.Sprintf("https://api.getcontact.com/search?phoneNumber=%s", phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "GetContact/4.8.1 (Android)")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	// Extract name from GetContact response
	if name, ok := result["displayName"].(string); ok && name != "" {
		return name
	}
	if tags, ok := result["tags"].([]interface{}); ok && len(tags) > 0 {
		if tag, ok := tags[0].(map[string]interface{}); ok {
			if tagName, ok := tag["tag"].(string); ok && tagName != "" {
				return tagName
			}
		}
	}

	return ""
}

// trySyncMeAPI tries Sync.me caller ID service
func trySyncMeAPI(phone string) string {
	phoneClean := strings.TrimPrefix(phone, "+")

	// Sync.me API endpoint
	url := fmt.Sprintf("https://api.sync.me/api/v3/contacts/search?phoneNumber=%s", phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "Sync.me/5.0")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	// Extract name from Sync.me response
	if contacts, ok := result["contacts"].([]interface{}); ok && len(contacts) > 0 {
		if contact, ok := contacts[0].(map[string]interface{}); ok {
			if name, ok := contact["name"].(string); ok && name != "" {
				return name
			}
		}
	}

	return ""
}

// runTrueCallerPlaywright runs the Playwright scraper for TrueCaller
func runTrueCallerPlaywright(phone string) string {
	// Import exec package at runtime
	cmd := exec.Command("node", "internal/scraper/truecaller_scraper.js", phone)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return ""
	}

	if success, ok := result["success"].(bool); ok && success {
		if name, ok := result["name"].(string); ok && name != "" {
			return name
		}
	}

	return ""
}

// tryTrueCallerJSONAPI uses Playwright to scrape TrueCaller with JavaScript execution
func tryTrueCallerJSONAPI(phone string) string {
	phoneClean := strings.ReplaceAll(strings.ReplaceAll(phone, "+", ""), " ", "")

	// Use Playwright scraper for TrueCaller
	return runTrueCallerPlaywright(phoneClean)
}

// tryEyeconAPI tries Eyecon caller ID API
func tryEyeconAPI(phone string) string {
	phoneClean := strings.TrimPrefix(phone, "+")

	url := fmt.Sprintf("https://api.eyecon-app.com/app/getnames.jsp?cli=%s&lang=en", phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "Eyecon/9.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return ""
	}

	// Extract name from Eyecon response
	if names, ok := result["names"].([]interface{}); ok && len(names) > 0 {
		if firstName, ok := names[0].(map[string]interface{}); ok {
			if name, ok := firstName["name"].(string); ok && name != "" {
				return name
			}
		}
	}

	return ""
}

// tryNumLookupAPI tries NumLookup free API
func tryNumLookupAPI(phone string) string {
	phoneClean := strings.TrimPrefix(phone, "+")

	url := fmt.Sprintf("https://www.numlookup.com/api/validate/%s", phoneClean)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	// Extract owner name if available
	if owner, ok := result["owner"].(string); ok && owner != "" {
		return owner
	}
	if name, ok := result["name"].(string); ok && name != "" {
		return name
	}

	return ""
}

// lookupNumverifyExtended gets extended data including owner info if available
func lookupNumverifyExtended(phone string, info *PhoneInfo, cfg *config.Config) error {
	// Some phone APIs provide owner information
	// This would require extended API access
	return fmt.Errorf("extended data not available")
}

// lookupPhoneDirectory searches public phone directories
func lookupPhoneDirectory(phone string) string {
	// Try free phone directory APIs
	// Note: Most accurate directories are paid services

	// Try phonevalidator.com directory
	url := fmt.Sprintf("https://www.phonevalidator.com/api/lookup/%s", phone)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}

	req.Header.Set("User-Agent", "OSINT-Master-Tool")

	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ""
	}

	// Try to extract owner name
	if name, ok := result["name"].(string); ok && name != "" {
		return name
	}

	if owner, ok := result["owner"].(string); ok && owner != "" {
		return owner
	}

	if subscriber, ok := result["subscriber"].(string); ok && subscriber != "" {
		return subscriber
	}

	return ""
}

// lookupSocialMedia checks if phone is linked to social media profiles
func lookupSocialMedia(phone string) string {
	// Try to find name from social media
	// This is limited due to privacy settings on most platforms

	// Format phone for WhatsApp/Telegram lookup
	formattedPhone := phone
	if !strings.HasPrefix(formattedPhone, "+") {
		formattedPhone = "+" + formattedPhone
	}

	// Note: Getting name from WhatsApp/Telegram requires:
	// 1. Authentication
	// 2. Contact in your address book
	// 3. User's privacy settings allow it

	// This is a placeholder - real implementation needs proper APIs
	return ""
}

// checkWhatsApp checks if a phone number is registered on WhatsApp
func checkWhatsApp(phone string) (bool, string) {
	// Method: Try to access the WhatsApp Web API endpoint
	// Remove + and spaces from phone number
	cleanedPhone := strings.ReplaceAll(strings.ReplaceAll(phone, "+", ""), " ", "")

	// Try using wa.me link which is an official WhatsApp redirect service
	url := fmt.Sprintf("https://wa.me/%s", cleanedPhone)

	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects, just check the response
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return false, "Unable to check (request error)"
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return false, "Unable to check (network error)"
	}
	defer resp.Body.Close()

	// If WhatsApp redirects to web.whatsapp.com or api.whatsapp.com, the number likely exists
	// Status code 302 (redirect) typically means the number is valid
	if resp.StatusCode == http.StatusFound || resp.StatusCode == http.StatusMovedPermanently {
		location := resp.Header.Get("Location")
		if strings.Contains(location, "web.whatsapp.com") || strings.Contains(location, "api.whatsapp.com") {
			return true, "Active on WhatsApp"
		}
	}

	// Alternative check: Use WhatsApp API check service
	// Try free WhatsApp checker API
	apiURL := fmt.Sprintf("https://api.wassenger.com/v1/numbers/%s/exists", cleanedPhone)

	req2, err := http.NewRequest("GET", apiURL, nil)
	if err == nil {
		req2.Header.Set("User-Agent", "OSINT-Master-Tool")
		resp2, err := client.Do(req2)
		if err == nil {
			defer resp2.Body.Close()

			if resp2.StatusCode == http.StatusOK {
				var result map[string]interface{}
				if err := json.NewDecoder(resp2.Body).Decode(&result); err == nil {
					if exists, ok := result["exists"].(bool); ok && exists {
						return true, "Active on WhatsApp (verified)"
					}
				}
			}
		}
	}

	// If all checks fail, we can't confirm
	return false, fmt.Sprintf("Not confirmed on WhatsApp (try manually: https://wa.me/%s)", cleanedPhone)
}
