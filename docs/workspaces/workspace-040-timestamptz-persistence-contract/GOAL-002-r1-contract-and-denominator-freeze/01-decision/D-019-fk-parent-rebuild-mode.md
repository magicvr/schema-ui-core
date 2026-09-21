---
id: D-019-fk-parent-rebuild-mode
doc: decision-entry
status: accepted
parent: GOAL-001-timestamptz-persistence-contract
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-019 · FK 父表重建模式与子表处置（用户 P-004 裁决）

## 决定的来源

- A-032（independent · grok-build grok-4.6 · high）独立复现了 E-028 的 F-1～F-4：重命名 FK 父表会**永久**改写子表 `REFERENCES`；`legacy_alter_table=ON` 不能救；事务内 `foreign_keys=OFF` 是 no-op；朴素重建会 CASCADE 删子行或 `DROP` 报 FK fail。A-032 新增 required **F-I-021**，并明确「须用户书面处置，本审**不**代选，仅给建议 = F-5」。
- 用户 2026-09-20 经 P-004 裁决两项：**模式 = 选项 A（F-5 子女先行 + TEMP 快照）**；**跨 descriptor 子表 = 选项 A（v74 先修 FK，v80/v81 再转类型）**。

## 决定

### 1. 重建模式 = F-5 子女先行 + TEMP 快照

对**任何有 FK 子表的父表**，descriptor 内按以下顺序（单事务）：

```text
1. 对每张子表：CREATE TEMP TABLE <child>_backup AS SELECT <列> FROM <child>;   -- 无 FK 子句
2. 对每张子表：DROP TABLE <child>;
3. ALTER TABLE <parent> RENAME TO <parent>_old;
4. CREATE TABLE <parent> ( ...目标形状... );
5. INSERT INTO <parent> (...) SELECT <转换表达式> FROM <parent>_old;
6. DROP TABLE <parent>_old;
7. 对每张子表：CREATE TABLE <child> ( ...live DDL 逐字... );
8. 对每张子表：INSERT INTO <child> (...) SELECT ... FROM <child>_backup;
9. 对每张子表：DROP TABLE <child>_backup;
10. CREATE INDEX ...（全部索引，含父表与子表，最后建）
```

- **无 FK 子表的表**仍可用裸四步（`ALTER … RENAME TO <t>_old` → CREATE → INSERT…SELECT → DROP → CREATE INDEX）。
- **禁止**把子表 rename 成 `<child>_old` 当数据源（父表 rename 会把该表的 FK 也改写，随后父表 DROP 会 CASCADE 删掉待拷数据）。
- **禁止**依赖 `PRAGMA foreign_keys=OFF` 或 `legacy_alter_table=ON`（实测均不可用/无效）。
- **不采用** runner 级官方 12 步（需改 `applyMigration` 为 Conn 钉住，属平台变更，超出本目标冻结面）。

### 2. 每张父表的**完整子表清单**（F-I-021 关闭要求 2；权威 = v72 已 apply 库的 `sqlite_master`）

| 父表 | descriptor | 子表 | 子表 FK 动作 | 子表是否含时间列 | 子表归属 |
|------|-----------|------|--------------|------------------|----------|
| `users` | v74 | `refresh_tokens`, `email_verification_challenges`, `password_recovery_challenges`, `login_failures`, `user_password_history`, `user_invites` | CASCADE / 部分无 | **是** → v74 一并转换 | 本 descriptor |
| `users` | v74 | **`user_roles`** | `users` CASCADE、`roles` RESTRICT | 否 | **本 descriptor 的 FK-preserve 集合** |
| `users` | v74 | `notifications`（v81 内） | CASCADE | 是（`read_at`/`created_at`） | 跨 descriptor |
| `users` | v74 | `user_mfa`, `mfa_proofs`（v80 内） | **无 CASCADE** → 不 drop 会 FK fail | 是 | 跨 descriptor |
| `roles` | v74 | **`user_roles`**（RESTRICT）、**`role_permissions`**（CASCADE）、**`role_menu_items`**（CASCADE） | 见左 | 否 | **本 descriptor 的 FK-preserve 集合** |
| `permissions` | v74 | **`role_permissions`**（RESTRICT） | 见左 | 否 | 同上 |
| `menu_items` | v74 | **`role_menu_items`**（RESTRICT） | 见左 | 否 | 同上 |
| `operation_log` | v75 | `operation_log_correlation`, `operation_log_session` | CASCADE | 否 | 同模块旁表 |
| `dict_types` | v77 | `dict_entries` | CASCADE | 是 | 本 descriptor |
| `scheduled_tasks` | v83 | `task_runs` | CASCADE | 是 | 本 descriptor |

- **v74 的表清单**：时间列转换仍为原 **12 张 / 31 列**（`D-014` 已接受 allocation 不变）；**另加** `user_roles` / `role_permissions` / `role_menu_items` 三张**无时间列**的联接表进入「FK-preserve identity recreate」集合。三张联接表**不进 90 列分母**。
- **`wallet_*`（v85）目前没有 `REFERENCES wallet_*` 子表**，但 v85 落码前仍须先做子表盘点（A-032 §G）。

### 3. 跨 descriptor 子表处置（用户裁决 2 = 选项 A：两次重建）

| 子表 | 步骤 1（v74） | 步骤 2（其原 descriptor） |
|------|---------------|---------------------------|
| `notifications` | F-5：TEMP 快照 → DROP → 重建父表 → 按 **live DDL**（**含 v37 的 `title_key`/`body_key`**）建回 → 回填；时间列仍是 **INTEGER** | v81 对该表做时间列类型转换（此时无人引用它，可裸 rename 四步） |
| `user_mfa`, `mfa_proofs` | 同上（live DDL，时间列仍 INTEGER） | v80 同 |

