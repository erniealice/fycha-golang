package doctemplate

import (
	"fmt"
	"regexp"
	"strings"
)

// ReplacePlaceholders replaces placeholders in an XML string with values from a data map.
// It skips special loop markers like {{#...}} and {{/...}}.
func ReplacePlaceholders(xmlContent string, data map[string]any) string {
	re := regexp.MustCompile(`{{(.*?)}}`)

	return re.ReplaceAllStringFunc(xmlContent, func(match string) string {
		// Trim the placeholder delimiters
		trimmedPlaceholder := strings.TrimSpace(match[2 : len(match)-2])

		// Skip loop markers
		if strings.HasPrefix(trimmedPlaceholder, "#") || strings.HasPrefix(trimmedPlaceholder, "/") {
			return match // Return the original marker
		}

		keys := parsePlaceholder(trimmedPlaceholder)
		if value, found := getNestedValue(keys, data); found {
			return fmt.Sprintf("%v", value)
		}

		// If the placeholder is not found, return the original placeholder
		return match
	})
}

// parsePlaceholder splits a placeholder key into its constituent parts for nested lookups.
func parsePlaceholder(placeholder string) []string {
	return strings.Split(placeholder, ".")
}

// getNestedValue retrieves a value from a nested map structure using a slice of keys.
// It correctly handles map[string]any types.
func getNestedValue(keys []string, data map[string]any) (any, bool) {
	if len(keys) == 0 {
		return nil, false
	}

	currentData := data
	for i, key := range keys {
		value, found := currentData[key]
		if !found {
			return nil, false
		}

		// If it's the last key, we've found our value.
		if i == len(keys)-1 {
			return value, true
		}

		// If it's not the last key, we expect the next level to be a map.
		if nestedData, isMap := value.(map[string]any); isMap {
			currentData = nestedData
		} else {
			// Path is not fully traversable
			return nil, false
		}
	}

	return nil, false
}
