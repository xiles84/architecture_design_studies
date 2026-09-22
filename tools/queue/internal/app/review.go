package app

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"adsqueue/internal/archive"
	"adsqueue/internal/gitx"
	"adsqueue/internal/model"
)

// ---------------------------------------------------------------------------
// review / approve / request-changes
// ---------------------------------------------------------------------------

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func reviewDirs(rec *model.TaskRecord) []string {
	paths, _ := filepath.Glob(filepath.Join(rec.Dir, "reviews", "*", "REVIEW.md"))
	return paths
}

func submitResultPath(rec *model.TaskRecord) string {
	matches, _ := filepath.Glob(filepath.Join(rec.Dir, "attempts", "*", "RESULT.md"))
	if len(matches) == 0 {
		return ""
	}
	return matches[len(matches)-1]
}

func cmdReview(args []string, out io.Writer) error {
	fs := newFlagSet("review")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	reviewID := fs.String("review-id", "", "review id")
	reviewFile := fs.String("review-file", "", "path to REVIEW.md to record (required)")
	summary := fs.String("summary", "", "one-line review summary")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if *taskID == "" || *reviewFile == "" {
		return failf(2, "review: --task and --review-file are required")
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "review: %v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	if rec.State != model.StateAwaitingReview {
		return failf(1, "review: task %s is %s, not awaiting_review", *taskID, rec.State)
	}
	rid := *reviewID
	if rid == "" {
		rid = "review-" + archive.CompactTime(c.Now) + "-" + randHex(2)
	}
	dst := filepath.Join(rec.Dir, "reviews", rid, "REVIEW.md")
	if err := copyFile(briefMustAbs(*reviewFile), dst); err != nil {
		return err
	}
	rel, err := c.Repo.RepoRelative(dst)
	if err != nil {
		return err
	}
	sum := *summary
	if sum == "" {
		sum = fmt.Sprintf("Review %s recorded by %s.", rid, id.Verbose())
	}
	ev, err := c.emitEvent(c.Repo, rec, eventParams{
		EventType:      "reviewed",
		State:          model.StateAwaitingReview,
		Actor:          id,
		Summary:        sum,
		EvidencePaths:  []string{rel},
		Related:        model.Related{ReviewID: rid},
		NextCapability: model.CapHIGH,
		NextWorkRole:   model.RoleReviewer,
		StatusNote:     "Review `" + rid + "` recorded; next: approve or request-changes.",
	})
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(ev)
	}
	fmt.Fprintf(out, "reviewed %s: %s\n", *taskID, rid)
	return nil
}

func cmdApprove(args []string, out io.Writer) error {
	fs := newFlagSet("approve")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	summary := fs.String("summary", "", "decision summary")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "approve: %v", err)
	}
	if id.SessionCapability != model.CapHIGH {
		return failf(2, "approve: only a HIGH session approves (got %s)", id.SessionCapability)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	if rec.State != model.StateAwaitingReview {
		return failf(1, "approve: task %s is %s, not awaiting_review", *taskID, rec.State)
	}
	if len(reviewDirs(rec)) == 0 {
		return failf(1, "approve: task %s has no recorded review under reviews/*/REVIEW.md", *taskID)
	}
	sum := *summary
	if sum == "" {
		sum = fmt.Sprintf("Approved for integration by %s.", id.Verbose())
	}
	ev, err := c.emitEvent(c.Repo, rec, eventParams{
		EventType:      "approved",
		State:          model.StateApprovedForIntegration,
		Actor:          id,
		Summary:        sum,
		NextCapability: model.CapHIGH,
		NextWorkRole:   model.RoleIntegrator,
		StatusNote:     "Approved for integration. NEXT: `queue integrate`.",
	})
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(ev)
	}
	fmt.Fprintf(out, "approved %s -> approved_for_integration\n", *taskID)
	return nil
}

