package handlers

import (
	"encoding/json"
	"net/http"
)

// respondError sends a JSON error response with a given status code and message.
func respondError(w http.ResponseWriter, statusCode int, message string) {
	respondJSON(w, statusCode, map[string]string{
		"error": message,
	})
}

// respondJSON sends any payload as a JSON response with a given status code.
func respondJSON(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(payload)
}
