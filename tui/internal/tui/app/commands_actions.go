package app

import (
	shareddto "crona/shared/dto"
	sharedtypes "crona/shared/types"
	"strings"

	"crona/tui/internal/api"
	"crona/tui/internal/logger"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
)

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
func cmdUpdateRepo(c *api.Client, repoID int64, name string) tea.Cmd {
	return func() tea.Msg {
		if err := c.UpdateRepo(repoID, name); err != nil {
			logger.Errorf("UpdateRepo: %v", err)
			return errMsg{err}
		}
		return tea.Batch(loadRepos(c), loadContext(c), loadAllIssues(c))()
	}
}
func cmdDeleteRepo(c *api.Client, repoID int64) tea.Cmd {
	return func() tea.Msg {
		if err := c.DeleteRepo(repoID); err != nil {
			logger.Errorf("DeleteRepo: %v", err)
			return errMsg{err}
		}
		return tea.Batch(loadRepos(c), loadAllIssues(c), loadDailySummary(c, ""), loadContext(c))()
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
func cmdUpdateStream(c *api.Client, repoID, streamID int64, name string) tea.Cmd {
	return func() tea.Msg {
		if err := c.UpdateStream(streamID, name); err != nil {
			logger.Errorf("UpdateStream: %v", err)
			return errMsg{err}
		}
		return tea.Batch(loadStreams(c, repoID), loadAllIssues(c), loadDailySummary(c, ""), loadContext(c))()
	}
}

func cmdDeleteStream(c *api.Client, repoID, streamID int64) tea.Cmd {
	return func() tea.Msg {
		if err := c.DeleteStream(streamID); err != nil {
			logger.Errorf("DeleteStream: %v", err)
			return errMsg{err}
		}
		cmds := []tea.Cmd{loadAllIssues(c), loadDailySummary(c, ""), loadContext(c)}
		if repoID != 0 {
			cmds = append(cmds, loadStreams(c, repoID))
		}
		return tea.Batch(cmds...)()
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

func cmdUpdateIssue(c *api.Client, issueID, streamID int64, title string, estimateMinutes *int, todoForDate *string, dashboardDate string) tea.Cmd {
	return func() tea.Msg {
		if err := c.UpdateIssue(issueID, title, estimateMinutes); err != nil {
			logger.Errorf("UpdateIssue: %v", err)
			return errMsg{err}
		}
		currentDue := ""
		if todoForDate != nil {
			currentDue = *todoForDate
		}
		if strings.TrimSpace(currentDue) == "" {
			if err := c.ClearIssueTodo(issueID); err != nil {
				logger.Errorf("ClearIssueTodo after UpdateIssue: %v", err)
				return errMsg{err}
			}
		} else if err := c.SetIssueTodoDate(issueID, currentDue); err != nil {
			logger.Errorf("SetIssueTodoDate after UpdateIssue: %v", err)
			return errMsg{err}
		}
		cmds := []tea.Cmd{loadAllIssues(c), loadDailySummary(c, dashboardDate), loadContext(c)}
		if streamID != 0 {
			cmds = append(cmds, loadIssues(c, streamID))
		}
		return tea.Batch(cmds...)()
	}
}

func cmdDeleteIssue(c *api.Client, issueID, streamID int64, dashboardDate string) tea.Cmd {
	return func() tea.Msg {
		if err := c.DeleteIssue(issueID); err != nil {
			logger.Errorf("DeleteIssue: %v", err)
			return errMsg{err}
		}
		cmds := []tea.Cmd{loadAllIssues(c), loadDailySummary(c, dashboardDate), loadContext(c), loadTimer(c)}
		if streamID != 0 {
			cmds = append(cmds, loadIssues(c, streamID))
		}
		return tea.Batch(cmds...)()
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

		return tea.Batch(loadRepos(c), loadStreams(c, repoID), loadIssues(c, streamID), loadAllIssues(c), loadDailySummary(c, ""))()
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

func cmdChangeIssueStatus(c *api.Client, issueID int64, status string, note *string, streamID int64, dashboardDate string) tea.Cmd {
	return func() tea.Msg {
		if err := c.ChangeIssueStatus(issueID, status, note); err != nil {
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

func cmdChangeIssueStatusAndEndSession(c *api.Client, issueID int64, status string, note *string, streamID int64, dashboardDate string, endInput shareddto.EndSessionRequest) tea.Cmd {
	return func() tea.Msg {
		if err := c.EndTimer(endInput); err != nil {
			logger.Errorf("EndTimer before ChangeIssueStatus: %v", err)
			return errMsg{err}
		}
		if err := c.ChangeIssueStatus(issueID, status, note); err != nil {
			logger.Errorf("ChangeIssueStatus after EndTimer: %v", err)
			return errMsg{err}
		}
		cmds := []tea.Cmd{loadAllIssues(c), loadDailySummary(c, dashboardDate), loadContext(c), loadTimer(c)}
		if streamID != 0 {
			cmds = append(cmds, loadIssues(c, streamID))
		}
		return tea.Batch(cmds...)()
	}
}

func cmdAmendSessionNote(c *api.Client, id string, note string) tea.Cmd {
	return func() tea.Msg {
		if err := c.AmendSessionNote(id, note); err != nil {
			logger.Errorf("AmendSessionNote(%s): %v", id, err)
			return errMsg{err}
		}
		return sessionAmendedMsg{id: id}
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

func cmdEndFocusSession(c *api.Client, streamID int64, dashboardDate string, endInput shareddto.EndSessionRequest) tea.Cmd {
	return func() tea.Msg {
		if err := c.EndTimer(endInput); err != nil {
			logger.Errorf("EndTimer: %v", err)
			return errMsg{err}
		}
		cmds := []tea.Cmd{loadAllIssues(c), loadDailySummary(c, dashboardDate), loadContext(c), loadTimer(c), loadSessionHistory(c, 200)}
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
		return tea.Batch(loadStashes(c), loadContext(c), loadTimer(c), loadSessionHistory(c, 200))()
	}
}

func cmdDropStash(c *api.Client, id string) tea.Cmd {
	return func() tea.Msg {
		if err := c.DropStash(id); err != nil {
			logger.Errorf("DropStash: %v", err)
			return errMsg{err}
		}
		return loadStashes(c)()
	}
}

type focusSessionChangedMsg struct {
	reloadContext bool
	reloadTimer   bool
}
