---
id: r1-c3-backup-recovery-boundary-v1.0-fc
doc_type: design-attachment
title: R1 C3 备份/回滚边界 v1.0（冻结候选）
status: freeze-candidate
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-002-r1-contract-and-denominator-freeze
version: 1.0.0
---

# R1 C3 备份/回滚边界 v1.0（冻结候选）

> **状态：`freeze-candidate`。** 本文件响应 **F-I-004** 的关闭要求（A-027 §F-I-004 / A-029 / A-032 复述，A-040 §G 更新）：
> 1. **唯一区分**（a）转换前 rollback artifact（旧合同）与（b）转换后 `CreateRecoveryPoint` artifact（新合同 + restore 校验）；
> 2. **PG restore 不得**再把预转换 `<artifact>` 当目标形状校验输入 → 双 token；
> 3. 写出 `CreateRecoveryPoint` 的**包路径**与 **before/after 调用点**（**含 PG 对称锚点**）；
> 4. restore-to-new-db **harness**（可执行，不只是命令模板）；
> 5. **旧 dump 不得通过新合同校验**，且必须钉**错误分类**。
>
> **本文件的 C3 权威范围**（响应 A-040 §G 第 2 项）：产物区分、token、包路径、调用点、harness、反向断言、错误分类以本文件为**唯一权威**。`r1-c3-backup-restore-runbook-v0.1.md` 已标 `superseded`；`r1-backup-port-contract-draft-v0.1.md` 的权威仅限 kernel Port 的**类型表面与后置条件**，其 L60/L71 已收口（不得再把 `snapshotBeforePending` 当 RecoveryPoint 源）。

## 1. 现状事实（实测，2026-09-20）

| 事实 | 证据 |
|------|------|
| `apps/` 中**不存在** `CreateRecoveryPoint` / `BackupService` / `RecoveryPoint` 任何标识符 | 全仓 `grep` `apps/api/**/*.go` → **0 匹配** |
| `kernel.Store` **没有** Backup 接口 | `kernel/store.go:30-47` 仅 `Dialect`/`Run`/`Ping`/`Close`/`WasFresh`/`MarkSystemDataReady`/`SystemDataReady` |
| `apps/api/internal/temporal` 包**不存在** | 目录探测 → absent |
| SQLite 现有快照机制 | `internal/store/migrate.go:81-103` `applyPending` 对 `version >= 2` 的每条 pending 迁移调 `snapshotBeforePending(version)`；实现 `:279-300`，`fresh` / `:memory:` / 无数据时返回 `nil` |
| SQLite 快照的语义 | 升级批次内**每条** per-migration 的 `VACUUM INTO`（`migrate.go:82-96` 注释与实现）——**是 rollback 点，不是 RecoveryPoint** |
| PG 现有失败路径 | 事务 rollback + ledger 写 `applied_at`（`postgres.go:154-171`）；**无** `pg_dump` |
| runner 批次后校验 | `PRAGMA foreign_key_check`（`migrate.go:349`）、`integrity_check`（`:367`） |

> **结论**：C3 的全部内容目前**均为设计**，`apps/` 中无一实现。本文件不得被读作任何已实施证据。

## 2. 三类产物的唯一区分（关闭要求 1、5）

| 代号 | 名称 | 产生时机 | 合同形状 | 用途 | **能否**满足 `CreateRecoveryPoint` 后置条件 |
|------|------|----------|----------|------|---------------------------------------------|
| **A** | **Pre-conversion rollback artifact** | 转换批次**开始之前**（**循环外，一次**） | **旧合同**：SQLite 时间列为 `INTEGER`；PG 为 `BIGINT` | 整批失败的**回滚/演练** | **否** |
| **B** | **Post-conversion RecoveryPoint artifact** | 转换批次**成功提交并校验之后** | **新合同**：SQLite fixed-6 UTC `TEXT`；PG `timestamptz(6)` | 目标合同的**可恢复点** | **是**（唯一可满足者） |
| **C** | **Per-migration batch-boundary snapshot**（现有 `snapshotBeforePending` 机制） | 批次**循环内**，每条 `version >= 2` 的 pending 迁移**之前** | **混合形状**：随批次推进而变化——批次早期为旧合同，中途含**部分**新合同列，批次末为全部新合同 | 单条迁移失败的**细粒度回滚点** | **否**（形状不稳定，见硬规则 1） |

