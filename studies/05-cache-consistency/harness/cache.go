package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"sync"
	"time"
)

// The cached representation, the canonical encoder that hashes it, the expiry
// policy and the store interface every backend implements.
//
// The content hash is the study's unit of truth. The oracle judges a read fresh or
// stale by comparing the hash the cache returned with the hash of the committed
// state the ledger says the key should hold. That choice matters: it makes
// freshness a statement about the DATA, not about a version number the cache and
// the database might disagree about, and it is the same test for the owned model
// (which has a version token) and the legacy model (which does not).

// hardTTL is the study's exact five-minute hard expiry. It is a constant, and it is
// never shortened to make a run finish faster; a scenario whose entries would
// outlive a shorter TTL is simply not exercising expiry, which the report says.
const (
	hardTTLSeconds = 300.0
	hardTTL        = 300 * time.Second
)

// PortalDonation is one element of the cached recent slice.
type PortalDonation struct {
	ID          int64  `json:"id"`
	AmountCents int64  `json:"amt"`
	Currency    string `json:"cur"`
	DonatedAt   string `json:"at"`
	Note        string `json:"note"`
}

// PortalContent is the coherent committed view the cache holds for one donor.
//
// It carries the donor identity and mutable metadata, the charity identity, the
// two aggregates and the newest twenty donations. It deliberately does NOT carry a
// version: the version is entry metadata, so two entries that describe the same
// data but were published at different times hash the same.
type PortalContent struct {
	PersonID           int64            `json:"person_id"`
	FullName           string           `json:"full_name"`
	Email              string           `json:"email"`
	JoinedAt           string           `json:"joined_at"`
	CharityID          int64            `json:"charity_id"`
	CharityName        string           `json:"charity_name"`
	CharityCountry     string           `json:"charity_country"`
	DonationCount      int64            `json:"donation_count"`
	DonationTotalCents int64            `json:"donation_total_cents"`
	Recent             []PortalDonation `json:"recent"`
}

// Normalize puts the recent slice into the study's deterministic order --
// (donated_at DESC, donation_id DESC) -- so a design that returns the same
// donations in a different order is not called wrong. Ordering is part of the
// cached representation's contract, and the gate checks it separately; the hash
// only has to compare membership and values.
func (c *PortalContent) Normalize() {
	if c.Recent == nil {
		c.Recent = []PortalDonation{}
	}
	sort.SliceStable(c.Recent, func(i, j int) bool {
		if c.Recent[i].DonatedAt != c.Recent[j].DonatedAt {
			return c.Recent[i].DonatedAt > c.Recent[j].DonatedAt
		}
		return c.Recent[i].ID > c.Recent[j].ID
	})
}

// Canonical is the one encoder every design's payload goes through, so the hash
// compares data rather than JSON renderers.
func (c PortalContent) Canonical() []byte {
	c.Normalize()
	b, err := json.Marshal(c)
	if err != nil {
		// The struct is ints, strings and slices; marshalling cannot fail. A panic
		// here is better than a silently empty hash.
		panic("canonical portal content: " + err.Error())
	}
	return b
}

// ContentHash is the identity of one committed state of one key.
func (c PortalContent) ContentHash() string {
	sum := sha256.Sum256(c.Canonical())
	return hex.EncodeToString(sum[:])
}

// Entry is one cache value: the committed content, plus the metadata the cache
// policy needs. The hard-expiry deadline is stored rather than inferred from the
// backend's own TTL so that memory and Redis implement the same semantics and the
// probabilistic early-expiry draw reads the same fields on both.
type Entry struct {
	// Fence is the key's invalidation fence at the moment this entry was read.
	// Publication is REFUSED unless the key's fence is still this value, which is
	// what stops a fill that began before an invalidation from republishing the
	// state it read. That race is real and is not closed by publishing only after a
	// commit: the reader's snapshot was taken before the write, so its content is a
	// committed but superseded state.
	Fence int64 `json:"f"`
	// Seq is the adapter's monotonic publication order for this key. It is used
	// for version fencing: an expired lease holder must not overwrite a newer
	// committed version. On the owned model Version (below) is authoritative.
	Seq int64 `json:"q"`
	// Version is the database cache_version of the state this entry holds, read in
	// the same snapshot as the content. Zero on the legacy model, which has no
	// token -- the harness's oracle generation is measurement metadata and is never
	// stored here as if the database had supplied it.
	Version int64 `json:"v"`
	// CreatedMS and ExpiresMS are wall-clock milliseconds.
	CreatedMS int64         `json:"c"`
	ExpiresMS int64         `json:"x"`
	Hash      string        `json:"h"`
	Content   PortalContent `json:"p"`
}

func newEntry(c PortalContent, version, seq, nowMS int64) *Entry {
	c.Normalize()
	return &Entry{
		Seq:       seq,
		Version:   version,
		CreatedMS: nowMS,
		ExpiresMS: nowMS + int64(hardTTL/time.Millisecond),
		Hash:      c.ContentHash(),
		Content:   c,
	}
}

