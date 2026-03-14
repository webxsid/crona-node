package commands

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"crona/kernel/internal/core"
	"crona/kernel/internal/sessionnotes"

	"github.com/google/uuid"

	sharedtypes "crona/shared/types"
)

func StartSession(ctx context.Context, c *core.Context, issueID int64) (sharedtypes.Session, error) {
	existing, err := c.Sessions.GetActiveSession(ctx, c.UserID)
	if err != nil {
		return sharedtypes.Session{}, err
	}
	if existing != nil {
		return sharedtypes.Session{}, errors.New("a session is already running")
	}
	now := c.Now()
	session := sharedtypes.Session{
		ID:        uuid.NewString(),
		IssueID:   issueID,
		StartTime: now,
	}
	created, err := c.Sessions.Start(ctx, session, c.UserID, c.DeviceID, now)
	if err != nil {
		return sharedtypes.Session{}, err
	}
	if _, err := c.SessionSegments.StartSegment(ctx, c.UserID, c.DeviceID, created.ID, sharedtypes.SessionSegmentWork); err != nil {
		return sharedtypes.Session{}, err
	}
	if err := c.Ops.Append(ctx, sharedtypes.Op{
		ID:        uuid.NewString(),
		Entity:    sharedtypes.OpEntitySession,
		EntityID:  created.ID,
		Action:    sharedtypes.OpActionCreate,
		Payload:   created,
		Timestamp: now,
		UserID:    c.UserID,
		DeviceID:  c.DeviceID,
	}); err != nil {
		return sharedtypes.Session{}, err
	}
	emit(c, sharedtypes.EventTypeSessionStarted, created)
	return created, nil
}

type SessionEndInput struct {
	CommitMessage *string
	WorkedOn      *string
	Outcome       *string
	NextStep      *string
	Blockers      *string
	Links         *string
}

func StopSession(ctx context.Context, c *core.Context, input SessionEndInput) (*sharedtypes.Session, error) {
	active, err := c.Sessions.GetActiveSession(ctx, c.UserID)
	if err != nil {
		return nil, err
	}
	if active == nil {
		return nil, nil
	}
	now := c.Now()
	if err := c.SessionSegments.EndActiveSegment(ctx, c.UserID, c.DeviceID, active.ID); err != nil {
		return nil, err
	}
	segments, err := c.SessionSegments.ListBySession(ctx, active.ID)
	if err != nil {
		return nil, err
	}
	workSummary := sessionnotes.ComputeWorkSummary(segments)
	workSummaryLines := sessionnotes.FormatWorkSummary(workSummary)
	parsedExisting := sessionnotes.Parse(active.Notes)

	activeContext, err := c.ActiveContext.Get(ctx, c.UserID, c.DeviceID)
	if err != nil {
		return nil, err
	}

	var mergedCommit *string
	existingCommit := parsedExisting[sharedtypes.SessionNoteSectionCommit]
	newCommit := strings.TrimSpace(strings.Join([]string{existingCommit, valueOrEmpty(input.CommitMessage)}, "\n"))
	if newCommit != "" {
		mergedCommit = &newCommit
	}

	workLines := workSummaryLines
	if existingWork := parsedExisting[sharedtypes.SessionNoteSectionWork]; existingWork != "" {
		workLines = append(workSummaryLines, "")
		workLines = append(workLines, strings.Split(existingWork, "\n")...)
	}

	userNotes := mergeSessionNotes(parsedExisting[sharedtypes.SessionNoteSectionNotes], formatSessionDetailNotes(input))

	var repoID, streamID, issueID *int64
	if activeContext != nil {
		repoID = activeContext.RepoID
		streamID = activeContext.StreamID
	}
	issueID = &active.IssueID

	notes := sessionnotes.GenerateDefaultSessionNotes(struct {
		Commit      *string
		RepoID      *int64
		StreamID    *int64
		IssueID     *int64
		WorkSummary []string
		Notes       *string
	}{
		Commit:      mergedCommit,
		RepoID:      repoID,
		StreamID:    streamID,
		IssueID:     issueID,
		WorkSummary: workLines,
		Notes:       userNotes,
	})
	notesPtr := &notes
	if err := sessionnotes.AssertCommitMessage(notesPtr); err != nil {
		return nil, err
	}

	offsetSeconds := 0
	for _, segment := range segments {
		if segment.ElapsedOffsetSeconds != nil {
			offsetSeconds += *segment.ElapsedOffsetSeconds
		}
	}

	stopped, err := c.Sessions.Stop(ctx, active.ID, struct {
		EndTime         string
		DurationSeconds int
		Notes           *string
	}{
		EndTime:         now,
		DurationSeconds: elapsedSeconds(active.StartTime, now) + offsetSeconds,
		Notes:           notesPtr,
	}, c.UserID, c.DeviceID, now)
	if err != nil {
		return nil, err
	}
	if err := c.Ops.Append(ctx, sharedtypes.Op{
		ID:        uuid.NewString(),
		Entity:    sharedtypes.OpEntitySession,
		EntityID:  stopped.ID,
		Action:    sharedtypes.OpActionUpdate,
		Payload:   map[string]any{"endTime": stopped.EndTime},
		Timestamp: now,
		UserID:    c.UserID,
		DeviceID:  c.DeviceID,
	}); err != nil {
		return nil, err
	}
	emit(c, sharedtypes.EventTypeSessionStopped, stopped)
	return stopped, nil
}

