package dialogs

import (
	shareddto "crona/shared/dto"
	sharedtypes "crona/shared/types"
	"strconv"
	"strings"
	"time"

	"crona/tui/internal/api"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type UpdateContext struct {
	Repos             []api.Repo
	Streams           []api.Stream
	AllIssues         []api.IssueWithMeta
	Context           *api.ActiveContext
	Stashes           []api.Stash
	SelectedIssueID   int64
	SelectedStreamID  int64
	HasSelectedIssue  bool
	ActiveIssueStream int64
	HasActiveIssue    bool
}

type Action struct {
	Kind       string
	ID         string
	RepoID     int64
	StreamID   int64
	IssueID    int64
	Name       string
	Path       string
	RepoName   string
	StreamName string
	Title      string
	Status     string
	Estimate   *int
	DueDate    *string
	Note       *string
	Payload    shareddto.EndSessionRequest
}

func OpenCreateScratchpad(state State) State {
	name := textinput.New()
	name.Placeholder = "My notes"
	name.Focus()
	name.CharLimit = 80
	name.Width = 40
	path := textinput.New()
	path.Placeholder = "notes/[[date]].md"
	path.CharLimit = 120
	path.Width = 40
	state = Close(state)
	state.Kind = "create_scratchpad"
	state.Inputs = []textinput.Model{name, path}
	return state
}

func OpenCreateRepo(state State) State {
	name := textinput.New()
	name.Placeholder = "Repo name"
	name.Focus()
	name.CharLimit = 80
	name.Width = 40
	state = Close(state)
	state.Kind = "create_repo"
	state.Inputs = []textinput.Model{name}
	return state
}

func OpenEditRepo(state State, repoID int64, name string) State {
	input := textinput.New()
	input.Placeholder = "Repo name"
	input.SetValue(name)
	input.Focus()
	input.CharLimit = 80
	input.Width = 40
	state = Close(state)
	state.Kind = "edit_repo"
	state.Inputs = []textinput.Model{input}
	state.RepoID = repoID
	return state
}

func OpenCreateStream(state State, repoID int64, repoName string) State {
	input := textinput.New()
	input.Placeholder = "Stream name"
	input.Focus()
	input.CharLimit = 80
	input.Width = 40
	state = Close(state)
	state.Kind = "create_stream"
	state.Inputs = []textinput.Model{input}
	state.RepoID = repoID
	state.RepoName = repoName
	return state
}

func OpenEditStream(state State, streamID, repoID int64, streamName, repoName string) State {
	input := textinput.New()
	input.Placeholder = "Stream name"
	input.SetValue(streamName)
	input.Focus()
	input.CharLimit = 80
	input.Width = 40
	state = Close(state)
	state.Kind = "edit_stream"
	state.Inputs = []textinput.Model{input}
	state.StreamID = streamID
	state.RepoID = repoID
	state.RepoName = repoName
	return state
}

func OpenCreateIssueMeta(state State, streamID int64, streamName, repoName string) State {
	title := textinput.New()
	title.Placeholder = "Issue title"
	title.Focus()
	title.CharLimit = 120
	title.Width = 52
	estimate := textinput.New()
	estimate.Placeholder = "Estimate minutes (optional)"
	estimate.CharLimit = 6
	estimate.Width = 20
	due := textinput.New()
	due.Placeholder = "Due date YYYY-MM-DD (optional)"
	due.CharLimit = 10
	due.Width = 22
	state = Close(state)
	state.Kind = "create_issue_meta"
	state.Inputs = []textinput.Model{title, estimate, due}
	state.StreamID = streamID
	state.StreamName = streamName
	state.RepoName = repoName
	return state
}

func OpenEditIssue(state State, issueID, streamID int64, title string, estimateMinutes *int, todoForDate *string) State {
	titleInput := textinput.New()
	titleInput.Placeholder = "Issue title"
	titleInput.SetValue(title)
	titleInput.Focus()
	titleInput.CharLimit = 120
	titleInput.Width = 52
	estimateInput := textinput.New()
	estimateInput.Placeholder = "Estimate minutes (optional)"
	estimateInput.CharLimit = 6
	estimateInput.Width = 20
	if estimateMinutes != nil {
		estimateInput.SetValue(strconv.Itoa(*estimateMinutes))
	}
	dueInput := textinput.New()
	dueInput.Placeholder = "Due date YYYY-MM-DD (optional)"
	dueInput.CharLimit = 10
	dueInput.Width = 22
	if todoForDate != nil {
		dueInput.SetValue(strings.TrimSpace(*todoForDate))
	}
	state = Close(state)
	state.Kind = "edit_issue"
	state.Inputs = []textinput.Model{titleInput, estimateInput, dueInput}
	state.IssueID = issueID
	state.StreamID = streamID
	return state
}

func OpenCreateIssueDefault(state State) State {
	repoFilter := textinput.New()
	repoFilter.Placeholder = "Search repo"
	repoFilter.CharLimit = 80
	repoFilter.Width = 36
	repoFilter.Focus()
	streamFilter := textinput.New()
	streamFilter.Placeholder = "Search stream"
	streamFilter.CharLimit = 80
	streamFilter.Width = 36
	title := textinput.New()
	title.Placeholder = "Issue title"
	title.CharLimit = 120
	title.Width = 52
	estimate := textinput.New()
	estimate.Placeholder = "Estimate minutes (optional)"
	estimate.CharLimit = 6
	estimate.Width = 20
	due := textinput.New()
	due.Placeholder = "Due date YYYY-MM-DD (optional)"
	due.CharLimit = 10
	due.Width = 22
	state = Close(state)
	state.Kind = "create_issue_default"
	state.Inputs = []textinput.Model{repoFilter, streamFilter, title, estimate, due}
	return state
}

func OpenConfirmDelete(state State, kind, id, label string, repoID, streamID int64) State {
	state = Close(state)
	state.Kind = "confirm_delete"
	state.DeleteKind = kind
	state.DeleteID = id
	state.DeleteLabel = label
	state.RepoID = repoID
	state.StreamID = streamID
	return state
}

func OpenStashList(state State) State {
	state = Close(state)
	state.Kind = "stash_list"
	return state
}

func OpenIssueStatus(state State, status string) State {
	state = Close(state)
	state.Kind = "issue_status"
	state.StatusItems = sharedtypes.AllowedIssueStatusTransitions(sharedtypes.IssueStatus(status))
	return state
}

func OpenIssueStatusNote(state State, issueID, streamID int64, status, label string, required bool) State {
	input := textinput.New()
	input.Placeholder = label
	input.Focus()
	input.CharLimit = 200
	input.Width = 48
	state = Close(state)
	state.Kind = "issue_status_note"
	state.Inputs = []textinput.Model{input}
	state.IssueID = issueID
	state.StreamID = streamID
	state.IssueStatus = status
	state.StatusLabel = label
	state.StatusRequired = required
	return state
}

func OpenSessionMessage(state State, kind string) State {
	input := textinput.New()
	if kind == "end_session" {
		input.Placeholder = "Commit message"
	} else {
		input.Placeholder = "Stash note"
	}
	input.Focus()
	input.CharLimit = 200
	input.Width = 48
	state = Close(state)
	state.Kind = kind
	state.Inputs = []textinput.Model{input}
	return state
}

func OpenIssueSessionTransition(state State, issueID int64, status string) State {
	state = Close(state)
	state.Kind = "issue_session_transition"
	state.IssueID = issueID
	state.IssueStatus = status
	switch status {
	case "done":
		input := textinput.New()
		input.Placeholder = "Completion note (optional)"
		input.Focus()
		input.CharLimit = 200
		input.Width = 48
		state.Inputs = []textinput.Model{input}
	case "abandoned":
		input := textinput.New()
		input.Placeholder = "Abandon reason"
		input.Focus()
		input.CharLimit = 200
		input.Width = 48
		state.Inputs = []textinput.Model{input}
	}
	return state
}

func OpenAmendSession(state State, sessionID string, commit string) State {
	input := textinput.New()
	input.Placeholder = "Commit message"
	input.SetValue(strings.TrimSpace(commit))
	input.Focus()
	input.CharLimit = 200
	input.Width = 48
	state = Close(state)
	state.Kind = "amend_session"
	state.SessionID = sessionID
	state.Inputs = []textinput.Model{input}
	return state
}

func OpenDatePicker(state State, parentDialog string, issueID int64, inputIndex int, initial *string, currentDate string) State {
	selected := ResolveDialogDate(initial, currentDate)
	monthStart := time.Date(selected.Year(), selected.Month(), 1, 0, 0, 0, 0, selected.Location())
	state.Parent = parentDialog
	state.IssueID = issueID
	state.Kind = "pick_date"
	state.DateCursorValue = selected.Format("2006-01-02")
	state.DateMonthValue = monthStart.Format("2006-01-02")
	state.FocusIdx = inputIndex
	return state
}

func Close(state State) State {
	state.Kind = ""
	state.Inputs = nil
	state.FocusIdx = 0
	state.DeleteID = ""
	state.DeleteKind = ""
	state.DeleteLabel = ""
	state.SessionID = ""
	state.IssueID = 0
	state.IssueStatus = ""
	state.RepoID = 0
	state.RepoName = ""
	state.StreamID = 0
	state.StreamName = ""
	state.RepoIndex = 0
	state.StreamIndex = 0
	state.Parent = ""
	state.DateMonthValue = ""
	state.DateCursorValue = ""
	state.StashCursor = 0
	state.StatusItems = nil
	state.StatusCursor = 0
	state.StatusLabel = ""
	state.StatusRequired = false
	return state
}

func ToggleEndSessionAdvanced(state State) State {
	if state.Kind != "end_session" {
		return state
	}
	if len(state.Inputs) > 1 {
		commit := state.Inputs[0].Value()
		input := newSessionDetailInput("Commit message")
		input.SetValue(commit)
		input.Focus()
		state.Inputs = []textinput.Model{input}
		state.FocusIdx = 0
		return state
	}
	commit := state.Inputs[0].Value()
	inputs := []textinput.Model{
		newSessionDetailInput("Commit message"),
		newSessionDetailInput("Worked on"),
		newSessionDetailInput("Outcome"),
		newSessionDetailInput("Next step"),
		newSessionDetailInput("Blockers"),
		newSessionDetailInput("Links"),
	}
	inputs[0].SetValue(commit)
	inputs[0].Focus()
	state.Inputs = inputs
	state.FocusIdx = 0
	return state
}

func Update(state State, ctx UpdateContext, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	switch state.Kind {
	case "create_repo":
		return updateSingleInput(state, msg, "Repo name is required", func(name string) *Action {
			return &Action{Kind: "create_repo", Name: name}
		})
	case "edit_repo":
		return updateSingleInput(state, msg, "Repo name is required", func(name string) *Action {
			return &Action{Kind: "edit_repo", RepoID: state.RepoID, Name: name}
		})
	case "create_stream":
		return updateSingleInput(state, msg, "Stream name is required", func(name string) *Action {
			return &Action{Kind: "create_stream", RepoID: state.RepoID, Name: name}
		})
	case "edit_stream":
		return updateSingleInput(state, msg, "Stream name is required", func(name string) *Action {
			return &Action{Kind: "edit_stream", RepoID: state.RepoID, StreamID: state.StreamID, Name: name}
		})
	case "create_scratchpad":
		return updateCreateScratchpad(state, msg)
	case "confirm_delete":
		return updateConfirmDelete(state, msg)
	case "stash_list":
		return updateStashList(state, ctx, msg)
	case "issue_status":
		return updateIssueStatus(state, ctx, currentDate, msg)
	case "issue_status_note":
		return updateIssueStatusNote(state, currentDate, msg)
	case "end_session", "stash_session":
		return updateSessionMessage(state, ctx, currentDate, msg)
	case "amend_session":
		return updateAmendSession(state, msg)
	case "issue_session_transition":
		return updateIssueSessionTransition(state, ctx, currentDate, msg)
	case "pick_date":
		return updateDatePicker(state, ctx, currentDate, msg)
	case "create_issue_meta":
		return updateCreateIssueMeta(state, currentDate, msg)
	case "create_issue_default":
		return updateCreateIssueDefault(state, ctx, currentDate, msg)
	case "edit_issue":
		return updateEditIssue(state, currentDate, msg)
	default:
		return state, nil, ""
	}
}

func newSessionDetailInput(placeholder string) textinput.Model {
	input := textinput.New()
	input.Placeholder = placeholder
	input.CharLimit = 200
	input.Width = 48
	return input
}

func updateSingleInput(state State, msg tea.KeyMsg, requiredMsg string, submit func(string) *Action) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "enter":
		name := strings.TrimSpace(state.Inputs[0].Value())
		if name == "" {
			return state, nil, requiredMsg
		}
		return Close(state), submit(name), ""
	}
	var cmd tea.Cmd
	state.Inputs[0], cmd = state.Inputs[0].Update(msg)
	_ = cmd
	return state, nil, ""
}

