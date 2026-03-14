package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func renderMetaView(theme Theme, state ContentState) string {
	topH, botH := splitVertical(state.Height, 7, 8, state.Height*30/100)
	streamsEmpty := "No streams — [a] create new"
	if state.Context == nil || state.Context.RepoID == nil {
		streamsEmpty = "No repo checked out — [1] then [c]"
	}
	issuesEmpty := "No issues — [a] create new"
	if state.Context == nil || state.Context.StreamID == nil {
		issuesEmpty = "No stream checked out — [2] then [c]"
	}
	leftW, rightW := splitHorizontal(state.Width, 18, 18, state.Width/2)
	repoPane := renderSimplePane(theme, "Repos [1]", state.Filters["repos"], state.Cursors["repos"], repoItems(state.Repos), state.Pane == "repos", leftW, topH, "No repos — [a] create new")
	streamPane := renderSimplePane(theme, "Streams [2]", state.Filters["streams"], state.Cursors["streams"], streamItems(state.Streams), state.Pane == "streams", rightW, topH, streamsEmpty)
	topRow := lipgloss.JoinHorizontal(lipgloss.Top, repoPane, streamPane)
	issuePane := renderMetaIssues(theme, state, state.Width, botH, issuesEmpty)
	return lipgloss.JoinVertical(lipgloss.Left, topRow, issuePane)
}

func renderMetaIssues(theme Theme, state ContentState, width, height int, emptyText string) string {
	active := state.Pane == "issues"
	cur := state.Cursors["issues"]
	var issues []apiIssue
	for _, issue := range state.Issues {
		issues = append(issues, newAPIIssue(issue.ID, issue.Title, issue.Status, issue.EstimateMinutes, issue.TodoForDate))
	}
	indices := filteredIssueIndices(issues, state.Filters["issues"])
	total := len(indices)
	inner := height - 5
	if inner < 1 {
		inner = 1
	}
	lines := []string{theme.StylePaneTitle.Render("Issues [3]"), renderFilterLine(theme, state.Filters["issues"], width-6)}
	if total == 0 {
		lines = append(lines, theme.StyleDim.Render(emptyText))
	} else {
		start, end := listWindow(cur, total, inner)
		if start > 0 {
			lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↑ %d more", start)))
		}
		for i := start; i < end; i++ {
			issue := issues[indices[i]]
			text := fmt.Sprintf("[%s] %s%s", plainIssueStatus(string(issue.Status)), issue.Title, issueDueSuffix(issue.TodoForDate))
			lines = append(lines, renderPaneRowStyled(theme, i, cur, active, text, issueStatusStyle(theme, string(issue.Status)), width))
		}
		if remaining := total - end; remaining > 0 {
			lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↓ %d more", remaining)))
		}
	}
	return renderPaneBox(theme, active, width, height, stringsJoin(lines))
}
