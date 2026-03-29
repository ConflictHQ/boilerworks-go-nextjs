package handler

import (
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/database/queries"
)

type DashboardHandler struct {
	items   *queries.ItemQueries
	categories *queries.CategoryQueries
	forms      *queries.FormQueries
	workflows  *queries.WorkflowQueries
}

func NewDashboardHandler(
	items *queries.ItemQueries,
	categories *queries.CategoryQueries,
	forms *queries.FormQueries,
	workflows *queries.WorkflowQueries,
) *DashboardHandler {
	return &DashboardHandler{
		items:   items,
		categories: categories,
		forms:      forms,
		workflows:  workflows,
	}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	itemCounts, _ := h.items.CountByStatus(r.Context())
	categoryCount, _ := h.categories.Count(r.Context())
	formCount, _ := h.forms.CountDefinitions(r.Context())
	submissionCount, _ := h.forms.CountSubmissions(r.Context())
	workflowCount, _ := h.workflows.CountDefinitions(r.Context())
	instanceCount, _ := h.workflows.CountInstances(r.Context())

	writeOK(w, map[string]interface{}{
		"items_by_status": itemCounts,
		"category_count":     categoryCount,
		"form_count":         formCount,
		"submission_count":   submissionCount,
		"workflow_count":     workflowCount,
		"instance_count":     instanceCount,
	})
}
