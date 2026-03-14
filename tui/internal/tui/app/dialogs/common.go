package dialogs

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

func modal(theme Theme, width, maxWidth int, border lipgloss.Color, rows []string) string {
	return lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(border).Padding(1, 3).Width(min(width-8, maxWidth)).Render(strings.Join(rows, "\n"))
}

func renderSingleInput(theme Theme, width int, title, label string, inputs []textinput.Model, border lipgloss.Color, hint string) string {
	rows := []string{theme.StylePaneTitle.Render(title), "", theme.StyleDim.Render(label), inputs[0].View(), "", theme.StyleDim.Render(hint)}
	return modal(theme, width, 52, border, rows)
}

func renderSelector(theme Theme, label string, active bool) string {
	style := theme.StyleNormal
	if active {
		style = theme.StyleCursor
	}
	return style.Render("[ " + label + " ]")
}

func plainIssueStatus(status string) string {
	switch status {
	case "in_progress":
		return "in progress"
	case "in_review":
		return "in review"
	default:
		return status
	}
}

func fallback(v, def string) string {
	if strings.TrimSpace(v) == "" {
		return def
	}
	return v
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
