package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"testing"

	bs "adsqueue/internal/brainstorm"
)

const brainstormID = "brainstorm-20260922T100000Z-demo"

func highArgs() []string {
	return []string{"--capability", "HIGH", "--role", "planner", "--model", "test", "--tool", "gotest", "--session-id", "brainstorm-test"}
}

func startBrainstorm(t *testing.T, repo string, positions, critiques int) {
	t.Helper()
	args := []string{"brainstorm-start", "--repo", repo, "--now", t0, "--id", brainstormID, "--title", "Demo decision", "--question", "Which option is best?", "--positions", strconv.Itoa(positions), "--critiques", strconv.Itoa(critiques)}
	args = append(args, highArgs()...)
	mustQ(t, args...)
}

func claimBrainstorm(t *testing.T, repo, now, lease string) bs.Claim {
	t.Helper()
	args := []string{"brainstorm-claim", "--repo", repo, "--now", now, "--brainstorm", brainstormID, "--lease", lease, "--json"}
	args = append(args, highArgs()...)
	r := mustQ(t, args...)
	var claim bs.Claim
	if err := json.Unmarshal([]byte(r.out), &claim); err != nil {
		t.Fatalf("decode claim: %v\n%s", err, r.out)
	}
	return claim
}

func submitBrainstorm(t *testing.T, repo, now string, claim bs.Claim, content, summary string) qres {
	t.Helper()
	draft := filepath.Join(t.TempDir(), claim.SlotID+".md")
	writeFile(t, draft, content)
	return mustQ(t, "brainstorm-submit", "--repo", repo, "--now", now, "--brainstorm", brainstormID,
		"--slot", claim.SlotID, "--claim-id", claim.ClaimID, "--claim-epoch", strconv.Itoa(claim.ClaimEpoch),
		"--file", draft, "--summary", summary, "--disagreements", "One trade-off remains.", "--confidence", "high")
}

func TestBrainstormLifecycleAndFreshClone(t *testing.T) {
	repo := newRepo(t)
	startBrainstorm(t, repo, 3, 2)

	if got := mustQ(t, "brainstorm-list", "--repo", repo, "--now", t0).out; !strings.Contains(got, "ONGOING") || !strings.Contains(got, brainstormID) || !strings.Contains(got, "CONCLUDED\n\n(none)") {
		t.Fatalf("unexpected initial list:\n%s", got)
	}
	for i, now := range []string{t1, t2, t3} {
		cl := claimBrainstorm(t, repo, now, "2h")
		if cl.SlotKind != bs.SlotPosition {
			t.Fatalf("claim %d kind = %s", i, cl.SlotKind)
		}
		submitBrainstorm(t, repo, now, cl, fmt.Sprintf("## Evidence examined\n\n- source-%d\n\n## Recommendation\n\nOption %d.", i, i), fmt.Sprintf("Position %d favors option %d.", i+1, i))
	}
	for i, now := range []string{t4, t5} {
		cl := claimBrainstorm(t, repo, now, "2h")
		if cl.SlotKind != bs.SlotCritique {
			t.Fatalf("critique claim %d kind = %s", i, cl.SlotKind)
		}
		submitBrainstorm(t, repo, now, cl, "## Evidence examined\n\nAll positions.\n\n## Critique\n\nCompare the risks.", fmt.Sprintf("Cross-review %d narrows the choice.", i+1))
	}
	cl := claimBrainstorm(t, repo, t6, "2h")
	if cl.SlotKind != bs.SlotSynthesis {
		t.Fatalf("synthesis claim kind = %s", cl.SlotKind)
	}
	ended := submitBrainstorm(t, repo, t6, cl, "## Decision\n\nChoose option 1.\n\n## Remaining dissent\n\nThe cost trade-off remains.\n\n## Task-ready next steps\n\nPrepare implementation tasks only when requested.", "Choose option 1, with a follow-up cost check.")
	if !strings.Contains(ended.out, "BRAINSTORM ENDED") || !strings.Contains(ended.out, "Tasks: NOT CREATED") || !strings.Contains(ended.out, "CREATE TASKS FROM BRAINSTORM "+brainstormID) {
		t.Fatalf("missing terminal output:\n%s", ended.out)
	}

	status := mustQ(t, "brainstorm-status", "--repo", repo, "--now", t6, "--brainstorm", brainstormID).out
	if !strings.Contains(status, "State: concluded") || !strings.Contains(status, "positions 3/3; critiques 2/2; synthesis complete") {
		t.Fatalf("unexpected status:\n%s", status)
	}
	mustQ(t, "brainstorm-audit", "--repo", repo, "--now", t6)

	recordDir, err := bs.Dir(repo, brainstormID)
	if err != nil {
		t.Fatal(err)
	}
	for _, rel := range []string{"brainstorm.json", "QUESTION.md", "STATE.md", "CONCLUSION.md", "positions", "critiques", "events"} {
		if _, err := os.Stat(filepath.Join(recordDir, rel)); err != nil {
			t.Fatalf("missing %s: %v", rel, err)
		}
	}

	clone := filepath.Join(t.TempDir(), "clone")
	git(t, repo, "clone", "-q", repo, clone)
	if got := mustQ(t, "brainstorm-status", "--repo", clone, "--now", t6, "--brainstorm", brainstormID).out; !strings.Contains(got, "State: concluded") {
		t.Fatalf("fresh clone did not reconstruct state:\n%s", got)
	}
	mustQ(t, "brainstorm-audit", "--repo", clone, "--now", t6)
}

