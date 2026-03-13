package tui

import (
	sharedtypes "crona/shared/types"
	"encoding/json"
	"strings"
	"time"

	"crona/tui/internal/api"
	"crona/tui/internal/logger"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

// ---------- Load commands ----------

func loadRepos(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		repos, err := c.ListRepos()
		if err != nil {
			logger.Errorf("loadRepos: %v", err)
			return errMsg{err}
		}
		return reposLoadedMsg{repos}
	}
}

func loadStreams(c *api.Client, repoID int64) tea.Cmd {
	return func() tea.Msg {
		streams, err := c.ListStreams(repoID)
		if err != nil {
			logger.Errorf("loadStreams(%s): %v", repoID, err)
			return errMsg{err}
		}
		return streamsLoadedMsg{streams}
	}
}

func loadIssues(c *api.Client, streamID int64) tea.Cmd {
	return func() tea.Msg {
		issues, err := c.ListIssues(streamID)
		if err != nil {
			logger.Errorf("loadIssues(%s): %v", streamID, err)
			return errMsg{err}
		}
		return issuesLoadedMsg{streamID: streamID, issues: issues}
	}
}

func loadAllIssues(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		issues, err := c.ListAllIssues()
		if err != nil {
			logger.Errorf("loadAllIssues: %v", err)
			return errMsg{err}
		}
		return allIssuesLoadedMsg{issues}
	}
}

func loadDailySummary(c *api.Client, date string) tea.Cmd {
	return func() tea.Msg {
		summary, err := c.GetDailySummary(date)
		if err != nil {
			logger.Errorf("loadDailySummary: %v", err)
			return errMsg{err}
		}
		return dailySummaryLoadedMsg{summary}
	}
}

func loadIssueSessions(c *api.Client, issueID int64) tea.Cmd {
	return func() tea.Msg {
		sessions, err := c.ListSessionsByIssue(issueID)
		if err != nil {
			logger.Errorf("loadIssueSessions(%s): %v", issueID, err)
			return errMsg{err}
		}
		return issueSessionsLoadedMsg{issueID: issueID, sessions: sessions}
	}
}

func loadSessionHistory(c *api.Client, limit int) tea.Cmd {
	return func() tea.Msg {
		sessions, err := c.ListSessionHistory(limit)
		if err != nil {
			logger.Errorf("loadSessionHistory: %v", err)
			return errMsg{err}
		}
		return sessionHistoryLoadedMsg{sessions}
	}
}

func loadScratchpads(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		pads, err := c.ListScratchpads()
		if err != nil {
			logger.Errorf("loadScratchpads: %v", err)
			return errMsg{err}
		}
		return scratchpadsLoadedMsg{pads}
	}
}

func loadStashes(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		stashes, err := c.ListStashes()
		if err != nil {
			logger.Errorf("loadStashes: %v", err)
			return errMsg{err}
		}
		return stashesLoadedMsg{stashes}
	}
}

func loadOps(c *api.Client, limit int) tea.Cmd {
	return func() tea.Msg {
		ops, err := c.ListOps(limit)
		if err != nil {
			logger.Errorf("loadOps: %v", err)
			return errMsg{err}
		}
		return opsLoadedMsg{ops}
	}
}

func loadContext(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		ctx, err := c.GetContext()
		if err != nil {
			logger.Errorf("loadContext: %v", err)
			return errMsg{err}
		}
		return contextLoadedMsg{ctx}
	}
}

func loadTimer(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		t, err := c.GetTimerState()
		if err != nil {
			logger.Errorf("loadTimer: %v", err)
			return errMsg{err}
		}
		return timerLoadedMsg{t}
	}
}

func loadHealth(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		h, err := c.GetHealth()
		if err != nil {
			logger.Errorf("loadHealth: %v", err)
			return errMsg{err}
		}
		return healthLoadedMsg{h}
	}
}

func loadSettings(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		settings, err := c.GetSettings()
		if err != nil {
			logger.Errorf("loadSettings: %v", err)
			return errMsg{err}
		}
		return settingsLoadedMsg{settings}
	}
}

