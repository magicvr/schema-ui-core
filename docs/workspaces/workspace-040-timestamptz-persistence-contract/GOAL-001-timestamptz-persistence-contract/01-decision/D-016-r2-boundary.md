---
id: D-016-r2-boundary
doc: decision-entry
status: accepted
parent: null
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
---

# D-016 · R2 阶段边界与信息门禁

## 决定的来源

- **R1 已关门**（`GOAL-002` `done · 4/4`，2026-09-20 用户书面确认；关门向 independent = A-046，开放 required = 0）。Root 检查点 R1 → `completed`，Root `progress: 1/3`。
- 按 `docs/architecture/workspace-protocol.md`，**新建阶段子目标前必须先在 Root 冻结该阶段边界与信息门禁**——本决策即 R2 的边界冻结。
- 既有约束：Root `D-004`（R2 按模块追加 migration，不修改 v1–v72）、child `D-017`（单 checksum / SQLite DDL 切片）、`D-018`（纯整数 SQL 行拷贝 + 无过渡期）、`D-019`（F-5 子女先行 + 跨 descriptor 两次重建）、`D-020`（可执行测试载体与重定向）、`D-021`（F-I-005 `accepted-residual` 的范围与复审触发）。

## 1. R2 目标

**双方言迁移 + Store 编解码**（Root 纲领路线图 R2 原话）：把 R1 冻结的合同落成可执行代码——15 个 conversion descriptor、共享时间 codec、各模块仓储的读写改造——并**首次记录真实 `MigrationChecksum`**。

## 2. R2 范围（in scope）

| # | 交付 | 载体 |
|--:|------|------|
| 1 | **SQLite 侧 15 个 conversion descriptor** 落码（v73–v87），按 `D-019` 的 F-5 子女先行与裸四步；含子表重建与回填 | `apps/api/modules/*/migration/` |
| 2 | **PG 侧 `ApplyPostgres` 显式 DDL**（禁止 `pgTimeColRe` 派生）；毫秒族用**整数拆分式**（原式已弃用） | 同上 |
| 3 | **canonical SQL 切片 + 真实 `MigrationChecksum` 记录**到 ledger；`transform_id` 按 `D-017` | `r1-c2-descriptor-ledger-v1.0-fc.md` §1/§2 |
| 4 | **共享时间 codec**（`apps/api/internal/temporal`，包名候选）：`FromUnix`/`FromUnixMilli`、fixed-6 formatter/parser、sentinel/NULL helper、微秒截断 | 新建内部包 |
| 5 | **仓储层读写改造**：各模块时间列 `int64` → `time.Time`/`sql.NullTime`；谓词按 exact SQL 表改造；**无过渡期**（`D-018`） | `apps/api/modules/*/store/`、`apps/api/internal/` |
| 6 | **测试改写**：`migrate_test.go` / `postgres_test.go` 的 v73+ 追加断言；**金额列断言拆分**（`wallet_accounts.balance_total` / `wallet_ledger_entries.amount_delta` 保持 `bigint`）；leftover 21 名补入 | `apps/api/internal/store/` |
| 7 | **边界测试重定向**：`apps/api/internal/w040contracttest/` 的用例改指真实迁移（`D-020`） | 既有测试包 |
| 8 | **`rebuildOperationLog` fail-closed 断言**（`D-019` §5 改动 1–3） | `apps/api/modules/operationlog/migration/` |
| 9 | **Backup Port 类型表面与 provider**（`D-019` §6 / C3 边界 §4.1）：`kernel.RecoveryPointPort` + `apps/api/internal/backup/`（含 `<recovery-artifact>` 校验、错误分类、restore harness）——**归属 M4 前**（用户 2026-09-20 裁决：**归 R2**，与 `D-021` residual 复审触发同批） | 新建 |
| 10 | **未发布 baseline 的发布前调整**：`D-014` 允许发布前有记录地拆分/调整 v73–v87 | — |

> **范围修正（用户 2026-09-20 裁决，响应 I-041-002）**：**公共 wire formatter 的实施（`apps/api/internal/handler/rfc3339.go` 由 milli 改固定 6 位）归 R3**，**不属 R2**。R2 只做 Store/持久化层（codec + 迁移 + 仓储），**不改 handler 的响应编码**。故上方第 1–8、10 项即 R2 的全部范围；第 9 项按 M4 前完成。

## 3. R2 **非目标**（out of scope，不得借 R2 实施）

