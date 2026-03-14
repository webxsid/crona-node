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
	tomorrow := time.Now().UTC().Add(24 * time.Hour).Format("2006-01-02")

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
	_, err = corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        appStream.ID,
		Title:           "Add CLI command surface",
		EstimateMinutes: devIntPtr(60),
		TodoForDate:     &tomorrow,
	})
	if err != nil {
		return err
	}

	readyIssue, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        infraStream.ID,
		Title:           "Prepare rollout checklist",
		EstimateMinutes: devIntPtr(45),
		Notes:           devStringPtr("Waiting only on final go/no-go review."),
		TodoForDate:     &today,
	})
	if err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, readyIssue.ID, sharedtypes.IssueStatusReady, nil); err != nil {
		return err
	}

	inProgressIssue, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        infraStream.ID,
		Title:           "Wire release packaging checks",
		EstimateMinutes: devIntPtr(75),
		Notes:           devStringPtr("Partially implemented; needs artifact verification."),
	})
	if err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, inProgressIssue.ID, sharedtypes.IssueStatusPlanned, nil); err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, inProgressIssue.ID, sharedtypes.IssueStatusInProgress, nil); err != nil {
		return err
	}

	blockedIssue, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        infraStream.ID,
		Title:           "Provision CI signing secrets",
		EstimateMinutes: devIntPtr(30),
	})
	if err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, blockedIssue.ID, sharedtypes.IssueStatusPlanned, nil); err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, blockedIssue.ID, sharedtypes.IssueStatusBlocked, devStringPtr("Awaiting access to the signing account")); err != nil {
		return err
	}

	reviewIssue, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        appStream.ID,
		Title:           "Review lifecycle UX copy",
		EstimateMinutes: devIntPtr(40),
	})
	if err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, reviewIssue.ID, sharedtypes.IssueStatusPlanned, nil); err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, reviewIssue.ID, sharedtypes.IssueStatusInProgress, nil); err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, reviewIssue.ID, sharedtypes.IssueStatusInReview, devStringPtr("Ready for wording and interaction review")); err != nil {
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
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, doneIssue.ID, sharedtypes.IssueStatusPlanned, nil); err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, doneIssue.ID, sharedtypes.IssueStatusInProgress, nil); err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, doneIssue.ID, sharedtypes.IssueStatusDone, nil); err != nil {
		return err
	}

	abandonedIssue, err := corecommands.CreateIssue(ctx, h.core, struct {
		StreamID        int64
		Title           string
		EstimateMinutes *int
		Notes           *string
		TodoForDate     *string
	}{
		StreamID:        homeStream.ID,
		Title:           "Research standing desk options",
		EstimateMinutes: devIntPtr(25),
	})
	if err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, abandonedIssue.ID, sharedtypes.IssueStatusPlanned, nil); err != nil {
		return err
	}
	if _, err := corecommands.ChangeIssueStatus(ctx, h.core, abandonedIssue.ID, sharedtypes.IssueStatusAbandoned, devStringPtr("Deferred until the room reorganization is done")); err != nil {
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
