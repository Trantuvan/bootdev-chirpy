package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port string = "8080"
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("./assets")))

	server := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: mux}

	log.Printf("Servering on port: %s\n", port)
	log.Fatalf("Server failed: %s", server.ListenAndServe())
}
