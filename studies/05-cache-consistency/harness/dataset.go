package main

import (
	"fmt"
	"math"
	"math/rand"
	"time"
)

// The dataset is deterministic and shared. Every scenario and every engine loads
// the same logical charity tree, so a difference between two results is a
// difference in the scenario and never in the data (methodology 6).
//
// Two properties matter here specifically:
//
//   * SKEWED, NOT UNIFORM. A uniform donor draw would make every key equally hot,
//     which hides both the cache's reason to exist and the contention that breaks
//     it. Eighty per cent of reads target the hot fifth of the donors.
//   * LONG-TAILED BUT BOUNDED. An unbounded power law turns the reference
//     embedding design into a strawman: its rewrite cost would grow with history
//     forever. The tail is bounded, exactly as study 01 bounds it.
//
// Timestamps are truncated to MICROSECONDS at generation time. PostgreSQL stores
// timestamptz at microsecond resolution, so a payload containing nanoseconds would
// hash differently after a database round trip, and the study's whole correctness
// check is a content-hash comparison. Truncating where the data is born keeps that
// comparison honest instead of papering over it in the encoder.

type Charity struct {
	ID        int64
	Name      string
	Country   string
	FoundedOn string
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
	// Note is nil when the donor left none. The payload encoder folds SQL NULL and
	// JSON null onto the same empty value, which is what lets the embedded
	// reference design's JSONB array hash identically.
	Note *string
}

type Dataset struct {
	Seed  int64
	Scale string

	Charities []Charity
	People    []Person
	// Donations is keyed by person and ordered by (DonatedAt, ID) ascending --
	// the order an append-only stream produces.
	Donations map[int64][]Donation

	TotalDonations  int
	SerializedBytes int64
	MaxDonationID   int64
	// ValueBytes/KeyBytes are the logical sizes the cache capacity is sized
	// against. They are measured from the generated payload, not guessed.
	PayloadBytes int

	byPerson map[int64]int
	byChar   map[int64]Charity
}

var charityNames = []string{
	"alight", "beacon", "cedar", "dobro", "eastwind", "fenn", "girasol", "harbour",
	"isla", "juniper", "kelvin", "lumen", "mira", "northgate", "otter", "palm",
	"quarry", "riverbend", "solstice", "tundra", "umber", "verde", "willow", "xylem",
}

var countries = []string{"GB", "US", "DE", "NL", "ES", "IT", "FR", "SE"}

var currencies = []string{"USD", "EUR", "GBP", "JPY", "AUD"}

var givenNames = []string{
	"ada", "bruno", "chidi", "dara", "elif", "farid", "gita", "hana", "ivan",
	"juno", "kofi", "lena", "mira", "noor", "omar", "pia", "quinn", "rosa",
	"sam", "tomas", "uma", "vera", "wren", "yusuf",
}

var familyNames = []string{
	"abadie", "bellini", "cabral", "dahl", "eriksen", "ferreira", "gruber",
	"haddad", "ibrahim", "jansen", "kowalski", "lindqvist", "moreau", "novak",
	"okafor", "petrov", "quintero", "rossi", "silva", "tanaka", "unger", "volkov",
}

// scaleShape returns the (people, donationsPerPerson, charities) triple for a
// named scale. `tiny` exists so every scenario, invariant and control can be
// exercised in minutes; its numbers are never reported.
func scaleShape(scale string) (people, perPerson, charities int) {
	switch scale {
	case "tiny":
		return 40, 8, 4
	case "medium":
		return 3000, 40, 24
	default: // small
		return 800, 25, 12
	}
}

// zipfIndex draws an index in [0,n) from a bounded power law. Uniform parents hide
// the effects this study is about; an unbounded law would make them into strawmen.
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

const donationInterval = 36 * time.Hour

var loadEpoch = time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

