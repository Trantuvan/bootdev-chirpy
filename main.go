package main

import (
	"log"
	"net/http"
	"sync/atomic"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const filepathRoot = "."
	const port = "8080"
	apiConfig := apiConfig{}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))
	mux.HandleFunc("GET /healthz", handlerReadiness)
	mux.HandleFunc("GET /metrics", apiConfig.handlerMetrics)
	mux.HandleFunc("POST /reset", apiConfig.handlerReset)

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	log.Printf("Serving files from %s on port: %s\n", filepathRoot, port)
	log.Fatal(srv.ListenAndServe())
}
