package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/model"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type FormsHandler struct {
	forms   *queries.FormQueries
	formSvc *service.FormService
}

func NewFormsHandler(forms *queries.FormQueries, formSvc *service.FormService) *FormsHandler {
	return &FormsHandler{forms: forms, formSvc: formSvc}
}

type formDefinitionRequest struct {
	Name        string            `json:"name"`
	Slug        string            `json:"slug"`
	Description string            `json:"description"`
	Status      string            `json:"status"`
	Schema      []model.FormField `json:"schema"`
}

type formSubmissionRequest struct {
	Data map[string]string `json:"data"`
}

func (h *FormsHandler) ListDefinitions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	forms, total, err := h.forms.ListDefinitions(r.Context(), perPage, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list forms")
		return
	}

	if forms == nil {
		forms = []model.FormDefinition{}
	}

	writeOK(w, map[string]interface{}{
		"forms":    forms,
		"total":    total,
		"page":     page,
		"per_page": perPage,
	})
}

func (h *FormsHandler) GetDefinition(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	form, err := h.forms.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Form not found")
		return
	}

	writeOK(w, form)
}

func (h *FormsHandler) CreateDefinition(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req formDefinitionRequest
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

	if req.Slug == "" {
		req.Slug = strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	}

	if req.Status == "" {
		req.Status = "draft"
	}

	schemaJSON, err := json.Marshal(req.Schema)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid schema")
		return
	}

	form, err := h.forms.CreateDefinition(r.Context(), req.Name, req.Slug, req.Description, req.Status, schemaJSON, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create form")
		return
	}

	writeCreated(w, form)
}

func (h *FormsHandler) UpdateDefinition(w http.ResponseWriter, r *http.Request) {
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

	var req formDefinitionRequest
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

	schemaJSON, err := json.Marshal(req.Schema)
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid schema")
		return
	}

	form, err := h.forms.UpdateDefinition(r.Context(), uid, req.Name, req.Slug, req.Description, req.Status, schemaJSON, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update form")
		return
	}

	writeOK(w, form)
}

func (h *FormsHandler) DeleteDefinition(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.forms.DeleteDefinition(r.Context(), uid); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete form")
		return
	}

	writeMutationOK(w)
}

func (h *FormsHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	form, err := h.forms.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Form not found")
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	submissions, total, err := h.forms.ListSubmissions(r.Context(), form.ID, perPage, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list submissions")
		return
	}

	if submissions == nil {
		submissions = []model.FormSubmission{}
	}

	writeOK(w, map[string]interface{}{
		"submissions": submissions,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
	})
}

func (h *FormsHandler) CreateSubmission(w http.ResponseWriter, r *http.Request) {
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

	form, err := h.forms.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Form not found")
		return
	}

	var req formSubmissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	jsonData, validationErrs := h.formSvc.ValidateSubmission(form, req.Data)
	if len(validationErrs) > 0 {
		writeJSON(w, http.StatusUnprocessableEntity, map[string]interface{}{
			"ok":     false,
			"errors": validationErrs,
		})
		return
	}

	submission, err := h.forms.CreateSubmission(r.Context(), form.ID, jsonData, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create submission")
		return
	}

	writeCreated(w, submission)
}
