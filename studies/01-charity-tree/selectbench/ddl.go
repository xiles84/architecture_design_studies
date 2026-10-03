package selectbench

import (
	"sort"
	"strconv"
	"strings"
)

// The additive teaching schema.
//
// The SELECT experiment runs against its own schema so it cannot disturb Study
// 01's existing tables and evidence, and so the logical data the oracle checks
// is exactly the data the loader inserted. Both tables are small on purpose:
// the point is to locate a crossover, not to produce a large throughput number.

// TeachingSchema is the schema the SELECT tables live in.
const TeachingSchema = "selectbench"

// SchemaSQL creates the fixture tables, dropping any previous copy. It is
// idempotent and is the only place the SELECT tables are defined, so a report
// and the oracle cannot disagree about the shape of the data.
func SchemaSQL() string {
	return `CREATE SCHEMA IF NOT EXISTS ` + TeachingSchema + `;
DROP TABLE IF EXISTS ` + TeachingSchema + `.orders;
DROP TABLE IF EXISTS ` + TeachingSchema + `.customers;
CREATE TABLE ` + TeachingSchema + `.customers (
  id    integer PRIMARY KEY,
  name  text    NOT NULL
);
CREATE TABLE ` + TeachingSchema + `.orders (
  id          integer PRIMARY KEY,
  customer_id integer REFERENCES ` + TeachingSchema + `.customers(id),
  amount      integer,
  payload     text
);
CREATE INDEX orders_amount_idx ON ` + TeachingSchema + `.orders (amount);
CREATE INDEX orders_customer_idx ON ` + TeachingSchema + `.orders (customer_id);
`
}

func quoteText(s string) string { return "'" + strings.ReplaceAll(s, "'", "''") + "'" }

func nullInt(p *int) string {
	if p == nil {
		return "NULL"
	}
	return strconv.Itoa(*p)
}

// LoadSQL is the deterministic INSERT script for a fixture. Rows are emitted in
// id order so the script is byte-identical for the same fixture, which lets a
// run's dataset be reproduced and fingerprinted.
func LoadSQL(fx Fixture) string {
	var b strings.Builder
	customers := append([]Customer(nil), fx.Customers...)
	sort.Slice(customers, func(i, j int) bool { return customers[i].ID < customers[j].ID })
	for _, c := range customers {
		b.WriteString("INSERT INTO " + TeachingSchema + ".customers (id, name) VALUES (")
		b.WriteString(strconv.Itoa(c.ID))
		b.WriteString(", ")
		b.WriteString(quoteText(c.Name))
		b.WriteString(");\n")
	}
	orders := append([]Order(nil), fx.Orders...)
	sort.Slice(orders, func(i, j int) bool { return orders[i].ID < orders[j].ID })
	for _, o := range orders {
		b.WriteString("INSERT INTO " + TeachingSchema + ".orders (id, customer_id, amount, payload) VALUES (")
		b.WriteString(strconv.Itoa(o.ID))
		b.WriteString(", ")
		b.WriteString(nullInt(o.CustomerID))
		b.WriteString(", ")
		b.WriteString(nullInt(o.Amount))
		b.WriteString(", ")
		if o.Payload == "" {
			b.WriteString("NULL")
		} else {
			b.WriteString(quoteText(o.Payload))
		}
		b.WriteString(");\n")
	}
	return b.String()
}

// Fingerprint is a short, stable identifier for a fixture's logical contents,
// recorded with every run so two results can be shown to share data or not.
func Fingerprint(fx Fixture) string {
	lines := []string{}
	for _, c := range fx.Customers {
		lines = append(lines, "c:"+strconv.Itoa(c.ID)+"="+c.Name)
	}
	for _, o := range fx.Orders {
		lines = append(lines, "o:"+strconv.Itoa(o.ID)+":"+nullInt(o.CustomerID)+":"+nullInt(o.Amount)+":"+o.Payload)
	}
	sort.Strings(lines)
	return hashLines(lines)
}
