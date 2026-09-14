package main

import (
	"context"
	"math/rand"
	"sort"
	"time"
)

// ---------------------------------------------------------------------------
// The buyer model, shared by the race, the lifecycle and the spread-demand write.
//
// A customer with a party of n reads the event's per-section availability, picks
// a section (the better ones more often), reads that section's seat map, lists
// every run of n adjacent free seats in a row, ranks them by seat quality and
// picks one of the best 16 (the best more often). If the hold finds a seat taken
// -- someone chose it a moment earlier -- the customer re-reads the map and
// chooses again, up to 20 times.
//
// The venue layout is static and cached on the client, like a seat-map image
// from a CDN; only the dynamic state is read from the database. The random
// stream a customer uses never depends on the design; the outcomes it meets do.
// ---------------------------------------------------------------------------

// partySize: 1 -> 20%, 2 -> 45%, 3 -> 10%, 4 -> 20%, 6 -> 5%.
func partySize(r *rand.Rand) int {
	x := r.Intn(100)
	switch {
	case x < 20:
		return 1
	case x < 65:
		return 2
	case x < 75:
		return 3
	case x < 95:
		return 4
	default:
		return 6
	}
}

// zipfIndex draws 0..k-1 with probability falling as a power s of the rank.
func zipfIndex(r *rand.Rand, s float64, k int) int {
	if k <= 1 {
		return 0
	}
	return int(rand.NewZipf(r, s, 1, uint64(k-1)).Uint64())
}

// blocksIn lists every run of n adjacent free seats in one row of a section,
// best first (lowest sum of quality ranks, then lowest first seat).
func blocksIn(v *Venue, sec Section, taken map[int32]uint8, n int) [][]int32 {
	type cand struct {
		first int32
		score int64
	}
	var cands []cand
	for row := 0; row < sec.Rows; row++ {
		base := sec.FirstSeat + int32(row*sec.PerRow)
		run := 0
		var score int64
		for i := 0; i < sec.PerRow; i++ {
			seat := base + int32(i)
			if _, t := taken[seat]; t {
				run, score = 0, 0
				continue
			}
			run++
			score += int64(v.seat(seat).Rank)
			if run > n {
				score -= int64(v.seat(seat - int32(n)).Rank)
				run = n
			}
			if run == n {
				cands = append(cands, cand{seat - int32(n) + 1, score})
			}
		}
	}
	sort.Slice(cands, func(i, j int) bool {
		if cands[i].score != cands[j].score {
			return cands[i].score < cands[j].score
		}
		return cands[i].first < cands[j].first
	})
	out := make([][]int32, len(cands))
	for i, c := range cands {
		b := make([]int32, n)
		for k := range b {
			b[k] = c.first + int32(k)
		}
		out[i] = b
	}
	return out
}

// Acquisition is how one customer's attempt to hold seats went.
type Acquisition struct {
	Block         Block
	Hold          HoldResult
	Granted       bool
	NoBlock       bool // no run of n free seats anywhere: the customer leaves
	GaveUp        bool // too many conflicts
	Conflicts     int
	MapReads      int
	SectionReads  int
	EngineRetries int
	HoldLatency   time.Duration // the successful hold statement's round trip
}

type Buyer struct {
	seller       *Seller
	node         int
	r            *rand.Rand
	maxConflicts int
	// gate, when set, is held around every database operation, so a probe can
	// pause every buyer between operations.
	gate interface {
		RLock()
		RUnlock()
	}
}

func (b *Buyer) do(f func()) {
	if b.gate != nil {
		b.gate.RLock()
		defer b.gate.RUnlock()
	}
	f()
}

// Acquire tries to hold n adjacent seats of ev for customer.
func (b *Buyer) Acquire(ctx context.Context, ev *Event, customer int64, n int, ttl time.Duration) (Acquisition, error) {
	var a Acquisition
	s := b.seller
	v := s.world.ds.venue(ev.VenueID)
	var secs []SectionAvail
	var err error
	b.do(func() { secs, err = s.Sections(ctx, b.node, ev) })
	a.SectionReads++
	if err != nil {
		return a, err
	}
	excluded := map[int32]bool{}
	dropped := 0
	for {
		var cands []int32
		for _, x := range secs {
			if x.Available >= int64(n) && !excluded[x.No] {
				cands = append(cands, x.No)
			}
		}
		if len(cands) == 0 {
			a.NoBlock = true
			return a, nil
		}
		secNo := cands[zipfIndex(b.r, 1.1, len(cands))]
		sec := v.section(secNo)
		for {
			if err := ctx.Err(); err != nil {
				return a, err
			}
			if dl, ok := ctx.Value(attemptDeadlineKey{}).(time.Time); ok && time.Now().After(dl) {
				return a, errAttemptDeadline
			}
			var taken map[int32]uint8
			b.do(func() { taken, err = s.SectionMap(ctx, b.node, ev, secNo) })
			a.MapReads++
			if err != nil {
				return a, err
			}
			blocks := blocksIn(v, sec, taken, n)
			if len(blocks) == 0 {
				excluded[secNo] = true
				if dropped++; dropped%3 == 0 {
					b.do(func() { secs, err = s.Sections(ctx, b.node, ev) })
					a.SectionReads++
					if err != nil {
						return a, err
					}
				}
				break
			}
			choice := blocks[zipfIndex(b.r, 1.2, min(16, len(blocks)))]
			blk := Block{Event: ev, Section: secNo, Seats: choice}
			t0 := time.Now()
			var hr HoldResult
			b.do(func() { hr, err = s.Hold(ctx, b.node, blk, customer, ttl) })
			a.EngineRetries += hr.Retries
			if err != nil {
				return a, err
			}
			if hr.Granted {
				a.Block, a.Hold, a.Granted, a.HoldLatency = blk, hr, true, time.Since(t0)
				return a, nil
			}
			if a.Conflicts++; a.Conflicts >= b.maxConflicts {
				a.GaveUp = true
				return a, nil
			}
		}
	}
}