func updateCreateScratchpad(state State, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "tab", "shift+tab", "down", "up":
		dir := 1
		if msg.String() == "shift+tab" || msg.String() == "up" {
			dir = -1
		}
		state.FocusIdx = (state.FocusIdx + dir + len(state.Inputs)) % len(state.Inputs)
		state.Inputs = SyncFocus(state.Inputs, state.FocusIdx)
		return state, nil, ""
	case "enter":
		name := strings.TrimSpace(state.Inputs[0].Value())
		path := strings.TrimSpace(state.Inputs[1].Value())
		if name == "" || path == "" {
			return state, nil, "Name and path are required"
		}
		return Close(state), &Action{Kind: "create_scratchpad", Name: name, Path: path}, ""
	}
	var cmd tea.Cmd
	state.Inputs[state.FocusIdx], cmd = state.Inputs[state.FocusIdx].Update(msg)
	_ = cmd
	return state, nil, ""
}

func updateConfirmDelete(state State, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "enter":
		action := &Action{Kind: "delete", ID: state.DeleteID, RepoID: state.RepoID, StreamID: state.StreamID}
		action.Name = state.DeleteKind
		return Close(state), action, ""
	}
	return state, nil, ""
}

func updateStashList(state State, ctx UpdateContext, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc", "q":
		return Close(state), nil, ""
	case "j", "down":
		if state.StashCursor < len(ctx.Stashes)-1 {
			state.StashCursor++
		}
	case "k", "up":
		if state.StashCursor > 0 {
			state.StashCursor--
		}
	case "enter":
		if len(ctx.Stashes) == 0 || state.StashCursor < 0 || state.StashCursor >= len(ctx.Stashes) {
			return state, nil, ""
		}
		return Close(state), &Action{Kind: "apply_stash", ID: ctx.Stashes[state.StashCursor].ID}, ""
	case "x":
		if len(ctx.Stashes) == 0 || state.StashCursor < 0 || state.StashCursor >= len(ctx.Stashes) {
			return state, nil, ""
		}
		return Close(state), &Action{Kind: "drop_stash", ID: ctx.Stashes[state.StashCursor].ID}, ""
	}
	return state, nil, ""
}

