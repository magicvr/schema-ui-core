package jobs

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"
)

type CommitFunc func(*sql.Tx) (json.RawMessage, error)

type Reporter interface {
	Progress(int) error
	Cancelled() bool
}

type Handler func(context.Context, Job, Reporter) (CommitFunc, error)

type runnerOutcome struct {
	commit CommitFunc
	err    error
}

type RunnerOptions struct {
	LeaseDuration     time.Duration
	HeartbeatInterval time.Duration
	ScanInterval      time.Duration
	ResultTTL         time.Duration
	BatchSize         int
	Now               func() time.Time
}

func DefaultRunnerOptions() RunnerOptions {
	return RunnerOptions{
		LeaseDuration: 30 * time.Second, HeartbeatInterval: 10 * time.Second,
		ScanInterval: 10 * time.Second, ResultTTL: 24 * time.Hour,
		BatchSize: 32, Now: time.Now,
	}
}

type Runner struct {
	repo     *Repository
	options  RunnerOptions
	instance string

	mu       sync.Mutex
	handlers map[string]Handler
	active   map[string]context.CancelFunc
	started  bool
	stopping bool
	stop     chan struct{}
	done     chan struct{}
	wake     chan struct{}
	workers  sync.WaitGroup
}

func NewRunner(repo *Repository, options RunnerOptions) (*Runner, error) {
	if repo == nil || options.LeaseDuration <= 0 || options.HeartbeatInterval <= 0 ||
		options.ScanInterval <= 0 || options.ResultTTL <= 0 || options.BatchSize <= 0 || options.Now == nil ||
		options.HeartbeatInterval >= options.LeaseDuration {
		return nil, ErrInvalid
	}
	instance, err := randomToken(12)
	if err != nil {
		return nil, err
	}
	return &Runner{
		repo: repo, options: options, instance: instance,
		handlers: map[string]Handler{}, active: map[string]context.CancelFunc{},
		wake: make(chan struct{}, 1),
	}, nil
}

func (r *Runner) Register(kind string, handler Handler) error {
	if kind == "" || handler == nil {
		return ErrInvalid
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.started {
		return errors.New("jobs: handlers must be registered before start")
	}
	if _, exists := r.handlers[kind]; exists {
		return fmt.Errorf("jobs: handler %q already registered", kind)
	}
	r.handlers[kind] = handler
	return nil
}

func (r *Runner) Start() error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if r.started {
		return nil
	}
	r.started = true
	r.stopping = false
	r.stop = make(chan struct{})
	r.done = make(chan struct{})
	go r.loop()
	return nil
}

