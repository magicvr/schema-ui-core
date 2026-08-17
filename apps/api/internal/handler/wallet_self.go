// GOAL-022 · 我的钱包自服务面 — identity-scoped, read-only self-service.
// GET /api/wallet/me and GET /api/wallet/me/entries derive the owner EXCLUSIVELY
// from the session identity (never from a client-supplied ownerId), so a user
// can only ever read their own wallet (D-002 §1/§2). Identity-only endpoints:
// no permission key — every authenticated user may view their own wallet, like
// the /api/account/* self-service surface. Lazy get-or-create reuses GOAL-020
// semantics; the first call auto-opens the zero-balance account (audited with
// the auto marker).
package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/account"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
)

// WalletSelfRoutes returns the self-service wallet surface (admin.wallet).
func WalletSelfRoutes(a *auth.Authenticator, service WalletService, operations operationlog.Recorder, moduleID string) []kernel.RouteContribution {
	var routes []kernel.RouteContribution
	add := func(method, pattern string, h http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              h,
		})
	}

	// Self summary: get-or-create the session user's account, returned in the
	// resourceList envelope so schema statCards can bind balance fields.
	add("GET", "/api/wallet/me", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate,
				`{"accountId":`+jsonQuote(account.ID)+`,"ownerId":`+jsonQuote(account.OwnerID)+`,"auto":true}`, now)
		}
		writeJSON(w, http.StatusOK, resourceList{Items: []map[string]any{accountToMap(*account)}, Total: 1, Page: 1, PageSize: 1})
	})))

	// Self entries: the session user's own ledger (paged), same lazy open.
	add("GET", "/api/wallet/me/entries", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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
			recordWalletEvent(operations, user, operationlog.EventWalletAccountCreate,
				`{"accountId":`+jsonQuote(account.ID)+`,"ownerId":`+jsonQuote(account.OwnerID)+`,"auto":true}`, now)
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

	return routes
}

// selfIdentity resolves the session user (401 when absent). The self-service
// surface is identity-only: no permission key, never a client-supplied owner.
func selfIdentity(w http.ResponseWriter, r *http.Request) (account.User, bool) {
	user, ok := auth.IdentityFrom(r.Context())
	if !ok {
		writeLocalizedError(w, r, http.StatusUnauthorized, "UNAUTHENTICATED", "no active session")
		return account.User{}, false
	}
	return user, true
}
