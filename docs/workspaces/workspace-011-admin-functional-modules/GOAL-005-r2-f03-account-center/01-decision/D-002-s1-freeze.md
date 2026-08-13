---
id: D-002
goal: GOAL-005-r2-f03-account-center
title: S1 · 方案冻结 — admin.account 模块（会话模型 + 启停 + 权限键 + Profile 扩展声明）
date: 2026-08-14
status: accepted
parent: GOAL-005-r2-f03-account-center
created: 2026-08-14
updated: 2026-08-14
version: 0.1.0
---

# D-002 · S1 方案冻结（F-03 个人中心与账户安全 + 账号启停）

> 依据：I-011-001 `3 F-03、`8 必办-5；GOAL-005 00-meta 边界与 S1 门禁。
> 本文件冻结以下内容；实现（S2）按此执行，偏差须回炉本决策。

## 1. 模块设计：`admin.account`（新标准 Admin 模块，单一模块承载 F-03）

| 项 | 冻结值 |
|----|--------|
| 模块 ID | `admin.account` |
| 依赖 | core.auth-session / core.navigation-capability / core.schema-render / core.operationlog |
| 能力 | StandardAdminCapabilities()（与既有标准模块一致） |
| 路由 | `GET /api/account/profile`、`PATCH /api/account/profile`、`POST /api/account/password`、`GET /api/account/sessions`、`POST /api/account/sessions/{id}/revoke`、`POST /api/users/{id}/enable`、`POST /api/users/{id}/disable`、`POST /api/users/{id}/unlock` |
| 页面 | `account`（个人中心） |
| 导航 | `menu_account` → navigation.user 区 |
| 权限键 | `users.enable`（启用+手动解锁）、`users.disable`（停用）— 均 PolicyAdmin |
| 持久化 | 迁移 `0013:account-enable-state`：`ALTER TABLE users ADD COLUMN enabled INTEGER NOT NULL DEFAULT 1`（owner = admin.account） |

**边界声明（I-011-001 `1 C-01）**：启停属 F-03 产品态能力，不并入 admin.users CRUD；但 `users.enable/disable` 键名与 `/api/users/{id}/enable|disable|unlock` 路径保持与 users 资源一致的命名空间（A-002 F-003 落地语义）。

## 2. 会话模型（I-001 关闭）

- **会话表 = `refresh_tokens`**（既有基建）：refresh token 是持久会话通道；access token 短时（accessTTL）不列示。
- **列表**：当前身份全部 refresh token（含已吊销），`created_at DESC` 分页（pageSize 上限 100）。行字段：`id`、`createdAt`、`expiresAt`、`revokedAt`、`status`（`active`｜`revoked` 计算字段）。
- **吊销粒度**：单条 `POST /api/account/sessions/{id}/revoke`，仅限本人（`id + user_id` 双条件）；幂等（已吊销 → 204）。
- **token_version 联动**：
  - 改密（自助 `POST /api/account/password` 或管理员改密）→ `token_version+1` → 全部已签发 access token 立即失效（复用 W4 P0-3 机制）。
  - **停用** → `enabled=0` + `token_version+1` + 吊销全部 refresh token → 即时全面失效（含活跃 access）。
  - 单条会话吊销 → 仅吊销该 refresh token；对应 access 在 accessTTL 内仍有效（吊销的是「会话通道」，短窗自过期——已文档化的残余）。
  - 锁定（C-11）→ 保持既有语义（锁定开窗时吊销 refresh）；enabled 正交不受影响。

## 3. 启停状态字段与权限（I-002 关闭）

| 操作 | 语义 | 守卫 | 错误码 |
|------|------|------|--------|
| disable | `enabled=0` + token_version+1 + 吊销全部 refresh | 不可停用自己；不可停用最后一名 admin | USER_NOT_FOUND / SELF_OPERATION / LAST_ADMIN（409） |
| enable | `enabled=1`（不清锁/不清失败计数） | 仅存在性 | USER_NOT_FOUND |
| unlock（手动解锁） | `locked_until=0` + `failed_login_count=0` | 仅存在性 | USER_NOT_FOUND |

- **登录拒绝**：Login 见 `enabled=0` → `403 ACCOUNT_DISABLED`（与 423 ACCOUNT_LOCKED 同类的显式终态；停用是管理员可见操作，不构成枚举 oracle 增量）。Refresh 见 `enabled=0` → 401 UNAUTHORIZED envelope（fail-closed）。中间件见 `enabled=0` → 401 UNAUTHENTICATED（与 token superseded 同 envelope，fail-closed——停用即时生效）。
- **操作日志**：新增事件 `users.enable` / `users.disable` / `users.unlock` / `account.password-change` / `account.session-revoke`（operationlog 事件常量为追加式，无白名单校验）。

## 4. 自助端点细节

| 端点 | 请求 | 响应 | 说明 |
|------|------|------|------|
| GET /api/account/profile | — | `{id, username, name, enabled, createdAt, updatedAt}` | 自服务视图；仅需身份，无权限键 |
| PATCH /api/account/profile | `{name}` | 更新后行 | name 非空校验（复用 INVALID_PATCH_FIELD 约定） |
| POST /api/account/password | `{currentPassword, newPassword}` | 204 | 校验当前密码 + 新密码 8–72 字节；成功后 token_version+1 + 吊销全部 refresh（复用 UpdateUser password 语义）；错误 `INVALID_PASSWORD` 400 |
| GET /api/account/sessions | — | `{items,total,page,pageSize}` | 会话列表（`2） |
| POST /api/account/sessions/{id}/revoke | — | 204 | 仅本人；未知/他人 → `SESSION_NOT_FOUND` 404 |

