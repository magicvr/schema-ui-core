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

// ---- dry-run ----

func TestDryRunPass(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("AUTH_JWT_SECRET", "test-key")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "test-pw")
	dir := t.TempDir()
	pkgPath := exportToFile(t, filepath.Join(dir, "pkg.yaml"))
	report, err := dryRun(pkgPath, "")
	if err != nil {
		t.Fatalf("dryRun: %v", err)
	}
	for _, c := range report.Checks {
		if c.Status != "ok" {
			t.Errorf("check %s = %s (%s)", c.Path, c.Status, c.Message)
		}
	}
	// 与默认目标一致 → 无变更
	if len(report.Changes) != 0 {
		t.Errorf("changes = %+v, want empty", report.Changes)
	}
	// 零副作用：dry-run 不得修改目标文件（快照对比）
	src := writeTemp(t, "target.yaml", "app:\n  env: ${APP_ENV:-development}\nhttp:\n  addr: \"0.0.0.0:27080\"\n")
	before, _ := os.ReadFile(src)
	if _, err := dryRun(pkgPath, src); err != nil {
		t.Fatalf("dryRun against src: %v", err)
	}
	after, _ := os.ReadFile(src)
	if string(before) != string(after) {
		t.Error("dry-run modified the target file (side effect)")
	}
}

func TestDryRunEnvMissingFailsClosed(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("AUTH_JWT_SECRET", "test-key")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "test-pw")
	dir := t.TempDir()
	pkgPath := exportToFile(t, filepath.Join(dir, "pkg.yaml"))
	os.Unsetenv("AUTH_JWT_SECRET") // 彻底缺失 → fail-closed
	report, err := dryRun(pkgPath, "")
	if err == nil {
		t.Fatal("expected fail-closed error when AUTH_JWT_SECRET missing")
	}
	failed := false
	for _, c := range report.Checks {
		if c.Status == "fail" && strings.Contains(c.Path, "auth.jwt_secret") {
			failed = true
		}
	}
	if !failed {
		t.Errorf("no fail check for auth.jwt_secret: %+v", report.Checks)
	}
	err = cmdConfigDryRun([]string{pkgPath})
	var ce *cliError
	if !errors.As(err, &ce) || ce.code != 1 {
		t.Errorf("cmdConfigDryRun err = %v, want cliError code 1", err)
	}
}

func TestDryRunChanges(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("AUTH_JWT_SECRET", "test-key")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "test-pw")
	dir := t.TempDir()
	pkgPath := exportToFile(t, filepath.Join(dir, "pkg.yaml"))
	src := writeTemp(t, "target.yaml", "app:\n  env: ${APP_ENV:-development}\nhttp:\n  addr: \"0.0.0.0:27080\"\n")
	report, err := dryRun(pkgPath, src)
	if err != nil {
		t.Fatalf("dryRun: %v", err)
	}
	if len(report.Changes) != 1 || report.Changes[0].Path != "http.addr" || report.Changes[0].Kind != "modify" {
		t.Fatalf("changes = %+v, want single modify http.addr", report.Changes)
	}
	if report.Changes[0].Old != "0.0.0.0:27080" || report.Changes[0].New != "127.0.0.1:25080" {
		t.Errorf("old/new direction wrong: %+v", report.Changes[0])
	}
}

func TestDryRunInvalidPackage(t *testing.T) {
	t.Setenv("APP_ENV", "development")
	t.Setenv("AUTH_JWT_SECRET", "test-key")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "test-pw")
	dir := t.TempDir()
	bad := writeTemp(t, "bad.yaml", "package:\n  format: other\nconfig:\n  bogus: 1\n")
	err := cmdConfigDryRun([]string{bad})
	var ce *cliError
	if !errors.As(err, &ce) || ce.code != 1 {
		t.Fatalf("err = %v, want cliError code 1 (precheck failure)", err)
	}
	missing := filepath.Join(dir, "missing.yaml")
	err = cmdConfigDryRun([]string{missing})
	if !errors.As(err, &ce) || ce.code != 2 {
		t.Fatalf("err = %v, want cliError code 2 (tool error)", err)
	}
}

// ---- import ----

func importEnv(t *testing.T) {
	t.Helper()
	t.Setenv("APP_ENV", "development")
	t.Setenv("AUTH_JWT_SECRET", "test-key")
	t.Setenv("ADMIN_INITIAL_PASSWORD", "test-pw")
}

