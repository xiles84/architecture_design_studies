package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"adsqueue/internal/archive"
	"adsqueue/internal/gitx"
	"adsqueue/internal/model"
)

const (
	t0 = "2026-09-22T10:00:00Z"
	t1 = "2026-09-22T10:01:00Z"
	t2 = "2026-09-22T10:02:00Z"
	t3 = "2026-09-22T10:03:00Z"
	t4 = "2026-09-22T10:04:00Z"
	t5 = "2026-09-22T10:05:00Z"
	t6 = "2026-09-22T10:06:00Z"
)

// TestHelperProcess is the standard re-exec hook: the test binary runs itself
// as the real queue CLI in a child process, which is how the parallel-claimer
// test exercises true multi-process atomicity.
func TestHelperProcess(t *testing.T) {
	if os.Getenv("ADS_QUEUE_HELPER") != "1" {
		return
	}
	idx := -1
	for i, a := range os.Args {
		if a == "--" {
			idx = i
			break
		}
	}
	if idx < 0 {
		os.Exit(2)
	}
	os.Exit(Run(os.Args[idx+1:], os.Stdout, os.Stderr))
}

type qres struct {
	code int
	out  string
	err  string
}

func qrun(args ...string) qres {
	var o, e bytes.Buffer
	code := Run(args, &o, &e)
	return qres{code, o.String(), e.String()}
}

func mustQ(t *testing.T, args ...string) qres {
	t.Helper()
	r := qrun(args...)
	if r.code != 0 {
		t.Fatalf("queue %s failed (%d): %s%s", strings.Join(args, " "), r.code, r.out, r.err)
	}
	return r
}

func git(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s in %s: %v\n%s", strings.Join(args, " "), dir, err, out)
	}
	return strings.TrimSpace(string(out))
}

func newRepo(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	// The cross-platform check needs a repository that a native Windows git can
	// also address, so it asks for a root under /mnt/<drive>.
	if base := os.Getenv("ADS_QUEUE_TEST_ROOT"); base != "" {
		if err := os.MkdirAll(base, 0o755); err != nil {
			t.Fatal(err)
		}
		d, err := os.MkdirTemp(base, "queuetest-")
		if err != nil {
			t.Fatal(err)
		}
		dir = d
	}
	git(t, dir, "init", "-q", "-b", "main")
	git(t, dir, "config", "user.email", "test@example.com")
	git(t, dir, "config", "user.name", "Queue Test")
	git(t, dir, "config", "commit.gpgsign", "false")
	git(t, dir, "config", "core.autocrlf", "false")
	return dir
}

func writeFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

func commitFile(t *testing.T, repo, rel, content, msg string) {
	t.Helper()
	writeFile(t, filepath.Join(repo, filepath.FromSlash(rel)), content)
	git(t, repo, "add", "--", rel)
	git(t, repo, "commit", "-q", "-m", msg, "--", rel)
}

func specJSON(id, minCap, branch, worktree, tag string) string {
	return fmt.Sprintf(`{
  "schema_version": 1,
  "task_id": %q,
  "goal_id": "goal-20260922T100000Z-demo",
  "root_task_id": %q,
  "title": "demo task",
  "kind": "implementation",
  "created_at": %q,
  "created_by": {"model":"test","tool":"gotest","effort":"unknown","session_capability":"HIGH","work_role":"planner","session_id":"s0"},
  "priority": 50,
  "not_before": %q,
  "minimum_capability": %q,
  "preferred_capability": %q,
  "work_role": "executor",
  "dependencies": [],
  "source_request": "docs/ai-work/goals/demo/REQUEST.md",
  "source_handoff": "docs/ai-work/tasks/demo/BRIEF.md",
  "canonical_branch": %q,
  "canonical_worktree": %q,
  "owned_paths": ["docs/ai-work"],
  "forbidden_paths": [],
  "acceptance_criteria": ["the demo works"],
  "expected_artifacts": ["a thing"],
  "validation": ["go test"],
  "review_policy": "HIGH required",
  "integration_required": true,
  "required_tag": %q,
  "benchmark_required": false,
  "revalidate_after": null
}`, id, id, t0, t0, minCap, minCap, branch, worktree, tag)
}

