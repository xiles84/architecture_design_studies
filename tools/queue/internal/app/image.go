package app

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// BuildSourceDigest is set by the Containerfile. An unmarked or stale binary
// must not silently apply obsolete claim semantics to current task records.
var BuildSourceDigest = "unknown"

func queueSourceDigest(root string) (string, error) {
	for _, name := range []string{"go.mod", "main.go", "Containerfile", "queue"} {
		if _, err := os.Stat(filepath.Join(root, name)); err != nil {
			return "", fmt.Errorf("queue source %s: %w", name, err)
		}
	}
	var names []string
	err := filepath.WalkDir(root, func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if strings.HasSuffix(rel, ".go") || rel == "go.mod" || rel == "go.sum" || rel == "Containerfile" || rel == "queue" {
			names = append(names, rel)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(names)
	h := sha256.New()
	for _, name := range names {
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(name)))
		if err != nil {
			return "", err
		}
		// Native Windows checkouts and Linux builds must identify the same source.
		data = bytes.ReplaceAll(data, []byte("\r\n"), []byte("\n"))
		fmt.Fprintf(h, "%d:%s:%d:", len(name), name, len(data))
		h.Write(data)
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func cmdImageSource(command string, args []string, out io.Writer) error {
	flags := newFlagSet(command)
	root := flags.String("source", "tools/queue", "queue source directory")
	if err := flags.Parse(args); err != nil {
		return failf(2, "%v", err)
	}
	digest, err := queueSourceDigest(*root)
	if err != nil {
		return failf(4, "cannot verify queue image source: %v", err)
	}
	if command == "image-check" && digest != BuildSourceDigest {
		return failf(4, "queue image/source mismatch: image=%s source=%s; run tools/queue/queue --build before continuing", BuildSourceDigest, digest)
	}
	fmt.Fprintln(out, digest)
	return nil
}