func cmdRequestChanges(args []string, out io.Writer) error {
	fs := newFlagSet("request-changes")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	reviewFile := fs.String("review-file", "", "path to REVIEW.md to record")
	summary := fs.String("summary", "", "what must change")
	reviewID := fs.String("review-id", "", "review id")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if *summary == "" {
		return failf(2, "request-changes: --summary is required")
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "request-changes: %v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	if rec.State != model.StateAwaitingReview {
		return failf(1, "request-changes: task %s is %s, not awaiting_review", *taskID, rec.State)
	}
	var evidence []string
	rid := *reviewID
	if *reviewFile != "" {
		if rid == "" {
			rid = "review-" + archive.CompactTime(c.Now) + "-" + randHex(2)
		}
		dst := filepath.Join(rec.Dir, "reviews", rid, "REVIEW.md")
		if err := copyFile(briefMustAbs(*reviewFile), dst); err != nil {
			return err
		}
		rel, err := c.Repo.RepoRelative(dst)
		if err != nil {
			return err
		}
		evidence = append(evidence, rel)
	}
	ev, err := c.emitEvent(c.Repo, rec, eventParams{
		EventType:      "changes_requested",
		State:          model.StateChangesRequested,
		Actor:          id,
		Summary:        "Changes requested: " + *summary,
		EvidencePaths:  evidence,
		Related:        model.Related{ReviewID: rid},
		NextCapability: model.CapLOW,
		NextWorkRole:   model.RoleExecutor,
		StatusNote:     "Corrections required: " + *summary + "\n\nClaim again with `queue claim --task " + *taskID + "`.",
	})
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(ev)
	}
	fmt.Fprintf(out, "request-changes %s -> changes_requested\n", *taskID)
	return nil
}

// ---------------------------------------------------------------------------
// integrate / complete
// ---------------------------------------------------------------------------

// acquireIntegrationLock takes the single local-main integration lock with a CAS
// create. It returns the observed object id so the caller can release it with a
// matching delete.
func (c *cmdCtx) acquireIntegrationLock(taskID, claimID string, worker model.Identity) (string, error) {
	lock := model.IntegrationLock{
		SchemaVersion: model.SchemaVersion,
		TaskID:        taskID,
		ClaimID:       claimID,
		Worker:        worker,
		AcquiredAt:    model.FormatTime(c.Now),
		Branch:        "main",
	}
	data, err := json.MarshalIndent(lock, "", "  ")
	if err != nil {
		return "", err
	}
	oid, err := c.Repo.WriteBlob(append(data, '\n'))
	if err != nil {
		return "", err
	}
	if err := c.Repo.CAS(IntegrationLockRef, oid, "", fmt.Sprintf("queue integrate %s", taskID)); err != nil {
		// Name the holder when the blob is readable.
		holder := "unknown"
		if oldOID, ok := c.Repo.RefOID(IntegrationLockRef); ok {
			if raw, rerr := c.Repo.ReadBlob(oldOID); rerr == nil {
				var h model.IntegrationLock
				if json.Unmarshal(raw, &h) == nil {
					holder = fmt.Sprintf("%s since %s", h.Worker.Verbose(), h.AcquiredAt)
				}
			}
		}
		return "", failf(1, "another integrator holds %s (%s)", IntegrationLockRef, holder)
	}
	return oid, nil
}

// mainWorktree finds the work tree checked out on local `main`.
func (c *cmdCtx) mainWorktree() (*gitx.Repo, error) {
	wts, err := c.Repo.Worktrees()
	if err != nil {
		return nil, err
	}
	for _, w := range wts {
		if w.Branch == "main" && w.Path != "" {
			return gitx.Open(w.Path)
		}
	}
	if br, _ := c.Repo.CurrentBranch(); br == "main" {
		return c.Repo, nil
	}
	return nil, fmt.Errorf("no worktree is checked out on local `main`")
}

func recordIn(root, taskID string) (*model.TaskRecord, error) {
	dir, err := archive.TaskDirAbs(root, taskID)
	if err != nil {
		return nil, err
	}
	return archive.LoadTaskRecord(dir)
}

