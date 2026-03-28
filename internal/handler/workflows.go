package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/service"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type WorkflowsHandler struct {
	workflows   *queries.WorkflowQueries
	workflowSvc *service.WorkflowService
}

func NewWorkflowsHandler(w *queries.WorkflowQueries, svc *service.WorkflowService) *WorkflowsHandler {
	return &WorkflowsHandler{workflows: w, workflowSvc: svc}
}

func (h *WorkflowsHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20

	workflows, total, err := h.workflows.ListDefinitions(r.Context(), perPage, (page-1)*perPage)
	if err != nil {
		http.Error(w, "Failed to load workflows", http.StatusInternalServerError)
		return
	}

	pagination := model.NewPagination(page, perPage, total)

	layout := templates.LayoutData{
		Title:       "Workflows",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.WorkflowsListPage(workflows, pagination, csrf, middleware.HasPermission(r.Context(), "workflows.create")))
}

func (h *WorkflowsHandler) New(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	layout := templates.LayoutData{
		Title:       "New Workflow",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.WorkflowDefinitionPage(nil, csrf, nil))
}

func (h *WorkflowsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	status := r.FormValue("status")
	statesStr := r.FormValue("states")
	transitionsStr := r.FormValue("transitions")

	var errs []string
	if name == "" {
		errs = append(errs, "Name is required")
	}

	var states json.RawMessage
	if err := json.Unmarshal([]byte(statesStr), &states); err != nil {
		errs = append(errs, "States must be valid JSON")
	}

	var transitions json.RawMessage
	if err := json.Unmarshal([]byte(transitionsStr), &transitions); err != nil {
		errs = append(errs, "Transitions must be valid JSON")
	}

	if len(errs) > 0 {
		perms := middleware.GetPermissions(r.Context())
		csrf := middleware.GetCSRFToken(r)
		layout := templates.LayoutData{Title: "New Workflow", User: user, Permissions: perms, CSRFToken: csrf}
		renderPage(w, r, layout, pages.WorkflowDefinitionPage(nil, csrf, errs))
		return
	}

	_, err := h.workflows.CreateDefinition(r.Context(), name, description, status, states, transitions, user.ID)
	if err != nil {
		http.Error(w, "Failed to create workflow", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/workflows", http.StatusSeeOther)
}

func (h *WorkflowsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	wf, err := h.workflows.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	layout := templates.LayoutData{
		Title:       "Edit Workflow",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.WorkflowDefinitionPage(wf, csrf, nil))
}

func (h *WorkflowsHandler) Update(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	status := r.FormValue("status")
	statesStr := r.FormValue("states")
	transitionsStr := r.FormValue("transitions")

	var states json.RawMessage
	if err := json.Unmarshal([]byte(statesStr), &states); err != nil {
		http.Error(w, "Invalid states JSON", http.StatusBadRequest)
		return
	}

	var transitions json.RawMessage
	if err := json.Unmarshal([]byte(transitionsStr), &transitions); err != nil {
		http.Error(w, "Invalid transitions JSON", http.StatusBadRequest)
		return
	}

	_, err = h.workflows.UpdateDefinition(r.Context(), uid, name, description, status, states, transitions, user.ID)
	if err != nil {
		http.Error(w, "Failed to update workflow", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/workflows", http.StatusSeeOther)
}

func (h *WorkflowsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.workflows.DeleteDefinition(r.Context(), uid); err != nil {
		http.Error(w, "Failed to delete workflow", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/workflows", http.StatusSeeOther)
}

func (h *WorkflowsHandler) ListInstances(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	wf, err := h.workflows.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20

	instances, total, err := h.workflows.ListInstances(r.Context(), wf.ID, perPage, (page-1)*perPage)
	if err != nil {
		http.Error(w, "Failed to load instances", http.StatusInternalServerError)
		return
	}

	pagination := model.NewPagination(page, perPage, total)

	layout := templates.LayoutData{
		Title:       wf.Name + " Instances",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.WorkflowInstancesPage(wf, instances, pagination, csrf))
}

func (h *WorkflowsHandler) CreateInstance(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	wf, err := h.workflows.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Workflow not found", http.StatusNotFound)
		return
	}

	initialState, err := h.workflowSvc.GetInitialState(wf)
	if err != nil {
		http.Error(w, "Workflow has no states", http.StatusBadRequest)
		return
	}

	_, err = h.workflows.CreateInstance(r.Context(), wf.ID, initialState, user.ID)
	if err != nil {
		http.Error(w, "Failed to create instance", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/workflows/"+uid.String()+"/instances", http.StatusSeeOther)
}

func (h *WorkflowsHandler) ShowInstance(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	instance, err := h.workflows.GetInstanceByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Instance not found", http.StatusNotFound)
		return
	}

	def, err := h.workflows.GetDefinitionByUUID(r.Context(), instance.WorkflowDefinitionID)
	if err != nil {
		// Try by ID instead
		http.Error(w, "Workflow definition not found", http.StatusNotFound)
		return
	}

	transitions := h.workflowSvc.GetAvailableTransitions(def, instance.CurrentState)
	logs, _ := h.workflows.GetTransitionLogs(r.Context(), instance.ID)

	layout := templates.LayoutData{
		Title:       "Workflow Instance",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.WorkflowInstancePage(instance, def, transitions, logs, csrf))
}

func (h *WorkflowsHandler) TransitionInstance(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	transitionName := r.FormValue("transition")

	instance, err := h.workflows.GetInstanceByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Instance not found", http.StatusNotFound)
		return
	}

	def, err := h.workflows.GetDefinitionByUUID(r.Context(), instance.WorkflowDefinitionID)
	if err != nil {
		http.Error(w, "Workflow definition not found", http.StatusNotFound)
		return
	}

	if err := h.workflowSvc.Transition(r.Context(), instance, def, transitionName, user.ID); err != nil {
		http.Error(w, "Transition failed: "+err.Error(), http.StatusBadRequest)
		return
	}

	http.Redirect(w, r, "/workflows/instances/"+uid.String(), http.StatusSeeOther)
}
