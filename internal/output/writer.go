package output

import (
	"fmt"
	"os"
	"time"
)

// SaveToFile saves the result string to a file
func SaveToFile(filename string, content string) error {
	