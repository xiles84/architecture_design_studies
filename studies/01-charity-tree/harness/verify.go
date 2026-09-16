package main

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ---------------------------------------------------------------------------
// Correctness gate
//
// A design that answers the question wrongly, quickly, is worth nothing. Before
// any timing is recorded, every design must reproduce the same answers, computed
// independently in Go from the generated dataset rather than from another query.
//
// This is what makes the comparison a comparison. Denormalised keys, trigger
// rollups and embedded documents are all opportunities to drift out of sync, and
// a benchmark that never checks would happily report a broken design as the
// fastest one.
//
// Ties are handled explicitly: with a million timestamps truncated to the
// millisecond, two donations CAN share the maximum instant, and the SQL is free
// to return either. Identity checks therefore accept any member of the tied set.
// ---------------------------------------------------------------------------

type Check struct {
	Query  string `json:"query"`
	OK     bool   `json:"ok"`
	Expect string `json:"expect,omitempty"`
	Got    string `json:"got,omitempty"`
	Detail string `json:"detail,omitempty"`
}

type VerifyReport struct {
	Passed int     `json:"passed"`
	Failed int     `json:"failed"`
	Checks []Check `json:"checks"`
}

func (v *VerifyReport) add(c Check) {
	if c.OK {
		v.Passed++
	} else {
		v.Failed++
	}
	v.Checks = append(v.Checks, c)
}

// truth holds the answers derived directly from the generated dataset.
type truth struct {
	charityID int64
	personID  int64

	globalLastIDs  map[int64]bool // donations tied at the global max instant
	globalLastAt   time.Time
	charLastIDs    map[int64]bool
	charLastAt     time.Time
	charLastPeople map[int64]bool

	topDonors     map[int64]bool // people tied at the charity's max total
	topDonorTotal int64

	personFirstAt time.Time
	personLastAt  time.Time
	personCount   int64

	totalAll     int64
	totalCharity int64

	donationID     int64
	donationAmount int64
	donationPerson int64
}

// pickPerson chooses a person who actually has donations and belongs to the
// charity under test. Verifying against a donor with an empty history would pass
// trivially in every design and prove nothing.
func pickPerson(ds *Dataset, charityID int64) int64 {
	best, bestN := int64(0), 0
	for i := range ds.People {
		if ds.People[i].CharityID != charityID {
			continue
		}
		if n := len(ds.DonationsByPerson[i]); n > bestN {
			best, bestN = ds.People[i].ID, n
		}
	}
	return best
}

func computeTruth(ds *Dataset, charityID int64) truth {
	t := truth{
		charityID:      charityID,
		globalLastIDs:  map[int64]bool{},
		charLastIDs:    map[int64]bool{},
		charLastPeople: map[int64]bool{},
		topDonors:      map[int64]bool{},
	}
	t.personID = pickPerson(ds, charityID)

	perPersonTotal := map[int64]int64{}
	for i := range ds.Donations {
		d := &ds.Donations[i]
		t.totalAll += d.AmountCents

		if d.DonatedAt.After(t.globalLastAt) {
			t.globalLastAt = d.DonatedAt
			t.globalLastIDs = map[int64]bool{d.ID: true}
		} else if d.DonatedAt.Equal(t.globalLastAt) {
			t.globalLastIDs[d.ID] = true
		}

		if d.CharityID == charityID {
			t.totalCharity += d.AmountCents
			perPersonTotal[d.PersonID] += d.AmountCents
			if d.DonatedAt.After(t.charLastAt) {
				t.charLastAt = d.DonatedAt
				t.charLastIDs = map[int64]bool{d.ID: true}
				t.charLastPeople = map[int64]bool{d.PersonID: true}
			} else if d.DonatedAt.Equal(t.charLastAt) {
				t.charLastIDs[d.ID] = true
				t.charLastPeople[d.PersonID] = true
			}
		}

		if d.PersonID == t.personID {
			t.personCount++
			if t.personFirstAt.IsZero() || d.DonatedAt.Before(t.personFirstAt) {
				t.personFirstAt = d.DonatedAt
			}
			if d.DonatedAt.After(t.personLastAt) {
				t.personLastAt = d.DonatedAt
			}
		}
	}

	for pid, total := range perPersonTotal {
		if total > t.topDonorTotal {
			t.topDonorTotal = total
			t.topDonors = map[int64]bool{pid: true}
		} else if total == t.topDonorTotal {
			t.topDonors[pid] = true
		}
	}

	// A mid-range donation id, so the point-lookup check is not testing the
	// first physical row.
	idx := len(ds.Donations) / 2
	t.donationID = ds.Donations[idx].ID
	t.donationAmount = ds.Donations[idx].AmountCents
	t.donationPerson = ds.Donations[idx].PersonID
	return t
}

