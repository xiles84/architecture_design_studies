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

// requireCurrentClaim is the guard rule as a reusable precondition: the live ref
// must exist, must still name this worker's claim id and epoch, the task spec
// must be unchanged, and the lease must not have expired. A failure is exit
// code 3 so a shell script can stop rather than continue with a stale claim.
func (c *cmdCtx) requireCurrentClaim(taskID, claimIDFlag string, epochFlag int) (*model.Claim, string, error) {
	cl, oid, err := c.readLiveClaim(taskID)
	if err != nil {
		return nil, "", err
	}
	if oid == "" || cl == nil {
		return nil, "", failf(3, "no live claim for %s", taskID)
	}
	if model.Terminal(cl.State) {
		return nil, "", failf(3, "claim for %s is in terminal state %s", taskID, cl.State)
	}
	claimID, epoch, err := c.resolveClaim(taskID, claimIDFlag, epochFlag)
	if err != nil {
		return nil, "", err
	}
	if cl.ClaimID != claimID || cl.ClaimEpoch != epoch {
		return nil, "", failf(3, "claim superseded: live is %s epoch %d, this worker holds %s epoch %d", cl.ClaimID, cl.ClaimEpoch, claimID, epoch)
	}
	rec, err := c.taskRecord(taskID)
	if err != nil {
		return nil, "", err
	}
	if _, raw, err := archive.LoadTask(rec.Dir); err == nil {
		if d := model.TaskDigest(raw); d != cl.TaskSpecDigest {
			return nil, "", failf(3, "task spec digest changed (%s != %s); stop and re-read the brief", cl.TaskSpecDigest, d)
		}
	}
	exp, err := model.ParseTime(cl.ExpiresAt)
	if err != nil {
		return nil, "", failf(3, "live claim has a malformed expires_at %q", cl.ExpiresAt)
	}
	if exp.Before(c.Now) {
		return nil, "", failf(3, "lease expired at %s; run `queue recover`", cl.ExpiresAt)
	}
	return cl, oid, nil
}

func (c *cmdCtx) dropClaimContext(taskID string) {
	_ = os.Remove(c.claimContextPath(taskID))
}

// ---------------------------------------------------------------------------
// claim
// ---------------------------------------------------------------------------

