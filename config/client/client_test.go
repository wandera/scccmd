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

	if c.Config().URI != tp.URI {
		t.Error("Incorrect URI ", assertString(c.Config().URI, tp.URI))
	}

	if c.Config().Application != tp.application {
		t.Error("Incorrect Application ", assertString(c.Config().Application, tp.application))
	}

	if c.Config().Profile != tp.profile {
		t.Error("Incorrect Profile ", assertString(c.Config().Profile, tp.profile))
	}

	if c.Config().Label != tp.label {
		t.Error("Incorrect Label ", assertString(c.Config().Label, tp.label))
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
		if r.RequestURI != tp.URI {
			t.Error(fmt.Sprintf("Expected call to '%s' but got '%s' instead.", tp.URI, r.RequestURI))
		}

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

	if gotContent := strings.TrimRight(string(cont), "\n"); gotContent != tp.testContent {
		t.Error("Content mismatch ", assertString(tp.testContent, gotContent))
	}
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
		if r.RequestURI != tp.URI {
			t.Error(fmt.Sprintf("Expected call to '%s' but got '%s' instead.", tp.URI, r.RequestURI))
		}

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

	if gotContent := strings.TrimRight(string(cont), "\n"); gotContent != tp.testContent {
		t.Error("Content mismatch ", assertString(tp.testContent, gotContent))
	}
}

func assertString(expected string, got string) string {
	return fmt.Sprintf("expected %s but got %s instead", expected, got)
}
