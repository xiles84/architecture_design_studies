package app

import (
	"path/filepath"
	"testing"

	"adsqueue/internal/gitx"
	"adsqueue/internal/model"
)

func TestRecoveredClaimCanSubmitWithoutExplicitFlags(t *testing.T) {
	repo := newRepo(t)
	id := "task-20260922T100000Z-recovered-submit"
	publishDemo(t, repo, id, "LOW")
	old := claimDemo(t, repo, id, t0)
	wt := filepath.Join(repo, ".worktrees", "demo")
	for _, now := range []string{"2026-09-22T13:00:00Z", "2026-09-22T13:10:01Z"} {
		mustQ(t, "recover", "--task", id, "--capability", "LOW", "--role", "executor",
			"--model", "test2", "--tool", "gotest", "--session-id", "s2", "--repo", wt, "--now", now)
	}
	gitRepo, err := gitx.Open(wt)
	if err != nil {
		t.Fatal(err)
	}
	ctx := &cmdCtx{Repo: gitRepo, Now: mustTime(t, "2026-09-22T13:10:01Z")}
	current, err := ctx.loadClaimContext(id)
	if err != nil {
		t.Fatal(err)
	}
	if current.ClaimEpoch != old.ClaimEpoch+1 || current.ClaimID == old.ClaimID || current.LastEvent == "" {
		t.Fatalf("recovery did not persist the current claim/event: %+v", current)
	}
	mustQ(t, "guard", "--task", id, "--repo", wt, "--now", "2026-09-22T13:10:01Z")
	rec, err := ctx.taskRecord(id)
	if err != nil {
		t.Fatal(err)
	}
	result := filepath.Join(rec.Dir, "attempts", current.AttemptID, "RESULT.md")
	rel, err := filepath.Rel(wt, result)
	if err != nil {
		t.Fatal(err)
	}
	commitFile(t, wt, rel, "Recovered worker result.\n", "result")
	mustQ(t, "submit", "--task", id, "--repo", wt, "--now", "2026-09-22T13:11:00Z")
	if got := stateIn(t, wt, id); got != model.StateAwaitingReview {
		t.Fatalf("recovered submission state: %s", got)
	}
}
