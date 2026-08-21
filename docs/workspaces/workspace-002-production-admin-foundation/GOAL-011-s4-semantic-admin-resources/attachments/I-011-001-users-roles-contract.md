---
title: I-011-001 · users/roles 语义资源领域契约与最小 IAM 边界
status: active
doc_type: contract
created: 2026-08-03
updated: 2026-08-03
parent: GOAL-011-s4-semantic-admin-resources
version: 0.2.0
related_info: I-011-001
related_decision: D-002, D-003
supersedes: null
---

# I-011-001 · users/roles 领域契约（冻结）

> **性质**：回答「users/roles 的精确资源契约与最小 IAM 边界是什么」。把 GOAL-011 S1 的领域与安全决策固化为可实施、可验收的版本化契约。冻结后 `I-011-001` 由 GOAL-011 **D-002** 置为 `verified`，解除 S1 方案冻结与 S2 实施门禁。
> **v0.2.0（2026-08-03 · A-002 响应，GOAL-011 D-003）**：修订 §7（`ResourceEntity` Create/Update/Delete 增传 `account.User` actor 通道，使 SELF_OPERATION 可诚实实现；DomainError 检查先于 ErrNotFound/INTERNAL）、§2.3（API 路径禁止复用 `linkUserRole`→`ensureRole` 的隐式建角色）、§3（冻结 roles 公开响应形状）。`I-011-001` 维持 `verified`（响应修订不改变冻结结论）。
> **依据**：现有 `apps/api/internal/{auth,account,store,handler}` 静态核对（users/roles/user_roles/permissions/menu/operation_log 表结构、`Resource` 通用工厂、`account.User` 身份隔离、`userWithRoles` 双写）；[I-010-001 通用资源契约 v0.2.1](../../GOAL-010-a002-schema-adapter/attachments/I-010-001-schema-resource-contract.md)（本契约的通用基础；§4/§5 在此做版本化扩展）；GOAL-006 I-006-001（RBAC 规范化）；GOAL-007 I-007-001（records 兼容基线）；GOAL-011 D-001（users + roles 立项）。
> **不是**：S2～S4 的实施成品（store 方法、handler、fixture、回归证据属实施）；完整 IAM 产品；SSO/SCIM；多租户。

## 1. 目的与范围

- 目标：users 替换 records 成为默认语义代表实体、roles 作为第二语义资源；二者在**通用资源工厂之上**完成 Schema 驱动标准 CRUD 闭环（GOAL-011 S2），前端仅改 Schema 接入（S4）。
- 范围：users/roles 的公开字段、敏感字段隔离、字段校验、角色分配、self/最后管理员保护、system role 与 grant 约束、权限键/菜单/操作日志、错误码、对通用工厂的最小契约扩展。
- 保持冻结（不重开）：I-010-001 的统一 list envelope `{items,total,page,pageSize}`、`{error,message}` 写 envelope 形状、`dataSource` 单斜杠同源规则、`rowKey` 不变量、4 KiB body 上限、`{ID}_NOT_FOUND` 404 形态、records 零 API 变更（I-007-001）历史事实。
- **通用工厂最小扩展（本契约 §7 冻结，S2 落地）**：`JSONFields`（任意 JSON 值字段透传）与 `DomainError`（实体返回带 status/code 的领域错误，工厂逐字映射）。此扩展不改 I-010-001 已有字段的语义，仅追加能力。

## 2. users 资源契约

| 项 | 值 |
|----|----|
| 资源 id | `users` |
| 端点 | `/api/users`（list/create/detail/update/delete 五路由） |
| 权限键 | `users.read` / `users.write` |
| NOT_FOUND | `USER_NOT_FOUND`（404） |
| list | `sortFields = [username, name, updatedAt]`；`qSearch = true`（username/name 子串，大小写不敏感） |
| create 字段 | `username`(必填非空 string)、`name`(必填非空 string)、`password`(必填非空 string，**仅写**) |
| JSON 字段 | `roles`（可选；JSON 字符串数组；缺省 = `[]`；每项必须是已存在角色 key） |
| patch 字段 | `name`、`password`（可选）、`roles`（可选）；**`username` 不可变** |
| id | 服务端生成 `usr-` + 16 位小写 hex（8 字节 crypto/rand） |

### 2.1 公开字段（响应形状）

list/detail/create/update 的响应行固定为：