**A 与 C 必须分开**（响应 A-040 §G 第 4 项）：A 是**批次开始前的一次性**旧合同产物；C 是**循环内 per-migration** 的产物，其形状**随批次推进而混合**。二者**不可合并为「批次前一次」**，也不可互替。

**硬规则**：

1. **A 与 C 永不作为**目标形状校验的输入；任何把 A/C 喂给「全部 90 列须为 TEXT / `timestamptz(6)`」校验的路径**必须失败**，且该失败要作为**正向证据**保留（见 §6）。
2. **B 只能由转换后的目标库产生**；生成 B 之前必须先有一次**成功提交并校验**的转换批次。
3. A/C 与 B 的**产物身份必须可区分**，且该区分必须是**机械可核对的**（见 §3.1），不允许仅靠文件名约定或人工判断。

## 3. PG artifact 独立 token（关闭要求 2）

**问题（A-029 点名）**：`r1-c3-backup-restore-runbook-v0.1.md` L32 把预转换 dump 定为 rollback artifact、L33 说转换后 BackupService 才能让 `CreateRecoveryPoint` 成功，**但 L34–L39 的 restore 命令仍绑定步骤 1 的 `<artifact>`**，而 L41 又要求 restore 后 `information_schema` 为 `timestamp with time zone` precision 6 → **预转换 dump 形状是 `BIGINT`，该 restore 必失败**。

**修正后的 token 约定**：

| token | 指代 | 生成命令 | 校验期望 |
|-------|------|----------|----------|
| `<rollback-artifact>` | §2 的 **A**（预转换 dump） | `pg_dump -F c --no-owner --file <rollback-artifact> <source-dsn>` | 仅**回滚演练**：restore 后**允许**为 `bigint`；**不得**对其断言 `timestamptz(6)` |
| `<recovery-artifact>` | §2 的 **B**（转换后 dump） | 转换提交成功后，对**同一库**执行 `pg_dump -F c --no-owner --file <recovery-artifact> <source-dsn>` | **必须**为 `timestamptz(6)` 精度 6；全部 90 列在场；catalog/checksum 匹配 |

- **`CreateRecoveryPoint` 只接受 `<recovery-artifact>` 路径**；`<rollback-artifact>` 传入时必须在**形状校验阶段**返回错误（不得返回未校验产物）。
- SQLite 侧同构：`<rollback-snapshot>`（§2 的 A，批次前一次）与 `<recovery-snapshot>`（§2 的 B，转换后）；**C 类**（`snapshotBeforePending` 的 per-migration 快照）既不是 A 也不是 B，**不得**用作任一侧的输入。

### 3.1 C → B 的机械身份（响应 A-040 §G 第 7 项）

仅靠文件名政策不足以区分 A/C/B。冻结要求：

| 产物 | 机械身份（必须落进 `RecoveryPoint` 元数据或 artifact 旁证） |
|------|-------------------------------------------------------------|
| A / C | **不得**生成 `RecoveryPoint`；其记录只含 `artifact_path` + `catalog_version` + `contract_shape = legacy`（SQLite: `INTEGER`；PG: `bigint`） |
| B | 元数据必须含 `TimeContract`（如 `vp040-timestamptz-0.1`）+ `CatalogVersion` + `ChecksumSet` + `Verification`（`VerificationSummary` 四项）+ 生成时的**批次末 version**；且其 `contract_shape` 由**对 restore 目标的实测**得出，不是由文件名推断 |

