package validator

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

// ValidateIP validates if a string is a valid IP address
func ValidateIP(ip string) error {
	