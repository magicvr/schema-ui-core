---
id: GOAL-004-r2-repository-and-predicate-rewrites
doc: audit
status: active
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.3.0
---

# 审计台账 · GOAL-004（R2 M3）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见的正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source` / 日期 / scope / `verdict`）。
> 独立审计默认只写意见，不修改 `status` / `progress` / 方案正文；响应归编排器。

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| A-001 | self | 2026-09-20 | GOAL-004 A/B/C（读写/谓词改造、双方言回归、边界测试重定向） | conditional | 成果 9 项可核对；7 项偏差提交 independent（读侧合同更正、PG ledger CAST、v7 seeder、部分升级、PublicView 类型、WithTx seam、子代理二手证据） | `03-audit/A-001-self-m3-implementation.md` |
| A-002 | independent | 2026-09-20 | GOAL-004 A/B/C + 待复审事项 1–8 + 未声明回归（commit `c69ee93d`） | conditional | 谓词/sentinel/金额列/leftover 21/读侧 fail-closed 可核对；新 required = PG 写路径未 Truncate（Root D-008）；待复审 1/4/5 判为不需修正 | `03-audit/A-002-independent-m3-implementation.md` |
| A-003 | self（编排器响应） | 2026-09-20 | 响应 A-002 全部意见 | **pass** | `F-I-001` **fixed**（PG 写路径 Truncate + 双方言对拍）；`F-I-004/005/006/007/008` fixed；`F-I-003` 部分 fixed + 记录；`F-I-002` **`user-overruled`**（用户 2026-09-20 书面接受 `*time.Time` + JSON null 属 R2 必要后果） | `03-audit/A-003-response-to-independent-m3.md` |

## 待复审事项（编排器登记，供独立审计取证）

| # | 事项 | 证据位置 | 说明 |
|--:|------|----------|------|
| 1 | ledger 写入的**列形状探测**偏离 M2 冻结附件 §6 字面（改为同事务探测） | `internal/store/identity.go`、`internal/store/migrate.go`、`internal/store/postgres.go` | 实测原字面写法会致 `NOT NULL constraint failed: schema_migrations.applied_at`；是否接受由 independent 判定 |
| 2 | v86 descriptor 的 `ModuleID` 由台账的 `admin.channel.telegram` 更正为模块实际的 `channel.telegram` | `modules/channel/telegram/migration/vp040_temporal.go`、`attachments/r2-v73-v87-generated-statements-v0.1.md` | 模块 ID 由 provider 机械决定；台账 §1 的 `ModuleID` 列不在 `D-014`/台账 §4 点名的「唯一允许输入」清单内 |
| 3 | v73–v87 语句由生成器从 **v72 live `sqlite_master`** 机械导出（含 m0/m4 查询语句进入 checksum 输入） | `attachments/r2-v73-v87-generated-statements-v0.1.md` | 逐列类型编辑的机械性与 m0/m4 进入哈希是否符合 `D-017` 由 independent 判定 |
| 4 | **PG ledger 写入使用两条显式 CAST 语句**（legacy `CAST(? AS bigint)` / canonical `CAST(? AS timestamptz)`），偏离「同一语句文本」的直觉写法 | `internal/store/identity.go` `postgresLedgerWrite` | 实测 pgx 按 SQL 文本缓存 prepared statement，跨 v73 复用同一文本会把参数类型钉死为 `bigint`（`22P02`）；是否接受由 independent 判定 |
| 5 | `modules/settings/migration/migration.go` 的 v7 seeder **保持 Unix 秒**，偏离冻结附件 §2.12「与 v84 同批改目标格式」 | 该文件 `migrate0007` / `migrate0007PG` | v7 在 v84 之前执行、列仍是 INTEGER；改绑 canonical 串会使 v84 的 `strftime` 得 NULL → `NOT NULL constraint failed`。是否接受该偏离由 independent 判定 |
| 6 | 读侧合同（`database/sql` 不解析 string → `*time.Time`）由 `internal/store/scan.go` 承担，D-001 §1 已同步更正 | `internal/store/scan.go`、`internal/store/scan_test.go` | 原 D-001 §1 表格写法经 Go 1.26 实测证伪；属**更正**而非新增设计，须记录 |
| 7 | 失败即部分升级：descriptor 各自事务，故 v85 失败时 v73–v84 已提交（非整批回滚） | `internal/store/migrate.go` `applyPending`；`internal/w040contracttest/migration_boundaries_test.go` | C3 恢复边界已覆盖；是否需要在 M4 前补「批级回滚/快照」由 `D-021` 复审一并判定 |
| 8 | 4 个子代理并行改造的一致性 | 各模块 `02-execution` 事实 + `git diff` | 判据：无方言分支、无驱动类型泄漏、`ORDER BY`/谓词逐字保留、无断言弱化 |

