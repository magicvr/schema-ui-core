// Package wallet provides the admin.wallet module surface as a kernel.Provider
// (S-14 · GOAL-019 D-002): wallet accounts, immutable ledger, reconciliation,
// wallet.read / wallet.write / wallet.adjust permission keys, the wallet and
// wallet-entries pages and wallet.* audit events.
package wallet

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/manifest"
	walletschema "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/schema"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/voucher"
)

// ModuleID is the stable admin.wallet module identifier.
const ModuleID = "admin.wallet"

// Service implements handler.WalletService over the wallet store.
type Service struct {
	repo     *walletstore.Repository
	subjects *subject.Store
	vouchers *voucher.Service
	now      func() time.Time
}

// NewService constructs the wallet domain service.
func NewService(repo *walletstore.Repository, runners ...walletstore.TxRunner) *Service {
	var runner walletstore.TxRunner
	if len(runners) > 0 {
		runner = runners[0]
	}
	var subStore *subject.Store
	var vSvc *voucher.Service
	if runner != nil {
		subStore = subject.NewStore(runner)
		vSvc = voucher.NewService(runner, repo, subStore)
	}
	return &Service{
		repo:     repo,
		subjects: subStore,
		vouchers: vSvc,
		now:      time.Now,
	}
}

// SubjectStore returns the external subjects store.
func (s *Service) SubjectStore() *subject.Store {
	return s.subjects
}

// VoucherService returns the prepaid voucher service.
func (s *Service) VoucherService() *voucher.Service {
	return s.vouchers
}

// newID returns a time-ordered hex id: a Unix-millisecond prefix, a
// per-process monotonic counter (8 hex) for entries created within the SAME
// millisecond, and a random suffix that prevents primary-key collisions. The
// D-002 §1 chain order is (created_at ASC, id ASC); the millisecond prefix
// keeps same-second entries ordered, and the counter restores CREATION order
// inside one millisecond — the historical random-only suffix made same-ms
// entries sort by chance, so replay could run a freeze before its funding
// adjust and report "insufficient balance" (GOAL-037 / F-008, 2026-08-23).
var entryIDSeq atomic.Uint64

func newID(now time.Time) (string, error) {
	buf := make([]byte, 12)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	// Counter is per-process; wraps only after 2^32 same-millisecond entries.
	seq := entryIDSeq.Add(1) & 0xFFFFFFFF
	return fmt.Sprintf("%016x%08x%s", now.UnixMilli(), seq, hex.EncodeToString(buf)), nil
}

// ListAccounts implements handler.WalletService.
func (s *Service) ListAccounts(q, ownerType string, page, pageSize int) ([]walletstore.Account, int, error) {
	return s.repo.ListAccounts(walletstore.ListFilter{Q: q, OwnerType: ownerType, Page: page, PageSize: pageSize})
}

// GetAccount implements handler.WalletService.
func (s *Service) GetAccount(id string) (*walletstore.Account, error) {
	return s.repo.GetAccount(id)
}

