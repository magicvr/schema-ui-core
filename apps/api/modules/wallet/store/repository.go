// Package store owns the admin.wallet persistence (S-14 · GOAL-019 D-002
// §1/§2): wallet accounts with the three-balance invariant, the immutable
// ledger (apply-table snapshots) and reconciliation runs. All balance
// mutations run inside the platform transaction boundary (WithTx) with an
// optimistic-lock version; ledger writes carry balance-after snapshots so the
// chain can be replayed by reconciliation.
package store

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pagination"
)

// TxRunner is the platform persistence boundary consumed by the repository.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// Repository owns the wallet domain queries.
type Repository struct {
	runner TxRunner
}

// NewRepository constructs the wallet repository over a platform transaction
// runner.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// escapeLikePattern escapes SQL LIKE metacharacters in user-supplied search
// input so the value matches literally (W13 F-011 · GOAL-013 A-001). Every
// clause using it must pin ESCAPE '\'.
func escapeLikePattern(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
}

// Domain sentinels mapped by the handler to frozen error codes.
var (
	ErrNotFound            = errors.New("wallet account not found")
	ErrOwnerTaken          = errors.New("wallet owner already exists")
	ErrDisabled            = errors.New("wallet account disabled")
	ErrInsufficient        = errors.New("insufficient balance")
	ErrVersionConflict     = errors.New("wallet ledger version conflict")
	ErrIdempotencyConflict = errors.New("idempotency key reused with different payload")
	ErrInvalidEntry        = errors.New("invalid ledger entry")
	errIdempotencyRace     = errors.New("wallet idempotency insert race")
)

// Entry types (D-002 §1 apply table).
const (
	EntryAdjust       = "adjust"
	EntryFreeze       = "freeze"
	EntryUnfreeze     = "unfreeze"
	EntryDeductFrozen = "deduct_frozen"
)

// Account statuses.
const (
	StatusActive   = "active"
	StatusDisabled = "disabled"
)

// Owner types.
const (
	OwnerUser     = "user"
	OwnerBusiness = "business"
	OwnerSystem   = "system"
)

// DefaultCurrency is the v1 single-currency default.
const DefaultCurrency = "CNY"

