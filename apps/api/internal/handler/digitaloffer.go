// GOAL-003 (workspace-031 · VP-031 R2 · GOAL-002 D-002 v1.0.0 §5.3/§7) —
// biz.digital-offer admin surface plus the public on-sale catalog. Offer
// writes are fail-closed audited inside the service transaction (§7); the
// purchase money path has NO HTTP surface in this wave (§4.3 — module API and
// Telegram only); purchases list is a read-only audit view.
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
	"github.com/magicvr/schema-ui-core/apps/api/internal/concurrency"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/service"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/store"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
)

// Public catalog rate limit (D-002 §8: bizoffer|list|<ip>, 60/min — request
// counting via AllowRecord, NEVER key-wide Clear).
const (
	bizOfferListRateLimiterWindow   = time.Minute
	bizOfferListRateLimiterMax      = 60
	bizOfferListRateLimiterCapacity = 1 << 16
)

// DigitalOfferRoutes returns the admin surface (module biz.digital-offer).
func DigitalOfferRoutes(a *auth.Authenticator, svc *service.Service, moduleID string, limiters kernel.RateLimiterProvider) []kernel.RouteContribution {
	var routes []kernel.RouteContribution
	add := func(method, pattern string, h http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              h,
		})
	}

	add("GET", "/api/digitaloffer/offers", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "digitaloffer.read"); !ok {
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
		offers, total, err := svc.ListOffers(store.OfferFilter{
			Q: r.URL.Query().Get("q"), Status: strings.TrimSpace(r.URL.Query().Get("status")),
			Page: page, PageSize: pageSize,
		})
		if err != nil {
			writeDigitalOfferError(w, r, err)
			return
		}
		rows := make([]map[string]any, 0, len(offers))
		for _, o := range offers {
			rows = append(rows, digitalOfferToMap(o))
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
	})))

	add("POST", "/api/digitaloffer/offers", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "digitaloffer.offer.manage")
		if !ok {
			return
		}
		var body struct {
			Name             string `json:"name"`
			Description      string `json:"description"`
			PriceAmount      int64  `json:"priceAmount"`
			Currency         string `json:"currency"`
			EntitlementForm  string `json:"entitlementForm"`
			DurationSeconds  int64  `json:"durationSeconds"`
			CountPerPurchase int64  `json:"countPerPurchase"`
			OnSale           bool   `json:"onSale"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BIZOFFER_REQUEST", "body must be JSON offer fields")
			return
		}
		offer, err := svc.CreateOffer(r.Context(), digitalOfferActor(user, r), service.CreateOfferInput{
			Name: body.Name, Description: body.Description, PriceAmount: body.PriceAmount,
			Currency: body.Currency, EntitlementForm: body.EntitlementForm,
			DurationSeconds: body.DurationSeconds, CountPerPurchase: body.CountPerPurchase,
			OnSale: body.OnSale,
		}, time.Now().UTC())
		if err != nil {
			writeDigitalOfferError(w, r, err)
			return
		}
		w.Header().Set("ETag", concurrency.ETag(offer.Version))
		writeJSON(w, http.StatusCreated, digitalOfferToMap(*offer))
	})))

	add("PATCH", "/api/digitaloffer/offers/{id}", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "digitaloffer.offer.manage")
		if !ok {
			return
		}
		var body struct {
			Name        *string `json:"name"`
			Description *string `json:"description"`
			PriceAmount *int64  `json:"priceAmount"`
			Status      *string `json:"status"`
			// Immutable per D-002 §2; any attempt is rejected with
			// BIZOFFER_FORM_CONFLICT rather than silently ignored.
			Currency         *string `json:"currency"`
			EntitlementForm  *string `json:"entitlementForm"`
			DurationSeconds  *int64  `json:"durationSeconds"`
			CountPerPurchase *int64  `json:"countPerPurchase"`
			ExpectedVersion  *int64  `json:"expectedVersion"`
			Version          *int64  `json:"version"`
		}
		r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BIZOFFER_REQUEST", "body must be JSON offer fields")
			return
		}
		if body.Currency != nil || body.EntitlementForm != nil || body.DurationSeconds != nil || body.CountPerPurchase != nil {
			writeLocalizedError(w, r, http.StatusConflict, "BIZOFFER_FORM_CONFLICT", "entitlement form fields are immutable")
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
		offer, err := svc.UpdateOffer(r.Context(), digitalOfferActor(user, r), r.PathValue("id"), service.UpdateOfferInput{
			Name: body.Name, Description: body.Description, PriceAmount: body.PriceAmount, Status: body.Status,
		}, expectedVersion, time.Now().UTC())
		if err != nil {
			writeDigitalOfferError(w, r, err)
			return
		}
		w.Header().Set("ETag", concurrency.ETag(offer.Version))
		writeJSON(w, http.StatusOK, digitalOfferToMap(*offer))
	})))

	add("GET", "/api/digitaloffer/purchases", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "digitaloffer.read"); !ok {
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
		purchases, total, err := svc.ListPurchases(store.PurchaseFilter{
			Q: r.URL.Query().Get("q"), SubjectID: strings.TrimSpace(r.URL.Query().Get("subjectId")),
			OfferID: strings.TrimSpace(r.URL.Query().Get("offerId")), Page: page, PageSize: pageSize,
		})
		if err != nil {
			writeDigitalOfferError(w, r, err)
			return
		}
		rows := make([]map[string]any, 0, len(purchases))
		for _, p := range purchases {
			rows = append(rows, digitalPurchaseToMap(p))
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
	})))

	add("GET", "/api/digitaloffer/entitlements", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := requirePermission(w, r, "digitaloffer.read"); !ok {
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
		entitlements, total, err := svc.ListEntitlements(store.EntitlementFilter{
			Q: r.URL.Query().Get("q"), SubjectID: strings.TrimSpace(r.URL.Query().Get("subjectId")),
			OfferID: strings.TrimSpace(r.URL.Query().Get("offerId")), Status: strings.TrimSpace(r.URL.Query().Get("status")),
			Page: page, PageSize: pageSize,
		})
		if err != nil {
			writeDigitalOfferError(w, r, err)
			return
		}
		rows := make([]map[string]any, 0, len(entitlements))
		for _, e := range entitlements {
			rows = append(rows, digitalEntitlementToMap(e))
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: total, Page: page, PageSize: pageSize})
	})))

	add("POST", "/api/digitaloffer/entitlements/{id}/void", a.Middleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := requirePermission(w, r, "digitaloffer.entitlement.void")
		if !ok {
			return
		}
		result, err := svc.VoidEntitlement(r.Context(), digitalOfferActor(user, r), r.PathValue("id"), time.Now().UTC())
		if err != nil {
			writeDigitalOfferError(w, r, err)
			return
		}
		body := digitalEntitlementToMap(*result.Entitlement)
		body["alreadyVoided"] = result.AlreadyVoided
		writeJSON(w, http.StatusOK, body)
	})))

	return routes
}

// DigitalOfferPublicRoutes returns the C-end catalog (D-002 §5.3): public
// read of on_sale offers, IP-bucketed request counting.
func DigitalOfferPublicRoutes(a *auth.Authenticator, svc *service.Service, moduleID string, limiters kernel.RateLimiterProvider) []kernel.RouteContribution {
	listLimiter := limiters.NewRateLimiter(bizOfferListRateLimiterWindow, bizOfferListRateLimiterMax, bizOfferListRateLimiterCapacity)
	var routes []kernel.RouteContribution
	add := func(method, pattern string, h http.Handler) {
		routes = append(routes, kernel.RouteContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: moduleID, Key: kernel.RouteKey(method, pattern)},
			Method:               method,
			Pattern:              pattern,
			Handler:              h,
		})
	}
	// Public catalog (D-002 §5.3): unauthenticated read, IP-bucketed.
	add("GET", "/api/biz/offers", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now().UTC()
		key := "bizoffer|list|" + clientIP(r)
		if !listLimiter.AllowRecord(key, now) {
			if sec := listLimiter.RetryAfterSeconds(key, now); sec > 0 {
				w.Header().Set("Retry-After", strconv.Itoa(sec))
			}
			writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many catalog requests; try again later")
			return
		}
		offers, err := svc.ListOnSale()
		if err != nil {
			writeDigitalOfferError(w, r, err)
			return
		}
		rows := make([]map[string]any, 0, len(offers))
		for _, o := range offers {
			rows = append(rows, digitalOfferPublicToMap(o))
		}
		writeJSON(w, http.StatusOK, resourceList{Items: rows, Total: len(rows), Page: 1, PageSize: len(rows)})
	}))
	return routes
}

func digitalOfferActor(user account.User, r *http.Request) service.Actor {
	return service.Actor{
		ID:            user.ID,
		Name:          user.Name,
		SessionID:     identitySession(user),
		CorrelationID: requestid.FromContext(r.Context()),
	}
}

func digitalOfferToMap(o store.Offer) map[string]any {
	out := map[string]any{
		"id":              o.ID,
		"name":            o.Name,
		"description":     o.Description,
		"priceAmount":     o.PriceAmount,
		"currency":        o.Currency,
		"entitlementForm": o.EntitlementForm,
		"status":          o.Status,
		"version":         o.Version,
		"createdAt":       o.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt":       o.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
	if o.EntitlementForm == store.FormDuration {
		out["durationSeconds"] = o.DurationSeconds
	}
	if o.EntitlementForm == store.FormCount {
		out["countPerPurchase"] = o.CountPerPurchase
	}
	return out
}

// digitalOfferPublicToMap renders the C-end catalog projection (D-002 §5.3):
// no version or audit-adjacent internals.
func digitalOfferPublicToMap(o store.Offer) map[string]any {
	out := map[string]any{
		"id":              o.ID,
		"name":            o.Name,
		"description":     o.Description,
		"priceAmount":     o.PriceAmount,
		"currency":        o.Currency,
		"entitlementForm": o.EntitlementForm,
	}
	if o.EntitlementForm == store.FormDuration {
		out["durationSeconds"] = o.DurationSeconds
	}
	if o.EntitlementForm == store.FormCount {
		out["countPerPurchase"] = o.CountPerPurchase
	}
	return out
}

func digitalPurchaseToMap(p store.Purchase) map[string]any {
	return map[string]any{
		"id":            p.ID,
		"subjectId":     p.SubjectID,
		"offerId":       p.OfferID,
		"offerName":     p.OfferName,
		"amount":        p.Amount,
		"currency":      p.Currency,
		"freezeEntryId": p.FreezeEntryID,
		"deductEntryId": p.DeductEntryID,
		"requestId":     p.RequestID,
		"status":        p.Status,
		"createdAt":     p.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
}

func digitalEntitlementToMap(e store.Entitlement) map[string]any {
	out := map[string]any{
		"id":         e.ID,
		"subjectId":  e.SubjectID,
		"offerId":    e.OfferID,
		"purchaseId": e.PurchaseID,
		"form":       e.Form,
		"status":     e.Status,
		"createdAt":  e.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
		"updatedAt":  e.UpdatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
	}
	if e.ExpiresAt != nil {
		out["expiresAt"] = e.ExpiresAt.UTC().Format("2006-01-02T15:04:05.000Z07:00")
	}
	if e.RemainingCount != nil {
		out["remainingCount"] = *e.RemainingCount
	}
	return out
}

// writeDigitalOfferError maps domain sentinels to the frozen D-002 §9 codes.
func writeDigitalOfferError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, store.ErrNotFound),
		errors.Is(err, store.ErrPurchaseNotFound),
		errors.Is(err, store.ErrEntitlementNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "BIZOFFER_NOT_FOUND", "digital offer not found")
	case errors.Is(err, store.ErrOfferNotOnSale):
		writeLocalizedError(w, r, http.StatusConflict, "BIZOFFER_NOT_ON_SALE", "digital offer is not on sale")
	case errors.Is(err, store.ErrCurrencyMismatch):
		writeLocalizedError(w, r, http.StatusConflict, "BIZOFFER_CURRENCY_MISMATCH", "offer currency does not match the wallet account")
	case errors.Is(err, walletstore.ErrInsufficient):
		writeLocalizedError(w, r, http.StatusConflict, "BIZOFFER_INSUFFICIENT_FUNDS", "insufficient balance for this purchase")
	case errors.Is(err, store.ErrSubjectNotFound):
		writeLocalizedError(w, r, http.StatusNotFound, "BIZOFFER_SUBJECT_NOT_FOUND", "subject not found")
	case errors.Is(err, store.ErrRequestIdConflict):
		writeLocalizedError(w, r, http.StatusConflict, "BIZOFFER_REQUEST_CONFLICT", "request id was already used with a different offer")
	case errors.Is(err, store.ErrEntitlementInsufficient):
		writeLocalizedError(w, r, http.StatusConflict, "BIZOFFER_ENTITLEMENT_INSUFFICIENT", "insufficient entitlement count")
	case errors.Is(err, store.ErrFormConflict):
		writeLocalizedError(w, r, http.StatusConflict, "BIZOFFER_FORM_CONFLICT", "entitlement form fields are immutable")
	case errors.Is(err, store.ErrVersionConflict):
		writeLocalizedError(w, r, http.StatusConflict, "BIZOFFER_VERSION_CONFLICT", "the offer changed concurrently; reload and retry")
	case errors.Is(err, store.ErrInvalidOffer):
		writeLocalizedError(w, r, http.StatusBadRequest, "INVALID_BIZOFFER_REQUEST", "invalid digital offer request")
	case errors.Is(err, service.ErrRateLimited):
		writeLocalizedError(w, r, http.StatusTooManyRequests, "RATE_LIMITED", "too many purchase attempts; try again later")
	default:
		writeLocalizedError(w, r, http.StatusInternalServerError, "INTERNAL", "digital offer operation failed")
	}
}
