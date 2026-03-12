package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
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

type Role struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	RolePath  string    `json:"role_path"`
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
			return m, tea.Cmd(func() tea.Msg { return ShowRoleFormMsg{} })
		case "d":
			if len(m.roles) > 0 {
				selected := m.list.Index()
				if selected >= 0 && selected < len(m.roles) {
					return m, tea.Cmd(func() tea.Msg { return DeleteRoleMsg{roleID: m.roles[selected].ID} })
				}
			}
		case "r":
			return m, tea.Cmd(func() tea.Msg { return RefreshRolesMsg{} })
		}
	case RefreshRolesMsg:
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
	case ShowRoleFormMsg:
		m.status = "Press n to add new role, d to delete, r to refresh"
	case DeleteRoleMsg:
		err := deleteRole(m.daemonURL, msg.roleID)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.status = "Role deleted"
			return m, tea.Cmd(func() tea.Msg { return RefreshRolesMsg{} })
		}
	case RoleFormModel:
		if msg.err != nil {
			m.status = fmt.Sprintf("Error: %v", msg.err)
		} else if msg.saved {
			m.status = "Role created successfully"
			return m, tea.Cmd(func() tea.Msg { return RefreshRolesMsg{} })
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

type RefreshRolesMsg struct{}
type ShowRoleFormMsg struct{}
type DeleteRoleMsg struct{ roleID uuid.UUID }

type RoleFormModel struct {
	nameInput textinput.Model
	pathInput textinput.Model
	daemonURL string
	focused   int
	err       error
	saved     bool
}

func NewRoleFormModel(daemonURL string) RoleFormModel {
	nameInput := textinput.New()
	nameInput.Placeholder = "Role name"
	nameInput.Focus()

	pathInput := textinput.New()
	pathInput.Placeholder = "/absolute/path/to/role/folder"

	return RoleFormModel{
		nameInput: nameInput,
		pathInput: pathInput,
		daemonURL: daemonURL,
		focused:   0,
	}
}

func (m RoleFormModel) DaemonURL() string {
	return m.daemonURL
}

func (m RoleFormModel) Init() tea.Cmd {
	return nil
}

func (m RoleFormModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "tab":
			m.focused = (m.focused + 1) % 2
			if m.focused == 0 {
				m.nameInput.Focus()
				m.pathInput.Blur()
			} else {
				m.nameInput.Blur()
				m.pathInput.Focus()
			}
			return m, nil
		case "enter":
			if m.focused == 0 {
				m.focused = 1
				m.nameInput.Blur()
				m.pathInput.Focus()
				return m, nil
			}
			return m, m.saveRole()
		case "esc":
			return m, tea.Cmd(func() tea.Msg { return RefreshRolesMsg{} })
		}
	case saveRoleMsg:
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.saved = true
		}
		return m, nil
	}

	var cmd tea.Cmd
	m.nameInput, cmd = m.nameInput.Update(msg)
	m.pathInput, _ = m.pathInput.Update(msg)
	return m, cmd
}

func (m RoleFormModel) View() string {
	return fmt.Sprintf(
		"Create New Role\n\n%s\n%s\n\n%s\n%s\n\n[Tab] Switch fields | [Enter] Save | [Esc] Cancel",
		m.nameInput.View(),
		m.pathInput.View(),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Name:"),
		lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Role folder path (absolute):"),
	)
}

func (m RoleFormModel) saveRole() tea.Cmd {
	return func() tea.Msg {
		role := Role{
			Name:     m.nameInput.Value(),
			RolePath: m.pathInput.Value(),
		}

		jsonData, err := json.Marshal(role)
		if err != nil {
			return saveRoleMsg{err: err}
		}

		resp, err := http.Post(m.daemonURL+"/roles", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return saveRoleMsg{err: err}
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			return saveRoleMsg{err: fmt.Errorf("API returned status %d", resp.StatusCode)}
		}

		return saveRoleMsg{}
	}
}

type saveRoleMsg struct {
	err   error
	saved bool
}
