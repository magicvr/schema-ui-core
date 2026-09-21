---
id: A-040-r1-independent-e033-c3-boundary-fi004
doc_type: goal-audit-entry
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-20
scope: GOAL-002-r1-contract-and-denominator-freeze · ad-hoc design-plan + finding-closure · E-033 / commit e2c0dac2 · attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md 对照 A-032/A-027/A-029 F-I-004 关闭要求原文 · 不是实施审计 · freeze-candidate ≠ 已实施
verdict: conditional
open_required: 3
status: recorded
created: 2026-09-20
updated: 2026-09-20
parent: GOAL-001-timestamptz-persistence-contract
version: 0.1.0
---

# A-040 · R1 independent · E-033 / C3 备份回滚边界对照 F-I-004

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：ad-hoc + finding-closure（用户指定核对 `attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md`，commit `e2c0dac2`，从未被 independent 复审，能否合法闭合 **F-I-004**。对照基线 A-032 §F-I-004 / A-027 / A-029 关闭要求原文。开放 required = 3：F-I-002 / F-I-004 / F-I-005。附件为设计材料，不是实施证据）。
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。
- 被审目标：`GOAL-002-r1-contract-and-denominator-freeze`。
- **未读其他工作区作为审计上下文**。未改 Charter / VP / Goal `status` / 检查点 / `progress` / 方案正文 / goal-tree / `apps/`。
- `git show --stat e2c0dac2`：6 个文件均在 workspace-040（`02-execution.md`、`E-033`、`03-audit.md`、`A-038`、`A-039`、C3 边界附件）；`apps/` 无变更。
- **证据窗口**：以 `e2c0dac2` 已提交的 C3 边界文件 + 现行 runbook / Port draft / conversion contract §0 + 本机 `apps/api` 源码为准。不把 freeze-candidate 当实施证据。

## 核对方法

1. 通读 C3 边界 v1.0-fc 全文、`r1-c3-backup-restore-runbook-v0.1.md`、`r1-backup-port-contract-draft-v0.1.md`、conversion contract §0、A-027/A-029/A-032 F-I-004 关闭要求原文、A-038/A-039/E-033、child `D-019`。
2. 对位源码：`kernel/store.go`；`internal/store/migrate.go`（`migrate` / `applyPending` / `applyMigration` / `snapshotBeforePending` / `verifyIntegrity`）；`internal/store/postgres.go`（`migrate` / `applyPendingPG` / `applyMigrationPG`）。
3. 全仓 `apps/**/*.go` 检索 `CreateRecoveryPoint` / `BackupService` / `RecoveryPoint`；目录探测 `internal/temporal`、`internal/backup`。
4. 按用户 A–H 逐项独立判定，不采信编排器转述。

## 成果（有证据）

1. **本轮未实施 Backup/RecoveryPoint。** `e2c0dac2` `--stat` 无 `apps/**`。边界文件 `status: freeze-candidate`。E-033 / A-039 自承未闭合任何 required。本审同意该实施边界。
2. **C3 边界文件首次进入证据窗。** A-038 正确把它排除在 `d1fdb4cc` 之外；本轮 `e2c0dac2` 将其纳入。这是 F-I-004 自 A-027 以来第一次有对位关闭要求的专用载体，**构成收窄，不构成闭合**。
3. **§2 三类产物在新文件内部把 (a)/(b) 写清楚了**（见 A）。A = 预转换旧合同 rollback；B = 转换后 RecoveryPoint；C = 现有 per-migration `VACUUM INTO`。三条硬规则方向正确。
4. **§4.1 包路径达到可落码具体度**（见 C）：`kernel` 仅中性 `RecoveryPointPort` + 三类；编排在 `apps/api/internal/backup/`。与 Root D-007 / Root D-010（child D-006 / D-009）同向。本审**接受**该候选包路径。
5. **§5 harness 规格超过「命令模板」**（见 D）：给出测试文件名、三个测试函数、六项同构断言、PG skip 须记录原因。作为 R1 设计冻结的 harness 合同足够；不是可执行脚本，也不应被当成 R2 已跑通。
6. **§7 与 `D-019` 不重叠**（见 H）：F-5 表序属 `D-019`；artifact 生成/校验属本文件。F-5 事务失败 → 旧合同、A/C 有效、B 不存在，与 `applyMigration` 单事务结构一致。
7. **F-I-004 整条仍不能闭合**（见 G）：冻结包仍有活的 `<artifact>` 矛盾；PG 调用点未写到 `postgres.go`；反向断言未钉错误分类；调用点 2 失败后 `actionNoop` 无重试。

