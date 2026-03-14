package dialogs

import (
	sharedtypes "crona/shared/types"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
)

type Theme struct {
	ColorCyan   lipgloss.Color
	ColorYellow lipgloss.Color
	ColorRed    lipgloss.Color
	ColorGreen  lipgloss.Color

	StylePaneTitle lipgloss.Style
	StyleDim       lipgloss.Style
	StyleCursor    lipgloss.Style
	StyleHeader    lipgloss.Style
	StyleError     lipgloss.Style
	StyleSelected  lipgloss.Style
	StyleNormal    lipgloss.Style
}

type State struct {
	Kind                string
	Width               int
	Inputs              []textinput.Model
	FocusIdx            int
	DeleteID            string
	DeleteKind          string
	DeleteLabel         string
	SessionID           string
	IssueID             int64
	StashCursor         int
	Stashes             []StashItem
	RepoID              int64
	StreamID            int64
	StatusItems         []sharedtypes.IssueStatus
	StatusCursor        int
	StatusLabel         string
	StatusRequired      bool
	IssueStatus         string
	RepoName            string
	StreamName          string
	RepoIndex           int
	StreamIndex         int
	Parent              string
	DateMonthValue      string
	DateCursorValue     string
	RepoSelectorLabel   string
	StreamSelectorLabel string
	DateTitle           string
	DateHeader          string
	DateMonth           string
	DateGrid            string
}

type StashItem struct {
	Label string
	Meta  string
}

func Render(theme Theme, state State) string {
	switch state.Kind {
	case "create_repo", "edit_repo", "create_stream", "edit_stream":
		return renderRepoStreamDialog(theme, state)
	case "create_issue_meta", "create_issue_default", "edit_issue", "issue_status", "issue_status_note":
		return renderIssueDialog(theme, state)
	case "end_session", "stash_session", "issue_session_transition", "stash_list", "amend_session":
		return renderSessionDialog(theme, state)
	case "confirm_delete", "pick_date", "create_scratchpad":
		return renderUtilityDialog(theme, state)
	default:
		return ""
	}
}
