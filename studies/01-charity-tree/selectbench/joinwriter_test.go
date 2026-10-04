package selectbench

import "testing"

// The writer's fixed key set must be deterministic and bounded by the fixture,
// so a mutation count is comparable across cells and arms.
func TestMaintenanceKeys(t *testing.T) {
	k := MaintenanceKeys(200000, 1000)
	if len(k) != 1000 {
		t.Fatalf("expected 1000 keys, got %d", len(k))
	}
	if k[0] != 1 || k[999] != 1000 {
		t.Fatalf("keys are not the lowest ids in order: first %d last %d", k[0], k[999])
	}
	if got := MaintenanceKeys(5, 100); len(got) != 5 || got[4] != 5 {
		t.Fatalf("keys must clamp to the fixture size, got %v", got)
	}
	if got := MaintenanceKeys(10, 0); len(got) != 1 {
		t.Fatalf("keys must keep at least one entry, got %v", got)
	}
}
