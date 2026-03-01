package handlers

import (
	"net/http"
	"time"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/ansible"
	"github.com/Neraverin/daos/pkg/models"
	"github.com/gin-gonic/gin"
)

func (s *Server) ListDeployments(ctx *gin.Context) {
	rows, err := s.db.Query(`
		SELECT
			d.id,
			d.host_id,
			d.package_id,
			d.status,
			d.created_at,
			d.updated_at,
			h.name as host_name,
			h.hostname as host_hostname,
			p.name as package_name
		FROM deployments d
		JOIN hosts h ON d.host_id = h.id
		JOIN packages p ON d.package_id = p.id
		ORDER BY d.created_at DESC
	`)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var deployments []api.Deployment
	for rows.Next() {
		var d api.Deployment
		if err := rows.Scan(&d.ID, &d.HostID, &d.PackageID, &d.Status, &d.CreatedAt, &d.UpdatedAt, &d.HostName, &d.HostHostname, &d.PackageName); err != nil {
			api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		deployments = append(deployments, d)
	}

	if deployments == nil {
		deployments = []api.Deployment{}
	}
	ctx.JSON(http.StatusOK, deployments)
}

func (s *Server) CreateDeployment(ctx *gin.Context) {
	var input api.DeploymentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	var hostCount int
	err := s.db.QueryRow("SELECT COUNT(*) FROM hosts WHERE id = ?", input.HostID).Scan(&hostCount)
	if err != nil || hostCount == 0 {
		api.ErrorJSON(ctx, http.StatusBadRequest, "host not found")
		return
	}

	var packageCount int
	err = s.db.QueryRow("SELECT COUNT(*) FROM packages WHERE id = ?", input.PackageID).Scan(&packageCount)
	if err != nil || packageCount == 0 {
		api.ErrorJSON(ctx, http.StatusBadRequest, "package not found")
		return
	}

	result, err := s.db.Exec(
		"INSERT INTO deployments (host_id, package_id, status) VALUES (?, ?, 'pending')",
		input.HostID, input.PackageID,
	)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	id, _ := result.LastInsertId()
	s.GetDeployment(ctx, int(id))
}

func (s *Server) GetDeployment(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid deployment id")
		return
	}

	var d api.Deployment
	err := s.db.QueryRow(`
		SELECT
			d.id,
			d.host_id,
			d.package_id,
			d.status,
			d.created_at,
			d.updated_at,
			h.name as host_name,
			h.hostname as host_hostname,
			p.name as package_name
		FROM deployments d
		JOIN hosts h ON d.host_id = h.id
		JOIN packages p ON d.package_id = p.id
		WHERE d.id = ?
	`, id).Scan(&d.ID, &d.HostID, &d.PackageID, &d.Status, &d.CreatedAt, &d.UpdatedAt, &d.HostName, &d.HostHostname, &d.PackageName)

	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	ctx.JSON(http.StatusOK, d)
}

func (s *Server) DeleteDeployment(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid deployment id")
		return
	}

	result, err := s.db.Exec("DELETE FROM deployments WHERE id = ?", id)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	rows, _ := result.RowsAffected()
	if rows == 0 {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (s *Server) RunDeployment(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid deployment id")
		return
	}

	var status string
	err := s.db.QueryRow("SELECT status FROM deployments WHERE id = ?", id).Scan(&status)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	if status == models.DeploymentStatusRunning {
		api.ErrorJSON(ctx, http.StatusBadRequest, "deployment already in progress")
		return
	}

	s.db.Exec("UPDATE deployments SET status = 'running', updated_at = ? WHERE id = ?", time.Now(), id)

	go s.runAnsibleDeployment(id)

	s.GetDeployment(ctx)
}

func (s *Server) runAnsibleDeployment(deploymentID int) {
	var hostID, packageID int64
	var hostname, username, sshKeyPath, composeContent string

	err := s.db.QueryRow(`
		SELECT d.host_id, d.package_id, h.hostname, h.username, h.ssh_key_path, p.compose_content
		FROM deployments d
		JOIN hosts h ON d.host_id = h.id
		JOIN packages p ON d.package_id = p.id
		WHERE d.id = ?
	`, deploymentID).Scan(&hostID, &packageID, &hostname, &username, &sshKeyPath, &composeContent)

	if err != nil {
		s.logMessage(deploymentID, "Failed to get deployment details: "+err.Error())
		s.updateStatus(deploymentID, models.DeploymentStatusFailed)
		return
	}

	s.logMessage(deploymentID, "Starting deployment to "+hostname)

	executor := ansible.NewExecutor(hostname, int(hostID), username, sshKeyPath, composeContent)
	err = executor.Run(func(line string) {
		s.logMessage(deploymentID, line)
	})

	if err != nil {
		s.logMessage(deploymentID, "Deployment failed: "+err.Error())
		s.updateStatus(deploymentID, models.DeploymentStatusFailed)
		return
	}

	s.logMessage(deploymentID, "Deployment completed successfully")
	s.updateStatus(deploymentID, models.DeploymentStatusSuccess)
}

func (s *Server) logMessage(deploymentID int, message string) {
	s.db.Exec("INSERT INTO logs (deployment_id, message) VALUES (?, ?)", deploymentID, message)
}

func (s *Server) updateStatus(deploymentID int, status string) {
	s.db.Exec("UPDATE deployments SET status = ?, updated_at = ? WHERE id = ?", status, time.Now(), deploymentID)
}

func (s *Server) GetDeploymentLogs(ctx *gin.Context) {
	id, ok := api.GetIDParam(ctx)
	if !ok {
		api.ErrorJSON(ctx, http.StatusBadRequest, "invalid deployment id")
		return
	}

	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM deployments WHERE id = ?", id).Scan(&count)
	if err != nil || count == 0 {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	rows, err := s.db.Query("SELECT id, deployment_id, timestamp, message FROM logs WHERE deployment_id = ? ORDER BY timestamp", id)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}
	defer rows.Close()

	var logs []api.Log
	for rows.Next() {
		var l api.Log
		if err := rows.Scan(&l.ID, &l.DeploymentID, &l.Timestamp, &l.Message); err != nil {
			api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
			return
		}
		logs = append(logs, l)
	}

	if logs == nil {
		logs = []api.Log{}
	}
	ctx.JSON(http.StatusOK, logs)
}
