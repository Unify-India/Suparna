package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// HandleError logs the given error and exits the program.
func HandleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

// GetHomeDirectory returns the user's home directory path.
func GetHomeDirectory() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return homeDir, nil
}

// JoinPaths joins multiple path segments into a single path string.
func JoinPaths(parts ...string) string {
	return filepath.Join(parts...)
}

// Capitalize capitalizes the first letter of a string.
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// ToTitleCase converts a string to title case (e.g., "hello world" becomes "Hello World").
func ToTitleCase(s string) string {
	return strings.Title(strings.ToLower(s))
}

// GetCurrentTimestamp returns the current timestamp.
func GetCurrentTimestamp() time.Time {
	return time.Now()
}

// LogInfo logs an informational message.
func LogInfo(msg string) {
	log.Printf("INFO: %s", msg)
}

// LogWarning logs a warning message.
func LogWarning(msg string) {
	log.Printf("WARNING: %s", msg)
}

// LogError logs an error message.
func LogError(err error) {
	log.Printf("ERROR: %v", err)
}

// SliceContains checks if a slice contains a specific value.
func SliceContains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// MapKeysToString converts a map's keys to a slice of strings.
func MapKeysToString(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
