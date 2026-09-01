package config

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/mail"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/magicvr/schema-ui-core/apps/api/kernel"
)

//go:embed config.default.yaml
var defaultConfigYAML []byte

// DefaultCacheMaxEntries is the fallback bounded-entry budget for the
// in-memory cache provider (VP-026 / workspace-026 GOAL-003 D-001): applied
// by Load when neither YAML nor env configures cache.max_entries, and by the
// composition root for zero-value (loader-bypassed) Configs.
const DefaultCacheMaxEntries = 10000

// DefaultEventBusBuffer is the fallback per-subscription buffer size for the
// in-memory event-bus provider (VP-028 / workspace-028 GOAL-003 D-001): applied
// by Load when neither YAML nor env configures eventbus.buffer_size, and by the
// composition root for zero-value (loader-bypassed) Configs. Must match
// kernel.DefaultEventBusBuffer.
const DefaultEventBusBuffer = 64

// MaxEventBusBuffer is the upper bound on eventbus.buffer_size (D-001): buffers
// larger than this are rejected fail-closed at config load to prevent unbounded
// memory growth or configuration typos from causing resource exhaustion.
const MaxEventBusBuffer = 4096

// Config is the R2 runtime configuration: HTTP + logging, plus the auth
// (JWT / refresh / SQLite) and dev-session surface defined by GOAL-005 D-004,
// and the upload surface (W7: UPLOAD_* moved from handler package vars into
// Config so YAML is the single configuration authority).
type Config struct {
	AppName      string
	AppEnv       string
	HTTPAddr     string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
	// HTTPShutdownTimeout is the graceful-shutdown drain budget (VP-021
	// contract §6): SIGINT/SIGTERM -> http.Server.Shutdown grace -> forced
	// exit. Default 10s; invalid (unparsable or <=0) values fail closed at
	// startup. Env override: HTTP_SHUTDOWN_TIMEOUT.
	HTTPShutdownTimeout time.Duration
	// HTTPCORSOrigins is the optional CORS allow-list (W15-F05). Empty means
	// no Access-Control headers (same-origin Nginx remains the default).
	HTTPCORSOrigins []string
	// HTTPTrustedProxies is the explicit reverse-proxy CIDR allow-list (W7
	// F-008). Only a direct peer within one of these CIDRs may supply the
	// X-Real-IP header used for login/captcha rate limiting. Empty means only
	// loopback is trusted (fail-safe); operators must add their proxy networks
	// explicitly rather than trusting all RFC1918 addresses.
	HTTPTrustedProxies []string
	LogLevelName       string

	AuthJWTSecret string
	// AuthJWTSecretPrevious is the optional previous (rotated-out) signing key
	// (VP-016 R1 / workspace-016 Root D-002): auth.jwt_secret_previous /
	// AUTH_JWT_SECRET_PREVIOUS. Empty (the default) keeps single-key behavior
	// identical in every environment. When set outside development it must
	// satisfy the same strength rule as the current key and must differ from
	// it; the overlap window lasts as long as the key stays configured and is
	// retired by removing it and restarting.
	AuthJWTSecretPrevious string
	AuthAccessTTL         time.Duration
	AuthRefreshTTL        time.Duration
	// AuthPublicBaseURL is the optional canonical external origin
	// (auth.public_base_url / AUTH_PUBLIC_BASE_URL; W13 F-006 · GOAL-013
	// A-001). When set (e.g. "https://ops.example.com") emailed invitation
	// links are always built from it — never from the request's Host /
	// X-Forwarded-Proto headers, which a client can influence. Empty keeps
	// the request-derived fallback (single-host dev / direct deployments).
	AuthPublicBaseURL string
	DBPath            string
	// DBDialect is the store dialect (VP-013 / R1 v1.4 §5): "" or "sqlite" or
	// "postgres". Load normalizes empty to "sqlite"; ValidateProd rejects
	// unknown values and enforces DSN/path pairing rules.
	DBDialect string
	// DBDSN is the postgres SQL connection string (DB_DSN; no default). It must
	// be empty when DBDialect is sqlite and non-empty for postgres. When set it
	// overrides the exploded db.host/port/name/user/password params below.
	DBDSN string
	// Exploded postgres connection params (db.host/port/name/user/password/
	// sslmode). host/port/sslmode have defaults; name/user/password come from
	// YAML/configs/.env/process env. The PASSWORD is a secret and must be
	// supplied via DB_PASSWORD (env / configs/.env), never hardcoded. A DSN is
	// built from these when db.dsn is empty.
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	// DBConnPool carries connection-pool bounds (0 = driver default for
	// postgres / the sqlite file-store default of 4). Wired through to
	// store.OpenOptions; sqlite uses PoolMaxOpenConns only.
	DBPoolMaxOpen  int
	DBPoolMaxIdle  int
	DBConnLifetime time.Duration

	AdminInitialPassword  string
	AuthDevSessionEnabled bool

	UploadAllowedTypes    string
	UploadMaxFilesPerUser int
	UploadMaxBytesPerUser int

	// W9 (GOAL-010): brand asset upload processing policy (config.yaml
	// branding section; env-overridable). Out-of-range values fall back to
	// the store defaults at construction.
	BrandingLogoMaxDimension int
	BrandingFaviconDimension int
	BrandingJPEGQuality      int
	BrandingMaxBytes         int

	// Object-storage surface (VP-014 / workspace-014 GOAL-002 D-001, R1).
	// ObjectsDriver selects the kernel.ObjectStore adapter: "" / "local" (the
	// embedded default; root derived from filepath.Dir(DBPath) at composition)
	// or "s3" (explicit production configuration; readyz extends only then).
	ObjectsDriver string
	// ObjectsLocalRoot optionally overrides the local disk root. Empty keeps
	// the historical derivation (the directory containing the sqlite file).
	ObjectsLocalRoot string
	// S3-compatible backend settings. Credentials are secrets: they must come
	// from env interpolation (configs/.env or process env), never literals.
	// driver=s3 requires endpoint/bucket/access_key_id/secret_access_key
	// (fail-closed at load); any non-empty s3.* key with driver=local fails.
	ObjectsS3Endpoint        string
	ObjectsS3Region          string
	ObjectsS3Bucket          string
	ObjectsS3AccessKeyID     string
	ObjectsS3SecretAccessKey string
	// ObjectsS3UsePathStyle defaults to true (MinIO/R2 need path-style);
	// virtual-host style can be enabled for AWS-compatible endpoints.
	ObjectsS3UsePathStyle bool

	// Cache surface (VP-026 / workspace-026 GOAL-003 D-001, R2).
	// CacheMaxEntries is the in-memory cache provider's bounded-entry budget:
	// after every Set the TOTAL stored entry count (including not-yet-swept
	// expired entries) is <= this value; the oldest entry is FIFO-evicted.
	// Default 10000; non-positive values fail closed at load (a typo must
	// never silently degrade to the default).
	CacheMaxEntries int

	// EventBus surface (VP-028 / workspace-028 GOAL-003 D-001, R2).
	// EventBusBufferSize is the in-memory event-bus provider's per-subscription
	// buffered-channel capacity. When a subscriber's buffer is full, Publish
	// blocks until space is available, the context is cancelled, or Stop drains.
	// Default 64 (kernel.DefaultEventBusBuffer); <= 0 falls back to that default;
	// > MaxEventBusBuffer (4096) fails closed at load.
	EventBusBufferSize int

	// Observability metrics surface (VP-015 / workspace-015 GOAL-002 D-001).
	// MetricsEnabled selects the DEDICATED Prometheus exposition listener —
	// never a route on the main mux. It is off by default so mvp/dev/Compose
	// keep the documented no-collector default (no extra port, no behavior
	// change). MetricsAddr defaults to loopback-only; any non-loopback bind
	// requires MetricsAuthToken (fail-closed at load/validate). The token is
	// a secret: supply it via OBSERVABILITY_METRICS_AUTH_TOKEN env /
	// configs/.env, never a YAML literal.
	MetricsEnabled   bool
	MetricsAddr      string
	MetricsAuthToken string

	// Observability traces surface (VP-015 / workspace-015 GOAL-004 D-001).
	// TracesEnabled selects OTLP/HTTP export; disabled (the default) keeps a
	// pure no-op tracer path — no provider, no spans, zero behavior change
	// (mvp/dev never require a collector). The endpoint is required when
	// enabled and must be an absolute http(s) URL (the SDK appends
	// /v1/traces). SampleRatio applies as ParentBased(TraceIDRatioBased)
	// with default 1.0. Export failures are logged, never fatal.
	TracesEnabled     bool
	TracesEndpoint    string
	TracesSampleRatio float64

	// Outbound-mail surface (VP-017 / workspace-017 GOAL-003 D-001, R2).
	// All keys empty = SMTP unconfigured: the process keeps the embedded
	// capture/log default and mvp/dev startup is unaffected. Any single key
	// makes the block explicit: host/username/password/from become REQUIRED
	// (fail-closed at ValidateProd in every environment; port defaults to
	// 465). The password is a secret: it must arrive via env interpolation
	// (MAIL_SMTP_PASSWORD via process env / configs/.env), never literals.
	MailSMTPHost     string
	MailSMTPPort     int // 0 = unset -> adapter applies the frozen 465
	MailSMTPUsername string
	MailSMTPPassword string
	MailSMTPFrom     string

	// Outbound-mail channel surface (VP-017 R6 / workspace-017 GOAL-007,
	// contract frozen by workspace-017 GOAL-006 D-002 §2/§4). MailChannel is
	// the explicit selector: "mock" | "resend" | "smtp" | "" (empty derives:
	// exactly one fully configured production block wins; both -> fail-closed
	// ambiguity; none -> mock, preserving the pre-R6 behavior for existing
	// mail.smtp deployments). The resend block mirrors the SMTP pairing
	// contract: touching ANY mail.resend.* key requires api-key + from
	// (fail-closed in every environment); api-key is a SECRET — env
	// interpolation only (MAIL_RESEND_API_KEY via process env /
	// configs/.env), never a YAML literal.
	MailChannel      string
	MailResendAPIKey string
	MailResendFrom   string

	// MailConfigMasterKey is the optional passphrase for the local master key
	// that encrypts admin-entered channel secrets at rest (VP-017 R7 / Root
	// D-007). ENV ONLY (MAIL_CONFIG_MASTER_KEY) — never a YAML literal. Empty
	// keeps the auto-generated key file under the data directory.
	MailConfigMasterKey string
	// MailMasterKeyPath optionally relocates the auto-generated master key
	// FILE (mail.master_key_path / MAIL_MASTER_KEY_PATH; W13 F-017 · GOAL-013
	// A-001). The default keeps the key beside the sqlite data directory,
	// which means one backup/ snapshot leaks both the ciphertext AND its key;
	// an operator who backs up that directory can point this at a separate
	// volume instead.
	MailMasterKeyPath string

	// NavigationOrder is the optional full navigation ordering (GOAL-013 D-002
	// §4): YAML navigation.order or NAVIGATION_ORDER env (comma-separated
	// NodeIDs). Empty means the built-in kernel default applies.
	NavigationOrder []string

	ProfileName       string
	ModulesEnabled    []string
	ProfileSource     string
	ProfilePrecedence []string
	// RuntimeMode is the startup-selected operational state consumed by the
	// bootstrap/status projections and the HTTP write gate (GOAL-006 R5).
	RuntimeMode  RuntimeMode
	ProfileError error
	// LoadError carries a fatal configuration-load failure (missing CONFIG_FILE,
	// unset ${VAR} without default, invalid YAML). ValidateProd surfaces it so
	// startup fails closed instead of silently running on fallback values.
	LoadError error
}

