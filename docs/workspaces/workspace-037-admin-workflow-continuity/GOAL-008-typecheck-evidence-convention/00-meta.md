---
id: GOAL-008-typecheck-evidence-convention
title: 类型检查证据约定纠偏与防复发
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

# GOAL-008 · 类型检查证据约定纠偏与防复发

## 概述

承接 R6 C6 审计 A-002 的 **F-005**（high required）：`apps/web/tsconfig.json` 是 solution-style 配置（`{"files": [], "references": [...]}`），因此**裸 `tsc --noEmit` 不检查任何源文件**，恒返回 exit 0。仓库内多处阶段性验证却以该命令作为「类型检查通过」的证据，属空转证据。

用户于 2026-09-18 按 P-004 裁决处置路径为 **方案 A：立独立目标系统性处置**（不追溯式加注、不记为残余）。本目标即该独立目标，承接「纠正口径 + 固化正确入口 + 防复发」，不重开 R6。

已于 2026-09-18 以 `done · 4/4` 完成：口径固化、正确入口、防复发守卫（含变异验证）与自审 `pass` 均已落盘，F-005 在本目标承接范围内闭环。

## 范围与边界

- 确立并固化**正确的类型检查口径**（`tsc -b` + e2e 项目 `tsc -p`，与 `apps/web/package.json` 的 `build` 脚本及 `apps/web/README.md` 既有约定一致）。
- 让正确入口**可发现、可复用**（脚本入口 + 文档 + CI 门禁），使后续阶段不会因为「随手敲的命令恰好空转」而再次产生失实证据。
- 登记受影响的历史条目范围，作为可追溯清单；**是否为其他工作区追溯更正**留待用户另行路由（`I-008-004` deferred）。
- 防复发机制形态见 `D-002`（结构守卫测试 + CI 显式门禁组合）。

明确非目标：不修改 `docs/workspaces/workspace-002|009|010|011` 等其他工作区的 canonical 台账（AGENTS §6c 禁止跨区写入）；不改 `tsconfig` 的项目引用结构或编译目标；不改 `apps/api`（Go）验证口径；不重开 R6 `GOAL-007` 的视觉范围。

## 高层路线图

1. **C1 · 影响面与正确口径基线**：确认 `tsconfig` 结构、证明空转、确定正确口径与受影响条目范围。已完成，证据见 `D-001`、`E-001`。
2. **C2 · 固化正确类型检查入口**：在 `apps/web` 提供显式 `typecheck` 脚本并在 README 记录约定。已完成，证据见 `E-002`。
3. **C3 · 防复发守卫**：实施结构守卫测试与 CI 显式门禁，并经变异验证。已完成，证据见 `D-002`、`E-003`、`E-004`。
4. **C4 · 审计与交付**：完成 self 审计（`A-001`，`pass`，开放 required = 0）并向 Root 投影。已完成，证据见 `A-001`、`E-005`。

## 成功检查点

- [x] C1：`tsconfig` solution-style 结构、空转证明（注入类型错误对比 `tsc --noEmit` exit 0 vs `tsc -b` exit 2）与受影响范围已形成可核对基线。
- [x] C2：`apps/web` 提供 `npm run typecheck` 并在 README 记录；实测该入口对真实类型错误返回非零。
- [x] C3：防复发守卫已实施并可核对（6 断言 + CI 门禁，5/5 变异捕获）；另发现并修复 e2e 项目未被检查的第二层同类缺口。
- [x] C4：self 审计 `A-001` verdict `pass`，开放 required = 0；Root 投影完成。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-008-001 | required | `tsconfig` 结构与裸 `tsc --noEmit` 是否真的空转？ | C1/C2 | C1 | 读取 `apps/web/tsconfig*.json`；注入类型错误对比两条命令的退出码与输出 | verified | 2026-09-18 已完成 | `D-001`、`E-001` |
| I-008-002 | required | 受影响的历史条目范围有多大？ | C1/C3 | C1 | 全仓检索裸 `tsc --noEmit`；追溯 `tsconfig.json` 引入时点 | verified | 2026-09-18 已完成；跨区条目只登记不代改 | `E-001` |
| I-008-003 | required | 防复发守卫应采用什么形态？ | C3 | C3 前 | 评估结构断言 / CI 引用 / 约定权威化；必要时问用户 | verified | 2026-09-18 决策为结构守卫测试 + CI 门禁组合并实施 | `D-002`、`E-003`、`E-004` |
| I-008-004 | non-blocking | 其他工作区历史条目是否追溯更正？ | 范围外 | 用户路由时 | 用户另行决定；本目标不跨区写入 | verified | 2026-09-18 用户授权追溯更正（Root `D-015`）；已加 11 处勘误注记、2 处复核确认有效 | `E-006`；`A-002` |

## 父目标

- `[workspace-037-admin-workflow-continuity]` `GOAL-001-admin-workflow-continuity`。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- 本目标**不是** Root 的纲领阶段，不改变 Root 六阶段分母与 `progress`（开设时为 `4/6`，R6 关闭投影后为 `5/6`）；它是 Root 下的整改子目标，与 R6 的视觉范围相互独立。
- 本目标完成使 F-005 在承接范围内闭环，从而解除用户为 R6 关门设定的前置条件（「等 F-005 处置落定后再关门」）。
- `I-008-004`：2026-09-18 用户授权追溯更正后，已在 workspace-002（复核确认有效，未改）、workspace-009、workspace-010、workspace-011 加 11 处勘误注记并落盘 `E-006`，状态转 `verified`。**新发现** `A-002 F-001`（守卫以 `-p` 令牌判定检查型调用，`-p tsconfig.json` 空转仍会通过）已由整改子目标 `GOAL-010-typecheck-guard-hardening` 以 `fixed` 路径闭合（2026-09-18，`GOAL-010` `done · 4/4`，其 `A-002`/`A-003` 独立审计 `pass`）；`A-001 F-002`（守卫未正向断言 CI 步骤存在）与 `A-002 F-002`（全仓 `tsc` 简写未逐条裁定）为 recommended 保持 open。
