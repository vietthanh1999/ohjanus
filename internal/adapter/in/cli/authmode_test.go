package cli

import (
	"testing"
)

func TestCheckAuthMode(t *testing.T) {
	t.Setenv("JANUS_ALLOW_NO_AUTH", "")
	if err := checkAuthMode("token"); err != nil {
		t.Fatalf("token mode: %v", err)
	}
	if err := checkAuthMode("none"); err == nil {
		t.Fatal("none mode without opt-in: want refusal")
	}
	t.Setenv("JANUS_ALLOW_NO_AUTH", "1")
	if err := checkAuthMode("none"); err != nil {
		t.Fatalf("none mode with opt-in: %v", err)
	}
	t.Setenv("JANUS_ALLOW_NO_AUTH", "0")
	if err := checkAuthMode("none"); err == nil {
		t.Fatal("none mode with opt-out value: want refusal")
	}
}