func updateIssueStatus(state State, ctx UpdateContext, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc", "q":
		return Close(state), nil, ""
	case "j", "down":
		if state.StatusCursor < len(state.StatusItems)-1 {
			state.StatusCursor++
		}
	case "k", "up":
		if state.StatusCursor > 0 {
			state.StatusCursor--
		}
	case "enter":
		if !ctx.HasSelectedIssue || len(state.StatusItems) == 0 || state.StatusCursor < 0 || state.StatusCursor >= len(state.StatusItems) {
			return state, nil, ""
		}
		status := string(state.StatusItems[state.StatusCursor])
		switch status {
		case "blocked":
			return OpenIssueStatusNote(state, ctx.SelectedIssueID, ctx.SelectedStreamID, status, "Blocker reason", true), nil, ""
		case "in_review":
			return OpenIssueStatusNote(state, ctx.SelectedIssueID, ctx.SelectedStreamID, status, "Review note (optional)", false), nil, ""
		case "done":
			return OpenIssueStatusNote(state, ctx.SelectedIssueID, ctx.SelectedStreamID, status, "Completion note (optional)", false), nil, ""
		case "abandoned":
			return OpenIssueStatusNote(state, ctx.SelectedIssueID, ctx.SelectedStreamID, status, "Abandon reason", true), nil, ""
		default:
			return Close(state), &Action{Kind: "change_issue_status", IssueID: ctx.SelectedIssueID, StreamID: ctx.SelectedStreamID, Status: status}, ""
		}
	}
	return state, nil, ""
}

