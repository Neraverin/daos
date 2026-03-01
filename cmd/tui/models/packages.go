package models

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

type PackageSummary struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PackagesList struct {
	list      list.Model
	daemonURL string
	packages  []PackageSummary
	status    string
}

func NewPackagesList(daemonURL string) PackagesList {
	l := list.New(nil, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Packages"
	l.SetShowStatusBar(true)
	l.SetFilteringEnabled(true)

	p := PackagesList{
		list:      l,
		daemonURL: daemonURL,
	}

	return p
}

func (m PackagesList) Update(msg tea.Msg) (PackagesList, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "n":
			return m, tea.Cmd(func() tea.Msg { return showPackageFormMsg{} })
		case "d":
			if len(m.packages) > 0 {
				selected := m.list.Index()
				if selected >= 0 && selected < len(m.packages) {
					return m, tea.Cmd(func() tea.Msg { return deletePackageMsg{packageID: m.packages[selected].ID} })
				}
			}
		case "r":
			return m, tea.Cmd(func() tea.Msg { return refreshPackagesMsg{} })
		}
	case refreshPackagesMsg:
		pkgs, err := fetchPackages(m.daemonURL)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.packages = pkgs
			m.status = fmt.Sprintf("%d packages", len(pkgs))
			items := make([]list.Item, len(pkgs))
			for i, p := range pkgs {
				items[i] = packageItem(p)
			}
			m.list.SetItems(items)
		}
	case showPackageFormMsg:
		m.status = "Press n to upload package, d to delete, r to refresh"
	case deletePackageMsg:
		err := deletePackage(m.daemonURL, msg.packageID)
		if err != nil {
			m.status = fmt.Sprintf("Error: %v", err)
		} else {
			m.status = "Package deleted"
			return m, tea.Cmd(func() tea.Msg { return refreshPackagesMsg{} })
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m PackagesList) View() string {
	return m.list.View()
}

func fetchPackages(url string) ([]PackageSummary, error) {
	resp, err := http.Get(url + "/packages")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var packages []PackageSummary
	if err := json.NewDecoder(resp.Body).Decode(&packages); err != nil {
		return nil, err
	}

	return packages, nil
}

func deletePackage(url string, packageID int) error {
	req, err := http.NewRequest("DELETE", fmt.Sprintf("%s/packages/%d", url, packageID), nil)
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

type packageItem PackageSummary

func (i packageItem) Title() string       { return i.Name }
func (i packageItem) Description() string { return fmt.Sprintf("Created: %s", i.CreatedAt.Format("2006-01-02 15:04")) }
func (i packageItem) FilterValue() string { return i.Name }

type refreshPackagesMsg struct{}
type showPackageFormMsg struct{}
type deletePackageMsg struct{ packageID int }
