package service

import (
	"context"
	"fmt"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/database/queries"
	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
	"github.com/google/uuid"
)

type WorkflowService struct {
	wq *queries.WorkflowQueries
}

func NewWorkflowService(wq *queries.WorkflowQueries) *WorkflowService {
	return &WorkflowService{wq: wq}
}

// GetAvailableTransitions returns valid transitions from the current state.
func (s *WorkflowService) GetAvailableTransitions(def *model.WorkflowDefinition, currentState string) []model.WorkflowTransition {
	var available []model.WorkflowTransition
	for _, t := range def.Transitions {
		if t.From == currentState {
			available = append(available, t)
		}
	}
	return available
}

// GetInitialState returns the first state defined in the workflow.
func (s *WorkflowService) GetInitialState(def *model.WorkflowDefinition) (string, error) {
	if len(def.States) == 0 {
		return "", fmt.Errorf("workflow has no states defined")
	}
	return def.States[0].Name, nil
}

// Transition moves a workflow instance from one state to another.
func (s *WorkflowService) Transition(ctx context.Context, instance *model.WorkflowInstance, def *model.WorkflowDefinition, transitionName string, userID uuid.UUID) error {
	var target *model.WorkflowTransition
	for _, t := range def.Transitions {
		if t.Name == transitionName && t.From == instance.CurrentState {
			target = &t
			break
		}
	}

	if target == nil {
		return fmt.Errorf("transition %q not available from state %q", transitionName, instance.CurrentState)
	}

	// Update instance state
	if err := s.wq.UpdateInstanceState(ctx, instance.ID, target.To); err != nil {
		return fmt.Errorf("update instance state: %w", err)
	}

	// Log the transition
	if err := s.wq.CreateTransitionLog(ctx, instance.ID, instance.CurrentState, target.To, userID); err != nil {
		return fmt.Errorf("create transition log: %w", err)
	}

	return nil
}
