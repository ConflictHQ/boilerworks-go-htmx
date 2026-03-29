package server

import (
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/handler"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/service"
	"github.com/go-chi/chi/v5"
)

func (s *Server) registerRoutes(
	r *chi.Mux,
	authSvc *service.AuthService,
	healthH *handler.HealthHandler,
	authH *handler.AuthHandler,
	dashboardH *handler.DashboardHandler,
	itemsH *handler.ItemsHandler,
	categoriesH *handler.CategoriesHandler,
	formsH *handler.FormsHandler,
	workflowsH *handler.WorkflowsHandler,
) {
	// Health check (no auth, no CSRF)
	r.Get("/health", healthH.Health)

	// Public routes with CSRF
	r.Group(func(r chi.Router) {
		r.Use(middleware.CSRF)

		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
		})
		r.Get("/login", authH.ShowLogin)
		r.Post("/login", authH.Login)
		r.Get("/register", authH.ShowRegister)
		r.Post("/register", authH.Register)
	})

	// Authenticated routes
	r.Group(func(r chi.Router) {
		r.Use(middleware.CSRF)
		r.Use(middleware.RequireAuth(authSvc))

		r.Post("/logout", authH.Logout)
		r.Get("/dashboard", dashboardH.Dashboard)

		// Items
		r.Route("/items", func(r chi.Router) {
			r.With(middleware.RequirePermission("items.view")).Get("/", itemsH.List)
			r.With(middleware.RequirePermission("items.create")).Get("/new", itemsH.New)
			r.With(middleware.RequirePermission("items.create")).Post("/", itemsH.Create)
			r.With(middleware.RequirePermission("items.edit")).Get("/{uuid}/edit", itemsH.Edit)
			r.With(middleware.RequirePermission("items.edit")).Post("/{uuid}", itemsH.Update)
			r.With(middleware.RequirePermission("items.delete")).Delete("/{uuid}", itemsH.Delete)
		})

		// Categories
		r.Route("/categories", func(r chi.Router) {
			r.With(middleware.RequirePermission("categories.view")).Get("/", categoriesH.List)
			r.With(middleware.RequirePermission("categories.create")).Get("/new", categoriesH.New)
			r.With(middleware.RequirePermission("categories.create")).Post("/", categoriesH.Create)
			r.With(middleware.RequirePermission("categories.edit")).Get("/{uuid}/edit", categoriesH.Edit)
			r.With(middleware.RequirePermission("categories.edit")).Post("/{uuid}", categoriesH.Update)
			r.With(middleware.RequirePermission("categories.delete")).Delete("/{uuid}", categoriesH.Delete)
		})

		// Forms
		r.Route("/forms", func(r chi.Router) {
			r.With(middleware.RequirePermission("forms.view")).Get("/", formsH.List)
			r.With(middleware.RequirePermission("forms.create")).Get("/new", formsH.New)
			r.With(middleware.RequirePermission("forms.create")).Post("/", formsH.Create)
			r.With(middleware.RequirePermission("forms.edit")).Get("/{uuid}/edit", formsH.Edit)
			r.With(middleware.RequirePermission("forms.edit")).Post("/{uuid}", formsH.Update)
			r.With(middleware.RequirePermission("forms.delete")).Delete("/{uuid}", formsH.Delete)
			r.With(middleware.RequirePermission("forms.view")).Get("/{uuid}/submit", formsH.ShowSubmitForm)
			r.With(middleware.RequirePermission("forms.create")).Post("/{uuid}/submit", formsH.SubmitForm)
			r.With(middleware.RequirePermission("forms.view")).Get("/{uuid}/submissions", formsH.ListSubmissions)
		})

		// Workflows
		r.Route("/workflows", func(r chi.Router) {
			r.With(middleware.RequirePermission("workflows.view")).Get("/", workflowsH.List)
			r.With(middleware.RequirePermission("workflows.create")).Get("/new", workflowsH.New)
			r.With(middleware.RequirePermission("workflows.create")).Post("/", workflowsH.Create)
			r.With(middleware.RequirePermission("workflows.edit")).Get("/{uuid}/edit", workflowsH.Edit)
			r.With(middleware.RequirePermission("workflows.edit")).Post("/{uuid}", workflowsH.Update)
			r.With(middleware.RequirePermission("workflows.delete")).Delete("/{uuid}", workflowsH.Delete)
			r.With(middleware.RequirePermission("workflows.view")).Get("/{uuid}/instances", workflowsH.ListInstances)
			r.With(middleware.RequirePermission("workflows.create")).Post("/{uuid}/instances", workflowsH.CreateInstance)
			r.With(middleware.RequirePermission("workflows.view")).Get("/instances/{uuid}", workflowsH.ShowInstance)
			r.With(middleware.RequirePermission("workflows.edit")).Post("/instances/{uuid}/transition", workflowsH.TransitionInstance)
		})
	})
}
