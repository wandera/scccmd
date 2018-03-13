package client

import (
	"testing"
	"fmt"
	"net/http/httptest"
	"net/http"
	"strings"
)

func TestNewClient(t *testing.T) {
	var tp = struct {
		application string
		profile     string
		label       string
		URI         string
	}{
		"service",
		"profile",
		"master",
		"http://localhost:8080",
	}

	c := NewClient(Config{
		URI:         tp.URI,
		Application: tp.application,
		Profile:     tp.profile,
		Label:       tp.label,
	})

	assertString(t, "Incorrect URI", tp.URI, c.Config().URI)
	assertString(t, "Incorrect Application", tp.application, c.Config().Application)
	assertString(t, "Incorrect Profile", tp.profile, c.Config().Profile)
	assertString(t, "Incorrect Label", tp.label, c.Config().Label)
}

func TestClient_FetchFile(t *testing.T) {
	var tp = struct {
		application string
		profile     string
		label       string
		URI         string
		testContent string
		fileName    string
	}{
		"service",
		"profile",
		"master",
		"/service/profile/master/File",
		"test",
		"File",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertString(t, "Incorrect URI call", tp.URI, r.RequestURI)
		fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	cont, err := NewClient(Config{
		URI:         ts.URL,
		Application: tp.application,
		Profile:     tp.profile,
		Label:       tp.label,
	}).FetchFile(tp.fileName)

	if err != nil {
		t.Error("FetchFile failed with: ", err)
	}

	assertString(t, "Content mismatch", tp.testContent, string(cont))
}

func TestClient_FetchAsYAML(t *testing.T) {
	var tp = struct {
		application string
		profile     string
		label       string
		URI         string
		testContent string
		fileName    string
	}{
		"service",
		"profile",
		"master",
		"/master/service-profile.yml",
		"test",
		"File",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assertString(t, "Incorrect URI call", tp.URI, r.RequestURI)
		fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	cont, err := NewClient(Config{
		URI:         ts.URL,
		Application: tp.application,
		Profile:     tp.profile,
		Label:       tp.label,
	}).FetchAsYAML()

	if err != nil {
		t.Error("FetchFile failed with: ", err)
	}

	assertString(t, "Content mismatch", tp.testContent, cont)
}

func assertString(t *testing.T, message string, expected string, got string) {
	expected = trimString(expected)
	got = trimString(got)
	if expected != got {
		t.Error(fmt.Sprintf("%s: expected %s but got %s instead", message, expected, got))
	}
}

func trimString(s string) string {
	return strings.TrimRight(s, "\n")
}
