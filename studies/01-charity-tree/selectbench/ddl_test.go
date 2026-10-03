package selectbench

import (
	"strings"
	"testing"
)

func TestSchemaSQLDefinesTheTeachingSchema(t *testing.T) {
	s := SchemaSQL()
	for _, want := range []string{"CREATE SCHEMA IF NOT EXISTS selectbench",
		"selectbench.customers", "selectbench.orders", "orders_amount_idx", "orders_customer_idx"} {
		if !strings.Contains(s, want) {
			t.Errorf("schema does not contain %q", want)
		}
	}
	// It must be re-runnable: the runner loads the same fixture repeatedly.
	if !strings.Contains(s, "DROP TABLE IF EXISTS") {
		t.Error("schema is not idempotent: no DROP TABLE IF EXISTS")
	}
}

func TestLoadSQLIsDeterministic(t *testing.T) {
	a := LoadSQL(Teaching())
	b := LoadSQL(Teaching())
	if a != b {
		t.Error("the loader SQL is not deterministic for the same fixture")
	}
	// Generation order must not matter either.
	fx := Teaching()
	r := fx.Orders
	fx.Orders = []Order{r[3], r[0], r[6], r[1], r[5], r[2], r[4]}
	if LoadSQL(fx) != a {
		t.Error("the loader SQL depends on the input order of orders")
	}
}

func TestLoadSQLCountsRowsAndRendersNulls(t *testing.T) {
	fx := Teaching()
	sql := LoadSQL(fx)
	if n := strings.Count(sql, "INSERT INTO selectbench.customers"); n != len(fx.Customers) {
		t.Errorf("customer inserts = %d, want %d", n, len(fx.Customers))
	}
	if n := strings.Count(sql, "INSERT INTO selectbench.orders"); n != len(fx.Orders) {
		t.Errorf("order inserts = %d, want %d", n, len(fx.Orders))
	}
	// Order 13 has a NULL amount and must be emitted as SQL NULL, not 'NULL'.
	if !strings.Contains(sql, "(13, 2, NULL, NULL);") {
		t.Errorf("order 13's NULLs were not rendered as SQL NULL:\n%s", sql)
	}
	if strings.Contains(sql, "'NULL'") {
		t.Error("a NULL was rendered as the string 'NULL'")
	}
	// Order 15's NULL foreign key must survive as NULL.
	if !strings.Contains(sql, "(15, NULL, 150, NULL);") {
		t.Errorf("order 15 row wrong:\n%s", sql)
	}
}

func TestFingerprintIsContentSensitive(t *testing.T) {
	fx := Teaching()
	base := Fingerprint(fx)
	if base == "" {
		t.Fatal("empty fingerprint")
	}
	// Reordering the input slices must not change the fingerprint.
	shuffled := fx
	shuffled.Orders = []Order{fx.Orders[2], fx.Orders[0], fx.Orders[1], fx.Orders[5], fx.Orders[3], fx.Orders[6], fx.Orders[4]}
	shuffled.Customers = []Customer{fx.Customers[1], fx.Customers[0], fx.Customers[3], fx.Customers[2]}
	if Fingerprint(shuffled) != base {
		t.Error("fingerprint depends on input order")
	}
	// Changing a value must change it.
	changed := fx
	changed.Orders = append([]Order(nil), fx.Orders...)
	changed.Orders[1].Amount = ip(999)
	if Fingerprint(changed) == base {
		t.Error("fingerprint did not change when the data changed")
	}
}
