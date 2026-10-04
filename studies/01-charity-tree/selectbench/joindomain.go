package selectbench

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
	"strings"
)

// Study 01 domain binding of the join/CTE family (review R2 / handoff C2).
//
// The family's reusable catalogue is written against the projection
// `customers(id, name)` and `facts(id, customer_id, rank, amount, payload)`.
// The domain arm loads the study's real entities — charity, person, donation —
// and exposes exactly that projection through two views, so the measured SQL is
// byte-identical to the synthetic arm and only the data's origin differs:
//
//   customers(id, name) = (person_id, the person's charity name)
//   facts(id, customer_id, rank, amount, payload) =
//       (donation_id, person_id, donation.rank, amount_cents, note)
//
// The projection is therefore (person_id, charity_name, donation_id,
// amount_cents): the parent is a person, the child is a donation, and the
// person's charity is the relationship the Charity donations case turns on.
//
// `donation.rank` is a deterministic shuffle of 1..N. Study 01's real selective
// predicate is a date or amount range, which cannot produce an exact selectivity;
// the reviewed matrix needs exact thresholds, so `rank <= threshold` is the
// documented proxy and the label travels with the result.
//
// Unlike the synthetic microfixture, the domain FKs are NOT NULL, matching Study
// 01's schema. The optional-match contract is exercised by persons with no
// qualifying donation; the synthetic microfixture keeps the orphan-NULL-FK case
// as its own evidence.

// JoinDomainSchema is the schema the domain tables and binding views live in.
const JoinDomainSchema = "selectjoin_domain"

// Charity is one Study 01 charity dimension row.
type Charity struct {
	ID   int
	Name string
}

// Person is one Study 01 person; CharityID is its charity relationship.
type Person struct {
	ID        int
	CharityID int
	FullName  string
}

// Donation is one Study 01 donation. Rank is the selectivity proxy; AmountCents
// is the projected measure.
type Donation struct {
	ID          int
	PersonID    int
	Rank        int
	AmountCents *int
	Note        string
}

// JoinDomainFixture is one domain dataset.
type JoinDomainFixture struct {
	Charities []Charity
	Persons   []Person
	Donations []Donation
}

// JoinDomainGenSpec is a deterministic domain shape.
type JoinDomainGenSpec struct {
	Charities       int
	Persons         int
	Donations       int
	NullAmountEvery int // 0 = no NULL amounts
	Wide            bool
}

// GenerateJoinDomain builds the deterministic domain fixture with the same
// seeded Fisher-Yates rank permutation the synthetic generator uses, so a
// threshold selects an exact, reproducible count.
func GenerateJoinDomain(seed int64, spec JoinDomainGenSpec) JoinDomainFixture {
	df := JoinDomainFixture{}
	for i := 1; i <= spec.Charities; i++ {
		df.Charities = append(df.Charities, Charity{ID: i, Name: "ch" + strconv.Itoa(i)})
	}
	for i := 1; i <= spec.Persons; i++ {
		cid := 1
		if spec.Charities > 0 {
			cid = ((i - 1) % spec.Charities) + 1
		}
		df.Persons = append(df.Persons, Person{ID: i, CharityID: cid, FullName: "p" + strconv.Itoa(i)})
	}
	n := spec.Donations
	perm := make([]int, n)
	for i := range perm {
		perm[i] = i + 1
	}
	r := rand.New(rand.NewSource(seed))
	r.Shuffle(n, func(i, j int) { perm[i], perm[j] = perm[j], perm[i] })
	df.Donations = make([]Donation, 0, n)
	for i := 1; i <= n; i++ {
		pid := 1
		if spec.Persons > 0 {
			pid = ((i - 1) % spec.Persons) + 1
		}
		var amt *int
		if spec.NullAmountEvery == 0 || i%spec.NullAmountEvery != 0 {
			amt = ip(1 + (i*37)%500)
		}
		d := Donation{ID: i, PersonID: pid, Rank: perm[i-1], AmountCents: amt}
		if spec.Wide {
			d.Note = "d" + strconv.Itoa(i) + "0123456789abcdef0123456789abcdef"
		}
		df.Donations = append(df.Donations, d)
	}
	return df
}

