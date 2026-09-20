---
id: VRev-104-vp040-timestamptz-persistence-contract-activation
doc_type: vision-review
title: VP-040 DB 时间列 timestamptz 持久化合同 · 激活就绪审视
source: self
date: 2026-09-20
scope: VP-040-timestamptz-persistence-contract · activation / P-005 默认候选 / 架构类 freshness / slug
verdict: pass
open_required: 0
status: recorded
auditor: /vision
created: 2026-09-20
updated: 2026-09-20
parent: null
version: 0.1.0
---

# VRev-104 · VP-040 激活就绪

## 背景与触发

用户 2026-09-20 指令：「/vision 走流程激活 vp-040，没有问题的话，交 /govern 开设工作区」。本轮核对 VP-039 前置波次、Charter 对齐、VP-040 的激活门禁、P-005 信息就绪、架构类 freshness 与工作区命名边界。

本 Review 只授权 VP-040 激活与实现层工作区/Root scaffold；不宣称时间列合同已冻结、迁移已实施或 SQLite/PG 已完成物理类型改造。

## 1. 前置、对齐与意图可判定性

**pass**。

- VP-039 已于 2026-09-19 `closed` v0.3.0；VRev-103 self `pass`，满足 VP-040「VP-039 波次之后」硬前置。
- VP-040 `vision_ref` 精确匹配现行 Charter `schema-ui-core-admin-foundation@0.4.0`；现行 Charter 唯一且 `status: active`。
- 六条方向级退出判据与 R1 → R2 → R3 纲领路线图可核对；R1 合同/分母冻结仍是后续实现门禁。
- VP-040 仍是架构 C1 独立波次，不并入 VP-039，不塞 VP-010，不重开 VP-013 或 VP-020。

## 2. P-005 信息就绪与激活门禁

| 信息项 | 当前状态 | 本轮结论 |
|--------|----------|----------|
| `I-040-001` SQLite 与 PG 合同平等的物理类型 | **collecting** | 已登记激活用默认候选：SQLite `INTEGER` / PostgreSQL `BIGINT`，值语义暂按 UTC Unix seconds；这是 R1 的候选，不是最终冻结。现有毫秒字段必须在 R1 单独枚举，不能静默转换。 |
| `I-040-002` 首波时间列分母 | open / required | 阻断 R1/R2；不阻断本轮激活。 |
| `I-040-003` 存量升级与备份 residual | open / required | 阻断 R1 方案与判据 4；不阻断本轮激活。 |
| `I-040-004` 与 VP-020 展示合同的回归矩阵 | open / required | 阻断 R3；不阻断本轮激活。 |
| `I-040-005` 架构 freshness 与 VP-039 前置 | **verified** | VP-039 `closed` + 本报告 freshness PASS；满足激活门禁。 |

本轮不把 `I-040-001` 的候选写成已验证事实；R1 仍须由实现层目标完成最终物理类型、精度、NULL/零值、编解码与例外分母冻结。

## 3. 架构类轻量 freshness（`6197e802` → `b0a6789b`）

**PASS**，不暂挂 VP-008 `go`。

| 域 | 核对 | 结果 |
|----|------|------|
| 协议 pin / provenance | `apps/web/src/protocol/upstream` 区间零 diff；现行仍为 `schema-ui-docs@v2.9.0` / pinned `81aa1d8` | PASS |
| 依赖锁 | `apps/api/go.mod` / `go.sum`、`apps/web/package.json` 与 lockfile 区间零 diff | PASS |
| 迁移台账 | `apps/api/internal/store` 与 `apps/api/modules/*/migration` 区间零 diff | PASS |
| Profile 默认集 / 装配 | `kernel/profile.go` 与 Profile/Manifest 装配规则区间零 diff；VP-039 仅增加 runtime-mode wiring，不改变模块集合 | PASS |
| 区间代码与红线 | 区间实现变更属于 VP-039 已审结的维护提示/版本诊断；无新 schema、ORM、第三库、Redis/MQ/A3 或协议面变更 | PASS |

本 VP 可在现有 Charter/VP-008 `go` 语境下进入实现层；实现中若改变 Profile 默认集、Manifest 装配或共同门禁，必须暂停并重新做 freshness/go 消费核对。

## 4. 组合对齐与 slug

**pass**。

按本轮「激活 VP-040 + 交 `/govern` 开设工作区」指令及 VP-013～039 的主干命名惯例，绑定如下：

| 项 | 值 |
|----|-----|
| workspace_id | `workspace-040-timestamptz-persistence-contract` |
| Root | `GOAL-001-timestamptz-persistence-contract` |
| vision_role | `delivery` |
| primary_plan | `VP-040-timestamptz-persistence-contract` |

该 delivery 工作区不改变 Charter `primary_workspace`（仍为 `workspace-001-mvp-admin-foundation`），也不与其他工作区建立 `parent` 关系。

## 5. 激活事务边界

| 项 | 结论 |
|----|------|
| Charter 目的/边界/非目标 | 不变；无 strategic 修订、无 re-align |
| VP-040 意图/退出判据 | 不改方向；仅从 `planned` 激活为 `active` |
| 代码与 schema | 不在本轮实施；不宣称任何迁移完成 |
| gated 基础设施 | Redis / MQ / 多实例 / A3 / ORM / 第三库保持 gated/排除 |
| 实现层交接 | 交 `/govern` scaffold 显式 delivery 工作区与 Root；Root 初始 `active · 0/3` |

## Verdict

**pass（open required = 0）**。

VP-040 可从 `planned` 激活为 `active` v0.2.0，并由 `/govern` 建立 `workspace-040-timestamptz-persistence-contract` 与 Root `GOAL-001-timestamptz-persistence-contract`。`I-040-001`～`004` 仍是后续 R1/R2/R3 的实现层信息门禁，不因激活而关闭。

## Findings

### 必改（required）

无。

### 建议（recommended）

无新增。VRev-101 的 `V-F132` 激活前默认候选要求已由本报告与 Root D-001 落盘满足；最终 R1 冻结仍由 `I-040-001` 承接，原 finding 不被改写。

## 声明

本意见为 `/vision` self Review，不冒充 independent Vision Review 或 Goal independent 审计。它不替代 Root `03-audit`，不宣称实现、迁移、验收或关门已完成。
