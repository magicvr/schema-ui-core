---
title: 决策 · 生产级可用 Admin 基架
status: active
created: 2026-08-01
updated: 2026-08-02
parent: null
version: 0.1.9
---

# 决策 · GOAL-001

## D-001 · 以新 delivery 工作区承接 VP-002

- **日期**：2026-08-01
- **决定**：建立 `workspace-002-production-admin-foundation` 与 Root `GOAL-001-production-admin-foundation`，将 VP-002 从 `planned` 激活为 `active`，并把本工作区设为其当前唯一 lead workspace；仓库 `primary_workspace` 保持 `workspace-001-mvp-admin-foundation`。
- **依据**：用户明确要求开启新工作区和 Root 承接 VP-002，并确认了工作区与 Root 命名。VP-001 与旧 Root 均已关闭/完成，新波次需要独立 canonical scope 和目标树。
- **边界**：不重开 VP-001 或旧 Root；不建立跨工作区 `parent`；只通过 Q2 路径引用已冻结历史基线。

### 未选方案

- **复用 workspace-001 并重开旧 Root**：会混合已关闭波次与新实施事实，破坏状态边界。
- **把新工作区设为 primary**：当前没有长期目的或仓库北极星换代，不应改写 Charter 的 `primary_workspace`。
- **把 VP-002 作为旧 Root 的子目标**：跨工作区 parent 被协议禁止，也无法提供独立交付树。

## D-002 · 采用五阶段串行纲领路线图

- **日期**：2026-08-01
- **决定**：Root 采用 `Renderer → 认证 → 持久化权限 → CRUD → 工程化与关门` 五个等权检查点，纲领阶段原则上串行，阶段内部可按依赖并行。
- **原因**：该顺序遵循 VP-002 的价值链，同时把原第三阶段拆成可独立验证的权限持久化、CRUD 和工程交付，避免一个阶段承载过多门禁。
- **执行约束**：只在当前阶段边界和 required 信息项就绪后创建具体子目标；`progress` 只按完成检查点数派生，不用于放行。

## D-003 · 先登记未知，再按阶段关闭

- **日期**：2026-08-01
- **决定**：将协议实施差量、认证机制、持久化模型、代表性 CRUD 实体及部署/fork 验收口径登记为阶段 required 信息项；操作日志保持 non-blocking。
- **原因**：这些未知不妨碍 Root 立项，但会改变对应阶段的方案或验收，必须在最晚门禁前以证据关闭或经用户书面接受有界 residual。

## D-004 · 采用 I-001 差量矩阵作为 R1 方案边界输入

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：
  1. 以附件 [I-001-implementation-gap-matrix.md](attachments/I-001-implementation-gap-matrix.md) 作为 `I-001` 的可核对答案，并将 `I-001` 标为 `verified`。
  2. **不改写** `I-PROTO-001 v0.1.3` 覆盖表；本波次仅继承其 include / include-partial / exclude。
  3. 将 R1 方案边界冻结为矩阵 §4：**In** = Schema 加载、默认 `RenderPage` 主路径、加载时结构校验、白名单节点/表单、`$context` reactions、代表性 Node 页、统一失败面、手写示例降级；**Reuse** = 现有 shell/fixture pin/演示 records API/静态 dev session（非生产身份）；**Out** = 真实认证、持久化 IAM、覆盖扩域、D-UPLOAD、`$deps` reactions、业务模块、fork/Docker 关门证据。
  4. 本回合**不**创建 R1 子目标、**不**修改产品代码；矩阵闭合只解除「R1 方案冻结前」的 `I-001` 信息门禁，**不**勾选 Root 路线图 R1、不宣称 Renderer 产品化完成。