func (r *Runner) Stop(ctx context.Context) error {
	r.mu.Lock()
	if !r.started {
		r.mu.Unlock()
		return nil
	}
	if !r.stopping {
		r.stopping = true
		close(r.stop)
		for _, cancel := range r.active {
			cancel()
		}
	}
	done := r.done
	r.mu.Unlock()

	select {
	case <-done:
		r.mu.Lock()
		r.started = false
		r.mu.Unlock()
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (r *Runner) Submit(ctx context.Context, input CreateInput) (*Job, error) {
	job, err := r.repo.Create(ctx, input)
	if err == nil {
		r.signal()
	}
	return job, err
}

func (r *Runner) Cancel(ctx context.Context, id, actorID string) (*Job, error) {
	job, err := r.repo.RequestCancel(ctx, id, actorID, r.options.Now().UTC())
	if err != nil {
		return nil, err
	}
	if job.Status == StatusRunning {
		r.mu.Lock()
		cancel := r.active[id]
		r.mu.Unlock()
		if cancel != nil {
			cancel()
		}
	}
	r.signal()
	return job, nil
}

func (r *Runner) Retry(ctx context.Context, id, actorID string) (*Job, error) {
	job, err := r.repo.Retry(ctx, id, actorID, r.options.Now().UTC())
	if err == nil {
		r.signal()
	}
	return job, err
}

func (r *Runner) ScanOnce(ctx context.Context) error {
	now := r.options.Now().UTC()
	if _, err := r.repo.RecoverCancelledDue(ctx, now); err != nil {
		return err
	}
	if _, err := r.repo.ExhaustExpired(ctx, now); err != nil {
		return err
	}
	if _, err := r.repo.ExpireDue(ctx, now); err != nil {
		return err
	}
	candidates, err := r.repo.ListRunnable(ctx, now, r.options.BatchSize)
	if err != nil {
		return err
	}
	for _, candidate := range candidates {
		r.dispatch(candidate.ID)
	}
	return nil
}

func (r *Runner) loop() {
	defer close(r.done)
	ticker := time.NewTicker(r.options.ScanInterval)
	defer ticker.Stop()
	_ = r.ScanOnce(context.Background())
	for {
		select {
		case <-r.stop:
			r.workers.Wait()
			return
		case <-ticker.C:
			_ = r.ScanOnce(context.Background())
		case <-r.wake:
			_ = r.ScanOnce(context.Background())
		}
	}
}

func (r *Runner) dispatch(id string) {
	r.mu.Lock()
	if r.stopping {
		r.mu.Unlock()
		return
	}
	if _, exists := r.active[id]; exists {
		r.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	r.active[id] = cancel
	r.workers.Add(1)
	r.mu.Unlock()

	go func() {
		defer r.workers.Done()
		defer func() {
			r.mu.Lock()
			delete(r.active, id)
			r.mu.Unlock()
		}()
		r.execute(ctx, cancel, id)
	}()
}

func (r *Runner) execute(ctx context.Context, cancel context.CancelFunc, id string) {
	ownerSuffix, err := randomToken(8)
	if err != nil {
		return
	}
	now := r.options.Now().UTC()
	job, lease, err := r.repo.Claim(context.Background(), id, r.instance+":"+ownerSuffix, now, r.options.LeaseDuration)
	if err != nil {
		return
	}
	r.mu.Lock()
	handler := r.handlers[job.Kind]
	r.mu.Unlock()
	if handler == nil {
		_ = r.repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "job handler not registered", r.options.Now().UTC())
		return
	}

	reporter := &runnerReporter{ctx: ctx, cancelFn: cancel, repo: r.repo, lease: lease, now: r.options.Now}
	result := make(chan runnerOutcome, 1)
	go func() {
		commit, err := handler(ctx, *job, reporter)
		result <- runnerOutcome{commit: commit, err: err}
	}()

	heartbeat := time.NewTicker(r.options.HeartbeatInterval)
	defer heartbeat.Stop()
	for {
		select {
		case <-heartbeat.C:
			requested, err := r.repo.IsCancelRequested(context.Background(), lease)
			if err != nil {
				return
			}
			if requested {
				reporter.cancelExecution()
			}
			if err := r.repo.Heartbeat(context.Background(), lease, r.options.Now().UTC(), r.options.LeaseDuration); err != nil {
				reporter.cancelExecution()
				return
			}
		case outcome := <-result:
			r.finish(lease, outcome)
			return
		}
	}
}

func (r *Runner) finish(lease Lease, outcome runnerOutcome) {
	if r.isStopping() {
		return
	}
	now := r.options.Now().UTC()
	requested, leaseErr := r.repo.IsCancelRequested(context.Background(), lease)
	if leaseErr != nil {
		return
	}
	if outcome.err != nil {
		if requested {
			_ = r.repo.FinalizeCancel(context.Background(), lease, now)
			return
		}
		_ = r.repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", outcome.err.Error(), now)
		return
	}
	if outcome.commit == nil {
		_ = r.repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "job handler returned no commit", now)
		return
	}
	_, _ = r.repo.CompleteWithCommit(context.Background(), lease, now, r.options.ResultTTL, outcome.commit)
}

func (r *Runner) isStopping() bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.stopping
}

func (r *Runner) signal() {
	select {
	case r.wake <- struct{}{}:
	default:
	}
}

func randomToken(size int) (string, error) {
	value := make([]byte, size)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return hex.EncodeToString(value), nil
}

type runnerReporter struct {
	ctx      context.Context
	cancelFn context.CancelFunc
	repo     *Repository
	lease    Lease
	now      func() time.Time
}

func (r *runnerReporter) Progress(value int) error {
	return r.repo.UpdateProgress(context.Background(), r.lease, value, r.now().UTC())
}

func (r *runnerReporter) Cancelled() bool {
	return r.ctx.Err() != nil
}

func (r *runnerReporter) cancelExecution() {
	r.cancelFn()
}