// Size is the resident size this entry is charged for, key included. The byte
// bound in memory and the Redis accounting both use it, so "capacity" means the
// same thing on both backends.
func (e *Entry) Size(key string) int64 {
	b, err := json.Marshal(e)
	if err != nil {
		return int64(len(key))
	}
	return int64(len(key) + len(b))
}

func (e *Entry) AgeMS(nowMS int64) int64 { return nowMS - e.CreatedMS }

func (e *Entry) TimeRemainingS(nowMS int64) float64 {
	rem := e.ExpiresMS - nowMS
	if rem < 0 {
		return 0
	}
	return float64(rem) / 1000.0
}

func (e *Entry) HardExpired(nowMS int64) bool { return nowMS >= e.ExpiresMS }

// ShouldExpireEarly is the study's probabilistic early expiration, exactly as the
// protocol writes it:
//
//	time_remaining = max(0, hard_expiry - now)
//	p_expire       = clamp(1 - time_remaining / 300 s, 0, 1)
//
// The draw comes from the trial's recorded random source, so the same seed
// reproduces the same decisions. A pure function on purpose: the unit tests pin it
// at ages 0, 75, 150, 225 and 300 seconds with a fake clock, and the harness
// additionally reports the observed expiration rate by age bucket so the
// implementation is checked against data as well as against arithmetic.
func ShouldExpireEarly(remainingS, ttlS, draw float64) bool {
	if ttlS <= 0 {
		return true
	}
	p := 1 - remainingS/ttlS
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	return draw < p
}

// AgeBucket labels an observation for the reported expiration rate. The buckets are
// the five ages the unit tests pin, plus the two ends.
func AgeBucket(ageMS int64) string {
	switch {
	case ageMS < 0:
		return "<0s"
	case ageMS < 75_000:
		return "0-75s"
	case ageMS < 150_000:
		return "75-150s"
	case ageMS < 225_000:
		return "150-225s"
	case ageMS < 300_000:
		return "225-300s"
	default:
		return ">=300s"
	}
}

// EntryStore is what a cache backend must provide. It is deliberately small: the
// study's cache policy lives in the adapter, above this line, so memory and Redis
// implement the SAME lease, fencing and expiry semantics and their comparison is
// about the backend and not about two different policy implementations.
type EntryStore interface {
	Kind() string
	Get(ctx context.Context, key string) (*Entry, bool, error)
	// Put publishes with BOTH fences. It returns false when the stored entry carries
	// a higher Seq (an expired lease holder must not overwrite a newer committed
	// version) or when the key's fence has moved since the entry was read (a fill
	// that began before an invalidation must not republish what it read).
	Put(ctx context.Context, key string, e *Entry, ttl time.Duration) (bool, error)
	// Fence advances the key's invalidation fence and removes any value. A strict
	// writer calls it BEFORE the authoritative mutation; a relaxed writer calls it
	// after the commit, as its best-effort invalidation. It is one operation because
	// "remove the value" and "refuse anything read before now" are one intent.
	Fence(ctx context.Context, key string) (int64, error)
	// FenceOf reads the current fence, 0 when the key has never been invalidated.
	FenceOf(ctx context.Context, key string) (int64, error)
	Delete(ctx context.Context, key string) error
	// Flush clears the whole cache. Used by the cold-cache and stampede phases and
	// by the fault injections; never by a steady-state read.
	Flush(ctx context.Context) error
	// TryLease takes the per-key fill lease with a unique token. It is atomic.
	TryLease(ctx context.Context, key, token string, ttl time.Duration) (bool, error)
	// ReleaseLease releases only if the token still matches, so a holder whose
	// lease expired cannot delete a successor's lease.
	ReleaseLease(ctx context.Context, key, token string) error
	Close(ctx context.Context) error
	// Stats is backend-level accounting: resident bytes, item count, evictions.
	Stats() map[string]int64
	// Info is what the backend reports about itself: Redis INFO fields, or the
	// in-process allocator's numbers.
	Info(ctx context.Context) map[string]string
}

// Instance is one logical application instance. Its store and its local generation
// counter are its own: with the memory backend, three instances have three
// independent LRUs, and NO shared memory. That is the property the multi-instance
// phase exists to expose, and it would be silently destroyed if the instances
// shared the harness's oracle.
type Instance struct {
	ID    int
	Store EntryStore

	gen map[int64]int64
	gm  sync.Mutex
}

func newInstance(id int, store EntryStore) *Instance {
	return &Instance{ID: id, Store: store, gen: map[int64]int64{}}
}

// BumpGen advances this instance's own view of a key's generation. Only writes
// that pass through THIS instance increment it. The legacy model has no database
// token, so this counter is all the instance knows about its own writes -- which is
// exactly why an instance cannot prove freshness for a peer's write.
func (in *Instance) BumpGen(keyID int64) int64 {
	in.gm.Lock()
	defer in.gm.Unlock()
	in.gen[keyID]++
	return in.gen[keyID]
}

func (in *Instance) Gen(keyID int64) int64 {
	in.gm.Lock()
	defer in.gm.Unlock()
	return in.gen[keyID]
}