- **理由**：用户确认 `/govern` 路径 A（仅文档关闭 I-001）。扫描显示库级/fixture 能力大量已有，但产品默认页面仍走 `EXAMPLE_PAGES`，`schemaUrl` 未驱动渲染——R1 主差量在主路径产品化，而非重新冻结协议。
- **关联信息项**：`I-001` → `verified`；证据路径见矩阵 frontmatter 与 §5。
- **后续**：可按矩阵候选拆分创建 R1 子目标并进入实施；I-002～I-005 仍分别阻断 R2～R5。

### 未选方案

- **跳过矩阵直接批量建 R1～R5 子目标**：违反 P-001/P-005；I-001 为 R1 方案冻结 required。
- **把现有手写示例与 fixture 通过直接标 R1 完成**：与 VP-002「Schema 驱动为默认主路径」成功标准冲突。
- **在本回合扩大 v0.1.3 或实现真实认证**：分别需要新覆盖决策与 R2 信息/方案，超出路径 A。

## D-005 · R1 拆为三个子目标（加载 / 主路径 / 代表性页）

- **日期**：2026-08-01
- **状态**：accepted
- **决定**：在 I-001 已 verified 且 D-004 方案边界生效后，于本工作区创建三个 **R1 阶段内** 子目标（均可 `active`，阶段内可并行准备）：
  1. `GOAL-002-r1-schema-load-validate` — Schema 加载、结构校验、统一错误面  
  2. `GOAL-003-r1-default-render-path` — 默认 `RenderPage` 主路径与 `EXAMPLE_PAGES` 降级（硬依赖 002）  
  3. `GOAL-004-r1-representative-node-pages` — 代表性列表/表单/组合 Node 页与回归证据（完整主路径证明依赖 002+003）  
- **理由**：用户明确要求按 D-004 创建「加载 + 主路径 + 代表性 Node 页」；拆分对齐矩阵 §4 候选，避免单目标混杂门禁。
- **边界**：不创建 R2～R5 子目标；不勾选 Root `progress` R1；不改产品代码（本决策仅立项）。
- **后续**：优先实施 GOAL-002；GOAL-003 切换默认分支须 002 可测；GOAL-004 资产可先行。

### 未选方案

- **合并为一个 R1 大子目标**：验收与并行困难，与矩阵建议拆分不一致。
- **只建加载、不做主路径/页面目标**：无法关闭 VP-002 阶段 1「默认 Schema 驱动」主张。

## D-006 · 认证方案的信息收集边界

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：R2 暂不冻结具体认证实现；将 `I-002` 置为 `collecting`，并以以下边界收集可供后续方案取舍的事实：
  1. 登录、登出、会话恢复、过期与撤销的生命周期和可验证行为；
  2. 凭据传递与存储边界，包括浏览器、API 与部署配置中的安全约束；
  3. 后端请求身份解析、中间件以及 `401` / `403` 的责任边界；
  4. 静态开发会话作为显式本地开发兜底的开关、隔离与禁止成为生产默认的约束；
  5. 认证方案与后续持久化身份模型的依赖边界，避免把 R3 的数据模型假设写成既定事实。
- **理由**：`I-002` 是 R2 方案冻结与实施的 required 信息门禁。现有 R1 证据只证明 Schema Renderer 产品化，不足以证明认证、凭据与请求级身份的方案已经就绪。
- **影响**：R1 通过不改变；R2 可进行有界信息收集，但不得据此冻结方案、创建认证实施范围或将任何认证机制表述为已选定、已验证。
- **后续**：收集当前前后端栈、部署约束与威胁边界的可核对事实，形成认证生命周期、请求身份和验收矩阵；再由用户确认具体方案取舍。

### 未选方案

- **立即选定 cookie、JWT 或第三方身份提供者**：当前缺少部署、凭据边界和现有栈约束的已验证输入，会把假设写成决策。
- **以 R1 的静态开发会话作为生产默认**：违反 VP-002 对真实认证和请求级身份的成功边界。
- **接受 `I-002` residual 后直接实施**：用户未作书面 residual 接受，且该未知直接影响 R2 的方案与安全边界。

