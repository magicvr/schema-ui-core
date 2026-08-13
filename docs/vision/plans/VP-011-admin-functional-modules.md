---
doc_type: vision-plan
id: VP-011-admin-functional-modules
title: 标准 Admin 功能模块（通用模块 + 常用业务领域 · 分档交付）
status: planned
vision_ref: schema-ui-core-admin-foundation@0.2.0
lead_workspace: workspace-011-admin-functional-modules
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
parent: null
---

# VP-011 · 标准 Admin 功能模块（通用模块 + 常用业务领域 · 分档交付）

## 意图

在已 closed 的基架波次（VP-001～008 固化的协议兼容 `I-PROTO-FULL-001`、模块架构与贡献 playbook、设计系统、locale/settings、全基架准入 `go`）之上，开始构建 **Admin 的实际功能模块**，作为 Charter 边界内首个「标准业务模块」交付波次。

功能范围分两步确定：

1. **有界调研（lead 工作区 Root 首个纲领阶段）**：收集业界通用 Admin 功能 + 常用业务领域（订单、钱包为典型代表，类目、通知等纳入候选池），形成候选池，并**对照已交付基架**区分「已覆盖」与「真·待建」，避免重复立项。
2. **三档分档（见下方法论）**：把候选池按优先级分成三档，作为后续波次立项依据。

**分层约定（方向级）**：本 VP 只固定**意图 + 三档方法论 + 方向级退出判据**；调研产出的**具体分档清单与逐功能路线图**落在 lead 工作区 Root 的**纲领路线图（P-001）**，由 Root 首个纲领阶段（有界调研）产出后回写 Root。**本 VP 正文不因调研结果逐条改写**；仅当调研结论动摇意图 / 边界 / 方向级退出判据时，才回 `/vision` 修订（极少数 strategic 走 re-align）。

### 三档方法论（方向级）

| 档 | 判定 | 交付节奏 |
|----|------|----------|
| **一等公民** | 业界普遍存在、几乎所有 Admin 都需要、且当前基架尚未覆盖的核心能力 | 第一批次（第一波） |
| **常用** | 高频使用但非普遍必备 | 第二批次（第二波） |
| **增补** | 低频、按需，可由后续 fork 项目按需启用 | 增补 backlog（按触发条件立项） |

### 与相邻 VP 的边界

| VP | 关系 |
|----|------|
| **VP-008**（准入 `go`） | **消费关系**：激活前必须完成并记录**消费前 freshness review**（见 §消费决策记录）；`go` 失效/暂挂时本 VP 实现门闩自动关闭，回流 VP-008 重验证或 P-004 裁决 |
| **VP-009**（生产加固） | 共享基架 Critical/High 缺陷回流 009 波次，不塞进本 VP 领域台账 |
| **VP-010**（设计—实现符合性） | blocker/major gap 回流 010 波次，不塞进本 VP 领域台账 |
| **VP-003/004** | 只读权威：模块契约与贡献 playbook 为模块实现对照 |
| **VP-005/006/007** | 只读基线：设计系统 / 整份协议面 / locale+settings 为已交付能力 |

- **不**修改 Charter `@0.2.0` 目的、成功边界或非目标（业务领域本就是 Charter 内的「后续 VP 候选能力」）。
- **不**重开历史 closed VP；**不**私增协议语义或跳过协议覆盖。

## 方向级退出判据

在同时满足下列方向时，本 VP **可以**有界或完整关门（证据必须在 lead 工作区目标内）：

1. 一等公民档模块已按 Root 分档清单交付，且每个纳入能力都有**协议驱动的范例页面 + 可验证路径**（符合 Charter「范例即验证」）。
2. 常用档模块已交付，或有明确、可追踪的第二波立项与验收证据。
3. 增补 backlog 已登记，并有明确的触发/增补机制（不由本 VP 强行关闭）。
4. 未改变 Charter 边界、未重开历史 VP、未私增协议语义。

## 消费决策记录（freshness review · 激活前必做）

> 依据 VP-008 §`go` 消费有效性：每个后续业务 VP **激活前**必须针对拟消费候选身份与解锁 scope 完成一次轻量 freshness review 并记录。本 VP 完成下列复核前**保持 `planned`，不得激活**。

