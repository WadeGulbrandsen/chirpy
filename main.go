package main

import (
	"log"
	"net/http"
)

const (
	port = "8080"
)

func main() {
	mux := http.NewServeMux()
	appDir := http.Dir(".")
	mux.Handle("/app/", http.StripPrefix("/app/", http.FileServer(appDir)))
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	log.Printf("Serving files from %s on port: %s\n", appDir, port)
	log.Fatal(server.ListenAndServe())
}
