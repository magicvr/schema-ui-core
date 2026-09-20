---
id: r1-c2-fk-parent-rebuild-finding-v1.0-fc
doc_type: design-attachment
title: R1 C2 FK 父表重建阻塞发现与已验证处置模式 v1.0
status: open-finding
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 1.0.0
---

# R1 C2 FK 父表重建阻塞发现 v1.0

> **状态：`open-finding`，待 independent 复审（用户 2026-09-20 裁决：先 `/audit` 再定方案）。** 本文件只记录**已实测的技术事实**与候选处置模式，**不**做结构决策。用户已明确：FK 父表族的处置方式（扩写 R1 / 拆子目标 / 其他）在 independent 审计之后决定。

## 1. 事实（本机 SQLite 3.51.2 实测）

### F-1 · 重命名 FK 父表会**永久**改写子表的 `REFERENCES` 文本

```text
ALTER TABLE users RENAME TO users_old;
→ sqlite_master 中 refresh_tokens 的 DDL 变为：
  CREATE TABLE refresh_tokens (..., user_id TEXT NOT NULL REFERENCES "users_old"(id) ON DELETE CASCADE)
```

- 该改写**不是瞬时的**：`DROP TABLE users_old` 完成后，子表 DDL 仍指向已不存在的 `users_old`。
- 后果：后续 `DELETE FROM users WHERE id='u1'` 报 `no such table: main.users_old`；子表 DML 同样失败。即**朴素重建会把该父表下所有 FK 静默损坏**。

### F-2 · `PRAGMA legacy_alter_table=ON` **不能**阻止该改写

实测：事务外设置该 pragma 后执行 rename，子表 DDL 仍被改写为 `REFERENCES "users_old"(id)`。该 pragma 只影响 rename 时对**视图/触发器**的处理，不影响 FK 文本改写。

### F-3 · `PRAGMA foreign_keys=OFF` 在本项目**不可用**

- DSN 固定 `_foreign_keys=on`（`apps/api/internal/store/store.go:57`）且 `SetMaxOpenConns(1)`（`:107`）；`Store.migrate` 启动即 `assertForeignKeysOn`（`migrate.go:32`，实现 `:253-265`）。
- `PRAGMA foreign_keys` 在**事务内是 no-op**；而 `applyMigration`（`migrate.go:108-132`）把每个 descriptor 的 `Apply` 与 ledger 写入包在**同一个事务**里。
- 故 SQLite 官方 12 步 rebuild 的前置条件（事务外关 FK）无法满足。

### F-4 · 朴素重建在子表有行时会**级联删除子表数据**

实测：`dict_types`（父）+ `dict_entries`（子，`ON DELETE CASCADE`）在 `foreign_keys=on` 下：
- 仅 `ALTER TABLE dict_types RENAME TO dict_types_old` → `stock` 子表 FK 被改写；
- `DROP TABLE dict_types_old` → **`dict_entries` 全部行被级联删除**，存活子表 schema 仍指向已删除的表。
- 若子表非空，`DROP TABLE <t>_old` 也可能直接报 `FOREIGN KEY constraint failed`。

### F-5 · 已验证可行的处置模式（候选，未定案）

**子女先行（drop children first）**，与仓库既有 `rebuildOperationLogWithSessions`（`operationlog/migration/migration.go:709-754`）同构：

```text
1. CREATE TEMP TABLE <child>_backup AS SELECT <cols> FROM <child>;   -- 无 FK 子句，父表 DROP 无法级联进 temp
2. DROP TABLE <child>;                                               -- 先摘掉子表
3. ALTER TABLE <parent> RENAME TO <parent>_old;                      -- 此刻无子表可被改写
4. CREATE TABLE <parent> ( ...目标形状... );
5. INSERT INTO <parent> (...) SELECT <转换表达式> FROM <parent>_old;
6. DROP TABLE <parent>_old;
7. CREATE TABLE <child> ( ...子表 DDL 逐字，FK 名此刻解析到新父表... );
8. INSERT INTO <child> (...) SELECT <转换后> FROM <child>_backup;
9. DROP TABLE <child>_backup;
10. CREATE INDEX ...（全部索引最后建）
```

实测结论（`dict_types`/`dict_entries` 整模块单事务）：2/2 行保留、子表 FK 文本恢复为 `REFERENCES dict_types(key)`、`PRAGMA foreign_key_check` 空、`integrity_check` = ok。

**关键点**：子表必须从 **TEMP 快照**回填，不能把子表 rename 成 `<child>_old` 当数据源——父表 rename 会把 `<child>_old` 的 FK 也改写，随后父表 `DROP` 会把你要拷的数据一起级联删掉。

## 2. 受影响的 FK 父表（workspace-040 范围内）

