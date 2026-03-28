package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ConflictHQ/boilerworks-go-htmx/internal/model"
)

type FormService struct{}

func NewFormService() *FormService {
	return &FormService{}
}

// ValidateSubmission validates form data against the form definition schema.
func (s *FormService) ValidateSubmission(def *model.FormDefinition, data map[string]string) (json.RawMessage, []string) {
	var errs []string
	result := make(map[string]interface{})

	for _, field := range def.Schema {
		val := strings.TrimSpace(data[field.Name])

		if field.Required && val == "" {
			errs = append(errs, fmt.Sprintf("%s is required", field.Label))
			continue
		}

		if val == "" {
			continue
		}

		switch field.Type {
		case "email":
			if !strings.Contains(val, "@") || !strings.Contains(val, ".") {
				errs = append(errs, fmt.Sprintf("%s must be a valid email", field.Label))
				continue
			}
		case "select":
			if len(field.Options) > 0 {
				found := false
				for _, opt := range field.Options {
					if opt == val {
						found = true
						break
					}
				}
				if !found {
					errs = append(errs, fmt.Sprintf("%s has an invalid option", field.Label))
					continue
				}
			}
		}

		result[field.Name] = val
	}

	if len(errs) > 0 {
		return nil, errs
	}

	jsonData, _ := json.Marshal(result)
	return jsonData, nil
}
