package app

import (
	"fmt"
	"strings"
	"time"

	sharedtypes "crona/shared/types"
	"crona/tui/internal/api"
)

func issueDueLabel(todoForDate *string) string {
	if todoForDate == nil {
		return ""
	}
	date := strings.TrimSpace(*todoForDate)
	if date == "" {
		return ""
	}
	if date == time.Now().Format("2006-01-02") {
		return "today"
	}
	return "due " + date
}

func deref(s *string) string {
	if s == nil {
		return "-"
	}
	return *s
}

func firstNonEmpty(a, b *string) string {
	if a != nil && *a != "" {
		return *a
	}
	return deref(b)
}

func truncate(s string, max int) string {
	if max < 4 {
		max = 4
	}
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max-3]) + "..."
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func sessionHistorySummary(entry api.SessionHistoryEntry) string {
	if entry.ParsedNotes != nil {
		if message := strings.TrimSpace(entry.ParsedNotes[sharedtypes.SessionNoteSectionCommit]); message != "" {
			return message
		}
		if note := strings.TrimSpace(entry.ParsedNotes[sharedtypes.SessionNoteSectionNotes]); note != "" {
			return note
		}
	}
	if entry.Notes != nil && strings.TrimSpace(*entry.Notes) != "" {
		return strings.TrimSpace(*entry.Notes)
	}
	return fmt.Sprintf("Issue #%d", entry.IssueID)
}
