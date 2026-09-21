---
id: A-002-independent-r2-closeout
doc: audit-entry
status: recorded
parent: GOAL-005-r2-backup-port-and-closeout
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
source: independent
verdict: conditional
open_required: 3
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-002 · independent · R2 阶段关门审计（M4 检查点 C）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **日期**：2026-09-20
- **类型** / **scope**：close-out · Root `D-016` §4 M1–M4 + C3 §4.2/§4.3/§3.1/§5.1 + 90 列分母 + 未声明回归
- **verdict**：**conditional**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`id` 匹配；`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。**未读其他工作区。**
- 被审目标：`GOAL-005-r2-backup-port-and-closeout`（R2 M4）。M1/M2 = `GOAL-003`（`done · 4/4`），M3 = `GOAL-004`（`done · 3/3`）。
- 对照：本目标 A-001（self · `conditional`）；Root `D-006`/`D-007`/`D-010`/`D-016`/`D-017`；C3 唯一权威 `GOAL-002/attachments/r1-c3-backup-recovery-boundary-v1.0-fc.md`；Port 表面草案 `r1-backup-port-contract-draft-v0.1.md`；分母 `r1-time-column-inventory-v0.3.md`；`D-021` residual 闭合 `GOAL-002/03-audit/A-048`；复审 `GOAL-003/03-audit/A-002`。
- 本轮可执行证据（`apps/api`，2026-09-20）：
  - `go test -count=1 ./internal/backup/ ./internal/store/ -run "C3|Restore|Legacy|Recovery|MissingArtifact|PerMigration"` → **ok**（backup 20.1s / store 29.4s）。
  - 真实 PG 15.4 + 容器客户端 `pg_dump (PostgreSQL) 15.19` 实际执行：`TestPGRestoreToNewDB`、`TestPGLegacyArtifactMustFail`（90 列 `bigint` 被 `TimeContractMismatch` 拒绝）、`TestC3RecoveryAnchorsOnPostgresUpgrade` **PASS**，未 skip。
  - `go test -count=1 ./internal/temporal/ ./internal/store/ -run "TestCompiledMigrationCatalogOwnership|TestFromUnix|TestFormatCanonical|TestTruncateTowardZero|TestPostgresWriteTruncates"` → **ok**。
  - 分母对照：inventory v0.3 与 `internal/temporalcontract/columns.go` **90 行 / 44 表逐行一致**（机械解析，无手抄差）。
- **未改** `status` / `progress` / 方案正文 / goal-tree / 任何 Go 代码。

## 成果（有证据）

1. **Port 表面符合 `D-007`/`D-010`。** `kernel/backup.go` 只导出 `CreateRecoveryPoint`；`RecoveryPointRequest`/`RecoveryPoint`/`VerificationSummary` 无驱动类型；`var _ kernel.RecoveryPointPort = (*Service)(nil)` 编译期锁定。`CreateRecoveryPoint` 在 `Passed()` 之前返回空 `RecoveryPoint`（`service.go:140-142,184-207`）。
2. **C3 §5.1 分类与反向断言在 SQLite 上正确。** `TestLegacyArtifactMustFail` 与 `TestPerMigrationSnapshotIsNotARecoveryPointSource` 断言分类属于 `{TimeContractMismatch, TemporalColumnSetIncomplete}` 且**不属于** `{ArtifactNotFound, ArtifactUnreadable, ToolFailure}`。`TestMissingArtifactIsNotAContractFailure` 把缺文件单独钉在 `ArtifactNotFound`。本轮全部 PASS。
3. **PG 预转换 dump 的反向断言同样是契约分类。** `TestPGLegacyArtifactMustFail` 观测 `TimeContractMismatch`，错误文本列出全部 90 个 `bigint` 列（本轮日志可核对）。
4. **§4.2 调用点在 store runner 内存在，且双方言测试覆盖了快乐路径。** SQLite：批前 A=`VACUUM INTO` `.batch-rollback-*`（恰好 1）、批次内 C=`.pre-v*`、批次后 B 一次、有 marker 再开 noop、删 marker 再开恰好补创一次（`TestC3RecoveryAnchorsOnUpgrade`）。PG：真实 dump A=1 / C=1 / B 经 `pg_dump`→新库 `pg_restore`→校验，DB 注释 marker 写入；清 `COMMENT` 后重开补创（`TestC3RecoveryAnchorsOnPostgresUpgrade`）。`applyMigration` / `applyMigrationPG` 事务内无 artifact I/O。
5. **§4.3 在 Port 已注入时语义成立。** `planStartup` 在 catalog 到 head 且 `RecoveryPointMissing` 时返回 `actionVerifyRecoveryPoint` 而非 `actionNoop`。补创每启动一次。B 失败不阻断启动、不写 marker、`HasPoint=false`（`TestC3RecoveryPointFailureNeverBlocksStartup`）。
6. **§3.1 A/C vs B 的区分是实测驱动，不是文件名。** `CreateRecoveryPoint` 对 restore 目标测物理类型；A/C 喂给同一校验必须失败。命名（`.batch-rollback-` / `.pre-v*`）只是辅助。marker 载体（SQLite sidecar JSON / PG `COMMENT ON DATABASE`）避免 catalog 88，**可接受**为「B 已记录」旁证——但 SQLite 探测只看文件存在，见 F-I-004。
7. **90 列分母完整。** `columns.go` 与 inventory v0.3 逐行 `(table, column, unit)` 一致：90 行、44 张表、`Count = 90`。SQLite 校验覆盖列在场 + 声明类型 `TEXT`；PG 校验覆盖列在场 + `timestamp with time zone` 且 `datetime_precision = 6`。样本（canonical 27 字符 / sentinel NULL）目前只在 SQLite 真正执行。
8. **`D-021` residual 三项已合法闭合。** 复审：`GOAL-003/03-audit/A-002` 判定 ①②③ 可 `fixed`。闭合记录：`GOAL-002/03-audit/A-048`（路径 `fixed`；`D-017` 未被修订）。GOAL-003 A-003 另将当时的 required `F-I-001`（v73 `records` 断言）按 `fixed` 闭合。GOAL-004 A-003 将 M3 required（PG Truncate）按 `fixed` 闭合。本审抽查 codec 单测与 `TestCompiledMigrationCatalogOwnership` 仍绿。
9. **`I-041-006` 已有用户书面裁决。** `GOAL-005/01-decision/D-001`：接受部分升级可续跑 + A/C 回滚，不做批级原子化；C3 §4.2 全部接线。`I-041-004` 仍为 R3 前 `non-blocking` deferred，不阻断 R2。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| M1 共享 codec + 单测 | **满足**（本轮抽查 ok） | `internal/temporal` 11 个 `Test*`；本轮 `-run TestFromUnix\|TestFormatCanonical\|TestTruncateTowardZero` ok |
| M2 15 descriptor + 真实 checksum | **满足**（兄弟目标已关门 + 本轮冻结表测试 ok） | `GOAL-003` `done · 4/4`；`TestCompiledMigrationCatalogOwnership` ok；闭合见 A-048 |
| M3 仓储/谓词 + 双方言回归 + 金额列 | **满足**（兄弟目标已关门；本轮未重跑全仓） | `GOAL-004` `done · 3/3` 的 A-003：PG Truncate `fixed`；金额列/leftover 21 源码仍在 `postgres_test.go` |
| M4 `D-021` residual 三项 + independent 复审 | **满足** | A-048 `fixed`；复审意见 = GOAL-003 A-002 |
| M4 Port + provider + harness | **部分** | 类型表面与 SQLite harness 完整；PG harness 缺 §5 样本（F-I-003）；PG C 重试不安全（F-I-001） |
| M4 C3 §4.2 调用点 | **部分** | store runner 已接线且测试绿；生产 `composition.openStore` 未注入，默认路径 §4.3 退回 `actionNoop`（F-I-002） |
| R2 self + independent 关门审计 | **本条即 independent**；self = A-001 | 开放 required = 3 → **不得**无条件放行 |
| I-041-006 | **已裁决**（D-001 accepted）；meta 未同步 | F-I-007 |
| I-041-004 | **不阻断 R2** | deferred，R3 前；实测组合 = server 15.4 + client 15.19 |

## R2 判据逐项（问题 1）

| 判据 | 本审判定 |
|------|----------|
| M1 | **满足。** codec 单测本轮抽查通过。 |
| M2 | **满足。** 15 个真实 checksum 已记录且经 independent 复审；v73 `records` 断言随后 `fixed`。 |
| M3 | **满足。** GOAL-004 required 已 `fixed` / `user-overruled`；本轮未重跑全仓 `go test ./...`，标 **unverified here**，但不另开 required。 |
| M4 residual 三项 + 复审 | **合法闭合。** A-048 路径 `fixed`，复审意见独立、范围未外扩、`D-017` 未修订。不得把该闭合读成「R2 已放行」。 |
| M4 Port / C3 调用点 / 关门审计 | **未满足无条件放行。** 见 Findings。 |

## C3 问题 2–5（可核对判断，不复述实现）

### §4.2 调用点

- **批前 A / 批次内 C / 批次后 B**：SQLite 与 PG 的 store 路径**都写了**，快乐路径测试**都绿**。
- **「无 B 不得当作已满足」**：仅当 `RecoveryPoints != nil` 时 `probeIdentity` 才置 `RecoveryPointMissing`（`identity.go:438-443,495-500`）。生产 `composition.openStore`（`composition.go:219-226`）**不注入** Port / creator / `ArtifactDir` → 默认启动分类在「到 head 且无 B」时仍是 `actionNoop`。这直接违反 C3 §4.3 第 1/4 项在**默认入口**上的要求（F-I-002）。
- **对称性**：SQLite 的 A/C 不依赖注入（原生 `VACUUM INTO`），生产文件库升级仍会拍 A/C。PG 的 A/C 在未配置 creator/dir 时 **return nil 跳过**（A 写 note，C **连 note 都不写**）。默认生产 PG **零** C3 产物。
- **事务内外部 I/O**：未发现。B 在 `verifyIntegrity`/`verifyIntegrityPG` 之后、事务外。

### §4.3 `actionVerifyRecoveryPoint`

- Port 已注入时：**显式动作 + 每启动最多一次 + B 失败不阻断 + 失败不写 marker** —— 测试锁定，成立。
- **漏判**：
  - 默认入口根本不探测（F-I-002）。
  - SQLite `HasPoint` = `os.Stat` 成功，损坏 JSON 仍当「已有 B」（F-I-004）。
  - 内存库 `probeRecoveryState` 直接 `Converted=false, HasPoint=false`；若误注入 Port，到 head 会走显式动作后因「未转换」**拒绝启动**（测试未覆盖；非生产主路径）。
  - PG 注释被清：已覆盖，会补创。
  - fresh 库：跳过 A/C；批次后仍会尝试 B（若 Port 注入）。合理。

### §3.1 机械身份

- A/C vs B：**实测驱动**，通过。不得仅靠文件名让 `CreateRecoveryPoint` 成功——实现也没有这条捷径。
- marker 载体（sidecar / `COMMENT ON DATABASE`）：在「禁止 catalog 88」约束下**可接受**。
- 脆弱点：SQLite 用**文件是否存在**当 `HasPoint`，不是解析 marker；这是「B 已记录」探测上的文件名/存在性捷径（F-I-004），不是 A/C 校验捷径。

### §5.1 错误分类

- SQLite A 与 C 反向断言：**正确**（本轮 PASS）。
- PG A 反向断言：**正确**（本轮 PASS，分类 `TimeContractMismatch`）。
- PG **没有** C（混合形状）反向断言（F-I-006，recommended）。
- 缺文件：SQLite 有独立测试归 `ArtifactNotFound`；PG `Restore` 源码同样分类（`provider_pg.go:152-154`），无对等测试。

## 90 列分母（问题 6）

机械对照结果：inventory 90 行 = `columns.go` 90 行；去重表 44 = 44；集合差为空。校验覆盖「列在场 + 类型/精度」：**是**（SQLite `TEXT`；PG `timestamptz` 精度 6）。SQLite 样本另检 canonical 27 字符与若干 sentinel 列 NULL。PG 样本被别名为类型检查（F-I-003）。

## 未声明回归 / 残余（问题 7）

| 项 | 判定 |
|----|------|
| `verifyIntegrityPG` 硬编码 87 | **真实**。`postgres.go:209` `head >= 87` 才因未转换报错。与「可续跑」一致，但不是「catalog 中最后一个 conversion descriptor」。见 F-I-005。 |
| PG 逐迁移 dump 成本 | **真实，不升 required。** D-001 已选全部接线；fresh 跳过。失败语义见 F-I-001（文件名碰撞），不是成本本身。 |
| `I-041-004` 跨版本矩阵 | **不阻断 R2。** 本轮记录组合 = 15.4 + 15.19；跨主版本仍 deferred。 |
| store 不得导入 backup | **当前成立。** 生产 `store` 只依赖 `kernel.RecoveryPointPort` + `RollbackArtifactCreator`。`package backup` 测试导入 `store`，`package store` 的 `c3_pg_anchors_test.go` 导入 `backup`——库循环尚未发生；若有人把 backup 写进 `store.go` 会与 backup 测试文件成环。保持窄接口即可。 |

## Findings

### F-I-001 · PG 类 C dump 固定文件名，失败重试会挡住可续跑

- **严重度**：high
- **建议**：required
- **状态**：open
- **影响门禁**：M4 / C3 §4.2 点 2 / D-001「可续跑」
- **描述**：SQLite 类 C 带毫秒时间戳，注释写明 D5「立即重试不得碰撞」。PG `snapshotBeforePendingPG` 写成 `$db.pre-v%04d.dump`（无时间戳）。`PgProvider.Create` 在目标已存在时返回 `InvalidRequest`（「already exists」）。按 D-001，批次中途失败后下次启动从该 pending 继续 → 会先再拍 C → 撞上旧 dump → **Open 失败**，可续跑变成不可启动。
- **证据**：`internal/store/postgres.go:230-235`；`internal/backup/provider_pg.go:125-127`；对比 `internal/store/migrate.go:337`。
- **可执行反例**（真实 PG + docker，沿用本机 `PG_TEST_*`）：

```text
# 1) 先跑绿的锚点测试，确认 ArtifactDir 语义
go test -count=1 ./internal/store/ -run TestC3RecoveryAnchorsOnPostgresUpgrade

