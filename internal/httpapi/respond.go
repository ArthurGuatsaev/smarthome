package httpapi

import (
	"encoding/json"
	"net/http"
)

type apiError struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, code int, msg string, details string) {
	writeJSON(w, code, apiError{Error: msg, Details: details})
}
