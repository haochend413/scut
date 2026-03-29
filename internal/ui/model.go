package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/table"

	// "github.com/haochend413/bubbles/v2/textinput"
	"github.com/haochend413/scut/internal/app"
)

type Model struct {
	app           *app.App
	shortcutTable table.Model
}

func NewModel(application *app.App) Model {
	// init model
	shortcut_column := []table.Column{
		{Title: "cmds", Width: 50},
	}
	shortcutTable := table.New(
		table.WithColumns(shortcut_column),
		table.WithFocused(true),
		table.WithHeight(40),
	)
	m := Model{
		app:           application,
		shortcutTable: shortcutTable,
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

// func (m *Model) updateShortcutTable() {
// 	shortcuts := m.app.DisplayCWDShortcuts()
// 	rows := make([]table.Row, 0, len(shortcuts))

// 	for _, sc := range shortcuts {
// 		cmdStr := sc.Command
// 		row := table.Row{
// 			fmt.Sprintf("%s", cmdStr),
// 		}
// 		rows = append(rows, row)
// 	}

// 	m.shortcutTable.SetRows(rows)
// }

func (m *Model) updateShortcutTable() {
	shortcuts := m.app.DisplayCWDShortcuts()

	rows := make([]table.Row, 0, len(shortcuts))
	for _, sc := range shortcuts {
		rows = append(rows, table.Row{sc.Command})
	}
	m.shortcutTable.SetRows(rows)
}
