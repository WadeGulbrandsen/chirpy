package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/WadeGulbrandsen/chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	port = "8080"
)

func main() {
	godotenv.Load()
	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open database connection")
	}
	apiCfg := apiConfig{}
	apiCfg.dbQueries = database.New(db)
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