func publishDemo(t *testing.T, repo, id, minCap string) {
	t.Helper()
	tmp := t.TempDir()
	spec := filepath.Join(tmp, "task.json")
	brief := filepath.Join(tmp, "BRIEF.md")
	writeFile(t, spec, specJSON(id, minCap, "repo/demo", ".worktrees/demo", "repo/demo"))
	writeFile(t, brief, "# brief\n\nDo the demo.\n")
	mustQ(t, "publish", "--spec", spec, "--brief", brief,
		"--capability", "HIGH", "--role", "planner", "--model", "test", "--tool", "gotest", "--session-id", "s0",
		"--repo", repo, "--now", t0)
}

func taskDirIn(root, id string) string {
	d, _ := archive.TaskDirAbs(root, id)
	return d
}

func stateIn(t *testing.T, root, id string) string {
	t.Helper()
	rec, err := archive.LoadTaskRecord(taskDirIn(root, id))
	if err != nil {
		t.Fatalf("load record: %v", err)
	}
	return rec.State
}

// assertStatusMatchesLastEvent guards the snapshot against lagging the archive:
// STATUS.md must name the newest event, its resulting state and its sequence.
func assertStatusMatchesLastEvent(t *testing.T, root, id string) {
	t.Helper()
	rec, err := archive.LoadTaskRecord(taskDirIn(root, id))
	if err != nil {
		t.Fatalf("load record: %v", err)
	}
	last := rec.LastEvent()
	if last == nil {
		t.Fatal("no events")
	}
	data, err := os.ReadFile(filepath.Join(taskDirIn(root, id), "STATUS.md"))
	if err != nil {
		t.Fatal(err)
	}
	status := string(data)
	for _, want := range []string{
		"| State | `" + last.ResultingState + "` |",
		fmt.Sprintf("| Sequence | `%d` |", last.Sequence),
		"| Last event | `" + last.EventID + "` |",
	} {
		if !strings.Contains(status, want) {
			t.Fatalf("STATUS.md lags the last event: missing %q in\n%s", want, status)
		}
	}
}

func claimDemo(t *testing.T, repo, id string, now string) model.Claim {
	t.Helper()
	r := mustQ(t, "claim", "--task", id, "--capability", "LOW", "--role", "executor",
		"--model", "test", "--tool", "gotest", "--session-id", "s1",
		"--repo", repo, "--now", now, "--json")
	var cl model.Claim
	if err := json.Unmarshal([]byte(r.out), &cl); err != nil {
		t.Fatalf("parse claim json: %v\n%s", err, r.out)
	}
	return cl
}

// ---------------------------------------------------------------------------

