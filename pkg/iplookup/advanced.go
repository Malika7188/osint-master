package iplookup

import (
	"fmt"
	"strings"
	"time"
)

// AdvancedLookupIP performs enhanced IP address analysis
func AdvancedLookupIP(ip string) (string, error) {
	var result strings.Builder
	// result.WriteString("⚠️  ADVANCED MODE: Enhanced IP address analysis\n")
	// result.WriteString("This mode performs additional checks and takes longer\n")
	// result.WriteString(strings.Repeat("=", 70) + "\n\n")

	// Perform standard lookup first
	startTime := time.Now()
	standardResult, err := LookupIP(ip)
	if err != nil {
		return "", err
	}

	result.WriteString(standardResult)
	result.WriteString("\n" + strings.Repeat("-", 70) + "\n")
	result.WriteString("ADVANCED CHECKS:\n")
	result.WriteString(strings.Repeat("-", 70) + "\n\n")

	result.WriteString("Security & Reputation Checks:\n")
	result.WriteString(fmt.Sprintf("  - AbuseIPDB: https://www.abuseipdb.com/check/%s\n", ip))
	result.WriteString(fmt.Sprintf("  - VirusTotal: https://www.virustotal.com/gui/ip-address/%s\n", ip))
	result.WriteString(fmt.Sprintf("  - IPVoid: https://www.ipvoid.com/ip-blacklist-check/\n"))
	result.WriteString(fmt.Sprintf("  - Shodan: https://www.shodan.io/host/%s\n", ip))
	result.WriteString(fmt.Sprintf("  - Censys: https://search.censys.io/hosts/%s\n", ip))

	result.WriteString("\nGeolocation & Network Info:\n")
	result.WriteString(fmt.Sprintf("  - IPInfo.io: https://ipinfo.io/%s\n", ip))
	result.WriteString(fmt.Sprintf("  - MaxMind: https://www.maxmind.com/en/geoip2-precision-demo\n"))
	result.WriteString(fmt.Sprintf("  - IP2Location: https://www.ip2location.com/%s\n", ip))

	result.WriteString("\nReverse DNS & WHOIS:\n")
	result.WriteString(fmt.Sprintf("  - ARIN WHOIS: https://search.arin.net/rdap/?query=%s\n", ip))
	result.WriteString(fmt.Sprintf("  - RIPE: https://apps.db.ripe.net/db-web-ui/query?searchtext=%s\n", ip))
	result.WriteString(fmt.Sprintf("  - MXToolbox: https://mxtoolbox.com/SuperTool.aspx?action=ptr%%3A%s\n", ip))

	result.WriteString("\nPort Scanning & Services:\n")
	result.WriteString("  ⚠️  Only scan IPs you own or have permission to scan\n")
	result.WriteString(fmt.Sprintf("  - Shodan Scan: https://www.shodan.io/host/%s\n", ip))
	result.WriteString(fmt.Sprintf("  - Censys Scan: https://search.censys.io/hosts/%s\n", ip))

	result.WriteString("\nThreat Intelligence:\n")
	result.WriteString(fmt.Sprintf("  - Talos Intelligence: https://www.talosintelligence.com/reputation_center/lookup?search=%s\n", ip))
	result.WriteString(fmt.Sprintf("  - AlienVault OTX: https://otx.alienvault.com/indicator/ip/%s\n", ip))
	result.WriteString(fmt.Sprintf("  - GreyNoise: https://viz.greynoise.io/ip/%s\n", ip))

	elapsed := time.Since(startTime)
	result.WriteString("\n" + strings.Repeat("=", 70) + "\n")
	result.WriteString(fmt.Sprintf("⏱️  Advanced search completed in %.2f seconds\n", elapsed.Seconds()))
	result.WriteString("Note: Advanced mode provides comprehensive security and threat analysis\n")

	return result.String(), nil
}
