package main

import (
	"fmt"
	"net/http"

	"github.com/trantuvan/chirpy/helpers"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf8")

	if cfg.platform != "dev" {
		w.Header().Add("Content-Type", "application/json")
		helpers.ResponseWithError(w, http.StatusForbidden, "handlerReset: Reset is only allow in dev environment", nil)
		return
	}

	if err := cfg.db.ResetUsers(r.Context()); err != nil {
		w.Header().Add("Content-Type", "application/json")
		helpers.ResponseWithError(w, http.StatusInternalServerError, fmt.Sprintf("handlerReset: cannot reset %s\n", err), err)
		return
	}

	w.WriteHeader(http.StatusOK)
	cfg.fileserverHits.Store(0)
	w.Write([]byte("Hits reset to 0"))
}
