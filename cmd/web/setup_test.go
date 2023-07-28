package main

import (
	"net/http"
	"os"
	"testing"
)

// THIS FILE WILL RUN BEFORE ALL OF OUR TESTS RUN
// main test function must be named as TestMain
func TestMain(m *testing.M) {
	// do something before

	// finally run the tests and exit
	os.Exit(m.Run())
}

// this must satisfy the http.Handler interface
// we must implement the same methods/functions that exist in the interface
type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
