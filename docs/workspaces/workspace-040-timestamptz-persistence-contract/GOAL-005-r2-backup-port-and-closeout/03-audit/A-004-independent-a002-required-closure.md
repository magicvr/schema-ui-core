---
id: A-004-independent-a002-required-closure
doc: audit-entry
status: recorded
parent: GOAL-005-r2-backup-port-and-closeout
created: 2026-09-20
updated: 2026-09-20
version: 0.1.0
source: independent
verdict: pass
open_required: 0
auditor: grok-build (grok-4.6 · reasoning high)
---

# A-004 · independent · A-002 required 闭合复审（检查点 C）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **日期**：2026-09-20
- **类型** / **scope**：finding-closure · 确认本目标 A-002 的 3 条 required（F-I-001 / F-I-002 / F-I-003）是否已合法闭合；并复核 recommended F-I-004～F-I-007 的处置。**不是**对 R2 的第三次全量关门审计。
- **verdict**：**pass**
- **open required**：**0**（检查点 C 的 independent 门禁在本 scope 内可成立）
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`id` 匹配；`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` 匹配；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。**未读其他工作区。**
- 被审目标：`GOAL-005-r2-backup-port-and-closeout`。
- 对照：A-002（independent · `conditional` · open required = 3）；A-003（self 响应，声称 3 required + 4 recommended 全部 `fixed`）；实现 commit `117d0e05`（其前 `d3dad0a1` 为 A-002 落盘）。
- P-005：`I-041-006` 在 `00-meta` 为 **verified**（用户 2026-09-20 `D-001` accepted）；`I-041-004` 仍为 R3 前 `non-blocking` deferred，不阻断本 scope。共享资料引用 none。
- 本轮可执行证据（`apps/api`，2026-09-20，真实 PG 15.4 + docker 客户端 `pg_dump (PostgreSQL) 15.19`，**未 skip**）：
  - `go test -count=1 ./internal/backup/ ./internal/store/ -run "C3|Restore|Legacy|Recovery|MidBatch"` → **ok**（backup 30.1s / store 32.7s）。
    - `TestPGRestoreToNewDB` PASS（11.91s）
    - `TestPGLegacyArtifactMustFail` PASS（`TimeContractMismatch`，90 列 `bigint`）
    - `TestPGMidBatchArtifactMustFail` PASS（`TimeContractMismatch`）
    - `TestC3RecoveryAnchorsOnPostgresUpgrade` PASS（19.66s）
    - SQLite：`TestSQLiteRestoreToNewDB` / `TestLegacyArtifactMustFail` / `TestPerMigrationSnapshotIsNotARecoveryPointSource` / `TestC3RecoveryAnchorsOnUpgrade` / `TestC3RecoveryPointFailureNeverBlocksStartup` 均 PASS
  - `go test -count=1 ./internal/composition/` → **ok**（63.9s）；抽查 `TestCompositionPostgresStartup` PASS（11.89s，`openStore` 实测 10.45s，与生产入口注入后走 B dump/restore/校验相符）。
- **未改** `status` / `progress` / 方案正文 / goal-tree / 任何 Go 代码。

## 成果（有证据）

1. **F-I-001 原缺陷已消失。** `snapshotBeforePendingPG` 目标名为 `<db>.pre-v%04d-<ts(ms)>-<hex4>.dump`（`postgres.go` `randomSuffix()` + `20060102T150405.000`）。`PgProvider.Create` 仍拒绝覆盖同名文件（`provider_pg.go`），但类 C 不再使用固定 per-version 名，失败重试会生成新路径，与 D-001「可续跑」不再冲突。类 A 为 `<db>.batch-rollback-<ts(ms)>.dump`（带 `Z` 的毫秒时间戳），**不是**固定名。
2. **F-I-002 生产入口已注入。** `composition.openStore` 调用 `recoveryWiring`：SQLite 注入 `backup.NewService`（`RollbackArtifacts=nil`，A/C 仍走原生 `VACUUM INTO`）；PG 注入同一 Service + 已 `RegisterProvider` 的 `backup.PgProvider`（同时作为 `RollbackArtifacts`）。`TestCompositionPostgresStartup` 走 `openStore` 且未 skip。
3. **F-I-003 不再把 `SampleVerified` 写成类型检查别名。** `service.go` PG 分支在置 `SampleVerified=true` 之前调用 `verifyPostgresSamples`；失败以 `SampleMismatch` 返回，不产出 RecoveryPoint。
4. **A-002 的 recommended 主体已落地**：marker 解析（SQLite）、`conversionCompletionVersion`、`TestPGMidBatchArtifactMustFail`、`00-meta` 的 `I-041-006=verified`。

## 对照 A-002 关闭要求

