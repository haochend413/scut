package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/key"
)

type keyMap struct {
	Quit    key.Binding
	Refresh key.Binding
	Up      key.Binding
	Down    key.Binding
}

var keys = keyMap{
	Quit:    key.NewBinding(key.WithKeys("ctrl+c", "q")),
	Refresh: key.NewBinding(key.WithKeys("r")),
	Up:      key.NewBinding(key.WithKeys("up", "k")),
	Down:    key.NewBinding(key.WithKeys("down", "j")),
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.shortcutTable.SetWidth(max(30, msg.Width-2))
		m.shortcutTable.SetHeight(max(5, 5))
		return m, nil

	case tea.KeyMsg:
		switch {
		case key.Matches(msg, keys.Quit):
			m.app.OnClose()
			return m, tea.Quit

		case key.Matches(msg, keys.Refresh):
			m.updateShortcutTable()
			return m, nil
		}
	}

	m.shortcutTable, cmd = m.shortcutTable.Update(msg)
	return m, cmd
}
