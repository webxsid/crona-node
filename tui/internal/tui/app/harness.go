package app

import (
	"crona/tui/internal/api"

	"github.com/charmbracelet/lipgloss"
)

func NewDailyRenderModel(width, height int) Model {
	repoName := "Work"
	streamName := "app"
	estimate := 60
	target := 15

	return Model{
		view:   ViewDaily,
		pane:   PaneIssues,
		width:  width,
		height: height,
		cursor: map[Pane]int{PaneIssues: 0, PaneHabits: 0},
		filters: map[Pane]string{
			PaneIssues: "",
			PaneHabits: "",
		},
		context: &api.ActiveContext{
			RepoName:   &repoName,
			StreamName: &streamName,
		},
		kernelInfo: &api.KernelInfo{Env: "Dev"},
		dailySummary: &api.DailyIssueSummary{
			Date: "2026-03-19",
			Issues: []api.Issue{
				{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
			},
		},
		allIssues: []api.IssueWithMeta{
			{
				Issue:      api.Issue{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
				RepoName:   "Work",
				StreamName: "app",
			},
		},
		dueHabits: []api.HabitDailyItem{
			{
				HabitWithMeta: api.HabitWithMeta{
					Habit: api.Habit{Name: "Inbox Zero Sweep", TargetMinutes: &target},
				},
				Status: "pending",
			},
		},
	}
}

func (m Model) RenderString() string { return m.View() }
func (m Model) BodyHeight() int      { return lipgloss.Height(m.renderBody()) }
func (m Model) ContentHeight() int   { return m.contentHeight() }

func MinimumSize() (int, int) {
	return minTUIWidth, minTUIHeight
}

func NewDailyHabitDeleteModel(habits []api.HabitDailyItem) Model {
	return Model{
		view:      ViewDaily,
		pane:      PaneHabits,
		cursor:    map[Pane]int{PaneHabits: 0},
		filters:   map[Pane]string{PaneHabits: ""},
		dueHabits: habits,
		timer:     &api.TimerState{State: "idle"},
	}
}

func OpenSelectedDeleteDialog(m Model) (Model, bool) {
	return m.openSelectedDeleteDialog()
}

func (m Model) DialogDeleteKind() string { return m.dialogDeleteKind }
func (m Model) DialogDeleteID() string   { return m.dialogDeleteID }
func (m Model) DialogStreamID() int64    { return m.dialogStreamID }
