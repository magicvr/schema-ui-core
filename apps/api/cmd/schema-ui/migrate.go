// migrate.go —— schema-ui migrate-fork：fork → 包迁移辅助（R6 · VP-024 判据 #7）。
//
// 非破坏性原则（D-001-r6）：除两处低侵入改写（go.mod require bump ·
// web/.npmrc 钉 npmjs，后者带备份）外零修改；用户代码一律检查 + 引导。
// 类型判定（指南 §1）：A 纯装配（薄/组合根无 kernel 覆盖）· B 轻度定制 ·
// C 深度定制（kernel 契约面 import）→ 建议保持 fork。
package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

const npmjsRegistryLine = "@magicvr:registry=https://registry.npmjs.org"

var requireRe = regexp.MustCompile(`(?m)^\s*require\s+` + apiModule + `\s+v?([0-9]+\.[0-9]+\.[0-9]+)`)
var requireBlockRe = regexp.MustCompile(`(?m)^\s*` + apiModule + `\s+v?([0-9]+\.[0-9]+\.[0-9]+)$`)
var kernelCoverRe = regexp.MustCompile(`(?m)(kernel\.|assembly\.OpenStore|"github.com/magicvr/schema-ui-core/apps/api/(kernel|assembly|modules)/)`)

type migrationReport struct {
	dir       string
	kind      string
	kindNote  string
	goVersion string // go.mod require 中的现行版本；空 = 未依赖
	npmrcLine string // .npmrc 中 @magicvr 映射行；空 = 无
	mainThin  bool   // cmd/server/main.go 已是薄封装
	mainCompose bool // main 为标准组合根（assembly.OpenStore 路径）
	steps     []string
}

func cmdMigrateFork(args []string) error {
	fs := flag.NewFlagSet("migrate-fork", flag.ContinueOnError)
	dir := fs.String("dir", ".", "目标 fork 仓根")
	dry := fs.Bool("dry-run", false, "只检查并输出步骤清单（不写任何文件）")
	if err := fs.Parse(args); err != nil {
		return err
	}
	return runMigrateFork(*dir, *dry)
}

func runMigrateFork(dir string, dry bool) error {
	abs, err := filepath.Abs(dir)
	if err != nil {
		return err
	}
	rep, err := inspectFork(abs)
	if err != nil {
		return err
	}
	fmt.Printf("schema-ui migrate-fork — 目标: %s\n", abs)
	fmt.Printf("类型判定: %s · %s\n", rep.kind, rep.kindNote)
	if rep.kind == "C" {
		fmt.Println("结论: 深度定制 fork（kernel 契约面被覆盖）→ 包形态暂不覆盖（Charter fork 并存）· 保持 fork 并登记包化承载面候选")
		fmt.Println("  （如需评估 kernel 扩展通道：assembly 扩展 + 六包 external 化的组合面 = 后续候选）")
		return nil
	}
	fmt.Println("迁移步骤清单（指南 §2 · 非破坏）:")
	for i, s := range rep.steps {
		fmt.Printf("  %d. %s\n", i+1, s)
	}
	if dry {
		fmt.Println("（--dry-run：未写任何文件）")
		return nil
	}
	// 实跑：1) go.mod require bump（registry 语义；无论是否已依赖都执行 go get @latest——F-002）
	fmt.Println("> go get " + apiModule + "@latest")
	if err := runIn(abs, "go", "get", apiModule+"@latest"); err != nil {
		return fmt.Errorf("go.mod bump 失败: %w", err)
	}
	if rep.goVersion != "" {
		fmt.Printf("  ok: require %s 已升至 registry @latest（原 %s）\n", apiModule, rep.goVersion)
	} else {
		fmt.Printf("  ok: require %s 已添加（registry @latest）\n", apiModule)
	}
	// 实跑：2) web/.npmrc 钉 npmjs（先备份）
	if rep.npmrcLine != "" && rep.npmrcLine != npmjsRegistryLine {
		p := filepath.Join(abs, "web", ".npmrc")
		bak := p + ".migrate.bak"
		if err := os.WriteFile(bak, []byte(rep.npmrcLine+"\n"), 0o644); err != nil {
			return err
		}
		rest := ""
		raw, err := os.ReadFile(p)
		if err == nil {
			rest = strings.ReplaceAll(strings.ReplaceAll(string(raw), rep.npmrcLine, ""), "\n\n", "\n")
		}
		if err := os.WriteFile(p, []byte(npmjsRegistryLine+"\n"+strings.TrimPrefix(rest, "\n")), 0o644); err != nil {
			return err
		}
		fmt.Printf("  ok: web/.npmrc 钉 npmjs（备份 %s）\n", filepath.Base(bak))
	}
	fmt.Println("引导（用户代码 · 未修改）:")
	if rep.mainThin {
		fmt.Println("  - cmd/server/main.go 已是薄封装（server.Serve）✓")
	} else {
		fmt.Println("  - cmd/server/main.go 为旧组合根 → 用 `schema-ui create <name>` 重建组合根，再迁入业务模块与 schema（指南 §2-2）")
	}
	fmt.Println("验证建议:")
	fmt.Println("  - go build ./cmd/server && go vet ./...")
	fmt.Println("  - go run ./cmd/server -dialect sqlite -dsn ./data.db（迁移台账核对）")
	if rep.npmrcLine != "" {
		fmt.Println("  - cd web && pnpm install && node probe.mjs && node probe-render.mjs && node probe-six.mjs && node token-check.mjs")
	}
	fmt.Println("migrate-fork 完成（A/B 型路径）")
	return nil
}

