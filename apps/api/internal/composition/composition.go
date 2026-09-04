package composition

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"time"

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/cache"
	telegraminternal "github.com/magicvr/schema-ui-core/apps/api/internal/channel/telegram"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/eventbus"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/jobs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/internal/manifest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/objectstore"
	"github.com/magicvr/schema-ui-core/apps/api/internal/obs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/internal/server"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	accountmodule "github.com/magicvr/schema-ui-core/apps/api/modules/account"
	activitymodule "github.com/magicvr/schema-ui-core/apps/api/modules/activity"
	authsession "github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	telegrammodule "github.com/magicvr/schema-ui-core/apps/api/modules/channel/telegram"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
	dashboardmodule "github.com/magicvr/schema-ui-core/apps/api/modules/dashboard"
	datadictionarymodule "github.com/magicvr/schema-ui-core/apps/api/modules/datadictionary"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/modules/datadictionary/store"
	datadictionarystore2 "github.com/magicvr/schema-ui-core/apps/api/modules/datadictionary/store"
	datapermissionmodule "github.com/magicvr/schema-ui-core/apps/api/modules/datapermission"
	datapermissionstore "github.com/magicvr/schema-ui-core/apps/api/modules/datapermission/store"
	datatransfermodule "github.com/magicvr/schema-ui-core/apps/api/modules/datatransfer"
	devexamplesmodule "github.com/magicvr/schema-ui-core/apps/api/modules/dev/examples"
	filelibrarymodule "github.com/magicvr/schema-ui-core/apps/api/modules/filelibrary"
	logincaptchamodule "github.com/magicvr/schema-ui-core/apps/api/modules/logincaptcha"
	logincaptchastore "github.com/magicvr/schema-ui-core/apps/api/modules/logincaptcha/store"
	mfamodule "github.com/magicvr/schema-ui-core/apps/api/modules/mfa"
	mfastore "github.com/magicvr/schema-ui-core/apps/api/modules/mfa/store"
	notificationsmodule "github.com/magicvr/schema-ui-core/apps/api/modules/notifications"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	recyclebinmodule "github.com/magicvr/schema-ui-core/apps/api/modules/recyclebin"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/modules/recyclebin/store"
	rolesmodule "github.com/magicvr/schema-ui-core/apps/api/modules/roles"
	scheduledtasksmodule "github.com/magicvr/schema-ui-core/apps/api/modules/scheduledtasks"
	scheduledtasksstore "github.com/magicvr/schema-ui-core/apps/api/modules/scheduledtasks/store"
	tasksstore2 "github.com/magicvr/schema-ui-core/apps/api/modules/scheduledtasks/store"
	schemarendermodule "github.com/magicvr/schema-ui-core/apps/api/modules/schemarender"
	settingsmodule "github.com/magicvr/schema-ui-core/apps/api/modules/settings"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/modules/settings/repository"
	systemmonitoringmodule "github.com/magicvr/schema-ui-core/apps/api/modules/systemmonitoring"
	usersmodule "github.com/magicvr/schema-ui-core/apps/api/modules/users"
	walletmodule "github.com/magicvr/schema-ui-core/apps/api/modules/wallet"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/modules/wallet/store"
	"github.com/magicvr/schema-ui-core/apps/api/modules/wallet/subject"
	"github.com/magicvr/schema-ui-core/apps/api/pkg/version"
)

type jwtSecret string
type seedPasswordHash string

// ResolvePlan expands the selected Profile and validates the complete compiled
// module graph before Fx starts any resource.
func ResolvePlan(cfg *config.Config) (kernel.Plan, error) {
	if cfg == nil {
		return kernel.Plan{}, fmt.Errorf("composition: config is required")
	}
	if cfg.ProfileError != nil {
		return kernel.Plan{}, fmt.Errorf("composition: invalid profile configuration: %w", cfg.ProfileError)
	}
	modules := append([]string(nil), cfg.ModulesEnabled...)
	if len(modules) == 0 {
		resolved, err := kernel.ResolveProfile(cfg.ProfileName, nil)
		if err != nil {
			return kernel.Plan{}, err
		}
		modules = resolved.Modules
	}
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		return kernel.Plan{}, err
	}
	plan, err := registry.Resolve(modules)
	if err != nil {
		return kernel.Plan{}, err
	}
	// GOAL-013 D-002 §4: the operator-provided navigation order rides on the
	// resolved plan so the kernel sorts navigation deterministically.
	if len(cfg.NavigationOrder) > 0 {
		plan.NavigationOrder = append([]string(nil), cfg.NavigationOrder...)
	}
	return plan, nil
}

// NewApp creates the Fx composition root. Fx types remain confined to this
// package; module descriptors and kernel contracts remain framework agnostic.
func NewApp(cfg *config.Config, secretValue, seedHash string, logger *slog.Logger) (*fx.App, error) {
	return newAppWithOptions(cfg, secretValue, seedHash, logger)
}

// newAppWithOptions is the composition-root test seam (F-001 / A-006): it builds
// the exact same Fx graph as NewApp plus caller-supplied fx.Options (e.g.
// fx.Populate) so a test can prove the injected *TelegramRuntime is the SAME
// instance the webhook mounts, without bypassing the Fx graph.
func newAppWithOptions(cfg *config.Config, secretValue, seedHash string, logger *slog.Logger, extra ...fx.Option) (*fx.App, error) {
	plan, err := ResolvePlan(cfg)
	if err != nil {
		return nil, err
	}
	if logger == nil {
		logger = slog.Default()
	}
	opts := []fx.Option{
		fx.Provide(
			func() *config.Config { return cfg },
			func() jwtSecret { return jwtSecret(secretValue) },
			func() seedPasswordHash { return seedPasswordHash(seedHash) },
			func() *slog.Logger { return logger },
			func() kernel.Plan { return plan },
			func() *readinessGate { return &readinessGate{} },
			openStore,
			newAuthSessionRepository,
			newOperationLogRepository,
			newSettingsRepository,
			newAuthenticator,
			newJobRuntime,
			newTracing,
			newObserver,
			newMetricsServer,
			newCache,
			newEventBus,
			newRateLimiters,
			newTelegramRuntime,
			func(tr *TelegramRuntime) kernel.TelegramDispatcher { return tr.Dispatcher },
			func(tr *TelegramRuntime) kernel.TelegramSender { return tr.Sender },
			newMux,
			newServer,
		),
		fx.Invoke(registerLifecycle),
	}
	opts = append(opts, extra...)
	return fx.New(opts...), nil
}