// RuntimeMode is the bounded operational state for cross-module availability.
type RuntimeMode string

const (
	RuntimeModeNormal      RuntimeMode = "normal"
	RuntimeModeMaintenance RuntimeMode = "maintenance"
	RuntimeModeDegraded    RuntimeMode = "degraded"
	RuntimeModeReadOnly    RuntimeMode = "read-only"
)

// ValidRuntimeMode reports whether a value belongs to the frozen R5 contract.
func ValidRuntimeMode(mode RuntimeMode) bool {
	switch mode {
	case RuntimeModeNormal, RuntimeModeMaintenance, RuntimeModeDegraded, RuntimeModeReadOnly:
		return true
	default:
		return false
	}
}

// yamlFile mirrors the on-disk config schema (see configs/config.yaml and the
// embedded config.default.yaml). Only recognized keys are read; unknown keys
// are rejected by yaml.v3 KnownFields so typos fail loudly.
type yamlFile struct {
	App struct {
		Name    *string `yaml:"name"`
		Env     *string `yaml:"env"`
		Profile *string `yaml:"profile"`
		// T-06 (GOAL-013 D-007): app.modules — preset (builtin name or preset
		// file path) OR inline list; mutually exclusive. This is the ONLY
		// authority for the enabled-module set (APP_PROFILE /
		// APP_MODULES_ENABLED and the legacy app.modules_enabled comma string
		// are no longer read).
		Modules yaml.Node `yaml:"modules"`
	} `yaml:"app"`
	HTTP struct {
		Addr            *string `yaml:"addr"`
		ReadTimeout     *string `yaml:"read_timeout"`
		WriteTimeout    *string `yaml:"write_timeout"`
		IdleTimeout     *string `yaml:"idle_timeout"`
		ShutdownTimeout *string `yaml:"shutdown_timeout"`
		CORSOrigins     *string `yaml:"cors_origins"`
		TrustedProxies  *string `yaml:"trusted_proxies"`
	} `yaml:"http"`
	Log struct {
		Level *string `yaml:"level"`
	} `yaml:"log"`
	Auth struct {
		JWTSecret         *string `yaml:"jwt_secret"`
		JWTSecretPrevious *string `yaml:"jwt_secret_previous"`
		AccessTTL         *string `yaml:"access_ttl"`
		RefreshTTL        *string `yaml:"refresh_ttl"`
		DevSessionEnabled *bool   `yaml:"dev_session_enabled"`
		PublicBaseURL     *string `yaml:"public_base_url"`
	} `yaml:"auth"`
	DB struct {
		Path         *string `yaml:"path"`
		Dialect      *string `yaml:"dialect"`
		DSN          *string `yaml:"dsn"`
		Host         *string `yaml:"host"`
		Port         *string `yaml:"port"`
		Name         *string `yaml:"name"`
		User         *string `yaml:"user"`
		Password     *string `yaml:"password"`
		SSLMode      *string `yaml:"sslmode"`
		PoolMaxOpen  *int    `yaml:"pool_max_open"`
		PoolMaxIdle  *int    `yaml:"pool_max_idle"`
		ConnLifetime *string `yaml:"conn_max_lifetime"`
	} `yaml:"db"`
	Admin struct {
		InitialPassword *string `yaml:"initial_password"`
	} `yaml:"admin"`
	Upload struct {
		AllowedTypes    *string `yaml:"allowed_types"`
		MaxFilesPerUser int     `yaml:"max_files_per_user"`
		MaxBytesPerUser int     `yaml:"max_bytes_per_user"`
	} `yaml:"upload"`
	Branding struct {
		LogoMaxDimension int `yaml:"logo_max_dimension"`
		FaviconDimension int `yaml:"favicon_dimension"`
		JPEGQuality      int `yaml:"jpeg_quality"`
		MaxBytes         int `yaml:"max_bytes"`
	} `yaml:"branding"`
	Storage struct {
		Objects struct {
			Driver *string `yaml:"driver"`
			Local  struct {
				Root *string `yaml:"root"`
			} `yaml:"local"`
			S3 struct {
				Endpoint        *string `yaml:"endpoint"`
				Region          *string `yaml:"region"`
				Bucket          *string `yaml:"bucket"`
				AccessKeyID     *string `yaml:"access_key_id"`
				SecretAccessKey *string `yaml:"secret_access_key"`
				UsePathStyle    *bool   `yaml:"use_path_style"`
			} `yaml:"s3"`
		} `yaml:"objects"`
	} `yaml:"storage"`
	Observability struct {
		Metrics struct {
			Enabled   *bool   `yaml:"enabled"`
			Addr      *string `yaml:"addr"`
			AuthToken *string `yaml:"auth_token"`
		} `yaml:"metrics"`
		Traces struct {
			Enabled     *bool    `yaml:"enabled"`
			Endpoint    *string  `yaml:"endpoint"`
			SampleRatio *float64 `yaml:"sample_ratio"`
		} `yaml:"traces"`
	} `yaml:"observability"`
	Mail struct {
		Channel       *string `yaml:"channel"`
		MasterKeyPath *string `yaml:"master_key_path"`
		SMTP          struct {
			Host     *string `yaml:"host"`
			Port     *int    `yaml:"port"`
			Username *string `yaml:"username"`
			Password *string `yaml:"password"`
			From     *string `yaml:"from"`
		} `yaml:"smtp"`
		Resend struct {
			APIKey *string `yaml:"api-key"`
			From   *string `yaml:"from"`
		} `yaml:"resend"`
	} `yaml:"mail"`
	Cache struct {
		// VP-026 / workspace-026 GOAL-003 D-001: bounded-entry budget of the
		// in-memory cache provider (respects only positive values; <= 0 fails
		// closed at load). Env: CACHE_MAX_ENTRIES.
		MaxEntries *int `yaml:"max_entries"`
	} `yaml:"cache"`
	EventBus struct {
		// VP-028 / workspace-028 GOAL-003 D-001: per-subscription buffer size
		// of the in-memory event-bus provider. <= 0 falls back to
		// DefaultEventBusBuffer (64); > MaxEventBusBuffer (4096) fails closed.
		// Env: EVENTBUS_BUFFER_SIZE.
		BufferSize *int `yaml:"buffer_size"`
	} `yaml:"eventbus"`
	Navigation struct {
		Order yaml.Node `yaml:"order"`
	} `yaml:"navigation"`
	Runtime struct {
		Mode *string `yaml:"mode"`
	} `yaml:"runtime"`
}

// defaultYAMLPath is the operator-editable config file loaded when CONFIG_FILE
// is not set. Missing is fine (embedded default applies); CONFIG_FILE set but
// missing is a startup error (fail-closed).
const defaultYAMLPath = "configs/config.yaml"