func loadKernelInfo(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		info, err := c.GetKernelInfo()
		if err != nil {
			logger.Errorf("loadKernelInfo: %v", err)
			return errMsg{err}
		}
		return kernelInfoLoadedMsg{info}
	}
}

func cmdPatchSetting(c *api.Client, key sharedtypes.CoreSettingsKey, value any) tea.Cmd {
	return func() tea.Msg {
		if err := c.PatchSetting(key, value); err != nil {
			logger.Errorf("PatchSetting(%s): %v", key, err)
			return errMsg{err}
		}
		return loadSettings(c)()
	}
}

func cmdShutdownKernel(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		if err := c.ShutdownKernel(); err != nil {
			logger.Errorf("ShutdownKernel: %v", err)
			return errMsg{err}
		}
		return kernelShutdownMsg{}
	}
}

func cmdSeedDevData(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		if err := c.SeedDevData(); err != nil {
			logger.Errorf("SeedDevData: %v", err)
			return errMsg{err}
		}
		return devSeededMsg{}
	}
}

func cmdClearDevData(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		if err := c.ClearDevData(); err != nil {
			logger.Errorf("ClearDevData: %v", err)
			return errMsg{err}
		}
		return devClearedMsg{}
	}
}

// tickAfter sends a timerTickMsg after 1 second, used for the live elapsed counter.
func tickAfter(seq int) tea.Cmd {
	return tea.Tick(time.Second, func(_ time.Time) tea.Msg {
		return timerTickMsg{seq: seq}
	})
}

func healthTickAfter() tea.Cmd {
	return tea.Tick(5*time.Second, func(_ time.Time) tea.Msg {
		return healthTickMsg{}
	})
}

// waitForEvent returns a Cmd that blocks until one kernel event arrives.
// Call it repeatedly from Update to keep draining the channel.
func waitForEvent(ch <-chan api.KernelEvent) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-ch
		if !ok {
			return nil
		}
		return kernelEventMsg{event}
	}
}

// ---------- Scratchpad commands ----------

func cmdCreateScratchpad(c *api.Client, name, path string) tea.Cmd {
	return func() tea.Msg {
		id := uuid.New().String()
		if err := c.RegisterScratchpad(id, name, path); err != nil {
			logger.Errorf("RegisterScratchpad: %v", err)
			return errMsg{err}
		}
		return loadScratchpads(c)()
	}
}

func cmdCreateRepoOnly(c *api.Client, name string) tea.Cmd {
	return func() tea.Msg {
		if _, err := c.CreateRepo(name); err != nil {
			logger.Errorf("CreateRepo: %v", err)
			return errMsg{err}
		}
		return loadRepos(c)()
	}
}

func cmdCreateStreamOnly(c *api.Client, repoID int64, name string) tea.Cmd {
	return func() tea.Msg {
		if _, err := c.CreateStream(repoID, name); err != nil {
			logger.Errorf("CreateStream: %v", err)
			return errMsg{err}
		}
		return tea.Batch(loadStreams(c, repoID), loadAllIssues(c), loadDailySummary(c, ""))()
	}
}

func cmdCreateIssueOnly(c *api.Client, streamID int64, title string, estimateMinutes *int, todoForDate *string) tea.Cmd {
	return func() tea.Msg {
		if _, err := c.CreateIssue(streamID, title, estimateMinutes, todoForDate); err != nil {
			logger.Errorf("CreateIssue: %v", err)
			return errMsg{err}
		}
		return tea.Batch(loadIssues(c, streamID), loadAllIssues(c), loadDailySummary(c, ""))()
	}
}

