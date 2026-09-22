// Package gitx is the queue's only interface to Git.
//
// The queue's atomicity comes from one primitive: `git update-ref <ref> <new>
// <observed-old>`. Git serialises ref updates with a lock file, so when several
// workers (or several git binaries: native Windows Git and WSL Git) act on the
// same common Git directory, exactly one of them wins a create with the
// zero object id as the expected old value. Every claim, lease extension,
// recovery step and integration lock is built on that primitive rather than on
// a central mutable file.
//
// The package deliberately shells out to `git`; there is no library dependency,
// so it behaves identically under the pinned Podman image, native Windows Git
// and WSL Git.
package gitx

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ErrCASConflict reports that a compare-and-swap lost: the ref no longer held
// the object id the caller observed. Callers treat this as "someone else won",
// never as a crash.
var ErrCASConflict = errors.New("compare-and-swap conflict: the ref changed since it was observed")

// ErrNotARepository reports that the start directory is not inside a Git work
// tree.
var ErrNotARepository = errors.New("not inside a Git work tree")

// Repo is one Git work tree sharing a common Git directory with every other
// work tree in the repository.
type Repo struct {
	// Root is the absolute working-tree root of the work tree the process is in.
	Root string
	// GitDir is the absolute git dir for this work tree (a per-worktree
	// directory under the common dir when Root is a linked work tree).
	GitDir string
	// CommonDir is the absolute shared git dir holding refs, objects and
	// worktree registrations. Two work trees on different branches share it.
	CommonDir string
	// Git is the git executable, overridable with ADS_QUEUE_GIT so a test can
	// drive a second binary (for example Windows git.exe) against the same
	// common directory.
	Git string
}

// runGit runs git without a Repo, for the discovery calls that establish one.
func runGit(git, dir string, args ...string) (string, error) {
	cmd := exec.Command(git, append([]string{"-c", "safe.directory=*"}, args...)...)
	cmd.Dir = dir
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("git %s: %w: %s", strings.Join(args, " "), err, strings.TrimSpace(errb.String()))
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

// Open locates the repository containing dir. It fails rather than guessing
// when dir is not inside a work tree.
func Open(dir string) (*Repo, error) {
	if dir == "" {
		var err error
		dir, err = os.Getwd()
		if err != nil {
			return nil, err
		}
	}
	git := os.Getenv("ADS_QUEUE_GIT")
	if git == "" {
		git = "git"
	}
	root, err := runGit(git, dir, "rev-parse", "--show-toplevel")
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrNotARepository, err)
	}
	root = strings.TrimSpace(root)
	gitDir, err := runGit(git, root, "rev-parse", "--absolute-git-dir")
	if err != nil {
		return nil, err
	}
	commonDir, err := runGit(git, root, "rev-parse", "--git-common-dir")
	if err != nil {
		return nil, err
	}
	if !filepath.IsAbs(commonDir) {
		commonDir = filepath.Join(root, commonDir)
	}
	return &Repo{
		Root:      filepath.Clean(root),
		GitDir:    filepath.Clean(strings.TrimSpace(gitDir)),
		CommonDir: filepath.Clean(commonDir),
		Git:       git,
	}, nil
}

// Git runs a git command in the repository root and returns trimmed stdout.
func (r *Repo) Run(args ...string) (string, error) {
	return r.RunIn(r.Root, args...)
}

// RunIn runs a git command in dir and returns trimmed stdout.
func (r *Repo) RunIn(dir string, args ...string) (string, error) {
	return runGit(r.Git, dir, args...)
}

// ZeroOID returns the all-zero object name for this repository's hash format,
// used as the expected old value when a ref must not yet exist.
func (r *Repo) ZeroOID() string {
	format, err := r.Run("rev-parse", "--show-object-format")
	if err != nil || strings.TrimSpace(format) == "sha256" {
		if strings.TrimSpace(format) == "sha256" {
			return strings.Repeat("0", 64)
		}
	}
	return strings.Repeat("0", 40)
}

// RefOID returns the object id a ref points at, or ok=false when it does not
// exist.
func (r *Repo) RefOID(ref string) (oid string, ok bool) {
	out, err := r.Run("rev-parse", "--verify", "--quiet", ref)
	if err != nil {
		return "", false
	}
	oid = strings.TrimSpace(out)
	if oid == "" {
		return "", false
	}
	return oid, true
}

// WriteBlob stores data as a Git blob and returns its object id. Live JSON is
// therefore a first-class Git object, addressable and garbage-collected like
// everything else.
func (r *Repo) WriteBlob(data []byte) (string, error) {
	cmd := exec.Command(r.Git, "-c", "safe.directory=*", "hash-object", "-w", "--stdin")
	cmd.Dir = r.Root
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	cmd.Stdin = bytes.NewReader(data)
	var out, errb bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("hash-object: %w: %s", err, strings.TrimSpace(errb.String()))
	}
	return strings.TrimSpace(out.String()), nil
}