// Account is one wallet_accounts row.
type Account struct {
	ID               string
	OwnerType        string
	OwnerID          string
	Currency         string
	BalanceTotal     int64
	BalanceAvailable int64
	BalanceFrozen    int64
	Status           string
	Version          int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

// LedgerEntry is one immutable wallet_ledger_entries row.
type LedgerEntry struct {
	ID                 string
	AccountID          string
	EntryType          string
	AmountDelta        int64
	BalanceAfterTotal  int64
	BalanceAfterAvail  int64
	BalanceAfterFrozen int64
	RefType            string
	RefID              string
	IdempotencyKey     string
	Memo               string
	ActorID            string
	ActorName          string
	CreatedAt          time.Time
}

// ReconciliationRun is one wallet_reconciliation_runs row.
type ReconciliationRun struct {
	ID            string
	AccountID     string
	Result        string
	MismatchCount int
	Details       string
	ActorID       string
	CreatedAt     time.Time
}

// Reconciliation results.
const (
	ResultConsistent   = "consistent"
	ResultInconsistent = "inconsistent"
)

// ListFilter carries account list parameters.
type ListFilter struct {
	Q         string
	OwnerType string
	Page      int
	PageSize  int
}

// LedgerEntryInput is the domain input for one balance mutation.
type LedgerEntryInput struct {
	EntryType      string
	AmountDelta    int64
	RefType        string
	RefID          string
	IdempotencyKey string
	Memo           string
	ActorID        string
	ActorName      string
}

// ListAccounts returns wallet accounts ordered by created_at desc.
func (r *Repository) ListAccounts(filter ListFilter) ([]Account, int, error) {
	accounts := []Account{}
	total := 0
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		where := "WHERE 1=1"
		args := []any{}
		if filter.Q != "" {
			// R4 S2: portable case-insensitive search (sqlite LIKE is CI for
			// ASCII; postgres LIKE is CS). LOWER(...) LIKE LOWER(?) restores
			// parity on both dialects. W13 F-011 (GOAL-013 A-001): q is escaped
			// and the pattern pins ESCAPE '\' so user input cannot inject %
			// / _ wildcards (scan amplification + character-by-character
			// probing).
			where += " AND (LOWER(owner_id) LIKE LOWER(?) ESCAPE '\\' OR LOWER(owner_type) LIKE LOWER(?) ESCAPE '\\' OR LOWER(currency) LIKE LOWER(?) ESCAPE '\\')"
			like := "%" + escapeLikePattern(filter.Q) + "%"
			args = append(args, like, like, like)
		}
		if filter.OwnerType != "" {
			where += " AND owner_type = ?"
			args = append(args, filter.OwnerType)
		}
		if err := tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM wallet_accounts "+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count wallet accounts: %w", err)
		}
		page := filter.Page
		if page < 1 {
			page = 1
		}
		pageSize := filter.PageSize
		if pageSize < 1 {
			pageSize = 20
		}
		query := "SELECT id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at FROM wallet_accounts " + where + " ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?"
		rows, err := tx.Query(context.Background(), query, append(args, pageSize, pagination.Offset(page, pageSize, total))...)
		if err != nil {
			return fmt.Errorf("list wallet accounts: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var a Account
			var created, updated int64
			if err := rows.Scan(&a.ID, &a.OwnerType, &a.OwnerID, &a.Currency, &a.BalanceTotal, &a.BalanceAvailable, &a.BalanceFrozen, &a.Status, &a.Version, &created, &updated); err != nil {
				return fmt.Errorf("scan wallet account: %w", err)
			}
			a.CreatedAt = time.Unix(created, 0)
			a.UpdatedAt = time.Unix(updated, 0)
			accounts = append(accounts, a)
		}
		return rows.Err()
	})
	return accounts, total, err
}

