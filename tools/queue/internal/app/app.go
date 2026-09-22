// Package app implements the queue commands on top of gitx, model and archive.
//
// The commands are thin: they read the committed archive, apply one rule, write
// one immutable event and one live-ref compare-and-swap, and report what
// happened. All durable meaning is in the event files; the live ref is only
// coordination.
package app

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"adsqueue/internal/archive"
	"adsqueue/internal/gitx"
	"adsqueue/internal/model"
)

// LiveRefPrefix is where one live claim blob per task is kept.
const LiveRefPrefix = "refs/ads-queue/live/"

// IntegrationLockRef serialises the final update of local `main`.
const IntegrationLockRef = "refs/ads-queue/locks/main-integration"

// exitError carries a process exit code through the error path.
type exitError struct {
	code int
	err  error
}

func (e *exitError) Error() string { return e.err.Error() }
func (e *exitError) Unwrap() error { return e.err }

func failf(code int, format string, args ...any) error {
	return &exitError{code: code, err: fmt.Errorf(format, args...)}
}

// Run dispatches a command and returns the process exit code.
func Run(args []string, stdout, stderr io.Writer) int {
	if len(args) == 0 {
		usage(stderr)
		return 2
	}
	cmd, rest := args[0], args[1:]
	var err error
	switch cmd {
	case "session-start":
		err = cmdSessionStart(rest, stdout)
	case "list":
		err = cmdList(rest, stdout)
	case "next":
		err = cmdNext(rest, stdout)
	case "publish":
		err = cmdPublish(rest, stdout)
	case "claim":
		err = cmdClaim(rest, stdout)
	case "guard":
		err = cmdGuard(rest, stdout)
	case "heartbeat":
		err = cmdHeartbeat(rest, stdout)
	case "checkpoint":
		err = cmdCheckpoint(rest, stdout)
	case "yield":
		err = cmdYield(rest, stdout)
	case "escalate":
		err = cmdEscalate(rest, stdout)
	case "submit":
		err = cmdSubmit(rest, stdout)
	case "review":
		err = cmdReview(rest, stdout)
	case "approve":
		err = cmdApprove(rest, stdout)
	case "request-changes":
		err = cmdRequestChanges(rest, stdout)
	case "integrate":
		err = cmdIntegrate(rest, stdout)
	case "complete":
		err = cmdComplete(rest, stdout)
	case "recover":
		err = cmdRecover(rest, stdout)
	case "audit":
		err = cmdAudit(rest, stdout)
	case "review-candidates":
		err = cmdReviewCandidates(rest, stdout)
	case "help", "-h", "--help":
		usage(stdout)
		return 0
	default:
		fmt.Fprintf(stderr, "queue: unknown command %q\n\n", cmd)
		usage(stderr)
		return 2
	}
	if err != nil {
		var ee *exitError
		if errors.As(err, &ee) {
			fmt.Fprintf(stderr, "queue: %v\n", ee.err)
			return ee.code
		}
		fmt.Fprintf(stderr, "queue: %v\n", err)
		return 1
	}
	return 0
}

func usage(w io.Writer) {
	fmt.Fprint(w, `queue — durable AI work queue

usage: queue <command> [flags]

commands:
  session-start     record this session's capability and role
  list              list tasks with their derived state
  next              print the next eligible task
  publish           publish a task record (HIGH)
  claim             claim a task (atomic; creates its worktree)
  guard             verify this worker's claim is still current
  heartbeat         extend the live lease
  checkpoint        record a checkpoint and move claimed -> in_progress
  yield             release a claim back to ready
  escalate          record ESCALATION REQUIRED and move to blocked_high
  submit            submit for review (awaiting_review)
  review            record a review (HIGH)
  approve           approve for integration (HIGH)
  request-changes   return a review with corrections (HIGH)
  integrate         merge the approved branch into local main and tag it
  complete          verify and mark completed
  recover           recover an expired claim (two CAS steps, grace period)
  audit             validate every event chain and live ref
  review-candidates list awaiting_review tasks submitted within --since

Common flags: --repo, --now, --json, --capability, --role,
              --model, --tool, --effort, --session-id
`)
}

// ---------------------------------------------------------------------------
// context
// ---------------------------------------------------------------------------

type cmdCtx struct {
	Repo *gitx.Repo
	Now  time.Time
	Out  io.Writer
}

func (c *cmdCtx) jsonOut(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	_, err = fmt.Fprintln(c.Out, string(data))
	return err
}

type commonFlags struct {
	repo string
	now  string
	json bool

	capability, role, modelID, tool, effort, sessionID string
}

