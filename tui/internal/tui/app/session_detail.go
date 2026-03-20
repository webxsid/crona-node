package app

import (
	"fmt"
	"strings"

	sharedtypes "crona/shared/types"
	helperpkg "crona/tui/internal/tui/app/helpers"

	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) sessionDetailContentLines() []string {
	if m.sessionDetail == nil {
		return []string{"Loading session detail...", "", "[esc] close"}
	}

	detail := m.sessionDetail
	ended := "-"
	if detail.EndTime != nil && strings.TrimSpace(*detail.EndTime) != "" {
		ended = *detail.EndTime
	}
	duration := helperpkg.FormatSessionDurationText(detail.DurationSeconds, detail.StartTime, detail.EndTime)
	lines := []string{
		fmt.Sprintf("Repo: %s", detail.RepoName),
		fmt.Sprintf("Stream: %s", detail.StreamName),
		fmt.Sprintf("Issue: #%d %s", detail.IssueID, detail.IssueTitle),
		fmt.Sprintf("Started: %s", detail.StartTime),
		fmt.Sprintf("Ended: %s", ended),
		fmt.Sprintf("Duration: %s", duration),
		"",
		fmt.Sprintf("Work: %s", helperpkg.FormatClockText(detail.WorkSummary.WorkSeconds)),
		fmt.Sprintf("Rest: %s", helperpkg.FormatClockText(detail.WorkSummary.RestSeconds)),
		fmt.Sprintf("Segments: %d work / %d rest", detail.WorkSummary.WorkSegments, detail.WorkSummary.RestSegments),
	}

	sectionOrder := []sharedtypes.SessionNoteSection{
		sharedtypes.SessionNoteSectionCommit,
		sharedtypes.SessionNoteSectionContext,
		sharedtypes.SessionNoteSectionWork,
		sharedtypes.SessionNoteSectionNotes,
	}
	labels := map[sharedtypes.SessionNoteSection]string{
		sharedtypes.SessionNoteSectionCommit:  "Commit",
		sharedtypes.SessionNoteSectionContext: "Context",
		sharedtypes.SessionNoteSectionWork:    "Work Summary",
		sharedtypes.SessionNoteSectionNotes:   "Notes",
	}
	for _, section := range sectionOrder {
		value := ""
		if detail.ParsedNotes != nil {
			value = strings.TrimSpace(detail.ParsedNotes[section])
		}
		if value == "" {
			continue
		}
		lines = append(lines, "", labels[section]+":")
		lines = append(lines, strings.Split(value, "\n")...)
	}
	return lines
}

func (m Model) sessionDetailViewportHeight() int {
	if m.height < 16 {
		return max(6, m.height-8)
	}
	return min(18, m.height-8)
}

func (m Model) sessionDetailMaxOffset() int {
	boxWidth := min(max(52, m.width-10), 96)
	innerWidth := boxWidth - 4
	lines := m.sessionDetailContentLines()
	wrapped := make([]string, 0, len(lines))
	for _, line := range lines {
		if line == "" {
			wrapped = append(wrapped, "")
			continue
		}
		wrapped = append(wrapped, wrapText(line, innerWidth)...)
	}
	return max(0, len(wrapped)-m.sessionDetailViewportHeight())
}

func (m Model) updateSessionDetailOverlay(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc", "q", "o", "enter":
		m.sessionDetailOpen = false
		m.sessionDetail = nil
		m.sessionDetailY = 0
		return m, nil
	case "j", "down":
		if m.sessionDetailY < m.sessionDetailMaxOffset() {
			m.sessionDetailY++
		}
		return m, nil
	case "k", "up":
		if m.sessionDetailY > 0 {
			m.sessionDetailY--
		}
		return m, nil
	case "e":
		if m.sessionDetail == nil {
			return m, nil
		}
		return m.openAmendSessionDialog(m.sessionDetail.ID, helperpkg.SessionCommit(m.sessionDetail)), nil
	}
	return m, nil
}
