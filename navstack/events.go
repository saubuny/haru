package navstack

import tea "github.com/charmbracelet/bubbletea"

type PopNavigation struct{}

type PushNavigation struct {
	Item tea.Model
}
