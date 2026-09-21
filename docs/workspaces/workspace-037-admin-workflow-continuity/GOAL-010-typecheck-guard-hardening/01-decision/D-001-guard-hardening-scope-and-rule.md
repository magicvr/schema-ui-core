---
id: D-001-guard-hardening-scope-and-rule
doc: decision
status: accepted
goal_id: GOAL-010-typecheck-guard-hardening
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# D-001 · 守卫加固范围与 `-p` 目标判定规则

## 问题

`GOAL-008` 的防复发守卫以「命令行中出现 `-b`/`--build`/`-p`/`--project` 令牌」判定检查型调用（`CHECKING_FLAG`）。该判定对 `-p` 是不充分的：`tsc --noEmit -p tsconfig.json` 的 `-p` 指向 solution-style 根配置（`files: []`、仅 `references`），非 build 模式下只编译空程序、恒 exit 0 —— 实测 exit 0（`GOAL-008 E-006`），却会被判合规。守卫是唯一的防复发装置，故该缺口属同类失效模式。

## 决定

**判定规则（收紧后）**：一个 `tsc` 调用构成「检查型」当且仅当满足其一——

1. 带 `-b` / `--build`（build 模式会走项目引用图）；或
2. 带 `-p` / `--project`，且其目标解析到**自身选择源文件**的配置（`selectsOwnSources`：`files` 非空，或 `include` 非空）。

否则为**违规**，包括：裸 `tsc` / `tsc --noEmit`、`-p` 指向 solution-style 根配置、`-p .`（解析为根配置）、`-p` 目标不存在或不可解析、`-p` 缺参数。

**fail closed**：目标无法读取/解析时计为违规而非放行——守卫的职责是防失实证据，宁可要求人工确认，不可静默放行空转命令。

**实现落点**：只改 `apps/web/src/typecheck-convention.guard.test.ts`（纯测试文件），不改任何产品代码、`package.json` 脚本或 CI 配置——正确口径未变，变的是守卫对错误形态的识别能力。

**合成变异用例**：以表驱动方式固定正/反例（含 `-p tsconfig.json`、`-p .`、`--project=`、不存在配置、`npm exec --` 前缀、以及「散文提及 `tsc` 不算调用」），使缺口回归时由测试直接指认。

## 未选方案

- **把 `-p` 一律视为有效（保持现状）**：即 `GOAL-008 A-002 F-001` 本身，不采纳。
- **只加 `-p tsconfig.json` 一条特例黑名单**：能挡住已知形态，但 `-p .`、`--project=tsconfig.json`、未来新增的 solution-style 子配置都会漏；不采纳，改为按配置内容判定。
- **要求 `-p` 目标必须是硬编码白名单（如仅 `tsconfig.app.json`/`e2e/tsconfig.json`）**：会阻止未来合法的项目配置（如新增 `tsconfig.test.json`），且白名单本身需要维护；不采纳。
- **顺带裁定全仓 300+ 行 `tsc` 简写**：范围与收益不成比例，`GOAL-008 A-002 F-002` 已登记为 recommended open；不并入本目标。
- **顺带修复其他已发现 recommended（`GOAL-008 A-001 F-002`、`GOAL-009 A-001 F-001/F-002`）**：本目标只承接一条 finding，保持范围小而可审计。

## 边界与不变量

- 不改正确口径：`npm run typecheck` = `tsc -b && tsc -p e2e/tsconfig.json` 仍为唯一入口，`build` 仍为 `tsc -b && vite build`。
- 不弱化既有 6 项断言：根配置仍须 solution-style、`typecheck` 仍须显式覆盖 e2e 项目、README 约定仍须在位。
- 不重开 `GOAL-008`；不改变 Root 六阶段分母与 `progress`；不关闭 `R5-I-004`。
- 独立审计按用户指令执行：本地 grok build · 模型 grok 4.6 · 思考强度 **xhigh**（项目级路径见 `docs/architecture/independent-audit-execution.md`，强度按本次用户指令提高）。
