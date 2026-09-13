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
// seed: the same bands, the same events of the same sizes, the same tickets
// already sold to the same customers in the same seats. Designs differ only in
// how that logical state is laid out in tables (a sold ticket row plus unsold
// ticket rows, or a sold ticket row plus a counter, or plus a pool of slots).
//
// Events come in three kinds:
//
//	catalogue  partially sold (20-60%), the backdrop for reads, spread-demand
//	           booking and cancellation
//	race       unsold; each is sold out by a crowd of concurrent buyers
//	churn      unsold; sold out while some buyers cancel -- the regime where a
//	           design can under-sell
//
// Event sizes are the study's tiers: 10, 100, 1 000, 10 000 and 100 000 seats.
// ---------------------------------------------------------------------------

var allTiers = []int{10, 100, 1000, 10000, 100000}

type Band struct {
	ID      int64
	Name    string
	Genre   string
	Country string
}

type Event struct {
	ID          int64
	BandID      int64
	Name        string
	Venue       string
	StartsAt    time.Time
	Capacity    int
	PriceCents  int64
	Description string
	Kind        string // catalogue | race | churn
	InitialSold int
	// FirstTicketID is the ticket id of seat 1. Ticket ids are allocated per
	// seat across every event, so a sold ticket has the same id in every design
	// whether or not its unsold neighbours exist as rows.
	FirstTicketID int64
}

type Ticket struct {
	ID         int64
	EventID    int64
	SeatNo     int
	CustomerID int64
	SoldAt     time.Time
	PriceCents int64
}

type Scale struct {
	Name      string
	Bands     int
	Customers int
	Catalogue map[int]int // tier -> events
	Race      map[int]int
	Churn     map[int]int
}

// Tier counts are chosen so each tier's total capacity is of the same order
// (the 100 000 tier is a single event), and so the small tiers get enough
// independent races for a rare overbooking to have many chances to appear.
var scales = map[string]Scale{
	"tiny": {
		Name: "tiny", Bands: 10, Customers: 2000,
		Catalogue: map[int]int{10: 40, 100: 8, 1000: 2, 10000: 1},
		Race:      map[int]int{10: 20, 100: 4, 1000: 1, 10000: 1},
		Churn:     map[int]int{10: 10, 100: 2, 1000: 1},
	},
	"small": {
		Name: "small", Bands: 40, Customers: 20000,
		Catalogue: map[int]int{10: 400, 100: 80, 1000: 16, 10000: 3, 100000: 1},
		Race:      map[int]int{10: 100, 100: 20, 1000: 5, 10000: 2, 100000: 1},
		Churn:     map[int]int{10: 50, 100: 10, 1000: 2},
	},
}

type Dataset struct {
	Scale     string
	Seed      int64
	Customers int64
	Bands     []Band
	Events    []Event
	Sold      []Ticket // tickets sold before the benchmark starts
	// NextTicketID is the first id never used by a loaded seat.
	NextTicketID int64
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
	r := rand.New(rand.NewSource(seed))
	ds := &Dataset{Scale: sc.Name, Seed: seed, Customers: int64(sc.Customers)}

	for i := 1; i <= sc.Bands; i++ {
		ds.Bands = append(ds.Bands, Band{
			ID: int64(i), Name: fmt.Sprintf("Band %03d", i),
			Genre: genres[r.Intn(len(genres))], Country: countries[r.Intn(len(countries))],
		})
	}
	// Bands are uneven: a few headline acts own many events. Zipf over bands,
	// like real rosters, so "a band's events" ranges from one to dozens.
	bandZipf := rand.NewZipf(r, 1.3, 2, uint64(sc.Bands-1))

	nextTicket := int64(1)
	eventID := int64(1)
	addEvents := func(kind string, counts map[int]int) {
		for _, tier := range allTiers {
			if len(keep) > 0 && !keep[tier] {
				continue
			}
			for n := 0; n < counts[tier]; n++ {
				e := Event{
					ID:            eventID,
					BandID:        int64(bandZipf.Uint64()) + 1,
					Name:          fmt.Sprintf("%s show %d", kind, eventID),
					Venue:         venueFor(tier),
					StartsAt:      showEpoch.Add(time.Duration(r.Intn(365*24)) * time.Hour),
					Capacity:      tier,
					PriceCents:    int64(2000 + 500*r.Intn(30)),
					Description:   "Doors open one hour before the show.",
					Kind:          kind,
					FirstTicketID: nextTicket,
				}
				if kind == "catalogue" {
					e.InitialSold = int(float64(tier) * (0.2 + 0.4*r.Float64()))
				}
				nextTicket += int64(tier)
				eventID++
				ds.Events = append(ds.Events, e)
			}
		}
	}
	addEvents("catalogue", sc.Catalogue)
	addEvents("race", sc.Race)
	addEvents("churn", sc.Churn)
	ds.NextTicketID = nextTicket

	// Customers are skewed too: a minority buy many tickets, so "my tickets"
	// ranges from one row to hundreds.
	custZipf := rand.NewZipf(r, 1.1, 5, uint64(sc.Customers-1))
	for i := range ds.Events {
		e := &ds.Events[i]
		for s := 1; s <= e.InitialSold; s++ {
			ds.Sold = append(ds.Sold, Ticket{
				ID:         e.FirstTicketID + int64(s) - 1,
				EventID:    e.ID,
				SeatNo:     s,
				CustomerID: int64(custZipf.Uint64()) + 1,
				SoldAt:     loadEpoch.Add(time.Duration(r.Intn(90*24*3600)) * time.Second),
				PriceCents: e.PriceCents,
			})
		}
	}
	return ds, nil
}

func venueFor(tier int) string {
	switch {
	case tier <= 10:
		return "back room"
	case tier <= 100:
		return "club"
	case tier <= 1000:
		return "theatre"
	case tier <= 10000:
		return "arena"
	default:
		return "stadium"
	}
}

func (ds *Dataset) Summary() map[string]any {
	perKind := map[string]int{}
	seats := map[string]int64{}
	for _, e := range ds.Events {
		perKind[e.Kind]++
		seats[e.Kind] += int64(e.Capacity)
	}
	tiers := map[string]int{}
	for _, e := range ds.Events {
		tiers[fmt.Sprintf("%s_%d", e.Kind, e.Capacity)]++
	}
	return map[string]any{
		"scale": ds.Scale, "seed": ds.Seed,
		"bands": len(ds.Bands), "customers": ds.Customers,
		"events": len(ds.Events), "events_by_kind": perKind, "events_by_kind_tier": tiers,
		"seats_by_kind": seats, "tickets_sold_at_load": len(ds.Sold),
		"seats_total": ds.NextTicketID - 1,
	}
}

// eventsOf returns events of a kind and tier, in id order.
func (ds *Dataset) eventsOf(kind string, tier int) []*Event {
	var out []*Event
	for i := range ds.Events {
		e := &ds.Events[i]
		if e.Kind == kind && (tier == 0 || e.Capacity == tier) {
			out = append(out, e)
		}
	}
	return out
}

func (ds *Dataset) tiersOf(kind string) []int {
	seen := map[int]bool{}
	var out []int
	for _, e := range ds.Events {
		if e.Kind == kind && !seen[e.Capacity] {
			seen[e.Capacity] = true
			out = append(out, e.Capacity)
		}
	}
	sort.Ints(out)
	return out
}