// Verify runs each answer-bearing query with fixed keys and compares the result
// against the independently computed truth.
func Verify(ctx context.Context, pool *pgxpool.Pool, d Design, ds *Dataset) (*VerifyReport, error) {
	src, err := readSQL(d.ID, "queries.sql")
	if err != nil {
		return nil, err
	}
	stmts, err := ParseCatalog(src)
	if err != nil {
		return nil, err
	}
	q := StmtMap(stmts)
	rep := &VerifyReport{}

	const charityID = int64(1) // the largest charity
	t := computeTruth(ds, charityID)

	vals := map[string]any{
		"charity_id":  charityID,
		"person_id":   t.personID,
		"donation_id": t.donationID,
	}

	// row fetches the first row of a query as a slice of generic values.
	row := func(name string) ([]any, error) {
		st, ok := q[name]
		if !ok {
			return nil, fmt.Errorf("query %s missing from catalogue", name)
		}
		args, err := bind(st.Params, vals)
		if err != nil {
			return nil, err
		}
		rows, err := pool.Query(ctx, st.SQL, args...)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		if !rows.Next() {
			if err := rows.Err(); err != nil {
				return nil, err
			}
			return nil, nil
		}
		return rows.Values()
	}

	countRows := func(name string) (int, []any, error) {
		st, ok := q[name]
		if !ok {
			return 0, nil, fmt.Errorf("query %s missing from catalogue", name)
		}
		args, err := bind(st.Params, vals)
		if err != nil {
			return 0, nil, err
		}
		rows, err := pool.Query(ctx, st.SQL, args...)
		if err != nil {
			return 0, nil, err
		}
		defer rows.Close()
		n := 0
		var first []any
		for rows.Next() {
			if n == 0 {
				first, _ = rows.Values()
			}
			n++
		}
		return n, first, rows.Err()
	}

	check := func(name string, fn func() (bool, string, string, error)) {
		ok, exp, got, err := fn()
		c := Check{Query: name, OK: ok, Expect: exp, Got: got}
		if err != nil {
			c.OK = false
			c.Detail = err.Error()
		}
		rep.add(c)
	}

	// q01 -- information of the last donation (global)
	check("q01_last_donation_global", func() (bool, string, string, error) {
		v, err := row("q01_last_donation_global")
		if err != nil || v == nil {
			return false, "", "", err
		}
		id := toI64(v[0])
		at := toTime(v[4])
		return t.globalLastIDs[id] && at.Equal(t.globalLastAt),
			fmt.Sprintf("one of %v at %s", keys(t.globalLastIDs), t.globalLastAt.UTC()),
			fmt.Sprintf("%d at %s", id, at.UTC()), nil
	})

	// q02 -- last donation of a charity
	check("q02_last_donation_charity", func() (bool, string, string, error) {
		v, err := row("q02_last_donation_charity")
		if err != nil || v == nil {
			return false, "", "", err
		}
		id, at := toI64(v[0]), toTime(v[4])
		return t.charLastIDs[id] && at.Equal(t.charLastAt),
			fmt.Sprintf("one of %v at %s", keys(t.charLastIDs), t.charLastAt.UTC()),
			fmt.Sprintf("%d at %s", id, at.UTC()), nil
	})

	// q03 -- who donates the most
	check("q03_top_donor_charity", func() (bool, string, string, error) {
		v, err := row("q03_top_donor_charity")
		if err != nil || v == nil {
			return false, "", "", err
		}
		pid, total := toI64(v[0]), toI64(v[2])
		return t.topDonors[pid] && total == t.topDonorTotal,
			fmt.Sprintf("one of %v with %d cents", keys(t.topDonors), t.topDonorTotal),
			fmt.Sprintf("%d with %d cents", pid, total), nil
	})

	// q04 -- leaderboard is ordered and complete
	check("q04_top_donors_leaderboard", func() (bool, string, string, error) {
		n, first, err := countRows("q04_top_donors_leaderboard")
		if err != nil {
			return false, "", "", err
		}
		if first == nil {
			return false, "10 rows", "0 rows", nil
		}
		return n == 10 && toI64(first[2]) == t.topDonorTotal,
			fmt.Sprintf("10 rows, leader at %d cents", t.topDonorTotal),
			fmt.Sprintf("%d rows, leader at %d cents", n, toI64(first[2])), nil
	})

	// q05 -- who donated last
	check("q05_last_donor_charity", func() (bool, string, string, error) {
		v, err := row("q05_last_donor_charity")
		if err != nil || v == nil {
			return false, "", "", err
		}
		pid, at := toI64(v[0]), toTime(v[2])
		return t.charLastPeople[pid] && at.Equal(t.charLastAt),
			fmt.Sprintf("one of %v at %s", keys(t.charLastPeople), t.charLastAt.UTC()),
			fmt.Sprintf("%d at %s", pid, at.UTC()), nil
	})

	// q06 -- first and last donation of a person
	check("q06_person_first_last", func() (bool, string, string, error) {
		v, err := row("q06_person_first_last")
		if err != nil || v == nil {
			return false, "", "", err
		}
		f, l := toTime(v[0]), toTime(v[1])
		return f.Equal(t.personFirstAt) && l.Equal(t.personLastAt),
			fmt.Sprintf("%s .. %s", t.personFirstAt.UTC(), t.personLastAt.UTC()),
			fmt.Sprintf("%s .. %s", f.UTC(), l.UTC()), nil
	})

	// q07 -- total donated overall
	check("q07_total_donated_global", func() (bool, string, string, error) {
		v, err := row("q07_total_donated_global")
		if err != nil || v == nil {
			return false, "", "", err
		}
		got := toI64(v[0])
		return got == t.totalAll, fmt.Sprint(t.totalAll), fmt.Sprint(got), nil
	})

	// q08 -- total donated to a charity
	check("q08_total_donated_charity", func() (bool, string, string, error) {
		v, err := row("q08_total_donated_charity")
		if err != nil || v == nil {
			return false, "", "", err
		}
		got := toI64(v[0])
		return got == t.totalCharity, fmt.Sprint(t.totalCharity), fmt.Sprint(got), nil
	})

	// q09 -- the person's page is newest-first and the right length
	check("q09_person_recent_donations", func() (bool, string, string, error) {
		n, first, err := countRows("q09_person_recent_donations")
		if err != nil {
			return false, "", "", err
		}
		want := int(t.personCount)
		if want > 20 {
			want = 20
		}
		if first == nil {
			return false, fmt.Sprintf("%d rows", want), "0 rows", nil
		}
		at := toTime(first[3])
		return n == want && at.Equal(t.personLastAt),
			fmt.Sprintf("%d rows, newest %s", want, t.personLastAt.UTC()),
			fmt.Sprintf("%d rows, newest %s", n, at.UTC()), nil
	})

	// q10 -- point lookup returns the right donation
	check("q10_donation_by_id", func() (bool, string, string, error) {
		v, err := row("q10_donation_by_id")
		if err != nil || v == nil {
			return false, "", "", err
		}
		return toI64(v[0]) == t.donationID && toI64(v[1]) == t.donationPerson && toI64(v[2]) == t.donationAmount,
			fmt.Sprintf("id=%d person=%d amount=%d", t.donationID, t.donationPerson, t.donationAmount),
			fmt.Sprintf("id=%v person=%v amount=%v", toI64(v[0]), toI64(v[1]), toI64(v[2])), nil
	})

	// q11 -- donation count of a person
	check("q11_person_donation_count", func() (bool, string, string, error) {
		v, err := row("q11_person_donation_count")
		if err != nil || v == nil {
			return false, "", "", err
		}
		got := toI64(v[0])
		return got == t.personCount, fmt.Sprint(t.personCount), fmt.Sprint(got), nil
	})

	// q12 -- charity feed is newest-first and full
	check("q12_charity_recent_feed", func() (bool, string, string, error) {
		n, first, err := countRows("q12_charity_recent_feed")
		if err != nil {
			return false, "", "", err
		}
		if first == nil {
			return false, "50 rows", "0 rows", nil
		}
		at := toTime(first[2])
		return n == 50 && at.Equal(t.charLastAt),
			fmt.Sprintf("50 rows, newest %s", t.charLastAt.UTC()),
			fmt.Sprintf("%d rows, newest %s", n, at.UTC()), nil
	})

	// -----------------------------------------------------------------------
	// q13-q16 -- "who made their LAST donation in this period" (RECENCY.md).
	//
	// A different shape from everything above: the answer is a SET defined by
	// a per-donor aggregate over the DONOR'S WHOLE HISTORY, not a row selected
	// by a key. Truth is computed once (independent of any window) as every
	// donor's own MAX(donated_at); each regime then filters that already-
	// computed aggregate, exactly mirroring what a correct design must do.
	// -----------------------------------------------------------------------

	lastByPerson := map[int64]time.Time{}
	personCharityID := map[int64]int64{}
	for i := range ds.People {
		personCharityID[ds.People[i].ID] = ds.People[i].CharityID
	}
	for i := range ds.Donations {
		d := &ds.Donations[i]
		if cur, ok := lastByPerson[d.PersonID]; !ok || d.DonatedAt.After(cur) {
			lastByPerson[d.PersonID] = d.DonatedAt
		}
	}

	for _, regime := range []string{"trailing", "historical"} {
		since, until := recencyWindowFor(epochEnd, regime)
		suffix := "@" + regime

		// q13/q15 -- ranked top-100, with the boundary-tie rule the rest of
		// this gate already uses: any donor sharing the exact cutoff instant
		// is an acceptable fill for the remaining slots.
		checkRanked := func(name string, scopeCharity int64) {
			check(name+suffix, func() (bool, string, string, error) {
				type row struct {
					id int64
					at time.Time
				}

				st, ok := q[name]
				if !ok {
					return false, "", "", fmt.Errorf("query %s missing from catalogue", name)
				}
				argVals := map[string]any{"charity_id": scopeCharity, "since": since, "until": until}
				args, err := bind(st.Params, argVals)
				if err != nil {
					return false, "", "", err
				}
				rows, err := pool.Query(ctx, st.SQL, args...)
				if err != nil {
					return false, "", "", err
				}
				defer rows.Close()
				var got []row
				for rows.Next() {
					v, err := rows.Values()
					if err != nil {
						return false, "", "", err
					}
					got = append(got, row{id: toI64(v[0]), at: toTime(v[2])})
				}
				if err := rows.Err(); err != nil {
					return false, "", "", err
				}
				n := len(got)

				expectedAll := recencyExpectedMatches(lastByPerson, personCharityID, since, until, scopeCharity)
				var expected []row
				for _, e := range expectedAll {
					expected = append(expected, row{id: e.PersonID, at: e.LastAt})
				}
				wantN := len(expected)
				if wantN > 100 {
					wantN = 100
				}

				if n != wantN {
					return false, fmt.Sprintf("%d rows", wantN), fmt.Sprintf("%d rows", n), nil
				}

				// Order: non-increasing last_at.
				for i := 1; i < len(got); i++ {
					if got[i].at.After(got[i-1].at) {
						return false, "non-increasing last_at", "out of order", nil
					}
				}

				// Every returned donor must genuinely match the predicate, and
				// the value returned must equal that donor's true last gift.
				for _, r := range got {
					want, ok := lastByPerson[r.id]
					if !ok || !want.Equal(r.at) || !recencyMatches(want, personCharityID[r.id], since, until, scopeCharity) {
						return false, "", fmt.Sprintf("donor %d at %s", r.id, r.at.UTC()),
							fmt.Errorf("returned donor does not match the window predicate")
					}
				}

				if wantN == 0 {
					return true, "0 rows", fmt.Sprintf("%d rows", n), nil
				}

				// Boundary tie handling: every expected donor STRICTLY newer
				// than the cutoff (the wantN-th expected value) must appear;
				// donors AT the cutoff instant are an acceptable substitute
				// for each other, as q04/q12's boundary already allows.
				boundary := expected[wantN-1].at
				gotSet := map[int64]bool{}
				for _, r := range got {
					gotSet[r.id] = true
				}
				for _, e := range expected[:wantN] {
					if e.at.After(boundary) && !gotSet[e.id] {
						return false, fmt.Sprintf("donor %d (above the tie boundary)", e.id), "missing", nil
					}
				}
				return true,
					fmt.Sprintf("%d rows, newest %s", wantN, expected[0].at.UTC()),
					fmt.Sprintf("%d rows, newest %s", n, got[0].at.UTC()), nil
			})
		}

		checkCount := func(name string, want int) {
			check(name+suffix, func() (bool, string, string, error) {
				st, ok := q[name]
				if !ok {
					return false, "", "", fmt.Errorf("query %s missing from catalogue", name)
				}
				argVals := map[string]any{"charity_id": charityID, "since": since, "until": until}
				args, err := bind(st.Params, argVals)
				if err != nil {
					return false, "", "", err
				}
				var got int64
				if err := pool.QueryRow(ctx, st.SQL, args...).Scan(&got); err != nil {
					return false, "", "", err
				}
				return got == int64(want), fmt.Sprint(want), fmt.Sprint(got), nil
			})
		}

		checkRanked("q13_donors_last_gift_window", 0)
		checkCount("q14_donors_last_gift_window_count", len(recencyExpectedMatches(lastByPerson, personCharityID, since, until, 0)))

		checkRanked("q15_charity_donors_last_gift_window", charityID)
		checkCount("q16_lapsed_donors_count", recencyExpectedLapsedCount(lastByPerson, since))
	}

	return rep, nil
}

