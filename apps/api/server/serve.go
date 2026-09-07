package server

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/assembly"
	"github.com/magicvr/schema-ui-core/apps/api/internal/auth"
	"github.com/magicvr/schema-ui-core/apps/api/internal/handler"
	"github.com/magicvr/schema-ui-core/apps/api/internal/mail"
	"github.com/magicvr/schema-ui-core/apps/api/internal/manifest"
	"github.com/magicvr/schema-ui-core/apps/api/internal/obs"
	"github.com/magicvr/schema-ui-core/apps/api/internal/ratelimit"
	"github.com/magicvr/schema-ui-core/apps/api/internal/requestid"
	"github.com/magicvr/schema-ui-core/apps/api/kernel"
	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession"
	"github.com/magicvr/schema-ui-core/apps/api/modules/authsession/systemdata"
	"github.com/magicvr/schema-ui-core/apps/api/modules/compiled"
	"github.com/magicvr/schema-ui-core/apps/api/modules/operationlog"
	"github.com/magicvr/schema-ui-core/apps/api/modules/users"
)

const bcryptCost = 10

// standardModules 是下游 serve 面的标准组合（D-001 冻结）：内核四核心 +
// operationlog + 一个标准管理模块（users）。与主仓 admin Profile 的差异
// （mfa/captcha/wallet/settings/jobs/metrics/... 不装配）为有界下游基线。
var standardModules = []string{
	"core.server-registration",
	"core.auth-session",
	"core.manifest-route",
	"core.navigation-capability",
	"core.schema-render",
	"core.operationlog",
	"admin.users",
}

// Options 装配 serve 面所需输入。字段构造全部经公开模块函数 + 类型推断
// （同 assembly 契约承诺：下游无需命名 internal 类型）。
type Options struct {
	Config *Config
	// Logger 可选；nil = slog.Default()。
	Logger *slog.Logger
	// Store 可选覆盖（测试注入）；nil = 依 Config 打开（compiled 迁移台账自动 apply）。
	Store kernel.Store
	// ExtraProviders 追加的下游 Provider（如后续标准模块），叠加到标准组合。
	ExtraProviders []kernel.Provider
}

