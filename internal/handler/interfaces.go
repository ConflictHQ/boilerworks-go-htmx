package handler

import (
	"context"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
)

// ItemStore defines the data access methods for items.
type ItemStore interface {
	List(ctx context.Context, limit, offset int) ([]model.Item, int, error)
	GetByUUID(ctx context.Context, uid uuid.UUID) (*model.Item, error)
	Create(ctx context.Context, name, description string, price float64, status string, categoryID *uuid.UUID, userID uuid.UUID) (*model.Item, error)
	Update(ctx context.Context, uid uuid.UUID, name, description string, price float64, status string, categoryID *uuid.UUID, userID uuid.UUID) (*model.Item, error)
	Delete(ctx context.Context, uid uuid.UUID) error
}

// CategoryStore defines the data access methods for categories.
type CategoryStore interface {
	List(ctx context.Context, limit, offset int) ([]model.Category, int, error)
	ListAll(ctx context.Context) ([]model.Category, error)
	GetByUUID(ctx context.Context, uid uuid.UUID) (*model.Category, error)
	Create(ctx context.Context, name, description string, userID uuid.UUID) (*model.Category, error)
	Update(ctx context.Context, uid uuid.UUID, name, description string, userID uuid.UUID) (*model.Category, error)
	Delete(ctx context.Context, uid uuid.UUID) error
}