// ---------------------------------------------------------------------------
// Recency truth (RECENCY.md section 2/5) -- pure functions, no DB, unit
// tested directly against a hand-built lastByPerson map in
// experiment_test.go. Extracted out of Verify's closures so the window
// predicate and the expected-set computation have exactly one definition,
// used both to check the gate and to build the negative-control examples.
// ---------------------------------------------------------------------------

// recencyMatches is the window-set predicate every correct design must
// implement, whichever table or column it reads it from: is this donor's
// last gift inside [since, until), and (if scopeCharity != 0) do they belong
// to that charity?
func recencyMatches(lastAt time.Time, donorCharity int64, since, until time.Time, scopeCharity int64) bool {
	if scopeCharity != 0 && donorCharity != scopeCharity {
		return false
	}
	return !lastAt.Before(since) && lastAt.Before(until)
}

type recencyDonorLast struct {
	PersonID int64
	LastAt   time.Time
}

// recencyExpectedMatches returns every donor matching the window predicate,
// sorted newest-first -- the full expected set for q13/q15 (the caller caps
// it at the LIMIT) and the exact count for q14.
func recencyExpectedMatches(lastByPerson map[int64]time.Time, personCharityID map[int64]int64, since, until time.Time, scopeCharity int64) []recencyDonorLast {
	var out []recencyDonorLast
	for pid, at := range lastByPerson {
		if recencyMatches(at, personCharityID[pid], since, until, scopeCharity) {
			out = append(out, recencyDonorLast{PersonID: pid, LastAt: at})
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LastAt.After(out[j].LastAt) })
	return out
}

