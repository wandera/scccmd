package cmd

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNoArgExecute(t *testing.T) {
	err := executeGetFiles(nil)
	if err != nil {
		t.Error("Execute failed with: ", err)
	}
}

func TestExecuteGetFiles(t *testing.T) {
	var testParams = []struct {
		testContent  string
		appName      string
		profile      string
		label        string
		srcFileName  string
		destFileName string
		requestURI   string
	}{
		{"{\"foo\":\"bar\"}",
			"app",
			"default",
			"master",
			"src",
			"dest",
			"/app/default/master/src"},
		{"{\"bar\":\"foo\"}",
			"app2",
			"default",
			"master",
			"src2",
			"destination",
			"/app2/default/master/src2"},
		{"{\"foo\":\"bar\"}",
			"app",
			"prod",
			"1.0.0",
			"app.yaml",
			"config.yaml",
			"/app/prod/1.0.0/app.yaml"},
	}

	for _, tp := range testParams {
		func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI != tp.requestURI {
					t.Error(fmt.Sprintf("Expected call to '%s' but got '%s' instead.", tp.requestURI, r.RequestURI))
				}

				fmt.Fprintln(w, tp.testContent)
			}))
			defer ts.Close()

			application = tp.appName
			profile = tp.profile
			label = tp.label
			source = ts.URL
			fileMappings = FileMappings{mappings: make([]FileMapping, 1)}
			fileMappings.mappings[0] = FileMapping{source: tp.srcFileName, destination: tp.destFileName}

			if err := executeGetFiles(nil); err != nil {
				t.Error("Execute failed with: ", err)
			}

			raw, err := ioutil.ReadFile(tp.destFileName)
			defer os.Remove(tp.destFileName)
			if err != nil {
				t.Error("Expected to download file: ", err)
			}

			if response := strings.TrimRight(string(raw[:]), "\n"); response != tp.testContent {
				t.Error(fmt.Sprintf("Expected response: '%s' got '%s' instead.", tp.testContent, response))
			}
		}()
	}
}

func TestExecuteGetValues(t *testing.T) {
	var testParams = []struct {
		testContent  string
		appName      string
		profile      string
		label        string
		destFileName string
		requestURI   string
	}{
		{"{\"foo\":\"bar\"}",
			"app",
			"default",
			"master",
			"dest",
			"/master/app-default.yml"},
		{"{\"bar\":\"foo\"}",
			"app2",
			"default",
			"master",
			"destination",
			"/master/app2-default.yml"},
		{"{\"foo\":\"bar\"}",
			"app",
			"prod",
			"1.0.0",
			"config.yaml",
			"/1.0.0/app-prod.yml"},
	}

	for _, tp := range testParams {
		func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI != tp.requestURI {
					t.Error(fmt.Sprintf("Expected call to '%s' but got '%s' instead.", tp.requestURI, r.RequestURI))
				}

				fmt.Fprintln(w, tp.testContent)
			}))
			defer ts.Close()

			application = tp.appName
			profile = tp.profile
			label = tp.label
			source = ts.URL
			destination = tp.destFileName

			if err := executeGetValues(nil); err != nil {
				t.Error("Execute failed with: ", err)
			}

			raw, err := ioutil.ReadFile(tp.destFileName)
			defer os.Remove(tp.destFileName)
			if err != nil {
				t.Error("Expected to download file: ", err)
			}

			if response := strings.TrimRight(string(raw[:]), "\n"); response != tp.testContent {
				t.Error(fmt.Sprintf("Expected response: '%s' got '%s' instead.", tp.testContent, response))
			}
		}()
	}

}
