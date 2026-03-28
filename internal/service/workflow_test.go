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

func TestIsTerminalState(t *testing.T) {
	def := newTestWorkflowDef()

	tests := []struct {
		name     string
		state    string
		terminal bool
	}{
		{"pending is not terminal", "pending", false},
		{"active is not terminal", "active", false},
		{"completed is terminal", "completed", true},
		{"unknown is not terminal", "unknown", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isEnd := false
			for _, s := range def.States {
				if s.Name == tt.state && s.IsEnd {
					isEnd = true
					break
				}
			}
			if isEnd != tt.terminal {
				t.Errorf("expected terminal=%v for state %q, got %v", tt.terminal, tt.state, isEnd)
			}
		})
	}
}

func TestTransitionValidPath(t *testing.T) {
	svc := NewWorkflowService(nil)
	def := newTestWorkflowDef()

	// Simulate walking the full workflow path: pending -> active -> completed
	state := "pending"
	transitions := svc.GetAvailableTransitions(def, state)
	if len(transitions) != 1 || transitions[0].To != "active" {
		t.Fatalf("expected transition to active from pending")
	}

	state = transitions[0].To
	transitions = svc.GetAvailableTransitions(def, state)
	if len(transitions) != 1 || transitions[0].To != "completed" {
		t.Fatalf("expected transition to completed from active")
	}

	state = transitions[0].To
	transitions = svc.GetAvailableTransitions(def, state)
	if len(transitions) != 0 {
		t.Error("expected no transitions from terminal state completed")
	}
}

func TestTransitionInvalidFromState(t *testing.T) {
	svc := NewWorkflowService(nil)
	def := newTestWorkflowDef()

	// "complete" transition should not be available from "pending"
	transitions := svc.GetAvailableTransitions(def, "pending")
	for _, tr := range transitions {
		if tr.Name == "complete" {
			t.Error("complete transition should not be available from pending")
		}
	}
}

func TestMultipleTransitionsFromOneState(t *testing.T) {
	svc := NewWorkflowService(nil)

	def := &model.WorkflowDefinition{
		States: []model.WorkflowState{
			{Name: "draft", Label: "Draft"},
			{Name: "review", Label: "Review"},
			{Name: "published", Label: "Published"},
			{Name: "archived", Label: "Archived", IsEnd: true},
		},
		Transitions: []model.WorkflowTransition{
			{Name: "submit", From: "draft", To: "review"},
			{Name: "publish", From: "review", To: "published"},
			{Name: "reject", From: "review", To: "draft"},
			{Name: "archive", From: "published", To: "archived"},
		},
	}

	transitions := svc.GetAvailableTransitions(def, "review")
	if len(transitions) != 2 {
		t.Fatalf("expected 2 transitions from review, got %d", len(transitions))
	}

	names := map[string]bool{}
	for _, tr := range transitions {
		names[tr.Name] = true
	}
	if !names["publish"] || !names["reject"] {
		t.Errorf("expected publish and reject transitions, got %v", names)
	}
}

func TestGetInitialStateReturnsFirst(t *testing.T) {
	svc := NewWorkflowService(nil)

	def := &model.WorkflowDefinition{
		States: []model.WorkflowState{
			{Name: "open", Label: "Open"},
			{Name: "closed", Label: "Closed", IsEnd: true},
		},
	}

	state, err := svc.GetInitialState(def)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if state != "open" {
		t.Errorf("expected initial state 'open', got '%s'", state)
	}
}
