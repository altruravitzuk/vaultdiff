package output

import (
	"fmt"
	"strings"
)

// ParseFormat converts a string to a Format value, returning an error for
// unrecognised format names.
func ParseFormat(s string) (Format, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "text", "":
		return FormatText, nil
	case "json":
		return FormatJSON, nil
	case "yaml", "yml":
		return FormatYAML, nil
	default:
		return "", fmt.Errorf("unsupported output format %q: must be one of text, json, yaml", s)
	}
}

// SupportedFormats returns a human-readable list of valid format names.
func SupportedFormats() string {
	return "text, json, yaml"
}
