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

	fullContent := header + content

	// Write to file with proper permissions
	err := os.WriteFile(filename, []byte(fullContent), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %v", err)
	}

	return nil
}
