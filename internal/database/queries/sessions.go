package queries

import (
	"context"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionQueries struct {
	pool *pgxpool.Pool
}

func NewSessionQueries(pool *pgxpool.Pool) *SessionQueries {
	return &SessionQueries{pool: pool}
}

func (q *SessionQueries) Create(ctx context.Context, userID uuid.UUID, tokenHash string) (*model.Session, error) {
	var s model.Session
	err := q.pool.QueryRow(ctx,
		`INSERT INTO sessions (id, user_id, token_hash, expires_at, created_at)
		 VALUES (gen_random_uuid(), $1, $2, NOW() + INTERVAL '30 days', NOW())
		 RETURNING id, user_id, token_hash, expires_at, created_at`,
		userID, tokenHash,
	).Scan(&s.ID, &s.UserID, &s.TokenHash, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (q *SessionQueries) GetByTokenHash(ctx context.Context, tokenHash string) (*model.Session, error) {
	var s model.Session
	err := q.pool.QueryRow(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at
		 FROM sessions
		 WHERE token_hash = $1 AND expires_at > NOW()`,
		tokenHash,
	).Scan(&s.ID, &s.UserID, &s.TokenHash, &s.ExpiresAt, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (q *SessionQueries) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := q.pool.Exec(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

func (q *SessionQueries) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	_, err := q.pool.Exec(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}

func (q *SessionQueries) DeleteExpired(ctx context.Context) error {
	_, err := q.pool.Exec(ctx, `DELETE FROM sessions WHERE expires_at < NOW()`)
	return err
}
