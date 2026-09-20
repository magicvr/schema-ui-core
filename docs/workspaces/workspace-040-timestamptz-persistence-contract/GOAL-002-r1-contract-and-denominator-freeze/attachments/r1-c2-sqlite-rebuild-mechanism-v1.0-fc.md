---
id: r1-c2-sqlite-rebuild-mechanism-v1.0-fc
doc_type: design-attachment
title: R1 C2 SQLite 表重建机制与 fixed-6 转换表达式 v1.0（冻结候选）
status: freeze-candidate
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 1.0.0
---

# R1 C2 SQLite 表重建机制与 fixed-6 转换表达式 v1.0（冻结候选）

> **状态：`freeze-candidate`。** 本文件承载 A-030 **F-I-002.1**（逐表 exact SQLite rebuild DDL）所需的**机制与表达式层**：重建路径、免 pragma 的可行性论证、两条经实测验证的整数转换表达式、以及 `#34`/`#72/#73` 特化后的唯一形态。逐表 `CREATE TABLE` / `CREATE INDEX` / `INSERT … SELECT` 正文在配套的逐表附件中，二者共同构成 F-I-002.1 的交付物。

## 1. 重建机制（无 pragma 可用；**仅对无 FK 子表的表**可裸用四步）

> **⚠️ 本节已按 A-032 更正（2026-09-20）。** v1.0.0 初稿在「被其他表 FK 引用的表能否这样重建」一行断言为「**能**」，理由是「`<t>_old` 被 DROP 后引用重新指向同名新表」。**该断言经独立审计实测否定，是错的**，已删除。正确结论见下：裸四步**只对没有 FK 子表的表安全**；有 FK 子表的父表必须走 **子女先行（F-5）**，见 `r1-c2-fk-parent-rebuild-finding-v1.0-fc.md` §F-5 / §F-5.1。

**裸四步（适用于无 FK 子表的表）**（与 `apps/api/modules/wallet/migration/migration.go:180-201` 同构）：

```text
1. ALTER TABLE "<t>" RENAME TO "<t>_old";
2. CREATE TABLE "<t>" ( ...目标形状... );
3. INSERT INTO "<t>" (<列清单>) SELECT <列清单/转换表达式> FROM "<t>_old";
4. DROP TABLE "<t>_old";
5. CREATE INDEX / CREATE UNIQUE INDEX ...（逐字重建，必须在 4 之后）
```

**先决条件（A-032 新增）= 该表没有任何存留 DDL 含 `REFERENCES <t>` 的子表。** 若不满足，必须改用 F-5 子女先行模式。

**为什么（更正后的事实）**：

| 关注点 | 事实（实测 / A-032 复现） | 结论 |
|--------|---------------------------|------|
| 经典 SQLite rebuild 要求 `PRAGMA foreign_keys=OFF` | DSN 固定 `_foreign_keys=on`（`apps/api/internal/store/store.go:57`）；`assertForeignKeysOn` 启动强制 ON（`migrate.go:32,253-265`）；`PRAGMA foreign_keys` 在**事务内是 no-op** | **不采用** `foreign_keys=OFF` 路径 |
| 连接池事实（**更正**） | `SetMaxOpenConns(1)` **只用于 in-memory**（`store.go:107`）；文件库默认 `sqlitePoolDefault = 4`（`:29,:109-113`）。pragma 是 per-connection，pool=4 上 `db.Exec("PRAGMA foreign_keys=OFF")` 后再 `db.Begin()` **不保证同连接** | 官方 12 步若要启用须把 `applyMigration` 改为 `Conn` 钉住后再 Begin——平台 runner 变更，**不是** descriptor 内可做的事 |
| `legacy_alter_table=ON` 能否救 | 事务外设 `legacy_alter_table=ON` + `foreign_keys=ON` 时子表**仍被改写**；`foreign_keys=OFF` **单独**也不能阻止改写；只有 `foreign_keys=OFF` **且** `legacy_alter_table=ON` 才保持 `REFERENCES <t>(id)`（A-032 的 A2c） | 不是可用解 |
| 裸四步用于有 FK 子表的父表 | 子表 DDL 被**永久**改写为 `REFERENCES "<t>_old"(id)`（重连、DROP 后仍在）；`DROP <t>_old` 时 CASCADE 子表会**丢行**、NO ACTION/RESTRICT 子表会 `FOREIGN KEY constraint failed` | **禁止**裸用；须 F-5 |
| 无 FK 子表的表 | `wallet/migration.go:182,356` 的既有 rename 重建当前**没有** `REFERENCES wallet_*` 子表，故未受损（但 v85 仍须先做子表盘点） | 裸四步可用 |
| 迁移事务边界 | `applyMigration` 一个 descriptor 一个事务，`Apply(sqlTx{tx})` 执行全部语句，随后同事务写 ledger（`migrate.go:108-132`）；F-5 的 TEMP 表绑在该连接/事务上，相容 | 全部步骤 + 索引 + ledger 在同一事务内，**失败即整体回滚** |
| 重建后校验 | runner 已在批次后跑 `PRAGMA foreign_key_check`（`migrate.go:349`）与 `integrity_check`（`:367`） | **复用**；逐表用例再加 `P-*-ORD` 与「子表 FK 文本解析回同名父表」断言 |

