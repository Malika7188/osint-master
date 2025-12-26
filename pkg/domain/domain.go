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
