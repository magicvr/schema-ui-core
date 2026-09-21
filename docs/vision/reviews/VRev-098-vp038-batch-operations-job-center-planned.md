---
id: VRev-098-vp038-batch-operations-job-center-planned
doc_type: vision-review
title: VP-038 Admin 批量操作与异步结果中心 · 计划阶段意图审视
source: self
scope: VP-038-batch-operations-and-job-center · planned
verdict: pass
open_required: 0
status: recorded
date: 2026-09-19
auditor: /vision
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.1.0
---

# VRev-098 · VP-038 计划阶段意图审视

## 审视范围

新立 `VP-038-batch-operations-and-job-center`（Admin 功能分支 · 批量操作与异步结果中心）的计划阶段意图审视：

- Charter `@0.4.0` 对齐
- 结构选型（新 VP vs VP-010 波次 vs 架构分支）
- 首波范围与七条退出判据的可判定性
- P-005 信息需求与激活门禁
- 与 VP-012 / VP-011 / VP-037 / VP-036 / VP-008 `go` / 架构与业务域分支的边界
- 本轮只读事实核对（Job 运行时与批量面现状）

## 本轮只读事实核对

| 主张 | 证据 | 结论 |
|------|------|------|
| Job 六态运行时已交付但无通用产品面 | `apps/api/internal/jobs/model.go`（六态常量 + `Job` 结构）；`apps/api/modules/jobs/migration/provider.go:19` `Register` 返回 `nil`；`modules/jobs/migration/migration.go:1`–`2` 注释「migration-only owner and is deliberately absent from runtime profiles」 | ✅ 属实 |
| 「通用 Job 管理页」是 VP-012 显式排除项 | [VP-012](../plans/VP-012-shared-cross-module-contracts.md) 首波冻结表「异步 Job / 长操作」行：已交付 = 六态/进度/重试/取消/结果读取过期 + wallet reconcile 202；**不进本 VP** = 通用 Job 管理页；外部队列 | ✅ 属实 |
| 现行 Job 读面仅限单一 kind 且绑定 actor | `internal/jobs/repository.go:66` `GetForActor(id, kind, actorID)`；`:292` `ListRunnable` 为调度器内部查询；无「按权限列出全部/按数据范围列出」查询 | ✅ 属实 |
| 唯一 Job 消费者是 wallet.reconcile，经模块自有路由暴露 | `modules/wallet/jobs.go:19` `ReconcileJobKind = "wallet.reconcile"`；`kernel/profile.go:202` `admin.wallet` 路由含 `GET /api/wallet/jobs/{id}`、`/cancel`、`/retry`、`/result` | ✅ 属实 |
| 批量面只有同步 `batch-delete`，无异步与结果中心 | `internal/handler/resources.go:815`–`824`（注释：whole-batch semantics，成功返回 `{"deleted": n}` 由客户端 reload）；四处注册点 `kernel/profile.go:166,167,189,193` | ✅ 属实 |
| 导出/导入为同步 | `kernel/profile.go:181` `admin.data-transfer` 路由 `GET /api/export/{resource}`、`POST /api/import/{resource}`、`GET /api/import/{resource}/template`；`internal/handler/import.go:133` `importResource()` 同步实现 | ✅ 属实 |
| 协议侧批量能力已存在 | `apps/web/src/protocol/host-support.json:13` `actions.batch.request`；`apps/web/src/protocol/conformance/request-construction.ts:654`–`700` `batchMapping` 构造与 `SELECTION_KEYS_BODY_ONLY` 等约束 | ✅ 属实 |
| roadmap 已把「批量结果中心」登记为未立项 | [roadmap.md](../roadmap.md)「体验增强」行；「未决项统一登记」§三：批量结果中心 = **未立项**，触发 = 新立 VP，责任人 `/vision` | ✅ 属实 |

## 审视结论

**verdict: `pass`**（0 required，1 recommended）。

VP-038 的意图落在现行 Charter 成功边界 #3（产品化 Admin 体验，工作导向）与 #5（模块可组合、避免 Shell 中央注册路径）内，是 `roadmap.md`「体验增强」清单中**唯一同时满足「基础设施前置已交付」与「未消耗任何 gated trigger」**的未立项项：Job 六态运行时由 VP-012 交付、协议侧批量能力由 VP-006/VP-011 交付，本 VP 只补产品面。

结构选择成立：它不是 VP-010 的 as-designed / as-built 偏差整改（VP-012 把「通用 Job 管理页」作为**显式非目标**排除，不是漏做），也不是架构分支项（不新建队列/存储/多实例能力），因此应为**新 VP + 新 delivery 工作区**。

本 Review 仅确认 `planned` 意图，**不授权激活、开工作区或进入实现**。

## Charter 对齐

| 项 | 检查结果 |
|----|---------|
| `vision_ref` = `schema-ui-core-admin-foundation@0.4.0` | ✅ 精确匹配唯一 active Charter |
| 落在成功边界内 | ✅ 成功边界 #3（产品化 Admin 体验）与 #5（模块贡献可组合）直接覆盖 |
| 不改变 Charter 目的/边界/非目标 | ✅ 不改 Charter；不新增业务域；不改变单主线模块策略；不预制有状态横切能力 |
| 与现有路线图一致 | ✅ 承接 Admin 功能分支「体验增强」未立项项「批量结果中心」；不解除任何 gated 行 |
| 不消耗 trigger | ✅ 不碰 Redis（`RT-Q03`）、外部队列（`RT-Q02`）、多实例（A3）、搜索引擎（`RT-X01`/`RT-X02`） |

