package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
)

type appItem string

func (a appItem) Title() string       { return string(a) }
func (a appItem) Description() string { return "" }
func (a appItem) FilterValue() string { return string(a) }

type model struct {
	list list.Model
}

func searchFlathub(query string) []list.Item {
	cmd := exec.Command("flatpak", "search", query)
	out, err := cmd.Output()
	if err != nil {
		return []list.Item{appItem("Error running flatpak command")}
	}

	lines := strings.Split(string(out), "\n")
	items := []list.Item{}

	for _, l := range lines[1:] { // skip header
		if strings.TrimSpace(l) == "" {
			continue
		}
		items = append(items, appItem(l))
	}
	return items
}

func initialModel() model {
	items := searchFlathub("browser")
	l := list.New(items, list.NewDefaultDelegate(), 50, 15)
	l.Title = "Flathub Search Results"
	return model{list: l}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {

		case "q", "ctrl+c":
			return m, tea.Quit

		case "enter":
			selected := m.list.SelectedItem().(appItem)
			fmt.Println("\nInstalling:", selected)
			parts := strings.Fields(string(selected))
			if len(parts) > 0 {
				exec.Command("flatpak", "install", "-y", parts[0]).Run()
			}
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return m.list.View()
}

func main() {
	if _, err := tea.NewProgram(initialModel()).Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