// readinessGate reports whether the module graph Start+Ready succeeded (R5 real
// readiness). readyz consults it; registerLifecycle flips it after Ready.
type readinessGate struct {
	ready atomic.Bool
}

func (g *readinessGate) Ready() bool { return g.ready.Load() }
func (g *readinessGate) setReady()   { g.ready.Store(true) }

type jobRuntime struct {
	repository *jobs.Repository
	runner     *jobs.Runner
	enabled    atomic.Bool
}

func newJobRuntime(st kernel.Store) (*jobRuntime, error) {
	repository := jobs.NewRepository(st)
	runner, err := jobs.NewRunner(repository, jobs.DefaultRunnerOptions())
	if err != nil {
		return nil, err
	}
	return &jobRuntime{repository: repository, runner: runner}, nil
}

func (r *jobRuntime) Start() error {
	if !r.enabled.Load() {
		return nil
	}
	return r.runner.Start()
}

func (r *jobRuntime) Stop(ctx context.Context) error {
	if !r.enabled.Load() {
		return nil
	}
	return r.runner.Stop(ctx)
}

func openStore(cfg *config.Config, seedHash seedPasswordHash) (kernel.Store, error) {
	// Persistence is compiled-global rather than profile-gated. The static
	// registry collects module-owned descriptors before store startup; store
	// only validates and executes the resulting catalog.
	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.persistence", Detail: fmt.Sprintf("collect persistence: %v", err)}
	}
	// VP-013 R1 v1.4: sqlite keeps applying the catalog on the default dev
	// path; postgres (R3 dual-dialect ledger) applies it too. From R4 the
	// composition root wires the kernel.Store interface, so a postgres DSN can
	// boot the full application (repositories speak kernel.Tx).
	dialect := kernel.DialectSQLite
	if cfg.DBDialect != "" {
		dialect = kernel.Dialect(cfg.DBDialect)
	}
	st, err := store.Open(context.Background(), store.OpenOptions{
		Dialect:          dialect,
		Path:             cfg.DBPath,
		DSN:              cfg.DBDSN,
		PoolMaxOpenConns: cfg.DBPoolMaxOpen,
		PoolMaxIdleConns: cfg.DBPoolMaxIdle,
		ConnMaxLifetime:  cfg.DBConnLifetime,
	}, catalog)
	if err != nil {
		return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("open store: %v", err)}
	}
	// (db.path) file-storage root derivation and upload/avatar dirs already use
	// cfg.DBPath below; postgres keeps the same file-path-shaped value.
	needsBootstrap, err := authsessiondata.NeedsBootstrap(context.Background(), st)
	if err != nil {
		_ = st.Close()
		return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("check bootstrap: %v", err)}
	}
	if needsBootstrap {
		if err := authsessiondata.Bootstrap(context.Background(), st, "admin", string(seedHash)); err != nil {
			_ = st.Close()
			return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("bootstrap auth data: %v", err)}
		}
	}
	// TEST_ADMIN_USERNAME / TEST_ADMIN_PASSWORD (optional): a stable test-only
	// admin credential upserted on every boot — does not touch the "admin"
	// bootstrap user. Empty TEST_ADMIN_PASSWORD = feature off.
	if cfg.TestAdminPassword != "" {
		username := cfg.TestAdminUsername
		if strings.TrimSpace(username) == "" {
			username = "testadmin"
		}
		testHash, err := auth.HashPassword(cfg.TestAdminPassword, 10)
		if err != nil {
			_ = st.Close()
			return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("hash test-admin password: %v", err)}
		}
		if err := authsessiondata.EnsureTestAdmin(context.Background(), st, username, testHash); err != nil {
			_ = st.Close()
			return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("upsert test admin: %v", err)}
		}
	}
	return st, nil
}

func newAuthSessionRepository(st kernel.Store) *authsession.Repository {
	return authsession.NewRepository(st)
}

func newOperationLogRepository(st kernel.Store) *operationlog.Repository {
	return operationlog.NewRepository(st)
}

func newSettingsRepository(st kernel.Store) *settingsrepository.Repository {
	return settingsrepository.New(st)
}

func newAuthenticator(cfg *config.Config, secret jwtSecret, repository *authsession.Repository) *auth.Authenticator {
	// VP-016 R2 (workspace-016 GOAL-003 D-001): an empty previous keeps exact
	// single-key behavior; a configured one opens the rotation overlap window
	// (verify current, fall back to previous). Strength/difference rules for
	// the pair are enforced earlier by config.ValidateProd.
	return auth.NewWithRepositoryAndPrevious([]byte(secret), []byte(cfg.AuthJWTSecretPrevious), cfg.AuthAccessTTL, cfg.AuthRefreshTTL, repository, cfg.AuthDevSessionEnabled)
}

// newTracing maps the observability.traces config surface onto the tracer
// path (GOAL-004 D-001 §3): disabled (the default) is a pure no-op.
func newTracing(cfg *config.Config) *obs.Tracing {
	return obs.NewTracing(obs.TracingOptions{
		Enabled:     cfg.TracesEnabled,
		Endpoint:    cfg.TracesEndpoint,
		SampleRatio: cfg.TracesSampleRatio,
		ServiceName: cfg.AppName,
		Version:     version.Version,
		Environment: cfg.AppEnv,
	}, slog.Default())
}

// newObserver builds the kernel metrics registry (VP-015 R2 / GOAL-003
// D-001 §5): static build identity plus one suc_kernel_modules_enabled line
// per enabled module in the resolved plan. Tracing is attached so Wrap emits
// both metrics and (when enabled) server spans from the same interception
// point.
func newObserver(cfg *config.Config, plan kernel.Plan, tracing *obs.Tracing) *obs.Observer {
	observer := obs.NewObserver(obs.BuildInfoFromVersion(cfg.ProfileName))
	observer.SetTracing(tracing)
	for _, id := range plan.IDs() {
		observer.RegisterModule(id)
	}
	return observer
}

// newMetricsServer maps the observability.metrics config surface onto the
// dedicated listener. Disabled (the default) means nothing listens.
func newMetricsServer(cfg *config.Config, observer *obs.Observer, logger *slog.Logger) *obs.Server {
	return obs.NewServer(obs.ServerOptions{
		Enabled:   cfg.MetricsEnabled,
		Addr:      cfg.MetricsAddr,
		AuthToken: cfg.MetricsAuthToken,
	}, observer, logger)
}

