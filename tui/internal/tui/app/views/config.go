package views

import (
	"fmt"

	"crona/tui/internal/api"
)

func renderConfigView(theme Theme, state ContentState) string {
	active := state.Pane == "config"
	cur := state.Cursors["config"]
	items := configItems(state.ExportAssets)
	indices := filteredStrings(items, state.Filters["config"])
	total := len(indices)
	lines := []string{theme.StylePaneTitle.Render("Config"), renderPaneActionLine(theme, state.Filters["config"], state.Width-6, paneActionsForState(theme, state, active))}
	if state.ExportAssets == nil {
		lines = append(lines, theme.StyleDim.Render("Loading export assets..."))
		return renderPaneBox(theme, active, state.Width, state.Height, stringsJoin(lines))
	}
	if total == 0 {
		lines = append(lines, theme.StyleDim.Render("No config items match the current filter"))
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
		lines = append(lines, renderPaneRowStyled(theme, i, cur, active, items[indices[i]], nil, state.Width))
	}
	if remaining := total - end; remaining > 0 {
		lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↓ %d more", remaining)))
	}
	return renderPaneBox(theme, active, state.Width, state.Height, stringsJoin(lines))
}

func configItems(status *api.ExportAssetStatus) []string {
	if status == nil {
		return nil
	}
	updateState := "up to date"
	if status.DefaultUpdateAvailable {
		updateState = "new default available"
	}
	customized := "default"
	if status.UserTemplateCustomized {
		customized = "customized"
	}
	pdfCustomized := "default"
	if status.PDFTemplateCustomized {
		pdfCustomized = "customized"
	}
	pdfRenderer := "unavailable"
	if status.PDFRendererAvailable {
		pdfRenderer = status.PDFRendererName
	}
	return []string{
		fmt.Sprintf("Daily report template    %s   [%s]", status.TemplateName, customized),
		fmt.Sprintf("PDF report template      %s   [%s]", status.PDFTemplateName, pdfCustomized),
		fmt.Sprintf("Template variables docs  %s", status.TemplateDocsPath),
		fmt.Sprintf("Reports directory        %s", status.ReportsDir),
		fmt.Sprintf("Template update status   %s", updateState),
		fmt.Sprintf("PDF renderer             %s", pdfRenderer),
	}
}