**判定规则**：`CreateRecoveryPoint` 在校验阶段**实测 restore 目标**的时间列物理类型；只要实测形状不是新合同（或 90 列不齐），一律归类为 **`TimeContractMismatch`** 并返回错误——**不得**因路径/文件名看起来像 recovery 就通过。这样 C→B 的区分是**实测驱动**而非命名驱动。

## 4. `CreateRecoveryPoint` 的包路径与调用点（关闭要求 3）

### 4.1 包路径（提议，冻结候选）

```text
apps/api/kernel/store.go                     -- 仅接口与中性类型（无驱动类型）
  type RecoveryPointPort interface { CreateRecoveryPoint(ctx, RecoveryPointRequest) (RecoveryPoint, error) }
  type RecoveryPointRequest / RecoveryPoint / VerificationSummary   -- 见 r1-backup-port-contract-draft-v0.1.md
apps/api/internal/backup/                     -- 编排与 provider（内部，不进 kernel）
  service.go        -- Service 实现 RecoveryPointPort；选择 provider、串起 create→restore→verify
  provider_sqlite.go-- SQLite provider（native snapshot + restore-to-new-file + checks）
  provider_pg.go    -- PG provider（pg_dump -F c + createdb + pg_restore --exit-on-error --no-owner）
  verify.go         -- 形状/样本/checksum 校验（含「A/C 必须失败」的反向断言）
```

- `kernel` **不得**出现 `pgx`/SQLite 驱动类型（Root 红线）；`RecoveryPoint` 的 `ArtifactRef` 是**不透明**引用。
- **包名上述为候选**；接受与否由 independent 复审判定。「kernel 只暴露 `CreateRecoveryPoint`、其余编排内部」对位 **Root** `D-007` + **Root** `D-010`（child `D-006` + child `D-009`）——响应 A-040 的编号限定要求。

### 4.2 before/after 调用点（关闭要求 3）

**SQLite 侧**：

| # | 调用点 | 时机 | 调用 | 失败行为 |
|--:|--------|------|------|----------|
| 1 | `Store.migrate` → `applyPending`（`internal/store/migrate.go:81-103`）**循环开始之前** | 批次前**一次** | 生成 **A**（旧合同、批次前） | 失败 → 批次不开始 |
| 2 | `applyPending` **循环内**，每条 `version >= 2` 的 pending 迁移之前（`:92-97`） | per-migration | 现有 `snapshotBeforePending` → 生成 **C**（**不得**当 A 或 B） | 失败 → 该迁移不执行 |
| 3 | `applyPending` **成功返回之后**（即 `verifyIntegrity()` `:102` 之后） | 批次已提交并校验 | 调 `RecoveryPointPort.CreateRecoveryPoint(...)` → 产出 **B** | 见 §4.3 |
| 4 | `Store.applyMigration`（`:108-132`） | 每条迁移的单事务内 | **禁止**调用（事务内不得做外部 artifact I/O） | — |

**PG 侧（对称锚点，响应 A-040 §G 第 3 项）**：

