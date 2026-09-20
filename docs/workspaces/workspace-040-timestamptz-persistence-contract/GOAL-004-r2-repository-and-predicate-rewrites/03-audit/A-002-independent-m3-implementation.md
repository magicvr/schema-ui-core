---
id: A-002-independent-m3-implementation
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-004-r2-repository-and-predicate-rewrites · 检查点 A/B/C + 待复审事项 1–8 + 未声明回归 · commit c69ee93d
verdict: conditional
open_required: 1
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-004-r2-repository-and-predicate-rewrites
version: 0.1.0
---

# A-002 · independent · GOAL-004 检查点 A/B/C

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：execution-facts · 检查点 A（读写/谓词）/ B（双方言回归 + 金额列拆分 + leftover 21）/ C（边界测试重定向）+ `03-audit.md` 待复审事项 1–8
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`root_goal` / `canonical_scope` / `shared_materials_catalog: none` / `primary_plan` 均匹配）。未读其他工作区。
- 被审目标：`GOAL-004-r2-repository-and-predicate-rewrites`。descriptor/checksum 主体在 sibling GOAL-003 A-002；本条只在交叉需要时引用，不改 GOAL-003 状态。
- 对照：本目标 A-001（self · `conditional`）；Root `D-008`/`D-012`/`D-013`/`D-015`/`D-016`/`D-017`；GOAL-002 冻结合同；本目标 `D-001`。
- 实现 commit `c69ee93d`。本轮未重跑全仓 `go test ./...` 或真实 PG（self 报 63/63 与 `TestFullCatalogPostgresBootstrapIntegration` ok；标 **unverified here**）。读适配器/绑定/谓词/边界测试/ledger 写入已直接读源码。
- **未改** status / progress / 方案正文 / Go 代码。

## 成果（有证据）

1. **读适配器 fail-closed。** `scan.go:111-127`：string/`[]byte` → `Parse`（宽）→ `Format`（严）→ **必须与存量文本字节相同**。3 位小数、`+08:00` 被拒（`scan_test.go:70-79`）。`int64` 存量 → `unsupported stored time value type`（不静默当 epoch）。PG 原生 `time.Time` 走 `v.UTC()`，不要求 canonical 文本（按设计）。
2. **SQLite 写适配器走 `NewValue`（Truncate + 27 字符）。** `store.go:230-254` 覆盖 `time.Time` / `*time.Time` / `sql.NullTime` / codec `Value`/`NullValue`。这是 Root `D-008` 在 SQLite 侧的正确落点。
3. **sentinel / voucher 运行时分档按 D-001 §2 落实。** #5 已锁定 `IS NOT NULL AND > ?`、未锁定 **加括号** `(IS NULL OR <= ?)`（`users_repository.go:505-509`）——冻结谓词原文无括号，AND 拼接时有优先级缺陷，**代码是对的，冻结文本应改**。#6 写 NULL、读 `IS NULL OR < ?`。#20 插入 NULL、`sql.NullTime`。#61 去掉 `COALESCE(finished_at,0)`。#72/#73 `Valid` 判存在、负值写路径 fail closed。#4/#11 微秒单调存在。
4. **D-021 residual ②③ 源码满足。** 金额列保持 `bigint`（`postgres_test.go:307-314`）；leftover 21 名 = allocation L38（`postgres_test.go:319-325`）；另有「凡此 21 名必须是 `timestamptz` 精度 6」查询（L336-344），可抓住仍为 `bigint` 的时间列。
5. **边界测试已改指真实迁移（D-020）。** `migration_boundaries_test.go`：v72 库 + legacy 整数播种 + 应用到 head；sentinel 0→NULL、999 ms 无进位、负毫秒 floor、公元 9999、词法序、voucher 负值 fail closed（v85 回滚、列仍 INTEGER）、普通负瞬间转换。
6. **`jobs.NewID` / `voucher.newID` 的 `UnixMilli` 是 ID 前缀，不是时间列。** 不在 90 列分母。**不是回归。**

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| A：时间列读写/谓词改造 | **部分** | 谓词/sentinel 落实；PG 写路径未 Truncate（F-I-001） |
| B：双方言回归 + 金额列 + leftover 21 | **源码满足；PG 本轮未复跑** | `postgres_test.go`；I-041-003 仍以 Root D-017 用户裁决为准 |
| C：边界重定向 + 全仓绿 | **边界测试源码满足；全仓绿 unverified here** | `w040contracttest`；self 报 63/63 |
| I-041-005 | **不能无条件关闭** | 适配器形态正确，但 PG 直通违反 Root D-008（F-I-001） |