func mergeSessionNotes(existing string, next *string) *string {
	if strings.TrimSpace(existing) == "" && (next == nil || strings.TrimSpace(*next) == "") {
		return nil
	}
	if next == nil || strings.TrimSpace(*next) == "" {
		trimmed := strings.TrimSpace(existing)
		return &trimmed
	}
	if strings.TrimSpace(existing) == "" {
		trimmed := strings.TrimSpace(*next)
		return &trimmed
	}
	merged := strings.TrimSpace(existing) + "\n\n" + strings.TrimSpace(*next)
	return &merged
}

func formatSessionDetailNotes(input SessionEndInput) *string {
	lines := []string{}
	appendDetail := func(label string, value *string) {
		if value == nil || strings.TrimSpace(*value) == "" {
			return
		}
		lines = append(lines, fmt.Sprintf("%s: %s", label, strings.TrimSpace(*value)))
	}

	appendDetail("Worked on", input.WorkedOn)
	appendDetail("Outcome", input.Outcome)
	appendDetail("Next step", input.NextStep)
	appendDetail("Blockers", input.Blockers)
	appendDetail("Links", input.Links)

	if len(lines) == 0 {
		return nil
	}
	joined := strings.Join(lines, "\n")
	return &joined
}

func AmendSessionNotes(ctx context.Context, c *core.Context, message string, sessionID *string) (*sharedtypes.Session, error) {
	var (
		session *sharedtypes.Session
		err     error
	)
	if sessionID != nil && *sessionID != "" {
		session, err = c.Sessions.GetByID(ctx, *sessionID, c.UserID)
	} else {
		session, err = c.Sessions.GetLastSessionForUser(ctx, c.UserID)
	}
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, errors.New("no session found to amend")
	}
	if session.EndTime == nil || strings.TrimSpace(*session.EndTime) == "" {
		return nil, errors.New("cannot amend an active session")
	}
	if strings.TrimSpace(message) == "" {
		return nil, errors.New("commit message is required")
	}
	merged := sessionnotes.AmendCommitMessage(session.Notes, message)
	mergedPtr := &merged
	if err := sessionnotes.AssertCommitMessage(mergedPtr); err != nil {
		return nil, err
	}
	now := c.Now()
	updated, err := c.Sessions.AmendSessionNotes(ctx, session.ID, merged, c.UserID, c.DeviceID, now)
	if err != nil {
		return nil, err
	}
	if err := c.Ops.Append(ctx, sharedtypes.Op{
		ID:        uuid.NewString(),
		Entity:    sharedtypes.OpEntitySession,
		EntityID:  updated.ID,
		Action:    sharedtypes.OpActionUpdate,
		Payload:   map[string]any{"notes": updated.Notes},
		Timestamp: now,
		UserID:    c.UserID,
		DeviceID:  c.DeviceID,
	}); err != nil {
		return nil, err
	}
	return updated, nil
}

