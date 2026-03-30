package handler

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/ConflictHQ/boilerworks-go-htmx/templates"
)

// renderPage renders a full page with layout, or just the content for HTMX requests.
func renderPage(w http.ResponseWriter, r *http.Request, layout templates.LayoutData, content templ.Component) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	if r.Header.Get("HX-Request") == "true" {
		_ = content.Render(r.Context(), w)
		return
	}

	wrapped := templates.LayoutWithContent(layout, content)
	_ = wrapped.Render(r.Context(), w)
}
