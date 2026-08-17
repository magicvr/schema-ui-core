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
	"sync/atomic"
	"time"

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/manifest"
	accountmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/account"
	activitymodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/activity"
	datatransfermodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datatransfer"
	dashboardmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/dashboard"
	datadictionarymodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary"
	systemmonitoringmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/systemmonitoring"
	scheduledtasksmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks"
	scheduledtasksstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
	logincaptchamodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/logincaptcha"
	datapermissionmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datapermission"
	datapermissionstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datapermission/store"
	mfamodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/mfa"
	mfastore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/mfa/store"
	logincaptchastore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/logincaptcha/store"
	recyclebinmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin"
	recyclestore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/recyclebin/store"
	walletmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet"
	walletstore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/wallet/store"
	datadictionarystore2 "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
	tasksstore2 "github.com/magicvr/schema-ui-core/apps/api/internal/modules/scheduledtasks/store"
	datadictionarystore "github.com/magicvr/schema-ui-core/apps/api/internal/modules/datadictionary/store"
	filelibrarymodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/filelibrary"
	notificationsmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/notifications"
	authsession "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession"
	authsessiondata "github.com/magicvr/schema-ui-core/apps/api/internal/modules/authsession/systemdata"
	compiledmodules "github.com/magicvr/schema-ui-core/apps/api/internal/modules/compiled"
	devexamplesmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/dev/examples"
	"github.com/magicvr/schema-ui-core/apps/api/internal/modules/operationlog"
	rolesmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/roles"
	schemarendermodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/schemarender"
	settingsmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings"
	settingsrepository "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings/repository"
	usersmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/users"
	"github.com/magicvr/schema-ui-core/apps/api/internal/server"
	"github.com/magicvr/schema-ui-core/apps/api/internal/store"
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
	plan, err := ResolvePlan(cfg)
	if err != nil {
		return nil, err
	}
	if logger == nil {
		logger = slog.Default()
	}
	return fx.New(
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
			newMux,
			newServer,
		),
		fx.Invoke(registerLifecycle),
	), nil
}

// readinessGate reports whether the module graph Start+Ready succeeded (R5 real
// readiness). readyz consults it; registerLifecycle flips it after Ready.
type readinessGate struct {
	ready atomic.Bool
}

func (g *readinessGate) Ready() bool { return g.ready.Load() }
func (g *readinessGate) setReady()   { g.ready.Store(true) }

func openStore(cfg *config.Config, seedHash seedPasswordHash) (*store.Store, error) {
	// Persistence is compiled-global rather than profile-gated. The static
	// registry collects module-owned descriptors before store startup; store
	// only validates and executes the resulting catalog.
	catalog, err := compiledmodules.PersistenceCatalog()
	if err != nil {
		return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.persistence", Detail: fmt.Sprintf("collect persistence: %v", err)}
	}
	st, err := store.OpenWithCatalog(cfg.DBPath, catalog)
	if err != nil {
		return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("open store: %v", err)}
	}
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
	return st, nil
}

func newAuthSessionRepository(st *store.Store) *authsession.Repository {
	return authsession.NewRepository(st)
}

func newOperationLogRepository(st *store.Store) *operationlog.Repository {
	return operationlog.NewRepository(st)
}

func newSettingsRepository(st *store.Store) *settingsrepository.Repository {
	return settingsrepository.New(st)
}

func newAuthenticator(cfg *config.Config, secret jwtSecret, repository *authsession.Repository) *auth.Authenticator {
	return auth.NewWithRepository([]byte(secret), cfg.AuthAccessTTL, cfg.AuthRefreshTTL, repository, cfg.AuthDevSessionEnabled)
}

func newMux(
	cfg *config.Config,
	a *auth.Authenticator,
	st *store.Store,
	authRepository *authsession.Repository,
	operations *operationlog.Repository,
	settingsRepository *settingsrepository.Repository,
	plan kernel.Plan,
	gate *readinessGate,
	secret jwtSecret,
) (*http.ServeMux, error) {
	return newMuxWithExtraProviders(cfg, a, st, authRepository, operations, settingsRepository, plan, gate, secret, nil)
}

