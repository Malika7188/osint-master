package output

import (
	"fmt"
	"os"
	"time"
)

// SaveToFile saves the result string to a file
func SaveToFile(filename string, content string) error {
	if filename == "" {
		return fmt.Errorf("filename cannot be empty")
	}

	// Add timestamp header
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// header := fmt.Sprintf("=== OSINT Master Report ===\n")
	header := fmt.Sprintf("Generated: %s\n", timestamp)
	// header += fmt.Sprintf("===========================\n\n")

	