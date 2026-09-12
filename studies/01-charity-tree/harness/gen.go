package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

// ---------------------------------------------------------------------------
// Deterministic dataset generation.
//
// Every design and every engine loads the SAME logical dataset, produced from a
// fixed seed. That is what makes cross-design numbers comparable: a difference
// in a result is a difference in the design, never a difference in the data.
// ---------------------------------------------------------------------------

type Charity struct {
	ID        int64
	Name      string
	Country   string
	FoundedOn time.Time
}

type Person struct {
	ID        int64
	CharityID int64
	FullName  string
	Email     string
	JoinedAt  time.Time
}

type Donation struct {
	ID          int64
	PersonID    int64
	CharityID   int64
	AmountCents int64
	Currency    string
	DonatedAt   time.Time
	Note        *string
}

type Dataset struct {
	Scale string
	Seed  int64
	// MaxPerPerson is the history cap (0 = uncapped). See GenerateProfile.
	MaxPerPerson int
	Charities    []Charity
	People       []Person
	Donations    []Donation
	// DonationsByPerson[i] holds indexes into Donations for People[i]. Built for
	// the embedded design's loader, which needs one array per person.
	DonationsByPerson [][]int32
}

type Scale struct {
	Name      string
	People    int
	Charities int
}

// Donation totals follow from People via donationCountFor (mean ~22/person):
// tiny ~11k, small ~110k, medium ~1.0M, large ~4.0M donations.
var scales = map[string]Scale{
	"tiny":   {Name: "tiny", People: 500, Charities: 10},
	"small":  {Name: "small", People: 5000, Charities: 10},
	"medium": {Name: "medium", People: 45000, Charities: 10},
	"large":  {Name: "large", People: 180000, Charities: 10},
}

var (
	epochStart = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)
	epochEnd   = time.Date(2026, 9, 1, 0, 0, 0, 0, time.UTC)
	joinEnd    = time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)
)

// charityWeights makes the charities deliberately uneven in size. Real portfolios
// never are uniform, and a uniform split would hide both the effect of skew on
// charity-scoped queries and the contention on the largest charity's rollup row.
var charityWeights = []int{30, 20, 14, 10, 8, 6, 5, 3, 2, 2}

var charityNames = []string{
	"Open Hands Foundation", "Clearwater Trust", "Northern Lights Relief",
	"Sunrise Education Fund", "Harbour Health Alliance", "Green Valley Initiative",
	"Stonebridge Shelter", "Meridian Arts Council", "Silver Pines Hospice",
	"Kestrel Wildlife Fund",
}

var countries = []string{"BR", "PT", "US", "GB", "DE", "FR", "ES", "NL", "CA", "IE"}

var firstNames = []string{
	"Ana", "Bruno", "Carla", "Diego", "Elena", "Felipe", "Greta", "Hugo", "Ines", "Joao",
	"Katia", "Lucas", "Marta", "Nuno", "Olivia", "Pedro", "Quentin", "Rita", "Sofia", "Tiago",
	"Ursula", "Victor", "Wanda", "Xavier", "Yara", "Zeca",
}

var lastNames = []string{
	"Alves", "Barbosa", "Costa", "Duarte", "Esteves", "Ferreira", "Gomes", "Henriques",
	"Iglesias", "Jardim", "Krause", "Lima", "Moreira", "Neves", "Oliveira", "Pereira",
	"Quintela", "Ribeiro", "Santos", "Teixeira", "Vasquez", "Whitaker",
}

var currencies = []string{"USD", "USD", "USD", "USD", "USD", "USD", "EUR", "EUR", "GBP", "BRL"}

var noteTemplates = []string{
	"monthly pledge", "in memory of a friend", "matched by employer",
	"gala pledge", "birthday fundraiser", "recurring gift", "one-off gift",
	"campaign response", "legacy contribution", "emergency appeal",
}

// donationCountFor draws a per-person lifetime donation count. The buckets give a
// long-tailed but BOUNDED distribution (mean ~22, max 500). An unbounded Zipf
// would put six-figure arrays into a single JSONB document and turn D6 into a
// strawman rather than a fair comparison.
func donationCountFor(r *rand.Rand) int {
	switch u := r.Float64(); {
	case u < 0.50:
		return 1 + r.Intn(5)
	case u < 0.85:
		return 6 + r.Intn(15)
	case u < 0.97:
		return 21 + r.Intn(80)
	default:
		return 101 + r.Intn(400)
	}
}