## 5. 前端设计

- **account 页**（`apps/api/internal/modules/account/schema/account.json`，schema 驱动）：
  - 个人信息 form：`recordSource` GET /api/account/profile（ADR-0021 prefill）+ name 字段 + `saveProfile`（PATCH）。
  - 修改密码 form：currentPassword/newPassword（password 控件）+ `changePassword`（POST）。
  - 会话 table：dataSource `/api/account/sessions`，列 createdAt/expiresAt/status；行操作 revoke（`POST /api/account/sessions/{id}/revoke` + confirm；`disabledWhen {field:"status", equals:"revoked"}`）。
  - 导航：navigation.user（右上用户区），可见性 PolicyAdminEditorViewer（自服务人人可用）。
- **users 页扩展**（admin.users 拥有 schema，跨模块动作引用，fail-open）：
  - 行新增 `enabled`（bool）列。
  - 行操作 enable / disable / unlock：
    - `permissions.edit` 表达式按权限键（`users.enable` / `users.disable`）——admin.account 未启用时键不存在 → 按钮隐藏（视觉 fail-open；服务端仍 403/404 fail-closed）。
    - `disabledWhen {field:"enabled", equals:true|false}` / `{field:"locked", equals:false}` 按行状态门控。
  - users API 行扩展：`enabled`（bool）、`locked`（计算：lockedUntil > now）。
- **i18n**：en-US.json / zh-CN.json 新增 `manifest.title.account`、`manifest.nav.account`、`schema.account.*`、`schema.users.column.enabled`、`schema.users.action.enable|disable|unlock`、`schema.users.confirm.*` 键。
- **StaticDevSession**：permissions += `users.enable`、`users.disable`；features += `menu_account`（dev 会话与种子权限对齐）。

## 6. Profile 内容扩展声明（防误触「不改变 Profile 默认集」门禁）

- **`admin.account` 加入 mvp + admin 默认启用集**（kernel/profile.go）。定性：**Profile 内容扩展**——自服务改密/会话管理是账号安全基线能力（与 users/roles 同级），经既有模块贡献机制（provider + fragment + reconcile）落地；**不改装配语义**：不改 Manifest 聚合规则、协议 pin（v2.8.0）、共同门禁、capability 语义。与 F-01 必办-3 同一模式。
- `adminFunctionalOrder` 尾部追加 `admin.account`（home 判定仍以 users 为首；仅当 users/roles/settings/activity 全禁用时 account 才可能成为 home——可接受、显式留痕）。
- 相应更新 `scripts/smoke.sh` SM-007 页面集：mvp = users roles account；admin = users roles settings activity account。

## 7. 必办核对（I-011-001 `8）

| 必办 | 适用 | 处置 |
|------|------|------|
| 必办-1（协议对照） | F-01/F-02 专属 | 不适用（本目标无协议面新增；fail-open 惯例沿 `2/`5 声明） |
| 必办-2（通知边界） | F-04 专属 | 不适用 |
| 必办-3（home 装配） | F-01 专属 | 不适用（`6 仅追加 order 尾部，不改变装配语义） |
| 必办-4（订单） | 领域模块 | 不适用 |
| **必办-5（启停端点+权限键+会话端点）** | **适用** | **✅ 本方案 `1/`3/`4 全部落地** |

## 8. 未选方案（留痕）

- 不引入「当前会话」标识（refresh 轮换后客户端持有的最新 token 即当前会话；列表仅列示+吊销，不标 current）。
- 不改 C-11 锁定语义：停用/启用/解锁与锁定窗口正交（I-003 关闭：disable 不重置失败计数与锁定窗口；enable 不清锁；unlock 才清）。
- 不把启停并入 admin.users 模块（模块边界：F-03 自持端点与权限键；users 页仅 schema 动作引用，跨模块依赖以权限键表达、fail-open）。
- 不加会话过期清理任务（保留期/清理策略归 R3 或 F-04 持久化策略统一处理）。
- 会话列表不暴露 token 明文或 hash（仅内部 id 与元数据）。

## 9. 实现范围（S2 清单）

1. 迁移 0013 + compiled provider 注册（admin.account 迁移加入 PersistenceProviders）。
2. authsession：`User.Enabled` + 全部 SELECT/INSERT 列清单更新 + `ListRefreshTokensForUser` + `RevokeRefreshTokenIfOwned` + `SetUserEnabled`（含守卫）+ `UnlockUser` + 相关 sentinel。
3. auth：Login/Refresh/Middleware 的 enabled 检查 + `ErrAccountDisabled`。
4. handler：account 端点 + users enable/disable/unlock 端点（复用 requirePermission）。
5. 模块 provider + manifest fragment + schema（account.json）+ users.json 扩展 + i18n + StaticDevSession。
6. kernel.BuiltinModules + profileDefaults（mvp/admin）+ adminFunctionalOrder + composition 装配。
7. smoke.sh SM-007 页面集更新。
8. 测试：Go 单测/集成（migration、auth 拒绝语义、account handler、启停守卫、权限门禁）+ Web 回归 + i18n 键。
