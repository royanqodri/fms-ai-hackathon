package util

import (
	"strings"
)

// IsEmpty function will check the string is empty or not
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

func CleanString(input string, escStr ...string) string {
	cleaned := input

	for _, s := range escStr {
		cleaned = strings.ReplaceAll(input, s, "")
	}

	// remove leading and trailing whitespace
	cleaned = strings.TrimSpace(cleaned)

	return cleaned
}