## 对照成功标准（若适用）

| 标准 | 状态 | 证据 |
|------|------|------|
| C3 原地转换 / 备份回滚冻结 | **仍不可冻结** | 新文件收窄 F-I-004；整条仍 open |
| F-I-004 关闭要求（A-029/A-032 原文） | **未满足** | 见 G 全部剩余项 |
| C2 / C4 / R2 放行 | **未满足** | F-I-002/005 未触及；freeze-candidate ≠ 实施 |

## 对用户 A–H 的直接判定

### A. 三类产物是否真正唯一

**新文件内部：方向成立，(a)/(b) 已分开写。冻结包级：尚未唯一。**

| 代号 | 新文件主张 | 本审 |
|------|------------|------|
| A | 转换前、旧合同、不得满足 `CreateRecoveryPoint` 后置条件 | **接受。** 形状是 INTEGER/BIGINT；§5 全 90 列 TEXT/`timestamptz(6)` 校验必失败 |
| B | 转换提交后、新合同、唯一可满足者 | **接受为定义。** 须由转换后目标库产生 |
| C | 现有 per-migration `snapshotBeforePending`，不得当 B | **部分接受。** 循环调用属实（`migrate.go:92-97`）；表内「合同形状 = 旧合同」不准确——C 随批次推进，转换中途/转换后的 C 可含部分或全部新合同列 |

**仍可把 A/C 当 B 用的路径：**

1. **冻结包仍授权旧 runbook 的单一 `<artifact>`**（见 B）。这是 A-029 原文点名的路径，**仍然活着**。
2. **Port draft L60** 仍允许 SQLite provider「use existing `snapshotBeforePending`/`VACUUM INTO` family **or** a C3-specific native snapshot」。这与硬规则 1（C 永不作为目标形状校验输入）直接矛盾。
3. **转换后的 C** 若被当作 `<recovery-snapshot>` 喂给同一校验，形状/catalog/checksum **可以过**。硬规则 1 是政策，不是机械身份。§3 只禁止「复用文件」，未规定 artifact 元数据必须声明 role（rollback vs recovery）与 `TimeContract`。文件名 token 可被改名/误传。
4. 硬规则 2（B 只能由转换后目标库产生）+ 硬规则 3（token 可区分）**足以阻断 A→B**（旧形状必失败），**不足以单独阻断转换后 C→B**。

**三条硬规则对 A-029 的 (a)/(b) 主缺口：在新文件内足够；在冻结包内不够。** 主缺口仍是「预转换 dump 当目标形状校验输入」，靠形状校验 + 双 token 可堵，前提是 runbook/Port 不再写反。

### B. PG token 修正是否成立；旧 runbook 是否仍矛盾

**新文件 §3 的双 token 设计成立，且对位 A-029 原文。旧 runbook 现行正文仍有该矛盾。这是 F-I-004 剩余项，不是旁注。**

本审逐行核对 `attachments/r1-c3-backup-restore-runbook-v0.1.md`（`status: proposed`，**未被标 superseded**）：

| 位置 | 现行正文 | 本审 |
|------|----------|------|
| L28 | `pg_dump -F c --no-owner --file <artifact> <source-dsn>` | 步骤 1 预转换 dump，token = `<artifact>` |
| L32 | 预转换 dump 是 rollback artifact，不是目标合同 RecoveryPoint | 散文区分，与 A-028/A-029 所见相同 |
| L38 | `pg_restore --exit-on-error --no-owner --dbname <restore-dsn> <artifact>` | **仍绑定步骤 1 的 `<artifact>`** |
| L41 | restore 后 `information_schema` 为 `timestamp with time zone` precision 6 | **预转换 custom dump 现行形状仍是 BIGINT；不重放转换则 L41 必失败** |