## 待复审事项 1–8（问题 8：逐项判定）

独立审计**不得**代用户走 `accepted-residual`。下列「建议路径」是给 `/govern` / P-004 的建议，不是闭合。

| # | 事项 | 本审判定 | 建议路径 |
|--:|------|----------|----------|
| 1 | ledger 列形状同事务探测 | **不是缺陷。** v73 与 v1–v72 共用 `applyMigration`；v73 之前列是 INTEGER，之后是 TEXT/`timestamptz`。按冻结「一律绑 canonical」会在 v1–v72 段把文本写入 INTEGER，v73 的 `strftime(...,'unixepoch')` 得 NULL → `NOT NULL constraint failed`（self 实测）。探测是 D-019 §4「目标类型」在无过渡期单批内的唯一可执行形态。Restore DDL 已是 TEXT / `timestamptz(6)`（`identity.go:62-76`）。 | 不需修正。若要对齐冻结附件字面，由用户书面 `accepted-residual`（范围=探测，不改目标类型） |
| 2 | v86 `ModuleID` → `channel.telegram` | **不是 checksum/transform_id 违约。** 台账 ModuleID 列不在 §4 唯一允许输入里；kernel 模块 ID 是 `channel.telegram`。 | 回写台账展示列（recommended，见 GOAL-003 F-I-005） |
| 3 | 生成器 + m0/m4 进哈希 | **m0/m4 符合 D-017/台账 §2**（GOAL-003 A-002）。生成器机械性：CREATE 为 live DDL + 逐列类型编辑，本审抽查通过。 | 不因「m0/m4 进哈希」重算 15 个 checksum |
| 4 | PG ledger 两条 CAST | **不是缺陷。** pgx 按 SQL 文本缓存 prepared statement，跨 v73 复用会把参数钉成 `bigint`（self 实测 `22P02`）。两条不同文本是驱动约束，不进 checksum。 | 不需修正；可选用户书面 residual |
| 5 | settings v7 seeder 保持 Unix 秒 | **不是缺陷；改了会破坏 v84。** v7 在 v84 之前执行，列仍是 INTEGER/`BIGINT`（`settings/migration/migration.go:52-55,69-72`）。绑 canonical 串 → v84 `strftime(updated_at,'unixepoch')` 得 NULL → `NOT NULL constraint failed`。且改 v7 绑定会动 v1–v72 Apply 语义（Root 红线）。冻结 §2.12「与 v84 同批改目标格式」应读成 **v84 之后的写入**，不是改写 v7。 | 不需修正 |
| 6 | `scan.go` 读侧更正 | **确认是更正。** `database/sql` 不把 string 扫进 `*time.Time`；D-001 §1 已改。严格往返拒绝非规范存量。 | 记为 `fixed`（设计更正，不是偏离） |
| 7 | 部分升级（v85 失败留下 v73–v84） | **确认是既有 runner 语义，不是本批新引入。** `applyPending` 一 descriptor 一事务（`migrate.go:93-101`）。边界测试已文档化（`migration_boundaries_test.go:287-289`）。C3 恢复边界覆盖该状态。混型库上对仍为 INTEGER 的列扫 `time.Time` 会 fail closed（`scan.go` 拒 int64），不是静默 epoch。 | **不**作为 M3 required。是否在 M4 前补批级快照交用户/M4 决策；独立审计不代接受 residual |
| 8 | 四子代理一致性 | **抽样通过，非穷尽。** 抽查 authsession/wallet/mail/telegram/recyclebin/scheduledtasks：无模块内方言分支、无 `pgtype`、wallet `ORDER BY created_at, id` 保留、时间列 `COALESCE(...,0)` 已不在。未逐行复核全部 diff。 | 保持抽样结论；不升 required |

## Findings

### F-I-001 · PG 写路径不对 `time.Time`/`sql.NullTime` 做微秒向零截断（违反 Root D-008）

