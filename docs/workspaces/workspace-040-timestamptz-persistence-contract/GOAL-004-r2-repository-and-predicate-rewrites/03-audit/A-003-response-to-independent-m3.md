---
id: A-003-response-to-independent-m3
doc: audit-entry
status: active
parent: GOAL-004-r2-repository-and-predicate-rewrites
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
source: self
verdict: pass
---

# A-003（self · 编排器响应）· GOAL-004 检查点 A/B/C

- **source**: self（编排器汇总响应）
- **日期**: 2026-09-20
- **scope**: 响应本目标 A-002（independent · grok-build grok-4.6 · high · `conditional`，open required = 1）的全部意见
- **verdict**: `pass`（A-002 的 required 已 `fixed`；`F-I-002` 经用户 2026-09-20 P-004 书面裁决 `user-overruled`）

## 必改项闭合

| finding | 级别 | 处置 | 证据 |
|---------|------|------|------|
| A-002 `F-I-001`（PG 写路径不截断 → `timestamptz(6)` round，违反 Root `D-008`） | required | **fixed** | `bindPostgresArg` 对 `time.Time` / `*time.Time` / `sql.NullTime` 统一 `temporal.Truncate`（与 SQLite 对称；codec `Value`/`NullValue` 本已截断）；**双方言回归**：`internal/store/temporal_truncate_test.go` 的 `TestPostgresWriteTruncatesToMicroseconds`（真实 PG 15.4：输入 `...123456789ns` → 存储 `...123456Z`，`to_char` 逐字核对）与 `TestSQLiteWriteTruncatesToMicroseconds` 同输入对拍 |

## recommended 项处置

| finding | 处置 | 证据 / 说明 |
|---------|------|-------------|
| `F-I-002`（`mail.PublicView.UpdatedAt` 改 `*time.Time` → JSON `null`） | **`user-overruled`（用户 2026-09-20 书面裁决：接受 `*time.Time` + JSON null 属 R2 必要后果）** | 事实：`updated_at` 为 NULL（未配置）时，此前读面会伪造 `time.UnixMilli(0)` 瞬时，现投影为 `null`；web 客户端已容忍非字符串（`mail-admin-tab.tsx:80`）。Root `D-016` 把「handler 响应编码 / `rfc3339.go`」划出 R2，本字段是模型投影而非 formatter。**等用户选择**：接受（R2 内必要后果）/ 退回 R3 与 wire formatter 同批 / 恢复伪造瞬时 |
| `F-I-003`（`Store.WithTx` 裸 seam 绕过适配器） | **部分 fixed + 记录** | 已把本批触及的测试写入改走 `kernel.Tx`（`internal/auth/auth_test.go`、`internal/handler/wallet_test.go`、`modules/authsession/repository_test.go` 的两个 helper、`modules/authsession/users_repository_test.go`）；该 seam 本身是 R4 记录在案的遗留债务，R2 不扩大改造。**未闭合的残余**：仓库内仍有多处测试经 `WithTx` 写非时间列（不受影响）；若后续有测试用它写时间列，会 fail closed 而非静默错误 |
| `F-I-004`（`user_mfa` 整数 seed） | **fixed** | 更正子代理的误述：`user_mfa.created_at/updated_at` **确由 v80 转换**；两处 seed 改为绑 `now`（`users_repository_test.go` 两个测试） |
| `F-I-005`（recyclebin payload 的 `timeField` 仍接受 Unix 数字） | **fixed（两种形状都冻结并各有一例）** | payload 是 JSON 快照、不在 90 列分母；`timeField` 同时接受 legacy Unix 数字与 canonical/RFC3339 字符串。新增 `TestRestoreTaskRoundTripCanonicalStringPayload`（canonical 字符串快照 → restore 后 `createdAt` 保留**快照瞬时**而非回退到 restore 时刻），与既有 `float64(now.Unix())` 用例共同冻结两形状 |
| `F-I-006`（leftover 查询只匹配 `integer`，漏 `bigint`） | **fixed** | `postgres_test.go` 改为 `data_type IN ('integer','bigint')`（第二查询仍要求 21 名全为 `timestamptz` 精度 6 作双保险） |
| `F-I-007`（真实迁移矩阵未覆盖 `#20`） | **fixed** | `TestRealMigrationsConvertBoundaryInstants` 新增 `login_failures` legacy `locked_until = 0` 播种，断言真实 v74 后为 SQL NULL |
| `F-I-008`（`handler/recyclebin_test.go` 经 `kernel.Tx` 绑 `now.Unix()`） | **fixed** | 改为绑 `now`（适配器归一为 canonical 串） |

## 待复审事项 1–8（A-002 判定为「非缺陷 / 必要偏离」者）

| # | 事项 | 处置 |
|--:|------|------|
| 1 | ledger 列形状同事务探测 | **无需修正**（A-002 判定为无过渡期单批内的唯一可执行形态）。若需对齐冻结附件字面，属用户书面 `accepted-residual` 的范围 |
| 2 | v86 `ModuleID` | **无需修正**（同一展示列更正，见 GOAL-003 A-003） |
| 3 | m0/m4 进哈希 | **无需重算**（A-002 判定符合台账 §2 + `D-017` 括号定义；把 m0/m4 排除会构成对 `D-017` 的修订，需 P-004） |
| 4 | PG ledger 两条 CAST | **无需修正**（pgx 语句缓存约束，实测 `22P02`） |
| 5 | settings v7 seeder 保持 Unix 秒 | **无需修正**（A-002 判定冻结 §2.12 应读作「v84 之后的写入」；改 v7 会破坏 v84 且改动 v1–v72 Apply 语义） |
| 6 | `scan.go` 读侧更正 | **fixed**（设计更正，已同步 `D-001`） |
| 7 | 部分升级（v85 失败留下 v73–v84） | **非 M3 required**；作为 **M4 决策输入**（是否补批级快照/整批回滚）——编排器保留该项，M4 立项时一并裁决 |
| 8 | 四子代理一致性 | 抽样结论保留；本响应不升 required |

## 未闭合项

- **A-002 `F-I-002`**：已由用户 2026-09-20 P-004 书面裁决 —— **接受** `*time.Time` + JSON `null` 属 R2 的必然后果（`mail_config.updated_at` NULL = 未配置，不再伪造 1970/零值瞬时；web 客户端仅在字段为 string 时使用，已容忍 null）。该 recommended finding 以 `user-overruled` 闭合，**不再要求改动**。\n- 至此 GOAL-004 无未闭合 required / recommended；`user-overruled` 依据见本文件与 `03-audit.md`。
- 其余 required / recommended 均已处置或按 A-002 判定记录。