func TestImportRoundtrip(t *testing.T) {
	importEnv(t)
	dir := t.TempDir()
	pkgPath := exportToFile(t, filepath.Join(dir, "pkg.yaml"))
	target := filepath.Join(dir, "config.yaml")
	if err := cmdConfigImport([]string{pkgPath, "-file", target}); err != nil {
		t.Fatalf("import: %v", err)
	}
	raw, err := os.ReadFile(target)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(raw), "generated by schema-ui config import") {
		t.Errorf("target missing generation header:\n%s", raw)
	}
	// 往返：再 export（-config target）→ 与源包 config 树 diff 无差。
	rePkgPath := filepath.Join(dir, "re.yaml")
	exportToFile(t, rePkgPath, "-config", target)
	la, _ := loadConfigLeaf(pkgPath)
	lb, _ := loadConfigLeaf(rePkgPath)
	if entries := diffLeafMaps(la, lb); len(entries) != 0 {
		t.Errorf("roundtrip diff = %+v, want empty", entries)
	}
}

func TestImportBackup(t *testing.T) {
	importEnv(t)
	dir := t.TempDir()
	pkgPath := exportToFile(t, filepath.Join(dir, "pkg.yaml"))
	target := writeTemp(t, "config.yaml", "app:\n  env: ${APP_ENV:-development}\nhttp:\n  addr: \"0.0.0.0:9999\"\n")
	old, _ := os.ReadFile(target)
	if err := cmdConfigImport([]string{pkgPath, "-file", target}); err != nil {
		t.Fatalf("import: %v", err)
	}
	backup := target + ".pre-import.bak"
	bakRaw, err := os.ReadFile(backup)
	if err != nil {
		t.Fatalf("backup missing: %v", err)
	}
	if string(bakRaw) != string(old) {
		t.Error("backup content differs from pre-import target")
	}
	cur, _ := os.ReadFile(target)
	if !strings.Contains(string(cur), "generated by schema-ui config import") {
		t.Errorf("target not replaced: %s", cur)
	}
}

func TestImportRejectsAndKeepsUntouched(t *testing.T) {
	importEnv(t)
	dir := t.TempDir()
	pkgPath := exportToFile(t, filepath.Join(dir, "pkg.yaml"))

	// 预检失败（env 缺失）→ 拒绝 + 目标零触碰。
	target := filepath.Join(dir, "never.yaml")
	os.Unsetenv("AUTH_JWT_SECRET")
	err := cmdConfigImport([]string{pkgPath, "-file", target})
	var ce *cliError
	if !errors.As(err, &ce) || ce.code != 1 {
		t.Fatalf("err = %v, want cliError code 1 (precheck)", err)
	}
	if _, statErr := os.Stat(target); !os.IsNotExist(statErr) {
		t.Errorf("target created despite precheck failure")
	}
	importEnv(t)

	// 坏包 → 预检失败 + 目标未创建。
	bad := writeTemp(t, "bad.yaml", "package:\n  format: other\nconfig:\n  bogus: 1\n")
	target2 := filepath.Join(dir, "never2.yaml")
	if err := cmdConfigImport([]string{bad, "-file", target2}); err == nil {
		t.Fatal("expected import rejection for invalid package")
	}
	if _, statErr := os.Stat(target2); !os.IsNotExist(statErr) {
		t.Error("target created despite invalid package")
	}

	// 装载校验失败（db.dialect 非法 → LoadConfig fail-closed）→ exit 1 + 目标原样 + 无 tmp 残留。
	badDialect := writeTemp(t, "baddialect.yaml", `package:
  format: schema-ui-config-package
  version: 1
  app: schema-ui-app
  env: development
  profile: admin
  exported_at: "2026-08-30T00:00:00Z"
config:
  app:
    env: ${APP_ENV:-development}
  db:
    dialect: oracle
  http:
    addr: "127.0.0.1:25080"
  profile: admin
secrets:
  exclude: []
`)
	exists := writeTemp(t, "keep.yaml", "app:\n  env: ${APP_ENV:-development}\n")
	keepBefore, _ := os.ReadFile(exists)
	err = cmdConfigImport([]string{badDialect, "-file", exists})
	if !errors.As(err, &ce) || ce.code != 1 {
		t.Fatalf("err = %v, want cliError code 1 (applicability)", err)
	}
	keepAfter, _ := os.ReadFile(exists)
	if string(keepBefore) != string(keepAfter) {
		t.Error("target modified by failed import")
	}
	if _, statErr := os.Stat(exists + ".tmp"); !os.IsNotExist(statErr) {
		t.Error("tmp file leaked after failed import")
	}
	// 但也别留下已生成备份：装载失败发生在备份之后？实现顺序 = 备份 → tmp 写 → 校验失败。
	// 备份保留为核对物（方案 A 语义）；此处断言目标未动即可。
}

func TestImportDefaultFile(t *testing.T) {
	importEnv(t)
	dir := t.TempDir()
	pkgPath := exportToFile(t, filepath.Join(dir, "pkg.yaml"))
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	defer func() { _ = os.Chdir(old) }()
	if err := cmdConfigImport([]string{pkgPath}); err != nil {
		t.Fatalf("import default file: %v", err)
	}
	if _, err := os.Stat("config.yaml"); err != nil {
		t.Errorf("default config.yaml not created: %v", err)
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
