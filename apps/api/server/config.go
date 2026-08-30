// Package server 提供下游 serve 面（assembly 的服务器面扩展 · VP-024 R1）。
//
// 背景：主仓服务面（internal/config · internal/composition · internal/server）
// 对下游模块不可见；生成骨架此前只能「装配冒烟」。本包把 config 装载 +
// 标准下游组合装配 + 中央面接线 + RT-D02 优雅停机封装为公开面：
// `schema-ui serve` 与生成骨架 cmd/server 共用同一实现。
//
// 下游基线（有界口径，D-001）：中央面 = healthz/readyz、登录、
// schema 文档、manifest/bootstrap、模块贡献路由；不装配
// jobs / metrics / tracing / objects / mail-admin / settings 等面
// （主仓形态保留）。RT-D02 合同（workspace-021 D-002 v0.1.1）§1 全序停机
// 按 §6 `http.shutdown_timeout` 预算执行，退出码 0/1 语义由调用方（CLI /
// 模板 main）落实。
package server

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"net/netip"
	"os"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/internal/config"
	"gopkg.in/yaml.v3"
)

//go:embed config.default.yaml
var defaultConfigYAML []byte

// Config 是下游 serve 面的公开配置（RT-K01 YAML+env 语义子集；RT-D02 §6 键）。
//
// 装载链：代码默认 →（显式 path 或内嵌默认）YAML（${VAR}/${VAR:-default} 插值，
// 裸 ${VAR} 未设置 fail-closed，同 internal/config 语义）→ 进程 env 定向覆盖。
type Config struct {
	AppName string
	AppEnv  string // "" = 未声明（validate 拒绝，refusing to guess；W15 F-001）；非 development 时密钥/种子类键 fail-closed
	LogLevel string

	HTTPAddr        string
	ReadTimeout     time.Duration
	WriteTimeout    time.Duration
	IdleTimeout     time.Duration
	ShutdownTimeout time.Duration // RT-D02 §6：http.shutdown_timeout（默认 10s；≤0 或解析失败 fail-closed）
	TrustedProxies  []string
	CORSOrigins     []string // 可选 CORS 白名单；空 = 无跨域头（同源反代默认）

	DBDialect string // "" | "sqlite" | "postgres"（缺省 sqlite）
	DBPath    string // sqlite 文件路径（缺省 ./data/schema-ui.db）
	DBDSN     string // postgres 连接串（dialect=postgres 必填；sqlite 必须为空）

	// ProfileName 缺省 "admin"；下游组合 = 本包固定装配的标准集
	// （server-registration / auth-session / manifest-route / navigation-capability
	//  / schema-render / operationlog / users），配置仅可收窄。
	ProfileName string

	AuthJWTSecret        string
	AuthAccessTTL        time.Duration
	AuthRefreshTTL       time.Duration
	AdminInitialPassword string
	PublicBaseURL        string
}

// yamlFile 是公开配置的 YAML 面（KnownFields 严格校验，未知键拒绝）。
type yamlFile struct {
	App struct {
		Name *string `yaml:"name"`
		Env  *string `yaml:"env"`
	} `yaml:"app"`
	HTTP struct {
		Addr            *string   `yaml:"addr"`
		ReadTimeout     *string   `yaml:"read_timeout"`
		WriteTimeout    *string   `yaml:"write_timeout"`
		IdleTimeout     *string   `yaml:"idle_timeout"`
		ShutdownTimeout *string   `yaml:"shutdown_timeout"`
		TrustedProxies  []string  `yaml:"trusted_proxies"`
		CORSOrigins     []string  `yaml:"cors_origins"`
	} `yaml:"http"`
	DB struct {
		Dialect *string `yaml:"dialect"`
		Path    *string `yaml:"path"`
		DSN     *string `yaml:"dsn"`
	} `yaml:"db"`
	Profile *string `yaml:"profile"`
	Auth    struct {
		JWTSecret     *string `yaml:"jwt_secret"`
		AccessTTL     *string `yaml:"access_ttl"`
		RefreshTTL    *string `yaml:"refresh_ttl"`
		PublicBaseURL *string `yaml:"public_base_url"`
	} `yaml:"auth"`
	Admin struct {
		InitialPassword *string `yaml:"initial_password"`
	} `yaml:"admin"`
	Log struct {
		Level *string `yaml:"level"`
	} `yaml:"log"`
}

