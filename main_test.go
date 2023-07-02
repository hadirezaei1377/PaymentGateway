package main_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBank(t *testing.T) {
	// Create a request to simulate the "/Bank{price}" endpoint
	req, err := http.NewRequest("GET", "/Bank5000", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response
	rr := httptest.NewRecorder()

	// Call the Bank handler function with the request and response recorder
	Bank(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}

	// Additional assertions for the response body or headers can be added here
}

func Bank(rr *httptest.ResponseRecorder, req *http.Request) {
	panic("unimplemented")
}

func TestCallBack(t *testing.T) {
	// Create a request to simulate the "/CallBack{price}" endpoint
	req, err := http.NewRequest("GET", "/CallBack5000", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Create a ResponseRecorder to record the response

	rr := httptest.NewRecorder()

	// Call the CallBack handler function with the request and response recorder
	CallBack(rr, req)

	// Check the response status code
	if rr.Code != http.StatusOK {
		t.Errorf("expected status OK; got %v", rr.Code)
	}

	// Additional assertions for the response body or headers can be added here
}

func CallBack(rr *httptest.ResponseRecorder, req *http.Request) {
	panic("unimplemented")
}