**逐表必须附带的校验步骤**（进 owner descriptor 的 `m4`）：`PRAGMA foreign_key_check` 无行、`PRAGMA integrity_check` = `ok`、`sqlite_master` 中 `<t>_old` 不存在、目标列 `PRAGMA table_info` 类型为 `TEXT`；**父表另加**：全部子表已重建且其 `sqlite_master.sql` 中 `REFERENCES` 指向同名新父表（不含 `_old`）。

## 2. 唯一转换表达式（实测验证，纯整数，禁用 `%f`）

### 2.1 为什么禁用 `strftime('%f', …)` 与 `/1000.0`

本机 SQLite 3.51.2 实测：

| 输入 | `strftime('%f', …)` 结果 | 问题 |
|------|--------------------------|------|
| `1758320000` | `20.000` | **变宽**：整秒只有 3 位小数 |
| `1758320000.123` | `20.123` | 3 位 |
| `1758320000.9999` | `22:13:21`（进位到下一秒） | **round，不是 truncate** |

结论：`%f` 既非定宽也不向零截断，**与 fixed-6 合同和 Root `D-015` 的截断规则冲突**，全包禁用；`/1000.0` 亦禁用（会引入二进制浮点与同一 round 问题）。

### 2.2 秒族（exact）

```sql
strftime('%Y-%m-%dT%H:%M:%S', <col>, 'unixepoch') || '.000000Z'
```

- 整数秒无小数部分，`%S` 恒为 2 位，拼接恒为 **27 字符**。
- 实测：`1758320000` → `2025-09-19T22:13:20.000000Z`；`-1` → `1969-12-31T23:59:59.000000Z`；`-86400` → `1969-12-31T00:00:00.000000Z`；`0` → `1970-01-01T00:00:00.000000Z`；`253402300799` → `9999-12-31T23:59:59.000000Z`。

### 2.3 毫秒族（exact，负值安全）

```sql
strftime('%Y-%m-%dT%H:%M:%S',
         CASE WHEN <col> >= 0 THEN <col>/1000 ELSE (<col>-999)/1000 END,
         'unixepoch')
  || '.' || printf('%03d', (<col>%1000 + 1000) % 1000)
  || '000Z'
```

**两处都必须按上面写，任何一处简化都会在负值上出错**（本机 SQLite 3.51.2 实测，R1 冻结前复证）：

| 事实 | 实测 | 后果 |
|------|------|------|
| SQLite 整数 `/` 是**向零截断**（`-86400000/1000 = -86400`），Go 整数除法同样向零截断 | `SELECT -86400000/1000` → `-86400` | 秒部分需要 **floor** 语义（Go `time.UnixMilli` 为 floor），故负值必须用 `(<col>-999)/1000` |
| SQLite 整数 `%` 的**符号跟随被除数**（`-1%1000 = -1`） | `SELECT -1%1000` → `-1` | 毫秒余数必须用 `(<col>%1000 + 1000) % 1000` 归一化，否则 `printf('%03d', -1)` 产出 `-01` 并撑到 28 字符 |

> **本候选的修正记录**：v1.0.0 初稿曾把秒部分写成 `CASE WHEN <col> >= 0 THEN <col>/1000 ELSE (<col>-999)/1000 END`（正确）但把余数写成 `CASE WHEN <col> >= 0 THEN <col>%1000 ELSE <col>%1000 + 1000 END`（**错误**：`-1` 时得 `999` 而非 `999`，`-1001` 时得 `-1`，实测产出 `...T00:00:00.-01000Z`）。经整表重建探针发现并改为上式；**修正后的整表探针结果见 §6**。

- 全部为整数算术，无二进制浮点。
- 实测（修正后）：`1758320000123` → `2025-09-19T22:13:20.123000Z`；`1758320000999` → `2025-09-19T22:13:20.999000Z`（**无进位**）；`-1` → `1969-12-31T23:59:59.999000Z`；`-999` → `1969-12-31T23:59:59.001000Z`；`-1000` → `1969-12-31T23:59:59.000000Z`；`-1001` → `1969-12-31T23:59:58.999000Z`；`-86400000` → `1969-12-31T00:00:00.000000Z`；`-1758320000123` → `1914-04-14T01:46:39.877000Z`；`0` → `1970-01-01T00:00:00.000000Z`；`1758316799999` → `2025-09-19T21:19:59.999000Z`。**全部 27 字符。**

