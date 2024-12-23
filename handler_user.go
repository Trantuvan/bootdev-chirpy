package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/trantuvan/chirpy/helpers"
	"github.com/trantuvan/chirpy/internal/auth"
	"github.com/trantuvan/chirpy/internal/database"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	params := parameter{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerCreateUser: failed to read params %s\n", err), err)
		return
	}

	hashedPass, err := auth.HashPassword(params.Password)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, "handlerCreateUser: bad password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Email:          params.Email,
		HashedPassword: sql.NullString{String: hashedPass, Valid: err == nil},
	})

	if err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerCreateUser: failed to create user %s\n", err), err)
		return
	}

	helpers.ResponseWithJson(w, http.StatusCreated, response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}

func (cfg *apiConfig) hanlderLogin(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Email     string         `json:"email"`
		Password  string         `json:"password"`
		ExpiresAt *time.Duration `json:"expires_in_seconds"` //optional pointer allow nil
	}
	type response struct {
		ID        uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email     string    `json:"email"`
	}

	params := parameter{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerLogin: failed to read params %s", err), err)
		return
	}

	expiresAtDefault := time.Hour

	if params.ExpiresAt == nil || *params.ExpiresAt > expiresAtDefault {
		params.ExpiresAt = &expiresAtDefault
	}

	user, err := cfg.db.GetUserByEmail(r.Context(), params.Email)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("handlerLogin: Incorrect email -  %s", params.Email), err)
		return
	}

	if errPass := auth.CheckPasswordHash(params.Password, user.HashedPassword.String); errPass != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, "handlerLogin: Incorrect password", err)
		return
	}

	helpers.ResponseWithJson(w, http.StatusOK, response{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	})
}
