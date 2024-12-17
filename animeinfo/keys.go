package animeinfo

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	ModifyEntry key.Binding
	Esc         key.Binding
	Help        key.Binding
}

// ShortHelp implements the KeyMap interface.
func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{km.Esc, km.ModifyEntry}
}

// FullHelp implements the KeyMap interface.
func (km KeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{km.Esc, km.ModifyEntry},
		{km.Help},
	}
}

var AnimeInfoKeyMap = KeyMap{
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
