package subject_test

import (
	"context"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
)

func now() time.Time { return time.Unix(1700000000, 0).UTC() }

func newSubjectStore(t *testing.T) *subject.Store {
	t.Helper()
	st, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { st.Close() })
	return subject.NewStore(st)
}

func TestGetOrCreateSubject(t *testing.T) {
	s := newSubjectStore(t)
	ctx := context.Background()

	// 1. First get-or-create should create the row.
	sub, created, err := s.GetOrCreateSubject(ctx, "telegram", "123456789", now())
	if err != nil {
		t.Fatalf("first get-or-create: %v", err)
	}
	if !created {
		t.Fatalf("expected created=true on first call")
	}
	if sub.ID == "" || sub.Issuer != "telegram" || sub.ExternalID != "123456789" {
		t.Fatalf("unexpected subject fields: %+v", sub)
	}

	// 2. Second get-or-create with same issuer + external_id should be idempotent.
	sub2, created2, err := s.GetOrCreateSubject(ctx, "telegram", "123456789", now().Add(time.Hour))
	if err != nil {
		t.Fatalf("second get-or-create: %v", err)
	}
	if created2 {
		t.Fatalf("expected created=false on duplicate call")
	}
	if sub2.ID != sub.ID {
		t.Fatalf("id mismatch: got %s, want %s", sub2.ID, sub.ID)
	}

	// 3. Different issuer, same external_id -> different subject.
	sub3, created3, err := s.GetOrCreateSubject(ctx, "discord", "123456789", now())
	if err != nil {
		t.Fatalf("third get-or-create: %v", err)
	}
	if !created3 || sub3.ID == sub.ID {
		t.Fatalf("expected new subject for different issuer")
	}

	// 4. Verify SubjectExists.
	exists, err := s.SubjectExists(ctx, sub.ID)
	if err != nil || !exists {
		t.Fatalf("expected exists=true for %s, got %v, err=%v", sub.ID, exists, err)
	}
	notExists, err := s.SubjectExists(ctx, "non-existent-id")
	if err != nil || notExists {
		t.Fatalf("expected exists=false, got %v, err=%v", notExists, err)
	}
}

func TestSubjectInvalidInputs(t *testing.T) {
	s := newSubjectStore(t)
	ctx := context.Background()

	if _, _, err := s.GetOrCreateSubject(ctx, "", "123", now()); err == nil {
		t.Fatal("expected error on empty issuer")
	}
	if _, _, err := s.GetOrCreateSubject(ctx, "tg", "", now()); err == nil {
		t.Fatal("expected error on empty external_id")
	}
}

func TestConcurrentGetOrCreateSubject(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "subject-race.db")
	st, err := testsupport.OpenStore(dbPath, "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	defer st.Close()

	s := subject.NewStore(st)
	ctx := context.Background()

	const concurrency = 15
	var wg sync.WaitGroup
	var createdCount atomic.Int32
	var ids sync.Map

	startGate := make(chan struct{})

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-startGate
			sub, created, err := s.GetOrCreateSubject(ctx, "telegram", "concurrent-user", now())
			if err != nil {
				t.Errorf("get-or-create err: %v", err)
				return
			}
			if created {
				createdCount.Add(1)
			}
			ids.Store(sub.ID, true)
		}()
	}
	close(startGate)
	wg.Wait()

	if createdCount.Load() != 1 {
		t.Fatalf("createdCount = %d, want EXACTLY 1", createdCount.Load())
	}
	var distinctIDs int
	ids.Range(func(key, value any) bool {
		distinctIDs++
		return true
	})
	if distinctIDs != 1 {
		t.Fatalf("distinctIDs = %d, want EXACTLY 1", distinctIDs)
	}
}
