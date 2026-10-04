package selectbench

import (
	"math/rand"
	"sort"
	"strconv"
)

// Customer and Order are the teaching analogue of Study 01's charity/person/
// donation rows. Amount and CustomerID are pointers so a NULL is a real NULL,
// not a sentinel zero: the NULL cases are the point of several controls.
type Customer struct {
	ID   int
	Name string
}

type Order struct {
	ID         int
	CustomerID *int
	Amount     *int
	Payload    string // non-empty in the wide-row variant
}

type Fixture struct {
	Customers []Customer
	Orders    []Order
}

func ip(v int) *int { return &v }

// Teaching is the edge-case fixture. It deliberately contains every case the
// protocol requires: a NULL amount (order 13), two qualifying children for one
// parent (orders 11 and 14 for customer 1) making a tie at 200, an unmatched
// parent (customer 3), an all-miss set (amount > 1000), and a NULL child foreign
// key (order 15) that poisons a nullable NOT IN. Customer 4 shares the name
// "Ana" with customer 1 so a name-based DISTINCT silently merges two identities.
func Teaching() Fixture {
	return Fixture{
		Customers: []Customer{{1, "Ana"}, {2, "Ben"}, {3, "Cara"}, {4, "Ana"}},
		Orders: []Order{
			{ID: 10, CustomerID: ip(1), Amount: ip(50)},
			{ID: 11, CustomerID: ip(1), Amount: ip(200)},
			{ID: 12, CustomerID: ip(2), Amount: ip(80)},
			{ID: 13, CustomerID: ip(2), Amount: nil}, // NULL amount
			{ID: 14, CustomerID: ip(1), Amount: ip(200)},
			{ID: 15, CustomerID: nil, Amount: ip(150)}, // NULL foreign key
			{ID: 16, CustomerID: ip(4), Amount: ip(300)},
		},
	}
}

// GenSpec is a deterministic dataset shape. Narrow/wide rows, fanout and skew
// are inputs, never results.
type GenSpec struct {
	Customers        int
	OrdersPerParent  int
	NullAmountEvery  int // 0 = no NULL amounts
	NullFKEvery      int // 0 = no NULL foreign keys
	HotCustomerShare float64
	Wide             bool
}

// Generate builds a deterministic fixture from a seed. It is the harness's only
// dataset source for the SELECT families, so an oracle and a query always see
// the same logical rows.
func Generate(seed int64, spec GenSpec) Fixture {
	r := rand.New(rand.NewSource(seed))
	fx := Fixture{}
	for i := 1; i <= spec.Customers; i++ {
		fx.Customers = append(fx.Customers, Customer{ID: i, Name: "c" + strconv.Itoa(i)})
	}
	nextID := 1000
	hot := 1
	for p := 1; p <= spec.Customers; p++ {
		n := spec.OrdersPerParent
		if spec.HotCustomerShare > 0 && p == hot {
			n = int(float64(n) * (1 + spec.HotCustomerShare*float64(spec.Customers)))
		}
		for k := 0; k < n; k++ {
			nextID++
			cid := p
			var amt *int
			if spec.NullAmountEvery == 0 || nextID%spec.NullAmountEvery != 0 {
				v := 1 + r.Intn(500)
				amt = ip(v)
			}
			cidp := &cid
			if spec.NullFKEvery > 0 && nextID%spec.NullFKEvery == 0 {
				cidp = nil
			}
			o := Order{ID: nextID, CustomerID: cidp, Amount: amt}
			if spec.Wide {
				o.Payload = "p" + strconv.Itoa(nextID) + "0123456789abcdef0123456789abcdef"
			}
			fx.Orders = append(fx.Orders, o)
		}
	}
	return fx
}

// Pred is the predicate under test.
type Pred func(Order) bool

// AmountGreaterThan is the fixture's canonical predicate. NULL never satisfies
// it, which is what makes the optional-match and absence contracts interesting.
func AmountGreaterThan(n int) Pred {
	return func(o Order) bool { return o.Amount != nil && *o.Amount > n }
}

// --- canonical row shapes -------------------------------------------------

type FullRow struct {
	CustomerID   int
	CustomerName string
	OrderID      int
	Amount       *int
}

type OptRow struct {
	CustomerID   int
	CustomerName string
	OrderID      *int
	Amount       *int
}

type AggRow struct {
	CustomerID   int
	CustomerName string
	Total        int // SUM over non-NULL amounts; 0 for a parent with no rows
	Children     int // COUNT(*)
}

type TopRow struct {
	CustomerID int
	OrderID    int
	Amount     *int
}

type PageRow struct {
	OrderID int
	Amount  *int
}

// --- oracles --------------------------------------------------------------

func byName(fx Fixture) map[int]string {
	m := map[int]string{}
	for _, c := range fx.Customers {
		m[c.ID] = c.Name
	}
	return m
}

