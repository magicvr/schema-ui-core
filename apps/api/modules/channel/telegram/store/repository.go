// Package store owns the channel.telegram R3 session and inbound receipt
// persistence. It deliberately depends only on the kernel transaction port.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/pagination"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

// TxRunner is the smallest transaction boundary needed by this module.
type TxRunner interface {
	Run(context.Context, func(kernel.Tx) error) error
}

// InboundMessage is the normalized data extracted from a Telegram update.
// Zero-valued optional IDs/strings are stored as SQL NULLs.
type InboundMessage struct {
	BotID           int64
	UpdateID        int64
	ChatID          int64
	ChatType        string
	ChatTitle       string
	ChatUsername    string
	UserID          int64
	MessageID       int64
	CallbackQueryID string
	MessageKind     string
	Text            string
	CallbackData    string
	SenderUsername  string
	ReceivedAt      time.Time
}

// Repository persists Telegram sessions and inbound idempotency receipts.
type Repository struct {
	runner TxRunner
}

var (
	// These sentinels are intentionally narrow store outcomes. The HTTP layer
	// maps them to the frozen Telegram operator error catalog without exposing
	// SQL or downstream Bot API details.
	ErrChatNotFound       = errors.New("telegram store: chat not found")
	ErrRequestNotFound    = errors.New("telegram store: request not found")
	ErrRequestInProgress  = errors.New("telegram store: request is already in progress")
	ErrRequestConflict    = errors.New("telegram store: request payload conflicts")
	ErrRetryNotAllowed    = errors.New("telegram store: retry is not allowed")
	ErrOutboundNotPending = errors.New("telegram store: outbound message is not pending")
)

var muxSafeRequestID = regexp.MustCompile(`^[A-Za-z0-9._-]{1,128}$`)

