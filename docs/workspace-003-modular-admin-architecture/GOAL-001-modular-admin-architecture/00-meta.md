---
id: GOAL-001-modular-admin-architecture
title: 单主线模块化 Admin 架构
status: active
parent: null
created: 2026-08-04
updated: 2026-08-05
version: 0.7.0
progress: 2/6
plan_refs:
  - VP-003-modular-admin-architecture
primary_plan: VP-003-modular-admin-architecture
serves_summary: 将现有生产级 Admin 基架在单一代码主线内演进为由薄内核、静态组合根、模块契约和启动时 Profile 驱动的模块化单体。
---

# GOAL-001 · 单主线模块化 Admin 架构

## 概述

本 Root 承接 [VP-003 · 单主线模块化 Admin 架构](../../vision/plans/VP-003-modular-admin-architecture.md)。本次仅建立独立的 delivery 工作区、路线图和信息门禁，作为模块化架构实施的可追溯范围；不把建区、文档或试点计划写成 R1-R6 的实现完成证据。

## 愿景对齐

| 字段 | 值 |
|------|----|
| `plan_refs` | `VP-003-modular-admin-architecture` |
| `primary_plan` | `VP-003-modular-admin-architecture` |
| Charter | `schema-ui-core-admin-foundation@0.2.0` |
| 工作区角色 | `delivery` |

方向、退出判据与非目标以 VP-003 和 [module-architecture.md](../../architecture/module-architecture.md) 为权威。本 Root 仅保留服务该意图所需的简短边界与路线图。

