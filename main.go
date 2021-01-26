package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	contextPath := mux.NewRouter().StrictSlash(true)
	router := contextPath.PathPrefix("/urlshort").Subrouter()

	router.HandleFunc("/", ArticlesCategoryHandler)

	http.ListenAndServe("0.0.0.0:5000", router)
}

// ArticlesCategoryHandler TEST
func ArticlesCategoryHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category: %v\n", vars["category"])
}
