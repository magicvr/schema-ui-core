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
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/concurrency"
	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/voucher"
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
	// Prepaid vouchers surface (VP-029 R3 · GOAL-003).
	GenerateVouchers(ctx context.Context, batchID string, count int, amount int64, currency string, expiresAt *time.Time, now time.Time) ([]voucher.GeneratedVoucher, error)
	ListVouchers(ctx context.Context, batchID, status string, page, pageSize int) ([]voucher.Voucher, int, error)
	GetVoucher(ctx context.Context, id string) (*voucher.Voucher, error)
	VoidVoucher(ctx context.Context, id string, now time.Time) error
	// RedeemForUser credits the session user's owner_type=user ledger
	// (VP-029 R5 · GOAL-005). Identity is supplied by the handler, never
	// from the request body.
	RedeemForUser(ctx context.Context, userID, actorName, code string, now time.Time) (*voucher.RedeemResult, error)
}

// WalletJobService is the actor-scoped async boundary consumed by wallet
// reconciliation routes.
type WalletJobService interface {
	SubmitReconcile(ctx context.Context, accountID string, actor account.User, correlationID string) (*jobs.Job, error)
	Get(ctx context.Context, id, actorID string) (*jobs.Job, error)
	Cancel(ctx context.Context, id, actorID string) (*jobs.Job, error)
	Retry(ctx context.Context, id, actorID string) (*jobs.Job, error)
}

// OwnerExistsFunc reports whether a user owner id refers to an existing
// account (W13 F-012 · GOAL-013 A-001). The auto-create paths must not mint
// ledger rows for nonexistent owners (orphan account books). nil disables the
// gate (bare test environments only).
type OwnerExistsFunc func(ownerID string) bool

