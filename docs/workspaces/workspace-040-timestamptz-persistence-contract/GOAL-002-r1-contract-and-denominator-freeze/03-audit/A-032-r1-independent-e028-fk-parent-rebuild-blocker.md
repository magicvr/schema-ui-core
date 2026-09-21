---
id: A-032-r1-independent-e028-fk-parent-rebuild-blocker
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · ad-hoc design-plan + finding-closure · E-028 / commit 2547ef25 FK 父表重建阻塞（F-1～F-8）对照 A-030 基线 · 不是实施审计 · 候选模式 ≠ 已实施
verdict: conditional
open_required: 6
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-032 · R1 independent · E-028 FK 父表重建阻塞

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：ad-hoc + finding-closure（用户指定先审该发现再定处置：commit `2547ef25` / E-028 的 FK 父表重建阻塞是否成立、结论是否正确、候选处置哪个更合适。对照 A-030，当时 open required = 4。两份附件为设计材料，不是实施证据）。
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / `apps/`。
- 实测：本机 `sqlite3` 3.51.2（`C:\Program Files\msys64\mingw64\bin\sqlite3.exe`），全部在 `:memory:` 或 `%TEMP%\a032-fk-audit\`，未触碰仓库数据。
- 本条**不**复审 A-031 对 F-I-006 谓词方向的修正；F-I-006 维持 open。

## 核对方法

1. 通读 E-028、`r1-c2-fk-parent-rebuild-finding-v1.0-fc.md`、`r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md`；对照 A-030 关闭要求。
2. 独立用 sqlite 3.51.2 复现 A1–A4、F-5 与反模式、官方 12 步、事务内 pragma、转换表达式边界。
3. 对位 `store.go` / `migrate.go` / `identity.go`、`authsession/migration.go`、`operationlog/migration.go`、`account` / `notifications` / `mfa` DDL、descriptor 台账 v74/v80/v81。
4. `git show --stat 2547ef25`：5 个文件均在 workspace-040；`apps/` 无变更。

## 成果（有证据）

1. **本轮未实施 DDL/codec，也未把候选模式写成已落地。** E-028 L32、commit `2547ef25` 的 `--stat` 均无 `apps/**`。finding 附件 `status: open-finding`；mechanism 附件 `status: freeze-candidate`。本审同意该实施边界。
2. **F-1～F-4 的技术事实成立**（本审独立复现，见下节 A）。finding 文档的阻塞结论正确；mechanism §1 表中「`<t>_old` 被 DROP 后引用重新指向同名新表」**不成立**，不得再作为 F-I-002.1 的安全论据。
3. **F-5 TEMP 子女先行在单事务内可行**（本审 1/1 行父 + 1/1 行子保留，子表 `REFERENCES parent(id)` 恢复，`foreign_key_check` 空，`integrity_check=ok`）。把子表 rename 成 `<child>_old` 当数据源会在 `DROP <parent>_old` 时被 CASCADE 删光。
4. **F-6 三列跨模块 ALTER 属实**；F-7 正则派生属实；F-8 无中途读/无自插入属实，但「restore 路径只依赖两个字面」不完整（见 F）。
5. **两条转换表达式在本机 3.51.2 抽测范围内正确**（fixed-6、27 字符、999 ms 无进位、epoch 0、负毫秒 floor、公元 9999 年）。`strftime('%f')` 对 `1758320000.9999` 进位到下一秒，禁用理由成立。
6. **commit `2547ef25` 卫生**：5 文件、workspace-040 only、`apps/` 未改。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C2 物理合同 / 逐表 SQLite rebuild | **仍不可冻结** | F-I-002 仍缺逐表正文；FK 处置未裁决（新 F-I-021） |
| C3 / C4 / R2 放行 | **未满足** | F-I-004/005/006 本轮未触及；A-031 谓词修正未复审 |
| 用户合同忠实 | **E-028 未把候选当实施；用户要求先 /audit 再定方案** | 本条只出意见 |

## 对用户 A–H 的直接判定

### A. 技术事实是否成立

**A1 成立，且是永久改写。** `ALTER TABLE users RENAME TO users_old` 立刻把子表 `sqlite_master.sql` 改成 `REFERENCES "users_old"(id)`；`COMMIT` 后仍在；重新打开同一文件库仍在；`DROP TABLE users_old` 之后子表 DDL **仍然**指向已不存在的 `users_old`。不是事务内瞬时缓存。

**A2 成立。** 事务外 `PRAGMA legacy_alter_table=ON` 且 `foreign_keys=ON` 时，子表仍被改写为 `REFERENCES "users_old"(id)`。补充：`foreign_keys=OFF` **单独**也不能阻止该改写（本审 A2c：OFF 后 rename，子表仍指向 `users_old`）。只有 `foreign_keys=OFF` **且** `legacy_alter_table=ON` 才保持 `REFERENCES users(id)`。故「开 legacy pragma」不是可用解。

**A3 在现行 Apply 合同下成立，但 finding 对 `SetMaxOpenConns(1)` 的引用不精确。**

- 本审实测：事务内 `PRAGMA foreign_keys=OFF` 是 no-op（`PRAGMA foreign_keys` 仍为 1）；事务外才能关掉。
- `applyMigration`（`migrate.go:108-132`）把 `Apply` 与 ledger `INSERT` 包在同一事务；descriptor 内无法关掉 FK。
- DSN `_foreign_keys=on`（`store.go:57`）+ `assertForeignKeysOn`（`migrate.go:32,253-265`）使每个连接默认 ON。
- **纠正**：`SetMaxOpenConns(1)` 只用于 **in-memory**（`store.go:107`）。文件库默认 `sqlitePoolDefault = 4`（`:29,:109-113`）。生产路径是连接池。pragma 是 per-connection：在 pool=4 上对 `db.Exec("PRAGMA foreign_keys=OFF")` 再 `db.Begin()` **不能保证**落在同一连接。官方 12 步若要启用，必须把 `applyMigration` 改成 `Conn` 钉住后再 Begin——这是平台 runner 变更，不是 descriptor 内可做的事。
- 本审另测：官方 12 步（`CREATE parent_new` → copy → `DROP parent` → `RENAME parent_new TO parent`）在 **事务外 FK=OFF** 时可行，子表 `REFERENCES parent(id)` 保持、行保留。同一 12 步在 **FK=ON** 时 `DROP parent` 会 CASCADE 删子行。故「官方 12 步」本身可用，但**前置 FK=OFF 在现行 runner 里不可用**。finding 的「不可用」结论按项目约束成立；不要读成「SQLite 12 步在引擎层无效」。

**A4 成立，分两种子表动作。**

| 子表 | `DROP <parent>_old` 结果（FK=ON） |
|------|-----------------------------------|
| `ON DELETE CASCADE` 且子表有行 | 子表行被级联删光；schema 仍指向 `_old`（本审 dict 样本：1 → 0 行） |
| 默认 NO ACTION / RESTRICT 且子表有行 | `FOREIGN KEY constraint failed`（exit 19） |
| CASCADE 但子表空 | DROP 成功；schema 仍悬挂 `_old` |

朴素「rename → create 同名新表 → copy → DROP `_old`」在子表有行时：CASCADE 路径会先把子行删光，新父表留下、子表 FK 损坏。finding 两支都写对了。

**mechanism §1 有一条错误主张**：声称 `rebuildOperationLog` 对已被 FK 引用的表「能」这样重建，因为「FK 解析按表名进行，`<t>_old` 被 DROP 后引用重新指向同名新表」。本审实测 **DROP 后引用不会回到新同名表**。该句必须从冻结候选中删除。仓库里真正安全的是 `rebuildOperationLogWithSessions`（先 drop 子表），不是裸 `rebuildOperationLog`。

### B. 候选处置模式是否安全

**为什么必须 TEMP 快照、不能 rename 子表当数据源（成立）。** 本审反模式：`ALTER TABLE child RENAME TO child_old` 后，再 `ALTER TABLE parent RENAME TO parent_old`，`child_old` 的 `REFERENCES` 被改写成 `parent_old`；随后 `DROP parent_old` 把 `child_old` 行 CASCADE 删光。TEMP / 无 FK 的 `CREATE TABLE … AS SELECT` 快照没有 `REFERENCES` 子句，父表 DROP 无法级联进去。

**单事务内安全（成立）。** 本审在 `BEGIN`…`COMMIT` 内走完 F-5 十步，提交后 FK 文本正确、行保留。SQLite DDL 可随事务回滚；TEMP 表绑在执行 `Apply` 的那条连接上，与 `applyMigration` 单 `*sql.Tx` 相容。

**推荐（本审建议，P-004 未代选）：采用 F-5（子女先行 + TEMP 快照），不要改 runner 去关 FK，也不要把 v80/v81 并进 v74。**

| 方案 | 判断 |
|------|------|
| **F-5 子女先行**（推荐） | 与已生产的 `rebuildOperationLogWithSessions`（`operationlog/migration/migration.go:709-754`）同构；FK 全程 ON；不改 store runner；失败整事务回滚。代价：必须冻结**完整子表清单**，跨 descriptor 子表用 **live** DDL 原样重建（见 C） |
| 官方 12 步 + runner 在 Begin 前 `foreign_keys=OFF` | 引擎层可行，且可避免 drop 跨 descriptor 子表；但要改 `applyMigration` 为 Conn 钉住、迁移窗口 FK 关闭、与 `assertForeignKeysOn` 对账。平台变更，超出本目标冻结面，除非用户书面选它 |
| 父表+全部子表并进同一 descriptor | `dict_types`/`scheduled_tasks` 已经如此，保持。把 `notifications`/`user_mfa` 并进 v74 会抢走 v81/v80 的时间列转换，破坏已接受 allocation。不推荐 |
| 不 rename、改用别的 DDL | SQLite 不能 `ALTER COLUMN TYPE`；类型变更必须重建。无更简单的第三 DDL |

F-5 的冻结要求（闭合 F-I-021 / 继续写 F-I-002.1 之前必须有）：每张被重建的父表附一张**完整子表清单**（含非本 descriptor、含无时间列的联接表）；子表重建 DDL 的权威是 **v72 已 apply 库的 `sqlite_master.sql` + 索引**，不是模块文件里的原始 `CREATE TABLE`（`notifications` 在 v37 追加了 `title_key`/`body_key`，本审确认 `sqlite_master` 会把 ALTER 列写进 `sql`）。漏一张子表 = 悬挂 FK 或静默丢列。

### C. descriptor 边界

**联接表 DDL 核对（authsession `migration.go`）属实，且不在 v74 的 12 表清单内。**

| 表 | 外键 | 时间列 | v74 清单 |
|----|------|--------|----------|
| `user_roles` L55–58 | `users(id) ON DELETE CASCADE`；`roles(id) ON DELETE RESTRICT` | 无 | 否 |
| `role_permissions` L68–71 | `roles(id) CASCADE`；`permissions(id) RESTRICT` | 无 | 否 |
| `role_menu_items` L83–86 | `roles(id) CASCADE`；`menu_items(id) RESTRICT` | 无 | 否 |

跨模块子表亦属实：`notifications`（v16/`v81`，`users(id) ON DELETE CASCADE`，且 v37 有 `title_key`/`body_key`）；`user_mfa` / `mfa_proofs`（v29/`v80`，`REFERENCES users(id)` **无 CASCADE**——朴素 DROP 会 FK fail 而不是删行）。

**不必把转换 ownership 重划，但 v74 的 Apply 必须 F-5 所有引用被重建父表的子表。** 否则 v74 重建 `users`/`roles`/`permissions`/`menu_items` 后，联接表与 v80/v81 表的 `REFERENCES` 会永久停在 `*_old`。

推荐切法：

1. **v74 转换表范围保持**（时间列仍是那 12 张 / 31 列）。**另外**把 `user_roles` / `role_permissions` / `role_menu_items` 写入 v74 的「FK-preserve identity recreate」集合（无时间转换、最后写者是 v74）。不并进 90 列分母。
2. **`notifications` / `user_mfa` / `mfa_proofs` 留在 v81/v80。** v74 只按 **当前** live DDL（含 `title_key`/`body_key`、INTEGER 时间列）TEMP→DROP→重建父→按 live SQL 建回子表→回填。v80/v81 再重建这些**子表自己**以转换时间列；重建子表（无人引用它们）可用朴素 rename，不需要再 drop `users`。
3. **不要**为了 FK 把 v74+v80+v81 合成一个 descriptor。
4. 同理 v75：`operation_log` 仍用已有 `rebuildOperationLogWithSessions` 舞步处理 `operation_log_correlation` / `operation_log_session`（它们不在 v75 表范围里，但必须被舞到）。v77/v83 父子已在同一 descriptor，F-5 即可。

漏处置的后果：悬挂 `REFERENCES *_old`、`DELETE FROM users` 报 `no such table: main.users_old`、或 CASCADE/FK fail 丢子行。这是 C2 冻结阻断，不是 R2 再发现的细节。

### D. F-6 跨模块列

**属实。** `account` v0013 `enabled`（`account/migration/migration.go:21`）；v0035 `avatar_url`（`:28`）；`notifications` v0017 `notifications_enabled`（`notifications/migration/migration.go:46`）。

finding 写「authsession 内的 6 个 ALTER 列」却列出 7 个名字。实际 authsession ALTER 列是 7 个：`token_version`(v11)、`failed_login_count`+`locked_until`(v12)、`must_change_password`(v38)、`email`+`email_status`(v54)、`last_login_failure_at`(v61)。加上 3 个跨模块列，live `users` 在基线 7 列之后共追加 10 列。

**权威列序证据（推荐，P-004 若有异议再改）：**

1. **主权威**：对一份 **catalog 已 apply 到 v72** 的 SQLite 库做 `PRAGMA table_info(users)`（测试夹具 / 新开 `OpenWithCatalog` 即可，不要用生产库）。SQLite 把 ALTER 列追加在末尾；这就是重建 `CREATE TABLE` / `INSERT SELECT` 必须对齐的序。
2. **交叉核对**：全局版本序推导 `v1 → v11 → v12 → v13 → v17 → v35 → v38 → v54 → v61`。与 live 快照不一致则 fail closed。
3. **禁止**按模块文件分组拼列序。

遗漏后果：`CREATE` 缺列 → `DROP users_old` 后该列数据永久消失；或 `INSERT SELECT` 不列该列 → 同样丢数据且无报错。`enabled` / `avatar_url` / `notifications_enabled` 均非时间列，但必须出现在 v74 的 `users` 新 DDL 里。

### E. F-7 PG 正则派生

**属实。** `pgTimeColRe`（`operationlog/migration/migration.go:250`）为 `(?m)^(\s*)(created_at|archived_at)\s+INTEGER NOT NULL`；`pgTimeDDL` 只把命中处换成 `BIGINT NOT NULL`。SQLite 字面改成 `TEXT` 后 Replace 为 **静默 no-op**，PG 侧继续用陈旧 BIGINT 映射且不报错。`pgRebuild`（`:296-301`）把派生结果送进同一套 rename 重建。

**该显式 PG DDL 属于 R1 冻结交付，不要拖到「R2 落码时再写」。** 它是 F-I-002 的 PG USING/DDL 正文的一部分；R2 只负责按冻结正文落码。D-017 单 checksum（PG 不进哈希）仍然成立——显式 PG DDL 用测试断言，不进 `MigrationChecksum`。v1–v72 的既有 SQLite 字面**不得**改；F-7 只约束 **v75 新切片**不得再走 `pgTimeDDL`。

### F. F-8 `schema_migrations`

**自引用/中途读：无危害，finding 成立。** ledger 在 `applyPending` 前读（`identity.go:334`、`migrate.go:63`）；循环内与 `applyMigration` 内不再读；`verifyIntegrity` 在批次后（`migrate.go:102`）。`applyMigration` 先 `Apply` 再 `INSERT`（`:113-123`）；v73 行在 copy 时尚不存在。

**「restore 只依赖 `identity.go:58/:65`」不完整。** `actionRestoreLedger` 的 **CREATE** 确实只用那两个字面（`:365-371` / PG `:390-396`），但写入 `applied_at` 的还有：

- `stampCatalog`（`identity.go:317`）`time.Now().UTC().Unix()`
- `applyMigration`（`migrate.go:122-123`）同样 Unix 秒
- PG `applyMigrationPG`（`postgres.go:166-167`）同样

v73 若把列改成 `TEXT` 而不改这三处写入，restore 库与 migrate 库形状分叉，且 v73 自己的 ledger 行会把整数写进 TEXT 列（SQLite 不报错，合同被静默破坏）。**必须与 v73 同批**：DDL 字面 + 三处 INSERT 改为 fixed-6 RFC3339（PG `timestamptz(6)`）。这**不是**改 v1 checksum：`identity.go` 的 restore 字面不在 `0001:r2-baseline` 的 `r2BaselineDDL` 哈希输入里；authsession `schemaMigrationsDDL`（`migration.go:18-23`）属于 v1，**禁止**改。

### G. 既有 `rebuildOperationLog` 生产缺陷

**函数层隐患成立；不等于「当前已 apply 到 v72 的库 schema 已被损坏」。**

- 裸 `rebuildOperationLog`（`:756-776`）是 rename-first 四步，与本审已证的损坏模式相同。
- 历史调用 `0005–0036` 发生在 `0041` correlation **之前**，当时没有子表，fresh bootstrap 路径安全。
- `0043/0045` 走 `WithCorrelation`；`0053/0071` 走 `WithSessions`。现行库在 0071 之后 FK 文本应仍指向 `operation_log`。
- 隐患是：**未来**任何在子表存在时调用裸 `rebuildOperationLog`（含 v75、含误用 `pgRebuild`）会损坏 FK。mechanism §1 把该函数当成「已被引用也能重建」的先例，是错的。
- `wallet` 的 rename 重建（`wallet/migration.go:182,356`）当前 **没有** `REFERENCES wallet_*` 子表，同类炸弹尚未上膛，但 v85 仍须先做子表盘点。
- `internal/store` **没有**同名 rebuild 函数；约束是 runner 事务边界，不是第二份损坏实现。

**最小安全改动（不改 0001–0072 canonical SQL / checksum）：**

1. `rebuildOperationLog`：在 rename 前 fail-closed 断言 `operation_log_correlation` / `operation_log_session` **不存在**（`sqlite_master` / PG `pg_class`）。`0005–0036` 仍走原路径；`WithCorrelation`/`WithSessions` 先 drop 再调用，断言仍通过。
2. 不要把 dance 内嵌进 `rebuildOperationLog`（会与 `WithSessions` 双重 drop）。
3. v75 必须调用 `rebuildOperationLogWithSessions`（或与 F-5 同一通用 helper），禁止 `pgRebuild`。
4. **禁止**改 `0004–0071` 的 DDL 切片字面（checksum 输入）。改 Go 控制流 / 断言 **不**进入 `MigrationChecksum`（哈希的是 stmts 切片，不是 Apply 函数体）。因此 **0001..0072 canonical SQL/checksum 不变性不受影响**——前提是不改那些切片、不改 transform_id。

**范围：属于 VP-040，不要另开 VP/目标。** 用户已书面「本轮一并修复」；v75 重建 `operation_log` 是本 VP 分母内工作；helper 断言是其安全前置。设计进 GOAL-002 冻结包，代码进 R2。另开目标只会把 C2 冻结再挡一道门。

### H. 转换表达式

**在抽测边界内正确，满足 fixed-6 与向零/floor 语义。** 本机 3.51.2 复现与附件 §2 声明一致，全部 27 字符：

| 输入 | 族 | 输出 |
|------|----|------|
| `0` | 秒/毫秒 | `1970-01-01T00:00:00.000000Z` |
| `1758320000` | 秒 | `2025-09-19T22:13:20.000000Z` |
| `1758320000999` | 毫秒 | `2025-09-19T22:13:20.999000Z`（无进位） |
| `253402300799` | 秒 | `9999-12-31T23:59:59.000000Z` |
| `-1` ms | 毫秒 | `1969-12-31T23:59:59.999000Z` |
| `-999` / `-1000` / `-1001` ms | 毫秒 | `…59.001000Z` / `…59.000000Z` / `…58.999000Z` |
| `-1758320000123` | 毫秒 | `1914-04-14T01:46:39.877000Z` |

`strftime('%Y-%m-%dT%H:%M:%f', 1758320000.9999, 'unixepoch')` → `2025-09-19T22:13:21.000`（round 进秒）。`/1000.0` 禁用理由成立。SQLite 整数 `/` 向零、`%` 跟随被除数：`-1%1000=-1`，余数必须 `(col%1000+1000)%1000`。负值秒部分必须 `(col-999)/1000` 才能得到 floor（与 Go `UnixMilli` 一致）。

生产路径上 `< 0` 由 `m0` 预检 fail closed、不进表达式；表达式自身对负值仍正确。本轮未见反例。R2 仍须用 `T-<#>-RT` 在目标 SQLite 上复证并对拍 Go codec（D-018）。

## Findings

### F-I-001 · 90 列 + catalog 72 + v1–v72 扫描

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-002 · 逐列 USING/rebuild 仍不足以为 C2 冻结

- **严重度**：high · **建议**：required
- **状态**：open（维持；**本轮收窄表达式层，不关闭**）
- **影响门禁**：C2/C3、R2；关联 `I-040-001`
- **本轮已修 / 收窄**：
  1. 秒族 / 毫秒族纯整数表达式本审复证通过（H）；`%f` / `/1000.0` 禁用成立。
  2. `#34` / `#72` / `#73` 在 mechanism §3 已收到与模板同一式（D0 去 NN；voucher `NULL OR =0`；CASE 无 `< 0`）。这是机制层收口，**不是**逐表 `CREATE`/`INSERT` 正文。
- **仍不闭合**：
  1. **逐表 exact SQLite rebuild DDL 仍未落盘**（E-028 L27 自承：已抽取、因 FK 方案未定而未附件化）。
  2. **FK 父表处置未冻结**（升级为 F-I-021；无处置则逐表正文无法写完）。
  3. mechanism §1 仍含错误的「DROP 后引用回到同名新表」句，不能当 F-I-002.1 证据。
  4. 用例仍是 ID；非法/越界可执行测试未发生。
- **关闭要求**：在 F-I-021 处置落盘后，写出每张受影响表的 exact `CREATE` / `INSERT SELECT` / 索引；父表走 F-5（或用户另选并留痕的等价）；PG 侧显式 DDL（见 E），不用 `pgTimeColRe`。freeze-candidate ≠ 实施。

### F-I-003 · 90 列 mapping

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-004 · Backup Port ≠ 可执行备份/回滚

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮未触及 C3）

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮无新哈希/测试）
- **本轮相关**：G 的 helper 修复**不得**改 0001–0072 stmts；F-7 显式 PG DDL 不进 checksum（D-017）。候选 descriptor 名 ≠ 已记录哈希。

### F-I-006 · CHECK / 部分索引 / 谓词列表仍未冻结

- **严重度**：med · **建议**：required · **状态**：open（维持）
- **本轮**：A-031 自称已改 `#5` 方向、readwrite superseded、ORDER BY 归属、jobs 索引、`#6` callsite。**本条 scope 不复审那些修正**；不得据此闭合。

### F-I-007 … F-I-020

- 维持 A-030 状态。F-I-018 / F-I-019 **closed** 维持。F-I-020 recommended 仍 open（本轮未核 A-031 卫生修正）。F-I-008 / F-I-009 recommended open 维持。

### F-I-021 · FK 父表重建阻塞：处置与子表清单未冻结

- **严重度**：high · **建议**：required · **状态**：open（本轮新增）
- **影响门禁**：C2 冻结、R2 实施；关联 `I-040-001`、`I-040-003`
- **描述**：E-028 / finding F-1～F-5 经独立实测成立。朴素 rename-first 会永久改写子表 `REFERENCES` 并 CASCADE 删子行或 FK fail。`legacy_alter_table` 不能救。现行 `applyMigration` 事务内不能 `foreign_keys=OFF`。F-5 候选可行但**未定案**；跨 descriptor 子表（`user_roles`/`role_permissions`/`role_menu_items`/`notifications`/`user_mfa`/`mfa_proofs`）未写入冻结清单；mechanism §1 的「DROP 后引用回同名表」为假。
- **关闭要求**（须用户书面处置，本审**不**代选，仅给建议）：
  1. P-004 选定重建模式（本审建议 F-5；备选 = runner 级官方 12 步）。
  2. 每张父表冻结完整子表清单 + live DDL 权威（v72 `sqlite_master`，含 `notifications.title_key/body_key`）。
  3. 写明 v74 Apply 可 F-5 非转换表，而不把 v80/v81 列抢进 v74。
  4. 删除/改写 mechanism §1 错误句与 D-018 中「含被 FK 引用的表也能裸重建」的理由句。
  5. 既有 `rebuildOperationLog` 的 fail-closed 断言进入冻结包（实施在 R2，见 G）。

### F-I-022 · v73 ledger 写入路径不只 identity.go 两个字面

- **严重度**：med · **建议**：required · **状态**：open（本轮新增，收窄 F-8）
- **影响门禁**：C2 中 v73 设计、R2 落码；关联 `I-040-001`
- **描述**：restore 的 CREATE 确为 `identity.go:58/:65`，但 `stampCatalog`、`applyMigration`、`applyMigrationPG` 仍写入 Unix 秒整数。只改两个 CREATE 字面会导致 restore/migrate 分叉，且 v73 行本身不是 RFC3339。
- **关闭要求**：v73 设计列出 CREATE 字面 **与** 三处 INSERT 的目标类型/格式；authsession v1 `schemaMigrationsDDL` 保持不变。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002 | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec | **收窄**：表达式复证、`#34/#72/#73` 机制层同一。**仍缺**逐表 DDL、FK 处置 |
| F-I-004 | C3、R2/R3 | 不得把 Port/runbook 当 C3 冻结 | **无新收窄** |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL | **无新收窄**；G 修复不得碰 0001–0072 stmts |
| F-I-006 | C2/C3、R2 | 不得在 exact old/new 未冻时 table-rebuild | **本轮未复审 A-031** |
| **F-I-021** | C2、R2 | 不得写完/接受逐表 rebuild DDL | **新增**；A1–A4 成立；待 P-004 处置 |
| **F-I-022** | C2 v73、R2 | 不得只改 restore CREATE 字面 | **新增**；F-8 写入路径不完整 |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、F-I-018、F-I-019 为 closed。F-I-008、F-I-009、F-I-020 为 recommended open。

**开放 required = 6**（F-I-002、F-I-004、F-I-005、F-I-006、**F-I-021**、**F-I-022**）。在这些合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-030 independent | E-028 / finding 自称 | A-032 independent（本条） |
|----|-------------------|----------------------|---------------------------|
| verdict | conditional | 未自证闭合；待本审 | **conditional** |
| F-I-002 | open；缺 exact SQLite rebuild | 表达式已测；DDL 因 FK 未落盘 | **维持 open**；表达式子项收窄；FK 阻断确认 |
| F-I-003/018/019 | closed | 未重开 | **维持 closed** |
| F-I-004/005/006 | open | 未声称触及 | **维持 open** |
| 新 required | 无 | F-1～F-8 待审 | **F-I-021、F-I-022** |
| mechanism §1「DROP 后引用回同名表」 | 当时未审 | finding 已否定；mechanism 仍写着 | **否定该句**；与 finding 一致、与 mechanism 冲突 |
| R2 | 禁止 | 禁止 | **禁止** |

无合同方向上的「一要一否」需要用户在 finding vs 本审之间裁。**需要 P-004 的是 F-I-021 的重建模式**（本审建议 F-5；不在本条选定）。D-018 行拷贝选项 C / 无过渡期窗口维持；其「裸 rebuild 也能处理被引用表」的理由句应在响应里划掉，不推翻选项 C 本身。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 表达式可核；逐表 DDL 与 FK 处置未闭（F-I-002/021/022） |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列分母仍可核对；不因联接表无时间列而扩大 90 列 |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放；FK 重建失败模式是新的转换失败面，归 F-I-021 |
| I-040-004 | required | R3 | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-001/003 仍开放，阻断 C2/C3/R2。无用户书面 residual。F-I-002 **仍不能闭合**：缺逐表正文、缺已裁决的 FK 处置与子表清单、mechanism §1 未改写。

## 结论 + 建议给编排器/用户的下一步

**conditional。** E-028 的 FK 阻塞发现成立；不是夸大。A-030 的 4 条 required 仍开放，并新增 F-I-021 / F-I-022。不足以冻结 C2/C3 或放行 R2。

建议 `/govern`：

1. 响应本 A-032；**不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
2. **P-004 选定 FK 重建模式**（建议 F-5 + 完整子表清单 + live `sqlite_master` 权威；不要合并 v74/v80/v81；不要在未改 runner 的前提下宣称官方 12 步可用）。
3. 接受 F-I-021 / F-I-022 为新 required；维持 F-I-002/004/005/006 open；**不要**把 mechanism 附件或 TEMP 舞步标成 F-I-002 closed。
4. 改写 mechanism §1 错误句与 D-018 过时理由；把 `rebuildOperationLog` fail-closed 断言与 v73 三处 INSERT 写入冻结包。
5. 处置落盘后再写逐表 exact DDL，然后才能再 `/audit` 看 F-I-002。
6. A-031 的 F-I-006 修正另安排复审，不要搭本条顺手闭合。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree / `apps/`。响应、finding 闭合与是否推进由 `/govern` 处理。
