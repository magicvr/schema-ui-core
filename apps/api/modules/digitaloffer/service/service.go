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
	repo           *store.Repository
	wallet         *walletstore.Repository
	subjects       *subject.Store
	operations     operationlog.TransactionalRecorder
	limiter        kernel.RateLimiter
	queryLimiter   kernel.RateLimiter
	telegramSender kernel.TelegramSender
	now            func() time.Time
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
	var limiter, queryLimiter kernel.RateLimiter
	if limiters != nil {
		limiter = limiters.NewRateLimiter(purchaseLimiterWindow, purchaseLimiterMax, purchaseLimiterCapacity)
		queryLimiter = limiters.NewRateLimiter(queryLimiterWindow, queryLimiterMax, queryLimiterCapacity)
	}
	return &Service{repo: repo, wallet: wallet, subjects: subjects, operations: operations, limiter: limiter, queryLimiter: queryLimiter, now: time.Now}
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

// Check reasons (D-002 §5.1): expired/exhausted are derived predicates; the
// aggregate is deterministic so every entry point reports the same state.
const (
	ReasonNone          = "none"
	ReasonNoEntitlement = "no_entitlement"
	ReasonExpired       = "expired"
	ReasonExhausted     = "exhausted"
	ReasonVoided        = "voided"
)

// CheckResult reports the unified entitlement validation outcome (§5.1).
type CheckResult struct {
	Valid  bool
	Reason string
}

// Check validates whether subjectID holds a currently valid entitlement for
// offerID (D-002 §5.1). Channels and services MUST call this before providing
// the capability. Aggregate priority: valid → no_entitlement → expired →
// exhausted → voided.
func (s *Service) Check(ctx context.Context, subjectID, offerID string, now time.Time) (CheckResult, error) {
	subjectID = strings.TrimSpace(subjectID)
	offerID = strings.TrimSpace(offerID)
	if subjectID == "" || offerID == "" {
		return CheckResult{}, store.ErrInvalidOffer
	}
	var rows []store.Entitlement
	err := s.runnerRun(ctx, func(tx kernel.Tx) error {
		var qErr error
		rows, qErr = s.repo.ListEntitlementsBySubjectOfferInTx(tx, subjectID, offerID)
		return qErr
	})
	if err != nil {
		return CheckResult{}, err
	}
	if len(rows) == 0 {
		return CheckResult{Valid: false, Reason: ReasonNoEntitlement}, nil
	}
	valid := false
	everValid := false
	anyVoided := false
	anyCount := false
	hasExpired := false
	hasExhausted := false
	for _, r := range rows {
		if r.Form == store.FormDuration {
			everValid = true
			if r.Status == store.EntitlementActive && r.ExpiresAt != nil && r.ExpiresAt.After(now) {
				valid = true
			}
			if r.Status == store.EntitlementActive && r.ExpiresAt != nil && !r.ExpiresAt.After(now) {
				hasExpired = true
			}
		} else {
			anyCount = true
			if r.Status == store.EntitlementActive && r.RemainingCount != nil && *r.RemainingCount > 0 {
				valid = true
			}
			if r.Status == store.EntitlementActive && r.RemainingCount != nil && *r.RemainingCount == 0 {
				hasExhausted = true
			}
		}
		if r.Status == store.EntitlementVoided {
			anyVoided = true
		}
	}
	if valid {
		return CheckResult{Valid: true, Reason: ReasonNone}, nil
	}
	switch {
	case hasExpired:
		return CheckResult{Valid: false, Reason: ReasonExpired}, nil
	case hasExhausted:
		return CheckResult{Valid: false, Reason: ReasonExhausted}, nil
	case anyVoided && !anyCount && !everValid:
		return CheckResult{Valid: false, Reason: ReasonVoided}, nil
	case anyVoided:
		return CheckResult{Valid: false, Reason: ReasonVoided}, nil
	default:
		return CheckResult{Valid: false, Reason: ReasonNoEntitlement}, nil
	}
}

// Consume takes n uses from the subject's count entitlements for one offer
// (D-002 §5.2): the attempt loop lives OUTSIDE store.Run; every attempt is a
// fresh transaction whose candidates are re-read; partial deductions roll back
// with the attempt; the final UPDATE re-checks all five predicates so a
// concurrent void or consumer linearizes per row.
func (s *Service) Consume(ctx context.Context, subjectID, offerID string, n int64, now time.Time) error {
	subjectID = strings.TrimSpace(subjectID)
	offerID = strings.TrimSpace(offerID)
	if subjectID == "" || offerID == "" || n < 1 {
		return store.ErrInvalidOffer
	}
	var lastErr error
	for attempt := 0; attempt < maxPurchaseAttempts; attempt++ {
		err := s.consumeAttempt(ctx, subjectID, offerID, n, now)
		if err == nil {
			return nil
		}
		if errors.Is(err, errConsumeInsufficient) {
			return store.ErrEntitlementInsufficient
		}
		lastErr = err
	}
	if lastErr == nil {
		lastErr = store.ErrEntitlementInsufficient
	}
	return fmt.Errorf("consume attempts exhausted: %w", lastErr)
}

