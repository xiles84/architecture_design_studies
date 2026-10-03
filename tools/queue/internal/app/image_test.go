package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestImageSourceRejectsDriftAndNormalizesCRLF(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"go.mod", "main.go", "Containerfile", "queue", "internal/app/example.go"} {
		writeFile(t, filepath.Join(root, name), name+"\n")
	}
	digest, err := queueSourceDigest(root)
	if err != nil {
		t.Fatal(err)
	}
	previous := BuildSourceDigest
	t.Cleanup(func() { BuildSourceDigest = previous })
	BuildSourceDigest = digest
	mustQ(t, "image-check", "--source", root)
	writeFile(t, filepath.Join(root, "main.go"), "main.go\r\n")
	mustQ(t, "image-check", "--source", root)
	writeFile(t, filepath.Join(root, "internal/app/example.go"), "changed semantics\n")
	r := qrun("image-check", "--source", root)
	if r.code != 4 || !strings.Contains(r.err, "--build") {
		t.Fatalf("stale image was not rejected: %+v", r)
	}
	BuildSourceDigest = "unknown"
	if r := qrun("image-check", "--source", root); r.code != 4 {
		t.Fatalf("unmarked image accepted: %+v", r)
	}
	if err := os.Remove(filepath.Join(root, "go.mod")); err != nil {
		t.Fatal(err)
	}
	if r := qrun("source-digest", "--source", root); r.code != 4 {
		t.Fatalf("incomplete source accepted: %+v", r)
	}
}