// FullRows is the bag of joined rows satisfying pred, sorted by (customer, order).
func FullRows(fx Fixture, pred Pred) []FullRow {
	names := byName(fx)
	var out []FullRow
	for _, o := range fx.Orders {
		if o.CustomerID == nil || !pred(o) {
			continue
		}
		n, ok := names[*o.CustomerID]
		if !ok {
			continue // orphan row: no parent, so it is in no joined bag
		}
		out = append(out, FullRow{*o.CustomerID, n, o.ID, o.Amount})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].CustomerID != out[j].CustomerID {
			return out[i].CustomerID < out[j].CustomerID
		}
		return out[i].OrderID < out[j].OrderID
	})
	return out
}

// Membership is the set of parent identities with at least one qualifying child.
func Membership(fx Fixture, pred Pred) []int {
	seen := map[int]bool{}
	for _, o := range fx.Orders {
		if o.CustomerID != nil && pred(o) {
			seen[*o.CustomerID] = true
		}
	}
	var out []int
	for id := range seen {
		out = append(out, id)
	}
	sort.Ints(out)
	return out
}

// MembershipWithMultiplicity is the deliberately-wrong membership grain: one
// entry per qualifying child row, so a parent with several qualifying children
// appears several times. It mirrors the x07 control (a plain join without
// DISTINCT) and must differ from Membership wherever a parent has >1 match.
func MembershipWithMultiplicity(fx Fixture, pred Pred) []int {
	var out []int
	for _, o := range fx.Orders {
		if o.CustomerID != nil && pred(o) {
			out = append(out, *o.CustomerID)
		}
	}
	sort.Ints(out)
	return out
}

// ExistenceShape reshapes a membership/absence fixture so the two controls can
// fire at every fanout level, without changing the total child count: the last
// parent's children are moved to the first parent (so one parent always has
// several children, making the x07 plain-join multiplier differ from the set),
// and the first parent's first two orders are given amount 400 (so they qualify
// under any positive threshold). The last parent is therefore childless, which
// keeps the absence set non-empty for the x02 nullable-NOT IN control. It is a
// fixture shape for the controls, reflected in the dataset fingerprint.
func ExistenceShape(fx Fixture) Fixture {
	if len(fx.Customers) < 2 || len(fx.Orders) == 0 {
		return fx
	}
	hot := fx.Customers[0].ID
	empty := fx.Customers[len(fx.Customers)-1].ID
	for i := range fx.Orders {
		if fx.Orders[i].CustomerID != nil && *fx.Orders[i].CustomerID == empty {
			fx.Orders[i].CustomerID = ip(hot)
		}
	}
	set := 0
	for i := range fx.Orders {
		if fx.Orders[i].CustomerID != nil && *fx.Orders[i].CustomerID == hot {
			fx.Orders[i].Amount = ip(400)
			set++
			if set == 2 {
				break
			}
		}
	}
	return fx
}

// Absence is the set of parent identities with no qualifying child.
func Absence(fx Fixture, pred Pred) []int {
	has := map[int]bool{}
	for _, id := range Membership(fx, pred) {
		has[id] = true
	}
	var out []int
	for _, c := range fx.Customers {
		if !has[c.ID] {
			out = append(out, c.ID)
		}
	}
	sort.Ints(out)
	return out
}

// --- predicate-alternative contract (P10) ---------------------------------
//
// A disjunction over two predicate arms. The contract is a BAG of child orders:
// an order that satisfies both arms still appears exactly once, because OR has
// no multiplicity. `UNION` restores that bag even when the arms overlap;
// `UNION ALL` is the same bag only while the arms are provably disjoint. The
// deliberately-wrong control is `UNION ALL` over overlapping arms, which returns
// a duplicated bag and must be detected rather than timed.

// PredAltRow is one order in the predicate-alternative bag.
type PredAltRow struct {
	OrderID int
	Amount  *int
}

// PredicateAlt is the OR bag over two predicate arms, ordered by order id.
// A nil arm contributes nothing, so the zero Args value yields the empty bag.
func PredicateAlt(fx Fixture, predA, predB Pred) []PredAltRow {
	var out []PredAltRow
	for _, o := range fx.Orders {
		if (predA != nil && predA(o)) || (predB != nil && predB(o)) {
			out = append(out, PredAltRow{o.ID, o.Amount})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OrderID < out[j].OrderID })
	return out
}

// PredicateAltOverlap mirrors the deliberately-wrong x05: `UNION ALL` of one arm
// with itself, so every matching order appears twice.
func PredicateAltOverlap(fx Fixture, pred Pred) []PredAltRow {
	var out []PredAltRow
	for pass := 0; pass < 2; pass++ {
		for _, o := range fx.Orders {
			if pred != nil && pred(o) {
				out = append(out, PredAltRow{o.ID, o.Amount})
			}
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OrderID < out[j].OrderID })
	return out
}

