// Wallet/ledger HTTP surface (S-14 · GOAL-019 D-002 §3): account lifecycle,
// balance mutations (adjust/freeze/unfreeze — the money-path permissions are
// wallet.adjust), the immutable ledger and reconciliation runs. All mutations
// are audited via operationlog.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/concurrency"
	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
)

// WalletService is the surface the wallet routes consume (satisfied
// structurally by the admin.wallet module Service — the direction is
// module → handler).
type WalletService interface {
	ListAccounts(q, ownerType string, page, pageSize int) ([]walletstore.Account, int, error)
	GetAccount(id string) (*walletstore.Account, error)
	CreateAccount(ownerType, ownerID, currency string, now time.Time) (*walletstore.Account, error)
	UpdateStatus(id, status string, version int64, now time.Time) (*walletstore.Account, error)
	ListEntries(accountID, entryType, q string, page, pageSize int) ([]walletstore.LedgerEntry, int, error)
	// GetOrCreateUserAccount returns the user account for ownerID, creating a
	// zero-balance account when absent (GOAL-020 D-001 §1). The bool reports
	// whether this call created the account (auto-audit marker).
	GetOrCreateUserAccount(ownerID string, now time.Time) (*walletstore.Account, bool, error)
	GetUserAccountByOwner(ownerID string) (*walletstore.Account, error)
	Mutate(id string, in walletstore.LedgerEntryInput, now time.Time) (*walletstore.Account, *walletstore.LedgerEntry, bool, error)
	Reconcile(accountID, actorID string, now time.Time) (*walletstore.ReconciliationRun, error)
	ListReconcileRuns(page, pageSize int) ([]walletstore.ReconciliationRun, int, error)
}

// WalletJobService is the actor-scoped async boundary consumed by wallet
// reconciliation routes.
type WalletJobService interface {
	SubmitReconcile(ctx context.Context, accountID string, actor account.User, correlationID string) (*jobs.Job, error)
	Get(ctx context.Context, id, actorID string) (*jobs.Job, error)
	Cancel(ctx context.Context, id, actorID string) (*jobs.Job, error)
	Retry(ctx context.Context, id, actorID string) (*jobs.Job, error)
}

