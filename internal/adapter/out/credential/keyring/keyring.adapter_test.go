package keyring

import (
	"context"
	"testing"

	"github.com/zalando/go-keyring"
)

func TestSplit(t *testing.T) {
	r := New("janus")
	cases := []struct {
		ref         string
		service     string
		key         string
		expectError bool
	}{
		{"keyring://janus/analytics", "janus", "analytics", false},
		{"keyring://svc/a/b", "svc", "a/b", false},
		{"keyring:analytics", "janus", "analytics", false},
		{"keyring://", "", "", true},
		{"keyring:///x", "", "", true},
		{"keyring://svc/", "", "", true},
		{"keyring:", "", "", true},
		{"env:X", "", "", true},
	}
	for _, tc := range cases {
		svc, key, err := r.split(tc.ref)
		if tc.expectError {
			if err == nil {
				t.Errorf("split(%q) should fail", tc.ref)
			}
			continue
		}
		if err != nil {
			t.Errorf("split(%q): %v", tc.ref, err)
			continue
		}
		if svc != tc.service || key != tc.key {
			t.Errorf("split(%q) = %q/%q, want %q/%q", tc.ref, svc, key, tc.service, tc.key)
		}
	}
}

func TestRoundtrip(t *testing.T) {
	const service = "ohjanus-test"
	const key = "roundtrip"
	const secret = "postgres://u:p@localhost/db"
	if err := keyring.Set(service, key, secret); err != nil {
		t.Skipf("keychain unavailable: %v", err)
	}
	t.Cleanup(func() { _ = keyring.Delete(service, key) })

	r := New("janus")
	creds, err := r.Resolve(context.Background(), "keyring://"+service+"/"+key)
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	if creds.DSN != secret {
		t.Errorf("DSN mismatch")
	}

	if _, err := r.Resolve(context.Background(), "keyring://"+service+"/missing"); err == nil {
		t.Error("missing key should fail")
	}
}
