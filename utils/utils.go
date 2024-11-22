package utils

import (
	"encoding/json"
	"log"
	"net/http"
)

func ParseJSON(r *http.Request, payload any) error {
	if r.Body == nil {
		log.Fatal("")
	}
	return json.NewDecoder(r.Body).Decode(payload)
}

func WrightJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

func WrightError(w http.ResponseWriter, status int, err error) {
	WrightJSON(w, status, map[string]string{"error": err.Error()})
}
