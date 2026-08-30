// GOAL-020 A-003 F-001 regression: concurrent get-or-create — exactly one
// creator, losers report created=false (no duplicate auto audit), one row.
package store_test

import (
	"sync"
	"testing"

	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
)

func TestGetOrCreateUserAccountConcurrent(t *testing.T) {
	repo := newRepo(t)
	const n = 8
	results := make([]struct {
		id      string
		created bool
		err     error
	}, n)
	var wg sync.WaitGroup
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			acct, created, err := repo.GetOrCreateUserAccount("u-concurrent", now())
			if acct != nil {
				results[i].id = acct.ID
			}
			results[i].created = created
			results[i].err = err
		}(i)
	}
	wg.Wait()

	creates := 0
	ids := map[string]int{}
	for _, r := range results {
		if r.err != nil {
			t.Fatalf("concurrent get-or-create: %v", r.err)
		}
		if r.created {
			creates++
		}
		ids[r.id]++
	}
	if creates != 1 {
		t.Fatalf("created=true count = %d, want exactly 1", creates)
	}
	if len(ids) != 1 {
		t.Fatalf("distinct account ids = %d, want 1", len(ids))
	}

	// Exactly one row in the table.
	accounts, total, err := repo.ListAccounts(store.ListFilter{Page: 1, PageSize: 50})
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(accounts) != 1 {
		t.Fatalf("account rows = %d/%d, want 1/1", len(accounts), total)
	}
}
