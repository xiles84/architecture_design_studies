package main

import (
	"fmt"

	"adsplatform/ports"
)

// Strategy names the booking flow the harness drives. The SQL each flow issues
// lives in the design's writes.sql; the flow decides which statements, in which
// transaction, and what "no row" means.
type Strategy string

const (
	LockFirst  Strategy = "lock_first"  // P1: FOR UPDATE on the lowest free seat
	SkipLocked Strategy = "skip_locked" // P2, P4: FOR UPDATE SKIP LOCKED + exact fallback
	CAS        Strategy = "cas"         // P3: unlocked read + compare-and-set update
	CountCheck Strategy = "count_check" // C1, C2, C3: count, then insert
	Counter    Strategy = "counter"     // C4, R1: guarded counter update, then insert
	SeatUnique Strategy = "seat_unique" // C5: highest+1 insert, unique index arbitrates
	Buckets    Strategy = "buckets"     // R2: sharded counter rows
	SeatPool   Strategy = "seat_pool"   // R3: consume a pre-created slot, insert ticket
	Hold       Strategy = "hold"        // H0, H1: reservation with expiry, then confirm
)

// Design describes one way of preventing overbooking. The flags tell the loader
// how to shape the logical dataset for this design's tables; everything a
// reader or a buyer does is driven by the design's SQL files.
type Design struct {
	ID      string
	Short   string
	Title   string
	Family  string
	Summary string

	Strategy  Strategy
	Isolation ports.Isolation
	// LockEvent takes a FOR UPDATE lock on the event row before counting (C3).
	LockEvent bool

	// Loader shape.
	Precreated     bool   // a ticket row per seat, with a status column
	CounterColumn  string // counter on the event row: "seats_sold" | "seats_taken" | ""
	InventoryRow   bool   // event_inventory(event_id, remaining)
	InventoryShard bool   // event_inventory_bucket rows
	SeatPool       bool   // seat_slot rows for unsold seats
	TicketCapacity bool   // ticket.event_capacity (C5's CHECK)
	Holds          bool   // reservation table; ticket.reservation_id

	// NegativeControl is non-empty for a design that is EXPECTED to violate the
	// invariant. It stays in the study because an audit that has never caught a
	// wrong design has not been shown to work. The report checks that it fired.
	NegativeControl string
	// ControlFires names the experiment in which the negative control is
	// expected to fire: "race" or "holds".
	ControlFires string
}

