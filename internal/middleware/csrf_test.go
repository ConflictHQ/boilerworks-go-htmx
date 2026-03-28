package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCSRFSetsCookie(t *testing.T) {
	handler := CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	cookies := w.Result().Cookies()
	found := false
	for _, c := range cookies {
		if c.Name == CSRFCookieName {
			found = true
			if c.Value == "" {
				t.Error("CSRF cookie value should not be empty")
			}
		}
	}
	if !found {
		t.Error("expected CSRF cookie to be set")
	}
}

func TestCSRFBlocksPostWithoutToken(t *testing.T) {
	handler := CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(""))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// Set a CSRF cookie but don't include the token in form/header
	req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "valid-token"})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusForbidden {
		t.Errorf("expected 403 without CSRF token, got %d", w.Code)
	}
}

func TestCSRFAllowsPostWithMatchingToken(t *testing.T) {
	handler := CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := "test-csrf-token"
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader("csrf_token="+token))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: token})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with matching CSRF token, got %d", w.Code)
	}
}

func TestCSRFAllowsPostWithHeader(t *testing.T) {
	handler := CSRF(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	token := "test-csrf-token"
	req := httptest.NewRequest(http.MethodPost, "/", nil)
	req.Header.Set(CSRFHeaderName, token)
	req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: token})
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("expected 200 with CSRF header, got %d", w.Code)
	}
}

func TestGetCSRFToken(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(&http.Cookie{Name: CSRFCookieName, Value: "my-token"})

	token := GetCSRFToken(req)
	if token != "my-token" {
		t.Errorf("expected 'my-token', got '%s'", token)
	}
}

func TestGetCSRFTokenMissing(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	token := GetCSRFToken(req)
	if token != "" {
		t.Errorf("expected empty string, got '%s'", token)
	}
}
