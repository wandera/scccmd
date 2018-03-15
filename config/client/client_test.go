package client

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

func TestErrorResponse(t *testing.T) {
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
		w.WriteHeader(406)
	}))
	defer ts.Close()

	_, err := NewClient(Config{
		URI:         ts.URL,
		Application: tp.application,
		Profile:     tp.profile,
		Label:       tp.label,
	}).FetchFile(tp.fileName)

	if err == nil {
		t.Error("FetchFile should have failed")
	}
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

func TestClient_Encrypt(t *testing.T) {
	var tp = struct {
		URI         string
		testContent string
	}{
		"/encrypt",
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

	cont, err := NewClient(Config{
		URI: ts.URL,
	}).Encrypt(tp.testContent)

	if err != nil {
		t.Error("Encrypt failed with: ", err)
	}

	assertString(t, "Content mismatch", tp.testContent, cont)
}

func TestClient_Decrypt(t *testing.T) {
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

	cont, err := NewClient(Config{
		URI: ts.URL,
	}).Decrypt(tp.testContent)

	if err != nil {
		t.Error("Decrypt failed with: ", err)
	}

	assertString(t, "Content mismatch", tp.testContent, cont)
}

func assertString(t *testing.T, message string, expected string, got string) {
	expected = trimString(expected)
	got = trimString(got)
	if expected != got {
		t.Errorf("%s: expected %s but got %s instead", message, expected, got)
	}
}

func trimString(s string) string {
	return strings.TrimRight(s, "\n")
}
