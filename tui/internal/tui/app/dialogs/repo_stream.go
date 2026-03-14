package dialogs

func renderRepoStreamDialog(theme Theme, state State) string {
	switch state.Kind {
	case "create_repo":
		return renderSingleInput(theme, state.Width, "New Repo", "Name", state.Inputs, theme.ColorCyan, "[enter] create   [esc] cancel")
	case "edit_repo":
		return renderSingleInput(theme, state.Width, "Edit Repo", "Name", state.Inputs, theme.ColorYellow, "[enter] save   [esc] cancel")
	case "create_stream":
		rows := []string{theme.StylePaneTitle.Render("New Stream"), "", theme.StyleDim.Render("Repo"), theme.StyleHeader.Render(state.RepoName), "", theme.StyleDim.Render("Name"), state.Inputs[0].View(), "", theme.StyleDim.Render("[enter] create   [esc] cancel")}
		return modal(theme, state.Width, 56, theme.ColorCyan, rows)
	case "edit_stream":
		rows := []string{theme.StylePaneTitle.Render("Edit Stream"), "", theme.StyleDim.Render("Repo"), theme.StyleHeader.Render(state.RepoName), "", theme.StyleDim.Render("Name"), state.Inputs[0].View(), "", theme.StyleDim.Render("[enter] save   [esc] cancel")}
		return modal(theme, state.Width, 56, theme.ColorYellow, rows)
	default:
		return ""
	}
}
