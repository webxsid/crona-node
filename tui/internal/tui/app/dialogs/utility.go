package dialogs

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func renderUtilityDialog(theme Theme, state State) string {
	switch state.Kind {
	case "confirm_delete":
		content := fmt.Sprintf("%s\n\nDelete %s?\n\n%s", theme.StylePaneTitle.Render("Confirm Delete"), theme.StyleError.Render(fallback(state.DeleteLabel, "this item")), theme.StyleDim.Render("[enter] delete   [esc] cancel"))
		return lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(theme.ColorRed).Padding(1, 3).Width(min(state.Width-8, 44)).Render(content)
	case "pick_date":
		rows := []string{theme.StylePaneTitle.Render(state.DateTitle), "", theme.StyleHeader.Render(state.DateHeader), theme.StyleDim.Render(state.DateMonth), "", state.DateGrid, "", theme.StyleDim.Render("[h/j/k/l] move   [,/.] month   [enter] choose   [c] clear   [esc] back")}
		return modal(theme, state.Width, 46, theme.ColorCyan, rows)
	case "create_scratchpad":
		rows := []string{theme.StylePaneTitle.Render("New Scratchpad"), "", theme.StyleDim.Render("Name"), state.Inputs[0].View(), "", theme.StyleDim.Render("Path  (supports [[date]], [[timestamp]])"), state.Inputs[1].View(), "", theme.StyleDim.Render("[tab] next field   [enter] create   [esc] cancel")}
		return modal(theme, state.Width, 54, theme.ColorCyan, rows)
	default:
		return ""
	}
}