func cmdClaim(args []string, out io.Writer) error {
	fs := newFlagSet("claim")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	base := fs.String("base", "main", "base revision for the task branch")
	leaseStr := fs.String("lease", "2h", "lease duration")
	noWorktree := fs.Bool("no-worktree", false, "do not create the canonical worktree")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if *taskID == "" {
		return failf(2, "claim: --task is required")
	}
	lease, err := model.ParseDuration(*leaseStr)
	if err != nil {
		return failf(2, "claim: --lease: %v", err)
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "claim: %v (pass --capability and --role)", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	if rec.State != model.StateReady && rec.State != model.StateChangesRequested {
		return failf(1, "task %s is %s, not claimable", *taskID, rec.State)
	}
	if !model.CapabilityEligible(id.SessionCapability, rec.Task.MinimumCapability) {
		return failf(1, "task %s requires %s capability; this session is %s", *taskID, rec.Task.MinimumCapability, id.SessionCapability)
	}
	if nb, err := model.ParseTime(rec.Task.NotBefore); err != nil || nb.After(c.Now) {
		return failf(1, "task %s is not_before %s", *taskID, rec.Task.NotBefore)
	}
	if missing := missingDependencies(c, rec.Task); len(missing) > 0 {
		return failf(1, "task %s has incomplete dependencies: %s", *taskID, strings.Join(missing, ", "))
	}

	oldClaim, oldOID, err := c.readLiveClaim(*taskID)
	if err != nil {
		return err
	}
	if oldOID != "" && oldClaim != nil && !model.Terminal(oldClaim.State) {
		exp, _ := model.ParseTime(oldClaim.ExpiresAt)
		if exp.After(c.Now) {
			return failf(1, "task %s is claimed by %s (claim %s, epoch %d) until %s",
				*taskID, oldClaim.Worker.Verbose(), oldClaim.ClaimID, oldClaim.ClaimEpoch, oldClaim.ExpiresAt)
		}
		return failf(1, "lease for task %s expired at %s; recovery is required (`queue recover --task %s`)", *taskID, oldClaim.ExpiresAt, *taskID)
	}

	baseOID, err := c.Repo.Run("rev-parse", "--verify", *base)
	if err != nil {
		return failf(1, "claim: base revision %q does not resolve", *base)
	}
	baseOID = strings.TrimSpace(baseOID)

	epoch := 1
	if oldClaim != nil {
		epoch = oldClaim.ClaimEpoch + 1
	}
	cl := &model.Claim{
		SchemaVersion:  model.SchemaVersion,
		TaskID:         *taskID,
		State:          model.StateClaimed,
		ClaimID:        "claim-" + randHex(8),
		AttemptID:      "attempt-" + archive.CompactTime(c.Now) + "-" + randHex(3),
		ClaimEpoch:     epoch,
		Worker:         id,
		ClaimedAt:      model.FormatTime(c.Now),
		HeartbeatAt:    model.FormatTime(c.Now),
		ExpiresAt:      model.FormatTime(c.Now.Add(lease)),
		Branch:         rec.Task.CanonicalBranch,
		Worktree:       rec.Task.CanonicalWorktree,
		RepoRoot:       filepath.ToSlash(c.Repo.CommonRoot()),
		BaseCommit:     baseOID,
		TaskSpecDigest: model.TaskDigest(mustBytes(c, rec.Task.TaskID)),
		NextCapability: model.CapLOW,
		NextWorkRole:   model.RoleExecutor,
	}
	if oldClaim != nil {
		cl.PreviousClaimID = oldClaim.ClaimID
	}
	if err := c.writeLiveClaim(cl, oldOID, fmt.Sprintf("queue claim %s by %s", *taskID, id.Verbose())); err != nil {
		return fmt.Errorf("claim lost to another worker: %w", err)
	}

	ev, err := c.emitEvent(c.Repo, rec, eventParams{
		EventType:      "claimed",
		State:          model.StateClaimed,
		Actor:          id,
		AttemptID:      cl.AttemptID,
		ClaimID:        cl.ClaimID,
		Summary:        fmt.Sprintf("Claimed by %s (epoch %d, lease %s).", id.Verbose(), cl.ClaimEpoch, *leaseStr),
		BaseCommit:     baseOID,
		NextCapability: model.CapLOW,
		NextWorkRole:   model.RoleExecutor,
		StatusNote:     fmt.Sprintf("Claimed by %s until %s. Worktree `%s`.", id.Verbose(), cl.ExpiresAt, rec.Task.CanonicalWorktree),
	})
	if err != nil {
		return err
	}
	rec.AppendEvent(ev)
	cl.LastEvent = ev.EventID

	if !*noWorktree {
		if _, err := c.ensureWorktree(rec.Task, *base); err != nil {
			return fmt.Errorf("claim %s is held but the worktree could not be created: %w (it is recoverable; run `queue recover`)", cl.ClaimID, err)
		}
	}
	if err := c.saveClaimContextFor(rec.Task.CanonicalWorktree, cl); err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(cl)
	}
	fmt.Fprintf(out, "claimed %s\n  claim_id:  %s\n  epoch:     %d\n  attempt:   %s\n  worktree:  %s\n  branch:    %s\n  expires:   %s\n",
		*taskID, cl.ClaimID, cl.ClaimEpoch, cl.AttemptID, rec.Task.CanonicalWorktree, rec.Task.CanonicalBranch, cl.ExpiresAt)
	return nil
}

func mustBytes(c *cmdCtx, taskID string) []byte {
	rec, err := c.taskRecord(taskID)
	if err != nil {
		return nil
	}
	_, raw, err := archive.LoadTask(rec.Dir)
	if err != nil {
		return nil
	}
	return raw
}

func missingDependencies(c *cmdCtx, t *model.Task) []string {
	if len(t.Dependencies) == 0 {
		return nil
	}
	all, _ := c.allRecords()
	byID := map[string]*model.TaskRecord{}
	for _, r := range all {
		byID[r.Task.TaskID] = r
	}
	var missing []string
	for _, dep := range t.Dependencies {
		r, ok := byID[dep]
		if !ok || r.State != model.StateCompleted {
			missing = append(missing, dep)
		}
	}
	return missing
}

// lastWorker returns the identity that holds (or last held) a task: the live
// claim if present, otherwise the last committed event's actor.
func lastWorker(rec *model.TaskRecord, cl *model.Claim) model.Identity {
	if cl != nil {
		return cl.Worker
	}
	if ev := rec.LastEvent(); ev != nil {
		return ev.Actor
	}
	return model.Identity{}
}

// isRegisteredWorktree reports whether path is a working tree of the same
// repository. It is deliberately filesystem- and git-based rather than relying
// on `git worktree list`: this repository stores the linked worktree's
// back-pointer as a relative path (portable between WSL Git and native Windows
// Git), which Git 2.43 cannot resolve in `worktree list` and reports as
// prunable. Operating inside the worktree itself still works.
func (c *cmdCtx) isRegisteredWorktree(path string) bool {
	if _, err := os.Stat(filepath.Join(path, ".git")); err != nil {
		return false
	}
	wt, err := gitx.Open(path)
	if err != nil {
		return false
	}
	return samePath(wt.CommonDir, c.Repo.CommonDir)
}

