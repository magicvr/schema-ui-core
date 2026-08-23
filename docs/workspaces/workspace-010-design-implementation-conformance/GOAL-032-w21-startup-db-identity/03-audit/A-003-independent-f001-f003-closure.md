---
id: GOAL-032-w21-startup-db-identity
doc: audit-entry
record_id: A-003
source: independent
scope: A-001 F-001～F-003 关闭证据复审（finding-closure；非 S5 全目标关门）
verdict: pass
audit_type: finding-closure
status: recorded
parent: GOAL-001-design-implementation-conformance
created: 2026-08-22
updated: 2026-08-22
version: 0.1.0
---

# A-003 · 独立交叉审计 · A-001 F-001～F-003 关闭证据（2026-08-22）

- **source**：independent
- **auditor**：grok-4.6（Grok Build `/audit`）
- **类型** / **scope**：finding-closure · A-001 F-001 / F-002 / F-003（不含 A-001 F-004～F-007，不含 GOAL-032 全目标再关门）
- **verdict**：**pass**
- **完整意见**：本文件（未超 32 KiB，无附件）

## 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-010-design-implementation-conformance` |
| Root | `GOAL-001-design-implementation-conformance` |
| canonical | `docs/workspaces/workspace-010-design-implementation-conformance/` |
| 资料目录 | `none` |
| 日期 | 2026-08-22 |
| covered | A-001 F-001～F-003 原文要求；A-002 主张；D-002；E-003；`identity.go` classify/plan；双方言 `migrate` refuse/restore/adopt 臂；本轮复跑的定向单测 |
| excluded | 不重审 A-001 已通过主张；不把 A-001 F-004b/F-006/F-007 标闭合；不改 `status`/`progress`；不读取其他工作区；未重跑活进程 `go run ./cmd/server` |
| 信息项 | I-001 / I-002 仍 verified（S1）。D-002 收紧 I-001 **执行**语义（盖章所需对象），未假装成新的开放 I-00N |

工作区绑定：`workspace.md` 的 `id` / `root_goal` / `canonical_scope` 与本目标路径一致。本波未改 Profile / 模块矩阵 / Manifest。无共享资料引用。

## 成果（有证据）

本轮只核对关闭声明是否可重复。下列**已核对**：

| 主张 | 证据 |
|------|------|
| 完整指纹不再只看四表 | `identity.go` `completeLostLedgerTables` = 原四表 + `service_credentials` + `operation_log_session`；`lostLedgerLooksComplete` 缺任一则 false |
| catalog 再涨则指纹测试 fail closed | `completeFingerprintCatalogHead = 48`；`TestCompleteFingerprintTracksCatalogHead` 本轮 **PASS**（compiled catalog max 未超过 48） |
| 四表/缺头对象不再 restore | `classifyIdentity`：未完整且 `hasPostV1CatalogTables` → `lost-ledger-unsafe`；`planStartup` → `refuse`，reason 含 `identity=lost-ledger-unsafe` |
| PG 反例：jobs 时代库丢 ledger + 缺头表 | `TestPostgresMigrateRefusesIncompleteLostLedger` 本轮 **PASS**（未 skip）：Open 失败且 **不**建 `schema_migrations` |
| 全量丢 ledger 仍可 restore | `TestPostgresMigrateRestoresLostLedger` 本轮 **PASS**（盖章行数 = catalog 长，用户行保留） |
| post-v1 脏库不再走 adopt-then-pending | `hasPostV1CatalogTables`；`TestClassifyIdentity`「four tables without catalog head」→ `lost-ledger-unsafe`；`TestPlanStartup` unsafe → refuse |
| sqlite V-MIG-03 未废 | `TestMigrateFailClosedPartialBaseline` 本轮 **PASS**：users-only fail closed，无 ledger |
| 合同写明方言 v1 差 | D-002 §3；sqlite `fingerprintR2` 仍要求恰好 `{users, refresh_tokens}`；PG `migrateBaselinePG` 仍可对 users-only 补 ledger/`refresh_tokens` |
| 双方言 Execute 对 unsafe 是 refuse | `migrate.go` / `postgres.go` `actionRefuse` 直接 return，不 stamp、不 Apply |

本轮命令（`apps/api`）：

```text
go test ./internal/store/ -count=1 -timeout 180s -v -run "TestClassifyIdentity|TestPlanStartup|TestCompleteFingerprintTracksCatalogHead|TestMigrateFailClosedPartialBaseline|TestMigrateFailClosedForeignSQLite|TestPostgresMigrateRefusesIncompleteLostLedger|TestPostgresMigrateRestoresLostLedger|TestPostgresMigrateAdoptsLedgerlessR2|TestMigrateExistingR2DB"
```

全部 PASS（含 PG 三条，未 skip）。

