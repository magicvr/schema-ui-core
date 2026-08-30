---
doc_type: vision-plan
id: VP-003-modular-admin-architecture
title: 单主线模块化 Admin 架构
status: closed
lead_workspace: workspace-003-modular-admin-architecture
vision_ref: schema-ui-core-admin-foundation@0.3.0
closed_under_vision_ref: schema-ui-core-admin-foundation@0.3.0
created: 2026-08-04
updated: 2026-08-06
version: 0.2.0
parent: null
---

# VP-003 · 单主线模块化 Admin 架构

## 最终意图

把当前生产级 Admin 基架改造成**单一代码主线上的模块化单体**：薄内核提供稳定平台契约，组合根静态接入一方模块，Uber Fx 只负责装配与生命周期；启动时 Profile 选择已编译模块集合，后端聚合 Manifest、Schema、导航、权限、迁移和系统数据贡献，同一前端 build 可承接不同 fork 起点。

本 VP 表达完整终态，不是 Activity/Settings 试点的“妥协版本”。试点只是在迭代路线图中验证模块切口、失败语义和迁移方式；即使试点通过，也不得据此关闭本 VP。

架构权威见 [单主线模块化 Admin 架构](../../architecture/module-architecture.md)。

## 范围修订 · 2026-08-05

用户已裁决：`records` 是历史范例演示的无语义实体，当前产品不恢复其 CRUD、API、
种子、权限、菜单或专属前端面。R4 中的 `records/Schema CRUD` 统一解释为仍存在的
Schema-driven Admin 能力；`0003`/`0006` 迁移账本、历史 operation-log 事件和历史治理
证据继续保留，不得改写。该修订由工作区 3 的 GOAL-005 D-003、GOAL-006 D-003
承接，作为后续 C2-C5 的范围约束。

## 继承的协议基线（I-PROTO-001 v0.1.3）

