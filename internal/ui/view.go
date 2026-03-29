package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/lipgloss/v2"
)

func (m Model) View() tea.View {
	baseContent := lipgloss.JoinVertical(lipgloss.Top,
		m.shortcutTable.View(),
		"helli",
	)
	v := tea.NewView(baseContent)
	v.AltScreen = false
	return v

}
