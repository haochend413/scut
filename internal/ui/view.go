package ui

import (
	"os"
	"strings"

	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/scut/internal/ui/styles"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))
	dimStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	hintStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
	emptyStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
)

func (m Model) View() tea.View {
	if m.quitting {
		v := tea.NewView("")
		v.AltScreen = false
		return v
	}

	m.shortcutTable.SetStyles(styles.TableStyle)
	m.historyTable.SetStyles(styles.TableStyle)

	var content string
	switch m.mode {
	case modeHistory:
		content = m.historyView()
	case modeInput:
		content = m.inputView()
	case modeEdit:
		content = m.editView()
	default:
		content = m.mainView()
	}

	v := tea.NewView(content)
	v.AltScreen = false
	return v
}

func (m Model) mainView() string {
	title := titleStyle.Render("scut") + "  " + dimStyle.Render(shortenPath(m.app.GetCWD()))

	var body string
	if len(m.shortcuts) == 0 {
		body = emptyStyle.Render("  (no shortcuts yet)")
	} else {
		body = dropHeader(m.shortcutTable.View())
	}

	hints := hintStyle.Render("  a add  i edit  h history  ⌫ delete  ↵ copy  q quit")
	return lipgloss.JoinVertical(lipgloss.Left, title, body, hints)
}

func (m Model) historyView() string {
	title := titleStyle.Render("scut") + "  " + dimStyle.Render("shell history")

	searchLine := hintStyle.Render("  / ") + m.historySearch.View()

	var body string
	if len(m.history) == 0 {
		body = emptyStyle.Render("  (no history found)")
	} else if len(m.filteredHistory) == 0 {
		body = emptyStyle.Render("  (no matches)")
	} else {
		body = dropHeader(m.historyTable.View())
	}

	hints := hintStyle.Render("  ↑↓ navigate  ↵ save to shortcuts  esc back")
	return lipgloss.JoinVertical(lipgloss.Left, title, searchLine, body, hints)
}

func (m Model) inputView() string {
	title := titleStyle.Render("scut") + "  " + dimStyle.Render(shortenPath(m.app.GetCWD()))

	var body string
	if len(m.shortcuts) == 0 {
		body = emptyStyle.Render("  (no shortcuts yet)")
	} else {
		body = dropHeader(m.shortcutTable.View())
	}

	inputLine := hintStyle.Render("  > ") + m.input.View()
	hints := hintStyle.Render("  ↵ save  esc cancel")
	return lipgloss.JoinVertical(lipgloss.Left, title, body, inputLine, hints)
}

func (m Model) editView() string {
	title := titleStyle.Render("scut") + "  " + dimStyle.Render(shortenPath(m.app.GetCWD()))

	var body string
	if len(m.shortcuts) == 0 {
		body = emptyStyle.Render("  (no shortcuts yet)")
	} else {
		body = dropHeader(m.shortcutTable.View())
	}

	inputLine := hintStyle.Render("  edit > ") + m.input.View()
	hints := hintStyle.Render("  ↵ save  esc cancel")
	return lipgloss.JoinVertical(lipgloss.Left, title, body, inputLine, hints)
}

// dropHeader removes the blank header row that the table always prepends.
func dropHeader(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[i+1:]
	}
	return s
}

func shortenPath(p string) string {
	if home, err := os.UserHomeDir(); err == nil && strings.HasPrefix(p, home) {
		return "~" + p[len(home):]
	}
	return p
}
