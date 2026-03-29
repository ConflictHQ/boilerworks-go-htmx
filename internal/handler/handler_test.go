package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// ---------------------------------------------------------------------------
// Mock stores
// ---------------------------------------------------------------------------

type mockItemStore struct {
	items []model.Item
	created  *model.Item
	updated  *model.Item
	deleted  bool
}

func (m *mockItemStore) List(_ context.Context, limit, offset int) ([]model.Item, int, error) {
	total := len(m.items)
	end := offset + limit
	if end > total {
		end = total
	}
	if offset > total {
		return nil, total, nil
	}
	return m.items[offset:end], total, nil
}

func (m *mockItemStore) GetByUUID(_ context.Context, uid uuid.UUID) (*model.Item, error) {
	for i := range m.items {
		if m.items[i].UUID == uid {
			return &m.items[i], nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockItemStore) Create(_ context.Context, name, description string, price float64, status string, categoryID *uuid.UUID, userID uuid.UUID) (*model.Item, error) {
	p := model.Item{
		ID:          uuid.New(),
		UUID:        uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		Status:      status,
		CategoryID:  categoryID,
		CreatedBy:   userID,
		UpdatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.items = append(m.items, p)
	m.created = &p
	return &p, nil
}

func (m *mockItemStore) Update(_ context.Context, uid uuid.UUID, name, description string, price float64, status string, categoryID *uuid.UUID, userID uuid.UUID) (*model.Item, error) {
	for i := range m.items {
		if m.items[i].UUID == uid {
			m.items[i].Name = name
			m.items[i].Description = description
			m.items[i].Price = price
			m.items[i].Status = status
			m.items[i].CategoryID = categoryID
			m.items[i].UpdatedBy = userID
			m.updated = &m.items[i]
			return &m.items[i], nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockItemStore) Delete(_ context.Context, uid uuid.UUID) error {
	for i := range m.items {
		if m.items[i].UUID == uid {
			m.deleted = true
			m.items = append(m.items[:i], m.items[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("not found")
}

type mockCategoryStore struct {
	categories []model.Category
	created    *model.Category
	updated    *model.Category
	deleted    bool
}

func (m *mockCategoryStore) List(_ context.Context, limit, offset int) ([]model.Category, int, error) {
	total := len(m.categories)
	end := offset + limit
	if end > total {
		end = total
	}
	if offset > total {
		return nil, total, nil
	}
	return m.categories[offset:end], total, nil
}

func (m *mockCategoryStore) ListAll(_ context.Context) ([]model.Category, error) {
	return m.categories, nil
}

func (m *mockCategoryStore) GetByUUID(_ context.Context, uid uuid.UUID) (*model.Category, error) {
	for i := range m.categories {
		if m.categories[i].UUID == uid {
			return &m.categories[i], nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockCategoryStore) Create(_ context.Context, name, description string, userID uuid.UUID) (*model.Category, error) {
	c := model.Category{
		ID:          uuid.New(),
		UUID:        uuid.New(),
		Name:        name,
		Description: description,
		CreatedBy:   userID,
		UpdatedBy:   userID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	m.categories = append(m.categories, c)
	m.created = &c
	return &c, nil
}

func (m *mockCategoryStore) Update(_ context.Context, uid uuid.UUID, name, description string, userID uuid.UUID) (*model.Category, error) {
	for i := range m.categories {
		if m.categories[i].UUID == uid {
			m.categories[i].Name = name
			m.categories[i].Description = description
			m.categories[i].UpdatedBy = userID
			m.updated = &m.categories[i]
			return &m.categories[i], nil
		}
	}
	return nil, fmt.Errorf("not found")
}

func (m *mockCategoryStore) Delete(_ context.Context, uid uuid.UUID) error {
	for i := range m.categories {
		if m.categories[i].UUID == uid {
			m.deleted = true
			m.categories = append(m.categories[:i], m.categories[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("not found")
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

func authedContext(ctx context.Context, perms []string) context.Context {
	user := &model.User{
		ID:    uuid.New(),
		Name:  "Test Admin",
		Email: "admin@test.com",
	}
	ctx = context.WithValue(ctx, middleware.UserContextKey, user)
	ctx = context.WithValue(ctx, middleware.PermissionsContextKey, perms)
	return ctx
}

func viewerContext(ctx context.Context) context.Context {
	return authedContext(ctx, []string{"items.view", "categories.view"})
}

func adminContext(ctx context.Context) context.Context {
	return authedContext(ctx, []string{
		"items.view", "items.create", "items.edit", "items.delete",
		"categories.view", "categories.create", "categories.edit", "categories.delete",
	})
}

func seedItem() (uuid.UUID, *mockItemStore) {
	uid := uuid.New()
	store := &mockItemStore{
		items: []model.Item{
			{
				ID:          uuid.New(),
				UUID:        uid,
				Name:        "Widget",
				Description: "A test widget",
				Price:       9.99,
				Status:      "active",
				CreatedBy:   uuid.New(),
				UpdatedBy:   uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}
	return uid, store
}

func seedCategory() (uuid.UUID, *mockCategoryStore) {
	uid := uuid.New()
	store := &mockCategoryStore{
		categories: []model.Category{
			{
				ID:          uuid.New(),
				UUID:        uid,
				Name:        "Electronics",
				Description: "Electronic goods",
				CreatedBy:   uuid.New(),
				UpdatedBy:   uuid.New(),
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
			},
		},
	}
	return uid, store
}

func buildItemRouter(ps ItemStore, cs CategoryStore) *chi.Mux {
	h := NewItemsHandler(ps, cs)
	r := chi.NewRouter()
	r.Route("/items", func(r chi.Router) {
		r.With(middleware.RequirePermission("items.view")).Get("/", h.List)
		r.With(middleware.RequirePermission("items.create")).Post("/", h.Create)
		r.With(middleware.RequirePermission("items.edit")).Get("/{uuid}/edit", h.Edit)
		r.With(middleware.RequirePermission("items.edit")).Post("/{uuid}", h.Update)
		r.With(middleware.RequirePermission("items.delete")).Delete("/{uuid}", h.Delete)
	})
	return r
}

func buildCategoryRouter(cs CategoryStore) *chi.Mux {
	h := NewCategoriesHandler(cs)
	r := chi.NewRouter()
	r.Route("/categories", func(r chi.Router) {
		r.With(middleware.RequirePermission("categories.view")).Get("/", h.List)
		r.With(middleware.RequirePermission("categories.create")).Post("/", h.Create)
		r.With(middleware.RequirePermission("categories.edit")).Get("/{uuid}/edit", h.Edit)
		r.With(middleware.RequirePermission("categories.edit")).Post("/{uuid}", h.Update)
		r.With(middleware.RequirePermission("categories.delete")).Delete("/{uuid}", h.Delete)
	})
	return r
}

// ---------------------------------------------------------------------------
// Item handler tests
// ---------------------------------------------------------------------------

func TestItemList(t *testing.T) {
	_, ps := seedItem()
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
	if ct := w.Header().Get("Content-Type"); !strings.Contains(ct, "text/html") {
		t.Errorf("expected text/html content type, got %s", ct)
	}
}

func TestItemCreate(t *testing.T) {
	ps := &mockItemStore{}
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	form := url.Values{}
	form.Set("name", "New Widget")
	form.Set("description", "A brand new widget")
	form.Set("price", "19.99")
	form.Set("status", "active")

	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303 redirect, got %d: %s", w.Code, w.Body.String())
	}
	if loc := w.Header().Get("Location"); loc != "/items" {
		t.Errorf("expected redirect to /items, got %s", loc)
	}
	if ps.created == nil {
		t.Fatal("expected item to be created")
	}
	if ps.created.Name != "New Widget" {
		t.Errorf("expected name 'New Widget', got '%s'", ps.created.Name)
	}
}

func TestItemGet(t *testing.T) {
	uid, ps := seedItem()
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	req := httptest.NewRequest(http.MethodGet, "/items/"+uid.String()+"/edit", nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestItemUpdate(t *testing.T) {
	uid, ps := seedItem()
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	form := url.Values{}
	form.Set("name", "Updated Widget")
	form.Set("description", "Updated description")
	form.Set("price", "29.99")
	form.Set("status", "inactive")

	req := httptest.NewRequest(http.MethodPost, "/items/"+uid.String(), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", w.Code, w.Body.String())
	}
	if ps.updated == nil {
		t.Fatal("expected item to be updated")
	}
	if ps.updated.Name != "Updated Widget" {
		t.Errorf("expected name 'Updated Widget', got '%s'", ps.updated.Name)
	}
}

func TestItemDelete(t *testing.T) {
	uid, ps := seedItem()
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	req := httptest.NewRequest(http.MethodDelete, "/items/"+uid.String(), nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Delete without HX-Request returns 303 redirect
	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", w.Code, w.Body.String())
	}
	if !ps.deleted {
		t.Fatal("expected item to be deleted from store")
	}
}

// ---------------------------------------------------------------------------
// Category handler tests
// ---------------------------------------------------------------------------

func TestCategoryList(t *testing.T) {
	_, cs := seedCategory()
	router := buildCategoryRouter(cs)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCategoryCreate(t *testing.T) {
	cs := &mockCategoryStore{}
	router := buildCategoryRouter(cs)

	form := url.Values{}
	form.Set("name", "Books")
	form.Set("description", "Book category")

	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", w.Code, w.Body.String())
	}
	if cs.created == nil {
		t.Fatal("expected category to be created")
	}
	if cs.created.Name != "Books" {
		t.Errorf("expected name 'Books', got '%s'", cs.created.Name)
	}
}

func TestCategoryGet(t *testing.T) {
	uid, cs := seedCategory()
	router := buildCategoryRouter(cs)

	req := httptest.NewRequest(http.MethodGet, "/categories/"+uid.String()+"/edit", nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", w.Code, w.Body.String())
	}
}

func TestCategoryUpdate(t *testing.T) {
	uid, cs := seedCategory()
	router := buildCategoryRouter(cs)

	form := url.Values{}
	form.Set("name", "Updated Electronics")
	form.Set("description", "Updated description")

	req := httptest.NewRequest(http.MethodPost, "/categories/"+uid.String(), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", w.Code, w.Body.String())
	}
	if cs.updated == nil {
		t.Fatal("expected category to be updated")
	}
	if cs.updated.Name != "Updated Electronics" {
		t.Errorf("expected name 'Updated Electronics', got '%s'", cs.updated.Name)
	}
}

func TestCategoryDelete(t *testing.T) {
	uid, cs := seedCategory()
	router := buildCategoryRouter(cs)

	req := httptest.NewRequest(http.MethodDelete, "/categories/"+uid.String(), nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusSeeOther {
		t.Fatalf("expected 303, got %d: %s", w.Code, w.Body.String())
	}
	if !cs.deleted {
		t.Fatal("expected category to be deleted from store")
	}
}

// ---------------------------------------------------------------------------
// Permission denial tests (viewer cannot create/update/delete)
// ---------------------------------------------------------------------------

func TestViewerCannotCreateItem(t *testing.T) {
	ps := &mockItemStore{}
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	form := url.Values{}
	form.Set("name", "Forbidden Widget")
	form.Set("price", "5.00")
	form.Set("status", "active")

	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(viewerContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for viewer creating item, got %d", w.Code)
	}
}

func TestViewerCannotUpdateItem(t *testing.T) {
	uid, ps := seedItem()
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	form := url.Values{}
	form.Set("name", "Hacked Widget")
	form.Set("price", "0.01")
	form.Set("status", "active")

	req := httptest.NewRequest(http.MethodPost, "/items/"+uid.String(), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(viewerContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for viewer updating item, got %d", w.Code)
	}
}

func TestViewerCannotDeleteCategory(t *testing.T) {
	uid, cs := seedCategory()
	router := buildCategoryRouter(cs)

	req := httptest.NewRequest(http.MethodDelete, "/categories/"+uid.String(), nil)
	req = req.WithContext(viewerContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for viewer deleting category, got %d", w.Code)
	}
	if cs.deleted {
		t.Fatal("category should not have been deleted")
	}
}

// ---------------------------------------------------------------------------
// HTMX-specific tests
// ---------------------------------------------------------------------------

func TestHTMXRequestReturnsFragment(t *testing.T) {
	_, ps := seedItem()
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	req.Header.Set("HX-Request", "true")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	// HTMX requests should return a fragment -- no <!DOCTYPE or <html> wrapper
	if strings.Contains(body, "<!DOCTYPE") || strings.Contains(body, "<html") {
		t.Error("HTMX response should be a fragment without full page layout")
	}
}

func TestNonHTMXRequestReturnsFullPage(t *testing.T) {
	_, ps := seedItem()
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	req := httptest.NewRequest(http.MethodGet, "/items", nil)
	// No HX-Request header -- normal browser request
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	// Normal requests should include full layout with DOCTYPE
	if !strings.Contains(body, "<!DOCTYPE") && !strings.Contains(body, "<html") {
		t.Error("non-HTMX response should include full page layout")
	}
}

// ---------------------------------------------------------------------------
// HTMX delete returns 200 (not redirect)
// ---------------------------------------------------------------------------

func TestHTMXDeleteReturnsOK(t *testing.T) {
	uid, ps := seedItem()
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	req := httptest.NewRequest(http.MethodDelete, "/items/"+uid.String(), nil)
	req.Header.Set("HX-Request", "true")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for HTMX delete, got %d", w.Code)
	}
	if !ps.deleted {
		t.Fatal("expected item to be deleted")
	}
}

func TestHTMXCategoryDeleteReturnsOK(t *testing.T) {
	uid, cs := seedCategory()
	router := buildCategoryRouter(cs)

	req := httptest.NewRequest(http.MethodDelete, "/categories/"+uid.String(), nil)
	req.Header.Set("HX-Request", "true")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for HTMX category delete, got %d", w.Code)
	}
	if !cs.deleted {
		t.Fatal("expected category to be deleted")
	}
}

// ---------------------------------------------------------------------------
// Validation tests
// ---------------------------------------------------------------------------

func TestItemCreateValidationRejectsEmptyName(t *testing.T) {
	ps := &mockItemStore{}
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	form := url.Values{}
	form.Set("name", "")
	form.Set("price", "10.00")
	form.Set("status", "active")

	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validation failure re-renders form (200), not a redirect
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for validation error, got %d", w.Code)
	}
	if ps.created != nil {
		t.Fatal("item should not have been created with empty name")
	}
}

func TestCategoryCreateValidationRejectsEmptyName(t *testing.T) {
	cs := &mockCategoryStore{}
	router := buildCategoryRouter(cs)

	form := url.Values{}
	form.Set("name", "")
	form.Set("description", "Something")

	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for validation error, got %d", w.Code)
	}
	if cs.created != nil {
		t.Fatal("category should not have been created with empty name")
	}
}

// ---------------------------------------------------------------------------
// Category HTMX tests
// ---------------------------------------------------------------------------

func TestHTMXCategoryListReturnsFragment(t *testing.T) {
	_, cs := seedCategory()
	router := buildCategoryRouter(cs)

	req := httptest.NewRequest(http.MethodGet, "/categories", nil)
	req.Header.Set("HX-Request", "true")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	body := w.Body.String()
	if strings.Contains(body, "<!DOCTYPE") || strings.Contains(body, "<html") {
		t.Error("HTMX category list should return a fragment, not full page")
	}
}

// ---------------------------------------------------------------------------
// Permission denial tests (viewer cannot create/update categories)
// ---------------------------------------------------------------------------

func TestViewerCannotCreateCategory(t *testing.T) {
	cs := &mockCategoryStore{}
	router := buildCategoryRouter(cs)

	form := url.Values{}
	form.Set("name", "Forbidden Category")
	form.Set("description", "No access")

	req := httptest.NewRequest(http.MethodPost, "/categories", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(viewerContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for viewer creating category, got %d", w.Code)
	}
	if cs.created != nil {
		t.Fatal("category should not have been created by viewer")
	}
}

func TestViewerCannotUpdateCategory(t *testing.T) {
	uid, cs := seedCategory()
	router := buildCategoryRouter(cs)

	form := url.Values{}
	form.Set("name", "Hacked Category")
	form.Set("description", "No access")

	req := httptest.NewRequest(http.MethodPost, "/categories/"+uid.String(), strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(viewerContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected 403 for viewer updating category, got %d", w.Code)
	}
	if cs.updated != nil {
		t.Fatal("category should not have been updated by viewer")
	}
}

// ---------------------------------------------------------------------------
// Invalid UUID tests
// ---------------------------------------------------------------------------

func TestItemEditInvalidUUID(t *testing.T) {
	ps := &mockItemStore{}
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	req := httptest.NewRequest(http.MethodGet, "/items/not-a-uuid/edit", nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid UUID, got %d", w.Code)
	}
}

func TestCategoryEditInvalidUUID(t *testing.T) {
	cs := &mockCategoryStore{}
	router := buildCategoryRouter(cs)

	req := httptest.NewRequest(http.MethodGet, "/categories/not-a-uuid/edit", nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected 400 for invalid UUID, got %d", w.Code)
	}
}

func TestItemDeleteNonExistent(t *testing.T) {
	ps := &mockItemStore{}
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	fakeUUID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/items/"+fakeUUID.String(), nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for deleting non-existent item, got %d", w.Code)
	}
}

func TestCategoryDeleteNonExistent(t *testing.T) {
	cs := &mockCategoryStore{}
	router := buildCategoryRouter(cs)

	fakeUUID := uuid.New()
	req := httptest.NewRequest(http.MethodDelete, "/categories/"+fakeUUID.String(), nil)
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected 500 for deleting non-existent category, got %d", w.Code)
	}
}

func TestItemCreateValidationRejectsInvalidPrice(t *testing.T) {
	ps := &mockItemStore{}
	cs := &mockCategoryStore{}
	router := buildItemRouter(ps, cs)

	form := url.Values{}
	form.Set("name", "Widget")
	form.Set("price", "not-a-number")
	form.Set("status", "active")

	req := httptest.NewRequest(http.MethodPost, "/items", strings.NewReader(form.Encode()))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req = req.WithContext(adminContext(req.Context()))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// Validation failure re-renders form (200)
	if w.Code != http.StatusOK {
		t.Fatalf("expected 200 for invalid price, got %d", w.Code)
	}
	if ps.created != nil {
		t.Fatal("item should not have been created with invalid price")
	}
}
