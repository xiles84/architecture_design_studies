package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// ---------------------------------------------------------------------------
// Deterministic dataset.
//
// Every design and every engine loads the same logical dataset from a fixed
// seed: the same venues, the same events, the same seats already sold to the
// same customers, the same live and expired holds on the same seats. Designs
// differ only in how that logical state is laid out in tables (seat rows, a cart
// per hold, claim rows, section documents).
//
// Events come in three kinds:
//
//	catalogue  partly sold, with live and expired holds: reads, isolated writes,
//	           and the correctness gate's hardest cases
//	race       unsold; each is sold out by a crowd choosing seats
//	lifecycle  unsold; the hold guarantee under compressed time
//
// Seat quality is ranked with no randomness -- section 1 first, row 1 first,
// centre of the row first -- so buyers want the seats people really want and
// conflicts concentrate where they would in reality.
// ---------------------------------------------------------------------------

var allTiers = []int{10, 100, 1000, 10000, 100000}

type venueShape struct {
	sections, rows, perRow int
	name                   string
}

// One venue per tier. Sections are sized so a seat map is one section of at most
// 1 000 seats, as a real seat map is.
var venueShapes = map[int]venueShape{
	10:     {1, 1, 10, "back room"},
	100:    {1, 5, 20, "club"},
	1000:   {4, 10, 25, "theatre"},
	10000:  {20, 20, 25, "arena"},
	100000: {100, 40, 25, "stadium"},
}

type Seat struct {
	ID      int32
	Section int32
	Row     int32
	No      int32
	Rank    int32
}

type Section struct {
	No        int32
	Rows      int
	PerRow    int
	Capacity  int
	FirstSeat int32
}

type Venue struct {
	ID       int64
	Tier     int
	Name     string
	Capacity int
	Sections []Section
	// Seats is indexed by seat_id - 1.
	Seats []Seat
}

func buildVenue(id int64, tier int) *Venue {
	sh := venueShapes[tier]
	v := &Venue{ID: id, Tier: tier, Name: sh.name, Capacity: sh.sections * sh.rows * sh.perRow}
	next := int32(1)
	for s := 1; s <= sh.sections; s++ {
		v.Sections = append(v.Sections, Section{No: int32(s), Rows: sh.rows, PerRow: sh.perRow, Capacity: sh.rows * sh.perRow, FirstSeat: next})
		for r := 1; r <= sh.rows; r++ {
			for n := 1; n <= sh.perRow; n++ {
				v.Seats = append(v.Seats, Seat{ID: next, Section: int32(s), Row: int32(r), No: int32(n)})
				next++
			}
		}
	}
	order := make([]int, len(v.Seats))
	for i := range order {
		order[i] = i
	}
	centre := func(s Seat) int { // twice the distance from the row centre, an integer
		d := 2*int(s.No) - (sh.perRow + 1)
		if d < 0 {
			d = -d
		}
		return d
	}
	sort.SliceStable(order, func(a, b int) bool {
		x, y := v.Seats[order[a]], v.Seats[order[b]]
		if x.Section != y.Section {
			return x.Section < y.Section
		}
		if x.Row != y.Row {
			return x.Row < y.Row
		}
		if centre(x) != centre(y) {
			return centre(x) < centre(y)
		}
		return x.No < y.No
	})
	for rank, i := range order {
		v.Seats[i].Rank = int32(rank + 1)
	}
	return v
}

func (v *Venue) seat(id int32) Seat { return v.Seats[id-1] }

func (v *Venue) section(no int32) Section { return v.Sections[no-1] }

// adjacent reports whether b is the seat right after a in the same row.
func (v *Venue) adjacent(a, b int32) bool {
	if b != a+1 || int(b) > len(v.Seats) {
		return false
	}
	sa, sb := v.seat(a), v.seat(b)
	return sa.Section == sb.Section && sa.Row == sb.Row
}

type Band struct {
	ID      int64
	Name    string
	Genre   string
	Country string
}

type Event struct {
	ID          int64
	BandID      int64
	VenueID     int64
	Tier        int
	Name        string
	StartsAt    time.Time
	PriceCents  int64
	Description string
	Kind        string // catalogue | race | lifecycle
}

type Ticket struct {
	ID         int64
	EventID    int64
	Section    int32
	SeatID     int32
	CustomerID int64
	SoldAt     time.Time
	PriceCents int64
}

// LoadedHold is a hold present at load. Its times are offsets from the database's
// now() read just before the bulk copy, so every engine gets holds that are live
// or expired by its own clock.
type LoadedHold struct {
	ID         int64
	EventID    int64
	Section    int32
	Seats      []int32
	CustomerID int64
	ExpiresOff time.Duration // ms precision; negative = already expired
}

func (h LoadedHold) expiresAt(loadNow time.Time) time.Time { return loadNow.Add(h.ExpiresOff) }

func (h LoadedHold) createdAt(loadNow time.Time) time.Time {
	return loadNow.Add(h.ExpiresOff - 40*time.Minute)
}