// Load reads configuration with safe local defaults. Priority (highest first):
//  1. process environment variables (already set -> override YAML)
//  2. CONFIG_FILE YAML (default configs/config.yaml)
//  3. embedded config.default.yaml (go:embed)
//
// YAML values support ${VAR} (fail-closed when unset) and ${VAR:-default}.
// CONFIG_ENV_FILE (default configs/.env) may supply secret values for the
// interpolation; copy configs/.env.example to create it. It never overrides
// an already-set process env. Load never returns an error: fatal load
// failures land in LoadError (and therefore ValidateProd) so existing call
// sites keep working.
func Load() *Config {
	cfg := &Config{
		AppName:      "schema-ui-core-api",
		AppEnv:       "",
		HTTPAddr:     ":25080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		// VP-021 contract §6: default drain budget = 10s (mirrors the legacy
		// hard-coded shutdown context in cmd/server/main.go).
		HTTPShutdownTimeout: 10 * time.Second,
		LogLevelName:        "info",

		AuthJWTSecret:            "",
		AuthAccessTTL:            15 * time.Minute,
		AuthRefreshTTL:           30 * 24 * time.Hour,
		DBPath:                   "./data/schema-ui.db",
		DBDialect:                "sqlite",
		DBDSN:                    "",
		DBHost:                   "127.0.0.1",
		DBPort:                   "5432",
		DBSSLMode:                "disable",
		DBPoolMaxOpen:            0,
		DBPoolMaxIdle:            0,
		DBConnLifetime:           0,
		AdminInitialPassword:     "",
		AuthDevSessionEnabled:    false,
		UploadMaxFilesPerUser:    1000,
		UploadMaxBytesPerUser:    256 << 20,
		BrandingLogoMaxDimension: 512,
		BrandingFaviconDimension: 64,
		BrandingJPEGQuality:      82,
		BrandingMaxBytes:         4 << 20,
		ObjectsDriver:            "local",
		ObjectsS3UsePathStyle:    true,
		CacheMaxEntries:          DefaultCacheMaxEntries,
		EventBusBufferSize:       DefaultEventBusBuffer,
		MetricsEnabled:           false,
		MetricsAddr:              "127.0.0.1:25081",
		MetricsAuthToken:         "",
		TracesEnabled:            false,
		TracesEndpoint:           "",
		TracesSampleRatio:        1.0,
		RuntimeMode:              RuntimeModeNormal,
	}

	// Optional env-file layer for secret values only (dev convenience).
	if err := loadEnvFile(cfg); err != nil {
		cfg.LoadError = err
		return cfg
	}

	// Pick the YAML source: CONFIG_FILE (explicit -> must exist), default
	// configs/config.yaml, else the embedded default.
	var yamlBytes []byte
	if p := strings.TrimSpace(os.Getenv("CONFIG_FILE")); p != "" {
		b, err := os.ReadFile(p)
		if err != nil {
			cfg.LoadError = fmt.Errorf("CONFIG_FILE %q: %w", p, err)
			return cfg
		}
		yamlBytes = b
	} else if b, err := os.ReadFile(defaultYAMLPath); err == nil {
		yamlBytes = b
	} else if !os.IsNotExist(err) {
		cfg.LoadError = fmt.Errorf("read %s: %w", defaultYAMLPath, err)
		return cfg
	} else {
		yamlBytes = defaultConfigYAML
	}

	// Interpolate ${VAR} / ${VAR:-default} (fail-closed on bare ${VAR}).
	interpolated, err := interpolateAll(string(yamlBytes))
	if err != nil {
		cfg.LoadError = err
		return cfg
	}

	var yf yamlFile
	dec := yaml.NewDecoder(strings.NewReader(interpolated))
	dec.KnownFields(true)
	if err := dec.Decode(&yf); err != nil {
		if errors.Is(err, io.EOF) {
			// F-005: an empty (or comment-only) file means "all defaults" —
			// no keys were supplied, so the zero yamlFile keeps every code
			// default below.
			yf = yamlFile{}
		} else {
			cfg.LoadError = fmt.Errorf("parse config YAML: %w", err)
			return cfg
		}
	}
	// F-005: reject multi-document YAML — a second document would silently
	// escape KnownFields validation (or introduce a different schema).
	var extra yamlFile
	if err := dec.Decode(&extra); err != io.EOF {
		cfg.LoadError = fmt.Errorf("parse config YAML: multiple YAML documents are not supported (%v)", err)
		return cfg
	}

	// YAML values (already interpolated) as the base layer. Pointer fields
	// keep the code default when the key is omitted (F-002); explicit empty
	// strings stay empty.
	cfg.AppName = strPtrOr(yf.App.Name, cfg.AppName)
	cfg.AppEnv = strPtrOr(yf.App.Env, cfg.AppEnv)
	cfg.HTTPAddr = strPtrOr(yf.HTTP.Addr, cfg.HTTPAddr)
	cfg.ReadTimeout = orDurationPtr(yf.HTTP.ReadTimeout, cfg.ReadTimeout)
	cfg.WriteTimeout = orDurationPtr(yf.HTTP.WriteTimeout, cfg.WriteTimeout)
	cfg.IdleTimeout = orDurationPtr(yf.HTTP.IdleTimeout, cfg.IdleTimeout)
	if cfg.HTTPShutdownTimeout, cfg.LoadError = strictDurationPtr(yf.HTTP.ShutdownTimeout, "http.shutdown_timeout", cfg.HTTPShutdownTimeout); cfg.LoadError != nil {
		return cfg
	}
	if yf.HTTP.CORSOrigins != nil {
		cfg.HTTPCORSOrigins = splitCSV(*yf.HTTP.CORSOrigins)
	}
	if yf.HTTP.TrustedProxies != nil {
		cfg.HTTPTrustedProxies = splitCSV(*yf.HTTP.TrustedProxies)
	}
	cfg.LogLevelName = strPtrOr(yf.Log.Level, cfg.LogLevelName)
	cfg.AuthJWTSecret = strPtrOr(yf.Auth.JWTSecret, cfg.AuthJWTSecret)
	cfg.AuthJWTSecretPrevious = strPtrOr(yf.Auth.JWTSecretPrevious, cfg.AuthJWTSecretPrevious)
	cfg.AuthAccessTTL = orDurationPtr(yf.Auth.AccessTTL, cfg.AuthAccessTTL)
	cfg.AuthRefreshTTL = orDurationPtr(yf.Auth.RefreshTTL, cfg.AuthRefreshTTL)
	if yf.Auth.DevSessionEnabled != nil {
		cfg.AuthDevSessionEnabled = *yf.Auth.DevSessionEnabled
	}
	cfg.AuthPublicBaseURL = strPtrOr(yf.Auth.PublicBaseURL, cfg.AuthPublicBaseURL)
	cfg.DBPath = strPtrOr(yf.DB.Path, cfg.DBPath)
	cfg.DBDialect = strPtrOr(yf.DB.Dialect, cfg.DBDialect)
	cfg.DBDSN = strPtrOr(yf.DB.DSN, cfg.DBDSN)
	cfg.DBHost = strPtrOr(yf.DB.Host, cfg.DBHost)
	cfg.DBPort = strPtrOr(yf.DB.Port, cfg.DBPort)
	cfg.DBName = strPtrOr(yf.DB.Name, cfg.DBName)
	cfg.DBUser = strPtrOr(yf.DB.User, cfg.DBUser)
	cfg.DBPassword = strPtrOr(yf.DB.Password, cfg.DBPassword)
	cfg.DBSSLMode = strPtrOr(yf.DB.SSLMode, cfg.DBSSLMode)
	if yf.DB.PoolMaxOpen != nil {
		cfg.DBPoolMaxOpen = *yf.DB.PoolMaxOpen
	}
	if yf.DB.PoolMaxIdle != nil {
		cfg.DBPoolMaxIdle = *yf.DB.PoolMaxIdle
	}
	cfg.DBConnLifetime = orDurationPtr(yf.DB.ConnLifetime, cfg.DBConnLifetime)
	cfg.AdminInitialPassword = strPtrOr(yf.Admin.InitialPassword, cfg.AdminInitialPassword)
	cfg.UploadAllowedTypes = strings.TrimSpace(strPtrOr(yf.Upload.AllowedTypes, cfg.UploadAllowedTypes))
	if yf.Upload.MaxFilesPerUser > 0 {
		cfg.UploadMaxFilesPerUser = yf.Upload.MaxFilesPerUser
	}
	if yf.Upload.MaxBytesPerUser > 0 {
		cfg.UploadMaxBytesPerUser = yf.Upload.MaxBytesPerUser
	}
	if yf.Branding.LogoMaxDimension > 0 {
		cfg.BrandingLogoMaxDimension = yf.Branding.LogoMaxDimension
	}
	if yf.Branding.FaviconDimension > 0 {
		cfg.BrandingFaviconDimension = yf.Branding.FaviconDimension
	}
	if yf.Branding.JPEGQuality > 0 && yf.Branding.JPEGQuality <= 100 {
		cfg.BrandingJPEGQuality = yf.Branding.JPEGQuality
	}
	if yf.Branding.MaxBytes > 0 {
		cfg.BrandingMaxBytes = yf.Branding.MaxBytes
	}
	cfg.ObjectsDriver = strings.ToLower(strings.TrimSpace(strPtrOr(yf.Storage.Objects.Driver, cfg.ObjectsDriver)))
	cfg.ObjectsLocalRoot = strPtrOr(yf.Storage.Objects.Local.Root, cfg.ObjectsLocalRoot)
	cfg.ObjectsS3Endpoint = strPtrOr(yf.Storage.Objects.S3.Endpoint, cfg.ObjectsS3Endpoint)
	cfg.ObjectsS3Region = strPtrOr(yf.Storage.Objects.S3.Region, cfg.ObjectsS3Region)
	cfg.ObjectsS3Bucket = strPtrOr(yf.Storage.Objects.S3.Bucket, cfg.ObjectsS3Bucket)
	cfg.ObjectsS3AccessKeyID = strPtrOr(yf.Storage.Objects.S3.AccessKeyID, cfg.ObjectsS3AccessKeyID)
	cfg.ObjectsS3SecretAccessKey = strPtrOr(yf.Storage.Objects.S3.SecretAccessKey, cfg.ObjectsS3SecretAccessKey)
	if yf.Storage.Objects.S3.UsePathStyle != nil {
		cfg.ObjectsS3UsePathStyle = *yf.Storage.Objects.S3.UsePathStyle
	}
	cfg.MailSMTPHost = strPtrOr(yf.Mail.SMTP.Host, cfg.MailSMTPHost)
	if yf.Mail.SMTP.Port != nil {
		cfg.MailSMTPPort = *yf.Mail.SMTP.Port
	}
	cfg.MailSMTPUsername = strPtrOr(yf.Mail.SMTP.Username, cfg.MailSMTPUsername)
	cfg.MailSMTPPassword = strPtrOr(yf.Mail.SMTP.Password, cfg.MailSMTPPassword)
	cfg.MailSMTPFrom = strPtrOr(yf.Mail.SMTP.From, cfg.MailSMTPFrom)
	cfg.MailChannel = strings.ToLower(strings.TrimSpace(strPtrOr(yf.Mail.Channel, cfg.MailChannel)))
	cfg.MailResendAPIKey = strPtrOr(yf.Mail.Resend.APIKey, cfg.MailResendAPIKey)
	cfg.MailResendFrom = strPtrOr(yf.Mail.Resend.From, cfg.MailResendFrom)
	cfg.MailMasterKeyPath = strings.TrimSpace(strPtrOr(yf.Mail.MasterKeyPath, cfg.MailMasterKeyPath))
	if yf.Cache.MaxEntries != nil {
		cfg.CacheMaxEntries = *yf.Cache.MaxEntries
	}
	if yf.EventBus.BufferSize != nil {
		cfg.EventBusBufferSize = *yf.EventBus.BufferSize
	}
	if yf.Observability.Metrics.Enabled != nil {
		cfg.MetricsEnabled = *yf.Observability.Metrics.Enabled
	}
	cfg.MetricsAddr = strPtrOr(yf.Observability.Metrics.Addr, cfg.MetricsAddr)
	cfg.MetricsAuthToken = strPtrOr(yf.Observability.Metrics.AuthToken, cfg.MetricsAuthToken)
	if yf.Observability.Traces.Enabled != nil {
		cfg.TracesEnabled = *yf.Observability.Traces.Enabled
	}
	cfg.TracesEndpoint = strPtrOr(yf.Observability.Traces.Endpoint, cfg.TracesEndpoint)
	if yf.Observability.Traces.SampleRatio != nil {
		cfg.TracesSampleRatio = *yf.Observability.Traces.SampleRatio
	}
	profile := strPtrOr(yf.App.Profile, string(kernel.ProfileMVP))

	// navigation.order (GOAL-013 D-002 §4): sequence of NodeIDs. A malformed
	// value falls back to the kernel default with a warning (never fail-closed
	// — ordering is UI structure, not a security gate). The env override
	// NAVIGATION_ORDER (comma-separated) wins when set.
	cfg.NavigationOrder = parseNavigationOrder(yf.Navigation.Order)
	if raw := strings.TrimSpace(os.Getenv("NAVIGATION_ORDER")); raw != "" {
		parts := strings.Split(raw, ",")
		list := make([]string, 0, len(parts))
		for _, p := range parts {
			if p = strings.TrimSpace(p); p != "" {
				list = append(list, p)
			}
		}
		if len(list) > 0 {
			cfg.NavigationOrder = list
		}
	}

	// Process env overrides YAML when set (existing env-only deployments keep
	// working with zero migration; empty values count as unset).
	cfg.AppName = envOr("APP_NAME", cfg.AppName)
	cfg.AppEnv = envOr("APP_ENV", cfg.AppEnv)
	cfg.HTTPAddr = envOr("HTTP_ADDR", cfg.HTTPAddr)
	cfg.ReadTimeout = durationEnv("HTTP_READ_TIMEOUT", cfg.ReadTimeout)
	cfg.WriteTimeout = durationEnv("HTTP_WRITE_TIMEOUT", cfg.WriteTimeout)
	cfg.IdleTimeout = durationEnv("HTTP_IDLE_TIMEOUT", cfg.IdleTimeout)
	if v := strings.TrimSpace(os.Getenv("HTTP_SHUTDOWN_TIMEOUT")); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			// VP-021 contract §6: an unparsable drain budget is a startup
			// error (fail-closed), not a silent fallback.
			cfg.LoadError = fmt.Errorf("HTTP_SHUTDOWN_TIMEOUT: invalid duration %q (fail-closed)", v)
			return cfg
		}
		cfg.HTTPShutdownTimeout = d
	}
	// VP-021 contract §6: a non-positive drain budget is invalid in every
	// environment (fail-closed at startup via ValidateProd -> LoadError).
	if cfg.HTTPShutdownTimeout <= 0 {
		cfg.LoadError = fmt.Errorf("HTTP_SHUTDOWN_TIMEOUT must be > 0 (drain budget), got %s", cfg.HTTPShutdownTimeout)
		return cfg
	}
	if raw := strings.TrimSpace(os.Getenv("HTTP_CORS_ORIGINS")); raw != "" {
		cfg.HTTPCORSOrigins = splitCSV(raw)
	}
	if raw := strings.TrimSpace(os.Getenv("HTTP_TRUSTED_PROXIES")); raw != "" {
		cfg.HTTPTrustedProxies = splitCSV(raw)
	}
	cfg.LogLevelName = envOr("LOG_LEVEL", cfg.LogLevelName)
	cfg.AuthJWTSecret = envOr("AUTH_JWT_SECRET", cfg.AuthJWTSecret)
	cfg.AuthJWTSecretPrevious = envOr("AUTH_JWT_SECRET_PREVIOUS", cfg.AuthJWTSecretPrevious)
	cfg.AuthAccessTTL = durationEnv("AUTH_ACCESS_TTL", cfg.AuthAccessTTL)
	cfg.AuthRefreshTTL = durationEnv("AUTH_REFRESH_TTL", cfg.AuthRefreshTTL)
	cfg.AuthDevSessionEnabled = boolEnv("AUTH_DEV_SESSION_ENABLED", cfg.AuthDevSessionEnabled)
	cfg.AuthPublicBaseURL = strings.TrimRight(strings.TrimSpace(envOr("AUTH_PUBLIC_BASE_URL", cfg.AuthPublicBaseURL)), "/")
	cfg.DBPath = envOr("DB_PATH", cfg.DBPath)
	cfg.DBDialect = envOr("DB_DIALECT", cfg.DBDialect)
	cfg.DBDSN = envOr("DB_DSN", cfg.DBDSN)
	cfg.DBHost = envOr("DB_HOST", cfg.DBHost)
	cfg.DBPort = envOr("DB_PORT", cfg.DBPort)
	cfg.DBName = envOr("DB_NAME", cfg.DBName)
	cfg.DBUser = envOr("DB_USER", cfg.DBUser)
	cfg.DBPassword = envOr("DB_PASSWORD", cfg.DBPassword)
	cfg.DBSSLMode = envOr("DB_SSLMODE", cfg.DBSSLMode)
	cfg.DBPoolMaxOpen = nonNegIntEnv("DB_POOL_MAX_OPEN", cfg.DBPoolMaxOpen)
	cfg.DBPoolMaxIdle = nonNegIntEnv("DB_POOL_MAX_IDLE", cfg.DBPoolMaxIdle)
	cfg.DBConnLifetime = durationEnv("DB_CONN_MAX_LIFETIME", cfg.DBConnLifetime)
	cfg.AdminInitialPassword = envOr("ADMIN_INITIAL_PASSWORD", cfg.AdminInitialPassword)
	cfg.ProfileName = profile
	cfg.UploadAllowedTypes = envOr("UPLOAD_ALLOWED_TYPES", cfg.UploadAllowedTypes)
	cfg.UploadMaxFilesPerUser = positiveIntEnv("UPLOAD_MAX_FILES_PER_USER", cfg.UploadMaxFilesPerUser)
	cfg.UploadMaxBytesPerUser = positiveIntEnv("UPLOAD_MAX_BYTES_PER_USER", cfg.UploadMaxBytesPerUser)
	cfg.BrandingMaxBytes = positiveIntEnv("BRANDING_MAX_BYTES", cfg.BrandingMaxBytes)
	cfg.BrandingLogoMaxDimension = positiveIntEnv("BRANDING_LOGO_MAX_DIMENSION", cfg.BrandingLogoMaxDimension)
	cfg.BrandingFaviconDimension = positiveIntEnv("BRANDING_FAVICON_DIMENSION", cfg.BrandingFaviconDimension)
	cfg.BrandingJPEGQuality = positiveIntEnv("BRANDING_JPEG_QUALITY", cfg.BrandingJPEGQuality)
	cfg.ObjectsDriver = strings.ToLower(envOr("STORAGE_OBJECTS_DRIVER", cfg.ObjectsDriver))
	cfg.ObjectsLocalRoot = envOr("STORAGE_OBJECTS_LOCAL_ROOT", cfg.ObjectsLocalRoot)
	cfg.ObjectsS3Endpoint = envOr("STORAGE_OBJECTS_S3_ENDPOINT", cfg.ObjectsS3Endpoint)
	cfg.ObjectsS3Region = envOr("STORAGE_OBJECTS_S3_REGION", cfg.ObjectsS3Region)
	cfg.ObjectsS3Bucket = envOr("STORAGE_OBJECTS_S3_BUCKET", cfg.ObjectsS3Bucket)
	cfg.ObjectsS3AccessKeyID = envOr("STORAGE_OBJECTS_S3_ACCESS_KEY_ID", cfg.ObjectsS3AccessKeyID)
	cfg.ObjectsS3SecretAccessKey = envOr("STORAGE_OBJECTS_S3_SECRET_ACCESS_KEY", cfg.ObjectsS3SecretAccessKey)
	if v, set := os.LookupEnv("STORAGE_OBJECTS_S3_USE_PATH_STYLE"); set && strings.TrimSpace(v) != "" {
		cfg.ObjectsS3UsePathStyle = boolEnv("STORAGE_OBJECTS_S3_USE_PATH_STYLE", cfg.ObjectsS3UsePathStyle)
	}
	cfg.MailSMTPHost = envOr("MAIL_SMTP_HOST", cfg.MailSMTPHost)
	cfg.MailSMTPUsername = envOr("MAIL_SMTP_USERNAME", cfg.MailSMTPUsername)
	cfg.MailSMTPPassword = envOr("MAIL_SMTP_PASSWORD", cfg.MailSMTPPassword)
	cfg.MailSMTPFrom = envOr("MAIL_SMTP_FROM", cfg.MailSMTPFrom)
	if raw := strings.TrimSpace(os.Getenv("MAIL_CHANNEL")); raw != "" {
		cfg.MailChannel = strings.ToLower(raw)
	}
	cfg.MailResendAPIKey = envOr("MAIL_RESEND_API_KEY", cfg.MailResendAPIKey)
	cfg.MailResendFrom = envOr("MAIL_RESEND_FROM", cfg.MailResendFrom)
	cfg.MailMasterKeyPath = strings.TrimSpace(envOr("MAIL_MASTER_KEY_PATH", cfg.MailMasterKeyPath))
	cfg.MailConfigMasterKey = envOr("MAIL_CONFIG_MASTER_KEY", cfg.MailConfigMasterKey)
	// cache.max_entries (VP-026 / workspace-026 GOAL-003 D-001): strict env
	// parse mirroring MAIL_SMTP_PORT — an explicitly supplied invalid value
	// fails closed instead of silently keeping the default.
	if raw := strings.TrimSpace(os.Getenv("CACHE_MAX_ENTRIES")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			cfg.LoadError = fmt.Errorf("config: CACHE_MAX_ENTRIES must be a positive integer")
			return cfg
		}
		cfg.CacheMaxEntries = n
	}
	// eventbus.buffer_size (VP-028 / workspace-028 GOAL-003 D-001): strict env
	// parse. <= 0 is acceptable (falls back to DefaultEventBusBuffer); unparsable
	// or > MaxEventBusBuffer fails closed.
	if raw := strings.TrimSpace(os.Getenv("EVENTBUS_BUFFER_SIZE")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil {
			cfg.LoadError = fmt.Errorf("config: EVENTBUS_BUFFER_SIZE must be an integer")
			return cfg
		}
		if n > MaxEventBusBuffer {
			cfg.LoadError = fmt.Errorf("config: EVENTBUS_BUFFER_SIZE must be <= %d (got %d)", MaxEventBusBuffer, n)
			return cfg
		}
		cfg.EventBusBufferSize = n
	}
	if raw := strings.TrimSpace(os.Getenv("MAIL_SMTP_PORT")); raw != "" {
		port, err := strconv.Atoi(raw)
		if err != nil || port < 1 || port > 65535 {
			// Fail closed instead of silently keeping the default: a typo in
			// an explicitly supplied port must never degrade to 465.
			cfg.LoadError = fmt.Errorf("config: MAIL_SMTP_PORT must be a number between 1 and 65535")
			return cfg
		}
		cfg.MailSMTPPort = port
	}
	cfg.MetricsEnabled = boolEnv("OBSERVABILITY_METRICS_ENABLED", cfg.MetricsEnabled)
	cfg.MetricsAddr = envOr("OBSERVABILITY_METRICS_ADDR", cfg.MetricsAddr)
	cfg.MetricsAuthToken = envOr("OBSERVABILITY_METRICS_AUTH_TOKEN", cfg.MetricsAuthToken)
	cfg.TracesEnabled = boolEnv("OBSERVABILITY_TRACES_ENABLED", cfg.TracesEnabled)
	cfg.TracesEndpoint = envOr("OBSERVABILITY_TRACES_ENDPOINT", cfg.TracesEndpoint)
	if raw := strings.TrimSpace(os.Getenv("OBSERVABILITY_TRACES_SAMPLE_RATIO")); raw != "" {
		if ratio, err := strconv.ParseFloat(raw, 64); err == nil {
			cfg.TracesSampleRatio = ratio
		} // invalid value keeps the default; range is enforced by validation
	}
	if yf.Runtime.Mode != nil {
		cfg.RuntimeMode = RuntimeMode(strings.TrimSpace(*yf.Runtime.Mode))
	}
	if raw, set := os.LookupEnv("RUNTIME_MODE"); set {
		// Unlike ordinary optional overrides, an explicitly empty runtime mode is
		// invalid: silently falling back to normal would defeat fail-closed ops.
		cfg.RuntimeMode = RuntimeMode(strings.TrimSpace(raw))
	}
	if !ValidRuntimeMode(cfg.RuntimeMode) {
		cfg.LoadError = fmt.Errorf("config: runtime.mode must be one of normal, maintenance, degraded, read-only")
		return cfg
	}

	// db.dialect / db.dsn (VP-013 R1 v1.4 §5): empty = sqlite; unknown dialect
	// fails closed at load time (mirrors runtime.mode). The DSN/path pairing and
	// file-path-shape rules are enforced again in ValidateProd for every
	// environment, including development.
	cfg.DBDialect = strings.ToLower(strings.TrimSpace(cfg.DBDialect))
	switch cfg.DBDialect {
	case "", "sqlite":
		cfg.DBDialect = "sqlite"
	case "postgres":
		cfg.DBDialect = "postgres"
	default:
		cfg.LoadError = fmt.Errorf("config: db.dialect must be one of sqlite or postgres (got %q)", cfg.DBDialect)
		return cfg
	}

	// storage.objects.driver (VP-014 GOAL-002 D-001): empty = local; unknown
	// drivers fail closed at load. S3 pairing rules mirror the postgres
	// credential rule above and are re-checked by ValidateProd.
	switch cfg.ObjectsDriver {
	case "", "local":
		cfg.ObjectsDriver = "local"
		// A-002 F-001: report only the offending KEY NAME — never interpolate
		// the value (a secret may be the only s3 key set, and this string is
		// logged verbatim by the startup error path).
		if key := firstSetS3Key(cfg); key != "" {
			cfg.LoadError = localS3KeyMisconfig(key)
			return cfg
		}
	case "s3":
		if err := validateObjectsS3(cfg); err != nil {
			cfg.LoadError = err
			return cfg
		}
	default:
		cfg.LoadError = fmt.Errorf("config: storage.objects.driver must be one of local or s3 (got %q)", cfg.ObjectsDriver)
		return cfg
	}

	// cache.max_entries (VP-026 / workspace-026 GOAL-003 D-001): the in-memory
	// cache provider's bounded-entry budget must be positive. Explicit YAML /
	// env values are checked above (pointer mapping / strict env parse); this
	// block is the fail-closed net for any value that reached the struct.
	if cfg.CacheMaxEntries <= 0 {
		cfg.LoadError = fmt.Errorf("config: cache.max_entries must be a positive integer (got %d)", cfg.CacheMaxEntries)
		return cfg
	}

	// eventbus.buffer_size (VP-028 / workspace-028 GOAL-003 D-001): the in-memory
	// event-bus provider's per-subscription buffer size. Explicit YAML / env values
	// are checked above; <= 0 is acceptable (composition falls back to default);
	// this block is the fail-closed net for out-of-range values.
	if cfg.EventBusBufferSize > MaxEventBusBuffer {
		cfg.LoadError = fmt.Errorf("config: eventbus.buffer_size must be <= %d (got %d)", MaxEventBusBuffer, cfg.EventBusBufferSize)
		return cfg
	}

	// observability.metrics (VP-015 / workspace-015 GOAL-002 D-001): pairing
	// and exposure rules fail closed at load time, mirroring runtime.mode /
	// db.dialect / storage.objects.driver above.
	if err := validateMetrics(cfg.MetricsEnabled, cfg.MetricsAddr, cfg.MetricsAuthToken); err != nil {
		cfg.LoadError = err
		return cfg
	}
	// observability.traces (VP-015 / workspace-015 GOAL-004 D-001): same
	// fail-closed posture for the OTLP export surface.
	if err := validateTraces(cfg.TracesEnabled, cfg.TracesEndpoint, cfg.TracesSampleRatio); err != nil {
		cfg.LoadError = err
		return cfg
	}

	// Compose a postgres DSN from the exploded params when no explicit db.dsn
	// was configured. The password is a secret supplied via DB_PASSWORD (env /
	// configs/.env) — never a hardcoded YAML literal — and is required (fail
	// closed) here and re-checked in validateDB.
	if cfg.DBDialect == "postgres" && strings.TrimSpace(cfg.DBDSN) == "" {
		if strings.TrimSpace(cfg.DBHost) == "" || strings.TrimSpace(cfg.DBPort) == "" ||
			strings.TrimSpace(cfg.DBName) == "" || strings.TrimSpace(cfg.DBUser) == "" ||
			strings.TrimSpace(cfg.DBPassword) == "" {
			cfg.LoadError = fmt.Errorf("config: db.dialect=postgres requires db.host/port/name/user/password — provide the password via DB_PASSWORD env (configs/.env) or set an explicit db.dsn")
			return cfg
		}
		cfg.DBDSN = buildPostgresDSN(cfg)
	}

	// T-06 (GOAL-013 D-007): the enabled-modules set is YAML-only.
	// app.modules (preset name/path or inline list) is authoritative; the
	// app.profile builtin defaults apply when it is absent.
	resolved, err := resolveModulesFromYAML(yf.App.Modules, cfg.ProfileName)
	if err != nil {
		cfg.ProfileError = err
		return cfg
	}
	cfg.ProfileName = string(resolved.Name)
	cfg.ModulesEnabled = append([]string(nil), resolved.Modules...)
	cfg.ProfileSource = resolved.Source
	cfg.ProfilePrecedence = append([]string(nil), resolved.Precedence...)
	return cfg
}

