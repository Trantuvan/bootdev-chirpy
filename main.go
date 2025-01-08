package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/trantuvan/chirpy/internal/database"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	db             *database.Queries
	platform       string
	secretKey      string
	polkaKey       string
}

func main() {
	err := godotenv.Load() //* load .env file into env of this process

	if err != nil {
		log.Fatalf("cannot load connection string: %s\n", err)
	}

	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")
	secretKey := os.Getenv("SECRET_KEY")
	polkaKey := os.Getenv("POLKA_KEY")

	if dbURL == "" || platform == "" || secretKey == "" || polkaKey == "" {
		log.Fatal("DB_URL & PLATFORM & SECRET_KEY & polkaKey must be set")
	}

	db, err := sql.Open("postgres", dbURL)

	if err != nil {
		log.Fatalf("cannot open database: %s\n", err)
	}

	const port string = "8080"
	apiConfig := apiConfig{
		fileserverHits: atomic.Int32{},
		db:             database.New(db),
		platform:       platform,
		secretKey:      secretKey,
		polkaKey:       polkaKey,
	}

	mux := http.NewServeMux()
	mux.Handle("/app/", apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(".")))))

	mux.HandleFunc("GET /api/healthz", handlerReadiness) // only GET
	mux.HandleFunc("POST /api/refresh", apiConfig.handlerGetUserFromRefreshToken)
	mux.HandleFunc("POST /api/revoke", apiConfig.handlderRevokeRefreshToken)
	mux.HandleFunc("POST /api/polka/webhooks", apiConfig.handlerUpdateUserToChirpyRed)
	mux.HandleFunc("POST /api/login", apiConfig.hanlderLogin)

	mux.HandleFunc("POST /api/users", apiConfig.handlerCreateUser)
	mux.HandleFunc("PUT /api/users", apiConfig.handlerUpdateUserEmailPassword)

	mux.HandleFunc("POST /api/chirps", apiConfig.handlerCreateChirp)
	mux.HandleFunc("GET /api/chirps", apiConfig.handlerGetChirps)
	mux.HandleFunc("GET /api/chirps/{chirpID}", apiConfig.handlerGetChirp)
	mux.HandleFunc("DELETE /api/chirps/{chirpID}", apiConfig.handlerDeleteChirp)

	mux.HandleFunc("GET /admin/metrics", apiConfig.handlerMetrics) // only GET
	mux.HandleFunc("POST /admin/reset", apiConfig.handlerReset)    // only POST

	server := http.Server{Addr: fmt.Sprintf(":%s", port), Handler: mux}

	log.Printf("Servering on port: %s\n", port)
	log.Fatalf("Server failed: %s", server.ListenAndServe())
}