| A-002 finding | 关闭要求（摘） | 本审判定 | 证据 |
|---------------|----------------|----------|------|
| F-I-001 required/high | 类 C 文件名加唯一后缀，或存在则覆盖/换名；加失败后同 ArtifactDir 重开仍能继续 pending 的用例 | **fixed**（唯一后缀已落地；专用碰撞用例仍缺，见残余） | `postgres.go` 文件名；`PgProvider.Create` 仍 fail-closed 于已存在路径；`TestC3RecoveryAnchorsOnPostgresUpgrade` 同 `ArtifactDir` 升级+重开绿 |
| F-I-002 required/med | 三选一：composition 注入 Port / 即使 Port 为 nil 也探测 marker / 用户书面 residual | **fixed**（选了注入） | `composition.go` `recoveryWiring`；SQLite 与 PG 两条方言；`TestCompositionPostgresStartup` 真实 PG `openStore` 10.45s |
| F-I-003 required/med | PG restore 目标至少断言一秒族、一毫秒族、一个 sentinel；或用户书面接受「类型即样本」 | **fixed**（别名已删除；样本函数真实执行。C3「各至少一例」仍未正向锁定，见 F-I-008） | `verify.go` `verifyPostgresSamples`；`service.go` 调用点；`TestPGRestoreToNewDB` PASS |
| F-I-004 recommended | 探测改为成功解析才算 `HasPoint` | **fixed**（SQLite 行为闭合；Detail 声称不完全，见 F-I-009） | `recovery.go` `probeRecoveryState` 走 `readRecoveryMarker` |
| F-I-005 recommended | 改为 catalog 中最后一个 conversion descriptor 的版本 | **fixed** | `conversionCompletionVersion` 认 `vp040_temporal_digital_offer`；`catalog[:72]` / `catalog[:86]` 不含该名 → 返回 0 |
| F-I-006 recommended | 对中途 PG dump 跑与 A 相同的分类断言 | **fixed** | `TestPGMidBatchArtifactMustFail` 本轮 PASS，`TimeContractMismatch` |
| F-I-007 recommended | 编排器同步 `00-meta` / 决策索引 | **partial**（`00-meta` 已 verified；`01-decision.md` 索引仍写 open，见 F-I-010） | `00-meta.md` 信息表 vs `01-decision.md` L17 |

## 逐条闭合判定（required）

### F-I-001 · PG 类 C dump 固定文件名 → **fixed**

- **文件名**：`fmt.Sprintf("%s.pre-v%04d-%s-%s.dump", db, version, time.Now().UTC().Format("20060102T150405.000"), randomSuffix())`。含毫秒时间戳 + 4 hex（`crypto/rand` 2 字节）。与 A-003 声称一致。
- **与 `PgProvider.Create` 拒绝覆盖的关系**：同名仍 `InvalidRequest`（「already exists」）。因为每次拍摄都换名，失败后下次启动再拍 C **不再**撞上 `$db.pre-v%04d.dump`。原可执行反例（预先放置固定名）对当前命名失效。
- **A 类是否唯一**：是。`snapshotBeforeBatchPG` 使用 `<db>.batch-rollback-<ts>.dump`（`20060102T150405.000Z`），不是固定名。A/C 前缀不同（`.batch-rollback-` / `.pre-v`），机械可分。
- **其他固定命名 artifact**：未发现第二处 PG dump 固定 per-version 名。SQLite 类 C 仍是 `<path>.pre-v%04d-<ts>.sqlite`（A-002 已接受）。
- **残余（不升 required）**：A-002 还要求「失败后同 ArtifactDir 重开仍能继续 pending」用例；本轮没有专门的碰撞/中途失败重开测试，快乐路径 `TestC3RecoveryAnchorsOnPostgresUpgrade` 不能替代它。`randomSuffix` 在 `rand.Read` 失败时退回 `"0000"`，与同毫秒重试叠加时理论上可撞名——概率极低，不构成开放 required。

### F-I-002 · composition 未注入 Port → **fixed**

