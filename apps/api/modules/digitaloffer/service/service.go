// Package service provides the biz.digital-offer domain surface (VP-031 ·
// GOAL-002 D-002 v1.0.0): sellable digital offers, thin purchase vouchers over
// the wallet freeze→deduct_frozen primitives and per-subject entitlements.
// The implementation is the R2 denominator of workspace-031 GOAL-003; every
// deviation from the contract requires a contract revision first. The package
// sits below the module root so the handler layer can consume it without an
// import cycle (wallet/voucher precedent).
package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/digitaloffer/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
)

// Operation-log event names frozen by D-002 §7. Every admin domain write
// pairs with its audit row inside the SAME caller-owned transaction.
const (
	EventOfferCreate     = "bizoffer.offer.create"
	EventOfferUpdate     = "bizoffer.offer.update"
	EventOfferStatus     = "bizoffer.offer.status"
	EventEntitlementVoid = "bizoffer.entitlement.void"
)

// purchaseRefType is the wallet ledger back-reference (D-002 §4.4): both the
// freeze and deduct_frozen entries point at the purchase voucher.
const purchaseRefType = "biz_offer_purchase"

// maxPurchaseAttempts is the D-002 §4.4 bounded retry budget. The loop lives
// OUTSIDE store.Run; each attempt is one fresh transaction.
const maxPurchaseAttempts = 3

// ErrRateLimited maps to the frozen generic RATE_LIMITED code (D-002 §8):
// the purchase request-count bucket denied this request.
var ErrRateLimited = errors.New("digital-offer purchase rate limited")

// Purchase bucket constants (D-002 §8): bizoffer|purchase|<subject_id>,
// 1 minute / 10 requests — request counting via AllowRecord, never Clear.
const (
	purchaseLimiterWindow   = time.Minute
	purchaseLimiterMax      = 10
	purchaseLimiterCapacity = 1 << 16
)

// Actor carries the admin-session identity for fail-closed audit rows.
type Actor struct {
	ID            string
	Name          string
	SessionID     string
	CorrelationID string
}

// Service implements the digital-offer domain surface.
type Service struct {
	repo       *store.Repository
	wallet     *walletstore.Repository
	subjects   *subject.Store
	operations operationlog.TransactionalRecorder
	limiter    kernel.RateLimiter
	now        func() time.Time
	// faultHook is a test-only injection point (D-002 §4.4 acceptance):
	// when set, it is invoked after each named purchase step inside the
	// transaction and a non-nil error aborts the attempt as retryable.
	faultHook func(step string) error
}

// NewService constructs the domain service. store and wallet must share the
// same platform runner so purchase attempts span both domains in ONE
// kernel.Tx (D-002 §4.2). limiters may be nil only in tests that never call
// Purchase — the §8 purchase bucket is mandatory on the production surface.
func NewService(repo *store.Repository, wallet *walletstore.Repository, subjects *subject.Store, operations operationlog.TransactionalRecorder, limiters kernel.RateLimiterProvider) *Service {
	var limiter kernel.RateLimiter
	if limiters != nil {
		limiter = limiters.NewRateLimiter(purchaseLimiterWindow, purchaseLimiterMax, purchaseLimiterCapacity)
	}
	return &Service{repo: repo, wallet: wallet, subjects: subjects, operations: operations, limiter: limiter, now: time.Now}
}

// SetPurchaseFaultHookForTest installs the §4.4 acceptance failure-injection
// seam. Steps: "freeze", "deduct", "purchase", "entitlement".
func (s *Service) SetPurchaseFaultHookForTest(hook func(step string) error) {
	s.faultHook = hook
}

// Offer returns the store package for the handler layer's read needs.
func (s *Service) Repo() *store.Repository { return s.repo }

var idSeq atomic.Uint64

// newID returns a time-ordered hex id (wallet precedent: millisecond prefix +
// per-process monotonic counter + random suffix) so (created_at, id) order
// matches creation order.
func newID(now time.Time) (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	seq := idSeq.Add(1) & 0xFFFFFFFF
	return fmt.Sprintf("%016x%08x%s", now.UnixMilli(), seq, hex.EncodeToString(buf)), nil
}