// Session is the operator-facing projection of one Telegram chat. All IDs are
// kept as int64 internally and serialized as decimal strings by the HTTP layer.
type Session struct {
	BotID         int64
	ChatID        int64
	ChatType      string
	Title         string
	Username      string
	LastMessageAt time.Time
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

// TimelineEntry is the normalized text-only operator transcript projection.
// Optional inbound identifiers are nil when Telegram did not provide them.
type TimelineEntry struct {
	BotID          int64
	ChatID         int64
	Direction      string
	Status         string
	OccurredAt     time.Time
	UpdateID       int64
	MessageID      *int64
	UserID         *int64
	SenderUsername string
	RequestID      string
	RetryOf        string
	Text           string
}

// OutboundMessage is the durable state-machine row for one manual send
// attempt. RetryOf is empty for the first attempt and points to RetryRoot for
// every explicit retry.
type OutboundMessage struct {
	BotID        int64
	RequestID    string
	RetryRoot    string
	RetryOf      string
	ChatID       int64
	Text         string
	Status       string
	ErrorMessage string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

// NewRepository constructs a Telegram R3 repository over a kernel store.
func NewRepository(runner TxRunner) *Repository {
	return &Repository{runner: runner}
}

// RecordInbound atomically records an inbound update and its session activity.
// It returns true only for a new receipt. Duplicate (bot_id, update_id) values
// are successful no-ops and do not update the session or trigger dispatch.
func (r *Repository) RecordInbound(ctx context.Context, msg InboundMessage) (bool, error) {
	if r == nil || r.runner == nil {
		return false, fmt.Errorf("telegram store: repository is unavailable")
	}
	if err := msg.validate(); err != nil {
		return false, err
	}

	inserted := false
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		result, err := tx.Exec(ctx, `INSERT INTO telegram_inbound_messages (
  bot_id, update_id, chat_id, user_id, message_id, callback_query_id,
  direction, message_kind, text, callback_data, sender_username, received_at
) VALUES (?, ?, ?, ?, ?, ?, 'inbound', ?, ?, ?, ?, ?)
ON CONFLICT (bot_id, update_id) DO NOTHING`,
			msg.BotID,
			msg.UpdateID,
			msg.ChatID,
			nullableInt64(msg.UserID),
			nullableInt64(msg.MessageID),
			nullableString(msg.CallbackQueryID),
			msg.MessageKind,
			nullableString(msg.Text),
			nullableString(msg.CallbackData),
			nullableString(msg.SenderUsername),
			msg.ReceivedAt.Unix(),
		)
		if err != nil {
			return fmt.Errorf("insert telegram inbound receipt: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read telegram inbound receipt result: %w", err)
		}
		switch affected {
		case 0:
			// Keep the transaction alive after a PostgreSQL conflict and avoid
			// touching the session row on the duplicate path.
			return nil
		case 1:
			inserted = true
		default:
			return fmt.Errorf("insert telegram inbound receipt affected %d rows", affected)
		}

		if _, err := tx.Exec(ctx, `INSERT INTO telegram_sessions (
  bot_id, chat_id, chat_type, title, username, last_message_at, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (bot_id, chat_id) DO UPDATE SET
  chat_type = CASE WHEN excluded.chat_type <> '' THEN excluded.chat_type ELSE telegram_sessions.chat_type END,
  title = CASE WHEN excluded.title <> '' THEN excluded.title ELSE telegram_sessions.title END,
  username = CASE WHEN excluded.username <> '' THEN excluded.username ELSE telegram_sessions.username END,
  last_message_at = CASE WHEN excluded.last_message_at > telegram_sessions.last_message_at
    THEN excluded.last_message_at ELSE telegram_sessions.last_message_at END,
  updated_at = CASE WHEN excluded.updated_at > telegram_sessions.updated_at
    THEN excluded.updated_at ELSE telegram_sessions.updated_at END`,
			msg.BotID,
			msg.ChatID,
			msg.ChatType,
			msg.ChatTitle,
			msg.ChatUsername,
			msg.ReceivedAt.Unix(),
			msg.ReceivedAt.Unix(),
			msg.ReceivedAt.Unix(),
		); err != nil {
			return fmt.Errorf("upsert telegram session: %w", err)
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return inserted, nil
}

const telegramSessionActivityPredicate = `
  s.bot_id = ?
  AND EXISTS (
    SELECT 1
      FROM telegram_inbound_messages m
     WHERE m.bot_id = s.bot_id
       AND m.chat_id = s.chat_id
       AND m.message_kind IN ('text', 'command')
       AND TRIM(COALESCE(m.text, '')) <> ''
  )`

func (r *Repository) validateAvailable() error {
	if r == nil || r.runner == nil {
		return fmt.Errorf("telegram store: repository is unavailable")
	}
	return nil
}

func validateBotID(botID int64) error {
	if botID <= 0 {
		return fmt.Errorf("telegram store: bot id must be positive")
	}
	return nil
}

func validateChatID(chatID int64) error {
	if chatID == 0 {
		return fmt.Errorf("telegram store: chat id is required")
	}
	return nil
}

func validatePage(page, pageSize int) error {
	if page < 1 {
		return fmt.Errorf("telegram store: page must be positive")
	}
	if pageSize < 1 || pageSize > 100 {
		return fmt.Errorf("telegram store: page size must be between 1 and 100")
	}
	return nil
}

// ListSessions returns only chats with persisted inbound text/command
// activity. Sorting is fixed server-side and offset calculation is overflow
// safe even for a very large requested page.
func (r *Repository) ListSessions(ctx context.Context, botID int64, page, pageSize int) ([]Session, int, error) {
	if err := r.validateAvailable(); err != nil {
		return nil, 0, err
	}
	if err := validateBotID(botID); err != nil {
		return nil, 0, err
	}
	if err := validatePage(page, pageSize); err != nil {
		return nil, 0, err
	}

	var sessions []Session
	var total int
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		countQuery := `SELECT COUNT(*) FROM telegram_sessions s WHERE ` + telegramSessionActivityPredicate
		if err := tx.QueryRow(ctx, countQuery, botID).Scan(&total); err != nil {
			return fmt.Errorf("count telegram sessions: %w", err)
		}

		offset := pagination.Offset(page, pageSize, total)
		rows, err := tx.Query(ctx, `SELECT s.bot_id, s.chat_id, s.chat_type, s.title, s.username,
  s.last_message_at, s.created_at, s.updated_at
FROM telegram_sessions s
WHERE `+telegramSessionActivityPredicate+`
ORDER BY s.last_message_at DESC, s.chat_id DESC
LIMIT ? OFFSET ?`, botID, pageSize, offset)
		if err != nil {
			return fmt.Errorf("list telegram sessions: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var session Session
			var lastMessageAt, createdAt, updatedAt int64
			if err := rows.Scan(&session.BotID, &session.ChatID, &session.ChatType, &session.Title, &session.Username, &lastMessageAt, &createdAt, &updatedAt); err != nil {
				return fmt.Errorf("scan telegram session: %w", err)
			}
			session.LastMessageAt = time.Unix(lastMessageAt, 0).UTC()
			session.CreatedAt = time.Unix(createdAt, 0).UTC()
			session.UpdatedAt = time.Unix(updatedAt, 0).UTC()
			sessions = append(sessions, session)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, 0, err
	}
	if sessions == nil {
		sessions = []Session{}
	}
	return sessions, total, nil
}

// SessionExists reports whether chat_id belongs to the current bot and has
// visible inbound text/command activity. Callback-only sessions intentionally
// do not satisfy this predicate.
func (r *Repository) SessionExists(ctx context.Context, botID, chatID int64) (bool, error) {
	if err := r.validateAvailable(); err != nil {
		return false, err
	}
	if err := validateBotID(botID); err != nil {
		return false, err
	}
	if err := validateChatID(chatID); err != nil {
		return false, err
	}
	var marker int
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		return tx.QueryRow(ctx, `SELECT 1 FROM telegram_sessions s WHERE s.bot_id = ? AND s.chat_id = ? AND EXISTS (
  SELECT 1 FROM telegram_inbound_messages m
   WHERE m.bot_id = s.bot_id AND m.chat_id = s.chat_id
     AND m.message_kind IN ('text', 'command')
     AND TRIM(COALESCE(m.text, '')) <> ''
) LIMIT 1`, botID, chatID).Scan(&marker)
	})
	if errors.Is(err, kernel.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check telegram session: %w", err)
	}
	return marker == 1, nil
}

const telegramTimelineSourceSQL = `
SELECT bot_id, chat_id, direction, status, occurred_at, sort_id,
       update_id, message_id, user_id, sender_username, request_id, retry_of, text
  FROM (
    SELECT bot_id, chat_id, 'inbound' AS direction, 'received' AS status,
           received_at AS occurred_at, CAST(update_id AS TEXT) AS sort_id,
           update_id, message_id, user_id, sender_username,
           NULL AS request_id, NULL AS retry_of, text
      FROM telegram_inbound_messages
     WHERE bot_id = ?
       AND chat_id = ?
       AND message_kind IN ('text', 'command')
       AND TRIM(COALESCE(text, '')) <> ''
    UNION ALL
    SELECT bot_id, chat_id, 'outbound' AS direction, status,
           created_at AS occurred_at, request_id AS sort_id,
           NULL AS update_id, NULL AS message_id, NULL AS user_id,
           NULL AS sender_username, request_id, retry_of, text
      FROM telegram_outbound_messages
     WHERE bot_id = ?
       AND chat_id = ?
  ) AS timeline`

// ListTimeline returns the unified text transcript for one current-bot chat.
// Callback rows, empty text and unmodeled media are excluded by construction.
func (r *Repository) ListTimeline(ctx context.Context, botID, chatID int64, page, pageSize int) ([]TimelineEntry, int, error) {
	if err := r.validateAvailable(); err != nil {
		return nil, 0, err
	}
	if err := validateBotID(botID); err != nil {
		return nil, 0, err
	}
	if err := validateChatID(chatID); err != nil {
		return nil, 0, err
	}
	if err := validatePage(page, pageSize); err != nil {
		return nil, 0, err
	}

	var entries []TimelineEntry
	var total int
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		countQuery := `SELECT COUNT(*) FROM (` + telegramTimelineSourceSQL + `) AS counted`
		args := []any{botID, chatID, botID, chatID}
		if err := tx.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
			return fmt.Errorf("count telegram timeline: %w", err)
		}
		offset := pagination.Offset(page, pageSize, total)
		selectQuery := telegramTimelineSourceSQL + `
ORDER BY occurred_at DESC, sort_id DESC, direction DESC
LIMIT ? OFFSET ?`
		args = append(args, pageSize, offset)
		rows, err := tx.Query(ctx, selectQuery, args...)
		if err != nil {
			return fmt.Errorf("list telegram timeline: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var entry TimelineEntry
			var occurredAt int64
			var updateID, messageID, userID sql.NullInt64
			var senderUsername, requestID, retryOf sql.NullString
			if err := rows.Scan(&entry.BotID, &entry.ChatID, &entry.Direction, &entry.Status, &occurredAt, new(string), &updateID, &messageID, &userID, &senderUsername, &requestID, &retryOf, &entry.Text); err != nil {
				return fmt.Errorf("scan telegram timeline: %w", err)
			}
			entry.OccurredAt = time.Unix(occurredAt, 0).UTC()
			if updateID.Valid {
				entry.UpdateID = updateID.Int64
			}
			if messageID.Valid {
				value := messageID.Int64
				entry.MessageID = &value
			}
			if userID.Valid {
				value := userID.Int64
				entry.UserID = &value
			}
			if senderUsername.Valid {
				entry.SenderUsername = senderUsername.String
			}
			if requestID.Valid {
				entry.RequestID = requestID.String
			}
			if retryOf.Valid {
				entry.RetryOf = retryOf.String
			}
			entries = append(entries, entry)
		}
		return rows.Err()
	})
	if err != nil {
		return nil, 0, err
	}
	if entries == nil {
		entries = []TimelineEntry{}
	}
	return entries, total, nil
}

