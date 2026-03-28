package handler

import (
	"net/http"
	"strconv"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type CategoriesHandler struct {
	categories *queries.CategoryQueries
}

func NewCategoriesHandler(c *queries.CategoryQueries) *CategoriesHandler {
	return &CategoriesHandler{categories: c}
}

func (h *CategoriesHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20

	categories, total, err := h.categories.List(r.Context(), perPage, (page-1)*perPage)
	if err != nil {
		http.Error(w, "Failed to load categories", http.StatusInternalServerError)
		return
	}

	pagination := model.NewPagination(page, perPage, total)

	layout := templates.LayoutData{
		Title:       "Categories",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.CategoriesListPage(categories, pagination, csrf, middleware.HasPermission(r.Context(), "categories.create")))
}

func (h *CategoriesHandler) New(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	layout := templates.LayoutData{
		Title:       "New Category",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.CategoryFormPage(nil, csrf, nil))
}

func (h *CategoriesHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")

	if name == "" {
		perms := middleware.GetPermissions(r.Context())
		csrf := middleware.GetCSRFToken(r)
		layout := templates.LayoutData{Title: "New Category", User: user, Permissions: perms, CSRFToken: csrf}
		renderPage(w, r, layout, pages.CategoryFormPage(nil, csrf, []string{"Name is required"}))
		return
	}

	_, err := h.categories.Create(r.Context(), name, description, user.ID)
	if err != nil {
		http.Error(w, "Failed to create category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (h *CategoriesHandler) Edit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	category, err := h.categories.GetByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	layout := templates.LayoutData{
		Title:       "Edit Category",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.CategoryFormPage(category, csrf, nil))
}

func (h *CategoriesHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.categories.Update(r.Context(), uid, name, description, user.ID)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}

func (h *CategoriesHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	if err := h.categories.Delete(r.Context(), uid); err != nil {
		http.Error(w, "Failed to delete category", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/categories", http.StatusSeeOther)
}
