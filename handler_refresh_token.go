package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/trantuvan/chirpy/helpers"
	"github.com/trantuvan/chirpy/internal/auth"
	"github.com/trantuvan/chirpy/internal/database"
)

func (cfg *apiConfig) handlerGetUserFromRefreshToken(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("GetUserFromRefreshToken: %s", err), err)
		return
	}

	user, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, "GetUserFromRefreshToken: token doesn't exist", err)
		return
	}

	tokenJWT, errJWT := auth.MakeJWT(user.ID, cfg.secretKey, ExpiresTime)

	if errJWT != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, "GetUserFromRefreshToken: Invalid token JWT", errJWT)
		return
	}

	helpers.ResponseWithJson(w, http.StatusOK, response{Token: tokenJWT})
}

func (cfg *apiConfig) handlderRevokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)

	if err != nil {
		helpers.ResponseWithError(w, http.StatusUnauthorized, fmt.Sprintf("RevokeRefreshToken: %s", err), err)
		return
	}

	errRevokeToken := cfg.db.RevokeRefreshToken(r.Context(), database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{Time: time.Now(), Valid: true},
		UpdatedAt: time.Now(),
		Token:     refreshToken,
	})

	if errRevokeToken != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("RevokeRefreshToken: %s", err), err)
		return
	}

	helpers.ResponseWithJson(w, http.StatusNoContent, "revoke success")
}