// recencyExpectedLapsedCount is q16's truth: donors whose last gift is
// strictly before since.
func recencyExpectedLapsedCount(lastByPerson map[int64]time.Time, since time.Time) int {
	n := 0
	for _, at := range lastByPerson {
		if at.Before(since) {
			n++
		}
	}
	return n
}

func keys(m map[int64]bool) []int64 {
	out := make([]int64, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	if len(out) > 5 {
		return out[:5]
	}
	return out
}

// toI64 normalises the several numeric shapes pgx can hand back for what is
// logically a count of cents.
//
// The one that matters: SUM(bigint) is NUMERIC in PostgreSQL, not BIGINT --
// because the sum of bigints can overflow a bigint. pgx therefore decodes it as
// pgtype.Numeric rather than int64, and a converter that only knows about int64
// silently yields zero. That is precisely how a verification pass can report a
// design as broken when only the type handling is.
func toI64(v any) int64 {
	if n, ok := v.(interface {
		Int64Value() (pgtype.Int8, error)
	}); ok {
		if i, err := n.Int64Value(); err == nil && i.Valid {
			return i.Int64
		}
	}
	switch x := v.(type) {
	case int64:
		return x
	case int32:
		return int64(x)
	case int:
		return int64(x)
	case float64:
		return int64(x)
	case string:
		var n int64
		fmt.Sscan(x, &n)
		return n
	case nil:
		return 0
	default:
		var n int64
		fmt.Sscan(fmt.Sprint(x), &n)
		return n
	}
}

func toTime(v any) time.Time {
	if t, ok := v.(time.Time); ok {
		return t.UTC()
	}
	return time.Time{}
}
