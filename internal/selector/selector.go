package selector

// All selectors i have found either have unneeded functionality or not enough. I will create my own. It is based off of the filepicker bubble.

// Lists selections in a simple list
// Can simply select an option and give you that option

import (
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var lastID int

func nextID() int {
	lastID += 1
	return lastID
}

type Model struct {
	id         int
	KeyMap     KeyMap
	selected   int
	selections []string
	Height     int
	Cursor     string
	Styles     Styles
}

func New(selections []string) Model {
	return Model{
		id:         nextID(),
		Cursor:     "> ",
		Height:     0,
		selected:   0,
		selections: selections,
		KeyMap:     DefaultKeyMap(),
		Styles:     DefaultStyles(),
	}
}

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Up, km.Down, km.Select}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{km.Up, km.Down}, {km.Select}}
}

func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("k", "up"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("j", "down"),
			key.WithHelp("↓/j", "move down"),
		),
		Select: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
	}
}

type Styles struct {
	Cursor   lipgloss.Style
	Selected lipgloss.Style
}

func DefaultStyles() Styles {
	return Styles{
		Cursor:   lipgloss.NewStyle().Foreground(lipgloss.Color("247")),
		Selected: lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Bold(true),
	}
}

type SelectedMsg string

func returnSelection(selection string) tea.Cmd {
	return func() tea.Msg {
		return SelectedMsg(selection)
	}
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.KeyMap.Up):
			m.selected = clamp(m.selected-1, 0, len(m.selections)-1)
			return m, nil
		case key.Matches(msg, m.KeyMap.Down):
			m.selected = clamp(m.selected+1, 0, len(m.selections)-1)
			return m, nil
		case key.Matches(msg, m.KeyMap.Select):
			// Why is this not working ???
			return m, returnSelection(m.selections[m.selected])
		}
	}
	return m, nil
}

func (m Model) View() string {
	var s strings.Builder
	for i, item := range m.selections {
		if m.selected == i {
			s.WriteString(m.Styles.Cursor.Render(m.Cursor) + item)
		} else {
			s.WriteString("  " + item)
		}
		s.WriteRune('\n')
	}

	// for i := lipgloss.Height(s.String()); i <= m.Height; i++ {
	// 	s.WriteRune('\n')
	// }

	return s.String()
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
