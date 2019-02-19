package testutil

import (
	"strings"
	"testing"
)

//AssertString asserts string values and prints the expected and received values if failed
func AssertString(t *testing.T, message string, expected string, got string) {
	expected = trimString(expected)
	got = trimString(got)
	if expected != got {
		t.Errorf("%s: expected '%s' but got '%s' instead", message, expected, got)
	}
}

func trimString(s string) string {
	return strings.TrimRight(s, "\n")
}
