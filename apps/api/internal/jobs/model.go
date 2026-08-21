// Package jobs provides the profile-independent durable async Job contract.
package jobs

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

type Status string

const (
	StatusQueued    Status = "queued"
	StatusRunning   Status = "running"
	StatusSucceeded Status = "succeeded"
	StatusFailed    Status = "failed"
	StatusCancelled Status = "cancelled"
	StatusExpired   Status = "expired"
)

const DefaultMaxAttempts = 3

var (
	ErrNotFound       = errors.New("job not found")
	ErrInvalid        = errors.New("invalid job input")
	ErrTransition     = errors.New("invalid job transition")
	ErrLeaseLost      = errors.New("job lease lost")
	ErrNotCancellable = errors.New("job is not cancellable")
	ErrNotRetryable   = errors.New("job is not retryable")
)

type Job struct {
	ID              string
	Kind            string
	Status          Status
	Payload         json.RawMessage
	Progress        int
	CancelRequested bool
	Attempt         int
	MaxAttempts     int
	LeaseOwner      string
	LeaseVersion    int64
	LeaseExpiresAt  *time.Time
	Result          json.RawMessage
	ErrorCode       string
	ErrorMessage    string
	ActorID         string
	CorrelationID   string
	CreatedAt       time.Time
	UpdatedAt       time.Time
	FinishedAt      *time.Time
	ResultExpiresAt *time.Time
}

type CreateInput struct {
	ID            string
	Kind          string
	Payload       json.RawMessage
	ActorID       string
	CorrelationID string
	MaxAttempts   int
	Now           time.Time
}

type Lease struct {
	JobID   string
	Owner   string
	Version int64
}

// NewID returns a sortable Job identifier without coupling callers to a module.
func NewID(now time.Time) (string, error) {
	random := make([]byte, 12)
	if _, err := rand.Read(random); err != nil {
		return "", err
	}
	return fmt.Sprintf("%016x%s", now.UTC().UnixMilli(), hex.EncodeToString(random)), nil
}

type rowScanner interface {
	Scan(...any) error
}

const jobColumns = `id, kind, status, payload, progress, cancel_requested,
attempt, max_attempts, lease_owner, lease_version, lease_expires_at,
result, error_code, error_message, actor_id, correlation_id,
created_at, updated_at, finished_at, expires_at`

func scanJob(row rowScanner) (*Job, error) {
	var job Job
	var status string
	var payload string
	var cancelRequested int
	var leaseOwner, result, errorCode, errorMessage sql.NullString
	var leaseExpiresAt, finishedAt, expiresAt sql.NullInt64
	var createdAt, updatedAt int64
	if err := row.Scan(
		&job.ID, &job.Kind, &status, &payload, &job.Progress, &cancelRequested,
		&job.Attempt, &job.MaxAttempts, &leaseOwner, &job.LeaseVersion, &leaseExpiresAt,
		&result, &errorCode, &errorMessage, &job.ActorID, &job.CorrelationID,
		&createdAt, &updatedAt, &finishedAt, &expiresAt,
	); err != nil {
		if errors.Is(err, kernel.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	job.Status = Status(status)
	job.Payload = json.RawMessage(payload)
	job.CancelRequested = cancelRequested != 0
	job.LeaseOwner = leaseOwner.String
	job.Result = nullableJSON(result)
	job.ErrorCode = errorCode.String
	job.ErrorMessage = errorMessage.String
	job.CreatedAt = fromMillis(createdAt)
	job.UpdatedAt = fromMillis(updatedAt)
	job.LeaseExpiresAt = nullableTime(leaseExpiresAt)
	job.FinishedAt = nullableTime(finishedAt)
	job.ResultExpiresAt = nullableTime(expiresAt)
	return &job, nil
}

func nullableJSON(value sql.NullString) json.RawMessage {
	if !value.Valid {
		return nil
	}
	return json.RawMessage(value.String)
}

func nullableTime(value sql.NullInt64) *time.Time {
	if !value.Valid {
		return nil
	}
	t := fromMillis(value.Int64)
	return &t
}

func toMillis(value time.Time) int64 { return value.UTC().UnixMilli() }

func fromMillis(value int64) time.Time { return time.UnixMilli(value).UTC() }