func (cf *commonFlags) register(fs *flag.FlagSet, withIdentity bool) {
	fs.StringVar(&cf.repo, "repo", "", "repository path (default: current directory)")
	fs.StringVar(&cf.now, "now", "", "freeze the clock at an RFC3339 UTC time (testing)")
	fs.BoolVar(&cf.json, "json", false, "emit JSON")
	if withIdentity {
		fs.StringVar(&cf.capability, "capability", "", "session capability HIGH or LOW (default: $ADS_QUEUE_CAPABILITY)")
		fs.StringVar(&cf.role, "role", "", "work role (default: $ADS_QUEUE_ROLE)")
		fs.StringVar(&cf.modelID, "model", "", "model id (default: $ADS_QUEUE_MODEL)")
		fs.StringVar(&cf.tool, "tool", "", "tool id (default: $ADS_QUEUE_TOOL)")
		fs.StringVar(&cf.effort, "effort", "", "effort setting (default: $ADS_QUEUE_EFFORT)")
		fs.StringVar(&cf.sessionID, "session-id", "", "session id (default: $ADS_QUEUE_SESSION_ID)")
	}
}

func envOr(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if v != "" {
			return v
		}
	}
	return ""
}

// identity resolves flags over environment, using "unknown" for anything the
// session does not expose.
func (cf *commonFlags) identity() model.Identity {
	return cf.actorFrom(model.Identity{})
}

// actorFrom resolves identity like identity(), but falls back to a live claim's
// recorded worker when the caller did not restate capability or role. A worker
// that already holds a claim should not have to repeat its identity on every
// command, and the claim is the authoritative record of who holds the task.
func (cf *commonFlags) actorFrom(def model.Identity) model.Identity {
	return model.Identity{
		Model:             firstNonEmpty(cf.modelID, os.Getenv("ADS_QUEUE_MODEL"), def.Model, "unknown"),
		Tool:              firstNonEmpty(cf.tool, os.Getenv("ADS_QUEUE_TOOL"), def.Tool, "unknown"),
		Effort:            firstNonEmpty(cf.effort, os.Getenv("ADS_QUEUE_EFFORT"), def.Effort, "unknown"),
		SessionCapability: firstNonEmpty(cf.capability, os.Getenv("ADS_QUEUE_CAPABILITY"), def.SessionCapability),
		WorkRole:          firstNonEmpty(cf.role, os.Getenv("ADS_QUEUE_ROLE"), def.WorkRole),
		SessionID:         firstNonEmpty(cf.sessionID, os.Getenv("ADS_QUEUE_SESSION_ID"), def.SessionID, "unknown"),
	}
}

func (cf *commonFlags) clock() (time.Time, error) {
	s := firstNonEmpty(cf.now, os.Getenv("ADS_QUEUE_NOW"))
	if s == "" {
		return time.Now().UTC().Truncate(time.Second), nil
	}
	t, err := model.ParseTime(s)
	if err != nil {
		return time.Time{}, err
	}
	return t, nil
}

// open resolves the repository and the clock for a command.
func (cf *commonFlags) open(out io.Writer) (*cmdCtx, error) {
	dir := firstNonEmpty(cf.repo, os.Getenv("ADS_QUEUE_REPO"))
	repo, err := gitx.Open(dir)
	if err != nil {
		return nil, err
	}
	now, err := cf.clock()
	if err != nil {
		return nil, err
	}
	return &cmdCtx{Repo: repo, Now: now, Out: out}, nil
}

// ---------------------------------------------------------------------------
// refs and claims
// ---------------------------------------------------------------------------

func liveRef(taskID string) string { return LiveRefPrefix + taskID }

// randHex returns n random bytes as hex, for claim and attempt ids.
func randHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		// A failure here would make two claims identical, which is worse than
		// stopping; fall back to time-based bytes and let CAS sort out races.
		copy(b, []byte(time.Now().UTC().Format(time.RFC3339Nano)))
	}
	return hex.EncodeToString(b)
}

// readLiveClaim reads and parses the live claim blob for a task.
func (c *cmdCtx) readLiveClaim(taskID string) (claim *model.Claim, oid string, err error) {
	oid, ok := c.Repo.RefOID(liveRef(taskID))
	if !ok {
		return nil, "", nil
	}
	data, err := c.Repo.ReadBlob(oid)
	if err != nil {
		return nil, oid, fmt.Errorf("read live claim %s: %w", taskID, err)
	}
	var cl model.Claim
	if err := json.Unmarshal(data, &cl); err != nil {
		return nil, oid, fmt.Errorf("parse live claim %s: %w", taskID, err)
	}
	return &cl, oid, nil
}

// writeLiveClaim writes the blob and compares-and-swaps the ref.
func (c *cmdCtx) writeLiveClaim(cl *model.Claim, oldOID, message string) error {
	data, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	oid, err := c.Repo.WriteBlob(data)
	if err != nil {
		return err
	}
	if err := c.Repo.CAS(liveRef(cl.TaskID), oid, oldOID, message); err != nil {
		return err
	}
	return nil
}

