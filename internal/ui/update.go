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
	mainQuitKey  = key.NewBinding(key.WithKeys("ctrl+c", "q"))
	histBackKey  = key.NewBinding(key.WithKeys("esc", "q"))
	backspaceKey = key.NewBinding(key.WithKeys("backspace"))
	addKey       = key.NewBinding(key.WithKeys("a"))
	histKey      = key.NewBinding(key.WithKeys("h"))
	refreshKey   = key.NewBinding(key.WithKeys("r"))
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		w := max(20, msg.Width-4)
		m.shortcutTable.SetWidth(w)
		m.historyTable.SetWidth(w)
		m.input.SetWidth(w)
		tableH := min(10, max(5, msg.Height/3))
		m.shortcutTable.SetHeight(tableH)
		m.historyTable.SetHeight(tableH)
		return m, nil

	case table.MoveSelectMsg:
		// cursor already updated inside the table; selectedShortcut() reads it dynamically
		return m, nil

	case tea.KeyMsg:
		switch m.mode {

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
			case key.Matches(msg, histBackKey):
				m.mode = modeMain
				return m, nil
			case key.Matches(msg, enterKey):
				if h := m.selectedHistoryCmd(); h != "" {
					cwd, _ := os.Getwd()
					_ = m.app.AddShortcut(models.Shortcut{WorkDirectory: cwd, Command: h})
					m.refreshShortcuts()
					m.shortcutTable.SetCursor(len(m.shortcuts) - 1)
					m.mode = modeMain
				}
				return m, nil
			}
			m.historyTable, cmd = m.historyTable.Update(msg)
			return m, cmd

		default: // modeMain
			switch {
			case key.Matches(msg, mainQuitKey):
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
			case key.Matches(msg, histKey):
				m.mode = modeHistory
				return m, nil
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
	case modeInput:
		m.input, cmd = m.input.Update(msg)
	case modeHistory:
		m.historyTable, cmd = m.historyTable.Update(msg)
	default:
		m.shortcutTable, cmd = m.shortcutTable.Update(msg)
	}
	return m, cmd
}
