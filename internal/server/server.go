package server

import (
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/handler"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/service"
	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Router *chi.Mux
	pool   *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Server {
	r := chi.NewRouter()

	// Global middleware
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(chimw.RealIP)
	r.Use(chimw.RequestID)

	// Queries
	userQ := queries.NewUserQueries(pool)
	sessionQ := queries.NewSessionQueries(pool)
	categoryQ := queries.NewCategoryQueries(pool)
	itemQ := queries.NewItemQueries(pool)
	formQ := queries.NewFormQueries(pool)
	workflowQ := queries.NewWorkflowQueries(pool)

	// Services
	authSvc := service.NewAuthService(userQ, sessionQ)
	formSvc := service.NewFormService()
	workflowSvc := service.NewWorkflowService(workflowQ)

	// Handlers
	healthH := handler.NewHealthHandler()
	authH := handler.NewAuthHandler(authSvc)
	dashboardH := handler.NewDashboardHandler(itemQ, categoryQ, formQ, workflowQ)
	itemsH := handler.NewItemsHandler(itemQ, categoryQ)
	categoriesH := handler.NewCategoriesHandler(categoryQ)
	formsH := handler.NewFormsHandler(formQ, formSvc)
	workflowsH := handler.NewWorkflowsHandler(workflowQ, workflowSvc)

	s := &Server{Router: r, pool: pool}
	s.registerRoutes(r, authSvc, healthH, authH, dashboardH, itemsH, categoriesH, formsH, workflowsH)

	return s
}
