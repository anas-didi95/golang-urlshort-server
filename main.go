package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	contextPath := mux.NewRouter().StrictSlash(true)
	router := contextPath.PathPrefix("/urlshort").Subrouter()

	router.HandleFunc("/hello/{name}", GetHelloWorld).Methods(http.MethodGet)

	log.Printf("[main] Server started at %s:%s", "0.0.0.0", "5000")
	log.Fatal(http.ListenAndServe("0.0.0.0:5000", router))
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

	data := map[string]interface{}{
		"value": value + name,
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
