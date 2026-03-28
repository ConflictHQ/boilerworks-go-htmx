package queries

import (
	"context"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CategoryQueries struct {
	pool *pgxpool.Pool
}

func NewCategoryQueries(pool *pgxpool.Pool) *CategoryQueries {
	return &CategoryQueries{pool: pool}
}

func (q *CategoryQueries) List(ctx context.Context, limit, offset int) ([]model.Category, int, error) {
	var total int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM categories WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.pool.Query(ctx,
		`SELECT id, uuid, name, description, created_by, updated_by, created_at, updated_at
		 FROM categories WHERE deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.UUID, &c.Name, &c.Description, &c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, 0, err
		}
		categories = append(categories, c)
	}
	return categories, total, nil
}

func (q *CategoryQueries) GetByUUID(ctx context.Context, uid uuid.UUID) (*model.Category, error) {
	var c model.Category
	err := q.pool.QueryRow(ctx,
		`SELECT id, uuid, name, description, created_by, updated_by, created_at, updated_at
		 FROM categories WHERE uuid = $1 AND deleted_at IS NULL`,
		uid,
	).Scan(&c.ID, &c.UUID, &c.Name, &c.Description, &c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (q *CategoryQueries) Create(ctx context.Context, name, description string, userID uuid.UUID) (*model.Category, error) {
	var c model.Category
	err := q.pool.QueryRow(ctx,
		`INSERT INTO categories (id, uuid, name, description, created_by, updated_by, created_at, updated_at)
		 VALUES (gen_random_uuid(), gen_random_uuid(), $1, $2, $3, $3, NOW(), NOW())
		 RETURNING id, uuid, name, description, created_by, updated_by, created_at, updated_at`,
		name, description, userID,
	).Scan(&c.ID, &c.UUID, &c.Name, &c.Description, &c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (q *CategoryQueries) Update(ctx context.Context, uid uuid.UUID, name, description string, userID uuid.UUID) (*model.Category, error) {
	var c model.Category
	err := q.pool.QueryRow(ctx,
		`UPDATE categories SET name = $1, description = $2, updated_by = $3, updated_at = NOW()
		 WHERE uuid = $4 AND deleted_at IS NULL
		 RETURNING id, uuid, name, description, created_by, updated_by, created_at, updated_at`,
		name, description, userID, uid,
	).Scan(&c.ID, &c.UUID, &c.Name, &c.Description, &c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (q *CategoryQueries) Delete(ctx context.Context, uid uuid.UUID) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE categories SET deleted_at = NOW() WHERE uuid = $1 AND deleted_at IS NULL`,
		uid,
	)
	return err
}

func (q *CategoryQueries) ListAll(ctx context.Context) ([]model.Category, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT id, uuid, name, description, created_by, updated_by, created_at, updated_at
		 FROM categories WHERE deleted_at IS NULL ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var categories []model.Category
	for rows.Next() {
		var c model.Category
		if err := rows.Scan(&c.ID, &c.UUID, &c.Name, &c.Description, &c.CreatedBy, &c.UpdatedBy, &c.CreatedAt, &c.UpdatedAt); err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}
	return categories, nil
}
