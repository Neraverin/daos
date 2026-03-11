package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
)

type RoleSummary struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type RolesList struct {
	list      list.Model
	daemonURL string
	roles     []RoleSummary
	status    string
}

func NewRolesList(daemonURL string) RolesList {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Background(lipgloss.Color("235")).Bold(true)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Background(lipgloss.Color("235"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	l := list.New(nil, delegate, 0, 0)
	l.Title = "Roles"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	r := RolesList{
		list:      l,
		daemonURL: daemonURL,
	}

	return r
}

func (m RolesList) Update(msg tea.Msg) (RolesList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, tea.Cmd(func() tea.Msg { return showRoleFormMsg{} })
		case "d":
			if len(m.roles) > 0 {
				selected := m.list.Index()
				if selected >= 0 && selected < len(m.roles) {
					return m, tea.Cmd(func() tea.Msg { return deleteRoleMsg{roleID: m.roles[selected].ID} })
				}
			}
		case "r":
			return m, tea.Cmd(func() tea.Msg { return refreshRolesMsg{} })
		}
	case refreshRolesMsg:
		roles, err := fetchRoles(m.daemonURL)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.roles = roles
			m.status = fmt.Sprintf("%d roles", len(roles))
			items := make([]list.Item, len(roles))
			for i, r := range roles {
				items[i] = roleItem(r)
			}
			m.list.SetItems(items)
		}
	case showRoleFormMsg:
		m.status = "Press n to upload role, d to delete, r to refresh"
	case deleteRoleMsg:
		err := deleteRole(m.daemonURL, msg.roleID)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.status = "Role deleted"
			return m, tea.Cmd(func() tea.Msg { return refreshRolesMsg{} })
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m RolesList) View() string {
	return m.list.View()
}

func fetchRoles(url string) ([]RoleSummary, error) {
	resp, err := http.Get(url + "/roles")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var roles []RoleSummary
	if err := json.NewDecoder(resp.Body).Decode(&roles); err != nil {
		return nil, err
	}

	return roles, nil
}

func deleteRole(url string, roleID uuid.UUID) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/roles/%s", url, roleID.String()), nil)
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

type roleItem RoleSummary

func (i roleItem) Title() string { return i.Name }
func (i roleItem) Description() string {
	return fmt.Sprintf("Created: %s", i.CreatedAt.Format("2006-01-02 15:04"))
}
func (i roleItem) FilterValue() string { return i.Name }

type refreshRolesMsg struct{}
type showRoleFormMsg struct{}
type deleteRoleMsg struct{ roleID uuid.UUID }