| # | 调用点 | 时机 | 调用 | 失败行为 |
|--:|--------|------|------|----------|
| 1 | `postgres.migrate` → `applyPendingPG`（`internal/store/postgres.go:120-127`）**循环之前** | 批次前一次 | 生成 **A`**（预转换 `pg_dump -F c`） | 失败 → 批次不开始 |
| 2 | `applyPendingPG` 循环内，每条 pending 迁移之前 | per-migration | 生成 **C`**（该迁移前的 dump） | 失败 → 该迁移不执行 |
| 3 | `applyPendingPG` **成功之后** | 批次已提交 | 调 `CreateRecoveryPoint` → 产出 **B`**（转换后 `<recovery-artifact>`） | 见 §4.3 |
| 4 | `applyMigrationPG`（`postgres.go:154-173`）事务内 | — | **禁止**调用 | — |

> **PG 侧的现状缺口（A-040 §G 第 3 项）**：`applyPendingPG` 现行**没有** snapshot 步骤，也**没有**对称的 `verifyIntegrity()`（SQLite 侧有 `migrate.go:345-367`，PG 侧缺失）。故 §4.2 的 PG 三处调用点是**新增设计**，R2 落码时须同时补 PG 的形状校验入口。

### 4.3 调用点 2/3 失败后的重试（响应 A-040 §G 第 5 项）

**问题**：B 生成失败时（如 PG 不可用），SQLite 路径会把进程置于 `actionNoop`（`migrate.go:52-53`，只跑 `verifyIntegrity()`），**不重试** `CreateRecoveryPoint` → **B 可能永不补创**；PG 路径（`postgres.go:94-95`）连 `verifyIntegrity` 都没跑。

**冻结要求**：

1. **B 缺失必须是可检测状态**，不是静默状态：启动分类（`planStartup` / PG 对应物）在「catalog 已到目标版本（含 v73+）但**无 B 记录**」时，必须产出**显式**动作（而非 `actionNoop`）。
2. 该动作**至少**要：跑形状校验；若形状已是新合同且无 B → **补创 B**（即重试 `CreateRecoveryPoint`）；若形状仍是旧合同 → 说明转换未完成，走**回滚**路径（A/C）。
3. **重试必须有界**：每次启动最多尝试一次补创，失败则记录并可继续服务（沿用现行 fail-closed 立场：不因补创失败而拒绝启动，除非用户另有裁决）；连续失败须留可核对记录。
4. **禁止**把「无 B」当作「已满足 RecoveryPoint 门禁」。

## 5. restore-to-new-db harness（关闭要求 4）

**现状**：runbook 给的是**命令模板**，不是可执行 harness（A-027 点名「命令模板 ≠ 可执行脚本」）。

**harness 规格（冻结候选）**：

```text
apps/api/internal/backup/restore_harness_test.go        -- 测试侧 harness（R2 落码）
  TestSQLiteRestoreToNewDB(t)      -- 从 <recovery-snapshot> restore 到新文件
  TestPGRestoreToNewDB(t)          -- 从 <recovery-artifact> restore 到新库（PG 不可用时 t.Skip 并记录 skip 原因，不得静默通过）
  TestLegacyArtifactMustFail(t)    -- 反向断言：A/C 不得通过新合同校验