func newMux(
	cfg *config.Config,
	a *auth.Authenticator,
	st kernel.Store,
	authRepository *authsession.Repository,
	operations *operationlog.Repository,
	settingsRepository *settingsrepository.Repository,
	plan kernel.Plan,
	gate *readinessGate,
	secret jwtSecret,
	jobRuntime *jobRuntime,
	observer *obs.Observer,
	logger *slog.Logger,
	cachePort kernel.Cache,
	eventBusPort kernel.EventBus,
	rateLimiters kernel.RateLimiterProvider,
	tr *TelegramRuntime,
) (*http.ServeMux, error) {
	return newMuxWithExtraProviders(cfg, a, st, authRepository, operations, settingsRepository, plan, gate, secret, jobRuntime, nil, observer, logger, cachePort, eventBusPort, rateLimiters, tr)
}

// newMuxWithExtraProviders is the composition-root assembly seam used by the S2
// access drill: it appends test-only providers (e.g. a probe module) to the
// compiled one-party set before registration, proving that a new standard module
// surfaces through the Provider contract without any handler/Web central business
// registration change. Production call sites use newMux (extra == nil).
func newMuxWithExtraProviders(
	cfg *config.Config,
	a *auth.Authenticator,
	st kernel.Store,
	authRepository *authsession.Repository,
	operations *operationlog.Repository,
	settingsRepository *settingsrepository.Repository,
	plan kernel.Plan,
	gate *readinessGate,
	secret jwtSecret,
	jobRuntime *jobRuntime,
	extra []kernel.Provider,
	observer *obs.Observer,
	logger *slog.Logger,
	cachePort kernel.Cache,
	eventBusPort kernel.EventBus,
	rateLimiters kernel.RateLimiterProvider,
	tr *TelegramRuntime,
) (*http.ServeMux, error) {
	// VP-015 R2 (GOAL-003 D-001 §1): the instrumented mux is the single
	// interception point — central handler registrations (Handle/HandleFunc)
	// and module-contributed routes are all measured with their owning
	// module_id. A nil observer (metrics disabled) passes through untouched.
	mux := obs.NewInstrumentedMux(observer)
	// W7 F-008: install the explicit trusted reverse-proxy CIDR allow-list for
	// login/captcha client-IP resolution (fail-closed on invalid CIDRs).
	if err := handler.SetTrustedProxyCIDRs(cfg.HTTPTrustedProxies); err != nil {
		return nil, err
	}
	a.SetServiceCredentialUseTransactionalRecorder(func(tx kernel.Tx, use auth.ServiceCredentialUse) error {
		detail, err := operationlog.NewDetail("service-credential-use", nil, map[string]any{
			"credentialId": use.CredentialID,
			"scopeCount":   use.ScopeCount,
			"method":       use.Method,
			"path":         use.Path,
		})
		if err != nil {
			return err
		}
		recordID := use.CredentialID
		return operations.RecordOperationTx(tx, operationlog.Operation{
			ID: "op-service-" + auth.NewServiceCredentialID(), Event: operationlog.EventServiceCredentialUse,
			ActorID: "service-credential:" + use.CredentialID, ActorName: use.Name,
			RecordID: &recordID, Detail: &detail, CorrelationID: use.CorrelationID, SessionID: use.CredentialID, CreatedAt: use.At,
		})
	})
	// S-11 (GOAL-011): one shared captcha service instance feeds both the login
	// gate verifier (handler.RegisterWithReadiness) and the module provider
	// surface (routes/settings). nil when the module is disabled.
	var captchaService *logincaptchamodule.Service
	// S-11 (GOAL-011): the login captcha gate is wired through the variadic
	// CaptchaVerifier — nil (module disabled) keeps the login contract
	// byte-identical to the pre-captcha behavior.
	var captchaVerifier handler.CaptchaVerifier
	if plan.HasModule("admin.login-captcha") {
		captchaService = logincaptchamodule.NewService(logincaptchastore.NewRepository(st))
		captchaVerifier = captchaService
	}
	// S-10 (GOAL-017 D-002 §3): the MFA service feeds the second-factor
	// login gate (nil when admin.mfa is disabled — login contract unchanged)
	// and the module surface (verify / self-service / admin reset). The
	// verifier is declared as the interface type so a disabled module yields
	// a TRUE nil (a typed-nil *Service through the variadic would satisfy the
	// != nil check and panic on use — testhelpers doc comment precedent).
	var mfaService *mfamodule.Service
	var mfaVerifier handler.MFAVerifier
	if plan.HasModule("admin.mfa") {
		// W11 F-004: the previous JWT secret (VP-016 rotation window) is passed
		// to MFA too, so a mid-rotation AUTH_JWT_SECRET change does not lock MFA
		// users into an undecryptable second factor (empty = single-key behavior).
		mfaService = mfamodule.NewService(mfastore.NewRepository(st), []byte(secret), []byte(cfg.AuthJWTSecretPrevious))
		mfaVerifier = mfaService
	}
	// VP-014 R3 (GOAL-004 D-001): ONE kernel.ObjectStore instance serves all
	// three first-party families — local disk by default (root derived from
	// db.path unless storage.objects.local.root overrides) or the explicitly
	// configured S3-compatible backend, whose HeadBucket probe also extends
	// readyz ("配置后 readyz 扩依赖", GOAL-003).
	objects, objectProbe, err := newObjectStore(cfg)
	if err != nil {
		return nil, err
	}
	// VP-017 R6+R7 (GOAL-007/008 D-001 over the GOAL-006 D-002 frozen
	// contract): ONE kernel.MailSender — the SwitchingSender. The runtime row
	// (migration 0052) seeds ONCE from the file/env resolution; admin saves
	// then hot-switch the active channel (single-process, secrets AES-GCM
	// encrypted at rest under a local master key, never read back). readyz
	// keeps its R4 semantics: only an explicitly configured SMTP boot channel
	// contributes the ESMTP Ping probe; production probes extend in R8. The
	// mock record read face and the admin config/test-send surface register
	// independently of the active channel.
	mailSender, mailProbe, err := newMailRuntime(cfg, st, logger)
	if err != nil {
		return nil, err
	}
	// VP-026 / workspace-026 GOAL-003/004 D-001 (R2/R3): the kernel.Cache
	// single instance arrives via dependency injection — the Fx container
	// (fx.Provide(newCache)) owns it for the process lifetime (F-002
	// disposition), and this seam is its explicit injection point: the FIRST
	// consumer (a future business-domain module or an explicit caching need)
	// requests it here. No probe: an in-process store has no external
	// dependency (Redis stays trigger-gated; the seam declaration lives in
	// docs/architecture/cache-redis-seam-and-track.md).
	_ = cachePort // accepted, not yet consumed (intentional until a consumer lands)
	logger.Info("kernel cache port ready", "provider", "memory", "max_entries", cfg.CacheMaxEntries)
	// VP-028 / workspace-028 GOAL-003 D-001 (R2): the kernel.EventBus single
	// instance arrives via dependency injection — the Fx container owns it for
	// the process lifetime, and this seam is its explicit injection point. No
	// probe: an in-process channel bus has no external dependency. Stop drain
	// happens in registerLifecycle OnStop (D-002 §5).
	_ = eventBusPort // accepted, not yet consumed (intentional until a consumer lands)
	logger.Info("kernel event-bus port ready", "provider", "memory", "buffer_size", cfg.EventBusBufferSize)
	handler.RegisterMailOutbox(mux, a, mail.NewOutboxSink(st, mail.DefaultOutboxCap))
	handler.RegisterMailAdmin(mux, a, mailSender, operations)
	handler.RegisterWithMFAProbes(mux, a, st, operations, plan, gate.Ready, rateLimiters, []handler.CaptchaVerifier{captchaVerifier}, mfaVerifier, objectProbe, mailProbe)
	// workspace-019 R2 (GOAL-003 D-001 §2): the self-recovery start/complete
	// pair is a CENTRAL pre-auth surface (same layer as login) so every
	// profile with core.auth-session gets it. The completion second-factor
	// gate reuses the MFA service; admin.mfa off keeps a TRUE nil interface,
	// which means no second factor is demanded (GOAL-002 D-001 §1).
	var recoveryGate handler.RecoverySecondFactor
	if plan.HasModule("admin.mfa") {
		recoveryGate = mfaService
	}
	handler.RegisterInviteAccept(mux, authRepository, rateLimiters)
	handler.RegisterRecovery(mux, operations, authRepository, authRepository, mailSender, recoveryGate, rateLimiters)
	// I-PROTO-FULL-001 D-UPLOAD: server-side upload contract (07 §7.2). The
	// uploads namespace is shared with admin.data-transfer (F-02 import reads
	// uploaded CSV files by id) and admin.file-library.
	handler.RegisterUpload(mux, a, objects,
		handler.WithAllowedTypes(cfg.UploadAllowedTypes),
		handler.WithUserLimits(cfg.UploadMaxFilesPerUser, cfg.UploadMaxBytesPerUser),
	)
	// W9 (GOAL-010): dedicated brand-assets store — NOT the shared upload
	// store (owner-gated reads) and NOT admin.file-library. Brand icons must
	// be publicly readable (login page / shell load pre-auth); every stored
	// object is a server-side re-encoded raster (never raw upload bytes).
	brandAssets := handler.NewBrandingAssetStore(
		objects,
		handler.BrandingAssetsOptions{
			MaxBytes:    cfg.BrandingMaxBytes,
			LogoMaxDim:  cfg.BrandingLogoMaxDimension,
			FaviconDim:  cfg.BrandingFaviconDimension,
			JPEGQuality: cfg.BrandingJPEGQuality,
		},
	)
	// I-004: startup GC — drop orphan assets not referenced by the current
	// settings singleton (crashed uploads / cancelled edits).
	if current, err := settingsRepository.GetSiteSettings(); err == nil {
		_ = brandAssets.GC([]string{current.LogoURL, current.LogoURLLight, current.LogoURLDark, current.FaviconURL})
	}
	// W13 T-05 (GOAL-014): account avatar store — same raster processing as
	// brand assets, dedicated directory, 256px longest-edge default. W7 F-004:
	// startup GC now reclaims unreferenced avatar files (crashed uploads,
	// uploads never assigned to a profile) by keeping only the avatar URLs
	// currently referenced by users.
	avatarAssets := handler.NewAvatarAssetStore(
		objects,
		handler.BrandingAssetsOptions{},
	)
	if users, _, err := authRepository.ListUsers(authsession.UserFilter{Page: 1, PageSize: 1_000_000}); err == nil {
		refs := make([]string, 0, len(users))
		for _, u := range users {
			if u.AvatarURL != "" {
				refs = append(refs, u.AvatarURL)
			}
		}
		_ = avatarAssets.GC(refs)
	}
	// R4 C3.3: admin.users / admin.roles HTTP surface comes from the module
	// kernel.Provider contract (freeze package §7 step 3). Core auth/accounts/
	// health/schema stay central; settings/activity migrate in C4.
	// W1 (GOAL-002 / workspace-010): core.schema-render and the optional
	// dev.examples module are both assembled by plan enablement (D-003 §3) —
	// production profiles default to no examples surface.
	// S-12 (GOAL-012 D-002 §2): the recycle-bin service is constructed before
	// the managed modules so its TrashRecorder can be injected into their
	// delete hooks. nil (module disabled) keeps delete semantics unchanged.
	var recycleService *recyclebinmodule.Service
	var trash handler.TrashRecorder
	if plan.HasModule("admin.recycle-bin") {
		recycleService = recyclebinmodule.NewService(recyclestore.NewRepository(st), datadictionarystore2.NewRepository(st), tasksstore2.NewRepository(st), st)
		trash = recycleService
	}
	var providers []kernel.Provider
	if plan.HasModule("core.schema-render") {
		providers = append(providers, schemarendermodule.New())
	}
	if plan.HasModule("dev.examples") {
		providers = append(providers, devexamplesmodule.New())
	}
	if plan.HasModule("admin.users") {
		// W13 F-006 (GOAL-013 A-001): the canonical public origin (when
		// configured) replaces client-influenceable Host/Proto headers in
		// emailed invitation links.
		providers = append(providers, usersmodule.New(a, authRepository, operations, mailSender, cfg.AuthPublicBaseURL))
	}
	if plan.HasModule("admin.roles") {
		providers = append(providers, rolesmodule.New(a, authRepository, operations))
	}
	if plan.HasModule("admin.settings") {
		providers = append(providers, settingsmodule.New(a, settingsRepository, operations, brandAssets, authRepository))
	} else {
		// Public bootstrap must work on mvp (and any profile without the
		// settings edit module). Edit/list/patch/reset stay admin.settings-only;
		// previously uploaded brand assets stay publicly readable.
		handler.RegisterPublicBranding(mux, settingsRepository)
		handler.RegisterPublicBrandingAssets(mux, brandAssets)
	}
	if plan.HasModule("admin.activity") {
		providers = append(providers, activitymodule.New(a, operations))
	}
	if plan.HasModule("admin.account") {
		providers = append(providers, accountmodule.New(a, authRepository, operations, avatarAssets, mailSender, rateLimiters))
	}
	if plan.HasModule("admin.data-transfer") {
		providers = append(providers, datatransfermodule.New(a, authRepository, operations, objects))
	}
	if plan.HasModule("admin.dashboard") {
		providers = append(providers, dashboardmodule.New())
	}
	if plan.HasModule("admin.file-library") {
		providers = append(providers, filelibrarymodule.New(a, operations, objects))
	}
	if plan.HasModule("admin.data-dictionary") {
		providers = append(providers, datadictionarymodule.New(a, datadictionarystore.NewRepository(st), operations, trash))
	}
	if plan.HasModule("admin.system-monitoring") {
		providers = append(providers, systemmonitoringmodule.New(a, st, plan, gate.Ready, cfg.DBPath, time.Now(), operations, string(cfg.RuntimeMode)))
	}
	if plan.HasModule("admin.scheduled-tasks") {
		providers = append(providers, scheduledtasksmodule.New(a, scheduledtasksstore.NewRepository(st), operations, trash))
	}
	if plan.HasModule("admin.login-captcha") {
		providers = append(providers, logincaptchamodule.New(a, captchaService, operations, rateLimiters))
	}
	// S-09 (GOAL-016 D-002 §2): the data-permission service feeds both the
	// management routes and the RowScopeProvider contract. v1 wires NO
	// production resource into the enforceable set — registration is left to
	// future domain modules (A-005: the gate lives on the policy PATCH).
	if plan.HasModule("admin.data-permission") {
		dataPermissionService := datapermissionmodule.NewService(datapermissionstore.NewRepository(st), nil)
		providers = append(providers, datapermissionmodule.New(a, dataPermissionService, operations))
	}
	// S-10 (GOAL-017 D-002 §4): admin.mfa — the MFAVerifier is already wired
	// into the login gate above; the provider mounts verify/self-service/
	// admin-reset routes and users.mfa-reset. The auth-session repository
	// satisfies handler.SessionRevoker for disable/reset session invalidation.
	if plan.HasModule("admin.mfa") {
		providers = append(providers, mfamodule.New(a, mfaService, operations, authRepository, rateLimiters))
	}
	if plan.HasModule("admin.recycle-bin") {
		providers = append(providers, recyclebinmodule.New(a, recycleService, operations))
	}
	// S-14 (GOAL-019 D-002 §3): admin.wallet — accounts + immutable ledger +
	// reconciliation. Money-path mutations are gated by wallet.adjust; the
	// module never touches the manifest/profile semantics (content extension).
	if plan.HasModule("admin.wallet") {
		walletService := walletmodule.NewService(walletstore.NewRepository(st), st)
		walletJobs, err := walletmodule.NewJobService(walletService, jobRuntime.repository, jobRuntime.runner, operations)
		if err != nil {
			return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: walletmodule.ModuleID, Detail: fmt.Sprintf("register wallet jobs: %v", err)}
		}
		jobRuntime.enabled.Store(true)
		// W13 F-012 (GOAL-013 A-001) / VP-029: the by-owner HTTP surface only
		// opens USER accounts, so its existence gate checks the live user table
		// ONLY — a registered external subject id must never mint an
		// owner_type=user book with no admin.users row (VP-029 A-005 F-001,
		// workspace-029 GOAL-001). Subject existence stays inside
		// CreateAccount(owner_type=subject) and voucher Redeem (SubjectExists).
		walletOwnerExists := handler.OwnerExistsFunc(func(ownerID string) bool {
			_, err := authRepository.UserByID(ownerID)
			return err == nil
		})
		providers = append(providers, walletmodule.New(a, walletService, walletJobs, operations, walletOwnerExists, rateLimiters))
	}
	if plan.HasModule("admin.notifications") {
		providers = append(providers, notificationsmodule.New(a, authRepository))
		// F-04 best-effort system-event hooks (lock/disable/unlock/password).
		a.OnLockOpened = func(userID string) {
			handler.NotifyAccountEvent(authRepository, userID, "account.locked", time.Now().UTC())
		}
	}

	// VP-030 (GOAL-003 R2 / GOAL-004 R3 / F-001 / F-002): channel.telegram — Telegram channel runtime.
	// Assembled by plan enablement (custom profile or explicit app.modules). The
	// injected `tr` is THE process instance provided by newTelegramRuntime in the
	// Fx graph — never reconstructed here (F-001 / A-006: variadic removed).
	if plan.HasModule("channel.telegram") && tr != nil && tr.Webhook != nil {
		tgSettings := a.Middleware(telegraminternal.NewSettingsHandler(tr.Manager))
		providers = append(providers, telegrammodule.New(tr.Webhook, tgSettings))
	}
	providers = append(providers, extra...)
	set, err := kernel.RegisterContributions(context.Background(), plan, providers)
	if err != nil {
		return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "", Detail: fmt.Sprintf("register contributions: %v", err)}
	}
	// W4 P0-2: upload is a centrally-registered shared-capability endpoint (no
	// owning module provider), so its files.write permission is contributed
	// centrally here. Default grant is admin-only (PolicyAdmin); a delegated
	// role holding files.write can still upload (RBAC委派对称).
	set.Permissions = append(set.Permissions,
		kernel.PermissionContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "core.server-registration", Key: "files.write"},
			Permission:           "files.write",
			Resource:             "files",
			Action:               "write",
			PolicyID:             authsessiondata.PolicyAdmin,
			SystemDataVersion:    authsessiondata.SystemDataVersion,
		},
		kernel.PermissionContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "core.auth-session", Key: "service-credentials.read"},
			Permission:           "service-credentials.read",
			Resource:             "service-credentials",
			Action:               "read",
			PolicyID:             authsessiondata.PolicyAdmin,
			SystemDataVersion:    authsessiondata.SystemDataVersion,
		},
		kernel.PermissionContribution{
			ContributionIdentity: kernel.ContributionIdentity{ModuleID: "core.auth-session", Key: "service-credentials.write"},
			Permission:           "service-credentials.write",
			Resource:             "service-credentials",
			Action:               "write",
			PolicyID:             authsessiondata.PolicyAdmin,
			SystemDataVersion:    authsessiondata.SystemDataVersion,
		},
	)
	if err := authsessiondata.Reconcile(context.Background(), st, set.Permissions, set.Navigation); err != nil {
		_ = st.Close()
		return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("reconcile system data: %v", err)}
	}
	st.MarkSystemDataReady()
	for _, route := range set.Routes {
		full := route.Method + " " + route.Pattern
		// VP-015 R2: contributed routes carry their owning module_id in the
		// metrics labels (R1 D-001 §6); central registrations default to core.
		mux.Own(full, route.ModuleID)
		mux.Handle(full, route.Handler)
	}
	for _, route := range handler.ServiceCredentialRoutes(a, authRepository, operations, "core.auth-session") {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	// R6 C6.3: finalized page contributions own both metadata and document bytes;
	// the handler has no static document or owner fallback.
	handler.RegisterSchemas(mux, a, set.Pages)
	if plan.HasModule("core.manifest-route") {
		moduleFragments := make([]manifest.Fragment, 0, len(set.Fragments))
		for _, fragment := range set.Fragments {
			moduleFragments = append(moduleFragments, manifest.Fragment{ModuleID: fragment.ModuleID, Raw: fragment.JSON})
		}
		// GOAL-013 D-002 §4: the navigation order (config override or the
		// product-frozen default) reorders the published manifest slots; the
		// kernel sort only drives system-data menu_items ordering. The order is
		// normalized against the registered NodeIDs so an invalid override
		// falls back to the default here as well (manifest aggregation has no
		// knowledge of the kernel's default list).
		knownNodeIDs := make([]string, 0, len(set.Navigation))
		for _, n := range set.Navigation {
			knownNodeIDs = append(knownNodeIDs, n.NodeID)
		}
		navOrder := kernel.NormalizeNavigationOrder(plan.NavigationOrder, knownNodeIDs)
		data, err := manifest.ForModulesWithFragments(plan.IDs(), moduleFragments, navOrder)
		if err != nil {
			return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "core.manifest-route", Detail: err.Error()}
		}
		// W1 (GOAL-002 / workspace-010, D-003 §1/§2): homePageRef is derived from
		// the enabled set at assembly — dev.examples enabled -> overview; else the
		// first enabled admin functional page; else the first enabled page; else
		// omitted — then stamped into the published app block.
		data, err = manifest.StampHomePageRef(data, deriveHomePageRef(plan))
		if err != nil {
			return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "core.manifest-route", Detail: err.Error()}
		}
		if err := handler.RegisterManifest(mux, data); err != nil {
			return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "core.manifest-route", Detail: fmt.Sprintf("register manifest: %v", err)}
		}
		// Host/App 互操作（ADR-0035）：bootstrap document 与 manifest 同字节组装，
		// 声明的 manifest.sha256 与真实响应一致。
		if err := handler.RegisterBootstrapWithAvailability(mux, data, string(cfg.RuntimeMode)); err != nil {
			return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "core.manifest-route", Detail: fmt.Sprintf("register bootstrap: %v", err)}
		}
	}
	return mux.ServeMux, nil
}

