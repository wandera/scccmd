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

func TestNoArgGetExecute(t *testing.T) {
	err := ExecuteGetFiles()
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
		{"{\"foo\":\"bar\"}",
			"app",
			"default",
			"master",
			"src",
			"-",
			"/app/default/master/src"},
	}

	for _, tp := range testParams {
		func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI != tp.requestURI {
					t.Errorf("Expected call to '%s' but got '%s' instead.", tp.requestURI, r.RequestURI)
				}

				fmt.Fprintln(w, tp.testContent)
			}))
			defer ts.Close()

			gp.application = tp.appName
			gp.profile = tp.profile
			gp.label = tp.label
			gp.source = ts.URL
			gp.fileMappings = FileMappings{mappings: make([]FileMapping, 1)}
			gp.fileMappings.mappings[0] = FileMapping{source: tp.srcFileName, destination: tp.destFileName}

			filename := ""
			var old *os.File
			var temp *os.File
			if tp.destFileName == stdoutPlaceholder {
				filename = "stdout"
				old = os.Stdout               // keep backup of the real stdout
				temp, _ = os.Create("stdout") // create temp file
				os.Stdout = temp
				defer func() {
					temp.Close()
					os.Stdout = old // restoring the real stdout
				}()
			} else {
				filename = tp.destFileName
			}
			if err := ExecuteGetFiles(); err != nil {
				t.Error("Execute failed with: ", err)
			}

			raw, err := ioutil.ReadFile(filename)
			defer os.Remove(filename)
			if err != nil {
				t.Error("Expected to download file: ", err)
			}

			if response := strings.TrimRight(string(raw[:]), "\n"); response != tp.testContent {
				t.Errorf("Expected response: '%s' got '%s' instead.", tp.testContent, response)
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
		format       string
	}{
		{"{\"foo\":\"bar\"}",
			"app",
			"default",
			"master",
			"dest",
			"/master/app-default.json",
			"json"},
		{"bar=foo",
			"app2",
			"default",
			"master",
			"destination",
			"/master/app2-default.properties",
			"properties"},
		{"\"foo\":\"bar\"",
			"app",
			"prod",
			"1.0.0",
			"config.yaml",
			"/1.0.0/app-prod.yml",
			"yaml"},
	}

	for _, tp := range testParams {
		func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI != tp.requestURI {
					t.Errorf("Expected call to '%s' but got '%s' instead.", tp.requestURI, r.RequestURI)
				}

				fmt.Fprintln(w, tp.testContent)
			}))
			defer ts.Close()

			gp.application = tp.appName
			gp.profile = tp.profile
			gp.label = tp.label
			gp.source = ts.URL
			gp.destination = tp.destFileName
			gp.format = tp.format

			if err := ExecuteGetValues(); err != nil {
				t.Error("Execute failed with: ", err)
			}

			raw, err := ioutil.ReadFile(tp.destFileName)
			defer os.Remove(tp.destFileName)
			if err != nil {
				t.Error("Expected to download file: ", err)
			}

			if response := strings.TrimRight(string(raw[:]), "\n"); response != tp.testContent {
				t.Errorf("Expected response: '%s' got '%s' instead.", tp.testContent, response)
			}
		}()
	}

}
