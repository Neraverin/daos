package validator

import (
	"fmt"
	"strings"

	"github.com/santhosh-tekuri/jsonschema/v6"
	"gopkg.in/yaml.v3"
)

var roleSchema = `
{
  "$schema": "http://json-schema.org/draft-07/schema#",
  "$id": "https://daos.local/schemas/role.schema.json",
  "title": "Role Schema",
  "description": "Schema for validating role YAML configuration files",
  "type": "object",
  "required": ["Version", "Role"],
  "properties": {
    "Version": {
      "type": "string",
      "pattern": "^v[0-9]+$",
      "description": "Schema version (must start with 'v' followed by number)"
    },
    "Role": {
      "type": "object",
      "required": ["TypeId", "Name", "Version"],
      "properties": {
        "TypeId": {
          "type": "string",
          "minLength": 1,
          "description": "Role type identifier"
        },
        "Name": {
          "type": "string",
          "minLength": 1,
          "description": "Role name"
        },
        "Version": {
          "type": "string",
          "minLength": 1,
          "description": "Role version"
        },
        "Description": {
          "type": "string",
          "description": "Role description (optional)"
        },
        "Definitions": {
          "type": "object",
          "properties": {
            "ImageFile": {
              "type": "string",
              "description": "Path to role image tar file"
            },
            "Images": {
              "type": "array",
              "items": {
                "type": "object",
                "required": ["Id"],
                "properties": {
                  "Id": {
                    "type": "string",
                    "minLength": 1,
                    "description": "Image identifier"
                  }
                }
              }
            },
            "TemplateFiles": {
              "type": "array",
              "items": {
                "type": "string",
                "description": "Template file path"
              }
            }
          }
        }
      }
    }
  }
}
`

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type ValidationResult struct {
	Valid  bool              `json:"valid"`
	Errors []ValidationError `json:"errors,omitempty"`
}

var compiledSchema *jsonschema.Schema

func init() {
	compiler := jsonschema.NewCompiler()
	schema, err := jsonschema.UnmarshalJSON(strings.NewReader(roleSchema))
	if err != nil {
		panic(fmt.Sprintf("failed to parse role schema: %v", err))
	}
	err = compiler.AddResource("role.schema.json", schema)
	if err != nil {
		panic(fmt.Sprintf("failed to add role schema: %v", err))
	}
	compiledSchema, err = compiler.Compile("role.schema.json")
	if err != nil {
		panic(fmt.Sprintf("failed to compile role schema: %v", err))
	}
}

func GetSchema() string {
	return roleSchema
}

func ValidateRole(yamlContent string) ValidationResult {
	var data interface{}
	decoder := yaml.NewDecoder(strings.NewReader(yamlContent))
	if err := decoder.Decode(&data); err != nil {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "",
					Message: fmt.Sprintf("failed to parse YAML: %v", err),
				},
			},
		}
	}

	err := compiledSchema.Validate(data)
	if err == nil {
		return ValidationResult{Valid: true}
	}

	validationErr, ok := err.(*jsonschema.ValidationError)
	if !ok {
		return ValidationResult{
			Valid: false,
			Errors: []ValidationError{
				{
					Field:   "",
					Message: err.Error(),
				},
			},
		}
	}

	validationErrors := []ValidationError{}
	if validationErr != nil {
		field := extractFieldFromError(validationErr)
		validationErrors = append(validationErrors, ValidationError{
			Field:   field,
			Message: validationErr.Error(),
		})
	}

	return ValidationResult{
		Valid:  false,
		Errors: validationErrors,
	}
}

func extractFieldFromError(err *jsonschema.ValidationError) string {
	instancePath := err.InstanceLocation
	if len(instancePath) == 0 {
		return "root"
	}
	return strings.Join(instancePath, ".")
}
