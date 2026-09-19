package jobs_test

import (
	"context"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
)

// GOAL-003 R2 (D-001 §1.3): the management-scope list must filter, paginate,
// sort and count consistently, and must not disturb the actor-scoped paths.
func TestListJobsFiltersPaginatesAndCounts(t *testing.T) {
	repo, _ := newRepository(t)
	now := testNow
	base := now.Add(-time.Hour)

	seed := []struct {
		id      string
		kind    string
		actorID string
		at      time.Time
	}{
		{"job-a", "wallet.reconcile", "user-1", base},
		{"job-b", "wallet.reconcile", "user-2", base.Add(1 * time.Minute)},
		{"job-c", "jobs.batch-export", "user-1", base.Add(2 * time.Minute)},
		{"job-d", "jobs.batch-export", "user-2", base.Add(3 * time.Minute)},
		{"job-e", "jobs.batch-export", "user-1", base.Add(4 * time.Minute)},
	}
	for _, s := range seed {
		if _, err := repo.Create(context.Background(), jobs.CreateInput{
			ID: s.id, Kind: s.kind, Payload: []byte(`{}`),
			ActorID: s.actorID, CorrelationID: "corr-" + s.id, Now: s.at,
		}); err != nil {
			t.Fatalf("create %s: %v", s.id, err)
		}
	}

	t.Run("default order is created_at DESC with total", func(t *testing.T) {
		items, total, err := repo.ListJobs(context.Background(), jobs.ListFilter{Page: 1, PageSize: 10})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if total != len(seed) {
			t.Fatalf("total = %d, want %d", total, len(seed))
		}
		want := []string{"job-e", "job-d", "job-c", "job-b", "job-a"}
		if len(items) != len(want) {
			t.Fatalf("items = %d, want %d", len(items), len(want))
		}
		for i, id := range want {
			if items[i].ID != id {
				t.Fatalf("items[%d] = %s, want %s (newest first)", i, items[i].ID, id)
			}
		}
	})

	t.Run("kind filter narrows both items and total", func(t *testing.T) {
		items, total, err := repo.ListJobs(context.Background(), jobs.ListFilter{
			Kind: "jobs.batch-export", Page: 1, PageSize: 10,
		})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if total != 3 || len(items) != 3 {
			t.Fatalf("kind filter total=%d items=%d, want 3/3", total, len(items))
		}
		for _, item := range items {
			if item.Kind != "jobs.batch-export" {
				t.Fatalf("kind = %s, want jobs.batch-export", item.Kind)
			}
		}
	})

	t.Run("actor filter is the management-scope cross-actor read", func(t *testing.T) {
		items, total, err := repo.ListJobs(context.Background(), jobs.ListFilter{
			ActorID: "user-2", Page: 1, PageSize: 10,
		})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if total != 2 || len(items) != 2 {
			t.Fatalf("actor filter total=%d items=%d, want 2/2", total, len(items))
		}
	})

	t.Run("time range bounds created_at inclusively", func(t *testing.T) {
		items, total, err := repo.ListJobs(context.Background(), jobs.ListFilter{
			From: base.Add(1 * time.Minute), To: base.Add(3 * time.Minute),
			Page: 1, PageSize: 10,
		})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if total != 3 || len(items) != 3 {
			t.Fatalf("range total=%d items=%d, want 3/3", total, len(items))
		}
	})

	t.Run("pagination keeps total stable across pages", func(t *testing.T) {
		first, total, err := repo.ListJobs(context.Background(), jobs.ListFilter{Page: 1, PageSize: 2})
		if err != nil {
			t.Fatalf("page 1: %v", err)
		}
		if total != 5 || len(first) != 2 {
			t.Fatalf("page1 total=%d items=%d, want 5/2", total, len(first))
		}
		second, total2, err := repo.ListJobs(context.Background(), jobs.ListFilter{Page: 2, PageSize: 2})
		if err != nil {
			t.Fatalf("page 2: %v", err)
		}
		if total2 != 5 || len(second) != 2 {
			t.Fatalf("page2 total=%d items=%d, want 5/2", total2, len(second))
		}
		if first[0].ID == second[0].ID {
			t.Fatalf("page 2 repeated page 1 item %s", first[0].ID)
		}
		// A page beyond the last page is an empty page, not an error.
		beyond, _, err := repo.ListJobs(context.Background(), jobs.ListFilter{Page: 99, PageSize: 2})
		if err != nil {
			t.Fatalf("page 99: %v", err)
		}
		if len(beyond) != 0 {
			t.Fatalf("beyond-last page items = %d, want 0", len(beyond))
		}
	})

	t.Run("unknown sort key falls back instead of reaching SQL", func(t *testing.T) {
		items, _, err := repo.ListJobs(context.Background(), jobs.ListFilter{
			Sort: "id; DROP TABLE jobs", Order: "desc", Page: 1, PageSize: 10,
		})
		if err != nil {
			t.Fatalf("list with hostile sort: %v", err)
		}
		if len(items) != 5 || items[0].ID != "job-e" {
			t.Fatalf("hostile sort must fall back to created_at DESC; got %d items", len(items))
		}
	})

	t.Run("status filter selects terminal states the runtime query cannot show", func(t *testing.T) {
		items, total, err := repo.ListJobs(context.Background(), jobs.ListFilter{
			Status: "queued", Page: 1, PageSize: 10,
		})
		if err != nil {
			t.Fatalf("list: %v", err)
		}
		if total != 5 || len(items) != 5 {
			t.Fatalf("queued total=%d items=%d, want 5/5", total, len(items))
		}
	})
}

