package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/key"
	"github.com/haochend413/bubbles/v2/table"
	"github.com/haochend413/scut/internal/utils"
)

type keyMap struct {
	Quit          key.Binding
	Refresh       key.Binding
	SelectAndQuit key.Binding
}

var keys = keyMap{
	Quit:          key.NewBinding(key.WithKeys("ctrl+c", "q")),
	Refresh:       key.NewBinding(key.WithKeys("r")),
	SelectAndQuit: key.NewBinding(key.WithKeys("enter")),
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
		case key.Matches(msg, keys.SelectAndQuit):
			utils.CopyToClipboard(m.currentShortcut.Command)
			m.app.OnClose()
			return m, tea.Quit
		}
	case table.MoveSelectMsg:
		// get the content shorcut. reset this for m.
		cursor := m.shortcutTable.Cursor()
		m.currentShortcut = m.app.ShortcutMgr.GetSelectedShortCut(cursor)
		return m, nil
	}

	m.shortcutTable, cmd = m.shortcutTable.Update(msg)
	return m, cmd
}
