---
id: E-024-open-typecheck-guard-hardening
doc: execution-entry
status: recorded
goal_id: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
parent: GOAL-001-admin-workflow-continuity
version: 1.0.0
---

# E-024 · 开设 GOAL-010 并完成守卫加固 C1～C3

2026-09-18，按用户 `D-016` 指示开设整改子目标 `GOAL-010-typecheck-guard-hardening`（非纲领，不计入 Root 分母），承接 `GOAL-008 A-002 F-001`，并于同日完成 C1～C3：

- C1：注入类型错误实测冻结空转边界——裸 `tsc --noEmit` 与 `-p tsconfig.json` 均 exit 0（空转），`-p tsconfig.app.json` 与 `tsc -b` 报 `TS2322`（真实检查）；根配置为 solution-style，其余 8 个 `tsconfig.*.json` 与 `e2e/tsconfig.json` 均通过 `include` 选择源文件。判定规则冻结为「`-b`，或 `-p` 目标自身选择源文件」，其余 fail closed（`D-001`、`E-001`）。
- C2：只改 `apps/web/src/typecheck-convention.guard.test.ts`（+239/−14），实现按配置内容判定、按位置判定引号的参数解析、插值目标通配解析，并新增 16 行合成变异用例表与 7 条目标解析单元断言；`typecheck`/`build` 改由同一行级检查器断言（强于原 `CHECKING_FLAG.test()`）（`E-002`）。
- C3：5/5 变异被捕获（含 CI workflow 与 `package.json` 真实面变异，旧守卫会放行）；`npm run typecheck` exit 0；Vitest **112 文件 / 1428 测试**通过；`git diff --check` 通过；`.github/**`、`package.json`、产品代码无差异（`E-003`）。

当前 `GOAL-010` 为 `active · 3/4`（C4 待 self 审计后的 grok build 独立审计与合并响应）。Root `progress` 仍为 `5/6`；`R5-I-004`、Root 与 VP-037 均未因本轮改动被关闭或推断。
