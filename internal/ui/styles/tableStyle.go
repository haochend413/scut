package styles

import (
	"github.com/haochend413/bubbles/v2/table"
	"github.com/haochend413/lipgloss/v2"
)

var TableStyle = table.Styles{
	// Header: lipgloss.NewStyle().
	// 	Bold(true).
	// 	Padding(0, 0).
	// 	Foreground(lipgloss.Color("252")),

	// Cell: lipgloss.NewStyle().
	// 	Padding(0, 0),
	Selected: lipgloss.NewStyle().
		Foreground(lipgloss.Color("46")),
}