// consumeAttempt outcome errors: errConsumeInsufficient is deterministic
// (every active row was drained inside this attempt); anything else is a
// race and the outer loop retries with a fresh transaction.
var errConsumeInsufficient = errors.New("digital-offer entitlement insufficient for this consume")

func (s *Service) consumeAttempt(ctx context.Context, subjectID, offerID string, n int64, now time.Time) error {
	return s.runnerRun(ctx, func(tx kernel.Tx) error {
		candidates, err := s.repo.ListConsumeCandidatesInTx(tx, subjectID, offerID)
		if err != nil {
			return err
		}
		needed := n
		for _, c := range candidates {
			if needed <= 0 {
				break
			}
			take := c.RemainingCount
			if take > needed {
				take = needed
			}
			ok, err := s.repo.DecrementEntitlementInTx(tx, c.ID, subjectID, offerID, take, now)
			if err != nil {
				return err
			}
			if ok {
				needed -= take
			}
		}
		if needed > 0 {
			return errConsumeInsufficient
		}
		return nil
	})
}

// Query bucket constants (D-002 §8): bizoffer|price|<subject_id> covers the
// Telegram price and entitlements queries — 1 minute / 30 requests.
const (
	queryLimiterWindow   = time.Minute
	queryLimiterMax      = 30
	queryLimiterCapacity = 1 << 16
)

// Telegram command names (D-001 I-031-004 / D-002 §6).
const (
	TelegramCommandPrice        = "price"
	TelegramCommandBuy          = "buy"
	TelegramCommandEntitlements = "entitlements"
)

// telegramUserTexts are the fail-closed reply fragments for §6 handlers.
const (
	telegramIdentityMissingText = "无法识别您的账户身份，请稍后重试或联系管理员完成绑定。"
	telegramBuyUsageText        = "用法：/buy <offer_id>（offer_id 见 /price 列表）"
	telegramNoEntitlementsText  = "您当前没有有效权益。"
)

// RegisterTelegram wires the §6 command handlers onto the channel dispatcher
// (kernel.TelegramDispatcher). When channel.telegram is disabled the
// composition injects the DisabledDispatcher no-op, so this call stays
// side-effect free and module tests never depend on the Bot API. Reply
// messages go through sender; a nil sender fails commands closed.
func (s *Service) RegisterTelegram(dispatcher kernel.TelegramDispatcher, sender kernel.TelegramSender) error {
	if dispatcher == nil {
		return errors.New("digitaloffer: telegram dispatcher is nil")
	}
	s.telegramSender = sender
	handlers := map[string]func(ctx context.Context, upd kernel.TelegramUpdate) error{
		TelegramCommandPrice:        s.telegramPrice,
		TelegramCommandBuy:          s.telegramBuy,
		TelegramCommandEntitlements: s.telegramEntitlements,
	}
	for name, h := range handlers {
		if err := dispatcher.RegisterCommand(name, h); err != nil {
			return fmt.Errorf("register telegram command %s: %w", name, err)
		}
	}
	return nil
}

func (s *Service) telegramPrice(ctx context.Context, upd kernel.TelegramUpdate) error {
	if err := s.queryAllowed(upd.SubjectID); err != nil {
		return s.telegramReply(ctx, upd, "请求过于频繁，请稍后再试。")
	}
	offers, err := s.ListOnSale()
	if err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("在售数字服务：\n")
	for i, o := range offers {
		fmt.Fprintf(&b, "%d. %s — %d %s（%s）\n   id: %s\n", i+1, o.Name, o.PriceAmount, o.Currency, formLabel(o.EntitlementForm), o.ID)
		if i == 19 {
			b.WriteString("…（仅显示前 20 项）\n")
			break
		}
	}
	if len(offers) == 0 {
		b.WriteString("（暂无在售项目）")
	}
	return s.telegramReply(ctx, upd, strings.TrimRight(b.String(), "\n"))
}

