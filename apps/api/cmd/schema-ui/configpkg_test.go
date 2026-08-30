package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func writeTemp(t *testing.T, name, content string) string {
	t.Helper()
	p := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func exportToFile(t *testing.T, path string, extra ...string) string {
	t.Helper()
	t.Setenv("APP_ENV", "development")
	args := append([]string{"-o", path}, extra...)
	if err := cmdConfigExport(args); err != nil {
		t.Fatalf("export %v: %v", args, err)
	}
	return path
}

// ---- export ----

func TestExportDefaultShape(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	pkg, err := buildPackage("")
	if err != nil {
		t.Fatalf("buildPackage: %v", err)
	}
	if pkg.Package.Format != configPackageFormat {
		t.Errorf("format = %q, want %q", pkg.Package.Format, configPackageFormat)
	}
	if pkg.Package.Version != 1 {
		t.Errorf("version = %d, want 1", pkg.Package.Version)
	}
	if pkg.Package.Env != "development" {
		t.Errorf("package.env = %q, want development", pkg.Package.Env)
	}
	if pkg.Package.Profile != "admin" {
		t.Errorf("package.profile = %q, want admin", pkg.Package.Profile)
	}
	// env 引用保留形态（不解析）
	if pkg.Config.App.Env != "${APP_ENV:-development}" {
		t.Errorf("config.app.env = %q, want ${APP_ENV:-development} (reference preserved)", pkg.Config.App.Env)
	}
	if pkg.Config.HTTP.Addr != "127.0.0.1:25080" {
		t.Errorf("config.http.addr = %q, want loopback default", pkg.Config.HTTP.Addr)
	}
	// 敏感键剔除
	if pkg.Config.Auth.JWTSecret != "" || pkg.Config.Admin.InitialPassword != "" {
		t.Error("sensitive keys not stripped from config tree")
	}
	if len(pkg.Secrets.Exclude) != 2 {
		t.Fatalf("secrets.exclude = %d entries, want 2", len(pkg.Secrets.Exclude))
	}
	if pkg.Secrets.Exclude[0].Key != "auth.jwt_secret" || pkg.Secrets.Exclude[0].Env != "AUTH_JWT_SECRET" {
		t.Errorf("exclude[0] = %+v, want auth.jwt_secret/AUTH_JWT_SECRET", pkg.Secrets.Exclude[0])
	}
	if pkg.Secrets.Exclude[1].Key != "admin.initial_password" || pkg.Secrets.Exclude[1].Env != "ADMIN_INITIAL_PASSWORD" {
		t.Errorf("exclude[1] = %+v, want admin.initial_password/ADMIN_INITIAL_PASSWORD", pkg.Secrets.Exclude[1])
	}
	// 无明文扫描：config 段不得出现敏感键（secrets.exclude 键名字样为规定产物）；
	// 整包文本不得出现任何凭据值形态。
	raw, _ := yaml.Marshal(&pkg)
	var root map[string]any
	if err := yaml.Unmarshal(raw, &root); err != nil {
		t.Fatal(err)
	}
	cfg, ok := root["config"].(map[string]any)
	if !ok {
		t.Fatal("package has no config section")
	}
	auth, _ := cfg["auth"].(map[string]any)
	admin, _ := cfg["admin"].(map[string]any)
	if _, leaked := auth["jwt_secret"]; leaked {
		t.Error("config.auth.jwt_secret leaked into package")
	}
	if _, leaked := admin["initial_password"]; leaked {
		t.Error("config.admin.initial_password leaked into package")
	}
	for _, banned := range []string{"supersecret", "changeme"} {
		if strings.Contains(string(raw), banned) {
			t.Errorf("package contains credential-shaped text %q:\n%s", banned, raw)
		}
	}
}

func TestExportOverlay(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	src := writeTemp(t, "config.yaml", "app:\n  env: ${APP_ENV:-development}\nhttp:\n  addr: \"0.0.0.0:9999\"\n")
	pkg, err := buildPackage(src)
	if err != nil {
		t.Fatalf("buildPackage(%s): %v", src, err)
	}
	if pkg.Config.HTTP.Addr != "0.0.0.0:9999" {
		t.Errorf("http.addr = %q, want explicit overlay", pkg.Config.HTTP.Addr)
	}
	if pkg.Config.HTTP.ReadTimeout != "5s" {
		t.Errorf("http.read_timeout = %q, want embedded default 5s", pkg.Config.HTTP.ReadTimeout)
	}
}

func TestExportJSON(t *testing.T) {
	out := filepath.Join(t.TempDir(), "pkg.json")
	exportToFile(t, out, "-f", "json")
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), `"format": "schema-ui-config-package"`) {
		t.Errorf("json output missing format marker:\n%s", raw)
	}
	if strings.Contains(string(raw), `"jwt_secret"`) || strings.Contains(string(raw), `"initial_password"`) {
		t.Errorf("json output leaks a sensitive key name:\n%s", raw)
	}
}

