package main

import (
	"math/rand"
	"sort"
	"time"
)

func nowMS() int64 { return time.Now().UnixNano() / int64(time.Millisecond) }

// The oracle: an independent model of every committed state of every key, built in
// Go from the generated dataset and the operations the client saw acknowledged.
//
// It is the reason this study can report a wrong-read rate without validating each
// cache hit with a hidden database read. Two things make it exact:
//
//  1. Every committed state of a key is recorded as the CONTENT HASH of that state,
//     in publication order. Freshness is therefore a statement about the data, not
//     about a version number two components might disagree about -- which matters
//     because the legacy model has no version number at all.
//  2. A read captures its requirement at the instant it begins: the hash of the
//     latest acknowledged write to that key before the read started. A value older
//     than that is a wrong read; a value newer than that is not (a read may
//     legitimately see a write that committed while it was in flight).
//
// The oracle is measurement metadata. The cache adapter never reads it, and no
// database schema is changed to give it a token: on the legacy model the oracle
// "version" exists only inside this file.

type stateStamp struct {
	Seq  int64
	Hash string
	// AtMS is when this state became the committed one, used for the
	// time-to-freshness measurement.
	AtMS int64
}

type oraclePerson struct {
	charityID int64
	joinedAt  string
	fullName  string
	email     string
	donations map[int64]Donation
	seq       int64
	curHash   string
	history   []stateStamp
}

type oracle struct {
	ds     *Dataset
	people map[int64]*oraclePerson
	mu     chanMutex
}

// chanMutex is a mutex that also gives the oracle a single lock-ordering point.
// A plain sync.Mutex would do; this keeps the locking explicit at every call site
// that spans more than one field read.
type chanMutex struct{ ch chan struct{} }

func newChanMutex() chanMutex { return chanMutex{ch: make(chan struct{}, 1)} }
func (m chanMutex) Lock()     { m.ch <- struct{}{} }
func (m chanMutex) Unlock()   { <-m.ch }

func newOracle(ds *Dataset) *oracle {
	o := &oracle{ds: ds, people: make(map[int64]*oraclePerson, len(ds.People)), mu: newChanMutex()}
	o.mu.Lock()
	defer o.mu.Unlock()
	for _, p := range ds.People {
		op := &oraclePerson{
			charityID: p.CharityID,
			joinedAt:  p.JoinedAt.UTC().Format(time.RFC3339Nano),
			fullName:  p.FullName,
			email:     p.Email,
			donations: map[int64]Donation{},
		}
		for _, d := range ds.Donations[p.ID] {
			op.donations[d.ID] = d
		}
		op.curHash = o.hashLocked(p.ID, op)
		op.history = []stateStamp{{Seq: 0, Hash: op.curHash, AtMS: nowMS()}}
		o.people[p.ID] = op
	}
	return o
}

// recentDescLocked returns the newest 20 donations in the study's deterministic
// order: (donated_at DESC, donation_id DESC).
func recentDescLocked(m map[int64]Donation) []Donation {
	out := make([]Donation, 0, len(m))
	for _, d := range m {
		out = append(out, d)
	}
	sort.Slice(out, func(i, j int) bool {
		if !out[i].DonatedAt.Equal(out[j].DonatedAt) {
			return out[i].DonatedAt.After(out[j].DonatedAt)
		}
		return out[i].ID > out[j].ID
	})
	if len(out) > 20 {
		out = out[:20]
	}
	return out
}

func (o *oracle) contentLocked(personID int64, op *oraclePerson) PortalContent {
	c, _ := o.ds.Charity(op.charityID)
	var total int64
	for _, d := range op.donations {
		total += d.AmountCents
	}
	recent := recentDescLocked(op.donations)
	recs := make([]PortalDonation, 0, len(recent))
	for _, d := range recent {
		recs = append(recs, PortalDonation{
			ID:          d.ID,
			AmountCents: d.AmountCents,
			Currency:    d.Currency,
			DonatedAt:   d.DonatedAt.UTC().Format(time.RFC3339Nano),
			Note:        d.Note2(),
		})
	}
	return PortalContent{
		PersonID:           personID,
		FullName:           op.fullName,
		Email:              op.email,
		JoinedAt:           op.joinedAt,
		CharityID:          op.charityID,
		CharityName:        c.Name,
		CharityCountry:     c.Country,
		DonationCount:      int64(len(op.donations)),
		DonationTotalCents: total,
		Recent:             recs,
	}
}

func (o *oracle) hashLocked(personID int64, op *oraclePerson) string {
	return o.contentLocked(personID, op).ContentHash()
}

// bumpLocked records that this key reached a new committed state.
func (o *oracle) bumpLocked(personID int64) {
	op := o.people[personID]
	if op == nil {
		return
	}
	op.seq++
	op.curHash = o.hashLocked(personID, op)
	op.history = append(op.history, stateStamp{Seq: op.seq, Hash: op.curHash, AtMS: nowMS()})
}