// JoinDomainEdge is the domain gate fixture: duplicate matches, a NULL amount
// and an unmatched person for the optional-match control, on NOT NULL FKs.
func JoinDomainEdge() JoinDomainFixture {
	return JoinDomainFixture{
		Charities: []Charity{{1, "Alpha"}, {2, "Beta"}},
		Persons: []Person{
			{ID: 1, CharityID: 1, FullName: "Ada"},
			{ID: 2, CharityID: 1, FullName: "Ben"},
			{ID: 3, CharityID: 2, FullName: "Cara"},
			{ID: 4, CharityID: 2, FullName: "Dan"},
		},
		Donations: []Donation{
			{ID: 101, PersonID: 1, Rank: 1, AmountCents: ip(50)},
			{ID: 102, PersonID: 1, Rank: 2, AmountCents: ip(50)},
			{ID: 103, PersonID: 2, Rank: 3, AmountCents: nil},
			{ID: 104, PersonID: 4, Rank: 5, AmountCents: ip(300)},
		},
	}
}

// AsJoinFixture maps the domain rows onto the projection the family's oracles
// and catalogue already use. The person's charity name is the projected parent
// name; the donation is the child. It is the single place the domain↔catalogue
// binding is defined.
func (df JoinDomainFixture) AsJoinFixture() JoinFixture {
	charityName := map[int]string{}
	for _, c := range df.Charities {
		charityName[c.ID] = c.Name
	}
	jf := JoinFixture{}
	for _, p := range df.Persons {
		jf.Customers = append(jf.Customers, Customer{ID: p.ID, Name: charityName[p.CharityID]})
	}
	for _, d := range df.Donations {
		pid := d.PersonID
		jf.Facts = append(jf.Facts, Fact{ID: d.ID, CustomerID: &pid, Rank: d.Rank, Amount: d.AmountCents, Payload: d.Note})
	}
	return jf
}

// JoinDomainSchemaSQL creates the domain tables and the two binding views,
// dropping any previous copy.
func JoinDomainSchemaSQL() string {
	return `CREATE SCHEMA IF NOT EXISTS ` + JoinDomainSchema + `;
DROP VIEW IF EXISTS ` + JoinDomainSchema + `.facts;
DROP VIEW IF EXISTS ` + JoinDomainSchema + `.customers;
DROP TABLE IF EXISTS ` + JoinDomainSchema + `.donation;
DROP TABLE IF EXISTS ` + JoinDomainSchema + `.person;
DROP TABLE IF EXISTS ` + JoinDomainSchema + `.charity;
CREATE TABLE ` + JoinDomainSchema + `.charity (
  charity_id integer PRIMARY KEY,
  name       text    NOT NULL
);
CREATE TABLE ` + JoinDomainSchema + `.person (
  person_id  integer PRIMARY KEY,
  charity_id integer NOT NULL REFERENCES ` + JoinDomainSchema + `.charity(charity_id),
  full_name  text    NOT NULL
);
CREATE TABLE ` + JoinDomainSchema + `.donation (
  donation_id  integer PRIMARY KEY,
  person_id    integer NOT NULL REFERENCES ` + JoinDomainSchema + `.person(person_id),
  rank         integer NOT NULL,
  amount_cents integer,
  note         text
);
CREATE INDEX donation_rank_idx ON ` + JoinDomainSchema + `.donation (rank);
CREATE INDEX donation_person_idx ON ` + JoinDomainSchema + `.donation (person_id);
CREATE VIEW ` + JoinDomainSchema + `.customers AS
  SELECT p.person_id AS id, ch.name AS name
  FROM ` + JoinDomainSchema + `.person p
  JOIN ` + JoinDomainSchema + `.charity ch ON ch.charity_id = p.charity_id;
CREATE VIEW ` + JoinDomainSchema + `.facts AS
  SELECT d.donation_id AS id, d.person_id AS customer_id, d.rank AS rank,
         d.amount_cents AS amount, d.note AS payload
  FROM ` + JoinDomainSchema + `.donation d;
`
}

