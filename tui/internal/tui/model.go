package tui

import (
	"crona/tui/internal/api"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ---------- View / Pane types ----------

type View string

const (
	ViewDefault        View = "default"
	ViewDaily          View = "daily"
	ViewMeta           View = "meta"
	ViewSessionHistory View = "session_history"
	ViewSessionActive  View = "session_active"
	ViewScratch        View = "scratchpads"
	ViewOps            View = "ops"
	ViewSettings       View = "settings"
)

// viewOrder only includes the tab-switchable views.
var viewOrder = []View{ViewSessionHistory, ViewDefault, ViewMeta, ViewScratch, ViewOps, ViewSettings, ViewDaily}

type Pane string

const (
	PaneRepos       Pane = "repos"
	PaneStreams     Pane = "streams"
	PaneIssues      Pane = "issues"
	PaneSessions    Pane = "sessions"
	PaneScratchpads Pane = "scratchpads"
	PaneOps         Pane = "ops"
	PaneSettings    Pane = "settings"
)

// viewPanes lists the focusable panes for each view.
var viewPanes = map[View][]Pane{
	ViewDefault:        {PaneIssues},
	ViewDaily:          {PaneIssues},
	ViewMeta:           {PaneRepos, PaneStreams, PaneIssues},
	ViewSessionHistory: {PaneSessions},
	ViewSessionActive:  {},
	ViewScratch:        {PaneScratchpads},
	ViewOps:            {PaneOps},
	ViewSettings:       {PaneSettings},
}

// viewDefaultPane is the initial focused pane when entering a view.
var viewDefaultPane = map[View]Pane{
	ViewDefault:        PaneIssues,
	ViewDaily:          PaneIssues,
	ViewMeta:           PaneRepos,
	ViewSessionHistory: PaneSessions,
	ViewSessionActive:  PaneIssues,
	ViewScratch:        PaneScratchpads,
	ViewOps:            PaneOps,
	ViewSettings:       PaneSettings,
}

// ---------- Messages ----------

type reposLoadedMsg struct{ repos []api.Repo }
type streamsLoadedMsg struct{ streams []api.Stream }
type issuesLoadedMsg struct {
	streamID int64
	issues   []api.Issue
}
type allIssuesLoadedMsg struct{ issues []api.IssueWithMeta }
type dailySummaryLoadedMsg struct{ summary *api.DailyIssueSummary }
type issueSessionsLoadedMsg struct {
	issueID  int64
	sessions []api.Session
}
type sessionHistoryLoadedMsg struct{ sessions []api.SessionHistoryEntry }
type scratchpadsLoadedMsg struct{ pads []api.ScratchPad }
type stashesLoadedMsg struct{ stashes []api.Stash }
type opsLoadedMsg struct{ ops []api.Op }
type contextLoadedMsg struct{ ctx *api.ActiveContext }
type timerLoadedMsg struct{ timer *api.TimerState }
type healthLoadedMsg struct{ health *api.Health }
type settingsLoadedMsg struct{ settings *api.CoreSettings }
type kernelInfoLoadedMsg struct{ info *api.KernelInfo }
type kernelEventMsg struct{ event api.KernelEvent }
type kernelShutdownMsg struct{}
type devSeededMsg struct{}
type devClearedMsg struct{}
type timerTickMsg struct{ seq int }
type healthTickMsg struct{}
type errMsg struct{ err error }

// ---------- Model ----------

type Model struct {
	// kernel client
	client *api.Client

	// kernel event stream
	eventStop chan struct{}

	// view / navigation
	view    View
	pane    Pane
	cursor  map[Pane]int
	filters map[Pane]string

	// pane-local search/filter input
	filterEditing  bool
	filterPane     Pane
	filterInput    textinput.Model
	opsLimit       int
	opsLimitPinned bool

	// data
	repos          []api.Repo
	streams        []api.Stream
	issues         []api.Issue // context-filtered (by active streamId)
	allIssues      []api.IssueWithMeta
	dailySummary   *api.DailyIssueSummary
	dashboardDate  string
	issueSessions  []api.Session
	sessionHistory []api.SessionHistoryEntry
	scratchpads    []api.ScratchPad
	stashes        []api.Stash
	ops            []api.Op
	context        *api.ActiveContext
	timer          *api.TimerState
	health         *api.Health
	settings       *api.CoreSettings
	kernelInfo     *api.KernelInfo
	elapsed        int // local seconds since last timer.state event
	timerTickSeq   int

	// terminal dimensions
	width  int
	height int

	// scratchpad reader state within the scratchpads pane
	scratchpadOpen     bool
	scratchpadMeta     *api.ScratchPad
	scratchpadFilePath string // resolved absolute path for $EDITOR
	scratchpadRendered string // glamour-rendered content
	scratchpadViewport viewport.Model

	// dialog state
	dialog            string // "" | "create_scratchpad" | "confirm_delete" | "stash_list"
	dialogInputs      []textinput.Model
	dialogFocusIdx    int
	dialogDeleteID    string // scratchpad id pending deletion
	dialogIssueID     int64
	dialogIssueStatus string
	dialogRepoID      int64
	dialogRepoName    string
	dialogStreamID    int64
	dialogStreamName  string
	dialogRepoIndex   int
	dialogStreamIndex int
	dialogParent      string
	dialogDateMonth   string
	dialogDateCursor  string
	dialogStashCursor int

	// status / error flash
	statusMsg string
}

// SetEventChannel provides the kernel event channel from main before the program starts.
func SetEventChannel(ch <-chan api.KernelEvent) {
	eventChannel = ch
}

func New(socketPath, scratchDir string, env string, done chan struct{}) Model {
	return Model{
		client:    api.NewClient(socketPath, scratchDir),
		eventStop: done,
		view:      ViewDefault,
		pane:      PaneIssues,
		cursor: map[Pane]int{
			PaneRepos:       0,
			PaneStreams:     0,
			PaneIssues:      0,
			PaneSessions:    0,
			PaneScratchpads: 0,
			PaneOps:         0,
			PaneSettings:    0,
		},
		filters: map[Pane]string{
			PaneRepos:       "",
			PaneStreams:     "",
			PaneIssues:      "",
			PaneSessions:    "",
			PaneScratchpads: "",
			PaneOps:         "",
			PaneSettings:    "",
		},
		kernelInfo: &api.KernelInfo{Env: env},
	}
}

// eventChannel receives kernel events forwarded from main.go.
var eventChannel <-chan api.KernelEvent

// ---------- Init ----------

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		loadRepos(m.client),
		loadAllIssues(m.client),
		loadDailySummary(m.client, ""),
		loadSessionHistory(m.client, 200),
		loadScratchpads(m.client),
		loadOps(m.client, m.currentOpsLimit()),
		loadContext(m.client),
		loadTimer(m.client),
		loadHealth(m.client),
		loadSettings(m.client),
		loadKernelInfo(m.client),
		healthTickAfter(),
		waitForEvent(eventChannel),
	)
}

