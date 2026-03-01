package handlers

import (
	"net/http"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/gin-gonic/gin"
)

func (s *Server) ListPackages(ctx *gin.Context) {
	rows, err := s.db.Query("SELECT id, name, created_at, updated_at FROM packages ORDER BY name")
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var packages []api.PackageSummary
	for rows.Next() {
		var p api.PackageSummary
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatedAt, &p.UpdatedAt); err != nil {
			api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		packages = append(packages, p)
	}

	if packages == nil {
		packages = []api.PackageSummary{}
	}
	ctx.JSON(http.StatusOK, packages)
}

func (s *Server) CreatePackage(ctx *gin.Context) {
	var input api.PackageInput
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

	result, err := s.db.Exec(
		"INSERT INTO packages (name, compose_content) VALUES (?, ?)",
		input.Name, input.ComposeContent,
	)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	s.GetPackage(ctx, int(id))
}

func (s *Server) GetPackage(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid package id")
		return
	}

	var p api.Package
	err := s.db.QueryRow(
		"SELECT id, name, compose_content, created_at, updated_at FROM packages WHERE id = ?",
		id,
	).Scan(&p.ID, &p.Name, &p.ComposeContent, &p.CreatedAt, &p.UpdatedAt)

	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "package not found")
		return
	}

	ctx.JSON(http.StatusOK, p)
}

func (s *Server) DeletePackage(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid package id")
		return
	}

	result, err := s.db.Exec("DELETE FROM packages WHERE id = ?", id)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		api.ErrorJSON(ctx, http.StatusNotFound, "package not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}
