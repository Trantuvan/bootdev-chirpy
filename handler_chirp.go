package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/trantuvan/chirpy/helpers"
	"github.com/trantuvan/chirpy/internal/auth"
	"github.com/trantuvan/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body string `json:"body"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	token, err := auth.GetBearerToken(r.Header)

	if token == "" {
		helpers.ResponseWithError(w, http.StatusUnauthorized, "there is no Authorization in headers", err)
		return
	}

	userID, errJWT := auth.ValidateJWT(token, cfg.secretKey)

	if errJWT != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, "failed to validate JWT", errJWT)
		return
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

	chirp, err := cfg.db.CreateChirps(r.Context(), database.CreateChirpsParams{Body: cleanedChirp, UserID: userID})

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

func (cfg *apiConfig) handlerDeleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("handlerDeleteChirp: %s", err), err)
		return
	}

	tokenJWT, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("handlerDeleteChirp: %s", err), err)
		return
	}

	userID, err := auth.ValidateJWT(tokenJWT, cfg.secretKey)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("handlerDeleteChirp: %s", err), err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)
	if err == sql.ErrNoRows {
		helpers.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("handlerDeleteChirp: chirp with ID - %s not exist", chirpID), err)
		return
	}
	if chirp.UserID != userID {
		helpers.ResponseWithError(w, http.StatusForbidden, fmt.Sprintf("handlerDeleteChirp: chirp ID - %s not belong to userID - %s", chirpID, userID), nil)
		return
	}

	errDel := cfg.db.DeleteChirpByID(r.Context(), database.DeleteChirpByIDParams{ID: chirpID, UserID: userID})
	if errDel != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerDeleteChirp: failed to del chirpID - %s", chirpID), err)
		return
	}

	helpers.ResponseWithJson(w, http.StatusNoContent, nil)
}

func (cfg *apiConfig) handlerGetChirp(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}
	chirpID, err := uuid.Parse(r.PathValue("chirpID"))

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("handlerGetChirp: %s\n", err), err)
		return
	}

	chirp, err := cfg.db.GetChirp(r.Context(), chirpID)

	if err == sql.ErrNoRows {
		helpers.ResponseWithError(w, http.StatusNotFound, fmt.Sprintf("handlerGetChirp: chirp with ID - %s not exist", chirpID), err)
		return
	}

	helpers.ResponseWithJson(w, http.StatusOK, response{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	})
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Body      string    `json:"body"`
		UserId    uuid.UUID `json:"user_id"`
	}

	sortQueryParam := r.URL.Query().Get("sort")

	if authorQueryParam := r.URL.Query().Get("author_id"); authorQueryParam != "" {
		userID, err := uuid.Parse(authorQueryParam)
		if err != nil {
			helpers.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("handlerGetChirps: failed parse userID %s\n", authorQueryParam), err)
			return
		}

		chirps, err := cfg.db.GetChirpsByUserID(r.Context(), userID)
		if err != nil {
			helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerGetChirps: failed to get chirps of userID %s\n", authorQueryParam), err)
			return
		}

		if sortQueryParam != "" && sortQueryParam == "desc" {
			// *func(i, j int) bool; Ask i < j true ?
			// *note order j index 0, i index 1
			// *if true swap i index 0, j index 1
			sort.Slice(chirps, func(i, j int) bool {
				return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
			})
		}

		//* map chirps to reponses
		responses := make([]response, len(chirps))

		for i, c := range chirps {
			responses[i] = response{
				ID:        c.ID,
				CreatedAt: c.CreatedAt,
				UpdatedAt: c.UpdatedAt,
				Body:      c.Body,
				UserId:    c.UserID,
			}
		}

		helpers.ResponseWithJson(w, http.StatusOK, responses)
		return
	}

	chirps, err := cfg.db.GetChirps(r.Context())

	if err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerGetChirps: failed to get chirps %s\n", err), err)
		return
	}

	if sortQueryParam != "" && sortQueryParam == "desc" {
		sort.Slice(chirps, func(i, j int) bool {
			return chirps[i].CreatedAt.After(chirps[j].CreatedAt)
		})
	}

	//* map chirps to reponses
	responses := make([]response, len(chirps))

	for i, c := range chirps {
		responses[i] = response{
			ID:        c.ID,
			CreatedAt: c.CreatedAt,
			UpdatedAt: c.UpdatedAt,
			Body:      c.Body,
			UserId:    c.UserID,
		}
	}

	helpers.ResponseWithJson(w, http.StatusOK, responses)
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