// CreateAccount implements handler.WalletService. ownerType must be a known
// owner kind, ownerID non-empty (A-007 F-003: production validation matches
// the test double); currency defaults to CNY.
func (s *Service) CreateAccount(ownerType, ownerID, currency string, now time.Time) (*walletstore.Account, error) {
	switch ownerType {
	case walletstore.OwnerUser, walletstore.OwnerBusiness, walletstore.OwnerSystem, walletstore.OwnerSubject:
	default:
		return nil, walletstore.ErrInvalidEntry
	}
	if ownerID == "" {
		return nil, walletstore.ErrInvalidEntry
	}
	// For subject account: must reference an existing external subject (W13 F-012 principle, no orphan ledgers).
	if ownerType == walletstore.OwnerSubject && s.subjects != nil {
		exists, err := s.subjects.SubjectExists(context.Background(), ownerID)
		if err != nil {
			return nil, err
		}
		if !exists {
			return nil, walletstore.ErrNotFound
		}
	}
	id, err := newID(now)
	if err != nil {
		return nil, err
	}
	if currency == "" {
		currency = walletstore.DefaultCurrency
	}
	acct := walletstore.Account{
		ID: id, OwnerType: ownerType, OwnerID: ownerID, Currency: currency,
		Status: walletstore.StatusActive, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.repo.CreateAccount(acct); err != nil {
		return nil, err
	}
	return &acct, nil
}

// UpdateStatus implements handler.WalletService (optimistic lock).
func (s *Service) UpdateStatus(id, status string, version int64, now time.Time) (*walletstore.Account, error) {
	return s.repo.UpdateStatus(id, status, version, now)
}

// ListEntries implements handler.WalletService.
func (s *Service) ListEntries(accountID, entryType, q string, page, pageSize int) ([]walletstore.LedgerEntry, int, error) {
	if _, err := s.repo.GetAccount(accountID); err != nil {
		return nil, 0, err
	}
	return s.repo.ListEntries(accountID, entryType, q, page, pageSize)
}

// GetOrCreateUserAccount implements handler.WalletService (GOAL-020 D-001 §1).
func (s *Service) GetOrCreateUserAccount(ownerID string, now time.Time) (*walletstore.Account, bool, error) {
	return s.repo.GetOrCreateUserAccount(ownerID, now)
}

// GetUserAccountByOwner implements handler.WalletService (W15-F11 read-only).
func (s *Service) GetUserAccountByOwner(ownerID string) (*walletstore.Account, error) {
	return s.repo.GetUserAccountByOwner(ownerID)
}

// Mutate implements handler.WalletService.
func (s *Service) Mutate(id string, in walletstore.LedgerEntryInput, now time.Time) (*walletstore.Account, *walletstore.LedgerEntry, bool, error) {
	if in.Memo == "" {
		return nil, nil, false, walletstore.ErrInvalidEntry
	}
	entryID, err := newID(now)
	if err != nil {
		return nil, nil, false, err
	}
	account, entry, err := s.repo.Mutate(id, in, entryID, now)
	if err != nil {
		return nil, nil, false, err
	}
	return account, entry, entry.ID != entryID, nil
}

// Reconcile implements handler.WalletService (read path, persisted run).
func (s *Service) Reconcile(accountID, actorID string, now time.Time) (*walletstore.ReconciliationRun, error) {
	runID, err := newID(now)
	if err != nil {
		return nil, err
	}
	return s.repo.ReconcileRun(accountID, runID, actorID, now)
}

// ReconcileOnceTx implements the Job consumer's atomic commit callback.
func (s *Service) ReconcileOnceTx(ctx context.Context, tx kernel.Tx, accountID, runID, actorID string, now time.Time) (*walletstore.ReconciliationRun, error) {
	return s.repo.ReconcileOnceTx(ctx, tx, accountID, runID, actorID, now)
}

// ListReconcileRuns implements handler.WalletService.
func (s *Service) ListReconcileRuns(page, pageSize int) ([]walletstore.ReconciliationRun, int, error) {
	return s.repo.ListReconcileRuns(page, pageSize)
}

// GenerateVouchers implements handler.WalletService (VP-029 R3).
func (s *Service) GenerateVouchers(ctx context.Context, batchID string, count int, amount int64, currency string, expiresAt *time.Time, now time.Time) ([]voucher.GeneratedVoucher, error) {
	if s.vouchers == nil {
		return nil, errors.New("voucher service not initialized")
	}
	return s.vouchers.GenerateBatch(ctx, batchID, count, amount, currency, expiresAt, now)
}

// ListVouchers implements handler.WalletService (VP-029 R3).
func (s *Service) ListVouchers(ctx context.Context, batchID, status string, page, pageSize int) ([]voucher.Voucher, int, error) {
	if s.vouchers == nil {
		return nil, 0, errors.New("voucher service not initialized")
	}
	return s.vouchers.ListVouchers(ctx, batchID, status, page, pageSize)
}

// GetVoucher implements handler.WalletService (VP-029 R3).
func (s *Service) GetVoucher(ctx context.Context, id string) (*voucher.Voucher, error) {
	if s.vouchers == nil {
		return nil, errors.New("voucher service not initialized")
	}
	return s.vouchers.GetVoucher(ctx, id)
}

// VoidVoucher implements handler.WalletService (VP-029 R3).
func (s *Service) VoidVoucher(ctx context.Context, id string, now time.Time) error {
	if s.vouchers == nil {
		return errors.New("voucher service not initialized")
	}
	return s.vouchers.VoidVoucher(ctx, id, now)
}

// RedeemForUser implements handler.WalletService (VP-029 R5 · GOAL-005).
func (s *Service) RedeemForUser(ctx context.Context, userID, actorName, code string, now time.Time) (*voucher.RedeemResult, error) {
	if s.vouchers == nil {
		return nil, errors.New("voucher service not initialized")
	}
	return s.vouchers.RedeemForUser(ctx, userID, actorName, code, now)
}

// Provider implements kernel.Provider for admin.wallet.
type Provider struct {
	a          *auth.Authenticator
	service    *Service
	jobs       *JobService
	operations operationlog.Recorder
	// ownerExists is the W13 F-012 gate (GOAL-013 A-001): auto-create wallet
	// paths verify the owner id against the live user table before opening an
	// account. Wired from the composition root (auth-session repository).
	ownerExists handler.OwnerExistsFunc
	limiters    kernel.RateLimiterProvider
}

// New constructs the wallet provider.
func New(a *auth.Authenticator, service *Service, jobs *JobService, operations operationlog.Recorder, ownerExists handler.OwnerExistsFunc, limiters kernel.RateLimiterProvider) *Provider {
	return &Provider{a: a, service: service, jobs: jobs, operations: operations, ownerExists: ownerExists, limiters: limiters}
}

func (p *Provider) Descriptor() kernel.Module {
	return kernel.Module{
		ID:             ModuleID,
		Version:        "2.0.0",
		KernelAPIRange: ">=2.0 <3.0",
		DependsOn:      []string{"core.auth-session", "core.navigation-capability", "core.schema-render", "core.operationlog"},
		Requires:       kernel.StandardAdminCapabilities(),
		Contributions: kernel.ContributionKeys{
			Routes: []string{
				"GET /api/wallet/accounts", "POST /api/wallet/accounts",
				"PATCH /api/wallet/accounts/{id}",
				"GET /api/wallet/by-owner/{ownerId}", "POST /api/wallet/by-owner/{ownerId}",
				"POST /api/wallet/by-owner/{ownerId}/adjust",
				"GET /api/wallet/accounts/{id}/entries", "GET /api/wallet/entries",
				"POST /api/wallet/accounts/{id}/adjust",
				"POST /api/wallet/accounts/{id}/freeze",
				"POST /api/wallet/accounts/{id}/unfreeze",
				"POST /api/wallet/accounts/{id}/deduct-frozen",
				"POST /api/wallet/reconcile", "GET /api/wallet/reconcile/runs",
				"GET /api/wallet/jobs/{id}", "POST /api/wallet/jobs/{id}/cancel",
				"POST /api/wallet/jobs/{id}/retry", "GET /api/wallet/jobs/{id}/result",
				// GOAL-022 (D-002 §2): identity-scoped self-service surface.
				"GET /api/wallet/me", "POST /api/wallet/me", "GET /api/wallet/me/entries",
				"POST /api/wallet/me/redeem",
				// VP-029 R3 (GOAL-003): prepaid vouchers surface.
				"POST /api/wallet/vouchers/batches", "GET /api/wallet/vouchers",
				"GET /api/wallet/vouchers/{id}", "POST /api/wallet/vouchers/{id}/void",
			},
			Pages:       []string{"wallet", "wallet-entries", "my-wallet", "wallet-vouchers"},
			Navigation:  []string{"menu_wallet", "menu_wallet_self", "menu_wallet_vouchers"},
			Permissions: []string{"wallet.read", "wallet.write", "wallet.adjust", "wallet.voucher.issue"},
			Fragments:   []string{"wallet"},
		},
	}
}

func (p *Provider) CompiledPersistence() ([]kernel.MigrationContribution, error) {
	return nil, nil // tables are owned by the wallet/migration provider (0031)
}

func (p *Provider) Register(ctx context.Context, reg kernel.Registrar) error {
	for _, route := range handler.WalletRoutes(p.a, p.service, p.jobs, p.operations, ModuleID, p.ownerExists) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	// GOAL-022 (D-002 §2): identity-scoped self-service routes.
	for _, route := range handler.WalletSelfRoutes(p.a, p.service, p.operations, ModuleID, p.limiters) {
		if err := reg.HTTP(route); err != nil {
			return err
		}
	}
	for _, pageID := range []string{"wallet", "wallet-entries"} {
		if err := reg.Schema(kernel.PageContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: pageID},
			PageID:               pageID,
			Resources:            []string{"wallet"},
			Actions:              []string{"list", "create", "update"},
			DataSource:           "/api/wallet/accounts",
			Owner:                ModuleID,
			Document:             walletschema.SchemaDocuments()[pageID],
		}); err != nil {
			return err
		}
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "my-wallet"},
		PageID:               "my-wallet",
		Resources:            []string{"wallet"},
		Actions:              []string{"list"},
		DataSource:           "/api/wallet/me",
		Owner:                ModuleID,
		Document:             walletschema.SchemaDocuments()["my-wallet"],
	}); err != nil {
		return err
	}
	if err := reg.Schema(kernel.PageContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet-vouchers"},
		PageID:               "wallet-vouchers",
		Resources:            []string{"wallet"},
		Actions:              []string{"list", "create", "update"},
		DataSource:           "/api/wallet/vouchers",
		Owner:                ModuleID,
		Document:             walletschema.SchemaDocuments()["wallet-vouchers"],
	}); err != nil {
		return err
	}
	for _, permission := range []kernel.PermissionContribution{
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet.read"}, Permission: "wallet.read", Resource: "wallet", Action: "read", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet.write"}, Permission: "wallet.write", Resource: "wallet", Action: "write", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet.adjust"}, Permission: "wallet.adjust", Resource: "wallet", Action: "adjust", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
		{ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet.voucher.issue"}, Permission: "wallet.voucher.issue", Resource: "wallet", Action: "voucher.issue", PolicyID: authsessiondata.PolicyAdmin, SystemDataVersion: authsessiondata.SystemDataVersion},
	} {
		if err := reg.Authorization(permission); err != nil {
			return err
		}
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_wallet"},
		NodeID:               "menu_wallet",
		PageID:               "wallet",
		Order:                10,
		Label:                "Wallet",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "wallet.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_wallet_vouchers"},
		NodeID:               "menu_wallet_vouchers",
		PageID:               "wallet-vouchers",
		Order:                11,
		Label:                "Prepaid vouchers",
		Visibility:           authsessiondata.PolicyAdmin,
		Permission:           "wallet.read",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	// GOAL-022 (D-002 §2): my-wallet self-service entry in the topbar user slot
	// (个人中心 → 我的钱包 → 设置). Identity-only: visible to every
	// authenticated role, no permission key — like menu_account.
	if err := reg.Navigation(kernel.NavigationContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "menu_wallet_self"},
		NodeID:               "menu_wallet_self",
		PageID:               "my-wallet",
		Order:                2,
		Label:                "My wallet",
		Visibility:           authsessiondata.PolicyAdminEditorViewer,
		Permission:           "",
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	}); err != nil {
		return err
	}
	return reg.Manifest(kernel.FragmentContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: ModuleID, Key: "wallet"},
		FragmentID:           "wallet",
		ProtocolVersion:      "2.7",
		RequiredCapabilities: []string{"manifest", "navigation"},
		JSON:                 manifest.FragmentJSON,
	})
}
