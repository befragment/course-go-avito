package handlers

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, httpStatus int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(data)
}

func respondWithError(w http.ResponseWriter, httpStatus int, message string) {
	respondWithJSON(w, httpStatus, map[string]string{"error": message})
}

