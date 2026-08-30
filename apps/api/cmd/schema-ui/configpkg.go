// configpkg —— schema-ui config 子命令族（VP-025 R2：export / diff）。
//
// 配置包合同 v0.1.0（workspace-025 GOAL-002 D-002）为责任分母：
//   - export：serve 壳配置树（内嵌默认 ∪ 显式文件）→ 配置包 v1
//     （package 元数据 + config 非敏感键 + secrets.exclude；env 引用保留
//     ${VAR} 形态不解析；敏感键剔除并记录所需 env）。
//   - diff：两包 / 包 vs 配置（--against）的键级差量（add/modify/remove +
//     路径 + old/new）；忽略 package 信息性元数据；退出码 0 无差 / 1 有差 /
//     2 错误（cliError）。
//   - dry-run / import：仅注册（R3 实现）。
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/magicvr/schema-ui-core/apps/api/server"
	"gopkg.in/yaml.v3"
)

// cliError 携带进程退出码（合同 §2.1：export 0/1；diff 0/1/2；错误一律可由
// 命令指定，默认 exit(1)）。
type cliError struct {
	code int
	err  error
}

func (e *cliError) Error() string { return e.err.Error() }
func (e *cliError) Unwrap() error { return e.err }

func errCli(code int, format string, args ...any) error {
	return &cliError{code: code, err: fmt.Errorf(format, args...)}
}

const (
	configPackageFormat = "schema-ui-config-package"
	configPackageVer    = 1
)

// ---- 包结构（导出树 = 与 config.default.yaml 树形同构的非敏感结构键） ----

// cfgTree 是配置包 v1 的 config 段（固定键序 —— 段序与 struct 字段序一致；
// omitempty：敏感键剔除后与空键不出现在产物中，缺失键 = 内嵌默认）。
// yaml/json tag 同构：两种输出面共享同一键命名。
type cfgTree struct {
	App struct {
		Name string `yaml:"name,omitempty" json:"name,omitempty"`
		Env  string `yaml:"env,omitempty" json:"env,omitempty"`
	} `yaml:"app,omitempty" json:"app,omitempty"`
	HTTP struct {
		Addr            string   `yaml:"addr,omitempty" json:"addr,omitempty"`
		ReadTimeout     string   `yaml:"read_timeout,omitempty" json:"read_timeout,omitempty"`
		WriteTimeout    string   `yaml:"write_timeout,omitempty" json:"write_timeout,omitempty"`
		IdleTimeout     string   `yaml:"idle_timeout,omitempty" json:"idle_timeout,omitempty"`
		ShutdownTimeout string   `yaml:"shutdown_timeout,omitempty" json:"shutdown_timeout,omitempty"`
		TrustedProxies  []string `yaml:"trusted_proxies,omitempty" json:"trusted_proxies,omitempty"`
		CORSOrigins     []string `yaml:"cors_origins,omitempty" json:"cors_origins,omitempty"`
	} `yaml:"http,omitempty" json:"http,omitempty"`
	DB struct {
		Dialect string `yaml:"dialect,omitempty" json:"dialect,omitempty"`
		Path    string `yaml:"path,omitempty" json:"path,omitempty"`
		DSN     string `yaml:"dsn,omitempty" json:"dsn,omitempty"`
	} `yaml:"db,omitempty" json:"db,omitempty"`
	Profile string `yaml:"profile,omitempty" json:"profile,omitempty"`
	Auth    struct {
		JWTSecret     string `yaml:"jwt_secret,omitempty" json:"jwt_secret,omitempty"`
		AccessTTL     string `yaml:"access_ttl,omitempty" json:"access_ttl,omitempty"`
		RefreshTTL    string `yaml:"refresh_ttl,omitempty" json:"refresh_ttl,omitempty"`
		PublicBaseURL string `yaml:"public_base_url,omitempty" json:"public_base_url,omitempty"`
	} `yaml:"auth,omitempty" json:"auth,omitempty"`
	Admin struct {
		InitialPassword string `yaml:"initial_password,omitempty" json:"initial_password,omitempty"`
	} `yaml:"admin,omitempty" json:"admin,omitempty"`
	Log struct {
		Level string `yaml:"level,omitempty" json:"level,omitempty"`
	} `yaml:"log,omitempty" json:"log,omitempty"`
}

