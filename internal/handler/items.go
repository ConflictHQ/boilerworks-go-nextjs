package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ItemsHandler struct {
	items   *queries.ItemQueries
	categories *queries.CategoryQueries
}

func NewItemsHandler(items *queries.ItemQueries, categories *queries.CategoryQueries) *ItemsHandler {
	return &ItemsHandler{items: items, categories: categories}
}

type itemRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	CategoryID  *string `json:"category_id,omitempty"`
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	items, total, err := h.items.List(r.Context(), perPage, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list items")
		return
	}

	if items == nil {
		items = []model.Item{}
	}

	writeOK(w, map[string]interface{}{
		"items": items,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *ItemsHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	item, err := h.items.GetByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Item not found")
		return
	}

	writeOK(w, item)
}

func (h *ItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req itemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		writeMutationErrors(w, []FieldError{
			{Field: "name", Messages: []string{"Name is required"}},
		})
		return
	}

	if req.Status == "" {
		req.Status = "draft"
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil && *req.CategoryID != "" {
		id, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid category ID")
			return
		}
		categoryID = &id
	}

	item, err := h.items.Create(r.Context(), req.Name, req.Description, req.Price, req.Status, categoryID, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create item")
		return
	}

	writeCreated(w, item)
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	var req itemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if req.Name == "" {
		writeMutationErrors(w, []FieldError{
			{Field: "name", Messages: []string{"Name is required"}},
		})
		return
	}

	var categoryID *uuid.UUID
	if req.CategoryID != nil && *req.CategoryID != "" {
		id, err := uuid.Parse(*req.CategoryID)
		if err != nil {
			writeError(w, http.StatusBadRequest, "Invalid category ID")
			return
		}
		categoryID = &id
	}

	item, err := h.items.Update(r.Context(), uid, req.Name, req.Description, req.Price, req.Status, categoryID, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update item")
		return
	}

	writeOK(w, item)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.items.Delete(r.Context(), uid); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete item")
		return
	}

	writeMutationOK(w)
}