## 结构选型

| 问题 | 判断 | 依据 |
|------|------|------|
| 改 Charter 目的/边界？ | 否 | 现有 Admin 体验方向内的有界增量 |
| 同愿景新纲领波次？ | 是 | roadmap「体验增强」已登记、且为「未决项统一登记」§三点名未立项项 |
| 是否属于 VP-010 普通波次？ | 否 | VP-012 把「通用 Job 管理页」写成**显式非目标**，属未做的产品能力，不是既有设计意图与实现的偏差 |
| 是否属于架构分支？ | 否 | 不新建队列/存储/多实例能力；消费既有进程内 Job 运行时与 Store |
| 是否需要独立 Goal 树？ | 是 | 需独立冻结 Job 种类×作用域矩阵、批量异步契约、权限/Profile 矩阵与跨模块回归 |
| 结论 | **新 VP + 新 delivery 工作区** | 计划阶段 0 区；激活后交 `/govern` scaffold（slug 须用户确认） |

## 首波边界与退出判据

七条退出判据均可判定，且每条都能指回工作区证据：

- 判据 1（分母与契约）可由 `I-038-001`～`003` 的矩阵承接，分母是**现存 Job 种类与已注册批量面**，不是「所有可能的长操作」。
- 判据 2（通用可见性）对应现行 `GetForActor` 的局限（仅单一 kind + actor），明确要求权限与作用域 fail-closed。
- 判据 3（异步承接）**显式保护已交付的同步 `batch-delete` 语义不退化**，避免把既有合同改成 breaking。
- 判据 4～5 复用 VP-005/007 体验基线与既有权限/Profile 守卫，可浏览器与自动化回归。
- 判据 6 与判据 7 分别锁死基础设施边界与审计闭合。

排除项已点名：实体全文检索、Saved Views 重做、Toast 全局重做、新业务域、组织/多租户权限、Redis/MQ/多实例/第二持久化栈。

## P-005 信息就绪

| 信息项 | 级别 | 状态 | 门禁 |
|--------|------|------|------|
| I-038-001：Job 分母与可见作用域 | required | open | 阻断 R1 范围冻结与 R2 读面 |
| I-038-002：批量异步契约与协议面影响 | required | open | 阻断 R1 方案冻结与 R3 实施 |
| I-038-003：首波批量操作分母 | required | open | 阻断 R1 范围冻结与 R3 实施 |
| I-038-004：Profile / 模块矩阵边界 | required | open | **阻断激活**（须用户 P-004 裁决），并影响判据 5/6 与 VP-008 `go` 消费有效性 |
| I-038-005：激活前 Admin 类 freshness | required | open | 阻断激活与开区，不阻断 `planned` 登记 |
| I-038-006：历史作业保留与清理策略 | non-blocking | deferred | 不影响首波；真实容量/合规需求出现时由 `/vision` 复核 |

无 required 信息项阻断 `planned` 立项。激活前必须完成 `I-038-004`（用户裁决）与 `I-038-005`；R1 冻结前必须完成 `I-038-001`～`003`。

## Findings

### V-F125（recommended · 非阻断）

`I-038-004` 的 Profile / 模块边界是本 VP 唯一触及 **VP-008 `go` 消费失效触发项**（「Profile 默认集、模块适用矩阵改变」）的地方：作业中心既可能作为新模块（如 `admin.jobs`）进入默认集，也可能挂在既有模块的页面/路由上。建议在激活包内就把两种方案的**默认集影响面**写明并请用户裁决，避免实施期以「加个页面」为由静默改变模块矩阵——那会同时触发 `go` 暂挂与 freshness 重验证。

状态：`open · recommended`。不阻断 `planned` 登记；应由激活包或 R1 决策记录承接（`I-038-004` 已按 required 登记该门禁）。

### V-F126（recommended · 非阻断）

首波只承诺「**至少一条**真实批量操作以异步 Job 承接」。建议在 R1 就把这条操作的**选择依据与分母**（例如导出/导入 vs 批量启停 vs 对账）与其余保持同步的项一并落成矩阵，避免关门时无法判定判据 3 的「至少一条」是否满足、或反向滑向「把所有批量都改异步」。

状态：`open · recommended`。不阻断 `planned` 登记；由 `I-038-003` 与 R1 决策记录承接。

## 激活门禁（后续 `/vision`）

1. **用户 P-004 裁决**：`I-038-004` 的 Profile / 模块矩阵边界（新增模块 vs 挂既有模块；是否进默认集）。
2. **Admin 类 freshness review**：核对协议 pin（`v2.9.0` / `81aa1d8`）、依赖锁、迁移台账、Profile 默认集与装配、provenance 及区间变更；不满足时不得消费 VP-008 `go`。
3. **激活就绪 self Review**：确认 `I-038-004`/`I-038-005` 已关闭、slug 与工作区绑定边界。
4. 用户确认 delivery workspace / Root slug 后，交 `/govern` 创建 workspace + Root；本 Review 不创建 Goal 五件套。

## 声明

- 本 Review source = `self`，不冒充 independent。
- 不改变 Charter、其它 VP、工作区或 Goal status/progress。
- open required = 0；`V-F125` / `V-F126` 为 recommended，不阻断 `planned` 登记或后续激活准备。
- 下一步：若继续推进，使用 `/vision` 完成 `I-038-004` 用户裁决 + freshness review 后激活 VP-038；激活后交 `/govern` scaffold 工作区。
