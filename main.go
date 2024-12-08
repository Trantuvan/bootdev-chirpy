package main

import (
	"fmt"
	"log"
	"net/http"
	"sync/atomic"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
}

func main() {
	const port string = "8080"
	apiConfig := apiConfig{atomic.Int32{}}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness) // only GET
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics) // only GET
	mux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)    // only POST

	server := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: mux}

	log.Printf("Servering on port: %s\n", port)
	log.Fatalf("Server failed: %s", server.ListenAndServe())
}