// GetAccount returns one account by id.
func (r *Repository) GetAccount(id string) (*Account, error) {
	var a Account
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var created, updated int64
		err := tx.QueryRow(context.Background(),
			`SELECT id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at FROM wallet_accounts WHERE id = ?`,
			id,
		).Scan(&a.ID, &a.OwnerType, &a.OwnerID, &a.Currency, &a.BalanceTotal, &a.BalanceAvailable, &a.BalanceFrozen, &a.Status, &a.Version, &created, &updated)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get wallet account: %w", err)
		}
		a.CreatedAt = time.Unix(created, 0)
		a.UpdatedAt = time.Unix(updated, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// CreateAccount inserts a zero-balance account. Duplicate owner → ErrOwnerTaken.
func (r *Repository) CreateAccount(a Account) error {
	return r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO wallet_accounts (id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at)
			 VALUES (?, ?, ?, ?, 0, 0, 0, ?, 0, ?, ?)`,
			a.ID, a.OwnerType, a.OwnerID, a.Currency, a.Status, a.CreatedAt.Unix(), a.UpdatedAt.Unix(),
		)
		if err != nil {
			if isUniqueViolation(err) {
				return ErrOwnerTaken
			}
			return fmt.Errorf("create wallet account: %w", err)
		}
		return nil
	})
}

// GetUserAccountByOwner is a read-only lookup (W15-F11). Missing → ErrNotFound.
func (r *Repository) GetUserAccountByOwner(ownerID string) (*Account, error) {
	var a Account
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var created, updated int64
		err := tx.QueryRow(context.Background(),
			`SELECT id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at FROM wallet_accounts WHERE owner_type = ? AND owner_id = ? AND currency = ?`,
			OwnerUser, ownerID, DefaultCurrency,
		).Scan(&a.ID, &a.OwnerType, &a.OwnerID, &a.Currency, &a.BalanceTotal, &a.BalanceAvailable, &a.BalanceFrozen, &a.Status, &a.Version, &created, &updated)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("get wallet account by owner: %w", err)
		}
		a.CreatedAt = time.Unix(created, 0)
		a.UpdatedAt = time.Unix(updated, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// GetOrCreateUserAccount returns the user account for ownerID, creating a
// zero-balance account when absent (GOAL-020 D-001 §1 get-or-create). The
// UNIQUE(owner_type, owner_id, currency) constraint makes concurrent creates
// safe: an INSERT conflict falls back to re-reading the existing row in a
// FRESH transaction — on postgres a failed INSERT aborts its own transaction,
// so the loser must never re-read inside it (W9 F-001).
func (r *Repository) GetOrCreateUserAccount(ownerID string, now time.Time) (*Account, bool, error) {
	// Fast path: the account already exists.
	if existing, err := r.GetUserAccountByOwner(ownerID); err == nil {
		return existing, false, nil
	} else if !errors.Is(err, ErrNotFound) {
		return nil, false, fmt.Errorf("get wallet account by owner: %w", err)
	}
	// Auto-create with a time-ordered id (millisecond prefix + random
	// suffix — same convention as the module newID; GOAL-020 A-003 F-003).
	randBytes := make([]byte, 12)
	if _, err := rand.Read(randBytes); err != nil {
		return nil, false, fmt.Errorf("auto-create wallet id: %w", err)
	}
	id := fmt.Sprintf("%016x%s", now.UnixMilli(), hex.EncodeToString(randBytes))
	insertErr := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		_, err := tx.Exec(context.Background(),
			`INSERT INTO wallet_accounts (id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at)
			 VALUES (?, ?, ?, ?, 0, 0, 0, ?, 0, ?, ?)`,
			id, OwnerUser, ownerID, DefaultCurrency, StatusActive, now.Unix(), now.Unix(),
		)
		return err
	})
	if insertErr != nil {
		if !isUniqueViolation(insertErr) {
			return nil, false, fmt.Errorf("auto-create wallet account: %w", insertErr)
		}
		// Concurrent create won: re-read the existing row in a fresh
		// transaction. The loser must NOT report a create (GOAL-020 A-003
		// F-001: no duplicate wallet.account-create on the shared row).
		existing, err := r.GetUserAccountByOwner(ownerID)
		if err != nil {
			return nil, false, fmt.Errorf("re-read wallet account after create conflict: %w", err)
		}
		return existing, false, nil
	}
	a := Account{ID: id, OwnerType: OwnerUser, OwnerID: ownerID, Currency: DefaultCurrency, Status: StatusActive, CreatedAt: now, UpdatedAt: now}
	return &a, true, nil
}

// UpdateStatus flips account status (active/disabled) with the optimistic
// lock: the caller passes the version observed when loading the row.
func (r *Repository) UpdateStatus(id, status string, version int64, now time.Time) (*Account, error) {
	var a Account
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		res, err := tx.Exec(context.Background(),
			`UPDATE wallet_accounts SET status = ?, version = version + 1, updated_at = ? WHERE id = ? AND version = ?`,
			status, now.Unix(), id, version,
		)
		if err != nil {
			return fmt.Errorf("update wallet status: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			// Either the account does not exist or the version moved.
			var exists int
			if err := tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM wallet_accounts WHERE id = ?", id).Scan(&exists); err != nil {
				return fmt.Errorf("check wallet account: %w", err)
			}
			if exists == 0 {
				return ErrNotFound
			}
			return ErrVersionConflict
		}
		var created, updated int64
		if err := tx.QueryRow(context.Background(),
			`SELECT id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at FROM wallet_accounts WHERE id = ?`,
			id,
		).Scan(&a.ID, &a.OwnerType, &a.OwnerID, &a.Currency, &a.BalanceTotal, &a.BalanceAvailable, &a.BalanceFrozen, &a.Status, &a.Version, &created, &updated); err != nil {
			return fmt.Errorf("reload wallet account: %w", err)
		}
		a.CreatedAt = time.Unix(created, 0)
		a.UpdatedAt = time.Unix(updated, 0)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &a, nil
}

// Apply computes the balance-after triplets for one entry per the D-002 §1
// apply table. It returns ErrInvalidEntry for malformed entries and
// ErrInsufficient when the mutation would drive any balance negative.
func Apply(prev Account, in LedgerEntryInput) (total, available, frozen int64, err error) {
	total = prev.BalanceTotal
	available = prev.BalanceAvailable
	frozen = prev.BalanceFrozen
	switch in.EntryType {
	case EntryAdjust:
		if in.AmountDelta == 0 {
			return 0, 0, 0, ErrInvalidEntry
		}
		total += in.AmountDelta
		available += in.AmountDelta
	case EntryFreeze:
		if in.AmountDelta <= 0 {
			return 0, 0, 0, ErrInvalidEntry
		}
		available -= in.AmountDelta
		frozen += in.AmountDelta
	case EntryUnfreeze:
		if in.AmountDelta <= 0 {
			return 0, 0, 0, ErrInvalidEntry
		}
		available += in.AmountDelta
		frozen -= in.AmountDelta
	case EntryDeductFrozen:
		// A-008 F-001: deduct from the frozen bucket atomically — total and
		// frozen drop by d, available stays untouched (the pre-authorized
		// money is consumed, never re-exposed as available).
		if in.AmountDelta <= 0 {
			return 0, 0, 0, ErrInvalidEntry
		}
		total -= in.AmountDelta
		frozen -= in.AmountDelta
	default:
		return 0, 0, 0, ErrInvalidEntry
	}
	if total < 0 || available < 0 || frozen < 0 {
		return 0, 0, 0, ErrInsufficient
	}
	if total != available+frozen {
		return 0, 0, 0, ErrInvalidEntry
	}
	return total, available, frozen, nil
}

// Mutate applies one balance mutation atomically: idempotency check → account
// load → apply-table validation → optimistic-lock update → immutable ledger
// insert. Returns the updated account and the recorded entry.
func (r *Repository) Mutate(id string, in LedgerEntryInput, entryID string, now time.Time) (*Account, *LedgerEntry, error) {
	var account Account
	var entry LedgerEntry
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		// Idempotency: same account + same key + same payload → return the
		// existing entry; same key + different payload → conflict. Lookups
		// always carry the account id (D-002 v1.1.0 §1: no bare-key reads).
		if in.IdempotencyKey != "" {
			existing, err := readIdempotentEntry(tx, id, in.IdempotencyKey)
			if err == nil {
				if sameIdempotencyPayload(existing, in) {
					entry = existing
					// Return the current account too so the caller can report
					// the idempotent replay without a second read.
					acct, err := readAccount(tx, id)
					if err != nil {
						return fmt.Errorf("reload wallet account: %w", err)
					}
					account = acct
					return nil
				}
				return ErrIdempotencyConflict
			}
			if !errors.Is(err, kernel.ErrNoRows) {
				return fmt.Errorf("check idempotency key: %w", err)
			}
		}
		// Load the account inside the transaction.
		var cr, up int64
		err := tx.QueryRow(context.Background(),
			`SELECT id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at FROM wallet_accounts WHERE id = ?`,
			id,
		).Scan(&account.ID, &account.OwnerType, &account.OwnerID, &account.Currency, &account.BalanceTotal, &account.BalanceAvailable, &account.BalanceFrozen, &account.Status, &account.Version, &cr, &up)
		if errors.Is(err, kernel.ErrNoRows) {
			return ErrNotFound
		}
		if err != nil {
			return fmt.Errorf("load wallet account: %w", err)
		}
		account.CreatedAt = time.Unix(cr, 0)
		account.UpdatedAt = time.Unix(up, 0)
		if account.Status != StatusActive {
			return ErrDisabled
		}
		total, available, frozen, err := Apply(account, in)
		if err != nil {
			return err
		}
		res, err := tx.Exec(context.Background(),
			`UPDATE wallet_accounts SET balance_total = ?, balance_available = ?, balance_frozen = ?, version = version + 1, updated_at = ? WHERE id = ? AND version = ?`,
			total, available, frozen, now.Unix(), id, account.Version,
		)
		if err != nil {
			return fmt.Errorf("update wallet balances: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return ErrVersionConflict
		}
		_, err = tx.Exec(context.Background(),
			`INSERT INTO wallet_ledger_entries (id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, ref_type, ref_id, idempotency_key, memo, actor_id, actor_name, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			entryID, id, in.EntryType, in.AmountDelta, total, available, frozen,
			nullIfEmpty(in.RefType), nullIfEmpty(in.RefID), nullIfEmpty(in.IdempotencyKey),
			in.Memo, in.ActorID, in.ActorName, now.Unix(),
		)
		if err != nil {
			if isUniqueViolation(err) {
				// Roll back our balance update before re-reading the winning
				// operation outside this transaction.
				return errIdempotencyRace
			}
			return fmt.Errorf("insert wallet ledger entry: %w", err)
		}
		account.BalanceTotal = total
		account.BalanceAvailable = available
		account.BalanceFrozen = frozen
		account.Version++
		account.UpdatedAt = now
		entry = LedgerEntry{
			ID: entryID, AccountID: id, EntryType: in.EntryType, AmountDelta: in.AmountDelta,
			BalanceAfterTotal: total, BalanceAfterAvail: available, BalanceAfterFrozen: frozen,
			RefType: in.RefType, RefID: in.RefID, IdempotencyKey: in.IdempotencyKey,
			Memo: in.Memo, ActorID: in.ActorID, ActorName: in.ActorName, CreatedAt: now,
		}
		return nil
	})
	if errors.Is(err, errIdempotencyRace) && in.IdempotencyKey != "" {
		return r.replayAfterIdempotencyRace(id, in)
	}
	if err != nil {
		return nil, nil, err
	}
	return &account, &entry, nil
}