func TestBrainstormClaimsUseIndependentSlotsAndRejectOldEpoch(t *testing.T) {
	repo := newRepo(t)
	startBrainstorm(t, repo, 2, 1)

	const claimers = 20
	results := make(chan qres, claimers)
	var wg sync.WaitGroup
	for i := 0; i < claimers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			args := []string{"brainstorm-claim", "--repo", repo, "--now", t1, "--brainstorm", brainstormID, "--lease", "1m", "--json", "--capability", "HIGH", "--role", "analyst", "--model", fmt.Sprintf("m%d", i), "--tool", "gotest", "--session-id", fmt.Sprintf("s%d", i)}
			results <- qrun(args...)
		}(i)
	}
	wg.Wait()
	close(results)
	var winners []bs.Claim
	for r := range results {
		if r.code == 0 {
			var cl bs.Claim
			if err := json.Unmarshal([]byte(r.out), &cl); err != nil {
				t.Fatal(err)
			}
			winners = append(winners, cl)
		}
	}
	if len(winners) != 2 || winners[0].SlotID == winners[1].SlotID {
		t.Fatalf("expected exactly two distinct slot winners, got %#v", winners)
	}

	// Once the leases expire, the first slot is recovered with a higher epoch.
	recovered := claimBrainstorm(t, repo, t3, "2h")
	var old bs.Claim
	for _, cl := range winners {
		if cl.SlotID == recovered.SlotID {
			old = cl
		}
	}
	if recovered.ClaimEpoch != old.ClaimEpoch+1 {
		t.Fatalf("recovery epoch = %d, old = %d", recovered.ClaimEpoch, old.ClaimEpoch)
	}
	r := qrun("brainstorm-guard", "--repo", repo, "--now", t3, "--brainstorm", brainstormID, "--slot", old.SlotID, "--claim-id", old.ClaimID, "--claim-epoch", strconv.Itoa(old.ClaimEpoch))
	if r.code == 0 || !strings.Contains(r.err, "superseded") {
		t.Fatalf("old claimant guard was not rejected: %#v", r)
	}
}

