package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// URL Document for urls collection
type URL struct {
	OriginalURL string `json:"originalUrl"`
	ShortURL    string `json:"shortUrl"`
}

func main() {
	contextPath := mux.NewRouter().StrictSlash(true)
	router := contextPath.PathPrefix("/urlshort").Subrouter()

	router.HandleFunc("/hello/{name}", GetHelloWorld).Methods(http.MethodGet)
	router.HandleFunc("/generate", PostGenerateShortURL).Methods(http.MethodPost)

	AppHost := os.Getenv("APP_HOST")
	AppPort := os.Getenv("APP_PORT")
	log.Printf("[main] Server started at %s:%s", AppHost, AppPort)
	log.Fatal(http.ListenAndServe(AppHost+":"+AppPort, router))
}

// GetHelloWorld Return greeting based on name and language given
func GetHelloWorld(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	query := r.URL.Query()
	lang := query.Get("lang")

	if len(lang) <= 0 {
		lang = "en"
	}

	value := ""
	switch lang {
	case "bm":
		value = "Selamat sejahtera, "
		break
	default:
		value = "Hello, "
		break
	}

	client := getMongoConnection()
	err := client.Ping(context.TODO(), nil)
	defer client.Disconnect(context.TODO())

	var mongoStatus string
	if err != nil {
		log.Fatal(err)
		mongoStatus = "Offline"
	} else {
		mongoStatus = "Online"
	}

	data := map[string]interface{}{
		"value": value + name,
		"mongo": mongoStatus,
	}

	sendResponse(w, http.StatusOK, data, true, "Response returned successfully.")
}

// PostGenerateShortURL Generate and return short url
func PostGenerateShortURL(w http.ResponseWriter, r *http.Request) {
	TAG := "PostGenerateShortURL"

	var requestBody struct {
		URL string `json:"url"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		log.Fatalf("[%s] Read request body failed! %v", TAG, err)
		sendResponse(w, http.StatusInternalServerError, nil, false, "Read request body failed!")
		return
	}

	client := getMongoConnection()
	collection := client.Database("urlshort").Collection("urls")

	document := URL{
		OriginalURL: requestBody.URL,
		ShortURL:    randSeq(7),
	}
	_, err = collection.InsertOne(context.TODO(), document)
	defer client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalf("[%s] Insert mongo document failed! %v", TAG, err)
		sendResponse(w, http.StatusInternalServerError, nil, false, "Insert mongo document failed!")
		return
	}

	responseBody := map[string]interface{}{
		"originalUrl": requestBody.URL,
		"shortUrl":    "https://api.anasdidi.dev/urlshort/s/" + document.ShortURL,
	}
	sendResponse(w, http.StatusOK, responseBody, true, "Endpoint reached.")
}

func sendResponse(w http.ResponseWriter, statusCode int, data map[string]interface{}, isSuccess bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	responseBody := map[string]interface{}{
		"status": map[string]interface{}{
			"isSuccess": isSuccess,
			"message":   message,
		},
		"data": data,
	}

	json.NewEncoder(w).Encode(responseBody)
}

func getMongoConnection() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING")))
	if err != nil {
		log.Fatalf("Mongo client connection failed! %v", err)
	}

	return client
}

func randSeq(n int) string {
	if os.Getenv("IS_TEST") == "true" {
		return "1234567"
	}

	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	b := make([]rune, n)

	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
