package main

import (
	"fmt"
	"os"

	"github.com/Neraverin/daos/cmd/tui/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("170")).
			Background(lipgloss.Color("57"))

	subtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type model struct {
	menu    models.Menu
	hosts   models.HostsList
	packages models.PackagesList
	deployments models.DeploymentsList
	current    string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "esc":
			m.current = "menu"
			return m, nil
		}
	}

	switch m.current {
	case "menu":
		m.menu, cmd = m.menu.Update(msg)
	case "hosts":
		m.hosts, cmd = m.hosts.Update(msg)
	case "packages":
		m.packages, cmd = m.packages.Update(msg)
	case "deployments":
		m.deployments, cmd = m.deployments.Update(msg)
	}

	return m, cmd
}

func (m model) View() string {
	var s string

	switch m.current {
	case "menu":
		s = m.menu.View()
	case "hosts":
		s = m.hosts.View()
	case "packages":
		s = m.packages.View()
	case "deployments":
		s = m.deployments.View()
	}

	return fmt.Sprintf("%s\n\n%s", titleStyle.Render("DAOS - Deployment and Orchestration Service"), s)
}

func main() {
	daemonURL := os.Getenv("DAOS_URL")
	if daemonURL == "" {
		daemonURL = "http://localhost:8080/api/v1"
	}

	p := tea.NewProgram(model{
		menu:       models.NewMenu(),
		hosts:      models.NewHostsList(daemonURL),
		packages:   models.NewPackagesList(daemonURL),
		deployments: models.NewDeploymentsList(daemonURL),
	})

	if err := p.Start(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