func cmdIntegrate(args []string, out io.Writer) error {
	fs := newFlagSet("integrate")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	tagMessage := fs.String("tag-message", "", "annotated tag message (default: generated)")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "integrate: %v", err)
	}
	if id.SessionCapability != model.CapHIGH {
		return failf(2, "integrate: only a HIGH session integrates (got %s)", id.SessionCapability)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	if rec.State != model.StateApprovedForIntegration {
		return failf(1, "integrate: task %s is %s, not approved_for_integration", *taskID, rec.State)
	}
	main, err := c.mainWorktree()
	if err != nil {
		return err
	}
	if st, err := main.StatusPorcelainTracked(); err != nil {
		return err
	} else if strings.TrimSpace(st) != "" {
		return failf(1, "integrate: the main worktree (%s) is not clean; another session may be working there:\n%s", main.Root, st)
	}

	lockOID, err := c.acquireIntegrationLock(*taskID, "", id)
	if err != nil {
		return err
	}
	released := false
	release := func() {
		if !released {
			_ = c.Repo.DeleteCAS(IntegrationLockRef, lockOID, fmt.Sprintf("queue integrate release %s", *taskID))
			released = true
		}
	}
	defer release()

	branch := rec.Task.CanonicalBranch
	branchOID, err := c.Repo.Run("rev-parse", "--verify", branch)
	if err != nil {
		return failf(1, "integrate: branch %s does not resolve", branch)
	}
	branchOID = strings.TrimSpace(branchOID)
	mainHead, err := main.Head()
	if err != nil {
		return err
	}
	switch {
	case main.IsAncestor(branchOID, mainHead):
		// already integrated
	case main.IsAncestor(mainHead, branchOID):
		if _, err := main.Run("merge", "--ff-only", branch); err != nil {
			return fmt.Errorf("integrate: fast-forward of main failed: %w", err)
		}
	default:
		if _, err := main.Run("merge", "--no-edit", branch); err != nil {
			_, _ = main.Run("merge", "--abort")
			return failf(1, "integrate: merging %s into main conflicted; resolve it in the main worktree and retry", branch)
		}
	}

	recMain, err := recordIn(main.Root, *taskID)
	if err != nil {
		return err
	}
	integrated, err := main.Head()
	if err != nil {
		return err
	}
	tag := rec.Task.RequiredTag
	msg := *tagMessage
	if msg == "" {
		msg = fmt.Sprintf("%s: integrated into local main from %s.", *taskID, branch)
	}
	ev, err := c.emitEvent(main, recMain, eventParams{
		EventType:      "integrated",
		State:          model.StateApprovedForIntegration,
		Actor:          id,
		Summary:        fmt.Sprintf("Integrated %s into local main at %s.", branch, integrated),
		ResultCommit:   integrated,
		Related:        model.Related{Tag: tag},
		NextCapability: model.CapHIGH,
		NextWorkRole:   model.RoleIntegrator,
		StatusNote:     "Integrated into local `main`. NEXT: `queue complete`.",
	})
	if err != nil {
		return err
	}
	final, err := main.Head()
	if err != nil {
		return err
	}
	if tag != "" && !main.TagExists(tag) {
		if err := main.CreateAnnotatedTag(tag, final, msg); err != nil {
			return err
		}
	}
	release()
	if cf.json {
		return c.jsonOut(map[string]any{"task_id": *taskID, "integrated_commit": integrated, "main_head": final, "tag": tag, "event": ev.EventID})
	}
	fmt.Fprintf(out, "integrated %s into main: %s (tag %s)\n", *taskID, final, tag)
	return nil
}

