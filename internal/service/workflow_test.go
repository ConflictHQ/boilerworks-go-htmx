package service

import (
	"testing"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
)

func newTestWorkflowDef() *model.WorkflowDefinition {
	return &model.WorkflowDefinition{
		States: []model.WorkflowState{
			{Name: "pending", Label: "Pending"},
			{Name: "active", Label: "Active"},
			{Name: "completed", Label: "Completed", IsEnd: true},
		},
		Transitions: []model.WorkflowTransition{
			{Name: "start", From: "pending", To: "active"},
			{Name: "complete", From: "active", To: "completed"},
		},
	}
}

func TestGetInitialState(t *testing.T) {
	svc := NewWorkflowService(nil)
	def := newTestWorkflowDef()

	state, err := svc.GetInitialState(def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != "pending" {
		t.Errorf("expected initial state 'pending', got '%s'", state)
	}
}

func TestGetInitialStateEmpty(t *testing.T) {
	svc := NewWorkflowService(nil)
	def := &model.WorkflowDefinition{States: []model.WorkflowState{}}

	_, err := svc.GetInitialState(def)
	if err == nil {
		t.Error("expected error for empty states")
	}
}

func TestGetAvailableTransitions(t *testing.T) {
	svc := NewWorkflowService(nil)
	def := newTestWorkflowDef()

	t.Run("from pending", func(t *testing.T) {
		transitions := svc.GetAvailableTransitions(def, "pending")
		if len(transitions) != 1 {
			t.Fatalf("expected 1 transition from pending, got %d", len(transitions))
		}
		if transitions[0].Name != "start" {
			t.Errorf("expected transition 'start', got '%s'", transitions[0].Name)
		}
	})

	t.Run("from active", func(t *testing.T) {
		transitions := svc.GetAvailableTransitions(def, "active")
		if len(transitions) != 1 {
			t.Fatalf("expected 1 transition from active, got %d", len(transitions))
		}
		if transitions[0].Name != "complete" {
			t.Errorf("expected transition 'complete', got '%s'", transitions[0].Name)
		}
	})

	t.Run("from completed (terminal)", func(t *testing.T) {
		transitions := svc.GetAvailableTransitions(def, "completed")
		if len(transitions) != 0 {
			t.Errorf("expected 0 transitions from completed, got %d", len(transitions))
		}
	})

	t.Run("from unknown state", func(t *testing.T) {
		transitions := svc.GetAvailableTransitions(def, "unknown")
		if len(transitions) != 0 {
			t.Errorf("expected 0 transitions from unknown, got %d", len(transitions))
		}
	})
}
