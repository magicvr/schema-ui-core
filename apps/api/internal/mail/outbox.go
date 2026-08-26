package mail

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"strings"
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

// Frozen delivery-status vocabulary (W26 · GOAL-038 D-001 §2.1). delivered =
// accepted by the in-app mock transport; sent = the real channel accepted the
// message (2xx / 250); failed = the channel adapter returned an error.
const (
	ChannelMock   = "mock"
	ChannelResend = "resend"
	ChannelSMTP   = "smtp"

	DeliveryDelivered = "delivered"
	DeliverySent      = "sent"
	DeliveryFailed    = "failed"
)

// OutboxRecord is one outbound message as surfaced through the admin
// retrieval API. Since W26 every channel records here; list views carry the
// body too (D-001 §2.1 contract revision — the declarative recordView drawer
// renders detail straight from the selected row).
type OutboxRecord struct {
	ID             string    `json:"id"`
	To             string    `json:"to"`
	Subject        string    `json:"subject"`
	Body           string    `json:"body,omitempty"`
	Channel        string    `json:"channel"`
	DeliveryStatus string    `json:"delivery_status"`
	CreatedAt      time.Time `json:"created_at"`
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

// publishOutboundRecord appends one outbound record with the frozen delivery
// status and enforces bounded retention inside the same transaction: either
// both succeed or neither does. Shared by the mock transport (Send below) and
// the runtime switcher's real-channel recording (runtime.go).
func publishOutboundRecord(store kernel.Store, retentionCap int, now time.Time, msg kernel.MailMessage, channel, status string) error {
	id := newOutboxID(now)
	err := store.Run(context.Background(), func(tx kernel.Tx) error {
		if _, err := tx.Exec(context.Background(),
			`INSERT INTO mail_outbox (id, to_addr, subject, body, channel, delivery_status, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			id, msg.To, msg.Subject, msg.TextBody, channel, status, now.UnixMilli(),
		); err != nil {
			return err
		}
		_, err := tx.Exec(context.Background(),
			`DELETE FROM mail_outbox WHERE id NOT IN (
				 SELECT id FROM mail_outbox ORDER BY created_at DESC, id DESC LIMIT ?
			 )`,
			retentionCap,
		)
		return err
	})
	if err != nil {
		return fmt.Errorf("mail: publish outbox record: %w", err)
	}
	return nil
}

// Send validates the message against the frozen port contract, then appends
// one record (mock → delivered) and enforces the bounded retention inside the
// same transaction. There is no delivery to fail.
func (s *OutboxSink) Send(ctx context.Context, msg kernel.MailMessage) error {
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("mail: %v", err)
	}
	now := s.now().UTC()
	if err := publishOutboundRecord(s.store, s.cap, now, msg, ChannelMock, DeliveryDelivered); err != nil {
		return err
	}
	return nil
}

// OutboxListQuery is the admin listing contract (W27 · GOAL-039 D-001 §2):
// standard page/pageSize pagination plus keyword / channel / delivery-status
// filters and a whitelist sort. Empty filter values mean "all"; unknown enum
// or sort values fall back to the defaults ("all" / created_at DESC,
// ParseInviteStatus fallback philosophy). The stable secondary key keeps
// pagination deterministic across dialects.
type OutboxListQuery struct {
	Page           int    // 1-based; <=0 → 1
	PageSize       int    // <=0 → DefaultOutboxPageSize; > MaxOutboxPageSize → capped
	Q              string // matches lower(to_addr) / lower(subject)
	Channel        string // mock | resend | smtp ("": all)
	DeliveryStatus string // delivered | sent | failed ("": all)
	Sort           string // created_at (only whitelist entry); unknown → created_at
	Order          string // asc | desc (default desc)
}

const (
	DefaultOutboxPageSize = 50
	MaxOutboxPageSize     = 200
)

func (q OutboxListQuery) normalized() OutboxListQuery {
	if q.Page < 1 {
		q.Page = 1
	}
	if q.PageSize < 1 {
		q.PageSize = DefaultOutboxPageSize
	}
	if q.PageSize > MaxOutboxPageSize {
		q.PageSize = MaxOutboxPageSize
	}
	switch q.Channel {
	case ChannelMock, ChannelResend, ChannelSMTP:
	default:
		q.Channel = ""
	}
	switch q.DeliveryStatus {
	case DeliveryDelivered, DeliverySent, DeliveryFailed:
	default:
		q.DeliveryStatus = ""
	}
	return q
}

// whereClause builds the portable filter predicate + args (usersWhere
// LOWER+LIKE precedent); empty clause = no filtering.
func (q OutboxListQuery) whereClause() (string, []any) {
	clauses := make([]string, 0, 3)
	args := make([]any, 0, 3)
	if needle := strings.ToLower(strings.TrimSpace(q.Q)); needle != "" {
		clauses = append(clauses, `(lower(to_addr) LIKE '%' || CAST(? AS TEXT) || '%' OR lower(subject) LIKE '%' || CAST(? AS TEXT) || '%')`)
		args = append(args, needle, needle)
	}
	if q.Channel != "" {
		clauses = append(clauses, `channel = ?`)
		args = append(args, q.Channel)
	}
	if q.DeliveryStatus != "" {
		clauses = append(clauses, `delivery_status = ?`)
		args = append(args, q.DeliveryStatus)
	}
	if len(clauses) == 0 {
		return "", nil
	}
	return ` WHERE ` + strings.Join(clauses, ` AND `), args
}

func (q OutboxListQuery) orderClause() string {
	direction := "DESC"
	if strings.EqualFold(strings.TrimSpace(q.Order), "asc") {
		direction = "ASC"
	}
	// created_at is the only sortable column (log-style list); the same-
	// direction id tiebreak keeps pagination stable across dialects.
	return "created_at " + direction + ", id " + direction
}

// List returns one page of records matching the query, newest-first by
// default, together with the filtered total for pagination.
func (s *OutboxSink) List(ctx context.Context, query OutboxListQuery) ([]OutboxRecord, int, error) {
	query = query.normalized()
	where, args := query.whereClause()
	var records []OutboxRecord
	var total int
	err := s.store.Run(ctx, func(tx kernel.Tx) error {
		if err := tx.QueryRow(context.Background(),
			`SELECT COUNT(*) FROM mail_outbox`+where, args...).Scan(&total); err != nil {
			return err
		}
		rows, err := tx.Query(context.Background(),
			`SELECT id, to_addr, subject, body, channel, delivery_status, created_at FROM mail_outbox`+where+
				` ORDER BY `+query.orderClause()+` LIMIT ? OFFSET ?`,
			append(append([]any{}, args...), query.PageSize, (query.Page-1)*query.PageSize)...,
		)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var rec OutboxRecord
			var ms int64
			if err := rows.Scan(&rec.ID, &rec.To, &rec.Subject, &rec.Body, &rec.Channel, &rec.DeliveryStatus, &ms); err != nil {
				return fmt.Errorf("scan outbox row: %w", err)
			}
			rec.CreatedAt = time.UnixMilli(ms).UTC()
			records = append(records, rec)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, 0, fmt.Errorf("mail: list outbox: %w", err)
	}
	return records, total, nil
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
			`SELECT id, to_addr, subject, body, channel, delivery_status, created_at FROM mail_outbox WHERE id = ?`, id,
		).Scan(&rec.ID, &rec.To, &rec.Subject, &rec.Body, &rec.Channel, &rec.DeliveryStatus, &ms); {
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
