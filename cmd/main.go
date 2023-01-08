package main

import (
	"log"
	"net/http"

	"code.crogge.rs/chris/lcwc_api/pkg/handlers"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/incidents", handlers.GetAllIncidents).Methods(http.MethodGet)

	log.Println("API is running!")
	http.ListenAndServe(":4000", router)
}