// adminFunctionalOrder is the frozen home-page priority (D-003 §2): the first
// enabled admin.* functional module in declaration order becomes the home page.
// F-03 (GOAL-005 D-002 §6): admin.account appended at the tail — home stays
// users-first; account only becomes home when every earlier admin module is
// disabled (explicit, documented edge).
// F-01 (GOAL-003 D-002 §3): admin.dashboard inserted at the HEAD so the
// production home becomes the dashboard (必办-3 content-extension semantics).
var adminFunctionalOrder = []string{"admin.dashboard", "admin.users", "admin.roles", "admin.settings", "admin.activity", "admin.account"}

// deriveHomePageRef implements the D-003 §2 decision table:
//
//  1. dev.examples enabled              -> "overview"
//  2. else first enabled admin module   -> that module's first declared page
//  3. else first enabled module with a page contribution -> that page
//  4. else "" (omit homePageRef)
func deriveHomePageRef(plan kernel.Plan) string {
	if plan.HasModule("dev.examples") {
		return "overview"
	}
	byID := make(map[string]kernel.Module, len(plan.Modules))
	for _, m := range plan.Modules {
		byID[m.ID] = m
	}
	for _, moduleID := range adminFunctionalOrder {
		if module, ok := byID[moduleID]; ok && len(module.Contributions.Pages) > 0 {
			return module.Contributions.Pages[0]
		}
	}
	for _, module := range plan.Modules {
		if len(module.Contributions.Pages) > 0 {
			return module.Contributions.Pages[0]
		}
	}
	return ""
}

