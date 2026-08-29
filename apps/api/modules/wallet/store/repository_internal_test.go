package store

import (
	"errors"
	"testing"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/testsupport"
)

func TestReplayAfterIdempotencyRace(t *testing.T) {
	platformStore, err := testsupport.OpenStore(":memory:", "admin", "hash", false)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { platformStore.Close() })
	repo := NewRepository(platformStore)
	now := time.Unix(1700000000, 0).UTC()
	account := Account{ID: "acct-race", OwnerType: OwnerUser, OwnerID: "race", Currency: DefaultCurrency, Status: StatusActive, CreatedAt: now, UpdatedAt: now}
	if err := repo.CreateAccount(account); err != nil {
		t.Fatal(err)
	}
	input := LedgerEntryInput{EntryType: EntryAdjust, AmountDelta: 25, Memo: "race", IdempotencyKey: "race-key", ActorID: "actor-1", ActorName: "Admin"}
	if _, _, err := repo.Mutate(account.ID, input, "winner", now); err != nil {
		t.Fatal(err)
	}
	replayedAccount, replayedEntry, err := repo.replayAfterIdempotencyRace(account.ID, input)
	if err != nil {
		t.Fatalf("replay winner: %v", err)
	}
	if replayedEntry.ID != "winner" || replayedAccount.BalanceTotal != 25 || replayedAccount.Version != 1 {
		t.Fatalf("race replay = account %+v entry %+v", replayedAccount, replayedEntry)
	}
	otherActor := input
	otherActor.ActorID = "actor-2"
	if _, _, err := repo.replayAfterIdempotencyRace(account.ID, otherActor); !errors.Is(err, ErrIdempotencyConflict) {
		t.Fatalf("race payload conflict = %v, want ErrIdempotencyConflict", err)
	}
}
