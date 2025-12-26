package domain

import (
	"fmt"
	"strings"
	"time"
)

// AdvancedEnumerateDomain performs enhanced domain enumeration with additional analysis
func AdvancedEnumerateDomain(domain string) (string, error) {
	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced domain analysis\n")
	// result.WriteString("This mode performs additional checks and takes longer\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard enumeration first
	startTime := time.Now()
	standardResult, err := EnumerateDomain(domain)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")

	// Clean domain
	cleanDomain := strings.TrimPrefix(domain, "http://")
	cleanDomain = strings.TrimPrefix(cleanDomain, "https://")
	cleanDomain = strings.TrimSuffix(cleanDomain, "/")

	result.WriteString("Additional Subdomain Enumeration Tools:\n")
	result.WriteString("  - Sublist3r: Automated subdomain enumeration tool\n")
	result.WriteString("  - Amass: In-depth DNS enumeration\n")
	result.WriteString("  - Subfinder: Fast passive subdomain enumeration\n")
	result.WriteString(fmt.Sprintf("  - SecurityTrails: https://securitytrails.com/domain/%s/dns\n", cleanDomain))
	result.WriteString(fmt.Sprintf("  - DNSDumpster: https://dnsdumpster.com/\n"))

	result.WriteString("\nDomain Intelligence:\n")
	result.WriteString(fmt.Sprintf("  - WHOIS: https://who.is/whois/%s\n", cleanDomain))
	result.WriteString(fmt.Sprintf("  - Domain History: https://whoisrequest.com/history/%s\n", cleanDomain))
	result.WriteString(fmt.Sprintf("  - Wayback Machine: https://web.archive.org/web/*/%s\n", cleanDomain))

	result.WriteString("\nSecurity & Reputation:\n")
	result.WriteString(fmt.Sprintf("  - VirusTotal: https://www.virustotal.com/gui/domain/%s\n", cleanDomain))
	result.WriteString(fmt.Sprintf("  - URLVoid: https://www.urlvoid.com/scan/%s\n", cleanDomain))
	result.WriteString(fmt.Sprintf("  - Google Safe Browsing: Check at transparencyreport.google.com\n"))
