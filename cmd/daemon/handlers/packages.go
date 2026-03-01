package handlers

import (
	"net/http"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/db"
	"github.com/gin-gonic/gin"
)

func (s *Server) ListPackages(ctx *gin.Context) {
	packages, err := s.db.GetAllPackages(ctx)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var result []api.PackageSummary
	for _, p := range packages {
		result = append(result, packageSummaryToAPI(p))
	}

	if result == nil {
		result = []api.PackageSummary{}
	}
	ctx.JSON(http.StatusOK, result)
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

	p, err := s.db.CreatePackage(ctx, db.CreatePackageParams{
		Name:           input.Name,
		ComposeContent: input.ComposeContent,
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, packageToAPI(p))
}

func (s *Server) GetPackage(ctx *gin.Context, id int) {
	p, err := s.db.GetPackage(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "package not found")
		return
	}

	ctx.JSON(http.StatusOK, packageToAPI(p))
}

func (s *Server) DeletePackage(ctx *gin.Context, id int) {
	err := s.db.DeletePackage(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "package not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func packageToAPI(p db.Package) api.Package {
	return api.Package{
		Id:             toPtr(int(p.ID)),
		Name:           toPtr(p.Name),
		ComposeContent: toPtr(p.ComposeContent),
		CreatedAt:      parseTime(p.CreatedAt),
		UpdatedAt:      parseTime(p.UpdatedAt),
	}
}

func packageSummaryToAPI(p db.GetAllPackagesRow) api.PackageSummary {
	return api.PackageSummary{
		Id:        toPtr(int(p.ID)),
		Name:      toPtr(p.Name),
		CreatedAt: parseTime(p.CreatedAt),
		UpdatedAt: parseTime(p.UpdatedAt),
	}
}
