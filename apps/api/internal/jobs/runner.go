package jobs

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"sync"
	"time"
)

type CommitFunc func(kernel.Tx) (json.RawMessage, error)

type Reporter interface {
	Progress(int) error
	Cancelled() bool
}

type Handler func(context.Context, Job, Reporter) (CommitFunc, error)

// TerminalHook observes a durable terminal transition after its transaction
// commits. Hooks are best-effort side effects such as audit recording; they
// must not be used for the consumer's durable result.
type TerminalHook func(Job)

type registration struct {
	handler  Handler
	terminal TerminalHook
}

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
	handlers map[string]registration
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
		handlers: map[string]registration{}, active: map[string]context.CancelFunc{},
		wake: make(chan struct{}, 1),
	}, nil
}

func (r *Runner) Register(kind string, handler Handler) error {
	return r.RegisterWithTerminalHook(kind, handler, nil)
}

// RegisterWithTerminalHook registers a handler and an optional observer for
// succeeded, failed, and cancelled terminal transitions.
func (r *Runner) RegisterWithTerminalHook(kind string, handler Handler, terminal TerminalHook) error {
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
	r.handlers[kind] = registration{handler: handler, terminal: terminal}
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
	cancelled, err := r.repo.RecoverCancelledDueJobs(ctx, now)
	if err != nil {
		return err
	}
	for _, job := range cancelled {
		r.notifyTerminalJob(job)
	}
	exhausted, err := r.repo.ExhaustExpiredJobs(ctx, now)
	if err != nil {
		return err
	}
	for _, job := range exhausted {
		r.notifyTerminalJob(job)
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
	job, lease, err := r.repo.Claim(ctx, id, r.instance+":"+ownerSuffix, now, r.options.LeaseDuration)
	if err != nil {
		return
	}
	r.mu.Lock()
	registered := r.handlers[job.Kind]
	r.mu.Unlock()
	if registered.handler == nil {
		_ = r.repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "job handler not registered", r.options.Now().UTC())
		return
	}

	reporter := &runnerReporter{ctx: ctx, cancelFn: cancel, repo: r.repo, lease: lease, now: r.options.Now}
	result := make(chan runnerOutcome, 1)
	go func() {
		// W9 F-007: a panicking handler must degrade to a durable job failure,
		// never crash the process. The panic is converted to the same outcome
		// path as a returned error so finish() records JOB_HANDLER_FAILED.
		defer func() {
			if recovered := recover(); recovered != nil {
				result <- runnerOutcome{err: fmt.Errorf("job handler panicked: %v", recovered)}
			}
		}()
		commit, err := registered.handler(ctx, *job, reporter)
		result <- runnerOutcome{commit: commit, err: err}
	}()

	heartbeat := time.NewTicker(r.options.HeartbeatInterval)
	defer heartbeat.Stop()
	for {
		select {
		case <-heartbeat.C:
			requested, err := r.repo.IsCancelRequested(ctx, lease)
			if err != nil {
				reporter.cancelExecution()
				if ctx.Err() == nil {
					r.abortLease(lease, err)
				}
				return
			}
			if requested {
				reporter.cancelExecution()
			}
			if err := r.repo.Heartbeat(ctx, lease, r.options.Now().UTC(), r.options.LeaseDuration); err != nil {
				reporter.cancelExecution()
				if ctx.Err() == nil {
					r.abortLease(lease, err)
				}
				return
			}
		case outcome := <-result:
			r.finish(lease, outcome)
			return
		}
	}
}

// abortLease cancels the handler and records a durable failure when the
// runner loses its heartbeat/cancellation database path. Cleanup uses a
// background context so an unrelated runner shutdown does not leave a lease
// silently running; shutdown cancellation is intentionally recoverable.
func (r *Runner) abortLease(lease Lease, cause error) {
	message := "job heartbeat failed"
	if cause != nil {
		message += ": " + cause.Error()
	}
	if err := r.repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", message, r.options.Now().UTC()); err == nil {
		r.notifyTerminal(lease.JobID)
	}
}

func (r *Runner) finish(lease Lease, outcome runnerOutcome) {
	if r.isStopping() {
		return
	}
	now := r.options.Now().UTC()
	requested, leaseErr := r.repo.IsCancelRequested(context.Background(), lease)
	if leaseErr != nil {
		r.abortLease(lease, leaseErr)
		return
	}
	if outcome.err != nil {
		if requested {
			if err := r.repo.FinalizeCancel(context.Background(), lease, now); err == nil {
				r.notifyTerminal(lease.JobID)
			}
			return
		}
		if err := r.repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", outcome.err.Error(), now); err == nil {
			r.notifyTerminal(lease.JobID)
		}
		return
	}
	if outcome.commit == nil {
		if err := r.repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", "job handler returned no commit", now); err == nil {
			r.notifyTerminal(lease.JobID)
		}
		return
	}
	if _, err := r.repo.CompleteWithCommit(context.Background(), lease, now, r.options.ResultTTL, outcome.commit); err == nil {
		r.notifyTerminal(lease.JobID)
	} else {
		requested, leaseErr := r.repo.IsCancelRequested(context.Background(), lease)
		if leaseErr != nil {
			r.abortLease(lease, leaseErr)
			return
		}
		if requested {
			if cancelErr := r.repo.FinalizeCancel(context.Background(), lease, now); cancelErr == nil {
				r.notifyTerminal(lease.JobID)
			}
			return
		}
		if failErr := r.repo.Fail(context.Background(), lease, "JOB_HANDLER_FAILED", err.Error(), now); failErr == nil {
			r.notifyTerminal(lease.JobID)
		}
	}
}

func (r *Runner) notifyTerminal(id string) {
	job, err := r.repo.Get(context.Background(), id)
	if err != nil {
		return
	}
	r.notifyTerminalJob(*job)
}

func (r *Runner) notifyTerminalJob(job Job) {
	r.mu.Lock()
	hook := r.handlers[job.Kind].terminal
	r.mu.Unlock()
	if hook != nil {
		hook(job)
	}
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
