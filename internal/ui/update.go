package ui

import (
	"os"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/bubbles/v2/key"
	"github.com/haochend413/bubbles/v2/table"
	"github.com/haochend413/bubbles/v2/textinput"
	"github.com/haochend413/scut/internal/models"
	"github.com/haochend413/scut/internal/utils"
)

var (
	ctrlCKey     = key.NewBinding(key.WithKeys("ctrl+c"))
	escKey       = key.NewBinding(key.WithKeys("esc"))
	enterKey     = key.NewBinding(key.WithKeys("enter"))
	backspaceKey = key.NewBinding(key.WithKeys("backspace"))
	addKey       = key.NewBinding(key.WithKeys("a"))
	dupKey       = key.NewBinding(key.WithKeys("d"))
	editKey      = key.NewBinding(key.WithKeys("i"))
	histKey      = key.NewBinding(key.WithKeys("h"))
	refreshKey   = key.NewBinding(key.WithKeys("r"))
	tableNavKey  = key.NewBinding(key.WithKeys("up", "down", "pgup", "pgdown"))
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		col := table.Column{Title: "", Width: msg.Width}
		m.shortcutTable.SetColumns([]table.Column{col})
		m.historyTable.SetColumns([]table.Column{col})
		m.shortcutTable.SetWidth(msg.Width)
		m.historyTable.SetWidth(msg.Width)
		m.input.SetWidth(msg.Width)
		m.historySearch.SetWidth(msg.Width)
		tableH := min(10, max(5, msg.Height/3))
		m.shortcutTable.SetHeight(tableH)
		m.historyTable.SetHeight(tableH)
		return m, nil

	case table.MoveSelectMsg:
		// cursor already updated inside the table; selectedShortcut() reads it dynamically
		return m, nil

	case tea.KeyMsg:
		switch m.mode {

		case modeEdit:
			switch {
			case key.Matches(msg, ctrlCKey):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, escKey):
				m.mode = modeMain
				m.input.Blur()
				m.input.SetValue("")
				m.editingID = 0
				return m, nil
			case key.Matches(msg, enterKey):
				if v := m.input.Value(); v != "" && m.editingID != 0 {
					cursor := m.shortcutTable.Cursor()
					_ = m.app.UpdateShortcut(m.editingID, v)
					m.refreshShortcuts()
					m.shortcutTable.SetCursor(cursor)
				}
				m.mode = modeMain
				m.input.Blur()
				m.input.SetValue("")
				m.editingID = 0
				return m, nil
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd

		case modeInput:
			switch {
			case key.Matches(msg, ctrlCKey):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, escKey):
				m.mode = modeMain
				m.input.Blur()
				m.input.SetValue("")
				return m, nil
			case key.Matches(msg, enterKey):
				if v := m.input.Value(); v != "" {
					cwd, _ := os.Getwd()
					_ = m.app.AddShortcut(models.Shortcut{WorkDirectory: cwd, Command: v})
					m.refreshShortcuts()
					m.shortcutTable.SetCursor(len(m.shortcuts) - 1)
				}
				m.mode = modeMain
				m.input.Blur()
				m.input.SetValue("")
				return m, nil
			}
			m.input, cmd = m.input.Update(msg)
			return m, cmd

		case modeHistory:
			switch {
			case key.Matches(msg, ctrlCKey):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, escKey):
				if m.historySearch.Value() != "" {
					m.historySearch.SetValue("")
					m.applyHistoryFilter()
					return m, nil
				}
				m.mode = modeMain
				m.historySearch.Blur()
				return m, nil
			case key.Matches(msg, enterKey):
				if h := m.selectedHistoryCmd(); h != "" {
					cwd, _ := os.Getwd()
					_ = m.app.AddShortcut(models.Shortcut{WorkDirectory: cwd, Command: h})
					m.refreshShortcuts()
					m.shortcutTable.SetCursor(len(m.shortcuts) - 1)
					m.mode = modeMain
					m.historySearch.SetValue("")
					m.historySearch.Blur()
					m.applyHistoryFilter()
				}
				return m, nil
			case key.Matches(msg, tableNavKey):
				m.historyTable, cmd = m.historyTable.Update(msg)
				return m, cmd
			}
			// All other keys go to the search input.
			m.historySearch, cmd = m.historySearch.Update(msg)
			m.applyHistoryFilter()
			return m, cmd

		default: // modeMain
			switch {
			case key.Matches(msg, ctrlCKey):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, enterKey):
				if sc := m.selectedShortcut(); sc != nil {
					_ = utils.CopyToClipboard(sc.Command)
				}
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, backspaceKey):
				if sc := m.selectedShortcut(); sc != nil {
					cursor := m.shortcutTable.Cursor()
					m.app.DeleteShortcut(sc.ID)
					m.refreshShortcuts()
					if cursor >= len(m.shortcuts) {
						cursor = max(0, len(m.shortcuts)-1)
					}
					m.shortcutTable.SetCursor(cursor)
				}
				return m, nil
			case key.Matches(msg, addKey):
				m.mode = modeInput
				focusCmd := m.input.Focus()
				return m, tea.Batch(focusCmd, textinput.Blink)
			case key.Matches(msg, dupKey):
				if sc := m.selectedShortcut(); sc != nil {
					cursor := m.shortcutTable.Cursor()
					cwd, _ := os.Getwd()
					inserted, err := m.app.DuplicateShortcut(models.Shortcut{WorkDirectory: cwd, Command: sc.Command})
					if err == nil {
						newList := make([]models.Shortcut, 0, len(m.shortcuts)+1)
						newList = append(newList, m.shortcuts[:cursor+1]...)
						newList = append(newList, inserted)
						newList = append(newList, m.shortcuts[cursor+1:]...)
						m.shortcuts = newList
						rows := make([]table.Row, 0, len(m.shortcuts))
						for _, s := range m.shortcuts {
							rows = append(rows, table.Row{s.Command})
						}
						m.shortcutTable.SetRows(rows)
						m.shortcutTable.SetCursor(cursor + 1)
					}
				}
				return m, nil
			case key.Matches(msg, editKey):
				if sc := m.selectedShortcut(); sc != nil {
					m.editingID = sc.ID
					m.input.SetValue(sc.Command)
					m.input.CursorEnd()
					m.mode = modeEdit
					focusCmd := m.input.Focus()
					return m, tea.Batch(focusCmd, textinput.Blink)
				}
			case key.Matches(msg, histKey):
				m.mode = modeHistory
				m.refreshHistory()
				focusCmd := m.historySearch.Focus()
				return m, tea.Batch(focusCmd, textinput.Blink)
			case key.Matches(msg, refreshKey):
				m.refreshShortcuts()
				m.refreshHistory()
				return m, nil
			}
			m.shortcutTable, cmd = m.shortcutTable.Update(msg)
			return m, cmd
		}
	}

	// Non-key messages: forward to the active component.
	switch m.mode {
	case modeInput, modeEdit:
		m.input, cmd = m.input.Update(msg)
	case modeHistory:
		m.historySearch, cmd = m.historySearch.Update(msg)
	default:
		m.shortcutTable, cmd = m.shortcutTable.Update(msg)
	}
	return m, cmd
}