VP-003 继承工作区 `workspace-001-mvp-admin-foundation` 中由 Root 决策 `D-009` 正式冻结的 `I-PROTO-001` 覆盖基线 `v0.1.3`。权威记录为 [I-PROTO-001 覆盖表](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)，其来源固定为 `schema-ui-docs v2.7.0` pinned commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`；对应 Root 决策见 [GOAL-001 `01-decision.md`](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/01-decision.md)。该基线是架构迁移的范围约束，不是 VP-003 的实现或验收证据。

| disposition | domain_id | VP-003 继承范围 |
|---|---|---|
| `include` | `D-NODE`, `D-EXPR`, `D-DATA`, `D-PERM`, `D-APP`, `D-VER`, `D-VAL` | 仅覆盖 v0.1.3 表中已列的语义、结构与验证入口 |
| `include-partial` | `D-COMP`, `D-ACT`, `D-TABLE`, `D-FORM` | 保留已冻结的白名单、非批量 action/request、基础表格交互与表单边界；不扩展为完整 registry 或批量语义 |
| `exclude` | `D-UPLOAD` | 上传 UI、端点及 `uploads` fixtures 整域不在本 VP 范围 |

任何新增 domain、扩大 `include-partial` 子集、改变 `D-UPLOAD` 排除或引入新的上游协议版本，都必须追加新的决策、递增覆盖表版本，并在受影响的 `/govern` 信息门禁前完成验证；不得静默改写 `v0.1.3`。

## 方向级退出判据

只有下列终态全部具备工作区 Q2 实施证据、验证结果与关门审计时，本 VP 才可提议 `closed`：

1. **单主线与 Profile 成立**
   MVP、完整 Admin 与 custom fork 起点由同一代码主线、同一模块契约和版本化 Profile 表达；不存在需要长期回灌的平行 MVP/Admin 代码线。Profile 解析、覆盖优先级、依赖闭包和 fail-closed 错误可验证。

2. **薄内核、组合根与模块能力契约完成**
   内核不导入业务模块；组合根静态汇集候选模块并通过 Uber Fx 管理依赖与生命周期；模块公共 API 不依赖 Fx 类型。现有一方 Admin 能力均迁入统一的模块描述与能力契约，不再依赖中央业务注册表接入。一方**标准 Admin 功能模块**必须具备核心六项贡献（HTTP、Schema、Authorization、Navigation、Manifest、Persistence）；Configuration / Lifecycle / Observability 及空依赖/配置/种子声明为按需，**「按需」不得覆盖核心六项**（权威见 [module-architecture.md](../../architecture/module-architecture.md) §2）。

3. **数据生命周期保持可升级、可恢复**
   全局迁移台账、顺序、checksum、事务、升级前快照与恢复边界保持成立；所有已编译一方模块迁移独立于启用状态执行，退役模块仍保留历史迁移 tombstone，未知已应用迁移 fail closed。fresh bootstrap 与 versioned system-data reconcile 分离，既有实例可补齐系统数据且不覆盖用户拥有字段。

4. **后端聚合运行时契约成为唯一生产路径**
   已启用模块的 Manifest、Schema、导航、权限和配置由后端确定性聚合、冲突校验并从 `/.well-known/schema-ui/app-manifest.json` 发布；端点登录前可读但不泄露秘密，Vite/Nginx 精确代理，协议/模块 API/前端 capability 不兼容时 fail closed，同一前端 build 在至少 `mvp` 与 `admin` Profile 下工作。生产静态 Manifest 兜底已移除。

5. **安全、横切能力与生命周期边界成立**
   认证和后端授权仍是最终权限边界；`operationlog` 始终记录关键写操作，`activity` 仅作为可选读取/UI 模块；Settings 不再依赖 Shell 私有通知。模块启动、就绪、停止、失败清理、健康诊断、日志与指标均有明确 `module_id` 语义（**指标属 Observability 按需能力**：当前基线不要求指标贡献契约，已交付范围为日志与健康诊断；若交付指标则须携带 `module_id`，权威见 [module-architecture.md](../../architecture/module-architecture.md) §2.2 与 workspace-003 Root D-011）。

6. **现有能力完成迁移且旧装配路径退出**
   当前用户、角色、仍存在的 Schema-driven Admin、Settings、Activity 等一方能力在新架构上保持既有行为和协议边界；历史 Records 不恢复为产品能力。旧中央路由/页面/导航注册、静态生产 Manifest、Shell 特例与已被替代的 host glue 被删除，而非无限期双轨兼容。

7. **可 fork、可运维、可回归**
   快速启动、Docker/生产代理、升级与恢复文档反映新架构；CI/本地矩阵覆盖模块契约、冲突失败、双 Profile、数据升级、认证授权、Activity 禁用场景和容器启动。新项目能通过选择 Profile、配置和模块贡献接入业务，不修改前端 Renderer/Shell 主路径。

## 迭代路线图

下列阶段是抵达终态的建议顺序，不缩减上述退出判据。具体工作区 Root 可在 `/govern` 中细化信息门禁、阶段检查点和子目标。

| 阶段 | 目的 | 阶段产物 / 门禁 | 与终态关系 |
|------|------|-----------------|------------|
| R1 | 契约与迁移基线冻结 | 盘点当前中央注册点、模块边界、迁移/seed 所有权、Profile 矩阵；冻结模块 API（含核心六项 / 按需能力口径）、capability 协商、迁移 tombstone、错误分类、现有兼容基线和回滚策略 | 只冻结实施边界，不算架构完成 |
| R2 | 内核与组合根基础 | 建立薄内核、框架无关模块契约、Fx 组合根、确定性图校验、全局迁移收集、Manifest 聚合骨架与 `/.well-known` 代理 | 提供可迁移平台，不关闭 VP |
| R3 | 有界试点 | 见下方 **R3 通过门闩**（继承固定历史评议输入 §4.5/§5）；通过仅允许进入 R4，不关闭本 VP | 验证切口与 Kernel 手术；**非**终态 |
| R4 | 全量一方模块迁移 | 将 users、roles、仍存在的 Schema-driven Admin 能力及其他现有 Admin 能力迁入统一能力契约；Records 仅保留历史迁移/审计兼容；清除模块对 Shell/中央注册表的特例依赖 | 达到功能覆盖面，但仍需退出旧路径与运维验收 |
| R5 | Profile、数据与运维收敛 | 完成 `mvp`/`admin`/custom 配置、fresh/reconcile、readyz/诊断、Vite/Nginx/Docker、升级恢复和 fork 文档 | 形成可发布候选，不自动等于 VP 关闭 |
| R6 | 旧路径移除与终态验收 | 删除双轨兼容和静态生产兜底；运行完整回归、双 Profile、升级/恢复、失败路径、容器/fork 验收与 close-out 审计 | 七条退出判据全部取证后方可提议关门 |

### R3 通过门闩（有界试点 · 继承固定历史评议输入）

R3 的目标是用 `operationlog`/`activity` 拆分与 `settings` 做手术刀，**证明 Kernel 切口正确**，不是「写出新模块就过关」。下列门闩全部满足后才可进入 R4；任一未满足则先加固 Kernel / 试点范围，**禁止盲目全量存量迁移**。

**A. 试点模块交付（5 项）**

1. `activity`（及拆分后的 `operationlog`）与 `settings` 按统一模块契约完整实现（标准 Admin 功能面满足核心六项；横切 `operationlog` 按架构说明豁免不适用 UI 项）。
2. **拆除**中心化 Register / 中央业务挂载中对试点模块的硬编码，改为模块自注册。
3. **将**试点 Schema 从全局 fixtures 迁入模块包，并由 Kernel 按模块加载。
4. **实现**后端 Manifest 聚合 API；试点页面/菜单经聚合暴露；前端至少完成对接联调（开发期双轨须有期限与告警，不得作为生产静默兜底）。
5. **排查并抽象** Host/Shell 对试点的硬编码特例（含 Settings 私有通知等），改为通用配置/事件贡献。

**B. 四个旧架构病灶（试点范围内必须切除）**

| 病灶 | R3 必须完成的动作 |
|------|-------------------|
| 中心化巨型业务 Register | 试点模块自注册；Register 中无试点硬编码挂载 |
| 全局 Schema fixtures | 试点 Schema 归属模块内；内核支持按模块加载 |
| 静态前端/生产 Manifest | 聚合 API 可用；试点片段注入；前端可消费聚合结果 |
| Host 层业务特例 | 试点相关特例删除或泛化为通用 hook |

**C. 验证项 V-1～V-4（必须全部通过）**

| 编号 | 验证项 | 验收标准 |
|------|--------|----------|
| V-1 | 模块注册与启动 | `enabled` 含试点模块时路由可访问、菜单出现 |
| V-2 | 模块禁用 | 从 `enabled` 移除并重启后，业务路由不可用、菜单消失；禁用≠删表 |
| V-3 | 前端零改动 | 同一前端 build 对接含/不含试点模块的后端，页面集随 Manifest 变化，前端代码无 diff |
| V-4 | Schema 贡献 | 试点列表/详情等标准页由后端 schema 驱动，Renderer 通用渲染 |

另须覆盖：依赖/冲突 fail-closed、配置变更事件（Settings）、同一 build 下至少 `mvp`/`admin` Profile 差异、既有实例升级路径上的系统数据 reconcile 不覆盖用户字段。

**D. 决策门（固定历史评议输入 §5.3）**

- 上述 A+B+C 全部通过 → 允许启动 R4 存量迁移计划。
- 卡在 Kernel 改造（尤其病灶 2/3 或 Manifest 聚合）→ 先加固 Kernel 再继续，**不得**用「模块代码已存在」放行 R4。
- R3 通过**仍不得**关闭本 VP；终态仍以七条退出判据与 R6 为准。

## 已接受的架构决策

| 主题 | 决策 |
|------|------|
| DI | 采用 Uber Fx；Google Wire 因已归档/不再维护且不直接解决运行时选择而不采用 |
| 选择模型 | 静态编译候选模块 + 启动时 Profile；不做运行时插件或热插拔 |
| API 边界 | 模块契约框架无关；一方标准 Admin 功能模块**核心六项必须**；Configuration/Lifecycle/Observability 按需；依赖显式；冲突与缺依赖 fail closed |
| 数据 | 全局迁移台账覆盖全部已编译一方模块；启用状态不决定迁移；bootstrap 与 reconcile 分离 |
| Activity | operationlog 始终启用；Activity 是可选查询/UI；Settings 同批验证通用 host hook |
| Observability | 指标属**按需能力**：当前基线不要求指标贡献契约，已交付范围为日志（带 `module_id`）与健康诊断（healthz/readyz 模块图门控）；若交付指标则须带 `module_id`（权威：module-architecture §2.2；workspace-003 Root D-011） |
| Manifest | 后端聚合、公共无秘密的 `/.well-known` 入口；生产不静默使用静态兜底 |
| 试点 | R3 必须完成 draft 病灶切除 + 5 交付 + V-1～V-4；失败则加固 Kernel，不盲目 R4；试点不是 VP 退出边界 |
| 演进 | 以单主线和 Profile 替代双线长期维护；历史 VP/工作区保持关闭事实，不伪造分支删除 |

## 信息门禁提示

VP 允许在 `planned` 状态携带实施未知，但未来激活/开区后至少应在对应门禁前登记并核验：

- 当前所有中央注册点、模块候选与跨模块依赖的完整清单；
- 既有 `0001` 起迁移链如何映射到全局模块所有权，以及升级/恢复夹具；
- Fx 和模块 API 的固定版本、Go 兼容范围及生命周期错误语义；
- `mvp`/`admin` Profile 的精确模块集合与配置覆盖顺序；
- Manifest 聚合的冲突规则、缓存、登录前加载与前端权限投影；
- 静态 Manifest 和 Shell 特例的删除清单、迁移期限与回滚触发条件。

这些信息项用于实施治理，不改变本 VP 的最终意图，也不能作为无限期保留双轨的理由。

## 工作区绑定

用户已于 2026-08-04 确认将本 VP 激活，并建立唯一 lead / delivery 工作区 [workspace-003-modular-admin-architecture](../../workspaces/workspace-003-modular-admin-architecture/workspace.md)。其 Root 为 [GOAL-001-modular-admin-architecture](../../workspaces/workspace-003-modular-admin-architecture/GOAL-001-modular-admin-architecture/00-meta.md)，`primary_plan` 与 `plan_refs` 均为本 VP。该绑定建立实现层治理范围，不构成 R1 或任何架构实现完成证据。

## 关门记录

仅在 `closed` 或 `abandoned` 时填写。

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-06 | **closed** | 七条方向级退出判据均以 lead 工作区 Q2 证据满足：① 单主线与 Profile（kernel `ResolveProfile` 优先级/依赖闭包 fail-closed，同一 `cmd/server` 双 Profile 运行，V-8/V-9）；② 薄内核、组合根与模块契约（kernel 零导入 modules、Fx 唯一组合根、核心六项强制、`RegisterContributions` 自注册，V-2）；③ 数据生命周期（全局台账/checksum/未知已应用迁移 fail-closed、升级前快照、编译全局迁移收集、`records_retire` tombstone、bootstrap/reconcile 分离，V-10/V-11）；④ 后端聚合运行时契约（确定性聚合、`/.well-known` 登录前无秘密、ETag/304、前端 `validateAppManifest` 精确匹配、生产无静态 Manifest 兜底、同一前端 build 双 Profile，V-4/V-6/V-7/V-12/V-14）；⑤ 安全、横切与生命周期（HS256 access + 旋转 refresh、operationlog 恒启用、activity 只读、Settings 通用事件、Start/Ready/Stop 带 `module_id`、healthz/readyz 门控；指标=按需已按 D-011/exit #5 澄清，V-5/V-8/V-9）；⑥ 现有能力迁移与旧路径退出（users/roles/settings/activity 为 kernel.Provider、中央业务 Register 仅剩 core、退役符号/静态 Manifest 无活动残留，V-13/V-14）；⑦ 可 fork、可运维、可回归（QUICKSTART/compose/双 Dockerfile/CI 矩阵/`scripts/smoke.sh`，V-15 系矩阵）。Root `GOAL-001` `done / 6/6`（A-018 self `pass` + A-019 Grok independent `pass` + A-020 `/govern` response；A-021 Grok independent 动态代码复审 `pass`，V-1～V-14，required 0；A-022 响应 recommended 全 `fixed`）；Root 03-audit **开放 required=0**；Vision Review **0 open required**（VRev-001～009；F-V014/F-V015 亦已 editorial `fixed`）。用户指令确认关门。 | [Root GOAL-001 00-meta](../../workspaces/workspace-003-modular-admin-architecture/GOAL-001-modular-admin-architecture/00-meta.md)；[goal-tree](../../workspaces/workspace-003-modular-admin-architecture/goal-tree.md)；[Root 03-audit](../../workspaces/workspace-003-modular-admin-architecture/GOAL-001-modular-admin-architecture/03-audit.md)；[A-021 动态复审](../../workspaces/workspace-003-modular-admin-architecture/GOAL-001-modular-admin-architecture/03-audit/A-021-vp003-apps-code-independent-reaudit.md)；[GOAL-013 00-meta](../../workspaces/workspace-003-modular-admin-architecture/GOAL-013-r6-old-path-removal/00-meta.md)；[GOAL-013 终态证据 C64-V01~V08](../../workspaces/workspace-003-modular-admin-architecture/GOAL-013-r6-old-path-removal/attachments/r6-c64-terminal-evidence.md)；[module-architecture.md](../../architecture/module-architecture.md) | 有界 closed（点名区/目标）：**R4-I004**（workspace-003 / GOAL-006 D-003 用户书面 `accepted-residual`：operationlog append 失败可能产生审计缺口，长期 duration/archive 未定义；scope=Users/Roles/Auth/Settings 写入与历史 events，owner=`magicvr`，review trigger=合规/运营 retention 要求、日志规模阈值、恢复演练缺口或 R5 数据生命周期决策；A-020 已字段级复核、未扩张，不视为 retention 已定义）。本 VP 非目标（运行时插件/`.so`/热插拔、I-PROTO-001 覆盖扩张、业务领域模块、微服务拆分、重写上游 Schema 语义、试点/能启动/文档完成充当终态证据）保持排除，不构成残余。 |

## Non-goals

- 不建设运行时插件市场、`.so` 加载、远程模块下载或运行中热启停。
- 不在本 VP 中扩张上述 `I-PROTO-001 v0.1.3` 的冻结协议范围。
- 不交付订单、钱包、类目、通知等具体业务领域模块。
- 不以微服务拆分替代模块化单体，也不在本项目内重写上游 Schema 语义。
- 不把试点、能启动、局部模块化或文档完成当作终态实现证据。

## 规划修订短史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-08-04 | `0.1.0` | 用户确认战略方向与全部工程建议；明确 VP 表达完整最终意图，Activity/Settings 等试点只进入迭代路线图，不构成妥协版退出边界。 |
| 2026-08-04 | `0.1.1` | `/vision` 响应 VRev-007：`F-V010` 核心六项必须 vs 按需能力口径；`F-V011` R3 继承 draft 病灶切除 + 5 交付 + V-1～V-4 门闩。未改 `planned`、未激活、未绑工作区。 |
| 2026-08-04 | `0.1.2` | `/vision` 响应 VRev-008：固定 `I-PROTO-001 v0.1.3` / `D-009` 的继承范围，并将现行 R3 对评议稿的引用改为固定历史评议输入。未改 `planned`、未激活、未绑工作区。 |
| 2026-08-04 | `0.1.3` | 用户确认激活并绑定唯一 lead / delivery 工作区 `workspace-003-modular-admin-architecture`；`/govern` 建立对应 Root。未将建区写成任何 R1-R6 的实现完成。 |
| 2026-08-05 | `0.1.4` | 用户裁决 `records` 范围修订：records 为历史范例演示的无语义实体，不恢复其 CRUD、API、种子、权限、菜单或专属前端面；`0003`/`0006` 迁移账本、历史 operation-log 事件与历史治理证据保留。承接 workspace-003 GOAL-005 D-003 与 GOAL-006 D-003。 |
| 2026-08-06 | `0.1.5` | `/vision` 响应 VRev-009（`pass` / `editorial`）：F-V014 `fixed`（补 0.1.4 短史行）、F-V015 `fixed`（exit #5 与已接受决策写明「指标 = Observability 按需，当前无指标贡献契约」）。未改 `vision_ref`、未改 `active` 状态。 |
| 2026-08-06 | `0.2.0` | 关门：七条方向级退出判据经 lead 工作区 Q2 证据满足（Root `GOAL-001` `done / 6/6` + A-018/A-019/A-020 关门链 + A-021/A-022 动态复审响应，Root 03-audit 开放 required=0；Vision Review 0 open required），用户指令确认关门 → `status` `active` → `closed`；关门记录 + roadmap/workspaces 同步；`closed_under_vision_ref = @0.2.0`。 |
