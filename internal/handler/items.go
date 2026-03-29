package handler

import (
	"net/http"
	"strconv"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates/pages"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ItemsHandler struct {
	items   ItemStore
	categories CategoryStore
}

func NewItemsHandler(p ItemStore, c CategoryStore) *ItemsHandler {
	return &ItemsHandler{items: p, categories: c}
}

func (h *ItemsHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20

	items, total, err := h.items.List(r.Context(), perPage, (page-1)*perPage)
	if err != nil {
		http.Error(w, "Failed to load items", http.StatusInternalServerError)
		return
	}

	pagination := model.NewPagination(page, perPage, total)

	layout := templates.LayoutData{
		Title:       "Items",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.ItemsListPage(items, pagination, csrf, middleware.HasPermission(r.Context(), "items.create")))
}

func (h *ItemsHandler) New(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	categories, _ := h.categories.ListAll(r.Context())

	layout := templates.LayoutData{
		Title:       "New Item",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.ItemFormPage(nil, categories, csrf, nil))
}

func (h *ItemsHandler) Create(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	description := r.FormValue("description")
	priceStr := r.FormValue("price")
	status := r.FormValue("status")
	categoryIDStr := r.FormValue("category_id")

	var errs []string
	if name == "" {
		errs = append(errs, "Name is required")
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		errs = append(errs, "Price must be a valid number")
	}

	var categoryID *uuid.UUID
	if categoryIDStr != "" {
		cid, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = &cid
		}
	}

	if len(errs) > 0 {
		perms := middleware.GetPermissions(r.Context())
		csrf := middleware.GetCSRFToken(r)
		categories, _ := h.categories.ListAll(r.Context())
		layout := templates.LayoutData{Title: "New Item", User: user, Permissions: perms, CSRFToken: csrf}
		renderPage(w, r, layout, pages.ItemFormPage(nil, categories, csrf, errs))
		return
	}

	_, err = h.items.Create(r.Context(), name, description, price, status, categoryID, user.ID)
	if err != nil {
		http.Error(w, "Failed to create item", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func (h *ItemsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	item, err := h.items.GetByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Item not found", http.StatusNotFound)
		return
	}

	categories, _ := h.categories.ListAll(r.Context())

	layout := templates.LayoutData{
		Title:       "Edit Item",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.ItemFormPage(item, categories, csrf, nil))
}

func (h *ItemsHandler) Update(w http.ResponseWriter, r *http.Request) {
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
	priceStr := r.FormValue("price")
	status := r.FormValue("status")
	categoryIDStr := r.FormValue("category_id")

	price, _ := strconv.ParseFloat(priceStr, 64)

	var categoryID *uuid.UUID
	if categoryIDStr != "" {
		cid, err := uuid.Parse(categoryIDStr)
		if err == nil {
			categoryID = &cid
		}
	}

	_, err = h.items.Update(r.Context(), uid, name, description, price, status, categoryID, user.ID)
	if err != nil {
		http.Error(w, "Failed to update item", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}

func (h *ItemsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.items.Delete(r.Context(), uid); err != nil {
		http.Error(w, "Failed to delete item", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/items", http.StatusSeeOther)
}
