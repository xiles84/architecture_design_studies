package app

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"adsqueue/internal/archive"
	"adsqueue/internal/model"
)

func newFlagSet(name string) *flag.FlagSet {
	fs := flag.NewFlagSet(name, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	return fs
}

// ---------------------------------------------------------------------------
// session-start
// ---------------------------------------------------------------------------

type sessionRecord struct {
	SchemaVersion          int    `json:"schema_version"`
	SessionID              string `json:"session_id"`
	Model                  string `json:"model"`
	Tool                   string `json:"tool"`
	Effort                 string `json:"effort"`
	SessionCapability      string `json:"session_capability"`
	SessionCapabilityInput string `json:"session_capability_input,omitempty"`
	WorkRole               string `json:"work_role"`
	StartedAt              string `json:"started_at"`
}

// cmdSessionStart records the answer to "is this session HIGH or LOW?". The
// workflow requires that question be asked once, not inferred from a model
// name, so the answer is an explicit input.
func cmdSessionStart(args []string, out io.Writer) error {
	fs := newFlagSet("session-start")
	var cf commonFlags
	cf.register(fs, true)
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "session-start: %v (pass --capability HIGH|LOW and --role <role>)", err)
	}
	if id.SessionID == "" || id.SessionID == "unknown" {
		id.SessionID = "session-" + archive.CompactTime(c.Now) + "-" + randHex(3)
	}
	rec := sessionRecord{
		SchemaVersion:          model.SchemaVersion,
		SessionID:              id.SessionID,
		Model:                  id.Model,
		Tool:                   id.Tool,
		Effort:                 id.Effort,
		SessionCapability:      id.SessionCapability,
		SessionCapabilityInput: id.SessionCapabilityInput,
		WorkRole:               id.WorkRole,
		StartedAt:              model.FormatTime(c.Now),
	}
	data, err := json.MarshalIndent(rec, "", "  ")
	if err != nil {
		return err
	}
	rel := filepath.ToSlash(filepath.Join(archive.SessionsRoot, id.SessionID+".json"))
	path := filepath.Join(c.Repo.Root, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return err
	}
	msg := fmt.Sprintf("Queue: record session %s (%s, %s)\n\nAgent: %s\n", id.SessionID, id.SessionCapability, id.WorkRole, id.Verbose())
	if _, err := c.Repo.CommitPaths(msg, rel); err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(rec)
	}
	fmt.Fprintf(out, "session %s: capability=%s role=%s model=%s tool=%s session_id=%s\n",
		id.SessionID, id.SessionCapability, id.WorkRole, id.Model, id.Tool, id.SessionID)
	return nil
}

// ---------------------------------------------------------------------------
// list / next / review-candidates
// ---------------------------------------------------------------------------

type taskRow struct {
	TaskID            string   `json:"task_id"`
	State             string   `json:"state"`
	Priority          int      `json:"priority"`
	MinimumCapability string   `json:"minimum_capability"`
	NotBefore         string   `json:"not_before"`
	Dependencies      []string `json:"dependencies"`
	Title             string   `json:"title"`
	GoalID            string   `json:"goal_id"`
	RootTaskID        string   `json:"root_task_id"`
	ParentTaskID      string   `json:"parent_task_id,omitempty"`
	Sequence          int      `json:"sequence"`
	Eligible          bool     `json:"eligible_now"`
	CanonicalBranch   string   `json:"canonical_branch"`
	CanonicalWorktree string   `json:"canonical_worktree"`
}

func rows(records []*model.TaskRecord, now interface{ Unix() int64 }, capability string) []taskRow {
	out := make([]taskRow, 0, len(records))
	for _, r := range records {
		row := taskRow{
			TaskID:            r.Task.TaskID,
			State:             r.State,
			Priority:          r.Task.Priority,
			MinimumCapability: r.Task.MinimumCapability,
			NotBefore:         r.Task.NotBefore,
			Dependencies:      r.Task.Dependencies,
			Title:             r.Task.Title,
			GoalID:            r.Task.GoalID,
			RootTaskID:        r.Task.RootTaskID,
			ParentTaskID:      r.Task.ParentTaskID,
			Sequence:          len(r.Events),
			CanonicalBranch:   r.Task.CanonicalBranch,
			CanonicalWorktree: r.Task.CanonicalWorktree,
		}
		out = append(out, row)
	}
	return out
}