```

harness 必须完成的断言（两侧同构）：

1. restore **到新目标**（新 SQLite 文件 / `createdb` 新库），**不**就地 restore；
2. SQLite：`PRAGMA integrity_check` = `ok`、`PRAGMA foreign_key_check` 无行；PG：`information_schema` 时间列 = `timestamp with time zone` 且 `datetime_precision = 6`；
3. 全部 **90 个时间列**在场且类型正确（分母取 `r1-time-column-inventory-v0.3.md`）；
4. 样本 round-trip（**仅对 B**）：秒、毫秒、sentinel `0`、可空缺失、fixed-6 词法序各至少一例。
   > **更正（响应 A-040 §D 末段）**：「负值非法」**不属于 B 的正向 round-trip**——转换是 fail-closed，B 内不应再有负值。负值断言应落在 `TestLegacyArtifactMustFail`（A/C）或转换预检（`m0`）上。
5. catalog/checksum 与源一致；retired `records` 表**不存在**；
6. 失败时产出**可核对的错误分类**（见 §5.1），**不得**返回成功产物。

**PG 版本兼容**：harness 记录 server/client major 版本；跨主版本 `pg_dump`/`pg_restore` 组合须显式记录为**已知不支持**或**已测通过**，不得默认通过。

### 5.1 错误分类表（响应 A-040 §G 第 6 项）

`TestLegacyArtifactMustFail` **不得**只断言 `err != nil`——那样缺文件、权限错误也会「绿」。冻结的错误分类与对应触发条件：

| 分类 | 触发条件 | 必须**不**归入此类 |
|------|----------|-------------------|
| `TimeContractMismatch` | restore 目标的实测时间列物理类型不是新合同（SQLite 非 fixed-6 `TEXT`；PG 非 `timestamptz(6)`） | 这是 A/C 被拒的**期望**分类 |
| `TemporalColumnSetIncomplete` | 实测时间列数量 ≠ 90 | — |
| `ChecksumMismatch` | catalog/checksum 与源不一致 | — |
| `ArtifactNotFound` | artifact 路径不存在 | **不得**用于 A/C 反向断言（否则缺文件也会绿） |
| `ArtifactUnreadable` | 权限/损坏导致无法读取 | 同上 |
| `ToolFailure` | `pg_dump`/`pg_restore`/`VACUUM INTO` 非零退出 | — |

**反向断言的硬要求**：`TestLegacyArtifactMustFail` 必须断言错误分类**属于** `{TimeContractMismatch, TemporalColumnSetIncomplete}`，且**明确断言不属于** `{ArtifactNotFound, ArtifactUnreadable, ToolFailure}`。证据落点：**测试名 + 观测到的分类**即足够，无需另发明附件格式（响应 A-040 §E）。

## 6. 旧 dump 不得通过新合同校验（关闭要求 5）

- **正向证据要求**：`TestLegacyArtifactMustFail` 必须以**预转换 artifact（A/C）**为输入跑第 5 节的形状校验，并断言其以 §5.1 的 `TimeContractMismatch`（或 `TemporalColumnSetIncomplete`）**失败**。该失败是**期望结果**，须收录为证据（证明新旧合同的区分真实生效），而非当作缺陷。
- C3 提交时必须同时具备：(i) B 通过全部校验；(ii) A/C 被同一校验以**正确分类**拒绝。只有 (i) 不构成 C3 冻结证据。

## 7. 与 `D-019` 的接口（FK 与备份的边界）

- F-5 的失败路径：descriptor 事务内任一步失败 → 整事务回滚 → 表层恢复为**旧合同**；此时 **A/C 仍有效**，B **不存在**。
- 事务**无法**回滚的情形（进程被杀、磁盘故障）由 A/C 承担恢复；此路径下 B 亦不存在。
- **不重叠**：本文件不规定表重建顺序（属 `D-019`）；`D-019` 不规定 artifact 生成与校验（属本文件）。

## 8. 本文件**未**闭合的项（诚实边界）

1. **无任何实现**：`apps/` 中不存在 §4.1 的包与类型；本文件是设计。
2. **PG harness 依赖真实 PG**：本机无 PG 时只能 `t.Skip`；跨版本兼容矩阵需单独信息项（当前 `I-040-003` 仍 `collecting`）。
3. **`ArtifactRef` 的清理/保留语义**（runbook Open 项）仍待定；本文件只规定 A/C 与 B 的区分与 token，不规定保留期。
4. 包名 `internal/backup`、类型名与 `TimeContract` 取值均为**候选**，待复审接受。

> **2026-09-20 补充（响应 A-040 §G）**：A-040 点名的其余剩余项已在本版收口——① 旧 runbook 已标 `superseded` 并就地标注两处错误点（§G1）；② conversion contract §0 的「尚未落盘」已收回，Port draft L60 已删「`snapshotBeforePending` or」、L71 已标已解决（§G2）；③ PG before/after 调用点已写出并标明 `applyPendingPG` 缺 snapshot 与 `verifyIntegrity` 的现状缺口（§G3）；④ A 与 C 已在 §2 分开并纠正 C 的混合形状表述（§G4）；⑤ 调用点失败重试见 §4.3（§G5）；⑥ 反向断言错误分类见 §5.1（§G6）；⑦ C→B 机械身份见 §3.1（§G7）。**这些是否构成闭合由 independent 复审判定。**

## 9. 声明

- 本文件 `status: freeze-candidate`；**不**声称任何实现；`apps/` 未修改。
- 是否构成 **F-I-004** 的合法闭合，由 independent 复审判定；本编排器不自证。
- 引用 `D-0NN` 一律限定 Root/child。