// GetOutbound reads one outbound attempt scoped to the current bot and chat.
// Scoping the lookup prevents request existence from leaking across chats.
func (r *Repository) GetOutbound(ctx context.Context, botID, chatID int64, requestID string) (OutboundMessage, error) {
	if err := r.validateAvailable(); err != nil {
		return OutboundMessage{}, err
	}
	if err := validateBotID(botID); err != nil {
		return OutboundMessage{}, err
	}
	if err := validateChatID(chatID); err != nil {
		return OutboundMessage{}, err
	}
	if err := validateRequestID(requestID); err != nil {
		return OutboundMessage{}, err
	}
	var message OutboundMessage
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		var err error
		message, err = readOutboundTx(ctx, tx, botID, chatID, requestID)
		return err
	})
	if errors.Is(err, ErrRequestNotFound) {
		return OutboundMessage{}, ErrRequestNotFound
	}
	return message, err
}

// CreatePending creates a first-attempt pending row, or returns a previously
// terminal row when the same request payload is replayed. Every conflict path
// uses ON CONFLICT DO NOTHING and then reads in the still-live transaction;
// this is safe for both SQLite and PostgreSQL, including the partial unique
// retry-root index.
func (r *Repository) CreatePending(ctx context.Context, botID, chatID int64, requestID, text string) (OutboundMessage, bool, error) {
	if err := r.validateAvailable(); err != nil {
		return OutboundMessage{}, false, err
	}
	if err := validateOutboundInput(botID, chatID, requestID, text); err != nil {
		return OutboundMessage{}, false, err
	}
	var message OutboundMessage
	created := false
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		exists, err := sessionExistsTx(ctx, tx, botID, chatID)
		if err != nil {
			return fmt.Errorf("check telegram session before send: %w", err)
		}
		if !exists {
			return ErrChatNotFound
		}
		message, created, err = createPendingTx(ctx, tx, botID, chatID, requestID, requestID, "", text)
		return err
	})
	return message, created, err
}

