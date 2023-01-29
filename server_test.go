package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	testCases := []struct {
		name       string
		path       string
		expCode    int
		expItems   int
		expContent string
	}{
		{name: "GetRoot", path: "/", expCode: http.StatusOK, expContent: content},
		{name: "NotFound", path: "/api/500", expCode: http.StatusNotFound},
	}

	url, cleanup := setupAPI(t)
	//test url :=  http://127.0.0.1:37377
	defer cleanup()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var (
				body []byte
				err  error
			)

			//append the test path to the test server url
			//http://127.0.0.1:37377/  --> expCode: ok
			//http://127.0.0.1:37377/api/500 --> expCode: not found
			w, err := http.Get(url + tc.path)
			if err != nil {
				t.Error(err)
			}
			defer func(Body io.ReadCloser) {
				err = Body.Close()
				if err != nil {
					t.Error(err)
				}
			}(w.Body)

			if w.StatusCode != tc.expCode {
				t.Fatalf("Expected %q, got %q.", http.StatusText(tc.expCode), http.StatusText(w.StatusCode))
			}

			switch {
			case strings.Contains(w.Header.Get("Content-Type"), "text/plain"):
				if body, err = io.ReadAll(w.Body); err != nil {
					t.Error(err)
				}
				if !strings.Contains(string(body), tc.expContent) {
					t.Errorf("Expected %q, got %q.", tc.expContent, string(body))
				}
			default:
				t.Fatalf("Unsupported Content-Type: %q", w.Header.Get("Content-Type"))
			}

		})
	}
}

// setupAPI creates a new test server and returns its URL and a cleanup function
func setupAPI(t *testing.T) (string, func()) {
	t.Helper()

	testSrv := httptest.NewServer(newMux(""))

	return testSrv.URL, func() {
		testSrv.Close()
	}
}