// ensureWorktree creates or reuses the canonical worktree. It never discards an
// existing working folder: that is where a crashed worker's uncommitted files
// live.
func (c *cmdCtx) ensureWorktree(t *model.Task, base string) (string, error) {
	path := c.Repo.ResolveCommonPath(t.CanonicalWorktree)
	if c.isRegisteredWorktree(path) {
		return path, nil
	}
	if _, err := os.Stat(path); err == nil {
		// A folder exists but git does not register it as a work tree. Reusing
		// it blindly risks confusing a plain directory with a worktree; report
		// it rather than deleting anything.
		return "", fmt.Errorf("%s exists but is not a registered work tree; inspect it (files are preserved) and repair or remove it deliberately", path)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return "", err
	}
	if c.Repo.BranchExists(t.CanonicalBranch) {
		if _, err := c.Repo.Run("worktree", "add", path, t.CanonicalBranch); err != nil {
			return "", err
		}
	} else if err := c.Repo.WorktreeAdd(path, t.CanonicalBranch, base); err != nil {
		return "", err
	}
	if err := c.Repo.MakeWorktreePathsRelative(path); err != nil {
		return "", fmt.Errorf("worktree created but making its registration portable failed: %w", err)
	}
	return path, nil
}

// ---------------------------------------------------------------------------
// guard
// ---------------------------------------------------------------------------

