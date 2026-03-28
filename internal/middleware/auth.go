package middleware

import (
	"context"
	"net/http"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/service"
)

type contextKey string

const (
	UserContextKey        contextKey = "user"
	PermissionsContextKey contextKey = "permissions"
)

// RequireAuth enforces that a valid session exists.
func RequireAuth(authSvc *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("session_token")
			if err != nil || cookie.Value == "" {
				if isHTMX(r) {
					w.Header().Set("HX-Redirect", "/login")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			user, err := authSvc.ValidateSession(r.Context(), cookie.Value)
			if err != nil {
				clearSessionCookie(w)
				if isHTMX(r) {
					w.Header().Set("HX-Redirect", "/login")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				http.Redirect(w, r, "/login", http.StatusSeeOther)
				return
			}

			perms, err := authSvc.GetUserPermissions(r.Context(), user.ID)
			if err != nil {
				perms = []string{}
			}

			ctx := context.WithValue(r.Context(), UserContextKey, user)
			ctx = context.WithValue(ctx, PermissionsContextKey, perms)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// RequirePermission checks that the authenticated user has the given permission.
func RequirePermission(perm string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			perms := GetPermissions(r.Context())
			for _, p := range perms {
				if p == perm {
					next.ServeHTTP(w, r)
					return
				}
			}
			http.Error(w, "Forbidden", http.StatusForbidden)
		})
	}
}

// GetUser extracts the user from context.
func GetUser(ctx context.Context) *model.User {
	user, ok := ctx.Value(UserContextKey).(*model.User)
	if !ok {
		return nil
	}
	return user
}

// GetPermissions extracts permissions from context.
func GetPermissions(ctx context.Context) []string {
	perms, ok := ctx.Value(PermissionsContextKey).([]string)
	if !ok {
		return nil
	}
	return perms
}

// HasPermission checks if the context user has a specific permission.
func HasPermission(ctx context.Context, perm string) bool {
	perms := GetPermissions(ctx)
	for _, p := range perms {
		if p == perm {
			return true
		}
	}
	return false
}

func isHTMX(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
}