func TestFullLifecycle(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-demo"
	publishDemo(t, repo, id, "LOW")
	if got := stateIn(t, repo, id); got != model.StateReady {
		t.Fatalf("after publish: %s", got)
	}

	cl := claimDemo(t, repo, id, t0)
	wt := filepath.Join(repo, ".worktrees", "demo")
	if got := stateIn(t, repo, id); got != model.StateClaimed {
		t.Fatalf("after claim: %s", got)
	}
	if _, err := os.Stat(wt); err != nil {
		t.Fatalf("worktree not created: %v", err)
	}

	// Worker writes the required result and commits it.
	resRel := filepath.ToSlash(filepath.Join("docs/ai-work/tasks/2026/09", id, "attempts", cl.AttemptID, "RESULT.md"))
	commitFile(t, wt, resRel, "# result\n\nDone.\n", "worker: result")

	mustQ(t, "checkpoint", "--task", id, "--repo", wt, "--now", t1)
	if got := stateIn(t, wt, id); got != model.StateInProgress {
		t.Fatalf("after checkpoint: %s", got)
	}
	mustQ(t, "submit", "--task", id, "--repo", wt, "--now", t2)
	if got := stateIn(t, wt, id); got != model.StateAwaitingReview {
		t.Fatalf("after submit: %s", got)
	}
	assertStatusMatchesLastEvent(t, wt, id)

	tmp := t.TempDir()
	review := filepath.Join(tmp, "REVIEW.md")
	writeFile(t, review, "# review\n\nApproved.\n")
	mustQ(t, "review", "--task", id, "--review-file", review,
		"--capability", "HIGH", "--role", "reviewer", "--model", "test", "--tool", "gotest", "--session-id", "s2",
		"--repo", wt, "--now", t3)
	mustQ(t, "approve", "--task", id,
		"--capability", "HIGH", "--role", "reviewer", "--model", "test", "--tool", "gotest", "--session-id", "s2",
		"--repo", wt, "--now", t4)
	mustQ(t, "integrate", "--task", id,
		"--capability", "HIGH", "--role", "integrator", "--model", "test", "--tool", "gotest", "--session-id", "s2",
		"--repo", wt, "--now", t5)
	mustQ(t, "complete", "--task", id,
		"--capability", "HIGH", "--role", "integrator", "--model", "test", "--tool", "gotest", "--session-id", "s2",
		"--repo", wt, "--now", t6)

	if got := stateIn(t, repo, id); got != model.StateCompleted {
		t.Fatalf("final state on main: %s", got)
	}
	assertStatusMatchesLastEvent(t, repo, id)
	if out := git(t, repo, "rev-parse", "--verify", "refs/tags/repo/demo"); out == "" {
		t.Fatal("required tag missing")
	}
	// The audit must be clean on the final state.
	mustQ(t, "audit", "--repo", repo, "--now", t6)
}

func TestCapabilityRules(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-highonly"
	publishDemo(t, repo, id, "HIGH")
	r := qrun("claim", "--task", id, "--capability", "LOW", "--role", "executor",
		"--model", "test", "--tool", "gotest", "--session-id", "s1", "--repo", repo, "--now", t0)
	if r.code == 0 {
		t.Fatal("LOW must not claim a HIGH-only task")
	}
	// A HIGH session may claim either.
	mustQ(t, "claim", "--task", id, "--capability", "HIGH", "--role", "executor",
		"--model", "test", "--tool", "gotest", "--session-id", "s1", "--repo", repo, "--now", t0)
}

func TestParallelClaimersExactlyOneWinner(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-race"
	publishDemo(t, repo, id, "LOW")

	const n = 20
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		winners int
	)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			cmd := exec.Command(os.Args[0], "-test.run=TestHelperProcess", "--",
				"claim", "--task", id, "--capability", "LOW", "--role", "executor",
				"--model", "test", "--tool", "gotest", "--session-id", fmt.Sprintf("s%d", i),
				"--repo", repo, "--now", t0, "--no-worktree")
			cmd.Env = append(os.Environ(), "ADS_QUEUE_HELPER=1")
			err := cmd.Run()
			if err == nil {
				mu.Lock()
				winners++
				mu.Unlock()
			}
		}(i)
	}
	wg.Wait()
	if winners != 1 {
		t.Fatalf("expected exactly one winner, got %d", winners)
	}
	if got := stateIn(t, repo, id); got != model.StateClaimed {
		t.Fatalf("state after race: %s", got)
	}
}

