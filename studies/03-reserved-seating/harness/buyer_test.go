package main

import (
	"math/rand"
	"testing"
)

func TestVenueShapes(t *testing.T) {
	for _, tier := range allTiers {
		v := buildVenue(1, tier)
		if v.Capacity != tier || len(v.Seats) != tier {
			t.Fatalf("tier %d: capacity %d, seats %d", tier, v.Capacity, len(v.Seats))
		}
		ranks := map[int32]bool{}
		for i, s := range v.Seats {
			if s.ID != int32(i+1) {
				t.Fatalf("tier %d: seat %d has id %d", tier, i+1, s.ID)
			}
			if ranks[s.Rank] || s.Rank < 1 || int(s.Rank) > tier {
				t.Fatalf("tier %d: rank %d repeated or out of range", tier, s.Rank)
			}
			ranks[s.Rank] = true
			sec := v.section(s.Section)
			want := sec.FirstSeat + (s.Row-1)*int32(sec.PerRow) + s.No - 1
			if s.ID != want {
				t.Fatalf("tier %d: seat %d not at its (section, row, number) position %d", tier, s.ID, want)
			}
		}
	}
	// Section 1, row 1, centre of the row is the best seat.
	v := buildVenue(1, 1000)
	best := v.Seats[0]
	for _, s := range v.Seats {
		if s.Rank == 1 {
			best = s
		}
	}
	if best.Section != 1 || best.Row != 1 || (best.No != 13) {
		t.Fatalf("best seat %+v, want section 1 row 1 seat 13", best)
	}
}

func TestBlocksIn(t *testing.T) {
	v := buildVenue(1, 100) // 1 section, 5 rows of 20
	sec := v.section(1)
	all := blocksIn(v, sec, map[int32]uint8{}, 2)
	if len(all) != 5*19 {
		t.Fatalf("pairs in an empty club: %d, want 95", len(all))
	}
	for _, b := range all {
		if !v.adjacent(b[0], b[1]) {
			t.Fatalf("block %v is not adjacent in one row", b)
		}
	}
	if s := v.seat(all[0][0]); s.Row != 1 {
		t.Fatalf("best pair should be in row 1, got %+v", s)
	}
	// Take every seat of row 1 except seats 1 and 2 and 20: row 1 offers only (1,2).
	taken := map[int32]uint8{}
	for n := int32(3); n <= 19; n++ {
		taken[n] = 1
	}
	blocks := blocksIn(v, sec, taken, 2)
	for _, b := range blocks {
		if v.seat(b[0]).Row == 1 && !(b[0] == 1 && b[1] == 2) {
			t.Fatalf("block %v crosses a taken seat", b)
		}
	}
	// A block never wraps from the end of one row to the start of the next.
	for _, b := range blocksIn(v, sec, map[int32]uint8{}, 3) {
		if v.seat(b[0]).Row != v.seat(b[2]).Row {
			t.Fatalf("block %v wraps rows", b)
		}
	}
	if n := len(blocksIn(v, sec, map[int32]uint8{}, 21)); n != 0 {
		t.Fatalf("no row holds 21 adjacent seats, got %d blocks", n)
	}
}

func TestPartySizesAndZipf(t *testing.T) {
	r := rand.New(rand.NewSource(1))
	counts := map[int]int{}
	for i := 0; i < 10000; i++ {
		counts[partySize(r)]++
	}
	if counts[2] < 4000 || counts[2] > 5000 || counts[5] != 0 {
		t.Fatalf("party sizes %v", counts)
	}
	low := 0
	for i := 0; i < 10000; i++ {
		k := zipfIndex(r, 1.2, 16)
		if k < 0 || k >= 16 {
			t.Fatalf("zipf index %d out of range", k)
		}
		if k == 0 {
			low++
		}
	}
	if low < 2000 {
		t.Fatalf("the best block should be chosen most often, got %d/10000", low)
	}
}

func TestDatasetIsDeterministicAndConsistent(t *testing.T) {
	a, _ := Generate("tiny", 42, nil)
	b, _ := Generate("tiny", 42, nil)
	if len(a.Sold) != len(b.Sold) || len(a.Holds) != len(b.Holds) || a.Sold[len(a.Sold)-1] != b.Sold[len(b.Sold)-1] {
		t.Fatal("same seed, different dataset")
	}
	for _, h := range a.Holds {
		st := a.State[h.EventID]
		v := a.venue(a.event(h.EventID).VenueID)
		for i, s := range h.Seats {
			want := stExpired
			if h.live() {
				want = stLive
			}
			if st[s-1] != want {
				t.Fatalf("hold %d seat %d state %d, want %d", h.ID, s, st[s-1], want)
			}
			if i > 0 && !v.adjacent(h.Seats[i-1], s) {
				t.Fatalf("hold %d seats %v not adjacent", h.ID, h.Seats)
			}
			if v.seat(s).Section != h.Section {
				t.Fatalf("hold %d seat %d outside its section", h.ID, s)
			}
		}
	}
	for _, tk := range a.Sold {
		if a.State[tk.EventID][tk.SeatID-1] != stSold {
			t.Fatalf("ticket %d seat not sold in the state", tk.ID)
		}
	}
	doc := Claims{3: {State: 1, Hold: 9, Customer: 4, ExpiresMS: 1_700_000_000_123}, 5: {State: 2, Customer: 8}}
	text, next := doc.Encode()
	back, err := DecodeClaims(text)
	if err != nil || len(back) != 2 || back[3] != doc[3] || back[5] != doc[5] || next == nil {
		t.Fatalf("section document round trip: %v %v %v", back, err, next)
	}
}