func updateIssueStatusNote(state State, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "enter":
		note := ValueToPointer(state.Inputs[0].Value())
		if state.StatusRequired && note == nil {
			return state, nil, state.StatusLabel + " is required"
		}
		return Close(state), &Action{Kind: "change_issue_status", IssueID: state.IssueID, StreamID: state.StreamID, Status: state.IssueStatus, Note: note}, ""
	}
	var cmd tea.Cmd
	state.Inputs[0], cmd = state.Inputs[0].Update(msg)
	_ = cmd
	return state, nil, ""
}

func updateSessionMessage(state State, ctx UpdateContext, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "ctrl+e", "f2":
		return ToggleEndSessionAdvanced(state), nil, ""
	case "tab", "shift+tab", "down", "up":
		if len(state.Inputs) == 1 {
			return state, nil, ""
		}
		dir := 1
		if msg.String() == "shift+tab" || msg.String() == "up" {
			dir = -1
		}
		state.FocusIdx = (state.FocusIdx + dir + len(state.Inputs)) % len(state.Inputs)
		state.Inputs = SyncFocus(state.Inputs, state.FocusIdx)
		return state, nil, ""
	case "enter":
		payload := EndSessionRequest(state.Inputs)
		kind := state.Kind
		state = Close(state)
		if kind == "end_session" {
			if !ctx.HasActiveIssue {
				return state, nil, "Active issue metadata unavailable"
			}
			return state, &Action{Kind: "end_session", StreamID: ctx.ActiveIssueStream, Payload: payload}, ""
		}
		return state, &Action{Kind: "stash_session", Note: payload.CommitMessage}, ""
	}
	var cmd tea.Cmd
	state.Inputs[state.FocusIdx], cmd = state.Inputs[state.FocusIdx].Update(msg)
	_ = cmd
	return state, nil, ""
}

