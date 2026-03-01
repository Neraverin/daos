package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/lipgloss"
)

type Host struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Hostname   string    `json:"hostname"`
	Port       int       `json:"port"`
	Username   string    `json:"username"`
	SSHKeyPath string    `json:"ssh_key_path"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type HostsList struct {
	list       list.Model
	daemonURL  string
	hosts      []Host
	status     string
}

func NewHostsList(daemonURL string) HostsList {
	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Background(lipgloss.Color("235")).Bold(true)
	delegate.Styles.SelectedDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Background(lipgloss.Color("235"))
	delegate.Styles.NormalTitle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	delegate.Styles.NormalDesc = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

	l := list.New(nil, delegate, 0, 0)
	l.Title = "Hosts"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	h := HostsList{
		list:      l,
		daemonURL: daemonURL,
	}

	return h
}

func (m HostsList) Update(msg tea.Msg) (HostsList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			return m, nil
		case "n":
			return m, tea.Cmd(func() tea.Msg { return showHostFormMsg{} })
		case "d":
			if len(m.hosts) > 0 {
				selected := m.list.Index()
				if selected >= 0 && selected < len(m.hosts) {
					return m, tea.Cmd(func() tea.Msg { return deleteHostMsg{hostID: m.hosts[selected].ID} })
				}
			}
		case "r":
			return m, tea.Cmd(func() tea.Msg { return refreshHostsMsg{} })
		}
	case refreshHostsMsg:
		hosts, err := fetchHosts(m.daemonURL)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.hosts = hosts
			m.status = fmt.Sprintf("%d hosts", len(hosts))
			items := make([]list.Item, len(hosts))
			for i, h := range hosts {
				items[i] = hostItem(h)
			}
			m.list.SetItems(items)
		}
	case showHostFormMsg:
		m.status = "Press n to add new host, d to delete, r to refresh"
	case deleteHostMsg:
		err := deleteHost(m.daemonURL, msg.hostID)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.status = "Host deleted"
			return m, tea.Cmd(func() tea.Msg { return refreshHostsMsg{} })
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m HostsList) View() string {
	if m.status != "" {
		m.list.SetStatusBarItemName("host", "hosts")
	}
	return m.list.View()
}

func fetchHosts(url string) ([]Host, error) {
	resp, err := http.Get(url + "/hosts")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var hosts []Host
	if err := json.NewDecoder(resp.Body).Decode(&hosts); err != nil {
		return nil, err
	}

	return hosts, nil
}

func deleteHost(url string, hostID int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/hosts/%d", url, hostID), nil)
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

type hostItem Host

func (i hostItem) Title() string       { return i.Name }
func (i hostItem) Description() string { return fmt.Sprintf("%s@%s:%d", i.Username, i.Hostname, i.Port) }
func (i hostItem) FilterValue() string  { return i.Name }

type refreshHostsMsg struct{}
type showHostFormMsg struct{}
type deleteHostMsg struct{ hostID int }
