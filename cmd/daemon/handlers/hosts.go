package handlers

import (
	"database/sql"
	"net/http"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/db"
	"github.com/gin-gonic/gin"
)

type Server struct {
	db     *db.Queries
	dbRaw  *sql.DB
}

func New(database *sql.DB) *Server {
	return &Server{db: db.New(database), dbRaw: database}
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

	h, err := s.db.CreateHost(ctx, db.CreateHostParams{
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

func (s *Server) GetHost(ctx *gin.Context, id int) {
	h, err := s.db.GetHost(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "host not found")
		return
	}

	ctx.JSON(http.StatusOK, hostToAPI(h))
}

func (s *Server) UpdateHost(ctx *gin.Context, id int) {
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
		ID:         int64(id),
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

func (s *Server) DeleteHost(ctx *gin.Context, id int) {
	err := s.db.DeleteHost(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "host not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func hostToAPI(h db.Host) api.Host {
	return api.Host{
		Id:         toPtr(int(h.ID)),
		Name:       toPtr(h.Name),
		Hostname:   toPtr(h.Hostname),
		Port:       toPtr(int(h.Port)),
		Username:   toPtr(h.Username),
		SshKeyPath: toPtr(h.SshKeyPath),
		CreatedAt:  parseTime(h.CreatedAt),
		UpdatedAt:  parseTime(h.UpdatedAt),
	}
}
