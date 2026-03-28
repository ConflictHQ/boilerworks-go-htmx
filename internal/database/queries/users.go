package queries

import (
	"context"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserQueries struct {
	pool *pgxpool.Pool
}

func NewUserQueries(pool *pgxpool.Pool) *UserQueries {
	return &UserQueries{pool: pool}
}

func (q *UserQueries) GetByID(ctx context.Context, id uuid.UUID) (*model.User, error) {
	var u model.User
	err := q.pool.QueryRow(ctx,
		`SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (q *UserQueries) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var u model.User
	err := q.pool.QueryRow(ctx,
		`SELECT id, name, email, password_hash, created_at, updated_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (q *UserQueries) Create(ctx context.Context, name, email, passwordHash string) (*model.User, error) {
	var u model.User
	err := q.pool.QueryRow(ctx,
		`INSERT INTO users (id, name, email, password_hash, created_at, updated_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, NOW(), NOW())
		 RETURNING id, name, email, password_hash, created_at, updated_at`,
		name, email, passwordHash,
	).Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

func (q *UserQueries) List(ctx context.Context, limit, offset int) ([]model.User, int, error) {
	var total int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM users`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.pool.Query(ctx,
		`SELECT id, name, email, password_hash, created_at, updated_at FROM users ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Name, &u.Email, &u.PasswordHash, &u.CreatedAt, &u.UpdatedAt); err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, nil
}

func (q *UserQueries) GetPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT DISTINCT p.name
		 FROM permissions p
		 JOIN group_permissions gp ON gp.permission_id = p.id
		 JOIN user_groups ug ON ug.group_id = gp.group_id
		 WHERE ug.user_id = $1`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}
	return perms, nil
}
