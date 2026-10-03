package selectbench

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

// Row canonicalization: the bridge between database rows and the oracles.
//
// The DB mode scans a result set into []any per row. Before a cell can be
// verified it must be expressed as the same canonical value the oracle returns
// for that contract, including the declared order. This file is that mapping, so
// the live check and the unit checks agree exactly on what "the same answer"
// means.

func toInt(v any) (int, bool, error) { // (value, isNull, error)
	switch x := v.(type) {
	case nil:
		return 0, true, nil
	case int:
		return x, false, nil
	case int32:
		return int(x), false, nil
	case int64:
		return int(x), false, nil
	case float64:
		if math.IsNaN(x) || math.IsInf(x, 0) || math.Trunc(x) != x || x >= math.Exp2(63) || x < -math.Exp2(63) {
			return 0, false, fmt.Errorf("not an exact integer: %v", x)
		}
		return int(x), false, nil
	case string:
		n, err := strconv.Atoi(x)
		if err != nil {
			return 0, false, fmt.Errorf("cannot read %q as an integer", x)
		}
		return n, false, nil
	default:
		return 0, false, fmt.Errorf("unsupported value type %T", v)
	}
}

func toString(v any) (string, error) {
	switch x := v.(type) {
	case string:
		return x, nil
	case []byte:
		return string(x), nil
	case nil:
		return "", nil
	default:
		return "", fmt.Errorf("unsupported text type %T", v)
	}
}

func need(cols []any, n int, group string) error {
	if len(cols) != n {
		return fmt.Errorf("%s: expected %d columns, got %d", group, n, len(cols))
	}
	return nil
}

// Canonicalize converts scanned rows into the canonical value for a contract.
// Rows are sorted exactly as the oracle sorts them, so verification compares
// answers rather than row arrival order.
func Canonicalize(group string, rows [][]any) (any, error) {
	switch group {
	case "sel_full_rows":
		out := make([]FullRow, 0, len(rows))
		for _, r := range rows {
			if err := need(r, 4, group); err != nil {
				return nil, err
			}
			cid, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			name, err := toString(r[1])
			if err != nil {
				return nil, err
			}
			oid, _, err := toInt(r[2])
			if err != nil {
				return nil, err
			}
			amt, isNull, err := toInt(r[3])
			if err != nil {
				return nil, err
			}
			var amount *int
			if !isNull {
				amount = ip(amt)
			}
			out = append(out, FullRow{cid, name, oid, amount})
		}
		sort.Slice(out, func(i, j int) bool {
			if out[i].CustomerID != out[j].CustomerID {
				return out[i].CustomerID < out[j].CustomerID
			}
			return out[i].OrderID < out[j].OrderID
		})
		return out, nil

	case "sel_membership", "sel_absence":
		out := make([]int, 0, len(rows))
		for _, r := range rows {
			if err := need(r, 2, group); err != nil {
				return nil, err
			}
			cid, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			out = append(out, cid)
		}
		sort.Ints(out)
		return out, nil

	case "sel_optional_match":
		out := make([]OptRow, 0, len(rows))
		for _, r := range rows {
			if err := need(r, 4, group); err != nil {
				return nil, err
			}
			cid, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			name, err := toString(r[1])
			if err != nil {
				return nil, err
			}
			oid, snull, err := toInt(r[2])
			if err != nil {
				return nil, err
			}
			amt, anull, err := toInt(r[3])
			if err != nil {
				return nil, err
			}
			row := OptRow{CustomerID: cid, CustomerName: name}
			if !snull {
				row.OrderID = ip(oid)
			}
			if !anull {
				row.Amount = ip(amt)
			}
			out = append(out, row)
		}
		sort.Slice(out, func(i, j int) bool {
			if out[i].CustomerID != out[j].CustomerID {
				return out[i].CustomerID < out[j].CustomerID
			}
			if out[i].OrderID == nil || out[j].OrderID == nil {
				return out[i].OrderID == nil && out[j].OrderID != nil
			}
			return *out[i].OrderID < *out[j].OrderID
		})
		return out, nil

	case "sel_aggregate":
		out := make([]AggRow, 0, len(rows))
		for _, r := range rows {
			if err := need(r, 4, group); err != nil {
				return nil, err
			}
			cid, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			name, err := toString(r[1])
			if err != nil {
				return nil, err
			}
			total, _, err := toInt(r[2])
			if err != nil {
				return nil, err
			}
			children, _, err := toInt(r[3])
			if err != nil {
				return nil, err
			}
			out = append(out, AggRow{cid, name, total, children})
		}
		sort.Slice(out, func(i, j int) bool { return out[i].CustomerID < out[j].CustomerID })
		return out, nil

	case "sel_topn":
		out := make([]TopRow, 0, len(rows))
		for _, r := range rows {
			if err := need(r, 3, group); err != nil {
				return nil, err
			}
			cid, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			oid, _, err := toInt(r[1])
			if err != nil {
				return nil, err
			}
			amt, anull, err := toInt(r[2])
			if err != nil {
				return nil, err
			}
			row := TopRow{CustomerID: cid, OrderID: oid}
			if !anull {
				row.Amount = ip(amt)
			}
			out = append(out, row)
		}
		// Both top-N templates declare parent/amount/id output order explicitly.
		return out, nil

	case "sel_pagination":
		out := make([]PageRow, 0, len(rows))
		for _, r := range rows {
			if err := need(r, 2, group); err != nil {
				return nil, err
			}
			oid, _, err := toInt(r[0])
			if err != nil {
				return nil, err
			}
			amt, anull, err := toInt(r[1])
			if err != nil {
				return nil, err
			}
			row := PageRow{OrderID: oid}
			if !anull {
				row.Amount = ip(amt)
			}
			out = append(out, row)
		}
		// The SQL promises a total order. Preserve arrival order so an incorrect
		// ORDER BY cannot pass verification after the harness sorts it correctly.
		return out, nil

	case ControlGroup:
		return nil, fmt.Errorf("a deliberately wrong control has no canonical answer to verify against")
	}
	return nil, fmt.Errorf("unknown semantic group %q", group)
}

// VerifyRows canonicalizes scanned database rows and verifies them against the
// statement's contract oracle. It returns an error for a different answer, a
// wrong column shape, or a control.
func VerifyRows(c *Catalog, name string, fx Fixture, a Args, rows [][]any) error {
	s, ok := c.Statements[name]
	if !ok {
		return fmt.Errorf("unknown statement %q", name)
	}
	got, err := Canonicalize(s.Group, rows)
	if err != nil {
		return fmt.Errorf("%s: %w", name, err)
	}
	return Verify(c, name, fx, a, got)
}
