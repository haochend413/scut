package ui

import (
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/table"
	"github.com/haochend413/bubbles/v2/textinput"
	"github.com/haochend413/scut/internal/app"
	"github.com/haochend413/scut/internal/models"
)

type viewMode int

const (
	modeMain    viewMode = iota
	modeHistory          // browse shell history to add a command
	modeInput            // type a command manually
	modeEdit             // edit the selected command in place
)

type Model struct {
	app             *app.App
	mode            viewMode
	shortcuts       []models.Shortcut
	history         []string
	filteredHistory []string
	shortcutTable   table.Model
	historyTable    table.Model
	input           textinput.Model
	historySearch   textinput.Model
	editingID       uint
	width           int
	height          int
	quitting        bool
}

func NewModel(application *app.App) Model {
	shortcutTable := table.New(
		table.WithColumns([]table.Column{{Title: "", Width: 80}}),
		table.WithFocused(true),
		table.WithHeight(8),
	)
	historyTable := table.New(
		table.WithColumns([]table.Column{{Title: "", Width: 80}}),
		table.WithFocused(true),
		table.WithHeight(8),
	)

	ti := textinput.New()
	ti.Placeholder = "type a command..."
	ti.CharLimit = 512

	hs := textinput.New()
	hs.Placeholder = "search..."
	hs.CharLimit = 0

	m := Model{
		app:           application,
		mode:          modeMain,
		shortcutTable: shortcutTable,
		historyTable:  historyTable,
		input:         ti,
		historySearch: hs,
	}
	m.refreshShortcuts()
	m.refreshHistory()
	return m
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m *Model) refreshShortcuts() {
	m.shortcuts = m.app.DisplayCWDShortcuts()
	rows := make([]table.Row, 0, len(m.shortcuts))
	for _, sc := range m.shortcuts {
		rows = append(rows, table.Row{sc.Command})
	}
	m.shortcutTable.SetRows(rows)
}

func (m *Model) refreshHistory() {
	m.history = m.app.GetHistory()
	m.filteredHistory = m.history
	rows := make([]table.Row, 0, len(m.history))
	for _, cmd := range m.history {
		rows = append(rows, table.Row{cmd})
	}
	m.historyTable.SetRows(rows)
}

func (m *Model) applyHistoryFilter() {
	query := m.historySearch.Value()
	if query == "" {
		m.filteredHistory = m.history
	} else {
		filtered := make([]string, 0)
		for _, cmd := range m.history {
			if fuzzyMatch(query, cmd) {
				filtered = append(filtered, cmd)
			}
		}
		m.filteredHistory = filtered
	}
	rows := make([]table.Row, 0, len(m.filteredHistory))
	for _, cmd := range m.filteredHistory {
		rows = append(rows, table.Row{cmd})
	}
	m.historyTable.SetRows(rows)
	m.historyTable.SetCursor(0)
}

// fuzzyMatch reports whether all runes of query appear in target in order.
func fuzzyMatch(query, target string) bool {
	query = strings.ToLower(query)
	target = strings.ToLower(target)
	qi := 0
	for i := 0; i < len(target) && qi < len(query); i++ {
		if target[i] == query[qi] {
			qi++
		}
	}
	return qi == len(query)
}

// selectedShortcut returns the shortcut under the cursor, or nil if the list is empty.
func (m *Model) selectedShortcut() *models.Shortcut {
	c := m.shortcutTable.Cursor()
	if c < 0 || c >= len(m.shortcuts) {
		return nil
	}
	return &m.shortcuts[c]
}

// selectedHistoryCmd returns the filtered history entry under the cursor, or "" if none.
func (m *Model) selectedHistoryCmd() string {
	c := m.historyTable.Cursor()
	if c < 0 || c >= len(m.filteredHistory) {
		return ""
	}
	return m.filteredHistory[c]
}