// WalletRoutes returns the admin.wallet HTTP surface.
func WalletRoutes(a *auth.Authenticator, service WalletService, jobService WalletJobService, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	var routes []kernel.RouteContribution
	add := func(method, pattern string, h http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              h,
		})
	}

	// Accounts: list (unified list envelope; amounts in integer min-units).
	add("GET", "/api/wallet/accounts", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "wallet.read"); !ok {
			return
		}
		page, ok := intParam(r.URL.Query().Get("page"), 1)
		if !ok {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		pageSize, ok := intParam(r.URL.Query().Get("pageSize"), DefaultPageSize)
		if !ok || pageSize > maxPageSize {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
			return
		}
		accounts, total, err := service.ListAccounts(r.URL.Query().Get("q"), r.URL.Query().Get("ownerType"), page, pageSize)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list wallet accounts")
			return
		}
		rows := make([]map[string]any, 0, len(accounts))
		for _, a := range accounts {
			rows = append(rows, accountToMap(a))
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
	})))

	// Accounts: read-only lookup by owner (W15-F11). Missing → 404.
	add("GET", "/api/wallet/by-owner/{ownerId}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "wallet.read"); !ok {
			return
		}
		ownerID := strings.TrimSpace(r.PathValue("ownerId"))
		if ownerID == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_WALLET_OWNER", "ownerId is required")
			return
		}
		account, err := service.GetUserAccountByOwner(ownerID)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		writeWalletAccount(w, http.StatusOK, *account)
	})))

	// Explicit create (POST). GET must stay read-only.
	add("POST", "/api/wallet/by-owner/{ownerId}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.adjust")
		if !ok {
			return
		}
		ownerID := strings.TrimSpace(r.PathValue("ownerId"))
		if ownerID == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_WALLET_OWNER", "ownerId is required")
			return
		}
		now := time.Now().UTC()
		account, created, err := service.GetOrCreateUserAccount(ownerID, now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		if created {
			recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate, "account-create", map[string]any{"accountId": account.ID, "ownerId": account.OwnerID, "auto": true}, now)
		}
		writeWalletAccount(w, http.StatusOK, *account)
	})))

	// Accounts: adjust a user account by owner, auto-creating it first
	// (GOAL-020 D-001 §1).
	add("POST", "/api/wallet/by-owner/{ownerId}/adjust", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.adjust")
		if !ok {
			return
		}
		ownerID := strings.TrimSpace(r.PathValue("ownerId"))
		if ownerID == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_WALLET_OWNER", "ownerId is required")
			return
		}
		var body struct {
			AmountDelta    int64  `json:"amountDelta"`
			Memo           string `json:"memo"`
			IdempotencyKey string `json:"idempotencyKey"`
			RefType        string `json:"refType"`
			RefID          string `json:"refId"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_LEDGER_ENTRY", "body must be JSON")
			return
		}
		now := time.Now().UTC()
		account, created, err := service.GetOrCreateUserAccount(ownerID, now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		// Audit the auto-open IMMEDIATELY (GOAL-020 A-003 F-002): the account
		// row exists from this point on, so the create must be on record even
		// when the subsequent adjust fails (invalid memo / zero amount /
		// over-draft) — otherwise the account could stay forever without its
		// wallet.account-create event.
		if created {
			recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate, "account-create", map[string]any{"accountId": account.ID, "ownerId": ownerID, "auto": true}, now)
		}
		account, entry, replayed, err := service.Mutate(account.ID, walletstore.LedgerEntryInput{
			EntryType: walletstore.EntryAdjust, AmountDelta: body.AmountDelta,
			RefType: strings.TrimSpace(body.RefType), RefID: strings.TrimSpace(body.RefID),
			IdempotencyKey: strings.TrimSpace(body.IdempotencyKey), Memo: strings.TrimSpace(body.Memo),
			ActorID: user.ID, ActorName: user.Name,
		}, now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		if !replayed {
			recordWalletEvent(operations, user, operationlog.EventWalletAdjust, "adjust", map[string]any{"accountId": account.ID, "entryId": entry.ID, "amountDelta": entry.AmountDelta}, now)
		}
		writeWalletMutation(w, *account, *entry, replayed)
	})))

	// Accounts: create (audited).
	add("POST", "/api/wallet/accounts", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.write")
		if !ok {
			return
		}
		var body struct {
			OwnerType string `json:"ownerType"`
			OwnerID   string `json:"ownerId"`
			Currency  string `json:"currency"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_WALLET_BODY", "body must be JSON with ownerType and ownerId")
			return
		}
		ownerType := strings.TrimSpace(body.OwnerType)
		// GOAL-020 D-001 §2: user wallet accounts are created automatically by
		// the system (get-or-create); manual creation is business/system only.
		if ownerType == walletstore.OwnerUser {
			writeLocalizedError(w, r, http.StatusConflict, "WALLET_USER_AUTO_ONLY", "user wallet accounts are created automatically")
			return
		}
		now := time.Now().UTC()
		account, err := service.CreateAccount(ownerType, strings.TrimSpace(body.OwnerID), strings.TrimSpace(body.Currency), now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate, "account-create", map[string]any{"accountId": account.ID, "ownerType": account.OwnerType, "ownerId": account.OwnerID}, now)
		writeWalletAccount(w, http.StatusCreated, *account)
	})))

	// Accounts: status update (audited). Optimistic lock: the caller passes
	// the version observed when loading the row.
	add("PATCH", "/api/wallet/accounts/{id}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.write")
		if !ok {
			return
		}
		var body struct {
			Status          string `json:"status"`
			ExpectedVersion *int64 `json:"expectedVersion"`
			Version         *int64 `json:"version"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_WALLET_BODY", "body must be JSON with status and version")
			return
		}
		status := strings.TrimSpace(body.Status)
		if status != walletstore.StatusActive && status != walletstore.StatusDisabled {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_WALLET_STATUS", "status must be active or disabled")
			return
		}
		expectedVersion, err := concurrency.ResolveExpectedVersion(r.Header.Values("If-Match"), body.ExpectedVersion, body.Version)
		if err != nil {
			if errors.Is(err, concurrency.ErrPreconditionRequired) {
				writeLocalizedError(w, r, http.StatusPreconditionRequired, "PRECONDITION_REQUIRED", "provide If-Match or expectedVersion")
				return
			}
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PRECONDITION", "version preconditions must be valid and agree")
			return
		}
		now := time.Now().UTC()
		account, err := service.UpdateStatus(r.PathValue("id"), status, expectedVersion, now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		recordWalletEvent(operations, user, operationlog.EventWalletAccountUpdate, "account-update", map[string]any{"accountId": account.ID, "status": account.Status}, now)
		writeWalletAccount(w, http.StatusOK, *account)
	})))

	// Entries: canonical REST path variant (D-002 §3).
	add("GET", "/api/wallet/accounts/{id}/entries", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		walletListEntries(w, r, service, r.PathValue("id"))
	})))

	// Entries: page-facing query variant (route-bound dataSource).
	add("GET", "/api/wallet/entries", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		walletListEntries(w, r, service, r.URL.Query().Get("accountId"))
	})))

	// Adjust: balance mutation with signed amount (audited).
	add("POST", "/api/wallet/accounts/{id}/adjust", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		walletMutate(w, r, service, operations, walletstore.EntryAdjust, "adjust")
	})))

	// Freeze: available → frozen transfer (audited).
	add("POST", "/api/wallet/accounts/{id}/freeze", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		walletMutate(w, r, service, operations, walletstore.EntryFreeze, "freeze")
	})))

	// Unfreeze: frozen → available transfer (audited).
	add("POST", "/api/wallet/accounts/{id}/unfreeze", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		walletMutate(w, r, service, operations, walletstore.EntryUnfreeze, "unfreeze")
	})))

	// Deduct-frozen: consume from the frozen bucket atomically (A-008 F-001 ·
	// GOAL-021 D-001 §1). available stays untouched — the pre-authorized money
	// is consumed, never re-exposed.
	add("POST", "/api/wallet/accounts/{id}/deduct-frozen", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		walletMutate(w, r, service, operations, walletstore.EntryDeductFrozen, "deduct-frozen")
	})))

	// Reconcile: durable async ledger chain check.
	add("POST", "/api/wallet/reconcile", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.read")
		if !ok {
			return
		}
		var body struct {
			AccountID string `json:"accountId"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		_ = json.NewDecoder(r.Body).Decode(&body)
		correlationID := requestid.FromContext(r.Context())
		if correlationID == "" {
			correlationID = requestid.New()
		}
		job, err := jobService.SubmitReconcile(r.Context(), strings.TrimSpace(body.AccountID), user, correlationID)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		writeJSON(w, http.StatusAccepted, walletJobToMap(*job))
	})))

	add("GET", "/api/wallet/jobs/{id}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.read")
		if !ok {
			return
		}
		job, err := jobService.Get(r.Context(), r.PathValue("id"), user.ID)
		if err != nil {
			writeWalletJobError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, walletJobToMap(*job))
	})))

	add("POST", "/api/wallet/jobs/{id}/cancel", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.read")
		if !ok {
			return
		}
		job, err := jobService.Cancel(r.Context(), r.PathValue("id"), user.ID)
		if err != nil {
			writeWalletJobError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, walletJobToMap(*job))
	})))

	add("POST", "/api/wallet/jobs/{id}/retry", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.read")
		if !ok {
			return
		}
		job, err := jobService.Retry(r.Context(), r.PathValue("id"), user.ID)
		if err != nil {
			writeWalletJobError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, walletJobToMap(*job))
	})))

	add("GET", "/api/wallet/jobs/{id}/result", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.read")
		if !ok {
			return
		}
		job, err := jobService.Get(r.Context(), r.PathValue("id"), user.ID)
		if err != nil {
			writeWalletJobError(w, r, err)
			return
		}
		switch job.Status {
		case jobs.StatusQueued, jobs.StatusRunning:
			writeLocalizedError(w, r, http.StatusConflict, "JOB_RESULT_NOT_READY", "job result is not ready")
		case jobs.StatusExpired:
			writeLocalizedError(w, r, http.StatusGone, "JOB_RESULT_EXPIRED", "job result has expired")
		case jobs.StatusFailed, jobs.StatusCancelled:
			writeJSON(w, http.StatusOK, walletJobToMap(*job))
		case jobs.StatusSucceeded:
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", "wallet-reconcile-"+job.ID+".json"))
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(job.Result)
		default:
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "invalid job status")
		}
	})))

	// Reconcile runs: list.
	add("GET", "/api/wallet/reconcile/runs", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "wallet.read"); !ok {
			return
		}
		page, ok := intParam(r.URL.Query().Get("page"), 1)
		if !ok {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
			return
		}
		pageSize, ok := intParam(r.URL.Query().Get("pageSize"), DefaultPageSize)
		if !ok || pageSize > maxPageSize {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
			return
		}
		runs, total, err := service.ListReconcileRuns(page, pageSize)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list reconciliation runs")
			return
		}
		rows := make([]map[string]any, 0, len(runs))
		for _, run := range runs {
			rows = append(rows, reconcileRunToMap(run))
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
	})))

	return routes
}

