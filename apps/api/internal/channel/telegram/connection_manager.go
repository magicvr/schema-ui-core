package telegram

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"
)

const (
	// PollingLeaseTTL is the default per-console-session lease. It covers two
	// 10-second heartbeat intervals as required by D-001.
	PollingLeaseTTL = 20 * time.Second
	// PollingLeaseSweepInterval bounds how quickly an expired lease drains the
	// receiver without coupling lease expiry to a long-poll response.
	PollingLeaseSweepInterval = time.Second
	telegramWebhookPath       = "/api/channel/telegram/webhook"
)

// ConnectionManager is the single Telegram receiver owner. It serializes
// startup, mode changes, lease-driven receiver changes, and shutdown while
// keeping all Bot API details inside this package.
type ConnectionManager struct {
	operationMu sync.Mutex
	stateMu     sync.Mutex

	runtime       *RuntimeManager
	dispatcher    *Dispatcher
	botAPI        *BotAPIClient
	pollingAPI    *BotAPIClient
	updateHandler func(context.Context, UpdatePayload) error
	now           func() time.Time

	started         bool
	lifecycleCtx    context.Context
	lifecycleCancel context.CancelFunc
	pollCancel      context.CancelFunc
	pollDone        chan struct{}
	leases          map[string]time.Time
}

// NewConnectionManager constructs the one process-level receiver owner.
func NewConnectionManager(runtime *RuntimeManager, dispatcher *Dispatcher, botAPI, pollingAPI *BotAPIClient, updateHandler func(context.Context, UpdatePayload) error) *ConnectionManager {
	if runtime != nil {
		if botAPI == nil {
			botAPI = NewBotAPIClient(runtime, nil, "")
		}
		if pollingAPI == nil {
			pollingAPI = NewPollingBotAPIClient(runtime, nil, "")
		}
	}
	manager := &ConnectionManager{
		runtime:       runtime,
		dispatcher:    dispatcher,
		botAPI:        botAPI,
		pollingAPI:    pollingAPI,
		updateHandler: updateHandler,
		now:           func() time.Time { return time.Now().UTC() },
		leases:        make(map[string]time.Time),
	}
	if runtime != nil {
		status := runtime.ConnectionStatus()
		if status.State == "" {
			if runtime.GetToken() == "" {
				status.State = ConnectionStateUnconfigured
			} else {
				status.State = ConnectionStateIdle
			}
			status.Receiver = ReceiverNone
			runtime.setConnectionStatus(status)
		}
	}
	return manager
}

// Start verifies the Bot API token, establishes the configured remote mode,
// and starts polling only when business handlers or a live lease require it.
func (m *ConnectionManager) Start(ctx context.Context) error {
	if m == nil {
		return fmt.Errorf("telegram: connection manager is unavailable")
	}
	ctx = nonNilContext(ctx)
	m.operationMu.Lock()
	defer m.operationMu.Unlock()

	m.stateMu.Lock()
	if m.started {
		m.stateMu.Unlock()
		return nil
	}
	m.started = true
	m.lifecycleCtx, m.lifecycleCancel = context.WithCancel(context.Background())
	m.stateMu.Unlock()

	if err := m.reconcileStarted(ctx); err != nil {
		m.stateMu.Lock()
		m.started = false
		cancel := m.lifecycleCancel
		m.lifecycleCancel = nil
		m.lifecycleCtx = nil
		m.stateMu.Unlock()
		if cancel != nil {
			cancel()
		}
		return err
	}

	m.stateMu.Lock()
	lifecycleCtx := m.lifecycleCtx
	m.stateMu.Unlock()
	if lifecycleCtx != nil {
		go m.watchDemand(lifecycleCtx)
	}
	return nil
}

// Ready returns a failure for a failed connection establishment while allowing
// unconfigured and idle states to start the application without a Bot API call.
func (m *ConnectionManager) Ready(ctx context.Context) error {
	if m == nil || m.runtime == nil {
		return fmt.Errorf("telegram: connection manager is unavailable")
	}
	status := m.runtime.ConnectionStatus()
	if status.State == ConnectionStateError {
		return fmt.Errorf("telegram: connection manager is in error: %s", status.LastError)
	}
	return nil
}

