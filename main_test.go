package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

func TestGetHelloWorld(t *testing.T) {
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest("GET", "/urlshort/hello/anas?lang=bm", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	// Need to create a router that we can pass the request through so that the vars will be added to the context
	router := mux.NewRouter()
	router.HandleFunc("/urlshort/hello/{name}", GetHelloWorld) // If error, it happen due linter problem

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	router.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	responseBody := map[string]interface{}{
		"status": map[string]interface{}{
			"isSuccess": true,
			"message":   "Response returned successfully.",
		},
		"data": map[string]interface{}{
			"value": "Selamat sejahtera, anas",
		},
	}
	expectedResponseBody, err := json.Marshal(responseBody)
	assertedResponseBody := strings.Trim(rr.Body.String(), "\n")

	if assertedResponseBody != string(expectedResponseBody) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			assertedResponseBody, string(expectedResponseBody))
	}

}