func TestBrainstormCapabilityIntentAndTaskBoundary(t *testing.T) {
	repo := newRepo(t)
	low := qrun("brainstorm-start", "--repo", repo, "--now", t0, "--id", brainstormID, "--title", "Demo", "--question", "Choose.", "--positions", "2", "--critiques", "1",
		"--capability", "LOW", "--role", "executor", "--model", "test", "--tool", "gotest", "--session-id", "low")
	if low.code == 0 || !strings.Contains(low.err, "only a HIGH session") {
		t.Fatalf("LOW start was not rejected: %#v", low)
	}
	startBrainstorm(t, repo, 2, 1)
	link := qrun("brainstorm-link-tasks", "--repo", repo, "--now", t1, "--brainstorm", brainstormID, "--task", "task-20260922T100000Z-demo",
		"--capability", "HIGH", "--role", "planner", "--model", "test", "--tool", "gotest", "--session-id", "high")
	if link.code == 0 || !strings.Contains(link.err, "only after conclusion") {
		t.Fatalf("pre-conclusion task linking was not rejected: %#v", link)
	}
	if got := mustQ(t, "brainstorm-intent", "--repo", repo, "--text", "What should we improve?").out; !strings.Contains(got, "no brainstorm intent") {
		t.Fatalf("ordinary question triggered brainstorm: %s", got)
	}
	if got := mustQ(t, "brainstorm-intent", "--repo", repo, "--text", "CONTINUE BRAINSTORM "+brainstormID).out; !strings.HasPrefix(got, "continue\t"+brainstormID) {
		t.Fatalf("explicit continue not parsed: %s", got)
	}
}

func TestBrainstormLinksOnlyCorrelatedPublishedTasks(t *testing.T) {
	repo := newRepo(t)
	startBrainstorm(t, repo, 2, 1)
	for i, now := range []string{t1, t2} {
		cl := claimBrainstorm(t, repo, now, "2h")
		submitBrainstorm(t, repo, now, cl, fmt.Sprintf("Position %d.", i+1), fmt.Sprintf("Position %d.", i+1))
	}
	cl := claimBrainstorm(t, repo, t3, "2h")
	submitBrainstorm(t, repo, t3, cl, "Both positions were checked.", "The cross-review retains one objection.")
	cl = claimBrainstorm(t, repo, t4, "2h")
	submitBrainstorm(t, repo, t4, cl, "Choose the first position and retain the objection.", "Choose the first position.")

	taskID := "task-20260922T100000Z-from-brainstorm"
	tmp := t.TempDir()
	spec := filepath.Join(tmp, "task.json")
	brief := filepath.Join(tmp, "BRIEF.md")
	raw := specJSON(taskID, "LOW", "repo/from-brainstorm", ".worktrees/from-brainstorm", "repo/from-brainstorm")
	raw = strings.Replace(raw, "\n}", fmt.Sprintf(",\n  \"originating_brainstorm_id\": %q\n}", brainstormID), 1)
	writeFile(t, spec, raw)
	writeFile(t, brief, "# Brief\n\nImplement the concluded decision.\n")
	args := []string{"publish", "--repo", repo, "--now", t5, "--spec", spec, "--brief", brief}
	args = append(args, highArgs()...)
	mustQ(t, args...)

	linked := mustQ(t, "brainstorm-link-tasks", "--repo", repo, "--now", t6, "--brainstorm", brainstormID, "--task", taskID,
		"--capability", "HIGH", "--role", "planner", "--model", "test", "--tool", "gotest", "--session-id", "high")
	if !strings.Contains(linked.out, "State: tasked") || !strings.Contains(linked.out, "Stage: tasked") {
		t.Fatalf("task was not linked:\n%s", linked.out)
	}
	record, err := bs.Load(repo, brainstormID)
	if err != nil {
		t.Fatal(err)
	}
	if got := record.LinkedTasks(); len(got) != 1 || got[0] != taskID {
		t.Fatalf("linked tasks = %v", got)
	}
	mustQ(t, "brainstorm-audit", "--repo", repo, "--now", t6)
}
