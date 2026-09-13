// Package filestore is the output adapter for results: JSON for machines, text
// for people. Plans are written as readable text rather than escaped JSON
// strings -- they are the most informative artefact a run produces.
package filestore

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func WriteJSON(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(b, '\n'), 0o644)
}

func WriteText(path, s string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(s), 0o644)
}

// ReadJSONTree decodes every *.json file under dir into a fresh T. Files that
// fail to decode are an error: a report silently skipping a result would
// present an incomplete matrix as a complete one.
func ReadJSONTree[T any](dir string) ([]T, error) {
	var out []T
	err := filepath.WalkDir(dir, func(path string, e os.DirEntry, err error) error {
		if err != nil || e.IsDir() || filepath.Ext(path) != ".json" {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var v T
		if err := json.Unmarshal(b, &v); err != nil {
			return &os.PathError{Op: "decode", Path: path, Err: err}
		}
		out = append(out, v)
		return nil
	})
	return out, err
}