A-029 点名的 L34–L39/L41 矛盾**一字未改**。边界文件只在**新**文件里引入 `<rollback-artifact>` / `<recovery-artifact>`，没有改 runbook，也没有把 runbook 标 `superseded`。

同时，conversion contract §0 仍写：该 C3 附件**尚未落盘**；「C3 权威在落盘前仍为 runbook + Port draft」。文件已经落盘，这句话现为假；且它把**仍含矛盾的 runbook** 指定为 C3 权威。A-031 为关 F-I-020.1 写下的回退句，落地后未收回。

**判定**：这构成 F-I-004 关闭要求 2 的**剩余项**（「PG restore 不得再把预转换 `<artifact>` 当目标形状校验输入」未在冻结包内兑现）。另见新 recommended **F-I-027**（卫生：权威未收口 / 未 superseded / 过期「尚未落盘」句）。

SQLite 段 L17–L19 的散文区分维持 A-029 已接受的收窄；本条不倒退。

### C. 包路径与调用点是否可执行

**包路径：可落码，本审接受候选。调用点：SQLite 行号对、结构可行；PG 调用点缺席；A/C 生成时机被挤在同一格；调用点 2 失败后无重试。**

#### 4.1 包路径

| 主张 | 本审 |
|------|------|
| `kernel/store.go` 只放 `RecoveryPointPort` + 中性类型，无驱动类型 | **可落码。** 与现行 `Store`（`:30-38`）并列，不把 Backup 塞进 `Store` 方法面。符合 Root D-007 / Root D-010 |
| `apps/api/internal/backup/{service,provider_sqlite,provider_pg,verify}.go` | **可落码。** 目录目前不存在（本审探测 absent），作为 R2 落点足够具体 |
| 类型见 Port draft | **可落码，但 Port draft 仍 `proposed` 且 L68–L73 Open 未收。** 边界文件引用它却未把它升格或标部分 superseded |

§4.1「用户 `D-010`/`D-007`」**无 Root/child 限定**，违反本文件自己的 §9。child D-007 = C2 精度/D0；child D-010 = `schema_migrations` owner。真正对应「kernel 只暴露 `CreateRecoveryPoint`」的是 **Root D-007 + Root D-010**（child D-006 + child D-009）。内容同向，编号会误导。见 F-I-027。

#### 4.2 调用点 vs 源码

独立核对 `migrate.go` / `postgres.go`：

| 文件主张 | 实测 | 本审 |
|----------|------|------|
| `applyPending` `:81-103` | **是。** `:81` 起；循环 `:92-101`；`:102` `return s.verifyIntegrity()` | 行号正确 |
| `applyMigration` `:108-132` | **是。** `:109` `Begin`；`:128` `Commit` | 行号正确 |
| `verifyIntegrity()` `:102` | **是调用点**，不是函数体（函数体 `:345-363`；`integrity_check` 在 `checkIntegrity` `:365-374`） | 作为「必须在批次校验之后」的锚点可接受 |
| `snapshotBeforePending` `:279-300` | **是。** `:279-301`；`fresh` / `:memory:` / 无数据 → `nil` | 与 §1 一致 |
| 禁止在 `applyMigration` 事务内调用 | **理由成立，且比原文更硬。** SQLite `VACUUM INTO` 不能在打开的事务内跑；`pg_dump` 从另一连接看不到未提交变更、或会拉长事务。`:108-132` 是单事务，禁止成立 | 接受 |
| 调用点 1「转换批次之前生成 A/C；失败 → 批次不开始」 | **A 与 C 被挤在同一格，不成立。** C 是循环内 per-migration（`:93-96`），v73 已提交后 v74 快照失败**不会**让「批次不开始」。A 若是循环前一次 dump，失败才阻止批次 | 剩余项 |
| 调用点 2「`applyPending` 成功返回之后」 | **SQLite 可行**：`migrate()` `:58` 现为 `return s.applyPending(...)`，须改成先 `applyPending` 再调 Port。位置在 `:102` 之后正确 | 结构可行 |
| 调用点 2 失败 → 报错但转换已提交 | **即时行为可接受。后续路径未写。** SQLite `actionNoop`（`:52-53`）只 `verifyIntegrity()`，**不**调 `CreateRecoveryPoint`。Port 失败 → 进程起不来 → 下次启动 `actionNoop` → **B 永不补创** | 剩余项 |
| **PG 调用点** | **未写。** PG 是独立类型：`postgres.migrate` `:72-118`；`applyPendingPG` `:120-127` **无 snapshot、无 `verifyIntegrity`，直接 `return nil`**；`applyMigrationPG` `:154-173` 经 `p.Run` 单事务。A-029 的主缺口在 PG，调用点表却只锚 SQLite `Store.migrate` | **关闭要求 3 未完成** |

