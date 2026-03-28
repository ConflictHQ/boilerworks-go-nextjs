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

type ProductsHandler struct {
	products   *queries.ProductQueries
	categories *queries.CategoryQueries
}

func NewProductsHandler(products *queries.ProductQueries, categories *queries.CategoryQueries) *ProductsHandler {
	return &ProductsHandler{products: products, categories: categories}
}

type productRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Status      string  `json:"status"`
	CategoryID  *string `json:"category_id,omitempty"`
}

func (h *ProductsHandler) List(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	products, total, err := h.products.List(r.Context(), perPage, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list products")
		return
	}

	if products == nil {
		products = []model.Product{}
	}

	writeOK(w, map[string]interface{}{
		"products": products,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *ProductsHandler) Get(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	product, err := h.products.GetByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Product not found")
		return
	}

	writeOK(w, product)
}

func (h *ProductsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req productRequest
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

	product, err := h.products.Create(r.Context(), req.Name, req.Description, req.Price, req.Status, categoryID, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	writeCreated(w, product)
}

func (h *ProductsHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	var req productRequest
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

	product, err := h.products.Update(r.Context(), uid, req.Name, req.Description, req.Price, req.Status, categoryID, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	writeOK(w, product)
}

func (h *ProductsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.products.Delete(r.Context(), uid); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	writeMutationOK(w)
}
