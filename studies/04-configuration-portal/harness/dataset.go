package main

import (
	"fmt"
	"math"
	"math/rand"
	"strings"
)

// The dataset is deterministic and shared. Every design and every engine loads
// the same logical configuration, so a difference between two results is a
// difference in the design and never in the data (methodology 6).
//
// It is also shaped like the real thing: skewed rather than uniform (a uniform
// key distribution hides both hot-key contention and hot-installation effects),
// long-tailed but bounded (an unbounded power law would turn the document designs
// into strawmen), and generated in arrival order.

type Definition struct {
	ID      int
	Name    string
	Version string
}

type Environment struct {
	ID   int
	Name string
}

type BusinessUnit struct {
	ID   int
	Name string
}

// Installation is one deployed instance of a product definition.
type Installation struct {
	ID           int64
	ProductDefID int
	EnvID        int
	BUID         int // 0 means no business unit
	DisplayName  string
}

// Entry is one configuration key and its value.
type Entry struct {
	Key   string
	Value string
}

// Dataset is the whole logical input, held in memory: the study's sizes are small
// enough that this is simpler and more auditable than a generator protocol, and
// the harness must be able to compute the expected answers independently of the
// database (methodology 5).
type Dataset struct {
	Seed          int64
	Definitions   []Definition
	Environments  []Environment
	BusinessUnits []BusinessUnit
	Installations []Installation
	// Entries is the complete initial configuration of each installed product,
	// in key order.
	Entries map[int64][]Entry

	Tier     int    // entries per installed product, for the exact tiers
	CardMode string // exact | constant-entries | constant-bytes | fixed-fleet | skewed
	Regime   string // small | large

	// Byte accounting, recorded rather than estimated. Key count is never
	// equated with document size; all four numbers are reported.
	KeyBytes        int
	ValueBytes      int
	TotalEntries    int
	SerializedBytes int64

	// counts, when CardMode is "skewed": entries per installation, bounded by Tier.
	counts map[int64]int
}

var sectionNames = []string{"db", "cache", "auth", "logging", "limits", "features"}
var entryKinds = []string{"pool_size", "ttl", "retry_limit", "timeout", "batch", "threshold"}

func definitionNames() []string {
	return []string{
		"ledger", "gateway", "scheduler", "ingest", "reporter",
		"identity", "billing", "catalog", "notifier", "search",
	}
}

func environmentNames() []string { return []string{"production", "staging", "development"} }

func businessUnitNames() []string { return []string{"payments", "identity", "reporting", "retail"} }

// makeKey is deterministic and section-prefixed: the section is the part before
// the first dot, which is what a section-sharded design would shard on.
func makeKey(i int) string {
	section := sectionNames[i%len(sectionNames)]
	kind := entryKinds[(i/len(sectionNames))%len(entryKinds)]
	return fmt.Sprintf("%s.%s_%03d", section, kind, i)
}

// makeValue produces a value of a fixed, recorded size in the requested regime.
// The two regimes exist because a configuration value is sometimes a scalar and
// sometimes a serialized document, and the document designs are expected to
// behave differently in each.
func makeValue(regime string, i int) string {
	if regime == "large" {
		// ~4 KiB of JSON. The padding is inside the value, so it is stored,
		// returned and rewritten exactly like a real serialized settings blob.
		const target = 4096
		head := fmt.Sprintf(`{"seq":%d,"enabled":true,"payload":"`, i)
		tail := `"}`
		pad := target - len(head) - len(tail)
		if pad < 0 {
			pad = 0
		}
		return head + strings.Repeat("a", pad) + tail
	}
	// ~96 bytes: a short scalar setting.
	const target = 96
	head := fmt.Sprintf("v%08d-", i)
	pad := target - len(head)
	if pad < 0 {
		pad = 0
	}
	return head + strings.Repeat("s", pad)
}

// zipfIndex draws an index in [0,n) from a bounded power law. Uniform parents
// hide the effects this study is about; an unbounded law would make them into
// strawmen. `cap` bounds the ratio between the hottest and the coldest draw.
func zipfIndex(r *rand.Rand, n int, s float64) int {
	if n <= 1 {
		return 0
	}
	h := 0.0
	for i := 1; i <= n; i++ {
		h += 1 / math.Pow(float64(i), s)
	}
	u := r.Float64() * h
	acc := 0.0
	for i := 1; i <= n; i++ {
		acc += 1 / math.Pow(float64(i), s)
		if u <= acc {
			return i - 1
		}
	}
	return n - 1
}