func cmdList(args []string, out io.Writer) error {
	fs := newFlagSet("list")
	var cf commonFlags
	cf.register(fs, true)
	state := fs.String("state", "", "only show tasks in this state")
	eligible := fs.Bool("eligible", false, "only show tasks this capability may claim now")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	records, loadErrs := c.allRecords()
	now := c.Now
	var list []*model.TaskRecord
	for _, r := range records {
		if *state != "" && r.State != *state {
			continue
		}
		list = append(list, r)
	}
	eligSet := map[string]bool{}
	capability := cf.canonicalCapability()
	if *eligible || capability != "" {
		for _, r := range model.Eligible(records, now, capability) {
			eligSet[r.Task.TaskID] = true
		}
	}
	if *eligible {
		var filtered []*model.TaskRecord
		for _, r := range list {
			if eligSet[r.Task.TaskID] {
				filtered = append(filtered, r)
			}
		}
		list = filtered
	}
	_ = loadErrs
	rs := rows(list, now, capability)
	for i := range rs {
		rs[i].Eligible = eligSet[rs[i].TaskID]
	}
	if cf.json {
		return c.jsonOut(rs)
	}
	fmt.Fprintf(out, "%-42s %-24s %4s %-5s %-20s %-30s %s\n", "TASK", "STATE", "PRI", "MIN", "NOT_BEFORE", "TITLE", "ELIGIBLE")
	for _, r := range rs {
		el := ""
		if r.Eligible {
			el = "yes"
		}
		fmt.Fprintf(out, "%-42s %-24s %4d %-5s %-20s %-30s %s\n",
			r.TaskID, r.State, r.Priority, r.MinimumCapability, r.NotBefore, truncate(r.Title, 30), el)
	}
	if len(rs) == 0 {
		fmt.Fprintln(out, "(no tasks)")
	}
	return nil
}