- **严重度**：high
- **建议**：required
- **状态**：open
- **影响门禁**：检查点 A；I-041-005；Root `D-008`（写入统一截断到微秒，**不得 round**）
- **描述**：SQLite `bindSQLiteArg` 经 `temporal.NewValue` Truncate。PG `bindPostgresArg`（`postgres.go:311-328`）对 `sql.NullTime` 返回 **未截断** 的 `v.Time`，对 `time.Time`/`*time.Time` **default 直通**。这与本目标 D-001「PG `time.Time` 直通」字面一致，但 **Root D-008 优先级更高**：`timestamptz(6)` 的 typmod 是 **round**，不能当向零截断。codec 注释自己写明 `Format` 不截断、裸 `time.Time` 会被 `time.Format` round（`temporal.go:96-99`）。生产仓储普遍绑 `time.Now().UTC()` 而不先 `Truncate`（mail `runtime.go` 更新路径、telegram、recyclebin `DeletedAt` 等）。
- **可执行反例**：

```text
t = 2025-09-19T22:13:20.123456789Z   # 123456789 ns
SQLite bindSQLiteArg → NewValue.Truncate → "2025-09-19T22:13:20.123456Z"
PG    bindPostgresArg → 直通 time.Time → timestamptz(6) ROUND → 123457 µs
     → 读回 "2025-09-19T22:13:20.123457Z"
```

复现：在 PG 测试库插入上述 `time.Time`，`SELECT` 该列与 SQLite 同一输入对比；或最小片段：

```go
t := time.Date(2025, 9, 19, 22, 13, 20, 123456789, time.UTC)
fmt.Println(temporal.NewValue(t).String()) // 123456Z
// 把 t 经 pgTx.Exec 写入任一 timestamptz(6) 列后 Scan 回来
```

- **关闭要求**：在 `bindPostgresArg` 对 `time.Time`/`*time.Time`/`sql.NullTime` 调用 `temporal.Truncate`（与 SQLite 对称），并加双方言对拍用例；**或**用户书面 `user-overruled`/`accepted-residual`（须写明「接受 PG typmod round vs SQLite truncate」的范围与复审触发）。在闭合前 **I-041-005 不得视为已验证**。

### F-I-002 · `mail.PublicView.UpdatedAt` 改为 `*time.Time` 改变 JSON 形状

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：`internal/mail/runtime.go:75-81,332-334`：NULL → JSON `null`，不再伪造瞬时。诚实对齐 #34 D0，但是对 `json:"updated_at"` 读面的形状变化。Root `D-016` 把「handler 响应编码 / rfc3339.go」划出 R2。本字段是模型投影不是 formatter，web 已容忍非 string。
- **关闭要求**：用户书面接受该投影变化属 R2 必要后果；或推迟到 R3 与 wire formatter 同批。独立审计**不**代接受。

### F-I-003 · `Store.WithTx` 仍绕过读写适配器

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：`store.go:179-191` 把裸 `*sql.Tx` 交给回调。绑 `time.Time` 会落入驱动 36 字符布局；读也不走 `scan.go`。生产模块走 `Run`/`kernel.Tx`。测试仍大量使用 `WithTx`（authsession/wallet/settings/testsupport 等）。self 已把部分测试改走 `st.Run`。遗留 seam，R2 未扩大。
- **反例**：经 `WithTx` 执行 `UPDATE site_settings SET updated_at = ?` 绑 `time.Now().UTC()`，再用 `st.Run`+`Scan(&time.Time)` —— 期望 `not canonical`。

### F-I-004 · `users_repository_test.go` 在 v80 转换后的 `user_mfa` 上仍插入整数 `1,1`

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：用户提示「该表无 descriptor 转换」**不成立**——v80 转换 `user_mfa.created_at`/`updated_at`。测试在全 catalog 之后 `INSERT … created_at, updated_at VALUES (…, 1, 1)`（`users_repository_test.go:123-124,158-159`），经 `WithTx`。SQLite TEXT affinity 把 `1` 存成 `"1"`。测试只 `COUNT(*)`，故仍绿。随后若有人按时间扫描这些夹具行，`scan.go` 会 fail closed。
- **关闭要求**：改绑 canonical 瞬时，或注明「仅测 FK 级联、不读时间列」。

