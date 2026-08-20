package config

import (
	_ "embed"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/magicvr/schema-ui-core/apps/api/internal/kernel"
)

//go:embed config.default.yaml
var defaultConfigYAML []byte

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

	AuthJWTSecret  string
	AuthAccessTTL  time.Duration
	AuthRefreshTTL time.Duration
	DBPath         string
	// DBDialect is the store dialect (VP-013 / R1 v1.4 §5): "" or "sqlite" or
	// "postgres". Load normalizes empty to "sqlite"; ValidateProd rejects
	// unknown values and enforces DSN/path pairing rules.
	DBDialect string
	// DBDSN is the postgres SQL connection string (DB_DSN; no default). It must
	// be empty when DBDialect is sqlite and non-empty for postgres.
	DBDSN                 string
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
		Addr           *string `yaml:"addr"`
		ReadTimeout    *string `yaml:"read_timeout"`
		WriteTimeout   *string `yaml:"write_timeout"`
		IdleTimeout    *string `yaml:"idle_timeout"`
		CORSOrigins    *string `yaml:"cors_origins"`
		TrustedProxies *string `yaml:"trusted_proxies"`
	} `yaml:"http"`
	Log struct {
		Level *string `yaml:"level"`
	} `yaml:"log"`
	Auth struct {
		JWTSecret         *string `yaml:"jwt_secret"`
		AccessTTL         *string `yaml:"access_ttl"`
		RefreshTTL        *string `yaml:"refresh_ttl"`
		DevSessionEnabled *bool   `yaml:"dev_session_enabled"`
	} `yaml:"auth"`
	DB struct {
		Path    *string `yaml:"path"`
		Dialect *string `yaml:"dialect"`
		DSN     *string `yaml:"dsn"`
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
// interpolation; it never overrides an already-set process env. Load never
// returns an error: fatal load failures land in LoadError (and therefore
// ValidateProd) so existing call sites keep working.
func Load() *Config {
	cfg := &Config{
		AppName:      "schema-ui-core-api",
		AppEnv:       "",
		HTTPAddr:     ":25080",
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
		LogLevelName: "info",

		AuthJWTSecret:            "",
		AuthAccessTTL:            15 * time.Minute,
		AuthRefreshTTL:           30 * 24 * time.Hour,
		DBPath:                   "./data/schema-ui.db",
		DBDialect:                "sqlite",
		DBDSN:                    "",
		AdminInitialPassword:     "",
		AuthDevSessionEnabled:    false,
		UploadMaxFilesPerUser:    1000,
		UploadMaxBytesPerUser:    256 << 20,
		BrandingLogoMaxDimension: 512,
		BrandingFaviconDimension: 64,
		BrandingJPEGQuality:      82,
		BrandingMaxBytes:         4 << 20,
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
	if yf.HTTP.CORSOrigins != nil {
		cfg.HTTPCORSOrigins = splitCSV(*yf.HTTP.CORSOrigins)
	}
	if yf.HTTP.TrustedProxies != nil {
		cfg.HTTPTrustedProxies = splitCSV(*yf.HTTP.TrustedProxies)
	}
	cfg.LogLevelName = strPtrOr(yf.Log.Level, cfg.LogLevelName)
	cfg.AuthJWTSecret = strPtrOr(yf.Auth.JWTSecret, cfg.AuthJWTSecret)
	cfg.AuthAccessTTL = orDurationPtr(yf.Auth.AccessTTL, cfg.AuthAccessTTL)
	cfg.AuthRefreshTTL = orDurationPtr(yf.Auth.RefreshTTL, cfg.AuthRefreshTTL)
	if yf.Auth.DevSessionEnabled != nil {
		cfg.AuthDevSessionEnabled = *yf.Auth.DevSessionEnabled
	}
	cfg.DBPath = strPtrOr(yf.DB.Path, cfg.DBPath)
	cfg.DBDialect = strPtrOr(yf.DB.Dialect, cfg.DBDialect)
	cfg.DBDSN = strPtrOr(yf.DB.DSN, cfg.DBDSN)
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
	if raw := strings.TrimSpace(os.Getenv("HTTP_CORS_ORIGINS")); raw != "" {
		cfg.HTTPCORSOrigins = splitCSV(raw)
	}
	if raw := strings.TrimSpace(os.Getenv("HTTP_TRUSTED_PROXIES")); raw != "" {
		cfg.HTTPTrustedProxies = splitCSV(raw)
	}
	cfg.LogLevelName = envOr("LOG_LEVEL", cfg.LogLevelName)
	cfg.AuthJWTSecret = envOr("AUTH_JWT_SECRET", cfg.AuthJWTSecret)
	cfg.AuthAccessTTL = durationEnv("AUTH_ACCESS_TTL", cfg.AuthAccessTTL)
	cfg.AuthRefreshTTL = durationEnv("AUTH_REFRESH_TTL", cfg.AuthRefreshTTL)
	cfg.AuthDevSessionEnabled = boolEnv("AUTH_DEV_SESSION_ENABLED", cfg.AuthDevSessionEnabled)
	cfg.DBPath = envOr("DB_PATH", cfg.DBPath)
	cfg.DBDialect = envOr("DB_DIALECT", cfg.DBDialect)
	cfg.DBDSN = envOr("DB_DSN", cfg.DBDSN)
	cfg.AdminInitialPassword = envOr("ADMIN_INITIAL_PASSWORD", cfg.AdminInitialPassword)
	cfg.ProfileName = profile
	cfg.UploadAllowedTypes = envOr("UPLOAD_ALLOWED_TYPES", cfg.UploadAllowedTypes)
	cfg.UploadMaxFilesPerUser = positiveIntEnv("UPLOAD_MAX_FILES_PER_USER", cfg.UploadMaxFilesPerUser)
	cfg.UploadMaxBytesPerUser = positiveIntEnv("UPLOAD_MAX_BYTES_PER_USER", cfg.UploadMaxBytesPerUser)
	cfg.BrandingMaxBytes = positiveIntEnv("BRANDING_MAX_BYTES", cfg.BrandingMaxBytes)
	cfg.BrandingLogoMaxDimension = positiveIntEnv("BRANDING_LOGO_MAX_DIMENSION", cfg.BrandingLogoMaxDimension)
	cfg.BrandingFaviconDimension = positiveIntEnv("BRANDING_FAVICON_DIMENSION", cfg.BrandingFaviconDimension)
	cfg.BrandingJPEGQuality = positiveIntEnv("BRANDING_JPEG_QUALITY", cfg.BrandingJPEGQuality)
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
	if c.AppEnv == "development" {
		return nil
	}
	if c.AuthDevSessionEnabled {
		return fmt.Errorf("AUTH_DEV_SESSION_ENABLED must be false when APP_ENV=%q", c.AppEnv)
	}
	if len(c.AuthJWTSecret) < minJWTSecretLen {
		return fmt.Errorf(
			"AUTH_JWT_SECRET must be at least %d characters when APP_ENV=%q",
			minJWTSecretLen, c.AppEnv,
		)
	}
	if !containsLettersAndDigits(c.AuthJWTSecret) {
		return fmt.Errorf(
			"AUTH_JWT_SECRET must contain both letters and digits when APP_ENV=%q",
			c.AppEnv,
		)
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

// minJWTSecretLen is the minimum HS256 signing-key length enforced outside
// development (A-002 F-002-005).
const minJWTSecretLen = 32

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
