// Command schema-ui —— schema-ui-core 的 cli+包 分发 CLI（VP-023 R2）。
//
// 形态（用户定案）：单命令 · Go 单二进制 · create/add/upgrade；零新增依赖
// （标准库解析 + go:embed 模板）；分发 = go install
// github.com/magicvr/schema-ui-core/apps/api/cmd/schema-ui@vX.Y.Z（复用模块 tag 链）。
package main

import (
	"embed"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/magicvr/schema-ui-core/apps/api/server"
)

//go:embed templates
var templateFS embed.FS

const (
	apiModule   = "github.com/magicvr/schema-ui-core/apps/api"
	apiVersion  = "v0.4.0" // 模板钉 = 含公开 serve 面（server pkg）的下一发布；R2 发布通道核销
	// 六包终值（R5 · A-002 F-002 响应）：create 骨架生成即装终值；upgrade 拉 npm/Go latest 保对齐
	protocolVer = "0.2.11"
	libVersion    = "0.1.10"
	shellVersion  = "0.1.4"
	themeVersion  = "0.1.4"
	uiVersion     = "0.1.8"
	rendererVer = "0.3.8"
)

type createOpts struct {
	name   string
	module string
	dir    string
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `schema-ui — schema-ui-core 分发 CLI（cli+包 形态）

用法:
  schema-ui create <name> [--module <path>] [--dir <dir>]   生成下游骨架（Go 组合根 + Web 骨架 + 探针）
  schema-ui serve [-config <path>] [-dialect sqlite|postgres] [-dsn <path|conn>] [-addr <addr>]   运行下游 serve 面（RT-D02 优雅停机）
  schema-ui add [module-id]                                  列出可用模块 / registry 装配
  schema-ui upgrade [--dry-run]                               registry 升级 + 探针回归
  schema-ui migrate-fork [--dir <path>] [--dry-run]           fork → 包 迁移辅助（非破坏性 · A/B 型）
  schema-ui config export [-config <path>] [-o <path>] [-f yaml|json]   导出配置包 v1（VP-025 R2）
  schema-ui config diff <pkg-a> <pkg-b> | <pkg> --against <config> [-f yaml|json]   键级差量（0 无差 / 1 有差 / 2 错误）
  schema-ui config dry-run <pkg> [-config <path>] [-f yaml|json]      导入前只读预检（VP-025 R3）
  schema-ui config import <pkg> [-file <path>]               导入配置包（仅注册 · R3）

示例:
  schema-ui create my-admin
  schema-ui serve -config config.yaml
  schema-ui serve -dialect postgres -dsn "postgres://user:pass@127.0.0.1:5432/db"
  schema-ui add                     # 列出 kernel.BuiltinModules
  schema-ui upgrade --dry-run
  schema-ui config export -o config.package.yaml
  schema-ui config diff a.yaml b.yaml
`)
	}
	flag.Parse()

	if flag.NArg() < 1 {
		flag.Usage()
		os.Exit(2)
	}
	var err error
	switch flag.Arg(0) {
	case "create":
		err = cmdCreate(flag.Args()[1:])
	case "serve":
		err = cmdServe(flag.Args()[1:])
	case "add":
		err = cmdAdd(flag.Args()[1:])
	case "upgrade":
		err = cmdUpgrade(flag.Args()[1:])
	case "migrate-fork":
		err = cmdMigrateFork(flag.Args()[1:])
	case "config":
		err = cmdConfig(flag.Args()[1:])
	default:
		flag.Usage()
		os.Exit(2)
	}
	if err != nil {
		var ce *cliError
		if errors.As(err, &ce) {
			fmt.Fprintln(os.Stderr, "schema-ui:", ce.err)
			os.Exit(ce.code)
		}
		fmt.Fprintln(os.Stderr, "schema-ui:", err)
		os.Exit(1)
	}
}

// ---- create ----

