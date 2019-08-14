package cmd

import (
	"bytes"
	"fmt"
	"github.com/WanderaOrg/scccmd/internal"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteEncrypt(t *testing.T) {
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
		buf.ReadFrom(r.Body)

		testutil.AssertString(t, "Incorrect Content received", tp.testContent, buf.String())
		fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	ep.source = ts.URL
	ep.value = tp.testContent
	err := ExecuteEncrypt()

	if err != nil {
		t.Error("Encrypt failed with: ", err)
	}
}

func ExampleExecuteEncrypt() {
	var tp = struct {
		URI         string
		testContent string
	}{
		"/encrypt",
		"test",
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, tp.testContent)
	}))
	defer ts.Close()

	ep.source = ts.URL
	ep.value = tp.testContent
	ExecuteEncrypt()
	// Output: test
}