// Reconcile applies the current persisted settings to an already-started
// manager. Settings are persisted before this callback runs, so a failed
// target remains visible as error and never as running.
func (m *ConnectionManager) Reconcile(ctx context.Context) error {
	if m == nil {
		return fmt.Errorf("telegram: connection manager is unavailable")
	}
	ctx = nonNilContext(ctx)
	m.operationMu.Lock()
	defer m.operationMu.Unlock()

	m.stateMu.Lock()
	started := m.started
	m.stateMu.Unlock()
	if !started {
		m.publishUnstartedState()
		return nil
	}
	return m.reconcileStarted(ctx)
}

// Stop cancels and waits for the polling receiver. A drain timeout is returned
// and retained as error state; the manager never reports a timed-out receiver
// as cleanly stopped.
func (m *ConnectionManager) Stop(ctx context.Context) error {
	if m == nil {
		return nil
	}
	ctx = nonNilContext(ctx)
	m.operationMu.Lock()
	defer m.operationMu.Unlock()

	m.stateMu.Lock()
	m.started = false
	rootCancel := m.lifecycleCancel
	m.stateMu.Unlock()

	err := m.stopReceiver(ctx)
	if rootCancel != nil {
		rootCancel()
	}
	if err != nil {
		return err
	}

	m.stateMu.Lock()
	m.lifecycleCancel = nil
	m.lifecycleCtx = nil
	m.stateMu.Unlock()
	m.publishUnstartedState()
	return nil
}

// Status returns the manager's non-secret connection state.
func (m *ConnectionManager) Status() ConnectionStatus {
	if m == nil || m.runtime == nil {
		return ConnectionStatus{State: ConnectionStateUnconfigured, Receiver: ReceiverNone}
	}
	return m.runtime.ConnectionStatus()
}

// AcquireLease records one console session and starts polling when the
// already-established mode is polling and no business handler owns it.
func (m *ConnectionManager) AcquireLease(ctx context.Context, sessionID string) error {
	return m.setLease(ctx, sessionID, false)
}

// HeartbeatLease refreshes one existing console session lease. An unknown id
// is treated as a first heartbeat so a reconnect can recover without a race
// between an initial acquire request and the browser timer.
func (m *ConnectionManager) HeartbeatLease(ctx context.Context, sessionID string) error {
	return m.setLease(ctx, sessionID, false)
}

// ReleaseLease removes one console session lease and drains polling when no
// business handler or other live lease requires it.
func (m *ConnectionManager) ReleaseLease(ctx context.Context, sessionID string) error {
	return m.setLease(ctx, sessionID, true)
}

// ActiveLeaseCount returns the number of unexpired console session leases.
func (m *ConnectionManager) ActiveLeaseCount() int {
	if m == nil {
		return 0
	}
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	m.pruneExpiredLeasesLocked(m.now())
	return len(m.leases)
}

func (m *ConnectionManager) setLease(ctx context.Context, sessionID string, release bool) error {
	if m == nil {
		return fmt.Errorf("telegram: connection manager is unavailable")
	}
	sessionID = strings.TrimSpace(sessionID)
	if sessionID == "" {
		return fmt.Errorf("telegram: lease session id is required")
	}
	ctx = nonNilContext(ctx)
	m.operationMu.Lock()
	defer m.operationMu.Unlock()
	m.stateMu.Lock()
	if release {
		delete(m.leases, sessionID)
	} else {
		m.leases[sessionID] = m.now().Add(PollingLeaseTTL)
	}
	started := m.started
	m.stateMu.Unlock()
	if !started {
		return nil
	}
	return m.reconcileDemand(ctx)
}

