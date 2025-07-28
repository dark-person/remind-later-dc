// Package to extract time.Duration from string.
package timeparse

import (
	"regexp"
	"strconv"
)

// Extract duration from given string.
// The string should be in strict format, e.g. "1h2m", "2m3s".
// If string cannot be parsed, then it will return all zero values with ok=false.
func ExtractDuration(input string) (hours, minutes, seconds int, ok bool) {
	// Regex pattern for strings like "1h2m3s", "2m3s", "1h2m", "2m"
	pattern := `^(?:(\d+)h)?(?:(\d+)m)?(?:(\d+)s)?$`
	re := regexp.MustCompile(pattern)

	// Check if the input matches the pattern
	if !re.MatchString(input) || input == "" {
		return 0, 0, 0, false
	}

	// Extract the captured groups
	matches := re.FindStringSubmatch(input)

	// Ensure at least one component is present
	if matches[1] == "" && matches[2] == "" && matches[3] == "" {
		return 0, 0, 0, false
	}

	// Convert captured strings to integers, defaulting to 0 if not present
	hours = 0
	if matches[1] != "" {
		hours, _ = strconv.Atoi(matches[1])
	}
	minutes = 0
	if matches[2] != "" {
		minutes, _ = strconv.Atoi(matches[2])
	}
	seconds = 0
	if matches[3] != "" {
		seconds, _ = strconv.Atoi(matches[3])
	}

	return hours, minutes, seconds, true
}
