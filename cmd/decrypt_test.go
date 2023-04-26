package cmd

import (
	"bytes"
	"fmt"
	"github.com/wandera/scccmd/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestExecuteDecrypt(t *testing.T) {
	tp := struct {
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
		buf.ReadFrom(r.Body)

		testutil.AssertString(t, "Incorrect Content received", tp.testContent, buf.String())
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
	tp := struct {
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
