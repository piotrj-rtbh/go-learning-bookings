package main

import (
	"fmt"
	"net/http"
	"testing"
)

// how to pass a "next" param? We have to create in setup_test.go as myHandler struct
func TestNoSurf(t *testing.T) {
	var myH myHandler
	h := NoSurf(&myH)

	// when we call NoSurf it should receive http.Handler or something that satisfies the interface
	// and should return something that is also http.Handler
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}

func TestSessionLoad(t *testing.T) {
	var myH myHandler
	h := SessionLoad(&myH)

	// when we call NoSurf it should receive http.Handler or something that satisfies the interface
	// and should return something that is also http.Handler
	switch v := h.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Error(fmt.Sprintf("type is not http.Handler, but is %T", v))
	}
}
