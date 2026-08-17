// Wallet/ledger HTTP surface (S-14 · GOAL-019 D-002 §3): account lifecycle,
// balance mutations (adjust/freeze/unfreeze — the money-path permissions are
// wallet.adjust), the immutable ledger and reconciliation runs. All mutations
// are audited via operationlog.
package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet/store"
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
	Mutate(id string, in walletstore.LedgerEntryInput, now time.Time) (*walletstore.Account, *walletstore.LedgerEntry, error)
	Reconcile(accountID, actorID string, now time.Time) (*walletstore.ReconciliationRun, error)
	ListReconcileRuns(page, pageSize int) ([]walletstore.ReconciliationRun, int, error)
}

// WalletRoutes returns the admin.wallet HTTP surface.
func WalletRoutes(a *auth.Authenticator, service WalletService, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
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
		pageSize, ok := intParam(r.URL.Query().Get("pageSize"), 20)
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


	// Accounts: get-or-create one user account by owner (GOAL-020 D-001 §1).
	// The auto-created zero-balance account is audited with an auto marker.
	add("GET", "/api/wallet/by-owner/{ownerId}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.read")
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
			detail := `{"accountId":` + jsonQuote(account.ID) + `,"ownerId":` + jsonQuote(account.OwnerID) + `,"auto":true}`
			recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate, detail, now)
		}
		writeJSON(w, http.StatusOK, accountToMap(*account))
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
			detail := `{"accountId":` + jsonQuote(account.ID) + `,"ownerId":` + jsonQuote(ownerID) + `,"auto":true}`
			recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate, detail, now)
		}
		account, entry, err := service.Mutate(account.ID, walletstore.LedgerEntryInput{
			EntryType: walletstore.EntryAdjust, AmountDelta: body.AmountDelta,
			RefType: strings.TrimSpace(body.RefType), RefID: strings.TrimSpace(body.RefID),
			IdempotencyKey: strings.TrimSpace(body.IdempotencyKey), Memo: strings.TrimSpace(body.Memo),
			ActorID: user.ID, ActorName: user.Name,
		}, now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		detail := `{"accountId":` + jsonQuote(account.ID) + `,"entryId":` + jsonQuote(entry.ID) + `,"amountDelta":` + strconv.FormatInt(entry.AmountDelta, 10) + `}`
		recordWalletEvent(operations, user, operationlog.EventWalletAdjust, detail, now)
		writeJSON(w, http.StatusOK, map[string]any{"account": accountToMap(*account), "entry": entryToMap(*entry)})
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
		recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate,
			`{"accountId":`+jsonQuote(account.ID)+`,"ownerType":`+jsonQuote(account.OwnerType)+`,"ownerId":`+jsonQuote(account.OwnerID)+`}`, now)
		writeJSON(w, http.StatusCreated, accountToMap(*account))
	})))

	// Accounts: status update (audited). Optimistic lock: the caller passes
	// the version observed when loading the row.
	add("PATCH", "/api/wallet/accounts/{id}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "wallet.write")
		if !ok {
			return
		}
		var body struct {
			Status  string `json:"status"`
			Version int64  `json:"version"`
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
		now := time.Now().UTC()
		account, err := service.UpdateStatus(r.PathValue("id"), status, body.Version, now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		recordWalletEvent(operations, user, operationlog.EventWalletAccountUpdate,
			`{"accountId":`+jsonQuote(account.ID)+`,"status":`+jsonQuote(account.Status)+`}`, now)
		writeJSON(w, http.StatusOK, accountToMap(*account))
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

	// Reconcile: ledger chain check (read path — does not change balances).
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
		now := time.Now().UTC()
		run, err := service.Reconcile(strings.TrimSpace(body.AccountID), user.ID, now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		detail := `{"runId":` + jsonQuote(run.ID) + `,"result":` + jsonQuote(run.Result) + `}`
		recordWalletEvent(operations, user, operationlog.EventWalletReconcile, detail, now)
		writeJSON(w, http.StatusOK, reconcileRunToMap(*run))
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
		pageSize, ok := intParam(r.URL.Query().Get("pageSize"), 20)
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
	pageSize, ok := intParam(r.URL.Query().Get("pageSize"), 20)
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
	account, entry, err := service.Mutate(r.PathValue("id"), walletstore.LedgerEntryInput{
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
	detail := `{"accountId":` + jsonQuote(account.ID) + `,"entryId":` + jsonQuote(entry.ID) + `,"amountDelta":` + strconv.FormatInt(entry.AmountDelta, 10) + `}`
	recordWalletEvent(operations, user, "wallet."+eventSuffix, detail, now)
	writeJSON(w, http.StatusOK, map[string]any{"account": accountToMap(*account), "entry": entryToMap(*entry)})
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

// recordWalletEvent writes a wallet audit row.
func recordWalletEvent(operations operationlog.Recorder, user account.User, event, detail string, now time.Time) {
	if operations == nil {
		return
	}
	_ = operations.RecordOperation(operationlog.Operation{
		ID: newOperationID(), Event: event,
		ActorID: user.ID, ActorName: user.Name, Detail: &detail, CreatedAt: now,
	})
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
	}
}

func entryToMap(e walletstore.LedgerEntry) map[string]any {
	return map[string]any{
		"id":                  e.ID,
		"accountId":           e.AccountID,
		"entryType":           e.EntryType,
		"amountDelta":         e.AmountDelta,
		"balanceAfterTotal":   e.BalanceAfterTotal,
		"balanceAfterAvail":   e.BalanceAfterAvail,
		"balanceAfterFrozen":  e.BalanceAfterFrozen,
		"refType":             e.RefType,
		"refId":               e.RefID,
		"memo":                e.Memo,
		"actorId":             e.ActorID,
		"actorName":           e.ActorName,
		"createdAt":           e.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
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