func amountFor(r *rand.Rand) int64 {
	switch u := r.Float64(); {
	case u < 0.60:
		return int64(500 + r.Intn(4500))
	case u < 0.85:
		return int64(5000 + r.Intn(20000))
	case u < 0.97:
		return int64(25000 + r.Intn(75000))
	default:
		return int64(100000 + r.Intn(4900000))
	}
}

func randTime(r *rand.Rand, from, to time.Time) time.Time {
	span := to.Sub(from)
	if span <= 0 {
		return from
	}
	return from.Add(time.Duration(r.Int63n(int64(span)))).Truncate(time.Millisecond)
}

// Generate builds the dataset for a scale. Donations are emitted in chronological
// order and numbered sequentially, mirroring how an append-only donation stream
// actually arrives: donation_id and donated_at are correlated, exactly as they
// would be in production. That correlation matters -- it is why a time-ordered
// index stays physically clustered, and a shuffled generator would quietly make
// every time-range read look worse than it is.
func Generate(scaleName string, seed int64) (*Dataset, error) {
	return GenerateProfile(scaleName, seed, Profile{})
}

// GenerateProfile is Generate with an optional cap on donations per person.
//
// Why a cap exists at all: several designs carry a limitation that only bites
// when one parent accumulates many children. D6 rewrites a donor's whole array on
// every append, so its write cost is O(history). D9's 20-element cache only
// answers "recent donations" in full if nobody asks for more than it holds. D4's
// delete trigger recomputes MIN/MAX over the donor's surviving history. In a
// domain where history is naturally bounded -- a subscription with a fixed term,
// a ticket with a capped number of comments, an order with a limited number of
// lines -- those limitations may simply not apply, and a design that loses the
// unbounded benchmark could be the right answer. The cap is how the study
// measures that regime instead of only asserting it.
//
// Holding VOLUME constant: capping per-person history on its own would also
// shrink the donation table, and every design would get faster for a reason that
// has nothing to do with the cap. So a capped dataset keeps adding donors until it
// reaches the same donation count as the uncapped dataset at the same scale and
// seed. What changes is the SHAPE of the tree -- more, shallower parents -- not
// the number of leaves.
func GenerateProfile(scaleName string, seed int64, prof Profile) (*Dataset, error) {
	sc, ok := scales[scaleName]
	if !ok {
		return nil, fmt.Errorf("unknown scale %q", scaleName)
	}
	if prof.Charities > 0 {
		sc.Charities = prof.Charities
	}
	maxPerPerson := prof.MaxPerPerson
	targetDonations := 0
	if maxPerPerson > 0 {
		base, err := GenerateProfile(scaleName, seed, Profile{Charities: prof.Charities})
		if err != nil {
			return nil, err
		}
		targetDonations = len(base.Donations)
	}
	r := rand.New(rand.NewSource(seed))
	ds := &Dataset{Scale: scaleName, Seed: seed, MaxPerPerson: maxPerPerson}

	for i := 0; i < sc.Charities; i++ {
		ds.Charities = append(ds.Charities, Charity{
			ID:        int64(i + 1),
			Name:      charityNames[i%len(charityNames)],
			Country:   countries[i%len(countries)],
			FoundedOn: epochStart.AddDate(-r.Intn(40), -r.Intn(12), -r.Intn(28)),
		})
	}

	totalWeight := 0
	for i := 0; i < sc.Charities; i++ {
		totalWeight += charityWeights[i%len(charityWeights)]
	}
	pickCharity := func() int64 {
		x := r.Intn(totalWeight)
		for i := 0; i < sc.Charities; i++ {
			w := charityWeights[i%len(charityWeights)]
			if x < w {
				return int64(i + 1)
			}
			x -= w
		}
		return int64(sc.Charities)
	}

	ds.People = make([]Person, 0, sc.People)
	ds.DonationsByPerson = make([][]int32, 0, sc.People)

	type pending struct {
		personIdx int32
		at        time.Time
	}
	var all []pending

	// Uncapped: exactly sc.People donors. Capped: keep adding donors until the
	// donation count matches the uncapped dataset, so only the shape changes.
	more := func(i int) bool {
		if maxPerPerson > 0 {
			return len(all) < targetDonations
		}
		return i < sc.People
	}

	for i := 0; more(i); i++ {
		ds.DonationsByPerson = append(ds.DonationsByPerson, nil)
		pid := int64(i + 1)
		joined := randTime(r, epochStart, joinEnd)
		ds.People = append(ds.People, Person{
			ID:        pid,
			CharityID: pickCharity(),
			FullName:  firstNames[r.Intn(len(firstNames))] + " " + lastNames[r.Intn(len(lastNames))],
			Email:     fmt.Sprintf("donor%d@example.org", pid),
			JoinedAt:  joined,
		})
		n := donationCountFor(r)
		// Clamping (rather than redrawing) keeps the same long-tailed draw and
		// simply truncates the tail at the cap -- heavy donors become donors at
		// the limit, which is what a bounded domain looks like in practice.
		if maxPerPerson > 0 && n > maxPerPerson {
			n = maxPerPerson
		}
		for j := 0; j < n; j++ {
			all = append(all, pending{personIdx: int32(i), at: randTime(r, joined, epochEnd)})
		}
	}

	sort.Slice(all, func(a, b int) bool {
		if all[a].at.Equal(all[b].at) {
			return all[a].personIdx < all[b].personIdx
		}
		return all[a].at.Before(all[b].at)
	})

	ds.Donations = make([]Donation, 0, len(all))
	for i, pd := range all {
		person := ds.People[pd.personIdx]
		var note *string
		if r.Float64() < 0.30 {
			s := noteTemplates[r.Intn(len(noteTemplates))]
			note = &s
		}
		ds.Donations = append(ds.Donations, Donation{
			ID:          int64(i + 1),
			PersonID:    person.ID,
			CharityID:   person.CharityID,
			AmountCents: amountFor(r),
			Currency:    currencies[r.Intn(len(currencies))],
			DonatedAt:   pd.at,
			Note:        note,
		})
		ds.DonationsByPerson[pd.personIdx] = append(ds.DonationsByPerson[pd.personIdx], int32(i))
	}
	return ds, nil
}

