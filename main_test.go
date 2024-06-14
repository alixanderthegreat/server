package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestIndex tests the index handler for various aspects
func TestIndex(t *testing.T) {

	for _, route := range routes {
		route := route // create a new instance of the route variable to avoid closure issues
		http.HandleFunc(
			route.Path, func(w http.ResponseWriter, r *http.Request) {
				serve(w, r, route.File, route.MimeType)
			},
		)
	}

	// Create a new test request and response recorder
	given := func(path string) ([]byte, *httptest.ResponseRecorder, error) {
		req, rr := newTestRequest(t, path)
		handler := http.DefaultServeMux // Use the default serve mux which has been set up
		handler.ServeHTTP(rr, req)

		// Read the body of the response
		body, err := io.ReadAll(rr.Body)
		if err != nil {
			t.Fatal(err)
		}

		return body, rr, err
	}

	body, recorded_response, _ := given("/")

	// Define the tests with expected values
	var tests = []struct {
		Name     string
		Expected any
		Check    func(t *testing.T, expected any)
	}{
		{
			Name:     "StatusCode",
			Expected: http.StatusOK,
			Check: func(t *testing.T, expected any) {
				if status := recorded_response.Code; status != expected.(int) {
					t.Errorf("handler returned wrong status code: got %v want %v", status, expected.(int))
				}
			},
		},
		{
			Name:     "ContentType",
			Expected: "text/html; charset=utf-8",
			Check: func(t *testing.T, expected any) {
				if contentType := recorded_response.Header().Get("Content-Type"); contentType != expected.(string) {
					t.Errorf("handler returned unexpected content type: got %v want %v", contentType, expected.(string))
				}
			},
		},
		{
			Name:     "Doctype",
			Expected: "<!DOCTYPE html>",
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: Doctype not found; got %v want %v", string(body), expected.(string))
				}
			},
		},
		{
			Name:     "LangAttribute",
			Expected: `lang="en"`,
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: Lang attribute not found; got %v want %v", string(body), expected.(string))
				}
			},
		},
		{
			Name:     "MetaCharset",
			Expected: `<meta charset="UTF-8">`,
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: Meta charset not found; got %v want %v", string(body), expected.(string))
				}
			},
		},
		{
			Name:     "MetaViewport",
			Expected: `<meta name="viewport" content="width=device-width, initial-scale=1.0">`,
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: Meta viewport not found; got %v want %v", string(body), expected.(string))
				}
			},
		},
		{
			Name:     "Body",
			Expected: "<p>Hello World!</p>",
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: got %v want %v", string(body), expected.(string))
				}
			},
		},
		{
			Name:     "JavaScript",
			Expected: `<script type="module" src="script.js"></script>`,
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: JavaScript not found; got %v want %v", string(body), expected.(string))
				}
			},
		},
		{
			Name:     "CSS",
			Expected: `<link rel="stylesheet" type="text/css" href="styles.css">`,
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: CSS not found; got %v want %v", string(body), expected.(string))
				}
			},
		},
		{
			Name:     "app",
			Expected: `<div id="app">`,
			Check: func(t *testing.T, expected any) {
				expectedStart := expected.(string)
				if !strings.Contains(string(body), expectedStart) {
					t.Errorf("handler returned unexpected body: Clock div not found; got %v want %v", string(body), expectedStart)
				}
				// Check for the presence of the current time
				now := time.Now()
				expectedTime := now.Format("15:04:05")
				nextSecond := now.Add(1 * time.Second).Format("15:04:05")

				// Check if the body contains the current or the next minute's time string
				if !strings.Contains(string(body), expectedTime) && !strings.Contains(string(body), nextSecond) {
					t.Errorf("handler returned unexpected body: Current time not found within range; got %v want %v or %v", string(body), expectedTime, nextSecond)
				}
			},
		},
		{
			Name:     "HTMLTitle",
			Expected: "<title>Hello World</title>",
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: HTML title not found; got %v want %v", string(body), expected.(string))
				}
			},
		},
		{
			Name:     "Favicon",
			Expected: "favicon.ico",
			Check: func(t *testing.T, expected any) {
				if !strings.Contains(string(body), expected.(string)) {
					t.Errorf("handler returned unexpected body: Favicon not found; got %v", string(body))
				}
			},
		},
	}

	// Run the tests
	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			test.Check(t, test.Expected)
		})
	}
}

// newTestRequest creates a new HTTP GET request and response recorder
func newTestRequest(t *testing.T, path string) (*http.Request, *httptest.ResponseRecorder) {
	req, err := http.NewRequest("GET", path, nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	return req, rr
}
