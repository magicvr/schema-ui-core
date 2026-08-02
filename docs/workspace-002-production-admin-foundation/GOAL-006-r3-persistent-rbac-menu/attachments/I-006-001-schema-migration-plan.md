---
title: I-006-001 · R3 SQLite 版本迁移、DDL 与增量种子计划
status: active
doc_type: info-collection
created: 2026-08-02
updated: 2026-08-02
parent: GOAL-006-r3-persistent-rbac-menu
version: 0.1.0
related_info: I-006-001
related_decision: D-002
---

# I-006-001 · R3 SQLite 版本迁移、DDL 与增量种子计划

> **结论**：本附件与 D-002 关闭 `I-006-001`，作为 S1/S2/S3 的实施输入。以下 SQL、迁移函数与测试名是冻结计划，不是已执行代码或数据库事实。

## 1. 当前兼容基线

| 事实 | 必须保持的边界 |
|------|----------------|
| `Store.Open()` 当前按 `migrate()` → `seedAdmin()` 启动，单个 `*sql.DB` 的最大连接数为 1 | 可替换迁移/seed 内部实现，但不破坏服务启动调用与错误返回 |
| R2 表为 `users` + `refresh_tokens`；`refresh_tokens.user_id` 指向 `users.id` | 用户 id、密码、refresh hash/撤销/过期关系在迁移后不变 |
| `store.User.Roles []string` 当前来自 `users.roles` JSON | 两步迁移期间保持 `UserByID/UserByUsername` 和 `account.User {id,name,roles}` 对外形状 |
| `CreateUser`、`UserByID`、`UserByUsername` 被 auth 登录、JWT middleware 与 refresh 使用 | 签名与 `ErrNotFound` 语义保持；阶段 A/B 只切换角色来源 |
| records GET 当前公开，写 gate 直接检查 admin 角色 | S4 才切 permission gate；DDL/seed 必须先提供可查询关系 |

证据：`apps/api/internal/store/store.go`（`Store` / `Open` / `migrate` / `seedAdmin` / user/refresh 方法）、`apps/api/internal/auth/auth.go`（`Login` / `Refresh` / `Middleware` / `accountFromUser`）、`apps/api/internal/handler/records.go`（routes / `writeGate`）。

## 2. 迁移台账与执行不变量

```sql
CREATE TABLE schema_migrations (
  version    INTEGER PRIMARY KEY,
  name       TEXT NOT NULL UNIQUE,
  checksum   TEXT NOT NULL CHECK (length(checksum) = 64),
  applied_at INTEGER NOT NULL
);
```

1. 编译期迁移清单按 `version` 严格递增；checksum = 规范化 SQL + 显式 data-transformer id 的 SHA-256 小写 hex。
2. 启动时先判断台账表是否存在：存在则读取全部已应用版本，未知版本、缺失中间版本、重复 name 或 checksum 不一致均拒绝启动；不存在则进入 `0001` bootstrap 事务，禁止在结构指纹核对前单独创建空台账。
3. 每个连接在任何迁移/查询前执行并核对 `PRAGMA foreign_keys=ON`；`SetMaxOpenConns(1)` 保留，但不替代该断言。
4. 每个迁移在一个事务内完成 DDL、数据变换和台账插入；任一步失败则整体回滚。
5. 非空文件库在首个待执行迁移前用 SQLite `VACUUM INTO` 产生一致性快照，路径形如 `<db>.pre-v0002-<UTC>.sqlite`；目标路径必须通过驱动安全绑定/引用，不经 shell 拼接。
6. 快照打开后必须 `PRAGMA integrity_check` 返回单行 `ok`；迁移后主库还须 `integrity_check=ok` 且 `foreign_key_check` 无行。

## 3. 迁移版本

| Version | Name | 空库 | 既有 R2 库 | 失败边界 |
|---------|------|------|------------|----------|
| `0001` | `r2_baseline` | 在一个 bootstrap 事务内创建现有 `users`、`refresh_tokens`、`idx_refresh_tokens_user_id`、迁移台账并登记版本 | 在同一 bootstrap 事务内先通过 `sqlite_master`、`table_info`、`foreign_key_list`、`index_list` 核对 R2 最小结构，再创建迁移台账并登记版本 | 部分表、缺列、冲突类型/约束或未知结构均回滚且不留下空台账 |
| `0002` | `rbac_expand` | 创建规范化关系；空数据回填 | 创建规范化关系；解析每个 `users.roles`，验证 role key、去重并回填 `roles` / `user_roles` | 非数组、空 key、重复 role、非法 key、唯一/FK 冲突或部分回填均事务回滚 |

旧角色 key 必须匹配 `^[a-z][a-z0-9_-]*$`；该正则由 `0002` 的 Go data transformer 在任何角色插入前逐项校验，DDL 的 `CHECK (key <> '')` 只承担最低数据库约束。派生 role id 为 `role-<key>`。迁移比较按**集合**判断相等；阶段 B 对外 `roles` 按 key 升序输出，保证确定性。

## 4. `0002` 精确 DDL

