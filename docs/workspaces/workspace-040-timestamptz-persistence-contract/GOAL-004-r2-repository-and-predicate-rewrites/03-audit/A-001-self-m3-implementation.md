---
id: A-001-self-m3-implementation
doc: audit-entry
status: active
parent: GOAL-004-r2-repository-and-predicate-rewrites
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
source: self
verdict: conditional
---

# A-001（self）· GOAL-004 检查点 A/B/C

- **source**: self（编排器自审）
- **日期**: 2026-09-20
- **scope**: GOAL-004 A（读写/谓词改造）、B（双方言回归 + 金额列断言拆分 + leftover 21 列）、C（边界测试重定向 + 全仓绿）；实现 commit `c69ee93d`。
- **verdict**: `conditional`

## 成果（可核对）

| 判据 | 证据 | 结果 |
|------|------|------|
| A：全部时间列读写/谓词改造 | `modules/**`、`internal/{jobs,mail,channel,auth,handler}` 的域类型化；`internal/store` 读写适配器（`bindSQLiteArgs`/`bindPostgresArgs` + `scan.go`）；`D-001` §2 逐列分档 | ✅ |
| #5/#6/#20/#34/#78 sentinel | `users.locked_until`、`users.last_login_failure_at`、`login_failures.locked_until`、`mail_config`/`telegram_config.updated_at` 谓词与 NULL 写入；`modules/authsession/temporal_sentinel_test.go`（P-5/P-6/P-20 回归） | ✅ |
| #61 `COALESCE(finished_at,0)` 去除 | `scheduledtasks/store` 去 COALESCE + `sql.NullTime` + 写 NULL | ✅ |
| #72/#73 voucher | `Valid` 判存在、负值 fail closed（读+写）、`TestNegativeExpiresAtFailsClosedAndIsNotAbsence` | ✅ |
| #4/#11 微秒单调 | `users_repository.go` / `roles_repository.go`（`Truncate(µs)` + `old + 1µs`）| ✅ |
| B：双方言回归 | SQLite：`go test -count=1 ./...` **63/63 ok**；PG：`TestFullCatalogPostgresBootstrapIntegration`、`TestPurchasePostgresAcceptance` 在真实 PG 15.4 通过 | ✅ |
| B：金额列断言拆分 + 21 名 | `postgres_test.go`（时间列 `timestamp with time zone`/precision 6；`balance_total`/`amount_delta` 保持 `bigint`；leftover 21 名 + 精度断言） | ✅ |
| C：边界测试重定向 | `internal/w040contracttest/migration_boundaries_test.go`（真实 v73–v87 迁移矩阵；含 voucher 负值 fail closed + v85 回滚证据） | ✅ |
| C：全仓绿后与 M2 一并提交 | commit `c69ee93d`（Root `D-017` §1） | ✅ |

## 偏差与开放项（编排器自评）

1. **读侧合同更正**：原 `D-001` §1 表格写「canonical TEXT 经 `database/sql` 直接解析」——经 Go 1.26 源码与实测证伪（`unsupported Scan, storing driver.Value type string into type *time.Time`）。已改为 `internal/store/scan.go` 承担，并同步修订 `D-001`。属**更正**，须记入台账（已列 `03-audit.md` 待复审事项 6）。
2. **ledger PG 写入两条 CAST 语句**：pgx 按 SQL 文本缓存 prepared statement，跨 v73 复用会把参数类型钉死为 `bigint`（实测 `22P02`）。修复为 legacy/canonical 两条不同文本（待复审事项 4）。
3. **`settings/migration` v7 seeder 保持 Unix 秒**：与冻结附件 §2.12「与 v84 同批改目标格式」字面冲突；实测若改绑 canonical 串，v84 的 `strftime` 会得 NULL → `NOT NULL constraint failed`（与 ledger 同型错误）。故 v7 保持绑定整数（待复审事项 5）。
4. **部分升级**：descriptor 各自事务，v85 失败时 v73–v84 已提交（非整批回滚）。C3 恢复边界覆盖该状态；是否需在 M4 前补批级快照待复审（待复审事项 7）。
5. **`PublicView.UpdatedAt` → `*time.Time`**（NULL 投影为 JSON null，替代原先的零值/epoch 伪造）。属 R2 范围内为「NULL = 未初始化」提供诚实读面；web 客户端已容忍非字符串（`typeof view.updated_at === "string"`）。审计可判为超出「不改 handler 响应编码」边界（该字段是模型投影，未改 formatter）。
6. **`Store.WithTx`（sqlite-only 原始 seam）绕过两个适配器**：测试中若经它写时间列会落成非规范值；本轮已把 `auth_test.go`、`handler/wallet_test.go` 的相关写入改走 `st.Run`。该 seam 本身为遗留债务（R4 记录在案），未在 R2 扩大改造。
7. **子代理二手证据**：M3 的模块改造由 4 个并行子代理完成，编排器以聚合测试 + 抽样核对（`ORDER BY`/谓词保留、voucher 负值、扫描适配器、PG 参数类型）验收，未逐行复核全部 diff。

## 自审边界（诚实声明）

- 本自审不闭合任何 required finding；上述 7 项与 `03-audit.md` 的 8 项待复审事项**全部**提交 independent。
- 「全仓绿」是**必要非充分**证据：不覆盖本地临时 PG 之外的部署形态、不覆盖 R3 的展示/输入时区矩阵、不覆盖 M4 的备份/还原。