func updateIssueSessionTransition(state State, ctx UpdateContext, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "enter":
		note := ValueToPointer("")
		if len(state.Inputs) > 0 {
			note = ValueToPointer(state.Inputs[0].Value())
		}
		if state.IssueStatus == "abandoned" && note == nil {
			return state, nil, "Abandon reason is required"
		}
		return Close(state), &Action{Kind: "change_issue_status_and_end_session", IssueID: state.IssueID, StreamID: ctx.ActiveIssueStream, Status: state.IssueStatus, Note: note, Payload: shareddto.EndSessionRequest{CommitMessage: note}}, ""
	case "n", "N":
		if state.IssueStatus != "done" && state.IssueStatus != "abandoned" {
			return Close(state), nil, ""
		}
	case "y", "Y":
		if state.IssueStatus != "done" && state.IssueStatus != "abandoned" {
			return Close(state), &Action{Kind: "change_issue_status_and_end_session", IssueID: state.IssueID, StreamID: ctx.ActiveIssueStream, Status: state.IssueStatus, Payload: shareddto.EndSessionRequest{}}, ""
		}
	}
	if (state.IssueStatus == "done" || state.IssueStatus == "abandoned") && len(state.Inputs) > 0 {
		var cmd tea.Cmd
		state.Inputs[0], cmd = state.Inputs[0].Update(msg)
		_ = cmd
	}
	return state, nil, ""
}

func updateAmendSession(state State, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "enter":
		note := strings.TrimSpace(state.Inputs[0].Value())
		if note == "" {
			return state, nil, "Commit message is required"
		}
		return Close(state), &Action{Kind: "amend_session", ID: state.SessionID, Note: ValueToPointer(note)}, ""
	}
	var cmd tea.Cmd
	state.Inputs[0], cmd = state.Inputs[0].Update(msg)
	_ = cmd
	return state, nil, ""
}

func updateCreateIssueMeta(state State, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	if state.FocusIdx == 2 && (msg.String() == "f2" || msg.String() == "ctrl+y") {
		return OpenDatePicker(state, "create_issue_meta", 0, 2, ValueToPointer(state.Inputs[2].Value()), currentDate), nil, ""
	}
	return updateMultiInputIssue(state, msg, 3, func(state State) (*Action, string) {
		title := strings.TrimSpace(state.Inputs[0].Value())
		if title == "" {
			return nil, "Issue title is required"
		}
		estimate, err := ParseEstimateInput(state.Inputs[1].Value())
		if err != nil {
			return nil, err.Error()
		}
		dueDate, err := ParseDueDateInput(state.Inputs[2].Value())
		if err != nil {
			return nil, err.Error()
		}
		return &Action{Kind: "create_issue_meta", StreamID: state.StreamID, Title: title, Estimate: estimate, DueDate: dueDate}, ""
	})
}

