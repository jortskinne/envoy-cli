package parser

import "strings"

// DefaultSecretPatterns are key substrings that indicate sensitive values.
var DefaultSecretPatterns = []string{
	"SECRET",
	"PASSWORD",
	"PASSWD",
	"TOKEN",
	"API_KEY",
	"PRIVATE_KEY",
	"CREDENTIALS",
	"AUTH",
}

// MaskOptions configures masking behaviour.
type MaskOptions struct {
	Patterns    []string
	MaskString  string
	RevealChars int // number of trailing chars to reveal (0 = fully masked)
}

// DefaultMaskOptions returns sensible defaults.
func DefaultMaskOptions() MaskOptions {
	return MaskOptions{
		Patterns:    DefaultSecretPatterns,
		MaskString:  "****",
		RevealChars: 0,
	}
}

// IsSensitive reports whether the key matches any secret pattern.
func IsSensitive(key string, patterns []string) bool {
	upper := strings.ToUpper(key)
	for _, p := range patterns {
		if strings.Contains(upper, strings.ToUpper(p)) {
			return true
		}
	}
	return false
}

// MaskValue masks a value according to the provided options.
func MaskValue(value string, opts MaskOptions) string {
	if opts.RevealChars > 0 && len(value) > opts.RevealChars {
		return opts.MaskString + value[len(value)-opts.RevealChars:]
	}
	return opts.MaskString
}

// MaskEntry returns the display value for an entry, masking if sensitive.
func MaskEntry(entry Entry, opts MaskOptions) string {
	if IsSensitive(entry.Key, opts.Patterns) {
		return MaskValue(entry.Value, opts)
	}
	return entry.Value
}
