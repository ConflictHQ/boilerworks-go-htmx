package queries

import (
	"context"
	"encoding/json"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WorkflowQueries struct {
	pool *pgxpool.Pool
}

func NewWorkflowQueries(pool *pgxpool.Pool) *WorkflowQueries {
	return &WorkflowQueries{pool: pool}
}

func (q *WorkflowQueries) ListDefinitions(ctx context.Context, limit, offset int) ([]model.WorkflowDefinition, int, error) {
	var total int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflow_definitions WHERE deleted_at IS NULL`).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.pool.Query(ctx,
		`SELECT id, uuid, name, description, status, states, transitions, created_by, updated_by, created_at, updated_at
		 FROM workflow_definitions WHERE deleted_at IS NULL
		 ORDER BY created_at DESC LIMIT $1 OFFSET $2`,
		limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var defs []model.WorkflowDefinition
	for rows.Next() {
		var d model.WorkflowDefinition
		if err := rows.Scan(&d.ID, &d.UUID, &d.Name, &d.Description, &d.Status,
			&d.StatesJSON, &d.TransitionsJSON, &d.CreatedBy, &d.UpdatedBy, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, 0, err
		}
		_ = json.Unmarshal(d.StatesJSON, &d.States)
		_ = json.Unmarshal(d.TransitionsJSON, &d.Transitions)
		defs = append(defs, d)
	}
	return defs, total, nil
}

func (q *WorkflowQueries) GetDefinitionByUUID(ctx context.Context, uid uuid.UUID) (*model.WorkflowDefinition, error) {
	var d model.WorkflowDefinition
	err := q.pool.QueryRow(ctx,
		`SELECT id, uuid, name, description, status, states, transitions, created_by, updated_by, created_at, updated_at
		 FROM workflow_definitions WHERE uuid = $1 AND deleted_at IS NULL`,
		uid,
	).Scan(&d.ID, &d.UUID, &d.Name, &d.Description, &d.Status,
		&d.StatesJSON, &d.TransitionsJSON, &d.CreatedBy, &d.UpdatedBy, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(d.StatesJSON, &d.States)
	_ = json.Unmarshal(d.TransitionsJSON, &d.Transitions)
	return &d, nil
}

func (q *WorkflowQueries) CreateDefinition(ctx context.Context, name, description, status string, states, transitions json.RawMessage, userID uuid.UUID) (*model.WorkflowDefinition, error) {
	var d model.WorkflowDefinition
	err := q.pool.QueryRow(ctx,
		`INSERT INTO workflow_definitions (id, uuid, name, description, status, states, transitions, created_by, updated_by, created_at, updated_at)
		 VALUES (gen_random_uuid(), gen_random_uuid(), $1, $2, $3, $4, $5, $6, $6, NOW(), NOW())
		 RETURNING id, uuid, name, description, status, states, transitions, created_by, updated_by, created_at, updated_at`,
		name, description, status, states, transitions, userID,
	).Scan(&d.ID, &d.UUID, &d.Name, &d.Description, &d.Status,
		&d.StatesJSON, &d.TransitionsJSON, &d.CreatedBy, &d.UpdatedBy, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(d.StatesJSON, &d.States)
	_ = json.Unmarshal(d.TransitionsJSON, &d.Transitions)
	return &d, nil
}

func (q *WorkflowQueries) UpdateDefinition(ctx context.Context, uid uuid.UUID, name, description, status string, states, transitions json.RawMessage, userID uuid.UUID) (*model.WorkflowDefinition, error) {
	var d model.WorkflowDefinition
	err := q.pool.QueryRow(ctx,
		`UPDATE workflow_definitions SET name = $1, description = $2, status = $3, states = $4, transitions = $5, updated_by = $6, updated_at = NOW()
		 WHERE uuid = $7 AND deleted_at IS NULL
		 RETURNING id, uuid, name, description, status, states, transitions, created_by, updated_by, created_at, updated_at`,
		name, description, status, states, transitions, userID, uid,
	).Scan(&d.ID, &d.UUID, &d.Name, &d.Description, &d.Status,
		&d.StatesJSON, &d.TransitionsJSON, &d.CreatedBy, &d.UpdatedBy, &d.CreatedAt, &d.UpdatedAt)
	if err != nil {
		return nil, err
	}
	_ = json.Unmarshal(d.StatesJSON, &d.States)
	_ = json.Unmarshal(d.TransitionsJSON, &d.Transitions)
	return &d, nil
}

func (q *WorkflowQueries) DeleteDefinition(ctx context.Context, uid uuid.UUID) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE workflow_definitions SET deleted_at = NOW() WHERE uuid = $1 AND deleted_at IS NULL`,
		uid,
	)
	return err
}

