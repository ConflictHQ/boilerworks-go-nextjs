package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/model"
	"github.com/google/uuid"
)

func TestCategoriesCreateNoAuth(t *testing.T) {
	h := NewCategoriesHandler(nil)

	payload := `{"name":"Test Category"}`
	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCategoriesCreateMissingName(t *testing.T) {
	h := NewCategoriesHandler(nil)

	user := &model.User{ID: uuid.New(), Name: "Admin", Email: "admin@test.com"}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, user)

	payload := `{"name":"","description":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	var result MutationResult
	json.NewDecoder(w.Body).Decode(&result)
	if result.OK {
		t.Error("expected ok=false for missing name")
	}
	if len(result.Errors) == 0 || result.Errors[0].Field != "name" {
		t.Error("expected error on 'name' field")
	}
}

func TestCategoriesCreateInvalidBody(t *testing.T) {
	h := NewCategoriesHandler(nil)

	user := &model.User{ID: uuid.New(), Name: "Admin", Email: "admin@test.com"}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, user)

	req := httptest.NewRequest(http.MethodPost, "/api/categories", bytes.NewBufferString("{invalid"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.Create(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCategoriesGetInvalidUUID(t *testing.T) {
	h := NewCategoriesHandler(nil)

	req := httptest.NewRequest(http.MethodGet, "/api/categories/bad-uuid", nil)
	w := httptest.NewRecorder()

	h.Get(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestCategoriesUpdateNoAuth(t *testing.T) {
	h := NewCategoriesHandler(nil)

	payload := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/categories/"+uuid.New().String(), bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	h.Update(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestCategoriesDeleteInvalidUUID(t *testing.T) {
	h := NewCategoriesHandler(nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/categories/not-valid", nil)
	w := httptest.NewRecorder()

	h.Delete(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}
