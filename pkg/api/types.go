package api

import (
	"time"
)

type Host struct {
	ID         int       `json:"id" yaml:"id"`
	Name       string    `json:"name" yaml:"name"`
	Hostname   string    `json:"hostname" yaml:"hostname"`
	Port       int       `json:"port" yaml:"port"`
	Username   string    `json:"username" yaml:"username"`
	SSHKeyPath string    `json:"ssh_key_path" yaml:"ssh_key_path"`
	CreatedAt  time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" yaml:"updated_at"`
}

type HostInput struct {
	Name       string `json:"name" yaml:"name"`
	Hostname   string `json:"hostname" yaml:"hostname"`
	Port       *int   `json:"port" yaml:"port"`
	Username   string `json:"username" yaml:"username"`
	SSHKeyPath string `json:"ssh_key_path" yaml:"ssh_key_path"`
}

type PackageSummary struct {
	ID        int       `json:"id" yaml:"id"`
	Name      string    `json:"name" yaml:"name"`
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

type Package struct {
	ID             int       `json:"id" yaml:"id"`
	Name           string    `json:"name" yaml:"name"`
	ComposeContent string    `json:"compose_content" yaml:"compose_content"`
	CreatedAt      time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" yaml:"updated_at"`
}

type PackageInput struct {
	Name           string `json:"name" yaml:"name"`
	ComposeContent string `json:"compose_content" yaml:"compose_content"`
}

type Deployment struct {
	ID             int       `json:"id" yaml:"id"`
	HostID         int       `json:"host_id" yaml:"host_id"`
	PackageID      int       `json:"package_id" yaml:"package_id"`
	Status         string    `json:"status" yaml:"status"`
	HostName       string    `json:"host_name" yaml:"host_name"`
	HostHostname   string    `json:"host_hostname" yaml:"host_hostname"`
	PackageName    string    `json:"package_name" yaml:"package_name"`
	CreatedAt      time.Time `json:"created_at" yaml:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" yaml:"updated_at"`
}

type DeploymentInput struct {
	HostID    int `json:"host_id" yaml:"host_id"`
	PackageID int `json:"package_id" yaml:"package_id"`
}

type Log struct {
	ID           int       `json:"id" yaml:"id"`
	DeploymentID int       `json:"deployment_id" yaml:"deployment_id"`
	Timestamp    time.Time `json:"timestamp" yaml:"timestamp"`
	Message      string    `json:"message" yaml:"message"`
}

type HealthStatus struct {
	Status string `json:"status" yaml:"status"`
}
