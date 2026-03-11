package handlers

import (
	"net/http"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	if input.Name == "" || input.ComposeContent == "" {
		api.ErrorJSON(ctx, http.StatusBadRequest, "name and compose_content are required")
		return
	}

	if len(input.ComposeContent) > 10*1024*1024 {
		api.ErrorJSON(ctx, http.StatusBadRequest, "compose_content exceeds 10MB limit")
		return
	}

	id := uuid.New().String()

	r, err := s.db.CreateRole(ctx, db.CreateRoleParams{
		ID:             id,
		Name:           input.Name,
		ComposeContent: input.ComposeContent,
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, roleToAPI(r))
}

func (s *Server) GetRole(ctx *gin.Context, id uuid.UUID) {
	r, err := s.db.GetRole(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "role not found")
		return
	}

	ctx.JSON(http.StatusOK, roleToAPI(r))
}

func (s *Server) DeleteRole(ctx *gin.Context, id uuid.UUID) {
	err := s.db.DeleteRole(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "role not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func roleToAPI(r db.Role) api.Role {
	parsedID, _ := uuid.Parse(r.ID)
	return api.Role{
		Id:             &parsedID,
		Name:           &r.Name,
		ComposeContent: &r.ComposeContent,
		CreatedAt:      parseTime(r.CreatedAt),
		UpdatedAt:      parseTime(r.UpdatedAt),
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