func cmdCreateIssueWithPath(c *api.Client, repoName, streamName, title string, estimateMinutes *int, todoForDate *string) tea.Cmd {
	return func() tea.Msg {
		repos, err := c.ListRepos()
		if err != nil {
			logger.Errorf("ListRepos before CreateIssueWithPath: %v", err)
			return errMsg{err}
		}

		var repoID int64
		for _, repo := range repos {
			if strings.EqualFold(strings.TrimSpace(repo.Name), strings.TrimSpace(repoName)) {
				repoID = repo.ID
				break
			}
		}
		if repoID == 0 {
			repo, err := c.CreateRepo(repoName)
			if err != nil {
				logger.Errorf("CreateRepo before CreateIssueWithPath: %v", err)
				return errMsg{err}
			}
			repoID = repo.ID
		}

		streams, err := c.ListStreams(repoID)
		if err != nil {
			logger.Errorf("ListStreams before CreateIssueWithPath: %v", err)
			return errMsg{err}
		}

		var streamID int64
		for _, stream := range streams {
			if strings.EqualFold(strings.TrimSpace(stream.Name), strings.TrimSpace(streamName)) {
				streamID = stream.ID
				break
			}
		}
		if streamID == 0 {
			stream, err := c.CreateStream(repoID, streamName)
			if err != nil {
				logger.Errorf("CreateStream before CreateIssueWithPath: %v", err)
				return errMsg{err}
			}
			streamID = stream.ID
		}

		if _, err := c.CreateIssue(streamID, title, estimateMinutes, todoForDate); err != nil {
			logger.Errorf("CreateIssue in CreateIssueWithPath: %v", err)
			return errMsg{err}
		}

		return tea.Batch(
			loadRepos(c),
			loadStreams(c, repoID),
			loadIssues(c, streamID),
			loadAllIssues(c),
			loadDailySummary(c, ""),
		)()
	}
}

func cmdDeleteScratchpad(c *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if err := c.DeleteScratchpad(id); err != nil {
			logger.Errorf("DeleteScratchpad(%s): %v", id, err)
			return errMsg{err}
		}
		return loadScratchpads(c)()
	}
}

func cmdOpenScratchpad(c *api.Client, scratchpads []api.ScratchPad, idx int) tea.Cmd {
	return func() tea.Msg {
		if idx >= len(scratchpads) {
			return nil
		}
		pad := scratchpads[idx]
		filePath, content, err := c.ReadScratchpad(pad.ID)
		if err != nil {
			logger.Errorf("ReadScratchpad(%s): %v", pad.ID, err)
			return errMsg{err}
		}
		return openScratchpadMsg{meta: pad, filePath: filePath, content: content}
	}
}

type openScratchpadMsg struct {
	meta     api.ScratchPad
	filePath string
	content  string
}

// ---------- Context switch commands ----------

func cmdCheckoutRepo(c *api.Client, repoID int64) tea.Cmd {
	return func() tea.Msg {
		if err := c.SwitchRepo(repoID); err != nil {
			logger.Errorf("SwitchRepo: %v", err)
			return errMsg{err}
		}
		return loadContext(c)()
	}
}

func cmdCheckoutStream(c *api.Client, streamID int64) tea.Cmd {
	return func() tea.Msg {
		if err := c.SwitchStream(streamID); err != nil {
			logger.Errorf("SwitchStream: %v", err)
			return errMsg{err}
		}
		return loadContext(c)()
	}
}

func cmdCheckoutIssue(c *api.Client, repoID, streamID, issueID int64) tea.Cmd {
	return func() tea.Msg {
		if err := c.SetFullContext(repoID, streamID, issueID); err != nil {
			logger.Errorf("SetFullContext: %v", err)
			return errMsg{err}
		}
		return loadContext(c)()
	}
}

func cmdCheckoutIssueOnly(c *api.Client, issueID int64) tea.Cmd {
	return func() tea.Msg {
		if err := c.SwitchIssue(issueID); err != nil {
			logger.Errorf("SwitchIssue: %v", err)
			return errMsg{err}
		}
		return loadContext(c)()
	}
}

func cmdChangeIssueStatus(c *api.Client, issueID int64, status string, streamID int64, dashboardDate string) tea.Cmd {
	return func() tea.Msg {
		if err := c.ChangeIssueStatus(issueID, status); err != nil {
			logger.Errorf("ChangeIssueStatus: %v", err)
			return errMsg{err}
		}

		cmds := []tea.Cmd{loadAllIssues(c), loadDailySummary(c, dashboardDate)}
		if streamID != 0 {
			cmds = append(cmds, loadIssues(c, streamID))
		}
		return tea.Batch(cmds...)()
	}
}

