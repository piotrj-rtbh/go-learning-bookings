package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

type postData struct {
	key   string
	value string
}

// creating table tests
type tableTests []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}

var theTests = tableTests{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gq", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	// use httptest server right for tests
	ts := httptest.NewTLSServer(routes)

	defer ts.Close() // when this function finishes then the the server is closed

	for _, e := range theTests {
		if e.method == "GET" {
			// create a fake web client
			resp, err := ts.Client().Get(ts.URL + e.url) // e.url is only a path so we mast prepend a real proto and domain
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}
			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		} else {

		}
	}

}
