package service

import (
	"strings"
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

func TestFormServiceMultipleRequiredFieldsMissing(t *testing.T) {
	svc := NewFormService()

	def := &model.FormDefinition{
		Schema: []model.FormField{
			{Name: "first_name", Label: "First Name", Type: "text", Required: true},
			{Name: "last_name", Label: "Last Name", Type: "text", Required: true},
			{Name: "email", Label: "Email", Type: "email", Required: true},
		},
	}

	data := map[string]string{
		"first_name": "",
		"last_name":  "",
		"email":      "",
	}

	_, errs := svc.ValidateSubmission(def, data)
	if len(errs) != 3 {
		t.Errorf("expected 3 validation errors, got %d: %v", len(errs), errs)
	}
}

func TestFormServiceWhitespaceOnlyRequired(t *testing.T) {
	svc := NewFormService()

	def := &model.FormDefinition{
		Schema: []model.FormField{
			{Name: "name", Label: "Name", Type: "text", Required: true},
		},
	}

	data := map[string]string{"name": "   "}

	_, errs := svc.ValidateSubmission(def, data)
	if len(errs) == 0 {
		t.Error("expected validation error for whitespace-only required field")
	}
}

func TestFormServiceEmailEdgeCases(t *testing.T) {
	svc := NewFormService()

	def := &model.FormDefinition{
		Schema: []model.FormField{
			{Name: "email", Label: "Email", Type: "email", Required: true},
		},
	}

	tests := []struct {
		name    string
		value   string
		wantErr bool
	}{
		{"missing at sign", "userexample.com", true},
		{"missing dot", "user@examplecom", true},
		{"valid email", "user@example.com", false},
		{"at sign only", "@", true},
		{"dot only", ".", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := map[string]string{"email": tt.value}
			_, errs := svc.ValidateSubmission(def, data)
			if tt.wantErr && len(errs) == 0 {
				t.Errorf("expected validation error for %q", tt.value)
			}
			if !tt.wantErr && len(errs) > 0 {
				t.Errorf("expected no error for %q, got: %v", tt.value, errs)
			}
		})
	}
}

func TestFormServiceValidDataReturnsJSON(t *testing.T) {
	svc := NewFormService()

	def := &model.FormDefinition{
		Schema: []model.FormField{
			{Name: "name", Label: "Name", Type: "text", Required: true},
			{Name: "email", Label: "Email", Type: "email", Required: true},
		},
	}

	data := map[string]string{
		"name":  "Jane",
		"email": "jane@example.com",
	}

	jsonData, errs := svc.ValidateSubmission(def, data)
	if len(errs) > 0 {
		t.Fatalf("expected no errors, got: %v", errs)
	}

	jsonStr := string(jsonData)
	if !strings.Contains(jsonStr, "Jane") {
		t.Error("expected JSON to contain 'Jane'")
	}
	if !strings.Contains(jsonStr, "jane@example.com") {
		t.Error("expected JSON to contain email")
	}
}

func TestFormServiceEmptySchema(t *testing.T) {
	svc := NewFormService()

	def := &model.FormDefinition{
		Schema: []model.FormField{},
	}

	data := map[string]string{"anything": "value"}

	jsonData, errs := svc.ValidateSubmission(def, data)
	if len(errs) > 0 {
		t.Errorf("expected no errors for empty schema, got: %v", errs)
	}
	if jsonData == nil {
		t.Error("expected JSON data, got nil")
	}
}