「禁止事务内调用」对 PG 同样成立（`applyMigrationPG` 的 `Run` 回调内不得 `pg_dump`）。须在 `applyPendingPG` **返回之后**（`postgres.migrate` `:100` 的 `return p.applyPendingPG(...)` 同样要拆开）写对称调用点。

### D. harness 是否够「可执行」

**作为 R1 设计冻结：够。作为 A-027「命令模板 ≠ 可执行脚本」的实施证据：不够，也不应在 R1 要求落地。**

- 已给出落点：`apps/api/internal/backup/restore_harness_test.go` + `TestSQLiteRestoreToNewDB` / `TestPGRestoreToNewDB` / `TestLegacyArtifactMustFail`。**where 已写，不必再另要文件名。**
- 六项断言覆盖 restore-to-new-db、两侧类型、90 列、样本、catalog/checksum、失败分类；PG 版本兼容显式记录；`t.Skip` 须记原因、不得静默通过。这比 runbook 命令模板具体。
- 按 F-I-003 / F-I-006 的同一把尺子（freeze-candidate 设计可关设计项，实施在 R2），harness **规格**不单独阻断 F-I-004。
- §5 断言 4 把「负值非法」放在 **B** 的正向 restore 上不贴切：转换 fail-closed 后 B 内不应仍有负值。负值应在 `TestLegacyArtifactMustFail`（A/C）或转换预检，不在 B 的 round-trip。不升格为 required。

### E. 关闭要求 5 / `TestLegacyArtifactMustFail`

**方向成立；「把失败当正向证据」尚未可操作。**

成立的部分：同一套 §5 形状校验吃 A/C，断言必须失败；C3 提交须同时具备 (i) B 全过 (ii) A/C 被拒。只有 (i) 不够。这正是 A-029 要的反向门。

不可操作的部分：未规定失败**分类**。若测试只断言 `err != nil`，缺文件、损坏 dump、权限错误也会绿。§5.6 要求「可核对的错误分类」，但 §6 没有把 `TestLegacyArtifactMustFail` 钉到例如 `TimeContractMismatch` / 旧物理类型，而不是 `ArtifactNotFound`。也未规定证据落点（测试名 + 失败分类即够，不必另发明附件格式，但必须写出）。

→ F-I-004 关闭要求 5 的剩余项：反向断言必须钉**形状/合同**错误类。

### F. 现状事实是否准确

四条主主张**成立**。有一处行号偏宽、两处遗漏。

| §1 主张 | 本审 |
|---------|------|
| `apps/` 中 `CreateRecoveryPoint` / `BackupService` / `RecoveryPoint` **0 匹配** | **成立。** 全仓 `apps/**/*.go` 检索 0 命中 |
| `kernel.Store` 无 Backup 面 | **成立。** `kernel/store.go:30-38` 仅 `Dialect`/`Run`/`Ping`/`Close`/`WasFresh`/`MarkSystemDataReady`/`SystemDataReady`。文件写 `:30-47`，`:40-47` 是 `Tx`，不是 Store 方法。主张仍真 |
| `internal/temporal` 不存在 | **成立。** 目录 absent。guardrails 草案仍把它当 codec 候选名，与 Backup 包无关，不构成 §1 虚假 |
| `snapshotBeforePending` 是 per-migration rollback 点 | **成立。** 循环 `:93-96` 对 `version >= 2` 每条 pending 调用。函数注释 `:267-270` 仍写「upgrade batch 的 first pending」（一次），与循环不一致；边界文件采循环语义，比源码注释更准 |
| PG 失败路径 = 事务 rollback + ledger `applied_at`，无 `pg_dump` | **成立。** `postgres.go:154-173` |
| runner 批次后 `foreign_key_check` `:349` / `integrity_check` `:367` | **SQLite 成立。** PG `applyPendingPG` **没有**这两步 |