// ReadBlob returns the bytes of a blob object.
func (r *Repo) ReadBlob(oid string) ([]byte, error) {
	out, err := r.Run("cat-file", "blob", oid)
	if err != nil {
		return nil, err
	}
	// Run trims trailing newlines; live JSON is written without a trailing
	// newline but a trimmed blob is still valid JSON. Return as bytes with a
	// newline restored is unnecessary, so hand back the trimmed text.
	return []byte(out), nil
}

// CAS atomically sets ref to newOID only if it currently holds oldOID. An empty
// oldOID means "must not exist". It creates a reflog entry so the live ref has
// the same auditability as a branch.
func (r *Repo) CAS(ref, newOID, oldOID, message string) error {
	expected := oldOID
	if expected == "" {
		expected = r.ZeroOID()
	}
	if _, err := r.Run("update-ref", "--create-reflog", "-m", message, ref, newOID, expected); err != nil {
		if looksLikeCASConflict(err) {
			return fmt.Errorf("%w: %s", ErrCASConflict, ref)
		}
		return err
	}
	return nil
}

// DeleteCASCAS removes ref only if it still holds oldOID.
func (r *Repo) DeleteCAS(ref, oldOID, message string) error {
	if _, err := r.Run("update-ref", "-m", message, "-d", ref, oldOID); err != nil {
		if looksLikeCASConflict(err) {
			return fmt.Errorf("%w: %s", ErrCASConflict, ref)
		}
		return err
	}
	return nil
}

func looksLikeCASConflict(err error) bool {
	if err == nil {
		return false
	}
	s := err.Error()
	return strings.Contains(s, "cannot lock ref") ||
		strings.Contains(s, "but expected") ||
		strings.Contains(s, "is at") ||
		strings.Contains(s, "unable to resolve") ||
		strings.Contains(s, "reference already exists")
}

// Head returns the current commit object id, or "" on an unborn branch.
func (r *Repo) Head() (string, error) {
	out, err := r.Run("rev-parse", "--verify", "--quiet", "HEAD")
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(out), nil
}

// CurrentBranch returns the checked-out branch name, or "" when detached.
func (r *Repo) CurrentBranch() (string, error) {
	out, err := r.Run("symbolic-ref", "--quiet", "--short", "HEAD")
	if err != nil {
		return "", nil
	}
	return strings.TrimSpace(out), nil
}

// CommitPaths stages exactly the named paths (repository-relative) and commits
// them, leaving any other session's uncommitted work untouched. It returns the
// new commit id.
func (r *Repo) CommitPaths(message string, paths ...string) (string, error) {
	if len(paths) == 0 {
		return "", errors.New("CommitPaths: no paths")
	}
	addArgs := append([]string{"add", "--"}, paths...)
	if _, err := r.Run(addArgs...); err != nil {
		return "", err
	}
	commitArgs := append([]string{"commit", "-m", message, "--"}, paths...)
	if _, err := r.Run(commitArgs...); err != nil {
		return "", err
	}
	return r.Head()
}

// StatusPorcelain returns `git status --porcelain` output. Empty means clean.
func (r *Repo) StatusPorcelain() (string, error) {
	return r.Run("status", "--porcelain")
}

// StatusPorcelainTracked is StatusPorcelain restricted to tracked files. It is
// what an integration precondition uses: another session's untracked working
// files in the main checkout are normal and must not block a merge, while a
// modified tracked file genuinely means the checkout is mid-edit.
func (r *Repo) StatusPorcelainTracked() (string, error) {
	return r.Run("status", "--porcelain", "--untracked-files=no")
}

// BranchExists reports whether a local branch exists.
func (r *Repo) BranchExists(branch string) bool {
	_, err := r.Run("show-ref", "--verify", "--quiet", "refs/heads/"+branch)
	return err == nil
}

// IsAncestor reports whether ancestor is reachable from descendant.
func (r *Repo) IsAncestor(ancestor, descendant string) bool {
	_, err := r.Run("merge-base", "--is-ancestor", ancestor, descendant)
	return err == nil
}

// BranchOID returns the commit at the tip of a local branch, resolved through
// refs/heads explicitly.
//
// A bare branch name is ambiguous when a tag of the same name exists, and this
// repository's tasks routinely set `required_tag` to the canonical branch name
// (`repo/<topic>`). In that case `git rev-parse <name>` resolves the tag object
// rather than the branch, so every branch lookup must go through refs/heads.
func (r *Repo) BranchOID(branch string) (string, error) {
	out, err := r.Run("rev-parse", "--verify", "refs/heads/"+branch)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(out), nil
}

// TagExists reports whether a tag ref exists.
func (r *Repo) TagExists(tag string) bool {
	_, err := r.Run("show-ref", "--verify", "--quiet", "refs/tags/"+tag)
	return err == nil
}

// CreateAnnotatedTag creates an annotated tag at target. An existing tag is an
// error: tags are never moved, deleted or reused.
func (r *Repo) CreateAnnotatedTag(tag, target, message string) error {
	if r.TagExists(tag) {
		return fmt.Errorf("tag %s already exists; tags are never moved or reused", tag)
	}
	_, err := r.Run("tag", "-a", "-m", message, tag, target)
	return err
}

