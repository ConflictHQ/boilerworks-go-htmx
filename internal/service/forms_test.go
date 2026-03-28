package service

import (
	"testing"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
)

func TestFormServiceValidateSubmission(t *testing.T) {
	svc := NewFormService()

	def := &model.FormDefinition{
		Schema: []model.FormField{
			{Name: "name", Label: "Name", Type: "text", Required: true},
			{Name: "email", Label: "Email", Type: "email", Required: true},
			{Name: "notes", Label: "Notes", Type: "textarea", Required: false},
		},
	}

	t.Run("valid submission", func(t *testing.T) {
		data := map[string]string{
			"name":  "John Doe",
			"email": "john@example.com",
			"notes": "Some notes",
		}

		jsonData, errs := svc.ValidateSubmission(def, data)
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
		if jsonData == nil {
			t.Error("expected JSON data, got nil")
		}
	})

	t.Run("missing required field", func(t *testing.T) {
		data := map[string]string{
			"name":  "",
			"email": "john@example.com",
		}

		_, errs := svc.ValidateSubmission(def, data)
		if len(errs) == 0 {
			t.Error("expected validation errors for missing required field")
		}
	})

	t.Run("invalid email", func(t *testing.T) {
		data := map[string]string{
			"name":  "John",
			"email": "not-an-email",
		}

		_, errs := svc.ValidateSubmission(def, data)
		if len(errs) == 0 {
			t.Error("expected validation error for invalid email")
		}
	})

	t.Run("optional field can be empty", func(t *testing.T) {
		data := map[string]string{
			"name":  "John",
			"email": "john@example.com",
			"notes": "",
		}

		jsonData, errs := svc.ValidateSubmission(def, data)
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
		if jsonData == nil {
			t.Error("expected JSON data, got nil")
		}
	})
}

func TestFormServiceSelectValidation(t *testing.T) {
	svc := NewFormService()

	def := &model.FormDefinition{
		Schema: []model.FormField{
			{Name: "color", Label: "Color", Type: "select", Required: true, Options: []string{"red", "blue", "green"}},
		},
	}

	t.Run("valid option", func(t *testing.T) {
		data := map[string]string{"color": "red"}
		_, errs := svc.ValidateSubmission(def, data)
		if len(errs) > 0 {
			t.Errorf("expected no errors, got: %v", errs)
		}
	})

	t.Run("invalid option", func(t *testing.T) {
		data := map[string]string{"color": "purple"}
		_, errs := svc.ValidateSubmission(def, data)
		if len(errs) == 0 {
			t.Error("expected validation error for invalid select option")
		}
	})
}