**遗漏（不影响四条主主张为真）：**

1. PG runner 是独立类型，无 snapshot、无 `verifyIntegrity`。§1 把「现状」写成 SQLite 图景。
2. `internal/backup` 同样 absent（与「无实现」同向，应点名）。
3. Port draft L60 仍允许把 `snapshotBeforePending` 当 RecoveryPoint 源——这是现行**文档**事实，§1 只写了代码事实。

### G. F-I-004 能否闭合

**不能。** 关闭证据不足。全部剩余项如下（不要只记一项）：

对照 A-029 / A-032 原文：「须唯一区分（a）转换前 rollback snapshot/dump（旧合同）与（b）转换后 RecoveryPoint（新合同 + restore 校验）；PG restore 不得再把预转换 `<artifact>` 当目标形状校验输入；写出 `CreateRecoveryPoint` 的包路径与 before/after 调用点。SQLite 散文区分 ≠ C3 冻结。」外加 A-027 的 harness / 旧 dump 不得通过新合同校验。

| # | 剩余项 | 对应关闭要求 | 证据 |
|--:|--------|--------------|------|
| 1 | 冻结包仍把预转换 `<artifact>` 当 PG 目标形状 restore 输入 | 2（A-029 原文） | runbook L28/L38/L41 未改；`status: proposed`；未被 superseded |
| 2 | C3 权威未收口到本边界文件 | 1、2（唯一区分） | conversion contract §0 仍称本文件「尚未落盘」，权威 = runbook + Port draft；Port draft L60 仍允许 C 当 RecoveryPoint 源，L71 仍 Open |
| 3 | PG before/after 调用点未写到 `postgres.go` | 3 | 只有 SQLite `Store.migrate`/`applyPending`/`applyMigration`；`applyPendingPG` `:120-127` 无对称锚点 |
| 4 | 调用点 1 把 A 与 C 挤成「批次前一次」 | 3 | C 是循环内 per-migration；「快照失败 → 批次不开始」只对循环前的 A 成立 |
| 5 | 调用点 2 失败后 `actionNoop` 无重试，B 可能永不补创 | 3 + 失败策略 | SQLite `:52-53`；PG `:94-95` 连 `verifyIntegrity` 都没有 |
| 6 | `TestLegacyArtifactMustFail` 未钉形状/合同错误分类 | 5 | §6 只要求失败；缺文件也会绿 |
| 7 | 转换后 C→B 无机械身份（仅文件名政策） | 1 | 硬规则 1 对后转换 C 不能靠形状挡住 |

**已满足、可维持的子项（不关闭整条）：**

- 新文件内部 (a)/(b) 三类表 + 双 token 设计；
- `CreateRecoveryPoint` 包路径（kernel 接口 + `internal/backup`）本审接受；
- SQLite 调用点行号与「禁止事务内调用」成立；
- harness 规格（含文件名）满足 R1 设计层「不是命令模板」；
- 反向断言的**方向**成立。

freeze-candidate ≠ 实施。无实现不单独阻断本条（与 F-I-003/F-I-006 同一尺子）；阻断的是冻结包不唯一与调用点不完整。

### H. 是否引入新 finding；§7 / §8

**§7 与 `D-019` 自洽。** `D-019` 管 F-5 子女先行顺序、跨 descriptor 两次重建、禁止 runner 级 12 步；不管 artifact。本文件不管表序。F-5 失败在 `applyMigration` 事务内 → 回滚到旧合同，A/C 仍有效、B 不存在。整批 `applyPending`（v73–v87 一次启动）之后才产 B，两次重建都结束后 90 列才齐，与调用点 2「`verifyIntegrity` 之后」同向。

**§8 自认四项真实，但不完整。**

| §8 自认 | 本审 |
|---------|------|
| 1. 无实现 | **真** |
| 2. PG harness 依赖真 PG；skip；版本矩阵归 `I-040-003` | **真** |
| 3. `ArtifactRef` 清理/保留仍待定 | **真**；A-027 Open 有此项，**不是** A-029 关闭要求原文，不单独挡 F-I-004 |
| 4. 包名/类型/`TimeContract` 为候选 | **真**；包路径本审现接受，类型仍挂在 `proposed` Port draft 上 |