// GetJob is the management-scope detail read; it must return jobs the
// actor-scoped GetForActor refuses, while GetForActor itself is unchanged.
func TestGetJobIsManagementScopeButGetForActorIsNot(t *testing.T) {
	repo, _ := newRepository(t)
	if _, err := repo.Create(context.Background(), jobs.CreateInput{
		ID: "job-mgmt", Kind: "wallet.reconcile", Payload: []byte(`{}`),
		ActorID: "owner", CorrelationID: "corr-1", Now: testNow,
	}); err != nil {
		t.Fatalf("create: %v", err)
	}

	got, err := repo.GetJob(context.Background(), "job-mgmt")
	if err != nil || got.ID != "job-mgmt" {
		t.Fatalf("GetJob = %+v err=%v, want the job", got, err)
	}
	// The actor-scoped path keeps its frozen isolation semantics.
	if _, err := repo.GetForActor(context.Background(), "job-mgmt", "wallet.reconcile", "other-user"); err == nil {
		t.Fatalf("GetForActor for another actor must still fail closed")
	}
	if _, err := repo.GetForActor(context.Background(), "job-mgmt", "wallet.reconcile", "owner"); err != nil {
		t.Fatalf("GetForActor for the owner = %v, want the job", err)
	}
	// A wrong kind still fails closed on the actor-scoped path.
	if _, err := repo.GetForActor(context.Background(), "job-mgmt", "jobs.batch-export", "owner"); err == nil {
		t.Fatalf("GetForActor with a mismatched kind must fail closed")
	}
}

// ResultURL is the shared derivation that replaces the hardcoded wallet path.
// The wallet module's base path must reproduce the historical string exactly.
func TestResultURLReproducesTheHistoricalWalletAddress(t *testing.T) {
	got := jobs.ResultURL("/api/wallet/jobs", "job-1")
	if got != "/api/wallet/jobs/job-1/result" {
		t.Fatalf("wallet result URL = %q, want the historical value", got)
	}
	if got := jobs.ResultURL("/api/jobs", "job-2"); got != "/api/jobs/job-2/result" {
		t.Fatalf("admin.jobs result URL = %q", got)
	}
	// A trailing slash on the declared base path must not double up.
	if got := jobs.ResultURL("/api/jobs/", "job-3"); got != "/api/jobs/job-3/result" {
		t.Fatalf("trailing-slash base = %q", got)
	}
	// The derivation must be byte-identical for ANY id, including hostile ones:
	// it is the drop-in replacement for "/api/wallet/jobs/" + id + "/result",
	// so no id may be normalized, escaped or trimmed (A-002 F-003).
	for _, id := range []string{"", "job with space", "a/b", "a?b#c", "任务-1"} {
		want := "/api/wallet/jobs/" + id + "/result"
		if got := jobs.ResultURL("/api/wallet/jobs", id); got != want {
			t.Fatalf("ResultURL(%q) = %q, want %q (must be byte-identical)", id, got, want)
		}
	}
}
