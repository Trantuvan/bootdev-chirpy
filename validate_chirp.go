package main

import (
	"encoding/json"
	"net/http"

	"github.com/trantuvan/chirpy/helpers"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	params := struct {
		Body string `json:"body"`
	}{}
	validResult := struct {
		Valid bool `json:"valid"`
	}{true}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(&params); err != nil {
		helpers.ResponseWithError(w, http.StatusInternalServerError, "Validate Chirp: cannot parse json")
		return
	}

	if len(params.Body) > 140 {
		helpers.ResponseWithError(w, http.StatusBadRequest, "Validate Chirp: Chirp is too long")
		return
	}

	helpers.ResponseWithJson(w, http.StatusOK, validResult)
}
