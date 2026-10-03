package regex

import (
	"strings"

	"github.com/vietthanh1999/ohjanus/internal/core/domain"
	"github.com/vietthanh1999/ohjanus/internal/core/port/out"
)

// Redactor applies regex patterns and per-column masking.
type Redactor struct {
	patterns []domain.RedactPattern
	columns  map[string]struct{}
}

var _ out.Redactor = (*Redactor)(nil)

// New builds a redactor. Columns match case-insensitively.
func New(patterns []domain.RedactPattern, columns []string) *Redactor {
	cols := make(map[string]struct{}, len(columns))
	for _, c := range columns {
		cols[strings.ToLower(c)] = struct{}{}
	}
	return &Redactor{patterns: patterns, columns: cols}
}

// RedactString replaces every pattern match with its replacement.
func (r *Redactor) RedactString(s string) string {
	for _, p := range r.patterns {
		s = p.Regex.ReplaceAllString(s, p.Replacement)
	}
	return s
}

// RedactValue masks the value when the column is sensitive.
func (r *Redactor) RedactValue(column, value string) string {
	if _, ok := r.columns[strings.ToLower(column)]; ok {
		return "[REDACTED]"
	}
	return r.RedactString(value)
}
