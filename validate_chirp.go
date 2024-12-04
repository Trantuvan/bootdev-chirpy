package main

import (
	"encoding/json"
	"net/http"

	"github.com/trantuvan/chirpy/helpers"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	const maxChirpLength = 140
	params := struct {
		Body string `json:"body"`
	}{}
	validResult := struct {
		Valid bool `json:"valid"`
	}{true}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, "Validate Chirp: cannot parse json", err)
		return
	}

	if len(params.Body) > maxChirpLength {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Validate Chirp: Chirp is too long", nil)
		return
	}

	helpers.ResponseWithJson(w, http.StatusOK, validResult)
}
