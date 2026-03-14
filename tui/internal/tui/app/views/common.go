package views

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func renderScratchpadPlaceholder(theme Theme, state ContentState) string {
	return renderSimplePane(theme, "Scratchpads", state.Filters["scratchpads"], state.Cursors["scratchpads"], scratchpadItems(state.Scratchpads), state.Pane == "scratchpads", state.Width, state.Height, "No scratchpads — [a] create new")
}

func renderSimplePane(theme Theme, title, filter string, cursor int, items []string, active bool, width, height int, empty string) string {
	indices := filteredStrings(items, filter)
	total := len(indices)
	inner := height - 5
	if inner < 1 {
		inner = 1
	}
	lines := []string{theme.StylePaneTitle.Render(title), renderFilterLine(theme, filter, width-6)}
	if total == 0 {
		lines = append(lines, theme.StyleDim.Render(empty))
	} else {
		start, end := listWindow(cursor, total, inner)
		if start > 0 {
			lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↑ %d more", start)))
		}
		for i := start; i < end; i++ {
			lines = append(lines, renderPaneRowStyled(theme, i, cursor, active, items[indices[i]], nil, width))
		}
		if remaining := total - end; remaining > 0 {
			lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↓ %d more", remaining)))
		}
	}
	return renderPaneBox(theme, active, width, height, stringsJoin(lines))
}

func renderFilterLine(theme Theme, filter string, width int) string {
	if strings.TrimSpace(filter) == "" {
		return theme.StyleDim.Render(truncate("[/] filter", width))
	}
	return theme.StyleHeader.Render(truncate("/ "+filter, width))
}

func renderPaneRowStyled(theme Theme, i, cur int, active bool, text string, contentStyle *lipgloss.Style, width int) string {
	line := truncate(text, width-6)
	if contentStyle != nil {
		line = contentStyle.Render(line)
	}
	if i == cur && active {
		return theme.StyleCursor.Render("▶ " + line)
	}
	if i == cur {
		return theme.StyleSelected.Render("  " + line)
	}
	return theme.StyleNormal.Render("  " + line)
}

func renderPaneBox(theme Theme, active bool, width, height int, content string) string {
	box := theme.StyleInactive
	if active {
		box = theme.StyleActive
	}
	return box.Width(width-2).Height(height-2).Padding(0, 1).Render(content)
}
