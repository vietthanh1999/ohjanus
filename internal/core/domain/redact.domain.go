package domain

import "regexp"

// RedactPattern replaces PII matches with a placeholder.
type RedactPattern struct {
	Name        string
	Regex       *regexp.Regexp
	Replacement string
}
