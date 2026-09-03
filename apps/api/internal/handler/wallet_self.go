// GOAL-022 · 我的钱包自服务面 — identity-scoped self-service.
// GET /api/wallet/me and GET /api/wallet/me/entries derive the owner EXCLUSIVELY
// from the session identity (never from a client-supplied ownerId). Identity-only
// endpoints: no permission key. Lazy get-or-create reuses GOAL-020 semantics.
// VP-029 R5 (GOAL-005): POST /api/wallet/me/redeem credits the session user's
// owner_type=user ledger; never Redeem(subjectID).
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
)

const (
	walletRedeemRateLimiterWindow   = 15 * time.Minute
	walletRedeemRateLimiterMax      = 10
	walletRedeemRateLimiterCapacity = 1 << 16
)

// WalletSelfRoutes returns the self-service wallet surface (admin.wallet).
func WalletSelfRoutes(a *auth.Authenticator, service WalletService, operations operationlog.Recorder, moduleID string, limiters kernel.RateLimiterProvider) []kernel.RouteContribution {
	redeemLimiter := limiters.NewRateLimiter(walletRedeemRateLimiterWindow, walletRedeemRateLimiterMax, walletRedeemRateLimiterCapacity)

	var routes []kernel.RouteContribution
	add := func(method, pattern string, h http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              h,
		})
	}

	// Self summary: read-only (W15-F11). Missing → 404.
	add("GET", "/api/wallet/me", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := selfIdentity(w, r)
		if !ok {
			return
		}
		account, err := service.GetUserAccountByOwner(user.ID)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		writeJSON(w, http.StatusOK, resourceList{Items: []map[string]any{accountToMap(*account)}, Total: 1, Page: 1, PageSize: 1})
	})))

	add("POST", "/api/wallet/me", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := selfIdentity(w, r)
		if !ok {
			return
		}
		now := time.Now().UTC()
		account, created, err := service.GetOrCreateUserAccount(user.ID, now)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		if created {
			recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate, "account-create", map[string]any{"accountId": account.ID, "ownerId": account.OwnerID, "auto": true}, now)
		}
		writeJSON(w, http.StatusOK, resourceList{Items: []map[string]any{accountToMap(*account)}, Total: 1, Page: 1, PageSize: 1})
	})))

	// Self entries: the session user's own ledger (paged), same lazy open.
	add("GET", "/api/wallet/me/entries", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := selfIdentity(w, r)
		if !ok {
			return
		}
		account, err := service.GetUserAccountByOwner(user.ID)
		if err != nil {
			writeWalletError(w, r, err)
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
		entries, total, err := service.ListEntries(account.ID, entryType, q, page, pageSize)
		if err != nil {
			writeWalletError(w, r, err)
			return
		}
		rows := make([]map[string]any, 0, len(entries))
		for _, e := range entries {
			rows = append(rows, entryToMap(e))
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
	})))

	add("POST", "/api/wallet/me/redeem", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := selfIdentity(w, r)
		if !ok {
			return
		}
		now := time.Now().UTC()
		if !redeemLimiter.AllowRecord(user.ID, now) {
			if sec := redeemLimiter.RetryAfterSeconds(user.ID, now); sec > 0 {
				w.Header().Set("Retry-After", strconv.Itoa(sec))
			}
			writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many voucher redeem attempts; try again later")
			return
		}
		var body struct {
			Code string `json:"code"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 4<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_BODY", "body must be JSON with code")
			return
		}
		code := strings.TrimSpace(body.Code)
		if code == "" {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_VOUCHER_BODY", "body must be JSON with code")
			return
		}
		result, err := service.RedeemForUser(r.Context(), user.ID, user.Name, code, now)
		if err != nil {
			writeVoucherRedeemError(w, r, err)
			return
		}
		redeemLimiter.Clear(user.ID)
		recordWalletEvent(operations, user, operationlog.EventWalletAdjust, "voucher-redeem", map[string]any{
			"voucherId":  result.VoucherID,
			"batchId":    result.BatchID,
			"codePrefix": result.CodePrefix,
			"amount":     result.Amount,
			"accountId":  result.AccountID,
			"entryId":    result.EntryID,
		}, now)
		writeJSON(w, http.StatusOK, map[string]any{
			"voucherId":    result.VoucherID,
			"batchId":      result.BatchID,
			"codePrefix":   result.CodePrefix,
			"amount":       result.Amount,
			"currency":     result.Currency,
			"accountId":    result.AccountID,
			"entryId":      result.EntryID,
			"balanceAfter": result.Balance,
		})
	})))

	return routes
}

// selfIdentity resolves the session user (401 when absent). The self-service
// surface is identity-only: no permission key, never a client-supplied owner.
func selfIdentity(w http.ResponseWriter, r *http.Request) (account.User, bool) {
	user, ok := auth.UserIdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	return user, true
}
