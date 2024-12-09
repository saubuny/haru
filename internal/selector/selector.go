package selector

// All selectors i have found either have unneeded functionality or not enough. I will create my own. It is based off of the filepicker bubble.

// Lists selections in a simple list
// Can simply select an option and give you that option

import (
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

func New() Model {
	return Model{
		id:         nextID(),
		Cursor:     ">",
		Height:     0,
		selected:   0,
		selections: []string{},
		KeyMap:     DefaultKeyMap(),
		Styles:     DefaultStyles(),
	}
}

type SelectedMsg string

const (
	marginBottom = 5
	paddingLeft  = 2
)

type KeyMap struct {
	Up     key.Binding
	Down   key.Binding
	Select key.Binding
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
			m.selected = clamp(m.selected+1, 0, len(m.selections))
		case key.Matches(msg, m.KeyMap.Down):
			m.selected = clamp(m.selected-1, 0, len(m.selections))
		case key.Matches(msg, m.KeyMap.Select):
			return m, returnSelection(m.selections[m.selected])
		}
	}
	return m, nil
}

func clamp(value, min, max int) int {
	if value < min {
		return min
	} else if value > max {
		return max
	}
	return value
}
