---
doc_type: vision-plan
id: VP-002-production-admin-foundation
title: 生产级可用 Admin 基架
status: closed
vision_ref: schema-ui-core-admin-foundation@0.2.0
closed_under_vision_ref: schema-ui-core-admin-foundation@0.1.0
lead_workspace: workspace-002-production-admin-foundation
created: 2026-08-01
updated: 2026-08-04
version: 0.3.1
parent: null
---

# VP-002 · 生产级可用 Admin 基架

## 意图

本 VP 继续以 `I-PROTO-001` 已冻结的协议子集为基础。阶段 1 的 Schema Renderer 必须完整支持该子集内的核心能力，包括 Node、基础 Form/Table、Reactions、Permissions、App Manifest 等。任何超出该子集的扩展，必须通过新的决策明确声明、版本化并验证。

本 VP 在 Charter `schema-ui-core-admin-foundation@0.1.0` 下完成并关闭，将协议验证型 MVP 基线升级为中小型业务项目可以直接 fork、启动并接入业务的生产级 Admin 基架。2026-08-04 仅将机读对齐链 re-align 到现行 `@0.2.0`；不重开本 VP，也不改写其历史退出边界。

本 VP 的最终判断标准不是治理记录是否闭环，而是：

> 业务开发者能否在不重写前端主路径的情况下，启动基架、完成登录、使用 Schema 驱动的 Admin 页面，并通过修改 Schema 新增业务页面与标准 CRUD 功能。

本 VP 继承现有 `schema-ui-docs v2.7.0` 固定协议边界，不重新定义上游协议语义。完整能力指本 VP 声明的 Admin 页面能力完整可用，不代表覆盖 `schema-ui-docs` 全部协议。

## 继承的协议基线（I-PROTO-001 v0.1.3）

