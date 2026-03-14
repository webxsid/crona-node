package dialogs

func renderIssueDialog(theme Theme, state State) string {
	switch state.Kind {
	case "create_issue_meta":
		rows := []string{theme.StylePaneTitle.Render("New Issue"), "", theme.StyleDim.Render("Repo / Stream"), theme.StyleHeader.Render(state.RepoName + " / " + state.StreamName), "", theme.StyleDim.Render("Title"), state.Inputs[0].View(), "", theme.StyleDim.Render("Estimate"), state.Inputs[1].View(), "", theme.StyleDim.Render("Due"), state.Inputs[2].View(), "", theme.StyleDim.Render("[f2] calendar   [tab] next   [enter] create   [esc] cancel")}
		return modal(theme, state.Width, 68, theme.ColorCyan, rows)
	case "create_issue_default":
		rows := []string{theme.StylePaneTitle.Render("New Issue"), "", theme.StyleDim.Render("Repo"), state.Inputs[0].View(), "", renderSelector(theme, state.RepoSelectorLabel, true), "", theme.StyleDim.Render("Stream"), state.Inputs[1].View(), "", renderSelector(theme, state.StreamSelectorLabel, false), "", theme.StyleDim.Render("Title"), state.Inputs[2].View(), "", theme.StyleDim.Render("Estimate"), state.Inputs[3].View(), "", theme.StyleDim.Render("Due"), state.Inputs[4].View(), "", theme.StyleDim.Render("[type] filter   [up/down] choose   [f2] calendar   [tab] next   [enter] create")}
		return modal(theme, state.Width, 72, theme.ColorCyan, rows)
	case "edit_issue":
		rows := []string{theme.StylePaneTitle.Render("Edit Issue"), "", theme.StyleDim.Render("Title"), state.Inputs[0].View(), "", theme.StyleDim.Render("Estimate"), state.Inputs[1].View(), "", theme.StyleDim.Render("Due"), state.Inputs[2].View(), "", theme.StyleDim.Render("[f2] calendar   [tab] next   [enter] save   [esc] cancel")}
		return modal(theme, state.Width, 68, theme.ColorYellow, rows)
	case "issue_status":
		rows := []string{theme.StylePaneTitle.Render("Set Issue Status"), ""}
		if len(state.StatusItems) == 0 {
			rows = append(rows, theme.StyleDim.Render("No valid status transitions"), "", theme.StyleDim.Render("[esc] close"))
		} else {
			for i, status := range state.StatusItems {
				label := plainIssueStatus(string(status))
				if i == state.StatusCursor {
					rows = append(rows, theme.StyleSelected.Render("> "+label))
				} else {
					rows = append(rows, "  "+label)
				}
			}
			rows = append(rows, "", theme.StyleDim.Render("[j/k] move   [enter] set   [esc] cancel"))
		}
		return modal(theme, state.Width, 48, theme.ColorYellow, rows)
	case "issue_status_note":
		title := map[string]string{"blocked": "Block Issue", "in_review": "Send To Review", "done": "Complete Issue", "abandoned": "Abandon Issue"}[state.IssueStatus]
		if title == "" {
			title = "Status Note"
		}
		rows := []string{theme.StylePaneTitle.Render(title), "", theme.StyleDim.Render(state.StatusLabel), state.Inputs[0].View(), "", theme.StyleDim.Render("[enter] set   [esc] cancel")}
		return modal(theme, state.Width, 60, theme.ColorYellow, rows)
	default:
		return ""
	}
}
