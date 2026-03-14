package views

import "fmt"

func renderSettingsView(theme Theme, state ContentState) string {
	active := state.Pane == "settings"
	cur := state.Cursors["settings"]
	indices := filteredSettingIndices(state.Filters["settings"], state.Settings)
	total := len(indices)
	lines := []string{theme.StylePaneTitle.Render("Settings"), theme.StyleDim.Render("Use [j/k] to move and [h/l] or [enter] to change values."), ""}
	if state.Settings == nil {
		lines = append(lines, theme.StyleDim.Render("Loading settings..."))
		return renderPaneBox(theme, active, state.Width, state.Height, stringsJoin(lines))
	}
	rows := []struct{ label, value string }{
		{"Timer Mode", string(state.Settings.TimerMode)},
		{"Breaks Enabled", onOff(state.Settings.BreaksEnabled)},
		{"Work Duration", fmt.Sprintf("%d min", state.Settings.WorkDurationMinutes)},
		{"Short Break", fmt.Sprintf("%d min", state.Settings.ShortBreakMinutes)},
		{"Long Break", fmt.Sprintf("%d min", state.Settings.LongBreakMinutes)},
		{"Long Break Enabled", onOff(state.Settings.LongBreakEnabled)},
		{"Cycles Before Long Break", fmt.Sprintf("%d", state.Settings.CyclesBeforeLongBreak)},
		{"Auto Start Breaks", onOff(state.Settings.AutoStartBreaks)},
		{"Auto Start Work", onOff(state.Settings.AutoStartWork)},
	}
	if total == 0 {
		lines = append(lines, theme.StyleDim.Render("No settings match the current filter"))
		return renderPaneBox(theme, active, state.Width, state.Height, stringsJoin(lines))
	}
	inner := state.Height - 5
	if inner < 1 {
		inner = 1
	}
	start, end := listWindow(cur, total, inner)
	if start > 0 {
		lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↑ %d more", start)))
	}
	for i := start; i < end; i++ {
		idx := indices[i]
		row := fmt.Sprintf("%-24s %s", rows[idx].label, rows[idx].value)
		if i == cur && active {
			lines = append(lines, theme.StyleCursor.Render("▶ "+row))
		} else if i == cur {
			lines = append(lines, theme.StyleSelected.Render("  "+row))
		} else {
			lines = append(lines, theme.StyleNormal.Render("  "+row))
		}
	}
	if remaining := total - end; remaining > 0 {
		lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↓ %d more", remaining)))
	}
	return renderPaneBox(theme, active, state.Width, state.Height, stringsJoin(lines))
}
