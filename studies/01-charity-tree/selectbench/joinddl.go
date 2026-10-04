package selectbench

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// The join/CTE schema and deterministic loader.
//
// The family keeps its own schema so it cannot disturb the teaching fixture the
// other SELECT families load, and so a run can be replayed from the same
// logical rows the oracle reads. The `rank` column carries the deterministic
// permutation: an index on it is part of the fixed physical arm every spelling
// pair shares.

// JoinSchema is the schema the join/CTE tables live in.
const JoinSchema = "selectjoin"

// JoinSchemaSQL creates the join fixture tables, dropping any previous copy.
func JoinSchemaSQL() string {
	return `CREATE SCHEMA IF NOT EXISTS ` + JoinSchema + `;
DROP TABLE IF EXISTS ` + JoinSchema + `.facts;
DROP TABLE IF EXISTS ` + JoinSchema + `.customers;
CREATE TABLE ` + JoinSchema + `.customers (
  id    integer PRIMARY KEY,
  name  text    NOT NULL
);
CREATE TABLE ` + JoinSchema + `.facts (
  id          integer PRIMARY KEY,
  customer_id integer REFERENCES ` + JoinSchema + `.customers(id),
  rank        integer NOT NULL,
  amount      integer,
  payload     text
);
CREATE INDEX facts_rank_idx ON ` + JoinSchema + `.facts (rank);
CREATE INDEX facts_customer_idx ON ` + JoinSchema + `.facts (customer_id);
`
}

// JoinLoadSQL is the deterministic loader, batched so a 200,000-row insert is
// a bounded number of statements rather than one statement per row. Rows are
// emitted in id order, so the script is byte-identical for the same fixture.
func JoinLoadSQL(jf JoinFixture, batch int) []string {
	if batch < 1 {
		batch = 1000
	}
	var out []string
	customers := append([]Customer(nil), jf.Customers...)
	sort.Slice(customers, func(i, j int) bool { return customers[i].ID < customers[j].ID })
	if len(customers) > 0 {
		var b strings.Builder
		b.WriteString("INSERT INTO " + JoinSchema + ".customers (id, name) VALUES ")
		for i, c := range customers {
			if i > 0 {
				b.WriteString(",")
			}
			fmt.Fprintf(&b, "(%d,%s)", c.ID, quoteText(c.Name))
		}
		out = append(out, b.String())
	}
	facts := append([]Fact(nil), jf.Facts...)
	sort.Slice(facts, func(i, j int) bool { return facts[i].ID < facts[j].ID })
	for start := 0; start < len(facts); start += batch {
		end := start + batch
		if end > len(facts) {
			end = len(facts)
		}
		var b strings.Builder
		b.WriteString("INSERT INTO " + JoinSchema + ".facts (id, customer_id, rank, amount, payload) VALUES ")
		for i := start; i < end; i++ {
			if i > start {
				b.WriteString(",")
			}
			f := facts[i]
			payload := "NULL"
			if f.Payload != "" {
				payload = quoteText(f.Payload)
			}
			fmt.Fprintf(&b, "(%d,%s,%d,%s,%s)", f.ID, nullInt(f.CustomerID), f.Rank, nullInt(f.Amount), payload)
		}
		out = append(out, b.String())
	}
	return out
}

// JoinFingerprint is a short, stable identifier for the logical join fixture.
func JoinFingerprint(jf JoinFixture) string {
	lines := make([]string, 0, len(jf.Customers)+len(jf.Facts))
	for _, c := range jf.Customers {
		lines = append(lines, "c:"+strconv.Itoa(c.ID)+"="+c.Name)
	}
	for _, f := range jf.Facts {
		lines = append(lines, "f:"+strconv.Itoa(f.ID)+":"+nullInt(f.CustomerID)+":"+
			strconv.Itoa(f.Rank)+":"+nullInt(f.Amount)+":"+f.Payload)
	}
	sort.Strings(lines)
	return hashLines(lines)
}