VP-002 继承工作区 `workspace-001-mvp-admin-foundation` 中由 Root 决策 `D-009` 正式冻结的 `I-PROTO-001` 覆盖基线 `v0.1.3`。权威记录为 [I-PROTO-001 覆盖表](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/attachments/I-PROTO-001-coverage-draft.md)，其来源固定为 `schema-ui-docs v2.7.0` pinned commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`；对应 Root 决策见 [GOAL-001 `01-decision.md`](../../workspaces/workspace-001-mvp-admin-foundation/GOAL-001-mvp-admin-foundation/01-decision.md)。该基线是范围约束，不是 VP-002 的实现或验收证据。

继承边界必须保持明确：

| disposition | domain_id | VP-002 继承范围 |
|---|---|---|
| `include` | `D-NODE`, `D-EXPR`, `D-DATA`, `D-PERM`, `D-APP`, `D-VER`, `D-VAL` | 仅覆盖 v0.1.3 表中已列的语义、结构与验证入口 |
| `include-partial` | `D-COMP`, `D-ACT`, `D-TABLE`, `D-FORM` | 保留已冻结的白名单、非批量 action/request、基础表格交互与表单边界；不扩展为完整 registry 或批量语义 |
| `exclude` | `D-UPLOAD` | 上传 UI、端点及 `uploads` fixtures 整域不在本 VP 范围 |

任何新增 domain、扩大 `include-partial` 子集、改变 `D-UPLOAD` 排除或引入新的上游协议版本，都必须追加新的决策、递增覆盖表版本，并在受影响的 `/govern` 信息门禁前完成验证；不得静默改写 `v0.1.3`。

## 与 VP-001 / GOAL-001 的关系

现有 [VP-001-mvp-admin-foundation.md](VP-001-mvp-admin-foundation.md) 与 `GOAL-001-mvp-admin-foundation` 的 `closed / done` 状态保留不变。

| 维度 | VP-001 / GOAL-001 | VP-002 |
|---|---|---|
| 主要价值 | 协议验证型 MVP 基线 | 产品可用性与业务接入能力 |
| 协议边界 | 固定 `schema-ui-docs@v2.7.0`，冻结 MVP 子集 | 继续使用同一固定来源，在必要范围内增加兼容扩展 |
| 页面主路径 | 手写示例页面和有限协议实现 | Schema Renderer 消费 Node 树页面，手写 React 不再是主路径 |
| 身份能力 | 骨架级账号/权限与开发会话 | 真实登录、会话/Token、请求级身份与后端鉴权 |
| 数据能力 | 协议示例和验证路径 | 持久化用户、角色、菜单及 Schema 驱动 CRUD |
| 结论 | 协议验证基线已完成 | 在其上继续交付可直接使用的 Admin 产品能力 |

VP-002 不重开或否定 VP-001，也不把 VP-001 的 `done` 改写成“生产级完成”。

## 产品级成功标准

VP-002 只有在以下方向级标准全部满足时，才可提出关闭：

1. **完整可用的 Schema Renderer**  
   存在稳定的 Renderer 主路径，能够消费符合 `schema-ui-docs v2.7` 约定的 Node 树页面，并完成页面组合、布局、导航、常用交互和状态呈现。新增页面的主要方式是修改 Schema，而不是手写 React 页面。

2. **真实认证体系**  
   支持真实登录、登出、会话或 Token 生命周期、失效与错误处理。静态开发会话只能作为明确 opt-in 的本地开发兜底，不能成为生产默认身份来源。

3. **持久化身份与最小权限模型**  
   - 用户、角色、权限或菜单关系可持久化保存；
   - 至少通过种子数据提供可用的超级管理员和基础角色；
   - 后端必须按请求身份执行真实鉴权，而不是依赖进程内静态会话；
   - 初期可以不提供完整的用户/角色可视化管理后台，但接口或种子机制必须足以支撑业务使用；
   - 至少能区分超级管理员、基础角色和未授权请求。

4. **Schema 驱动的标准 CRUD 闭环**  
   列表、搜索、详情、编辑、删除五类标准生命周期能够通过 Schema 页面与后端真实接口完整跑通，并具有输入校验、空状态、加载状态、成功反馈和统一错误处理。

5. **可重复的种子数据**  
   提供默认超级管理员、基础角色和示例菜单；初始化过程可重复执行或明确处理已存在数据，不依赖手工插库。

6. **Fork 后可直接接业务**  
   按文档配置环境变量并启动后，开发者可以在短时间内进入系统、使用种子账号，并通过修改 Schema 新增业务页面，而无需修改前端 Renderer 主路径。“按文档完成环境配置并成功登录进入系统，目标时间 ≤15 分钟”作为可验证体验指标。

7. **基础工程化可交付**  
   提供清晰的环境变量说明、Docker 一键启动、健康检查或等价启动验证、统一错误响应与日志边界，以及本地开发和生产配置的明确区分。

方向级验收证据应落在绑定工作区的 Goal 与审计记录中，不在 VP 中维护 progress 百分比。

## 历史双线演进策略（关闭时语境）

以下内容保留 VP-002 在 `@0.1.0` 下规划和验收时的语境。它不再定义后续架构：Charter `@0.2.0` 已以单主线 + Profile 替代长期双线意图，未来方向由 [VP-003](VP-003-modular-admin-architecture.md) 与 [module-architecture.md](../../architecture/module-architecture.md) 约束。

### A. 最小可扩展 MVP 基线

继续保留当前协议冻结、可 fork 骨架、示例和验证路径，作为低成本起点：

- 保留固定上游版本和已声明的协议子集；
- 保留开发会话等低门槛本地开发能力，但明确其非生产属性；
- 不因生产能力增加而扩大为“完整协议支持”；
- 新扩展必须说明来源、兼容性和验证范围。

### B. 完整 Admin 能力线

以业务开发者直接使用为目标，按以下优先级逐步补齐：

`Schema Renderer → 真实认证 → 持久化与最小权限模型 → Schema 驱动 CRUD → Fork 与工程化交付`

两条线共享协议边界和基础组件，但验收标准分开记录，不能用 MVP 基线的完成状态替代完整 Admin 能力验收。

## 建议阶段

### 阶段 1：Schema Renderer 产品化

目标：

- 建立 Schema 加载、校验、Node 树解释和渲染主路径；
- 接入现有 Admin shell、路由和菜单；
- 覆盖 VP 声明范围内的常用页面节点、交互和状态；
- 明确无效 Schema、未知节点和运行时错误的处理方式。

验收建议：

- 至少一组代表性列表、详情、表单和组合页面由 Node 树直接渲染；
- 新增示例页面只修改 Schema 即可完成；
- 手写 React 页面退居次要或兼容路径，Schema 驱动成为新增业务页面的默认和推荐方式；
- 无效或不兼容 Schema 被确定性拒绝并给出统一错误。

### 阶段 2：真实认证与请求级身份

目标：

- 实现真实登录、登出、会话/Token、过期和撤销；
- 引入可替换的认证实现与安全配置；
- 建立后端请求级身份中间件；
- 将静态开发会话限制为显式开发模式。

验收建议：

- 登录成功、失败、登出、过期和撤销行为可验证；
- 重启服务后身份状态符合持久化策略；
- 未登录请求返回 `401`，无权限请求返回 `403`；
- 生产配置不会默认启用静态开发会话。

### 阶段 3：持久化权限、Schema CRUD 与交付体验

目标：

- 持久化用户、角色、权限/菜单和种子数据；
- 完成 Schema 驱动的列表、搜索、详情、编辑、删除闭环；
- 加入统一错误处理、输入校验和基本运行诊断；
- 提供环境变量、Docker 和 fork 快速启动文档。

验收建议：

- 使用种子超级管理员登录后，可在 Schema 页面完成完整 CRUD；
- 基础角色受到后端真实权限约束；
- 数据重启后仍符合持久化预期；
- Docker 一键启动成功；
- 按文档 fork 后可通过修改 Schema 新增页面，不修改 Renderer 主路径；
- 最小操作日志（记录关键写操作）作为加分项，不作为 VP 硬性关门条件。

## Non-goals

本 VP 明确不包含：

- 订单、钱包、通知等完整业务模块；
- `schema-ui-docs` 全量协议覆盖；
- 在本项目内重新定义或替代上游 Schema 语义；
- 复杂多租户、组织层级、跨租户隔离；
- 完整 IAM 产品，包括企业 SSO、SCIM、复杂策略编排等；
- 复杂工作流、审批流和领域事件平台；
- 以 Web 页面取代治理仓库中的 Markdown、Goal Tree 或审计台账；
- 将“能启动”或“有示例”表述为生产级完成。

## 治理与落盘边界

- 唯一现行 Charter 已 re-align 为 `schema-ui-core-admin-foundation@0.2.0`；本 VP 的历史关门版本由 `closed_under_vision_ref` 保留。
- VP-001 与 GOAL-001 保持历史关闭状态。
- 本 VP 于 2026-08-01 从 `planned` 激活为 `active`，并绑定 `workspace-002-production-admin-foundation` 为当前唯一 lead workspace。
- 新工作区角色为 `delivery`；仓库级 `primary_workspace` 仍为 `workspace-001-mvp-admin-foundation`，不改写 Charter 或历史交付树。
- 实现层由 `/govern` 维护 Root 路线图、阶段信息门禁、Goal 五件套与审计证据；VP 不维护 progress。
- 跨工作区只通过 Q2 路径引用冻结基线，不建立跨工作区 `parent`。

## 工作区绑定

| workspace | root_goal | role | 绑定日期 | 说明 |
|-----------|-----------|------|----------|------|
| `workspace-002-production-admin-foundation` | `GOAL-001-production-admin-foundation` | `delivery` / lead | 2026-08-01 | 用户确认新工作区与 Root 命名；以独立实现树承接本 VP |

## 关门记录

仅在 `closed` 或 `abandoned` 时填写。

| date | outcome | summary | evidence_links | residuals |
|------|---------|---------|----------------|-----------|
| 2026-08-04 | **closed** | 七条方向级产品成功标准均以 lead 工作区 Q2 证据满足：① Schema Renderer 主路径（`GOAL-002/003/004` done：加载/校验/统一错误面 + 默认渲染路径 + 代表性 Node 页面）；② 真实认证（`GOAL-005` done：登录/登出/刷新/撤销 + 请求级身份中间件 + 401/403 断言）；③ 持久化身份与最小权限（`GOAL-006` done：迁移链 + 增量幂等种子 + permission key 后端门禁 + `me.features` 投影）；④ Schema 驱动 CRUD 闭环（`GOAL-007` done：`list-edit-lifecycle` + SQLite 持久化 + L1/L2 重启证据）；⑤ 可重复种子数据（`GOAL-006` S3 + 复现协议）；⑥ fork 后直接接业务（`GOAL-008` done：QUICKSTART + 无编译缓存复现 64.833s ≤ 900s + `smoke.sh` + CI）；⑦ 基础工程化（`GOAL-008` done：env 清单/健康检查/Docker Compose/dev-prod 区分）。后续加固目标 `GOAL-009/010/011/012/013` 全部 done（A-002/A-005 required 全 fixed、语义化 users/roles 双实体接入、Shell 导航 fixture 洁净、Settings/Activity 品牌与操作日志只读面）。Root `GOAL-001` `done / 5/5`（A-004 关门 → A-005 回退 → `GOAL-012` 闭合 F-001 → A-007 self close-out `pass` 再关门）；Root 03-audit **开放 required=0**（A-002/A-005 全部 `fixed`，A-006 `pass`，A-007 self `pass`）。Vision Review **0 open required**（VRev-001～004）。用户指令确认关门。 | [Root GOAL-001 00-meta](../../workspaces/workspace-002-production-admin-foundation/GOAL-001-production-admin-foundation/00-meta.md)；[goal-tree](../../workspaces/workspace-002-production-admin-foundation/goal-tree.md)；[Root 03-audit](../../workspaces/workspace-002-production-admin-foundation/GOAL-001-production-admin-foundation/03-audit.md)；[GOAL-005 00-meta](../../workspaces/workspace-002-production-admin-foundation/GOAL-005-r2-auth-session/00-meta.md)；[GOAL-006 00-meta](../../workspaces/workspace-002-production-admin-foundation/GOAL-006-r3-persistent-rbac-menu/00-meta.md)；[GOAL-007 00-meta](../../workspaces/workspace-002-production-admin-foundation/GOAL-007-r4-schema-crud/00-meta.md)；[GOAL-008 00-meta](../../workspaces/workspace-002-production-admin-foundation/GOAL-008-r5-engineering-fork/00-meta.md)；[GOAL-011 00-meta](../../workspaces/workspace-002-production-admin-foundation/GOAL-011-s4-semantic-admin-resources/00-meta.md) | 有界 closed（非阻断，点名区/目标）：vision 层 `F-V003`（recommended open——双线分支维护契约，按 VRev-003/004 约定推迟至后续双线 VP 建立前落盘）；`GOAL-011` `F-006`（recommended open / non-blocking）；Root A-006 `R-005` residual-by-design（短 access TTL，已 handled）；本 VP 非目标（业务模块、全量协议覆盖、IAM/SSO/SCIM、多租户、复杂工作流）保持排除，不构成残余。 |

## 规划修订短史

| 日期 | 版本 | 变更 |
|------|------|------|
| 2026-08-01 | `0.2.0` | 经用户确认完成结构选型：VP 从 `planned` 激活为 `active`，绑定新 delivery 工作区与 Root；不改变 Charter primary workspace。 |
| 2026-08-04 | `0.3.0` | `/vision`：七条方向级产品成功标准经 lead 工作区 Q2 证据满足（Root `GOAL-001` `done / 5/5` + `GOAL-002`～`GOAL-013` 全部 `done`，Root 03-audit 开放 required=0，A-007 self close-out `pass`；Vision Review 0 open required），用户指令关门 → `status` `active` → `closed`；关门记录 + roadmap/workspaces/README 同步。 |
| 2026-08-04 | `0.3.1` | `/vision` strategic re-align：机读 `vision_ref` 更新为 Charter `@0.2.0`，以 `closed_under_vision_ref` 保留 `@0.1.0` 关门语境；历史双线章节不再约束未来，后续由 VP-003 的单主线模块化终态承接。未重开 VP、未改关门证据。 |