// modulesPresetNames are the compiled built-in presets (D-007 §2: mvp /
// admin / demo; their module sets stay identical to kernel.ResolveProfile).
var modulesPresetNames = map[string]bool{"mvp": true, "admin": true, "demo": true}

// resolveModulesFromYAML resolves the enabled-module set from the app.modules
// YAML node (T-06 · GOAL-013 D-007):
//   - preset: <name>  — builtin compiled preset (mvp | admin | demo)
//   - preset: <path>  — custom preset YAML file declaring a modules: list
//   - list: [a, b]    — inline module ids (the "custom" form)
//
// preset and list are mutually exclusive; both present fails closed. A node
// kind of 0 means the section is absent — the app.profile builtin defaults
// apply.
func resolveModulesFromYAML(node yaml.Node, profileName string) (kernel.ProfileResolution, error) {
	if node.Kind == 0 {
		return kernel.ResolveProfile(profileName, nil)
	}
	if node.Kind != yaml.MappingNode {
		return kernel.ProfileResolution{}, fmt.Errorf("config: app.modules must be a mapping with preset or list")
	}
	var preset, list *yaml.Node
	for i := 0; i+1 < len(node.Content); i += 2 {
		key := node.Content[i]
		if key.Kind != yaml.ScalarNode {
			continue
		}
		switch key.Value {
		case "preset":
			preset = node.Content[i+1]
		case "list":
			list = node.Content[i+1]
		}
	}
	if preset != nil && list != nil {
		return kernel.ProfileResolution{}, fmt.Errorf("config: app.modules preset and list are mutually exclusive")
	}
	if preset != nil {
		name := strings.TrimSpace(preset.Value)
		if name == "" {
			return kernel.ProfileResolution{}, fmt.Errorf("config: app.modules.preset must not be empty")
		}
		var resolved kernel.ProfileResolution
		var err error
		if modulesPresetNames[name] {
			resolved, err = kernel.ResolveProfile(name, nil)
		} else {
			// Custom preset file: a YAML document declaring a modules: list.
			var modules []string
			modules, err = loadPresetFile(name)
			if err == nil {
				resolved, err = kernel.ResolveProfile(string(kernel.ProfileCustom), modules)
			}
		}
		if err != nil {
			return kernel.ProfileResolution{}, err
		}
		// The preset form is one YAML authority; surface it as modules.preset
		// regardless of whether the preset was a builtin or a custom file.
		resolved.Source = "modules.preset"
		resolved.Precedence = []string{"compiled-profile-default", "modules.preset"}
		return resolved, nil
	}
	if list != nil {
		if list.Kind != yaml.SequenceNode {
			return kernel.ProfileResolution{}, fmt.Errorf("config: app.modules.list must be a YAML list")
		}
		modules := make([]string, 0, len(list.Content))
		for _, item := range list.Content {
			if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
				return kernel.ProfileResolution{}, fmt.Errorf("config: app.modules.list entries must be strings")
			}
			modules = append(modules, item.Value)
		}
		if len(modules) == 0 {
			return kernel.ProfileResolution{}, fmt.Errorf("config: app.modules.list must not be empty")
		}
		return kernel.ResolveProfile(string(kernel.ProfileCustom), modules)
	}
	return kernel.ProfileResolution{}, fmt.Errorf("config: app.modules must declare preset or list")
}