### 复核结果（2026-08-14 首轮 · **暂挂 pending re-verification**）

| 字段 | 结果 |
|------|------|
| consumer_vp | VP-011-admin-functional-modules |
| go 候选身份 | `ed99e88`（VP-008 S5 候选，clean）——有效 commit，**但非当前 HEAD** |
| 当前 HEAD | `ac81d44`（fix(security): W5 leftover lows after zero mid/high scan） |
| 候选漂移 | **是**：`ed99e88` 落后 HEAD ~15 提交（含 VP-009 W1–W5 安全修复、VP-010 W1–W4 符合性整改） |
| patch/manifest/input digest | 原候选 `ed99e88` 已非运行面；须以最新 wave-close 候选重验 |
| 冻结命令 / 关键证据可执行性 | 未重跑（候选已漂移） |
| 协议 pin | **漂移**：VP-010 W3 已将 `schema-ui-docs` 从 `v2.7.0`（`ca9e5fe`）升级并 pin 到 `v2.8.0`（tag `521cff8` / content `4fae4605`）；**Charter 协议来源仍写 `v2.7.0`** |
| 外部输入 / 环境可用性 | 未验证（上游 v2.8.0 制品已 vendor/固定，见 VP-010 W3 E-005） |
| 最新 finding + residual 投影 | VP-009 W1–W4 done、VP-010 W1–W4 done；但 HEAD 含**未归档 W5**（VP-009 安全 W5、VP-010 recordView 声明字段） |
| 结论 | **failed → `go` 暂挂**（VP-008 §`go` 消费有效性：协议 pin 改变 + 候选漂移属失效触发） |

### 阻断项（激活前须先闭合）

1. **Charter / 协议来源漂移（`/vision` 层）** — ✅ **已闭合（2026-08-14 · VR-020 · editorial）**：用户裁决 A（additive pin bump）；Charter 协议来源 / 目标语义 / 成功边界 1 / H-001 已升至 `v2.8.0`（`521cff8`），`I-PROTO-FULL-001` v1.0.1 保留为 v2.7.0 历史分母、被 v2.8.0 覆盖。
2. **消费候选重验证**：`go` 原候选 `ed99e88` 已非 HEAD；须以最新 wave-close 候选重跑冻结命令/关键证据并更新 identity+digest，或按 VP-008 失效规则回流重验证。
3. **未归档 W5**：HEAD 含 VP-009 W5（安全遗留 low）、VP-010 W5（recordView 声明字段）未写入两区 goal-tree；须先归档并确认 `go` 判定（无影响/不暂挂）再消费。

> 上述三项闭合前，本 VP **保持 `planned`，不得激活**；共享基架/`go` 语义问题由 `/vision` 决定（重开 VP-008 或新建准入 VP），不塞进本 VP。

## 工作区绑定

| workspace_id | root_goal | role | joined | notes |
|--------------|-----------|------|--------|-------|
| workspace-011-admin-functional-modules | GOAL-001-admin-functional-modules | lead | — | `planned`，0 区；完成 freshness review 并激活后由 `/govern` scaffold；Root 首个纲领阶段 = 有界调研（收集 + 分档） |

## 关门记录

（仅 `closed` / `abandoned` 时填写。）

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| — | — | — | — | — |

## 规划修订短史

| date | change |
|------|--------|
| 2026-08-14 | 初创（`planned`）：用户确认结构选型 = 新 VP-011 + 新 workspace-011；Root 首阶段 = 有界调研；调研回写 Root 纲领路线图而非 VP；VP 只留意图 + 三档方法论；范围 = 标准 Admin 通用模块 + 常用业务领域（订单/钱包为典型）一起入候选池分档 |
| 2026-08-14 | 消费前 freshness review 首轮：**failed → `go` 暂挂**（候选 `ed99e88` 落后 HEAD `ac81d44`；协议 pin `v2.7.0`→`v2.8.0` 且 Charter 未同步；未归档 W5）；阻断项见 §消费决策记录 |
| 2026-08-14 | 阻断项 1 闭合（VR-020 editorial pin bump）；剩余阻断项 2（候选重验证）、3（W5 归档）交 `/govern` |
