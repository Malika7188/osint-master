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
