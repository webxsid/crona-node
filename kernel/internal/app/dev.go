package app

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	corecommands "crona/kernel/internal/core/commands"
	"crona/kernel/internal/scratchfile"
	"crona/kernel/internal/store"
	sharedtypes "crona/shared/types"
)

func (h *Handler) clearDevData(ctx context.Context) error {
	h.timer.ClearBoundary()

	if err := store.ClearAllData(ctx, h.core.Store.DB()); err != nil {
		return fmt.Errorf("clear sqlite data: %w", err)
	}
	if err := scratchfile.ClearAll(h.core.ScratchDir); err != nil {
		return fmt.Errorf("clear scratch files: %w", err)
	}
	if err := h.core.InitDefaults(ctx); err != nil {
		return fmt.Errorf("reinitialize defaults: %w", err)
	}

	payload, _ := json.Marshal(sharedtypes.TimerState{State: "idle"})
	h.core.Events.Emit(sharedtypes.KernelEvent{
		Type:    sharedtypes.EventTypeTimerState,
		Payload: payload,
	})
	return nil
}

func (h *Handler) seedDevData(ctx context.Context) error {
	if err := h.clearDevData(ctx); err != nil {
		return err
	}

	today := time.Now().UTC().Format("2006-01-02")

	workRepo, err := corecommands.CreateRepo(ctx, h.core, struct {
		Name  string
		Color *string
	}{
		Name:  "Work",
		Color: devStringPtr("blue"),
	})
	if err != nil {
		return err
	}
	personalRepo, err := corecommands.CreateRepo(ctx, h.core, struct {
		Name  string
		Color *string
	}{
		Name:  "Personal",
		Color: devStringPtr("green"),
	})
	if err != nil {
		return err
	}

	appStream, err := corecommands.CreateStream(ctx, h.core, struct {
		RepoID     int64
		Name       string
		Visibility *sharedtypes.StreamVisibility
	}{
		RepoID: workRepo.ID,
		Name:   "app",
	})
	if err != nil {
		return err
	}
	infraStream, err := corecommands.CreateStream(ctx, h.core, struct {
		RepoID     int64
		Name       string
		Visibility *sharedtypes.StreamVisibility
	}{
		RepoID: workRepo.ID,
		Name:   "infra",
	})
	if err != nil {
		return err
	}
	homeStream, err := corecommands.CreateStream(ctx, h.core, struct {
		RepoID     int64
		Name       string
		Visibility *sharedtypes.StreamVisibility
	}{
		RepoID: personalRepo.ID,
		Name:   "home",
	})
	if err != nil {
		return err
	}

	focusIssue, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        appStream.ID,
		Title:           "Port dev tooling to Go",
		EstimateMinutes: devIntPtr(90),
		Notes:           devStringPtr("Validate the new IPC-first workflow."),
		TodoForDate:     &today,
	})
	if err != nil {
		return err
	}
	if _, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        appStream.ID,
		Title:           "Add CLI command surface",
		EstimateMinutes: devIntPtr(60),
	}); err != nil {
		return err
	}
	doneIssue, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        infraStream.ID,
		Title:           "Replace HTTP transport",
		EstimateMinutes: devIntPtr(120),
	})
	if err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, doneIssue.ID, sharedtypes.IssueStatusActive); err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, doneIssue.ID, sharedtypes.IssueStatusDone); err != nil {
		return err
	}
	if _, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        homeStream.ID,
		Title:           "Plan weekend errands",
		EstimateMinutes: devIntPtr(30),
		TodoForDate:     &today,
	}); err != nil {
		return err
	}

	scratchPath, err := corecommands.RegisterScratchpad(ctx, h.core, sharedtypes.ScratchPadMeta{
		Name: "Daily Notes",
		Path: "dev/[[date]]-notes.md",
	})
	if err != nil {
		return err
	}
	if _, err := scratchfile.Create(h.core.ScratchDir, scratchPath, "Daily Notes"); err != nil {
		return err
	}

	if _, err := corecommands.SetContext(ctx, h.core, corecommands.ContextPatch{
		RepoSet:   true,
		RepoID:    &workRepo.ID,
		StreamSet: true,
		StreamID:  &appStream.ID,
		IssueSet:  true,
		IssueID:   &focusIssue.ID,
	}); err != nil {
		return err
	}

	return nil
}

func devStringPtr(value string) *string {
	return &value
}

func devIntPtr(value int) *int {
	return &value
}