## D-007 · R2 认证方案：短 JWT Access + Opaque Refresh + SQLite（C+B 混合）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. R2 认证采用 **短时效 JWT Access Token（`Authorization: Bearer`）+ Opaque Refresh Token** 混合方案（附件 §5 方案 C+B）。
  2. 会话/凭据存储引入 **SQLite**（先支持 sqlite）：刷新令牌哈希存储、可撤销；登录凭据落 SQLite。
  3. **接受引入 JWT 库**：Go 侧后端依赖从当前零第三方依赖扩展为含 JWT 库 + SQLite 驱动（具体库选型与版本在 R2 实施时定稿并留痕）。
- **理由**：用户裁决候选方案与开放前提。短 JWT access 提供无状态、可水平扩展的请求身份载体；opaque refresh 服务端可撤销并支持过期/登出语义，弥补纯 JWT 撤销难；SQLite 满足 MVP 轻量持久化，且 R3 身份模型可在此之上平滑扩展；Bearer 方案不耦合同源 cookie，部署形态灵活。
- **边界**：
  - 访问令牌为短时效 JWT；刷新令牌为不透明随机串、哈希存储于 SQLite，登出/刷新时撤销。
  - 静态开发会话（`StaticDevSession`）仅保留为显式本地 dev 兜底，生产默认不启用（验收 M9）。
  - 本决策只冻结**机制 / 存储 / 依赖**；具体参数（access/refresh TTL、env 配置键、前端令牌存储策略、CORS 与同源/跨源托管、SQLite 表结构、种子用户/凭据边界、密码哈希）在 R2 子目标方案/计划中定稿并留痕，**不在此静默假设**。
  - **与 R3 边界**：SQLite 用户/凭据在 R2 为满足真实登录的最小种子形态；用户—角色—菜单持久化与权限模型属 R3（`I-003`），本决策不实施 R3。
- **影响**：`I-002` → `verified`（方案已裁决；证据 = 本决策 + [I-002-auth-collection.md](attachments/I-002-auth-collection.md) §5/§6）。R2「方案冻结与实施」信息门禁解除，可进入 R2 子目标立项。
- **后续**：冻结 R2 方案边界 → 创建 R2 子目标（登录/登出/会话恢复/过期/撤销/请求身份中间件 + `401`/`403` + dev 兜底开关）→ 实施时定稿 TTL / env / 前端存储 / CORS / 表结构并留痕。

### 未选方案

- **方案 A · HttpOnly 会话 Cookie（同源）**：需同源托管 + CSRF 防护，与当前独立 SPA/API 进程形态耦合；本轮不选。
- **纯 Opaque Bearer（B 不含 JWT）**：撤销直接，但缺可扩展的 access 令牌形态；用户选择 JWT 承载请求身份。
- **纯签名 JWT（C 不含 opaque refresh）**：撤销难，与「短 access + 可撤销会话」目标冲突。
- **保持 Go 零依赖 / 进程内会话**：无法满足 SQLite 持久化与刷新令牌撤销；用户明确接受 JWT 库与 DB。

## D-008 · R3 持久化权限模型的信息收集边界

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 将 `I-003` 从 `open` 置为 `collecting`；先收集 R2 SQLite 占位模型、迁移/种子机制、用户—角色—菜单关系、后端授权点、前端导航投影与恢复验证事实。
  2. 以附件 [I-003-persistence-permission-collection.md](attachments/I-003-persistence-permission-collection.md) 作为本轮收集产物；其中方案 A/B/C、迁移计划与 M-R3-01～12 均为**候选**，不构成模型选定或方案冻结。
  3. 保持 R2 D-006 的兼容边界：规范化持久化必须继续输出 `account.User {id,name,roles}`，不得破坏 `/api/accounts/me`、`$context.user`、JWT subject 或 refresh-token 用户关系。
  4. 本轮不创建 R3 子目标、不修改产品代码、不勾选 R3 路线图检查点；`I-003` 在用户完成模型、菜单投影、迁移兼容期、读授权与恢复证据裁决前继续阻断 R3 方案冻结与实施。
