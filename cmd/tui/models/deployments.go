package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

type Deployment struct {
	ID             int       `json:"id"`
	HostID         int       `json:"host_id"`
	PackageID      int       `json:"package_id"`
	Status         string    `json:"status"`
	HostName       string    `json:"host_name"`
	HostHostname   string    `json:"host_hostname"`
	PackageName    string    `json:"package_name"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type DeploymentsList struct {
	list        *list.List
	daemonURL    string
	deployments  []Deployment
	status       string
}

func NewDeploymentsList(daemonURL string) DeploymentsList {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Deployments"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	d := DeploymentsList{
		list:      l,
		daemonURL: daemonURL,
	}

	return d
}

func (m DeploymentsList) Update(msg tea.Msg) (DeploymentsList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, tea.Cmd(func() tea.Msg { return showDeploymentFormMsg{} })
		case "r":
			if len(m.deployments) > 0 {
				selected := m.list.Index()
				if selected >= 0 && selected < len(m.deployments) {
					d := m.deployments[selected]
					if d.Status == "pending" || d.Status == "failed" {
						return m, tea.Cmd(func() tea.Msg { return runDeploymentMsg{deploymentID: d.ID} })
					}
				}
			}
			return m, tea.Cmd(func() tea.Msg { return refreshDeploymentsMsg{} })
		case "d":
			if len(m.deployments) > 0 {
				selected := m.list.Index()
				if selected >= 0 && selected < len(m.deployments) {
					return m, tea.Cmd(func() tea.Msg { return deleteDeploymentMsg{deploymentID: m.deployments[selected].ID} })
				}
			}
		case "l":
			if len(m.deployments) > 0 {
				selected := m.list.Index()
				if selected >= 0 && selected < len(m.deployments) {
					return m, tea.Cmd(func() tea.Msg { return showLogsMsg{deploymentID: m.deployments[selected].ID} })
				}
			}
		case "R":
			return m, tea.Cmd(func() tea.Msg { return refreshDeploymentsMsg{} })
		}
	case refreshDeploymentsMsg:
		deps, err := fetchDeployments(m.daemonURL)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.deployments = deps
			m.status = fmt.Sprintf("%d deployments", len(deps))
			items := make([]list.Item, len(deps))
			for i, d := range deps {
				items[i] = deploymentItem(d)
			}
			m.list.SetItems(items)
		}
	case showDeploymentFormMsg:
		m.status = "Press n to create, r to run, d to delete, l for logs, R to refresh"
	case runDeploymentMsg:
		err := runDeployment(m.daemonURL, msg.deploymentID)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.status = "Deployment started"
			return m, tea.Cmd(func() tea.Msg { return refreshDeploymentsMsg{} })
		}
	case deleteDeploymentMsg:
		err := deleteDeployment(m.daemonURL, msg.deploymentID)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.status = "Deployment deleted"
			return m, tea.Cmd(func() tea.Msg { return refreshDeploymentsMsg{} })
		}
	case showLogsMsg:
		logs, err := fetchDeploymentLogs(m.daemonURL, msg.deploymentID)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			for _, l := range logs {
				fmt.Printf("[%s] %s\n", l.Timestamp.Format("15:04:05"), l.Message)
			}
			m.status = fmt.Sprintf("Showing %d log entries", len(logs))
		}
	}

	m.list, _ = m.list.Update(msg)
	return m, nil
}

func (m DeploymentsList) View() string {
	return m.list.View()
}

func fetchDeployments(url string) ([]Deployment, error) {
	resp, err := http.Get(url + "/deployments")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var deployments []Deployment
	if err := json.NewDecoder(resp.Body).Decode(&deployments); err != nil {
		return nil, err
	}

	return deployments, nil
}

func runDeployment(url string, deploymentID int) error {
	resp, err := http.Post(fmt.Sprintf("%s/deployments/%d/run", url, deploymentID), "application/json", nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func deleteDeployment(url string, deploymentID int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/deployments/%d", url, deploymentID), nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

type DeploymentLog struct {
	ID           int       `json:"id"`
	DeploymentID int       `json:"deployment_id"`
	Timestamp    time.Time `json:"timestamp"`
	Message      string    `json:"message"`
}

func fetchDeploymentLogs(url string, deploymentID int) ([]DeploymentLog, error) {
	resp, err := http.Get(fmt.Sprintf("%s/deployments/%d/logs", url, deploymentID))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var logs []DeploymentLog
	if err := json.NewDecoder(resp.Body).Decode(&logs); err != nil {
		return nil, err
	}

	return logs, nil
}

type deploymentItem Deployment

func (i deploymentItem) Title() string       { return fmt.Sprintf("%s -> %s", i.PackageName, i.HostName) }
func (i deploymentItem) Description() string { return fmt.Sprintf("Status: %s | %s", i.Status, i.UpdatedAt.Format("2006-01-02 15:04")) }
func (i deploymentItem) FilterValue() string  { return i.PackageName + " " + i.HostName }

type refreshDeploymentsMsg struct{}
type showDeploymentFormMsg struct{}
type runDeploymentMsg struct{ deploymentID int }
type deleteDeploymentMsg struct{ deploymentID int }
type showLogsMsg struct{ deploymentID int }