func cmdComplete(args []string, out io.Writer) error {
	fs := newFlagSet("complete")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	summary := fs.String("summary", "", "completion note")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "complete: %v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	if rec.State == model.StateCompleted {
		fmt.Fprintf(out, "%s is already completed\n", *taskID)
		return nil
	}
	if rec.State != model.StateApprovedForIntegration {
		return failf(1, "complete: task %s is %s, not approved_for_integration", *taskID, rec.State)
	}
	main, err := c.mainWorktree()
	if err != nil {
		return err
	}
	// Required checks before the task may be called complete.
	if p := submitResultPath(rec); p == "" {
		return failf(1, "complete: no attempts/*/RESULT.md exists for %s", *taskID)
	}
	if len(reviewDirs(rec)) == 0 {
		return failf(1, "complete: no reviews/*/REVIEW.md exists for %s", *taskID)
	}
	branchOID, err := c.Repo.Run("rev-parse", "--verify", rec.Task.CanonicalBranch)
	if err != nil {
		return failf(1, "complete: branch %s does not resolve", rec.Task.CanonicalBranch)
	}
	mainHead, _ := main.Head()
	if !main.IsAncestor(strings.TrimSpace(branchOID), mainHead) {
		return failf(1, "complete: branch %s is not reachable from local main", rec.Task.CanonicalBranch)
	}
	if rec.Task.RequiredTag != "" && !main.TagExists(rec.Task.RequiredTag) {
		return failf(1, "complete: required tag %s does not exist", rec.Task.RequiredTag)
	}
	sum := *summary
	if sum == "" {
		sum = fmt.Sprintf("Completed by %s; every required check passed and the branch is reachable from main.", id.Verbose())
	}
	recMain, err := recordIn(main.Root, *taskID)
	if err != nil {
		return err
	}
	ev, err := c.emitEvent(main, recMain, eventParams{
		EventType:      "completed",
		State:          model.StateCompleted,
		Actor:          id,
		Summary:        sum,
		ResultCommit:   mainHead,
		Related:        model.Related{Tag: rec.Task.RequiredTag},
		NextCapability: "",
		NextWorkRole:   "",
		StatusNote:     "Task complete. The immutable event files remain the record.",
	})
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(ev)
	}
	fmt.Fprintf(out, "completed %s\n", *taskID)
	return nil
}

// ---------------------------------------------------------------------------
// recover
// ---------------------------------------------------------------------------