func (m *ConnectionManager) reconcileStarted(ctx context.Context) error {
	if err := m.stopReceiver(ctx); err != nil {
		return err
	}

	if m.runtime == nil {
		return m.fail(fmt.Errorf("telegram: runtime is unavailable"), BotUser{}, ReceiverNone)
	}
	if strings.TrimSpace(m.runtime.GetToken()) == "" {
		m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateUnconfigured, Receiver: ReceiverNone})
		return nil
	}

	m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateStarting, Receiver: ReceiverNone})
	bot, err := m.botAPI.GetMe(ctx)
	if err != nil {
		return m.fail(err, BotUser{}, ReceiverNone)
	}
	m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateStarting, Receiver: ReceiverNone, BotID: bot.ID, BotUsername: bot.Username})

	switch m.runtime.GetMode() {
	case TelegramModeWebhook:
		secret := strings.TrimSpace(m.runtime.GetSecret())
		baseURL := strings.TrimSpace(m.runtime.GetWebhookPublicBaseURL())
		if secret == "" {
			return m.fail(fmt.Errorf("telegram: webhook mode requires a webhook secret"), bot, ReceiverNone)
		}
		if baseURL == "" {
			return m.fail(fmt.Errorf("telegram: webhook mode requires webhook_public_base_url"), bot, ReceiverNone)
		}
		webhookURL := strings.TrimRight(baseURL, "/") + telegramWebhookPath
		if err := m.botAPI.SetWebhook(ctx, webhookURL, secret); err != nil {
			return m.fail(err, bot, ReceiverNone)
		}
		m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateRunning, Receiver: ReceiverWebhook, BotID: bot.ID, BotUsername: bot.Username})
		return nil

	case TelegramModePolling:
		// Mode establishment is independent from receiver demand. Even an idle
		// polling manager must remove a stale remote webhook before returning.
		if err := m.botAPI.DeleteWebhook(ctx); err != nil {
			return m.fail(err, bot, ReceiverNone)
		}
		if m.hasPollingDemand() {
			m.startPolling(bot)
		} else {
			m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateIdle, Receiver: ReceiverNone, BotID: bot.ID, BotUsername: bot.Username})
		}
		return nil

	default:
		return m.fail(fmt.Errorf("telegram: unsupported mode %q", m.runtime.GetMode()), bot, ReceiverNone)
	}
}

func (m *ConnectionManager) reconcileDemand(ctx context.Context) error {
	if m.runtime == nil || m.runtime.GetMode() != TelegramModePolling {
		return nil
	}
	status := m.Status()
	demand := m.hasPollingDemand()
	if demand && status.State == ConnectionStateIdle && status.Receiver == ReceiverNone {
		m.startPolling(BotUser{ID: status.BotID, Username: status.BotUsername})
		return nil
	}
	if !demand && status.State == ConnectionStateRunning && status.Receiver == ReceiverPolling {
		if err := m.stopReceiver(ctx); err != nil {
			return err
		}
		m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateIdle, Receiver: ReceiverNone, BotID: status.BotID, BotUsername: status.BotUsername})
	}
	return nil
}

func (m *ConnectionManager) startPolling(bot BotUser) {
	m.stateMu.Lock()
	if m.pollCancel != nil || !m.started || m.lifecycleCtx == nil || m.lifecycleCtx.Err() != nil {
		m.stateMu.Unlock()
		return
	}
	lifecycleCtx := m.lifecycleCtx
	pollCtx, cancel := context.WithCancel(lifecycleCtx)
	done := make(chan struct{})
	m.pollCancel = cancel
	m.pollDone = done
	m.stateMu.Unlock()

	m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateRunning, Receiver: ReceiverPolling, BotID: bot.ID, BotUsername: bot.Username})
	go m.runPolling(pollCtx, done)
}

