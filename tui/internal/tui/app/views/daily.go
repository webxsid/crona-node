package views

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func renderDailyView(theme Theme, state ContentState) string {
	summaryH, listH := splitVertical(state.Height, 9, 8, state.Height/3)
	return lipgloss.JoinVertical(lipgloss.Left, renderDailySummary(theme, state, state.Width, summaryH), renderDailyIssues(theme, state, state.Width, listH))
}

func renderDailySummary(theme Theme, state ContentState, width, height int) string {
	dateText := currentDashboardDate(state)
	totalIssues, totalEstimate, completedCount, abandonedCount, workedSeconds := 0, 0, 0, 0, 0
	if state.DailySummary != nil {
		dateText = state.DailySummary.Date
		totalIssues = state.DailySummary.TotalIssues
		totalEstimate = state.DailySummary.TotalEstimatedMinutes
		completedCount = state.DailySummary.CompletedIssues
		abandonedCount = state.DailySummary.AbandonedIssues
		workedSeconds = state.DailySummary.WorkedSeconds
	}
	resolvedCount := completedCount + abandonedCount
	progressText := fmt.Sprintf("%d/%d planned resolved", resolvedCount, totalIssues)
	workedEstimateText := fmt.Sprintf("%s / %dm", formatClock(workedSeconds), totalEstimate)
	lines := []string{
		theme.StylePaneTitle.Render("Daily Dashboard"),
		theme.StyleDim.Render(fmt.Sprintf("date: %s   [,] prev   [.] next   [g] today", dateText)),
		"",
		progressText,
		renderProgressBar(theme, completedCount, abandonedCount, max(0, totalIssues-resolvedCount), max(20, width-10)),
		"",
		fmt.Sprintf("%s   %s   %s",
			theme.StyleHeader.Render(fmt.Sprintf("planned %d", totalIssues)),
			theme.StyleHeader.Render(fmt.Sprintf("worked %s", workedEstimateText)),
			issueStatusStyle(theme, "done").Render(fmt.Sprintf("done %d", completedCount))+"  "+issueStatusStyle(theme, "abandoned").Render(fmt.Sprintf("abandoned %d", abandonedCount)),
		),
	}
	return lipgloss.NewStyle().BorderStyle(lipgloss.RoundedBorder()).BorderForeground(theme.ColorCyan).Padding(1, 2).Width(width - 2).Height(max(7, height-2)).Render(stringsJoin(lines))
}

func renderDailyIssues(theme Theme, state ContentState, width, height int) string {
	active := state.Pane == "issues"
	cur := state.Cursors["issues"]
	var issues []apiIssue
	if state.DailySummary != nil {
		for _, issue := range state.DailySummary.Issues {
			issues = append(issues, newAPIIssue(issue.ID, issue.Title, issue.Status, issue.EstimateMinutes, issue.TodoForDate))
		}
	}
	indices := filteredIssueIndices(issues, state.Filters["issues"])
	total := len(indices)
	inner := height - 5
	if inner < 1 {
		inner = 1
	}
	lines := []string{theme.StylePaneTitle.Render("Planned Tasks"), renderFilterLine(theme, state.Filters["issues"], width-6)}
	if len(issues) == 0 || total == 0 {
		lines = append(lines, theme.StyleDim.Render("No planned tasks for this date"))
		return renderPaneBox(theme, active, width, height, stringsJoin(lines))
	}
	statusW := 11
	estimateW := 8
	repoW := max(10, width/7)
	streamW := max(10, width/7)
	titleW := width - statusW - estimateW - repoW - streamW - 16
	if titleW < 14 {
		titleW = 14
	}
	header := fmt.Sprintf("%-2s %-*s %-*s %-*s %-*s %-*s", "", titleW, "Issue", statusW, "Status", estimateW, "Estimate", repoW, "Repo", streamW, "Stream")
	lines = append(lines, theme.StyleDim.Render(truncate(header, width-6)))
	start, end := listWindow(cur, total, inner)
	if start > 0 {
		lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↑ %d more", start)))
	}
	for i := start; i < end; i++ {
		issue := issues[indices[i]]
		meta := issueMetaByID(state.AllIssues, issue.ID)
		repoName, streamName := "-", "-"
		if meta != nil {
			repoName = meta.RepoName
			streamName = meta.StreamName
		}
		estimate := "-"
		if issue.EstimateMinutes != nil {
			estimate = fmt.Sprintf("%dm", *issue.EstimateMinutes)
		}
		title := issue.Title + issueDueSuffix(issue.TodoForDate)
		row := fmt.Sprintf("%-2s %-*s %-*s %-*s %-*s %-*s", "", titleW, truncate(title, titleW), statusW, truncate(plainIssueStatus(string(issue.Status)), statusW), estimateW, estimate, repoW, truncate(repoName, repoW), streamW, truncate(streamName, streamW))
		lines = append(lines, renderPaneRowStyled(theme, i, cur, active, row, issueStatusStyle(theme, string(issue.Status)), width))
	}
	if remaining := total - end; remaining > 0 {
		lines = append(lines, theme.StyleDim.Render(fmt.Sprintf("↓ %d more", remaining)))
	}
	return renderPaneBox(theme, active, width, height, stringsJoin(lines))
}
