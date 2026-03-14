package views

import (
	sharedtypes "crona/shared/types"
)

type ActionsState struct {
	View           string
	Pane           string
	ScratchpadOpen bool
	TimerState     string
	IsDevMode      bool
}

func PaneActions(theme Theme, state ActionsState) []string {
	base := []string{
		theme.StyleHeader.Render("[tab]") + theme.StyleDim.Render(" next pane"),
	}

	if state.View == "session_active" {
		if state.TimerState == "" || state.TimerState == "idle" {
			return []string{
				theme.StyleHeader.Render("[/]") + theme.StyleDim.Render(" filter"),
				theme.StyleHeader.Render("[f]") + theme.StyleDim.Render(" focus"),
				theme.StyleHeader.Render("[Z]") + theme.StyleDim.Render(" stashes"),
			}
		}
		return []string{
			theme.StyleHeader.Render("[p/r]") + theme.StyleDim.Render(" pause/resume"),
			theme.StyleHeader.Render("[x]") + theme.StyleDim.Render(" end"),
			theme.StyleHeader.Render("[z]") + theme.StyleDim.Render(" stash"),
			theme.StyleHeader.Render("[s/A]") + theme.StyleDim.Render(" issue"),
			theme.StyleHeader.Render("[[] []]") + theme.StyleDim.Render(" session/scratchpads"),
		}
	}
	if state.View == "session_history" {
		return []string{
			theme.StyleHeader.Render("[/]") + theme.StyleDim.Render(" filter"),
			theme.StyleHeader.Render("[enter]") + theme.StyleDim.Render(" details"),
			theme.StyleHeader.Render("[f]") + theme.StyleDim.Render(" focus"),
			theme.StyleHeader.Render("[Z]") + theme.StyleDim.Render(" stashes"),
		}
	}

	switch state.Pane {
	case "repos", "streams":
		return append(base,
			theme.StyleHeader.Render("[c]")+theme.StyleDim.Render(" checkout"),
			theme.StyleHeader.Render("[/]")+theme.StyleDim.Render(" filter"),
			theme.StyleHeader.Render("[a]")+theme.StyleDim.Render(" new"),
			theme.StyleHeader.Render("[e]")+theme.StyleDim.Render(" edit"),
			theme.StyleHeader.Render("[d]")+theme.StyleDim.Render(" delete"),
		)
	case "issues":
		if state.View == "daily" {
			return append(base,
				theme.StyleHeader.Render("[f]")+theme.StyleDim.Render(" focus"),
				theme.StyleHeader.Render("[s]")+theme.StyleDim.Render(" status"),
				theme.StyleHeader.Render("[A]")+theme.StyleDim.Render(" abandon"),
				theme.StyleHeader.Render("[y]")+theme.StyleDim.Render(" today"),
				theme.StyleHeader.Render("[D]")+theme.StyleDim.Render(" due date"),
				theme.StyleHeader.Render("[e/d]")+theme.StyleDim.Render(" edit/delete"),
				theme.StyleHeader.Render("[,/.]")+theme.StyleDim.Render(" date"),
				theme.StyleHeader.Render("[/]")+theme.StyleDim.Render(" filter"),
			)
		}
		return append(base,
			theme.StyleHeader.Render("[1/2]")+theme.StyleDim.Render(" active/completed"),
			theme.StyleHeader.Render("[f]")+theme.StyleDim.Render(" focus"),
			theme.StyleHeader.Render("[s]")+theme.StyleDim.Render(" status"),
			theme.StyleHeader.Render("[A]")+theme.StyleDim.Render(" abandon"),
			theme.StyleHeader.Render("[y]")+theme.StyleDim.Render(" today"),
			theme.StyleHeader.Render("[D]")+theme.StyleDim.Render(" due date"),
			theme.StyleHeader.Render("[e/d]")+theme.StyleDim.Render(" edit/delete"),
			theme.StyleHeader.Render("[/]")+theme.StyleDim.Render(" filter"),
			theme.StyleHeader.Render("[a]")+theme.StyleDim.Render(" new"),
		)
	case "scratchpads":
		if state.ScratchpadOpen {
			actions := append(base,
				theme.StyleHeader.Render("[h/l]")+theme.StyleDim.Render(" switch"),
				theme.StyleHeader.Render("[e]")+theme.StyleDim.Render(" edit"),
				theme.StyleHeader.Render("[esc]")+theme.StyleDim.Render(" close"),
			)
			if state.TimerState != "" && state.TimerState != "idle" {
				actions = append(actions, theme.StyleHeader.Render("[s/A]")+theme.StyleDim.Render(" issue"))
			}
			return actions
		}
		actions := append(base,
			theme.StyleHeader.Render("[enter]")+theme.StyleDim.Render(" open"),
			theme.StyleHeader.Render("[/]")+theme.StyleDim.Render(" filter"),
			theme.StyleHeader.Render("[a]")+theme.StyleDim.Render(" new"),
			theme.StyleHeader.Render("[d]")+theme.StyleDim.Render(" delete"),
		)
		if state.TimerState != "" && state.TimerState != "idle" {
			actions = append(actions, theme.StyleHeader.Render("[s/A]")+theme.StyleDim.Render(" issue"))
		}
		return actions
	case "ops":
		return append(base,
			theme.StyleHeader.Render("[+]")+theme.StyleDim.Render(" more"),
			theme.StyleHeader.Render("[-]")+theme.StyleDim.Render(" less"),
			theme.StyleHeader.Render("[/]")+theme.StyleDim.Render(" filter"),
		)
	case "settings":
		return append(base,
			theme.StyleHeader.Render("[h/l]")+theme.StyleDim.Render(" change"),
			theme.StyleHeader.Render("[enter]")+theme.StyleDim.Render(" toggle/advance"),
		)
	}
	return base
}

func SettingsItemLabels(settings *sharedtypes.CoreSettings) []string {
	if settings == nil {
		return nil
	}
	return []string{
		"Timer Mode",
		"Breaks Enabled",
		"Work Duration",
		"Short Break",
		"Long Break",
		"Long Break Enabled",
		"Cycles Before Long Break",
		"Auto Start Breaks",
		"Auto Start Work",
	}
}