```json
{
  "id": "usr-…",
  "username": "alice",
  "name": "Alice",
  "roles": ["admin"],
  "createdAt": "2026-08-03T00:00:00.000Z",
  "updatedAt": "2026-08-03T00:00:00.000Z"
}
```

- `createdAt`/`updatedAt`：RFC3339 固定 3 位毫秒形状（与 records `updatedAt` D-004 一致）；存储为 unix 秒，映射时统一毫秒格式。
- `roles`：规范化 role key 数组（以 `user_roles` 关系为准，GOAL-006 阶段 B 双写集合一致）。

### 2.2 敏感字段隔离（硬性不变量）

- **`password_hash`（及任何内部列）不得出现在任何响应**（list/detail/create/update/error）。这是 S2「敏感字段隔离」检查点的强制验收项，负向路径必有断言。
- `password` 是**仅写**输入：create/patch 接收明文，服务端 bcrypt 哈希后写入 `password_hash`；响应永不含 `password` 或 `password_hash`。
- 操作日志 detail 只含 `username`，**永不**含密码或令牌（I-008-003 §3 边界延续）。

### 2.3 角色分配

- create/patch 的 `roles` 数组只接受**已存在角色 key**（校验 `roles` 表）；未注册 key → `INVALID_ROLE_REF`(400)。users 写操作**不**隐式创建角色（角色由 roles 资源管理）。
- **实现约束（v0.2.0 · A-002 F-004）**：既有 `linkUserRole`→`ensureRole` 会隐式 `INSERT … system=0` 创建角色行，**API 写路径禁止复用**该方法；S2 须另写「仅链接已存在角色」的事务（先校验 roles 表存在，再双写），确保管理 API 永不隐式造角色。
- 分配实现沿用 GOAL-006 阶段 B 双写：`users.roles` legacy JSON + `user_roles` 关系在同一事务写入，读侧 `userWithRoles` 集合一致性校验保持。

### 2.4 self / 最后管理员保护

- **self 保护**：不能删除自己 → `SELF_OPERATION`(409)；不能移除自己的 `admin` 角色 → `SELF_OPERATION`(409)。
- **最后管理员保护**：删除用户或移除某用户的 `admin` 角色，若会导致**零个拥有 admin 角色的用户** → `LAST_ADMIN`(409)。
- 种子 `user-admin` 是普通数据行（无独立 system 标志），由上述不变量保护。
- **实现依托（v0.2.0 · A-002 F-001）**：操作者身份经 §7 actor 通道传入实体，S2 在 store 事务内按 `actor.ID` 判定 self、按 admin 用户计数判定 last-admin。

### 2.5 用户删除的级联边界

- 删除用户必须**原子**清理：`refresh_tokens`（该用户的行，撤销会话）、`user_roles`（ON DELETE CASCADE）、`users` 行；失败整体回滚。
- `operation_log.actor_id` 不设外键，历史操作日志不受用户删除影响（保留审计链）。

### 2.6 username 唯一

- `users.username` 为 DB UNIQUE（既有约束）；重复 → `USERNAME_TAKEN`(409)。trim 后空 → 既有 `INVALID_CREATE_FIELD`/`INVALID_PATCH_FIELD` 路径。

## 3. roles 资源契约

| 项 | 值 |
|----|----|
| 资源 id | `roles` |
| 端点 | `/api/roles`（五路由） |
| 权限键 | `roles.read` / `roles.write` |
| NOT_FOUND | `ROLE_NOT_FOUND`（404） |
| list | `sortFields = [key, name, updatedAt]`；`qSearch = true`（key/name 子串） |
| create 字段 | `key`(必填)、`name`(必填) |
| patch 字段 | `name`（可选）；**`key` 不可变** |
| id | `role-<key>`（与既有 `ensureRole` 派生一致） |
| system | **API 不可设置**；用户创建角色固定 `system = 0` |

### 3.0 公开字段（响应形状 · v0.2.0 · A-002 F-003）

list/detail/create/update 的响应行固定为（与 users §2.1 同精度）：

```json
{
  "id": "role-admin",
  "key": "admin",
  "name": "admin",
  "system": true,
  "createdAt": "2026-08-03T00:00:00.000Z",
  "updatedAt": "2026-08-03T00:00:00.000Z"
}
```

- `system`：**boolean**（JSON `true`/`false`，映射 `roles.system` 0/1）；S2/S4 前端可据此禁用 system 角色的编辑/删除控件，无需依赖 409。
- `createdAt`/`updatedAt`：RFC3339 固定 3 位毫秒形状（与 users/records D-004 一致）。
- `key` 不可变；`system` 只读（API 不可设置）。

