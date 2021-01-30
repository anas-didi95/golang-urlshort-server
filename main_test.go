package main

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestGetEnvVar(t *testing.T) {
	AppHost := os.Getenv("APP_HOST")
	if len(AppHost) == 0 {
		t.Errorf("Env var not defined! key %s", "APP_HOST")
	}

	AppPort := os.Getenv("APP_PORT")
	if len(AppPort) == 0 {
		t.Errorf("Env var not defined! key %s", "APP_PORT")
	}

	MongoConnectionString := os.Getenv("MONGO_CONNECTION_STRING")
	if len(MongoConnectionString) == 0 {
		t.Errorf("Env var not defined! key %s", "MONGO_CONNECTION_STRING")
	}

	BaseURL := os.Getenv("BASE_URL")
	if len(BaseURL) == 0 {
		t.Errorf("Env var not defined! key %s", "BASE_URL")
	}

	ShortURLExpireInSeconds := os.Getenv("SHORT_URL_EXPIRE_IN_SECONDS")
	if len(ShortURLExpireInSeconds) == 0 {
		t.Errorf("Env var not defined! key %s", "SHORT_URL_EXPIRE_IN_SECONDS")
	}
}

func TestPingMongoClient(t *testing.T) {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING")))
	if err != nil {
		t.Errorf("Mongo client connection failed! %v", err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		t.Errorf("Mongo client ping failed! %v", err)
	}

	err = client.Database("urlshort").Drop(context.TODO())
	if err != nil {
		t.Errorf("Mongo client drop database failed! %v", err)
	}

	defer client.Disconnect(context.TODO())
}

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
			"mongo": "Online",
		},
	}
	expectedResponseBody, err := json.Marshal(responseBody)
	assertedResponseBody := strings.Trim(rr.Body.String(), "\n")

	if assertedResponseBody != string(expectedResponseBody) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			assertedResponseBody, string(expectedResponseBody))
	}
}

func TestPostGenerateShortURL(t *testing.T) {
	os.Setenv("IS_TEST", "true")
	requestBody := map[string]string{
		"url": "https://www.google.com",
	}
	var buf bytes.Buffer
	json.NewEncoder(&buf).Encode(requestBody)

	req, err := http.NewRequest(http.MethodPost, "/urlshort/generate", &buf)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/urlshort/generate", PostGenerateShortURL)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	responseBody := map[string]interface{}{
		"status": map[string]interface{}{
			"isSuccess": true,
			"message":   "Short URL generated successfully.",
		},
		"data": map[string]interface{}{
			"originalURL": "https://www.google.com",
			"shortURL":    os.Getenv("BASE_URL") + "/s/1234567",
			"shortID":     "1234567",
		},
	}
	expectedResponseBody, err := json.Marshal(responseBody)
	assertedResponseBody := strings.Trim(rr.Body.String(), "\n")

	if assertedResponseBody != string(expectedResponseBody) {
		t.Errorf("handler returned unexpected body: got %v want %v",
			assertedResponseBody, string(expectedResponseBody))
	}

	os.Unsetenv("IS_TEST")
}

func TestGetRedirectShortURL(t *testing.T) {
	os.Setenv("IS_TEST", "true")
	req, err := http.NewRequest(http.MethodGet, "/urlshort/s/1234567", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/urlshort/s/{shortID}", GetRedirectShortURL)
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusSeeOther {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusSeeOther)
	}

	os.Unsetenv("IS_TEST")
}
