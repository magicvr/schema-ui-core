---
title: I-003 · 持久化身份、角色、菜单与最小权限模型信息收集
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-001-production-admin-foundation
version: 0.2.0
related_info: I-003
related_decision: D-009
---

# I-003 · R3 持久化权限模型信息收集与裁决（D-008 / D-009）

> **性质**：回答「数据存储、迁移、种子和用户—角色—菜单关系的最小模型是什么」所需的当前实现事实、候选比较、用户裁决与验证矩阵。
> **裁决状态**：用户于 2026-08-02 确认方案 B、`features` 菜单投影、两步迁移、读写权限边界与恢复证据口径；Root D-009 已将 `I-003` 置为 `verified`。
> **不是**：任何数据表、迁移、授权或菜单投影已经实现；`I-003` 的关闭不勾选 Root R3。
> **扫描日期**：2026-08-02（仓库只读对照；无产品代码变更）。工作区 `shared_materials_catalog: none`，本附件未使用共享资料作为事实或关闭依据。

## 0. 总览结论

| 维度 | 当前事实 | R3 必须解决的差量 |
|------|----------|-------------------|
| 身份持久化 | SQLite 已保存 `users` 与 `refresh_tokens`；登录、刷新、登出已使用真实存储 | 保留 R2 身份/会话行为，避免迁移破坏现有 token 生命周期 |
| 角色 | `users.roles` 是 JSON 数组占位；`account.User {id,name,roles}` 已是 `/api/accounts/me` 与 `$context.user` 契约 | 将角色规范化，同时保持对外 `roles: string[]` 形状不变 |
| 权限 | API 仅在 records 写路由用固定 `admin` 角色表达式做真实 `403` 授权 | 明确稳定 permission key、角色授权关系与后端强制执行点 |
| 菜单 | 真实 App manifest 是静态 JSON；当前菜单没有持久化表，真实 manifest 也没有权限门控项 | 明确菜单持久化与静态 manifest 的职责边界，并让权限结果进入真实导航路径 |
| 迁移 | 启动时一次执行 `CREATE TABLE IF NOT EXISTS`；没有 schema version 或顺序迁移 | 为已有 R2 数据库提供可重复、可失败回滚/诊断的版本迁移路径 |
| 种子 | 任意用户存在即跳过 admin seed；不是按稳定实体逐项幂等 | 定义 admin、基础角色、权限与菜单的可重复增量种子 |
| 恢复验证 | 有 store/auth 生命周期单测；没有历史库升级、约束、备份恢复或 schema 完整性测试 | 建立迁移前后、重启、恢复、孤儿关系和失败路径验证 |

## 1. 当前 SQLite、迁移与种子事实

| 事实 | 证据 |
|------|------|
| `Store` 使用 `modernc.org/sqlite` 和单一 `*sql.DB`；连接上限为 1 | `apps/api/internal/store/store.go` `Store` / `Open` |
| schema 只有 `users`、`refresh_tokens` 与 refresh-token 用户索引 | `apps/api/internal/store/store.go` `const schema` |
| `users.roles` 为 `TEXT NOT NULL`，代码注释明确为 `JSON array; R3 normalizes` | `apps/api/internal/store/store.go` `users` DDL |
| refresh token 只保存 SHA-256 hash，带用户外键、过期时间和撤销时间 | `apps/api/internal/store/store.go` `refresh_tokens` DDL；`apps/api/internal/auth/auth.go` `issue` / `Refresh` / `Logout` |
| `migrate()` 仅执行整段 schema；未发现版本表、迁移目录、迁移编号或旧数据变换 | `apps/api/internal/store/store.go` `migrate()`；仓库路径扫描 |
| `seedAdmin()` 先 `COUNT(*) FROM users`；只要已有任意用户就整体跳过，否则写入固定 `user-admin` 与 `roles=["admin","editor"]` | `apps/api/internal/store/store.go` `seedAdmin()` |
| 服务启动总是 `store.Open(DB_PATH, "admin", seedHash, true)`；开发缺省密码为 `admin`，非开发缺失种子密码会 fail closed | `apps/api/cmd/server/main.go` `run()` / `resolveSeedHash()`；`apps/api/internal/config/config.go` |
| store 单测覆盖种子幂等、用户查找与 refresh 生命周期，但未覆盖历史 schema 升级、外键启用、备份/恢复或迁移失败 | `apps/api/internal/store/store_test.go` |

## 2. 当前身份、权限与菜单链路

### 2.1 必须保持的 R2 契约

- R2 D-006 已冻结 `account.User {id,name,roles}` 为 `$context.user` 与 `/api/accounts/me` 的契约形状；R3 只能替换其持久化来源，不能破坏响应形状。
- `users.roles` JSON 是明确的 R3 占位。规范化后，store/auth 必须通过关联查询重建相同的 `roles: string[]` 快照。
- `refresh_tokens.user_id -> users.id` 与现有登录、轮换、撤销语义必须在迁移后继续成立。

