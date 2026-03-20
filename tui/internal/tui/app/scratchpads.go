package app

import (
	"crona/tui/internal/api"
	"crona/tui/internal/logger"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
)

func (m Model) scratchpadPaneSize() (int, int) {
	availH := m.contentHeight()
	if availH < 4 {
		availH = 4
	}
	return m.mainContentWidth(), availH
}

func (m *Model) syncScratchpadViewport() {
	paneW, paneH := m.scratchpadPaneSize()
	contentW := max(20, paneW-6)
	contentH := paneH - 7
	if contentH < 1 {
		contentH = 1
	}
	m.scratchpadViewport.Width = contentW
	m.scratchpadViewport.Height = contentH
	if m.scratchpadRendered != "" {
		m.scratchpadViewport.SetContent(m.scratchpadRendered)
	}
}

func (m *Model) setActiveScratchpadByIndex(idx int) {
	if idx < 0 || idx >= len(m.scratchpads) {
		return
	}
	pad := m.scratchpads[idx]
	m.scratchpadMeta = &api.ScratchPad{
		ID:           pad.ID,
		Name:         pad.Name,
		Path:         pad.Path,
		Pinned:       pad.Pinned,
		LastOpenedAt: pad.LastOpenedAt,
	}
}

func (m Model) enterScratchpadPane(msg openScratchpadMsg) Model {
	m.scratchpadOpen = true
	m.scratchpadMeta = &api.ScratchPad{
		ID:           msg.Meta.ID,
		Name:         msg.Meta.Name,
		Path:         msg.Meta.Path,
		Pinned:       msg.Meta.Pinned,
		LastOpenedAt: msg.Meta.LastOpenedAt,
	}
	m.scratchpadFilePath = msg.FilePath
	rendered, err := renderMarkdown(msg.Content, m.scratchpadRenderWidth())
	if err != nil {
		logger.Errorf("glamour render: %v", err)
		rendered = msg.Content
	}
	m.scratchpadRendered = rendered
	m.syncScratchpadViewport()
	m.scratchpadViewport.GotoTop()
	return m
}

func (m Model) scratchpadRenderWidth() int {
	paneW, _ := m.scratchpadPaneSize()
	return max(20, paneW-6)
}

func renderMarkdown(content string, width int) (string, error) {
	w := width - 2
	if w < 24 {
		w = 24
	}
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(w),
	)
	if err != nil {
		return "", err
	}
	return r.Render(content)
}

func (m Model) updateScratchpadPane(msg tea.KeyMsg) (Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		m.scratchpadOpen = false
		m.scratchpadMeta = nil
		m.scratchpadFilePath = ""
		m.scratchpadRendered = ""
		return m, nil
	case "e":
		if m.scratchpadFilePath == "" {
			return m, nil
		}
		return m, openEditor(m.scratchpadFilePath)
	case "s":
		if m.timer != nil && m.timer.State != "idle" {
			return m, m.setStatus("End or stash the active session before changing issue status", true)
		}
		return m, nil
	case "A":
		if m.timer != nil && m.timer.State != "idle" {
			next, cmd := m.abandonSelectedIssue()
			return next.(Model), cmd
		}
		return m, nil
	case "left", "h":
		if m.cursor[PaneScratchpads] > 0 {
			m.cursor[PaneScratchpads]--
			rawIdx := m.filteredIndexAtCursor(PaneScratchpads)
			if rawIdx >= 0 {
				m.setActiveScratchpadByIndex(rawIdx)
				return m, cmdOpenScratchpad(m.client, m.scratchpads, rawIdx)
			}
		}
		return m, nil
	case "right", "l":
		if m.cursor[PaneScratchpads] < m.listLen(PaneScratchpads)-1 {
			m.cursor[PaneScratchpads]++
			rawIdx := m.filteredIndexAtCursor(PaneScratchpads)
			if rawIdx >= 0 {
				m.setActiveScratchpadByIndex(rawIdx)
				return m, cmdOpenScratchpad(m.client, m.scratchpads, rawIdx)
			}
		}
		return m, nil
	case "j", "down":
		m.scratchpadViewport.ScrollDown(1)
		return m, nil
	case "k", "up":
		m.scratchpadViewport.ScrollUp(1)
		return m, nil
	case "d", "ctrl+d":
		m.scratchpadViewport.HalfPageDown()
		return m, nil
	case "u", "ctrl+u":
		m.scratchpadViewport.HalfPageUp()
		return m, nil
	case "g":
		m.scratchpadViewport.GotoTop()
		return m, nil
	case "G":
		m.scratchpadViewport.GotoBottom()
		return m, nil
	}
	var cmd tea.Cmd
	m.scratchpadViewport, cmd = m.scratchpadViewport.Update(msg)
	return m, cmd
}

func cmdReloadScratchpad(c *api.Client, meta *api.ScratchPad, width int) tea.Cmd {
	return func() tea.Msg {
		if meta == nil {
			return nil
		}
		filePath, content, err := c.ReadScratchpad(meta.ID)
		if err != nil {
			return errMsg{Err: err}
		}
		rendered, err := renderMarkdown(content, width)
		if err != nil {
			logger.Errorf("glamour render: %v", err)
			rendered = content
		}
		return scratchpadReloadedMsg{Rendered: rendered, FilePath: filePath}
	}
}

func (m Model) scratchpadTabIndexByID(id string) int {
	for i, pad := range m.scratchpads {
		if pad.ID == id {
			return i
		}
	}
	return -1
}
