package main

import (
	"fmt"

	"adsplatform/ports"
)

// Strategy names the flow the harness drives for a hold and its confirmation.
// The SQL each flow issues lives in the design's writes.sql; the flow decides
// which statements, in which transaction, and what a short row count means.
type Strategy string

const (
	Guarded      Strategy = "guarded"        // S1, E0, E1, K0, K1, L3: one conditional multi-row UPDATE
	LockWait     Strategy = "lock_wait"      // S2: SELECT ... FOR UPDATE, then the conditional UPDATE
	LockNowait   Strategy = "lock_nowait"    // S3: the same with NOWAIT
	CheckThenAct Strategy = "check_then_act" // S0 (RC), S4 (SER): plain read, unconditional UPDATE
	Cart         Strategy = "cart"           // E2: cart row + seats pointing to it
	Claim        Strategy = "claim"          // L1: claim rows created on hold, unique key arbitrates
	Document     Strategy = "document"       // L2: a JSONB document per section, version CAS
)

// Layout tells the loader which inventory tables the design has.
type Layout string

const (
	SeatRows  Layout = "seat_rows"  // event_seat, expiry on every seat row
	CartRows  Layout = "cart_rows"  // hold + event_seat
	ClaimRows Layout = "claim_rows" // seat_claim, no per-event seat rows
	Documents Layout = "documents"  // event_section with JSONB claims
)

// Design describes one way of keeping a chosen seat for its buyer. The flags tell
// the loader how to lay out the logical dataset and the seller which flow to run;
// everything the database does is in the design's SQL files.
type Design struct {
	ID      string
	Short   string
	Title   string
	Family  string
	Summary string

	Strategy  Strategy
	Isolation ports.Isolation
	Layout    Layout

	// ExpiryOnSweeper: expiry takes effect only when the sweeper releases the hold
	// (E1). The harness runs the sweeper to completion before verification.
	ExpiryOnSweeper bool
	// AppClock: the design judges expiry by $app_now from simulated application
	// nodes, one of them running ahead (E0).
	AppClock bool
	// PaymentWindow: starting checkout extends the hold (K1).
	PaymentWindow bool
	// YBOnly: the SQL uses YugabyteDB-only syntax (L3).
	YBOnly bool

	// NegativeControl is non-empty for a design EXPECTED to violate an invariant.
	// It stays in the study because an audit that has never caught a wrong design
	// has not been shown to work. The report checks that it fired.
	NegativeControl string
	// ControlFires names the experiments in which the control is expected to fire.
	ControlFires []string
}

// Family labels, in report order.
const (
	famArbitration = "arbitration"
	famExpiry      = "expiry"
	famCheckout    = "checkout"
	famLayout      = "layout"
)

