package handlers

import (
	"net/http"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/db"
	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

type Server struct {
	db *db.DB
}

func New(database *db.DB) *Server {
	return &Server{db: database}
}

func (s *Server) HealthCheck(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, api.HealthStatus{Status: "healthy"})
}

func (s *Server) ListHosts(ctx *gin.Context) {
	rows, err := s.db.Query("SELECT id, name, hostname, port, username, ssh_key_path, created_at, updated_at FROM hosts ORDER BY name")
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var hosts []api.Host
	for rows.Next() {
		var h api.Host
		if err := rows.Scan(&h.ID, &h.Name, &h.Hostname, &h.Port, &h.Username, &h.SSHKeyPath, &h.CreatedAt, &h.UpdatedAt); err != nil {
			api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		hosts = append(hosts, h)
	}

	if hosts == nil {
		hosts = []api.Host{}
	}
	ctx.JSON(http.StatusOK, hosts)
}

func (s *Server) CreateHost(ctx *gin.Context) {
	var input api.HostInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	port := 22
	if input.Port != nil {
		port = *input.Port
	}

	result, err := s.db.Exec(
		"INSERT INTO hosts (name, hostname, port, username, ssh_key_path) VALUES (?, ?, ?, ?, ?)",
		input.Name, input.Hostname, port, input.Username, input.SSHKeyPath,
	)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	s.GetHost(ctx, int(id))
}

func (s *Server) GetHost(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid host id")
		return
	}

	var h api.Host
	err := s.db.QueryRow(
		"SELECT id, name, hostname, port, username, ssh_key_path, created_at, updated_at FROM hosts WHERE id = ?",
		id,
	).Scan(&h.ID, &h.Name, &h.Hostname, &h.Port, &h.Username, &h.SSHKeyPath, &h.CreatedAt, &h.UpdatedAt)

	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "host not found")
		return
	}

	ctx.JSON(http.StatusOK, h)
}

func (s *Server) UpdateHost(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid host id")
		return
	}

	var input api.HostInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	port := 22
	if input.Port != nil {
		port = *input.Port
	}

	_, err := s.db.Exec(
		"UPDATE hosts SET name = ?, hostname = ?, port = ?, username = ?, ssh_key_path = ?, updated_at = datetime('now') WHERE id = ?",
		input.Name, input.Hostname, port, input.Username, input.SSHKeyPath, id,
	)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	s.GetHost(ctx)
}

func (s *Server) DeleteHost(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid host id")
		return
	}

	result, err := s.db.Exec("DELETE FROM hosts WHERE id = ?", id)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		api.ErrorJSON(ctx, http.StatusNotFound, "host not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}
