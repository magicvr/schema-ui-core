package mail

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

// ErrOutboxRecordNotFound is returned by OutboxSink.Get for unknown ids.
var ErrOutboxRecordNotFound = errors.New("mail: outbox record not found")

// OutboxSink is the mock-channel kernel.MailSender adapter (VP-017 R6 /
// workspace-017 GOAL-007; contract frozen by workspace-017 GOAL-006 D-002 §3).
// It publishes every accepted message to the admin-visible in-app outbound
// record table (mail_outbox, migration 0051) — NOT a user-facing notification
// transport. Records survive restarts and are shared across instances.
//
// Retention is bounded: the most recent DefaultOutboxCap records are kept
// (inserts beyond the cap evict the oldest). The cap is a mock-channel
// configuration value; the settings-surface override lands with R7 — until
// then composition passes DefaultOutboxCap.
//
// Like CaptureSink it never fails for configuration reasons and never blocks
// on network I/O; errors mirror the port rules plus storage failures only.
type OutboxSink struct {
	store kernel.Store
	cap   int
	now   func() time.Time // test seam
}

// DefaultOutboxCap is the frozen default bounded-retention cap
// (GOAL-006 D-002 §3.3): keep the most recent 500 records.
const DefaultOutboxCap = 500

// OutboxRecord is one mock outbound message as surfaced through the admin
// retrieval API. List views omit Body; Get fills it.
type OutboxRecord struct {
	ID        string    `json:"id"`
	To        string    `json:"to"`
	Subject   string    `json:"subject"`
	Body      string    `json:"body,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

// outboxIDSeq guarantees process-local uniqueness regardless of entropy
// quality (same contract as the scheduled-tasks run id): monotonic sequence +
// UnixNano; rand bytes only perturb the string.
var outboxIDSeq atomic.Uint64

func newOutboxID(now time.Time) string {
	seq := outboxIDSeq.Add(1)
	var b [4]byte
	if _, err := rand.Read(b[:]); err == nil {
		return fmt.Sprintf("outbox-%x-%x-%x", now.UnixNano(), seq, b)
	}
	return fmt.Sprintf("outbox-%x-%x", now.UnixNano(), seq)
}

// NewOutboxSink returns the mock-channel publisher over the platform store.
// retentionCap <= 0 falls back to DefaultOutboxCap.
func NewOutboxSink(store kernel.Store, retentionCap int) *OutboxSink {
	if retentionCap <= 0 {
		retentionCap = DefaultOutboxCap
	}
	return &OutboxSink{store: store, cap: retentionCap, now: time.Now}
}

// Send validates the message against the frozen port contract, then appends
// one record and enforces the bounded retention inside the same transaction:
// either both succeed or neither does. There is no delivery to fail.
func (s *OutboxSink) Send(ctx context.Context, msg kernel.MailMessage) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("mail: %v", err)
	}
	now := s.now().UTC()
	id := newOutboxID(now)
	err := s.store.Run(ctx, func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO mail_outbox (id, to_addr, subject, body, created_at)
			 VALUES (?, ?, ?, ?, ?)`,
			id, msg.To, msg.Subject, msg.TextBody, now.UnixMilli(),
		); err != nil {
			return err
		}
		_, err := tx.Exec(context.Background(),
			`DELETE FROM mail_outbox WHERE id NOT IN (
				 SELECT id FROM mail_outbox ORDER BY created_at DESC, id DESC LIMIT ?
			 )`,
			s.cap,
		)
		return err
	})
	if err != nil {
		return fmt.Errorf("mail: publish outbox record: %w", err)
	}
	return nil
}

// List returns at most limit records, newest first, without bodies. limit<=0
// falls back to a page of 50; offset<0 is treated as 0.
func (s *OutboxSink) List(ctx context.Context, limit, offset int) ([]OutboxRecord, error) {
	if limit <= 0 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	var records []OutboxRecord
	err := s.store.Run(ctx, func(tx kernel.Tx) error {
		rows, err := tx.Query(context.Background(),
			`SELECT id, to_addr, subject, created_at FROM mail_outbox
			 ORDER BY created_at DESC, id DESC LIMIT ? OFFSET ?`,
			limit, offset,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var rec OutboxRecord
			var ms int64
			if err := rows.Scan(&rec.ID, &rec.To, &rec.Subject, &ms); err != nil {
				return fmt.Errorf("scan outbox row: %w", err)
			}
			rec.CreatedAt = time.UnixMilli(ms).UTC()
			records = append(records, rec)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, fmt.Errorf("mail: list outbox: %w", err)
	}
	return records, nil
}

// Count returns the current number of stored records (pagination total).
func (s *OutboxSink) Count(ctx context.Context) (int64, error) {
	var n int64
	err := s.store.Run(ctx, func(tx kernel.Tx) error {
		return tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM mail_outbox`).Scan(&n)
	})
	if err != nil {
		return 0, fmt.Errorf("mail: count outbox: %w", err)
	}
	return n, nil
}

// Get returns one full record (including the body) by id.
func (s *OutboxSink) Get(ctx context.Context, id string) (OutboxRecord, error) {
	var rec OutboxRecord
	var ms int64
	err := s.store.Run(ctx, func(tx kernel.Tx) error {
		switch err := tx.QueryRow(context.Background(),
			`SELECT id, to_addr, subject, body, created_at FROM mail_outbox WHERE id = ?`, id,
		).Scan(&rec.ID, &rec.To, &rec.Subject, &rec.Body, &ms); {
		case errors.Is(err, kernel.ErrNoRows):
			return ErrOutboxRecordNotFound
		default:
			return err
		}
	})
	if err != nil {
		if errors.Is(err, ErrOutboxRecordNotFound) {
			return OutboxRecord{}, err
		}
		return OutboxRecord{}, fmt.Errorf("mail: get outbox record: %w", err)
	}
	rec.CreatedAt = time.UnixMilli(ms).UTC()
	return rec, nil
}
