package composition

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"

	"go.uber.org/fx"

	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/internal/manifest"
	activitymodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/activity"
	rolesmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/roles"
	settingsmodule "github.com/magicvr/schema-ui-core/apps/api/internal/modules/settings"
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
			openStore,
			newAuthenticator,
			newMux,
			newServer,
		),
		fx.Invoke(registerLifecycle),
	), nil
}

func openStore(cfg *config.Config, seedHash seedPasswordHash) (*store.Store, error) {
	st, err := store.Open(cfg.DBPath, "admin", string(seedHash), true)
	if err != nil {
		return nil, &kernel.Error{Code: kernel.CodeLifecycleStartFailed, ModuleID: "core.auth-session", Detail: fmt.Sprintf("open store: %v", err)}
	}
	return st, nil
}

func newAuthenticator(cfg *config.Config, secret jwtSecret, st *store.Store) *auth.Authenticator {
	return auth.New([]byte(secret), cfg.AuthAccessTTL, cfg.AuthRefreshTTL, st, cfg.AuthDevSessionEnabled)
}

func newMux(a *auth.Authenticator, st *store.Store, plan kernel.Plan) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	handler.Register(mux, a, st, plan)
	// R4 C3.3: admin.users / admin.roles HTTP surface comes from the module
	// kernel.Provider contract (freeze package §7 step 3). Core auth/accounts/
	// health/schema stay central; settings/activity migrate in C4.
	var providers []kernel.Provider
	if plan.HasModule("admin.users") {
		providers = append(providers, usersmodule.New(a, st))
	}
	if plan.HasModule("admin.roles") {
		providers = append(providers, rolesmodule.New(a, st))
	}
	set, err := kernel.RegisterContributions(context.Background(), plan, providers)
	if err != nil {
		return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "admin.users", Detail: fmt.Sprintf("register contributions: %v", err)}
	}
	for _, route := range set.Routes {
		mux.Handle(route.Method+" "+route.Pattern, route.Handler)
	}
	if plan.HasModule("admin.settings") {
		settingsmodule.Register(mux, a, st)
	}
	if plan.HasModule("admin.activity") {
		activitymodule.Register(mux, a, st)
	}
	if plan.HasModule("core.manifest-route") {
		moduleFragments := make([]manifest.Fragment, 0, len(set.Fragments))
		for _, fragment := range set.Fragments {
			moduleFragments = append(moduleFragments, manifest.Fragment{ModuleID: fragment.ModuleID, Raw: fragment.JSON})
		}
		data, err := manifest.ForModulesWithFragments(plan.IDs(), moduleFragments)
		if err != nil {
			return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "core.manifest-route", Detail: err.Error()}
		}
		if err := handler.RegisterManifest(mux, data); err != nil {
			return nil, &kernel.Error{Code: kernel.CodeModuleInvalid, ModuleID: "core.manifest-route", Detail: fmt.Sprintf("register manifest: %v", err)}
		}
	}
	return mux, nil
}

func newServer(cfg *config.Config, mux *http.ServeMux, logger *slog.Logger) *http.Server {
	return server.New(cfg, mux, logger)
}

func registerLifecycle(lc fx.Lifecycle, srv *http.Server, st *store.Store, logger *slog.Logger, cfg *config.Config, plan kernel.Plan) {
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
				stopErr := runtime.Stop(ctx)
				_ = ln.Close()
				listener = nil
				_ = st.Close()
				return errors.Join(err, stopErr)
			}
			logger.Info("server starting",
				"addr", cfg.HTTPAddr,
				"profile", cfg.ProfileName,
				"modules", plan.IDs(),
				"capabilities", plan.Capabilities,
			)
			go func() {
				if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
					logger.Error("server failed", "err", err)
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