- **理由**：R3 的 required 信息门禁已到达最晚阶段。当前实现只有 `users.roles` JSON 占位、启动时幂等建表和整体跳过式 admin seed；权限在后端固定角色检查与前端展示门控之间分裂，真实 manifest 无权限菜单项。直接立项会把数据关系、迁移和安全边界假设写成实施承诺。
- **未选方案**：
  - **直接沿用 `users.roles` JSON 开工**：无法交付规范化角色/菜单关系，也没有版本迁移与恢复证据。
  - **立即选定通用策略表达式模型**：复杂度超过最小权限闭环，且用户尚未裁决。
  - **把前端菜单隐藏当作授权完成**：不构成后端安全边界，不能替代 `401` / `403` 负向证据。
- **影响**：`I-003` → `collecting`；Root 仍为 `active / 2/5`；R3 仍未立项、未放行。
- **后续**：用户裁决候选方案与五个边界问题后，记录 R3 方案决策并判断 `I-003` 是否可转 `verified`；之后才冻结成功标准并决定子目标拆分。

## D-009 · R3 采用规范化 RBAC、features 菜单投影与两步迁移

- **日期**：2026-08-02
- **状态**：accepted
- **用户裁决**：确认推荐方案 B、`features` 菜单投影、两步迁移、读写权限边界及 R3 恢复证据口径。
- **决定**：
  1. **数据模型**采用方案 B：在现有 `users` / `refresh_tokens` 上增加 `roles`、`user_roles`、`permissions`、`role_permissions`、`menu_items`、`role_menu_items` 与 `schema_migrations`；API 授权依赖稳定 permission key，而不是硬编码角色名。
  2. **菜单投影**保留静态 App manifest 作为页面、路由、标签和导航结构来源；数据库保存 `page_ref`、显式且唯一的 `feature_key` 与角色 grant。`/api/accounts/me` 保持 `user.id/name/roles` 形状不变，通过现有 `features: Record<string,bool>` 投影菜单可见性；真实 manifest 使用 `$context.features.<feature_key>` 的 `visibleWhen`，前端隐藏不替代 API 授权。
  3. **迁移采用两步兼容期**：第一步建立版本表和规范化关系、回填 `users.roles` JSON，并在旧/新读结果间执行一致性核对；第二步切换规范化读写路径。旧 `users.roles` 在验证完成前保留，删除/停用必须是后续显式迁移，不与切换合并为不可逆动作。
  4. **读写权限边界**统一保护 records 路由：稳定 key 为 `records.read` 与 `records.write`；基础 viewer 角色仅获读权限，admin 获读写权限。匿名读写返回 `401`，已认证但缺相应 permission 返回 `403`。
  5. **R3 恢复证据**至少自动覆盖迁移前数据库副本、该副本的恢复、迁移后 `PRAGMA integrity_check`，以及用户身份、角色/权限、菜单 grant 与 refresh 关系的关键查询；完整生产备份运维流程仍归 R5。
  6. R3 建立一个端到端子目标 `GOAL-006-r3-persistent-rbac-menu`，用顺序检查点承载 migration → RBAC → seed → 读写授权 → `features` 菜单投影 → 恢复/回归，避免把强耦合闭环机械拆成多个目标。
- **信息门禁**：`I-003` → `verified`。证据为本决策及 [I-003-persistence-permission-collection.md](attachments/I-003-persistence-permission-collection.md) 的当前事实、候选比较与 M-R3-01～12；这只关闭 R3 方案冻结/立项目门禁，不构成实现或验收事实。
- **边界**：不引入通用策略表达式/IAM；不实施 R4 CRUD 扩域；不以菜单隐藏代替后端授权；不勾选 Root R3，Root 保持 `active / 2/5`。

