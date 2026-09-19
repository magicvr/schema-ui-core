---
id: GOAL-044-w32-r4-residual-seams
title: W32 · R4 三项残余修复（列值本地化 · 表格定向刷新 seam · 空闲不轮询）
status: done
parent: GOAL-001-design-implementation-conformance
created: 2026-09-19
updated: 2026-09-19
version: 1.0.0
progress: 4/4
plan_refs:
  - VP-010-design-implementation-conformance
primary_plan: VP-010-design-implementation-conformance
vision_ref: schema-ui-core-admin-foundation@0.4.0
---
# GOAL-044 · W32 · R4 三项残余修复

## 概述

`[workspace-038-batch-operations-and-job-center]` 的 R4（结果中心）以 **cross** 审计（self `A-001` + grok-build independent `A-002`，均 `pass`、开放 required = 0）关门时留下 3 条 low 级 recommended，性质**全部是跨页面的通用能力缺口**，不是 VP-038 交付范围内的缺陷：

| # | 残余（来源） | 为何不在 038 内修 |
|---|--------------|-------------------|
| ① | 状态列以原始状态码呈现，未逐值本地化（038 `A-001` F-002） | 需要**通用列值本地化能力**（任何含枚举列的页面同样受益：钱包状态、定时任务 enabled、用户 MFA…）；在 038 内只能做成 jobs 专用特例 |
| ② | 自动刷新依赖 `reloadList()`，而该 seam 会清空本页表选择（038 `A-001` F-003） | 需要**渲染器「定向刷新指定表格且不清空选择」seam**，改动 ADR-0022 D2 的现有语义 → 跨全部列表页 |
| ③ | 无进行中作业时仍按档位轮询（038 `A-001` F-004） | 依赖 ② 提供的行可见性/表格定位能力；单独在 038 内做会另造一次性机制 |

用户 2026-09-19 指令：

> 三条都修，但是判断一下是直接修，还是需要再工作区开启一个子目标来承载治理上下文，如果后者比较好，则先开再工作区10开启一个承载这三项修复的子目标，然后本工作区有界接受，转由工作区10的新子目标执行修正——反之则直接进行修正。

**结构选型判定**：三项均为**跨工作区、跨页面**的通用能力（渲染器/表格层），且 ① 涉及协议面的本地扩展口径、② 涉及 ADR-0022 既有语义边界 —— 按 AGENTS §6e 属「独立树/跨边界」类，落在 038 会与其「不触碰 pinned 工件、不扩渲染器能力」的冻结非目标冲突。故选择**后者**：在 `workspace-010-design-implementation-conformance`（VP-010 持续符合性程序）开一个波次子目标承载，038 侧**有界接受**（`accepted-residual`，含范围与复核触发 = 本目标交付）并移交。先例：`GOAL-043-w31-cross-workspace-residual-closeout`（VP-037 关门后残余统一收口）。

## 范围

- **① 通用列值本地化（本地扩展）**：渲染器表格列支持「值 → i18n 键」映射（拟沿用既有本地扩展姿态，如 `badgeStyleField`/`truncate`，**不改 pinned 工件**、不新增 pinned 能力 id）；`jobs` 结果中心六态列接入。
- **② 表格定向刷新 seam**：渲染器/CRUD seam 支持刷新**指定表格**且**不清空表选择**，并保持既有 in-flight 合并语义；`jobs-auto-refresh` 改用它。
- **③ 空闲不轮询**：`jobs-auto-refresh` 仅在目标表格存在非终态行时 tick（依赖 ② 的行可见性/定位能力）。
- **④ 回填**：`[workspace-038…] GOAL-005` `A-003` 的 3 条 `accepted-residual` 按 P-003 回填为 `fixed`（附证据），只加闭合注记，不改其 status/progress。

## 明确非目标

- 不重开 VP-038 / workspace-038；不改 Job 六态合同、同步 `batch-delete`、任何 pinned 协议工件。
- 不做逐页全量设计审视（那是新波次）；不把 ①②③ 扩展成通用 BI/主题/权限改动。
- 不为「本地扩展」新增 pinned 能力 id 或修改 `docs/schemas/**`、`apps/web/src/protocol/upstream/**`。

## 高层路线图（P-001）

