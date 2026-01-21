package utils

import (
	"log"
	"net/http"
)

const (
	ErrInternalServer = "Internal server error"
)

func RespondWithError(w http.ResponseWriter, httpStatus int, message string) {
	RespondWithJSON(w, httpStatus, map[string]string{"error": message})
}

func RespondInternalServerError(w http.ResponseWriter, err error) {
	log.Printf("Internal server error: %v\n", err)
	RespondWithError(w, http.StatusInternalServerError, ErrInternalServer)
}