### 3.1 key 格式

- create 的 `key` 必须匹配 `^[a-z][a-z0-9_-]*$`（与 `roleKeyRe` 一致）→ 否则 `INVALID_ROLE_KEY`(400)。
- 重复 key → `ROLE_KEY_TAKEN`(409)。

### 3.2 system role 保护

- `system = 1` 的角色（种子 `admin`/`editor`/`viewer`）为种子管理：**不可 delete、不可改 key、不可改 name** → `ROLE_SYSTEM`(409)。只读。
- `system = 0` 的用户自定义角色：可 update `name`、可 delete（无关联用户时）。

### 3.3 删除保护

- 有关联用户（`user_roles` 非空）的角色不可删除 → `ROLE_IN_USE`(409)；DB `ON DELETE RESTRICT` 为兜底。
- system 角色不可删除 → `ROLE_SYSTEM`(409)。

### 3.4 grant 约束

- 本资源**不**管理 `role_permissions` / `role_menu_items`（grant 的 CRUD 界面非目标——完整 IAM 为 Out）。
- 既有种子 grants（admin/editor/viewer → 各权限/菜单）保持冻结；`system=1` 角色的 grants 由种子唯一权威，API 无修改面，自然满足「system 角色 grant 约束」。
- 用户自定义角色当前无 grants（`roles.*` 权限仅由种子授予系统角色）；后续 grant 管理属扩展，非本契约范围。

## 4. 最小 IAM 边界与权限投影

- **新权限键**：`users.read`、`users.write`、`roles.read`、`roles.write`。
- **种子 grants**（S3 起替换 records grants）：
  - `admin` → `users.read`+`users.write`+`roles.read`+`roles.write` + 菜单 `menu_users`、`menu_roles`
  - `editor` → `users.read`、`roles.read`（只读）
  - `viewer` → `users.read`、`roles.read`（只读）
- **菜单项**（seed 新增）：`menu_users`（page_ref = users 页，feature_key `menu_users`）、`menu_roles`（page_ref = roles 页，feature_key `menu_roles`）；admin 唯一授权。
- **401/403 负向路径必测**：匿名访问 users/roles 五路由 → 401 `UNAUTHENTICATED`；缺权限（如 viewer 写）→ 403 `FORBIDDEN`。
- `account.StaticDevSession`（dev 兜底）的 permissions/features 同步 users/roles 键与菜单（否则 dev 兜底与种子不一致）。

## 5. 操作日志

- **事件扩展（migration 0005）**：`operation_log.event` CHECK 新增 `users.create`/`users.update`/`users.delete`、`roles.create`/`roles.update`/`roles.delete`；**保留** `records.*`（历史行）与 `auth.*`（登录态事件仍活跃）。
- users/roles 写路径挂 `OnWrite`，事件与 detail：users → detail `{"username":"…"}`；roles → detail `{"key":"…"}`；**永不**含密码/令牌。
- 现有 auth.login/logout/refresh 事件保持不动。

## 6. 错误码（对 I-010-001 §5 的限定扩展）

- **Envelope 形状不变**：`{error, message}`；不新增字段。
- **既有通用码保持**：`UNAUTHENTICATED`(401)、`FORBIDDEN`(403)、`INVALID_SORT_FIELD`/`INVALID_SORT_ORDER`/`INVALID_PAGE`/`INVALID_PAGE_SIZE`(400)、`INVALID_CREATE_BODY`/`INVALID_CREATE_FIELD`(400)、`INVALID_PATCH_BODY`/`INVALID_PATCH_FIELD`(400)、`INTERNAL`(500)。
- **新增资源特定码**（users/roles 领域冲突，409 为对 I-010-001 §5「不引入 409/业务唯一冲突」的**显式、限定范围偏离**，仅账号域）：
  | 资源 | 码 | 状态 | 触发 |
  |------|----|------|------|
  | users | `USER_NOT_FOUND` | 404 | id 不存在 |
  | users | `USERNAME_TAKEN` | 409 | username 重复 |
  | users | `LAST_ADMIN` | 409 | 删除/降级将导致零 admin 用户 |
  | users | `SELF_OPERATION` | 409 | 删除自己 / 移除自己 admin |
  | users | `INVALID_ROLE_REF` | 400 | roles 含未注册 key |
  | roles | `ROLE_NOT_FOUND` | 404 | id 不存在 |
  | roles | `ROLE_KEY_TAKEN` | 409 | key 重复 |
  | roles | `ROLE_IN_USE` | 409 | 删除有关联用户的角色 |
  | roles | `ROLE_SYSTEM` | 409 | 对 system 角色做禁用修改/删除 |
  | roles | `INVALID_ROLE_KEY` | 400 | key 格式非法 |

