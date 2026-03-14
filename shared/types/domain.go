package types

// Shared domain and wire types used across the Go workspace.

type IssueStatus string

const (
	IssueStatusBacklog    IssueStatus = "backlog"
	IssueStatusPlanned    IssueStatus = "planned"
	IssueStatusReady      IssueStatus = "ready"
	IssueStatusInProgress IssueStatus = "in_progress"
	IssueStatusBlocked    IssueStatus = "blocked"
	IssueStatusInReview   IssueStatus = "in_review"
	IssueStatusDone       IssueStatus = "done"
	IssueStatusAbandoned  IssueStatus = "abandoned"
)

type StreamVisibility string

const (
	StreamVisibilityPersonal StreamVisibility = "personal"
	StreamVisibilityShared   StreamVisibility = "shared"
)

type SessionSegmentType string

const (
	SessionSegmentWork       SessionSegmentType = "work"
	SessionSegmentShortBreak SessionSegmentType = "short_break"
	SessionSegmentLongBreak  SessionSegmentType = "long_break"
	SessionSegmentRest       SessionSegmentType = "rest"
)

type TimerMode string

const (
	TimerModeStopwatch  TimerMode = "stopwatch"
	TimerModeStructured TimerMode = "structured"
)

type OpEntity string

const (
	OpEntityRepo           OpEntity = "repo"
	OpEntityStream         OpEntity = "stream"
	OpEntityIssue          OpEntity = "issue"
	OpEntitySession        OpEntity = "session"
	OpEntitySessionSegment OpEntity = "session_segment"
	OpEntityActiveContext  OpEntity = "active_context"
	OpEntityStash          OpEntity = "stash"
)

type OpAction string

const (
	OpActionCreate  OpAction = "create"
	OpActionUpdate  OpAction = "update"
	OpActionDelete  OpAction = "delete"
	OpActionRestore OpAction = "restore"
)

type CoreSettingsKey string

const (
	CoreSettingsKeyTimerMode             CoreSettingsKey = "timerMode"
	CoreSettingsKeyBreaksEnabled         CoreSettingsKey = "breaksEnabled"
	CoreSettingsKeyWorkDurationMinutes   CoreSettingsKey = "workDurationMinutes"
	CoreSettingsKeyShortBreakMinutes     CoreSettingsKey = "shortBreakMinutes"
	CoreSettingsKeyLongBreakMinutes      CoreSettingsKey = "longBreakMinutes"
	CoreSettingsKeyLongBreakEnabled      CoreSettingsKey = "longBreakEnabled"
	CoreSettingsKeyCyclesBeforeLongBreak CoreSettingsKey = "cyclesBeforeLongBreak"
	CoreSettingsKeyAutoStartBreaks       CoreSettingsKey = "autoStartBreaks"
	CoreSettingsKeyAutoStartWork         CoreSettingsKey = "autoStartWork"
)

type SessionNoteSection string

const (
	SessionNoteSectionCommit  SessionNoteSection = "commit"
	SessionNoteSectionContext SessionNoteSection = "context"
	SessionNoteSectionWork    SessionNoteSection = "work"
	SessionNoteSectionNotes   SessionNoteSection = "notes"
)

type Repo struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Color *string `json:"color,omitempty"`
}

type Stream struct {
	ID         int64            `json:"id"`
	RepoID     int64            `json:"repoId"`
	Name       string           `json:"name"`
	Visibility StreamVisibility `json:"visibility"`
}

type Issue struct {
	ID              int64       `json:"id"`
	StreamID        int64       `json:"streamId"`
	Title           string      `json:"title"`
	Status          IssueStatus `json:"status"`
	EstimateMinutes *int        `json:"estimateMinutes,omitempty"`
	Notes           *string     `json:"notes,omitempty"`
	TodoForDate     *string     `json:"todoForDate,omitempty"`
	CompletedAt     *string     `json:"completedAt,omitempty"`
	AbandonedAt     *string     `json:"abandonedAt,omitempty"`
}

type IssueWithMeta struct {
	Issue
	RepoID     int64  `json:"repoId"`
	RepoName   string `json:"repoName"`
	StreamName string `json:"streamName"`
}

type DailyIssueSummary struct {
	Date                  string  `json:"date"`
	TotalIssues           int     `json:"totalIssues"`
	Issues                []Issue `json:"issues"`
	TotalEstimatedMinutes int     `json:"totalEstimatedMinutes"`
	CompletedIssues       int     `json:"completedIssues"`
	AbandonedIssues       int     `json:"abandonedIssues"`
	WorkedSeconds         int     `json:"workedSeconds"`
}

type Session struct {
	ID              string  `json:"id"`
	IssueID         int64   `json:"issueId"`
	StartTime       string  `json:"startTime"`
	EndTime         *string `json:"endTime,omitempty"`
	DurationSeconds *int    `json:"durationSeconds,omitempty"`
	Notes           *string `json:"notes,omitempty"`
}

type ParsedSessionNotes map[SessionNoteSection]string

type SessionHistoryEntry struct {
	Session
	ParsedNotes ParsedSessionNotes `json:"parsedNotes,omitempty"`
}