## A-002 · independent · 检查点 A/B/C（2026-09-20）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **verdict**：conditional
- **完整意见**：[`03-audit/A-002-independent-m3-implementation.md`](A-002-independent-m3-implementation.md)

独立判定（**不是**闭合；闭合归 `/govern` / P-004）：

| # | 独立判定 |
|--:|----------|
| 1 | 不需修正（同事务探测是无过渡期单批的必要形态） |
| 2 | 不是 checksum 违约；台账展示列建议回写 |
| 3 | m0/m4 进哈希 **符合** D-017；不因该项重算 15 个 checksum |
| 4 | 不需修正（pgx prepared-statement 缓存约束） |
| 5 | 不需修正（改 v7 会破坏 v84 且动 v1–v72 Apply 语义） |
| 6 | 记为设计更正 `fixed` |
| 7 | 确认既有 runner 语义；不作为 M3 required；是否补批级快照交 M4/用户 |
| 8 | 抽样通过，非穷尽 |

新开放 required：`A-002` **F-I-001**（PG 写路径未对 `time.Time`/`sql.NullTime` 做微秒向零截断，违反 Root `D-008`）。**R2 未放行**。

## A-003 · self（编排器响应）· 2026-09-20

- **verdict**：conditional（唯一未闭合项 = `F-I-002`，需用户 P-004）
- **完整响应**：[`03-audit/A-003-response-to-independent-m3.md`](A-003-response-to-independent-m3.md)

| finding | 处置 |
|---------|------|
| `F-I-001`（required，PG Truncate） | **fixed**：`bindPostgresArg` 统一 `temporal.Truncate`；`temporal_truncate_test.go` 双方言对拍（真实 PG 15.4 逐字核对 `...123456Z`） |
| `F-I-002`（`PublicView.UpdatedAt` → JSON null） | **待用户 P-004 裁决**（接受 / 退回 R3 / 恢复伪造瞬时） |
| `F-I-003`（`WithTx` 绕过适配器） | 部分 fixed（本批触及的测试改走 `kernel.Tx`）+ 记录遗留 seam |
| `F-I-004`（`user_mfa` 整数 seed） | **fixed**（并更正「该表无转换」的误述：v80 确实转换 `user_mfa`） |
| `F-I-005`（recyclebin payload） | **fixed**：两形状各一例（新增 canonical 字符串快照 restore 对拍） |
| `F-I-006`（leftover 查询漏 bigint） | **fixed** |
| `F-I-007`（`#20` 真实迁移 0→NULL 未覆盖） | **fixed**（边界矩阵新增 `login_failures.locked_until = 0`） |
| `F-I-008`（recyclebin handler 测试绑 `Unix()`） | **fixed** |
| 待复审 1/2/3/4/5/6/8 | 按 A-002 判定记录（1/2/4/5 不需修正；3 不重算；6 fixed；8 抽样） |
| 待复审 7（部分升级） | 保留为 **M4 决策输入**（是否补批级快照/整批回滚） |

## 关门记录（2026-09-20）

A/B/C 检查点全部 `completed`；`03-audit` 的 self（A-001）与 independent（A-002）意见均已落盘并由 A-003 响应（**pass**）：唯一 required（PG 写截断，Root D-008）已 `fixed` 并有双方言对拍证据；`F-I-002` 经用户书面裁决 `user-overruled`；其余 recommended 项均已 fixed 或按独立判定记录。全仓 `go test -count=1 ./...` 63/63 包 ok（含真实 PostgreSQL 15.4 路径）。据此按既有用户裁决（非关键子目标可经交叉审计后静默关门）**静默关门**（`done · 3/3`）。

**边界**：本关门只覆盖 R2 的 M3（仓储读写/谓词/回归/边界重定向）；**不**等于 R2 完成 —— M4（Backup Port 类型表面与 provider、`D-021` 第 9 项、R2 关门审计）由 `GOAL-005-r2-backup-port-and-closeout` 承担。
