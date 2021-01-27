package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var clientOptions *options.ClientOptions

func main() {
	clientOptions = options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING"))
	contextPath := mux.NewRouter().StrictSlash(true)
	router := contextPath.PathPrefix("/urlshort").Subrouter()

	router.HandleFunc("/hello/{name}", GetHelloWorld).Methods(http.MethodGet)

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
	client, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGO_CONNECTION_STRING")))
	if err != nil {
		log.Fatalf("Mongo client failed! %v", err)
	}

	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatalf("Mongo client connection failed! %v", err)
	}

	return client
}
