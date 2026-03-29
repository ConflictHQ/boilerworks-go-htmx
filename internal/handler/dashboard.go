package handler

import (
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates/pages"
)

type DashboardHandler struct {
	items   *queries.ItemQueries
	categories *queries.CategoryQueries
	forms      *queries.FormQueries
	workflows  *queries.WorkflowQueries
}

func NewDashboardHandler(p *queries.ItemQueries, c *queries.CategoryQueries, f *queries.FormQueries, w *queries.WorkflowQueries) *DashboardHandler {
	return &DashboardHandler{items: p, categories: c, forms: f, workflows: w}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	_, totalItems, _ := h.items.List(r.Context(), 1, 0)
	_, totalCategories, _ := h.categories.List(r.Context(), 1, 0)
	totalForms, _ := h.forms.CountDefinitions(r.Context())
	totalSubmissions, _ := h.forms.CountSubmissions(r.Context())
	totalWorkflows, _ := h.workflows.CountDefinitions(r.Context())
	totalInstances, _ := h.workflows.CountInstances(r.Context())
	itemsByStatus, _ := h.items.CountByStatus(r.Context())

	data := pages.DashboardData{
		TotalItems:    totalItems,
		TotalCategories:  totalCategories,
		TotalForms:       totalForms,
		TotalSubmissions: totalSubmissions,
		TotalWorkflows:   totalWorkflows,
		TotalInstances:   totalInstances,
		ItemsByStatus: itemsByStatus,
	}

	layout := templates.LayoutData{
		Title:       "Dashboard",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.DashboardPage(data))
}