func (s *Service) telegramBuy(ctx context.Context, upd kernel.TelegramUpdate) error {
	// Fail-closed identity gate (D-002 §6): no mapped subject, no purchase.
	if strings.TrimSpace(upd.SubjectID) == "" {
		return s.telegramReply(ctx, upd, telegramIdentityMissingText)
	}
	fields := strings.Fields(strings.TrimPrefix(upd.Text, "/buy"))
	if len(fields) < 1 || strings.TrimSpace(fields[0]) == "" {
		return s.telegramReply(ctx, upd, telegramBuyUsageText)
	}
	offerID := fields[0]
	requestID := fmt.Sprintf("tg:%s:%d", upd.SubjectID, upd.UpdateID)
	res, err := s.Purchase(ctx, upd.SubjectID, offerID, requestID, time.Now().UTC())
	if err != nil {
		return s.telegramReply(ctx, upd, telegramErrorText(err))
	}
	var b strings.Builder
	if res.Replayed {
		b.WriteString("该购买已完成（幂等重放），未重复扣款。\n")
	} else {
		b.WriteString("购买成功！\n")
	}
	fmt.Fprintf(&b, "凭证：%s\n服务：%s\n金额：%d %s\n", res.Purchase.ID, res.Purchase.OfferName, res.Purchase.Amount, res.Purchase.Currency)
	switch res.Entitlement.Form {
	case store.FormDuration:
		fmt.Fprintf(&b, "权益：时长型，有效期至 %s", res.Entitlement.ExpiresAt.UTC().Format("2006-01-02 15:04"))
	case store.FormCount:
		fmt.Fprintf(&b, "权益：次数型，剩余 %d 次", *res.Entitlement.RemainingCount)
	}
	return s.telegramReply(ctx, upd, b.String())
}

func (s *Service) telegramEntitlements(ctx context.Context, upd kernel.TelegramUpdate) error {
	if err := s.queryAllowed(upd.SubjectID); err != nil {
		return s.telegramReply(ctx, upd, "请求过于频繁，请稍后再试。")
	}
	if strings.TrimSpace(upd.SubjectID) == "" {
		return s.telegramReply(ctx, upd, telegramIdentityMissingText)
	}
	now := time.Now().UTC()
	rows, total, err := s.ListEntitlements(store.EntitlementFilter{SubjectID: upd.SubjectID, Status: store.EntitlementActive, Page: 1, PageSize: 20})
	if err != nil {
		return err
	}
	if total == 0 {
		return s.telegramReply(ctx, upd, telegramNoEntitlementsText)
	}
	var b strings.Builder
	b.WriteString("我的有效权益：\n")
	shown := 0
	for _, r := range rows {
		if r.Form == store.FormDuration && r.ExpiresAt != nil && !r.ExpiresAt.After(now) {
			continue // expired rows are a derived predicate, never listed
		}
		if r.Form == store.FormCount && (r.RemainingCount == nil || *r.RemainingCount <= 0) {
			continue // exhausted rows are a derived predicate, never listed
		}
		shown++
		switch r.Form {
		case store.FormDuration:
			fmt.Fprintf(&b, "%d. %s — 时长型，有效期至 %s\n", shown, r.OfferID, r.ExpiresAt.UTC().Format("2006-01-02 15:04"))
		case store.FormCount:
			fmt.Fprintf(&b, "%d. %s — 次数型，剩余 %d 次\n", shown, r.OfferID, *r.RemainingCount)
		}
	}
	if shown == 0 {
		return s.telegramReply(ctx, upd, telegramNoEntitlementsText)
	}
	return s.telegramReply(ctx, upd, strings.TrimRight(b.String(), "\n"))
}

// queryAllowed applies the §8 price/entitlements query bucket.
func (s *Service) queryAllowed(subjectID string) error {
	if s.queryLimiter == nil {
		return nil
	}
	if !s.queryLimiter.AllowRecord("bizoffer|price|"+subjectID, time.Now().UTC()) {
		return ErrRateLimited
	}
	return nil
}

func (s *Service) telegramReply(ctx context.Context, upd kernel.TelegramUpdate, text string) error {
	if s.telegramSender == nil {
		return errors.New("digitaloffer: telegram sender is nil")
	}
	return s.telegramSender.Send(ctx, kernel.TelegramMessage{ChatID: upd.ChatID, Text: text})
}

// telegramErrorText maps purchase failures to user-facing reply text (§9).
func telegramErrorText(err error) string {
	switch {
	case errors.Is(err, walletstore.ErrInsufficient):
		return "余额不足，无法完成购买。"
	case errors.Is(err, store.ErrOfferNotOnSale):
		return "该服务当前未上架。"
	case errors.Is(err, store.ErrNotFound):
		return "未找到该 offer，请用 /price 查看在售列表。"
	case errors.Is(err, store.ErrCurrencyMismatch):
		return "服务币种与您的钱包账户不一致，无法购买。"
	case errors.Is(err, store.ErrSubjectNotFound):
		return "账户身份尚未注册，请稍后重试。"
	case errors.Is(err, store.ErrRequestIdConflict), errors.Is(err, walletstore.ErrIdempotencyConflict):
		return "请求冲突，请重新发送购买指令。"
	case errors.Is(err, ErrRateLimited):
		return "购买过于频繁，请稍后再试。"
	default:
		return "购买失败，请稍后重试。"
	}
}

func formLabel(form string) string {
	if form == store.FormDuration {
		return "时长型"
	}
	return "次数型"
}
