package store

import (
	"errors"
	"path/filepath"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func newTestRepository(t *testing.T) *Repository {
	t.Helper()
	hash, err := auth.HashPassword("test-password", 4)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	st, err := testsupport.OpenStore(filepath.Join(t.TempDir(), "test.db"), "admin", hash, true)
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return NewRepository(st)
}

func TestRepositoryPolicyLifecycle(t *testing.T) {
	r := newTestRepository(t)
	now := time.Now().UTC()

	if _, err := r.GetPolicy("orders"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing policy err = %v, want ErrNotFound", err)
	}
	if err := r.UpsertPolicy("orders", "owner_id", "self", true, now); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if err := r.UpsertPolicy("orders", "owner_id", "team", true, now); !errors.Is(err, ErrInvalidScope) {
		t.Fatalf("invalid scope err = %v, want ErrInvalidScope", err)
	}
	p, err := r.GetPolicy("orders")
	if err != nil || p.OwnerColumn != "owner_id" || p.DefaultScope != "self" || !p.Enabled {
		t.Fatalf("policy = %+v, %v", p, err)
	}
	// Upsert is idempotent (update path).
	if err := r.UpsertPolicy("orders", "owner_id", "all", false, now); err != nil {
		t.Fatalf("re-upsert: %v", err)
	}
	p, err = r.GetPolicy("orders")
	if err != nil || p.DefaultScope != "all" || p.Enabled {
		t.Fatalf("updated policy = %+v, %v", p, err)
	}
	policies, err := r.ListPolicies()
	if err != nil || len(policies) != 1 {
		t.Fatalf("policies = %+v, %v", policies, err)
	}
}

func TestRepositoryAssignmentLifecycle(t *testing.T) {
	r := newTestRepository(t)
	now := time.Now().UTC()

	if err := r.UpsertAssignments("u1", map[string]string{"orders": "all", "catalog": "self"}, now); err != nil {
		t.Fatalf("upsert: %v", err)
	}
	if err := r.UpsertAssignments("u1", map[string]string{"orders": "team"}, now); !errors.Is(err, ErrInvalidScope) {
		t.Fatalf("invalid scope err = %v", err)
	}
	// Idempotent re-upsert.
	if err := r.UpsertAssignments("u1", map[string]string{"orders": "self"}, now); err != nil {
		t.Fatalf("re-upsert: %v", err)
	}
	a, err := r.GetAssignment("u1", "orders")
	if err != nil || a.ScopeType != "self" {
		t.Fatalf("assignment = %+v, %v", a, err)
	}
	if _, err := r.GetAssignment("u1", "missing"); !errors.Is(err, ErrNotFound) {
		t.Fatalf("missing assignment err = %v", err)
	}
	list, err := r.ListAssignments("u1")
	if err != nil || len(list) != 2 {
		t.Fatalf("assignments = %+v, %v", list, err)
	}
	list2, err := r.ListAssignments("u2")
	if err != nil || len(list2) != 0 {
		t.Fatalf("other user assignments = %+v, %v", list2, err)
	}
}
