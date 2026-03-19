package testsuite

import (
	"fmt"
	"strings"
	"testing"

	"crona/tui/internal/api"
	"crona/tui/internal/tui/app"
	"crona/tui/internal/tui/app/views"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

func TestPaneActionLineWrapsInsteadOfDroppingActions(t *testing.T) {
	rendered := views.RenderPaneActionLine(
		views.TestingTheme(),
		"",
		20,
		[]string{"[enter] view", "[a] new", "[c] context"},
	)
	lines := strings.Split(rendered, "\n")
	if len(lines) < 2 {
		t.Fatalf("expected wrapped action lines, got %q", rendered)
	}
	if !strings.Contains(rendered, "[c] context") {
		t.Fatalf("expected final action to be preserved, got %q", rendered)
	}
	for _, line := range lines {
		if got := lipgloss.Width(line); got > 20 {
			t.Fatalf("line width %d exceeds max width 20: %q", got, line)
		}
	}
}

func TestPaneBoxClipsOverflowingContent(t *testing.T) {
	rendered := views.RenderPaneBoxForTesting(views.TestingTheme(), true, 40, 8, strings.Join([]string{
		"line1",
		"line2",
		"line3",
		"line4",
		"line5",
		"line6",
		"line7",
		"line8",
	}, "\n"))
	if got := lipgloss.Height(rendered); got > 8 {
		t.Fatalf("pane box height %d exceeds allotted height 8", got)
	}
}

func TestDailyViewStacksOnNarrowWidths(t *testing.T) {
	state := views.ContentState{
		View:   "daily",
		Pane:   "issues",
		Width:  32,
		Height: 22,
		Cursors: map[string]int{
			"issues": 0,
			"habits": 0,
		},
		Filters: map[string]string{
			"issues": "",
			"habits": "",
		},
	}

	rendered := views.RenderDailyForTesting(views.TestingTheme(), state)
	for _, line := range strings.Split(rendered, "\n") {
		if got := lipgloss.Width(line); got > state.Width {
			t.Fatalf("daily view line width %d exceeds content width %d: %q", got, state.Width, line)
		}
	}
}

func TestDailyViewDoesNotExceedAllocatedHeight(t *testing.T) {
	state := views.ContentState{
		View:   "daily",
		Pane:   "issues",
		Width:  80,
		Height: 14,
		Cursors: map[string]int{
			"issues": 0,
			"habits": 0,
		},
		Filters: map[string]string{
			"issues": "",
			"habits": "",
		},
	}

	rendered := views.RenderDailyForTesting(views.TestingTheme(), state)
	if got := lipgloss.Height(rendered); got > state.Height {
		t.Fatalf("daily view height %d exceeds allocated height %d", got, state.Height)
	}
}

func TestDailyViewReportedHeightRangeFitsAllocation(t *testing.T) {
	estimate := 60
	target := 15
	state := views.ContentState{
		View:   "daily",
		Pane:   "issues",
		Width:  70,
		Height: 43,
		Cursors: map[string]int{
			"issues": 0,
			"habits": 0,
		},
		Filters: map[string]string{
			"issues": "",
			"habits": "",
		},
		DailySummary: &api.DailyIssueSummary{
			Date: "2026-03-19",
			Issues: []api.Issue{
				{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
			},
		},
		DailyIssues: []api.Issue{
			{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
		},
		AllIssues: []api.IssueWithMeta{
			{Issue: api.Issue{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate}, RepoName: "Work", StreamName: "app"},
		},
		DueHabits: []api.HabitDailyItem{
			{HabitWithMeta: api.HabitWithMeta{Habit: api.Habit{Name: "Inbox Zero Sweep", TargetMinutes: &target}}, Status: "pending"},
		},
		Context: &api.ActiveContext{
			RepoName:   strPtr("Work"),
			StreamName: strPtr("app"),
		},
	}

	rendered := views.RenderDailyForTesting(views.TestingTheme(), state)
	if got := lipgloss.Height(rendered); got > state.Height {
		t.Fatalf("daily view height %d exceeds allocated height %d", got, state.Height)
	}
}

func TestDailyViewDoesNotExceedTerminalHeightInReportedRange(t *testing.T) {
	for height := 46; height <= 54; height++ {
		model := app.NewDailyTestModel(92, height)
		if got, want := model.BodyHeightForTesting(), model.ContentHeightForTesting(); got > want {
			t.Fatalf("daily body height %d exceeds content height %d at terminal height %d", got, want, height)
		}
		rendered := model.RenderForTesting()
		if got := lipgloss.Height(rendered); got > height {
			t.Fatalf("daily view height %d exceeds terminal height %d", got, height)
		}
	}
}

func TestDailySummaryUsesCompactInlineModeBelowHeight55(t *testing.T) {
	estimate := 60
	target := 15
	state := views.ContentState{
		View:   "daily",
		Pane:   "issues",
		Width:  70,
		Height: 54,
		Cursors: map[string]int{
			"issues": 0,
			"habits": 0,
		},
		Filters: map[string]string{
			"issues": "",
			"habits": "",
		},
		DailySummary: &api.DailyIssueSummary{
			Date: "2026-03-19",
			Issues: []api.Issue{
				{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
			},
		},
		DailyIssues: []api.Issue{
			{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
		},
		DueHabits: []api.HabitDailyItem{
			{HabitWithMeta: api.HabitWithMeta{Habit: api.Habit{Name: "Inbox Zero Sweep", TargetMinutes: &target}}, Status: "pending"},
		},
		Context: &api.ActiveContext{
			RepoName:   strPtr("Work"),
			StreamName: strPtr("app"),
		},
	}

	rendered := views.RenderDailyForTesting(views.TestingTheme(), state)
	if !strings.Contains(rendered, "Issues  0/1 resolved") {
		t.Fatalf("expected compact inline issues row below height 55")
	}
	if !strings.Contains(rendered, "Habits  0/1 completed") {
		t.Fatalf("expected compact inline habits row below height 55")
	}
	if !strings.Contains(rendered, "planned 1") {
		t.Fatalf("expected compact legend text below height 55")
	}
	if !strings.Contains(rendered, "logged 0m / target 15m") {
		t.Fatalf("expected compact habit meta below height 55")
	}
	if !strings.Contains(rendered, "█") {
		t.Fatalf("expected inline bars below height 55")
	}
}

func TestDailySummaryShowsBarsAtHeight55AndAbove(t *testing.T) {
	estimate := 60
	target := 15
	state := views.ContentState{
		View:   "daily",
		Pane:   "issues",
		Width:  70,
		Height: 55,
		Cursors: map[string]int{
			"issues": 0,
			"habits": 0,
		},
		Filters: map[string]string{
			"issues": "",
			"habits": "",
		},
		DailySummary: &api.DailyIssueSummary{
			Date: "2026-03-19",
			Issues: []api.Issue{
				{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
			},
		},
		DailyIssues: []api.Issue{
			{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
		},
		DueHabits: []api.HabitDailyItem{
			{HabitWithMeta: api.HabitWithMeta{Habit: api.Habit{Name: "Inbox Zero Sweep", TargetMinutes: &target}}, Status: "pending"},
		},
		Context: &api.ActiveContext{
			RepoName:   strPtr("Work"),
			StreamName: strPtr("app"),
		},
	}

	rendered := views.RenderDailyForTesting(views.TestingTheme(), state)
	if !strings.Contains(rendered, "████") {
		t.Fatalf("expected bars to remain visible at height 55 and above")
	}
}

func TestDailySummaryUsesUltraCompactModeBelowHeight48(t *testing.T) {
	estimate := 60
	target := 15
	state := views.ContentState{
		View:   "daily",
		Pane:   "issues",
		Width:  70,
		Height: 46,
		Cursors: map[string]int{
			"issues": 0,
			"habits": 0,
		},
		Filters: map[string]string{
			"issues": "",
			"habits": "",
		},
		DailySummary: &api.DailyIssueSummary{
			Date: "2026-03-19",
			Issues: []api.Issue{
				{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
			},
		},
		DailyIssues: []api.Issue{
			{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
		},
		DueHabits: []api.HabitDailyItem{
			{HabitWithMeta: api.HabitWithMeta{Habit: api.Habit{Name: "Inbox Zero Sweep", TargetMinutes: &target}}, Status: "pending"},
		},
		Context: &api.ActiveContext{
			RepoName:   strPtr("Work"),
			StreamName: strPtr("app"),
		},
	}

	rendered := views.RenderDailyForTesting(views.TestingTheme(), state)
	if !strings.Contains(rendered, "Issues  0/1 resolved") || !strings.Contains(rendered, "Habits  0/1 completed") {
		t.Fatalf("expected ultra-compact rows for both issues and habits")
	}
	if strings.Contains(rendered, "planned 1") {
		t.Fatalf("expected issue legend row to be omitted below height 48")
	}
	if strings.Contains(rendered, "failed 0   remaining 1") {
		t.Fatalf("expected habit meta row to be omitted below height 48")
	}
	if !strings.Contains(rendered, "█") {
		t.Fatalf("expected inline bars to remain in ultra-compact mode")
	}
}

func TestDailySummaryUsesTinyHeightModeAt36(t *testing.T) {
	estimate := 60
	target := 15
	state := views.ContentState{
		View:   "daily",
		Pane:   "issues",
		Width:  70,
		Height: 36,
		Cursors: map[string]int{
			"issues": 0,
			"habits": 0,
		},
		Filters: map[string]string{
			"issues": "",
			"habits": "",
		},
		DailySummary: &api.DailyIssueSummary{
			Date: "2026-03-19",
			Issues: []api.Issue{
				{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
			},
		},
		DailyIssues: []api.Issue{
			{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
		},
		DueHabits: []api.HabitDailyItem{
			{HabitWithMeta: api.HabitWithMeta{Habit: api.Habit{Name: "Inbox Zero Sweep", TargetMinutes: &target}}, Status: "pending"},
		},
		Context: &api.ActiveContext{
			RepoName:   strPtr("Work"),
			StreamName: strPtr("app"),
		},
	}

	rendered := views.RenderDailyForTesting(views.TestingTheme(), state)
	assertTinySummary(t, rendered)
}

func TestDailySummaryUsesTinyHeightModeAt30(t *testing.T) {
	estimate := 60
	target := 15
	state := views.ContentState{
		View:   "daily",
		Pane:   "issues",
		Width:  70,
		Height: 30,
		Cursors: map[string]int{
			"issues": 0,
			"habits": 0,
		},
		Filters: map[string]string{
			"issues": "",
			"habits": "",
		},
		DailySummary: &api.DailyIssueSummary{
			Date: "2026-03-19",
			Issues: []api.Issue{
				{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
			},
		},
		DailyIssues: []api.Issue{
			{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate},
		},
		DueHabits: []api.HabitDailyItem{
			{HabitWithMeta: api.HabitWithMeta{Habit: api.Habit{Name: "Inbox Zero Sweep", TargetMinutes: &target}}, Status: "pending"},
		},
		Context: &api.ActiveContext{
			RepoName:   strPtr("Work"),
			StreamName: strPtr("app"),
		},
	}

	rendered := views.RenderDailyForTesting(views.TestingTheme(), state)
	assertTinySummary(t, rendered)
}

func TestDefaultViewUsesCompactModeAt36(t *testing.T) {
	state := compactDefaultState(36)
	rendered := views.RenderDefaultForTesting(views.TestingTheme(), state)
	assertCompactDefault(t, rendered, state.Height)
}

func TestDefaultViewUsesCompactModeAt30(t *testing.T) {
	state := compactDefaultState(30)
	rendered := views.RenderDefaultForTesting(views.TestingTheme(), state)
	assertCompactDefault(t, rendered, state.Height)
}

func TestWellbeingViewUsesCompactModeAt36(t *testing.T) {
	state := compactWellbeingState(36)
	rendered := views.RenderWellbeingForTesting(views.TestingTheme(), state)
	assertCompactWellbeing(t, rendered, state.Height)
}

func TestWellbeingViewUsesCompactModeAt30(t *testing.T) {
	state := compactWellbeingState(30)
	rendered := views.RenderWellbeingForTesting(views.TestingTheme(), state)
	assertCompactWellbeing(t, rendered, state.Height)
}

func TestUndersizedWidthShowsMinimumSizeWarning(t *testing.T) {
	minWidth, minHeight := app.MinimumSizeForTesting()
	model := app.NewDailyTestModel(minWidth-1, minHeight)
	rendered := model.RenderForTesting()
	assertMinimumSizeWarning(t, rendered, minWidth-1, minHeight, minWidth, minHeight)
}

func TestUndersizedHeightShowsMinimumSizeWarning(t *testing.T) {
	minWidth, minHeight := app.MinimumSizeForTesting()
	model := app.NewDailyTestModel(minWidth, minHeight-1)
	rendered := model.RenderForTesting()
	assertMinimumSizeWarning(t, rendered, minWidth, minHeight-1, minWidth, minHeight)
}

func TestUndersizedBothDimensionsShowMinimumSizeWarning(t *testing.T) {
	minWidth, minHeight := app.MinimumSizeForTesting()
	model := app.NewDailyTestModel(minWidth-5, minHeight-2)
	rendered := model.RenderForTesting()
	assertMinimumSizeWarning(t, rendered, minWidth-5, minHeight-2, minWidth, minHeight)
}

func TestMinimumSizeThresholdRendersNormalUI(t *testing.T) {
	minWidth, minHeight := app.MinimumSizeForTesting()
	model := app.NewDailyTestModel(minWidth, minHeight)
	rendered := model.RenderForTesting()
	if strings.Contains(rendered, "Terminal Too Small") {
		t.Fatalf("expected normal UI at minimum size")
	}
	if !strings.Contains(rendered, "Daily Dashboard") {
		t.Fatalf("expected daily UI at minimum size")
	}
	if got := lipgloss.Height(rendered); got > minHeight {
		t.Fatalf("rendered height %d exceeds terminal height %d", got, minHeight)
	}
}

func TestAboveMinimumSizeRendersNormalUI(t *testing.T) {
	minWidth, minHeight := app.MinimumSizeForTesting()
	model := app.NewDailyTestModel(minWidth+1, minHeight+1)
	rendered := model.RenderForTesting()
	if strings.Contains(rendered, "Terminal Too Small") {
		t.Fatalf("expected normal UI above minimum size")
	}
	if !strings.Contains(rendered, "Daily Dashboard") {
		t.Fatalf("expected daily UI above minimum size")
	}
}

func assertMinimumSizeWarning(t *testing.T, rendered string, currentWidth, currentHeight, minWidth, minHeight int) {
	t.Helper()
	if !strings.Contains(rendered, "Terminal Too Small") {
		t.Fatalf("expected undersized warning, got %q", rendered)
	}
	if !strings.Contains(rendered, fmt.Sprintf("Current: %dx%d", currentWidth, currentHeight)) {
		t.Fatalf("expected current dimensions in warning")
	}
	if !strings.Contains(rendered, fmt.Sprintf("Required: %dx%d", minWidth, minHeight)) {
		t.Fatalf("expected required dimensions in warning")
	}
	if !strings.Contains(rendered, "Resize the terminal to continue.") {
		t.Fatalf("expected resize instruction in warning")
	}
	if got := lipgloss.Width(rendered); got > currentWidth {
		t.Fatalf("warning width %d exceeds viewport width %d", got, currentWidth)
	}
	if got := lipgloss.Height(rendered); got > currentHeight {
		t.Fatalf("warning height %d exceeds viewport height %d", got, currentHeight)
	}
	if strings.Contains(rendered, "Daily Dashboard") {
		t.Fatalf("expected normal UI to be suppressed while undersized")
	}
}

func strPtr(v string) *string { return &v }

func assertTinySummary(t *testing.T, rendered string) {
	t.Helper()
	if !strings.Contains(rendered, "Daily Dashboard") {
		t.Fatalf("expected dashboard title in tiny summary")
	}
	if !strings.Contains(rendered, "2026-03-19") {
		t.Fatalf("expected date in tiny summary")
	}
	if !strings.Contains(rendered, "Scope: Work > app") {
		t.Fatalf("expected scope in tiny summary")
	}
	if !strings.Contains(rendered, "[,] [.] [g]") {
		t.Fatalf("expected compact date hints in tiny summary")
	}
	if !strings.Contains(rendered, "Issues  0/1") || !strings.Contains(rendered, "Habits  0/1") {
		t.Fatalf("expected both issue and habit summary rows in tiny summary")
	}
	if !strings.Contains(rendered, "p1") {
		t.Fatalf("expected abbreviated issue legend in tiny summary")
	}
	if !strings.Contains(rendered, "f0 r1") {
		t.Fatalf("expected abbreviated habit tail in tiny summary")
	}
	if !strings.Contains(rendered, "█") {
		t.Fatalf("expected micro-bars in tiny summary")
	}
}

func compactDefaultState(height int) views.ContentState {
	estimate1, estimate2, estimate3 := 60, 35, 25
	today := "2026-03-19"
	return views.ContentState{
		View:   "default",
		Pane:   "issues",
		Width:  92,
		Height: height,
		Cursors: map[string]int{
			"issues": 0,
		},
		Filters: map[string]string{
			"issues": "",
		},
		DefaultIssueSection: "open",
		DefaultIssues: []api.IssueWithMeta{
			{Issue: api.Issue{ID: 1, Title: "Add keyboard-first command palette", Status: "planned", EstimateMinutes: &estimate1, TodoForDate: &today}, RepoName: "Work", StreamName: "app"},
			{Issue: api.Issue{ID: 2, Title: "Improve install docs for Linux", Status: "planned", EstimateMinutes: &estimate2, TodoForDate: &today}, RepoName: "OSS", StreamName: "cli"},
			{Issue: api.Issue{ID: 3, Title: "Research standing desk options", Status: "abandoned", EstimateMinutes: &estimate3}, RepoName: "Personal", StreamName: "home"},
		},
		Context: &api.ActiveContext{},
	}
}

func compactWellbeingState(height int) views.ContentState {
	avgMood, avgEnergy := 4.0, 3.7
	return views.ContentState{
		View:   "wellbeing",
		Pane:   "issues",
		Width:  92,
		Height: height,
		WellbeingDate: "2026-03-19",
		MetricsRollup: &api.MetricsRollup{
			Days:          7,
			CheckInDays:   6,
			FocusDays:     1,
			WorkedSeconds: 4956,
			RestSeconds:   2,
			AverageMood:   &avgMood,
			AverageEnergy: &avgEnergy,
			LatestBurnout: &api.BurnoutIndicator{
				Level:  "low",
				Score:  31,
				Factors: map[string]float64{"breakCompliance": 0.99},
			},
		},
		Streaks: &api.StreakSummary{
			CurrentCheckInDays: 0,
			LongestCheckInDays: 6,
			CurrentFocusDays:   0,
			LongestFocusDays:   1,
		},
	}
}

func assertCompactDefault(t *testing.T, rendered string, height int) {
	t.Helper()
	plain := ansi.Strip(rendered)
	if !strings.Contains(plain, "Due 2") || !strings.Contains(plain, "Open 2") || !strings.Contains(plain, "Closed 1") {
		t.Fatalf("expected compact stats header in default view")
	}
	if !strings.Contains(plain, "Active Issues [1]") {
		t.Fatalf("expected primary issue list in compact default view")
	}
	if !strings.Contains(plain, "Closed") {
		t.Fatalf("expected compact completed footer in default view")
	}
	if !strings.Contains(plain, "Add keyboard-first command palette") {
		t.Fatalf("expected open issue rows in compact default view")
	}
	if got := lipgloss.Height(rendered); got > height {
		t.Fatalf("default compact view height %d exceeds allocated height %d", got, height)
	}
}

func assertCompactWellbeing(t *testing.T, rendered string, height int) {
	t.Helper()
	plain := ansi.Strip(rendered)
	if !strings.Contains(plain, "Wellbeing") || !strings.Contains(plain, "2026-03-19") {
		t.Fatalf("expected compact wellbeing header")
	}
	if !strings.Contains(plain, "[,/.]") && !strings.Contains(plain, "[a/e]") {
		t.Fatalf("expected action hints in compact wellbeing view")
	}
	if !strings.Contains(plain, "No check-in recorded for this date") {
		t.Fatalf("expected current day summary in compact wellbeing view")
	}
	if !strings.Contains(plain, "Burnout") || !strings.Contains(plain, "31 LOW") {
		t.Fatalf("expected burnout summary in compact wellbeing view")
	}
	if !strings.Contains(plain, "Metrics Window") || !strings.Contains(plain, "Days  7") {
		t.Fatalf("expected compact metrics block in wellbeing view")
	}
	if got := lipgloss.Height(rendered); got > height {
		t.Fatalf("wellbeing compact view height %d exceeds allocated height %d", got, height)
	}
}