// ---------------------------------------------------------------- mutations
//
// The harness calls exactly one of these after a mutation's database transaction
// has committed, with the values it committed. A mutation that rolled back, or
// whose COMMIT outcome was ambiguous and which was therefore rolled back, is
// applied through RollbackMutation instead.

func (o *oracle) ApplyInsert(personID int64, d Donation) {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return
	}
	op.donations[d.ID] = d
	o.bumpLocked(personID)
}

func (o *oracle) ApplyCorrect(personID int64, donationID, amountCents int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return
	}
	d, ok := op.donations[donationID]
	if !ok {
		return
	}
	d.AmountCents = amountCents
	op.donations[donationID] = d
	o.bumpLocked(personID)
}

func (o *oracle) ApplyDelete(personID int64, donationID int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return
	}
	delete(op.donations, donationID)
	o.bumpLocked(personID)
}

func (o *oracle) ApplyPersonUpdate(personID int64, fullName, email string) {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return
	}
	op.fullName = fullName
	op.email = email
	o.bumpLocked(personID)
}

// ApplyReassign moves a donation between two people. Both keys reach a new
// committed state, so both are bumped -- which is why reassignment is the mutation
// that breaks caches which only ever think about one key.
func (o *oracle) ApplyReassign(donationID, fromPerson, toPerson, newCharityID int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	from := o.people[fromPerson]
	to := o.people[toPerson]
	if from == nil || to == nil {
		return
	}
	d, ok := from.donations[donationID]
	if !ok {
		return
	}
	delete(from.donations, donationID)
	d.PersonID = toPerson
	d.CharityID = newCharityID
	to.donations[donationID] = d
	o.bumpLocked(fromPerson)
	o.bumpLocked(toPerson)
}

// ExtBump records an acknowledged EXTERNAL write: one that bypassed the cache
// adapter and committed to the database. The adapter cannot see it; the oracle
// must, or the strict contract would be judged against an incomplete ledger.
func (o *oracle) ExtBump(personID int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.bumpLocked(personID)
}

// ---------------------------------------------------------------- reads

// Required returns the freshness requirement for a read that begins now: the hash
// and sequence of the latest acknowledged write to this key.
func (o *oracle) Required(personID int64) (hash string, seq int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return "", 0
	}
	return op.curHash, op.seq
}

func (o *oracle) CurrentSeq(personID int64) int64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	if op := o.people[personID]; op != nil {
		return op.seq
	}
	return 0
}

// Classify turns a returned content hash into the study's read vocabulary.
//
//	fresh      the value is the required committed state (or the content is
//	           identical to it, whatever publication it came from)
//	ahead      the value is a LATER committed state than the read required, which
//	           happens legitimately when a write commits while the read is in
//	           flight. It is not a wrong read.
//	stale      the value is an EARLIER committed state than the read required --
//	           the violation the strict contract forbids
//	impossible the value corresponds to no committed state this key ever had
func (o *oracle) Classify(personID int64, returnedHash string, reqHash string, reqSeq int64) (kind string, returnedSeq int64, behind int64, atMS int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return KindImpossible, -1, 0, 0
	}
	if returnedHash == reqHash {
		return KindFresh, reqSeq, 0, 0
	}
	for i := len(op.history) - 1; i >= 0; i-- {
		if op.history[i].Hash == returnedHash {
			s := op.history[i].Seq
			if s > reqSeq {
				return KindAhead, s, 0, op.history[i].AtMS
			}
			return KindStale, s, reqSeq - s, op.history[i].AtMS
		}
	}
	return KindImpossible, -1, 0, 0
}

// ---------------------------------------------------------------- audit views

// ExpectedContent is the oracle's view of a key's committed state, used by the
// audit phase to compare the database against something computed independently.
func (o *oracle) ExpectedContent(personID int64) (PortalContent, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return PortalContent{}, false
	}
	return o.contentLocked(personID, op), true
}

func (o *oracle) ExpectedDonationIDs(personID int64) []int64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return nil
	}
	out := make([]int64, 0, len(op.donations))
	for id := range op.donations {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func (o *oracle) ExpectedRecentIDs(personID int64) []int64 {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil {
		return nil
	}
	recent := recentDescLocked(op.donations)
	out := make([]int64, 0, len(recent))
	for _, d := range recent {
		out = append(out, d.ID)
	}
	return out
}

func (o *oracle) Person(personID int64) (Person, bool) { return o.ds.Person(personID) }

// PickDonation chooses one of a donor's CURRENT donations from the oracle's state.
// The workload uses it to decide what to correct, delete or reassign, so the
// operation always names a row the ledger believes exists.
func (o *oracle) PickDonation(personID int64, r *rand.Rand) (Donation, bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	op := o.people[personID]
	if op == nil || len(op.donations) == 0 {
		return Donation{}, false
	}
	idx := r.Intn(len(op.donations))
	i := 0
	for _, d := range op.donations {
		if i == idx {
			return d, true
		}
		i++
	}
	return Donation{}, false
}
