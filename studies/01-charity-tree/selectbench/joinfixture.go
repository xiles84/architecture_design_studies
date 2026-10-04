package selectbench

import (
	"crypto/md5"
	"encoding/hex"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// Join/filter fixture and independent Go oracles (task select-join-cte).
//
// The join/CTE family needs a fact table large enough that a 0.01% selective
// predicate returns rows and a 90% predicate is a broad scan, while every
// selectivity level is exactly reproducible. It therefore carries an explicit
// `rank` column holding a deterministic shuffle of 1..N: the predicate is
// `rank <= threshold`, so a threshold count maps to an exact selectivity and a
// fixed index set can be held across every spelling pair. The harness labels
// that mapping rather than pretending the teaching `amount > n` predicate is
// the same question.
//
// The oracle is computed in Go from the same in-memory fixture the loader
// inserts, so a query and its truth can never see different rows.

// Fact is one child row of the join fixture. Rank is a permutation of 1..N so a
// threshold selects an exact count; CustomerID is a pointer so an orphan row is
// a real NULL, which is what the left-join and anti-join contracts turn on.
type Fact struct {
	ID         int
	CustomerID *int
	Rank       int
	Amount     *int
	Payload    string
}

// JoinFixture is the join/CTE dataset: parents from the teaching domain and a
// separate child fact table that carries the rank column.
type JoinFixture struct {
	Customers []Customer
	Facts     []Fact
}

// JoinGenSpec is a deterministic dataset shape. FactCount is the fact-table
// scale; the rank permutation covers exactly 1..FactCount.
type JoinGenSpec struct {
	Customers       int
	FactCount       int
	NullAmountEvery int // 0 = no NULL amounts
	NullFKEvery     int // 0 = no NULL foreign keys
	Wide            bool
}

// GenerateJoin builds the deterministic join fixture. Ranks are a seeded
// Fisher-Yates permutation of 1..N, so `rank <= t` selects exactly the rows
// whose permuted position is <= t and the selectivity is a rational number with
// no sampling error. Parent keys are a fixed modulo assignment, so the join
// fanout is deterministic and identical across every spelling.
func GenerateJoin(seed int64, spec JoinGenSpec) JoinFixture {
	jf := JoinFixture{}
	for i := 1; i <= spec.Customers; i++ {
		jf.Customers = append(jf.Customers, Customer{ID: i, Name: "c" + strconv.Itoa(i)})
	}
	n := spec.FactCount
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i + 1
	}
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(n, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
	jf.Facts = make([]Fact, 0, n)
	for i := 1; i <= n; i++ {
		cid := 1
		if spec.Customers > 0 {
			cid = ((i - 1) % spec.Customers) + 1
		}
		cidp := &cid
		if spec.NullFKEvery > 0 && i%spec.NullFKEvery == 0 {
			cidp = nil
		}
		var amt *int
		if spec.NullAmountEvery == 0 || i%spec.NullAmountEvery != 0 {
			v := 1 + (i*37)%500
			amt = ip(v)
		}
		f := Fact{ID: i, CustomerID: cidp, Rank: perm[i-1], Amount: amt}
		if spec.Wide {
			f.Payload = "p" + strconv.Itoa(i) + "0123456789abcdef0123456789abcdef"
		}
		jf.Facts = append(jf.Facts, f)
	}
	return jf
}

// JoinEdge is the small edge-case fixture the bag and optional-match controls
// are proved on: a parent with two qualifying children (duplicate matches), a
// NULL amount, an orphan fact with a NULL parent key, and an unmatched parent.
func JoinEdge() JoinFixture {
	return JoinFixture{
		Customers: []Customer{{1, "Ana"}, {2, "Ben"}, {3, "Cara"}, {4, "Ana"}},
		Facts: []Fact{
			{ID: 101, CustomerID: ip(1), Rank: 1, Amount: ip(50)},
			{ID: 102, CustomerID: ip(1), Rank: 2, Amount: ip(50)},
			{ID: 103, CustomerID: ip(2), Rank: 3, Amount: nil},
			{ID: 104, CustomerID: nil, Rank: 4, Amount: ip(150)},
			{ID: 105, CustomerID: ip(4), Rank: 5, Amount: ip(300)},
		},
	}
}

// --- canonical row shapes --------------------------------------------------

// JoinRow is one joined bag row: (customer, child). The projection order is
// fixed across every spelling so a difference is never an output-shape artefact.
type JoinRow struct {
	CustomerID   int
	CustomerName string
	FactID       int
	Amount       *int
}

// JoinOptRow is one optional-match row: every parent kept, with a qualifying
// child or NULLs.
type JoinOptRow struct {
	CustomerID   int
	CustomerName string
	FactID       *int
	Amount       *int
}

// ReuseRow is one CTE reference's answer in the reuse contract.
type ReuseRow struct {
	Ref      int
	Checksum int64
}

// --- oracles ---------------------------------------------------------------

func joinNames(jf JoinFixture) map[int]string {
	m := map[int]string{}
	for _, c := range jf.Customers {
		m[c.ID] = c.Name
	}
	return m
}

// JoinBag is the bag of joined rows whose child rank is <= threshold, sorted by
// (customer, fact). An orphan fact belongs to no parent and is absent; a NULL
// amount is still a joined row, so the bag must carry it.
func JoinBag(jf JoinFixture, threshold int) []JoinRow {
	names := joinNames(jf)
	var out []JoinRow
	for _, f := range jf.Facts {
		if f.Rank < 1 || f.Rank > threshold || f.CustomerID == nil {
			continue
		}
		n, ok := names[*f.CustomerID]
		if !ok {
			continue
		}
		out = append(out, JoinRow{*f.CustomerID, n, f.ID, f.Amount})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CustomerID != out[j].CustomerID {
			return out[i].CustomerID < out[j].CustomerID
		}
		return out[i].FactID < out[j].FactID
	})
	return out
}

