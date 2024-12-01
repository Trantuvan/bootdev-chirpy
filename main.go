package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port string = "8080"
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", hanlderReadiness)
	mux.HandleFunc("/metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("/reset", apiConfig.handlerReset)
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!\n")
	})

	server := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: mux}

	log.Printf("Servering on port: %s\n", port)
	log.Fatalf("Server failed: %s", server.ListenAndServe())
}