func TestExportBadConfigFails(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope.yaml")
	err := cmdConfigExport([]string{"-config", missing})
	var ce *cliError
	if !errors.As(err, &ce) || ce.code != 1 {
		t.Fatalf("err = %v, want cliError code 1", err)
	}
}

// ---- diff ----

func TestDiffIdenticalAndIgnoredMeta(t *testing.T) {
	dir := t.TempDir()
	a := exportToFile(t, filepath.Join(dir, "a.yaml"))
	b := exportToFile(t, filepath.Join(dir, "b.yaml"))
	// exported_at 必不同（两次导出）；把 package.env 也改成不同 —— 都应被忽略。
	raw, err := os.ReadFile(b)
	if err != nil {
		t.Fatal(err)
	}
	raw = []byte(strings.Replace(string(raw), "env: development", "env: production", 1))
	if err := os.WriteFile(b, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := cmdConfigDiff([]string{a, b}); err != nil {
		t.Errorf("identical packages (ignoring metadata) differ: %v", err)
	}
	la, _ := loadConfigLeaf(a)
	lb, _ := loadConfigLeaf(b)
	if len(diffLeafMaps(la, lb)) != 0 {
		t.Errorf("leaf diff not empty: %v", diffLeafMaps(la, lb))
	}
}

func TestDiffModify(t *testing.T) {
	dir := t.TempDir()
	a := exportToFile(t, filepath.Join(dir, "a.yaml"))
	b := exportToFile(t, filepath.Join(dir, "b.yaml"))
	// TestDiffModify：导出 → 文本子串替换地址 → 命令层退出码 1 + diff 条目 modify http.addr。
	raw, err := os.ReadFile(b)
	if err != nil {
		t.Fatal(err)
	}
	raw = []byte(strings.Replace(string(raw), "127.0.0.1:25080", "0.0.0.0:25080", 1))
	if err := os.WriteFile(b, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	la, _ := loadConfigLeaf(a)
	lb, _ := loadConfigLeaf(b)
	entries := diffLeafMaps(la, lb)
	if len(entries) != 1 || entries[0].Path != "http.addr" || entries[0].Kind != "modify" {
		t.Fatalf("entries = %+v, want single modify http.addr", entries)
	}
	err = cmdConfigDiff([]string{a, b})
	var ce *cliError
	if !errors.As(err, &ce) || ce.code != 1 {
		t.Fatalf("err = %v, want cliError code 1 (differences found)", err)
	}
}

func TestDiffAddRemove(t *testing.T) {
	mk := func(addr, level string, withProxy bool) string {
		proxy := ""
		if withProxy {
			proxy = "    trusted_proxies: [\"10.0.0.0/8\"]\n"
		}
		return `package:
  format: schema-ui-config-package
  version: 1
  app: schema-ui-app
  env: development
  profile: admin
  exported_at: "2026-08-30T00:00:00Z"
config:
  app:
    env: ${APP_ENV:-development}
  http:
    addr: "` + addr + `"
` + proxy + `  profile: admin
  log:
    level: ` + level + `
secrets:
  exclude: []
`
	}
	a := writeTemp(t, "a.yaml", mk("127.0.0.1:25080", "info", false))
	b := writeTemp(t, "b.yaml", mk("0.0.0.0:9999", "debug", true))
	la, _ := loadConfigLeaf(a)
	lb, _ := loadConfigLeaf(b)
	entries := diffLeafMaps(la, lb)
	got := map[string]string{}
	for _, e := range entries {
		got[e.Kind+" "+e.Path] = e.Path
	}
	if _, ok := got["modify http.addr"]; !ok {
		t.Errorf("missing modify http.addr: %+v", entries)
	}
	if _, ok := got["modify log.level"]; !ok {
		t.Errorf("missing modify log.level: %+v", entries)
	}
	if _, ok := got["add http.trusted_proxies"]; !ok {
		t.Errorf("missing add http.trusted_proxies: %+v", entries)
	}
}

func TestDiffAgainst(t *testing.T) {
	dir := t.TempDir()
	a := exportToFile(t, filepath.Join(dir, "a.yaml"))
	src := writeTemp(t, "src.yaml", "app:\n  env: ${APP_ENV:-development}\nhttp:\n  addr: \"0.0.0.0:25080\"\n")
	la, _ := loadConfigLeaf(a)
	lb, err := loadConfigLeaf(src)
	if err != nil {
		t.Fatal(err)
	}
	entries := diffLeafMaps(la, lb)
	if len(entries) == 0 || entries[0].Path != "http.addr" {
		t.Fatalf("--against entries = %+v, want modify http.addr", entries)
	}
	err = cmdConfigDiff([]string{a, "--against", src})
	var ce *cliError
	if !errors.As(err, &ce) || ce.code != 1 {
		t.Fatalf("err = %v, want cliError code 1", err)
	}
}

func TestDiffErrors(t *testing.T) {
	dir := t.TempDir()
	missing := filepath.Join(dir, "missing.yaml")
	other := exportToFile(t, filepath.Join(dir, "other.yaml"))
	err := cmdConfigDiff([]string{missing, other})
	var ce *cliError
	if !errors.As(err, &ce) || ce.code != 2 {
		t.Fatalf("err = %v, want cliError code 2", err)
	}
	// 参数形态错误
	if err := cmdConfigDiff([]string{"only-one"}); err == nil {
		t.Error("expected arg-shape error")
	} else if ce2, ok := err.(*cliError); !ok || ce2.code != 2 {
		t.Errorf("arg error = %#v, want cliError code 2", err)
	}
	// dry-run / import 占位（R3）
	for _, sub := range []string{"dry-run", "import"} {
		err := cmdConfig([]string{sub, "x.yaml"})
		if ce2, ok := err.(*cliError); !ok || ce2.code != 2 {
			t.Errorf("%s placeholder = %#v, want cliError code 2", sub, err)
		}
	}
	if err := cmdConfig(nil); err == nil {
		t.Error("config without subcommand must fail")
	}
}

func TestRoundtrip(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	dir := t.TempDir()
	path := exportToFile(t, filepath.Join(dir, "pkg.yaml"))
	leaf, err := loadConfigLeaf(path)
	if err != nil {
		t.Fatal(err)
	}
	tree, _, err := buildExportTree("")
	if err != nil {
		t.Fatal(err)
	}
	treeRaw, _ := yaml.Marshal(&tree)
	var m map[string]any
	if err := yaml.Unmarshal(treeRaw, &m); err != nil {
		t.Fatal(err)
	}
	expected := map[string]string{}
	flattenLeaf("", m, expected)
	if len(leaf) != len(expected) {
		t.Fatalf("roundtrip leaf count = %d, want %d (%v vs %v)", len(leaf), len(expected), leaf, expected)
	}
	for k, v := range expected {
		if leaf[k] != v {
			t.Errorf("roundtrip mismatch at %s: %q vs %q", k, leaf[k], v)
		}
	}
}