// CreateRetry creates a new pending attempt for a failed source row. The new
// request id is distinct, while retry_of always points to the original root.
func (r *Repository) CreateRetry(ctx context.Context, botID, chatID int64, sourceRequestID, requestID string) (OutboundMessage, bool, error) {
	if err := r.validateAvailable(); err != nil {
		return OutboundMessage{}, false, err
	}
	if err := validateBotID(botID); err != nil {
		return OutboundMessage{}, false, err
	}
	if err := validateChatID(chatID); err != nil {
		return OutboundMessage{}, false, err
	}
	if err := validateRequestID(sourceRequestID); err != nil {
		return OutboundMessage{}, false, err
	}
	if err := validateRequestID(requestID); err != nil {
		return OutboundMessage{}, false, err
	}
	var message OutboundMessage
	created := false
	err := r.runner.Run(ctx, func(tx kernel.Tx) error {
		source, err := readOutboundTx(ctx, tx, botID, chatID, sourceRequestID)
		if errors.Is(err, ErrRequestNotFound) {
			return ErrRequestNotFound
		}
		if err != nil {
			return err
		}
		if source.Status != "failed" {
			return ErrRetryNotAllowed
		}
		sent, err := hasOutboundStatusTx(ctx, tx, botID, source.RetryRoot, "sent")
		if err != nil {
			return err
		}
		if sent {
			return ErrRetryNotAllowed
		}
		pending, err := hasOutboundStatusTx(ctx, tx, botID, source.RetryRoot, "pending")
		if err != nil {
			return err
		}
		if pending {
			return ErrRequestInProgress
		}
		message, created, err = createPendingTx(ctx, tx, botID, chatID, requestID, source.RetryRoot, source.RetryRoot, source.Text)
		return err
	})
	return message, created, err
}