// OverlapCount counts orders satisfying BOTH arms. Zero is the disjointness
// proof required before UNION ALL may be compared with OR as the same bag.
func OverlapCount(fx Fixture, predA, predB Pred) int {
	if predA == nil || predB == nil {
		return 0
	}
	n := 0
	for _, o := range fx.Orders {
		if predA(o) && predB(o) {
			n++
		}
	}
	return n
}

// AmountLE / AmountGE are the disjunction's arms. NULL never satisfies either,
// matching SQL's three-valued logic.
func AmountLE(n int) Pred { return func(o Order) bool { return o.Amount != nil && *o.Amount <= n } }
func AmountGE(n int) Pred { return func(o Order) bool { return o.Amount != nil && *o.Amount >= n } }

// OptionalMatch is the bag of parents with qualifying children or NULLs, sorted
// by (customer, order-with-NULL-last).
func OptionalMatch(fx Fixture, pred Pred) []OptRow {
	names := byName(fx)
	var out []OptRow
	for _, c := range fx.Customers {
		matched := false
		for _, o := range fx.Orders {
			if o.CustomerID == nil || *o.CustomerID != c.ID || !pred(o) {
				continue
			}
			matched = true
			out = append(out, OptRow{c.ID, names[c.ID], ip(o.ID), o.Amount})
		}
		if !matched {
			out = append(out, OptRow{c.ID, names[c.ID], nil, nil})
		}
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
	return out
}

// Aggregate is one row per parent: SUM of non-NULL amounts and COUNT(*).
func Aggregate(fx Fixture) []AggRow {
	names := byName(fx)
	type acc struct{ total, n int }
	m := map[int]*acc{}
	for _, o := range fx.Orders {
		if o.CustomerID == nil {
			continue
		}
		a := m[*o.CustomerID]
		if a == nil {
			a = &acc{}
			m[*o.CustomerID] = a
		}
		a.n++
		if o.Amount != nil {
			a.total += *o.Amount
		}
	}
	var out []AggRow
	for _, c := range fx.Customers {
		a := m[c.ID]
		if a == nil {
			out = append(out, AggRow{c.ID, names[c.ID], 0, 0})
			continue
		}
		out = append(out, AggRow{c.ID, names[c.ID], a.total, a.n})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CustomerID < out[j].CustomerID })
	return out
}

// orderLess is the declared total order: amount DESC NULLS LAST, then id ASC.
func orderLess(ai *int, aid int, bi *int, bid int) bool {
	an, bn := ai == nil, bi == nil
	switch {
	case an && bn:
		return aid < bid
	case an:
		return false // NULL sorts last
	case bn:
		return true
	case *ai != *bi:
		return *ai > *bi
	default:
		return aid < bid
	}
}

// TopN is at most n rows per parent under the declared total order. A parent
// with no orders contributes nothing; a parent whose only orders have NULL
// amounts still contributes, with a NULL amount, because NULL sorts last.
func TopN(fx Fixture, n int) []TopRow {
	byCust := map[int][]Order{}
	for _, o := range fx.Orders {
		if o.CustomerID == nil {
			continue
		}
		byCust[*o.CustomerID] = append(byCust[*o.CustomerID], o)
	}
	var out []TopRow
	for _, c := range fx.Customers {
		rows := byCust[c.ID]
		sort.Slice(rows, func(i, j int) bool { return orderLess(rows[i].Amount, rows[i].ID, rows[j].Amount, rows[j].ID) })
		if len(rows) > n {
			rows = rows[:n]
		}
		for _, o := range rows {
			out = append(out, TopRow{c.ID, o.ID, o.Amount})
		}
	}
	return out
}

// Page is the ordered page of all orders (including NULL amounts) under the
// declared total order.
func Page(fx Fixture, limit, offset int) []PageRow {
	rows := append([]Order(nil), fx.Orders...)
	sort.Slice(rows, func(i, j int) bool { return orderLess(rows[i].Amount, rows[i].ID, rows[j].Amount, rows[j].ID) })
	if offset > len(rows) {
		offset = len(rows)
	}
	rows = rows[offset:]
	if limit < len(rows) {
		rows = rows[:limit]
	}
	var out []PageRow
	for _, o := range rows {
		out = append(out, PageRow{o.ID, o.Amount})
	}
	return out
}

// --- aggregate / top-N sweep oracles (task select-aggregate-topn) ----------
//
// The sweep restricts the answer to a prefix of parent ids (`maxParent`, 0 =
// all), so one deterministic fixture serves parent-count levels 1..10,000
// without reloading. The unsuffixed Aggregate/TopN oracles remain the
// harness-validation contracts; these are the sweep contracts.