```sql
CREATE TABLE roles (
  id         TEXT PRIMARY KEY,
  key        TEXT NOT NULL UNIQUE CHECK (key <> ''),
  name       TEXT NOT NULL,
  system     INTEGER NOT NULL DEFAULT 0 CHECK (system IN (0, 1)),
  created_at INTEGER NOT NULL,
  updated_at INTEGER NOT NULL
);

CREATE TABLE user_roles (
  user_id TEXT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  role_id TEXT NOT NULL REFERENCES roles(id) ON DELETE RESTRICT,
  PRIMARY KEY (user_id, role_id)
);
CREATE INDEX idx_user_roles_role_id ON user_roles(role_id);

CREATE TABLE permissions (
  id          TEXT PRIMARY KEY,
  key         TEXT NOT NULL UNIQUE CHECK (key <> ''),
  description TEXT NOT NULL DEFAULT '',
  created_at  INTEGER NOT NULL,
  updated_at  INTEGER NOT NULL
);

CREATE TABLE role_permissions (
  role_id       TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  permission_id TEXT NOT NULL REFERENCES permissions(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, permission_id)
);
CREATE INDEX idx_role_permissions_permission_id ON role_permissions(permission_id);

CREATE TABLE menu_items (
  id          TEXT PRIMARY KEY,
  page_ref    TEXT NOT NULL UNIQUE CHECK (page_ref <> ''),
  feature_key TEXT NOT NULL UNIQUE CHECK (feature_key <> ''),
  sort_order  INTEGER NOT NULL DEFAULT 0,
  enabled     INTEGER NOT NULL DEFAULT 1 CHECK (enabled IN (0, 1)),
  created_at  INTEGER NOT NULL,
  updated_at  INTEGER NOT NULL
);

CREATE TABLE role_menu_items (
  role_id     TEXT NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
  menu_item_id TEXT NOT NULL REFERENCES menu_items(id) ON DELETE RESTRICT,
  PRIMARY KEY (role_id, menu_item_id)
);
CREATE INDEX idx_role_menu_items_menu_item_id ON role_menu_items(menu_item_id);
```

删除语义：删除用户级联其角色关系；删除角色前必须先解除用户关系，删除成功后级联其 permission/menu grants；删除 permission 或 menu item 前必须先解除所有 grants，避免静默改变有效策略。

## 5. 两步读写切换

### 阶段 A · 扩展与核对

- 应用 `0001/0002`，回填旧角色并运行增量 seed。
- `CreateUser` 在同一事务内写 JSON 与 `user_roles`；`UserByID/UserByUsername` 仍返回旧 JSON，但同时查询规范化角色并做集合比对，不一致即返回可诊断错误。
- 不改变 auth/account 对外签名，不删除旧列。

### 阶段 B · 规范化权威

- `UserByID/UserByUsername` 从 join 关系重建角色，按 role key 升序输出；仍读取旧 JSON 做集合比对。
- `CreateUser` 继续双写，直到独立后续迁移显式删除/停用旧列。
- `refresh_tokens.user_id`、JWT subject、密码字段与用户主键不变。

## 6. 稳定增量 seed

| 类型 | Stable id / key | 关系 |
|------|-----------------|------|
| role | `role-admin` / `admin` | `system=1`；`records.read` + `records.write`；菜单 `menu-list-edit-lifecycle` |
| role | `role-editor` / `editor` | `system=1`；仅 `records.read`，保持 R2 非 admin 不可写 |
| role | `role-viewer` / `viewer` | `system=1`；仅 `records.read` |
| permission | `perm-records-read` / `records.read` | records GET gate |
| permission | `perm-records-write` / `records.write` | records PATCH/DELETE gate |
| menu | `menu-list-edit-lifecycle` / `list-edit-lifecycle` / `menu_list_edit_lifecycle` | `enabled=1`；仅 admin grant |
| user | `user-admin` / `admin` | 保留 `admin` + `editor` 角色与现有密码/身份边界 |

每个实体与关系独立 ensure；已有任意用户不会跳过角色、权限、菜单或 grants。seed 不覆盖非种子用户字段，也不把 editor 自动升级为写权限。

## 7. 实施验证矩阵

| ID | 必须证明 |
|----|----------|
| `V-MIG-01` | 空库顺序应用 0001/0002；重复启动不重复执行 |
| `V-MIG-02` | 既有 R2 库登记 0001、应用 0002；用户/密码/refresh id 与关系不变 |
| `V-MIG-03` | 非法/重复 roles、部分基线、未知版本、checksum 漂移均 fail closed 且无部分提交 |
| `V-MIG-04` | `foreign_keys=ON`、unique/FK/delete 正反行为与反向索引存在 |
| `V-MIG-05` | 阶段 A/B 的 legacy/normalized 集合核对；规范化输出顺序确定 |
| `V-SEED-01` | 任意已有用户场景仍补齐 stable entities/grants；重复 seed 无重复且不覆盖用户字段 |
| `V-REC-01` | pre-v0002 快照可在新路径打开，`integrity_check=ok`，原用户/refresh 查询可复现 |
| `V-REC-02` | 迁移后主库 `integrity_check=ok`、`foreign_key_check` 空，身份/权限/菜单关键查询通过 |
| `V-REG-01` | 现有 `TestSeedAdminIdempotent`、`TestCreateUserAndLookup`、`TestRefreshTokenLifecycle` 与 auth/account 回归保持 |

实现时新增测试名可按行为命名；本附件不把计划中的测试名或 SQL 记为已通过事实。
