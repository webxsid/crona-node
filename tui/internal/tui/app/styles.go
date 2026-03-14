package app

import "github.com/charmbracelet/lipgloss"

var (
	colorBlue    = lipgloss.Color("12")
	colorCyan    = lipgloss.Color("14")
	colorGreen   = lipgloss.Color("10")
	colorMagenta = lipgloss.Color("13")
	colorSubtle  = lipgloss.Color("7")
	colorYellow  = lipgloss.Color("11")
	colorRed     = lipgloss.Color("9")
	colorDim     = lipgloss.Color("8")
	colorWhite   = lipgloss.Color("15")

	styleActive = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorCyan)

	styleInactive = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(colorDim)

	stylePaneTitle = lipgloss.NewStyle().Bold(true).Foreground(colorCyan)
	styleDim       = lipgloss.NewStyle().Foreground(colorDim)
	styleCursor    = lipgloss.NewStyle().Foreground(colorGreen).Bold(true)
	styleHeader    = lipgloss.NewStyle().Foreground(colorCyan)
	styleError     = lipgloss.NewStyle().Foreground(colorRed)
	styleSelected  = lipgloss.NewStyle().Foreground(colorGreen)
	styleNormal    = lipgloss.NewStyle()
)
