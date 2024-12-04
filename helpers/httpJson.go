package helpers

import (
	"encoding/json"
	"net/http"
)

func ResponseWithJson(w http.ResponseWriter, statusCode int, payload interface{}) error {
	response, err := json.Marshal(payload)

	if err != nil {
		return err
	}

	w.Header().Add("Content-Type", "application/json")
	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.WriteHeader(statusCode)
	w.Write(response)
	return nil
}

func ResponseWithError(w http.ResponseWriter, statusCode int, msg string) error {
	return ResponseWithJson(w, statusCode, map[string]string{"error": msg})
}