// BuildDataset constructs the logical dataset. fleet is the requested number of
// installed products; the cardinality mode may reduce it so that total volume is
// held approximately constant while entries per installation changes
// (methodology 8a: change the shape, hold the volume).
func BuildDataset(seed int64, fleet, tier int, regime, cardMode string) *Dataset {
	r := rand.New(rand.NewSource(seed))
	d := &Dataset{
		Seed:        seed,
		Tier:        tier,
		CardMode:    cardMode,
		Regime:      regime,
		Entries:     map[int64][]Entry{},
		counts:      map[int64]int{},
		KeyBytes:    len(makeKey(0)),
		ValueBytes:  len(makeValue(regime, 0)),
	}

	for i, n := range definitionNames() {
		d.Definitions = append(d.Definitions, Definition{ID: i + 1, Name: n, Version: fmt.Sprintf("%d.%d", 1+i%3, i%5)})
	}
	for i, n := range environmentNames() {
		d.Environments = append(d.Environments, Environment{ID: i + 1, Name: n})
	}
	for i, n := range businessUnitNames() {
		d.BusinessUnits = append(d.BusinessUnits, BusinessUnit{ID: i + 1, Name: n})
	}

	// The reference volume is 60 entries on a 60-installation fleet: the study's
	// mandatory maximum, used as the unit of "one constant".
	const refTier = 60
	const refFleet = 60
	avgEntry := int64(d.KeyBytes + d.ValueBytes)

	effectiveFleet := fleet
	switch cardMode {
	case "constant-entries":
		// Total logical entries held constant while entries per installation changes.
		if tier > 0 {
			effectiveFleet = (refTier * refFleet) / tier
		}
	case "constant-bytes":
		// Total serialized bytes held constant instead of entry count.
		if tier > 0 && avgEntry > 0 {
			effectiveFleet = int((int64(refTier*refFleet) * avgEntry) / (int64(tier) * avgEntry))
		}
	}
	if effectiveFleet < 1 {
		effectiveFleet = 1
	}

	// Installations: one per (definition, environment, business unit) combination,
	// cycled deterministically until the fleet is reached. Not every installation
	// has a business unit -- it is an optional dimension, and INV-11 has to cope
	// with its absence.
	for i := 0; i < effectiveFleet; i++ {
		ip := Installation{
			ID:           int64(i + 1),
			ProductDefID: d.Definitions[i%len(d.Definitions)].ID,
			EnvID:        d.Environments[(i/len(d.Definitions))%len(d.Environments)].ID,
			DisplayName:  fmt.Sprintf("%s-%d", d.Definitions[i%len(d.Definitions)].Name, i+1),
		}
		if i%3 != 2 {
			ip.BUID = d.BusinessUnits[i%len(d.BusinessUnits)].ID
		}
		d.Installations = append(d.Installations, ip)
	}

	for _, ip := range d.Installations {
		n := tier
		if cardMode == "skewed" {
			// A realistic distribution bounded by the configured maximum: the
			// tier is a ceiling, not a uniform count.
			n = 1 + zipfIndex(r, max(tier, 1), 1.0)
		}
		d.counts[ip.ID] = n
		entries := make([]Entry, 0, n)
		for k := 0; k < n; k++ {
			key := makeKey(k)
			val := makeValue(regime, int(ip.ID)*1000+k)
			entries = append(entries, Entry{Key: key, Value: val})
		}
		d.Entries[ip.ID] = entries
		d.TotalEntries += len(entries)
		for _, e := range entries {
			d.SerializedBytes += int64(len(e.Key) + len(e.Value))
		}
	}
	return d
}

// EntryCount is the configured number of entries for one installation.
func (d *Dataset) EntryCount(id int64) int { return d.counts[id] }

// PickInstallation draws an installed product for a read. skew selects a bounded
// power law so a small set of installations is genuinely hot: 80% of reads target
// the hot fifth, which is what makes contention and cache effects visible.
func (d *Dataset) PickInstallation(r *rand.Rand, hot bool) Installation {
	n := len(d.Installations)
	if n == 0 {
		return Installation{}
	}
	if hot {
		hotN := max(n/5, 1)
		return d.Installations[zipfIndex(r, hotN, 1.1)]
	}
	return d.Installations[zipfIndex(r, n, 0.8)]
}

// PickKey draws a key index for one installation. Key popularity is bounded so
// that hot keys exist without making the distribution degenerate.
func (d *Dataset) PickKey(r *rand.Rand, id int64) int {
	n := d.EntryCount(id)
	if n <= 0 {
		return 0
	}
	return zipfIndex(r, n, 1.1)
}

// SerializedSize is the logical size of one installation's configuration: the
// bytes a complete read would have to transfer, before any encoding overhead.
func (d *Dataset) SerializedSize(id int64) int64 {
	var n int64
	for _, e := range d.Entries[id] {
		n += int64(len(e.Key) + len(e.Value))
	}
	return n
}
