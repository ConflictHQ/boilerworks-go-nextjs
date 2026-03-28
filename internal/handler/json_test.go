package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWriteOK(t *testing.T) {
	w := httptest.NewRecorder()
	writeOK(w, map[string]string{"key": "value"})

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}
	if ct := w.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["key"] != "value" {
		t.Errorf("expected key=value, got %s", body["key"])
	}
}

func TestWriteCreated(t *testing.T) {
	w := httptest.NewRecorder()
	writeCreated(w, map[string]int{"id": 1})

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
}

func TestWriteError(t *testing.T) {
	w := httptest.NewRecorder()
	writeError(w, http.StatusNotFound, "Not found")

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["error"] != "Not found" {
		t.Errorf("expected 'Not found', got '%s'", body["error"])
	}
}

func TestWriteMutationOK(t *testing.T) {
	w := httptest.NewRecorder()
	writeMutationOK(w)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result MutationResult
	json.NewDecoder(w.Body).Decode(&result)
	if !result.OK {
		t.Error("expected ok=true")
	}
	if len(result.Errors) != 0 {
		t.Errorf("expected no errors, got %d", len(result.Errors))
	}
}

func TestWriteMutationErrors(t *testing.T) {
	w := httptest.NewRecorder()
	writeMutationErrors(w, []FieldError{
		{Field: "name", Messages: []string{"Name is required"}},
		{Field: "email", Messages: []string{"Email is required", "Email must be valid"}},
	})

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	var result MutationResult
	json.NewDecoder(w.Body).Decode(&result)
	if result.OK {
		t.Error("expected ok=false")
	}
	if len(result.Errors) != 2 {
		t.Fatalf("expected 2 field errors, got %d", len(result.Errors))
	}
	if result.Errors[1].Field != "email" {
		t.Errorf("expected second error on 'email', got '%s'", result.Errors[1].Field)
	}
	if len(result.Errors[1].Messages) != 2 {
		t.Errorf("expected 2 messages on email error, got %d", len(result.Errors[1].Messages))
	}
}
