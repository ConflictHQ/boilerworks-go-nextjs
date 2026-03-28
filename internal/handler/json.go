package handler

import (
	"encoding/json"
	"net/http"
)

// MutationResult is the standard response for create/update/delete operations.
type MutationResult struct {
	OK     bool         `json:"ok"`
	Errors []FieldError `json:"errors,omitempty"`
}

// FieldError represents a field-level validation error.
type FieldError struct {
	Field    string   `json:"field"`
	Messages []string `json:"messages"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeOK(w http.ResponseWriter, v interface{}) {
	writeJSON(w, http.StatusOK, v)
}

func writeCreated(w http.ResponseWriter, v interface{}) {
	writeJSON(w, http.StatusCreated, v)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

func writeMutationOK(w http.ResponseWriter) {
	writeJSON(w, http.StatusOK, MutationResult{OK: true})
}

func writeMutationErrors(w http.ResponseWriter, errors []FieldError) {
	writeJSON(w, http.StatusUnprocessableEntity, MutationResult{OK: false, Errors: errors})
}