// WalletRoutes returns the admin.wallet HTTP surface.
func WalletRoutes(a *auth.Authenticator, service WalletService, jobService WalletJobService, operations operationlog.Recorder, moduleID string, ownerExists OwnerExistsFunc) []kernel.RouteContribution {
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
		// W13 F-012: an explicit create must reference a REAL user — an
		// unknown owner would otherwise open an orphan account book.
		if ownerExists != nil && !ownerExists(ownerID) {
			writeLocalizedError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "no user with that ownerId")
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
		// W13 F-012: the adjust path auto-creates on first sight — same
		// existence gate as the explicit create so a typo'd ownerId can never
		// silently open an orphan ledger.
		if ownerExists != nil && !ownerExists(ownerID) {
			writeLocalizedError(w, r, http.StatusNotFound, "USER_NOT_FOUND", "no user with that ownerId")
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

	// Reconcile: durable async ledger chain check. Submit/cancel/retry are
	// write operations (they queue jobs against the whole ledger) and require
	// wallet.write — a read-only role must not be able to trigger or mutate
	// reconciliation jobs (W11 F-006).
	add("POST", "/api/wallet/reconcile", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.write")
		if !ok {
			return
		}
		var body struct {
			AccountID string `json:"accountId"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		// W11 F-006: a garbage body used to silently fall through to the
		// empty-accountId sentinel and queue a FULL-ledger reconcile. Decode
		// failures are now a hard 400 — an empty body (io.EOF) stays valid
		// and means the documented all-accounts reconcile; anything else
		// that is not a JSON object is rejected.
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil && !errors.Is(err, io.EOF) {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BODY", "expected a JSON object with an optional accountId")
			return
		}
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
		user, ok := requirePermission(w, r, "wallet.write")
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
		user, ok := requirePermission(w, r, "wallet.write")
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

	// Vouchers: generate batch (audited). Requires wallet.voucher.issue.
	// E-008: batchId is OPTIONAL — omitted, the server generates a unique
	// VB-… id (operators no longer invent ids); amount is entered in CNY yuan
	// with up to 2 decimal places and converted to min units (分) internally.
	add("POST", "/api/wallet/vouchers/batches", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.voucher.issue")
		if !ok {
			return
		}
		var body struct {
			BatchID   string          `json:"batchId"`
			Count     int             `json:"count"`
			Amount    json.RawMessage `json:"amount"`
			Currency  string          `json:"currency"`
			ExpiresAt json.RawMessage `json:"expiresAt"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_BODY", "body must be JSON with count and amount")
			return
		}
		currency := strings.TrimSpace(body.Currency)
		if currency != "" && currency != walletstore.DefaultCurrency {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_PARAMS", "currency must be CNY")
			return
		}
		if body.Count <= 0 || body.Count > 1000 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_PARAMS", "invalid count (1-1000)")
			return
		}
		amountCents, okAmount := parseYuanToCents(body.Amount)
		if !okAmount || amountCents <= 0 {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_PARAMS", "amount must be a CNY amount in yuan with up to 2 decimal places and greater than zero (e.g. 12.5)")
			return
		}
		now := time.Now().UTC()
		batchID := strings.TrimSpace(body.BatchID)
		if batchID == "" {
			id, err := voucher.NewBatchID(now)
			if err != nil {
				writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not generate voucher batch id")
				return
			}
			batchID = id
		}
		// E-009: expiresAt accepts legacy Unix seconds or a UTC date
		// (YYYY-MM-DD) picked in the admin form — the server converts the date
		// to end-of-day seconds (valid through 23:59:59 of the chosen day).
		expiry, okExpiry := parseVoucherExpiry(body.ExpiresAt)
		if !okExpiry {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_PARAMS", "expiresAt must be omitted, Unix seconds, or a YYYY-MM-DD UTC date (valid through 23:59:59 that day; range 2001-09-09..2099-12-31)")
			return
		}
		var exp *time.Time
		if expiry != nil {
			t := time.Unix(*expiry, 0).UTC()
			exp = &t
		}
		generated, err := service.GenerateVouchers(r.Context(), batchID, body.Count, amountCents, currency, exp, now)
		if err != nil {
			// A-005 F-004 (A-008): repeated batchId is a conflict (0065
			// registry), not an internal error. Explicit batchId remains
			// supported for API callers; the admin form omits it.
			if errors.Is(err, voucher.ErrVoucherBatchExists) {
				writeLocalizedError(w, r, http.StatusConflict, "VOUCHER_BATCH_EXISTS", "a batch with that batchId already exists")
				return
			}
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not generate voucher batch")
			return
		}
		// Audited without plaintext codes!
		recordWalletEvent(operations, user, "records.create", "batch-generate", map[string]any{
			"batchId":  batchID,
			"count":    body.Count,
			"amount":   amountCents,
			"currency": currency,
		}, now)

		items := make([]map[string]any, len(generated))
		for i, g := range generated {
			item := map[string]any{
				"id":         g.Voucher.ID,
				"batchId":    g.Voucher.BatchID,
				"codePrefix": g.Voucher.CodePrefix,
				"amount":     g.Voucher.Amount,
				"currency":   g.Voucher.Currency,
				"status":     string(g.Voucher.Status),
				"code":       g.Code, // One-time plaintext returned only here
				"createdAt":  g.Voucher.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			}
			if g.Voucher.ExpiresAt != nil {
				item["expiresAt"] = g.Voucher.ExpiresAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
			}
			items[i] = item
		}
		writeJSON(w, http.StatusCreated, map[string]any{
			"items": items,
			"total": len(items),
		})
	})))

	// Vouchers: list (requires wallet.read). Never returns plaintext codes.
	add("GET", "/api/wallet/vouchers", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
		batchID := r.URL.Query().Get("batchId")
		status := r.URL.Query().Get("status")
		vouchers, total, err := service.ListVouchers(r.Context(), batchID, status, page, pageSize)
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not list vouchers")
			return
		}
		rows := make([]map[string]any, len(vouchers))
		for i, v := range vouchers {
			row := map[string]any{
				"id":         v.ID,
				"batchId":    v.BatchID,
				"codePrefix": v.CodePrefix,
				"amount":     v.Amount,
				"currency":   v.Currency,
				"status":     string(v.Status),
				"createdAt":  v.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
				"updatedAt":  v.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			}
			if v.ExpiresAt != nil {
				row["expiresAt"] = v.ExpiresAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
			}
			if v.RedeemedBy != nil {
				row["redeemedBy"] = *v.RedeemedBy
			}
			if v.RedeemedAt != nil {
				row["redeemedAt"] = v.RedeemedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
			}
			rows[i] = row
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
	})))

	// Vouchers: get by id (requires wallet.read).
	add("GET", "/api/wallet/vouchers/{id}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "wallet.read"); !ok {
			return
		}
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_ID", "id is required")
			return
		}
		v, err := service.GetVoucher(r.Context(), id)
		if errors.Is(err, voucher.ErrNotFound) {
			writeLocalizedError(w, r, http.StatusNotFound, "VOUCHER_NOT_FOUND", "voucher not found")
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not get voucher")
			return
		}
		row := map[string]any{
			"id":         v.ID,
			"batchId":    v.BatchID,
			"codePrefix": v.CodePrefix,
			"amount":     v.Amount,
			"currency":   v.Currency,
			"status":     string(v.Status),
			"createdAt":  v.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
			"updatedAt":  v.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		}
		if v.ExpiresAt != nil {
			row["expiresAt"] = v.ExpiresAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
		}
		if v.RedeemedBy != nil {
			row["redeemedBy"] = *v.RedeemedBy
		}
		if v.RedeemedAt != nil {
			row["redeemedAt"] = v.RedeemedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
		}
		writeJSON(w, http.StatusOK, row)
	})))

	// Vouchers: void voucher (audited). Requires wallet.voucher.issue.
	add("POST", "/api/wallet/vouchers/{id}/void", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.voucher.issue")
		if !ok {
			return
		}
		id := strings.TrimSpace(r.PathValue("id"))
		if id == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_ID", "id is required")
			return
		}
		now := time.Now().UTC()
		err := service.VoidVoucher(r.Context(), id, now)
		if errors.Is(err, voucher.ErrNotFound) {
			writeLocalizedError(w, r, http.StatusNotFound, "VOUCHER_NOT_FOUND", "voucher not found")
			return
		}
		if errors.Is(err, voucher.ErrVoucherAlreadyRedeemed) {
			writeLocalizedError(w, r, http.StatusConflict, "VOUCHER_ALREADY_REDEEMED", "cannot void already redeemed voucher")
			return
		}
		if err != nil {
			writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "could not void voucher")
			return
		}
		recordWalletEvent(operations, user, "records.update", "voucher-void", map[string]any{"voucherId": id}, now)
		writeJSON(w, http.StatusOK, map[string]any{"id": id, "status": "void"})
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

