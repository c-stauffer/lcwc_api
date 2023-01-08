package main

import (
	"log"
	"net/http"

	"code.crogge.rs/chris/lcwc_api/pkg/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	//router.HandleFunc("/books", func(w http.ResponseWriter, r *http.Request) {
	//	json.NewEncoder(w).Encode("Hello World")
	//})
	router.HandleFunc("/books", handlers.GetAllBooks).Methods(http.MethodGet)

	log.Println("API is running!")
	http.ListenAndServe(":4000", router)
}
