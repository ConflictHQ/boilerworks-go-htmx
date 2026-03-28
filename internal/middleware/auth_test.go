package middleware

import (
	"context"
	"testing"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
)

func TestGetUserFromContext(t *testing.T) {
	user := &model.User{
		ID:    uuid.New(),
		Name:  "Test User",
		Email: "test@example.com",
	}

	ctx := context.WithValue(context.Background(), UserContextKey, user)

	got := GetUser(ctx)
	if got == nil {
		t.Fatal("expected user, got nil")
	}
	if got.ID != user.ID {
		t.Errorf("expected user ID %s, got %s", user.ID, got.ID)
	}
}

func TestGetUserFromEmptyContext(t *testing.T) {
	got := GetUser(context.Background())
	if got != nil {
		t.Errorf("expected nil user from empty context, got %v", got)
	}
}

func TestGetPermissionsFromContext(t *testing.T) {
	perms := []string{"products.view", "products.create"}

	ctx := context.WithValue(context.Background(), PermissionsContextKey, perms)

	got := GetPermissions(ctx)
	if len(got) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(got))
	}
}

func TestHasPermission(t *testing.T) {
	perms := []string{"products.view", "products.create", "categories.view"}
	ctx := context.WithValue(context.Background(), PermissionsContextKey, perms)

	if !HasPermission(ctx, "products.view") {
		t.Error("expected HasPermission to return true for products.view")
	}

	if HasPermission(ctx, "products.delete") {
		t.Error("expected HasPermission to return false for products.delete")
	}
}

func TestHasPermissionEmptyContext(t *testing.T) {
	if HasPermission(context.Background(), "anything") {
		t.Error("expected false for empty context")
	}
}