§8 **未**自认、但本审认为必须写出的：旧 runbook 仍矛盾且仍是 conversion contract 点名的 C3 权威；Port draft L60/L71；PG 调用点缺席；A/C 时机混淆；`actionNoop` 无重试；反向断言未钉错误类。不是隐瞒实现（文件明确是设计），是关闭要求对位不完整。

**新 finding：**

- **F-I-027 recommended**（C3 冻结包卫生，见下）。不升格 required，以免与 F-I-004 剩余项双计。
- 无新 required。F-I-002 / F-I-005 本轮未触及，维持 open。F-I-025 recommended 维持 open（「20 张时间列表」）。

## Findings

### F-I-002 · 逐列 USING / rebuild / codec

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮未触及）
- A-038 三项设计剩余 `fixed`、整条因可执行测试仍 open——本审不重开、不闭合。

### F-I-003 · 90 列 mapping

- **严重度**：high · **建议**：required · **状态**：**closed**（维持）

### F-I-004 · Backup Port ≠ 可执行备份/回滚方案

- **严重度**：high
- **建议**：required
- **状态**：open（**收窄，不关闭**）
- **影响门禁**：C3、R2/R3；关联 `I-040-003`
- **本轮收窄**：专用边界文件落盘；三类产物 + 双 token 在**新文件内**对位 A-029；(SQLite) 包路径与行号级调用点可核对；harness 规格含文件名；反向断言有方向。
- **仍不闭合**：见 G 全部 7 项。关键句：SQLite 散文 + 新文件双 token **≠** 冻结包已唯一；旧 runbook `<artifact>` 仍在。
- **关闭要求**（更新，减去本审已接受的包路径候选）：把 runbook PG restore 改成 `<recovery-artifact>`（或标 superseded 并声明本边界文件为唯一 C3 权威）；Port draft L60 删掉「`snapshotBeforePending` or」；写出 `postgres.go` 对称 before/after 调用点；分开 A（循环前）与 C（循环内）；写明调用点 2 失败后的 `actionNoop` 重试或 fail-closed；`TestLegacyArtifactMustFail` 钉形状/合同错误类。freeze-candidate ≠ 实施。

### F-I-005 · checksum / append-only 仍不是可执行硬门

- **严重度**：high · **建议**：required · **状态**：open（维持；本轮无新哈希/测试）

### F-I-006 … F-I-026

- 维持既有 closed/recommended。F-I-006 / F-I-020 / F-I-021 / F-I-022 / F-I-023 / F-I-024 / F-I-026 **closed** 维持。F-I-025 recommended 仍 open。

### F-I-027 · C3 落地后冻结包未收口（runbook/Port/conversion §0）

- **严重度**：med
- **建议**：recommended
- **状态**：open（本轮新增）
- **影响门禁**：不单独挡 C3（挡门的是 F-I-004 剩余项 1–2）；卫生不收口会使 F-I-004 复审再次失败
- **描述**：边界文件落地后，冻结包仍三源并行且互相否定：
  1. conversion contract §0：本文件「尚未落盘」，C3 权威 = runbook + Port draft（现为假）。
  2. runbook `status: proposed`，L28/L38/L41 仍用同一 `<artifact>` 做预转换 dump 与目标形状 restore。
  3. Port draft L60 仍允许 `snapshotBeforePending` 当 RecoveryPoint 源；L71 仍把「SQLite conversion snapshot vs Backup Port artifact relationship」列为 Open。
  4. 边界 §4.1 无限定 `D-010`/`D-007`（child 号含义不同；应为 Root D-007 + Root D-010，或 child D-006 + D-009）。
- **关闭要求**：conversion contract §0 改为「已落盘，C3 权威 = 本边界文件」；runbook 标 superseded 或改双 token；Port draft L60/L71 收到与硬规则 1 同一；D-NNN 加 Root/child 限定。

## 必改项汇总

