package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// runManifest writes the explicit primary-input manifest for a correction run
// set: every result JSON with a normalized per-file SHA-256 and one corpus
// digest over the framed sorted paths, lengths and bytes. Pilots, failed cells
// and old fixed-order runs are excluded by the caller choosing the directory.
func runManifest(dir, reportOut string) error {
	if dir == "" {
		return fmt.Errorf("-cmd manifest needs -results")
	}
	paths, err := filepath.Glob(filepath.Join(dir, "*", "*.json"))
	if err != nil {
		return err
	}
	sort.Strings(paths)
	type entry struct {
		rel, tag, sha string
		bytes         int
	}
	var entries []entry
	h := sha256.New()
	corpus := 0
	for _, p := range paths {
		base := filepath.Base(p)
		if base != "access.json" && base != "aggregate.json" && base != "existence.json" {
			continue
		}
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		data := []byte(strings.ReplaceAll(string(raw), "\r\n", "\n"))
		sum := sha256.Sum256(data)
		rel := filepath.ToSlash(filepath.Join(filepath.Base(filepath.Dir(p)), base))
		tag := extractString(data, "run_tag")
		entries = append(entries, entry{rel: rel, tag: tag, sha: fmt.Sprintf("%x", sum), bytes: len(data)})
		fmt.Fprintf(h, "%s\n%d\n", rel, len(data))
		h.Write(data)
		corpus += len(data)
	}
	if len(entries) == 0 {
		return fmt.Errorf("no primary result JSONs under %s", dir)
	}
	var b strings.Builder
	b.WriteString("# Primary-input manifest — SELECT order correction\n\n")
	b.WriteString("Generated facts only. SHA-256 is over the file with CRLF normalized to LF. ")
	b.WriteString("The corpus digest frames sorted repository-relative paths and normalized byte lengths followed by the file bytes.\n\n")
	b.WriteString("| run | result file | producing run tag | SHA-256 | bytes |\n|---|---|---|---|---|\n")
	for _, e := range entries {
		fmt.Fprintf(&b, "| %s | `%s` | `%s` | `%s` | %d |\n",
			strings.SplitN(e.rel, "/", 2)[0], e.rel, e.tag, e.sha, e.bytes)
	}
	fmt.Fprintf(&b, "\nCorpus inputs_digest: `%x`\n\n%d files, %d normalized bytes.\n",
		h.Sum(nil), len(entries), corpus)
	if reportOut == "" {
		fmt.Print(b.String())
		return nil
	}
	return os.WriteFile(reportOut, []byte(b.String()), 0o644)
}

var runTagRE = regexp.MustCompile(`"run_tag"\s*:\s*"([^"]*)"`)

func extractString(data []byte, key string) string {
	if key == "run_tag" {
		if m := runTagRE.FindSubmatch(data); m != nil {
			return string(m[1])
		}
	}
	// Fall back to a generic decode for other keys.
	var probe map[string]any
	if err := json.Unmarshal(data, &probe); err == nil {
		if v, ok := probe[key].(string); ok {
			return v
		}
	}
	return ""
}