// CreateOfferInput is the admin create payload (D-002 §2). Form parameters
// are frozen at create time and immutable afterwards.
type CreateOfferInput struct {
	Name             string
	Description      string
	PriceAmount      int64
	Currency         string
	EntitlementForm  string
	DurationSeconds  int64
	CountPerPurchase int64
	OnSale           bool
}

func (in CreateOfferInput) validate() error {
	if strings.TrimSpace(in.Name) == "" {
		return store.ErrInvalidOffer
	}
	if in.PriceAmount <= 0 {
		return store.ErrInvalidOffer
	}
	if strings.TrimSpace(in.Currency) == "" {
		return store.ErrInvalidOffer
	}
	switch in.EntitlementForm {
	case store.FormDuration:
		if in.DurationSeconds < 1 || in.CountPerPurchase != 0 {
			return store.ErrInvalidOffer
		}
	case store.FormCount:
		if in.CountPerPurchase < 1 || in.DurationSeconds != 0 {
			return store.ErrInvalidOffer
		}
	default:
		return store.ErrInvalidOffer
	}
	return nil
}

// CreateOffer inserts one offer with a fail-closed audit row (D-002 §7).
func (s *Service) CreateOffer(ctx context.Context, actor Actor, in CreateOfferInput, now time.Time) (*store.Offer, error) {
	if err := in.validate(); err != nil {
		return nil, err
	}
	id, err := newID(now)
	if err != nil {
		return nil, err
	}
	status := store.StatusDraft
	if in.OnSale {
		status = store.StatusOnSale
	}
	offer := store.Offer{
		ID: id, Name: strings.TrimSpace(in.Name), Description: in.Description,
		PriceAmount: in.PriceAmount, Currency: strings.TrimSpace(in.Currency),
		EntitlementForm: in.EntitlementForm, DurationSeconds: in.DurationSeconds,
		CountPerPurchase: in.CountPerPurchase, Status: status,
		Version: 0, CreatedAt: now, UpdatedAt: now,
	}
	detail := auditDetail("offer-create", map[string]any{
		"offerId": id, "name": offer.Name, "price": offer.PriceAmount,
		"currency": offer.Currency, "form": offer.EntitlementForm, "status": status,
	})
	err = s.runnerRun(ctx, func(tx kernel.Tx) error {
		if err := s.repo.InsertOfferInTx(tx, offer); err != nil {
			return err
		}
		return s.recordOperationTx(tx, actor, EventOfferCreate, id, detail, now)
	})
	if err != nil {
		return nil, err
	}
	return &offer, nil
}

// UpdateOfferInput carries the mutable fields (D-002 §2). Nil = unchanged.
// Immutable form fields are NOT addressable; the handler maps any attempt to
// ErrFormConflict before reaching this point.
type UpdateOfferInput struct {
	Name        *string
	Description *string
	PriceAmount *int64
	Status      *string
}

// UpdateOffer applies mutable fields with the version guard and a fail-closed
// audit row. Status changes are audited as bizoffer.offer.status.
func (s *Service) UpdateOffer(ctx context.Context, actor Actor, id string, in UpdateOfferInput, expectedVersion int64, now time.Time) (*store.Offer, error) {
	if in.Name == nil && in.Description == nil && in.PriceAmount == nil && in.Status == nil {
		return nil, store.ErrInvalidOffer
	}
	if in.Name != nil && strings.TrimSpace(*in.Name) == "" {
		return nil, store.ErrInvalidOffer
	}
	if in.PriceAmount != nil && *in.PriceAmount <= 0 {
		return nil, store.ErrInvalidOffer
	}
	if in.Status != nil {
		switch *in.Status {
		case store.StatusDraft, store.StatusOnSale, store.StatusOffSale:
		default:
			return nil, store.ErrInvalidOffer
		}
	}
	var updated *store.Offer
	err := s.runnerRun(ctx, func(tx kernel.Tx) error {
		offer, err := s.repo.UpdateOfferInTx(tx, id, in.Name, in.Description, in.PriceAmount, in.Status, expectedVersion, now)
		if err != nil {
			return err
		}
		updated = offer
		event := EventOfferUpdate
		fields := map[string]any{"offerId": id}
		if in.Name != nil {
			fields["name"] = *in.Name
		}
		if in.Description != nil {
			fields["descriptionChanged"] = true
		}
		if in.PriceAmount != nil {
			fields["price"] = *in.PriceAmount
		}
		if in.Status != nil {
			event = EventOfferStatus
			fields["status"] = *in.Status
		}
		return s.recordOperationTx(tx, actor, event, id, auditDetail("offer-update", fields), now)
	})
	if err != nil {
		return nil, err
	}
	return updated, nil
}