// newMuxWithExtraProviders is the composition-root assembly seam used by the S2
// access drill: it appends test-only providers (e.g. a probe module) to the
// compiled one-party set before registration, proving that a new standard module
// surfaces through the Provider contract without any handler/Web central business
// registration change. Production call sites use newMux (extra == nil).
func newMuxWithExtraProviders(
	cfg *config.Config,
	a *auth.Authenticator,
	st *store.Store,
	authRepository *authsession.Repository,
	operations *operationlog.Repository,
	settingsRepository *settingsrepository.Repository,
	plan kernel.Plan,
	gate *readinessGate,
	secret jwtSecret,
	extra []kernel.Provider,
) (*http.ServeMux, error) {
	mux := http.NewServeMux()
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
		mfaService = mfamodule.NewService(mfastore.NewRepository(st), []byte(secret))
		mfaVerifier = mfaService
	}
	handler.RegisterWithMFA(mux, a, st, operations, plan, gate.Ready, []handler.CaptchaVerifier{captchaVerifier}, mfaVerifier)
	// I-PROTO-FULL-001 D-UPLOAD: server-side upload contract (07 §7.2). The
	// uploads directory is shared with admin.data-transfer (F-02 import reads
	// uploaded CSV files by id).
	uploadDir := filepath.Join(filepath.Dir(cfg.DBPath), "uploads")
	handler.RegisterUpload(mux, a, uploadDir,
		handler.WithAllowedTypes(cfg.UploadAllowedTypes),
		handler.WithUserLimits(cfg.UploadMaxFilesPerUser, cfg.UploadMaxBytesPerUser),
	)
	// W9 (GOAL-010): dedicated brand-assets store — NOT the shared upload
	// store (owner-gated reads) and NOT admin.file-library. Brand icons must
	// be publicly readable (login page / shell load pre-auth); every stored
	// object is a server-side re-encoded raster (never raw upload bytes).
	brandAssets := handler.NewBrandingAssetStore(
		filepath.Join(filepath.Dir(cfg.DBPath), "brand-assets"),
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
	// brand assets, dedicated directory, 256px longest-edge default. No startup
	// GC: the referenced set lives in the users table and the replace/clear
	// paths delete previous files best-effort (see account_avatar.go).
	avatarAssets := handler.NewAvatarAssetStore(
		filepath.Join(filepath.Dir(cfg.DBPath), "avatars"),
		handler.BrandingAssetsOptions{},
	)
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
		recycleService = recyclebinmodule.NewService(recyclestore.NewRepository(st), datadictionarystore2.NewRepository(st), tasksstore2.NewRepository(st))
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
		providers = append(providers, usersmodule.New(a, authRepository, operations))
	}
	if plan.HasModule("admin.roles") {
		providers = append(providers, rolesmodule.New(a, authRepository, operations))
	}
	if plan.HasModule("admin.settings") {
		providers = append(providers, settingsmodule.New(a, settingsRepository, operations, brandAssets))
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
		providers = append(providers, accountmodule.New(a, authRepository, operations, avatarAssets))
	}
	if plan.HasModule("admin.data-transfer") {
		providers = append(providers, datatransfermodule.New(a, authRepository, operations, uploadDir))
	}
	if plan.HasModule("admin.dashboard") {
		providers = append(providers, dashboardmodule.New())
	}
	if plan.HasModule("admin.file-library") {
		providers = append(providers, filelibrarymodule.New(a, operations, uploadDir))
	}
	if plan.HasModule("admin.data-dictionary") {
		providers = append(providers, datadictionarymodule.New(a, datadictionarystore.NewRepository(st), operations, trash))
	}
	if plan.HasModule("admin.system-monitoring") {
		providers = append(providers, systemmonitoringmodule.New(a, st, plan, gate.Ready, cfg.DBPath, time.Now(), operations))
	}
	if plan.HasModule("admin.scheduled-tasks") {
		providers = append(providers, scheduledtasksmodule.New(a, scheduledtasksstore.NewRepository(st), operations, trash))
	}
	if plan.HasModule("admin.login-captcha") {
		providers = append(providers, logincaptchamodule.New(a, captchaService, operations))
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
		providers = append(providers, mfamodule.New(a, mfaService, operations, authRepository))
	}
	if plan.HasModule("admin.recycle-bin") {
		providers = append(providers, recyclebinmodule.New(a, recycleService, operations))
	}
	// S-14 (GOAL-019 D-002 §3): admin.wallet — accounts + immutable ledger +
	// reconciliation. Money-path mutations are gated by wallet.adjust; the
	// module never touches the manifest/profile semantics (content extension).
	if plan.HasModule("admin.wallet") {
		providers = append(providers, walletmodule.New(a, walletmodule.NewService(walletstore.NewRepository(st)), operations))
	}
	if plan.HasModule("admin.notifications") {
		providers = append(providers, notificationsmodule.New(a, authRepository))
		// F-04 best-effort system-event hooks (lock/disable/unlock/password).
		a.OnLockOpened = func(userID string) {
			handler.NotifyAccountEvent(authRepository, userID, "account.locked", time.Now().UTC())
		}
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
	set.Permissions = append(set.Permissions, kernel.PermissionContribution{
		ContributionIdentity: kernel.ContributionIdentity{ModuleID: "core.server-registration", Key: "files.write"},
		Permission:           "files.write",
		Resource:             "files",
		Action:               "write",
		PolicyID:             authsessiondata.PolicyAdmin,
		SystemDataVersion:    authsessiondata.SystemDataVersion,
	})
	if err := authsessiondata.Reconcile(context.Background(), st, set.Permissions, set.Navigation); err != nil {
		_ = st.Close()
		return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("reconcile system data: %v", err)}
	}
	st.MarkSystemDataReady()
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	// R6 C6.3: finalized page contributions own both metadata and document bytes;
	// the handler has no static document or owner fallback.
	handler.RegisterSchemas(mux, set.Pages)
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
		if err := handler.RegisterBootstrap(mux, data); err != nil {
			return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "core.manifest-route", Detail: fmt.Sprintf("register bootstrap: %v", err)}
		}
	}
	return mux, nil
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
//	1. dev.examples enabled              -> "overview"
//	2. else first enabled admin module   -> that module's first declared page
//	3. else first enabled module with a page contribution -> that page
//	4. else "" (omit homePageRef)
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

func newServer(cfg *config.Config, mux *http.ServeMux, logger *slog.Logger) *http.Server {
	return server.New(cfg, handler.WithJSONRouteErrors(mux), logger)
}

func registerLifecycle(lc fx.Lifecycle, srv *http.Server, st *store.Store, logger *slog.Logger, cfg *config.Config, plan kernel.Plan, gate *readinessGate) {
	var listener net.Listener
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
			runtimeErr := runtime.Stop(ctx)
			closeErr := st.Close()
			return errors.Join(shutdownErr, runtimeErr, closeErr)
		},
	})
}

func withLifecycleHooks(plan kernel.Plan, st *store.Store, logger *slog.Logger, listenerReady func() bool) kernel.Plan {
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