### 未选方案

- **方案 A · 角色 + 菜单、API 继续硬编码角色**：安全语义与角色名耦合，不能形成稳定 permission 边界。
- **方案 C · 通用持久化策略表达式**：版本、校验、审计和运行时复杂度超过 R3 最小闭环。
- **一次迁移同时删除 `users.roles`**：压缩了核对与恢复窗口，失败时放大不可逆风险。
- **只保护写路由或只做前端菜单隐藏**：无法证明管理数据读边界与后端真实授权闭环。

## D-010 · R4 采用 records、SQLite 持久化与统一错误 envelope

- **日期**：2026-08-02
- **状态**：accepted
- **用户裁决**：采用 `records` 作为 R4 代表实体；R4 必须使用 SQLite 持久化并验证重启保持；沿用统一错误 envelope，精确 `code` 在 R4 子目标方案中冻结。
- **决定**：
  1. **代表实体**采用现有 `records`，继承 `id`、`name`、`status`、`owner`、`updatedAt` 候选字段与现有 list/search/detail/PATCH/DELETE、`records.read` / `records.write` 基线；R4 必须补齐 create 与 Schema 驱动的真实详情、编辑、删除闭环。
  2. **持久化**必须把 records 从进程内静态切片迁入 SQLite，并纳入版本迁移、可重复 seed、读写一致性与失败诊断；不得把现有静态数据路径保留为生产默认实现。
  3. **重启保持**是 R4 required 验收：至少自动证明 create/update/delete 后重启，list/detail 结果与持久化预期一致；同时覆盖 migration/seed 重复执行与关键失败路径。
  4. **错误契约**沿用现有统一 envelope（HTTP status + 稳定 `code` + message）；保留已实现且适用的 `UNAUTHENTICATED`、`FORBIDDEN`、列表参数错误、`RECORD_NOT_FOUND`、PATCH body/field 错误语义。
  5. **精确实施契约下沉**：create、字段校验、持久化/并发冲突及 Schema action 失败的 HTTP status、精确 `code`、请求/响应字段映射，必须在 R4 子目标方案中作为实施前 required 信息项登记、冻结并以测试验证；本 Root 决策不臆造未实现 code。
- **理由**：I-004 收集显示 `records` 是唯一已有 API、Schema fixture、loading/empty/error、权限与集中回归基线的业务候选，差量小于 users/RBAC 或新建业务域；SQLite 与重启保持直接对应 VP-002 的生产级与数据重启预期；沿用统一 envelope 可保持现有 API 错误边界，同时避免 Root 层过早冻结未设计的精确 code。
- **信息门禁**：`I-004` → `verified`。证据为本决策及 [I-004-schema-crud-collection.md](attachments/I-004-schema-crud-collection.md) 的现状、候选比较与 M-R4-01～11；只解除 Root 的 R4 方向冻结与子目标立项目门禁，不构成 R4 详细方案、实现或验收事实。
- **后续**：下一拍按 P-001 创建 R4 子目标；创建时必须登记上述实施前 required 信息项并建立可枚举路线图，在精确契约冻结前不得实施受影响的 create/persistence/Schema write 范围。
- **边界**：不扩展为 users/RBAC 管理后台或新业务域；不扩大 `I-PROTO-001 v0.1.3`；不把前端权限隐藏当后端授权；不勾选 Root R4，Root 保持 `active / 3/5`。

### 未选方案

- **采用 users/RBAC 作为代表实体**：虽已有持久化与重启证据，但没有完整 CRUD API 或 Schema CRUD 页面，会把身份域管理范围引入 R4。
- **新建另一业务实体**：会重复建立 records 已有的 API、Schema、权限和测试基线，缺少额外产品价值证据。
- **保留进程内 records 作为 R4 终态**：无法证明生产持久化和重启保持，与用户裁决及 VP-002 验收边界冲突。
- **在 Root 直接枚举未实现 error code**：会把尚未设计和验证的精确语义写成事实；应由 R4 子目标方案在实施前冻结。