### 2.4 空值与 sentinel 包裹

| 列类别 | `INSERT … SELECT` 中的唯一形态 |
|--------|-------------------------------|
| `NN`（NOT NULL，无 sentinel） | 直接使用 §2.2 / §2.3 表达式 |
| `N`（可空，NULL = 缺失） | `CASE WHEN <col> IS NULL THEN NULL ELSE <表达式> END` |
| `D0`（`NOT NULL DEFAULT 0`，0 = 缺失） | `CASE WHEN <col> = 0 THEN NULL ELSE <表达式> END`；新 DDL **去 `NOT NULL`、去 `DEFAULT 0`** |
| `#72/#73 vouchers`（可空且 legacy `0` 亦为缺失） | `CASE WHEN <col> IS NULL OR <col> = 0 THEN NULL ELSE <表达式> END` |

**负值 `< 0` 一律不进表达式**（不进 `USING`/rebuild 的任何 `CASE` 分支）。但**政策按列分档**（经 A-042 更正，与 Root `D-012`/`D-015` 一致）：

| 列 | 负值处置 | 依据 |
|----|----------|------|
| **仅** `vouchers.expires_at` / `vouchers.redeemed_at`（`#72`/`#73`） | `m0` 预检 **fail closed**（数据损坏） | **Root** `D-012`（voucher 专属） |
| **其余全部时间列** | **正常转换**（负 epoch 是合法 instant，epoch 之前） | **Root** `D-015`：负 epoch **不是** sentinel |

- 本文件所有 `CASE` 分支**不得**出现 `< 0` 条件（转换表达式本身对负值已是正确的 floor 语义，见 §2.3）。
- 可执行断言见 `apps/api/internal/w040contracttest/contract_boundaries_test.go` 的 `TestNegativeMustFailClosed`（voucher 子测试必须失败；`ordinary_negative_is_valid_instant` 子测试必须通过）。

## 3. `#34` / `#72` / `#73` 的唯一化结果（响应 A-030 F-I-002.2）

| # | 列 | 新 DDL 形状 | copy 表达式 |
|--:|----|-------------|-------------|
| 34 | `mail_config.updated_at` | `TEXT`，**可空、无默认**（去 `NOT NULL` / `DEFAULT 0`） | `CASE WHEN updated_at = 0 THEN NULL ELSE <毫秒表达式> END` |
| 72 | `vouchers.expires_at` | `TEXT`，可空（原即 NULL） | `CASE WHEN expires_at IS NULL OR expires_at = 0 THEN NULL ELSE <秒表达式> END` |
| 73 | `vouchers.redeemed_at` | `TEXT`，可空（原即 NULL） | `CASE WHEN redeemed_at IS NULL OR redeemed_at = 0 THEN NULL ELSE <秒表达式> END` |

三者与 §2 模板、Root `D-015` 毫秒族（毫秒用整数 interval 语义）**同一**，不再存在第二种机制。

## 4. 行拷贝机制的用户裁决（2026-09-20 · P-004 选项 C）

- **决定**：行拷贝采用**纯整数 SQL 表达式**（§2）写进单条 `INSERT … SELECT`；**同时**把「Go codec `FromUnix` / `FromUnixMilli` 逐行等价比对」写成 R2 **强制验收步骤**（`T-<#>-RT` 用例族；任一不等价即 fail closed）。
- **理由**：转换逻辑完全落在 checksum 所覆盖的 `stmts` 切片内（契合 `D-017` 的单 checksum 约定），单条 SQL 可逐表核对，且不依赖尚未落码的 `internal/temporal`；对拍步骤同时满足「经 codec 语义」的方向要求。
- **已知偏差（须在交付时明示，不得掩盖）**：A-029 原文要求「copy rows through codec」，本决定下 codec 不在拷贝路径上，而是作为**独立对拍验收**。该偏差已由用户书面裁决接受；independent 复审如有异议需按 P-004 回到用户。
- **兼容窗口**：用户裁决**不留过渡期**——列形状改 `TEXT` 与仓储层扫描/绑定改造在**同一发布**内完成；迁移期间无「TEXT 列 + 旧 int64 扫描器」并存状态。

## 5. 声明

- 本文件 `status: freeze-candidate`；不改变 `status` / `progress` / goal-tree；**不**声称任何 DDL、迁移或测试已实施。
- §2 的两条表达式为**本机 SQLite 3.51.2 实测**结论；R2 落码时须以 `T-<#>-RT` 用例在目标 SQLite 版本上复证，并以 Go codec 对拍（§4）。
- 逐表正文见配套逐表附件；本文件与逐表附件是否构成 F-I-002.1 的合法闭合，由 independent 复审判定。