func parentSelected(id, maxParent int) bool {
	return maxParent <= 0 || id <= maxParent
}

// selectParents returns a fixture restricted to parents with id <= maxParent
// and the orders that belong to them. Orphan orders (NULL parent) are excluded:
// they belong to no parent, so no parent aggregate can see them.
func selectParents(fx Fixture, maxParent int) Fixture {
	sub := Fixture{}
	for _, c := range fx.Customers {
		if parentSelected(c.ID, maxParent) {
			sub.Customers = append(sub.Customers, c)
		}
	}
	for _, o := range fx.Orders {
		if o.CustomerID != nil && parentSelected(*o.CustomerID, maxParent) {
			sub.Orders = append(sub.Orders, o)
		}
	}
	return sub
}

// AggregateUpTo is Aggregate restricted to a prefix of parent ids.
func AggregateUpTo(fx Fixture, maxParent int) []AggRow {
	return Aggregate(selectParents(fx, maxParent))
}

// TopNUpTo is TopN restricted to a prefix of parent ids.
func TopNUpTo(fx Fixture, n, maxParent int) []TopRow {
	return TopN(selectParents(fx, maxParent), n)
}

// DoubleCountAggregate is the deliberately-wrong pre-aggregation grain: it
// answers as if orders had been joined to a second one-to-many collection, so
// every parent's SUM and COUNT are multiplied by its child count. It is a
// control: it must differ from the aggregate oracle wherever a parent has more
// than one child, and must never be shown as a faster aggregate.
func DoubleCountAggregate(fx Fixture, maxParent int) []AggRow {
	base := AggregateUpTo(fx, maxParent)
	counts := map[int]int{}
	for _, o := range fx.Orders {
		if o.CustomerID != nil && parentSelected(*o.CustomerID, maxParent) {
			counts[*o.CustomerID]++
		}
	}
	out := make([]AggRow, 0, len(base))
	for _, r := range base {
		n := counts[r.CustomerID]
		out = append(out, AggRow{CustomerID: r.CustomerID, CustomerName: r.CustomerName,
			Total: r.Total * n, Children: r.Children * n})
	}
	return out
}

// RankTiesTopN is the deliberately-wrong top-N tie behaviour: RANK keeps every
// tied row (a separate contract), so a parent whose tied group starts at rank 1
// returns more than N rows. It must be detected wherever the fixture holds a
// tie, which GuaranteedTie ensures.
func RankTiesTopN(fx Fixture, n, maxParent int) []TopRow {
	byCust := map[int][]Order{}
	for _, o := range fx.Orders {
		if o.CustomerID == nil || !parentSelected(*o.CustomerID, maxParent) {
			continue
		}
		byCust[*o.CustomerID] = append(byCust[*o.CustomerID], o)
	}
	var out []TopRow
	for _, c := range fx.Customers {
		if !parentSelected(c.ID, maxParent) {
			continue
		}
		rows := byCust[c.ID]
		sort.Slice(rows, func(i, j int) bool {
			return orderLess(rows[i].Amount, rows[i].ID, rows[j].Amount, rows[j].ID)
		})
		rank := 0
		havePrev := false
		var prev *int
		for _, o := range rows {
			if !havePrev || !sameAmount(prev, o.Amount) {
				rank++
			}
			prev, havePrev = o.Amount, true
			if rank <= n {
				out = append(out, TopRow{CustomerID: c.ID, OrderID: o.ID, Amount: o.Amount})
			}
		}
	}
	return out
}

func sameAmount(a, b *int) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

// GuaranteedTie makes the smallest parent with at least two children hold two
// equal amounts, so the tie control (x06) always has a tie to violate. It is a
// fixture shape for the control, not a workload change, and it is reflected in
// the dataset fingerprint.
func GuaranteedTie(fx Fixture) Fixture {
	byCust := map[int][]int{}
	for i, o := range fx.Orders {
		if o.CustomerID != nil {
			byCust[*o.CustomerID] = append(byCust[*o.CustomerID], i)
		}
	}
	best := 0
	for id, idx := range byCust {
		if len(idx) >= 2 && (best == 0 || id < best) {
			best = id
		}
	}
	if best == 0 {
		return fx
	}
	idx := byCust[best]
	// Tie at the TOP of the parent's order: give the first two orders the
	// parent's maximum amount. A tie anywhere lower would not make RANK keep an
	// extra row at rank 1, so the control would not fire.
	max := 0
	for _, i := range idx {
		if a := fx.Orders[i].Amount; a != nil && *a > max {
			max = *a
		}
	}
	if max == 0 {
		max = 100
	}
	fx.Orders[idx[0]].Amount = ip(max)
	fx.Orders[idx[1]].Amount = ip(max)
	return fx
}