// voucher expiry window (A-005 F-005 / A-008): 2001-09-09T00:00:00Z ..
// 2100-01-01T00:00:00Z in Unix seconds.
const (
	voucherExpiryMinUnix int64 = 1_000_000_000
	voucherExpiryMaxUnix int64 = 4_102_444_800
)

// parseVoucherExpiry normalizes an expiresAt payload to Unix seconds:
//   - absent / "" / <=0 seconds → nil (no expiry, historical semantics)
//   - numeric Unix seconds (JSON int, exponent form, or quoted digits)
//   - a "YYYY-MM-DD" UTC date (E-009): converted to 23:59:59 UTC of that day,
//     so the whole chosen day stays redeemable
//
// Out-of-window values and malformed input return ok=false (fail-closed).
func parseVoucherExpiry(raw json.RawMessage) (*int64, bool) {
	if len(raw) == 0 {
		return nil, true
	}
	s := strings.TrimSpace(string(raw))
	if s == "" {
		return nil, true
	}
	var sec int64
	if s[0] == '"' {
		var str string
		if err := json.Unmarshal(raw, &str); err != nil {
			return nil, false
		}
		s = strings.TrimSpace(str)
		if s == "" {
			return nil, true
		}
		if endOfDay, ok := voucherDateToEndOfDaySeconds(s); ok {
			sec = endOfDay
		} else {
			n, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, false
			}
			sec = n
		}
	} else {
		f, err := strconv.ParseFloat(s, 64)
		if err != nil || f != math.Trunc(f) {
			return nil, false
		}
		sec = int64(f)
	}
	if sec <= 0 {
		return nil, true // historical "absent" semantics
	}
	if sec < voucherExpiryMinUnix || sec > voucherExpiryMaxUnix {
		return nil, false
	}
	return &sec, true
}

