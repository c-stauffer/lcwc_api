package main

import (
	"log"
	"net/http"

	"code.crogge.rs/chris/lcwc_api/pkg/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/incidents", handlers.GetAllIncidents).Methods(http.MethodGet)

	c := cors.Default()
	handler := c.Handler(router)

	log.Println("API is running!")
	http.ListenAndServe(":4000", handler)
}
