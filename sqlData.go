package utils

import (
	"regexp"
	"strings"
)

func SafeColumnName(value string) string {
	match := regexp.MustCompile("[^a-z|A-Z|0-9|_|$]")
	return match.ReplaceAllString(value, "")
}

func Safe(value string) string {
	return strings.ReplaceAll(value, "'", "''")
}
