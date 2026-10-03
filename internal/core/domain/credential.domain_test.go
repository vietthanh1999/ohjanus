package domain

import "testing"

func TestRedactedDSN(t *testing.T) {
	cases := []struct{ in, want string }{
		{"postgres://u:p@h/db", "postgres://u:***@h/db"},
		{"postgres://u:p@ss@h/db", "postgres://u:***@h/db"},
		{"postgres://u@h/db", "postgres://u@h/db"},
		{"postgres://u:p@h/db?sslmode=disable", "postgres://u:***@h/db?sslmode=disable"},
		{"user:pass@tcp(h)/db", "user:***@tcp(h)/db"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := (Credentials{DSN: tc.in}).RedactedDSN(); got != tc.want {
			t.Errorf("RedactedDSN(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}
