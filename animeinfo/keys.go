package animeinfo

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	Exit        key.Binding
	ModifyEntry key.Binding
	Esc         key.Binding
	Help        key.Binding
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
	ModifyEntry: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "modify entry"),
	),
}
