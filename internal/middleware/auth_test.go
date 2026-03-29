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
	perms := []string{"items.view", "items.create"}

	ctx := context.WithValue(context.Background(), PermissionsContextKey, perms)

	got := GetPermissions(ctx)
	if len(got) != 2 {
		t.Fatalf("expected 2 permissions, got %d", len(got))
	}
}

func TestHasPermission(t *testing.T) {
	perms := []string{"items.view", "items.create", "categories.view"}
	ctx := context.WithValue(context.Background(), PermissionsContextKey, perms)

	if !HasPermission(ctx, "items.view") {
		t.Error("expected HasPermission to return true for items.view")
	}

	if HasPermission(ctx, "items.delete") {
		t.Error("expected HasPermission to return false for items.delete")
	}
}

func TestHasPermissionEmptyContext(t *testing.T) {
	if HasPermission(context.Background(), "anything") {
		t.Error("expected false for empty context")
	}
}
