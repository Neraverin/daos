package handlers

import (
	"context"
	"net/http"

	"github.com/Neraverin/daos/pkg/ansible"
	"github.com/Neraverin/daos/pkg/api"
	"github.com/Neraverin/daos/pkg/db"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	_, err := s.db.GetHost(ctx, input.HostId.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, "host not found")
		return
	}

	_, err = s.db.GetRole(ctx, input.RoleId.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, "role not found")
		return
	}

	id := uuid.New().String()

	d, err := s.db.CreateDeployment(ctx, db.CreateDeploymentParams{
		ID:     id,
		HostID: input.HostId.String(),
		RoleID: input.RoleId.String(),
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusBadRequest, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, deploymentToAPI(d))
}

func (s *Server) GetDeployment(ctx *gin.Context, id uuid.UUID) {
	d, err := s.db.GetDeployment(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	ctx.JSON(http.StatusOK, deploymentGetRowToAPI(d))
}

func (s *Server) DeleteDeployment(ctx *gin.Context, id uuid.UUID) {
	err := s.db.DeleteDeployment(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (s *Server) RunDeployment(ctx *gin.Context, id uuid.UUID) {
	d, err := s.db.GetDeployment(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	if d.Status == "running" {
		api.ErrorJSON(ctx, http.StatusBadRequest, "deployment already in progress")
		return
	}

	_, err = s.db.UpdateDeploymentStatus(ctx, db.UpdateDeploymentStatusParams{
		ID:     id.String(),
		Status: "running",
	})
	if err != nil {
		api.ErrorJSON(ctx, http.StatusInternalServerError, err.Error())
		return
	}

	go s.runAnsibleDeployment(id.String())

	s.GetDeployment(ctx, id)
}

type deploymentDetails struct {
	hostID         string
	roleID         string
	hostname       string
	username       string
	sshKeyPath     string
	composeContent string
}

func (s *Server) getDeploymentDetails(ctx context.Context, deploymentID string) (*deploymentDetails, error) {
	details, err := s.db.GetDeploymentDetails(ctx, deploymentID)
	if err != nil {
		return nil, err
	}

	return &deploymentDetails{
		hostID:         details.HostID,
		roleID:         details.RoleID,
		hostname:       details.Hostname,
		username:       details.Username,
		sshKeyPath:     details.SshKeyPath,
		composeContent: details.ComposeContent,
	}, nil
}

func (s *Server) runAnsibleDeployment(deploymentID string) {
	details, err := s.getDeploymentDetails(context.Background(), deploymentID)
	if err != nil {
		s.logMessage(deploymentID, "Failed to get deployment details: "+err.Error())
		s.updateStatus(deploymentID, "failed")
		return
	}

	s.logMessage(deploymentID, "Starting deployment to "+details.hostname)

	executor := ansible.NewExecutor(details.hostname, 22, details.username, details.sshKeyPath, details.composeContent)
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

func (s *Server) logMessage(deploymentID string, message string) {
	id := uuid.New().String()
	s.db.CreateLog(context.Background(), db.CreateLogParams{
		ID:           id,
		DeploymentID: deploymentID,
		Message:      message,
	})
}

func (s *Server) updateStatus(deploymentID string, status string) {
	s.db.UpdateDeploymentStatus(context.Background(), db.UpdateDeploymentStatusParams{
		ID:     deploymentID,
		Status: status,
	})
}

func (s *Server) GetDeploymentLogs(ctx *gin.Context, id uuid.UUID) {
	_, err := s.db.GetDeployment(ctx, id.String())
	if err != nil {
		api.ErrorJSON(ctx, http.StatusNotFound, "deployment not found")
		return
	}

	logs, err := s.db.GetLogsByDeployment(ctx, id.String())
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
	parsedID, _ := uuid.Parse(d.ID)
	parsedHostID, _ := uuid.Parse(d.HostID)
	parsedRoleID, _ := uuid.Parse(d.RoleID)
	return api.Deployment{
		Id:        &parsedID,
		HostId:    &parsedHostID,
		RoleId:    &parsedRoleID,
		Status:    toPtr(api.DeploymentStatus(d.Status)),
		CreatedAt: parseTime(d.CreatedAt),
		UpdatedAt: parseTime(d.UpdatedAt),
	}
}

func deploymentRowToAPI(d db.GetAllDeploymentsRow) api.Deployment {
	parsedID, _ := uuid.Parse(d.ID)
	parsedHostID, _ := uuid.Parse(d.HostID)
	parsedRoleID, _ := uuid.Parse(d.RoleID)
	return api.Deployment{
		Id:           &parsedID,
		HostId:       &parsedHostID,
		RoleId:       &parsedRoleID,
		Status:       toPtr(api.DeploymentStatus(d.Status)),
		CreatedAt:    parseTime(d.CreatedAt),
		UpdatedAt:    parseTime(d.UpdatedAt),
		HostName:     toPtr(d.HostName),
		HostHostname: toPtr(d.HostHostname),
		RoleName:     toPtr(d.RoleName),
	}
}

func deploymentGetRowToAPI(d db.GetDeploymentRow) api.Deployment {
	parsedID, _ := uuid.Parse(d.ID)
	parsedHostID, _ := uuid.Parse(d.HostID)
	parsedRoleID, _ := uuid.Parse(d.RoleID)
	return api.Deployment{
		Id:           &parsedID,
		HostId:       &parsedHostID,
		RoleId:       &parsedRoleID,
		Status:       toPtr(api.DeploymentStatus(d.Status)),
		CreatedAt:    parseTime(d.CreatedAt),
		UpdatedAt:    parseTime(d.UpdatedAt),
		HostName:     toPtr(d.HostName),
		HostHostname: toPtr(d.HostHostname),
		RoleName:     toPtr(d.RoleName),
	}
}

func logToAPI(l db.Log) api.Log {
	parsedID, _ := uuid.Parse(l.ID)
	parsedDeploymentID, _ := uuid.Parse(l.DeploymentID)
	return api.Log{
		Id:           &parsedID,
		DeploymentId: &parsedDeploymentID,
		Timestamp:    parseTime(l.Timestamp),
		Message:      toPtr(l.Message),
	}
}