- **SQLite**：`recoveryWiring` 返回 `(service, nil, artifactDir)`。`RecoveryPoints != nil` → `probeIdentity` 会置 `RecoveryPointMissing`；到 head 且无 B 走 `actionVerifyRecoveryPoint`，不再是生产入口上的 `actionNoop`。A/C 不依赖 creator（原生 `VACUUM INTO`），与注释一致。
- **PG**：返回 `(service, provider, artifactDir)`，`PgProvider` 同时满足 Port 的 dump/restore 与 `RollbackArtifactCreator`。A/C/B 在 creator/dir 非空时不再 `return nil` 跳过。
- **artifact 目录**：`recoveryArtifactsDir` **优先** `filepath.Dir(cfg.DBPath)/recovery`。默认 `DBPath` 为 `./data/schema-ui.db`（postgres 也保留该文件路径形值），因此生产 PG **实际**写 DB 同级 `recovery/`，而不是 A-003 正文说的「用户缓存目录按 DSN 派生」。DSN 派生路径仅在 `DBPath` 为空时启用；`UserCacheDir` 失败则回落 `os.TempDir()`。目录由 `PgProvider.Create` 的 `MkdirAll` 创建，可写。
- **`ArtifactDir` 为空**：`recoveryWiring` 显式 `return nil, nil, ""`（`DBPath` 与 `DSN` 都空，或 DSN 净化后为空）。这是显式停用，不是「假装已有 B」。此时 `RecoveryPointMissing` 仍不会被赋值（`recoveryPoints == nil`），`planStartup` 到 head 仍是 `actionNoop`——只覆盖**非法/空配置**，不是默认生产路径。
- **Nil 端口残留**：PG `snapshotBeforeBatchPG` 在未配置时写 `recoveryNote`；`snapshotBeforePendingPG` 在未配置时仍 **silent `return nil`（不写 note）**。生产 composition 已注入，该分支不是默认入口。不重开 required。
- **生产证据**：`TestCompositionPostgresStartup` 配置了 `DBPath`（temp）+ `DBDSN`，`openStore` 10.45s 后 Start 成功——与注入后 fresh 库跳过 A/C、批次后拍 B（dump→restore→校验）的成本相符。

### F-I-003 · PG `SampleVerified` 类型检查别名 → **fixed**

- **不再无条件用类型检查赋值。** 旧语句 `summary.SampleVerified = summary.TypeContractVerified` 已删除。现路径：`verifyPostgresSamples` 成功 → `SampleVerified = true`；失败则返回，不产出 point。`Passed()` 仍要求四项全真。
- **覆盖对照 C3 §5 item 4**（「秒、毫秒、sentinel `0`、可空缺失、fixed-6 词法序各至少一例」）：
  | 项 | PG 实现 | 本审判定 |
  |----|---------|----------|
  | 秒 | 对所有分母列抽样 `LIMIT 50` 非空值，要求 `Nanosecond()%1000 == 0`。fresh 转换库的 `schema_migrations.applied_at`（unit=sec）会被读到 | **有执行**（依赖库内已有行） |
  | 毫秒 | 同一抽样循环，不区分 unit。`TestPGRestoreToNewDB` **不播种** jobs/mail_outbox/mail_config，毫秒族可能 0 行 | **未正向锁定「至少一例」** |
  | sentinel `0` | `mail_config.updated_at` / `telegram_config.updated_at` / `users.locked_until` / `login_failures.locked_until` 要求 **全部 NULL** | **有检查**；空表 COUNT=0 也通过，不能证明 `0→NULL` 转换发生过 |
  | 可空缺失 | 查询失败 / 无行 → `continue`，不失败 | **合理** |
  | fixed-6 词法序 | PG 原生 `timestamptz`，用微秒精度代替 canonical 27 字符 | **可接受的方言同构**（A-002 允许 `to_char` 或 `time.Time`） |
- **「无条件置真」路径**：仅 `verifyPostgresSamples` 返回 nil 之后。空表使函数返回 nil 时，`SampleVerified` 仍为 true（与 SQLite `verifySQLiteSamples` 丢弃 `total` 相同）。这不是类型检查别名。
- **空表行为**：合理（可空缺失不是失败）。`TestPGRestoreToNewDB` 仍不播种秒/毫秒/sentinel 样本（SQLite `convertedStore` 会插 `updated_at=0` 再转换）。原关闭要求的「至少断言一例」由代码在「若有数据则检查」层面满足秒族（账本行），毫秒族与 sentinel 转换例在该测试中仍是空操作。不把这一点升回 required：原缺陷（别名）已消失，且 SQLite 侧 A-002 已接受同等「有则检、无则过」模型。残余见 F-I-008。

## recommended 复核

