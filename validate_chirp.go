package main

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/trantuvan/chirpy/helpers"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameter struct {
		Body string `json:"body"`
	}
	type result struct {
		CleanedBody string `json:"cleaned_body"`
	}
	const maxChirpLength = 140
	profaneWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	params := parameter{}

	if err := decoder.Decode(&params); err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, "Validate Chirp: cannot parse json", err)
		return
	}

	if len(params.Body) > maxChirpLength {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Validate Chirp: Chirp is too long", nil)
		return
	}

	origins := strings.Split(params.Body, " ")

	for i, w := range origins {
		if _, ok := profaneWords[strings.ToLower(w)]; ok {
			origins[i] = "****"
		}
	}

	helpers.ResponseWithJson(w, http.StatusOK, result{strings.Join(origins, " ")})
}
