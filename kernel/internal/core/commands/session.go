package commands

import (
	"context"
	"errors"
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

func StopSession(ctx context.Context, c *core.Context, commitMessage *string) (*sharedtypes.Session, error) {
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
	newCommit := strings.TrimSpace(strings.Join([]string{existingCommit, valueOrEmpty(commitMessage)}, "\n"))
	if newCommit != "" {
		mergedCommit = &newCommit
	}

	workLines := workSummaryLines
	if existingWork := parsedExisting[sharedtypes.SessionNoteSectionWork]; existingWork != "" {
		workLines = append(workSummaryLines, "")
		workLines = append(workLines, strings.Split(existingWork, "\n")...)
	}

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
	}{
		Commit:      mergedCommit,
		RepoID:      repoID,
		StreamID:    streamID,
		IssueID:     issueID,
		WorkSummary: workLines,
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
		return nil, errors.New("no session found to ammend")
	}
	merged := sessionnotes.AmendCommitMessage(session.Notes, message)
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
