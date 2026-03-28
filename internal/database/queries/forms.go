package queries

import (
	"context"
	"encoding/json"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type FormQueries struct {
	pool *pgxpool.Pool
}

func NewFormQueries(pool *pgxpool.Pool) *FormQueries {
	return &FormQueries{pool: pool}
}

func (q *FormQueries) ListDefinitions(ctx context.Context, limit, offset int) ([]model.FormDefinition, int, error) {
	var total int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM form_definitions WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.pool.Query(ctx,
		`SELECT id, uuid, name, slug, description, status, schema, created_by, updated_by, created_at, updated_at
		 FROM form_definitions WHERE deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var forms []model.FormDefinition
	for rows.Next() {
		var f model.FormDefinition
		if err := rows.Scan(&f.ID, &f.UUID, &f.Name, &f.Slug, &f.Description, &f.Status,
			&f.SchemaJSON, &f.CreatedBy, &f.UpdatedBy, &f.CreatedAt, &f.UpdatedAt); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(f.SchemaJSON, &f.Schema)
		forms = append(forms, f)
	}
	return forms, total, nil
}

func (q *FormQueries) GetDefinitionByUUID(ctx context.Context, uid uuid.UUID) (*model.FormDefinition, error) {
	var f model.FormDefinition
	err := q.pool.QueryRow(ctx,
		`SELECT id, uuid, name, slug, description, status, schema, created_by, updated_by, created_at, updated_at
		 FROM form_definitions WHERE uuid = $1 AND deleted_at IS NULL`,
		uid,
	).Scan(&f.ID, &f.UUID, &f.Name, &f.Slug, &f.Description, &f.Status,
		&f.SchemaJSON, &f.CreatedBy, &f.UpdatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(f.SchemaJSON, &f.Schema)
	return &f, nil
}

func (q *FormQueries) GetDefinitionBySlug(ctx context.Context, slug string) (*model.FormDefinition, error) {
	var f model.FormDefinition
	err := q.pool.QueryRow(ctx,
		`SELECT id, uuid, name, slug, description, status, schema, created_by, updated_by, created_at, updated_at
		 FROM form_definitions WHERE slug = $1 AND deleted_at IS NULL`,
		slug,
	).Scan(&f.ID, &f.UUID, &f.Name, &f.Slug, &f.Description, &f.Status,
		&f.SchemaJSON, &f.CreatedBy, &f.UpdatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(f.SchemaJSON, &f.Schema)
	return &f, nil
}

func (q *FormQueries) CreateDefinition(ctx context.Context, name, slug, description, status string, schema json.RawMessage, userID uuid.UUID) (*model.FormDefinition, error) {
	var f model.FormDefinition
	err := q.pool.QueryRow(ctx,
		`INSERT INTO form_definitions (id, uuid, name, slug, description, status, schema, created_by, updated_by, created_at, updated_at)
		 VALUES (gen_random_uuid(), gen_random_uuid(), $1, $2, $3, $4, $5, $6, $6, NOW(), NOW())
		 RETURNING id, uuid, name, slug, description, status, schema, created_by, updated_by, created_at, updated_at`,
		name, slug, description, status, schema, userID,
	).Scan(&f.ID, &f.UUID, &f.Name, &f.Slug, &f.Description, &f.Status,
		&f.SchemaJSON, &f.CreatedBy, &f.UpdatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(f.SchemaJSON, &f.Schema)
	return &f, nil
}

func (q *FormQueries) UpdateDefinition(ctx context.Context, uid uuid.UUID, name, slug, description, status string, schema json.RawMessage, userID uuid.UUID) (*model.FormDefinition, error) {
	var f model.FormDefinition
	err := q.pool.QueryRow(ctx,
		`UPDATE form_definitions SET name = $1, slug = $2, description = $3, status = $4, schema = $5, updated_by = $6, updated_at = NOW()
		 WHERE uuid = $7 AND deleted_at IS NULL
		 RETURNING id, uuid, name, slug, description, status, schema, created_by, updated_by, created_at, updated_at`,
		name, slug, description, status, schema, userID, uid,
	).Scan(&f.ID, &f.UUID, &f.Name, &f.Slug, &f.Description, &f.Status,
		&f.SchemaJSON, &f.CreatedBy, &f.UpdatedBy, &f.CreatedAt, &f.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(f.SchemaJSON, &f.Schema)
	return &f, nil
}

func (q *FormQueries) DeleteDefinition(ctx context.Context, uid uuid.UUID) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE form_definitions SET deleted_at = NOW() WHERE uuid = $1 AND deleted_at IS NULL`,
		uid,
	)
	return err
}

func (q *FormQueries) ListSubmissions(ctx context.Context, formDefID uuid.UUID, limit, offset int) ([]model.FormSubmission, int, error) {
	var total int
	err := q.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM form_submissions WHERE form_definition_id = $1 AND deleted_at IS NULL`,
		formDefID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.pool.Query(ctx,
		`SELECT fs.id, fs.uuid, fs.form_definition_id, fs.data, fs.created_by, fs.created_at,
		        fd.name as form_name
		 FROM form_submissions fs
		 JOIN form_definitions fd ON fd.id = fs.form_definition_id
		 WHERE fs.form_definition_id = $1 AND fs.deleted_at IS NULL
		 ORDER BY fs.created_at DESC LIMIT $2 OFFSET $3`,
		formDefID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var submissions []model.FormSubmission
	for rows.Next() {
		var s model.FormSubmission
		if err := rows.Scan(&s.ID, &s.UUID, &s.FormDefinitionID, &s.Data, &s.CreatedBy, &s.CreatedAt, &s.FormName); err != nil {
			return nil, 0, err
		}
		submissions = append(submissions, s)
	}
	return submissions, total, nil
}

func (q *FormQueries) CreateSubmission(ctx context.Context, formDefID uuid.UUID, data json.RawMessage, userID uuid.UUID) (*model.FormSubmission, error) {
	var s model.FormSubmission
	err := q.pool.QueryRow(ctx,
		`INSERT INTO form_submissions (id, uuid, form_definition_id, data, created_by, created_at)
		 VALUES (gen_random_uuid(), gen_random_uuid(), $1, $2, $3, NOW())
		 RETURNING id, uuid, form_definition_id, data, created_by, created_at`,
		formDefID, data, userID,
	).Scan(&s.ID, &s.UUID, &s.FormDefinitionID, &s.Data, &s.CreatedBy, &s.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (q *FormQueries) CountDefinitions(ctx context.Context) (int, error) {
	var count int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM form_definitions WHERE deleted_at IS NULL`).Scan(&count)
	return count, err
}

func (q *FormQueries) CountSubmissions(ctx context.Context) (int, error) {
	var count int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM form_submissions WHERE deleted_at IS NULL`).Scan(&count)
	return count, err
}