func cmdToggleIssueToday(c *api.Client, issueID int64, markedForToday bool, streamID int64, dashboardDate string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if markedForToday {
			err = c.ClearIssueTodo(issueID)
		} else {
			err = c.MarkIssueTodoForToday(issueID)
		}
		if err != nil {
			logger.Errorf("ToggleIssueToday: %v", err)
			return errMsg{err}
		}

		cmds := []tea.Cmd{loadAllIssues(c), loadDailySummary(c, dashboardDate)}
		if streamID != 0 {
			cmds = append(cmds, loadIssues(c, streamID))
		}
		return tea.Batch(cmds...)()
	}
}

func cmdSetIssueTodoDate(c *api.Client, issueID int64, date string, streamID int64, dashboardDate string) tea.Cmd {
	return func() tea.Msg {
		var err error
		if strings.TrimSpace(date) == "" {
			err = c.ClearIssueTodo(issueID)
		} else {
			err = c.SetIssueTodoDate(issueID, date)
		}
		if err != nil {
			logger.Errorf("SetIssueTodoDate: %v", err)
			return errMsg{err}
		}

		cmds := []tea.Cmd{loadAllIssues(c), loadDailySummary(c, dashboardDate)}
		if streamID != 0 {
			cmds = append(cmds, loadIssues(c, streamID))
		}
		return tea.Batch(cmds...)()
	}
}

func cmdChangeIssueStatusAndEndSession(c *api.Client, issueID int64, status string, streamID int64, dashboardDate, message string) tea.Cmd {
	return func() tea.Msg {
		if err := c.ChangeIssueStatus(issueID, status); err != nil {
			logger.Errorf("ChangeIssueStatus before EndTimer: %v", err)
			return errMsg{err}
		}
		if err := c.EndTimer(message); err != nil {
			logger.Errorf("EndTimer after ChangeIssueStatus: %v", err)
			return errMsg{err}
		}

		cmds := []tea.Cmd{
			loadAllIssues(c),
			loadDailySummary(c, dashboardDate),
			loadContext(c),
			loadTimer(c),
		}
		if streamID != 0 {
			cmds = append(cmds, loadIssues(c, streamID))
		}
		return tea.Batch(cmds...)()
	}
}

func cmdStartFocusSession(c *api.Client, repoID, streamID, issueID int64) tea.Cmd {
	return func() tea.Msg {
		if repoID != 0 && streamID != 0 && issueID != 0 {
			if err := c.SetFullContext(repoID, streamID, issueID); err != nil {
				logger.Errorf("SetFullContext before StartTimer: %v", err)
				return errMsg{err}
			}
		} else if issueID != 0 {
			if err := c.SwitchIssue(issueID); err != nil {
				logger.Errorf("SwitchIssue before StartTimer: %v", err)
				return errMsg{err}
			}
		}

		if err := c.StartTimer(0); err != nil {
			logger.Errorf("StartTimer: %v", err)
			return errMsg{err}
		}
		return focusSessionChangedMsg{reloadContext: true, reloadTimer: true}
	}
}

func cmdPauseFocusSession(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		if err := c.PauseTimer(); err != nil {
			logger.Errorf("PauseTimer: %v", err)
			return errMsg{err}
		}
		return focusSessionChangedMsg{reloadTimer: true}
	}
}

func cmdResumeFocusSession(c *api.Client) tea.Cmd {
	return func() tea.Msg {
		if err := c.ResumeTimer(); err != nil {
			logger.Errorf("ResumeTimer: %v", err)
			return errMsg{err}
		}
		return focusSessionChangedMsg{reloadTimer: true}
	}
}

