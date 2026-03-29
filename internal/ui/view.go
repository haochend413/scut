package ui

import (
	tea "charm.land/bubbletea/v2"
	"github.com/haochend413/lipgloss/v2"
	"github.com/haochend413/scut/internal/ui/styles"
)

func (m Model) View() tea.View {
	// currContent := ">>> " + m.currentShortcut.Command
	m.shortcutTable.SetStyles(styles.TableStyle)
	baseContent := lipgloss.JoinVertical(lipgloss.Top,
		m.shortcutTable.View(),
	)
	v := tea.NewView(baseContent)
	v.AltScreen = false
	return v
}
