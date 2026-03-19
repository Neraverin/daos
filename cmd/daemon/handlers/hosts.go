package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/config"
	"github.com/Neraverin/daos/pkg/db"
	"github.com/Neraverin/daos/pkg/image"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Server struct {
	db             *db.Queries
	dbRaw          *sql.DB
	imageProcessor *image.Processor
}

func New(database *sql.DB, cfg *config.Config) *Server {
	imageProcessor := image.NewProcessor(image.Config{
		Registry:         cfg.Docker.Registry,
		ImageLoadTimeout: cfg.Docker.ImageLoadTimeout,
	})
	return &Server{db: db.New(database), dbRaw: database, imageProcessor: imageProcessor}
}

func (s *Server) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, api.HealthStatus{Status: "healthy"})
}

func (s *Server) ListHosts(ctx *gin.Context) {
	hosts, err := s.db.GetAllHosts(ctx)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var result []api.Host
	for _, h := range hosts {
		result = append(result, hostToAPI(h))
	}

	if result == nil {
		result = []api.Host{}
	}
	ctx.JSON(http.StatusOK, result)
}

func (s *Server) CreateHost(ctx *gin.Context) {
	var input api.HostInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	port := int64(22)
	if input.Port != nil {
		port = int64(*input.Port)
	}

	id := uuid.New().String()

	h, err := s.db.CreateHost(ctx, db.CreateHostParams{
		ID:         id,
		Name:       input.Name,
		Hostname:   input.Hostname,
		Port:       port,
		Username:   input.Username,
		SshKeyPath: input.SshKeyPath,
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, hostToAPI(h))
}

func (s *Server) GetHost(ctx *gin.Context, id uuid.UUID) {
	h, err := s.db.GetHost(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "host not found")
		return
	}

	ctx.JSON(http.StatusOK, hostToAPI(h))
}

func (s *Server) UpdateHost(ctx *gin.Context, id uuid.UUID) {
	var input api.HostInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	port := int64(22)
	if input.Port != nil {
		port = int64(*input.Port)
	}

	_, err := s.db.UpdateHost(ctx, db.UpdateHostParams{
		ID:         id.String(),
		Name:       input.Name,
		Hostname:   input.Hostname,
		Port:       port,
		Username:   input.Username,
		SshKeyPath: input.SshKeyPath,
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	s.GetHost(ctx, id)
}

func (s *Server) DeleteHost(ctx *gin.Context, id uuid.UUID) {
	err := s.db.DeleteHost(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "host not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func hostToAPI(h db.Host) api.Host {
	id, _ := uuid.Parse(h.ID)
	return api.Host{
		Id:         &id,
		Name:       &h.Name,
		Hostname:   &h.Hostname,
		Port:       toPtr(int(h.Port)),
		Username:   &h.Username,
		SshKeyPath: &h.SshKeyPath,
		CreatedAt:  parseTime(h.CreatedAt),
		UpdatedAt:  parseTime(h.UpdatedAt),
	}
}
