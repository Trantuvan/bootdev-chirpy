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

const ExpiresTime = time.Second * 3600

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

func (cfg *apiConfig) handlerUpdateUserEmailPassword(w http.ResponseWriter, r *http.Request) {
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

	tokenJWT, err := auth.GetBearerToken(r.Header)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("UpdateUserEmailPassword: %s", err), err)
		return
	}

	userID, err := auth.ValidateJWT(tokenJWT, cfg.secretKey)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("UpdateUserEmailPassword: %s", err), err)
		return
	}

	params := parameter{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("UpdateUserEmailPassword: failed to read params %s", err), err)
		return
	}

	hasedPass, err := auth.HashPassword(params.Password)
	if err != nil {
		helpers.ResponseWithError(w, http.StatusBadRequest, fmt.Sprintf("UpdateUserEmailPassword: bad password %s", err), err)
		return
	}

	updatedUser, err := cfg.db.UpdateUserEmailPassword(r.Context(), database.UpdateUserEmailPasswordParams{
		Email:          params.Email,
		HashedPassword: sql.NullString{String: hasedPass, Valid: err != nil},
		ID:             userID,
	})
	if err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("UpdateUserEmailPassword: failed to update user %s", err), err)
		return
	}

	helpers.ResponseWithJson(w, http.StatusOK, response{
		ID:        updatedUser.ID,
		CreatedAt: updatedUser.UpdatedAt,
		UpdatedAt: updatedUser.UpdatedAt,
		Email:     updatedUser.Email,
	})
}

func (cfg *apiConfig) hanlderLogin(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	type response struct {
		ID           uuid.UUID `json:"id"`
		CreatedAt    time.Time `json:"created_at"`
		UpdatedAt    time.Time `json:"updated_at"`
		Email        string    `json:"email"`
		Token        string    `json:"token"`
		RefreshToken string    `json:"refresh_token"`
	}

	params := parameter{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerLogin: failed to read params %s", err), err)
		return
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

	tokenJWT, errTokenJWT := auth.MakeJWT(user.ID, cfg.secretKey, ExpiresTime)

	if errTokenJWT != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, "handlerLogin: Invalid token JWT", errTokenJWT)
		return
	}

	tokenRefresh, errRefresh := auth.MakeRefreshToken()

	if errRefresh != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, "handlerLogin: Invalid token Refresh", errRefresh)
		return
	}

	refreshToken, errRefreshToken := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		Token:     tokenRefresh,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		UserID:    user.ID,
		ExpiresAt: time.Now().Add(60 * 24 * time.Hour), // 60 days
	})

	if errRefreshToken != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("handlerLogin: failed to create token %s", errRefreshToken), errRefreshToken)
		return
	}

	helpers.ResponseWithJson(w, http.StatusOK, response{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAt:    user.UpdatedAt,
		Email:        user.Email,
		Token:        tokenJWT,
		RefreshToken: refreshToken.Token,
	})
}
