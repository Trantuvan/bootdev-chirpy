package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/trantuvan/chirpy/helpers"
)

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Email string `json:"email"`
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

	user, err := cfg.db.CreateUser(r.Context(), params.Email)

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