// walletListEntries serves one account's ledger (shared by the path and query
// route variants).
func walletListEntries(w http.ResponseWriter, r *http.Request, service WalletService, accountID string) {
	if _, ok := requirePermission(w, r, "wallet.read"); !ok {
		return
	}
	if strings.TrimSpace(accountID) == "" {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_WALLET_ACCOUNT", "accountId is required")
		return
	}
	page, ok := intParam(r.URL.Query().Get("page"), 1)
	if !ok {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE", "page must be a positive integer")
		return
	}
	pageSize, ok := intParam(r.URL.Query().Get("pageSize"), DefaultPageSize)
	if !ok || pageSize > maxPageSize {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_PAGE_SIZE", "pageSize must be a positive integer not exceeding 100")
		return
	}
	entryType := strings.TrimSpace(r.URL.Query().Get("entryType"))
	q := strings.TrimSpace(r.URL.Query().Get("q"))
	entries, total, err := service.ListEntries(accountID, entryType, q, page, pageSize)
	if err != nil {
		if errors.Is(err, walletstore.ErrNotFound) {
			writeLocalizedError(w, r, http.StatusNotFound, "WALLET_NOT_FOUND", "wallet account not found")
			return
		}
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list wallet entries")
		return
	}
	rows := make([]map[string]any, 0, len(entries))
	for _, e := range entries {
		rows = append(rows, entryToMap(e))
	}
	writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
}