func cmdNext(args []string, out io.Writer) error {
	fs := newFlagSet("next")
	var cf commonFlags
	cf.register(fs, true)
	claimCmd := fs.Bool("claim-command", false, "also print the canonical claim command")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	capability := cf.canonicalCapability()
	if capability == "" {
		return failf(2, "next: --capability HIGH|LOW is required (or set ADS_QUEUE_CAPABILITY)")
	}
	if capability != model.CapHIGH && capability != model.CapLOW {
		return failf(2, "next: %v", model.Identity{SessionCapability: capability, WorkRole: model.RoleExecutor}.Validate())
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	records, loadErrs := c.allRecords()
	if len(loadErrs) > 0 {
		return fmt.Errorf("cannot compute the queue: %v", loadErrs[0])
	}
	rec := model.Next(records, c.Now, capability)
	if rec == nil {
		if cf.json {
			return c.jsonOut(map[string]any{"task_id": nil})
		}
		fmt.Fprintln(out, "(no eligible task)")
		return nil
	}
	if cf.json {
		return c.jsonOut(taskRow{
			TaskID: rec.Task.TaskID, State: rec.State, Priority: rec.Task.Priority,
			MinimumCapability: rec.Task.MinimumCapability, NotBefore: rec.Task.NotBefore,
			Dependencies: rec.Task.Dependencies, Title: rec.Task.Title,
			CanonicalBranch: rec.Task.CanonicalBranch, CanonicalWorktree: rec.Task.CanonicalWorktree,
		})
	}
	fmt.Fprintf(out, "%s\t%s\t%d\t%s\n", rec.Task.TaskID, rec.State, rec.Task.Priority, rec.Task.Title)
	if *claimCmd {
		fmt.Fprintf(out, "claim: queue claim --task %s --capability %s --role %s\n",
			rec.Task.TaskID, capability, firstNonEmpty(cf.role, model.RoleExecutor))
	}
	return nil
}

func cmdReviewCandidates(args []string, out io.Writer) error {
	fs := newFlagSet("review-candidates")
	var cf commonFlags
	cf.register(fs, false)
	since := fs.String("since", "7d", "only candidates submitted within this duration (e.g. 7d, 36h)")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	dur, err := model.ParseDuration(*since)
	if err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	records, loadErrs := c.allRecords()
	if len(loadErrs) > 0 {
		return fmt.Errorf("cannot compute review candidates: %v", loadErrs[0])
	}
	cutoff := c.Now.Add(-dur)
	type candidate struct {
		TaskID       string `json:"task_id"`
		SubmittedAt  string `json:"submitted_at"`
		EventID      string `json:"event_id"`
		Title        string `json:"title"`
		ReviewPolicy string `json:"review_policy"`
	}
	var cands []candidate
	for _, r := range records {
		if r.State != model.StateAwaitingReview {
			continue
		}
		ev := r.LastEvent()
		if ev == nil {
			continue
		}
		at, err := model.ParseTime(ev.OccurredAt)
		if err != nil || at.Before(cutoff) {
			continue
		}
		cands = append(cands, candidate{TaskID: r.Task.TaskID, SubmittedAt: ev.OccurredAt, EventID: ev.EventID, Title: r.Task.Title, ReviewPolicy: r.Task.ReviewPolicy})
	}
	sort.SliceStable(cands, func(i, j int) bool { return cands[i].SubmittedAt < cands[j].SubmittedAt })
	if cf.json {
		return c.jsonOut(cands)
	}
	if len(cands) == 0 {
		fmt.Fprintf(out, "(no review candidates in the last %s)\n", *since)
		return nil
	}
	for _, cand := range cands {
		fmt.Fprintf(out, "%s\t%s\t%s\n", cand.TaskID, cand.SubmittedAt, cand.Title)
	}
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

// ---------------------------------------------------------------------------
// publish
// ---------------------------------------------------------------------------

func cmdPublish(args []string, out io.Writer) error {
	fs := newFlagSet("publish")
	var cf commonFlags
	cf.register(fs, true)
	spec := fs.String("spec", "", "path to the task.json to publish (required unless --release)")
	brief := fs.String("brief", "", "path to BRIEF.md (required unless --release)")
	state := fs.String("state", model.StateReady, "resulting state for a new task: ready or proposed")
	release := fs.Bool("release", false, "release an existing proposed task to ready")
	taskID := fs.String("task", "", "existing task id (with --release)")
	summary := fs.String("summary", "", "release summary")
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	if *state != model.StateReady && *state != model.StateProposed {
		return failf(2, "publish: --state must be ready or proposed")
	}
	id := cf.identity()
	if err := id.Validate(); err != nil {
		return failf(2, "publish: %v", err)
	}
	if id.SessionCapability != model.CapHIGH {
		return failf(2, "publish: only a HIGH session publishes or releases tasks (got %s)", id.SessionCapability)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}

	// Releasing a pre-created `proposed` task is the `proposed -> ready` arrow
	// of the state machine. Whole delivery tasks are published once, as
	// proposed, and released when their turn arrives; without this, no command
	// could move them and the queue could never advance past its first task.
	if *release {
		if *taskID == "" {
			return failf(2, "publish --release: --task is required")
		}
		rec, err := c.taskRecord(*taskID)
		if err != nil {
			return err
		}
		if rec.State != model.StateProposed {
			return failf(1, "publish --release: task %s is %s, not proposed", *taskID, rec.State)
		}
		sum := *summary
		if sum == "" {
			sum = fmt.Sprintf("Released to ready by %s.", id.Verbose())
		}
		ev, err := c.emitEvent(c.Repo, rec, eventParams{
			EventType:      "released",
			State:          model.StateReady,
			Actor:          id,
			Summary:        sum,
			NextCapability: nextCapabilityForTask(rec.Task),
			NextWorkRole:   rec.Task.WorkRole,
			StatusNote:     "Released to ready. " + sum,
		})
		if err != nil {
			return err
		}
		if cf.json {
			return c.jsonOut(ev)
		}
		fmt.Fprintf(out, "released %s -> ready (%s)\n", *taskID, ev.EventID)
		return nil
	}

	if *spec == "" || *brief == "" {
		return failf(2, "publish: --spec and --brief are required (or use --release --task)")
	}
	specPath, err := filepath.Abs(*spec)
	if err != nil {
		return err
	}
	specBytes, err := os.ReadFile(specPath)
	if err != nil {
		return err
	}
	var t model.Task
	if err := json.Unmarshal(specBytes, &t); err != nil {
		return fmt.Errorf("%s: %w", specPath, err)
	}
	if err := t.Validate(); err != nil {
		return err
	}
	dir, err := archive.TaskDirAbs(c.Repo.Root, t.TaskID)
	if err != nil {
		return err
	}
	if _, err := os.Stat(filepath.Join(dir, "task.json")); err == nil {
		return failf(1, "task %s is already published at %s; use an amendment or a superseding task", t.TaskID, dir)
	}
	briefBytes, err := os.ReadFile(briefMustAbs(*brief))
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Join(dir, "events"), 0o755); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "task.json"), ensureNL(specBytes), 0o644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(dir, "BRIEF.md"), ensureNL(briefBytes), 0o644); err != nil {
		return err
	}

	rec := &model.TaskRecord{Task: &t, Dir: dir, State: model.StateProposed}
	taskRel, _ := c.Repo.RepoRelative(filepath.Join(dir, "task.json"))
	briefRel, _ := c.Repo.RepoRelative(filepath.Join(dir, "BRIEF.md"))
	_, err = c.emitEvent(c.Repo, rec, eventParams{
		EventType:      "published",
		State:          *state,
		Actor:          id,
		Summary:        fmt.Sprintf("Published %s as %s.", t.TaskID, *state),
		BaseCommit:     headOrEmpty(c),
		EvidencePaths:  []string{taskRel, briefRel, "docs/ai-work/WORKFLOW.md", "docs/ai-work/SCHEMA.md"},
		NextCapability: nextCapabilityForTask(&t),
		NextWorkRole:   t.WorkRole,
		StatusNote:     fmt.Sprintf("Published by %s. The immutable event files are authoritative for history.", id.Verbose()),
	})
	if err != nil {
		return err
	}
	if cf.json {
		return c.jsonOut(map[string]any{"task_id": t.TaskID, "dir": dir, "state": *state, "digest": model.TaskDigest(specBytes)})
	}
	fmt.Fprintf(out, "published %s state=%s dir=%s\n", t.TaskID, *state, dir)
	return nil
}