// MarkSent transitions one pending row to sent. A terminal row is never
// silently treated as success: the caller must not risk a second external send
// after a state transition race or database failure.
func (r *Repository) MarkSent(ctx context.Context, botID int64, requestID string) error {
	return r.markStatus(ctx, botID, requestID, "sent", "")
}

// MarkFailed transitions one pending row to failed and stores only one of the
// fixed, non-sensitive diagnostic categories accepted by safeErrorMessage.
func (r *Repository) MarkFailed(ctx context.Context, botID int64, requestID, reason string) error {
	return r.markStatus(ctx, botID, requestID, "failed", safeErrorMessage(reason))
}

func (r *Repository) markStatus(ctx context.Context, botID int64, requestID, status, reason string) error {
	if err := r.validateAvailable(); err != nil {
		return err
	}
	if err := validateBotID(botID); err != nil {
		return err
	}
	if err := validateRequestID(requestID); err != nil {
		return err
	}
	return r.runner.Run(ctx, func(tx kernel.Tx) error {
		var result kernel.Result
		var err error
		if status == "sent" {
			result, err = tx.Exec(ctx, `UPDATE telegram_outbound_messages
SET status = 'sent', error_message = NULL, updated_at = ?
WHERE bot_id = ? AND request_id = ? AND status = 'pending'`, time.Now().UTC().Unix(), botID, requestID)
		} else {
			result, err = tx.Exec(ctx, `UPDATE telegram_outbound_messages
SET status = 'failed', error_message = ?, updated_at = ?
WHERE bot_id = ? AND request_id = ? AND status = 'pending'`, reason, time.Now().UTC().Unix(), botID, requestID)
		}
		if err != nil {
			return fmt.Errorf("update telegram outbound status: %w", err)
		}
		affected, err := result.RowsAffected()
		if err != nil {
			return fmt.Errorf("read telegram outbound status result: %w", err)
		}
		if affected == 1 {
			return nil
		}
		current, err := readOutboundByRequestTx(ctx, tx, botID, requestID)
		if errors.Is(err, ErrRequestNotFound) {
			return ErrRequestNotFound
		}
		if err != nil {
			return err
		}
		if current.Status != "pending" {
			return ErrOutboundNotPending
		}
		return fmt.Errorf("telegram store: status update affected %d rows", affected)
	})
}

func validateRequestID(requestID string) error {
	if !muxSafeRequestID.MatchString(requestID) {
		return fmt.Errorf("telegram store: request id must match [A-Za-z0-9._-]{1,128}")
	}
	return nil
}

