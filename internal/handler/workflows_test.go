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

func TestWorkflowsCreateDefinitionNoAuth(t *testing.T) {
	h := NewWorkflowsHandler(nil, nil)

	payload := `{"name":"Test Workflow"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workflows", bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	h.CreateDefinition(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestWorkflowsCreateDefinitionMissingName(t *testing.T) {
	h := NewWorkflowsHandler(nil, nil)

	user := &model.User{ID: uuid.New(), Name: "Admin", Email: "admin@test.com"}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, user)

	payload := `{"name":"","description":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/workflows", bytes.NewBufferString(payload))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.CreateDefinition(w, req)

	if w.Code != http.StatusUnprocessableEntity {
		t.Fatalf("expected 422, got %d", w.Code)
	}

	var result MutationResult
	_ = json.NewDecoder(w.Body).Decode(&result)
	if result.OK {
		t.Error("expected ok=false")
	}
}

func TestWorkflowsCreateDefinitionInvalidBody(t *testing.T) {
	h := NewWorkflowsHandler(nil, nil)

	user := &model.User{ID: uuid.New(), Name: "Admin", Email: "admin@test.com"}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, user)

	req := httptest.NewRequest(http.MethodPost, "/api/workflows", bytes.NewBufferString("bad json"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.CreateDefinition(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestWorkflowsGetDefinitionInvalidUUID(t *testing.T) {
	h := NewWorkflowsHandler(nil, nil)

	req := httptest.NewRequest(http.MethodGet, "/api/workflows/not-valid", nil)
	w := httptest.NewRecorder()

	h.GetDefinition(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestWorkflowsTransitionInvalidBody(t *testing.T) {
	h := NewWorkflowsHandler(nil, nil)

	user := &model.User{ID: uuid.New(), Name: "Admin", Email: "admin@test.com"}
	ctx := context.WithValue(context.Background(), middleware.UserContextKey, user)

	req := httptest.NewRequest(http.MethodPost, "/api/workflows/instances/some/transition", bytes.NewBufferString("not json"))
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	h.TransitionInstance(w, req)

	// chi.URLParam returns "" → uuid.Parse("") fails → 400
	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestWorkflowsUpdateDefinitionNoAuth(t *testing.T) {
	h := NewWorkflowsHandler(nil, nil)

	payload := `{"name":"Updated"}`
	req := httptest.NewRequest(http.MethodPut, "/api/workflows/"+uuid.New().String(), bytes.NewBufferString(payload))
	w := httptest.NewRecorder()

	h.UpdateDefinition(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}

func TestWorkflowsDeleteDefinitionInvalidUUID(t *testing.T) {
	h := NewWorkflowsHandler(nil, nil)

	req := httptest.NewRequest(http.MethodDelete, "/api/workflows/invalid", nil)
	w := httptest.NewRecorder()

	h.DeleteDefinition(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", w.Code)
	}
}

func TestWorkflowsCreateInstanceNoAuth(t *testing.T) {
	h := NewWorkflowsHandler(nil, nil)

	req := httptest.NewRequest(http.MethodPost, "/api/workflows/"+uuid.New().String()+"/instances", nil)
	w := httptest.NewRecorder()

	h.CreateInstance(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", w.Code)
	}
}