// newObjectStore builds THE shared kernel.ObjectStore instance (VP-014
// GOAL-004 D-001): local disk default rooted at db.path dir (or the
// explicit override), or the S3-compatible adapter when driver=s3. The
// second return is the optional readyz probe (HeadBucket on s3; nil for
// local so readyz semantics stay unchanged).
func newObjectStore(cfg *config.Config) (kernel.ObjectStore, func(context.Context) error, error) {
	if cfg.ObjectsDriver == "s3" {
		objStore, err := objectstore.NewS3(cfg.ObjectsS3Endpoint, cfg.ObjectsS3Region,
			cfg.ObjectsS3Bucket, cfg.ObjectsS3AccessKeyID, cfg.ObjectsS3SecretAccessKey,
			cfg.ObjectsS3UsePathStyle)
		if err != nil {
			return nil, nil, err
		}
		return objStore, objStore.Ping, nil
	}
	// Defense-in-depth (GOAL-003 A-002 N-005): Load already rejects unknown
	// drivers fail-closed; re-check here so a hand-built Config cannot silently
	// fall through to the local adapter.
	if cfg.ObjectsDriver != "local" && cfg.ObjectsDriver != "" {
		return nil, nil, fmt.Errorf("composition: unknown storage.objects.driver %q", cfg.ObjectsDriver)
	}
	root := cfg.ObjectsLocalRoot
	if strings.TrimSpace(root) == "" {
		root = filepath.Dir(cfg.DBPath)
	}
	return objectstore.NewLocal(root), nil, nil
}