| ID | 门禁 | 闭合前禁止 | 本轮 |
|----|------|------------|------|
| F-I-002 | C2/C3、R2 | 不得冻结 C2；不得实施 schema/codec | **未触及** |
| F-I-004 | C3、R2/R3 | 不得把 Port/runbook 当 C3 冻结 | **收窄**：专用载体+三类+双 token+包路径。**仍缺** G 的 7 项 |
| F-I-005 | C2、R2 | 不得改历史 checksum/DDL | **未触及** |
| F-I-027 | C3 卫生 | 不单独挡冻结 | **新增 recommended** |

F-I-001、F-I-003、F-I-007、F-I-010（planning）、F-I-011、F-I-012、F-I-013、F-I-014、F-I-015、F-I-016、F-I-017、F-I-018、F-I-019、F-I-006、F-I-020、F-I-021、F-I-022、F-I-023、F-I-024、F-I-026 为 closed。F-I-008、F-I-009、F-I-025、**F-I-027** 为 recommended open。

**开放 required = 3**（F-I-002、F-I-004、F-I-005）。在这些合法闭合前：不得冻结 C2、不得冻结 C3、不得修改 migration DDL/公共 formatter、不得放行 R2、不得将 GOAL-002 或 Root R1 标 `done`。

## 与既有意见的异同

| 项 | A-038 independent | A-039 / E-033 自称 | A-040 independent（本条） |
|----|-------------------|--------------------|---------------------------|
| verdict | conditional；open required=3 | 不自证闭合；待本审 | **conditional**；open required=**3** |
| F-I-004 | open；C3 文件不在证据窗 | 首次落盘，待复审 | **仍 open**；收窄不闭合 |
| 包路径 | 未审 | §4.1 已写 | **接受候选** |
| 双 token | 未审 | 声称修掉 L34–L39/L41 | **新文件内成立；旧 runbook 仍矛盾** |
| harness | 未审 | 规格已写 | **R1 设计层足够** |
| 新 finding | 无 | 无 | **F-I-027 recommended** |
| R2 | 禁止 | 禁止 | **禁止** |

无「一要一否」需用户在 finding 之间裁。无需本轮 P-004。F-I-004 不得走 residual 绕过 runbook 矛盾（那是可修的文档不唯一，不是不可验证风险）。

## 信息门禁（P-005）

| ID | 级别 | 最晚阶段 | 当前状态 | 本审 |
|----|------|----------|----------|------|
| I-040-001 | required | C2/R2 | collecting | 本轮未触及；可执行测试仍开，阻断 C2 |
| I-040-002 | required | C1/C2/R2 | collecting | 90 列分母不因本轮扩大 |
| I-040-003 | required | C3/R2/R3 | collecting | F-I-004 仍开放；本轮有设计收窄，信息项未关 |
| I-040-004 | required | R3 | open | F-I-009 仍开放 |
| 共享资料 | — | — | `none` | 无固定引用被当成关闭证据 |

到期且影响本 scope 的 required 信息项：I-040-003 仍开放，阻断 C3/R2。无用户书面 residual。

## 结论 + 建议给编排器/用户的下一步

**conditional。** `e2c0dac2` 把 F-I-004 从「没有边界文件」收窄为「有冻结候选，但冻结包仍三源、PG 调用点未写完、反向断言未钉错误类」。不足以闭合 F-I-004，不足以冻结 C3 或放行 R2。A-039 未声称本条已闭，标记准确。

建议 `/govern`：

1. 响应本 A-040；**不要**把 F-I-004 标 closed；接受包路径候选与「禁止事务内调用」；接受 G 的 7 项剩余为关闭清单。
2. 最小文档收口（仍在 R1 设计面，不改 `apps/`）：runbook PG 双 token 或 superseded；conversion contract §0 收回「尚未落盘」；Port draft L60/L71 与硬规则 1 同一；补 `postgres.go` 调用点；A 与 C 分开；`actionNoop` 重试；反向断言钉错误类；D-NNN 加限定。
3. **不要**冻结 C2/C3，**不要**启动 R2，**不要**改 formatter/DDL。
4. F-I-002 测试口径仍按 A-039 §3 另取 P-004；本条不代裁。
5. 顺手改 rebuild 附件「20 张时间列表」（F-I-025）与 conversion §0（F-I-027）。

## 声明

本意见 `source: independent`，不修改 status / progress / 方案决策 / goal-tree / `apps/`。响应、finding 闭合与是否推进由 `/govern` 处理。
