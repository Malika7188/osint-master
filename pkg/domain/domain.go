package domain

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
)

// Subdomain represents information about a subdomain
type Subdomain struct {
	Name        string
	IP          string
	SSLCert     string
	IsTakeover  bool
	TakeoverMsg string
}

// DomainInfo holds all information about a domain
type DomainInfo struct {
	MainDomain string
	Subdomains []Subdomain
}

// EnumerateDomain enumerates subdomains and checks for takeover risks
func EnumerateDomain(domain string) (string, error) {
	if domain == "" {
		return "", fmt.Errorf("domain cannot be empty")
	}

	// Remove protocol if present
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimSuffix(domain, "/")

	fmt.Println("\nEnumerating subdomains... This may take a moment.")

	// Get subdomains from Certificate Transparency logs
	subdomains, err := getSubdomainsFromCrtSh(domain)
	if err != nil {
		return "", fmt.Errorf("failed to enumerate subdomains: %v", err)
	}

	// Check each subdomain for details
	domainInfo := &DomainInfo{
		MainDomain: domain,
		Subdomains: make([]Subdomain, 0),
	}

	for _, sub := range subdomains {
		info := checkSubdomain(sub)
		domainInfo.Subdomains = append(domainInfo.Subdomains, info)
	}

	result := formatDomainInfo(domainInfo)
	return result, nil
}

// getSubdomainsFromCrtSh queries crt.sh for subdomains via Certificate Transparency
func getSubdomainsFromCrtSh(domain string) ([]string, error) {
	url := fmt.Sprintf("https://crt.sh/?q=%%25.%s&output=json", domain)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("crt.sh returned status: %d", resp.StatusCode)
	}

	var certs []struct {
		NameValue string `json:"name_value"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&certs); err != nil {
		return nil, err
	}

	// Extract unique subdomains
	subdomainMap := make(map[string]bool)
	for _, cert := range certs {
		names := strings.Split(cert.NameValue, "\n")
		for _, name := range names {
			name = strings.TrimSpace(name)
			// Skip wildcards and duplicates
			if !strings.HasPrefix(name, "*") && name != "" {
				subdomainMap[name] = true
			}
		}
	}

	// Convert map to slice
	subdomains := make([]string, 0, len(subdomainMap))
	for sub := range subdomainMap {
		subdomains = append(subdomains, sub)
	}

	// Limit to first 10 for demo purposes
	if len(subdomains) > 10 {
		subdomains = subdomains[:10]
	}

	return subdomains, nil
}

// checkSubdomain checks a subdomain for IP, SSL, and takeover risks
func checkSubdomain(subdomain string) Subdomain {
	info := Subdomain{
		Name:   subdomain,
		IP:     "Unknown",
		SSLCert: "Not checked",
		IsTakeover: false,
	}

	// Resolve IP address
	ips, err := net.LookupIP(subdomain)
	if err == nil && len(ips) > 0 {
		info.IP = ips[0].String()
	}

	// Check SSL certificate
	info.SSLCert = checkSSLCert(subdomain)

	// Check for potential subdomain takeover
	info.IsTakeover, info.TakeoverMsg = checkTakeoverRisk(subdomain)

	return info
}