// ---------- Helpers: clamp cursor ----------

func (m *Model) clamp(p Pane, max int) {
	if max == 0 {
		m.cursor[p] = 0
		return
	}
	if m.cursor[p] >= max {
		m.cursor[p] = max - 1
	}
}

func (m *Model) listLen(p Pane) int {
	return len(m.filteredIndices(p))
}

func (m *Model) defaultOpsLimit() int {
	availableHeight := m.contentHeight()
	if availableHeight < 4 {
		availableHeight = 4
	}
	visibleRows := availableHeight - 6
	if visibleRows < 10 {
		visibleRows = 10
	}
	return visibleRows
}

func (m *Model) currentOpsLimit() int {
	if m.opsLimit > 0 {
		return m.opsLimit
	}
	return m.defaultOpsLimit()
}

func (m Model) contentHeight() int {
	headerH := 4
	if m.width > 0 {
		headerH = lipgloss.Height(m.renderHeader())
	}
	helpH := 1
	if m.width > 0 {
		helpH = lipgloss.Height(m.renderHelpBar())
	}
	availableHeight := m.height - headerH - helpH
	if m.statusMsg != "" {
		availableHeight--
	}
	if availableHeight < 4 {
		availableHeight = 4
	}
	return availableHeight
}

func (m Model) sidebarWidth() int {
	if m.width < 90 {
		return 18
	}
	return 22
}

func (m Model) mainContentWidth() int {
	width := m.width - m.sidebarWidth() - 4
	if width < 40 {
		return 40
	}
	return width
}

func (m Model) isDevMode() bool {
	return m.kernelInfo != nil && m.kernelInfo.Env == "Dev"
}
