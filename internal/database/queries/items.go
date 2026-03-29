package queries

import (
	"context"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ItemQueries struct {
	pool *pgxpool.Pool
}

func NewItemQueries(pool *pgxpool.Pool) *ItemQueries {
	return &ItemQueries{pool: pool}
}

func (q *ItemQueries) List(ctx context.Context, limit, offset int) ([]model.Item, int, error) {
	var total int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM items WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.pool.Query(ctx,
		`SELECT p.id, p.uuid, p.name, p.description, p.price, p.status, p.category_id,
		        p.created_by, p.updated_by, p.created_at, p.updated_at,
		        COALESCE(c.name, '') as category_name
		 FROM items p
		 LEFT JOIN categories c ON c.id = p.category_id AND c.deleted_at IS NULL
		 WHERE p.deleted_at IS NULL
		 ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var items []model.Item
	for rows.Next() {
		var p model.Item
		if err := rows.Scan(&p.ID, &p.UUID, &p.Name, &p.Description, &p.Price, &p.Status,
			&p.CategoryID, &p.CreatedBy, &p.UpdatedBy, &p.CreatedAt, &p.UpdatedAt,
			&p.CategoryName); err != nil {
			return nil, 0, err
		}
		items = append(items, p)
	}
	return items, total, nil
}

func (q *ItemQueries) GetByUUID(ctx context.Context, uid uuid.UUID) (*model.Item, error) {
	var p model.Item
	err := q.pool.QueryRow(ctx,
		`SELECT p.id, p.uuid, p.name, p.description, p.price, p.status, p.category_id,
		        p.created_by, p.updated_by, p.created_at, p.updated_at,
		        COALESCE(c.name, '') as category_name
		 FROM items p
		 LEFT JOIN categories c ON c.id = p.category_id AND c.deleted_at IS NULL
		 WHERE p.uuid = $1 AND p.deleted_at IS NULL`,
		uid,
	).Scan(&p.ID, &p.UUID, &p.Name, &p.Description, &p.Price, &p.Status,
		&p.CategoryID, &p.CreatedBy, &p.UpdatedBy, &p.CreatedAt, &p.UpdatedAt,
		&p.CategoryName)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (q *ItemQueries) Create(ctx context.Context, name, description string, price float64, status string, categoryID *uuid.UUID, userID uuid.UUID) (*model.Item, error) {
	var p model.Item
	err := q.pool.QueryRow(ctx,
		`INSERT INTO items (id, uuid, name, description, price, status, category_id, created_by, updated_by, created_at, updated_at)
		 VALUES (gen_random_uuid(), gen_random_uuid(), $1, $2, $3, $4, $5, $6, $6, NOW(), NOW())
		 RETURNING id, uuid, name, description, price, status, category_id, created_by, updated_by, created_at, updated_at`,
		name, description, price, status, categoryID, userID,
	).Scan(&p.ID, &p.UUID, &p.Name, &p.Description, &p.Price, &p.Status,
		&p.CategoryID, &p.CreatedBy, &p.UpdatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (q *ItemQueries) Update(ctx context.Context, uid uuid.UUID, name, description string, price float64, status string, categoryID *uuid.UUID, userID uuid.UUID) (*model.Item, error) {
	var p model.Item
	err := q.pool.QueryRow(ctx,
		`UPDATE items SET name = $1, description = $2, price = $3, status = $4, category_id = $5, updated_by = $6, updated_at = NOW()
		 WHERE uuid = $7 AND deleted_at IS NULL
		 RETURNING id, uuid, name, description, price, status, category_id, created_by, updated_by, created_at, updated_at`,
		name, description, price, status, categoryID, userID, uid,
	).Scan(&p.ID, &p.UUID, &p.Name, &p.Description, &p.Price, &p.Status,
		&p.CategoryID, &p.CreatedBy, &p.UpdatedBy, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (q *ItemQueries) Delete(ctx context.Context, uid uuid.UUID) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE items SET deleted_at = NOW() WHERE uuid = $1 AND deleted_at IS NULL`,
		uid,
	)
	return err
}

func (q *ItemQueries) CountByStatus(ctx context.Context) (map[string]int, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT status, COUNT(*) FROM items WHERE deleted_at IS NULL GROUP BY status`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	counts := make(map[string]int)
	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		counts[status] = count
	}
	return counts, nil
}