func cmdCreate(args []string) error {
	var module, dir string
	var positional []string
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--module" || a == "-module" || a == "--module=" || a == "-module=":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				module = args[i+1]
				i++
			}
		case strings.HasPrefix(a, "--module=") || strings.HasPrefix(a, "-module="):
			module = strings.TrimPrefix(strings.TrimPrefix(a, "--module="), "-module=")
		case a == "--dir" || a == "-dir" || a == "--dir=" || a == "-dir=":
			if i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
				dir = args[i+1]
				i++
			}
		case strings.HasPrefix(a, "--dir=") || strings.HasPrefix(a, "-dir="):
			dir = strings.TrimPrefix(strings.TrimPrefix(a, "--dir="), "-dir=")
		default:
			positional = append(positional, a)
		}
	}
	if len(positional) < 1 {
		return fmt.Errorf("create 需要 <name>")
	}
	opts := createOpts{name: positional[0]}
	if module != "" {
		opts.module = module
	} else {
		opts.module = "github.com/magicvr/" + opts.name
	}
	if dir != "" {
		opts.dir = dir
	} else {
		opts.dir = opts.name
	}
	return generate(opts)
}

func generate(o createOpts) error {
	if _, err := os.Stat(o.dir); err == nil {
		return fmt.Errorf("目标目录已存在: %s", o.dir)
	}
	data := map[string]string{
		"Name":           o.name,
		"LibVersion":      libVersion,
		"ShellVersion":    shellVersion,
		"ThemeVersion":    themeVersion,
		"UiVersion":       uiVersion,
		"Module":         o.module,
		"APIVersion":     apiVersion,
		"ProtocolVersion": protocolVer,
		"RendererVersion": rendererVer,
	}
	files := map[string]string{
		"templates/go.mod.tmpl":              filepath.Join(o.dir, "go.mod"),
		"templates/main.go.tmpl":             filepath.Join(o.dir, "cmd", "server", "main.go"),
		"templates/config.yaml.tmpl":         filepath.Join(o.dir, "config.example.yaml"),
		"templates/README.md.tmpl":           filepath.Join(o.dir, "README.md"),
		"templates/gitignore.tmpl":           filepath.Join(o.dir, ".gitignore"),
		"templates/web/package.json.tmpl":    filepath.Join(o.dir, "web", "package.json"),
		"templates/web/npmrc.tmpl":           filepath.Join(o.dir, "web", ".npmrc"),
		"templates/probe.mjs":            filepath.Join(o.dir, "web", "probe.mjs"),
		"templates/probe-render.mjs":     filepath.Join(o.dir, "web", "probe-render.mjs"),
		"templates/token-check.mjs":      filepath.Join(o.dir, "web", "token-check.mjs"),
		"templates/probe-six.mjs":            filepath.Join(o.dir, "web", "probe-six.mjs"),
		"templates/index.css":            filepath.Join(o.dir, "web", "index.css"),
		"templates/brand.css":            filepath.Join(o.dir, "web", "brand.css"),
	}
	for src, dst := range files {
		if err := writeTemplate(src, dst, data); err != nil {
			return err
		}
	}
	fmt.Printf("created %s\n  module: %s\n  api:    %s@%s\n\n下一步:\n  cd %s && go mod tidy && go run ./cmd/server\n  cd %s/web && pnpm install && node probe.mjs\n",
		o.dir, o.module, apiModule, apiVersion, o.dir, o.dir)
	return nil
}

// ---- serve ----

// cmdServe 运行下游 serve 面（public server 包 + RT-D02 优雅停机）。
func cmdServe(args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	configPath := fs.String("config", "", "serve 配置路径（缺省 = 内嵌默认 + env）")
	dialect := fs.String("dialect", "", "覆盖方言（sqlite|postgres）")
	dsn := fs.String("dsn", "", "覆盖连接串（sqlite=文件路径；postgres=DSN）")
	addr := fs.String("addr", "", "覆盖监听地址")
	if err := fs.Parse(args); err != nil {
		return err
	}
	cfg, err := server.LoadConfig(*configPath)
	if err != nil {
		return err
	}
	if *dialect != "" {
		cfg.DBDialect = *dialect
	}
	if *dsn != "" {
		if cfg.DBDialect == "postgres" {
			cfg.DBDSN = *dsn
		} else {
			cfg.DBPath = *dsn
		}
	}
	if *addr != "" {
		cfg.HTTPAddr = *addr
	}
	return server.Serve(server.Options{Config: cfg})
}

func writeTemplate(src, dst string, data map[string]string) error {
	raw, err := templateFS.ReadFile(src)
	if err != nil {
		return err
	}
	if strings.HasSuffix(src, ".tmpl") {
		tmpl, err := template.New(src).Parse(string(raw))
		if err != nil {
			return err
		}
		var buf strings.Builder
		if err := tmpl.Execute(&buf, data); err != nil {
			return err
		}
		raw = []byte(buf.String())
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, raw, 0o644)
}