### F-I-005 · recyclebin payload 的 `timeField` 仍把 JSON 数字当 Unix 秒

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：列 `deleted_at`/`restored_at` 已域类型化。payload 快照仍是 `map[string]any`。`timeField`（`recyclebin/service.go:270-284`）对 `float64`/`int64` 走 `time.Unix`。这对**迁移前**写入的旧 payload 是兼容，不是列转换遗漏。生产快照在 R2 之后是否改为 RFC3339 **unverified**（测试仍种 unix float64）。
- **关闭要求**：冻结 payload 时间字段的新旧两种形状，并各给一例 restore；或声明 payload 不在 90 列分母、保持 unix 数字。

### F-I-006 · leftover 整数查询只匹配 `integer`，漏 `bigint`（被第二条查询兜住）

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`postgres_test.go:327-328` `data_type = 'integer'` 抓不住历史 PG 时间列（`bigint`）。随后的 unprecise 查询（L336-338）要求 21 名全部是 `timestamp with time zone` 精度 6，**会**抓住残留 bigint。residual ③ 仍满足。建议把第一条改为 integer **或** bigint，避免误读。

### F-I-007 · 真实迁移 0→NULL 矩阵未覆盖 `#20 login_failures.locked_until`

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`TestRealMigrationsConvertBoundaryInstants` 覆盖 #5/#6/#34/#78 与 voucher，不播种 `login_failures`。转换 SQL 与 #5 同构（`vp040_temporal.go:171-172,350-352`）。运行时 NULL 插入有 `temporal_sentinel_test.go`。缺的是 **legacy `locked_until=0` 行经真实 v74 变 NULL** 的可执行证据。
- **反例**：在 `atV72` 库 `INSERT INTO login_failures (…, locked_until, …) VALUES (…, 0, …)`，`toHead` 后 `SELECT locked_until` 应为 SQL NULL。

### F-I-008 · `handler/recyclebin_test.go` 仍把 `now.Unix()` 写入 `deleted_at`

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：测试双经 `kernel.Tx` 绑 `int64`。适配器不转换 int64。SQLite 存 `"1758…"` 文本，随后时间扫描 fail closed。不是生产服务路径。

## 未声明回归（问题 9）

| 项 | 判定 |
|----|------|
| `mail.PublicView.UpdatedAt` 语义 | **已声明偏离**（self 待复审 5 / A-001 项 5）。诚实 NULL；JSON 形状变了。见 F-I-002 |
| `jobs.NewID` / `voucher.newID` UnixMilli | **不是回归**（ID 形状，非时间列） |
| recyclebin payload 整数字段 | **payload 兼容层仍在**（F-I-005）；列本身已改 |
| `user_mfa` 整数 seed | **测试夹具味道**（F-I-004）；表**有** v80 转换 |

## 必改项汇总

1. **F-I-001（required）**：PG 写适配器 Truncate，或用户书面接受「PG round vs SQLite truncate」残余。

待复审 1/2/3/4/5/6 不构成本条 required。7 交 M4。F-I-002 若用户把 D-016「不改响应编码」读成字面禁止 JSON null，可在 `/govern` 升为 required——本审保持 recommended 并 **等用户裁决**。

## 与 A-001（self）的异同

| 项 | self | 本条 |
|----|------|------|
| 读侧 scan.go | 更正 | 同意，记 fixed |
| PG ledger CAST / 列探测 / v7 seeder | 待判 | **不需修正**（必要偏离） |
| 部分升级 | 待判 | 确认；M4 议题，不阻 M3 required |
| PublicView `*time.Time` | 待判 | recommended F-I-002，需用户 |
| WithTx | 待判 | recommended F-I-003 |
| PG Truncate | **未提** | **新 required F-I-001** |
| 四子代理 | 二手证据 | 抽样同意 |

## 结论 + 建议给编排器/用户的下一步

M3 的谓词、sentinel、金额列拆分、leftover 21、边界重定向与读侧 fail-closed **可核对**。开放 required = **F-I-001**（PG 写截断）。在该条 `fixed` / 用户书面 residual / overruled 之前，不得关闭 I-041-005，不得把 GOAL-004 标 `done`。

建议 `/govern` 下一句：响应 GOAL-004 A-002 F-I-001（建议在 `bindPostgresArg` Truncate）；对 F-I-002 与待复审 7 做 P-004；GOAL-003 侧见该目标 A-002（F-I-005 residual 可 `fixed`，另有 v73 `records` 断言 required）。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
