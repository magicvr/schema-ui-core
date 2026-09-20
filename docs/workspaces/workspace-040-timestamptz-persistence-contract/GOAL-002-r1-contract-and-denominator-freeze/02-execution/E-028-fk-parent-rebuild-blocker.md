---
id: E-028-fk-parent-rebuild-blocker
doc_type: goal-execution-entry
status: recorded
date: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# E-028 · 逐表 rebuild DDL 前置调查发现 FK 父表阻塞

## 事实

1. **用户 2026-09-20 P-004 裁决两项**（落盘 child `01-decision/D-018-r2-row-copy-and-compat-window.md`）：
   - 行拷贝机制 = **选项 C**：纯整数 SQL 表达式写进单条 `INSERT … SELECT`，**并**把 Go codec 逐行等价比对列为 R2 强制验收步骤；
   - 扫描器兼容窗口 = **选项 A**：无过渡期，列形状变更与仓储层扫描/绑定改造同一发布完成。
2. **冻结表达式实测修正**：`strftime('%f', …)` 实测 **round**（`1758320000.9999` → 下一秒）且输出变宽（整秒仅 3 位小数），全包禁用；`/1000.0` 同样禁用。秒族与毫秒族改用纯整数式；毫秒余数必须归一化为 `(<col>%1000 + 1000) % 1000`（SQLite `%` 符号跟随被除数，`-1%1000 = -1`，旧 CASE 形式会产出 `-01000`）。落盘 `attachments/r1-c2-sqlite-rebuild-mechanism-v1.0-fc.md`。
3. **整表重建探针通过**：4 行样本（含 NULL、D0、负秒、负毫秒、9999 年）重建后 `bad_len_rows=0`（全部 27 字符）、`PRAGMA foreign_key_check` 空、D0 `0 → NULL` 正确。
4. **发现 PK 前置阻塞（本轮主要产出）**：重命名 FK 父表会**永久**改写子表 `REFERENCES` 文本；`legacy_alter_table=ON` 无效；`foreign_keys=OFF` 因 DSN 固定 + 事务内 no-op 而不可用；朴素重建在子表有行时会**级联删除子表数据**或在 `DROP <parent>_old` 时报 FK 失败。已验证可行的候选模式为「子表 TEMP 快照 → DROP 子表 → 重建父表 → 重建子表 → 回填 → 建索引」，与仓库既有 `rebuildOperationLogWithSessions`（`operationlog/migration/migration.go:709-754`）同构。全部落盘 `attachments/r1-c2-fk-parent-rebuild-finding-v1.0-fc.md`。
5. **附带发现两项**：
   - `users` 有 3 个**跨模块** ALTER 列（`enabled` v0013、`notifications_enabled` v0017、`avatar_url` v0035），v74 的 `users` 新 DDL 若遗漏即静默丢列；且真实列序由全局版本序决定，不是按文件分组。
   - `operationlog/migration/migration.go:250` 的 PG DDL 由正则 `INTEGER NOT NULL` 派生；SQLite 字面改 `TEXT` 后该派生**静默失效**，PG DDL 必须改为显式书写。
6. **`schema_migrations` 自引用核对无危害**：ledger 只在 `applyPending` 前读取，`applyMigration` 先 `Apply` 后同事务插 ledger 行，v73 行在重建 copy 时尚不存在；残留项为 `actionRestoreLedger` 路径只依赖 `identity.go:58/:65` 两个字面，须同批改为 `TEXT applied_at`。
7. **逐表 DDL 已抽取但未落盘**：15 个 owner / 20 张受影响表的 exact legacy DDL、目标 DDL、索引与 copy 表达式已逐表抽取并在内存 SQLite 3.51.2 上试跑（含 `jobs` 六态 CHECK、`dict_entries.badge_style`、`telegram_config` 两个 v0067 ALTER 列、`wallet_ledger_entries` 14 列、`site_settings` 15 列等有效形状）。**因用户裁决先审计 FK 方案再定结构，该正文尚未落盘为附件。**

## 证据

- 实测环境：本机 `sqlite3` 3.51.2（`C:\Program Files\msys64\mingw64\bin\sqlite3.exe`），全部在 `:memory:` 或临时库内执行，未触碰仓库数据。
- 代码对位：`internal/store/store.go:57,107`；`internal/store/migrate.go:32,81-132,221-251,253-265,345-363`；`internal/store/identity.go:58,65,329-372`；`modules/operationlog/migration/migration.go:250,544-562,703-776,782-786`；`modules/wallet/migration/migration.go:180-201`；`modules/datadictionary/migration/migration.go:32,76`；`modules/account/migration/migration.go:21,28`；`modules/notifications/migration/migration.go:46`。
- 本条目**未修改** `apps/**` 任何文件。

## 状态评估

- 本条目记录的 FK 阻塞**尚未定案**：用户已裁决「先让 grok 独立审计该发现再定处置方案」（`/audit` → A-032），并裁决「既有 `rebuildOperationLog` 的 FK 隐患本轮一并修复」（修复待审计结论后实施）。
- F-I-002 的「逐表 exact SQLite rebuild DDL」子项**仍未闭合**：正文已具备但受 FK 处置方案阻塞。
- 开放 required 仍为 4（F-I-002、F-I-004、F-I-005、F-I-006）；C2/C3 未冻结，R2 未放行。
