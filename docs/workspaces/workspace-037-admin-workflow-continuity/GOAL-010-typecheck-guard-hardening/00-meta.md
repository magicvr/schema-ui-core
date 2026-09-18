---
id: GOAL-010-typecheck-guard-hardening
title: 类型检查守卫加固（`-p` 目标有效性）
status: done
parent: GOAL-001-admin-workflow-continuity
created: 2026-09-18
updated: 2026-09-18
version: 1.2.0
progress: 4/4
plan_refs:
  - VP-037-admin-workflow-continuity
primary_plan: VP-037-admin-workflow-continuity
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-010 · 类型检查守卫加固（`-p` 目标有效性）

## 概述

承接 `GOAL-008` 自审 `A-002` 的 **F-001**（recommended）：`apps/web/src/typecheck-convention.guard.test.ts` 判定「是否为检查型调用」时，只检查命令行里是否出现 `-b` / `--build` / `-p` / `--project` **令牌**，不校验 `-p` 所指配置自身是否选择源文件。因此

```
tsc --noEmit -p tsconfig.json
```

会被守卫判为合规，而 `apps/web/tsconfig.json` 是 solution-style（`files: []`、仅 `references`）；非 build 模式下 TypeScript 只编译该配置选中的文件，于是该命令**编译空程序、恒 exit 0**。2026-09-18 以注入类型错误实测确认（`GOAL-008 E-006`）：裸 `tsc --noEmit` 与 `-p tsconfig.json` 均 exit 0，`-p tsconfig.app.json` 与 `tsc -b` 报 `TS2322`。

这正是 `GOAL-008` 要防的失效模式换了一个 `-p` 外壳——守卫是唯一的防复发装置，所以用户于 2026-09-18 指示「开一个小整改子目标（GOAL-010）加固守卫并补 `-p tsconfig.json` 变异用例」，并在其后执行交叉审计（本地 grok build · grok 4.6 · 思考强度 xhigh）后关门。

本目标已于 2026-09-18 以 `done · 4/4` 完成：判定改为按 `-p` 目标配置内容、补齐 18 行合成变异用例与动态目标解析，5 种自设变异 + 审计员 N3 变异均被捕获；self `A-001` `pass`、grok 独立审计 `A-002` `pass`、finding-closure 复审 `A-003` `pass`，开放 required = 0。

## 范围与边界

- 加固 `typecheck-convention.guard.test.ts`：`-p` / `--project` 目标必须解析到**自身选择源文件**的配置（`selectsOwnSources`）才计为检查型调用；无法解析为有效配置时 fail closed（计为违规）。
- 为 `-p tsconfig.json`、`-p .`、`--project tsconfig.json`、不存在/不可解析配置等形态补**合成变异用例**（表驱动的正/反例），使该缺口一旦回归就被测试直接指出。
- 用真实可执行面变异（CI workflow / `package.json` 脚本）验证扫描路径确实使用加固后的判定。
- 保持既有 6 项断言的语义不弱化：根配置仍须 solution-style、`typecheck` 仍须覆盖 e2e 项目、README 约定仍须在位。

明确非目标：不改 `tsc` 口径本身（正确口径仍是 `tsc -b` + `tsc -p e2e/tsconfig.json`）；不改 `tsconfig` 项目结构或编译目标；不逐条裁定全仓 300+ 行「tsc 简写」（`GOAL-008 A-002 F-002` 保持 open）；不改 `apps/api`；不改跨工作区历史台账；不重开 `GOAL-008`（本目标是其 finding 的整改承接，非其重开）。

## 高层路线图

1. **C1 · 口径收紧与证据基线**：确认空转形态边界（哪些 `-p` 有效、哪些无效）并冻结判定规则。已完成，证据见 `D-001`、`E-001`。
2. **C2 · 守卫加固与合成用例**：实现目标配置解析与 fail-closed 判定，补齐正/反例表。已完成，证据见 `E-002`。
3. **C3 · 变异验证与全量回归**：以真实可执行面变异证明扫描路径有效，并复跑全量测试与类型检查。已完成，证据见 `E-003`。
4. **C4 · 双审与投影**：self 审计 → 本地 grok build 独立审计（grok 4.6 · xhigh）→ 合并响应与 Root 投影。证据见 `A-001`、`A-002`、`E-004`。

## 成功检查点

- [x] C1：空转形态边界与判定规则冻结；`-p tsconfig.json` 空转有注入错误实测证据（`E-001`）。
- [x] C2：守卫按「目标配置必须自身选择源文件」判定，含 `-p tsconfig.json` 等 16 行合成用例表与动态目标解析（`E-002`）。
- [x] C3：5/5 变异被捕获（含 CI/`package.json` 真实面变异）；全量 Vitest 112/1428 与 `npm run typecheck` 通过（`E-003`）。
- [x] C4：self 审计 `A-001` `pass`；grok build（grok 4.6 · xhigh）独立审计 `A-002` `pass`（0 required）与 finding-closure 复审 `A-003` `pass`（F-001/F-002 均 fixed）；开放 required = 0，已于 2026-09-18 关门（`E-004`）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-010-001 | required | `-p tsconfig.json` 是否真空转？边界在哪？ | C1/C2 | C1 | 注入类型错误，对比 `tsc --noEmit` / `-p tsconfig.json` / `-p tsconfig.app.json` / `-b` 的退出码与输出 | verified | 2026-09-18 已完成（`GOAL-008 E-006` 实测，本目标复算） | `E-001` |
| I-010-002 | required | 守卫应如何判定 `-p` 目标有效，且不产生既有可执行面误报？ | C2 | C2 | 设计解析规则并对全仓可执行面复算 | verified | 2026-09-18 已完成 | `D-001`、`E-002` |
| I-010-003 | required | 独立审计的 provider、模型与思考强度？ | C4 | C4 | 用户 2026-09-18 指令：本地 grok build · grok 4.6 · 思考强度 xhigh | verified | 2026-09-18 已按指令执行 | `A-002` |

## 父目标

- `[workspace-037-admin-workflow-continuity]` `GOAL-001-admin-workflow-continuity`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- 本目标**不是** Root 的纲领阶段，不改变 Root 六阶段分母与 `progress: 5/6`；它是 Root 下的整改子目标，与 `GOAL-008`/`GOAL-009` 平行。
- 目标已关闭 `GOAL-008 A-002 F-001`（recommended）为 `fixed`。该 finding 不阻断任何门禁，本目标是主动补齐防复发覆盖。
- 不改变 `GOAL-008` 的 `done · 4/4`；`GOAL-008 A-002 F-002`（全仓 `tsc` 简写未逐条裁定）与 `GOAL-009 A-001 F-001/F-002` 仍为 recommended open。
- R5 `GOAL-006` 的 `R5-I-004` 用户书面关门确认仍开放，本目标关门不替代、不关闭它，也不关闭 Root 或 VP-037。