func BuildDataset(seed int64, scale string) *Dataset {
	people, perPerson, nCharities := scaleShape(scale)
	r := rand.New(rand.NewSource(seed))
	nextDonationID := int64(1)

	d := &Dataset{
		Seed:      seed,
		Scale:     scale,
		Donations: map[int64][]Donation{},
		byPerson:  map[int64]int{},
		byChar:    map[int64]Charity{},
	}

	for i := 0; i < nCharities; i++ {
		c := Charity{
			ID:        int64(i + 1),
			Name:      fmt.Sprintf("%s-foundation", charityNames[i%len(charityNames)]),
			Country:   countries[i%len(countries)],
			FoundedOn: time.Date(1960+i*2, time.Month(1+i%12), 1, 0, 0, 0, 0, time.UTC).Format("2006-01-02"),
		}
		d.Charities = append(d.Charities, c)
		d.byChar[c.ID] = c
	}

	for i := 0; i < people; i++ {
		giving := perPerson
		// Bounded long tail: most donors give a few times, a small minority give
		// up to twice the tier. Bounded, so the reference embedding design is not a
		// strawman and the capacity calculation stays meaningful.
		if r.Intn(10) == 0 {
			giving = perPerson + r.Intn(perPerson+1)
		}
		p := Person{
			ID:        int64(i + 1),
			CharityID: int64(1 + i%nCharities),
			FullName:  fmt.Sprintf("%s %s", givenNames[i%len(givenNames)], familyNames[(i/len(givenNames))%len(familyNames)]),
			Email:     fmt.Sprintf("donor%06d@example.org", i+1),
			JoinedAt:  loadEpoch.AddDate(0, 0, -r.Intn(1500)).Truncate(time.Microsecond),
		}
		d.People = append(d.People, p)
		d.byPerson[p.ID] = len(d.People) - 1

		var id int64
		ds := make([]Donation, 0, giving)
		for k := 0; k < giving; k++ {
			// A single running counter, so donation ids are dense, unique and in
			// the same order as the append-only stream they represent. The workload
			// continues this counter when it inserts, so a new donation's id is
			// above every loaded id and cannot collide with one.
			id = nextDonationID
			nextDonationID++
			amt := int64(500 + r.Intn(500_000))
			var note *string
			if r.Intn(10) < 6 {
				n := fmt.Sprintf("gift %d for %s", k+1, charityNames[i%len(charityNames)])
				note = &n
			}
			ds = append(ds, Donation{
				ID:          id,
				PersonID:    p.ID,
				CharityID:   p.CharityID,
				AmountCents: amt,
				Currency:    currencies[r.Intn(len(currencies))],
				DonatedAt:   loadEpoch.Add(time.Duration(k) * donationInterval).Truncate(time.Microsecond),
				Note:        note,
			})
			if id > d.MaxDonationID {
				d.MaxDonationID = id
			}
		}
		d.Donations[p.ID] = ds
		d.TotalDonations += len(ds)
		for _, dn := range ds {
			d.SerializedBytes += int64(len(dn.Currency) + len(dn.Note2()))
		}
	}

	// The cached representation's logical size, measured from a real payload once
	// the encoder exists. A capacity chosen without it would be a guess.
	d.PayloadBytes = 0
	return d
}

// Note2 is a small helper so the byte accounting does not have to unwrap the
// pointer at every call site.
func (dn Donation) Note2() string {
	if dn.Note == nil {
		return ""
	}
	return *dn.Note
}

func (d *Dataset) Person(id int64) (Person, bool) {
	i, ok := d.byPerson[id]
	if !ok {
		return Person{}, false
	}
	return d.People[i], true
}

func (d *Dataset) Charity(id int64) (Charity, bool) {
	c, ok := d.byChar[id]
	return c, ok
}

// PickPerson draws a donor for a read. 80 % of the traffic goes to the hot fifth,
// which is what makes cache effects and contention visible at all.
func (d *Dataset) PickPerson(r *rand.Rand) Person {
	n := len(d.People)
	if n == 0 {
		return Person{}
	}
	if r.Intn(100) < 80 {
		hot := max(n/5, 1)
		return d.People[zipfIndex(r, hot, 1.1)]
	}
	return d.People[zipfIndex(r, n, 0.8)]
}

// HotPeople returns the n hottest donors, used by the hotspot, stampede and
// contention phases so they target a set the access distribution already favours.
func (d *Dataset) HotPeople(n int) []Person {
	if n > len(d.People) {
		n = len(d.People)
	}
	return append([]Person(nil), d.People[:n]...)
}

// DonationCount is the loaded donation count for one donor.
func (d *Dataset) DonationCount(personID int64) int { return len(d.Donations[personID]) }

// CacheKey is the study's primary key. One function so every artefact agrees on
// the key format, including the analysis.
func CacheKey(personID int64) string { return fmt.Sprintf("donor:%d:portal", personID) }