证据：[GOAL-005 D-006](../../GOAL-005-r2-auth-session/01-decision.md)；`apps/api/internal/store/store.go`；`apps/api/internal/auth/auth.go`。

### 2.2 后端安全边界

| 路径 | 当前行为 | 边界判断 |
|------|----------|----------|
| Bearer middleware | access 缺失、无效、过期或 subject 不存在 → `401 UNAUTHENTICATED` | 已是真实请求身份边界，可复用 |
| records GET | list/detail 当前不经过认证中间件 | R3 是否保护读操作须在方案中明确，不得静默扩大 |
| records PATCH/DELETE | middleware + `writeGate`；无身份 `401`，无 `admin` 角色 `403 FORBIDDEN` | 已有最小正反授权证据，但授权语义仍硬编码为角色字符串 |
| 权限求值 | `account.Allow` 对非法表达式或未声明 `$context` 路径 fail closed | 可复用表达式行为；不能把前端隐藏当后端授权 |

证据：`apps/api/internal/auth/auth.go` `Middleware`；`apps/api/internal/handler/records.go` `routes` / `writeGate`；`apps/api/internal/account/permission.go` `Allow`。

### 2.3 前端导航与菜单现状

- Web 从静态 `apps/web/public/.well-known/schema-ui/app-manifest.json` 加载页面与导航结构。
- `projectNavigation()` 会依据 manifest 的 `permissions.view` 与 `visibleWhen.when` 过滤项目并剪枝空组；这是展示门控，不是 API 安全边界。
- 导航单测已有 admin 可见、viewer 隐藏与空组剪枝；真实 manifest 当前没有任何 permission-gated navigation item，因此没有真实菜单端到端权限证据。
- `/api/accounts/me` 提供的 `user` 与 `features` 会进入 `NavigationContext`；D-009 已选择复用现有 `features` 通道，不新增独立导航 entitlement API。

证据：`apps/web/src/main.tsx`；`apps/web/src/protocol/app-manifest.ts` `isNavigationItemVisible`；`apps/web/src/app/navigation.ts` `projectItems`；`apps/web/src/app/navigation.test.ts`；真实 `app-manifest.json`。

## 3. 候选最小模型（D-009 采用方案 B）

### 方案 A · 仅规范化角色 + 角色菜单授权

候选表：`roles`、`user_roles`、`menu_items`、`role_menu_items`；API 写授权继续直接判断角色字符串。

- **优点**：表少；可最小替换 `users.roles` JSON；容易保持 `/me` 契约。
- **风险**：后端授权仍与角色名耦合；菜单授权与 API 授权可能漂移；缺少稳定 permission key，不利于后续实体扩展。
- **适用判断**：只在 R3 明确接受“最小角色模型、权限仍为角色检查”时可选。

### 方案 B · 规范化 RBAC + 稳定 permission key + 菜单授权（推荐）

候选表：

```text
users                 # 保留 R2 身份与凭据字段，迁移后移除/停用 roles JSON
roles                 # id, key, name, system, timestamps
user_roles            # user_id, role_id
permissions           # id, key, description
role_permissions      # role_id, permission_id
menu_items            # id, page_ref, feature_key, sort_order, enabled
role_menu_items       # role_id, menu_item_id
schema_migrations     # version, applied_at, checksum
refresh_tokens        # 保持现有会话关系
```

- **优点**：后端路由以稳定 permission key 强制授权；角色名称不再直接充当全部安全策略；菜单授权与 API 授权可分别验证；满足用户、角色、菜单和最小权限关系的持久化边界。
- **风险**：需要定义 manifest `pageRef` 与 `menu_items` 的唯一映射，以及 entitlement 如何投影到 Web；比方案 A 多两组关联表。
- **推荐理由**：这是在不引入完整 IAM 的前提下，最小化“前端菜单隐藏”和“后端真实授权”漂移的方案。
- **已选投影（D-009）**：静态 manifest 继续负责页面、标签、路由和协议结构；数据库保存角色可见的 `page_ref`、显式唯一且兼容 `$context.features.<key>` 的 `feature_key` 及 grants。`/api/accounts/me` 保持 `user` 形状不变，通过既有 `features: Record<string,bool>` 提供可见性输入；真实 manifest 用 `visibleWhen` 消费该投影。

### 方案 C · 持久化通用权限表达式 / 策略

在角色、菜单或路由上存储 `$context...` 表达式或通用 policy 文本，运行时统一求值。

- **优点**：表达能力强，可直接复用现有表达式引擎。
- **风险**：把协议表达式、数据库策略和路由安全边界耦合；需要策略版本、校验、变更审计与更复杂的 fail-closed 迁移，超出“最小权限模型”。
- **建议**：R3 不选；除非用户明确要求策略化并接受额外复杂度。

## 4. D-009 冻结的迁移与增量种子边界

> 本节已作为 D-009 的方案边界输入，但仍不表示已经实现；精确 DDL、迁移版本与约束在 GOAL-006 内定稿。