func validateOutboundInput(botID, chatID int64, requestID, text string) error {
	if err := validateBotID(botID); err != nil {
		return err
	}
	if err := validateChatID(chatID); err != nil {
		return err
	}
	if err := validateRequestID(requestID); err != nil {
		return err
	}
	if strings.TrimSpace(text) == "" {
		return fmt.Errorf("telegram store: text is required")
	}
	if len([]byte(text)) > 4096 {
		return fmt.Errorf("telegram store: text must not exceed 4096 UTF-8 bytes")
	}
	return nil
}

func createPendingTx(ctx context.Context, tx kernel.Tx, botID, chatID int64, requestID, retryRoot, retryOf, text string) (OutboundMessage, bool, error) {
	now := time.Now().UTC()
	result, err := tx.Exec(ctx, `INSERT INTO telegram_outbound_messages (
  bot_id, request_id, retry_root, retry_of, chat_id, text, status, error_message, created_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, 'pending', NULL, ?, ?)
ON CONFLICT DO NOTHING`, botID, requestID, retryRoot, nullableString(retryOf), chatID, text, now.Unix(), now.Unix())
	if err != nil {
		return OutboundMessage{}, false, fmt.Errorf("insert telegram outbound pending: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return OutboundMessage{}, false, fmt.Errorf("read telegram outbound insert result: %w", err)
	}
	switch affected {
	case 1:
		return OutboundMessage{
			BotID: botID, RequestID: requestID, RetryRoot: retryRoot, RetryOf: retryOf,
			ChatID: chatID, Text: text, Status: "pending", CreatedAt: now, UpdatedAt: now,
		}, true, nil
	case 0:
		existing, readErr := readOutboundTx(ctx, tx, botID, chatID, requestID)
		if readErr == nil {
			if !sameOutboundPayload(existing, chatID, text, retryOf) {
				return OutboundMessage{}, false, ErrRequestConflict
			}
			if existing.Status == "pending" {
				return OutboundMessage{}, false, ErrRequestInProgress
			}
			return existing, false, nil
		}
		if !errors.Is(readErr, ErrRequestNotFound) {
			return OutboundMessage{}, false, readErr
		}
		// The primary key is bot-wide, while the first read is deliberately
		// chat-scoped to avoid leaking another chat's row. A request id already
		// owned by another chat is nevertheless a payload conflict, including
		// when its existing row is terminal; never turn that case into a 500.
		if _, requestErr := readOutboundByRequestTx(ctx, tx, botID, requestID); requestErr == nil {
			return OutboundMessage{}, false, ErrRequestConflict
		} else if !errors.Is(requestErr, ErrRequestNotFound) {
			return OutboundMessage{}, false, requestErr
		}
		pending, statusErr := hasOutboundStatusTx(ctx, tx, botID, retryRoot, "pending")
		if statusErr != nil {
			return OutboundMessage{}, false, statusErr
		}
		if pending {
			return OutboundMessage{}, false, ErrRequestInProgress
		}
		return OutboundMessage{}, false, fmt.Errorf("telegram store: outbound insert conflict could not be resolved")
	default:
		return OutboundMessage{}, false, fmt.Errorf("telegram store: outbound insert affected %d rows", affected)
	}
}

func sessionExistsTx(ctx context.Context, tx kernel.Tx, botID, chatID int64) (bool, error) {
	var marker int
	err := tx.QueryRow(ctx, `SELECT 1 FROM telegram_sessions s WHERE s.bot_id = ? AND s.chat_id = ? AND EXISTS (
  SELECT 1 FROM telegram_inbound_messages m
   WHERE m.bot_id = s.bot_id AND m.chat_id = s.chat_id
     AND m.message_kind IN ('text', 'command')
     AND TRIM(COALESCE(m.text, '')) <> ''
) LIMIT 1`, botID, chatID).Scan(&marker)
	if errors.Is(err, kernel.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return marker == 1, nil
}

func hasOutboundStatusTx(ctx context.Context, tx kernel.Tx, botID int64, retryRoot, status string) (bool, error) {
	var marker int
	err := tx.QueryRow(ctx, `SELECT 1 FROM telegram_outbound_messages WHERE bot_id = ? AND retry_root = ? AND status = ? LIMIT 1`, botID, retryRoot, status).Scan(&marker)
	if errors.Is(err, kernel.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("check telegram outbound status: %w", err)
	}
	return marker == 1, nil
}

func readOutboundByRequestTx(ctx context.Context, tx kernel.Tx, botID int64, requestID string) (OutboundMessage, error) {
	return scanOutbound(tx.QueryRow(ctx, `SELECT bot_id, request_id, retry_root, retry_of, chat_id, text, status, error_message, created_at, updated_at
FROM telegram_outbound_messages WHERE bot_id = ? AND request_id = ?`, botID, requestID))
}

func readOutboundTx(ctx context.Context, tx kernel.Tx, botID, chatID int64, requestID string) (OutboundMessage, error) {
	message, err := scanOutbound(tx.QueryRow(ctx, `SELECT bot_id, request_id, retry_root, retry_of, chat_id, text, status, error_message, created_at, updated_at
FROM telegram_outbound_messages WHERE bot_id = ? AND chat_id = ? AND request_id = ?`, botID, chatID, requestID))
	if errors.Is(err, kernel.ErrNoRows) {
		return OutboundMessage{}, ErrRequestNotFound
	}
	return message, err
}

func scanOutbound(row kernel.Row) (OutboundMessage, error) {
	var message OutboundMessage
	var retryOf, errorMessage sql.NullString
	var createdAt, updatedAt int64
	if err := row.Scan(&message.BotID, &message.RequestID, &message.RetryRoot, &retryOf, &message.ChatID, &message.Text, &message.Status, &errorMessage, &createdAt, &updatedAt); err != nil {
		if errors.Is(err, kernel.ErrNoRows) {
			return OutboundMessage{}, ErrRequestNotFound
		}
		return OutboundMessage{}, fmt.Errorf("scan telegram outbound message: %w", err)
	}
	if retryOf.Valid {
		message.RetryOf = retryOf.String
	}
	if errorMessage.Valid {
		message.ErrorMessage = errorMessage.String
	}
	message.CreatedAt = time.Unix(createdAt, 0).UTC()
	message.UpdatedAt = time.Unix(updatedAt, 0).UTC()
	return message, nil
}

func sameOutboundPayload(message OutboundMessage, chatID int64, text, retryOf string) bool {
	return message.ChatID == chatID && message.Text == text && message.RetryOf == retryOf
}

func safeErrorMessage(reason string) string {
	clean := strings.Join(strings.Fields(reason), " ")
	switch clean {
	case "telegram send failed", "telegram operator became unavailable":
		return clean
	default:
		// Downstream Bot API errors can contain request metadata or secrets.
		// The handler supplies only these fixed categories, and the repository
		// keeps the same allow-list as a defense-in-depth boundary for callers.
		return "telegram send failed"
	}
}

func (m InboundMessage) validate() error {
	if m.BotID <= 0 {
		return fmt.Errorf("telegram store: bot id must be positive")
	}
	if m.UpdateID <= 0 {
		return fmt.Errorf("telegram store: update id must be positive")
	}
	if m.ChatID == 0 {
		return fmt.Errorf("telegram store: chat id is required")
	}
	if m.ReceivedAt.IsZero() {
		return fmt.Errorf("telegram store: received time is required")
	}
	switch m.MessageKind {
	case "text", "command", "callback":
	default:
		return fmt.Errorf("telegram store: unsupported message kind %q", m.MessageKind)
	}
	if (m.MessageKind == "text" || m.MessageKind == "command") && strings.TrimSpace(m.Text) == "" {
		return fmt.Errorf("telegram store: text is required for %s", m.MessageKind)
	}
	return nil
}

func nullableInt64(value int64) any {
	if value == 0 {
		return nil
	}
	return value
}

func nullableString(value string) any {
	if value == "" {
		return nil
	}
	return value
}
