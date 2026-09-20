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

> **状态：`freeze-candidate`。** 本文件响应 **F-I-004** 的关闭要求（A-027 §F-I-004 / A-029 / A-032 复述）：
> 1. **唯一区分**（a）转换前 rollback snapshot/dump（旧合同）与（b）转换后 `CreateRecoveryPoint` artifact（新合同 + restore 校验）；
> 2. **PG restore 不得**再把预转换 `<artifact>` 当目标形状校验输入 → 给转换后 PG artifact 一个**独立 token**；
> 3. 写出 `CreateRecoveryPoint` 的**包路径**与 **before/after 调用点**；
> 4. restore-to-new-db **harness**（可执行，不只是命令模板）；
> 5. **旧 dump 不得通过新合同校验**。
>
> 本文件之前不存在（A-030 **F-I-020.1** 曾点名转换合同引用了它却未落盘；该悬空引用已在 `A-031` 修正为「尚未落盘」）。本文件即该缺口的落盘。

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
| **A** | **Pre-conversion rollback artifact** | 转换**前** | **旧合同**：SQLite 时间列为 `INTEGER`；PG 为 `BIGINT` | 升级失败的**回滚/演练** | **否**（形状是旧合同，必须失败） |
| **B** | **Post-conversion RecoveryPoint artifact** | 转换**后**（同一批次成功提交之后） | **新合同**：SQLite fixed-6 UTC `TEXT`；PG `timestamptz(6)` | 目标合同的**可恢复点** | **是**（唯一可满足者） |
| **C** | **Batch boundary snapshot**（现状机制） | 每条 pending 迁移前 | 旧合同（随批次推进而不同） | 现有 per-migration 回滚点 | **否** |

**硬规则**：

1. **A 与 C 永不作为**目标形状校验的输入；任何把 A/C 喂给「全部 90 列须为 TEXT / `timestamptz(6)`」校验的路径**必须失败**，且失败要作为**正向证据**保留（见 §6）。
2. **B 只能由转换后的目标库产生**；生成 B 之前必须先有一次**成功提交**的转换。
3. A/C 与 B 的**产物身份必须可区分**（token，见 §3），不允许用同一个占位符指代两者。

## 3. PG artifact 独立 token（关闭要求 2）

**问题（A-029 点名）**：`r1-c3-backup-restore-runbook-v0.1.md` L32 把预转换 dump 定为 rollback artifact、L33 说转换后 BackupService 才能让 `CreateRecoveryPoint` 成功，**但 L34–L39 的 restore 命令仍绑定步骤 1 的 `<artifact>`**，而 L41 又要求 restore 后 `information_schema` 为 `timestamp with time zone` precision 6 → **预转换 dump 形状是 `BIGINT`，该 restore 必失败**。

**修正后的 token 约定**：

| token | 指代 | 生成命令 | 校验期望 |
|-------|------|----------|----------|
| `<rollback-artifact>` | §2 的 **A**（预转换 dump） | `pg_dump -F c --no-owner --file <rollback-artifact> <source-dsn>` | 仅**回滚演练**：restore 后**允许**为 `bigint`；**不得**对其断言 `timestamptz(6)` |
| `<recovery-artifact>` | §2 的 **B**（转换后 dump） | 转换提交成功后，对**同一库**执行 `pg_dump -F c --no-owner --file <recovery-artifact> <source-dsn>` | **必须**为 `timestamptz(6)` 精度 6；全部 90 列在场；catalog/checksum 匹配 |

- **`CreateRecoveryPoint` 只接受 `<recovery-artifact>` 路径**；`<rollback-artifact>` 传入时必须在**形状校验阶段**返回错误（不得返回未校验产物）。
- SQLite 侧同构：`<rollback-snapshot>`（预转换）与 `<recovery-snapshot>`（转换后）；转换后快照**不得**复用 `snapshotBeforePending` 的文件（那是 per-migration 批次点，§2 的 C）。

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
- **包名上述为候选**；接受与否由 independent 复审判定。用户 `D-010`/`D-007` 已定「kernel 只暴露 `CreateRecoveryPoint`，其余编排内部」。