// JoinDomainAnalyzeSQL stats the domain tables after load.
func JoinDomainAnalyzeSQL() string {
	return "ANALYZE " + JoinDomainSchema + ".charity; ANALYZE " + JoinDomainSchema + ".person; ANALYZE " + JoinDomainSchema + ".donation;"
}

// JoinDomainLoadSQL is the deterministic, batched domain loader. Rows are
// emitted in id order, so the script is byte-identical for the same fixture.
func JoinDomainLoadSQL(df JoinDomainFixture, batch int) []string {
	if batch < 1 {
		batch = 1000
	}
	var out []string
	charities := append([]Charity(nil), df.Charities...)
	sort.Slice(charities, func(i, j int) bool { return charities[i].ID < charities[j].ID })
	if len(charities) > 0 {
		var b strings.Builder
		b.WriteString("INSERT INTO " + JoinDomainSchema + ".charity (charity_id, name) VALUES ")
		for i, c := range charities {
			if i > 0 {
				b.WriteString(",")
			}
			fmt.Fprintf(&b, "(%d,%s)", c.ID, quoteText(c.Name))
		}
		out = append(out, b.String())
	}
	persons := append([]Person(nil), df.Persons...)
	sort.Slice(persons, func(i, j int) bool { return persons[i].ID < persons[j].ID })
	for start := 0; start < len(persons); start += batch {
		end := start + batch
		if end > len(persons) {
			end = len(persons)
		}
		var b strings.Builder
		b.WriteString("INSERT INTO " + JoinDomainSchema + ".person (person_id, charity_id, full_name) VALUES ")
		for i := start; i < end; i++ {
			if i > start {
				b.WriteString(",")
			}
			p := persons[i]
			fmt.Fprintf(&b, "(%d,%d,%s)", p.ID, p.CharityID, quoteText(p.FullName))
		}
		out = append(out, b.String())
	}
	donations := append([]Donation(nil), df.Donations...)
	sort.Slice(donations, func(i, j int) bool { return donations[i].ID < donations[j].ID })
	for start := 0; start < len(donations); start += batch {
		end := start + batch
		if end > len(donations) {
			end = len(donations)
		}
		var b strings.Builder
		b.WriteString("INSERT INTO " + JoinDomainSchema + ".donation (donation_id, person_id, rank, amount_cents, note) VALUES ")
		for i := start; i < end; i++ {
			if i > start {
				b.WriteString(",")
			}
			d := donations[i]
			note := "NULL"
			if d.Note != "" {
				note = quoteText(d.Note)
			}
			fmt.Fprintf(&b, "(%d,%d,%d,%s,%s)", d.ID, d.PersonID, d.Rank, nullInt(d.AmountCents), note)
		}
		out = append(out, b.String())
	}
	return out
}

// JoinDomainFingerprint is a short, stable identifier for the domain fixture.
func JoinDomainFingerprint(df JoinDomainFixture) string {
	lines := make([]string, 0, len(df.Charities)+len(df.Persons)+len(df.Donations))
	for _, c := range df.Charities {
		lines = append(lines, "ch:"+strconv.Itoa(c.ID)+"="+c.Name)
	}
	for _, p := range df.Persons {
		lines = append(lines, "p:"+strconv.Itoa(p.ID)+":"+strconv.Itoa(p.CharityID)+":"+p.FullName)
	}
	for _, d := range df.Donations {
		lines = append(lines, "d:"+strconv.Itoa(d.ID)+":"+strconv.Itoa(d.PersonID)+":"+
			strconv.Itoa(d.Rank)+":"+nullInt(d.AmountCents)+":"+d.Note)
	}
	sort.Strings(lines)
	return hashLines(lines)
}
