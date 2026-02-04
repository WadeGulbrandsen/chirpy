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
	dir := http.Dir(".")
	mux.Handle("/", http.FileServer(dir))
	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}
	log.Printf("Serving files from %s on port: %s\n", dir, port)
	log.Fatal(server.ListenAndServe())
}