| 父表（本 VP 内重建） | 引用它的子表 | 子表是否在本 VP 列范围内 |
|----------------------|--------------|--------------------------|
| `users` | `refresh_tokens`, `email_verification_challenges`, `password_recovery_challenges`, `login_failures`, `user_password_history`, `user_invites` | 是（v74 同 descriptor） |
| `users` | `user_roles`, `notifications`, `user_mfa`, `mfa_proofs` | **部分否**：`notifications`(v81)、`user_mfa`/`mfa_proofs`(v80) 在其他 descriptor；`user_roles` 不在 v74 的 12 表清单内 |
| `roles` | `user_roles`, `role_permissions`, `role_menu_items` | **否**（联接表不在 v74 清单内） |
| `permissions` | `role_permissions` | **否** |
| `menu_items` | `role_menu_items` | **否** |
| `operation_log` | `operation_log_session`, `operation_log_correlation` | 否（同模块旁表；已有 dance 可复用） |
| `dict_types` | `dict_entries` | 是（v77 同 descriptor） |
| `scheduled_tasks` | `task_runs` | 是（v83 同 descriptor） |

> `users.roles` 是 **JSON TEXT 数组**，`users` 本身**没有** FK 列；`users` 只是被引用方。

## 3. 附带的两个新发现（同轮实测/核对）

### F-6 · `users` 有 3 个**跨模块** ALTER 列，遗漏即丢数据

| 列 | 引入方 | 位置 |
|----|--------|------|
| `enabled` | `account` 模块 v0013 | `apps/api/modules/account/migration/migration.go:21` |
| `notifications_enabled` | `notifications` 模块 v0017 | `apps/api/modules/notifications/migration/migration.go:46` |
| `avatar_url` | `account` 模块 v0035 | `apps/api/modules/account/migration/migration.go:28` |

v74 的 `users` 新 DDL **必须**包含这三列（连同 authsession 内的 6 个 ALTER 列：`token_version`、`failed_login_count`、`locked_until`、`must_change_password`、`email`、`email_status`、`last_login_failure_at`），否则重建静默丢列。**注意最终列序以 live 库为准**：SQLite 把 ALTER 追加列放在末尾，跨模块 ALTER 的**真实顺序**由全局版本序 v0012 → v0013 → v0017 → v0035 → v0038 → v0054 → v0061 决定，**不是**按文件分组。

### F-7 · `operation_log` 的 PG DDL 由正则派生，改 SQLite 字面会**静默失效**

`operationlog/migration/migration.go:250` 的 `pgTimeColRe` 匹配 `(created_at|archived_at)\s+INTEGER NOT NULL`；一旦 SQLite 字面改为 `TEXT NOT NULL`，该派生变成**静默 no-op**（PG 侧保留陈旧 BIGINT 映射且不报错）。→ 这些表的 PG DDL 必须改为**显式书写**，与 `D-017` 的单 checksum 约定一致（PG 不进 checksum，但必须显式落盘）。

### F-8 · `schema_migrations` 自引用无危害，但 `restoreLedger` 路径有残留

- **无中途读**：ledger 只在 `applyPending` 之前读取（`identity.go:334`、`migrate.go:63`），循环内与 `applyMigration` 内都不再读；`verifyIntegrity` 在批次后（`migrate.go:102`）。
- **无自插入**：`applyMigration` 先 `Apply` 再 `INSERT INTO schema_migrations (version,…)`（`migrate.go:113-123`），同事务、插入在后；v73 行在重建 copy 时**尚不存在**，故不会被自我复制。
- **残留**：`actionRestoreLedger`（`migrate.go:71-75`）**不重放** pending descriptor，此类库的 ledger 形状**只**来自 `identity.go:58` / `:65` 两个字面。→ 这两个字面必须与 v73 同批改为 `TEXT applied_at`，否则恢复路径的库与迁移路径的库形状分叉。

## 4. 待 independent 复审的问题

1. F-1～F-5 的技术结论是否成立、是否有更优处置（例如以显式重建联接表替代子女先行、或改变 descriptor 边界使父表与其全部子表同属一个 descriptor）。
2. 受影响子表跨 descriptor（`user_roles` 不在 v74 清单；`notifications`/`mfa` 在其他 descriptor）时，**descriptor 边界**应如何调整才不残留悬挂 FK。
3. 既有 `rebuildOperationLog` 的 FK 隐患是否为**生产缺陷**、修复应归属哪个目标。
4. F-6 的跨模块列序应以什么为权威证据（live `PRAGMA table_info` 快照 vs 全局版本序推导）。
5. F-7 的 PG 显式 DDL 是否属于 R1 冻结交付，或可与 R2 落码合并。

## 5. 声明

- 本文件 `status: open-finding`；**不**做结构决策、不改变任何 `status` / `progress` / goal-tree。
- 本文件所有"实测"结论均在本机 SQLite 3.51.2 上执行得出；`apps/` 代码本轮**未**修改。
- 逐表 exact rebuild DDL 正文仍待处置方案定案后落盘；本文件是其前置阻塞说明。