1. 建立 `schema_migrations(version, applied_at, checksum)`；每个迁移在事务内按版本顺序执行，未知/校验不匹配版本 fail closed。
2. 将当前 R2 schema 视为已存在基线：创建规范化表后，解析每个 `users.roles` JSON；为角色去重建表并写入 `user_roles`。
3. 迁移期间保持 `users.id` 与 `refresh_tokens.user_id` 不变；迁移后通过关联查询重建 `account.User.Roles`。
4. 对空、非法或非数组 roles JSON 不静默丢弃：迁移失败并保留原数据库；提供可定位错误。
5. 种子改为按稳定 key 逐项 upsert/ensure：`admin`、基础角色、最小 permission keys、代表性 menu grants；已有业务用户不得阻止 bootstrap 关系修复。
6. 在读路径切换和验证通过前保留 roles JSON；删除或停用旧列应作为后续显式迁移，不在同一步不可逆完成。
7. SQLite 外键显式启用并验证；关联删除规则、唯一约束和种子升级语义必须在 R3 方案中定稿。

## 5. R3 验证矩阵（D-009 采用的最低证据）

| ID | 验证项 | pass 判定 |
|----|--------|-----------|
| `M-R3-01` | 空库初始化 | admin、基础角色、permission/menu grants 与 refresh schema 一次建立；可登录 |
| `M-R3-02` | 种子重复执行 | 连续启动/重复 seed 不产生重复关系，不覆盖用户修改的非种子字段 |
| `M-R3-03` | 旧 R2 数据迁移 | `roles=["admin","editor"]` 正确转为规范化关系；用户 id、密码与 refresh 关系不变 |
| `M-R3-04` | 非法旧数据 fail closed | 非法 roles JSON、重复 key 或孤儿关系使迁移可诊断失败，不部分提交 |
| `M-R3-05` | schema 版本幂等 | 已应用迁移不重复执行；未知/校验不符版本拒绝启动 |
| `M-R3-06` | 身份契约兼容 | `/api/accounts/me` 仍返回相同 `user.id/name/roles` 形状；前端登录与恢复不回归 |
| `M-R3-07` | 后端正反授权 | admin/具备 permission 的角色成功；基础角色 `403`；匿名 `401`；前端隐藏不替代 API 拒绝 |
| `M-R3-08` | 真实菜单投影 | 真实 manifest/运行路径至少一项受持久化 grant 控制；admin 可见，基础角色隐藏，空组剪枝 |
| `M-R3-09` | 重启持久化 | 用户、角色、权限、菜单关系和 refresh 状态在服务重启后保持 |
| `M-R3-10` | 约束与删除语义 | FK/unique/cascade-or-restrict 行为有正反测试；不存在静默孤儿关系 |
| `M-R3-11` | 恢复验证 | 迁移前副本可恢复；迁移后 SQLite `integrity_check` 与关键计数/身份查询通过 |
| `M-R3-12` | 回归 | store/auth/handler/Web navigation/Renderer 既有自动化及 build/vet 继续通过 |

## 6. 用户裁决与实施边界（2026-08-02）

1. **模型主路径**：采用方案 B；稳定 permission key 与菜单 grants 都持久化，不以角色名直接替代 permission。
2. **菜单投影**：采用现有 `features` 通道；数据库 `page_ref`/`feature_key` grants 投影为布尔值，静态 manifest 保持结构权威。
3. **两步迁移**：先建表、回填并双读核对，再切规范化读写；旧 `users.roles` 的删除/停用留给后续显式迁移。
4. **读写权限**：records 读写均要求认证和 permission；`records.read` 与 `records.write` 分离，viewer 只读，admin 读写；匿名 `401`，缺权限 `403`。
5. **恢复证据**：自动验证迁移前副本可恢复、迁移后 `PRAGMA integrity_check` 与身份/授权/菜单/refresh 关键查询；完整生产备份运维流程留 R5。
6. **目标结构**：创建一个端到端 `GOAL-006-r3-persistent-rbac-menu`，内部以顺序检查点承载强耦合闭环。

以上裁决由 Root D-009 固化，`I-003` → `verified`。精确 DDL、迁移编号、约束和首个真实菜单 key 属 GOAL-006 的实施前信息项；本附件不宣称实现已发生。

## 7. 证据索引

- Root/R2 边界：`GOAL-001-production-admin-foundation/00-meta.md` I-003；`GOAL-005-r2-auth-session/01-decision.md` D-003 / D-006
- 后端存储：`apps/api/internal/store/store.go`、`store_test.go`
- 认证：`apps/api/internal/auth/auth.go`、`auth_test.go`、`apps/api/cmd/server/main.go`、`apps/api/internal/config/config.go`
- 后端授权：`apps/api/internal/account/permission.go`、`permission_test.go`、`apps/api/internal/handler/records.go`、`records_test.go`
- 前端身份/导航：`apps/web/src/main.tsx`、`apps/web/src/account/context.ts`、`apps/web/src/protocol/app-manifest.ts`、`apps/web/src/app/navigation.ts`、`navigation.test.ts`
- 真实 manifest：`apps/web/public/.well-known/schema-ui/app-manifest.json`
- API 当前边界：`apps/api/README.md`