// inspectFork 读取目标仓并判定类型与步骤（不写文件）。
func inspectFork(dir string) (*migrationReport, error) {
	rep := &migrationReport{dir: dir}
	// 1) go.mod（单行 require 或 require ( … ) 块形态——F-002 补块解析）
	goMod := filepath.Join(dir, "go.mod")
	if raw, err := os.ReadFile(goMod); err == nil {
		if m := requireRe.FindSubmatch(raw); m != nil {
			rep.goVersion = strings.TrimSpace(string(m[1]))
		} else if m := requireBlockRe.FindSubmatch(raw); m != nil {
			rep.goVersion = strings.TrimSpace(string(m[1]))
		}
	}
	// 2) web/.npmrc
	npmrc := filepath.Join(dir, "web", ".npmrc")
	if raw, err := os.ReadFile(npmrc); err == nil {
		for _, line := range strings.Split(string(raw), "\n") {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "@magicvr:") || strings.Contains(line, "registry=https://npm.pkg.github.com") {
				rep.npmrcLine = line
				break
			}
		}
	}
	// 3) cmd/server/main.go 形态
	mainPath := filepath.Join(dir, "cmd", "server", "main.go")
	kernelCover := false
	composeRoot := false
	sawMain := false
	if raw, err := os.ReadFile(mainPath); err == nil {
		sawMain = true
		text := string(raw)
		composeRoot = strings.Contains(text, "assembly.OpenStore")
		// C 型 = 手搓 kernel 面且未走标准组合路径（模板组合根的标准件不算覆盖）
		kernelCover = kernelCoverRe.Match(raw) && !composeRoot
		rep.mainThin = strings.Contains(text, "server.Serve(")
	}
	rep.mainCompose = composeRoot
	// 类型判定（指南 §1 · 校准：标准组合根 = 可迁移）
	switch {
	case !sawMain && rep.goVersion == "":
		rep.kind = "?"
		rep.kindNote = "未识别 fork 形态（缺 go.mod 依赖 / cmd/server）"
	case kernelCover:
		rep.kind = "C"
		rep.kindNote = "cmd/server 手搓 kernel 面且未走标准组合路径（深度定制）"
	case rep.goVersion == "" && rep.npmrcLine == "":
		rep.kind = "A"
		rep.kindNote = "纯装配型（尚未依赖包面）"
	default:
		rep.kind = "A"
		if !rep.mainThin {
			rep.kind = "B"
			rep.kindNote = "轻度定制（旧组合根形态，无 kernel 覆盖）"
		} else {
			rep.kindNote = "纯装配/薄封装组合根"
		}
	}
	// 步骤清单
	if rep.goVersion == "" {
		rep.steps = append(rep.steps, "go.mod 尚无 "+apiModule+" 依赖 → 添加 require（go get "+apiModule+"@latest；registry 语义）")
	} else {
		rep.steps = append(rep.steps, "go.mod require "+apiModule+" 当前 "+rep.goVersion+" → bump 至 registry @latest")
	}
	if rep.npmrcLine == "" {
		rep.steps = append(rep.steps, "web/.npmrc 无 scope 映射 →（可选）写入 "+npmjsRegistryLine)
	} else if rep.npmrcLine == npmjsRegistryLine {
		rep.steps = append(rep.steps, "web/.npmrc 已钉 npmjs ✓（无动作）")
	} else {
		rep.steps = append(rep.steps, "web/.npmrc 含 "+rep.npmrcLine+" → 替换为 "+npmjsRegistryLine+"（备份 .npmrc.migrate.bak）")
	}
	if rep.mainThin {
		rep.steps = append(rep.steps, "cmd/server/main.go 薄封装 ✓（无动作）")
	} else if rep.mainCompose {
		rep.steps = append(rep.steps, "cmd/server/main.go 为旧组合根 → 换薄封装（schema-ui create 引导 · 业务模块迁入 · 不覆盖）")
	} else {
		rep.steps = append(rep.steps, "cmd/server/main.go 需重建（schema-ui create 引导 · 不覆盖）")
	}
	rep.steps = append(rep.steps, "验证：go build ./cmd/server + serve 冒烟 + web 四探针")
	return rep, nil
}

func runIn(dir, name string, arg ...string) error {
	c := exec.Command(name, arg...)
	c.Dir = dir
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

var _ = runtime.GOOS // 保留：若未来用 shell 语义