## D-011 · 建立一个 R4 端到端子目标并下沉 required 信息门禁

- **日期**：2026-08-02
- **状态**：accepted
- **用户指令**：按 D-010 创建 R4 子目标并登记实施前 required 信息项。
- **决定**：
  1. 在本工作区创建 `GOAL-007-r4-schema-crud`，以六个顺序检查点承载精确契约、SQLite 结构/迁移/seed、持久化 CRUD API、Schema 读写主路径、交互/权限负向和重启回归。
  2. 不为每个未知机械创建信息目标；在 GOAL-007 内登记四项 required：`I-007-001`（API/错误）、`I-007-002`（SQLite/迁移/seed/并发）、`I-007-003`（Schema action/状态/权限）、`I-007-004`（重启与端到端验收协议）。
  3. 四项均须在各自首个受影响实施或验收动作前由证据关闭；当前只允许收集和方案冻结，不放行受影响产品代码。
- **理由**：D-010 已验证 Root 层方向，足以立项；精确契约仍是子目标实施门禁。单一端到端目标保留业务生命周期的完整验收边界，required 信息项则防止强耦合范围被隐式实现决定。
- **影响**：R4 子目标立项完成，但 Root R4 检查点仍未完成；Root 保持 `active / 3/5`，不产生实现或验收主张。

### 未选方案

- **把四项信息分别建成四个目标**：当前收集工作没有独立交付价值，且会把同一 CRUD 生命周期机械切碎。
- **立项后立即实施**：四项 required 尚未关闭，会越过 D-010 明示的实施前门禁。

## D-012 · R5 I-005/I-006 信息收集边界（工程化 / fork / 操作日志）

- **日期**：2026-08-02
- **状态**：accepted
- **决定**：
  1. 将 `I-005` 从 `open` 置为 `collecting`；以附件 [I-005-engineering-fork-collection.md](attachments/I-005-engineering-fork-collection.md)（v0.1.0）作为本轮收集产物，固定 R5 方案冻结所需的四类输入：**部署基线**（候选 A/B/C）、**15 分钟 fork 计时口径**、**可复现实验方法**（文档步骤 + smoke 清单 + 独立复现记录）、**容器/部署边界**（dev/prod、同源反代、DB volume）。候选均为选项，不构成选定或冻结。
  2. **复核 `I-006`**（non-blocking · R5 范围取舍）：现状无 operation/audit log；VP-002 将最小操作日志列为加分项非硬关门条件。附件 §5 给出方案甲（纳入 R5 可选加分 checkpoint）/ 方案乙（本波次排除）及建议（甲），最终范围取舍由用户裁决后写入 R5 方案。
  3. 本轮**不**创建 R5 子目标、不修改产品代码、不勾选 R5 检查点；`I-005` 在用户裁决附件 §6 四类边界前继续阻断 R5 方案冻结与立项；Root 保持 `4/5`。
- **理由**：`I-005` 为 R5 required 信息项且已到达最晚阶段（R5 方案冻结前）。扫描事实：全仓无 Dockerfile / docker-compose / 容器路径，无生产静态托管与 `/api` 反代文档，无 15 分钟 fork 计时口径与可复现实验方法；直接立项会把部署形态与验收口径假设写成实施承诺，且 VP-002 #7（Docker 一键启动、健康检查、dev/prod 区分）与 #6（≤15 分钟 fork）无法验收。操作日志与本波次是否纳入是范围取舍，须用户裁决。
- **影响**：`I-005` → `collecting`（required，仍阻断 R5 方案冻结/立项）；`I-006` → `collecting`（non-blocking，复核中待裁决）；Root `status: active`、派生进度 `4/5` 不变。
- **后续**：用户裁决附件 §6 四类边界与 I-006 取舍后，记录 R5 方案决策并判断 `I-005` 是否可置 `verified`；之后按 P-001 创建 R5 子目标（`GOAL-008-…`）并在其方案中登记实施前 required 信息项。