### 4.2 before/after 调用点（关闭要求 3）

| # | 调用点 | 时机 | 调用 | 失败行为 |
|--:|--------|------|------|----------|
| 1 | `Store.migrate` → `applyPending`（`internal/store/migrate.go:81-103`） | 转换批次**之前** | 生成 **A/C**（rollback artifact）；**不**调用 `CreateRecoveryPoint` | 快照失败 → 批次不开始 |
| 2 | `Store.migrate` → `applyPending` **成功返回之后** | 转换**已提交** | 调用 `RecoveryPointPort.CreateRecoveryPoint(ctx, {Dialect, SourceID, CatalogVersion, TimeContract, …})` | 失败 → **报错但转换已提交**；不得回滚已提交的转换（回滚走 A/C） |
| 3 | `Store.applyMigration`（`:108-132`） | 每条迁移的单事务内 | **禁止**调用（事务内不得做外部 artifact I/O） | — |

- 调用点 2 的位置**必须在 `verifyIntegrity()`（`:102`）之后**，即在批次完整校验通过后才产出 B。
- **禁止**在 `applyMigration` 事务内调用（§4.2 #3）：`pg_dump` / `VACUUM INTO` 都是外部 I/O，会拉长事务并与 SQLite 单连接/写锁冲突。

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
4. 样本 round-trip：秒、毫秒、sentinel `0`、可空缺失、负值非法、fixed-6 词法序各至少一例；
5. catalog/checksum 与源一致；retired `records` 表**不存在**；
6. 失败时产出可核对的错误分类，**不得**返回成功产物。

**PG 版本兼容**：harness 记录 server/client major 版本；跨主版本 `pg_dump`/`pg_restore` 组合须显式记录为**已知不支持**或**已测通过**，不得默认通过。

## 6. 旧 dump 不得通过新合同校验（关闭要求 5）

- **正向证据要求**：`TestLegacyArtifactMustFail` 必须以**预转换 artifact（A/C）**为输入跑第 5 节的形状校验，并断言其**失败**。该失败是**期望结果**，须收录为证据（证明新旧合同的区分真实生效），而非当作缺陷。
- C3 提交时必须同时具备：(i) B 通过全部校验；(ii) A/C 被同一校验**拒绝**。只有 (i) 不构成 C3 冻结证据。

## 7. 与 `D-019` 的接口（FK 与备份的边界）

- F-5 的失败路径：descriptor 事务内任一步失败 → 整事务回滚 → 表层恢复为**旧合同**；此时 **A/C 仍有效**，B **不存在**。
- 事务**无法**回滚的情形（进程被杀、磁盘故障）由 A/C 承担恢复；此路径下 B 亦不存在。
- **不重叠**：本文件不规定表重建顺序（属 `D-019`）；`D-019` 不规定 artifact 生成与校验（属本文件）。

## 8. 本文件**未**闭合的项（诚实边界）

1. **无任何实现**：`apps/` 中不存在 §4.1 的包与类型；本文件是设计。
2. **PG harness 依赖真实 PG**：本机无 PG 时只能 `t.Skip`；跨版本兼容矩阵需单独信息项（当前 `I-040-003` 仍 `collecting`）。
3. **`ArtifactRef` 的清理/保留语义**（runbook Open 项）仍待定；本文件只规定 A/C 与 B 的区分与 token，不规定保留期。
4. 包名 `internal/backup`、类型名与 `TimeContract` 取值均为**候选**，待复审接受。

## 9. 声明

- 本文件 `status: freeze-candidate`；**不**声称任何实现；`apps/` 未修改。
- 是否构成 **F-I-004** 的合法闭合，由 independent 复审判定；本编排器不自证。
- 引用 `D-0NN` 一律限定 Root/child。