// presetFile is the shape of a custom preset YAML file (app.modules.preset
// pointing at a path): a single modules: list.
type presetFile struct {
	Modules []string `yaml:"modules"`
}

// loadPresetFile reads a custom preset YAML file and returns its module list.
func loadPresetFile(path string) ([]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config: app.modules.preset %q: %w", path, err)
	}
	var pf presetFile
	dec := yaml.NewDecoder(strings.NewReader(string(raw)))
	dec.KnownFields(true)
	if err := dec.Decode(&pf); err != nil {
		if errors.Is(err, io.EOF) {
			return nil, fmt.Errorf("config: preset %q is empty", path)
		}
		return nil, fmt.Errorf("config: parse preset %q: %w", path, err)
	}
	if len(pf.Modules) == 0 {
		return nil, fmt.Errorf("config: preset %q declares no modules", path)
	}
	return pf.Modules, nil
}

func splitCSV(raw string) []string {
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if p = strings.TrimSpace(p); p != "" {
			out = append(out, p)
		}
	}
	return out
}

// parseNavigationOrder extracts a string sequence from the YAML navigation
// node. A missing/empty/null order yields nil (kernel default applies). A
// sequence containing non-scalar or non-string items is invalid: the whole
// order falls back to nil with a warning (GOAL-013 D-002 §4).
func parseNavigationOrder(node yaml.Node) []string {
	if node.Kind != yaml.SequenceNode {
		if node.Kind != 0 {
			slog.Warn("config: navigation.order is not a list; using the default navigation order")
		}
		return nil
	}
	list := make([]string, 0, len(node.Content))
	for _, item := range node.Content {
		if item.Kind != yaml.ScalarNode || item.Tag != "!!str" {
			slog.Warn("config: navigation.order contains a non-string entry; using the default navigation order")
			return nil
		}
		list = append(list, item.Value)
	}
	if len(list) == 0 {
		return nil // empty list = default
	}
	return list
}

