package handler

import (
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates/pages"
)

type DashboardHandler struct {
	products   *queries.ProductQueries
	categories *queries.CategoryQueries
	forms      *queries.FormQueries
	workflows  *queries.WorkflowQueries
}

func NewDashboardHandler(p *queries.ProductQueries, c *queries.CategoryQueries, f *queries.FormQueries, w *queries.WorkflowQueries) *DashboardHandler {
	return &DashboardHandler{products: p, categories: c, forms: f, workflows: w}
}

func (h *DashboardHandler) Dashboard(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	_, totalProducts, _ := h.products.List(r.Context(), 1, 0)
	_, totalCategories, _ := h.categories.List(r.Context(), 1, 0)
	totalForms, _ := h.forms.CountDefinitions(r.Context())
	totalSubmissions, _ := h.forms.CountSubmissions(r.Context())
	totalWorkflows, _ := h.workflows.CountDefinitions(r.Context())
	totalInstances, _ := h.workflows.CountInstances(r.Context())
	productsByStatus, _ := h.products.CountByStatus(r.Context())

	data := pages.DashboardData{
		TotalProducts:    totalProducts,
		TotalCategories:  totalCategories,
		TotalForms:       totalForms,
		TotalSubmissions: totalSubmissions,
		TotalWorkflows:   totalWorkflows,
		TotalInstances:   totalInstances,
		ProductsByStatus: productsByStatus,
	}

	layout := templates.LayoutData{
		Title:       "Dashboard",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.DashboardPage(data))
}