func cmdRecover(args []string, out io.Writer) error {
	fs := newFlagSet("recover")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	leaseStr := fs.String("lease", "2h", "new lease duration")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if *taskID == "" {
		return failf(2, "recover: --task is required")
	}
	lease, err := model.ParseDuration(*leaseStr)
	if err != nil {
		return failf(2, "recover: --lease: %v", err)
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "recover: %v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, oid, err := c.readLiveClaim(*taskID)
	if err != nil {
		return err
	}
	if oid == "" || cl == nil {
		return failf(1, "recover: no live claim for %s", *taskID)
	}
	if model.Terminal(cl.State) {
		return failf(1, "recover: task %s is already %s", *taskID, cl.State)
	}
	exp, err := model.ParseTime(cl.ExpiresAt)
	if err != nil {
		return failf(1, "recover: live claim has a malformed expires_at")
	}
	if exp.After(c.Now) {
		return failf(1, "recover: lease for %s has not expired (expires %s)", *taskID, cl.ExpiresAt)
	}
	// The benchmark lock is the only measurement lock; recovery must not run
	// while this task's worktree still owns it.
	worktree := c.Repo.ResolveRepoPath(firstNonEmpty(cl.Worktree, ".worktrees/"+*taskID))
	if benchmarkLockHeldBy(worktree) {
		return failf(1, "recover: %s owns the benchmark lock (ads-run-lock); wait for the measurement to finish", *taskID)
	}

	if cl.RecoveryStartedAt == "" {
		// First compare-and-swap: mark recovery_pending and start the grace
		// period. The old claim id and epoch are preserved so a late heartbeat
		// from the old worker is still rejected only after the second step.
		cl.RecoveryStartedAt = model.FormatTime(c.Now)
		cl.RecoveryReadyAt = model.FormatTime(c.Now.Add(model.RecoveryGrace))
		if err := c.writeLiveClaim(cl, oid, fmt.Sprintf("queue recovery pending %s", *taskID)); err != nil {
			return fmt.Errorf("recover: lost the compare-and-swap: %w", err)
		}
		if cf.json {
			return c.jsonOut(cl)
		}
		fmt.Fprintf(out, "recover: %s is recovery_pending; grace period ends %s. Re-run to complete recovery.\n", *taskID, cl.RecoveryReadyAt)
		return nil
	}

	ready, err := model.ParseTime(cl.RecoveryReadyAt)
	if err != nil {
		return failf(1, "recover: malformed recovery_ready_at %q", cl.RecoveryReadyAt)
	}
	if c.Now.Before(ready) {
		return failf(1, "recover: grace period ends %s; retry after that", cl.RecoveryReadyAt)
	}
	// Second state check: the blob must not have changed since the first step.
	cl2, oid2, err := c.readLiveClaim(*taskID)
	if err != nil {
		return err
	}
	if oid2 != oid || cl2 == nil || cl2.RecoveryStartedAt != cl.RecoveryStartedAt {
		return failf(1, "recover: state changed during the grace period; re-run recovery")
	}

	branch := firstNonEmpty(cl.Branch, "")
	base := firstNonEmpty(cl.BaseCommit, "main")
	newCl := &model.Claim{
		SchemaVersion:   model.SchemaVersion,
		TaskID:          *taskID,
		State:           model.StateClaimed,
		ClaimID:         "claim-" + randHex(8),
		AttemptID:       "attempt-" + archive.CompactTime(c.Now) + "-" + randHex(3),
		ClaimEpoch:      cl.ClaimEpoch + 1,
		Worker:          id,
		ClaimedAt:       model.FormatTime(c.Now),
		HeartbeatAt:     model.FormatTime(c.Now),
		ExpiresAt:       model.FormatTime(c.Now.Add(lease)),
		Branch:          branch,
		Worktree:        cl.Worktree,
		RepoRoot:        cl.RepoRoot,
		BaseCommit:      base,
		TaskSpecDigest:  cl.TaskSpecDigest,
		PreviousClaimID: cl.ClaimID,
		NextCapability:  model.CapLOW,
		NextWorkRole:    model.RoleExecutor,
	}
	if err := c.writeLiveClaim(newCl, oid, fmt.Sprintf("queue recover %s epoch %d", *taskID, newCl.ClaimEpoch)); err != nil {
		return fmt.Errorf("recover: lost the compare-and-swap to another recoverer: %w", err)
	}
	// Recover into the committed chain too, so a fresh clone sees the epoch.
	if rec, err := c.taskRecord(*taskID); err == nil {
		ev, eerr := c.emitEvent(c.Repo, rec, eventParams{
			EventType:      "recovered",
			State:          model.StateClaimed,
			Actor:          id,
			AttemptID:      newCl.AttemptID,
			ClaimID:        newCl.ClaimID,
			Summary:        fmt.Sprintf("Recovered expired claim %s (epoch %d) by %s.", cl.ClaimID, cl.ClaimEpoch, id.Verbose()),
			BaseCommit:     base,
			NextCapability: model.CapLOW,
			NextWorkRole:   model.RoleExecutor,
			StatusNote:     fmt.Sprintf("Recovered by %s as claim %s epoch %d until %s. Existing worktree files preserved.", id.Verbose(), newCl.ClaimID, newCl.ClaimEpoch, newCl.ExpiresAt),
		})
		if eerr != nil {
			return eerr
		}
		newCl.LastEvent = ev.EventID
		oidNow, _ := c.Repo.RefOID(liveRef(*taskID))
		if err := c.writeLiveClaim(newCl, oidNow, fmt.Sprintf("queue recover %s record event", *taskID)); err != nil {
			return err
		}
	}
	// Reuse and preserve the existing worktree; create it only if it is gone.
	if _, err := c.ensureWorktree(&model.Task{CanonicalBranch: branch, CanonicalWorktree: newCl.Worktree}, base); err != nil {
		return err
	}
	if err := c.saveClaimContextFor(newCl.Worktree, newCl); err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(newCl)
	}
	fmt.Fprintf(out, "recovered %s: claim %s epoch %d expires %s\n", *taskID, newCl.ClaimID, newCl.ClaimEpoch, newCl.ExpiresAt)
	return nil
}