var designs = []Design{
	{
		ID: "s0_check_then_hold_rc", Short: "S0 check/RC ✗", Title: "check-then-hold-rc", Family: famArbitration,
		Summary:  "Read the block's seats, then mark them held unconditionally, at READ COMMITTED. Negative control: expected to steal holds.",
		Strategy: CheckThenAct, Isolation: ports.ReadCommitted, Layout: SeatRows,
		NegativeControl: "check-then-act at READ COMMITTED: a second buyer's unconditional update overwrites the first hold after waiting for it",
		ControlFires:    []string{"race", "lifecycle"},
	},
	{
		ID: "s1_conditional_update", Short: "S1 conditional", Title: "conditional-update", Family: famArbitration,
		Summary:  "The reference: one conditional multi-row UPDATE takes the block or fewer rows come back and it rolls back. Lazy expiry, checked confirmation.",
		Strategy: Guarded, Isolation: ports.ReadCommitted, Layout: SeatRows,
	},
	{
		ID: "s2_lock_then_update", Short: "S2 lock", Title: "lock-then-update", Family: famArbitration,
		Summary:  "SELECT ... ORDER BY seat_id FOR UPDATE on the block, check, then S1's update.",
		Strategy: LockWait, Isolation: ports.ReadCommitted, Layout: SeatRows,
	},
	{
		ID: "s3_lock_nowait", Short: "S3 nowait", Title: "lock-nowait", Family: famArbitration,
		Summary:  "S2 with FOR UPDATE NOWAIT: a seat another buyer is taking fails at once and the buyer chooses again.",
		Strategy: LockNowait, Isolation: ports.ReadCommitted, Layout: SeatRows,
	},
	{
		ID: "s4_check_then_hold_serializable", Short: "S4 check/SER", Title: "check-then-hold-serializable", Family: famArbitration,
		Summary:  "S0's SQL at SERIALIZABLE, retrying serialization failures.",
		Strategy: CheckThenAct, Isolation: ports.Serializable, Layout: SeatRows,
	},
	{
		ID: "e0_app_clock_expiry", Short: "E0 app clock ✗", Title: "app-clock-expiry", Family: famExpiry,
		Summary:  "S1 judging expiry by the application's clock, one simulated node running ahead. Negative control: expected to steal valid holds.",
		Strategy: Guarded, Isolation: ports.ReadCommitted, Layout: SeatRows, AppClock: true,
		NegativeControl: "expiry judged by an application clock running ahead (injected skew): valid holds are taken before they expire",
		ControlFires:    []string{"lifecycle"},
	},
	{
		ID: "e1_sweeper_expiry", Short: "E1 sweeper", Title: "sweeper-expiry", Family: famExpiry,
		Summary:  "Expiry takes effect only when a sweeper releases the hold; holds require status = 'available'.",
		Strategy: Guarded, Isolation: ports.ReadCommitted, Layout: SeatRows, ExpiryOnSweeper: true,
	},
	{
		ID: "e2_cart_expiry", Short: "E2 cart", Title: "cart-expiry", Family: famExpiry,
		Summary:  "Expiry stored once on a cart row; seat rows point to the cart and steals test its expiry through a subquery.",
		Strategy: Cart, Isolation: ports.ReadCommitted, Layout: CartRows,
	},
	{
		ID: "k0_naive_confirm", Short: "K0 naive ✗", Title: "naive-confirm", Family: famCheckout,
		Summary:  "S1 with a confirmation that marks the seats sold without checking the hold. Negative control: expected to sell after expiry.",
		Strategy: Guarded, Isolation: ports.ReadCommitted, Layout: SeatRows,
		NegativeControl: "confirms by seat, trusting a pre-payment check made before the hold expired",
		ControlFires:    []string{"lifecycle"},
	},
	{
		ID: "k1_payment_window", Short: "K1 window", Title: "payment-window", Family: famCheckout,
		Summary:  "S1 plus a checkout that extends a valid hold once by a bounded payment window.",
		Strategy: Guarded, Isolation: ports.ReadCommitted, Layout: SeatRows, PaymentWindow: true,
	},
	{
		ID: "l1_claim_rows", Short: "L1 claims", Title: "claim-rows", Family: famLayout,
		Summary:  "No per-event seat rows: a claim row is inserted on hold, the primary key arbitrates, an upsert steals expired claims.",
		Strategy: Claim, Isolation: ports.ReadCommitted, Layout: ClaimRows,
	},
	{
		ID: "l2_section_document", Short: "L2 document", Title: "section-document", Family: famLayout,
		Summary:  "One JSONB document of claimed seats per event section, rewritten by version compare-and-set.",
		Strategy: Document, Isolation: ports.ReadCommitted, Layout: Documents,
	},
	{
		ID: "l3_section_sharded", Short: "L3 sharded", Title: "section-sharded", Family: famLayout,
		Summary:  "S1 with event_seat hash-sharded by (event_id, section_no). YugabyteDB only.",
		Strategy: Guarded, Isolation: ports.ReadCommitted, Layout: SeatRows, YBOnly: true,
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

func (d Design) controlFiresIn(experiment string) bool {
	for _, e := range d.ControlFires {
		if e == experiment {
			return true
		}
	}
	return false
}
