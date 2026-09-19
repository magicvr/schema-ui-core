---
id: GOAL-006-r5-evidence-and-closeout
title: R5 证据与关门
status: active
parent: GOAL-001-batch-operations-and-job-center
created: 2026-09-19
updated: 2026-09-19
version: 0.3.0
progress: 3/4
plan_refs:
  - VP-038-batch-operations-and-job-center
primary_plan: VP-038-batch-operations-and-job-center
vision_ref: schema-ui-core-admin-foundation@0.4.0
---

# GOAL-006 · R5 证据与关门

## 概述

承接 Root `GOAL-001` 的**最后一个**纲领阶段 **R5**：把 R1～R4 的交付事实收敛成 VP-038 的**方向级退出证据**，并完成组合投影与关门裁决的入口。

R1～R4 均已 `done · 4/4`（R4 的 cross 审计两腿 `pass`、开放 required = 0，其三条残余已由 `[workspace-010] GOAL-044-w32-r4-residual-seams` 当日交付并回填 `fixed`）。本阶段**不做新功能**，只做证据、回归、独立意见、残余登记与投影。

## 范围与非目标

### 本目标范围

- **退出矩阵**：逐条核对 `VP-038` §方向级退出判据 1～7，每条给出**可核对证据**（工作区内的 E/A 条目、测试命令与结果、代码位置），不得用叙事代替证据。
- **浏览器/自动化回归**：e2e（Playwright）在 admin profile 上跑通；记录命令、结果与已知旁路。
- **独立意见**：按项目级决策 [independent-audit-execution.md](../../../../architecture/independent-audit-execution.md) 调用本地 grok build 做**关门独立审计**（`source: independent`），落盘 `03-audit/A-000`；由编排器合并响应。
- **残余登记**：本区残余（含 R5 自身 finding 与跨区移交项）在 `docs/vision/roadmap.md`「未决项统一登记」可查。
- **组合投影同步**：`goal-tree.md`、`workspace.md`、Root `GOAL-001` 检查点与 `VP-038` 判据状态同步；**VP-038 关门须用户书面确认**（P-004 / P-006）。

### 明确非目标

- 不新增功能、不改 Job 六态合同、不动同步 `batch-delete`、不触碰 pinned 协议工件。
- 不重开 VP-012/011/037/036；不解除 Redis/MQ/多实例/搜索 gated。
- **不自行宣布 VP-038 `closed`**：那需要用户书面确认；本目标只准备证据与提请。

## 成功检查点

- [x] **C1 退出矩阵**：VP-038 判据 1～7 逐条附证据；未达成或有界项显式标注（含残余登记指向）。证据：`02-execution/E-001-r5-exit-matrix.md` §1（判据 1～6 **达成**；判据 7 依赖本目标的 C3/C4）。
- [x] **C2 浏览器/自动化回归**：e2e 结果落盘（命令 + 通过/失败明细 + 覆盖说明）。证据：`E-001` §2/§3 —— 首次运行暴露**既有挂具顺序缺陷**（VP-036 新增 spec 消费 fresh-seed 前提），隔离复验确认与本 VP 无关，重命名修复后 **16 passed / 4 skipped / 0 failed（exit 0）**。
- [x] **C3 独立意见与残余登记**：independent 审计落盘并合并响应（开放 required = 0）；残余在 roadmap 登记节可查。证据：`03-audit/A-002`（independent · grok-4.6 high · **pass**，0 required + 4 recommended；独立复跑 Go 点名测试/jobs 包/迁移 checksum/6 个 vitest 文件 64 例/**全量 e2e 16 passed·4 skipped·0 failed**）；`03-audit/A-003`（响应：F-001 登记 roadmap bounded residual；F-002 登记 → 经**用户 2026-09-19 指令**补测后转 `fixed`；F-003 由登记闭合；F-004 索引纠偏）；`03-audit/A-004`（补测：新增 `e2e/jobs-result-center.spec.ts`，admin profile 端到端跑通；双 profile 全量 mvp 16 passed·5 skipped·0 failed / admin **17 passed·4 skipped·0 failed**）。
- [ ] **C4 组合投影与关门提请**：目标树/工作区/Root 同步；向用户提请 VP-038 关门（**等用户书面确认**，不静默）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-038-017 | required | e2e 覆盖是否足以支撑「浏览器/自动化回归」判据？作业/结果中心是否被真实浏览器路径触及？ | C2 | C2 前 | 读 `apps/web/e2e/**` 的用例清单与断言，判定覆盖面与缺口 | **verified** | — | `E-001` §3 + `A-002` 独立复跑：11 个 spec **无** jobs 结果中心端到端路径，且默认 `APP_PROFILE=mvp` 不含 `admin.jobs`；缺口已登记 roadmap（bounded residual），判据 4/5 的浏览器侧证据来自真实 `jobs.json` 的交互级 24 例 + HTTP 契约测试 |
| I-038-018 | non-blocking | `I-038-006`（历史作业保留/清理）在关门时的处置口径 | C4 | C4 | 确认其 `deferred · non-blocking` 与 roadmap 登记一致 | **verified** | — | 保持 `deferred · non-blocking`（责任人 `/vision`，触发 = 作业表容量/保留期出现真实需求）；与 VP-038 信息表一致，关门不因此受阻 |
| I-038-019 | required | VP-038 判据 7「组合投影同步且经用户书面确认」的关门形态 | C4 | C4 | 按 VP-037 先例（VRev + 用户书面确认）给出提请文本 | **open** | — | 待用户书面确认（`A-002` 独立腿明确：本条未决前 **VP-038 不具备关门条件**，不得把 `pass` 读成 `closed`） |

## 父目标

- `GOAL-001-batch-operations-and-job-center`（Root 的 R5 纲领阶段子目标；Root 现为 `active · 4/5`）。

## 台账布局

本目标使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/` 作为证据载体。

## 备注

- `progress: 0/4` 只由上方 4 个显式检查点派生；不放行阶段、不关闭 finding、不覆盖 status。
- 本目标是**最后一个纲领阶段**：其完成使 Root `GOAL-001` 具备关门条件；VP-038 的 `closed` 状态属愿景层，须 `/vision` + 用户确认。
