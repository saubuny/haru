package main

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Up        key.Binding
	Down      key.Binding
	Exit      key.Binding
	Select    key.Binding
	Esc       key.Binding
	Help      key.Binding
	ChangeTab key.Binding
	AnimeInfo bool
}

// ShortHelp implements the KeyMap interface.
func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Up, km.Down}
}

// FullHelp implements the KeyMap interface.
func (km KeyMap) FullHelp() [][]key.Binding {
	if km.AnimeInfo {
		return [][]key.Binding{
			{km.Exit, km.Esc, km.Help},
			{},
		}
	}
	return [][]key.Binding{
		{km.Up, km.Down, km.Esc, km.ChangeTab},
		{km.Exit, km.Select, km.Help},
	}
}

var AnimeInfoKeyMap = KeyMap{
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "exit"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "go back"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	AnimeInfo: true, // Probably a better way to implement this
}

var DefaultKeyMap = KeyMap{
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "move up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "move down"),
	),
	Exit: key.NewBinding(
		key.WithKeys("ctrl+c"),
		key.WithHelp("ctrl+c", "exit"),
	),
	Select: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "select"),
	),
	Esc: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "toggle search"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	ChangeTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "change tab"),
	),
}