// ListOffers / GetOffer / ListOnSale / ListPurchases / ListEntitlements are
// read paths delegating to the repository.

func (s *Service) ListOffers(filter store.OfferFilter) ([]store.Offer, int, error) {
	return s.repo.ListOffers(filter)
}

func (s *Service) GetOffer(id string) (*store.Offer, error) { return s.repo.GetOffer(id) }

func (s *Service) ListOnSale() ([]store.Offer, error) { return s.repo.ListOnSale() }

func (s *Service) ListPurchases(filter store.PurchaseFilter) ([]store.Purchase, int, error) {
	return s.repo.ListPurchases(filter)
}

func (s *Service) ListEntitlements(filter store.EntitlementFilter) ([]store.Entitlement, int, error) {
	return s.repo.ListEntitlements(filter)
}

func (s *Service) GetEntitlement(id string) (*store.Entitlement, error) {
	return s.repo.GetEntitlement(id)
}

// PurchaseResult reports one fulfilled purchase (fresh or idempotent replay).
type PurchaseResult struct {
	Purchase    *store.Purchase
	Entitlement *store.Entitlement
	Replayed    bool
}

// Purchase executes the D-002 §4.2 single-transaction purchase with the
// §4.4 bounded retry protocol: the attempt loop lives OUTSIDE store.Run,
// every attempt starts with the (subject_id, request_id) read-back, terminal
// errors are never retried, and everything else retries with a fresh
// transaction (convergence: the loser reads back the committed winner).
func (s *Service) Purchase(ctx context.Context, subjectID, offerID, requestID string, now time.Time) (*PurchaseResult, error) {
	subjectID = strings.TrimSpace(subjectID)
	offerID = strings.TrimSpace(offerID)
	requestID = strings.TrimSpace(requestID)
	if subjectID == "" || offerID == "" || requestID == "" {
		return nil, store.ErrInvalidOffer
	}
	// §8 purchase bucket: request counting per call — AllowRecord, and the
	// key-wide Clear is never invoked (V-F119).
	if s.limiter != nil {
		if !s.limiter.AllowRecord("bizoffer|purchase|"+subjectID, now) {
			return nil, ErrRateLimited
		}
	}
	purchaseID, err := newID(now)
	if err != nil {
		return nil, err
	}
	entryFreezeID, err := newID(now)
	if err != nil {
		return nil, err
	}
	entryDeductID, err := newID(now)
	if err != nil {
		return nil, err
	}
	entitlementID, err := newID(now)
	if err != nil {
		return nil, err
	}
	var lastErr error
	for attempt := 0; attempt < maxPurchaseAttempts; attempt++ {
		result, err := s.purchaseAttempt(ctx, subjectID, offerID, requestID, purchaseID, entryFreezeID, entryDeductID, entitlementID, now)
		if err == nil {
			return result, nil
		}
		if isTerminalPurchaseError(err) {
			return nil, err
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = store.ErrInvalidOffer
	}
	return nil, fmt.Errorf("purchase attempts exhausted: %w", lastErr)
}

// isTerminalPurchaseError classifies D-002 §4.4 terminal errors. Everything
// else (wallet version conflicts, unique races — including the wallet-internal
// unexported sentinel — and commit failures) is retryable.
func isTerminalPurchaseError(err error) bool {
	switch {
	case errors.Is(err, store.ErrInvalidOffer),
		errors.Is(err, store.ErrNotFound),
		errors.Is(err, store.ErrOfferNotOnSale),
		errors.Is(err, store.ErrSubjectNotFound),
		errors.Is(err, store.ErrCurrencyMismatch),
		errors.Is(err, store.ErrRequestIdConflict),
		errors.Is(err, walletstore.ErrInsufficient),
		errors.Is(err, walletstore.ErrInvalidEntry),
		errors.Is(err, walletstore.ErrIdempotencyConflict),
		errors.Is(err, walletstore.ErrDisabled):
		return true
	}
	return false
}

func (s *Service) purchaseAttempt(ctx context.Context, subjectID, offerID, requestID, purchaseID, entryFreezeID, entryDeductID, entitlementID string, now time.Time) (*PurchaseResult, error) {
	var (
		result *PurchaseResult
	)
	err := s.runnerRun(ctx, func(tx kernel.Tx) error {
		// Attempt step 1: idempotency read-back (D-002 §4.4).
		existing, err := s.repo.GetPurchaseByRequestInTx(tx, subjectID, requestID)
		if err == nil {
			if existing.OfferID != offerID {
				return store.ErrRequestIdConflict
			}
			ent, entErr := s.repo.GetEntitlementByPurchaseInTx(tx, existing.ID)
			if entErr != nil {
				return entErr
			}
			result = &PurchaseResult{Purchase: existing, Entitlement: ent, Replayed: true}
			return nil
		}
		if !errors.Is(err, store.ErrPurchaseNotFound) {
			return err
		}

		// Attempt step 2: offer / subject / account validation.
		offer, err := s.repo.GetOfferInTx(tx, offerID)
		if err != nil {
			return err
		}
		if offer.Status != store.StatusOnSale {
			return store.ErrOfferNotOnSale
		}
		exists, err := s.subjects.SubjectExistsInTx(ctx, tx, subjectID)
		if err != nil {
			return err
		}
		if !exists {
			return store.ErrSubjectNotFound
		}
		account, _, err := s.wallet.GetOrCreateSubjectAccountInTx(tx, subjectID, now)
		if err != nil {
			return err
		}
		if account.Currency != offer.Currency {
			return store.ErrCurrencyMismatch
		}

		// Attempt step 3 (D-002 v1.2.0 §4.2 reorder): the purchase voucher is
		// inserted FIRST so the UNIQUE(subject_id, request_id) guard resolves
		// concurrent replays BEFORE any wallet mutation — a loser never
		// reaches the ledger with its divergent pre-generated purchase id.
		purchase := &store.Purchase{
			ID: purchaseID, SubjectID: subjectID, OfferID: offer.ID, OfferName: offer.Name,
			Amount: offer.PriceAmount, Currency: offer.Currency,
			FreezeEntryID: entryFreezeID, DeductEntryID: entryDeductID,
			RequestID: requestID, Status: "fulfilled", CreatedAt: now,
		}
		if err := s.repo.InsertPurchaseInTx(tx, *purchase); err != nil {
			return err
		}
		if err := s.fault("purchase"); err != nil {
			return err
		}

		// Attempt wallet steps: freeze → deduct_frozen inside the SAME tx.
		freezeIn := walletstore.LedgerEntryInput{
			EntryType:      walletstore.EntryFreeze,
			AmountDelta:    offer.PriceAmount,
			RefType:        purchaseRefType,
			RefID:          purchaseID,
			IdempotencyKey: requestID + ":freeze",
			Memo:           fmt.Sprintf("digital-offer purchase %s", purchaseID),
			ActorID:        subjectID,
			ActorName:      "Subject",
		}
		if _, _, err := s.wallet.MutateInTx(tx, account.ID, freezeIn, entryFreezeID, now); err != nil {
			return err
		}
		if err := s.fault("freeze"); err != nil {
			return err
		}
		deductIn := walletstore.LedgerEntryInput{
			EntryType:      walletstore.EntryDeductFrozen,
			AmountDelta:    offer.PriceAmount,
			RefType:        purchaseRefType,
			RefID:          purchaseID,
			IdempotencyKey: requestID + ":deduct",
			Memo:           fmt.Sprintf("digital-offer purchase %s", purchaseID),
			ActorID:        subjectID,
			ActorName:      "Subject",
		}
		if _, _, err := s.wallet.MutateInTx(tx, account.ID, deductIn, entryDeductID, now); err != nil {
			return err
		}
		if err := s.fault("deduct"); err != nil {
			return err
		}

		// Attempt final step: the entitlement grant.
		ent := &store.Entitlement{
			ID: entitlementID, SubjectID: subjectID, OfferID: offer.ID,
			PurchaseID: purchase.ID, Form: offer.EntitlementForm,
			Status: store.EntitlementActive, CreatedAt: now, UpdatedAt: now,
		}
		switch offer.EntitlementForm {
		case store.FormDuration:
			exp := now.Add(time.Duration(offer.DurationSeconds) * time.Second)
			ent.ExpiresAt = &exp
		case store.FormCount:
			count := offer.CountPerPurchase
			ent.RemainingCount = &count
		}
		if err := s.repo.InsertEntitlementInTx(tx, *ent); err != nil {
			return err
		}
		if err := s.fault("entitlement"); err != nil {
			return err
		}
		result = &PurchaseResult{Purchase: purchase, Entitlement: ent}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// VoidEntitlementResult reports the void outcome (idempotent: a repeated void
// is a success that appends no audit row, D-002 §7).
type VoidEntitlementResult struct {
	Entitlement   *store.Entitlement
	AlreadyVoided bool
}

// VoidEntitlement transitions active → voided with a fail-closed audit row.
func (s *Service) VoidEntitlement(ctx context.Context, actor Actor, id string, now time.Time) (*VoidEntitlementResult, error) {
	var result *VoidEntitlementResult
	err := s.runnerRun(ctx, func(tx kernel.Tx) error {
		alreadyVoided, err := s.repo.VoidEntitlementInTx(tx, id, now)
		if err != nil {
			return err
		}
		ent, err := s.repo.GetEntitlementInTx(tx, id)
		if err != nil {
			return err
		}
		if !alreadyVoided {
			if err := s.recordOperationTx(tx, actor, EventEntitlementVoid, id, auditDetail("entitlement-void", map[string]any{
				"entitlementId": id, "subjectId": ent.SubjectID, "offerId": ent.OfferID,
			}), now); err != nil {
				return err
			}
		}
		result = &VoidEntitlementResult{Entitlement: ent, AlreadyVoided: alreadyVoided}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// runnerRun opens one caller-owned transaction (D-002 §4.4: the retry loop
// lives here, NOT inside the callback).
func (s *Service) runnerRun(ctx context.Context, fn func(tx kernel.Tx) error) error {
	return s.repo.Runner().Run(ctx, fn)
}

// fault invokes the test-only failure injection for one purchase step; nil
// hook (production) is always a no-op.
func (s *Service) fault(step string) error {
	if s.faultHook == nil {
		return nil
	}
	return s.faultHook(step)
}

// recordOperationTx pairs the audit row with the domain write inside the
// caller-owned transaction — audit failure rolls the domain write back
// (D-002 §7 fail-closed).
func (s *Service) recordOperationTx(tx kernel.Tx, actor Actor, event, recordID string, detail *string, now time.Time) error {
	if s.operations == nil {
		return errors.New("digitaloffer: operation recorder not wired")
	}
	op := operationlog.Operation{
		ID:        newOperationID(),
		Event:     event,
		ActorID:   actor.ID,
		ActorName: actor.Name,
		Detail:    detail,
		SessionID: actor.SessionID,
		CreatedAt: now,
	}
	if recordID != "" {
		op.RecordID = &recordID
	}
	if actor.CorrelationID != "" {
		op.CorrelationID = actor.CorrelationID
	}
	return s.operations.RecordOperationTx(tx, op)
}

var opSeq atomic.Uint64

func newOperationID() string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%016x%08x", time.Now().UnixMilli(), opSeq.Add(1))
	}
	return fmt.Sprintf("%016x%s", time.Now().UnixMilli(), hex.EncodeToString(buf))
}

func auditDetail(action string, fields map[string]any) *string {
	payload := map[string]any{"action": action}
	for k, v := range fields {
		payload[k] = v
	}
	raw, err := json.Marshal(payload)
	if err != nil {
		return nil
	}
	s := string(raw)
	return &s
}