# 2) 复现碰撞：在即将升级的库的 ArtifactDir 里预先放同名 C dump
#    打开 v86 库后，升级到 v87 时 snapshotBeforePendingPG 目标为
#    <ArtifactDir>/<db>.pre-v0087.dump
#    预先 os.WriteFile 该路径（任意字节即可）再 Open(catalog[:87])
#    期望：CreateRollbackArtifact 因 already exists 失败，升级不执行。
```

最小片段：对已存在路径调 `PgProvider{WorkDir: dir}.Create(ctx, dsn, filepath.Join(dir, db+".pre-v0087.dump"))` → 非 nil，`KindOf` = `InvalidRequest`。
- **关闭要求**：类 C 文件名加唯一后缀（对齐 SQLite 时间戳），或存在则覆盖/换名；加「失败后同 ArtifactDir 重开仍能继续 pending」用例。**不要**只靠测试用 `t.TempDir()` 来假装没有碰撞。

### F-I-002 · 默认 composition 不注入 Port；§4.3 在生产入口上是空操作

- **严重度**：med
- **建议**：required
- **状态**：open
- **影响门禁**：M4 / C3 §4.3 第 1、4 项 / D-007「供 composition 调用」
- **描述**：store 调用点是「可注入」的，不是「默认生效」的。`composition.openStore` 的 `OpenOptions` 只有连接字段。于是：
  1. 生产 SQLite 仍拍 A/C（原生），但**从不**拍 B，也**从不**把缺 B 当成显式动作；
  2. 生产 PG 的 A/C/B **全部跳过**；
  3. `RecoveryPointMissing` 只在 `recoveryPoints != nil` 时赋值，默认永远是 false → `planStartup` 走到 head 就是 `actionNoop`。这正是 C3 §4.3 禁止的「无 B = 已满足」。
- **证据**：`apps/api/internal/composition/composition.go:219-226`；`identity.go:438-443,495-500`；全仓 `RecoveryPoints:` 赋值只出现在 `c3_anchors_test.go` / `c3_pg_anchors_test.go`。
- **可执行反例**：对已转换的文件库 `store.OpenWithCatalog(path, catalog)`（不传 Port）→ 无 sidecar、`RecoveryPointState().HasPoint == false`，启动成功且无 `verify-recovery-point`。对比同库注入 fake Port 后删 sidecar 再开 → 会补创（已有测试）。
- **关闭要求**（三选一，须留痕）：
  1. composition 注入 Port（SQLite 至少；PG 须同时修 F-I-001 并写明 dump 成本）；或
  2. **即使 Port 为 nil 也探测 marker**，缺 B 不得 `actionNoop`，只记 note、不假装门禁满足；或
  3. 用户书面 `accepted-residual`：R2 只交付可注入能力，生产默认关闭，范围=composition 未接线，复审触发=第一次在生产 Open 注入 Port 时。

### F-I-003 · PG `SampleVerified` 被写成类型检查的别名，§5 样本未跑

- **严重度**：med
- **建议**：required
- **状态**：open
- **影响门禁**：M4 / C3 §5 项 4 / Root `D-010` 四项验证后置条件
- **描述**：`VerificationSummary.Passed()` 要求四项全真。SQLite 在测完形状后调用 `verifySQLiteSamples`（秒/毫秒/sentinel NULL/canonical 27 字符）。PG 分支写 `summary.SampleVerified = summary.TypeContractVerified`（`service.go:180-183`），没有读任何瞬时样本。`TestPGRestoreToNewDB` 也不播种、不断言秒/毫秒/sentinel。C3 §5 写明「两侧同构」。当前 PG 成功返回的 RecoveryPoint 在 `SampleVerified` 上是**口头通过**。
- **证据**：`internal/backup/service.go:180-183` vs `137-139`；`provider_pg_test.go:87-137` 无样本断言。本轮 `TestPGRestoreToNewDB` PASS 只证明形状+账本指纹。
- **关闭要求**：PG restore 目标上至少断言：一秒族列、一毫秒族列、一个 sentinel `0→NULL` 列（可用 `to_char(..., 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')` 或扫 `time.Time` 再 `temporal.MustFormat`）；**或**用户书面接受「PG 无 canonical 文本，类型+精度 6 即样本」并写复审触发。不得继续用类型检查给 `SampleVerified` 赋值而不留痕。

### F-I-004 · SQLite `HasPoint` 只看 sidecar 文件存在，不解析 marker

- **严重度**：med
- **建议**：recommended
- **状态**：open
- **描述**：`probeRecoveryState` 用 `os.Stat(sqliteMarkerPath)` 置 `HasPoint`（`recovery.go:315-317`）。`readRecoveryMarker` 能解析 JSON，但探测路径没用它。损坏/空文件会让缺 B 看起来已满足，§4.3 补创被跳过。PG 侧会解析；前缀在但 base64/JSON 坏了会 **error 阻断启动**（与 §4.3「补创失败不阻断」不对称）。
- **反例**：对已转换库写入 `path+".recovery-point.json"` 内容 `{`，注入 Port 再开 → `HasPoint=true`，不会补创。
- **关闭要求**：探测改为成功解析才算 `HasPoint`；损坏 marker 视为缺失并走有界补创（或记 note 不阻断）。

### F-I-005 · `verifyIntegrityPG` 把转换完成阈值写死为 87

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`head >= 87` 才在未转换时报错。当前 catalog head 恰好是 87，与可续跑一致。下一支非 conversion 迁移会让「部分升级但 head>87」或「转换结束版本上移」失真。应改为 catalog 中最后一个 temporal conversion descriptor 的版本（或传入的 catalog head，并另判形状）。
- **证据**：`postgres.go:193-213`；全仓仅此一处 `>= 87`。

### F-I-006 · PG 缺少类 C 混合形状的 §6 反向断言

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：C3 §6 要求 `TestLegacyArtifactMustFail` 以 A/**C** 为输入。SQLite 有独立的 `TestPerMigrationSnapshotIsNotARecoveryPointSource`（v73 中途）。PG 只有预转换 A。快乐路径的 C dump 存在，但没有「把该 dump 喂给 CreateRecoveryPoint 必须契约失败」的锁。
- **关闭要求**：对 v73-only（或任一中途）PG dump 跑与 A 相同的分类断言。

### F-I-007 · `I-041-006` 已在 D-001 裁决，但 00-meta / 决策索引仍写 open

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`D-001` 用户书面接受部分升级语义，检查点 B 已标 completed。`00-meta.md` 信息表与 `01-decision.md` 索引仍写 `I-041-006` **open（须用户 P-004）**。不否定裁决本身；关门前应由编排器同步元数据（本审不改）。`workspace.md` 仍写「R2/R3 尚未创建」、R1 仍 active，同属过程漂移。

## 必改项汇总

1. **F-I-001（required / high）**：PG 类 C dump 文件名必须可重试，否则 D-001 可续跑在注入 creator 后会变成无法 Open。
2. **F-I-002（required / med）**：默认 composition 不得把「无 B」做成 `actionNoop`；接线、探测与 Port 解耦、或用户书面 residual。
3. **F-I-003（required / med）**：PG `SampleVerified` 必须真正做样本，或用户书面接受「类型即样本」。

开放 required = **3**。recommended 不阻断，但 F-I-004 与 F-I-001 同类（探测捷径），建议同批修。

## 与 A-001（self）的异同

| 项 | self A-001 | 本条 |
|----|------------|------|
| M1–M3 / D-021 residual | ✅ | **同意**（residual 闭合合法；M3 全仓本轮未复跑） |
| Port + 双方言 harness 快乐路径 | ✅ | **同意**（本轮真实 PG 复跑绿） |
| §5.1 反向断言 | ✅ | **同意 SQLite A/C 与 PG A**；指出 PG 无 C 锁、PG 样本是别名 |
| marker 载体 | 交复审 | **接受** sidecar / COMMENT；SQLite Stat 探测不接受为充分 |
| PG 逐迁移成本 | 交复审 | **不升 required**（D-001 已知） |
| 硬编码 87 | 交复审 | recommended F-I-005 |
| backup/store 测试包方向 | 交复审 | **当前无环**；保持窄接口 |
| composition 未注入 | 未提 | **新 required F-I-002** |
| PG C 文件名碰撞 | 未提 | **新 required F-I-001** |
| independent 关门审计 | pending | **本条** |

## 结论 + 建议给编排器/用户的下一步

R2 的 codec / 15 个 descriptor / 仓储改造 / `D-021` residual **可以**视为已满足 M1–M3 与 M4 的 residual 子句。Backup Port 类型表面、SQLite harness、PG 快乐路径与 §5.1 A 类反向断言**可核对**。

不能无条件放行 R2：store 的 C3 接线在默认进程里是关掉的；一旦打开 PG 接线，类 C 重试会把 D-001 的可续跑打成启动失败；PG 成功 RecoveryPoint 的 `SampleVerified` 不是样本证据。

建议 `/govern` 下一句：响应 GOAL-005 A-002；先修 F-I-001，再在 F-I-002/F-I-003 上选「接线 / 探测解耦 / 用户 residual」并留痕。在 3 条 required 合法闭合之前，**不得**把 GOAL-005 标 `done`，**不得**把 Root R2 标完成。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
