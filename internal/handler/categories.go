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

type CategoriesHandler struct {
	categories *queries.CategoryQueries
}

func NewCategoriesHandler(categories *queries.CategoryQueries) *CategoriesHandler {
	return &CategoriesHandler{categories: categories}
}

type categoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (h *CategoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	categories, total, err := h.categories.List(r.Context(), perPage, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list categories")
		return
	}

	if categories == nil {
		categories = []model.Category{}
	}

	writeOK(w, map[string]interface{}{
		"categories": categories,
		"total":      total,
		"page":       page,
		"per_page":   perPage,
	})
}

func (h *CategoriesHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	category, err := h.categories.GetByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Category not found")
		return
	}

	writeOK(w, category)
}

func (h *CategoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req categoryRequest
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

	category, err := h.categories.Create(r.Context(), req.Name, req.Description, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create category")
		return
	}

	writeCreated(w, category)
}

func (h *CategoriesHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req categoryRequest
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

	category, err := h.categories.Update(r.Context(), uid, req.Name, req.Description, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update category")
		return
	}

	writeOK(w, category)
}

func (h *CategoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.categories.Delete(r.Context(), uid); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete category")
		return
	}

	writeMutationOK(w)
}