// voucherDateToEndOfDaySeconds converts a strict YYYY-MM-DD date (interpreted
// in UTC, the store's clock) to 23:59:59 UTC of the same day.
func voucherDateToEndOfDaySeconds(s string) (int64, bool) {
	if len(s) != len("2006-01-02") {
		return 0, false
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return 0, false
	}
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 0, time.UTC).Unix(), true
}

// parseYuanToCents converts a CNY amount given in yuan (JSON number or quoted
// string, up to 2 decimal places, e.g. 12.5 or "12.50") to integer min units
// (分). Returns ok=false for absent, malformed, non-positive or >2-decimal
// input. Floats are normalized through their shortest decimal representation
// (strconv.FormatFloat -1) so binary drift can never mint a wrong cent value
// (E-008: voucher batch generation inputs amounts in yuan).
func parseYuanToCents(raw json.RawMessage) (int64, bool) {
	s := strings.TrimSpace(string(raw))
	if len(s) >= 2 && s[0] == '"' {
		var str string
		if err := json.Unmarshal(raw, &str); err != nil {
			return 0, false
		}
		s = strings.TrimSpace(str)
	}
	if s == "" {
		return 0, false
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || f <= 0 {
		return 0, false
	}
	// Shortest round-trip decimal: reject more than two fractional digits
	// instead of silently rounding (e.g. 12.345).
	rep := strconv.FormatFloat(f, 'f', -1, 64)
	whole := rep
	frac := ""
	if dot := strings.IndexByte(rep, '.'); dot >= 0 {
		whole, frac = rep[:dot], rep[dot+1:]
	}
	if strings.HasPrefix(whole, "-") || len(whole) > 15 || len(frac) > 2 {
		return 0, false
	}
	yuan, err := strconv.ParseInt(whole, 10, 64)
	if err != nil {
		return 0, false
	}
	cents := yuan * 100
	if len(frac) == 1 {
		cents += int64(frac[0]-'0') * 10
	} else if len(frac) == 2 {
		cents += int64(frac[0]-'0')*10 + int64(frac[1]-'0')
	}
	return cents, true
}

func writeVoucherRedeemError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, voucher.ErrNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "VOUCHER_NOT_FOUND", "voucher not found")
	case errors.Is(err, voucher.ErrVoucherAlreadyRedeemed), errors.Is(err, voucher.ErrVoucherConflict):
		writeLocalizedError(w, r, http.StatusConflict, "VOUCHER_ALREADY_REDEEMED", "voucher has already been redeemed")
	case errors.Is(err, voucher.ErrVoucherVoid):
		writeLocalizedError(w, r, http.StatusConflict, "VOUCHER_VOID", "voucher has been voided")
	case errors.Is(err, voucher.ErrVoucherExpired):
		writeLocalizedError(w, r, http.StatusConflict, "VOUCHER_EXPIRED", "voucher has expired")
	case errors.Is(err, voucher.ErrInvalidInput), errors.Is(err, voucher.ErrVoucherInvalid):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_BODY", "body must be JSON with code")
	case errors.Is(err, voucher.ErrCurrencyMismatch):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_PARAMS", "voucher currency mismatch")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "voucher redeem failed")
	}
}

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
