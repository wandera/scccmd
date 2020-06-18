package client

import (
	"bytes"
	"fmt"
	"github.com/wandera/scccmd/internal"
	"net/http"
	"net/http/httptest"
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

	testutil.AssertString(t, "Incorrect URI", tp.URI, c.Config().URI)
	testutil.AssertString(t, "Incorrect Application", tp.application, c.Config().Application)
	testutil.AssertString(t, "Incorrect Profile", tp.profile, c.Config().Profile)
	testutil.AssertString(t, "Incorrect Label", tp.label, c.Config().Label)
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
		testutil.AssertString(t, "Incorrect URI call", tp.URI, r.RequestURI)
		w.WriteHeader(406)
	}))
	defer ts.Close()

	_, err := NewClient(Config{
		URI:         ts.URL,
		Application: tp.application,
		Profile:     tp.profile,
		Label:       tp.label,
	}).FetchFileE(tp.fileName)

	if err == nil {
		t.Error("FetchFile should have failed")
	}
}

func TestRedirectResponse(t *testing.T) {
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
		testutil.AssertString(t, "Incorrect URI call", tp.URI, r.RequestURI)
		w.Header().Add("location", "http://somwhere.el.se")
		w.WriteHeader(301)
	}))
	defer ts.Close()

	_, err := NewClient(Config{
		URI:         ts.URL,
		Application: tp.application,
		Profile:     tp.profile,
		Label:       tp.label,
	}).FetchFileE(tp.fileName)

	if err == nil {
		t.Error("FetchFile should have failed")
	}
}

func Test503Response(t *testing.T) {
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
		testutil.AssertString(t, "Incorrect URI call", tp.URI, r.RequestURI)
		w.WriteHeader(503)
		return
	}))
	defer ts.Close()

	_, err := NewClient(Config{
		URI:         ts.URL,
		Application: tp.application,
		Profile:     tp.profile,
		Label:       tp.label,
	}).FetchFileE(tp.fileName)

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
		testutil.AssertString(t, "Incorrect URI call", tp.URI, r.RequestURI)
		_, _ = fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	cont, err := NewClient(Config{
		URI:         ts.URL,
		Application: tp.application,
		Profile:     tp.profile,
		Label:       tp.label,
	}).FetchFileE(tp.fileName)

	if err != nil {
		t.Error("FetchFile failed with: ", err)
	}

	testutil.AssertString(t, "Content mismatch", tp.testContent, string(cont))
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
		testutil.AssertString(t, "Incorrect URI call", tp.URI, r.RequestURI)
		_, _ = fmt.Fprintln(w, tp.testContent)
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

	testutil.AssertString(t, "Content mismatch", tp.testContent, cont)
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
		testutil.AssertString(t, "Incorrect Method", "POST", r.Method)
		testutil.AssertString(t, "Incorrect URI call", tp.URI, r.RequestURI)

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(r.Body)

		testutil.AssertString(t, "Incorrect Content received", tp.testContent, buf.String())
		_, _ = fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	cont, err := NewClient(Config{
		URI: ts.URL,
	}).Encrypt(tp.testContent)

	if err != nil {
		t.Error("Encrypt failed with: ", err)
	}

	testutil.AssertString(t, "Content mismatch", tp.testContent, cont)
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
		testutil.AssertString(t, "Incorrect Method", "POST", r.Method)
		testutil.AssertString(t, "Incorrect URI call", tp.URI, r.RequestURI)

		buf := new(bytes.Buffer)
		_, _ = buf.ReadFrom(r.Body)

		testutil.AssertString(t, "Incorrect Content received", tp.testContent, buf.String())
		_, _ = fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	cont, err := NewClient(Config{
		URI: ts.URL,
	}).Decrypt(tp.testContent)

	if err != nil {
		t.Error("Decrypt failed with: ", err)
	}

	testutil.AssertString(t, "Content mismatch", tp.testContent, cont)
}