协议覆盖基线继承：本 Root 实施默认不扩大 [I-PROTO-001 v0.1.3](../../vision/plans/VP-003-modular-admin-architecture.md#继承的协议基线i-proto-001-v013)（Q2 权威覆盖表见 VP 该节链出的 workspace-001 附件）；扩大 domain / 改 exclude / 升上游协议版本须新决策并递增覆盖表版本。登记见 [I-007](#信息需求与阶段门禁) 与 [D-002](01-decision/D-002-a002-design-response.md)。

## 成功边界

成功边界分两层，**不得**互相替代：

1. **阶段可验收结果**（下表「阶段层」）：用于 R1–R6 阶段放行与子目标验收，对应纲领检查点。
2. **VP 关门必证七条**：以 [VP-003 方向级退出判据](../../vision/plans/VP-003-modular-admin-architecture.md#方向级退出判据) 全文为权威；本区关门须对 exit #1–#7 **逐条**具备实施证据、验证结果与关门审计。

**硬约束**：派生 `progress=6/6`（R1–R6 检查点全勾选）**不得**推导本 Root `done`，**不得**推导 VP-003 `closed`，也不得跳过 exit #1–#7 逐条取证。与 [goal-tree.md](../goal-tree.md) 维护说明一致。

### 阶段层（可验收）

- [x] **R1**：在完整模块/中央注册/迁移所有权清单与协议继承核对基础上，冻结可实施的模块、迁移、错误和回滚边界（含 Profile **候选/依赖盘点**，不含 Profile 精确集合冻结）。
- [x] **R2**：薄内核、框架无关模块契约、Fx 组合根与确定性后端聚合骨架就位；`mvp`/`admin` Profile **精确模块集合与配置覆盖顺序**已按 I-004 冻结，I-005 聚合/代理边界已 verified。
- [ ] **R3**：进行中；先收集并冻结 I-006 的静态 Manifest/Shell 兼容、移除和回滚边界，再完成 [VP-003 R3 A+B+C+D 门闩](../../vision/plans/VP-003-modular-admin-architecture.md#r3-通过门闩有界试点--继承固定历史评议输入) 后才可进入 R4；**禁止**以「试点模块已写出」替代门闩。
- [ ] **R4–R5**：现有一方 Admin 能力迁入统一模块契约；Profile、数据升级/恢复、健康诊断、Docker/代理与 fork 路径具有可核对证据。
- [ ] **R6 / 关门**：旧装配路径退出；exit #1–#7 均有本区证据与审计，且无开放 required 信息项或必改 finding。

### R 阶段 ↔ VP 退出判据映射

| R 阶段检查点 | 主要服务的 VP 退出判据 | 预期证据类型（摘要） |
|--------------|------------------------|----------------------|
| R1 契约与迁移基线冻结 | #2（契约/能力口径）、#3（迁移/tombstone 边界）、#1 的**盘点侧**（Profile 矩阵候选，非精确冻结） | 模块/注册/迁移清单；API/错误/回滚决策；I-PROTO 继承核对（I-007）；I-001～I-003 verified |
| R2 内核与组合根基础 | #2（薄内核/Fx/契约实现）、#4（聚合骨架/代理）、#1（Profile 精确集 + 覆盖顺序，I-004） | 可运行骨架、图校验、迁移收集、Manifest 聚合骨架；I-004/I-005 方案冻结证据 |
| R3 有界试点 | #2/#4/#5 的试点范围切片；#6 的试点路径（非全量） | VP R3 A+B+C+D 全过；V-1～V-4；四病灶切除；I-006 开发期兼容边界 |
| R4 全量一方模块迁移 | #6（能力迁入 + 退出中央特例）、#2（一方能力入统一契约） | 各模块迁移清单与行为/协议边界对照；Shell/中央注册特例清除证据 |
| R5 Profile、数据与运维收敛 | #1（Profile 运维/配置收敛）、#3（fresh/reconcile/升级恢复）、#7（文档/容器/诊断） | 双 Profile 案例、升级/恢复、readyz/诊断、fork 文档；**不**回写否定 R2 已冻结 Profile 集，除非新决策 |
| R6 旧路径移除与终态验收 | **#1–#7 全部** | 双轨/静态兜底删除；完整回归矩阵；exit 逐条证据包 + 关门审计 |

阶段层第 4 条（R4–R5）对应上表 R4+R5 行；阶段层第 5 条兜底要求 exit 七条全部取证，与 R6 行一致。

## 纲领路线图

六个检查点默认等权并原则上串行；同一阶段内可在相应信息门禁已满足后创建并行子目标。R1、R2 已完成并通过对应 close-out audit，Root 派生进度为 `2/6`；R3 正在由 GOAL-004 承接，R4-R6 尚未完成。R3 仍须先完成 I-006 信息边界，不能把已有模块声明或 R2 骨架当作试点通过。

| 阶段 | 名称 | 状态 | 说明 |
|------|------|------|------|
| R1 | 契约与迁移基线冻结 | 已完成 | 由 [GOAL-002-r1-contract-migration-baseline](../GOAL-002-r1-contract-migration-baseline/00-meta.md) 承接；C1-C4 证据、D-003～D-005、Grok A-004 independent 与 A-005 response 已落盘；冻结模块 API、迁移 tombstone、错误分类、兼容基线和回滚策略。**Profile 精确模块集合与覆盖顺序未在本阶段冻结**（I-004 留给 R2）。协议范围默认不扩大 I-PROTO-001 v0.1.3（I-007）。 |
| R2 | 内核与组合根基础 | 已完成 | 由 [GOAL-003-r2-kernel-composition-root](../GOAL-003-r2-kernel-composition-root/00-meta.md) 承接并以 `done 5/5` 收束；I-004/I-005 verified，Root R2 stage close-out 已落盘，Root progress 为 `2/6`。 |
| R3 | 有界试点 | 进行中 | 由 [GOAL-004-r3-bounded-pilot](../GOAL-004-r3-bounded-pilot/00-meta.md) 承接；先收集 I-006 删除/兼容/回滚边界，再用 operationlog/activity 与 settings 验证 Kernel 切口、四个病灶切除和 V-1～V-4。**未满足 VP-003 R3 A+B+C+D 不得进入 R4**；「试点模块写出」不等于通过。 |
| R4 | 全量一方模块迁移 | 未开始 | 将 users、roles、Schema CRUD 及其他现有 Admin 能力迁入统一能力契约，移除 Shell/中央注册表特例。范围受 I-PROTO-001 v0.1.3 继承约束。 |
| R5 | Profile、数据与运维收敛 | 未开始 | 完成 Profile **运维/配置收敛**与文档、fresh/reconcile、readyz/诊断、代理/容器、升级恢复与 fork 文档。R5 **不**否定 R2 已冻结的精确 Profile 集，除非新决策书面改写。 |
| R6 | 旧路径移除与终态验收 | 未开始 | 删除双轨与静态生产兜底，完成完整回归、双 Profile、升级/恢复、失败路径、容器/fork 验收；对 exit #1–#7 逐条取证后关门审计。 |

### 阶段审计模式（预置建议）

实施前仍可按风险确认；下列为 Root 默认建议，降低首次推进时的模式歧义（D-001 仅覆盖建区 `self`）。

| Scope | 建议模式 | 说明 |
|-------|----------|------|
| 建区 / 纯文档治理补强 | `self` 或 `none` | 已完成：建区 A-001；A-002 响应为文档修正 |
| R1 契约/迁移策略冻结 | 至少 `independent` | compatibility / migration 边界 |
| R2 内核与聚合骨架（常规实现） | `self`；触及生产路径契约时升 `independent` | 边界清楚的平台骨架 |
| R3 试点门闩放行 | `independent` 或 `cross` | 高影响门闩；禁止静默降级 |
| R4 全量迁移切片 | `independent` | migration / production |
| R5 运维与 Profile 收敛 | `self`；升级/恢复/容器放行倾向 `independent` | 按切片风险 |
| R6 旧路径删除与关门 | 至少 `independent` | release / close-out |

### 子目标拆分约定（最小）

- **默认**：按 R 阶段建立实施子目标（一阶段一主实施目标，阶段内可并行有独立交付价值的切片）。
- **信息项**：仅当收集本身有独立范围、依赖、交付证据或并行价值时，才升格为信息澄清/收集子目标。
- **禁止**：为 I-001～I-007 各机械创建两个子目标。

## 信息需求与阶段门禁

开放信息项不阻断本 Root 的建立或 A-002 设计补强；它们分别阻断列明的后续阶段，且不得在没有证据或用户书面 residual 的情况下被写成 `verified`。

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 验证 / 收集动作 | 状态 | 延期 / 复核 | 证据 / 结论 |
|----|------|-----------------|----------|--------------|-----------------|------|-------------|-------------|
| I-001 | required | 当前中央注册点、模块候选、能力归属与跨模块依赖的完整清单是什么？ | R1 方案冻结与实施 | R1 方案冻结前 | 对 API、Web、Shell、路由、导航和现有注册路径做可追溯盘点，并形成模块清单。 | verified | 2026-08-04 复核；R2 仅承接精确 Profile/I-004 | [GOAL-002 C1](../GOAL-002-r1-contract-migration-baseline/attachments/r1-c1-module-profile-inventory.md)；D-003；Grok A-004 independent + A-005 response；Root D-004。 |
| I-002 | required | 既有 `0001` 起迁移链与 seed 如何映射到全局模块所有权，升级/恢复夹具覆盖什么？ | R1 迁移策略冻结与 R2 实施 | R1 方案冻结前 | 核对迁移台账、checksum、快照/恢复、tombstone 与系统数据 reconcile 的现状和目标边界。 | verified | 2026-08-04 复核；R2/R4 承接实现与深度验证 | [GOAL-002 C2](../GOAL-002-r1-contract-migration-baseline/attachments/r1-c2-migration-seed-boundary.md)；D-003；Grok A-004 independent + A-005 response；当前 rollback=事务失败 rollback + pre-upgrade snapshot 恢复后备，seed 不等于 versioned reconcile，tombstone 为目标边界。 |
| I-003 | required | Fx、模块 API、Go 兼容范围和启动/就绪/停止/失败清理的错误语义如何固定？ | R1 契约冻结与 R2 实施 | R1 方案冻结前 | 固定候选版本、模块 API 边界和生命周期/错误分类，并记录取舍。 | verified | 2026-08-04 复核；具体 Fx 版本、Go type surface、stable error codes 和实现留给 R2 | [GOAL-002 C3](../GOAL-002-r1-contract-migration-baseline/attachments/r1-c3-lifecycle-contract.md)；D-004；Grok A-004 independent + A-005 response。R1 只冻结 Uber Fx 组合根候选、框架无关模块语义、核心六项/按需能力和 fail-closed 生命周期边界。 |
| I-004 | required | `mvp`、`admin` Profile 的精确模块集合与配置覆盖顺序是什么？ | R2 Profile 方案冻结与实施 | R2 方案冻结前 | R2 C1 已冻结并通过正反例：`mvp` 为 core.server-registration、core.auth-session、core.manifest-route、core.navigation-capability、core.schema-render、core.operationlog、admin.users、admin.roles；`admin` 在此基础上增加 admin.settings、admin.activity；custom 必须显式提供模块。显式 `APP_MODULES_ENABLED` 覆盖编译 Profile 默认，解析来源和优先级为 compiled-profile-default → modules.enabled → environment。 | verified | 2026-08-05 复核；R3/R5 只验证运行时运维，不改写本精确集合，除非新决策 | GOAL-003 C1 `attachments/r2-c1-profile-graph-evidence.md`、D-002、A-002、A-003；Root D-006/A-005。 |
| I-005 | required | Manifest 聚合的冲突规则、缓存、登录前加载与前端权限投影边界是什么？ | R2 聚合 API 方案冻结与 R3 联调 | R2 方案冻结前 | R2 C4 已定义并验证：Fragment 按 ModuleID 确定性排序；app/protocol/page/navigation/贡献冲突 fail closed；Profile 选择只投影启用模块；API 在登录前提供 `GET /.well-known/schema-ui/app-manifest.json` 和精确 ETag/304；Vite/Nginx 走 API，production image 删除静态 manifest，权限/可见性仍由前端运行时投影。 | verified | 2026-08-05 复核；R3 联调验证消费端，不扩大聚合冲突或登录前无秘密边界 | GOAL-003 C4 `attachments/r2-c4-aggregation-proxy-evidence.md`、C5 snapshot、D-004、A-002、A-003；Root D-006/A-005。 |
| I-006 | required | 静态 Manifest 和 Shell 特例的删除清单、迁移期限与回滚触发条件是什么？ | R3 试点门闩与 R6 旧路径移除 | R3 方案冻结前 | 盘点双轨入口与 Shell 特例，确定有期限的开发期兼容、告警和回滚条件。 | open | 不延期；R3 方案冻结前复核 | VP-003 信息门禁提示；待收集。 |
| I-007 | required | 本 Root 是否可读并遵守 VP-003 继承的 `I-PROTO-001 v0.1.3` 范围？与 R1 迁移模块清单是否一致？扩大范围的决策门槛是否明确？ | R1 契约/范围冻结与 R4 全量迁移范围 | R1 方案冻结前 | 核对 [VP-003 继承节](../../vision/plans/VP-003-modular-admin-architecture.md#继承的协议基线i-proto-001-v013) 与 Q2 覆盖表路径；对照 I-001 模块清单；扩大 domain/改 exclude 须新决策 + 覆盖表升版（D-002 已冻结「默认不扩大」约束）。 | verified | 2026-08-04 复核；范围扩张仍须新决策、覆盖表升版和验证 | [GOAL-002 C4](../GOAL-002-r1-contract-migration-baseline/attachments/r1-c4-protocol-matrix.md)；D-005；Grok A-004 independent + A-005 response。D-EXPR/D-VER、partial boundaries、D-UPLOAD exclude 与 version gate 均保留。 |

## 台账布局

本 Root 使用平铺 `01-decision/`、`02-execution/`、`03-audit/` 目录承载独立记录；索引文件只维护摘要与链接。当前没有共享资料引用。
