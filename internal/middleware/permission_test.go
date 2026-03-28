package middleware

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequirePermissionGranted(t *testing.T) {
	perms := []string{"products.view", "products.create", "categories.view"}
	ctx := context.WithValue(context.Background(), PermissionsContextKey, perms)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"ok":true}`))
	})

	handler := RequirePermission("products.view")(inner)

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 when permission granted, got %d", w.Code)
	}
}

func TestRequirePermissionDenied(t *testing.T) {
	// Viewer permissions — no create or delete
	perms := []string{"products.view", "categories.view"}
	ctx := context.WithValue(context.Background(), PermissionsContextKey, perms)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("inner handler should not be called when permission denied")
	})

	handler := RequirePermission("products.create")(inner)

	req := httptest.NewRequest(http.MethodPost, "/api/products", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 when permission denied, got %d", w.Code)
	}

	var body map[string]string
	json.NewDecoder(w.Body).Decode(&body)
	if body["error"] != "Forbidden" {
		t.Errorf("expected error 'Forbidden', got '%s'", body["error"])
	}
}

func TestRequirePermissionDeleteDenied(t *testing.T) {
	perms := []string{"products.view", "categories.view"}
	ctx := context.WithValue(context.Background(), PermissionsContextKey, perms)

	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("inner handler should not be called")
	})

	handler := RequirePermission("products.delete")(inner)

	req := httptest.NewRequest(http.MethodDelete, "/api/products/some-uuid", nil)
	req = req.WithContext(ctx)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403, got %d", w.Code)
	}
}

func TestRequirePermissionNoPermissions(t *testing.T) {
	// Empty context — no permissions at all
	inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Fatal("inner handler should not be called")
	})

	handler := RequirePermission("products.view")(inner)

	req := httptest.NewRequest(http.MethodGet, "/api/products", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for nil permissions, got %d", w.Code)
	}
}