func (h LoadedHold) live() bool { return h.ExpiresOff > 0 }

// Seat states in the generated ground truth.
const (
	stFree    uint8 = 0
	stLive    uint8 = 1 // held, valid
	stSold    uint8 = 2
	stExpired uint8 = 3 // held, expired: logically available
)

type Scale struct {
	Name      string
	Bands     int
	Customers int
	Catalogue map[int]int
	Race      map[int]int
	Lifecycle map[int]int
}

var scales = map[string]Scale{
	"tiny": {
		Name: "tiny", Bands: 10, Customers: 2000,
		Catalogue: map[int]int{10: 40, 100: 8, 1000: 2, 10000: 1},
		Race:      map[int]int{10: 20, 100: 4, 1000: 1, 10000: 1},
		Lifecycle: map[int]int{10: 4, 100: 2, 1000: 1},
	},
	"small": {
		Name: "small", Bands: 40, Customers: 20000,
		Catalogue: map[int]int{10: 400, 100: 80, 1000: 16, 10000: 3, 100000: 1},
		Race:      map[int]int{10: 100, 100: 20, 1000: 5, 10000: 2, 100000: 1},
		Lifecycle: map[int]int{10: 10, 100: 10, 1000: 10},
	},
}

type Dataset struct {
	Scale     string
	Seed      int64
	Customers int64
	Bands     []Band
	Venues    []*Venue
	venueByID map[int64]*Venue
	Events    []Event
	eventByID map[int64]*Event
	Sold      []Ticket
	Holds     []LoadedHold
	// State holds the generated seat states of catalogue events, indexed by
	// seat_id - 1. Race and lifecycle events are entirely free.
	State        map[int64][]uint8
	NextTicketID int64
	NextHoldID   int64
}

var (
	genres    = []string{"rock", "jazz", "electronic", "folk", "metal", "pop", "hip-hop", "classical"}
	countries = []string{"BR", "PT", "UK", "DE", "US", "JP", "AR", "FR"}
	loadEpoch = time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
	showEpoch = time.Date(2026, 10, 1, 20, 0, 0, 0, time.UTC)
)