// newRateLimiters builds THE shared kernel.RateLimiterProvider (VP-027 /
// workspace-027 GOAL-003 D-001, R2): the in-memory factory. The Fx container
// owns it for the process lifetime (combination-root single holder, seam doc
// §2.4); a future Redis-tier provider replaces it at composition when RT-Q05
// triggers. No probe: an in-process store has no external dependency.
func newRateLimiters() kernel.RateLimiterProvider {
	return ratelimit.NewProvider()
}

// newCache builds THE shared kernel.Cache instance (VP-026 / workspace-026
// GOAL-003 D-001): the in-memory provider with the configured bounded-entry
// budget. There is no readyz probe — an in-process store has no external
// dependency. Zero on a loader-bypassed (zero-value) Config means "use the
// load default" (mirrors the db rules and newObjectStore's empty-driver
// handling); a negative budget is a programming error and fails closed.
func newCache(cfg *config.Config) (kernel.Cache, error) {
	budget := cfg.CacheMaxEntries
	if budget < 0 {
		return nil, fmt.Errorf("composition: cache.max_entries must be positive (got %d)", cfg.CacheMaxEntries)
	}
	if budget == 0 {
		budget = config.DefaultCacheMaxEntries
	}
	return cache.NewMemory(budget)
}

// newEventBus builds THE shared kernel.EventBus instance (VP-028 / workspace-028
// GOAL-003 D-001): the in-memory provider with the configured per-subscription
// buffer size. There is no readyz probe — an in-process channel bus has no
// external dependency. <= 0 on a loader-bypassed Config or explicit YAML/env
// means "use the default" (mirrors the cache zero-value handling); the provider
// falls back to kernel.DefaultEventBusBuffer.
func newEventBus(cfg *config.Config, logger *slog.Logger) kernel.EventBus {
	buffer := cfg.EventBusBufferSize
	if buffer == 0 {
		buffer = config.DefaultEventBusBuffer
	}
	return eventbus.NewMemory(buffer, logger)
}

