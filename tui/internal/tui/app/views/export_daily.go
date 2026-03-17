package views

import (
	"fmt"
	"path/filepath"

	"crona/tui/internal/api"
)

func renderExportDailyView(theme Theme, state ContentState) string {
	active := state.Pane == "export_reports"
	cur := state.Cursors["export_reports"]
	items := exportReportItems(state.ExportReports)
	indices := filteredStrings(items, state.Filters["export_reports"])
	total := len(indices)
	lines := []string{theme.StylePaneTitle.Render("Daily Exports"), renderPaneActionLine(theme, state.Filters["export_reports"], state.Width-6, paneActionsForState(theme, state, active))}
	if state.ExportAssets == nil {
		lines = append(lines, theme.StyleDim.Render("Loading export configuration..."))
		return renderPaneBox(theme, active, state.Width, state.Height, stringsJoin(lines))
	}
	lines = append(lines, theme.StyleDim.Render("Dir: "+state.ExportAssets.ReportsDir))
	if total == 0 {
		lines = append(lines, "")
		lines = append(lines, theme.StyleDim.Render("No exported markdown reports found"))
		return renderPaneBox(theme, active, state.Width, state.Height, stringsJoin(lines))
	}
	inner := state.Height - 6
	if inner < 1 {
		inner = 1
	}
	start, end := listWindow(cur, total, inner)
	if start > 0 {
		lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↑ %d more", start)))
	}
	for i := start; i < end; i++ {
		lines = append(lines, renderPaneRowStyled(theme, i, cur, active, items[indices[i]], nil, state.Width))
	}
	if remaining := total - end; remaining > 0 {
		lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↓ %d more", remaining)))
	}
	return renderPaneBox(theme, active, state.Width, state.Height, stringsJoin(lines))
}

func exportReportItems(reports []api.ExportReportFile) []string {
	items := make([]string, 0, len(reports))
	for _, report := range reports {
		items = append(items, fmt.Sprintf("%s    [%s] %s    %d B", report.Date, report.Format, filepath.Base(report.Path), report.SizeBytes))
	}
	return items
}