// Generate builds the dataset for a scale, optionally restricted to some tiers.
func Generate(scaleName string, seed int64, tiers []int) (*Dataset, error) {
	sc, ok := scales[scaleName]
	if !ok {
		return nil, fmt.Errorf("unknown scale %q", scaleName)
	}
	keep := map[int]bool{}
	for _, t := range tiers {
		keep[t] = true
	}
	want := func(t int) bool { return len(keep) == 0 || keep[t] }
	r := rand.New(rand.NewSource(seed))
	ds := &Dataset{Scale: sc.Name, Seed: seed, Customers: int64(sc.Customers),
		venueByID: map[int64]*Venue{}, eventByID: map[int64]*Event{}, State: map[int64][]uint8{}}

	for i := 1; i <= sc.Bands; i++ {
		ds.Bands = append(ds.Bands, Band{
			ID: int64(i), Name: fmt.Sprintf("Band %03d", i),
			Genre: genres[r.Intn(len(genres))], Country: countries[r.Intn(len(countries))],
		})
	}
	for i, t := range allTiers {
		if !want(t) {
			continue
		}
		v := buildVenue(int64(i+1), t)
		ds.Venues = append(ds.Venues, v)
		ds.venueByID[v.ID] = v
	}
	venueOfTier := func(t int) *Venue {
		for _, v := range ds.Venues {
			if v.Tier == t {
				return v
			}
		}
		return nil
	}

	bandZipf := rand.NewZipf(r, 1.3, 2, uint64(sc.Bands-1))
	eventID := int64(1)
	addEvents := func(kind string, counts map[int]int) {
		for _, tier := range allTiers {
			if !want(tier) {
				continue
			}
			v := venueOfTier(tier)
			for n := 0; n < counts[tier]; n++ {
				ds.Events = append(ds.Events, Event{
					ID: eventID, BandID: int64(bandZipf.Uint64()) + 1, VenueID: v.ID, Tier: tier,
					Name:        fmt.Sprintf("%s show %d", kind, eventID),
					StartsAt:    showEpoch.Add(time.Duration(r.Intn(365*24)) * time.Hour),
					PriceCents:  int64(2000 + 500*r.Intn(30)),
					Description: "Doors open one hour before the show.",
					Kind:        kind,
				})
				eventID++
			}
		}
	}
	addEvents("catalogue", sc.Catalogue)
	addEvents("race", sc.Race)
	addEvents("lifecycle", sc.Lifecycle)
	for i := range ds.Events {
		ds.eventByID[ds.Events[i].ID] = &ds.Events[i]
	}

	custZipf := rand.NewZipf(r, 1.1, 5, uint64(sc.Customers-1))
	customer := func() int64 { return int64(custZipf.Uint64()) + 1 }
	nextTicket := int64(1)
	nextHold := int64(1)
	for i := range ds.Events {
		e := &ds.Events[i]
		if e.Kind != "catalogue" {
			continue
		}
		v := ds.venueByID[e.VenueID]
		n := v.Capacity
		state := make([]uint8, n)

		// Sold: the best seats first, with noise.
		f := 0.2 + 0.4*r.Float64()
		type keyed struct {
			id  int32
			key float64
		}
		ks := make([]keyed, n)
		for j, s := range v.Seats {
			ks[j] = keyed{s.ID, float64(s.Rank) + r.Float64()*0.3*float64(n)}
		}
		sort.SliceStable(ks, func(a, b int) bool {
			if ks[a].key != ks[b].key {
				return ks[a].key < ks[b].key
			}
			return ks[a].id < ks[b].id
		})
		nSold := int(float64(n) * f)
		for j := 0; j < nSold; j++ {
			id := ks[j].id
			state[id-1] = stSold
			ds.Sold = append(ds.Sold, Ticket{
				ID: nextTicket, EventID: e.ID, Section: v.seat(id).Section, SeatID: id, CustomerID: customer(),
				SoldAt: loadEpoch.Add(time.Duration(r.Intn(90*24*3600)) * time.Second), PriceCents: e.PriceCents,
			})
			nextTicket++
		}

		// Holds: live, then expired, as blocks of 1-4 adjacent free seats in a row.
		free := make([]int32, 0, n-nSold)
		for _, s := range v.Seats {
			if state[s.ID-1] == stFree {
				free = append(free, s.ID)
			}
		}
		remaining := len(free)
		r.Shuffle(len(free), func(a, b int) { free[a], free[b] = free[b], free[a] })
		place := func(target int, live bool) {
			placed := 0
			for _, start := range free {
				if placed >= target {
					return
				}
				if state[start-1] != stFree {
					continue
				}
				k := 1 + r.Intn(4)
				block := []int32{start}
				for len(block) < k && placed+len(block) < target {
					nxt := block[len(block)-1] + 1
					if !v.adjacent(block[len(block)-1], nxt) || state[nxt-1] != stFree {
						break
					}
					block = append(block, nxt)
				}
				h := LoadedHold{ID: nextHold, EventID: e.ID, Section: v.seat(start).Section, Seats: block, CustomerID: customer()}
				st := stLive
				if live {
					h.ExpiresOff = 10*time.Minute + time.Duration(r.Int63n(int64(30*time.Minute)))
				} else {
					h.ExpiresOff = -(time.Minute + time.Duration(r.Int63n(int64(59*time.Minute))))
					st = stExpired
				}
				h.ExpiresOff = h.ExpiresOff.Truncate(time.Millisecond)
				for _, s := range block {
					state[s-1] = st
				}
				placed += len(block)
				ds.Holds = append(ds.Holds, h)
				nextHold++
			}
		}
		place(int(0.04*float64(remaining)), true)
		place(int(0.01*float64(remaining)), false)
		ds.State[e.ID] = state
	}
	ds.NextTicketID = nextTicket
	ds.NextHoldID = nextHold
	return ds, nil
}

func (ds *Dataset) venue(id int64) *Venue { return ds.venueByID[id] }

func (ds *Dataset) event(id int64) *Event { return ds.eventByID[id] }

func (ds *Dataset) Summary() map[string]any {
	perKind := map[string]int{}
	seats := map[string]int64{}
	tiers := map[string]int{}
	for _, e := range ds.Events {
		perKind[e.Kind]++
		seats[e.Kind] += int64(ds.venue(e.VenueID).Capacity)
		tiers[fmt.Sprintf("%s_%d", e.Kind, e.Tier)]++
	}
	live, expired, heldSeats := 0, 0, 0
	for _, h := range ds.Holds {
		if h.live() {
			live++
		} else {
			expired++
		}
		heldSeats += len(h.Seats)
	}
	return map[string]any{
		"scale": ds.Scale, "seed": ds.Seed,
		"bands": len(ds.Bands), "customers": ds.Customers, "venues": len(ds.Venues),
		"events": len(ds.Events), "events_by_kind": perKind, "events_by_kind_tier": tiers,
		"seats_by_kind": seats, "tickets_sold_at_load": len(ds.Sold),
		"live_holds_at_load": live, "expired_holds_at_load": expired, "held_seats_at_load": heldSeats,
	}
}

// eventsOf returns events of a kind and tier, in id order.
func (ds *Dataset) eventsOf(kind string, tier int) []*Event {
	var out []*Event
	for i := range ds.Events {
		e := &ds.Events[i]
		if e.Kind == kind && (tier == 0 || e.Tier == tier) {
			out = append(out, e)
		}
	}
	return out
}

func (ds *Dataset) tiersOf(kind string) []int {
	seen := map[int]bool{}
	var out []int
	for _, e := range ds.Events {
		if e.Kind == kind && !seen[e.Tier] {
			seen[e.Tier] = true
			out = append(out, e.Tier)
		}
	}
	sort.Ints(out)
	return out
}
