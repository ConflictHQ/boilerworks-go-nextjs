package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthLoginMissingBody(t *testing.T) {
	// AuthHandler.Login reads JSON body; sending empty body should return 400
	h := &AuthHandler{authSvc: nil}

	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewReader([]byte("")))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["error"] != "Invalid request body" {
		t.Errorf("unexpected error message: %s", body["error"])
	}
}

func TestAuthLoginEmptyFields(t *testing.T) {
	h := &AuthHandler{authSvc: nil}

	payload := `{"email":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/login", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	h.Login(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	var result MutationResult
	json.NewDecoder(w.Body).Decode(&result)
	if result.OK {
		t.Error("expected ok=false for empty fields")
	}
	if len(result.Errors) < 2 {
		t.Errorf("expected at least 2 field errors, got %d", len(result.Errors))
	}
}

func TestAuthRegisterMissingBody(t *testing.T) {
	h := &AuthHandler{authSvc: nil}

	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewReader([]byte("not json")))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestAuthRegisterEmptyFields(t *testing.T) {
	h := &AuthHandler{authSvc: nil}

	payload := `{"name":"","email":"","password":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/auth/register", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	h.Register(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	var result MutationResult
	json.NewDecoder(w.Body).Decode(&result)
	if result.OK {
		t.Error("expected ok=false for empty fields")
	}
	if len(result.Errors) < 3 {
		t.Errorf("expected at least 3 field errors, got %d", len(result.Errors))
	}
}

func TestAuthLogoutNoSession(t *testing.T) {
	// Logout without a cookie should still return ok
	h := &AuthHandler{authSvc: nil}

	req := httptest.NewRequest(http.MethodPost, "/api/auth/logout", nil)
	w := httptest.NewRecorder()

	h.Logout(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	var result MutationResult
	json.NewDecoder(w.Body).Decode(&result)
	if !result.OK {
		t.Error("expected ok=true for logout without session")
	}

	// Check that session cookie was cleared
	cookies := w.Result().Cookies()
	var found bool
	for _, c := range cookies {
		if c.Name == "session_token" && c.MaxAge < 0 {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected session_token cookie to be cleared")
	}
}

func TestAuthMeNoUser(t *testing.T) {
	h := &AuthHandler{authSvc: nil}

	// Request without user in context
	req := httptest.NewRequest(http.MethodGet, "/api/auth/me", nil)
	w := httptest.NewRecorder()

	h.Me(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
