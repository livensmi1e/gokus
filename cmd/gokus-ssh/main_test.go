package main

import "testing"

func TestEnvOrDefault(t *testing.T) {
	const name = "GOKUS_TEST_VALUE"

	t.Run("uses fallback when unset", func(t *testing.T) {
		t.Setenv(name, "")
		if got := envOrDefault(name, "fallback"); got != "fallback" {
			t.Fatalf("envOrDefault() = %q, want %q", got, "fallback")
		}
	})

	t.Run("uses environment value", func(t *testing.T) {
		t.Setenv(name, "configured")
		if got := envOrDefault(name, "fallback"); got != "configured" {
			t.Fatalf("envOrDefault() = %q, want %q", got, "configured")
		}
	})
}
