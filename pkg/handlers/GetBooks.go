package handlers

import (
	"encoding/json"
	"net/http"

	"code.crogge.rs/chris/lcwc_api/pkg/mocks"
)

func GetAllBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(mocks.Books)
}