- **不**重开 Root `D-004` 的归属裁决（不改成单一 platform migration；不修改历史 DDL）。
- **不**修改 v1–v72 的 canonical SQL / checksum / identity / Apply 语义。
- **不**引入 ORM、第三数据库、Redis/MQ/多实例/A3（Root 红线）。
- **不**把 `pgx`/SQLite 驱动类型泄漏到 handler/模块公共契约。
- **不**做 VP-020 的展示/输入时区回归矩阵（属 **R3**）与其用例 ID（recommended F-I-009）。
- **不**做备份的调度、鉴权、远端存储、保留策略、KMS/TLS、UI（C3 边界已排除）。
- **不**把 `accepted-residual` 或本地临时容器验证结果当作已验证事实。

## 4. R2 完成判据（checkpoint，用于 progress 派生）

| 检查点 | 判据 |
|--------|------|
| **M1** | 共享 codec 落码并有可执行单测（含 fixed-6、负毫秒 floor、999 ms 无进位、公元 9999 年、sentinel/ NULL） |
| **M2** | 15 个 descriptor 的 SQLite `Apply` + PG `ApplyPostgres` 落码，canonical SQL 与**真实 `MigrationChecksum` 已记录**；v1–v72 checksum 不变 |
| **M3** | 仓储读写与谓词改造完成，**双方言**回归通过（至少一条 PG 路径）；金额列断言保持 `bigint` |
| **M4** | `D-021` residual 三项全部完成，并经 **independent 复审**（`D-021` 复审触发）→ residual 关闭；R2 self + grok independent 关门审计通过 |

> `progress` 由 M1～M4 等权派生；**不**放行阶段、**不**关闭 finding、**不**推导 `done`。

## 5. R2 信息需求与门禁（P-005）

| ID | 级别 | 所需信息 / 问题 | 影响门禁 | 最晚需要阶段 | 状态 |
|----|------|-----------------|----------|--------------|------|
| I-041-001 | required | **Go codec 的公共 API 形态**：函数签名、错误分类、是否导出 `Truncate` 与 sentinel helper——须与 `D-018` 的 Go 对拍验收一致 | M1/M3 | M1 前 | **collecting**（在 `GOAL-003` 内定稿并落盘；未经用户裁决的技术细节由审计复审） |
| I-041-002 | required | **公共 wire formatter 的改造落点与是否属 R2** | M3 | M3 前 | **verified（用户裁决 2026-09-20）**：**归 R3**，R2 不动 handler |
| I-041-003 | required | **Backup provider 的 PG 侧可执行验证环境**：本机无 `psql`/`pg_dump`/`pg_restore`、无常驻 PG；`D-021` 只授权**临时**容器。R2 的 PG 回归是否需要常驻 PG 或 CI 方案 | M3/M4 | M3 前 | **open**（`GOAL-003` 之后、M3 前的子目标须先行关闭；可能触发 P-004） |
| I-041-004 | non-blocking | PG 15/16/17 跨版本 `pg_restore` 兼容矩阵 | R3 | R3 前 | deferred（`I-040-003` 已登记的 R3 侧 residual） |

- **到期 open required 阻断对应门禁**；`deferred` 保留级别并须在 R3 前复核。
- 本表**新增**编号段 `I-041-NNN`（Root 级），不与 child `I-040-NNN` 混用。

## 6. 渐进子目标策略

- 按 `D-004` 与 `D-019`，R2 **不预创建全部子目标**；先按本决策冻结边界，随后**按检查点渐进立项**。
- 首个候选子目标：**`GOAL-003-r2-codec-and-descriptor-m1-m2`**（codec + 15 descriptor + checksum 记录），对应 M1/M2。
- M3 的仓储改造与 M4 的关门审计**待 M1/M2 完成后再立项**，以避免过早细粒度拆分（AGENTS §6 P-001）。
- **子目标命名与 slug 须经用户确认后落盘**（AGENTS §11 硬约束：禁止静默默认 slug）。

## 未选方案

- **一次性预创建 R2 全部子目标**：违反 P-001「按阶段创建」，且 M3 的范围取决于 M1/M2 的实际形态。未采用。
- **R2 不设检查点、直接以「迁移完成」为唯一判据**：无法派生可核对 progress，也无法在 M2 处设门禁。未采用。

## 影响与边界

- 本决策只冻结 **R2 的边界与门禁**；**不**放行任何 schema 变更——R2 的每次迁移仍须遵守 `D-017`–`D-021` 与 Root `D-004` 的不变量。
- 本决策**不**改变 `D-014`（未发布 baseline）与 `D-021`（F-I-005 residual 的范围与复审触发）。
- 三个 `required` 信息项（I-041-001～003）**在对应门禁前关闭**；其中 I-041-002 可能触发 P-004（范围归属），届时**询问用户**。
