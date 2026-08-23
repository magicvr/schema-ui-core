package wallet

// GOAL-037 / F-008 regression: newID must order SAME-millisecond entries by
// creation order (the D-002 §1 replay sorts (created_at ASC, id ASC)); the
// historical random-only suffix made same-ms ordering arbitrary, so replay
// could run a freeze before its funding adjust ("replay apply failed:
// insufficient balance", inconsistent reconcile).
import (
	"testing"
	"time"
)

func TestNewIDSameMillisecondOrdering(t *testing.T) {
	now := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	ids := make([]string, 0, 1000)
	for i := 0; i < 1000; i++ {
		id, err := newID(now)
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, id)
	}
	for i := 1; i < len(ids); i++ {
		if !(ids[i-1] < ids[i]) {
			t.Fatalf("same-millisecond ids out of creation order at %d: %q vs %q", i, ids[i-1], ids[i])
		}
	}
	// Uniqueness holds even with a constant clock (counter + random suffix).
	seen := map[string]bool{}
	for _, id := range ids {
		if seen[id] {
			t.Fatalf("duplicate id %q", id)
		}
		seen[id] = true
	}
}

func TestNewIDCrossMillisecondOrdering(t *testing.T) {
	base := time.Date(2026, 8, 23, 10, 0, 0, 0, time.UTC)
	id1, err := newID(base)
	if err != nil {
		t.Fatal(err)
	}
	id2, err := newID(base.Add(time.Millisecond))
	if err != nil {
		t.Fatal(err)
	}
	if !(id1 < id2) {
		t.Fatalf("cross-millisecond ids not ordered: %q vs %q", id1, id2)
	}
}