func GetSessionDetail(ctx context.Context, c *core.Context, sessionID string) (*sharedtypes.SessionDetail, error) {
	detail, err := c.Sessions.GetDetail(ctx, sessionID, c.UserID)
	if err != nil {
		return nil, err
	}
	if detail == nil {
		return nil, nil
	}
	segments, err := c.SessionSegments.ListBySession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	detail.ParsedNotes = sessionnotes.Parse(detail.Notes)
	detail.WorkSummary = sessionnotes.ComputeWorkSummary(segments)
	return detail, nil
}

func PauseSession(ctx context.Context, c *core.Context, nextSegmentType sharedtypes.SessionSegmentType) error {
	active, err := c.Sessions.GetActiveSession(ctx, c.UserID)
	if err != nil || active == nil {
		return err
	}
	current, err := c.SessionSegments.GetActive(ctx, c.UserID, c.DeviceID, active.ID)
	if err != nil {
		return err
	}
	if current != nil && current.SegmentType == sharedtypes.SessionSegmentRest {
		return nil
	}
	_, err = c.SessionSegments.StartSegment(ctx, c.UserID, c.DeviceID, active.ID, nextSegmentType)
	return err
}

func ResumeSession(ctx context.Context, c *core.Context) error {
	active, err := c.Sessions.GetActiveSession(ctx, c.UserID)
	if err != nil || active == nil {
		return err
	}
	current, err := c.SessionSegments.GetActive(ctx, c.UserID, c.DeviceID, active.ID)
	if err != nil {
		return err
	}
	if current != nil && current.SegmentType == sharedtypes.SessionSegmentWork {
		return nil
	}
	_, err = c.SessionSegments.StartSegment(ctx, c.UserID, c.DeviceID, active.ID, sharedtypes.SessionSegmentWork)
	return err
}

func ListSessionHistory(ctx context.Context, c *core.Context, query struct {
	RepoID   *int64
	StreamID *int64
	IssueID  *int64
	Since    *string
	Until    *string
	Limit    *int
	Offset   *int
}, useContext bool) ([]sharedtypes.SessionHistoryEntry, error) {
	if useContext {
		activeContext, err := c.ActiveContext.Get(ctx, c.UserID, c.DeviceID)
		if err != nil {
			return nil, err
		}
		if activeContext != nil {
			if activeContext.RepoID != nil {
				query.RepoID = activeContext.RepoID
			}
			if activeContext.StreamID != nil {
				query.StreamID = activeContext.StreamID
			}
			if activeContext.IssueID != nil {
				query.IssueID = activeContext.IssueID
			}
		}
	}
	if query.Limit == nil {
		limit := 100
		query.Limit = &limit
	}
	return c.Sessions.ListEnded(ctx, struct {
		UserID   string
		RepoID   *int64
		StreamID *int64
		IssueID  *int64
		Since    *string
		Until    *string
		Limit    *int
		Offset   *int
	}{
		UserID:   c.UserID,
		RepoID:   query.RepoID,
		StreamID: query.StreamID,
		IssueID:  query.IssueID,
		Since:    query.Since,
		Until:    query.Until,
		Limit:    query.Limit,
		Offset:   query.Offset,
	})
}

func valueOrEmpty(value *string) string {
	if value == nil {
		return ""
	}
	return *value
}