type rowQueryer interface {
	QueryRow(context.Context, string, ...any) kernel.Row
}

func readIdempotentEntry(q rowQueryer, accountID, key string) (LedgerEntry, error) {
	var entry LedgerEntry
	var created int64
	var refType, refID, idemKey sql.NullString
	err := q.QueryRow(context.Background(),
		`SELECT id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, ref_type, ref_id, idempotency_key, memo, actor_id, actor_name, created_at
		 FROM wallet_ledger_entries WHERE account_id = ? AND idempotency_key = ?`,
		accountID, key,
	).Scan(&entry.ID, &entry.AccountID, &entry.EntryType, &entry.AmountDelta, &entry.BalanceAfterTotal, &entry.BalanceAfterAvail, &entry.BalanceAfterFrozen, &refType, &refID, &idemKey, &entry.Memo, &entry.ActorID, &entry.ActorName, &created)
	if err != nil {
		return LedgerEntry{}, err
	}
	entry.RefType = refType.String
	entry.RefID = refID.String
	entry.IdempotencyKey = idemKey.String
	entry.CreatedAt = time.Unix(created, 0)
	return entry, nil
}

func readAccount(q rowQueryer, id string) (Account, error) {
	var account Account
	var created, updated int64
	err := q.QueryRow(context.Background(),
		`SELECT id, owner_type, owner_id, currency, balance_total, balance_available, balance_frozen, status, version, created_at, updated_at FROM wallet_accounts WHERE id = ?`,
		id,
	).Scan(&account.ID, &account.OwnerType, &account.OwnerID, &account.Currency, &account.BalanceTotal, &account.BalanceAvailable, &account.BalanceFrozen, &account.Status, &account.Version, &created, &updated)
	if err != nil {
		return Account{}, err
	}
	account.CreatedAt = time.Unix(created, 0)
	account.UpdatedAt = time.Unix(updated, 0)
	return account, nil
}

