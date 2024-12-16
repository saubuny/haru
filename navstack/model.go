package navstack

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Any model can implement the closable interface for cleanup
type Closeable interface {
	Close() error
}

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
	return item.Init()
}

// Pops an item off the stack, performs any cleanup required, and renders the new top item. Does not do anything if there is only one item left on the stack.
func (m *Model) Pop() tea.Cmd {
	top := m.Top()

	// Don't do anything if trying to pop off an empty stack
	if top == nil {
		return nil
	}

	// Do not allow popping off the last item in a stack
	if len(m.stack) <= 1 {
		return nil
	}

	m.stack = m.stack[:len(m.stack)-1]
	top = m.Top()

	return top.Init()
}

// Returns the top item on the stack
func (m Model) Top() tea.Model {
	if len(m.stack) == 0 {
		return nil
	}

	return m.stack[len(m.stack)-1]
}

func (m *Model) Update(msg tea.Msg) tea.Cmd {
	top := m.Top()
	switch msg := msg.(type) {
	case PopNavigation:
		return m.Pop()
	case PushNavigation:
		return m.Push(msg.Item)
	}

	if top == nil {
		return nil
	}

	updatedModel, cmd := top.Update(msg)
	m.stack[len(m.stack)-1] = updatedModel
	return cmd
}

func (m Model) View() string {
	top := m.Top()
	if top == nil {
		return ""
	}

	return top.View()
}
