package selectbench

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

// Parameter binding for the embedded catalogue.
//
// The catalogue uses `:name` placeholders (`:pred`, `:n`, `:limit`, ...). The
// harness must bind them without ambiguity: PostgreSQL's `::type` cast is not a
// placeholder, and a statement that silently loses a parameter would run with the
// wrong predicate. Bind is therefore strict — every placeholder must have a
// value and every supplied value must be used — and it is the single place quoted
// literals are produced, so the DB mode and the checks agree.

// castShield stands in for `::` while placeholders are matched, because Go's
// regexp has no lookbehind and `::text` would otherwise look like `:text`.
const castShield = "\x00ads-cast\x00"

var placeholderRe = regexp.MustCompile(`:([a-zA-Z_][a-zA-Z0-9_]*)`)

// expressionParams are inlined verbatim: they are SQL expressions (a predicate),
// not values, so quoting them would turn `WHERE :pred` into a constant. Every
// other parameter is a value and is rendered as a literal.
var expressionParams = map[string]bool{"pred": true}

func sqlLiteral(v string) string {
	if _, err := strconv.ParseFloat(v, 64); err == nil {
		return v
	}
	return "'" + strings.ReplaceAll(v, "'", "''") + "'"
}

// Bind substitutes every `:name` placeholder with a SQL literal and returns the
// statement plus the sorted list of parameters it consumed. It errors on a
// missing value, an unused value, or an empty statement.
func Bind(sql string, params map[string]string) (string, []string, error) {
	if strings.TrimSpace(sql) == "" {
		return "", nil, fmt.Errorf("empty statement")
	}
	shielded := strings.ReplaceAll(sql, "::", castShield)

	used := map[string]bool{}
	var missing []string
	out := placeholderRe.ReplaceAllStringFunc(shielded, func(m string) string {
		name := m[1:]
		v, ok := params[name]
		if !ok {
			missing = append(missing, name)
			return m
		}
		used[name] = true
		if expressionParams[name] {
			return v
		}
		return sqlLiteral(v)
	})
	if len(missing) > 0 {
		sort.Strings(missing)
		return "", nil, fmt.Errorf("unbound parameter(s): %s", strings.Join(dedupe(missing), ", "))
	}
	var unused []string
	for k := range params {
		if !used[k] {
			unused = append(unused, k)
		}
	}
	if len(unused) > 0 {
		sort.Strings(unused)
		return "", nil, fmt.Errorf("supplied but unused parameter(s): %s", strings.Join(unused, ", "))
	}
	names := make([]string, 0, len(used))
	for k := range used {
		names = append(names, k)
	}
	sort.Strings(names)
	return strings.ReplaceAll(out, castShield, "::"), names, nil
}

func dedupe(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