func sameIdempotencyPayload(entry LedgerEntry, in LedgerEntryInput) bool {
	return entry.EntryType == in.EntryType &&
		entry.AmountDelta == in.AmountDelta &&
		entry.Memo == in.Memo &&
		entry.RefType == in.RefType &&
		entry.RefID == in.RefID &&
		entry.ActorID == in.ActorID
}

func (r *Repository) replayAfterIdempotencyRace(accountID string, in LedgerEntryInput) (*Account, *LedgerEntry, error) {
	var account Account
	var entry LedgerEntry
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		existing, err := readIdempotentEntry(tx, accountID, in.IdempotencyKey)
		if err != nil || !sameIdempotencyPayload(existing, in) {
			return ErrIdempotencyConflict
		}
		current, err := readAccount(tx, accountID)
		if err != nil {
			return ErrIdempotencyConflict
		}
		entry = existing
		account = current
		return nil
	})
	if err != nil {
		return nil, nil, err
	}
	return &account, &entry, nil
}

// ListEntries returns the ledger of one account, newest first. entryType, when
// non-empty, filters to one entry type; q searches memo/ref fields (W14 F-07).
func (r *Repository) ListEntries(accountID, entryType, q string, page, pageSize int) ([]LedgerEntry, int, error) {
	entries := []LedgerEntry{}
	total := 0
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		where := "WHERE account_id = ?"
		args := []any{accountID}
		if entryType != "" {
			where += " AND entry_type = ?"
			args = append(args, entryType)
		}
		if q != "" {
			where += " AND (lower(COALESCE(memo,'')) LIKE '%' || CAST(? AS TEXT) || '%' ESCAPE '\\' OR lower(COALESCE(ref_type,'')) LIKE '%' || CAST(? AS TEXT) || '%' ESCAPE '\\' OR lower(COALESCE(ref_id,'')) LIKE '%' || CAST(? AS TEXT) || '%' ESCAPE '\\')"
			lq := escapeLikePattern(strings.ToLower(q))
			args = append(args, lq, lq, lq)
		}
		if err := tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM wallet_ledger_entries "+where, args...).Scan(&total); err != nil {
			return fmt.Errorf("count wallet entries: %w", err)
		}
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 20
		}
		queryArgs := append(append([]any{}, args...), pageSize, pagination.Offset(page, pageSize, total))
		rows, err := tx.Query(context.Background(),
			`SELECT id, account_id, entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen, ref_type, ref_id, idempotency_key, memo, actor_id, actor_name, created_at
			 FROM wallet_ledger_entries `+where+` ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`,
			queryArgs...,
		)
		if err != nil {
			return fmt.Errorf("list wallet entries: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var e LedgerEntry
			var created int64
			var refType, refID, idemKey sql.NullString
			if err := rows.Scan(&e.ID, &e.AccountID, &e.EntryType, &e.AmountDelta, &e.BalanceAfterTotal, &e.BalanceAfterAvail, &e.BalanceAfterFrozen, &refType, &refID, &idemKey, &e.Memo, &e.ActorID, &e.ActorName, &created); err != nil {
				return fmt.Errorf("scan wallet entry: %w", err)
			}
			e.RefType = refType.String
			e.RefID = refID.String
			e.IdempotencyKey = idemKey.String
			e.CreatedAt = time.Unix(created, 0)
			entries = append(entries, e)
		}
		return rows.Err()
	})
	return entries, total, err
}

