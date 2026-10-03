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
