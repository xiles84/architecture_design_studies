package evidence

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 8; i++ {
		if _, err := os.Stat(filepath.Join(dir, "book", "evidence", "claims.schema.json")); err == nil {
			return dir
		}
		dir = filepath.Dir(dir)
	}
	t.Fatal("could not locate the repository root")
	return ""
}

func loadSchema(t *testing.T, root string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(root, "book", "evidence", "claims.schema.json"))
	if err != nil {
		t.Fatal(err)
	}
	var schema map[string]any
	if err := json.Unmarshal(data, &schema); err != nil {
		t.Fatal(err)
	}
	return schema
}

func TestRealRegistryValidates(t *testing.T) {
	root := repoRoot(t)
	schema := loadSchema(t, root)
	doc, raw, err := Load(filepath.Join(root, "book", "evidence", "claims.json"))
	if err != nil {
		t.Fatal(err)
	}
	if errs := Validate(root, schema, doc, raw); len(errs) > 0 {
		for _, e := range errs {
			t.Errorf("unexpected error: %v", e)
		}
		t.Fatalf("the committed registry must validate cleanly")
	}
	if len(doc.BookInputs()) != len(doc.Claims) {
		t.Fatalf("no claim is superseded, so all %d must be book inputs", len(doc.Claims))
	}
}

func TestNegativeFixtures(t *testing.T) {
	root := repoRoot(t)
	schema := loadSchema(t, root)
	cases := []struct {
		file string
		want string
	}{
		{"failed-cell-winner.json", "cannot be a winner"},
		{"stale-digest.json", "stale digest"},
		{"cross-study-without-record.json", "comparability record"},
		{"missing-topology.json", "schema"},
		{"superseded-not-marked.json", "superseded_by"},
	}
	for _, tc := range cases {
		tc := tc
		t.Run(tc.file, func(t *testing.T) {
			doc, raw, err := Load(filepath.Join("..", "..", "testdata", tc.file))
			if err != nil {
				t.Fatal(err)
			}
			errs := Validate(root, schema, doc, raw)
			if len(errs) == 0 {
				t.Fatalf("fixture %s validated but must fail", tc.file)
			}
			for _, e := range errs {
				if strings.Contains(e.Error(), tc.want) {
					return
				}
			}
			t.Fatalf("fixture %s: no error mentioned %q; got %v", tc.file, tc.want, errs)
		})
	}
}

func TestSchemaRejectsUnknownStrength(t *testing.T) {
	schema := map[string]any{
		"type": "object",
		"properties": map[string]any{
			"strength": map[string]any{"enum": []any{"gap"}},
		},
	}
	doc := map[string]any{"strength": "invented"}
	if errs := ValidateSchema(schema, doc); len(errs) == 0 {
		t.Fatal("an out-of-enum value must fail")
	}
}