func cmdGuard(args []string, out io.Writer) error {
	fs := newFlagSet("guard")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	claimID := fs.String("claim-id", "", "claim id (default: this worktree's saved claim)")
	epoch := fs.Int("claim-epoch", 0, "claim epoch (default: this worktree's saved claim)")
	quiet := fs.Bool("quiet", false, "print nothing on success")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if *taskID == "" {
		return failf(2, "guard: --task is required")
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, _, err := c.requireCurrentClaim(*taskID, *claimID, *epoch)
	if err != nil {
		return err
	}
	if !*quiet {
		fmt.Fprintf(out, "guard ok: %s claim %s epoch %d expires %s\n", *taskID, cl.ClaimID, cl.ClaimEpoch, cl.ExpiresAt)
	}
	return nil
}

// ---------------------------------------------------------------------------
// heartbeat
// ---------------------------------------------------------------------------

func cmdHeartbeat(args []string, out io.Writer) error {
	fs := newFlagSet("heartbeat")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	claimID := fs.String("claim-id", "", "claim id")
	epoch := fs.Int("claim-epoch", 0, "claim epoch")
	leaseStr := fs.String("lease", "2h", "new lease duration")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	lease, err := model.ParseDuration(*leaseStr)
	if err != nil {
		return failf(2, "heartbeat: --lease: %v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, oid, err := c.requireCurrentClaim(*taskID, *claimID, *epoch)
	if err != nil {
		return err
	}
	cl.HeartbeatAt = model.FormatTime(c.Now)
	cl.ExpiresAt = model.FormatTime(c.Now.Add(lease))
	if err := c.writeLiveClaim(cl, oid, fmt.Sprintf("queue heartbeat %s claim %s", *taskID, cl.ClaimID)); err != nil {
		return fmt.Errorf("heartbeat rejected: %w", err)
	}
	if err := c.saveClaimContext(cl); err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(cl)
	}
	fmt.Fprintf(out, "heartbeat %s: claim %s epoch %d expires %s\n", *taskID, cl.ClaimID, cl.ClaimEpoch, cl.ExpiresAt)
	return nil
}

// ---------------------------------------------------------------------------
// checkpoint
// ---------------------------------------------------------------------------

func cmdCheckpoint(args []string, out io.Writer) error {
	fs := newFlagSet("checkpoint")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	claimID := fs.String("claim-id", "", "claim id")
	epoch := fs.Int("claim-epoch", 0, "claim epoch")
	summary := fs.String("summary", "", "what this checkpoint recorded")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, oid, err := c.requireCurrentClaim(*taskID, *claimID, *epoch)
	if err != nil {
		return err
	}
	id := cf.actorFrom(cl.Worker)
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	head, err := c.Repo.Head()
	if err != nil {
		return err
	}
	newState := rec.State
	if rec.State == model.StateClaimed {
		newState = model.StateInProgress
	}
	sum := *summary
	if sum == "" {
		sum = fmt.Sprintf("Checkpoint at %s by %s.", head, id.Verbose())
	}
	ev, err := c.emitEvent(c.Repo, rec, eventParams{
		EventType:        "checkpoint",
		State:            newState,
		Actor:            id,
		AttemptID:        cl.AttemptID,
		ClaimID:          cl.ClaimID,
		Summary:          sum,
		CheckpointCommit: head,
		BaseCommit:       cl.BaseCommit,
		NextCapability:   model.CapLOW,
		NextWorkRole:     model.RoleExecutor,
		StatusNote:       fmt.Sprintf("Checkpoint %s; claim %s epoch %d.", head, cl.ClaimID, cl.ClaimEpoch),
	})
	if err != nil {
		return err
	}
	cl.State = newState
	cl.CheckpointCommit = head
	cl.LastEvent = ev.EventID
	if err := c.writeLiveClaim(cl, oid, fmt.Sprintf("queue checkpoint %s claim %s", *taskID, cl.ClaimID)); err != nil {
		return err
	}
	if err := c.saveClaimContext(cl); err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(cl)
	}
	fmt.Fprintf(out, "checkpoint %s: %s (%s)\n", *taskID, head, newState)
	return nil
}

// ---------------------------------------------------------------------------
// yield
// ---------------------------------------------------------------------------

func cmdYield(args []string, out io.Writer) error {
	fs := newFlagSet("yield")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	claimID := fs.String("claim-id", "", "claim id")
	epoch := fs.Int("claim-epoch", 0, "claim epoch")
	reason := fs.String("reason", "", "why the claim is being released")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	cl, oid, err := c.readLiveClaim(*taskID)
	if err != nil {
		return err
	}
	id := cf.actorFrom(lastWorker(rec, cl))
	if err := id.Validate(); err != nil {
		return failf(2, "yield: %v", err)
	}
	// blocked_high has no live claim (escalate released it); every other state
	// must still be this worker's claim.
	if rec.State != model.StateBlockedHigh || oid != "" {
		if _, _, err := c.requireCurrentClaim(*taskID, *claimID, *epoch); err != nil {
			return err
		}
	}
	sum := *reason
	if sum == "" {
		sum = fmt.Sprintf("Released to ready by %s.", id.Verbose())
	}
	ev, err := c.emitEvent(c.Repo, rec, eventParams{
		EventType:      "yielded",
		State:          model.StateReady,
		Actor:          id,
		ClaimID:        claimIDOf(cl),
		Summary:        sum,
		NextCapability: nextCapabilityForTask(rec.Task),
		NextWorkRole:   model.RoleExecutor,
		StatusNote:     "Returned to the ready pool. " + sum,
	})
	if err != nil {
		return err
	}
	if err := c.releaseLiveClaim(*taskID, oid, fmt.Sprintf("queue yield %s", *taskID)); err != nil {
		return err
	}
	c.dropClaimContext(*taskID)
	if cf.json {
		return c.jsonOut(ev)
	}
	fmt.Fprintf(out, "yielded %s -> ready (%s)\n", *taskID, ev.EventID)
	return nil
}

func claimIDOf(cl *model.Claim) string {
	if cl == nil {
		return ""
	}
	return cl.ClaimID
}

// ---------------------------------------------------------------------------
// escalate
// ---------------------------------------------------------------------------

func cmdEscalate(args []string, out io.Writer) error {
	fs := newFlagSet("escalate")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	claimID := fs.String("claim-id", "", "claim id")
	epoch := fs.Int("claim-epoch", 0, "claim epoch")
	decision := fs.String("decision-needed", "", "the decision HIGH must make (required)")
	escalationID := fs.String("escalation-id", "", "escalation id")
	body := fs.String("detail", "", "facts, options and what work continues")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if *decision == "" {
		return failf(2, "escalate: --decision-needed is required")
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, oid, err := c.requireCurrentClaim(*taskID, *claimID, *epoch)
	if err != nil {
		return err
	}
	id := cf.actorFrom(cl.Worker)
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	escID := *escalationID
	if escID == "" {
		escID = "ESC-" + archive.CompactTime(c.Now) + "-" + randHex(2)
	}
	attemptDir := filepath.Join(rec.Dir, "attempts", cl.AttemptID)
	if err := os.MkdirAll(attemptDir, 0o755); err != nil {
		return err
	}
	escPath := filepath.Join(attemptDir, "ESCALATIONS.md")
	entry := fmt.Sprintf("\n## %s — ESCALATION REQUIRED\n\n**Decision needed:** %s\n\n**Raised by:** %s at %s\n\n%s\n",
		escID, *decision, id.Verbose(), model.FormatTime(c.Now), strings.TrimSpace(*body))
	f, err := os.OpenFile(escPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.WriteString(entry); err != nil {
		f.Close()
		return err
	}
	f.Close()
	relEsc, err := c.Repo.RepoRelative(escPath)
	if err != nil {
		return err
	}
	head, _ := c.Repo.Head()
	ev, err := c.emitEvent(c.Repo, rec, eventParams{
		EventType:        "escalated",
		State:            model.StateBlockedHigh,
		Actor:            id,
		AttemptID:        cl.AttemptID,
		ClaimID:          cl.ClaimID,
		Summary:          "ESCALATION REQUIRED: " + *decision,
		CheckpointCommit: head,
		EvidencePaths:    []string{relEsc},
		Related:          model.Related{EscalationID: escID},
		NextCapability:   model.CapHIGH,
		NextWorkRole:     model.RolePlanner,
		StatusNote:       "Blocked on HIGH. See `" + relEsc + "`.",
	})
	if err != nil {
		return err
	}
	if err := c.releaseLiveClaim(*taskID, oid, fmt.Sprintf("queue escalate %s", *taskID)); err != nil {
		return err
	}
	c.dropClaimContext(*taskID)
	if cf.json {
		return c.jsonOut(ev)
	}
	fmt.Fprintf(out, "escalated %s -> blocked_high (%s); HIGH must answer %s\n", *taskID, escID, *decision)
	return nil
}

// ---------------------------------------------------------------------------
// submit
// ---------------------------------------------------------------------------

func cmdSubmit(args []string, out io.Writer) error {
	fs := newFlagSet("submit")
	var cf commonFlags
	cf.register(fs, true)
	taskID := fs.String("task", "", "task id (required)")
	claimID := fs.String("claim-id", "", "claim id")
	epoch := fs.Int("claim-epoch", 0, "claim epoch")
	resultFlag := fs.String("result", "", "path to RESULT.md (default: this attempt's RESULT.md)")
	summary := fs.String("summary", "", "what was delivered")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	cl, oid, err := c.requireCurrentClaim(*taskID, *claimID, *epoch)
	if err != nil {
		return err
	}
	rec, err := c.taskRecord(*taskID)
	if err != nil {
		return err
	}
	if rec.State != model.StateClaimed && rec.State != model.StateInProgress {
		return failf(1, "submit: task %s is %s, not submittable", *taskID, rec.State)
	}
	resultPath := *resultFlag
	if resultPath == "" {
		resultPath = filepath.Join(rec.Dir, "attempts", cl.AttemptID, "RESULT.md")
	} else {
		resultPath = briefMustAbs(resultPath)
	}
	if _, err := os.Stat(resultPath); err != nil {
		return failf(1, "submit: result file %s does not exist; write and commit it first", resultPath)
	}
	// The worker's own work must be committed: the event is the durable record,
	// and an uncommitted tree would leave the result outside the archive.
	if st, err := c.Repo.StatusPorcelain(); err != nil {
		return err
	} else if strings.TrimSpace(st) != "" {
		return failf(1, "submit: working tree is not clean; commit or remove:\n%s", st)
	}
	relResult, err := c.Repo.RepoRelative(resultPath)
	if err != nil {
		return err
	}
	head, _ := c.Repo.Head()
	id := cf.actorFrom(cl.Worker)
	sum := *summary
	if sum == "" {
		sum = fmt.Sprintf("Submitted for review by %s (%s).", id.Verbose(), head)
	}
	ev, err := c.emitEvent(c.Repo, rec, eventParams{
		EventType:      "submitted",
		State:          model.StateAwaitingReview,
		Actor:          id,
		AttemptID:      cl.AttemptID,
		ClaimID:        cl.ClaimID,
		Summary:        sum,
		ResultCommit:   head,
		EvidencePaths:  []string{relResult},
		NextCapability: model.CapHIGH,
		NextWorkRole:   model.RoleReviewer,
		StatusNote:     "Awaiting HIGH review. Result: `" + relResult + "`.",
	})
	if err != nil {
		return err
	}
	if err := c.releaseLiveClaim(*taskID, oid, fmt.Sprintf("queue submit %s", *taskID)); err != nil {
		return err
	}
	c.dropClaimContext(*taskID)
	if cf.json {
		return c.jsonOut(ev)
	}
	fmt.Fprintf(out, "submitted %s for review (%s); result_commit=%s\n", *taskID, ev.EventID, head)
	return nil
}

// emitJSON is a tiny helper used by tests and JSON output paths.
func emitJSON(out io.Writer, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, string(data))
	return err
}