func updateCreateIssueDefault(state State, ctx UpdateContext, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "f2", "ctrl+y":
		if state.FocusIdx == 4 {
			return OpenDatePicker(state, "create_issue_default", 0, 4, ValueToPointer(state.Inputs[4].Value()), currentDate), nil, ""
		}
	case "tab", "shift+tab", "down", "up":
		if (msg.String() == "down" || msg.String() == "up") && (state.FocusIdx == 0 || state.FocusIdx == 1) {
			if state.FocusIdx == 0 {
				state.RepoIndex = ShiftSelection(state.RepoIndex, len(DefaultRepoOptions(state.Inputs, ctx.Repos)), ternaryDir(msg.String()))
				state.StreamIndex = 0
			} else {
				state.StreamIndex = ShiftSelection(state.StreamIndex, len(DefaultStreamOptions(state.Inputs, state.RepoIndex, ctx.Repos, ctx.AllIssues, ctx.Streams, ctx.Context)), ternaryDir(msg.String()))
			}
			return state, nil, ""
		}
		state.FocusIdx = (state.FocusIdx + ternaryDir(msg.String()) + 5) % 5
		state.Inputs = SyncFocus(state.Inputs, state.FocusIdx)
		return state, nil, ""
	case "left":
		if state.FocusIdx == 0 {
			state.RepoIndex = ShiftSelection(state.RepoIndex, len(DefaultRepoOptions(state.Inputs, ctx.Repos)), -1)
			state.StreamIndex = 0
			return state, nil, ""
		}
		if state.FocusIdx == 1 {
			state.StreamIndex = ShiftSelection(state.StreamIndex, len(DefaultStreamOptions(state.Inputs, state.RepoIndex, ctx.Repos, ctx.AllIssues, ctx.Streams, ctx.Context)), -1)
			return state, nil, ""
		}
	case "right":
		if state.FocusIdx == 0 {
			state.RepoIndex = ShiftSelection(state.RepoIndex, len(DefaultRepoOptions(state.Inputs, ctx.Repos)), 1)
			state.StreamIndex = 0
			return state, nil, ""
		}
		if state.FocusIdx == 1 {
			state.StreamIndex = ShiftSelection(state.StreamIndex, len(DefaultStreamOptions(state.Inputs, state.RepoIndex, ctx.Repos, ctx.AllIssues, ctx.Streams, ctx.Context)), 1)
			return state, nil, ""
		}
	case "enter":
		repoName, streamName := DefaultIssueDialogNames(state.Inputs, state.RepoIndex, state.StreamIndex, ctx.Repos, ctx.AllIssues, ctx.Streams, ctx.Context)
		title := strings.TrimSpace(state.Inputs[2].Value())
		if repoName == "" || streamName == "" || title == "" {
			return state, nil, "Repo, stream, and issue title are required"
		}
		estimate, err := ParseEstimateInput(state.Inputs[3].Value())
		if err != nil {
			return state, nil, err.Error()
		}
		dueDate, err := ParseDueDateInput(state.Inputs[4].Value())
		if err != nil {
			return state, nil, err.Error()
		}
		return Close(state), &Action{Kind: "create_issue_default", RepoName: repoName, StreamName: streamName, Title: title, Estimate: estimate, DueDate: dueDate}, ""
	}
	var cmd tea.Cmd
	state.Inputs[state.FocusIdx], cmd = state.Inputs[state.FocusIdx].Update(msg)
	_ = cmd
	if state.FocusIdx == 0 {
		state.RepoIndex = 0
		state.StreamIndex = 0
	}
	if state.FocusIdx == 1 {
		state.StreamIndex = 0
	}
	return state, nil, ""
}

func updateEditIssue(state State, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	if state.FocusIdx == 2 && (msg.String() == "f2" || msg.String() == "ctrl+y") {
		return OpenDatePicker(state, "edit_issue", state.IssueID, 2, ValueToPointer(state.Inputs[2].Value()), currentDate), nil, ""
	}
	return updateMultiInputIssue(state, msg, 3, func(state State) (*Action, string) {
		title := strings.TrimSpace(state.Inputs[0].Value())
		if title == "" {
			return nil, "Issue title is required"
		}
		estimate, err := ParseEstimateInput(state.Inputs[1].Value())
		if err != nil {
			return nil, err.Error()
		}
		dueDate, err := ParseDueDateInput(state.Inputs[2].Value())
		if err != nil {
			return nil, err.Error()
		}
		return &Action{Kind: "edit_issue", IssueID: state.IssueID, StreamID: state.StreamID, Title: title, Estimate: estimate, DueDate: dueDate}, ""
	})
}

