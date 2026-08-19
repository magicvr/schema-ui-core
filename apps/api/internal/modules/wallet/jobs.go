package wallet

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet/store"
)

const ReconcileJobKind = "wallet.reconcile"

type reconcileJobPayload struct {
	AccountID string `json:"accountId,omitempty"`
}

// JobService binds the profile-independent Job runtime to wallet reconcile.
type JobService struct {
	service    *Service
	repository *jobs.Repository
	runner     *jobs.Runner
	operations operationlog.Recorder
	now        func() time.Time
}

func NewJobService(service *Service, repository *jobs.Repository, runner *jobs.Runner, operations operationlog.Recorder) (*JobService, error) {
	if service == nil || repository == nil || runner == nil {
		return nil, jobs.ErrInvalid
	}
	s := &JobService{service: service, repository: repository, runner: runner, operations: operations, now: time.Now}
	if err := runner.RegisterWithTerminalHook(ReconcileJobKind, s.runReconcile, s.recordTerminal); err != nil {
		return nil, err
	}
	return s, nil
}

func (s *JobService) SubmitReconcile(ctx context.Context, accountID string, actor account.User, correlationID string) (*jobs.Job, error) {
	now := s.now().UTC()
	id, err := jobs.NewID(now)
	if err != nil {
		return nil, err
	}
	payload, err := json.Marshal(reconcileJobPayload{AccountID: accountID})
	if err != nil {
		return nil, err
	}
	job, err := s.runner.Submit(ctx, jobs.CreateInput{
		ID: id, Kind: ReconcileJobKind, Payload: payload, ActorID: actor.ID,
		CorrelationID: correlationID, MaxAttempts: jobs.DefaultMaxAttempts, Now: now,
	})
	if err != nil {
		return nil, err
	}
	s.recordOperation(operationlog.EventWalletReconcileQueued, *job, actor.Name,
		map[string]any{"jobId": job.ID, "correlationId": job.CorrelationID}, now)
	return job, nil
}

func (s *JobService) Get(ctx context.Context, id, actorID string) (*jobs.Job, error) {
	job, err := s.repository.GetForActor(ctx, id, ReconcileJobKind, actorID)
	if err != nil {
		return nil, err
	}
	if job.Status == jobs.StatusSucceeded {
		if _, err := s.repository.ExpireIfDue(ctx, id, s.now().UTC()); err != nil {
			return nil, err
		}
		return s.repository.GetForActor(ctx, id, ReconcileJobKind, actorID)
	}
	return job, nil
}

func (s *JobService) Cancel(ctx context.Context, id, actorID string) (*jobs.Job, error) {
	if _, err := s.repository.GetForActor(ctx, id, ReconcileJobKind, actorID); err != nil {
		return nil, err
	}
	job, err := s.runner.Cancel(ctx, id, actorID)
	if err != nil {
		return nil, err
	}
	if job.Status == jobs.StatusCancelled {
		s.recordTerminal(*job)
	}
	return job, nil
}

func (s *JobService) Retry(ctx context.Context, id, actorID string) (*jobs.Job, error) {
	if _, err := s.repository.GetForActor(ctx, id, ReconcileJobKind, actorID); err != nil {
		return nil, err
	}
	return s.runner.Retry(ctx, id, actorID)
}

func (s *JobService) runReconcile(ctx context.Context, job jobs.Job, reporter jobs.Reporter) (jobs.CommitFunc, error) {
	var payload reconcileJobPayload
	if err := json.Unmarshal(job.Payload, &payload); err != nil {
		return nil, fmt.Errorf("decode wallet reconcile payload: %w", err)
	}
	if err := reporter.Progress(10); err != nil {
		return nil, err
	}
	if reporter.Cancelled() {
		return nil, ctx.Err()
	}
	return func(tx *sql.Tx) (json.RawMessage, error) {
		run, err := s.service.ReconcileOnceTx(context.Background(), tx, payload.AccountID, job.ID, job.ActorID, s.now().UTC())
		if err != nil {
			return nil, err
		}
		return json.Marshal(reconciliationResult(*run))
	}, nil
}

func (s *JobService) recordTerminal(job jobs.Job) {
	now := s.now().UTC()
	switch job.Status {
	case jobs.StatusSucceeded:
		var result struct {
			ID     string `json:"id"`
			Result string `json:"result"`
		}
		_ = json.Unmarshal(job.Result, &result)
		s.recordOperation(operationlog.EventWalletReconcile, job, job.ActorID,
			map[string]any{"jobId": job.ID, "runId": result.ID, "result": result.Result}, now)
	case jobs.StatusFailed:
		s.recordOperation(operationlog.EventWalletReconcileFailed, job, job.ActorID,
			map[string]any{"jobId": job.ID, "errorCode": job.ErrorCode}, now)
	case jobs.StatusCancelled:
		s.recordOperation(operationlog.EventWalletReconcileCancelled, job, job.ActorID,
			map[string]any{"jobId": job.ID}, now)
	}
}

func (s *JobService) recordOperation(event string, job jobs.Job, actorName string, detailValue map[string]any, now time.Time) {
	if s.operations == nil {
		return
	}
	id, err := newID(now)
	if err != nil {
		return
	}
	action := strings.TrimPrefix(event, "wallet.")
	detail, err := operationlog.NewDetail(action, nil, detailValue)
	if err != nil {
		return
	}
	_ = s.operations.RecordOperation(operationlog.Operation{
		ID: id, Event: event, ActorID: job.ActorID, ActorName: actorName,
		RecordID: &job.ID, Detail: &detail, CorrelationID: job.CorrelationID, CreatedAt: now,
	})
}

func reconciliationResult(run walletstore.ReconciliationRun) map[string]any {
	return map[string]any{
		"id": run.ID, "accountId": run.AccountID, "result": run.Result,
		"mismatchCount": run.MismatchCount, "details": run.Details,
		"actorId":   run.ActorID,
		"createdAt": run.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}
