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

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/manifest"
	accountmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/account"
	activitymodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/activity"
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
	return registry.Resolve(modules)
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
) (*http.ServeMux, error) {
	return newMuxWithExtraProviders(cfg, a, st, authRepository, operations, settingsRepository, plan, gate, nil)
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
	extra []kernel.Provider,
) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	handler.RegisterWithReadiness(mux, a, st, operations, plan, gate.Ready)
	// I-PROTO-FULL-001 D-UPLOAD: server-side upload contract (07 §7.2).
	handler.RegisterUpload(mux, a, filepath.Join(filepath.Dir(cfg.DBPath), "uploads"))
	// R4 C3.3: admin.users / admin.roles HTTP surface comes from the module
	// kernel.Provider contract (freeze package §7 step 3). Core auth/accounts/
	// health/schema stay central; settings/activity migrate in C4.
	// W1 (GOAL-002 / workspace-010): core.schema-render and the optional
	// dev.examples module are both assembled by plan enablement (D-003 §3) —
	// production profiles default to no examples surface.
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
		providers = append(providers, settingsmodule.New(a, settingsRepository, operations))
	} else {
		// Public bootstrap must work on mvp (and any profile without the
		// settings edit module). Edit/list/patch/reset stay admin.settings-only.
		handler.RegisterPublicBranding(mux, settingsRepository)
	}
	if plan.HasModule("admin.activity") {
		providers = append(providers, activitymodule.New(a, operations))
	}
	if plan.HasModule("admin.account") {
		providers = append(providers, accountmodule.New(a, authRepository, operations))
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
		data, err := manifest.ForModulesWithFragments(plan.IDs(), moduleFragments)
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
var adminFunctionalOrder = []string{"admin.users", "admin.roles", "admin.settings", "admin.activity", "admin.account"}

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
	return server.New(cfg, mux, logger)
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