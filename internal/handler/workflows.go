package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/model"
	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type WorkflowsHandler struct {
	workflows   *queries.WorkflowQueries
	workflowSvc *service.WorkflowService
}

func NewWorkflowsHandler(workflows *queries.WorkflowQueries, workflowSvc *service.WorkflowService) *WorkflowsHandler {
	return &WorkflowsHandler{workflows: workflows, workflowSvc: workflowSvc}
}

type workflowDefinitionRequest struct {
	Name        string                     `json:"name"`
	Description string                     `json:"description"`
	Status      string                     `json:"status"`
	States      []model.WorkflowState      `json:"states"`
	Transitions []model.WorkflowTransition `json:"transitions"`
}

type transitionRequest struct {
	Transition string `json:"transition"`
}

func (h *WorkflowsHandler) ListDefinitions(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage, _ := strconv.Atoi(r.URL.Query().Get("per_page"))
	if perPage < 1 {
		perPage = 20
	}
	offset := (page - 1) * perPage

	defs, total, err := h.workflows.ListDefinitions(r.Context(), perPage, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list workflows")
		return
	}

	if defs == nil {
		defs = []model.WorkflowDefinition{}
	}

	writeOK(w, map[string]interface{}{
		"workflows": defs,
		"total":     total,
		"page":      page,
		"per_page":  perPage,
	})
}

func (h *WorkflowsHandler) GetDefinition(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	def, err := h.workflows.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Workflow not found")
		return
	}

	writeOK(w, def)
}

func (h *WorkflowsHandler) CreateDefinition(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if user == nil {
		writeError(w, http.StatusUnauthorized, "Authentication required")
		return
	}

	var req workflowDefinitionRequest
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

	statesJSON, _ := json.Marshal(req.States)
	transitionsJSON, _ := json.Marshal(req.Transitions)

	def, err := h.workflows.CreateDefinition(r.Context(), req.Name, req.Description, req.Status, statesJSON, transitionsJSON, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create workflow")
		return
	}

	writeCreated(w, def)
}

func (h *WorkflowsHandler) UpdateDefinition(w http.ResponseWriter, r *http.Request) {
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

	var req workflowDefinitionRequest
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

	statesJSON, _ := json.Marshal(req.States)
	transitionsJSON, _ := json.Marshal(req.Transitions)

	def, err := h.workflows.UpdateDefinition(r.Context(), uid, req.Name, req.Description, req.Status, statesJSON, transitionsJSON, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to update workflow")
		return
	}

	writeOK(w, def)
}

func (h *WorkflowsHandler) DeleteDefinition(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	if err := h.workflows.DeleteDefinition(r.Context(), uid); err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to delete workflow")
		return
	}

	writeMutationOK(w)
}

func (h *WorkflowsHandler) ListInstances(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	def, err := h.workflows.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Workflow not found")
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

	instances, total, err := h.workflows.ListInstances(r.Context(), def.ID, perPage, offset)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to list instances")
		return
	}

	if instances == nil {
		instances = []model.WorkflowInstance{}
	}

	writeOK(w, map[string]interface{}{
		"instances": instances,
		"total":     total,
		"page":      page,
		"per_page":  perPage,
	})
}

func (h *WorkflowsHandler) CreateInstance(w http.ResponseWriter, r *http.Request) {
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

	def, err := h.workflows.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Workflow not found")
		return
	}

	initialState, err := h.workflowSvc.GetInitialState(def)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	instance, err := h.workflows.CreateInstance(r.Context(), def.ID, initialState, user.ID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to create instance")
		return
	}

	writeCreated(w, instance)
}

func (h *WorkflowsHandler) GetInstance(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "Invalid UUID")
		return
	}

	instance, err := h.workflows.GetInstanceByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Instance not found")
		return
	}

	def, err := h.workflows.GetDefinitionByID(r.Context(), instance.WorkflowDefinitionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get workflow definition")
		return
	}

	// Note: GetDefinitionByUUID uses the definition's UUID, but instance stores the internal ID.
	// We need to look up by internal ID. For now, fetch available transitions separately.
	availableTransitions := h.workflowSvc.GetAvailableTransitions(def, instance.CurrentState)

	logs, _ := h.workflows.GetTransitionLogs(r.Context(), instance.ID)
	if logs == nil {
		logs = []model.TransitionLog{}
	}

	writeOK(w, map[string]interface{}{
		"instance":              instance,
		"available_transitions": availableTransitions,
		"transition_logs":       logs,
	})
}

func (h *WorkflowsHandler) TransitionInstance(w http.ResponseWriter, r *http.Request) {
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

	var req transitionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	instance, err := h.workflows.GetInstanceByUUID(r.Context(), uid)
	if err != nil {
		writeError(w, http.StatusNotFound, "Instance not found")
		return
	}

	def, err := h.workflows.GetDefinitionByID(r.Context(), instance.WorkflowDefinitionID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "Failed to get workflow definition")
		return
	}

	if err := h.workflowSvc.Transition(r.Context(), instance, def, req.Transition, user.ID); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeMutationOK(w)
}

func (h *WorkflowsHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	itemCounts, _ := h.workflows.CountDefinitions(r.Context())
	instanceCount, _ := h.workflows.CountInstances(r.Context())

	writeOK(w, map[string]interface{}{
		"workflow_definitions": itemCounts,
		"workflow_instances":   instanceCount,
	})
}