// loadEnvFile reads the optional CONFIG_ENV_FILE (default configs/.env) of
// KEY=VALUE lines and exports them WITHOUT overriding an already-set process
// env. Missing default file is fine; an explicitly configured file that does
// not exist is a startup error (fail-closed).
func loadEnvFile(cfg *Config) error {
	path := strings.TrimSpace(os.Getenv("CONFIG_ENV_FILE"))
	explicit := path != ""
	if path == "" {
		path = "configs/.env"
	}
	raw, err := os.ReadFile(path)
	if err != nil {
		if explicit {
			return fmt.Errorf("CONFIG_ENV_FILE %q: %w", path, err)
		}
		if os.IsNotExist(err) {
			return nil // optional dev convenience file
		}
		return fmt.Errorf("read %s: %w", path, err)
	}
	for i, line := range strings.Split(string(raw), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.TrimPrefix(line, "export ")
		k, v, ok := strings.Cut(line, "=")
		if !ok || strings.TrimSpace(k) == "" {
			return fmt.Errorf("%s line %d: expected KEY=VALUE", path, i+1)
		}
		k = strings.TrimSpace(k)
		v = strings.TrimSpace(v)
		v = strings.Trim(v, `"'`)
		if _, exists := os.LookupEnv(k); exists {
			continue // process env wins; never override
		}
		if err := os.Setenv(k, v); err != nil {
			return fmt.Errorf("%s line %d: %w", path, i+1, err)
		}
	}
	return nil
}