// Summary describes the generated data for the run manifest, so a report can
// state exactly what was measured without re-deriving it.
func (ds *Dataset) Summary() map[string]any {
	maxPer, withZero := 0, 0
	var total int64
	for i := range ds.People {
		n := len(ds.DonationsByPerson[i])
		if n > maxPer {
			maxPer = n
		}
		if n == 0 {
			withZero++
		}
	}
	for i := range ds.Donations {
		total += ds.Donations[i].AmountCents
	}
	perCharity := map[string]int{}
	for i := range ds.People {
		perCharity[fmt.Sprint(ds.People[i].CharityID)]++
	}
	return map[string]any{
		"scale":                   ds.Scale,
		"max_per_person_cap":      ds.MaxPerPerson,
		"profile":                 profileName(ds.MaxPerPerson, len(ds.Charities)),
		"seed":                    ds.Seed,
		"charities":               len(ds.Charities),
		"people":                  len(ds.People),
		"donations":               len(ds.Donations),
		"avg_donations_person":    float64(len(ds.Donations)) / float64(len(ds.People)),
		"max_donations_person":    maxPer,
		"people_without_donation": withZero,
		"total_donated_cents":     total,
		"people_per_charity":      perCharity,
	}
}

// profileName labels a dataset shape for results and reports.
func profileName(maxPerPerson, charities int) string {
	name := "unbounded"
	if maxPerPerson > 0 {
		name = fmt.Sprintf("cap%d", maxPerPerson)
	}
	if charities != 10 {
		name += fmt.Sprintf("-charities%d", charities)
	}
	return name
}

// Profile describes the SHAPE of the generated tree, as opposed to its size.
//
// Both knobs exist to test the regime in which some design's limitation stops
// mattering, at each level of the tree:
//
//	MaxPerPerson  caps children per person (donations per donor). Targets the
//	              designs whose cost grows with one parent's history: D6's
//	              whole-array rewrite, D9's cache coverage, D4's delete-time
//	              MIN/MAX recompute, D7's per-person tablet.
//	Charities     widens the top level (persons per charity shrink as it grows).
//	              Targets the designs whose cost grows with one charity's size:
//	              D4/D5's single hot rollup row per charity, and the charity-
//	              scoped queries that D6 must answer by unnesting every array in
//	              the charity.
//
// Either way the number of donations is held constant, so a faster result means
// the shape helped -- not that the table got smaller.
type Profile struct {
	MaxPerPerson int
	Charities    int
}