type pkgMeta struct {
	Format     string `yaml:"format" json:"format"`
	Version    int    `yaml:"version" json:"version"`
	App        string `yaml:"app" json:"app"`
	Env        string `yaml:"env" json:"env"`
	Profile    string `yaml:"profile" json:"profile"`
	ExportedAt string `yaml:"exported_at" json:"exported_at"`
}

type secretEntry struct {
	Key string `yaml:"key" json:"key"`
	Env string `yaml:"env" json:"env"`
}

type configPackage struct {
	Package pkgMeta `yaml:"package" json:"package"`
	Config  cfgTree `yaml:"config" json:"config"`
	Secrets struct {
		Exclude []secretEntry `yaml:"exclude" json:"exclude"`
	} `yaml:"secrets" json:"secrets"`
}

// treeFile 是源配置 YAML 的已知字段面（与 server.yamlFile 同构；指针字段
// 保留「键是否显式出现」，供 默认 ∪ 显式 合并）。
type treeFile struct {
	App struct {
		Name *string `yaml:"name"`
		Env  *string `yaml:"env"`
	} `yaml:"app"`
	HTTP struct {
		Addr            *string  `yaml:"addr"`
		ReadTimeout     *string  `yaml:"read_timeout"`
		WriteTimeout    *string  `yaml:"write_timeout"`
		IdleTimeout     *string  `yaml:"idle_timeout"`
		ShutdownTimeout *string  `yaml:"shutdown_timeout"`
		TrustedProxies  []string `yaml:"trusted_proxies"`
		CORSOrigins     []string `yaml:"cors_origins"`
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

// sensitiveNameRe 是实现宽规则的保守匹配（合同 §1：键名含 secret/password/
// token 即视为敏感）。
var sensitiveNameRe = regexp.MustCompile(`(?i)secret|password|token`)

// envRefRe 从源值文本中提取首个 ${VAR} 引用的变量名（${VAR} 或 ${VAR:-def}）。
var envRefRe = regexp.MustCompile(`\$\{([A-Za-z_][A-Za-z0-9_]*)(?::-|})`)

func envNameFrom(value string) string {
	if m := envRefRe.FindStringSubmatch(value); m != nil {
		return m[1]
	}
	return ""
}

// parseTreeFile 用 KnownFields 严格解析一段源配置 YAML（未知键 / 多文档拒绝，
// 与 server.LoadConfig 同纪律）。
func parseTreeFile(raw []byte) (treeFile, error) {
	var tf treeFile
	dec := yaml.NewDecoder(strings.NewReader(string(raw)))
	dec.KnownFields(true)
	if err := dec.Decode(&tf); err != nil {
		if errors.Is(err, io.EOF) {
			return tf, errors.New("config: source YAML is empty")
		}
		return tf, fmt.Errorf("config: parse source YAML: %w", err)
	}
	var extra treeFile
	if err := dec.Decode(&extra); err != io.EOF {
		return tf, errors.New("config: multiple YAML documents are not supported")
	}
	return tf, nil
}

// defaultTreeFile 解析内嵌 serve 默认配置（go:embed）为已知字段树。
func defaultTreeFile() (treeFile, error) {
	return parseTreeFile(server.DefaultConfigYAML())
}

// buildExportTree 把源配置（内嵌默认 ∪ 显式文件按键覆盖）合并为导出树，
// 敏感键剔除并收集 secrets.exclude（所需 env 名从源值 ${VAR} 提取）。
//
// path 为空 → 纯内嵌默认（保留 ${VAR} 形态）；path 非空 → 文件必须可解析。
func buildExportTree(path string) (cfgTree, []secretEntry, error) {
	def, err := defaultTreeFile()
	if err != nil {
		return cfgTree{}, nil, err
	}
	var file treeFile
	if strings.TrimSpace(path) != "" {
		b, err := os.ReadFile(strings.TrimSpace(path))
		if err != nil {
			return cfgTree{}, nil, fmt.Errorf("config export: read %q: %w", path, err)
		}
		file, err = parseTreeFile(b)
		if err != nil {
			return cfgTree{}, nil, err
		}
	}

	tree := cfgTree{}
	tree.App.Name = strOr(file.App.Name, def.App.Name)
	tree.App.Env = strOr(file.App.Env, def.App.Env)
	tree.HTTP.Addr = strOr(file.HTTP.Addr, def.HTTP.Addr)
	tree.HTTP.ReadTimeout = strOr(file.HTTP.ReadTimeout, def.HTTP.ReadTimeout)
	tree.HTTP.WriteTimeout = strOr(file.HTTP.WriteTimeout, def.HTTP.WriteTimeout)
	tree.HTTP.IdleTimeout = strOr(file.HTTP.IdleTimeout, def.HTTP.IdleTimeout)
	tree.HTTP.ShutdownTimeout = strOr(file.HTTP.ShutdownTimeout, def.HTTP.ShutdownTimeout)
	tree.HTTP.TrustedProxies = listOr(file.HTTP.TrustedProxies, def.HTTP.TrustedProxies)
	tree.HTTP.CORSOrigins = listOr(file.HTTP.CORSOrigins, def.HTTP.CORSOrigins)
	tree.DB.Dialect = strOr(file.DB.Dialect, def.DB.Dialect)
	tree.DB.Path = strOr(file.DB.Path, def.DB.Path)
	tree.DB.DSN = strOr(file.DB.DSN, def.DB.DSN)
	tree.Profile = strOr(file.Profile, def.Profile)
	tree.Auth.JWTSecret = strOr(file.Auth.JWTSecret, def.Auth.JWTSecret)
	tree.Auth.AccessTTL = strOr(file.Auth.AccessTTL, def.Auth.AccessTTL)
	tree.Auth.RefreshTTL = strOr(file.Auth.RefreshTTL, def.Auth.RefreshTTL)
	tree.Auth.PublicBaseURL = strOr(file.Auth.PublicBaseURL, def.Auth.PublicBaseURL)
	tree.Admin.InitialPassword = strOr(file.Admin.InitialPassword, def.Admin.InitialPassword)
	tree.Log.Level = strOr(file.Log.Level, def.Log.Level)

	// 敏感键剔除（宽规则匹配字段名）+ secrets.exclude（键路径 + 所需 env）。
	// sensitiveFields 登记表 = 当前 serve 面敏感字段全集；新增字段若命中宽规则
	// （名字含 secret/password/token）必须在此登记，否则不变量断言失败。
	type sensitiveField struct {
		path       string
		srcFile    *string
		srcDefault *string
	}
	fields := []sensitiveField{
		{path: "auth.jwt_secret", srcFile: file.Auth.JWTSecret, srcDefault: def.Auth.JWTSecret},
		{path: "admin.initial_password", srcFile: file.Admin.InitialPassword, srcDefault: def.Admin.InitialPassword},
	}
	exclude := []secretEntry{}
	for _, f := range fields {
		if !sensitiveNameRe.MatchString(f.path) {
			return cfgTree{}, nil, fmt.Errorf("internal: sensitive field %q not matched by conservative rule", f.path)
		}
		src := f.srcFile
		if src == nil {
			src = f.srcDefault
		}
		if src != nil {
			exclude = append(exclude, secretEntry{Key: f.path, Env: envNameFrom(*src)})
		}
	}
	tree.Auth.JWTSecret = ""
	tree.Admin.InitialPassword = ""
	return tree, exclude, nil
}

func strOr(v *string, fallback *string) string {
	if v != nil {
		return *v
	}
	if fallback != nil {
		return *fallback
	}
	return ""
}

func listOr(v, fallback []string) []string {
	if len(v) > 0 {
		return append([]string(nil), v...)
	}
	return append([]string(nil), fallback...)
}

// buildPackage 构建配置包 v1（装载校验 + 元数据 + 导出树 + secrets.exclude）。
// cmdConfigExport 与测试共用；path 为空 = 内嵌默认 + env 覆盖。
func buildPackage(path string) (configPackage, error) {
	cfg, err := server.LoadConfig(path)
	if err != nil {
		return configPackage{}, err
	}
	tree, exclude, err := buildExportTree(path)
	if err != nil {
		return configPackage{}, err
	}
	pkg := configPackage{
		Package: pkgMeta{
			Format:     configPackageFormat,
			Version:    configPackageVer,
			App:        cfg.AppName,
			Env:        cfg.AppEnv,
			Profile:    cfg.ProfileName,
			ExportedAt: time.Now().UTC().Format(time.RFC3339),
		},
		Config: tree,
	}
	pkg.Secrets.Exclude = exclude
	return pkg, nil
}

// ---- config 分派 ----

func cmdConfig(args []string) error {
	if len(args) < 1 {
		return errCli(2, "config 需要子命令: export | diff | dry-run | import")
	}
	switch args[0] {
	case "export":
		return cmdConfigExport(args[1:])
	case "diff":
		return cmdConfigDiff(args[1:])
	case "dry-run":
		return errCli(2, "config dry-run: 尚未实现（VP-025 R3）")
	case "import":
		return errCli(2, "config import: 尚未实现（VP-025 R3）")
	default:
		return errCli(2, "未知 config 子命令 %q（export | diff | dry-run | import）", args[0])
	}
}

// ---- export ----

func cmdConfigExport(args []string) error {
	fs := flag.NewFlagSet("config export", flag.ContinueOnError)
	configPath := fs.String("config", "", "serve 配置路径（缺省 = 内嵌默认 + env 覆盖）")
	outPath := fs.String("o", "", "输出文件（缺省 stdout）")
	format := fs.String("f", "yaml", "输出格式: yaml|json")
	if err := fs.Parse(args); err != nil {
		return errCli(2, "%v", err)
	}
	if *format != "yaml" && *format != "json" {
		return errCli(2, "config export: 未知格式 %q（yaml|json）", *format)
	}

	pkg, err := buildPackage(*configPath)
	if err != nil {
		return errCli(1, "%v", err)
	}

	var out []byte
	if *format == "json" {
		out, err = json.MarshalIndent(pkg, "", "  ")
		if err != nil {
			return errCli(1, "config export: %v", err)
		}
		out = append(out, '\n')
	} else {
		out, err = yaml.Marshal(&pkg)
		if err != nil {
			return errCli(1, "config export: %v", err)
		}
	}

	if strings.TrimSpace(*outPath) != "" {
		if err := os.WriteFile(*outPath, out, 0o644); err != nil {
			return errCli(1, "config export: write %q: %v", *outPath, err)
		}
	} else {
		if _, err := os.Stdout.Write(out); err != nil {
			return errCli(1, "config export: %v", err)
		}
	}
	return nil
}

// ---- diff ----

type diffEntry struct {
	Path string `yaml:"path" json:"path"`
	Kind string `yaml:"kind" json:"kind"` // add | modify | remove
	Old  string `yaml:"old,omitempty" json:"old,omitempty"`
	New  string `yaml:"new,omitempty" json:"new,omitempty"`
}

// flattenLeaf 递归把配置树扁平化为「点分路径 → 显示串」叶子（列表元素 join）。
func flattenLeaf(prefix string, v any, out map[string]string) {
	switch t := v.(type) {
	case map[string]any:
		for k, vv := range t {
			flat := k
			if prefix != "" {
				flat = prefix + "." + k
			}
			flattenLeaf(flat, vv, out)
		}
	case []any:
		parts := make([]string, 0, len(t))
		for _, e := range t {
			parts = append(parts, fmt.Sprint(e))
		}
		out[prefix] = strings.Join(parts, ", ")
	default:
		out[prefix] = fmt.Sprint(v)
	}
}

// loadConfigSubtree 解析包文件（或 --against 源）为 config 段扁平叶子。
// pkg 文件 → configPackage YAML 的 config 子树；源配置 → buildExportTree。
func loadConfigLeaf(path string) (map[string]string, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("config diff: read %q: %w", path, err)
	}
	var root map[string]any
	if err := yaml.Unmarshal(raw, &root); err != nil {
		return nil, fmt.Errorf("config diff: parse %q: %w", path, err)
	}
	sub, ok := root["config"]
	if !ok {
		// 允许 --against 语义：非包文件按 export 管线处理（默认 ∪ 显式）。
		tree, _, err := buildExportTree(path)
		if err != nil {
			return nil, err
		}
		treeRaw, err := yaml.Marshal(&tree)
		if err != nil {
			return nil, err
		}
		sub = map[string]any{}
		if err := yaml.Unmarshal(treeRaw, &sub); err != nil {
			return nil, fmt.Errorf("config diff: normalize %q: %w", path, err)
		}
	}
	leaf := map[string]string{}
	flattenLeaf("", sub, leaf)
	return leaf, nil
}

func diffLeafMaps(a, b map[string]string) []diffEntry {
	paths := make(map[string]bool, len(a)+len(b))
	for p := range a {
		paths[p] = true
	}
	for p := range b {
		paths[p] = true
	}
	sorted := make([]string, 0, len(paths))
	for p := range paths {
		sorted = append(sorted, p)
	}
	sort.Strings(sorted)

	out := make([]diffEntry, 0)
	for _, p := range sorted {
		av, aok := a[p]
		bv, bok := b[p]
		switch {
		case aok && !bok:
			out = append(out, diffEntry{Path: p, Kind: "remove", Old: av})
		case !aok && bok:
			out = append(out, diffEntry{Path: p, Kind: "add", New: bv})
		case av != bv:
			out = append(out, diffEntry{Path: p, Kind: "modify", Old: av, New: bv})
		}
	}
	return out
}

func cmdConfigDiff(args []string) error {
	// 手写参数解析（同 cmdCreate 风格）：允许 <pkg> --against <src> 任意顺序。
	var against, format string
	var pos []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--against" || a == "-against":
			if i+1 >= len(args) {
				return errCli(2, "config diff: --against 需要值")
			}
			against = args[i+1]
			i++
		case strings.HasPrefix(a, "--against="):
			against = strings.TrimPrefix(a, "--against=")
		case a == "-f" || a == "--f":
			if i+1 >= len(args) {
				return errCli(2, "config diff: -f 需要值")
			}
			format = args[i+1]
			i++
		case strings.HasPrefix(a, "-f="):
			format = strings.TrimPrefix(a, "-f=")
		case strings.HasPrefix(a, "-") && a != "-":
			return errCli(2, "config diff: 未知参数 %q", a)
		default:
			pos = append(pos, a)
		}
	}
	if format == "" {
		format = "yaml"
	}
	if format != "yaml" && format != "json" {
		return errCli(2, "config diff: 未知格式 %q（yaml|json）", format)
	}

	var a, b map[string]string
	var err error
	if strings.TrimSpace(against) != "" {
		if len(pos) != 1 {
			return errCli(2, "config diff: --against 模式需要恰好一个包路径")
		}
		if a, err = loadConfigLeaf(pos[0]); err != nil {
			return errCli(2, "%v", err)
		}
		if b, err = loadConfigLeaf(against); err != nil {
			return errCli(2, "%v", err)
		}
	} else {
		if len(pos) != 2 {
			return errCli(2, "config diff 需要 <pkg-a> <pkg-b> 或 <pkg> --against <config>")
		}
		if a, err = loadConfigLeaf(pos[0]); err != nil {
			return errCli(2, "%v", err)
		}
		if b, err = loadConfigLeaf(pos[1]); err != nil {
			return errCli(2, "%v", err)
		}
	}

	entries := diffLeafMaps(a, b)
	if entries == nil {
		entries = []diffEntry{}
	}

	var out []byte
	if format == "json" {
		out, err = json.MarshalIndent(entries, "", "  ")
		if err != nil {
			return errCli(2, "config diff: %v", err)
		}
		out = append(out, '\n')
	} else {
		out, err = yaml.Marshal(entries)
		if err != nil {
			return errCli(2, "config diff: %v", err)
		}
	}
	if _, err := os.Stdout.Write(out); err != nil {
		return errCli(2, "config diff: %v", err)
	}
	if len(entries) > 0 {
		return errCli(1, "config diff: %d difference(s) found", len(entries))
	}
	return nil
}
