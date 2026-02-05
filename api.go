package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync/atomic"

	"github.com/WadeGulbrandsen/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
}

func handleHealthz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	template := `<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited %d times!</p>
  </body>
</html>`
	html_text := fmt.Sprintf(template, cfg.fileserverHits.Load())
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(html_text))
}

func (cfg *apiConfig) resetMetricsHandler(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

func validateChirpHandler(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}
	type jsonError struct {
		ErrorMsg string `json:"error"`
	}
	type jsonSuccess struct {
		CleanedBody string `json:"cleaned_body"`
	}

	bad_words := map[string]struct{}{
		"kerfuffle": struct{}{},
		"sharbert":  struct{}{},
		"fornax":    struct{}{},
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err == nil && len(params.Body) > 140 {
		err = fmt.Errorf("Chirp is too long")
	}
	var data []byte
	var status int
	if err != nil {
		status = http.StatusBadRequest
		data, err = json.Marshal(jsonError{ErrorMsg: err.Error()})
	} else {
		status = http.StatusOK
		cleaned_words := []string{}
		for _, word := range strings.Split(params.Body, " ") {
			if _, ok := bad_words[strings.ToLower(word)]; ok {
				cleaned_words = append(cleaned_words, "****")
			} else {
				cleaned_words = append(cleaned_words, word)
			}
		}
		data, err = json.Marshal(
			jsonSuccess{CleanedBody: strings.Join(cleaned_words, " ")},
		)
	}
	if err != nil {
		log.Panicf("Error marshalling JSON: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(data)
}
