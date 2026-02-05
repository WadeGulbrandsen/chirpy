package main

import (
	"log"
	"net/http"
)

const (
	port = "8080"
)

func main() {
	apiCfg := apiConfig{}
	apiCfg.fileserverHits.Store(0)
	mux := http.NewServeMux()
	mux.Handle(appPrefix, appHandler(&apiCfg))
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHandler)
	mux.HandleFunc("POST /admin/reset", apiCfg.resetMetricsHandler)
	mux.HandleFunc("GET /api/healthz", handleHealthz)
	mux.HandleFunc("POST /api/validate_chirp", validateChirpHandler)
	log.Printf("Serving files from %s on port: %s\n", appDir, port)
	log.Fatal(server.ListenAndServe())
}