// LoadConfig 装载 serve 配置。
//
//	path 为空：内嵌默认（可被 env 覆盖）；显式 path 必须存在（fail-closed）。
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{
		AppName:           "schema-ui-app",
		AppEnv:            "",
		LogLevel:          "info",
		HTTPAddr:          "127.0.0.1:25080", // W15 F-001: loopback default; LAN exposure requires explicit config
		ReadTimeout:       5 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
		ShutdownTimeout:   10 * time.Second, // RT-D02 §6 默认
		DBDialect:         "sqlite",
		DBPath:            "./data/schema-ui.db",
		ProfileName:       "admin",
		AuthAccessTTL:     15 * time.Minute,
		AuthRefreshTTL:    30 * 24 * time.Hour,
		AdminInitialPassword: "",
	}

	var yamlBytes []byte
	if strings.TrimSpace(path) != "" {
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, fmt.Errorf("server: read config %q: %w", path, err)
		}
		yamlBytes = b
	} else {
		yamlBytes = defaultConfigYAML
	}

	interpolated, err := interpolateAll(string(yamlBytes))
	if err != nil {
		return nil, err
	}
	var yf yamlFile
	dec := yaml.NewDecoder(strings.NewReader(interpolated))
	dec.KnownFields(true)
	if err := dec.Decode(&yf); err != nil {
		if errors.Is(err, io.EOF) {
			// 空（或仅注释）文件 = 全默认
			yf = yamlFile{}
		} else {
			return nil, fmt.Errorf("server: parse config YAML: %w", err)
		}
	}
	var extra yamlFile
	if err := dec.Decode(&extra); err != io.EOF {
		return nil, fmt.Errorf("server: parse config YAML: multiple YAML documents are not supported (%v)", err)
	}

	if yf.App.Name != nil {
		cfg.AppName = *yf.App.Name
	}
	if yf.App.Env != nil {
		cfg.AppEnv = *yf.App.Env
	}
	if yf.Log.Level != nil {
		cfg.LogLevel = *yf.Log.Level
	}
	if yf.HTTP.Addr != nil {
		cfg.HTTPAddr = *yf.HTTP.Addr
	}
	if yf.HTTP.ReadTimeout != nil {
		if cfg.ReadTimeout, err = time.ParseDuration(*yf.HTTP.ReadTimeout); err != nil {
			return nil, fmt.Errorf("server: http.read_timeout %q: %w", *yf.HTTP.ReadTimeout, err)
		}
	}
	if yf.HTTP.WriteTimeout != nil {
		if cfg.WriteTimeout, err = time.ParseDuration(*yf.HTTP.WriteTimeout); err != nil {
			return nil, fmt.Errorf("server: http.write_timeout %q: %w", *yf.HTTP.WriteTimeout, err)
		}
	}
	if yf.HTTP.IdleTimeout != nil {
		if cfg.IdleTimeout, err = time.ParseDuration(*yf.HTTP.IdleTimeout); err != nil {
			return nil, fmt.Errorf("server: http.idle_timeout %q: %w", *yf.HTTP.IdleTimeout, err)
		}
	}
	if yf.HTTP.ShutdownTimeout != nil {
		if cfg.ShutdownTimeout, err = time.ParseDuration(*yf.HTTP.ShutdownTimeout); err != nil {
			return nil, fmt.Errorf("server: http.shutdown_timeout %q: %w", *yf.HTTP.ShutdownTimeout, err)
		}
	}
	if len(yf.HTTP.TrustedProxies) > 0 {
		cfg.TrustedProxies = append([]string(nil), yf.HTTP.TrustedProxies...)
	}
	if len(yf.HTTP.CORSOrigins) > 0 {
		cfg.CORSOrigins = append([]string(nil), yf.HTTP.CORSOrigins...)
	}
	if yf.DB.Dialect != nil {
		cfg.DBDialect = *yf.DB.Dialect
	}
	if yf.DB.Path != nil {
		cfg.DBPath = *yf.DB.Path
	}
	if yf.DB.DSN != nil {
		cfg.DBDSN = *yf.DB.DSN
	}
	if yf.Profile != nil {
		cfg.ProfileName = *yf.Profile
	}
	if yf.Auth.JWTSecret != nil {
		cfg.AuthJWTSecret = *yf.Auth.JWTSecret
	}
	if yf.Auth.AccessTTL != nil {
		if cfg.AuthAccessTTL, err = time.ParseDuration(*yf.Auth.AccessTTL); err != nil {
			return nil, fmt.Errorf("server: auth.access_ttl %q: %w", *yf.Auth.AccessTTL, err)
		}
	}
	if yf.Auth.RefreshTTL != nil {
		if cfg.AuthRefreshTTL, err = time.ParseDuration(*yf.Auth.RefreshTTL); err != nil {
			return nil, fmt.Errorf("server: auth.refresh_ttl %q: %w", *yf.Auth.RefreshTTL, err)
		}
	}
	if yf.Auth.PublicBaseURL != nil {
		cfg.PublicBaseURL = *yf.Auth.PublicBaseURL
	}
	if yf.Admin.InitialPassword != nil {
		cfg.AdminInitialPassword = *yf.Admin.InitialPassword
	}

	// 进程 env 定向覆盖（秘密注入路径，与主仓 env 键对齐）。
	applyEnv := func(key string, apply func(string)) {
		if v, ok := os.LookupEnv(key); ok && strings.TrimSpace(v) != "" {
			apply(strings.TrimSpace(v))
		}
	}
	applyEnv("APP_ENV", func(v string) { cfg.AppEnv = v })
	applyEnv("LOG_LEVEL", func(v string) { cfg.LogLevel = v })
	applyEnv("HTTP_ADDR", func(v string) { cfg.HTTPAddr = v })
	applyEnv("HTTP_SHUTDOWN_TIMEOUT", func(v string) {
		cfg.ShutdownTimeout, _ = parseDurationStrict(v)
	})
	applyEnv("DB_DIALECT", func(v string) { cfg.DBDialect = v })
	applyEnv("DB_PATH", func(v string) { cfg.DBPath = v })
	applyEnv("DB_DSN", func(v string) { cfg.DBDSN = v })
	applyEnv("PROFILE_NAME", func(v string) { cfg.ProfileName = v })
	applyEnv("AUTH_JWT_SECRET", func(v string) { cfg.AuthJWTSecret = v })
	applyEnv("ADMIN_INITIAL_PASSWORD", func(v string) { cfg.AdminInitialPassword = v })
	applyEnv("AUTH_PUBLIC_BASE_URL", func(v string) { cfg.PublicBaseURL = v })

	if err := cfg.validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// parseDurationStrict 解析 env 时长；非法值返回 0（validate 层 fail-closed）。
