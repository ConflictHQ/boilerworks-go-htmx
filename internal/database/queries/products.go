package queries

import (
	"context"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductQueries struct {
	pool *pgxpool.Pool
}

func NewProductQueries(pool *pgxpool.Pool) *ProductQueries {
	return &ProductQueries{pool: pool}
}

func (q *ProductQueries) List(ctx context.Context, limit, offset int) ([]model.Product, int, error) {
	var total int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM products WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.pool.Query(ctx,
		`SELECT p.id, p.uuid, p.name, p.description, p.price, p.status, p.category_id,
		        p.created_by, p.updated_by, p.created_at, p.updated_at,
		        COALESCE(c.name, '') as category_name
		 FROM products p
		 LEFT JOIN categories c ON c.id = p.category_id AND c.deleted_at IS NULL
		 WHERE p.deleted_at IS NULL
		 ORDER BY p.created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.UUID, &p.Name, &p.Description, &p.Price, &p.Status,
			&p.CategoryID, &p.CreatedBy, &p.UpdatedBy, &p.CreatedAt, &p.UpdatedAt,
			&p.CategoryName); err != nil {
			return nil, 0, err
		}
		products = append(products, p)
	}
	return products, total, nil
}

func (q *ProductQueries) GetByUUID(ctx context.Context, uid uuid.UUID) (*model.Product, error) {
	var p model.Product
	err := q.pool.QueryRow(ctx,
		`SELECT p.id, p.uuid, p.name, p.description, p.price, p.status, p.category_id,
		        p.created_by, p.updated_by, p.created_at, p.updated_at,
		        COALESCE(c.name, '') as category_name
		 FROM products p
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

func (q *ProductQueries) Create(ctx context.Context, name, description string, price float64, status string, categoryID *uuid.UUID, userID uuid.UUID) (*model.Product, error) {
	var p model.Product
	err := q.pool.QueryRow(ctx,
		`INSERT INTO products (id, uuid, name, description, price, status, category_id, created_by, updated_by, created_at, updated_at)
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

func (q *ProductQueries) Update(ctx context.Context, uid uuid.UUID, name, description string, price float64, status string, categoryID *uuid.UUID, userID uuid.UUID) (*model.Product, error) {
	var p model.Product
	err := q.pool.QueryRow(ctx,
		`UPDATE products SET name = $1, description = $2, price = $3, status = $4, category_id = $5, updated_by = $6, updated_at = NOW()
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

func (q *ProductQueries) Delete(ctx context.Context, uid uuid.UUID) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE products SET deleted_at = NOW() WHERE uuid = $1 AND deleted_at IS NULL`,
		uid,
	)
	return err
}

func (q *ProductQueries) CountByStatus(ctx context.Context) (map[string]int, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT status, COUNT(*) FROM products WHERE deleted_at IS NULL GROUP BY status`,
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
