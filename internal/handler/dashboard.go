package handler

import (
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-nextjs/internal/database/queries"
)

type DashboardHandler struct {
	products   *queries.ProductQueries
	categories *queries.CategoryQueries
	forms      *queries.FormQueries
	workflows  *queries.WorkflowQueries
}

func NewDashboardHandler(
	products *queries.ProductQueries,
	categories *queries.CategoryQueries,
	forms *queries.FormQueries,
	workflows *queries.WorkflowQueries,
) *DashboardHandler {
	return &DashboardHandler{
		products:   products,
		categories: categories,
		forms:      forms,
		workflows:  workflows,
	}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	productCounts, _ := h.products.CountByStatus(r.Context())
	categoryCount, _ := h.categories.Count(r.Context())
	formCount, _ := h.forms.CountDefinitions(r.Context())
	submissionCount, _ := h.forms.CountSubmissions(r.Context())
	workflowCount, _ := h.workflows.CountDefinitions(r.Context())
	instanceCount, _ := h.workflows.CountInstances(r.Context())

	writeOK(w, map[string]interface{}{
		"products_by_status": productCounts,
		"category_count":     categoryCount,
		"form_count":         formCount,
		"submission_count":   submissionCount,
		"workflow_count":     workflowCount,
		"instance_count":     instanceCount,
	})
}
