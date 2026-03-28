package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const (
	CSRFCookieName = "csrf_token"
	CSRFFieldName  = "csrf_token"
	CSRFHeaderName = "X-CSRF-Token"
)

// CSRF middleware generates and validates CSRF tokens.
func CSRF(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Ensure CSRF cookie exists
		cookie, err := r.Cookie(CSRFCookieName)
		if err != nil || cookie.Value == "" {
			token := generateCSRFToken()
			http.SetCookie(w, &http.Cookie{
				Name:     CSRFCookieName,
				Value:    token,
				Path:     "/",
				HttpOnly: false, // needs to be readable by HTMX
				SameSite: http.SameSiteLaxMode,
			})
			cookie = &http.Cookie{Name: CSRFCookieName, Value: token}
		}

		// Validate on mutating requests
		if r.Method == http.MethodPost || r.Method == http.MethodPut || r.Method == http.MethodDelete || r.Method == http.MethodPatch {
			// Check header first (HTMX), then form field
			token := r.Header.Get(CSRFHeaderName)
			if token == "" {
				token = r.FormValue(CSRFFieldName)
			}

			if token != cookie.Value {
				http.Error(w, "Invalid CSRF token", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func generateCSRFToken() string {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// GetCSRFToken retrieves the current CSRF token from the request.
func GetCSRFToken(r *http.Request) string {
	cookie, err := r.Cookie(CSRFCookieName)
	if err != nil {
		return ""
	}
	return cookie.Value
}
