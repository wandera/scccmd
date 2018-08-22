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

func TestExecuteDiffFiles(t *testing.T) {
	var testParams = []struct {
		appName       string
		fileName      string
		difftext      string
		testContentA  string
		responseCodeA int
		profileA      string
		labelA        string
		requestURIA   string
		testContentB  string
		responseCodeB int
		profileB      string
		labelB        string
		requestURIB   string
	}{
		{"app",
			"src",
			"@@ -1,3 +1,3 @@\n foo\n-bar\n+baz\n ",
			"foo\nbar",
			200,
			"default",
			"master",
			"/app/default/master/src",
			"foo\nbaz",
			200,
			"default",
			"develop",
			"/app/default/develop/src",
		},
		{"app",
			"src",
			"@@ -1,3 +1,3 @@\n foo\n-bar\n+baz\n ",
			"foo\nbar",
			200,
			"development",
			"master",
			"/app/development/master/src",
			"foo\nbaz",
			200,
			"qa",
			"master",
			"/app/qa/master/src",
		},
		{"app",
			"src",
			"",
			"foo\nbar",
			200,
			"default",
			"master",
			"/app/default/master/src",
			"foo\nbar",
			200,
			"default",
			"develop",
			"/app/default/develop/src",
		},
		{"app",
			"src",
			"",
			"foo\nbar",
			200,
			"default",
			"master",
			"/app/default/master/src",
			"foo\nbar",
			200,
			"default",
			"master",
			"/app/default/master/src",
		},
		{"app",
			"src",
			"@@ -1,3 +1 @@\n-foo\n-bar\n ",
			"foo\nbar",
			200,
			"default",
			"master",
			"/app/default/master/src",
			"error",
			404,
			"default",
			"develop",
			"/app/default/develop/src",
		},
		{"app",
			"src",
			"@@ -1 +1,3 @@\n+foo\n+bar\n ",
			"error",
			404,
			"default",
			"master",
			"/app/default/master/src",
			"foo\nbar",
			200,
			"default",
			"develop",
			"/app/default/develop/src",
		},
		{"app",
			"src",
			"",
			"error",
			404,
			"default",
			"master",
			"/app/default/master/src",
			"error",
			404,
			"default",
			"develop",
			"/app/default/develop/src",
		},
	}

	for _, tp := range testParams {
		func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI == tp.requestURIA {
					w.WriteHeader(tp.responseCodeA)
					fmt.Fprintln(w, tp.testContentA)
				} else if r.RequestURI == tp.requestURIB {
					w.WriteHeader(tp.responseCodeB)
					fmt.Fprintln(w, tp.testContentB)
				} else {
					t.Errorf("Expected call to '%s' or '%s', but got '%s' instead.", tp.requestURIA, tp.requestURIB, r.RequestURI)
				}
			}))
			defer ts.Close()

			diffp.application = tp.appName
			diffp.profile = tp.profileA
			diffp.label = tp.labelA
			diffp.targetLabel = tp.labelB
			diffp.targetProfile = tp.profileB
			diffp.source = ts.URL
			diffp.files = tp.fileName

			filename := "stdout"
			old := os.Stdout               // keep backup of the real stdout
			temp, _ := os.Create(filename) // create temp file
			os.Stdout = temp
			defer func() {
				temp.Close()
				os.Stdout = old // restoring the real stdout
			}()

			if err := ExecuteDiffFiles(nil); err != nil {
				t.Error("Execute failed with: ", err)
			}

			raw, err := ioutil.ReadFile(filename)
			defer os.Remove(filename)
			if err != nil {
				t.Error("Expected to download file: ", err)
			}

			if response := strings.TrimRight(string(raw[:]), "\n"); response != tp.difftext {
				t.Errorf("Expected response: '%s' got '%s' instead.", tp.difftext, response)
			}
		}()
	}
}

func TestExecuteDiffValues(t *testing.T) {
	var testParams = []struct {
		appName      string
		difftext     string
		testContentA string
		profileA     string
		labelA       string
		requestURIA  string
		testContentB string
		profileB     string
		labelB       string
		requestURIB  string
		format       string
	}{
		{"app",
			"@@ -1,2 +1,2 @@\n foo: 1\n-bar: 2\n+baz: 2",
			"foo: 1\nbar: 2",
			"default",
			"master",
			"/master/app-default.yml",
			"foo: 1\nbaz: 2",
			"default",
			"develop",
			"/develop/app-default.yml",
			"yaml",
		},
		{"app",
			"",
			"foo=bar",
			"default",
			"master",
			"/master/app-default.properties",
			"foo=bar",
			"default",
			"develop",
			"/develop/app-default.properties",
			"properties",
		},
		{"app",
			"@@ -1 +1 @@\n-{\"foo\":\"bar\", \"foo\":\"bar\"}\n+{\"foo\":\"bar\", \"foo\":\"baz\"}",
			"{\"foo\":\"bar\", \"foo\":\"bar\"}",
			"qa",
			"master",
			"/master/app-qa.json",
			"{\"foo\":\"bar\", \"foo\":\"baz\"}",
			"development",
			"develop",
			"/develop/app-development.json",
			"json",
		},
	}

	for _, tp := range testParams {
		func() {
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.RequestURI == tp.requestURIA {
					fmt.Fprintln(w, tp.testContentA)
				} else if r.RequestURI == tp.requestURIB {
					fmt.Fprintln(w, tp.testContentB)
				} else {
					t.Errorf("Expected call to '%s' or '%s', but got '%s' instead.", tp.requestURIA, tp.requestURIB, r.RequestURI)
				}
			}))
			defer ts.Close()

			diffp.application = tp.appName
			diffp.profile = tp.profileA
			diffp.label = tp.labelA
			diffp.targetProfile = tp.profileB
			diffp.targetLabel = tp.labelB
			diffp.source = ts.URL
			diffp.format = tp.format

			filename := "stdout"
			old := os.Stdout               // keep backup of the real stdout
			temp, _ := os.Create(filename) // create temp file
			os.Stdout = temp
			defer func() {
				temp.Close()
				os.Stdout = old // restoring the real stdout
			}()

			if err := ExecuteDiffValues(nil); err != nil {
				t.Error("Execute failed with: ", err)
			}

			raw, err := ioutil.ReadFile(filename)
			defer os.Remove(filename)
			if err != nil {
				t.Error("Expected to download file: ", err)
			}

			if response := strings.TrimRight(string(raw[:]), "\n"); response != tp.difftext {
				t.Errorf("Expected response: '%s' got '%s' instead.", tp.difftext, response)
			}
		}()
	}

}