type SessionDetail struct {
	SessionHistoryEntry
	RepoID      int64              `json:"repoId"`
	RepoName    string             `json:"repoName"`
	StreamID    int64              `json:"streamId"`
	StreamName  string             `json:"streamName"`
	IssueTitle  string             `json:"issueTitle"`
	WorkSummary SessionWorkSummary `json:"workSummary"`
}

type SessionWorkSummary struct {
	WorkSeconds  int `json:"workSeconds"`
	RestSeconds  int `json:"restSeconds"`
	WorkSegments int `json:"workSegments"`
	RestSegments int `json:"restSegments"`
	TotalSeconds int `json:"totalSeconds"`
}

type SessionSegment struct {
	ID                   string             `json:"id"`
	UserID               string             `json:"userId"`
	DeviceID             string             `json:"deviceId"`
	SessionID            string             `json:"sessionId"`
	SegmentType          SessionSegmentType `json:"segmentType"`
	StartTime            string             `json:"startTime"`
	EndTime              *string            `json:"endTime,omitempty"`
	ElapsedOffsetSeconds *int               `json:"elapsedOffsetSeconds,omitempty"`
	CreatedAt            string             `json:"createdAt"`
}

type ActiveContext struct {
	UserID     string  `json:"userId"`
	DeviceID   string  `json:"deviceId"`
	RepoID     *int64  `json:"repoId,omitempty"`
	RepoName   *string `json:"repoName,omitempty"`
	StreamID   *int64  `json:"streamId,omitempty"`
	StreamName *string `json:"streamName,omitempty"`
	IssueID    *int64  `json:"issueId,omitempty"`
	IssueTitle *string `json:"issueTitle,omitempty"`
	UpdatedAt  *string `json:"updatedAt,omitempty"`
}

type CoreSettings struct {
	UserID                string    `json:"userId"`
	DeviceID              string    `json:"deviceId"`
	TimerMode             TimerMode `json:"timerMode"`
	BreaksEnabled         bool      `json:"breaksEnabled"`
	WorkDurationMinutes   int       `json:"workDurationMinutes"`
	ShortBreakMinutes     int       `json:"shortBreakMinutes"`
	LongBreakMinutes      int       `json:"longBreakMinutes"`
	LongBreakEnabled      bool      `json:"longBreakEnabled"`
	CyclesBeforeLongBreak int       `json:"cyclesBeforeLongBreak"`
	AutoStartBreaks       bool      `json:"autoStartBreaks"`
	AutoStartWork         bool      `json:"autoStartWork"`
	CreatedAt             string    `json:"createdAt"`
	UpdatedAt             string    `json:"updatedAt"`
}

type TimerState struct {
	State          string              `json:"state"`
	SessionID      *string             `json:"sessionId,omitempty"`
	IssueID        *int64              `json:"issueId,omitempty"`
	SegmentType    *SessionSegmentType `json:"segmentType,omitempty"`
	ElapsedSeconds int                 `json:"elapsedSeconds,omitempty"`
}

type Stash struct {
	ID                string              `json:"id"`
	UserID            string              `json:"userId"`
	DeviceID          string              `json:"deviceId"`
	RepoID            *int64              `json:"repoId,omitempty"`
	StreamID          *int64              `json:"streamId,omitempty"`
	IssueID           *int64              `json:"issueId,omitempty"`
	SessionID         *string             `json:"sessionId,omitempty"`
	PausedSegmentType *SessionSegmentType `json:"pausedSegmentType,omitempty"`
	ElapsedSeconds    *int                `json:"elapsedSeconds,omitempty"`
	Note              *string             `json:"note,omitempty"`
	CreatedAt         string              `json:"createdAt"`
	UpdatedAt         string              `json:"updatedAt"`
}

type ScratchPadMeta struct {
	ID           string `json:"id"`
	Path         string `json:"path"`
	Name         string `json:"name"`
	LastOpenedAt string `json:"lastOpenedAt"`
	Pinned       bool   `json:"pinned"`
}

type ScratchPadRead struct {
	OK      bool            `json:"ok"`
	Error   *string         `json:"error,omitempty"`
	Meta    *ScratchPadMeta `json:"meta,omitempty"`
	Content *string         `json:"content,omitempty"`
}

type Op struct {
	ID        string   `json:"id"`
	Entity    OpEntity `json:"entity"`
	EntityID  string   `json:"entityId"`
	Action    OpAction `json:"action"`
	Payload   any      `json:"payload,omitempty"`
	Timestamp string   `json:"timestamp"`
	UserID    string   `json:"userId"`
	DeviceID  string   `json:"deviceId"`
}

type Health struct {
	Status string  `json:"status"`
	DB     bool    `json:"db"`
	OK     int     `json:"ok"`
	Uptime float64 `json:"uptime"`
}

type KernelInfo struct {
	PID        int    `json:"pid"`
	Port       int    `json:"port,omitempty"`
	SocketPath string `json:"socketPath,omitempty"`
	Token      string `json:"token"`
	StartedAt  string `json:"startedAt"`
	ScratchDir string `json:"scratchDir"`
	Env        string `json:"env"`
}