func cmdEndFocusSession(c *api.Client, streamID int64, dashboardDate, message string) tea.Cmd {
	return func() tea.Msg {
		if err := c.EndTimer(message); err != nil {
			logger.Errorf("EndTimer: %v", err)
			return errMsg{err}
		}
		cmds := []tea.Cmd{
			loadAllIssues(c),
			loadDailySummary(c, dashboardDate),
			loadContext(c),
			loadTimer(c),
			loadSessionHistory(c, 200),
		}
		if streamID != 0 {
			cmds = append(cmds, loadIssues(c, streamID))
		}
		return tea.Batch(cmds...)()
	}
}

func cmdStashFocusSession(c *api.Client, note string) tea.Cmd {
	return func() tea.Msg {
		if err := c.StashPush(note); err != nil {
			logger.Errorf("StashPush: %v", err)
			return errMsg{err}
		}
		return focusSessionChangedMsg{reloadContext: true, reloadTimer: true}
	}
}

func cmdApplyStash(c *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if err := c.ApplyStash(id); err != nil {
			logger.Errorf("ApplyStash: %v", err)
			return errMsg{err}
		}
		return tea.Batch(
			loadStashes(c),
			loadContext(c),
			loadTimer(c),
			loadSessionHistory(c, 200),
		)()
	}
}

type focusSessionChangedMsg struct {
	reloadContext bool
	reloadTimer   bool
}

// ---------- Kernel event -> data refresh ----------

// handleKernelEvent maps a raw kernel event to the appropriate reload command(s).
func handleKernelEvent(m Model, event api.KernelEvent) (Model, tea.Cmd) {
	logger.Infof("kernel event: %s", event.Type)

	switch event.Type {

	case "repo.created", "repo.updated", "repo.deleted":
		return m, loadRepos(m.client)

	case "stream.created", "stream.updated", "stream.deleted":
		if m.context != nil && m.context.RepoID != nil {
			return m, loadStreams(m.client, *m.context.RepoID)
		}

	case "issue.created", "issue.updated", "issue.deleted":
		cmds := []tea.Cmd{loadAllIssues(m.client), loadDailySummary(m.client, m.dashboardDate)}
		if m.context != nil && m.context.StreamID != nil {
			cmds = append(cmds, loadIssues(m.client, *m.context.StreamID))
		}
		return m, tea.Batch(cmds...)

	case "scratchpad.created", "scratchpad.updated", "scratchpad.deleted":
		return m, loadScratchpads(m.client)

	case "session.started", "session.stopped":
		return m, tea.Batch(loadTimer(m.client), loadContext(m.client), loadSessionHistory(m.client, 200))

	case "stash.created", "stash.applied", "stash.dropped":
		return m, tea.Batch(loadStashes(m.client), loadContext(m.client), loadTimer(m.client))

	case "context.repo.changed", "context.stream.changed", "context.issue.changed", "context.cleared":
		var payload sharedtypes.ContextChangedPayload
		_ = json.Unmarshal(event.Payload, &payload)

		cmds := []tea.Cmd{loadContext(m.client)}
		if payload.RepoID != nil {
			cmds = append(cmds, loadStreams(m.client, *payload.RepoID))
		} else {
			m.streams = nil
			m.issues = nil
			m.cursor[PaneStreams] = 0
			m.cursor[PaneIssues] = 0
		}
		if payload.StreamID != nil {
			cmds = append(cmds, loadIssues(m.client, *payload.StreamID))
		} else if payload.RepoID != nil {
			m.issues = nil
			m.cursor[PaneIssues] = 0
		}
		return m, tea.Batch(cmds...)

	case "timer.state":
		var timer api.TimerState
		if err := json.Unmarshal(event.Payload, &timer); err == nil {
			m.timer = &timer
			m.elapsed = 0
			m.timerTickSeq++
			if timer.State != "idle" {
				if m.view != ViewScratch {
					m.view = ViewSessionActive
				}
				m.pane = viewDefaultPane[m.view]
				return m, tickAfter(m.timerTickSeq)
			} else if m.view == ViewSessionActive {
				m.view = ViewDefault
				m.pane = viewDefaultPane[m.view]
			}
		}

	case "timer.boundary":
		m.elapsed = 0
		return m, loadTimer(m.client)

	case "ops.created":
		return m, loadOps(m.client, m.currentOpsLimit())
	}

	return m, nil
}
