package handler

import (
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/middleware"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/service"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates/pages"
)

type AuthHandler struct {
	authSvc *service.AuthService
}

func NewAuthHandler(authSvc *service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

func (h *AuthHandler) LoginPage(w http.ResponseWriter, r *http.Request) {
	csrf := middleware.GetCSRFToken(r)
	layout := templates.LayoutData{Title: "Login", CSRFToken: csrf}
	_ = templates.Layout(layout).Render(r.Context(), w)
	_ = pages.LoginPage(csrf, "").Render(r.Context(), w)
}

func (h *AuthHandler) LoginPageFull(w http.ResponseWriter, r *http.Request) {
	csrf := middleware.GetCSRFToken(r)
	component := pages.LoginPage(csrf, "")
	layout := templates.LayoutData{Title: "Login", CSRFToken: csrf}
	_ = templates.Layout(layout).Render(r.Context(), w)
	_ = component
}

func (h *AuthHandler) ShowLogin(w http.ResponseWriter, r *http.Request) {
	csrf := middleware.GetCSRFToken(r)
	renderPage(w, r, templates.LayoutData{Title: "Login", CSRFToken: csrf}, pages.LoginPage(csrf, ""))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	email := r.FormValue("email")
	password := r.FormValue("password")
	csrf := middleware.GetCSRFToken(r)

	_, token, err := h.authSvc.Login(r.Context(), email, password)
	if err != nil {
		renderPage(w, r, templates.LayoutData{Title: "Login", CSRFToken: csrf}, pages.LoginPage(csrf, "Invalid email or password"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) ShowRegister(w http.ResponseWriter, r *http.Request) {
	csrf := middleware.GetCSRFToken(r)
	renderPage(w, r, templates.LayoutData{Title: "Register", CSRFToken: csrf}, pages.RegisterPage(csrf, ""))
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20)
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	email := r.FormValue("email")
	password := r.FormValue("password")
	csrf := middleware.GetCSRFToken(r)

	if len(password) < 8 {
		renderPage(w, r, templates.LayoutData{Title: "Register", CSRFToken: csrf}, pages.RegisterPage(csrf, "Password must be at least 8 characters"))
		return
	}

	_, token, err := h.authSvc.Register(r.Context(), name, email, password)
	if err != nil {
		renderPage(w, r, templates.LayoutData{Title: "Register", CSRFToken: csrf}, pages.RegisterPage(csrf, "Email already in use"))
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    token,
		Path:     "/",
		MaxAge:   30 * 24 * 60 * 60,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_token")
	if err == nil {
		_ = h.authSvc.Logout(r.Context(), cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