// JoinOptional is the optional-match bag: every parent preserved, with its
// qualifying children or one NULL row. The predicate stays in ON; moving it to
// WHERE answers the inner bag instead, which is the deliberately-wrong control.
func JoinOptional(jf JoinFixture, threshold int) []JoinOptRow {
	names := joinNames(jf)
	var out []JoinOptRow
	for _, c := range jf.Customers {
		matched := false
		for _, f := range jf.Facts {
			if f.CustomerID == nil || *f.CustomerID != c.ID || f.Rank < 1 || f.Rank > threshold {
				continue
			}
			matched = true
			out = append(out, JoinOptRow{c.ID, names[c.ID], ip(f.ID), f.Amount})
		}
		if !matched {
			out = append(out, JoinOptRow{CustomerID: c.ID, CustomerName: names[c.ID]})
		}
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CustomerID != out[j].CustomerID {
			return out[i].CustomerID < out[j].CustomerID
		}
		if out[i].FactID == nil || out[j].FactID == nil {
			return out[i].FactID == nil && out[j].FactID != nil
		}
		return *out[i].FactID < *out[j].FactID
	})
	return out
}

// ControlJoinLimitBefore is the deliberately-wrong bag control: the limit is
// applied to the child side before the multiplicative join, so the joined
// result is truncated arbitrarily. It must be *detected*, never timed as a
// spelling.
func ControlJoinLimitBefore(jf JoinFixture, threshold, limit int) []JoinRow {
	if limit <= 0 {
		return nil
	}
	names := joinNames(jf)
	kids := make([]Fact, 0, len(jf.Facts))
	for _, f := range jf.Facts {
		if f.CustomerID != nil && f.Rank >= 1 && f.Rank <= threshold {
			kids = append(kids, f)
		}
	}
	sort.Slice(kids, func(i, j int) bool { return kids[i].ID < kids[j].ID })
	if len(kids) > limit {
		kids = kids[:limit]
	}
	var out []JoinRow
	for _, f := range kids {
		out = append(out, JoinRow{*f.CustomerID, names[*f.CustomerID], f.ID, f.Amount})
	}
	return out
}

// --- reuse contract --------------------------------------------------------

// ReuseExpr is the side-effect-free expensive token: md5 of the id text
// repeated 32 times consumes real per-row work without changing data, and its
// first hex character's ASCII code is a cheap, deterministic checksum. It must
// be verified as supported on each pinned engine before the reuse matrix runs.
const ReuseExpr = `ascii(substr(md5(repeat(f.id::text, 32)), 1, 1))`

// reuseChecksum is the Go mirror of ReuseExpr for one fact id.
func reuseChecksum(id int) int64 {
	sum := md5.Sum([]byte(strings.Repeat(strconv.Itoa(id), 32)))
	h := hex.EncodeToString(sum[:])
	if h == "" {
		return 0
	}
	return int64(h[0])
}

// ReuseAnswer is the answer of a CTE referenced refs times: each reference
// consumes the common CTE over the rows the per-consumer filter keeps (id <=
// consumerMax). Materialization is a policy, not a promise; the oracle only
// fixes the answer, which must be identical under every policy.
func ReuseAnswer(jf JoinFixture, threshold, refs, consumerMax int) []ReuseRow {
	var base int64
	for _, f := range jf.Facts {
		if f.Rank >= 1 && f.Rank <= threshold && f.ID <= consumerMax {
			base += reuseChecksum(f.ID)
		}
	}
	out := make([]ReuseRow, 0, refs)
	for r := 1; r <= refs; r++ {
		out = append(out, ReuseRow{Ref: r, Checksum: base})
	}
	return out
}
