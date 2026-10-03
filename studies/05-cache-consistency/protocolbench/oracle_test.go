package main

import (
	"context"
	"testing"
	"time"
)

func TestLedgerRejectsUncommittedAndDistinguishesAcknowledgement(t *testing.T) {
	ctx := context.Background()
	l := newLedger(1)
	old := expected(1, 1, false)
	fresh := expected(1, 2, false)
	p := l.begin(1, fresh)
	l.finish(1, p, false)
	if e := requireClass(ctx, l, 1, fresh, "impossible"); e != nil {
		t.Fatal(e)
	}
	p = l.begin(1, fresh)
	l.finish(1, p, true)
	if e := requireClass(ctx, l, 1, old, "committed"); e != nil {
		t.Fatal(e)
	}
	l.acknowledge(1, 2)
	if e := requireClass(ctx, l, 1, old, "stale-after-ack"); e != nil {
		t.Fatal(e)
	}
	forged := fresh
	forged.Value++
	if e := requireClass(ctx, l, 1, forged, "impossible"); e != nil {
		t.Fatal(e)
	}
}

func TestAheadSnapshotWaitsForCommitOutcome(t *testing.T) {
	for _, committed := range []bool{true, false} {
		l := newLedger(1)
		s := expected(1, 2, false)
		p := l.begin(1, s)
		done := make(chan string, 1)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		go func() {
			got, e := l.classify(ctx, 1, s, 1)
			if e != nil {
				done <- e.Error()
			} else {
				done <- got
			}
		}()
		select {
		case got := <-done:
			t.Fatalf("pending state admitted before outcome: %s", got)
		case <-time.After(10 * time.Millisecond):
		}
		l.finish(1, p, committed)
		want := "impossible"
		if committed {
			want = "committed"
		}
		if got := <-done; got != want {
			t.Fatalf("outcome %t: got %s want %s", committed, got, want)
		}
		cancel()
	}
}
func TestLedgerRetainsDeletionAcrossRecreation(t *testing.T) {
	l := newLedger(1)
	ctx := context.Background()
	gone := expected(1, 2, true)
	p := l.begin(1, gone)
	l.finish(1, p, true)
	l.acknowledge(1, 2)
	if e := requireClass(ctx, l, 1, gone, "committed"); e != nil {
		t.Fatal(e)
	}
	reborn := expected(1, 3, false)
	p = l.begin(1, reborn)
	l.finish(1, p, true)
	l.acknowledge(1, 3)
	if e := requireClass(ctx, l, 1, gone, "stale-after-ack"); e != nil {
		t.Fatal(e)
	}
}
