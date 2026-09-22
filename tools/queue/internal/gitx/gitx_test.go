package gitx

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func gitRun(t *testing.T, dir string, args ...string) string {
	t.Helper()
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_CONFIG_GLOBAL=/dev/null", "GIT_CONFIG_SYSTEM=/dev/null")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %s: %v\n%s", strings.Join(args, " "), err, out)
	}
	return strings.TrimSpace(string(out))
}

// TestBranchOIDDisambiguatesTag covers the queue's own naming convention:
// required_tag is usually the canonical branch name, so a bare `rev-parse
// repo/<topic>` resolves the tag object. Every branch lookup must go through
// refs/heads, which BranchOID does.
func TestBranchOIDDisambiguatesTag(t *testing.T) {
	dir := t.TempDir()
	gitRun(t, dir, "init", "-q", "-b", "main")
	gitRun(t, dir, "config", "user.email", "test@example.com")
	gitRun(t, dir, "config", "user.name", "Queue Test")
	gitRun(t, dir, "config", "commit.gpgsign", "false")
	if err := os.WriteFile(filepath.Join(dir, "a.txt"), []byte("a\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "a.txt")
	gitRun(t, dir, "commit", "-q", "-m", "base")

	// A branch and an annotated tag share the name, and their tips differ.
	gitRun(t, dir, "branch", "repo/topic")
	gitRun(t, dir, "tag", "-a", "-m", "integrated", "repo/topic")
	gitRun(t, dir, "checkout", "-q", "repo/topic")
	if err := os.WriteFile(filepath.Join(dir, "b.txt"), []byte("b\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitRun(t, dir, "add", "b.txt")
	gitRun(t, dir, "commit", "-q", "-m", "branch work")
	gitRun(t, dir, "checkout", "-q", "main")

	r, err := Open(dir)
	if err != nil {
		t.Fatal(err)
	}
	want := gitRun(t, dir, "rev-parse", "refs/heads/repo/topic")
	got, err := r.BranchOID("repo/topic")
	if err != nil {
		t.Fatalf("BranchOID: %v", err)
	}
	if got != want {
		t.Fatalf("BranchOID = %s, want the branch tip %s", got, want)
	}
	if !r.TagExists("repo/topic") {
		t.Fatal("TagExists did not see the tag")
	}
	// The bare name must not silently be treated as the branch: whatever it
	// resolves to (the tag object here), BranchOID still returns the branch.
	bare := gitRun(t, dir, "rev-parse", "refs/tags/repo/topic")
	if bare == want {
		t.Fatal("test setup: tag and branch tips are identical, ambiguity not exercised")
	}
	if got == bare {
		t.Fatalf("BranchOID returned the tag object %s", bare)
	}
}