// ReconcileRun checks the ledger chain of one account (or every account when
// accountID is empty) against the D-002 §1 replay rules: per-entry apply
// equality, snapshot invariant, and last-snapshot == current balances. The
// run is persisted as a reconciliation_runs row.
func (r *Repository) ReconcileRun(accountID, runID, actorID string, now time.Time) (*ReconciliationRun, error) {
	var run *ReconciliationRun
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		var err error
		run, err = r.ReconcileOnceTx(context.Background(), tx, accountID, runID, actorID, now)
		return err
	})
	return run, err
}

// ReconcileOnceTx persists one idempotent reconciliation run inside the
// caller's transaction. Job consumers use this to atomically commit the
// wallet run and the Job succeeded transition.
func (r *Repository) ReconcileOnceTx(ctx context.Context, tx kernel.Tx, accountID, runID, actorID string, now time.Time) (*ReconciliationRun, error) {
	if tx == nil || runID == "" || actorID == "" {
		return nil, ErrInvalidEntry
	}
	if existing, err := getReconciliationRunTx(ctx, tx, runID); err == nil {
		return existing, nil
	} else if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	type mismatch struct {
		AccountID string `json:"accountId"`
		Reason    string `json:"reason"`
	}
	var mismatches []mismatch
	ids := []string{}
	if accountID == "" {
		rows, err := tx.Query(ctx, "SELECT id FROM wallet_accounts ORDER BY id")
		if err != nil {
			return nil, fmt.Errorf("list wallet accounts for reconcile: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				return nil, fmt.Errorf("scan wallet account id: %w", err)
			}
			ids = append(ids, id)
		}
		if err := rows.Err(); err != nil {
			return nil, err
		}
	} else {
		var exists int
		if err := tx.QueryRow(ctx, "SELECT COUNT(*) FROM wallet_accounts WHERE id = ?", accountID).Scan(&exists); err != nil {
			return nil, fmt.Errorf("check wallet account: %w", err)
		}
		if exists == 0 {
			return nil, ErrNotFound
		}
		ids = append(ids, accountID)
	}
	for _, id := range ids {
		if reason, ok := checkAccountChain(tx, id); !ok {
			mismatches = append(mismatches, mismatch{AccountID: id, Reason: reason})
		}
	}
	result := ResultConsistent
	if len(mismatches) > 0 {
		result = ResultInconsistent
	}
	details := "{}"
	if len(mismatches) > 0 {
		raw, err := json.Marshal(map[string]any{"mismatches": mismatches})
		if err == nil {
			details = string(raw)
		}
	}
	var acctID any
	if accountID != "" {
		acctID = accountID
	}
	if _, err := tx.Exec(ctx,
		`INSERT INTO wallet_reconciliation_runs (id, account_id, result, mismatch_count, details, actor_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		runID, acctID, result, len(mismatches), details, actorID, now.Unix(),
	); err != nil {
		return nil, fmt.Errorf("insert wallet reconciliation run: %w", err)
	}
	run := &ReconciliationRun{ID: runID, AccountID: accountID, Result: result, MismatchCount: len(mismatches), Details: details, ActorID: actorID, CreatedAt: now.UTC()}
	return run, nil
}

func getReconciliationRunTx(ctx context.Context, tx kernel.Tx, id string) (*ReconciliationRun, error) {
	var run ReconciliationRun
	var accountID sql.NullString
	var created int64
	err := tx.QueryRow(ctx,
		`SELECT id, account_id, result, mismatch_count, details, actor_id, created_at FROM wallet_reconciliation_runs WHERE id = ?`, id,
	).Scan(&run.ID, &accountID, &run.Result, &run.MismatchCount, &run.Details, &run.ActorID, &created)
	if errors.Is(err, kernel.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get wallet reconciliation run: %w", err)
	}
	run.AccountID = accountID.String
	run.CreatedAt = time.Unix(created, 0).UTC()
	return &run, nil
}

// checkAccountChain replays one account's ledger inside the given transaction.
// Chain order is (created_at ASC, id ASC); the first entry starts from
// (0,0,0); the last snapshot must equal the account's current balances; every
// snapshot must satisfy the invariant. Returns a human-readable reason when
// the chain is inconsistent.
func checkAccountChain(tx kernel.Tx, accountID string) (string, bool) {
	var total, available, frozen int64
	rows, err := tx.Query(context.Background(),
		`SELECT entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen FROM wallet_ledger_entries WHERE account_id = ? ORDER BY created_at ASC, id ASC`,
		accountID,
	)
	if err != nil {
		return "ledger read failed", false
	}
	defer rows.Close()
	first := true
	for rows.Next() {
		var entryType string
		var delta, afterTotal, afterAvail, afterFrozen int64
		if err := rows.Scan(&entryType, &delta, &afterTotal, &afterAvail, &afterFrozen); err != nil {
			return "ledger scan failed", false
		}
		if afterTotal != afterAvail+afterFrozen {
			return "snapshot invariant violated", false
		}
		prev := Account{BalanceTotal: total, BalanceAvailable: available, BalanceFrozen: frozen}
		if first {
			prev = Account{BalanceTotal: 0, BalanceAvailable: 0, BalanceFrozen: 0}
			first = false
		}
		wantTotal, wantAvail, wantFrozen, err := Apply(prev, LedgerEntryInput{EntryType: entryType, AmountDelta: delta})
		if err != nil {
			return "replay apply failed: " + err.Error(), false
		}
		if wantTotal != afterTotal || wantAvail != afterAvail || wantFrozen != afterFrozen {
			return "snapshot does not match replay", false
		}
		total, available, frozen = afterTotal, afterAvail, afterFrozen
	}
	if err := rows.Err(); err != nil {
		return "ledger iteration failed", false
	}
	var curTotal, curAvail, curFrozen int64
	if err := tx.QueryRow(context.Background(),
		`SELECT balance_total, balance_available, balance_frozen FROM wallet_accounts WHERE id = ?`,
		accountID,
	).Scan(&curTotal, &curAvail, &curFrozen); err != nil {
		return "account read failed", false
	}
	if curTotal != curAvail+curFrozen {
		return "account invariant violated", false
	}
	if total != curTotal || available != curAvail || frozen != curFrozen {
		return "last snapshot != account balances", false
	}
	return "", true
}

// ListReconcileRuns returns reconciliation runs, newest first.
func (r *Repository) ListReconcileRuns(page, pageSize int) ([]ReconciliationRun, int, error) {
	runs := []ReconciliationRun{}
	total := 0
	err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
		if err := tx.QueryRow(context.Background(), "SELECT COUNT(*) FROM wallet_reconciliation_runs").Scan(&total); err != nil {
			return fmt.Errorf("count reconciliation runs: %w", err)
		}
		if page < 1 {
			page = 1
		}
		if pageSize < 1 {
			pageSize = 20
		}
		rows, err := tx.Query(context.Background(),
			`SELECT id, account_id, result, mismatch_count, details, actor_id, created_at FROM wallet_reconciliation_runs ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`,
			pageSize, pagination.Offset(page, pageSize, total),
		)
		if err != nil {
			return fmt.Errorf("list reconciliation runs: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var run ReconciliationRun
			var created int64
			var acctID sql.NullString
			if err := rows.Scan(&run.ID, &acctID, &run.Result, &run.MismatchCount, &run.Details, &run.ActorID, &created); err != nil {
				return fmt.Errorf("scan reconciliation run: %w", err)
			}
			run.AccountID = acctID.String
			run.CreatedAt = time.Unix(created, 0)
			runs = append(runs, run)
		}
		return rows.Err()
	})
	return runs, total, err
}

// isUniqueViolation is dialect-agnostic (W9 F-001): the SQLite-only message
// match previously never fired on postgres (SQLSTATE 23505), silently breaking
// ErrOwnerTaken mapping, the get-or-create conflict fallback and the ledger
// idempotency-race replay on the postgres dialect.
func isUniqueViolation(err error) bool {
	return kernel.IsUniqueViolation(err)
}

func contains(s, sub string) bool {
	return len(sub) == 0 || (len(s) >= len(sub) && (s == sub || len(sub) > 0 && (len(s) > len(sub) && indexOf(s, sub) >= 0)))
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}

func nullIfEmpty(s string) any {
	if s == "" {
		return nil
	}
	return s
}
