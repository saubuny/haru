package navstack

// This is a little overkill since we only ever go like 3 layers deep but its fine :)

import (
	tea "github.com/charmbracelet/bubbletea"
)

func New(model tea.Model) Model {
	return Model{stack: []tea.Model{model}}
}

// Helper function to easily return messages as a tea.Cmd
func Cmd(msg interface{}) tea.Cmd {
	return func() tea.Msg {
		return msg
	}
}

type Model struct {
	stack []tea.Model
}

func (m Model) Init() tea.Cmd {
	if m.Top() == nil {
		return nil
	}
	return m.Top().Init()
}

// Pushes an item onto the stack. Performs any cleanup required by the previous top item and renders the new top item.
func (m *Model) Push(item tea.Model) tea.Cmd {
	m.stack = append(m.stack, item)
	return tea.Sequence(item.Init(), tea.WindowSize())
}

// Pops an item off the stack, performs any cleanup required, and renders the new top item. Does not do anything if there is only one item left on the stack.
func (m *Model) Pop() tea.Cmd {
	top := m.Top()

	// Don't do anything if trying to pop off an empty stack (shouldn't be possible anyway)
	if top == nil {
		return nil
	}

	// Do not allow popping off the last item in a stack
	if len(m.stack) <= 1 {
		return nil
	}

	m.stack = m.stack[:len(m.stack)-1]
	return nil
}

// Returns the top item on the stack
func (m Model) Top() tea.Model {
	if len(m.stack) == 0 {
		return nil
	}

	return m.stack[len(m.stack)-1]
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	top := m.Top()
	switch msg := msg.(type) {
	case PopNavigation:
		return m, m.Pop()
	case PushNavigation:
		return m, m.Push(msg.Item)
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		}
	}

	if top == nil {
		return m, nil
	}

	updatedModel, cmd := top.Update(msg)
	m.stack[len(m.stack)-1] = updatedModel
	return m, cmd
}

func (m Model) View() string {
	top := m.Top()
	if top == nil {
		return ""
	}

	return top.View()
}
