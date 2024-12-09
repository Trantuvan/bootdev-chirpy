package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/trantuvan/chirpy/helpers"
	"github.com/trantuvan/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	params := parameter{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerCreateChirp: failed to read params %s\n", err), err)
		return
	}

	cleanedChirp, err := validateChirp(params.Body)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("handlerCreateChirp: %s\n", err), err)
		return
	}

	chirp, err := cfg.db.CreateChirps(r.Context(), database.CreateChirpsParams{Body: cleanedChirp, UserID: params.UserId})

	if err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerCreateChirp: failed to create chirp %s\n", err), err)
		return
	}

	helpers.ResponseWithJson(w, http.StatusCreated, response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func validateChirp(msg string) (string, error) {
	const maxChirpLength = 140
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	if len(msg) > maxChirpLength {
		return "", errors.New("chirp is too long")
	}

	origins := strings.Split(msg, " ")

	for i, w := range origins {
		if _, ok := profaneWords[strings.ToLower(w)]; ok {
			origins[i] = "****"
		}
	}

	return strings.Join(origins, " "), nil
}