func (m *ConnectionManager) runPolling(ctx context.Context, done chan struct{}) {
	defer func() {
		close(done)
		m.stateMu.Lock()
		if m.pollDone == done {
			m.pollCancel = nil
			m.pollDone = nil
		}
		m.stateMu.Unlock()
	}()
	var offset int64
	for {
		updates, err := m.pollingAPI.GetUpdates(ctx, offset)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			status := m.Status()
			m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateError, Receiver: ReceiverNone, BotID: status.BotID, BotUsername: status.BotUsername, LastError: err.Error()})
			return
		}
		for _, payload := range updates {
			if m.updateHandler == nil {
				if payload.UpdateID >= offset {
					offset = payload.UpdateID + 1
				}
				continue
			}
			if err := m.updateHandler(ctx, payload); err != nil {
				var limitErr *rateLimitExceededError
				if errors.As(err, &limitErr) {
					// A rate-limited update is deliberately rejected rather than
					// persisted. Preserve the existing drop/no-auto-retry
					// semantics while keeping ordinary failures retryable.
					if payload.UpdateID >= offset {
						offset = payload.UpdateID + 1
					}
					continue
				}
				status := m.Status()
				m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateError, Receiver: ReceiverNone, BotID: status.BotID, BotUsername: status.BotUsername, LastError: err.Error()})
				return
			}
			if payload.UpdateID >= offset {
				offset = payload.UpdateID + 1
			}
		}
	}
}

func (m *ConnectionManager) watchDemand(ctx context.Context) {
	ticker := time.NewTicker(PollingLeaseSweepInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			if ctx.Err() != nil {
				return
			}
			m.operationMu.Lock()
			m.stateMu.Lock()
			started := m.started
			m.stateMu.Unlock()
			if !started || ctx.Err() != nil {
				m.operationMu.Unlock()
				return
			}
			err := m.reconcileDemand(context.Background())
			m.operationMu.Unlock()
			if err != nil {
				return
			}
		}
	}
}

func (m *ConnectionManager) stopReceiver(ctx context.Context) error {
	m.stateMu.Lock()
	cancel := m.pollCancel
	done := m.pollDone
	status := m.Status()
	if cancel != nil {
		m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateStopping, Receiver: ReceiverPolling, BotID: status.BotID, BotUsername: status.BotUsername})
	}
	m.stateMu.Unlock()
	if cancel == nil {
		return nil
	}
	cancel()
	select {
	case <-done:
		m.stateMu.Lock()
		if m.pollDone == done {
			m.pollCancel = nil
			m.pollDone = nil
		}
		m.stateMu.Unlock()
		return nil
	case <-ctx.Done():
		m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateError, Receiver: ReceiverPolling, BotID: status.BotID, BotUsername: status.BotUsername, LastError: fmt.Sprintf("receiver drain: %v", ctx.Err())})
		return fmt.Errorf("telegram: receiver drain: %w", ctx.Err())
	}
}

func (m *ConnectionManager) hasPollingDemand() bool {
	if m.dispatcher != nil && m.dispatcher.HasBusinessHandlers() {
		return true
	}
	m.stateMu.Lock()
	defer m.stateMu.Unlock()
	m.pruneExpiredLeasesLocked(m.now())
	return len(m.leases) > 0
}

func (m *ConnectionManager) pruneExpiredLeasesLocked(now time.Time) {
	for id, expiry := range m.leases {
		if !expiry.After(now) {
			delete(m.leases, id)
		}
	}
}

func (m *ConnectionManager) publishUnstartedState() {
	if m.runtime == nil {
		return
	}
	if strings.TrimSpace(m.runtime.GetToken()) == "" {
		m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateUnconfigured, Receiver: ReceiverNone})
		return
	}
	status := m.runtime.ConnectionStatus()
	m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateIdle, Receiver: ReceiverNone, BotID: status.BotID, BotUsername: status.BotUsername})
}

func (m *ConnectionManager) fail(err error, bot BotUser, receiver string) error {
	if m.runtime != nil {
		m.runtime.setConnectionStatus(ConnectionStatus{State: ConnectionStateError, Receiver: receiver, BotID: bot.ID, BotUsername: bot.Username, LastError: err.Error()})
	}
	return err
}

func nonNilContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}