## 对照成功标准（关闭复审）

| A-001 要求 | 状态 | 证据 |
|------------|------|------|
| F-001：不得仅四表整表盖章；核验对象或缩小盖章 / 绑定 catalog 头；补「四表在、后继缺」测试 | **达到**（D-002 选 catalog 头绑定） | 指纹 + unsafe refuse + PG 反例本轮绿；原四表因含 `jobs`/`operation_log` 也会走 unsafe，不会 restore |
| F-002：明确 partial 可执行子集；有 post-v1 catalog 表则 refuse | **达到**（相对 A-001 点名场景） | D-002；`lost-ledger-unsafe`；classify/plan/PG Open 一致 |
| F-003：消解 D-001 ↔ V-MIG-03 ↔ PG users-only 合同并留痕 | **达到** | D-002 选择保留 sqlite fail closed + PG 可 adopt；sqlite 测试语义未改 |

D-001 正文身份表仍写「四表 = complete」。D-002 写明 restore/partial **以 D-002 为准**；`01-decision.md` 索引指向 D-002。不把这份文档滞后当成 F-003 未闭合。

## 关闭证据表

| Finding | A-002 主张 | 本轮 | 说明 |
|---------|------------|------|------|
| A-001 F-001 | fixed | **fixed** | 原 high：静默跳过 v43–v48。现缺头对象 refuse；完整库丢 ledger 仍 restore。残余见本条 F-001（recommended），不重开本 finding |
| A-001 F-002 | fixed | **fixed** | 原 med：`roles`/`operation_log`/jobs 时代库走 adopt 撞 42P07。现该类 → unsafe refuse。残余见本条 F-002（recommended） |
| A-001 F-003 | fixed | **fixed** | 原 med：三岔合同。D-002 书面收窄；sqlite users-only 仍 fail closed |

无 P-004 residual / overruled。无意见冲突。

## Findings（本条新开）

### F-001 · `TestCompleteFingerprintTracksCatalogHead` 只锁版本号，不锁表名

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**

测试条件是 `compiled catalog max > 48`。若有人把常量改成新头版本、但忘记把新头表写入 `completeLostLedgerTables`，测试仍绿，F-001 类静默盖章会复发。非阻断：当前 max=48 与常量一致，且本轮 PG 反例锁住了「缺 `service_credentials`/`operation_log_session` 不得 stamp」。

### F-002 · `postV1CatalogTables` 是抽样，不是 compiled catalog 派生

- 严重度：**low**
- 建议：**recommended**
- 状态：**open**
- 关联：A-001 F-002（已 fixed，本条不重开）

列表：`roles` / `permissions` / `menu_items` / `operation_log` / `jobs` / `service_credentials` / `notifications`。顺序 apply 的丢失 ledger 库只要走到 v2+ 就会有 `roles` 等，会被 refuse。

未列入的 catalog 表（如 `records`、`dict_types`、`wallet_accounts`）若与我方 `users` 一起出现、且没有列表中任一项，仍会标 `ours-partial-no-ledger`。PG 随后 `CREATE TABLE` 可能 42P07。这不是 A-001 点名的 jobs/RBAC 现场路径，故保持 recommended。

## 必改项汇总

开放 required：**0**（相对 A-001 F-001～F-003）。

本条新 recommended：F-001、F-002。A-001 仍开放 recommended（本 scope 外）：F-004b（sqlite 完整丢 ledger Open）、F-006、F-007。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 台账状态 | 本轮 |
|----|------|----------|----------|------|
| I-001 | required | S1 | verified | D-002 执行语义（四表 + 头表）与代码一致；足以关闭 F-001 的「四表盖章」洞，不等于核验 catalog 全对象 |
| I-002 | required | S1 | verified | 仍沿用 `schema_migrations` |

无到期未答 required 信息项。无用户书面 `accepted-residual`。

## 与既有意见的异同

- **A-001**：conditional，3 required 开放。本条只复审这三条关闭证据，不推翻 A-001 对其它范围的描述。
- **A-002**：self 主张 F-001～F-003 fixed。本轮代码 + 定向测试（含 PG 反例实跑）支持该主张。
- 无 verdict 冲突需 P-004。

## 结论 + 建议给编排器/用户的下一步

A-001 三条 required 的关闭证据**充分、可重复核对**。本意见 **pass**。这不是 GOAL-032 `done`，也不是第二次全目标 close-out。

建议 `/govern`：

1. 把 A-001 F-001～F-003 记为 `fixed`（本条确认）。
2. 再决定是否做 S5 关门；剩余 recommended 不阻断 required 门禁，是否带进关门由编排器/用户定。
3. 不因本波改 VP-008 `go`。

## 声明

本意见不修改 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。响应由 `/govern` 处理。