| finding | A-003 声称 | 本审判定 |
|---------|------------|----------|
| F-I-004 | 解析失败不算已满足，Detail 写入 `recovery marker is unreadable` | **行为 fixed**（`HasPoint` 仅在 `has==true` 时置位）。随后 `state.Detail =` 形状探测结果，**覆盖** unreadable 文案。无损坏 sidecar 测试。PG `readPostgresRecoveryMarker` 解码失败仍 `return err`，`probeIdentity` 会阻断 Open（A-002 已记录的不对称，本轮未改）。见 F-I-009。 |
| F-I-005 | 以最后一个转换描述符为门禁，裁短历史合法，全量 catalog 强制 90/90 | **fixed**。实现按**名字** `vp040_temporal_digital_offer`（冻结台账 v87，确为最后一条 conversion），不是「catalog 中最后一个 `vp040_temporal_*`」的通用扫描。`catalog[:72]` / `catalog[:86]` 不含该名 → `needConverted=0` → 未转换不报错（`TestPGLegacyArtifactMustFail` 能 Open v72；`TestC3RecoveryAnchorsOnPostgresUpgrade` 能 Open v86 再升 v87）。全量 catalog 含该名 → `needConverted=87`；`postgresConvertedShape` 要求 `measured==90 && converted==90 && records==0`，`head >= 87` 且未转换则报错。 |
| F-I-006 | `TestPGMidBatchArtifactMustFail` | **fixed**。本轮真实 PG PASS，分类 `TimeContractMismatch`，并排除 NotFound/Unreadable/ToolFailure。输入为 `catalog[:v73Head]`（v1–v73 混合形状）。 |
| F-I-007 | `00-meta` 同步为 verified | **partial**。`00-meta.md` 信息表已是 **verified**。`01-decision.md` 信息表仍写 `I-041-006` **open（须用户 P-004）**，决策索引仍写「尚无决策」，尽管 `01-decision/D-001-*.md` 已存在且 `status: accepted`。见 F-I-010。 |

## Findings（本轮新开；均为 recommended）

本轮 **不重开** A-002 的 3 条 required。下列为闭合后仍值得编排器处理的残余，**不阻断**检查点 C。

### F-I-008 · PG 样本仍未正向锁定 C3「各至少一例」

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：`verifyPostgresSamples` 不区分秒/毫秒族，不要求 `measured > 0`，sentinel 检查在空表上空过。`TestPGRestoreToNewDB` 仍只证明形状+账本指纹+（空表）sentinel COUNT=0。SQLite harness 有 `convertedStore` 播种，PG 没有对等播种。另：sentinel 列被做成「整列必须全 NULL」——`mail.Switcher.ensureSeeded` 会写入 `updated_at = time.Now()`；若 marker 丢失后补创 B，可能 `SampleMismatch` 而永远补不上（§4.3 不阻断启动，但 B 门禁会一直可见为未满足）。这与 SQLite 样本函数同构，故不升 required。
- **关闭要求**：PG harness 播种一秒族、一毫秒族、一个 `0→NULL` sentinel 并断言读回；或把 sentinel 检查改成「抽样已转换的 0 行」而不是「该列永远不许有值」。

### F-I-009 · 损坏 marker 的 Detail 与 PG 阻断不对称仍在

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：SQLite 损坏 JSON 已不再把 `HasPoint` 置真（F-I-004 核心已闭）。`Detail` 的 unreadable 文案会被形状探测覆盖；无用例锁住「写 `{` 再开 → 补创」。PG 前缀存在但 base64/JSON 坏了仍 error 阻断 Open。
- **关闭要求**：损坏 marker 视为缺失并走有界补创（两侧对称）；补一条 sidecar/`COMMENT` 损坏用例。

### F-I-010 · `01-decision.md` 索引未同步 `I-041-006`

- **严重度**：low
- **建议**：recommended
- **状态**：open
- **描述**：F-I-007 点名的两处之一（`00-meta`）已修；`01-decision.md` 信息表与「尚无决策」索引未修。不否定 `D-001` 裁决本身。
- **关闭要求**：编排器同步决策索引（本审不改）。

## 必改项汇总

**无。** 开放 required = **0**。

## 与既有意见的异同

| 项 | A-002 | A-003（声称） | 本条 |
|----|-------|---------------|------|
| F-I-001 | open required | fixed | **同意 fixed** |
| F-I-002 | open required | fixed | **同意 fixed**（生产 PG artifact 目录实际跟 DBPath，不是 A-003 写的 cache/DSN） |
| F-I-003 | open required | fixed | **同意 fixed**（别名已删除）；C3「各至少一例」未锁 → 新 recommended F-I-008 |
| F-I-004～006 | open recommended | fixed | **同意主体 fixed**；F-I-004 Detail/PG 阻断 → F-I-009 |
| F-I-007 | open recommended | fixed | **00-meta fixed；01-decision 未同步** → F-I-010 |
| 检查点 C | 不得无条件放行 | 待本复审 | **independent 门禁可成立**；status/`done` 仍由 `/govern` 改 |

## 结论 + 建议给编排器/用户的下一步

A-002 的 3 条 required 均已合法闭合（路径 `fixed`，证据可重复核对）。**open required = 0**。检查点 C 的 independent 关门意见在本 scope 内为 **pass**。

建议 `/govern` 下一句：响应 GOAL-005 A-004；将检查点 C 标完成（派生 `progress` 3/3）并在用户确认后把 GOAL-005 标 `done`、Root R2 标完成。F-I-008～F-I-010 为 recommended，不阻断关门；若一并修，优先 F-I-010（元数据）与 F-I-008 的 PG 播种。

## 声明

本意见不修改 status/progress；响应由 /govern 处理。
