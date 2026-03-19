package app

import (
	"crona/tui/internal/api"

	"github.com/charmbracelet/lipgloss"
)

// NewDailyTestModel returns a model seeded with enough daily-view state for
// isolated rendering tests outside this package.
func NewDailyTestModel(width, height int) Model {
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

func (m Model) RenderForTesting() string {
	return m.View()
}

func (m Model) BodyHeightForTesting() int {
	return lipgloss.Height(m.renderBody())
}

func (m Model) ContentHeightForTesting() int {
	return m.contentHeight()
}

func MinimumSizeForTesting() (int, int) {
	return minTUIWidth, minTUIHeight
}