// walletMutate handles adjust/freeze/unfreeze with the shared audit path.
func walletMutate(w http.ResponseWriter, r *http.Request, service WalletService, operations operationlog.Recorder, entryType, eventSuffix string) {
	user, ok := requirePermission(w, r, "wallet.adjust")
	if !ok {
		return
	}
	var body struct {
		AmountDelta    int64  `json:"amountDelta"`
		Amount         int64  `json:"amount"`
		Memo           string `json:"memo"`
		IdempotencyKey string `json:"idempotencyKey"`
		RefType        string `json:"refType"`
		RefID          string `json:"refId"`
	}
	r.Body = http.MaxBytesReader(w, r.Body, 8<<10)
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_LEDGER_ENTRY", "body must be JSON")
		return
	}
	delta := body.AmountDelta
	if entryType != walletstore.EntryAdjust {
		delta = body.Amount
	}
	now := time.Now().UTC()
	account, entry, replayed, err := service.Mutate(r.PathValue("id"), walletstore.LedgerEntryInput{
		EntryType:      entryType,
		AmountDelta:    delta,
		RefType:        strings.TrimSpace(body.RefType),
		RefID:          strings.TrimSpace(body.RefID),
		IdempotencyKey: strings.TrimSpace(body.IdempotencyKey),
		Memo:           strings.TrimSpace(body.Memo),
		ActorID:        user.ID,
		ActorName:      user.Name,
	}, now)
	if err != nil {
		writeWalletError(w, r, err)
		return
	}
	if !replayed {
		recordWalletEvent(operations, user, "wallet."+eventSuffix, eventSuffix, map[string]any{"accountId": account.ID, "entryId": entry.ID, "amountDelta": entry.AmountDelta}, now)
	}
	writeWalletMutation(w, *account, *entry, replayed)
}

// writeWalletError maps wallet domain errors to the frozen wire codes.
func writeWalletError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, walletstore.ErrNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "WALLET_NOT_FOUND", "wallet account not found")
	case errors.Is(err, walletstore.ErrOwnerTaken):
		writeLocalizedError(w, r, http.StatusConflict, "WALLET_OWNER_TAKEN", "an account for that owner already exists")
	case errors.Is(err, walletstore.ErrDisabled):
		writeLocalizedError(w, r, http.StatusConflict, "WALLET_DISABLED", "wallet account is disabled")
	case errors.Is(err, walletstore.ErrInsufficient):
		writeLocalizedError(w, r, http.StatusConflict, "INSUFFICIENT_BALANCE", "insufficient balance for this mutation")
	case errors.Is(err, walletstore.ErrVersionConflict):
		writeLocalizedError(w, r, http.StatusConflict, "LEDGER_VERSION_CONFLICT", "the account changed concurrently; reload and retry")
	case errors.Is(err, walletstore.ErrIdempotencyConflict):
		writeLocalizedError(w, r, http.StatusConflict, "LEDGER_IDEMPOTENCY_CONFLICT", "idempotency key was already used with a different payload")
	case errors.Is(err, walletstore.ErrInvalidEntry):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_LEDGER_ENTRY", "invalid ledger entry")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "wallet operation failed")
	}
}