// Serve 运行 serve 面直到 SIGINT/SIGTERM。RT-D02 §1 全序停机完成后返回 nil
// （调用方进程 exit 0）；停机错误/预算耗尽返回错误（调用方 exit 1）。
func Serve(opts Options) error {
	if opts.Config == nil {
		return errors.New("server: config is required")
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	_, err := Run(context.Background(), opts, signals)
	return err
}

type readinessGate struct {
	ready atomic.Bool
}

func (g *readinessGate) Ready() bool { return g.ready.Load() }

func (g *readinessGate) setReady() { g.ready.Store(true) }

// Run 运行 serve 面生命周期（可注入 ctx / signal 通道以便测试与嵌入）。
// 返回监听地址（"127.0.0.1:port"）。signals 为 nil 时仅 ctx 取消触发停机。
// 停机语义（RT-D02）：signal/ctx → shutdown.starting → http.Server.Shutdown
// （预算 = Config.ShutdownTimeout）→ kernel runtime 逆序 Stop → Store Close →
// shutdown.complete / 错误（exit 语义交调用方）。
func Run(ctx context.Context, opts Options, signals <-chan os.Signal) (string, error) {
	cfg := opts.Config
	if cfg == nil {
		return "", errors.New("server: config is required")
	}
	logger := opts.Logger
	if logger == nil {
		logger = slog.Default()
	}

	// 1. Store（可选覆盖；否则依 Config 打开，compiled 迁移台账自动 apply）
	//    并完成 admin 种子引导（同主仓 composition.openStore 语义）。
	st := opts.Store
	if st == nil {
		catalog, err := compiled.PersistenceCatalog()
		if err != nil {
			return "", fmt.Errorf("server: catalog: %w", err)
		}
		st, err = assembly.OpenStore(ctx, kernel.Dialect(cfg.DBDialect), cfg.DBPath, cfg.DBDSN, catalog)
		if err != nil {
			return "", fmt.Errorf("server: store: %w", err)
		}
		if err := bootstrapAdmin(ctx, st, cfg); err != nil {
			_ = st.Close()
			return "", err
		}
	}

	// 2. 仓库 / 认证 / 邮件（下游形态：站内 outbox sink）。
	repo := authsession.NewRepository(st)
	ops := operationlog.NewRepository(st)
	mailer := mail.NewOutboxSink(st, mail.DefaultOutboxCap)
	authn := auth.NewWithRepositoryAndPrevious([]byte(resolveSecret(cfg)), nil,
		cfg.AuthAccessTTL, cfg.AuthRefreshTTL, repo, false)

	// 3. 标准下游组合（D-001）：registry 校验后固定装配。
	registry, err := kernel.NewRegistry(kernel.BuiltinModules())
	if err != nil {
		return "", fmt.Errorf("server: registry: %w", err)
	}
	plan, err := registry.Resolve(standardModules)
	if err != nil {
		return "", fmt.Errorf("server: resolve standard plan: %w", err)
	}

	// 4. 模块 Provider：users + 下游追加；贡献集校验后 Reconcile 系统数据。
	providers := []kernel.Provider{users.New(authn, repo, ops, mailer, cfg.PublicBaseURL)}
	providers = append(providers, opts.ExtraProviders...)
	set, err := kernel.RegisterContributions(ctx, plan, providers)
	if err != nil {
		_ = st.Close()
		return "", fmt.Errorf("server: register contributions: %w", err)
	}
	if err := systemdata.Reconcile(ctx, st, set.Permissions, set.Navigation); err != nil {
		_ = st.Close()
		return "", fmt.Errorf("server: reconcile system data: %w", err)
	}
	st.MarkSystemDataReady()

	// 5. mux + 中央面（healthz/readyz/登录、schema 文档、manifest/bootstrap、
	//    模块贡献路由）。captcha/mfa 未装配 → nil verifier（有界基线）。
	gate := &readinessGate{}
	mux := obs.NewInstrumentedMux(nil)
	if err := handler.SetTrustedProxyCIDRs(cfg.TrustedProxies); err != nil {
		_ = st.Close()
		return "", fmt.Errorf("server: trusted proxies: %w", err)
	}
	handler.RegisterWithMFAProbes(mux, authn, st, ops, plan, gate.Ready, ratelimit.NewProvider(), nil, nil)
	handler.RegisterSchemas(mux, authn, set.Pages)
	if plan.HasModule("core.manifest-route") {
		moduleFragments := make([]manifest.Fragment, 0, len(set.Fragments))
		for _, fragment := range set.Fragments {
			moduleFragments = append(moduleFragments, manifest.Fragment{ModuleID: fragment.ModuleID, Raw: fragment.JSON})
		}
		knownNodeIDs := make([]string, 0, len(set.Navigation))
		for _, n := range set.Navigation {
			knownNodeIDs = append(knownNodeIDs, n.NodeID)
		}
		navOrder := kernel.NormalizeNavigationOrder(plan.NavigationOrder, knownNodeIDs)
		data, err := manifest.ForModulesWithFragments(plan.IDs(), moduleFragments, navOrder)
		if err != nil {
			_ = st.Close()
			return "", fmt.Errorf("server: build manifest: %w", err)
		}
		data, err = manifest.StampHomePageRef(data, deriveHomePageRef(plan))
		if err != nil {
			_ = st.Close()
			return "", fmt.Errorf("server: stamp manifest home page: %w", err)
		}
		if err := handler.RegisterManifest(mux, data); err != nil {
			_ = st.Close()
			return "", fmt.Errorf("server: register manifest: %w", err)
		}
		if err := handler.RegisterBootstrapWithAvailability(mux, data, "normal"); err != nil {
			_ = st.Close()
			return "", fmt.Errorf("server: register bootstrap: %w", err)
		}
	}
	for _, route := range set.Routes {
		full := route.Method + " " + route.Pattern
		mux.Own(full, route.ModuleID)
		mux.Handle(full, route.Handler)
	}

	// 6. http.Server（timeouts + request-id + nosniff/CORS，镜像 internal/server）。
	handler2 := handler.WithJSONRouteErrors(mux.ServeMux)
	headerTimeout := cfg.ReadTimeout
	if headerTimeout <= 0 {
		headerTimeout = 5 * time.Second
	}
	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           requestid.Middleware(wrapSecurity(cfg, handler2)),
		ReadTimeout:       cfg.ReadTimeout,
		ReadHeaderTimeout: headerTimeout,
		WriteTimeout:      cfg.WriteTimeout,
		IdleTimeout:       cfg.IdleTimeout,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	ln, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		_ = st.Close()
		return "", fmt.Errorf("server: listen %s: %w", srv.Addr, err)
	}

	runtime := kernel.NewRuntime(plan)
	startCtx, startCancel := context.WithTimeout(context.Background(), 15*time.Second)
	if err := runtime.Start(startCtx); err != nil {
		startCancel()
		_ = ln.Close()
		_ = st.Close()
		return "", err
	}
	if err := runtime.Ready(startCtx); err != nil {
		startCancel()
		_ = runtime.Stop(startCtx)
		_ = ln.Close()
		_ = st.Close()
		return "", err
	}
	startCancel()
	gate.setReady()
	logger.Info("server starting",
		"addr", cfg.HTTPAddr,
		"profile", cfg.ProfileName,
		"modules", plan.IDs(),
	)
	go func() {
		if err := srv.Serve(ln); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("server failed; exiting", "err", err)
			os.Exit(1)
		}
	}()

	// 7. RT-D02 §1 停机全序（signal/ctx → Shutdown → runtime → store）。
	sigName := "context"
	select {
	case <-ctx.Done():
	case s, ok := <-signals:
		if ok {
			sigName = s.String()
		} else {
			sigName = "signal-unknown"
		}
	}
	logger.Info("shutdown.starting", "signal", sigName)
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.ShutdownTimeout)
	defer shutdownCancel()
	shutdownErr := srv.Shutdown(shutdownCtx)
	if ln != nil {
		_ = ln.Close()
	}
	runtimeErr := runtime.Stop(shutdownCtx)
	closeErr := st.Close()
	joined := errors.Join(shutdownErr, runtimeErr, closeErr)
	if joined != nil {
		if shutdownCtx.Err() != nil {
			logger.Error("shutdown.timeout", "err", joined)
		} else {
			logger.Error("shutdown.error", "err", joined)
		}
		return ln.Addr().String(), joined
	}
	logger.Info("shutdown.complete")
	return ln.Addr().String(), nil
}

