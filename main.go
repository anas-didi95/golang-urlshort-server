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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// URL Document for urls collection
type URL struct {
	OriginalURL      string    `json:"originalURL,omitempty"`
	ShortID          string    `json:"shortID,omitempty"`
	LastModifiedDate time.Time `json:"lastModifiedDate"`
}

func main() {
	setupMongoDatabase()

	contextPath := mux.NewRouter().StrictSlash(true)
	router := contextPath.PathPrefix("/urlshort").Subrouter()

	router.HandleFunc("/hello/{name}", GetHelloWorld).Methods(http.MethodGet)
	router.HandleFunc("/generate", PostGenerateShortURL).Methods(http.MethodPost)
	router.HandleFunc("/s/{shortID}", GetRedirectShortURL).Methods(http.MethodGet)

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
		OriginalURL:      requestBody.URL,
		ShortID:          randSeq(7),
		LastModifiedDate: time.Now(),
	}
	_, err = collection.InsertOne(context.TODO(), document)
	defer client.Disconnect(context.TODO())
	if err != nil {
		log.Fatalf("[%s] Insert mongo document failed! %v", TAG, err)
		sendResponse(w, http.StatusInternalServerError, nil, false, "Insert mongo document failed!")
		return
	}

	responseBody := map[string]interface{}{
		"originalURL": requestBody.URL,
		"shortURL":    os.Getenv("BASE_URL") + "/s/" + document.ShortID,
		"shortID":     document.ShortID,
	}
	sendResponse(w, http.StatusOK, responseBody, true, "Short URL generated successfully.")
}

// GetRedirectShortURL Redirect short URL
func GetRedirectShortURL(w http.ResponseWriter, r *http.Request) {
	TAG := "GetRedirectShortURL"

	vars := mux.Vars(r)
	shortID := vars["shortID"]
	if len(shortID) == 0 {
		log.Fatalf("[%s] Short ID not sent!", TAG)
		sendResponse(w, http.StatusInternalServerError, nil, false, "Short ID not sent!")
		return
	}

	client := getMongoConnection()
	filter := bson.D{primitive.E{Key: "shortid", Value: shortID}}
	var URL URL
	err := client.Database("urlshort").Collection("urls").FindOne(context.TODO(), filter).Decode(&URL)
	if err != nil {
		log.Fatalf("[%s] Get mongo document failed! %v", TAG, err)
		sendResponse(w, http.StatusInternalServerError, nil, false, "Get mongo document failed!")
		return
	}
	defer client.Disconnect(context.TODO())

	http.Redirect(w, r, URL.OriginalURL, http.StatusSeeOther)
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

func setupMongoDatabase() {
	TAG := "setupMongoDatabase"

	client := getMongoConnection()
	collections, err := client.Database("urlshort").ListCollectionNames(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatalf("[%s] Get collection names failed! %v", TAG, err)
	}

	hasURLS := contains(collections, "urls")
	if !hasURLS {
		client.Database("urlshort").CreateCollection(context.TODO(), "urls")
		log.Printf("[%s] Collection urls created successfully", TAG)
	}

	client.Database("urlshort").Collection("urls").Indexes().DropAll(context.TODO())

	idxLastModifiedDateTTL := mongo.IndexModel{Keys: bson.M{
		"lastmodifieddate": 1,
	}, Options: options.Index().SetExpireAfterSeconds(1 * 24 * 60 * 60).SetName("ttl_lastmodifieddate")}
	client.Database("urlshort").Collection("urls").Indexes().CreateOne(context.TODO(), idxLastModifiedDateTTL)
	log.Printf("[%s][%s] Index ttl_lastmodifieddate created successfully", TAG, "urls")

	defer client.Disconnect(context.TODO())
}

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