// TelegramRuntime holds the process-level Telegram ports and state (F-001).
type TelegramRuntime struct {
	Dispatcher kernel.TelegramDispatcher
	Sender     kernel.TelegramSender
	Manager    *telegraminternal.RuntimeManager
	Webhook    *telegraminternal.WebhookHandler
}

// newTelegramRuntime builds the shared TelegramRuntime for the process (F-001).
// The at-rest master key is resolved the same way as the mail channel
// (F-002 / A-006): operator-passphrase env (TELEGRAM_MASTER_KEY) or an
// auto-generated key file beside the database — never a source constant.
func newTelegramRuntime(plan kernel.Plan, cfg *config.Config, st kernel.Store, rateLimiters kernel.RateLimiterProvider) (*TelegramRuntime, error) {
	if plan.HasModule("channel.telegram") {
		masterKeyPath := cfg.TelegramMasterKeyPath
		if strings.TrimSpace(masterKeyPath) == "" {
			masterKeyPath = filepath.Join(filepath.Dir(cfg.DBPath), "telegram-master.key")
		}
		masterKey, err := mail.LoadOrCreateMasterKey(cfg.TelegramMasterKey, masterKeyPath)
		if err != nil {
			return nil, fmt.Errorf("composition: telegram master key: %w", err)
		}
		subStore := subject.NewStore(st)
		disp := telegraminternal.NewDispatcher()
		mockSender := telegraminternal.NewCaptureSender()
		rt, err := telegraminternal.NewRuntimeManagerWithSettings(cfg.TelegramBotToken, cfg.TelegramWebhookSecret, cfg.TelegramMode, cfg.TelegramWebhookPublicBaseURL, mockSender, masterKey, st)
		if err != nil {
			return nil, fmt.Errorf("composition: telegram runtime: %w", err)
		}
		sender := telegraminternal.NewHTTPSender(rt, nil, "")
		webhook := telegraminternal.NewWebhookHandler(telegraminternal.HandlerConfig{
			TokenGetter:  rt.GetToken,
			SecretGetter: rt.GetSecret,
			RateLimiters: rateLimiters,
			SubjectStore: subStore,
			Dispatcher:   disp,
			Sender:       sender,
		})
		return &TelegramRuntime{
			Dispatcher: disp,
			Sender:     sender,
			Manager:    rt,
			Webhook:    webhook,
		}, nil
	}
	return &TelegramRuntime{
		Dispatcher: telegraminternal.NewDisabledDispatcher(),
		Sender:     telegraminternal.NewDisabledSender(),
	}, nil
}

// ResolveTelegramPorts returns the process-level TelegramDispatcher and TelegramSender (D-002 §1 / F-001).
// It is a standalone helper for non-Fx consumers (tests, external harnesses);
// the Fx production graph injects the single *TelegramRuntime directly.
func ResolveTelegramPorts(plan kernel.Plan, cfg *config.Config, st kernel.Store) (kernel.TelegramDispatcher, kernel.TelegramSender, error) {
	tr, err := newTelegramRuntime(plan, cfg, st, newRateLimiters())
	if err != nil {
		return nil, nil, err
	}
	return tr.Dispatcher, tr.Sender, nil
}

