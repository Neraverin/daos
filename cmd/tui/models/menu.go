package models

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/list"
)

type Menu struct {
	list *list.List
}

func NewMenu() Menu {
	items := []list.Item{
		menuItem{title: "Hosts", desc: "Manage remote hosts"},
		menuItem{title: "Packages", desc: "Manage Docker Compose packages"},
		menuItem{title: "Deployments", desc: "Manage deployments"},
		menuItem{title: "Quit", desc: "Exit application"},
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Main Menu"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)

	return Menu{list: l}
}

func (m Menu) Update(msg tea.Msg) (Menu, tea.Cmd) {
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Menu) View() string {
	return m.list.View()
}

func (m Menu) SelectedItem() (string, int) {
	if i, ok := m.list.SelectedItem().(menuItem); ok {
		switch i.title {
		case "Hosts":
			return "hosts", 0
		case "Packages":
			return "packages", 0
		case "Deployments":
			return "deployments", 0
		case "Quit":
			return "quit", 0
		}
	}
	return "menu", 0
}

type menuItem struct {
	title, desc string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string  { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type MenuDelegate struct{}

func (d MenuDelegate) Render(m *list.Model, sb *string, i list.Item) {
	switch item := i.(type) {
	case menuItem:
		*sb += fmt.Sprintf("  %s\n", item.title)
		*sb += fmt.Sprintf("     %s\n", item.desc)
	}
}

func (d MenuDelegate) Height(*list.Model, int) int       { return 2 }
func (d MenuDelegate) Spacing(*list.Model) int           { return 0 }
func (d MenuDelegate) Update(*list.Model, *list.Msg)     {}
