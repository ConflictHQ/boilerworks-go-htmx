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

type ProductsHandler struct {
	products   *queries.ProductQueries
	categories *queries.CategoryQueries
}

func NewProductsHandler(p *queries.ProductQueries, c *queries.CategoryQueries) *ProductsHandler {
	return &ProductsHandler{products: p, categories: c}
}

func (h *ProductsHandler) List(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	perPage := 20

	products, total, err := h.products.List(r.Context(), perPage, (page-1)*perPage)
	if err != nil {
		http.Error(w, "Failed to load products", http.StatusInternalServerError)
		return
	}

	pagination := model.NewPagination(page, perPage, total)

	layout := templates.LayoutData{
		Title:       "Products",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.ProductsListPage(products, pagination, csrf, middleware.HasPermission(r.Context(), "products.create")))
}

func (h *ProductsHandler) New(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	categories, _ := h.categories.ListAll(r.Context())

	layout := templates.LayoutData{
		Title:       "New Product",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.ProductFormPage(nil, categories, csrf, nil))
}

func (h *ProductsHandler) Create(w http.ResponseWriter, r *http.Request) {
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
		layout := templates.LayoutData{Title: "New Product", User: user, Permissions: perms, CSRFToken: csrf}
		renderPage(w, r, layout, pages.ProductFormPage(nil, categories, csrf, errs))
		return
	}

	_, err = h.products.Create(r.Context(), name, description, price, status, categoryID, user.ID)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func (h *ProductsHandler) Edit(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r.Context())
	perms := middleware.GetPermissions(r.Context())
	csrf := middleware.GetCSRFToken(r)

	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	product, err := h.products.GetByUUID(r.Context(), uid)
	if err != nil {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	categories, _ := h.categories.ListAll(r.Context())

	layout := templates.LayoutData{
		Title:       "Edit Product",
		User:        user,
		Permissions: perms,
		CSRFToken:   csrf,
	}

	renderPage(w, r, layout, pages.ProductFormPage(product, categories, csrf, nil))
}

func (h *ProductsHandler) Update(w http.ResponseWriter, r *http.Request) {
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

	_, err = h.products.Update(r.Context(), uid, name, description, price, status, categoryID, user.ID)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}

func (h *ProductsHandler) Delete(w http.ResponseWriter, r *http.Request) {
	uid, err := uuid.Parse(chi.URLParam(r, "uuid"))
	if err != nil {
		http.Error(w, "Invalid UUID", http.StatusBadRequest)
		return
	}

	if err := h.products.Delete(r.Context(), uid); err != nil {
		http.Error(w, "Failed to delete product", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/products", http.StatusSeeOther)
}