func writeWalletJobError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, jobs.ErrNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "JOB_NOT_FOUND", "job not found")
	case errors.Is(err, jobs.ErrNotCancellable):
		writeLocalizedError(w, r, http.StatusConflict, "JOB_NOT_CANCELLABLE", "job cannot be cancelled")
	case errors.Is(err, jobs.ErrNotRetryable):
		writeLocalizedError(w, r, http.StatusConflict, "JOB_NOT_RETRYABLE", "job cannot be retried")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "job operation failed")
	}
}

func recordWalletEvent(operations operationlog.Recorder, user account.User, event, action string, fields map[string]any, now time.Time) {
	recordAudit(operations, user, event, "", auditDetail(action, fields), now, nil)
}

func accountToMap(a walletstore.Account) map[string]any {
	return map[string]any{
		"id":               a.ID,
		"ownerType":        a.OwnerType,
		"ownerId":          a.OwnerID,
		"currency":         a.Currency,
		"balanceTotal":     a.BalanceTotal,
		"balanceAvailable": a.BalanceAvailable,
		"balanceFrozen":    a.BalanceFrozen,
		"status":           a.Status,
		"version":          a.Version,
		"updatedAt":        a.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"decimals":         2,
	}
}

func writeWalletAccount(w http.ResponseWriter, status int, account walletstore.Account) {
	w.Header().Set("ETag", concurrency.ETag(account.Version))
	writeJSON(w, status, accountToMap(account))
}

func writeWalletMutation(w http.ResponseWriter, account walletstore.Account, entry walletstore.LedgerEntry, replayed bool) {
	w.Header().Set("ETag", concurrency.ETag(account.Version))
	operation := map[string]any{
		"operationId":     entry.ID,
		"state":           "succeeded",
		"replayed":        replayed,
		"resourceVersion": account.Version,
	}
	if entry.IdempotencyKey != "" {
		operation["idempotencyKey"] = entry.IdempotencyKey
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"account":   accountToMap(account),
		"entry":     entryToMap(entry),
		"operation": operation,
	})
}

func entryToMap(e walletstore.LedgerEntry) map[string]any {
	return map[string]any{
		"id":                 e.ID,
		"accountId":          e.AccountID,
		"entryType":          e.EntryType,
		"amountDelta":        e.AmountDelta,
		"balanceAfterTotal":  e.BalanceAfterTotal,
		"balanceAfterAvail":  e.BalanceAfterAvail,
		"balanceAfterFrozen": e.BalanceAfterFrozen,
		"refType":            e.RefType,
		"refId":              e.RefID,
		"memo":               e.Memo,
		"actorId":            e.ActorID,
		"actorName":          e.ActorName,
		"createdAt":          e.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func reconcileRunToMap(r walletstore.ReconciliationRun) map[string]any {
	return map[string]any{
		"id":            r.ID,
		"accountId":     r.AccountID,
		"result":        r.Result,
		"mismatchCount": r.MismatchCount,
		"details":       r.Details,
		"actorId":       r.ActorID,
		"createdAt":     r.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func walletJobToMap(job jobs.Job) map[string]any {
	row := map[string]any{
		"id": job.ID, "kind": job.Kind, "status": job.Status,
		"progress": job.Progress, "attempt": job.Attempt, "maxAttempts": job.MaxAttempts,
		"cancelRequested": job.CancelRequested, "createdAt": job.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt": job.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
	if job.ErrorCode != "" {
		row["error"] = map[string]any{"code": job.ErrorCode, "message": job.ErrorMessage}
	}
	if job.FinishedAt != nil {
		row["finishedAt"] = job.FinishedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	}
	if job.Status == jobs.StatusSucceeded {
		row["resultUrl"] = "/api/wallet/jobs/" + job.ID + "/result"
	}
	return row
}
