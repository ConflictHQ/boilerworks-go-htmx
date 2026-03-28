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

type FormsHandler struct {
	forms   *queries.FormQueries
	formSvc *service.FormService
}

func NewFormsHandler(f *queries.FormQueries, svc *service.FormService) *FormsHandler {
	return &FormsHandler{forms: f, formSvc: svc}
}

func (h *FormsHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20

	forms, total, err := h.forms.ListDefinitions(r.Context(), perPage, (page-1)*perPage)
	if err != nil {
		http.Error(w, "Failed to load forms", http.StatusInternalServerError)
		return
	}

	pagination := model.NewPagination(page, perPage, total)

	layout := templates.LayoutData{
		Title:       "Forms",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.FormsListPage(forms, pagination, csrf, middleware.HasPermission(r.Context(), "forms.create")))
}

func (h *FormsHandler) New(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	layout := templates.LayoutData{
		Title:       "New Form",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.FormDefinitionPage(nil, csrf, nil))
}

func (h *FormsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	slug := r.FormValue("slug")
	description := r.FormValue("description")
	status := r.FormValue("status")
	schemaStr := r.FormValue("schema")

	var errs []string
	if name == "" {
		errs = append(errs, "Name is required")
	}
	if slug == "" {
		errs = append(errs, "Slug is required")
	}

	var schema json.RawMessage
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		errs = append(errs, "Schema must be valid JSON")
	}

	if len(errs) > 0 {
		perms := middleware.GetPermissions(r.Context())
		csrf := middleware.GetCSRFToken(r)
		layout := templates.LayoutData{Title: "New Form", User: user, Permissions: perms, CSRFToken: csrf}
		renderPage(w, r, layout, pages.FormDefinitionPage(nil, csrf, errs))
		return
	}

	_, err := h.forms.CreateDefinition(r.Context(), name, slug, description, status, schema, user.ID)
	if err != nil {
		http.Error(w, "Failed to create form", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/forms", http.StatusSeeOther)
}

func (h *FormsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	form, err := h.forms.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	layout := templates.LayoutData{
		Title:       "Edit Form",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.FormDefinitionPage(form, csrf, nil))
}

func (h *FormsHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	slug := r.FormValue("slug")
	description := r.FormValue("description")
	status := r.FormValue("status")
	schemaStr := r.FormValue("schema")

	var schema json.RawMessage
	if err := json.Unmarshal([]byte(schemaStr), &schema); err != nil {
		http.Error(w, "Invalid schema JSON", http.StatusBadRequest)
		return
	}

	_, err = h.forms.UpdateDefinition(r.Context(), uid, name, slug, description, status, schema, user.ID)
	if err != nil {
		http.Error(w, "Failed to update form", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/forms", http.StatusSeeOther)
}

func (h *FormsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.forms.DeleteDefinition(r.Context(), uid); err != nil {
		http.Error(w, "Failed to delete form", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/forms", http.StatusSeeOther)
}

func (h *FormsHandler) ShowSubmitForm(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	form, err := h.forms.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	layout := templates.LayoutData{
		Title:       form.Name,
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.FormSubmitPage(form, csrf, nil))
}

func (h *FormsHandler) SubmitForm(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	form, err := h.forms.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	// Collect form data
	data := make(map[string]string)
	for _, field := range form.Schema {
		data[field.Name] = r.FormValue(field.Name)
	}

	jsonData, validationErrs := h.formSvc.ValidateSubmission(form, data)
	if len(validationErrs) > 0 {
		layout := templates.LayoutData{Title: form.Name, User: user, Permissions: perms, CSRFToken: csrf}
		renderPage(w, r, layout, pages.FormSubmitPage(form, csrf, validationErrs))
		return
	}

	_, err = h.forms.CreateSubmission(r.Context(), form.ID, jsonData, user.ID)
	if err != nil {
		http.Error(w, "Failed to submit form", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/forms", http.StatusSeeOther)
}

func (h *FormsHandler) ListSubmissions(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	form, err := h.forms.GetDefinitionByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Form not found", http.StatusNotFound)
		return
	}

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20

	submissions, total, err := h.forms.ListSubmissions(r.Context(), form.ID, perPage, (page-1)*perPage)
	if err != nil {
		http.Error(w, "Failed to load submissions", http.StatusInternalServerError)
		return
	}

	pagination := model.NewPagination(page, perPage, total)

	layout := templates.LayoutData{
		Title:       form.Name + " Submissions",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.FormSubmissionsPage(form, submissions, pagination))
}
