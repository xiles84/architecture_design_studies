package selectbench

import (
	"strings"
	"testing"
)

func TestBindSubstitutesExpression(t *testing.T) {
	got, used, err := Bind("SELECT id FROM orders WHERE (:pred)", map[string]string{"pred": "amount > 100"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "SELECT id FROM orders WHERE (amount > 100)" {
		t.Errorf("got %q", got)
	}
	if len(used) != 1 || used[0] != "pred" {
		t.Errorf("used = %v, want [pred]", used)
	}
}

func TestBindQuotesStringsAndKeepsNumbers(t *testing.T) {
	got, _, err := Bind("LIMIT :n AND name = :name", map[string]string{"n": "3", "name": "O'Brien"})
	if err != nil {
		t.Fatal(err)
	}
	want := "LIMIT 3 AND name = 'O''Brien'"
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// A `::type` cast must never be mistaken for a placeholder.
func TestBindLeavesCastsAlone(t *testing.T) {
	got, _, err := Bind("x::text = :v::int", map[string]string{"v": "3"})
	if err != nil {
		t.Fatal(err)
	}
	if got != "x::text = 3::int" {
		t.Errorf("got %q, want %q", got, "x::text = 3::int")
	}
}

func TestBindErrorsOnMissingAndUnused(t *testing.T) {
	if _, _, err := Bind("WHERE :pred", map[string]string{}); err == nil {
		t.Error("a missing parameter did not error")
	}
	if _, _, err := Bind("WHERE 1=1", map[string]string{"pred": "x"}); err == nil {
		t.Error("an unused parameter did not error")
	}
	if _, _, err := Bind("   ", nil); err == nil {
		t.Error("an empty statement did not error")
	}
}

// Every catalogue statement's declared `params:` must match its placeholders:
// a mismatch would mean the harness binds a statement the survey never intended.
func TestCatalogueParamsMatchPlaceholders(t *testing.T) {
	c := load(t)
	vals := map[string]string{
		"pred":        "amount > 100",
		"n":           "3",
		"limit":       "2",
		"offset":      "0",
		"last_amount": "100",
		"last_id":     "5",
	}
	for _, name := range c.Order {
		s := c.Statements[name]
		params := map[string]string{}
		for _, p := range s.Params {
			v, ok := vals[p]
			if !ok {
				v = "1"
			}
			params[p] = v
		}
		if _, _, err := Bind(s.SQL, params); err != nil {
			t.Errorf("%s: %v", name, err)
		}
		if len(s.Params) == 0 && strings.Contains(s.SQL, ":") {
			// no declared params, but the SQL has a placeholder-like token
			for _, tok := range []string{":pred", ":n", ":limit", ":offset"} {
				if strings.Contains(s.SQL, tok) {
					t.Errorf("%s: SQL contains %s but declares no params", name, tok)
				}
			}
		}
	}
}
