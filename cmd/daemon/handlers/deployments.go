package handlers

import (
	"context"
	"net/http"

	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/ansible"
	"github.com/Neraverin/daos/pkg/db"
	"github.com/gin-gonic/gin"
)

func (s *Server) ListDeployments(ctx *gin.Context) {
	deployments, err := s.db.GetAllDeployments(ctx)
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var result []api.Deployment
	for _, d := range deployments {
		result = append(result, deploymentRowToAPI(d))
	}

	if result == nil {
		result = []api.Deployment{}
	}
	ctx.JSON(http.StatusOK, result)
}

func (s *Server) CreateDeployment(ctx *gin.Context) {
	var input api.DeploymentInput
	if err := ctx.ShouldBindJSON(&input); err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	_, err := s.db.GetHost(ctx, int64(input.HostId))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, "host not found")
		return
	}

	_, err = s.db.GetPackage(ctx, int64(input.PackageId))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, "package not found")
		return
	}

	d, err := s.db.CreateDeployment(ctx, db.CreateDeploymentParams{
		HostID:    int64(input.HostId),
		PackageID: int64(input.PackageId),
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, deploymentToAPI(d))
}

func (s *Server) GetDeployment(ctx *gin.Context, id int) {
	d, err := s.db.GetDeployment(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	ctx.JSON(http.StatusOK, deploymentGetRowToAPI(d))
}

func (s *Server) DeleteDeployment(ctx *gin.Context, id int) {
	err := s.db.DeleteDeployment(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (s *Server) RunDeployment(ctx *gin.Context, id int) {
	d, err := s.db.GetDeployment(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	if d.Status == "running" {
		api.ErrorJSON(ctx, http.StatusBadRequest, "deployment already in progress")
		return
	}

	_, err = s.db.UpdateDeploymentStatus(ctx, db.UpdateDeploymentStatusParams{
		ID:     int64(id),
		Status: "running",
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go s.runAnsibleDeployment(int64(id))

	s.GetDeployment(ctx, id)
}

type deploymentDetails struct {
	hostID        int64
	packageID     int64
	hostname      string
	username      string
	sshKeyPath    string
	composeContent string
}

func (s *Server) getDeploymentDetails(deploymentID int64) (*deploymentDetails, error) {
	rows, err := s.dbRaw.QueryContext(context.Background(), `
		SELECT d.host_id, d.package_id, h.hostname, h.username, h.ssh_key_path, p.compose_content
		FROM deployments d
		JOIN hosts h ON d.host_id = h.id
		JOIN packages p ON d.package_id = p.id
		WHERE d.id = ?
	`, deploymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var details deploymentDetails
	for rows.Next() {
		if err := rows.Scan(&details.hostID, &details.packageID, &details.hostname, &details.username, &details.sshKeyPath, &details.composeContent); err != nil {
			return nil, err
		}
	}
	return &details, rows.Err()
}

func (s *Server) runAnsibleDeployment(deploymentID int64) {
	details, err := s.getDeploymentDetails(deploymentID)
	if err != nil {
		s.logMessage(deploymentID, "Failed to get deployment details: "+err.Error())
		s.updateStatus(deploymentID, "failed")
		return
	}

	s.logMessage(deploymentID, "Starting deployment to "+details.hostname)

	executor := ansible.NewExecutor(details.hostname, int(details.hostID), details.username, details.sshKeyPath, details.composeContent)
	err = executor.Run(func(line string) {
		s.logMessage(deploymentID, line)
	})

	if err != nil {
		s.logMessage(deploymentID, "Deployment failed: "+err.Error())
		s.updateStatus(deploymentID, "failed")
		return
	}

	s.logMessage(deploymentID, "Deployment completed successfully")
	s.updateStatus(deploymentID, "success")
}

func (s *Server) logMessage(deploymentID int64, message string) {
	s.db.CreateLog(context.Background(), db.CreateLogParams{
		DeploymentID: deploymentID,
		Message:       message,
	})
}

func (s *Server) updateStatus(deploymentID int64, status string) {
	s.db.UpdateDeploymentStatus(context.Background(), db.UpdateDeploymentStatusParams{
		ID:     deploymentID,
		Status: status,
	})
}

func (s *Server) GetDeploymentLogs(ctx *gin.Context, id int) {
	_, err := s.db.GetDeployment(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	logs, err := s.db.GetLogsByDeployment(ctx, int64(id))
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	var result []api.Log
	for _, l := range logs {
		result = append(result, logToAPI(l))
	}

	if result == nil {
		result = []api.Log{}
	}
	ctx.JSON(http.StatusOK, result)
}

func deploymentToAPI(d db.Deployment) api.Deployment {
	return api.Deployment{
		Id:        toPtr(int(d.ID)),
		HostId:    toPtr(int(d.HostID)),
		PackageId: toPtr(int(d.PackageID)),
		Status:    toPtr(api.DeploymentStatus(d.Status)),
		CreatedAt: parseTime(d.CreatedAt),
		UpdatedAt: parseTime(d.UpdatedAt),
	}
}

func deploymentRowToAPI(d db.GetAllDeploymentsRow) api.Deployment {
	return api.Deployment{
		Id:           toPtr(int(d.ID)),
		HostId:       toPtr(int(d.HostID)),
		PackageId:    toPtr(int(d.PackageID)),
		Status:       toPtr(api.DeploymentStatus(d.Status)),
		CreatedAt:    parseTime(d.CreatedAt),
		UpdatedAt:    parseTime(d.UpdatedAt),
		HostName:     toPtr(d.HostName),
		HostHostname: toPtr(d.HostHostname),
		PackageName:  toPtr(d.PackageName),
	}
}

func deploymentGetRowToAPI(d db.GetDeploymentRow) api.Deployment {
	return api.Deployment{
		Id:           toPtr(int(d.ID)),
		HostId:       toPtr(int(d.HostID)),
		PackageId:    toPtr(int(d.PackageID)),
		Status:       toPtr(api.DeploymentStatus(d.Status)),
		CreatedAt:    parseTime(d.CreatedAt),
		UpdatedAt:    parseTime(d.UpdatedAt),
		HostName:     toPtr(d.HostName),
		HostHostname: toPtr(d.HostHostname),
		PackageName:  toPtr(d.PackageName),
	}
}

func logToAPI(l db.Log) api.Log {
	return api.Log{
		Id:           toPtr(int(l.ID)),
		DeploymentId: toPtr(int(l.DeploymentID)),
		Timestamp:    parseTime(l.Timestamp),
		Message:      toPtr(l.Message),
	}
}
