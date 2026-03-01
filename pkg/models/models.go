package models

import "time"

type Host struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	Hostname   string    `json:"hostname"`
	Port       int       `json:"port"`
	Username   string    `json:"username"`
	SSHKeyPath string    `json:"ssh_key_path"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type Package struct {
	ID             int64     `json:"id"`
	Name           string    `json:"name"`
	ComposeContent string    `json:"compose_content,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Deployment struct {
	ID          int64     `json:"id"`
	HostID      int64     `json:"host_id"`
	PackageID   int64     `json:"package_id"`
	Status      string    `json:"status"`
	HostName    string    `json:"host_name,omitempty"`
	HostHostname string   `json:"host_hostname,omitempty"`
	PackageName string    `json:"package_name,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Log struct {
	ID           int64     `json:"id"`
	DeploymentID int64     `json:"deployment_id"`
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message"`
}

const (
	DeploymentStatusPending = "pending"
	DeploymentStatusRunning = "running"
	DeploymentStatusSuccess = "success"
	DeploymentStatusFailed  = "failed"
)