## 7. 通用工厂最小扩展（S2 落地）

对 [I-010-001 §4/§5](../../GOAL-010-a002-schema-adapter/attachments/I-010-001-schema-resource-contract.md) 的**追加式**扩展（不改既有字段语义）：

1. `Resource.JSONFields []string`：新增字段声明。当 body 中出现这些字段时，以**原始 JSON 值**解码并透传 entity（不强制 string）；缺省视为未提供（create → entity 用默认值，patch → 不触碰）。与 `CreateFields`/`PatchFields`（string、必填/trim 校验）互补。users 用 `roles`。
2. **actor 通道（v0.2.0 · A-002 F-001）**：`ResourceEntity` 的 `Create`/`Update`/`Delete` 签名增传 `account.User`（请求操作者，工厂已从 `requirePermission` 取得），使实体可诚实实现 `SELF_OPERATION`/`LAST_ADMIN` 等按操作者判定的领域不变量。`List`/`Get` 不传（只读无需）。既有 records 实体的对应方法签名补齐 `_ account.User` 参数并忽略——**零对外行为变化**，records 写路径语义不变。
3. `DomainError{Status, Code, Message}`：entity 可返回的类型化领域错误；工厂在 create/update/delete/detail 的 `errors.As` 识别并**逐字映射**（status + `{error,message}`）。**检查顺序**：先识别 `DomainError`，再 `store.ErrNotFound`（→ 404），未识别错误才走 `INTERNAL`(500)（v0.2.0 · A-002 F-006 / A-001 F-003）。当前 create 路径的 `ErrRecordExists` 重试循环保持在 `DomainError` 检查之后。
4. 既有 records 资源行为不变（零对外变更历史事实保持）。

## 8. 非目标

- 密码复杂度策略、密码找回/重置流程；SSO/SCIM；完整 IAM 产品；grant 管理界面；多租户；乐观锁/软删除；批量操作；扩大 `I-PROTO-001 v0.1.3` 覆盖。
- 不把 users/roles 当作无约束平面 CRUD（本契约 §2.2/§2.4/§3.2/§3.3 为硬性领域不变量）。

## 9. 证据索引

- `apps/api/internal/handler/resources.go`（`Resource`/`ResourceEntity`/`registerResource`——§7 扩展对象）
- `apps/api/internal/handler/records.go`（records 注册实例——§7 兼容基线）
- `apps/api/internal/store/store.go`（users/roles/user_roles、`CreateUser`、`userWithRoles`、`linkUserRole`、`ensureRole`）
- `apps/api/internal/store/migrate.go`（users/roles/permissions/menu 表 DDL；0005 扩展对象）
- `apps/api/internal/auth/auth.go`（`accountFromUser` 权限投影、`account.User` 隔离）
- `apps/api/internal/account/session.go`（`User`/`StaticDevSession`——§4 dev 兜底同步对象）
- `apps/api/internal/store/operations.go` + `migrate.go`（operation_log event CHECK——§5 扩展对象）
- `apps/api/internal/store/seed.go`（seedRBAC——§4 权限/菜单种子替换对象）
- [I-010-001 v0.2.1](../../GOAL-010-a002-schema-adapter/attachments/I-010-001-schema-resource-contract.md)（通用基础；§4/§5 扩展来源）

## 10. 修订记录

| 版本 | 日期 | 变更 |
|------|------|------|
| 0.1.0 | 2026-08-03 | 冻结（GOAL-011 D-002，用户裁决三项：通用工厂+最小契约扩展、操作日志纳入、records 硬退场）；关闭 `I-011-001` |
| 0.2.0 | 2026-08-03 | A-002 响应（GOAL-011 D-003）：§7 增 `account.User` actor 通道（F-001，SELF_OPERATION 可诚实实现）+ DomainError 检查优先级（F-006）；§2.3 禁 API 路径 `ensureRole` 隐式建角色（F-004）；§3.0 冻结 roles 公开响应形状（F-003）。`I-011-001` 维持 `verified` |
