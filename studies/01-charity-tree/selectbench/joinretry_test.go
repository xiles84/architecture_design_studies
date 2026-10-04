package selectbench

import (
	"context"
	"fmt"
	"testing"
)

// The blended phase's reader retries transient serialization conflicts. A retry
// must not become a way to accept a wrong answer: it only re-runs the same
// statement, and the caller still verifies against the oracle.
type flakyQuerier struct {
	n         int
	failFirst int
	permanent string
}

func (f *flakyQuerier) Exec(ctx context.Context, sql string) error { return nil }
func (f *flakyQuerier) Explain(ctx context.Context, sql string) (string, error) {
	return "plan", nil
}
func (f *flakyQuerier) Query(ctx context.Context, sql string) ([][]any, error) {
	f.n++
	if f.permanent != "" {
		return nil, fmt.Errorf("%s", f.permanent)
	}
	if f.n <= f.failFirst {
		return nil, fmt.Errorf("ERROR: Restart read required (SQLSTATE 40001)")
	}
	return [][]any{{1}}, nil
}

func TestQueryRetryTransientConflict(t *testing.T) {
	q := &flakyQuerier{failFirst: 2}
	rows, retries, err := queryRetry(context.Background(), q, "SELECT 1")
	if err != nil {
		t.Fatalf("retryable conflict should be retried, got %v", err)
	}
	if retries != 2 {
		t.Fatalf("expected 2 retries, got %d", retries)
	}
	if len(rows) != 1 {
		t.Fatalf("expected the retried answer, got %v", rows)
	}
}

func TestQueryRetryDoesNotRetryPermanentError(t *testing.T) {
	q := &flakyQuerier{permanent: "syntax error at or near \"SELEC\""}
	_, retries, err := queryRetry(context.Background(), q, "SELEC 1")
	if err == nil {
		t.Fatal("a permanent error must not be swallowed")
	}
	if retries != 0 {
		t.Fatalf("a permanent error must not be retried, got %d retries", retries)
	}
}