func parseDurationStrict(v string) (time.Duration, error) {
	d, err := time.ParseDuration(v)
	if err != nil || d <= 0 {
		return 0, fmt.Errorf("invalid duration %q", v)
	}
	return d, nil
}

// validate 实施 fail-closed 规则（镜像主仓 ValidateProd 的必要子集）。
func (c *Config) validate() error {
	// W15 F-001: an explicitly declared APP_ENV is mandatory; the embedded
	// default pins "development", so only a custom YAML that omits app.env
	// reaches this gate. Never guess an environment from silence.
	if c.AppEnv == "" {
		return errors.New("server: APP_ENV must be set explicitly (development for local runs, production for deployments); refusing to guess")
	}
	if c.HTTPAddr == "" {
		return errors.New("server: http.addr must not be empty")
	}
	if c.ShutdownTimeout <= 0 {
		return fmt.Errorf("server: http.shutdown_timeout must be > 0 (got %s)", c.ShutdownTimeout)
	}
	if c.ReadTimeout < 0 || c.WriteTimeout < 0 || c.IdleTimeout < 0 {
		return errors.New("server: http timeouts must not be negative")
	}
	switch c.DBDialect {
	case "", "sqlite":
		c.DBDialect = "sqlite"
		if c.DBDSN != "" {
			return errors.New("server: db.dsn must be empty when dialect is sqlite")
		}
		if c.DBPath == "" {
			return errors.New("server: db.path must not be empty when dialect is sqlite")
		}
	case "postgres":
		if c.DBDSN == "" {
			return errors.New("server: db.dsn is required when dialect is postgres")
		}
	default:
		return fmt.Errorf("server: unknown db.dialect %q (sqlite | postgres)", c.DBDialect)
	}
	for _, cidr := range c.TrustedProxies {
		if _, err := netip.ParsePrefix(cidr); err != nil {
			return fmt.Errorf("server: invalid http.trusted_proxies entry %q", cidr)
		}
	}
	dev := c.AppEnv == "development"
	if !dev && c.AuthJWTSecret == "" {
		return errors.New("server: AUTH_JWT_SECRET must be set in non-development environment")
	}
	// W15 F-002: reuse the main production HS256 bar (≥32 chars + letters and
	// digits) as a single source — a short or guessable key must not silently
	// start the public serve surface.
	if !dev {
		if err := config.ValidateJWTSecretStrength("AUTH_JWT_SECRET", c.AuthJWTSecret); err != nil {
			return fmt.Errorf("server: %w", err)
		}
	}
	if !dev && c.AdminInitialPassword == "" {
		return errors.New("server: ADMIN_INITIAL_PASSWORD must be set to seed the initial admin user in non-development environment")
	}
	return nil
}

// interpolateAll 展开 YAML 文本中的 ${VAR} / ${VAR:-default}。
// 裸 ${VAR} 未设置（无默认）→ fail-closed（同 internal/config 语义）。
func interpolateAll(text string) (string, error) {
	var sb strings.Builder
	rest := text
	for {
		start := strings.Index(rest, "${")
		if start < 0 {
			sb.WriteString(rest)
			return sb.String(), nil
		}
		sb.WriteString(rest[:start])
		end := strings.Index(rest[start:], "}")
		if end < 0 {
			return "", fmt.Errorf("server: config interpolation: unterminated ${ in %q", rest[start:start+24])
		}
		expr := rest[start+2 : start+end]
		rest = rest[start+end+1:]
		name, def, hasDef := strings.Cut(expr, ":-")
		name = strings.TrimSpace(name)
		if v, ok := os.LookupEnv(name); ok {
			sb.WriteString(v)
			continue
		}
		if hasDef {
			sb.WriteString(def)
			continue
		}
		return "", fmt.Errorf("server: config interpolation: ${%s} is not set and has no default (fail-closed)", name)
	}
}