// bootstrapAdmin 在 needs-bootstrap 时种入 admin 用户（dev 缺省 admin/admin；
// 非 dev 已由 Config.validate fail-closed 强制密码，且种子必须满足冻结策略）。
func bootstrapAdmin(ctx context.Context, st kernel.Store, cfg *Config) error {
	needs, err := systemdata.NeedsBootstrap(ctx, st)
	if err != nil {
		return fmt.Errorf("server: check bootstrap: %w", err)
	}
	if !needs {
		return nil
	}
	seed := cfg.AdminInitialPassword
	if seed == "" {
		seed = "admin"
	}
	// W15 F-003: production bootstrap enforces the frozen 8–72 byte policy
	// before hashing (dev keeps the documented "admin" fallback).
	if cfg.AppEnv != "development" {
		if err := authsession.ValidateSeedPassword(seed); err != nil {
			return fmt.Errorf("server: bootstrap seed password: %w", err)
		}
	}
	hash, err := auth.HashPassword(seed, bcryptCost)
	if err != nil {
		return fmt.Errorf("server: hash seed password: %w", err)
	}
	if err := systemdata.Bootstrap(ctx, st, "admin", hash); err != nil {
		return fmt.Errorf("server: bootstrap auth data: %w", err)
	}
	return nil
}

// resolveSecret 返回签名密钥；非 dev 缺空已由 validate fail-closed（dev 缺省
// 与主仓 cmd/server 一致的开发密钥）。
func resolveSecret(cfg *Config) string {
	if cfg.AuthJWTSecret != "" {
		return cfg.AuthJWTSecret
	}
	return "dev-only-insecure-jwt-secret-change-me"
}

// deriveHomePageRef 镜像主仓 composition 决策表（D-003 §2）：
// 首个启用的 admin.* 功能模块的首个页面；无则首个带页面贡献的模块；再否则 ""。
func deriveHomePageRef(plan kernel.Plan) string {
	byID := make(map[string]kernel.Module, len(plan.Modules))
	for _, m := range plan.Modules {
		byID[m.ID] = m
	}
	adminFunctionalOrder := []string{"admin.dashboard", "admin.users", "admin.roles", "admin.settings", "admin.activity", "admin.account"}
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

// wrapSecurity applies security headers and CORS policy (W16 F-002 hardening:
// origin validation, credential policy, restricted methods/headers, preflight
// cache). Rejects null or malformed origins, logs denied requests.
func wrapSecurity(cfg *Config, next http.Handler) http.Handler {
	allow := map[string]struct{}{}
	for _, origin := range cfg.CORSOrigins {
		allow[origin] = struct{}{}
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		origin := r.Header.Get("Origin")
		
		// W16 F-002: validate origin before reflecting it in ACAO header
		if origin != "" {
			// Reject null origin (often from sandboxed iframe or data: scheme)
			if origin == "null" {
				// Don't set CORS headers; let browser deny
				next.ServeHTTP(w, r)
				return
			}
			
			// Validate origin is well-formed URL
			if !isValidOrigin(origin) {
				// Malformed origin: deny silently (no CORS headers)
				next.ServeHTTP(w, r)
				return
			}
			
			// Check whitelist
			if _, ok := allow[origin]; ok {
				// Allowed origin: set CORS headers
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Vary", "Origin")
				// W16 F-002: restrict headers to minimum required set
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type, Accept-Language, X-Refresh-Token")
				// W16 F-002: restrict methods to actually used ones
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
				// W16 F-002: credentials policy (true because we use Authorization header)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				// W16 F-002: preflight cache (24 hours)
				w.Header().Set("Access-Control-Max-Age", "86400")
				
				if r.Method == http.MethodOptions {
					w.WriteHeader(http.StatusNoContent)
					return
				}
			}
			// else: origin not in whitelist, no CORS headers, browser will deny
		}
		next.ServeHTTP(w, r)
	})
}

// isValidOrigin checks if the origin string is a well-formed URL with scheme
// and host (W16 F-002: reject malformed origins before reflecting them).
func isValidOrigin(origin string) bool {
	// Empty already handled by caller
	if origin == "" {
		return false
	}
	// Parse as URL
	u, err := url.Parse(origin)
	if err != nil {
		return false
	}
	// Must have scheme (http/https) and host
	if u.Scheme == "" || u.Host == "" {
		return false
	}
	// Scheme must be http or https
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	// Origin must not have path, query, or fragment (per CORS spec, origin is scheme + host + port)
	if u.Path != "" && u.Path != "/" {
		return false
	}
	if u.RawQuery != "" || u.Fragment != "" {
		return false
	}
	return true
}