func briefMustAbs(p string) string {
	if abs, err := filepath.Abs(p); err == nil {
		return abs
	}
	return p
}

func ensureNL(b []byte) []byte {
	if len(b) > 0 && b[len(b)-1] == '\n' {
		return b
	}
	return append(b, '\n')
}

func headOrEmpty(c *cmdCtx) string {
	h, _ := c.Repo.Head()
	return h
}

func nextCapabilityForTask(t *model.Task) string {
	if t.MinimumCapability == model.CapHIGH {
		return model.CapHIGH
	}
	return model.CapLOW
}

// ---------------------------------------------------------------------------
// audit
// ---------------------------------------------------------------------------

type auditIssue struct {
	TaskID   string `json:"task_id,omitempty"`
	Severity string `json:"severity"`
	Kind     string `json:"kind"`
	Detail   string `json:"detail"`
}

type auditReport struct {
	GeneratedAt string       `json:"generated_at"`
	Tasks       int          `json:"tasks"`
	Errors      int          `json:"errors"`
	Warnings    int          `json:"warnings"`
	Issues      []auditIssue `json:"issues"`
}

func cmdAudit(args []string, out io.Writer) error {
	fs := newFlagSet("audit")
	var cf commonFlags
	cf.register(fs, false)
	if err := fs.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	c, err := cf.open(out)
	if err != nil {
		return err
	}
	report := auditReport{GeneratedAt: model.FormatTime(c.Now), Issues: []auditIssue{}}
	records, loadErrs := c.allRecords()
	report.Tasks = len(records)
	for _, e := range loadErrs {
		report.Issues = append(report.Issues, auditIssue{Severity: "error", Kind: "unreadable_record", Detail: e.Error()})
	}
	byID := map[string]*model.TaskRecord{}
	for _, r := range records {
		byID[r.Task.TaskID] = r
	}
	for _, r := range records {
		for _, e := range archive.ValidateChain(r.Task, r.Events, c.Now) {
			report.Issues = append(report.Issues, auditIssue{TaskID: r.Task.TaskID, Severity: "error", Kind: "event_chain", Detail: e.Error()})
		}
	}

	// Live refs: every one must name an existing task, and its spec digest must
	// match the committed task.json.
	out2, err := c.Repo.Run("for-each-ref", "--format=%(refname)", "refs/ads-queue/live/")
	if err != nil {
		return err
	}
	for _, ref := range strings.Split(strings.TrimSpace(out2), "\n") {
		if strings.TrimSpace(ref) == "" {
			continue
		}
		taskID := strings.TrimPrefix(ref, LiveRefPrefix)
		rec, ok := byID[taskID]
		if !ok {
			report.Issues = append(report.Issues, auditIssue{TaskID: taskID, Severity: "error", Kind: "orphan_live_ref", Detail: ref + " names no committed task"})
			continue
		}
		cl, _, err := c.readLiveClaim(taskID)
		if err != nil {
			report.Issues = append(report.Issues, auditIssue{TaskID: taskID, Severity: "error", Kind: "malformed_live_ref", Detail: err.Error()})
			continue
		}
		_, raw, err := archive.LoadTask(rec.Dir)
		if err == nil {
			if d := model.TaskDigest(raw); d != cl.TaskSpecDigest {
				report.Issues = append(report.Issues, auditIssue{TaskID: taskID, Severity: "error", Kind: "task_digest_changed", Detail: fmt.Sprintf("live %s != committed %s", cl.TaskSpecDigest, d)})
			}
		}
		if model.Terminal(cl.State) {
			report.Issues = append(report.Issues, auditIssue{TaskID: taskID, Severity: "warning", Kind: "stale_live_ref", Detail: "live claim exists in terminal state " + cl.State})
		}
		if exp, err := model.ParseTime(cl.ExpiresAt); err == nil && exp.Before(c.Now) && !model.Terminal(cl.State) {
			report.Issues = append(report.Issues, auditIssue{TaskID: taskID, Severity: "warning", Kind: "expired_lease", Detail: "lease expired " + cl.ExpiresAt})
		}
		if last := rec.LastEvent(); last != nil && cl.LastEvent != last.EventID {
			report.Issues = append(report.Issues, auditIssue{TaskID: taskID, Severity: "warning", Kind: "live_ref_drift", Detail: fmt.Sprintf("live last_event %s != committed %s", cl.LastEvent, last.EventID)})
		}
	}

	// Integration lock names must resolve too.
	if _, ok := c.Repo.RefOID(IntegrationLockRef); ok {
		oid, _ := c.Repo.RefOID(IntegrationLockRef)
		data, err := c.Repo.ReadBlob(oid)
		if err != nil {
			report.Issues = append(report.Issues, auditIssue{Severity: "error", Kind: "malformed_integration_lock", Detail: err.Error()})
		} else {
			var lock model.IntegrationLock
			if err := json.Unmarshal(data, &lock); err != nil {
				report.Issues = append(report.Issues, auditIssue{Severity: "error", Kind: "malformed_integration_lock", Detail: err.Error()})
			} else if _, ok := byID[lock.TaskID]; !ok {
				report.Issues = append(report.Issues, auditIssue{TaskID: lock.TaskID, Severity: "error", Kind: "orphan_integration_lock", Detail: "integration lock names no committed task"})
			}
		}
	}

	for _, i := range report.Issues {
		if i.Severity == "error" {
			report.Errors++
		} else {
			report.Warnings++
		}
	}
	if cf.json {
		if err := c.jsonOut(report); err != nil {
			return err
		}
	} else {
		fmt.Fprintf(out, "audit: %d tasks, %d errors, %d warnings\n", report.Tasks, report.Errors, report.Warnings)
		for _, i := range report.Issues {
			fmt.Fprintf(out, "  [%s] %s: %s %s\n", strings.ToUpper(i.Severity), i.Kind, i.TaskID, i.Detail)
		}
	}
	if report.Errors > 0 {
		return failf(1, "audit found %d error(s)", report.Errors)
	}
	return nil
}