// ---- add ----

func cmdAdd(args []string) error {
	fmt.Printf("可用标准模块（kernel.BuiltinModules）:\n")
	mods, err := listModules()
	if err != nil {
		return err
	}
	for _, m := range mods {
		fmt.Printf("  %-24s v%-8s kernel-api %s\n", m.id, m.version, m.kernelAPI)
	}
	if len(args) == 0 {
		return nil
	}
	// add <module> = 对当前仓执行 registry 装配（go get / pnpm add）
	if err := registryAdd(args[0]); err != nil {
		return err
	}
	return nil
}

type moduleInfo struct{ id, version, kernelAPI string }

// listModules 经同模块 import 列出可用模块（CLI 在 apps/api 模块内，直接读 kernel）。
func listModules() ([]moduleInfo, error) {
	mods := builtinModules()
	return mods, nil
}

// registryAdd 对目标仓（探测 go.mod / web）执行 registry 语义装配。
func registryAdd(pkg string) error {
	shell := "sh"
	if runtime.GOOS == "windows" {
		shell = "cmd"
	}
	// 探测：仓根有 go.mod → go get；web/ 有 package.json → pnpm add
	if _, err := os.Stat("web/package.json"); err == nil {
		return runCmd(shell, buildShellCmd("pnpm add "+pkg))
	}
	if _, err := os.Stat("go.mod"); err == nil {
		return runCmd(shell, buildShellCmd("go get "+pkg+"@latest"))
	}
	return fmt.Errorf("未识别仓类型（缺 go.mod / web/package.json）")
}

// ---- upgrade ----

func cmdUpgrade(args []string) error {
	dry := false
	for _, a := range args {
		if a == "--dry-run" {
			dry = true
		}
	}
	shell := "sh"
	if runtime.GOOS == "windows" {
		shell = "cmd"
	}
	fmt.Println("schema-ui upgrade：registry 语义升级 + 探针回归")
	if _, err := os.Stat("go.mod"); err == nil {
		cmd := "go get " + apiModule + "@latest"
		if dry {
			cmd = "go list -m -u " + apiModule
		}
		fmt.Println("> " + cmd)
		if !dry {
			if err := runCmd(shell, buildShellCmd(cmd)); err != nil {
				return err
			}
		}
	}
	if _, err := os.Stat("web/package.json"); err == nil {
		cmd := "pnpm add @magicvr/schema-ui-lib@latest @magicvr/schema-ui-protocol@latest @magicvr/schema-ui-renderer@latest @magicvr/schema-ui-shell@latest @magicvr/schema-ui-theme@latest @magicvr/schema-ui-ui@latest"
		if dry {
			cmd = "pnpm outdated @magicvr/schema-ui-lib @magicvr/schema-ui-protocol @magicvr/schema-ui-renderer @magicvr/schema-ui-shell @magicvr/schema-ui-theme @magicvr/schema-ui-ui"
		}
		fmt.Println("> " + cmd)
		if !dry {
			if err := runCmd(shell, buildShellCmd(cmd)); err != nil {
				return err
			}
		}
	}
	if !dry {
		// 探针回归
		if _, err := os.Stat("web/probe.mjs"); err == nil {
			fmt.Println("> node web/probe.mjs")
			if err := runCmd(shell, buildShellCmd("node web/probe.mjs")); err != nil {
				return fmt.Errorf("探针失败: %w", err)
			}
		}
		if _, err := os.Stat("cmd/server"); err == nil {
			// A-002 F-006：serve 为常驻进程，upgrade 不再同步执行（防挂起）；改为提示手动冒烟
			fmt.Println("（组合根冒烟：另起终端 `go run ./cmd/server -dialect sqlite -dsn ./data.db`，随后验 /healthz）")
		}
	}
	fmt.Println("upgrade" + map[bool]string{true: "（dry-run）", false: ""}[dry] + " 完成")
	return nil
}

func buildShellCmd(cmd string) string {
	if runtime.GOOS == "windows" {
		return "/c " + cmd
	}
	return "-c " + cmd
}

func runCmd(shell, args string) error {
	c := exec.Command(shell, args)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}