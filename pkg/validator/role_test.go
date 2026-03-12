package validator

import (
	"testing"
)

func TestValidateRole_ValidFull(t *testing.T) {
	yamlContent := `Version: v2
Role:
  TypeId: SqlStorage
  Name: sqlstorage
  Version: 101.20.18991
  Description: SqlStorage
  Definitions:
    ImageFile: "images/SqlStorage_101.20.18991.tar"
    Images:
      - Id: storage-postgres
      - Id: storage-pgadmin
    TemplateFiles:
      - storage-postgres/config/default.env
`

	result := ValidateRole(yamlContent)
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidateRole_ValidMinimal(t *testing.T) {
	yamlContent := `Version: v2
Role:
  TypeId: SqlStorage
  Name: sqlstorage
  Version: 101.20.18991
`

	result := ValidateRole(yamlContent)
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidateRole_MissingVersion(t *testing.T) {
	yamlContent := `Role:
  TypeId: SqlStorage
  Name: sqlstorage
  Version: 101.20.18991
`

	result := ValidateRole(yamlContent)
	if result.Valid {
		t.Errorf("expected invalid, got valid")
	}
	if len(result.Errors) == 0 {
		t.Errorf("expected at least one error")
	}
}

func TestValidateRole_MissingRoleTypeId(t *testing.T) {
	yamlContent := `Version: v2
Role:
  Name: sqlstorage
  Version: 101.20.18991
`

	result := ValidateRole(yamlContent)
	if result.Valid {
		t.Errorf("expected invalid, got valid")
	}
}

func TestValidateRole_MissingRoleName(t *testing.T) {
	yamlContent := `Version: v2
Role:
  TypeId: SqlStorage
  Version: 101.20.18991
`

	result := ValidateRole(yamlContent)
	if result.Valid {
		t.Errorf("expected invalid, got valid")
	}
}

func TestValidateRole_MissingRoleVersion(t *testing.T) {
	yamlContent := `Version: v2
Role:
  TypeId: SqlStorage
  Name: sqlstorage
`

	result := ValidateRole(yamlContent)
	if result.Valid {
		t.Errorf("expected invalid, got valid")
	}
}

func TestValidateRole_InvalidYaml(t *testing.T) {
	yamlContent := `invalid: yaml: content:`

	result := ValidateRole(yamlContent)
	if result.Valid {
		t.Errorf("expected invalid, got valid")
	}
	if len(result.Errors) == 0 {
		t.Errorf("expected at least one error")
	}
}

func TestValidateRole_EmptyDefinitions(t *testing.T) {
	yamlContent := `Version: v2
Role:
  TypeId: SqlStorage
  Name: sqlstorage
  Version: 101.20.18991
  Definitions: {}
`

	result := ValidateRole(yamlContent)
	if !result.Valid {
		t.Errorf("expected valid, got errors: %v", result.Errors)
	}
}

func TestValidateRole_InvalidImageFileType(t *testing.T) {
	yamlContent := `Version: v2
Role:
  TypeId: SqlStorage
  Name: sqlstorage
  Version: 101.20.18991
  Definitions:
    ImageFile: 12345
`

	result := ValidateRole(yamlContent)
	if result.Valid {
		t.Errorf("expected invalid, got valid")
	}
}

func TestValidateRole_InvalidImagesArray(t *testing.T) {
	yamlContent := `Version: v2
Role:
  TypeId: SqlStorage
  Name: sqlstorage
  Version: 101.20.18991
  Definitions:
    Images:
      - name: storage-postgres
`

	result := ValidateRole(yamlContent)
	if result.Valid {
		t.Errorf("expected invalid, got valid")
	}
}
