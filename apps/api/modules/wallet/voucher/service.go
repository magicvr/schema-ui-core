package voucher

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
)

// SubjectVerifier checks whether a subject exists.
type SubjectVerifier interface {
	SubjectExists(ctx context.Context, id string) (bool, error)
}

// Service coordinates voucher lifecycle operations and atomic single-transaction redemption.
type Service struct {
	runner     walletstore.TxRunner
	walletRepo *walletstore.Repository
	subjects   SubjectVerifier
}

// NewService constructs a voucher Service.
func NewService(runner walletstore.TxRunner, walletRepo *walletstore.Repository, subjects SubjectVerifier) *Service {
	return &Service{
		runner:     runner,
		walletRepo: walletRepo,
		subjects:   subjects,
	}
}

// GenerateBatch creates a batch of vouchers, persisting their SHA-256 hashes and returning
// one-time plaintext codes.
func (s *Service) GenerateBatch(ctx context.Context, batchID string, count int, amount int64, currency string, expiresAt *time.Time, now time.Time) ([]GeneratedVoucher, error) {
	batchID = strings.TrimSpace(batchID)
	currency = strings.TrimSpace(currency)
	if batchID == "" || count <= 0 || count > 1000 || amount <= 0 {
		return nil, ErrInvalidInput
	}
	if currency != "" && currency != walletstore.DefaultCurrency {
		return nil, ErrCurrencyMismatch
	}
	if currency == "" {
		currency = walletstore.DefaultCurrency
	}

	generated := make([]GeneratedVoucher, count)
	for i := 0; i < count; i++ {
		code, prefix, hash, err := GenerateCode()
		if err != nil {
			return nil, err
		}
		id, err := newID(now)
		if err != nil {
			return nil, err
		}
		v := Voucher{
			ID:         id,
			BatchID:    batchID,
			CodeHash:   hash,
			CodePrefix: prefix,
			Amount:     amount,
			Currency:   currency,
			Status:     StatusUnused,
			ExpiresAt:  expiresAt,
			CreatedAt:  now.Truncate(time.Second),
			UpdatedAt:  now.Truncate(time.Second),
		}
		generated[i] = GeneratedVoucher{
			Voucher: v,
			Code:    code,
		}
	}

	err := s.runner.Run(ctx, func(tx kernel.Tx) error {
		// A-005 F-004 (0065): batch_id is a first-class unique identity — the
		// same tx registers the batch BEFORE its codes land, so a repeated id
		// (0 affected rows) aborts the whole run and never mixes two lists.
		res, err := tx.Exec(ctx,
			`INSERT INTO voucher_batches (batch_id, created_at, updated_at) VALUES (?, ?, ?)
			 ON CONFLICT (batch_id) DO NOTHING`,
			batchID, now.Unix(), now.Unix(),
		)
		if err != nil {
			return fmt.Errorf("register voucher batch: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return ErrVoucherBatchExists
		}
		for _, g := range generated {
			v := g.Voucher
			var exp sql.NullInt64
			if v.ExpiresAt != nil {
				exp = sql.NullInt64{Int64: v.ExpiresAt.Unix(), Valid: true}
			}
			_, err := tx.Exec(ctx,
				`INSERT INTO vouchers (id, batch_id, code_hash, code_prefix, amount, currency, status, expires_at, created_at, updated_at)
				 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
				v.ID, v.BatchID, v.CodeHash, v.CodePrefix, v.Amount, v.Currency, string(v.Status), exp, v.CreatedAt.Unix(), v.UpdatedAt.Unix(),
			)
			if err != nil {
				return fmt.Errorf("insert voucher: %w", err)
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return generated, nil
}

// Redeem performs atomic, idempotent single-transaction voucher redemption
// into an external subject's ledger (owner_type=subject).
func (s *Service) Redeem(ctx context.Context, subjectID string, code string, now time.Time) (*RedeemResult, error) {
	subjectID = strings.TrimSpace(subjectID)
	code = strings.TrimSpace(code)
	if subjectID == "" || code == "" {
		return nil, ErrInvalidInput
	}

	// External subject gate: unknown subject cannot redeem or auto-open account.
	if s.subjects != nil {
		exists, err := s.subjects.SubjectExists(ctx, subjectID)
		if err != nil {
			return nil, fmt.Errorf("verify subject: %w", err)
		}
		if !exists {
			return nil, ErrSubjectNotFound
		}
	}

	return s.redeemInto(ctx, code, now, subjectID, "Subject", func(tx kernel.Tx) (*walletstore.Account, error) {
		account, _, err := s.walletRepo.GetOrCreateSubjectAccountInTx(tx, subjectID, now)
		if err != nil {
			return nil, fmt.Errorf("wallet account get-or-create: %w", err)
		}
		return account, nil
	})
}

// RedeemForUser credits the Admin user's own owner_type=user ledger
// (VP-029 R5 · GOAL-005 D-002). It MUST NOT call Redeem(subjectID) and MUST
// NOT open a subject account.
func (s *Service) RedeemForUser(ctx context.Context, userID, actorName, code string, now time.Time) (*RedeemResult, error) {
	userID = strings.TrimSpace(userID)
	code = strings.TrimSpace(code)
	if userID == "" || code == "" {
		return nil, ErrInvalidInput
	}
	if strings.TrimSpace(actorName) == "" {
		actorName = "User"
	}
	return s.redeemInto(ctx, code, now, userID, actorName, func(tx kernel.Tx) (*walletstore.Account, error) {
		account, _, err := s.walletRepo.GetOrCreateUserAccountInTx(tx, userID, now)
		if err != nil {
			return nil, fmt.Errorf("user wallet account get-or-create: %w", err)
		}
		return account, nil
	})
}

type accountOpener func(tx kernel.Tx) (*walletstore.Account, error)

func (s *Service) redeemInto(ctx context.Context, code string, now time.Time, redeemedByID, actorName string, open accountOpener) (*RedeemResult, error) {
	codeHash := HashCode(code)

	var result RedeemResult
	err := s.runner.Run(ctx, func(tx kernel.Tx) error {
		var (
			id, batchID, storedHash, prefix, currency, status string
			amount                                            int64
			expiresAt, redeemedAt                             sql.NullInt64
			redeemedBy                                        sql.NullString
			cr, up                                            int64
		)
		row := tx.QueryRow(ctx,
			`SELECT id, batch_id, code_hash, code_prefix, amount, currency, status, expires_at, redeemed_by, redeemed_at, created_at, updated_at
			 FROM vouchers WHERE code_hash = ?`,
			codeHash,
		)
		if err := row.Scan(&id, &batchID, &storedHash, &prefix, &amount, &currency, &status, &expiresAt, &redeemedBy, &redeemedAt, &cr, &up); err != nil {
			if errors.Is(err, kernel.ErrNoRows) {
				return ErrNotFound
			}
			return fmt.Errorf("query voucher: %w", err)
		}

		// Constant-time compare for timing safety.
		if !ConstantTimeCompare(storedHash, codeHash) {
			return ErrNotFound
		}

		if status == string(StatusRedeemed) {
			return ErrVoucherAlreadyRedeemed
		}
		if status == string(StatusVoid) {
			return ErrVoucherVoid
		}
		if status != string(StatusUnused) {
			return ErrVoucherInvalid
		}
		if currency != "" && currency != walletstore.DefaultCurrency {
			return ErrCurrencyMismatch
		}
		if expiresAt.Valid && expiresAt.Int64 > 0 && expiresAt.Int64 < now.Unix() {
			return ErrVoucherExpired
		}

		// CAS: only transition from 'unused' to 'redeemed'.
		res, err := tx.Exec(ctx,
			`UPDATE vouchers SET status = ?, redeemed_by = ?, redeemed_at = ?, updated_at = ? WHERE id = ? AND status = ?`,
			string(StatusRedeemed), redeemedByID, now.Unix(), now.Unix(), id, string(StatusUnused),
		)
		if err != nil {
			return fmt.Errorf("update voucher status: %w", err)
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			// Race condition: another concurrent caller redeemed it first.
			return ErrVoucherConflict
		}

		account, err := open(tx)
		if err != nil {
			return err
		}

		entryID, err := walletstore.NewEntryID(now)
		if err != nil {
			return fmt.Errorf("generate entry id: %w", err)
		}
		in := walletstore.LedgerEntryInput{
			EntryType:      walletstore.EntryAdjust,
			AmountDelta:    amount,
			RefType:        "voucher",
			RefID:          id,
			IdempotencyKey: id, // strictly idempotent on voucher ID
			Memo:           fmt.Sprintf("Voucher redeem %s", prefix),
			ActorID:        redeemedByID,
			ActorName:      actorName,
		}
		updatedAcct, entry, err := s.walletRepo.MutateInTx(tx, account.ID, in, entryID, now)
		if err != nil {
			return fmt.Errorf("credit wallet account: %w", err)
		}

		result = RedeemResult{
			VoucherID:  id,
			BatchID:    batchID,
			CodePrefix: prefix,
			Amount:     amount,
			Currency:   currency,
			AccountID:  updatedAcct.ID,
			EntryID:    entry.ID,
			Balance:    updatedAcct.BalanceTotal,
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return &result, nil
}

// VoidVoucher invalidates an unused voucher so it can never be redeemed.
func (s *Service) VoidVoucher(ctx context.Context, voucherID string, now time.Time) error {
	voucherID = strings.TrimSpace(voucherID)
	if voucherID == "" {
		return ErrInvalidInput
	}
	return s.runner.Run(ctx, func(tx kernel.Tx) error {
		var status string
		row := tx.QueryRow(ctx, `SELECT status FROM vouchers WHERE id = ?`, voucherID)
		if err := row.Scan(&status); err != nil {
			if errors.Is(err, kernel.ErrNoRows) {
				return ErrNotFound
			}
			return err
		}
		if status == string(StatusRedeemed) {
			return ErrVoucherAlreadyRedeemed
		}
		if status == string(StatusVoid) {
			return nil // idempotent
		}
		res, err := tx.Exec(ctx,
			`UPDATE vouchers SET status = ?, updated_at = ? WHERE id = ? AND status = ?`,
			string(StatusVoid), now.Unix(), voucherID, string(StatusUnused),
		)
		if err != nil {
			return err
		}
		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}
		if affected == 0 {
			return ErrVoucherConflict
		}
		return nil
	})
}

// GetVoucher fetches a voucher by its ID.
func (s *Service) GetVoucher(ctx context.Context, voucherID string) (*Voucher, error) {
	voucherID = strings.TrimSpace(voucherID)
	if voucherID == "" {
		return nil, ErrNotFound
	}
	var v Voucher
	var exp, redAt sql.NullInt64
	var redBy sql.NullString
	var cr, up int64
	err := s.runner.Run(ctx, func(tx kernel.Tx) error {
		row := tx.QueryRow(ctx,
			`SELECT id, batch_id, code_prefix, amount, currency, status, expires_at, redeemed_by, redeemed_at, created_at, updated_at
			 FROM vouchers WHERE id = ?`,
			voucherID,
		)
		if err := row.Scan(&v.ID, &v.BatchID, &v.CodePrefix, &v.Amount, &v.Currency, &v.Status, &exp, &redBy, &redAt, &cr, &up); err != nil {
			if errors.Is(err, kernel.ErrNoRows) {
				return ErrNotFound
			}
			return err
		}
		if exp.Valid && exp.Int64 > 0 {
			t := time.Unix(exp.Int64, 0).UTC()
			v.ExpiresAt = &t
		}
		if redBy.Valid {
			v.RedeemedBy = &redBy.String
		}
		if redAt.Valid && redAt.Int64 > 0 {
			t := time.Unix(redAt.Int64, 0).UTC()
			v.RedeemedAt = &t
		}
		v.CreatedAt = time.Unix(cr, 0).UTC()
		v.UpdatedAt = time.Unix(up, 0).UTC()
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &v, nil
}

// ListVouchers returns a paginated list of vouchers, optionally filtered by batchID and status.
func (s *Service) ListVouchers(ctx context.Context, batchID string, status string, page, pageSize int) ([]Voucher, int, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	offset := (page - 1) * pageSize

	var list []Voucher
	var total int
	err := s.runner.Run(ctx, func(tx kernel.Tx) error {
		where := "WHERE 1=1"
		args := []any{}
		if batchID != "" {
			where += " AND batch_id = ?"
			args = append(args, batchID)
		}
		if status != "" {
			where += " AND status = ?"
			args = append(args, status)
		}

		countQuery := fmt.Sprintf("SELECT COUNT(*) FROM vouchers %s", where)
		if err := tx.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
			return fmt.Errorf("count vouchers: %w", err)
		}

		selectQuery := fmt.Sprintf(`SELECT id, batch_id, code_prefix, amount, currency, status, expires_at, redeemed_by, redeemed_at, created_at, updated_at
			FROM vouchers %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, where)
		queryArgs := append(args, pageSize, offset)
		rows, err := tx.Query(ctx, selectQuery, queryArgs...)
		if err != nil {
			return fmt.Errorf("list vouchers: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var v Voucher
			var exp, redAt sql.NullInt64
			var redBy sql.NullString
			var cr, up int64
			if err := rows.Scan(&v.ID, &v.BatchID, &v.CodePrefix, &v.Amount, &v.Currency, &v.Status, &exp, &redBy, &redAt, &cr, &up); err != nil {
				return err
			}
			if exp.Valid && exp.Int64 > 0 {
				t := time.Unix(exp.Int64, 0).UTC()
				v.ExpiresAt = &t
			}
			if redBy.Valid {
				v.RedeemedBy = &redBy.String
			}
			if redAt.Valid && redAt.Int64 > 0 {
				t := time.Unix(redAt.Int64, 0).UTC()
				v.RedeemedAt = &t
			}
			v.CreatedAt = time.Unix(cr, 0).UTC()
			v.UpdatedAt = time.Unix(up, 0).UTC()
			list = append(list, v)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, 0, err
	}
	return list, total, nil
}