func (c *cmdCtx) releaseLiveClaim(taskID, oldOID, message string) error {
	if oldOID == "" {
		return nil
	}
	if err := c.Repo.DeleteCAS(liveRef(taskID), oldOID, message); err != nil {
		return err
	}
	return nil
}

// claimContextPath records this worktree's current claim id and epoch inside
// the git dir, so it is never committed. It is per worktree: a claim made from
// the main checkout and then followed by work in the task worktree writes a
// second copy into that worktree's git dir (see saveClaimContextFor).
func (c *cmdCtx) claimContextPath(taskID string) string {
	return filepath.Join(c.Repo.GitDir, "ads-queue", "claims", taskID+".json")
}

// saveClaimContextFor writes the claim context into the git dir of the named
// worktree, falling back to the current repo when the worktree is not a repo
// yet. A worker in the task worktree then finds its own claim.
func (c *cmdCtx) saveClaimContextFor(worktreeRel string, cl *model.Claim) error {
	repo := c.Repo
	if worktreeRel != "" {
		if wt, err := gitx.Open(c.Repo.ResolveCommonPath(worktreeRel)); err == nil {
			repo = wt
		}
	}
	return (&cmdCtx{Repo: repo}).saveClaimContext(cl)
}

func (c *cmdCtx) saveClaimContext(cl *model.Claim) error {
	path := c.claimContextPath(cl.TaskID)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cl, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func (c *cmdCtx) loadClaimContext(taskID string) (*model.Claim, error) {
	data, err := os.ReadFile(c.claimContextPath(taskID))
	if err != nil {
		return nil, err
	}
	var cl model.Claim
	if err := json.Unmarshal(data, &cl); err != nil {
		return nil, err
	}
	return &cl, nil
}

// resolveClaim finds this worker's claim id and epoch from flags first, then the
// saved context file.
func (c *cmdCtx) resolveClaim(taskID, claimID string, epoch int) (string, int, error) {
	if claimID != "" {
		return claimID, epoch, nil
	}
	if envID := os.Getenv("ADS_QUEUE_CLAIM_ID"); envID != "" {
		e := epoch
		if e == 0 {
			e = 1
		}
		return envID, e, nil
	}
	cl, err := c.loadClaimContext(taskID)
	if err != nil {
		return "", 0, failf(3, "no claim context for %s; pass --claim-id (run `queue claim` first)", taskID)
	}
	return cl.ClaimID, cl.ClaimEpoch, nil
}

// ---------------------------------------------------------------------------
// task loading
// ---------------------------------------------------------------------------

func (c *cmdCtx) taskRecord(taskID string) (*model.TaskRecord, error) {
	if !model.ValidID(taskID, "task-") {
		return nil, fmt.Errorf("invalid task id %q", taskID)
	}
	dir, err := archive.TaskDirAbs(c.Repo.Root, taskID)
	if err != nil {
		return nil, err
	}
	rec, err := archive.LoadTaskRecord(dir)
	if err != nil {
		return nil, err
	}
	if errs := archive.ValidateChain(rec.Task, rec.Events, c.Now); len(errs) > 0 {
		return nil, fmt.Errorf("task %s has an invalid event chain: %v", taskID, errs[0])
	}
	return rec, nil
}

func (c *cmdCtx) allRecords() ([]*model.TaskRecord, []error) {
	return archive.LoadAll(c.Repo.Root)
}

// ---------------------------------------------------------------------------
// event emission
// ---------------------------------------------------------------------------

type eventParams struct {
	EventType string
	State     string
	Actor     model.Identity
	AttemptID string
	ClaimID   string
	Summary   string

	BaseCommit       string
	CheckpointCommit string
	ResultCommit     string

	EvidencePaths []string
	Related       model.Related

	NextCapability string
	NextWorkRole   string
	StatusNote     string
}

// emitEvent writes one immutable event plus the STATUS snapshot, commits
// exactly those two files on repo's current branch, and returns the event. The
// commit is the durable record; callers update live refs separately.
func (c *cmdCtx) emitEvent(repo *gitx.Repo, rec *model.TaskRecord, p eventParams) (*model.Event, error) {
	if !model.AllowedTransition(rec.State, p.State) {
		return nil, fmt.Errorf("impossible transition %s -> %s for %s", rec.State, p.State, rec.Task.TaskID)
	}
	if err := p.Actor.Validate(); err != nil {
		return nil, fmt.Errorf("actor: %w", err)
	}
	taken, err := archive.EventIDs(rec.Dir)
	if err != nil {
		return nil, err
	}
	ev := &model.Event{
		SchemaVersion:    model.SchemaVersion,
		EventID:          archive.EventID(p.EventType, c.Now, func(id string) bool { return taken[id] }),
		Sequence:         len(rec.Events) + 1,
		TaskID:           rec.Task.TaskID,
		GoalID:           rec.Task.GoalID,
		RootTaskID:       rec.Task.RootTaskID,
		ParentTaskID:     rec.Task.ParentTaskID,
		OccurredAt:       model.FormatTime(c.Now),
		Actor:            p.Actor,
		EventType:        p.EventType,
		ResultingState:   p.State,
		Summary:          p.Summary,
		EvidencePaths:    p.EvidencePaths,
		Related:          p.Related,
		NextCapability:   p.NextCapability,
		NextWorkRole:     p.NextWorkRole,
		BaseCommit:       optString(p.BaseCommit),
		CheckpointCommit: optString(p.CheckpointCommit),
		ResultCommit:     optString(p.ResultCommit),
	}
	if last := rec.LastEvent(); last != nil {
		id := last.EventID
		ev.PreviousEventID = &id
	}
	if p.AttemptID != "" {
		id := p.AttemptID
		ev.AttemptID = &id
	}
	if p.ClaimID != "" {
		id := p.ClaimID
		ev.ClaimID = &id
	}

	path, err := archive.WriteEvent(rec.Dir, ev)
	if err != nil {
		return nil, err
	}
	// Render the snapshot from the state *after* this event. Rendering from rec
	// would leave STATUS.md one event behind (it would still name the previous
	// event and state), which is exactly the kind of quiet drift the archive is
	// supposed to make impossible.
	view := *rec
	view.Events = append(append([]*model.Event{}, rec.Events...), ev)
	view.State = ev.ResultingState
	statusPath, err := archive.WriteStatus(rec.Dir, archive.RenderStatus(&view, c.Now, p.NextCapability, p.NextWorkRole, p.StatusNote))
	if err != nil {
		return nil, err
	}
	relEvent, err := repo.RepoRelative(path)
	if err != nil {
		return nil, err
	}
	relStatus, err := repo.RepoRelative(statusPath)
	if err != nil {
		return nil, err
	}
	msg := fmt.Sprintf("Queue: %s %s -> %s (%s)\n\n%s\n\nAgent: %s\n",
		rec.Task.TaskID, rec.State, p.State, p.EventType, p.Summary, p.Actor.Verbose())
	files := []string{relEvent, relStatus}
	for _, ep := range p.EvidencePaths {
		if ep == "" {
			continue
		}
		abs := repo.ResolveRepoPath(ep)
		fi, statErr := os.Stat(abs)
		if statErr != nil || fi.IsDir() {
			continue
		}
		rel, relErr := repo.RepoRelative(abs)
		if relErr != nil || strings.HasPrefix(rel, "..") {
			continue
		}
		files = append(files, rel)
	}
	if _, err := repo.CommitPaths(msg, files...); err != nil {
		return nil, err
	}
	return ev, nil
}

func optString(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func summarize(ev *model.Event) string {
	if ev == nil {
		return "(no events)"
	}
	prev := "none"
	if ev.PreviousEventID != nil {
		prev = *ev.PreviousEventID
	}
	return fmt.Sprintf("%s -> %s (%s) at %s", prev, ev.ResultingState, ev.EventType, ev.OccurredAt)
}

// ---------------------------------------------------------------------------
// benchmark lock
// ---------------------------------------------------------------------------

// benchmarkLockHeldBy reports whether the benchmark lock is held by the task's
// worktree. Recovery must not run while the crashed worker still owns
// ads-run-lock: another process may be measuring on its behalf.
//
// The lock is a Podman named volume (infra/lib.sh). Podman may be unavailable
// (for example during a pure unit test); in that case only the environment
// variable is consulted, which is what the runner exports.
func benchmarkLockHeldBy(worktree string) bool {
	if os.Getenv("ADS_RUN_LOCK_HELD") != "" {
		return true
	}
	podmanBin := firstNonEmpty(os.Getenv("ADS_PODMAN"), os.Getenv("PODMAN"))
	if podmanBin == "" {
		podmanBin = "podman"
	}
	out, err := exec.Command(podmanBin, "volume", "inspect", "ads-run-lock",
		"--format", `{{index .Labels "ads.worktree"}}`).Output()
	if err != nil {
		return false
	}
	holder := strings.TrimSpace(string(out))
	if holder == "" {
		return false
	}
	return samePath(holder, worktree)
}

func samePath(a, b string) bool {
	na := strings.ToLower(strings.ReplaceAll(filepath.ToSlash(a), "\\", "/"))
	nb := strings.ToLower(strings.ReplaceAll(filepath.ToSlash(b), "\\", "/"))
	return strings.TrimSuffix(na, "/") == strings.TrimSuffix(nb, "/")
}