### 未选方案

- **直接按既有本地双进程命令立项 R5**：无容器/部署边界、无 15 分钟口径与复现方法，无法验收 VP-002 #6/#7；会把口径假设写成承诺。
- **立即选定 Docker Compose 并写进 R5 方案**：用户尚未裁决部署形态与计时口径，会跳过 P-004 用户裁决点。
- **把操作日志直接排除或直接纳入 R5**：范围取舍未做价值/成本裁决即静默决定，违反 P-004。

## D-013 · R5 方案边界冻结（部署基线 A + 建议口径 + I-006 方案甲）

- **日期**：2026-08-02
- **状态**：accepted
- **用户裁决**（P-004，书面确认附件 [I-005-engineering-fork-collection.md](attachments/I-005-engineering-fork-collection.md) §6 四类边界）：
  1. **部署基线 A**：文档化本地双进程为默认 + **可选 Docker Compose 一键启动**——api 多阶段构建镜像 + web 静态构建（nginx/Caddy 服务 `dist/` 并把 `/api` 反代到 api service），单源同源免 CORS，探针复用 `/healthz`，DB 走 volume。
  2. **15 分钟 fork 计时口径**：终点 = 登录成功 + 后台首页可交互（列表加载）；**不含依赖下载**（依赖就绪后文档首条命令起计时）；要求 **≥1 次独立复现**并记录（日期/ref/版本/耗时）。
  3. **可复现验收**：文档步骤 + smoke 清单（`/healthz` → login → `/me` → 代表页 → 种子可重复）；R5 落地 `scripts/smoke.sh`。
  4. **I-006 方案甲**：最小操作日志纳入 R5 为**可选加分 checkpoint**——若实施则记入 SQLite（`operation_log` 迁移 + repository），覆盖 records 写 + auth 关键事件；**不阻断** R5 核心验收；工时紧张可降级为留作后续。
- **决定**：以上四项作为 **R5 方案边界**冻结。`I-005` → **`verified`**（R5 方案冻结/立项目门禁解除，证据 = 本决策 + 附件 v0.2.0）；`I-006` → **`closed`**（非阻断，方案甲）。
- **理由**：用户按 P-004 书面裁决四类边界；A 匹配 VP-002 #6（≤15 分钟 fork）与 #7（Docker 一键启动、健康检查、dev/prod 区分）；口径可复现且不含依赖下载以弱化网络耦合；操作日志保持 VP-002 加分项定位并给出可降级路径。
- **影响**：Root R5 可立项；Root `status: active`、派生进度 `4/5` 不变（**不**据此勾选 R5，R5 检查点须待 R5 子目标完成证据）。
- **后续**：按 P-001 创建 R5 子目标 `GOAL-008-r5-engineering-fork`，登记实施前 required 信息项（精确 env/配置/部署基线、compose/镜像/volume/探针/反代契约、15 分钟计时复现协议 + smoke.sh 判据、operation_log 事件/存储契约[若实施]），未关闭前不得实施受影响范围或验收 R5。

### 未选方案

- **部署基线 B（仅文档双进程）**：会把 VP-002 #7 的 Docker 一键启动直接留在 R5 未交付。
- **部署基线 C（API 托管静态单进程）**：需改 api 静态服务与 SPA fallback，与「Schema Renderer 主路径不改写」目标交叉、成本中高。
- **更严格/更宽松计时口径**：含依赖下载与网络强耦合、复现成本高；不含首页可交互则对 VP-002 #6 可验证性最弱。
- **I-006 方案乙（本波次排除）**：放弃开箱可审计写路径的加分价值；用户选择保留为可选加分项。


