package utils

import (
	"regexp"
	"strings"
)

// SafeColumnName sanitizes a given string to be used as a column name in a database.
// It removes any characters that are not letters, numbers, underscores, or dollar signs.
// The sanitized string is returned.
func SafeColumnName(value string) string {
	// Create a regular expression pattern that matches any characters that are not letters, numbers, underscores, or dollar signs.
	pattern := "[^a-zA-Z0-9_$]"

	// Compile the regular expression pattern.
	match := regexp.MustCompile(pattern)

	// Replace any matches in the input string with an empty string to remove them.
	return match.ReplaceAllString(value, "")
}

// Safe sanitizes a given string to be used in a SQL query as a value.
// It replaces any occurrences of the single quote character with two single quote characters to escape them.
// The sanitized string is returned.
func Safe(value string) string {
	// Replace any occurrences of the single quote character with two single quote characters to escape them.
	return strings.ReplaceAll(value, "'", "''")
}