// varPattern matches ${NAME} and ${NAME:-default}.
var varPattern = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?::-([^}]*))?\}`)

// inlineCommentIndex returns the index of the first YAML inline comment marker
// (" #" or a hash preceded by whitespace) that is NOT inside a quoted value
// (F-003: values like "My App #1" must survive). A hash without preceding
// whitespace is not a comment per YAML.
func inlineCommentIndex(line string) int {
	inSingle, inDouble := false, false
	escaped := false
	for i := 0; i < len(line); i++ {
		c := line[i]
		if escaped {
			escaped = false
			continue
		}
		if c == '\\' && inDouble {
			escaped = true
			continue
		}
		if c == '\'' && !inDouble {
			inSingle = !inSingle
			continue
		}
		if c == '"' && !inSingle {
			inDouble = !inDouble
			continue
		}
		if c == '#' && !inSingle && !inDouble {
			if i == 0 || line[i-1] == ' ' || line[i-1] == '\t' {
				return i
			}
		}
	}
	return -1
}

// interpolateAll expands every ${...} reference in the YAML text. A bare
// ${VAR} whose env is unset is a hard error (fail-closed); ${VAR:-default}
// falls back to default. Values are substituted verbatim (no re-parsing), so
// the env value may contain YAML characters.
//
// Interpolation is line-scoped: whole-line comments (starting with #) and
// inline comments (after " #") are excluded, so documentation examples such
// as "  ${VAR} -> fail-closed" never count as live references.
func interpolateAll(text string) (string, error) {
	lines := strings.Split(text, "\n")
	var missing string
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}
		if idx := inlineCommentIndex(line); idx >= 0 {
			line = line[:idx]
		}
		lines[i] = varPattern.ReplaceAllStringFunc(line, func(m string) string {
			if missing != "" {
				return m // already failing; keep the rest intact for the error
			}
			sub := varPattern.FindStringSubmatch(m)
			name := sub[1]
			if v, ok := os.LookupEnv(name); ok {
				return v
			}
			// ${VAR:-default}: the default follows ":-" inside the braces. An empty
			// default (${VAR:-}) is still a valid default, so detect the marker
			// rather than testing the captured group for emptiness.
			if _, def, hasDefault := strings.Cut(m, ":-"); hasDefault {
				return strings.TrimSuffix(def, "}")
			}
			missing = name
			return m
		})
	}
	if missing != "" {
		return "", fmt.Errorf("config interpolation: ${%s} is not set and has no default (fail-closed)", missing)
	}
	return strings.Join(lines, "\n"), nil
}

// ValidateProd fails startup for non-development environments when the static
// development session fallback is enabled (GOAL-008 A-005 F-002). The dev
// session substitutes a high-privilege static identity for unauthenticated
// requests, so it may only ever run in local development; any other APP_ENV
// with AUTH_DEV_SESSION_ENABLED=true is a hard startup error, not a warning.
//
// A-002 F-002-005 (GOAL-009 S5): non-development environments additionally
// require a JWT signing secret with a minimum length and both letters and
// digits, so a short or guessable HS256 key cannot silently start production.
// Development keeps the explicit low bar (documented insecure dev key).
//
// VP-016 R1 (workspace-016 Root D-002): an explicitly configured previous
// signing key (AUTH_JWT_SECRET_PREVIOUS) must satisfy the same strength rule
// and must differ from the current key; an absent previous keeps single-key
// behavior in every environment.
//
// W7: a LoadError (bad CONFIG_FILE, unset ${VAR}, invalid YAML) fails startup
// regardless of environment — configuration must never silently fall back.
func (c *Config) ValidateProd() error {
	if c.LoadError != nil {
		return fmt.Errorf("configuration load failed: %w", c.LoadError)
	}
	if c.ProfileError != nil {
		return fmt.Errorf("invalid module profile: %w", c.ProfileError)
	}
	// A zero-value Config is used by focused unit tests and means the loader
	// was bypassed; Load itself rejects an explicitly empty runtime mode.
	if c.RuntimeMode != "" && !ValidRuntimeMode(c.RuntimeMode) {
		return fmt.Errorf("invalid runtime mode %q", c.RuntimeMode)
	}
	if c.AppEnv == "" {
		return fmt.Errorf("APP_ENV must be set explicitly (development for local runs, production for deployments); refusing to guess")
	}
	// VP-013 R1 v1.4 §5: db dialect/DSN/path-shape rules are startup gates for
	// every environment (including development). An empty DBDialect / DBPath on
	// a zero-value or test Config means "use load defaults" and is skipped.
	if err := c.validateDB(); err != nil {
		return err
	}
	// VP-014 GOAL-002 D-001: object-storage pairing rules are startup gates
	// for every environment (including development), like the db rules.
	if err := c.validateObjects(); err != nil {
		return err
	}
	// VP-015 GOAL-002 D-001: metrics exposure/pairing rules are startup gates
	// for every environment (including development) as well.
	if err := c.validateObservability(); err != nil {
		return err
	}
	// VP-026 (workspace-026 GOAL-003 D-001): the cache entry budget is a
	// startup gate for every environment; zero on a bypassed/zero-value
	// Config means "use load defaults" and is skipped (mirrors the db rules).
	if c.CacheMaxEntries < 0 {
		return fmt.Errorf("cache.max_entries must be positive (got %d)", c.CacheMaxEntries)
	}
	// VP-028 (workspace-028 GOAL-003 D-001): the event-bus buffer size is a
	// startup gate for every environment; <= 0 means "use default" and is
	// skipped; > MaxEventBusBuffer is rejected.
	if c.EventBusBufferSize > MaxEventBusBuffer {
		return fmt.Errorf("eventbus.buffer_size must be <= %d (got %d)", MaxEventBusBuffer, c.EventBusBufferSize)
	}
	// VP-017 GOAL-003 D-001: mail.smtp pairing rules are startup gates for
	// every environment (including development), like the db/objects rules.
	if err := c.validateMail(); err != nil {
		return err
	}
	// W13 F-006 (GOAL-013 A-001): a malformed public base URL is a startup
	// gate in every environment — a bad value would corrupt every emailed
	// invitation link.
	if c.AuthPublicBaseURL != "" {
		parsed, perr := url.Parse(c.AuthPublicBaseURL)
		if perr != nil || (parsed.Scheme != "http" && parsed.Scheme != "https") || parsed.Host == "" ||
			parsed.Path != "" || parsed.RawQuery != "" || parsed.Fragment != "" ||
			strings.ContainsAny(c.AuthPublicBaseURL, " \t\r\n") {
			return fmt.Errorf("AUTH_PUBLIC_BASE_URL must be an absolute origin like https://host (no path, query or fragment); got %q", c.AuthPublicBaseURL)
		}
	}
	if c.AppEnv == "development" {
		return nil
	}
	if c.AuthDevSessionEnabled {
		return fmt.Errorf("AUTH_DEV_SESSION_ENABLED must be false when APP_ENV=%q", c.AppEnv)
	}
	// W15 F-002 (GOAL-016 A-001): the strength rule is exported so the public
	// serve surface (apps/api/server) enforces the same bar, single-source.
	if err := ValidateJWTSecretStrength("AUTH_JWT_SECRET", c.AuthJWTSecret); err != nil {
		return fmt.Errorf("%w when APP_ENV=%q", err, c.AppEnv)
	}
	// VP-016 R1 (workspace-016 Root D-002): an explicitly configured previous
	// signing key follows the same strength rule as the current key — during
	// the overlap window it still verifies signatures, so a weak previous key
	// widens the forgery surface exactly like a weak current one. A configured
	// previous equal to the current key is a no-op "rotation": fail closed so
	// operators cannot mistake it for a real overlap window. An absent
	// previous keeps single-key behavior in every environment.
	if prev := c.AuthJWTSecretPrevious; prev != "" {
		if err := ValidateJWTSecretStrength("AUTH_JWT_SECRET_PREVIOUS", prev); err != nil {
			return fmt.Errorf("%w when APP_ENV=%q", err, c.AppEnv)
		}
		if prev == c.AuthJWTSecret {
			return fmt.Errorf("AUTH_JWT_SECRET_PREVIOUS must differ from AUTH_JWT_SECRET (a no-op rotation is a misconfiguration)")
		}
	}
	return nil
}

// firstSetS3Key returns the NAME of the first non-empty storage.objects.s3.*
// field in a stable order, so misconfiguration errors can name the key
// without ever carrying a credential value into logs (A-002 F-001).
func firstSetS3Key(c *Config) string {
	for _, pair := range []struct{ key, value string }{
		{"endpoint", c.ObjectsS3Endpoint},
		{"region", c.ObjectsS3Region},
		{"bucket", c.ObjectsS3Bucket},
		{"access_key_id", c.ObjectsS3AccessKeyID},
		{"secret_access_key", c.ObjectsS3SecretAccessKey},
	} {
		if strings.TrimSpace(pair.value) != "" {
			return pair.key
		}
	}
	return ""
}

// localS3KeyMisconfig is the shared driver=local + s3.* misconfig error
// (A-002 R-001): one wording for Load and ValidateProd.
func localS3KeyMisconfig(key string) error {
	return fmt.Errorf("config: storage.objects.s3.%s is set but storage.objects.driver is local — set driver to s3 or remove the s3 keys", key)
}

// validateObjects enforces the storage.objects pairing rules (VP-014
// workspace-014 GOAL-002 D-001) on every Config that reaches ValidateProd,
// including zero-value/test configs that bypass Load's own checks.
func (c *Config) validateObjects() error {
	switch c.ObjectsDriver {
	case "", "local":
		// A-002 R-001: mirror Load's local-misconfig recheck so the contract
		// holds for configs that bypass Load, exactly like validateDB.
		if key := firstSetS3Key(c); key != "" {
			return localS3KeyMisconfig(key)
		}
	case "s3":
		return validateObjectsS3(c)
	default:
		return fmt.Errorf("config: storage.objects.driver must be one of local or s3 (got %q)", c.ObjectsDriver)
	}
	return nil
}

// validateObjectsS3 fails closed when an explicitly selected S3-compatible
// backend is missing any required setting. The secret must arrive via env
// interpolation (configs/.env or process env); a literal in YAML is possible
// but documented as forbidden — the same contract as DB_PASSWORD.
func validateObjectsS3(c *Config) error {
	for _, pair := range []struct{ name, value string }{
		{"endpoint", c.ObjectsS3Endpoint},
		{"bucket", c.ObjectsS3Bucket},
		{"access_key_id", c.ObjectsS3AccessKeyID},
		{"secret_access_key", c.ObjectsS3SecretAccessKey},
	} {
		if strings.TrimSpace(pair.value) == "" {
			return fmt.Errorf("config: storage.objects.driver=s3 requires storage.objects.s3.%s (provide secrets via STORAGE_OBJECTS_S3_* env / configs/.env)", pair.name)
		}
	}
	return nil
}

// MailSMTPConfigured reports whether the operator touched the mail.smtp block
// at all (workspace-017 GOAL-003 D-001). False keeps the embedded capture/log
// default; true means the fail-closed pairing rules below apply.
func (c *Config) MailSMTPConfigured() bool {
	return strings.TrimSpace(c.MailSMTPHost) != "" ||
		strings.TrimSpace(c.MailSMTPUsername) != "" ||
		strings.TrimSpace(c.MailSMTPPassword) != "" ||
		strings.TrimSpace(c.MailSMTPFrom) != "" ||
		c.MailSMTPPort != 0
}

// Outbound channel identifiers (workspace-017 GOAL-006 D-002 §1): the frozen
// first-wave set. The set is closed; a new provider requires a new decision.
const (
	MailChannelMock   = "mock"
	MailChannelResend = "resend"
	MailChannelSMTP   = "smtp"
)

// MailResendConfigured reports whether the operator touched the mail.resend
// block (workspace-017 GOAL-006 D-002 §4). Touching any key makes api-key and
// from REQUIRED (mirror of the SMTP pairing contract).
func (c *Config) MailResendConfigured() bool {
	return strings.TrimSpace(c.MailResendAPIKey) != "" ||
		strings.TrimSpace(c.MailResendFrom) != ""
}

// ResolveMailChannel implements the frozen resolution algorithm (workspace-017
// GOAL-006 D-002 §2):
//
//   - An explicit mail.channel wins. The named channel's block must carry its
//     required keys ("mock" always resolves); an incomplete production block
//     fails closed naming the missing key.
//   - Empty selector derives: exactly ONE fully configured production block
//     (resend / smtp) wins; BOTH fully configured is ambiguous -> fail closed
//     so the operator must state intent; NEITHER keeps the mock default,
//     which preserves the pre-R6 behavior of existing mail.smtp deployments.
//
// The error names configuration KEYS only, never values.
func (c *Config) ResolveMailChannel() (string, error) {
	switch strings.TrimSpace(c.MailChannel) {
	case "":
		// derive below
	case MailChannelMock:
		return MailChannelMock, nil
	case MailChannelResend:
		if err := c.validateResendBlock(); err != nil {
			return "", err
		}
		return MailChannelResend, nil
	case MailChannelSMTP:
		if !c.MailSMTPConfigured() {
			return "", fmt.Errorf("config: mail.channel=smtp requires an explicit mail.smtp.* block (provide secrets via MAIL_SMTP_* env / configs/.env)")
		}
		return MailChannelSMTP, nil
	default:
		return "", fmt.Errorf("config: mail.channel must be one of mock, resend, smtp (got %q)", c.MailChannel)
	}
	resendConfigured := c.MailResendConfigured()
	smtpConfigured := c.MailSMTPConfigured()
	switch {
	case resendConfigured && smtpConfigured:
		return "", fmt.Errorf("config: both mail.resend and mail.smtp are configured — set mail.channel explicitly to pick the outbound channel")
	case resendConfigured:
		if err := c.validateResendBlock(); err != nil {
			return "", err
		}
		return MailChannelResend, nil
	case smtpConfigured:
		return MailChannelSMTP, nil
	default:
		return MailChannelMock, nil
	}
}

// validateResendBlock enforces the mail.resend pairing contract
// (workspace-017 GOAL-006 D-002 §4): touching the block requires api-key and
// from; from must be a bare address. Errors name keys only.
func (c *Config) validateResendBlock() error {
	for _, pair := range []struct{ name, value string }{
		{"api-key", c.MailResendAPIKey},
		{"from", c.MailResendFrom},
	} {
		if strings.TrimSpace(pair.value) == "" {
			return fmt.Errorf("config: an explicit mail.resend block requires mail.resend.%s (provide secrets via MAIL_RESEND_* env / configs/.env)", pair.name)
		}
	}
	from := strings.TrimSpace(c.MailResendFrom)
	parsed, err := mail.ParseAddress(from)
	if err != nil || parsed.Address != from {
		return fmt.Errorf("config: mail.resend.from must be a bare address (got an unusable form)")
	}
	return nil
}

// validateMail enforces the outbound-mail contracts (VP-017 / workspace-017
// GOAL-003 D-001; R6 extension per GOAL-006 D-002 §2/§4) on every Config that
// reaches ValidateProd, including zero-value/test configs that bypass Load.
// The untouched surface is the valid embedded-default path — mock default,
// startup unaffected. Rules:
//
//   - A touched production block (mail.smtp OR mail.resend) must be COMPLETE
//     regardless of the selected channel — a half-filled block always fails
//     closed and names the first missing KEY (never a value).
//   - mail.channel must be one of mock|resend|smtp|"", and the frozen
//     resolution algorithm (ResolveMailChannel) must succeed: an explicit
//     channel needs its block, and two fully configured production blocks
//     without an explicit selector are ambiguous.
func (c *Config) validateMail() error {
	if c.MailSMTPConfigured() {
		for _, pair := range []struct{ name, value string }{
			{"host", c.MailSMTPHost},
			{"username", c.MailSMTPUsername},
			{"password", c.MailSMTPPassword},
			{"from", c.MailSMTPFrom},
		} {
			if strings.TrimSpace(pair.value) == "" {
				return fmt.Errorf("config: an explicit mail.smtp block requires mail.smtp.%s (provide secrets via MAIL_SMTP_* env / configs/.env)", pair.name)
			}
		}
		if c.MailSMTPPort < 0 || c.MailSMTPPort > 65535 {
			return fmt.Errorf("config: mail.smtp.port must be between 1 and 65535")
		}
		if from := strings.TrimSpace(c.MailSMTPFrom); from != "" {
			parsed, err := mail.ParseAddress(from)
			if err != nil || parsed.Address != from {
				return fmt.Errorf("config: mail.smtp.from must be a bare address (got an unusable form)")
			}
		}
	}
	if c.MailResendConfigured() {
		if err := c.validateResendBlock(); err != nil {
			return err
		}
	}
	if _, err := c.ResolveMailChannel(); err != nil {
		return err
	}
	return nil
}

// minMetricsTokenLen is the minimum length for an observability.metrics
// bearer token (VP-015 / workspace-015 GOAL-002 D-001 §2). A single-factor
// gate deserves an entropy floor; the JWT secret keeps its own stricter bar.
const minMetricsTokenLen = 16

// validateObservability enforces the observability.metrics contract on every
// Config that reaches ValidateProd, including zero-value/test configs that
// bypass Load's own checks. An untouched surface (disabled, no addr, no
// token) is skipped so focused unit Configs keep working.
func (c *Config) validateObservability() error {
	metricsUntouched := !c.MetricsEnabled && strings.TrimSpace(c.MetricsAddr) == "" && strings.TrimSpace(c.MetricsAuthToken) == ""
	tracesUntouched := !c.TracesEnabled && strings.TrimSpace(c.TracesEndpoint) == "" && c.TracesSampleRatio == 0
	if metricsUntouched && tracesUntouched {
		return nil
	}
	addr := c.MetricsAddr
	if strings.TrimSpace(addr) == "" {
		addr = "127.0.0.1:25081" // zero-value Config keeps the documented default
	}
	if err := validateMetrics(c.MetricsEnabled, addr, c.MetricsAuthToken); err != nil {
		return err
	}
	return validateTraces(c.TracesEnabled, c.TracesEndpoint, c.TracesSampleRatio)
}

// validateMetrics enforces the observability.metrics contract (VP-015 /
// workspace-015 GOAL-002 D-001 §2): a dead auth_token is rejected (mirroring
// the s3-keys-with-local-driver precedent), the bind address must parse as
// host:port with a numeric port, set tokens need a minimal length, and any
// non-loopback bind requires a bearer token in EVERY environment (exposure
// risk does not depend on APP_ENV). Error messages name keys only — never a
// token value.
func validateMetrics(enabled bool, addr, token string) error {
	token = strings.TrimSpace(token)
	if !enabled {
		if token != "" {
			return fmt.Errorf("config: observability.metrics.auth_token is set but observability.metrics.enabled is false — enable metrics or remove the token")
		}
		return nil
	}
	host, port, err := net.SplitHostPort(strings.TrimSpace(addr))
	if err != nil {
		return fmt.Errorf("config: observability.metrics.addr %q must be host:port: %v", addr, err)
	}
	portNum, err := strconv.Atoi(port)
	if err != nil || portNum < 1 || portNum > 65535 {
		return fmt.Errorf("config: observability.metrics.addr %q must use a numeric port in 1-65535", addr)
	}
	if token != "" && len(token) < minMetricsTokenLen {
		return fmt.Errorf("config: observability.metrics.auth_token must be at least %d characters", minMetricsTokenLen)
	}
	if !isLoopbackHost(host) && token == "" {
		return fmt.Errorf("config: observability.metrics.addr %q binds beyond loopback — set observability.metrics.auth_token via OBSERVABILITY_METRICS_AUTH_TOKEN (env / configs/.env) before exposing metrics", addr)
	}
	return nil
}

// isLoopbackHost reports whether a host:port host part binds only the local
// loopback interface. An empty host (":25081") means ALL interfaces and is
// therefore NOT loopback; only "localhost" and loopback IPs qualify.
func isLoopbackHost(host string) bool {
	if strings.EqualFold(host, "localhost") {
		return true
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip.IsLoopback()
	}
	return false
}

// validateTraces enforces the observability.traces contract (VP-015 /
// workspace-015 GOAL-004 D-001 §2): a dead endpoint is rejected (mirroring
// the metrics token precedent), an enabled surface requires an absolute
// http(s) endpoint, and an explicitly provided sample ratio must lie in
// (0,1]. A zero ratio on hand-built/test Configs means "not set" and is
// tolerated (the loader default is 1.0). Error messages never carry secrets.
func validateTraces(enabled bool, endpoint string, ratio float64) error {
	endpoint = strings.TrimSpace(endpoint)
	if !enabled {
		if endpoint != "" {
			return fmt.Errorf("config: observability.traces.endpoint is set but observability.traces.enabled is false — enable traces or remove the endpoint")
		}
		return nil
	}
	if endpoint == "" {
		return fmt.Errorf("config: observability.traces.endpoint is required when observability.traces.enabled is true (OTLP/HTTP base URL, e.g. http://localhost:4318)")
	}
	u, err := url.Parse(endpoint)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return fmt.Errorf("config: observability.traces.endpoint %q must be an absolute http(s) URL", endpoint)
	}
	if ratio < 0 || ratio > 1 {
		return fmt.Errorf("config: observability.traces.sample_ratio must be within (0, 1] (got %v)", ratio)
	}
	return nil
}

// validateDB enforces the VP-013 R1 v1.4 §5 pairing and path-shape rules.
// An empty DBDialect means sqlite (zero-value / test Config); an empty DBPath
// means "load default applies" and is skipped so focused unit Configs keep
// working. Real configs from Load always carry non-empty defaults.
func (c *Config) validateDB() error {
	dialect := strings.ToLower(strings.TrimSpace(c.DBDialect))
	if dialect == "" {
		dialect = "sqlite"
	}
	switch dialect {
	case "sqlite":
		if strings.TrimSpace(c.DBDSN) != "" {
			return fmt.Errorf("config: db.dsn must be empty when db.dialect is sqlite (DB_DSN must not be set)")
		}
		if strings.TrimSpace(c.DBPath) != "" {
			return validateDBPathShape(c.DBPath)
		}
	case "postgres":
		if strings.TrimSpace(c.DBDSN) == "" {
			return fmt.Errorf("config: db.dsn is required when db.dialect is postgres (DB_DSN must be set)")
		}
		if strings.TrimSpace(c.DBPath) != "" {
			return validateDBPathShape(c.DBPath)
		}
	default:
		return fmt.Errorf("config: db.dialect must be one of sqlite or postgres (got %q)", c.DBDialect)
	}
	return nil
}

// validateDBPathShape rejects values that are not a file path with an
// extension, so filepath.Dir(db.path) derives a stable file-storage root
// (R1 v1.4 §5 predicates 1–6). The default ./data/schema-ui.db passes; a
// trailing separator or an existing directory (e.g. ./data) is rejected.
func validateDBPathShape(path string) error {
	trimmed := strings.TrimSpace(path)
	if trimmed == "" {
		return fmt.Errorf("config: db.path must not be empty")
	}
	// Predicate 1: reject a trailing separator (/ \ or filepath.Separator).
	if strings.HasSuffix(trimmed, "/") || strings.HasSuffix(trimmed, "\\") {
		return fmt.Errorf("config: db.path %q must be a file path, not a directory (reject trailing separator)", trimmed)
	}
	// Predicate 2: an already-existing path must not be a directory.
	if st, err := os.Stat(trimmed); err == nil {
		if st.IsDir() {
			return fmt.Errorf("config: db.path %q is an existing directory, want a database file path", trimmed)
		}
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("config: stat db.path %q: %w", trimmed, err)
	}
	// Predicate 3: base must be non-empty and not "." or "..".
	base := filepath.Base(trimmed)
	if base == "" || base == "." || base == ".." {
		return fmt.Errorf("config: db.path %q must name a file, not a directory", trimmed)
	}
	// Predicate 4/5: don't reject on Dir=="." alone (cwd-relative files are
	// valid); the file-shape guarantee comes from the extension predicate.
	// Predicate 6: a non-empty extension is required (blocks a not-yet-existing
	// "./data" / ".\\data" directory component masquerading as a file path).
	if ext := filepath.Ext(base); ext == "" {
		return fmt.Errorf("config: db.path %q must include a file extension (e.g. .db)", trimmed)
	}
	return nil
}

// buildPostgresDSN composes a postgres URL DSN from the exploded connection
// params. User/password are URL-escaped so special characters in the secret
// survive; sslmode is applied as a query param (default disable).
func buildPostgresDSN(c *Config) string {
	u := url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(c.DBUser, c.DBPassword),
		Host:   net.JoinHostPort(c.DBHost, c.DBPort),
		Path:   "/" + c.DBName,
	}
	q := u.Query()
	q.Set("sslmode", c.DBSSLMode)
	u.RawQuery = q.Encode()
	return u.String()
}

// minJWTSecretLen is the minimum HS256 signing-key length enforced outside
// development (A-002 F-002-005).
const minJWTSecretLen = 32

// ValidateJWTSecretStrength enforces the production HS256 signing-key bar
// shared by the main configuration and the public serve surface (W15 F-002,
// GOAL-016 A-001): a minimum length plus both letters and digits, so a short
// or guessable key cannot silently start signing tokens outside development.
// The error names the key but never carries the secret value.
func ValidateJWTSecretStrength(keyName, secret string) error {
	if len(secret) < minJWTSecretLen {
		return fmt.Errorf("%s must be at least %d characters", keyName, minJWTSecretLen)
	}
	if !containsLettersAndDigits(secret) {
		return fmt.Errorf("%s must contain both letters and digits", keyName)
	}
	return nil
}

// containsLettersAndDigits rejects all-digit / all-letter / single-class keys
// (weak entropy) while keeping the rule simple and checkable.
func containsLettersAndDigits(s string) bool {
	hasLetter := false
	hasDigit := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') {
			hasLetter = true
		}
		if r >= '0' && r <= '9' {
			hasDigit = true
		}
		if hasLetter && hasDigit {
			return true
		}
	}
	return false
}

// LogLevel maps LOG_LEVEL to slog.
func (c *Config) LogLevel() slog.Level {
	switch strings.ToLower(c.LogLevelName) {
	case "debug":
		return slog.LevelDebug
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
func envOr(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func durationEnv(key string, fallback time.Duration) time.Duration {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

// strPtrOr returns the dereferenced value or the fallback when the pointer is
// nil (key omitted in YAML, F-002).
func strPtrOr(v *string, fallback string) string {
	if v == nil {
		return fallback
	}
	return *v
}

// orDuration parses a YAML duration string; on empty/invalid it returns the
// caller-provided default (same leniency as durationEnv).
func orDuration(v string, fallback time.Duration) time.Duration {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return fallback
	}
	return d
}

// orDurationPtr is orDuration for a possibly-omitted YAML key.
func orDurationPtr(v *string, fallback time.Duration) time.Duration {
	if v == nil {
		return fallback
	}
	return orDuration(*v, fallback)
}

// strictDurationPtr parses a YAML duration string strictly (VP-021 contract
// §6 / config http.shutdown_timeout): a nil pointer keeps the fallback; an
// explicitly set empty or unparsable value is a startup error (fail-closed),
// unlike the lenient orDurationPtr used by request-level timeouts.
func strictDurationPtr(v *string, key string, fallback time.Duration) (time.Duration, error) {
	if v == nil {
		return fallback, nil
	}
	s := strings.TrimSpace(*v)
	if s == "" {
		return fallback, fmt.Errorf("%s: must not be empty when set (fail-closed)", key)
	}
	d, err := time.ParseDuration(s)
	if err != nil {
		return fallback, fmt.Errorf("%s: invalid duration %q (fail-closed)", key, s)
	}
	return d, nil
}

func boolEnv(key string, fallback bool) bool {
	v := strings.TrimSpace(os.Getenv(key))
	if v == "" {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}

// positiveIntEnv mirrors the upload handler's old envPositiveInt semantics:
// unset or invalid keeps the fallback (never a zero/negative quota).
func positiveIntEnv(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n <= 0 {
		return fallback
	}
	return n
}

// nonNegIntEnv mirrors positiveIntEnv but accepts 0 (pool knobs use 0 as
// "driver default / no cap"); unset or invalid keeps the fallback.
func nonNegIntEnv(key string, fallback int) int {
	raw := strings.TrimSpace(os.Getenv(key))
	if raw == "" {
		return fallback
	}
	n, err := strconv.Atoi(raw)
	if err != nil || n < 0 {
		return fallback
	}
	return n
}
