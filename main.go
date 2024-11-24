package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port string = "8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	})
	server := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: mux}

	log.Printf("Servering on port: %s\n", port)
	log.Fatalf("Server failed: %s", server.ListenAndServe())
}