- **权威 DDL 来源** = 一份 catalog 已 apply 到 **v72** 的 SQLite 库的 `sqlite_master.sql` + 索引清单（测试夹具 / `OpenWithCatalog`，**不用生产库**）；**不是**模块文件里的原始 `CREATE TABLE`（`notifications` 的 `title_key`/`body_key` 由 v37 ALTER 追加，只有 `sqlite_master` 才反映）。
- 明确**不**把 v80/v81 的时间列转换并进 v74（会抢走 ownership、破坏已接受 allocation）。代价：这三张表会被重建两次（先修 FK、后转类型），已知且接受。
- **禁止**为 FK 把 v74+v80+v81 合成一个 descriptor。

### 4. v73 ledger 写入路径（F-I-022 关闭要求）

v73 设计必须同时列出 **DDL 字面**与**四处 `applied_at` 写入路径**的目标类型/格式：

| 位置 | 现状 | 目标 |
|------|------|------|
| `internal/store/identity.go:58`（`sqliteLedgerDDL`，restore 用 CREATE） | `applied_at INTEGER NOT NULL` | `TEXT NOT NULL` |
| `internal/store/identity.go:65`（`postgresLedgerDDL`，同） | `applied_at BIGINT NOT NULL` | `timestamptz(6) NOT NULL` |
| `internal/store/identity.go:317`（`stampCatalog` INSERT） | `time.Now().UTC().Unix()` | fixed-6 UTC RFC3339（PG `timestamptz(6)`） |
| `internal/store/migrate.go:122-123`（`applyMigration` INSERT） | 同上 | 同上 |
| `internal/store/postgres.go:166-167`（`applyMigrationPG` INSERT） | 同上 | 同上 |

- **禁止改** authsession `schemaMigrationsDDL`（`authsession/migration/migration.go:18-23`）——它属 **v1**。
- **不影响 v1 checksum**：`identity.go` 的两个 restore 字面**不在** `0001:r2-baseline` 的 `r2BaselineDDL` 哈希输入内。

### 5. 既有 `rebuildOperationLog` 的最小安全改动（用户 2026-09-20 裁决「本轮一并修复」+ A-032 §G）

- **函数层隐患成立**，但**当前 v72 库未必已损坏**（`0043/0045` 走 `WithCorrelation`、`0053/0071` 走 `WithSessions`；`0005–0036` 调用时无子表）。
- **最小改动**：在 `rebuildOperationLog`（`operationlog/migration/migration.go:756-776`）rename **之前**加 fail-closed 断言——`operation_log_correlation` / `operation_log_session` **不存在**（查 `sqlite_master` / PG `pg_class`）；`WithCorrelation`/`WithSessions` 先 drop 再调用，断言仍通过。
- **不要**把 dance 内嵌进 `rebuildOperationLog`（会与 `WithSessions` 双重 drop）。
- **v75 必须调用 `rebuildOperationLogWithSessions`**（或与 F-5 同一通用 helper），**禁止** `pgRebuild`。
- **禁止**改 `0004–0071` 的 DDL 切片字面（checksum 输入）。改 Go 控制流/断言**不进入** `MigrationChecksum`（哈希的是 `stmts` 切片，不是 `Apply` 函数体）→ **0001–0072 canonical SQL/checksum 不变性不受影响**。
- **范围**：设计进 GOAL-002 冻结包（R1），**代码进 R2**。**不另开 VP/目标**（用户与 A-032 §G 一致）。

### 6. PG 侧显式 DDL（F-I-002 的 E 项）

- `operationlog/migration/migration.go:250` 的 `pgTimeColRe`（`(created_at|archived_at)\s+INTEGER NOT NULL`）在 SQLite 字面改 `TEXT` 后**静默 no-op**。
- → v75 起**禁止**再走 `pgTimeDDL` 派生；PG DDL 必须**显式书写**，并作为 **R1 冻结交付**（不拖到 R2）。
- `D-017` 单 checksum 约定不变：PG 显式 DDL **不进** `MigrationChecksum`，由测试断言覆盖。

## 未选方案

- **runner 级官方 12 步**：引擎层可行，但需改 `applyMigration` 为 `Conn` 钉住 + 迁移窗口全局关 FK + 与 `assertForeignKeysOn`（`migrate.go:32,253-265`）及 pool=4（`store.go:29,109-113`）对账。属平台变更，超出本目标冻结面。用户未选。
- **父表 + 全部子表并进同一 descriptor**：把 `notifications`/`user_mfa`/`mfa_proofs` 并进 v74 会抢走 v81/v80 的时间列转换、破坏已接受 allocation。用户未选（改选两次重建）。
- **拆出独立子目标承载 FK 父表族**：A-032 明确不建议（会把 C2 冻结再挡一道门）。用户未选。

## 影响与边界

- 本决策**落盘即满足 F-I-021 关闭要求的第 1/3/4 项**（P-004 选定模式、写明 v74 可 F-5 非转换表、删除机制附件错误句——后者已在 `76d829d6` 完成）；第 2 项（冻结完整子表清单 + live DDL 权威）以本文件 §2/§3 为载体，**是否接受仍由 independent 复审判定**；第 5 项（`rebuildOperationLog` 断言进冻结包）见 §5。
- 本决策**不闭合 F-I-002**：逐表 exact `CREATE` / `INSERT SELECT` / 索引正文仍须在下一轮落盘。
- 本决策**不闭合 F-I-022**：v73 的 DDL 与写入路径正文仍须落盘（§4 已定目标形态）。