// Worktree is one registered work tree.
type Worktree struct {
	Path   string
	Branch string
	Head   string
	Bare   bool
}

// Worktrees lists registered work trees.
func (r *Repo) Worktrees() ([]Worktree, error) {
	out, err := r.Run("worktree", "list", "--porcelain")
	if err != nil {
		return nil, err
	}
	var (
		list []Worktree
		cur  Worktree
	)
	flush := func() {
		if cur.Path != "" {
			list = append(list, cur)
		}
		cur = Worktree{}
	}
	for _, line := range strings.Split(out, "\n") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			flush()
			cur.Path = strings.TrimPrefix(line, "worktree ")
		case strings.HasPrefix(line, "HEAD "):
			cur.Head = strings.TrimPrefix(line, "HEAD ")
		case strings.HasPrefix(line, "branch "):
			cur.Branch = strings.TrimPrefix(strings.TrimPrefix(line, "branch "), "refs/heads/")
		case line == "bare":
			cur.Bare = true
		case line == "detached":
			// no branch
		}
	}
	flush()
	return list, nil
}

// WorktreeAdd creates a linked work tree at path on a new branch, based at
// base. It is the caller's job to make the new registration portable with
// MakeWorktreePathsRelative.
func (r *Repo) WorktreeAdd(path, branch, base string) error {
	if _, err := r.Run("worktree", "add", path, "-b", branch, base); err != nil {
		return err
	}
	return nil
}

// MakeWorktreePathsRelative rewrites the two files git uses to link a linked
// work tree to its common dir so that they hold relative paths.
//
// Git 2.43 -- the version on this machine's WSL side -- has no
// `worktree.useRelativePaths` support, so `git worktree add` writes absolute
// paths. An absolute path is only meaningful to the machine (and the symlink
// layout) that created it: the same registration then fails under native
// Windows Git or inside a Podman bind mount. Git reads either form; relative is
// the portable one. Both git binaries on this machine were tested against the
// result.
func (r *Repo) MakeWorktreePathsRelative(worktreePath string) error {
	worktreePath = filepath.Clean(worktreePath)
	gitFile := filepath.Join(worktreePath, ".git")
	data, err := os.ReadFile(gitFile)
	if err != nil {
		return fmt.Errorf("read %s: %w", gitFile, err)
	}
	const prefix = "gitdir: "
	text := strings.TrimSpace(string(data))
	if !strings.HasPrefix(text, prefix) {
		return fmt.Errorf("%s does not look like a worktree link file", gitFile)
	}
	gitdir := strings.TrimSpace(strings.TrimPrefix(text, prefix))
	if !filepath.IsAbs(gitdir) {
		gitdir = filepath.Join(worktreePath, gitdir)
	}
	gitdir = filepath.Clean(gitdir)

	// The worktree link file points at the per-worktree git dir.
	relGitdir, err := filepath.Rel(worktreePath, gitdir)
	if err != nil {
		return err
	}
	if err := os.WriteFile(gitFile, []byte(prefix+filepath.ToSlash(relGitdir)+"\n"), 0o644); err != nil {
		return err
	}

	// The per-worktree git dir's `gitdir` file points back at the link file.
	back, err := filepath.Rel(gitdir, gitFile)
	if err != nil {
		return err
	}
	backFile := filepath.Join(gitdir, "gitdir")
	if _, err := os.Stat(backFile); err != nil {
		return fmt.Errorf("expected %s to exist: %w", backFile, err)
	}
	if err := os.WriteFile(backFile, []byte(filepath.ToSlash(back)+"\n"), 0o644); err != nil {
		return err
	}
	return nil
}

// CommonRoot is the directory containing the shared git dir. Task records store
// `canonical_worktree` relative to this root, so the CLI must resolve it here
// rather than against the current (possibly linked) worktree root.
func (r *Repo) CommonRoot() string {
	return filepath.Dir(r.CommonDir)
}

// ResolveCommonPath resolves a path relative to the common repository root.
func (r *Repo) ResolveCommonPath(rel string) string {
	if filepath.IsAbs(rel) {
		return filepath.Clean(rel)
	}
	return filepath.Join(r.CommonRoot(), filepath.FromSlash(rel))
}

// RepoRelative converts an absolute or cwd-relative path into one relative to
// the repository root, which is what task records store and what a host wrapper
// can execute from the root.
func (r *Repo) RepoRelative(path string) (string, error) {
	if !filepath.IsAbs(path) {
		abs, err := filepath.Abs(path)
		if err != nil {
			return "", err
		}
		path = abs
	}
	return filepath.Rel(r.Root, filepath.Clean(path))
}

// ResolveRepoPath resolves a repository-relative path against the root.
func (r *Repo) ResolveRepoPath(rel string) string {
	if filepath.IsAbs(rel) {
		return filepath.Clean(rel)
	}
	return filepath.Join(r.Root, rel)
}
