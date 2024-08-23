package main

import (
	"log"
	"net/http"
)

func main() {
	const addr = "localhost:8080"
	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Printf("Serving on port: %s\n", addr)
	log.Fatal(server.ListenAndServe())
}
