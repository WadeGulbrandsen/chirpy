package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, status_code int, message string, err error) {
	if err != nil {
		log.Println(err)
	}
	if status_code > 499 {
		log.Printf("Responding with 5XX error: %s", message)
	}
	type jsonError struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, status_code, jsonError{
		Error: message,
	})
}

func respondWithJSON(w http.ResponseWriter, status_code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status_code)
	w.Write(data)
}