func (q *WorkflowQueries) ListInstances(ctx context.Context, defID uuid.UUID, limit, offset int) ([]model.WorkflowInstance, int, error) {
	var total int
	err := q.pool.QueryRow(ctx,
		`SELECT COUNT(*) FROM workflow_instances WHERE workflow_definition_id = $1`,
		defID,
	).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	rows, err := q.pool.Query(ctx,
		`SELECT wi.id, wi.uuid, wi.workflow_definition_id, wi.current_state, wi.created_by, wi.created_at,
		        wd.name as workflow_name
		 FROM workflow_instances wi
		 JOIN workflow_definitions wd ON wd.id = wi.workflow_definition_id
		 WHERE wi.workflow_definition_id = $1
		 ORDER BY wi.created_at DESC LIMIT $2 OFFSET $3`,
		defID, limit, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var instances []model.WorkflowInstance
	for rows.Next() {
		var i model.WorkflowInstance
		if err := rows.Scan(&i.ID, &i.UUID, &i.WorkflowDefinitionID, &i.CurrentState, &i.CreatedBy, &i.CreatedAt, &i.WorkflowName); err != nil {
			return nil, 0, err
		}
		instances = append(instances, i)
	}
	return instances, total, nil
}

func (q *WorkflowQueries) GetInstanceByUUID(ctx context.Context, uid uuid.UUID) (*model.WorkflowInstance, error) {
	var i model.WorkflowInstance
	err := q.pool.QueryRow(ctx,
		`SELECT wi.id, wi.uuid, wi.workflow_definition_id, wi.current_state, wi.created_by, wi.created_at,
		        wd.name as workflow_name
		 FROM workflow_instances wi
		 JOIN workflow_definitions wd ON wd.id = wi.workflow_definition_id
		 WHERE wi.uuid = $1`,
		uid,
	).Scan(&i.ID, &i.UUID, &i.WorkflowDefinitionID, &i.CurrentState, &i.CreatedBy, &i.CreatedAt, &i.WorkflowName)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (q *WorkflowQueries) CreateInstance(ctx context.Context, defID uuid.UUID, initialState string, userID uuid.UUID) (*model.WorkflowInstance, error) {
	var i model.WorkflowInstance
	err := q.pool.QueryRow(ctx,
		`INSERT INTO workflow_instances (id, uuid, workflow_definition_id, current_state, created_by, created_at)
		 VALUES (gen_random_uuid(), gen_random_uuid(), $1, $2, $3, NOW())
		 RETURNING id, uuid, workflow_definition_id, current_state, created_by, created_at`,
		defID, initialState, userID,
	).Scan(&i.ID, &i.UUID, &i.WorkflowDefinitionID, &i.CurrentState, &i.CreatedBy, &i.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (q *WorkflowQueries) UpdateInstanceState(ctx context.Context, instanceID uuid.UUID, newState string) error {
	_, err := q.pool.Exec(ctx,
		`UPDATE workflow_instances SET current_state = $1 WHERE id = $2`,
		newState, instanceID,
	)
	return err
}

func (q *WorkflowQueries) CreateTransitionLog(ctx context.Context, instanceID uuid.UUID, fromState, toState string, userID uuid.UUID) error {
	_, err := q.pool.Exec(ctx,
		`INSERT INTO transition_logs (id, workflow_instance_id, from_state, to_state, performed_by, created_at)
		 VALUES (gen_random_uuid(), $1, $2, $3, $4, NOW())`,
		instanceID, fromState, toState, userID,
	)
	return err
}

func (q *WorkflowQueries) GetTransitionLogs(ctx context.Context, instanceID uuid.UUID) ([]model.TransitionLog, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT tl.id, tl.workflow_instance_id, tl.from_state, tl.to_state, tl.performed_by, tl.created_at,
		        u.name as performed_by_name
		 FROM transition_logs tl
		 JOIN users u ON u.id = tl.performed_by
		 WHERE tl.workflow_instance_id = $1
		 ORDER BY tl.created_at ASC`,
		instanceID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var logs []model.TransitionLog
	for rows.Next() {
		var l model.TransitionLog
		if err := rows.Scan(&l.ID, &l.WorkflowInstanceID, &l.FromState, &l.ToState, &l.PerformedBy, &l.CreatedAt, &l.PerformedByName); err != nil {
			return nil, err
		}
		logs = append(logs, l)
	}
	return logs, nil
}

func (q *WorkflowQueries) CountDefinitions(ctx context.Context) (int, error) {
	var count int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflow_definitions WHERE deleted_at IS NULL`).Scan(&count)
	return count, err
}

func (q *WorkflowQueries) CountInstances(ctx context.Context) (int, error) {
	var count int
	err := q.pool.QueryRow(ctx, `SELECT COUNT(*) FROM workflow_instances`).Scan(&count)
	return count, err
}

func (q *WorkflowQueries) ListAllDefinitions(ctx context.Context) ([]model.WorkflowDefinition, error) {
	rows, err := q.pool.Query(ctx,
		`SELECT id, uuid, name, description, status, states, transitions, created_by, updated_by, created_at, updated_at
		 FROM workflow_definitions WHERE deleted_at IS NULL AND status = 'active'
		 ORDER BY name ASC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var defs []model.WorkflowDefinition
	for rows.Next() {
		var d model.WorkflowDefinition
		if err := rows.Scan(&d.ID, &d.UUID, &d.Name, &d.Description, &d.Status,
			&d.StatesJSON, &d.TransitionsJSON, &d.CreatedBy, &d.UpdatedBy, &d.CreatedAt, &d.UpdatedAt); err != nil {
			return nil, err
		}
		_ = json.Unmarshal(d.StatesJSON, &d.States)
		_ = json.Unmarshal(d.TransitionsJSON, &d.Transitions)
		defs = append(defs, d)
	}
	return defs, nil
}