var designs = []Design{
	{
		ID: "p1_precreated_lock_first", Short: "P1 lock-first", Title: "precreated-lock-first", Family: "pre-created tickets",
		Summary:  "A ticket row per seat; booking locks the lowest available seat with FOR UPDATE, then sells it.",
		Strategy: LockFirst, Isolation: ports.ReadCommitted, Precreated: true,
	},
	{
		ID: "p2_precreated_skip_locked", Short: "P2 skip-locked", Title: "precreated-skip-locked", Family: "pre-created tickets",
		Summary:  "P1 with FOR UPDATE SKIP LOCKED: buyers pass over seats others are taking, with an exact sold-out fallback.",
		Strategy: SkipLocked, Isolation: ports.ReadCommitted, Precreated: true,
	},
	{
		ID: "p3_precreated_cas", Short: "P3 CAS", Title: "precreated-cas", Family: "pre-created tickets",
		Summary:  "No locks: read a candidate from a random start seat, sell it with a compare-and-set UPDATE, retry on loss.",
		Strategy: CAS, Isolation: ports.ReadCommitted, Precreated: true,
	},
	{
		ID: "p4_precreated_counter", Short: "P4 skip+counter", Title: "precreated-counter", Family: "pre-created tickets",
		Summary:  "P2 plus a seats_sold counter on the event row, maintained in the same transaction.",
		Strategy: SkipLocked, Isolation: ports.ReadCommitted, Precreated: true, CounterColumn: "seats_sold",
	},
	{
		ID: "c1_count_naive", Short: "C1 count (RC) ✗", Title: "count-naive", Family: "tickets created on booking",
		Summary:  "Count tickets, insert if below capacity, at READ COMMITTED. Negative control: expected to overbook.",
		Strategy: CountCheck, Isolation: ports.ReadCommitted,
		NegativeControl: "check-then-insert at READ COMMITTED: two buyers can both see the last seat free",
		ControlFires:    "race",
	},
	{
		ID: "c2_count_serializable", Short: "C2 count (SER)", Title: "count-serializable", Family: "tickets created on booking",
		Summary:  "C1's byte-identical SQL at SERIALIZABLE isolation, retrying serialization failures.",
		Strategy: CountCheck, Isolation: ports.Serializable,
	},
	{
		ID: "c3_count_lock_event", Short: "C3 count+lock", Title: "count-lock-event", Family: "tickets created on booking",
		Summary:  "C1 preceded by SELECT ... FOR UPDATE on the event row: a per-event critical section around the count.",
		Strategy: CountCheck, Isolation: ports.ReadCommitted, LockEvent: true,
	},
	{
		ID: "c4_counter_guard", Short: "C4 counter", Title: "counter-guard", Family: "tickets created on booking",
		Summary:  "seats_sold on the event row, claimed with one guarded UPDATE ... WHERE seats_sold < capacity.",
		Strategy: Counter, Isolation: ports.ReadCommitted, CounterColumn: "seats_sold",
	},
	{
		ID: "c5_seat_unique", Short: "C5 seat-unique", Title: "seat-unique", Family: "tickets created on booking",
		Summary:  "No lock or counter: insert seat highest+1; UNIQUE(event, seat) and CHECK(seat <= capacity) arbitrate.",
		Strategy: SeatUnique, Isolation: ports.ReadCommitted, TicketCapacity: true,
	},
	{
		ID: "r1_inventory_row", Short: "R1 inventory", Title: "inventory-row", Family: "extra inventory table",
		Summary:  "C4's guarded counter moved off the event row into a separate event_inventory table.",
		Strategy: Counter, Isolation: ports.ReadCommitted, InventoryRow: true,
	},
	{
		ID: "r2_inventory_buckets", Short: "R2 buckets", Title: "inventory-buckets", Family: "extra inventory table",
		Summary:  "R1 split into up to 8 bucket rows per event, taken with SKIP LOCKED: a sharded counter.",
		Strategy: Buckets, Isolation: ports.ReadCommitted, InventoryShard: true,
	},
	{
		ID: "r3_seat_pool", Short: "R3 seat-pool", Title: "seat-pool", Family: "extra inventory table",
		Summary:  "Unsold seats as narrow pre-created slot rows; booking deletes a slot (SKIP LOCKED) and inserts the ticket.",
		Strategy: SeatPool, Isolation: ports.ReadCommitted, SeatPool: true,
	},
	{
		ID: "h0_hold_naive_confirm", Short: "H0 hold/naive ✗", Title: "hold-naive-confirm", Family: "reservation with expiry",
		Summary:  "Hold with expiry then confirm; confirmation trusts the pre-payment check. Negative control for late payments.",
		Strategy: Hold, Isolation: ports.ReadCommitted, CounterColumn: "seats_taken", Holds: true,
		NegativeControl: "confirms a hold that expired during payment, after the seat may have been resold",
		ControlFires:    "holds",
	},
	{
		ID: "h1_hold_checked_confirm", Short: "H1 hold/checked", Title: "hold-checked-confirm", Family: "reservation with expiry",
		Summary:  "H0 with the confirmation conditional on the hold still being held and unexpired, by the database clock.",
		Strategy: Hold, Isolation: ports.ReadCommitted, CounterColumn: "seats_taken", Holds: true,
	},
}

func designByID(id string) (Design, error) {
	for _, d := range designs {
		if d.ID == id {
			return d, nil
		}
	}
	return Design{}, fmt.Errorf("unknown design %q (see -cmd list)", id)
}

func designShort(id string) string {
	for _, d := range designs {
		if d.ID == id {
			return d.Short
		}
	}
	return id
}