1. **C1 · 勘察与方案冻结**：三项的落点、语义与边界；① 的本地扩展命名与 fail-closed 口径；② 与 ADR-0022 D2（reload 清空选择）及 `refreshList`（display-only）的关系裁定；③ 的判定数据来源；跨工作区可写范围与用户授权登记。证据 → `01-decision/D-001`。
2. **C2 · 实施**：渲染器 ① 列值映射 + ② 定向刷新 seam；`jobs.json` 六态接入；`jobs-auto-refresh` 改用 ② 并实现 ③。
3. **C3 · 回归与证据**：交互级测试（含 ② 不清空选择、③ 空闲零请求、① 六态本地化）+ 变异验证；全量回归（Go + vitest + typecheck）；038 残余闭合回填。
4. **C4 · 审计与投影**：self 审计（按风险判定是否追加 independent）；`goal-tree.md`/`workspace.md` 同步；Root 保持 active 程序容器。

## 成功检查点

- [x] **C1 方案冻结**：`D-001` 落盘（含本地扩展口径、seam 语义、用户授权与跨区可写范围）；`I-044-001`～`003` 关闭为 `verified`。
- [x] **C2 三项实施**：① 列级 `valueLabels` + jobs 六态接入；② `refreshTable` 定向刷新 seam（保留选择）；③ `activeStatuses` 空闲不轮询。证据 `E-001` §2，checkpoint `c2ea042b`。
- [x] **C3 回归与回填**：三处变异验证（去 `valueLabels` / 让 `refreshTable` 清空选择 / 禁用空闲判定，均实测变红后还原）；`go test ./...` 全绿、`vitest` 119 files / 1468 tests 全绿、`typecheck` exit 0；038 `A-003` 三条回填 `fixed`。
- [x] **C4 审计与投影**：self `A-001` `pass`（0 required + 3 recommended）→ `A-002` 响应三条全 `fixed`，开放 required = 0；`goal-tree.md`/`workspace.md` 同步；Root 保持 active 程序容器。

## 关门（2026-09-19）

`done · 4/4`。三项残余以**通用能力**形态交付（渲染器列值本地化 / 定向刷新 seam / 行可见性驱动轮询），既闭合了 `[workspace-038] GOAL-005` 的 F-002/F-003/F-004，也为其它枚举列页面与轮询控件提供同一能力；`reloadList()` 的 ADR-0022 D2 语义以对照测试钉住未变。审计模式 `self`（改动为渲染器本地扩展与内部 seam，不触安全/数据/迁移/发布面；跨工作区的用户可见效果由 038 的 R5 浏览器/自动化回归复核）。

## 信息就绪与未知项（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-044-001 | required | ① 的本地扩展形态：复用既有列属性命名族还是新增 `valueLabels`？是否会与 pinned `component-registry.json` 的 `format: tag`+`tagMap` 语义重叠？ | C2 | C1 | 读 `docs/schemas/component-registry.json` 的列定义与仓库既有本地扩展（`badgeStyleField` 等）先例 | **verified** | — | `D-001` §1：新增列级 `valueLabels`（值 → i18n 键）；pinned `tagMap` 是字面量映射且本仓库未实现，二者划清边界 |
| I-044-002 | required | ② 的 seam 语义：刷新指定表格时是否清空选择？与 `reloadList()`（ADR-0022 D2 清空全部选择）和 `refreshList()`（display-only）的关系与边界？ | C2 | C1 | 读 `render.tsx` 的 `reloadList`/`refreshList`/`fetchList`/`selections` 与 SchemaTable 的取数路径 | **verified** | — | `D-001` §2：三条 seam 分工表 + 保留选择 + 不删 in-flight 键的取舍理由；对照测试钉住 |
| I-044-003 | non-blocking | ③ 的判定数据来源：组件如何得知目标表格是否还有非终态行（组件无行数据 seam）？ | C2 | C1 | 评估「② 暴露的行可见性」与「组件自持查询」两条路径 | **verified** | — | `D-001` §3：采用 ② 的行注册表（`publishTableRows`/`tableRows`，ref 支撑、按需读取）；行不可得时保守刷新 |

## 父目标

- `[workspace-010-design-implementation-conformance]` `GOAL-001-design-implementation-conformance`（长期程序容器，保持 active）。

## 台账布局

本目标从第一条记录起使用平铺 ledger：`01-decision/`、`02-execution/`、`03-audit/`，并保留 `attachments/`。

## 备注

- 本目标是 VP-010 持续符合性程序的一个**波次子目标**，不是新 VP、不改变 Charter；Root 保持 active 程序容器。
- 跨工作区写入（`apps/web` 渲染器为多工作区共享代码；闭合回填进 workspace-038 台账）由用户 2026-09-19 指令显式授权，范围见 `D-001`。
- 038 侧的 `accepted-residual` 以本目标交付为**复核触发**：本目标 C3 完成后回填为 `fixed`。