func TestParallelIntegrationLockExactlyOneHolder(t *testing.T) {
	repo := newRepo(t)
	// An initial commit so the repo has a HEAD and refs work.
	commitFile(t, repo, "seed.txt", "seed\n", "seed")
	c, err := gitx.Open(repo)
	if err != nil {
		t.Fatal(err)
	}
	ctx := &cmdCtx{Repo: c, Now: mustTime(t, t0), Out: &bytes.Buffer{}}
	worker := model.Identity{Model: "test", Tool: "gotest", Effort: "unknown", SessionCapability: model.CapHIGH, WorkRole: model.RoleIntegrator, SessionID: "s"}

	const n = 20
	var (
		wg      sync.WaitGroup
		mu      sync.Mutex
		winners int
	)
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if _, err := ctx.acquireIntegrationLock("task-20260922T100000Z-x", "claim-x", worker); err == nil {
				mu.Lock()
				winners++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	if winners != 1 {
		t.Fatalf("expected exactly one integration-lock holder, got %d", winners)
	}
}

func mustTime(t *testing.T, s string) time.Time {
	t.Helper()
	tt, err := model.ParseTime(s)
	if err != nil {
		t.Fatal(err)
	}
	return tt
}

func TestGuardRejectsSupersededClaim(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-recover"
	publishDemo(t, repo, id, "LOW")
	cl := claimDemo(t, repo, id, t0)
	wt := filepath.Join(repo, ".worktrees", "demo")

	// Lease expires; a second worker recovers in two CAS steps.
	expired := "2026-09-22T13:00:00Z" // 3h after t0 (2h lease)
	mustQ(t, "recover", "--task", id, "--capability", "LOW", "--role", "executor",
		"--model", "test2", "--tool", "gotest", "--session-id", "s2", "--repo", repo, "--now", expired)
	afterGrace := "2026-09-22T13:10:01Z"
	mustQ(t, "recover", "--task", id, "--capability", "LOW", "--role", "executor",
		"--model", "test2", "--tool", "gotest", "--session-id", "s2", "--repo", repo, "--now", afterGrace)

	// The original worker's next heartbeat/guard must be rejected: it presents
	// its own (now stale) epoch to the guard.
	r := qrun("guard", "--task", id, "--claim-id", cl.ClaimID, "--claim-epoch", fmt.Sprint(cl.ClaimEpoch),
		"--repo", repo, "--now", afterGrace)
	if r.code != 3 {
		t.Fatalf("stale guard must fail with code 3, got %d: %s%s", r.code, r.out, r.err)
	}
	if !strings.Contains(r.err, "superseded") {
		t.Fatalf("stale guard message: %s", r.err)
	}
	// The new claim, held in the task worktree, is current.
	mustQ(t, "guard", "--task", id, "--repo", wt, "--now", afterGrace)
}

func TestRecoveryPreservesDirtyWorktree(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-dirty"
	publishDemo(t, repo, id, "LOW")
	claimDemo(t, repo, id, t0)
	wt := filepath.Join(repo, ".worktrees", "demo")

	// An uncommitted file the crashed worker left behind.
	if err := os.WriteFile(filepath.Join(wt, "unfinished.txt"), []byte("do not lose me\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	expired := "2026-09-22T13:00:00Z"
	afterGrace := "2026-09-22T13:10:01Z"
	mustQ(t, "recover", "--task", id, "--capability", "LOW", "--role", "executor", "--model", "test2", "--tool", "gotest", "--session-id", "s2", "--repo", repo, "--now", expired)
	mustQ(t, "recover", "--task", id, "--capability", "LOW", "--role", "executor", "--model", "test2", "--tool", "gotest", "--session-id", "s2", "--repo", repo, "--now", afterGrace)

	data, err := os.ReadFile(filepath.Join(wt, "unfinished.txt"))
	if err != nil || string(data) != "do not lose me\n" {
		t.Fatalf("recovery discarded an unfinished file: %v %q", err, data)
	}
}

func TestRecoveryRefusedWhileBenchmarkLockHeld(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-locked"
	publishDemo(t, repo, id, "LOW")
	claimDemo(t, repo, id, t0)
	expired := "2026-09-22T13:00:00Z"
	t.Setenv("ADS_RUN_LOCK_HELD", "1")
	r := qrun("recover", "--task", id, "--capability", "LOW", "--role", "executor",
		"--model", "test2", "--tool", "gotest", "--session-id", "s2", "--repo", repo, "--now", expired)
	if r.code == 0 {
		t.Fatal("recovery must be refused while the task owns the benchmark lock")
	}
	if !strings.Contains(r.err, "benchmark lock") {
		t.Fatalf("expected a benchmark-lock refusal, got: %s", r.err)
	}
}

func TestCrashBeforeWorktreeIsRecoverable(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-crash"
	publishDemo(t, repo, id, "LOW")
	// Claim without creating the worktree: this is the state a crash between
	// the CAS and `git worktree add` leaves behind.
	claimDemoNoWorktree(t, repo, id, t0)
	if _, err := os.Stat(filepath.Join(repo, ".worktrees", "demo")); err == nil {
		t.Fatal("worktree should not exist yet")
	}
	expired := "2026-09-22T13:00:00Z"
	afterGrace := "2026-09-22T13:10:01Z"
	mustQ(t, "recover", "--task", id, "--capability", "LOW", "--role", "executor", "--model", "test2", "--tool", "gotest", "--session-id", "s2", "--repo", repo, "--now", expired)
	mustQ(t, "recover", "--task", id, "--capability", "LOW", "--role", "executor", "--model", "test2", "--tool", "gotest", "--session-id", "s2", "--repo", repo, "--now", afterGrace)
	if _, err := os.Stat(filepath.Join(repo, ".worktrees", "demo")); err != nil {
		t.Fatalf("recovery must recreate the missing worktree: %v", err)
	}
}

func claimDemoNoWorktree(t *testing.T, repo, id, now string) {
	t.Helper()
	mustQ(t, "claim", "--task", id, "--capability", "LOW", "--role", "executor",
		"--model", "test", "--tool", "gotest", "--session-id", "s1",
		"--repo", repo, "--now", now, "--no-worktree")
}

func TestFreshCloneReconstruction(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-clone"
	publishDemo(t, repo, id, "LOW")
	claimDemo(t, repo, id, t0)
	// Commit the claimed event is on main; the branch also has the worktree.
	clone := filepath.Join(t.TempDir(), "clone")
	git(t, filepath.Dir(clone), "clone", "-q", repo, clone)
	git(t, clone, "config", "user.email", "test@example.com")
	git(t, clone, "config", "user.name", "Queue Test")
	// Custom live refs are not cloned: the committed events must reconstruct.
	if got := stateIn(t, clone, id); got != model.StateClaimed {
		t.Fatalf("clone state = %s, want claimed", got)
	}
	r := mustQ(t, "list", "--repo", clone, "--json", "--capability", "LOW")
	if !strings.Contains(r.out, id) || !strings.Contains(r.out, "claimed") {
		t.Fatalf("clone list did not reconstruct: %s", r.out)
	}
}

func TestWorktreeRegistrationIsPortable(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-portable"
	publishDemo(t, repo, id, "LOW")
	claimDemo(t, repo, id, t0)
	wt := filepath.Join(repo, ".worktrees", "demo")

	link, err := os.ReadFile(filepath.Join(wt, ".git"))
	if err != nil {
		t.Fatal(err)
	}
	target := strings.TrimSpace(strings.TrimPrefix(string(link), "gitdir: "))
	if filepath.IsAbs(target) || strings.Contains(target, ":") {
		t.Fatalf("worktree link file is not relative: %q", string(link))
	}
	// The per-worktree gitdir back-pointer must be relative too.
	gitdir := filepath.Join(repo, ".git", "worktrees", "demo", "gitdir")
	back, err := os.ReadFile(gitdir)
	if err != nil {
		t.Fatal(err)
	}
	if filepath.IsAbs(strings.TrimSpace(string(back))) {
		t.Fatalf("gitdir back-pointer is absolute: %q", string(back))
	}
	// A second git binary (native Windows Git) must be able to resolve it when
	// one is available on this host.
	if second := secondGit(); second != "" && hostVisible(repo) {
		cmd := exec.Command(second, "-C", toWindowsPath(repo), "worktree", "list", "--porcelain")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("second git could not use the shared common dir: %v\n%s", err, out)
		}
		if !strings.Contains(string(out), "demo") {
			t.Fatalf("second git did not list the worktree: %s", out)
		}
		// And both gits must see the same live coordination ref.
		git(t, repo, "update-ref", "refs/ads-queue/live/task-20260922T100000Z-portable", headOf(t, repo))
		if oid := strings.TrimSpace(string(mustOutput(t, second, "-C", toWindowsPath(repo), "rev-parse", "refs/ads-queue/live/task-20260922T100000Z-portable"))); oid != headOf(t, repo) {
			t.Fatalf("second git read a different live ref: %q", oid)
		}
	}
}

func headOf(t *testing.T, repo string) string {
	t.Helper()
	return strings.TrimSpace(git(t, repo, "rev-parse", "HEAD"))
}

func mustOutput(t *testing.T, bin string, args ...string) []byte {
	t.Helper()
	out, err := exec.Command(bin, args...).CombinedOutput()
	if err != nil {
		t.Fatalf("%s %s: %v\n%s", bin, strings.Join(args, " "), err, out)
	}
	return out
}

// toWindowsPath converts a WSL /mnt/<drive>/... path into C:/... so a native
// Windows binary can address it.
func toWindowsPath(p string) string {
	if len(p) > 6 && strings.HasPrefix(p, "/mnt/") && p[6] == '/' {
		return strings.ToUpper(p[5:6]) + ":" + p[6:]
	}
	return p
}

func secondGit() string {
	candidates := []string{"/mnt/c/Program Files/Git/cmd/git.exe"}
	if v := os.Getenv("ADS_QUEUE_GIT2"); v != "" {
		candidates = append([]string{v}, candidates...)
	}
	for _, c := range candidates {
		if fi, err := os.Stat(c); err == nil && !fi.IsDir() {
			return c
		}
	}
	return ""
}

// hostVisible reports whether a path can be addressed by a native Windows
// binary (i.e. it lives under /mnt/<drive>/).
func hostVisible(p string) bool {
	abs, _ := filepath.Abs(p)
	return strings.HasPrefix(abs, "/mnt/") && len(abs) > 6 && abs[6] == '/'
}

func TestReviewCandidatesSince(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-cand"
	publishDemo(t, repo, id, "LOW")
	cl := claimDemo(t, repo, id, t0)
	wt := filepath.Join(repo, ".worktrees", "demo")
	resRel := filepath.ToSlash(filepath.Join("docs/ai-work/tasks/2026/09", id, "attempts", cl.AttemptID, "RESULT.md"))
	commitFile(t, wt, resRel, "# result\n", "worker: result")
	mustQ(t, "checkpoint", "--task", id, "--repo", wt, "--now", t1)
	mustQ(t, "submit", "--task", id, "--repo", wt, "--now", t2)

	// Submitted at t2 (2026-09-22T10:02Z). At t2+6d it is inside 7d; at t2+8d it
	// is outside.
	inside := "2026-09-28T10:00:00Z"
	outside := "2026-09-30T10:03:00Z"
	r := mustQ(t, "review-candidates", "--repo", wt, "--now", inside, "--since", "7d", "--json")
	if !strings.Contains(r.out, id) {
		t.Fatalf("expected candidate within 7d: %s", r.out)
	}
	r = mustQ(t, "review-candidates", "--repo", wt, "--now", outside, "--since", "7d", "--json")
	if strings.Contains(r.out, id) {
		t.Fatalf("candidate must fall outside 7d: %s", r.out)
	}
}

func TestAuditFindsOrphanLiveRef(t *testing.T) {
	repo := newRepo(t)
	commitFile(t, repo, "seed.txt", "seed\n", "seed")
	// Fabricate a live ref with no task behind it.
	blob := exec.Command("git", "hash-object", "-w", "--stdin")
	blob.Dir = repo
	blob.Stdin = strings.NewReader(`{"schema_version":1,"task_id":"task-20260922T100000Z-ghost","state":"claimed"}`)
	out, err := blob.Output()
	if err != nil {
		t.Fatal(err)
	}
	git(t, repo, "update-ref", "refs/ads-queue/live/task-20260922T100000Z-ghost", strings.TrimSpace(string(out)))
	r := qrun("audit", "--repo", repo, "--now", t0)
	if r.code == 0 {
		t.Fatalf("audit must fail on an orphan ref: %s%s", r.out, r.err)
	}
	if !strings.Contains(r.out+r.err, "orphan_live_ref") {
		t.Fatalf("audit did not name the orphan: %s%s", r.out, r.err)
	}
}

// TestCapabilityAliases implements the requirement added by the queue-v1
// amendment 20260922T101113Z-capability-aliases: every accepted declaration
// normalizes to canonical HIGH/LOW, the raw input is preserved separately, and
// the excluded datastore terms are rejected.
func TestCapabilityAliases(t *testing.T) {
	cases := []struct{ input, want string }{
		{"HIGH", model.CapHIGH},
		{"high", model.CapHIGH},
		{"Leader", model.CapHIGH},
		{"  master  ", model.CapHIGH},
		{"LOW", model.CapLOW},
		{"worker", model.CapLOW},
		{"Follower", model.CapLOW},
		{"SLAVE", model.CapLOW},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(strings.TrimSpace(tc.input), func(t *testing.T) {
			repo := newRepo(t)
			id := "task-20260922T100000Z-alias"
			publishDemo(t, repo, id, "LOW")
			sid := "session-20260922T100000Z-alias"
			mustQ(t, "session-start", "--capability", tc.input, "--role", "executor",
				"--model", "test", "--tool", "gotest", "--session-id", sid, "--repo", repo, "--now", t0)

			data, err := os.ReadFile(filepath.Join(repo, "docs", "ai-work", "sessions", sid+".json"))
			if err != nil {
				t.Fatal(err)
			}
			var sr struct {
				SessionCapability      string `json:"session_capability"`
				SessionCapabilityInput string `json:"session_capability_input"`
				WorkRole               string `json:"work_role"`
			}
			if err := json.Unmarshal(data, &sr); err != nil {
				t.Fatal(err)
			}
			if sr.SessionCapability != tc.want || sr.SessionCapabilityInput != tc.input {
				t.Fatalf("session record capability=%q input=%q; want %q and %q",
					sr.SessionCapability, sr.SessionCapabilityInput, tc.want, tc.input)
			}
			if sr.WorkRole != "executor" {
				t.Fatalf("capability alias changed the work role: %q", sr.WorkRole)
			}

			r := mustQ(t, "claim", "--task", id, "--capability", tc.input, "--role", "executor",
				"--model", "test", "--tool", "gotest", "--session-id", sid, "--repo", repo, "--now", t0, "--json")
			var cl model.Claim
			if err := json.Unmarshal([]byte(r.out), &cl); err != nil {
				t.Fatalf("parse claim json: %v\n%s", err, r.out)
			}
			if cl.Worker.SessionCapability != tc.want {
				t.Fatalf("live claim stored capability %q; want canonical %q", cl.Worker.SessionCapability, tc.want)
			}
			if cl.Worker.SessionCapabilityInput != tc.input {
				t.Fatalf("live claim lost the raw declaration: %q vs %q", cl.Worker.SessionCapabilityInput, tc.input)
			}
			// Eligibility uses the canonical value: a LOW-eligible task is
			// claimable by an alias of LOW, not by an alias of HIGH.
			if got := stateIn(t, repo, id); got != model.StateClaimed {
				t.Fatalf("state after alias claim: %s", got)
			}
		})
	}

	// Excluded datastore terms are rejected as capability declarations.
	for _, bad := range []string{"primary", "replica"} {
		repo := newRepo(t)
		id := "task-20260922T100000Z-bad"
		publishDemo(t, repo, id, "LOW")
		r := qrun("claim", "--task", id, "--capability", bad, "--role", "executor",
			"--model", "test", "--tool", "gotest", "--session-id", "s", "--repo", repo, "--now", t0)
		if r.code == 0 {
			t.Fatalf("capability %q must be rejected", bad)
		}
		if !strings.Contains(r.err, "session_capability") {
			t.Fatalf("rejection for %q should name session_capability: %s", bad, r.err)
		}
		if strings.Contains(r.err, "master") || strings.Contains(r.err, "slave") {
			t.Fatalf("rejection must not emit legacy terms: %s", r.err)
		}
	}

	// A HIGH alias may claim a HIGH-only task; work_role stays independent of
	// the capability declaration.
	repo := newRepo(t)
	id := "task-20260922T100000Z-highonly-alias"
	publishDemo(t, repo, id, "HIGH")
	mustQ(t, "claim", "--task", id, "--capability", "leader", "--role", "reviewer",
		"--model", "test", "--tool", "gotest", "--session-id", "s", "--repo", repo, "--now", t0)
	if got := stateIn(t, repo, id); got != model.StateClaimed {
		t.Fatalf("HIGH alias failed to claim a HIGH task: %s", got)
	}
}