// newMailRuntime builds THE kernel.MailSender for the process (VP-017 R7 /
// workspace-017 GOAL-008; Root D-007): a *mail.Switcher over the mail_config
// runtime row. The row seeds once from the file/env layer resolution
// (mock by default, explicit resend/smtp when configured — fail-closed on
// ambiguity or incomplete blocks via ResolveMailChannel). The second return
// is the optional readyz probe: nil except when the BOOT channel is SMTP
// (its ESMTP Ping extends readyz exactly as frozen in R4; Resend joins R8).
func newMailRuntime(cfg *config.Config, st kernel.Store, logger *slog.Logger) (*mail.Switcher, func(context.Context) error, error) {
	channel, err := cfg.ResolveMailChannel()
	if err != nil {
		return nil, nil, fmt.Errorf("composition: resolve mail.channel: %w", err)
	}
	masterKeyPath := cfg.MailMasterKeyPath
	if strings.TrimSpace(masterKeyPath) == "" {
		// W13 F-017 (GOAL-013 A-001): the historical default co-locates the
		// key with the data directory, so one backup/ snapshot leaks both the
		// encrypted channel secrets AND the key that unlocks them. Operators
		// can relocate it via mail.master_key_path / MAIL_MASTER_KEY_PATH.
		masterKeyPath = filepath.Join(filepath.Dir(cfg.DBPath), "mail-master.key")
	}
	masterKey, err := mail.LoadOrCreateMasterKey(cfg.MailConfigMasterKey, masterKeyPath)
	if err != nil {
		return nil, nil, err
	}
	switcher, err := mail.NewSwitcher(st, masterKey, mail.SeedConfig{
		Channel:       channel,
		MockRetention: mail.DefaultOutboxCap,
		ResendFrom:    cfg.MailResendFrom,
		ResendAPIKey:  cfg.MailResendAPIKey,
		SMTPHost:      cfg.MailSMTPHost,
		SMTPPort:      cfg.MailSMTPPort,
		SMTPUsername:  cfg.MailSMTPUsername,
		SMTPPassword:  cfg.MailSMTPPassword,
		SMTPFrom:      cfg.MailSMTPFrom,
	}, logger)
	if err != nil {
		return nil, nil, err
	}
	var probe func(context.Context) error
	switch channel {
	case config.MailChannelSMTP:
		if cfg.MailSMTPConfigured() {
			sender, err := mail.NewSMTP(mail.SMTPOptions{
				Host:     cfg.MailSMTPHost,
				Port:     cfg.MailSMTPPort,
				Username: cfg.MailSMTPUsername,
				Password: cfg.MailSMTPPassword,
				From:     cfg.MailSMTPFrom,
			})
			if err != nil {
				return nil, nil, fmt.Errorf("composition: invalid mail.smtp configuration: %w", err)
			}
			probe = sender.Ping
		}
	case config.MailChannelResend:
		// VP-017 R8 (GOAL-009): the explicitly configured production Resend
		// channel extends readyz with its availability probe, mirroring the
		// SMTP ESMTP Ping precedent.
		sender, err := mail.NewResend(mail.ResendOptions{
			APIKey: cfg.MailResendAPIKey,
			From:   cfg.MailResendFrom,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("composition: invalid mail.resend configuration: %w", err)
		}
		probe = sender.Ping
	}
	return switcher, probe, nil
}

func newServer(cfg *config.Config, mux *http.ServeMux, logger *slog.Logger) *http.Server {
	routes := handler.WithJSONRouteErrors(mux)
	return server.New(cfg, handler.WithOperationalGate(cfg, mux, routes), logger)
}

func registerLifecycle(lc fx.Lifecycle, srv *http.Server, st kernel.Store, logger *slog.Logger, cfg *config.Config, plan kernel.Plan, gate *readinessGate, jobs *jobRuntime, operations *operationlog.Repository, settingsRepository *settingsrepository.Repository, metrics *obs.Server, tracing *obs.Tracing, eventBusPort kernel.EventBus) {
	var listener net.Listener
	var stopRetention func()
	runtime := kernel.NewRuntime(withLifecycleHooks(plan, st, logger, func() bool { return listener != nil }))
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			ln, err := net.Listen("tcp", srv.Addr)
			if err != nil {
				_ = st.Close()
				return &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.server-registration", Detail: fmt.Sprintf("listen %s: %v", srv.Addr, err)}
			}
			listener = ln
			if err := runtime.Start(ctx); err != nil {
				_ = ln.Close()
				listener = nil
				_ = st.Close()
				return err
			}
			if err := runtime.Ready(ctx); err != nil {
				_ = ln.Close()
				listener = nil
				_ = st.Close()
				return err
			}
			if err := jobs.Start(); err != nil {
				_ = runtime.Stop(ctx)
				_ = ln.Close()
				listener = nil
				_ = st.Close()
				return &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: walletmodule.ModuleID, Detail: fmt.Sprintf("start job runner: %v", err)}
			}
			stopRetention = operationlog.StartRetentionSweep(operations, func() (operationlog.RetentionPolicy, error) {
				settings, err := settingsRepository.GetSiteSettings()
				if err != nil {
					return operationlog.RetentionPolicy{}, err
				}
				return operationlog.RetentionPolicy{
					Days:   settings.OperationLogRetentionDays,
					Action: settings.OperationLogExpirationAction,
				}, nil
			}, time.Hour, logger)
			// VP-015 R2 (GOAL-003 D-001 §3): the dedicated metrics listener is a
			// bypass face — it never gates readiness, but an explicitly enabled
			// listener that cannot bind fails startup (fail-closed).
			if err := metrics.Start(ctx); err != nil {
				stopRetention()
				_ = jobs.Stop(ctx)
				_ = runtime.Stop(ctx)
				_ = ln.Close()
				listener = nil
				_ = st.Close()
				return &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.server-registration", Detail: fmt.Sprintf("start metrics listener: %v", err)}
			}
			// R5 real readiness: only after every module Start + Ready succeeds
			// does /readyz report ready.
			gate.setReady()
			logger.Info("server starting",
				"addr", cfg.HTTPAddr,
				"profile", cfg.ProfileName,
				"modules", plan.IDs(),
				"capabilities", plan.Capabilities,
			)
			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					// D-001 P2 fail-closed: an unexpected Serve failure (e.g.
					// listen/accept loop breaking after startup) must take the
					// process down so the compose restart policy brings it back,
					// instead of leaving a half-dead instance serving nothing.
					logger.Error("server failed; exiting", "err", err)
					os.Exit(1)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			shutdownErr := srv.Shutdown(ctx)
			if listener != nil {
				_ = listener.Close()
			}
			if stopRetention != nil {
				stopRetention()
			}
			metricsErr := metrics.Stop(ctx)
			jobsErr := jobs.Stop(ctx)
			// VP-028 / workspace-028 GOAL-003 D-001 (R2 / D-002 §5): Stop drains
			// buffered events, waits in-flight handlers against ctx, then rejects
			// further Publish/Register/Subscribe.
			eventBusErr := eventBusPort.Stop(ctx)
			runtimeErr := runtime.Stop(ctx)
			closeErr := st.Close()
			// GOAL-004 D-001 §6: shutdown flushes pending spans through the
			// OTLP exporter before the process exits.
			tracingErr := tracing.Shutdown(ctx)
			return errors.Join(shutdownErr, metricsErr, jobsErr, eventBusErr, runtimeErr, closeErr, tracingErr)
		},
	})
}

func withLifecycleHooks(plan kernel.Plan, st kernel.Store, logger *slog.Logger, listenerReady func() bool) kernel.Plan {
	for i := range plan.Modules {
		moduleID := plan.Modules[i].ID
		plan.Modules[i].Hooks = kernel.Hooks{
			Start: func(ctx context.Context) error {
				switch moduleID {
				case "core.server-registration":
					if !listenerReady() {
						return &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: moduleID, Detail: "HTTP listener is not ready"}
					}
				case "core.auth-session":
					if err := st.Ping(ctx); err != nil {
						return &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: moduleID, Detail: fmt.Sprintf("store is not available: %v", err)}
					}
				}
				logger.Debug("module started", "module_id", moduleID)
				return nil
			},
			Ready: func(ctx context.Context) error {
				if moduleID == "core.auth-session" {
					if err := st.Ping(ctx); err != nil {
						return &kernel.Error{Code: kernel.CodeLifecycleReadyFailed, ModuleID: moduleID, Detail: fmt.Sprintf("store readiness failed: %v", err)}
					}
					if err := st.SystemDataReady(); err != nil {
						return &kernel.Error{Code: kernel.CodeLifecycleReadyFailed, ModuleID: moduleID, Detail: fmt.Sprintf("system-data readiness failed: %v", err)}
					}
				}
				return nil
			},
			Stop: func(context.Context) error {
				logger.Debug("module stopped", "module_id", moduleID)
				return nil
			},
		}
	}
	return plan
}