func updateMultiInputIssue(state State, msg tea.KeyMsg, inputCount int, submit func(State) (*Action, string)) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return Close(state), nil, ""
	case "tab", "shift+tab", "down", "up":
		state.FocusIdx = (state.FocusIdx + ternaryDir(msg.String()) + inputCount) % inputCount
		state.Inputs = SyncFocus(state.Inputs, state.FocusIdx)
		return state, nil, ""
	case "enter":
		action, status := submit(state)
		if action == nil {
			return state, nil, status
		}
		return Close(state), action, status
	}
	var cmd tea.Cmd
	state.Inputs[state.FocusIdx], cmd = state.Inputs[state.FocusIdx].Update(msg)
	_ = cmd
	return state, nil, ""
}

func updateDatePicker(state State, ctx UpdateContext, currentDate string, msg tea.KeyMsg) (State, *Action, string) {
	switch msg.String() {
	case "esc":
		return closeDatePicker(state), nil, ""
	case "enter", " ":
		selected := state.DateCursorValue
		if state.Parent == "create_issue_meta" || state.Parent == "create_issue_default" || state.Parent == "edit_issue" {
			if idx := state.FocusIdx; idx >= 0 && idx < len(state.Inputs) {
				state.Inputs[idx].SetValue(selected)
			}
			return closeDatePicker(state), nil, ""
		}
		return Close(state), &Action{Kind: "set_issue_todo_date", IssueID: state.IssueID, StreamID: ctx.ActiveIssueStream, DueDate: ValueToPointer(selected)}, ""
	case "backspace", "delete", "c":
		if state.Parent == "create_issue_meta" || state.Parent == "create_issue_default" || state.Parent == "edit_issue" {
			if idx := state.FocusIdx; idx >= 0 && idx < len(state.Inputs) {
				state.Inputs[idx].SetValue("")
			}
			return closeDatePicker(state), nil, ""
		}
		return Close(state), &Action{Kind: "set_issue_todo_date", IssueID: state.IssueID, StreamID: ctx.ActiveIssueStream, DueDate: ValueToPointer("")}, ""
	case "left", "h":
		return shiftDatePicker(state, 0, 0, -1), nil, ""
	case "right", "l":
		return shiftDatePicker(state, 0, 0, 1), nil, ""
	case "up", "k":
		return shiftDatePicker(state, 0, 0, -7), nil, ""
	case "down", "j":
		return shiftDatePicker(state, 0, 0, 7), nil, ""
	case ",":
		return shiftDatePicker(state, 0, -1, 0), nil, ""
	case ".":
		return shiftDatePicker(state, 0, 1, 0), nil, ""
	case "g":
		return OpenDatePicker(state, state.Parent, state.IssueID, state.FocusIdx, ValueToPointer(time.Now().Format("2006-01-02")), currentDate), nil, ""
	}
	return state, nil, ""
}

func ternaryDir(key string) int {
	if key == "shift+tab" || key == "up" {
		return -1
	}
	return 1
}

func closeDatePicker(state State) State {
	parent := state.Parent
	dateField := state.FocusIdx
	state.Kind = parent
	state.Parent = ""
	state.DateMonthValue = ""
	state.DateCursorValue = ""
	if parent == "create_issue_meta" || parent == "create_issue_default" {
		state.FocusIdx = dateField
		state.Inputs = SyncFocus(state.Inputs, state.FocusIdx)
		return state
	}
	for i := range state.Inputs {
		if i == state.FocusIdx {
			state.Inputs[i].Focus()
		} else {
			state.Inputs[i].Blur()
		}
	}
	return state
}

func shiftDatePicker(state State, years, months, days int) State {
	selected := DialogDate(state, time.Now().Format("2006-01-02")).AddDate(years, months, days)
	monthStart := time.Date(selected.Year(), selected.Month(), 1, 0, 0, 0, 0, selected.Location())
	state.DateCursorValue = selected.Format("2006-01-02")
	state.DateMonthValue = monthStart.Format("2006-01-02")
	return state
}
