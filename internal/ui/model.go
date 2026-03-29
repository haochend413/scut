package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/table"

	// "github.com/haochend413/bubbles/v2/textinput"
	"github.com/haochend413/scut/internal/app"
	"github.com/haochend413/scut/internal/models"
)

type Model struct {
	app             *app.App
	shortcutTable   table.Model
	currentShortcut *models.Shortcut
}

func NewModel(application *app.App) Model {
	// init model
	shortcut_column := []table.Column{
		{Title: "   ", Width: 200},
	}
	shortcutTable := table.New(
		table.WithColumns(shortcut_column),
		table.WithFocused(true),
		table.WithHeight(40),
	)
	curr := application.ShortcutMgr.GetSelectedShortCut(0)
	m := Model{
		app:             application,
		shortcutTable:   shortcutTable,
		currentShortcut: curr,
	}
	m.updateShortcutTable()
	return m
}

// type tickMsg time.Time

// func tick() tea.Cmd {
// 	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
// 		return tickMsg(t)
// 	})
// }

// // Init initializes the Bubble Tea model
func (m Model) Init() tea.Cmd {
	// return tea.Batch(textinput.Blink, tick())
	return nil
}
func (m *Model) RemoveCurrent() {
	cursor := m.shortcutTable.Cursor()
	m.app.DeleteShortcut(m.currentShortcut.ID)
	new_cursor := max(0, cursor-1)
	m.shortcutTable.SetCursor(new_cursor)
	m.currentShortcut = m.app.ShortcutMgr.GetSelectedShortCut(new_cursor)
	m.updateShortcutTable()
}

func (m *Model) updateShortcutTable() {
	shortcuts := m.app.DisplayCWDShortcuts()

	rows := make([]table.Row, 0, len(shortcuts))
	for _, sc := range shortcuts {
		rows = append(rows, table.Row{sc.Command})
	}
	m.shortcutTable.SetRows(rows)
}
