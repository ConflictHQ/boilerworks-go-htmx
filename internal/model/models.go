package model

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID  `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

type Session struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

type Group struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Permission struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID          uuid.UUID  `json:"id"`
	UUID        uuid.UUID  `json:"uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedBy   uuid.UUID  `json:"updated_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type Item struct {
	ID          uuid.UUID  `json:"id"`
	UUID        uuid.UUID  `json:"uuid"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Price       float64    `json:"price"`
	Status      string     `json:"status"`
	CategoryID  *uuid.UUID `json:"category_id,omitempty"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedBy   uuid.UUID  `json:"updated_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`

	// Joined fields
	CategoryName string `json:"category_name,omitempty"`
}

type FormField struct {
	Name     string   `json:"name"`
	Label    string   `json:"label"`
	Type     string   `json:"type"`
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"`
}

type FormDefinition struct {
	ID          uuid.UUID  `json:"id"`
	UUID        uuid.UUID  `json:"uuid"`
	Name        string     `json:"name"`
	Slug        string     `json:"slug"`
	Description string     `json:"description"`
	Status      string     `json:"status"`
	Schema      []FormField `json:"schema"`
	SchemaJSON  json.RawMessage `json:"-"`
	CreatedBy   uuid.UUID  `json:"created_by"`
	UpdatedBy   uuid.UUID  `json:"updated_by"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
}

type FormSubmission struct {
	ID               uuid.UUID       `json:"id"`
	UUID             uuid.UUID       `json:"uuid"`
	FormDefinitionID uuid.UUID       `json:"form_definition_id"`
	Data             json.RawMessage `json:"data"`
	CreatedBy        uuid.UUID       `json:"created_by"`
	CreatedAt        time.Time       `json:"created_at"`
	DeletedAt        *time.Time      `json:"deleted_at,omitempty"`

	// Joined fields
	FormName string `json:"form_name,omitempty"`
}

type WorkflowState struct {
	Name   string `json:"name"`
	Label  string `json:"label"`
	IsEnd  bool   `json:"is_end,omitempty"`
}

type WorkflowTransition struct {
	Name string `json:"name"`
	From string `json:"from"`
	To   string `json:"to"`
}

type WorkflowDefinition struct {
	ID          uuid.UUID            `json:"id"`
	UUID        uuid.UUID            `json:"uuid"`
	Name        string               `json:"name"`
	Description string               `json:"description"`
	Status      string               `json:"status"`
	States      []WorkflowState      `json:"states"`
	Transitions []WorkflowTransition `json:"transitions"`
	StatesJSON      json.RawMessage  `json:"-"`
	TransitionsJSON json.RawMessage  `json:"-"`
	CreatedBy   uuid.UUID            `json:"created_by"`
	UpdatedBy   uuid.UUID            `json:"updated_by"`
	CreatedAt   time.Time            `json:"created_at"`
	UpdatedAt   time.Time            `json:"updated_at"`
	DeletedAt   *time.Time           `json:"deleted_at,omitempty"`
}

type WorkflowInstance struct {
	ID                   uuid.UUID `json:"id"`
	UUID                 uuid.UUID `json:"uuid"`
	WorkflowDefinitionID uuid.UUID `json:"workflow_definition_id"`
	CurrentState         string    `json:"current_state"`
	CreatedBy            uuid.UUID `json:"created_by"`
	CreatedAt            time.Time `json:"created_at"`

	// Joined fields
	WorkflowName string `json:"workflow_name,omitempty"`
}

type TransitionLog struct {
	ID                 uuid.UUID `json:"id"`
	WorkflowInstanceID uuid.UUID `json:"workflow_instance_id"`
	FromState          string    `json:"from_state"`
	ToState            string    `json:"to_state"`
	PerformedBy        uuid.UUID `json:"performed_by"`
	CreatedAt          time.Time `json:"created_at"`

	// Joined fields
	PerformedByName string `json:"performed_by_name,omitempty"`
}

// Pagination holds pagination metadata.
type Pagination struct {
	Page       int
	PerPage    int
	Total      int
	TotalPages int
}

func NewPagination(page, perPage, total int) Pagination {
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 20
	}
	totalPages := total / perPage
	if total%perPage > 0 {
		totalPages++
	}
	return Pagination{
		Page:       page,
		PerPage:    perPage,
		Total:      total,
		TotalPages: totalPages,
	}
}
