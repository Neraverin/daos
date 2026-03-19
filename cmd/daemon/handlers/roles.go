package handlers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gopkg.in/yaml.v3"
)

func (s *Server) ListRoles(ctx *gin.Context) {
	roles, err := s.db.GetAllRoles(ctx)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var result []api.RoleSummary
	for _, r := range roles {
		result = append(result, roleSummaryToAPI(r))
	}

	if result == nil {
		result = []api.RoleSummary{}
	}
	ctx.JSON(http.StatusOK, result)
}

func (s *Server) CreateRole(ctx *gin.Context) {
	var input api.RoleInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	if input.Name == "" || input.RolePath == "" {
		api.ErrorJSON(ctx, http.StatusBadRequest, "name and role_path are required")
		return
	}

	if !filepath.IsAbs(input.RolePath) {
		api.ErrorJSON(ctx, http.StatusBadRequest, "folder path must be absolute")
		return
	}

	info, err := os.Stat(input.RolePath)
	if os.IsNotExist(err) {
		api.ErrorJSON(ctx, http.StatusNotFound, "role folder not found at path")
		return
	}
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	if !info.IsDir() {
		api.ErrorJSON(ctx, http.StatusBadRequest, "role_path must be a directory")
		return
	}

	if err := checkFilePermissions(input.RolePath); err != nil {
		api.ErrorJSON(ctx, http.StatusForbidden, "permission denied")
		return
	}

	imageFile, err := extractImageFile(input.RolePath)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, "failed to parse role YAML: "+err.Error())
		return
	}

	if imageFile != "" {
		err = s.imageProcessor.ValidateImageFile(input.RolePath, imageFile)
		if err != nil {
			api.ErrorJSON(ctx, http.StatusBadRequest, "invalid image file: "+err.Error())
			return
		}

		result, err := s.imageProcessor.ProcessImage(input.RolePath, imageFile)
		if err != nil {
			api.ErrorJSON(ctx, http.StatusInternalServerError, "failed to process image: "+err.Error())
			return
		}
		if result != nil {
			log.Printf("Processed image: %s (ID: %s)", result.ImageName, result.ImageID)
		}
	}

	id := uuid.New().String()

	r, err := s.db.CreateRole(ctx, db.CreateRoleParams{
		ID:       id,
		Name:     input.Name,
		RolePath: input.RolePath,
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, createRoleRowToAPI(r))
}

func checkFilePermissions(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	file.Close()
	return nil
}

type roleYAML struct {
	Version string `yaml:"Version"`
	Role    struct {
		TypeId      string `yaml:"TypeId"`
		Definitions struct {
			ImageFile string `yaml:"ImageFile"`
		} `yaml:"Definitions"`
	} `yaml:"Role"`
}

func findRoleYAML(folderPath string) (string, error) {
	entries, err := os.ReadDir(folderPath)
	if err != nil {
		return "", err
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		name := entry.Name()
		ext := filepath.Ext(name)
		if ext != ".yaml" && ext != ".yml" {
			continue
		}

		filePath := filepath.Join(folderPath, name)
		data, err := os.ReadFile(filePath)
		if err != nil {
			continue
		}

		var role roleYAML
		if err := yaml.Unmarshal(data, &role); err != nil {
			continue
		}

		if role.Role.TypeId != "" {
			return filePath, nil
		}
	}

	return "", nil
}

func extractImageFile(rolePath string) (string, error) {
	yamlPath, err := findRoleYAML(rolePath)
	if err != nil {
		return "", err
	}

	if yamlPath == "" {
		return "", nil
	}

	data, err := os.ReadFile(yamlPath)
	if err != nil {
		return "", err
	}

	var role roleYAML
	if err := yaml.Unmarshal(data, &role); err != nil {
		return "", err
	}

	return role.Role.Definitions.ImageFile, nil
}

func (s *Server) GetRole(ctx *gin.Context, id uuid.UUID) {
	r, err := s.db.GetRole(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "role not found")
		return
	}

	ctx.JSON(http.StatusOK, getRoleRowToAPI(r))
}

func (s *Server) DeleteRole(ctx *gin.Context, id uuid.UUID) {
	err := s.db.DeleteRole(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "role not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func createRoleRowToAPI(r db.CreateRoleRow) api.Role {
	parsedID, _ := uuid.Parse(r.ID)
	return api.Role{
		Id:        &parsedID,
		Name:      &r.Name,
		RolePath:  &r.RolePath,
		CreatedAt: parseTime(r.CreatedAt),
		UpdatedAt: parseTime(r.UpdatedAt),
	}
}

func getRoleRowToAPI(r db.GetRoleRow) api.Role {
	parsedID, _ := uuid.Parse(r.ID)
	return api.Role{
		Id:        &parsedID,
		Name:      &r.Name,
		RolePath:  &r.RolePath,
		CreatedAt: parseTime(r.CreatedAt),
		UpdatedAt: parseTime(r.UpdatedAt),
	}
}

func roleSummaryToAPI(r db.GetAllRolesRow) api.RoleSummary {
	parsedID, _ := uuid.Parse(r.ID)
	return api.RoleSummary{
		Id:        &parsedID,
		Name:      &r.Name,
		CreatedAt: parseTime(r.CreatedAt),
		UpdatedAt: parseTime(r.UpdatedAt),
	}
}
