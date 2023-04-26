package cmd

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestExecuteDecrypt(t *testing.T) {
	var tp = struct {
		URI         string
		testContent string
	}{
		"/decrypt",
		"test",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertString(t, "Incorrect Method", "POST", r.Method)
		assertString(t, "Incorrect URI call", tp.URI, r.RequestURI)

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)

		assertString(t, "Incorrect Content received", tp.testContent, buf.String())
		fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	dp.source = ts.URL
	dp.value = tp.testContent
	err := ExecuteDecrypt()

	if err != nil {
		t.Error("Decrypt failed with: ", err)
	}
}

func ExampleExecuteDecrypt() {
	var tp = struct {
		URI         string
		testContent string
	}{
		"/decrypt",
		"test",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	dp.source = ts.URL
	dp.value = tp.testContent
	ExecuteDecrypt()
	// Output: test
}

// AssertString asserts string values and prints the expected and received values if failed.
func assertString(t *testing.T, message string, expected string, got string) {
	expected = trimString(expected)
	got = trimString(got)
	if expected != got {
		t.Errorf("%s: expected '%s' but got '%s' instead", message, expected, got)
	}
}

func trimString(s string) string {
	return strings.TrimRight(s, "\n")
}
