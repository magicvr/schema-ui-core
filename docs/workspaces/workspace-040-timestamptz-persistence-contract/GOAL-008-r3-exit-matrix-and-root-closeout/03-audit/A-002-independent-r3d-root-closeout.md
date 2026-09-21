---
id: A-002-independent-r3d-root-closeout
doc: audit-entry
status: active
parent: GOAL-008-r3-exit-matrix-and-root-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---
# A-002 · independent 关门审计：R3-D 退出矩阵与 Root 关门就绪度

- **source**: `independent`
- **provider / auditor**: grok build · model grok-4.6 · reasoning high
- **日期**: 2026-09-21
- **scope**: `GOAL-008-r3-exit-matrix-and-root-closeout` 检查点 B（退出判据矩阵、判据 5 反向核验、残留清账、self `A-001` 复核）+ **Root** `GOAL-001-timestamptz-persistence-contract` 关门前状态（六条成功标准是否可交给用户确认关门）
- **audit_type**: close-out
- **verdict**: `conditional`
- **开放 required**: **0**
- **是否同意把 Root 交给用户确认关门**: **同意**（判据 1–5 实质成立；判据 6 按设计仍待 `I-041-010`）
- **是否同意把 Root 标 `done`**: **不同意**（本意见不代替用户确认）

本意见不修改 `status` / `progress` / 方案正文 / goal-tree。响应归 `/govern`。仓库内文件本轮**未改**。

---

## 逐判据结论表

| # | 矩阵结论 | 本审判定 | 依据（独立取证，非采信叙述） | 是否同意矩阵表述 |
|--:|----------|----------|------------------------------|------------------|
| 1 | 满足 | **实质满足** | Root `D-002`（PG `timestamptz(6)` / SQLite fixed-6 UTC RFC3339 TEXT / sentinel `0→NULL` / 不提供 SQLite→PG 搬运器）；Root `D-003`/`D-005`/`D-009`/`D-015`；GOAL-002 `D-001` 承接；inventory v0.3（90 列）；`internal/temporalcontract.Count = 90`；`internal/temporal`（`Layout`/`Truncate`/`FormatWire`/`Parse`）；handler `ParseWireTime` 拒绝非零 offset 与逗号小数。**反例未站住**：`users.go` 仍有「3-digit-millisecond」注释，但 `userToMap` 走 `FormatWireTime`（fixed-6），不构成未落码语义。 | **同意结论，不同意引用路径**：矩阵把用户裁决写成不存在的 `GOAL-002/01-decision/D-002-r1-contract-freeze-user-decisions.md`（见 `F-I-001`）。GOAL-002 的 `D-002` 实际是 R2 migration ownership。 |
| 2 | 满足 | **满足** | `migrate_test.go` 冻结表 **v73–v87 共 15 行**哈希（逐条点名，见复跑记录）；`TestCompiledMigrationCatalogOwnership` **PASS**（编译期 descriptor 与冻结哈希一致）；15 个 `modules/*/migration/vp040_temporal.go`；`columns.go` 机械计数 **90 列 / 44 表**；PG leftover 守卫 = 21 个时间列名不得再为 `integer`/`bigint`，金额列 `balance_total`/`amount_delta` 保持 `bigint`；`duration_seconds` / `last_used_step` / `balance_total` 在重建 DDL 中仍为 INTEGER（排除是真的，不是空话）。PG 变体不进哈希，与 child `D-017` 一致（同意 `F-S-003`）。**反例未站住**：v74 对 `mfa_proofs`/`notifications`/`user_mfa` 的 F-5 子女重建仍写 INTEGER，属 D-019 两次重建中间态，终态由后续 descriptor 转换；leftover=0 守卫覆盖终态。 | **同意**。分母外时间列静默留下、分母内未迁移：本轮未构造出成立的反例。排除清单是**类别 + leftover 名集 + 金额列对拍**，不是全库每一条 INTEGER 的穷尽点名表；对已冻结分母足够。 |
| 3 | 满足（本环境） | **满足（本环境）** | 本轮复跑四测试 **全 PASS、无 SKIP**（见 §反例/复跑）。写入截断：`temporal.Truncate` + `bindPostgresArg`；读侧 `scan.go` fail-closed；wire `FormatWire`/`ParseWireTime`；单位族矩阵用 `requireDenominatorColumn` 对冻结分母；VP-020 测试文件存在（Go `w040_r3b_timezone_roundtrip_test.go` + Web `utc-roundtrip.test.tsx`）。`L-1` 成立：缺 `PG_TEST_*` 时 `TestC3RecoveryAnchorsOnPostgresUpgrade` 会 `t.Skip`。 | **同意**，含「满足绑定本环境」限定。 |
| 4 | 满足 | **满足** | SQLite `TestSQLiteRestoreToNewDB` PASS；PG `TestPGRestoreToNewDB` PASS（本轮 `pg_dump (PostgreSQL) 15.19`）；升级路径 `TestC3RecoveryAnchorsOnPostgresUpgrade` PASS（v86→v87 有界升级 + A/C/B 锚点）；负例 `TestPGLegacyArtifactMustFail` / `TestPGMidBatchArtifactMustFail` 本轮 backup 包 PASS。R3-C 附件与 `GOAL-007/D-002` §4 **明示**容器矩阵 = CI/reproducibility，不是生产就绪。`F-I-005`：R1 先 `accepted-residual`（`D-021`），复审触发后由 `GOAL-002/A-048` 按 **`fixed`** 闭合；矩阵未把它重记为 open。 | **同意**。 |
| 5 | 满足 | **满足** | 本轮复跑四类扫描（见 §反例/复跑）：依赖清单 **零变更**；`go.mod`/`package.json` 无 ORM/Redis/MQ/第三库字符串；`pgtype.` **零命中**；`pgx.` 仅 `cmd/e2e-pgset`；`modernc.org/sqlite` 仅适配器/测试/cmd 空导入；`kernel/` 本工作区只改 `backup.go`；生产非测试 `.go` **108** 个，范围与矩阵一致。`sql.NullTime` 与 `*sql.Tx` 见 `F-S-001`/`F-S-002`。未发现本工作区新增的锁服务/选举/pub-sub/外部缓存。composition 仍装配既有 `kernel.Cache`/`EventBus`，但 `go.mod` 无新依赖，不是本工作区引入的第三库。 | **同意**，含 `L-2`（变更范围口径，不是全仓架构复审）。 |
| 6 | 部分满足 | **部分满足** | 六个已关门目标关门向意见核过：开放 required 均为 0（见下）。矩阵与 `A-001` **没有**把 Root 标 `done`，也没有用矩阵代替 `I-041-010`。Root `status: active` / `progress: 2/3` 与自身规则一致（R1/R2 completed，R3 未完成；progress 不推导 `done`）。成功标准六条仍为 `[ ]`——在用户确认前这是正确形态，不是漏勾。 | **同意**「部分满足 / 用户确认仍待」。 |

---

## Findings

### F-I-001 · 矩阵把合同冻结裁决指到不存在的文件
- **级别**: recommended
- **问题**: 退出矩阵 §1「用户方案裁决」路径为 `GOAL-002/01-decision/D-002-r1-contract-freeze-user-decisions.md`。该文件**不存在**。Root `D-002-r1-contract-freeze-user-decisions.md` 才是用户裁决；GOAL-002 承接是 `D-001-r1-contract-freeze.md`。GOAL-002 的 `D-002` 是 R2 migration ownership。`A-001` 复述了同一误引（`GOAL-002 D-002/D-003`）。
- **证据**: 直接打开该路径失败；GOAL-002 `01-decision/` 目录清单；Root `D-002` 正文五条裁决与 codec/inventory 一致。
- **关闭要求**: 改正矩阵（及若重发 self）的引用，使每一行都指向真实产物。不要求重开合同。
- **为何不是 required**: 同表其它指针（inventory v0.3、`temporalcontract`、codec、Root `D-003`/`D-005`/`D-009`/`D-015`）可独立核对；判据 1 实质成立。

### F-I-002 · Root 信息表 I-040-001/002/003 仍为 `collecting`
- **级别**: recommended
- **问题**: Root `00-meta.md` 与 `01-decision.md` 的 P-005 表仍把 `I-040-001`/`002`/`003` 标为 `collecting`（最晚阶段 R1）。GOAL-002 `00-meta` / `A-047` 已将其改为 **`verified`**，且 R1 已关门。用户若只读 Root 信息表，会以为 required 信息项仍开放。
- **证据**: Root `00-meta.md` 信息表；GOAL-002 `00-meta.md` 信息表；`GOAL-002/03-audit/A-047` §3。
- **关闭要求**: 检查点 C 把 Root 两处信息表同步为 `verified`（证据指向 GOAL-002 A-047 / inventory / 冻结附件），或在确认包中显式说明「子目标已 verified、Root 表未刷新」。成功标准 `[ ]` **不要**提前勾掉。
- **为何不是 required**: 信息已被收集并在子目标 verified；这是 Root 台账未刷新，不是真的开放信息门禁。

### F-I-003 · `goal-tree.md` 说明段与树/表矛盾
- **级别**: recommended
- **问题**: 树、状态表、纲领段均写 `GOAL-008 active · 1/3`。同文件「说明」末段仍写「仅剩 R3-D……**待立项**」。
- **证据**: `goal-tree.md` L19–L29、L36–L39 vs L56。
- **关闭要求**: 改写说明段，与树/表一致。属编排器检查点 C 卫生，不必另开 P-004。

### F-I-004 · GOAL-004 `F-I-002` 的 `user-overruled` 书面载体偏弱
- **级别**: recommended
- **问题**: 三路径标签用在一条 **recommended** finding 上，本身不产生开放 required。但用户书面依据只有编排器 `A-003` 的自述（「用户 2026-09-20 书面裁决」），`01-decision/` 无对应 D 条目；且 `A-003` 表格仍残留「**等用户选择**：接受 / 退回 R3 / 恢复伪造瞬时」，与「已 overruled」互相矛盾。
- **证据**: GOAL-004 `03-audit.md` 关门记录；`A-003-response-to-independent-m3.md` L30 与 L53；`01-decision.md` 决策索引只有 `D-001`；`02-execution.md` 无该裁决。
- **关闭要求**: 在 `I-041-010` 确认包中顺带请用户确认该 overruled 记录；或补一条 D；至少删掉「等用户选择」残留句。
- **为何不是 required**: 原 finding 为 recommended；JSON `null` 投影已被后续 R3-A/B formatter 路径覆盖为既成契约。不阻断 Root 确认。

### F-I-005 · GOAL-002 意见索引表未登记 A-048
- **级别**: recommended
- **问题**: `A-048` 文件存在，正文有链接，**索引表止于 A-047**。矩阵把 `A-048` 当作 `F-I-005` 已 `fixed` 的权威闭合点，索引缺口会让「开放 required = 0」的核对路径不完整。
- **证据**: GOAL-002 `03-audit.md` 索引表 L26–L73 vs 正文 L118；`03-audit/A-048-*.md` 存在，`verdict: pass`，闭合路径 `fixed`。
- **关闭要求**: 把 A-048 补进索引表（source/日期/scope/verdict）。不重开 residual。

### F-I-006 · vision/workspace 投影过时（已知差异）
- **级别**: recommended（note）
- **问题**: 同意 `F-S-005`：`docs/vision/roadmap.md` 仍投影 Root `active · 0/3`、`I-040-001 collecting`。另外 `workspace.md` 纲领段仍写 R1 `active · 0/4`、R2/R3 pending——比 roadmap 更旧。
- **证据**: `docs/vision/roadmap.md` L60/L260/L414/L470；`workspace.md` L23–L52。
- **关闭要求**: 属 `/vision` 层与 workspace 上下文刷新；**不是**判据矛盾，也**不是**第二套目标状态。检查点 C 可同步或在确认包中声明「`docs/vision` / `workspace.md` 叙述不是关门权威」。不必为此否决 `I-041-010`。

**必改项汇总**: 无。开放 required = **0**。

---

## 逐条结论表（待复审 10 条 + A-001 提交项）

| 项 | 本审判定 |
|----|----------|
| 待复审 1 · 矩阵是否指向真实产物 | **大体是**；判据 1 有一条死链（`F-I-001`）。其余指向的台账/测试/附件可打开。未用 `progress` 代替判据。 |
| 待复审 2 · 判据 1 与冻结物一致 | **一致**（Root D-002 + codec + inventory）。未发现「已冻结未落码」的合同语义。误引不推翻实质。 |
| 待复审 3 · 判据 2 checksum / 排除 | **成立**。15/15 哈希文本已数过并与 `TestCompiledMigrationCatalogOwnership` 对拍；90/44 机械计数成立；排除有类别清单 + leftover 名集 + 金额列断言。 |
| 待复审 4 · 判据 3 含真实 PG | **成立（本环境）**。四测试无 SKIP。 |
| 待复审 5 · 判据 4 双路径 + residual | **成立**。SQLite/PG restore + 升级锚点；容器矩阵定位正确；`F-I-005` 为 `fixed`（A-048），未重记 open。 |
| 待复审 6 · 判据 5 反向核验 | **成立**。本轮亲自复跑，结论与矩阵一致。 |
| 待复审 7 · 残留清账 | **成立**。`F-I-005` 未重开；`I-041-004`/`008`/`I-040-004` 在对应目标为 verified。Root 信息表未同步是 `F-I-002`，不是把已关闭项留在 open 清单里。 |
| 待复审 8 · 跨目标开放 required | **0**。见下节逐目标。GOAL-004 `F-I-002` 的 overruled 载体偏弱（`F-I-004`），不构成开放 required。 |
| 待复审 9 · 判据 6 门禁 | **被尊重**。Root 仍 `active · 2/3`；`I-041-010` open；矩阵写「部分满足」；self 未放行。 |
| 待复审 10 · 台账/投影 | 树/表/`GOAL-008` `1/3`/Root `2/3` **自洽**。说明段（`F-I-003`）、Root 信息表（`F-I-002`）、vision/workspace 投影（`F-I-006`）是已知/应修差异，不是「矩阵与权威状态互否」。 |
| `F-S-001` · `sql.NullTime` 是否红线 | **不构成 Root 红线违反**。红线禁止的是 `pgtype` / 驱动时间类型 / `*sql.Tx` 进入公共契约。`sql.NullTime` 是 `database/sql` 可空域表示，GOAL-004 `D-001` §1/§2#5 已冻结。handler（`users.go:108`、`users_state.go:111`）只读 `.Valid/.Time` 派生布尔 `locked`，wire 不出现 `LockedUntil`。 |
| `F-S-002` · testsupport `*sql.Tx` | **不越界**。`Store.WithTx` 源码注释写明不是方言中立模块契约（走 `kernel.Tx`）；`testsupport` 用它清 `must_change_password`，属测试支撑。不是 handler/模块公共面泄漏。遗留 seam 已在 GOAL-004 `F-I-003` 记录，R2 不扩大改造。 |
| `F-S-003` · PG 变体不进哈希 | **与已冻结 `D-017` 一致**。矩阵未把 PG 变体说成「已哈希」。checksum 算法注释与冻结表均为 SQLite canonical 切片。 |
| `F-S-004` · 未越权代替用户确认 | **同意**。矩阵/self/本审均未把 Root 标 `done`。 |
| `F-S-005` · vision 投影过时 | **同意：已知差异，不是矛盾**。`docs/vision` 不是 goal-tree / 关门权威。见 `F-I-006`。 |
| `L-1` · 真实路径绑定本环境 | **同意，且本轮独立复现**。 |
| `L-2` · 多实例判断基于变更范围 | **同意**。108 个生产文件；不是全仓架构复审。 |

### 跨目标开放 required（逐个核对）

| 目标 | 关门向意见 | 开放 required | 三路径抽查 |
|------|------------|---------------|------------|
| GOAL-002 | A-046 independent `conditional`/0 → A-047 self `pass`；A-048 将 `F-I-005` 由 residual 复审为 **`fixed`** | **0** | `F-I-002` `fixed`；`F-I-005` 先 `accepted-residual`（`D-021` 用户书面范围+触发）再 `fixed`（触发已发生且 independent 复审）——路径用对 |
| GOAL-003 | A-002 → A-003 `pass` | **0** | A-002 的 required `F-I-001` `fixed` |
| GOAL-004 | A-002（required=`F-I-001` PG 截断）→ A-003 `pass` | **0** | `F-I-001` `fixed`（有双方言对拍）；`F-I-002` recommended 标 `user-overruled`（载体偏弱，见 `F-I-004`） |
| GOAL-005 | A-002 required=3 → A-003 → **A-004 `pass`/0** → A-005 关门 | **0** | 三条 required 均 `fixed`，A-004 复审确认 |
| GOAL-006 | A-002 `fail`/3 required → A-003 → **A-004 `pass`/0** → A-005 关门 | **0** | 三条 required `fixed` |
| GOAL-007 | A-002 `conditional`/**0 required**（独立复跑矩阵）→ A-003 关门 | **0** | 仅 recommended，A-003 称已 `fixed` |
| GOAL-008 | A-001 self `conditional`/0 required | **0**（本审亦 0） | — |
| Root | 无独立 `03-audit` 条目（关门意见走 GOAL-008，正确） | 无未闭合 required finding；`I-041-010` 是**待用户确认的信息门禁**，不是 finding | — |

---

## 反例/复跑记录

### 基线
- HEAD: `26e04231`（`docs(w040-r3d): land the Root exit-criteria evidence matrix (checkpoint A)`）
- 审计**开始时**工作树**不干净**（与任务说明「应 clean」不符）：已有编排器写入、未提交的 `GOAL-008/03-audit.md` 修改 + 未跟踪的 `03-audit/A-001-self-r3d-exit-matrix.md`。本审**只读**这些文件，**未再写入**任何仓库路径。
- 本审结束：`git status` 仍仅上述两处（编排器自审），无本审临时文件。

### 判据 2 · 15 个冻结 checksum（原文）
```
v73 vp040_temporal_core_persistence     4dd07092330cb3b344143f49ec57edd1f89e92cf2786446c0635c4a320eb8a9f
v74 vp040_temporal_authsession          ae1aefe89925f1759e6e0154f64ddf73b9eb04472c1111e1400b473d7c477f55
v75 vp040_temporal_operationlog         5b038c6cf66f0e01f2446721a246917fc4c0b459dc877242f374aea971322b3a
v76 vp040_temporal_jobs                 6b3649579cc6aedc713fab8c7f9f6dafbec548317f7395082d9e5ddb730aca56
v77 vp040_temporal_dictionary           0266f2937603f3cbed07b33ee460964ad1fe083e2c0f426c4946b0f922104edb
v78 vp040_temporal_data_permission      b1fa8aa94597a8f48a061efe43d7285300016070d391e38b8df95f49dbcdf7ed
v79 vp040_temporal_captcha              c6a661ebfb90158f3a712ad149084f3e84f996ca6773f04cb8743b0e2e41142f
v80 vp040_temporal_mfa                  2088626f0bdf5c2b9aba9a5faf297762cba553c0cf750c1255e1c8633ed6c2fc
v81 vp040_temporal_notifications        a7565dc641f3c3291ff25cbefa52199ca06a8f94f7a978efac5b907615442b37
v82 vp040_temporal_recycle              0132f6a873dd427b42c3a668bc88badc9b50a6c6729601e7cd1c5d5e4a1568fe
v83 vp040_temporal_scheduled_tasks      b6c4f115e54d163a2ec9a1dfab0c74ca5d10c09c9cf1c41c14b02ad5d1186e91
v84 vp040_temporal_settings             bb3041a3d3fbeb5b3d706209f53cc578dc0e5d15016502919aac040b6bec2112
v85 vp040_temporal_wallet               5c004e0c643f46de5e72035022642f9674ca92306784782014f4e285c07efeb8
v86 vp040_temporal_telegram             07a9a0ba61110b94acd4f4f7e87144fd297e8a99e44a273059080e07f1d172f1
v87 vp040_temporal_digital_offer        753b22027027066bd54b8909974b2f867ac5340bd513ce861bf1e2553f1fc7e4
```
`TestCompiledMigrationCatalogOwnership` → `ok store 0.329s`。

### 判据 3/4 · 本轮复跑（`apps/api`，`-count=1 -v`）
```
--- PASS: TestPGRestoreToNewDB (13.56s)             无 SKIP；pg client 15.19；sample coverage map[ms:2 sec:90]
--- PASS: TestSQLiteRestoreToNewDB (0.94s)          无 SKIP
--- PASS: TestC3RecoveryAnchorsOnPostgresUpgrade (20.09s)  无 SKIP（docker + PG_TEST_* 均可用）
--- PASS: TestCompositionPostgresStartup (14.60s)   无 SKIP
EXIT=0
```
另：`TestPGLegacyArtifactMustFail` / `TestPGMidBatchArtifactMustFail` → `ok backup 14.058s`。

### 判据 5 · 本轮扫描
| 检查 | 结果 |
|------|------|
| `git diff --name-only b0a6789b HEAD` 对 `go.mod`/`go.sum`/`package.json`/lockfile | **空** |
| `go.mod`：`gorm\|entgo\|sqlx\|go-redis\|redigo\|sarama\|nats\|amqp\|mongo\|mysql\|lib/pq\|pgx/v4` | **无命中**（既有 `jackc/pgx/v5`、`modernc.org/sqlite` 为双方言驱动，非本工作区新增） |
| `apps/web/package.json`：`redis\|amqp\|mqtt\|kafkajs\|mongodb\|mysql\|prisma\|typeorm\|sequelize` | **无命中** |
| `pgtype.` | **零命中** |
| `pgx.`（非测试） | 仅 `cmd/e2e-pgset/main.go` |
| `*sql.Tx` | store 适配器内部 + `testsupport` + 测试；`kernel/contribution.go:16` 为注释（公共面用 `kernel.Tx`） |
| `kernel/` diff | **仅** `apps/api/kernel/backup.go` |
| 生产非测试 `.go` | **108** |

### 主动构造、未站住的反例
1. `authsession.User.LockedUntil sql.NullTime` 被 handler 读取 → **不是**驱动类型泄漏（见 `F-S-001`）。
2. `testsupport` 使用 `*sql.Tx` → **不是**公共契约越界（见 `F-S-002`）。
3. v74 F-5 子女表仍 `INTEGER` 时间列 → **中间态**，不是分母内未迁移。
4. `users.go`「3-digit-millisecond」注释 → **注释过时**，序列化走 `FormatWireTime`。
5. composition 装配 `Cache`/`EventBus` → **既有** kernel 端口，本工作区未引入 Redis/MQ（依赖零变更）。

---

## 明确结论

- **开放 required 数量**: **0**（跨 GOAL-002～008 与 Root finding 台账）。
- **判据 1–5**: 实质满足；矩阵对判据 1 的一处引用是死链（recommended）。
- **判据 6**: 部分满足——独立意见将由本条落盘；用户确认 **未发生**。
- **是否同意把 Root 交给用户确认关门**: **同意**。
- **是否同意现在把 Root 标 `done`**: **不同意**。
- **最小补齐（不阻断提交用户，但应进入确认包）**: 改矩阵死链（`F-I-001`）；同步 Root 信息表（`F-I-002`）；改 goal-tree 说明段（`F-I-003`）；补 GOAL-002 索引 A-048（`F-I-005`）。这些是编排器检查点 C 卫生，**不是**新的 P-004，除非用户否决同步。

### 交给用户的裁决点

| ID | 级别 | 问题 | 本审建议 |
|----|------|------|----------|
| `I-041-010` | required | 是否接受本退出矩阵（含本独立意见与残留清账）后**关闭 Root** | **建议接受**。判据 1–5 有可复跑证据；开放 required = 0；容器矩阵不得读成生产就绪；真实 PG 路径绑定本环境（常驻 15.4 + `PG_TEST_*`）。 |
| `I-041-011` | non-blocking | 产品级「PG 客户端/服务端支持组合」是否写入发布说明或退出矩阵 | **同意 self**：矩阵/R3-C 附件已逐格记录即可；若要进产品发布说明，请写明范围与受众。17 client 不能恢复到 15/16 是已测事实，不是本工作区的支持承诺。 |
| （可选，本审新增）GOAL-004 `F-I-002` overruled 记录 | non-blocking | 是否确认编排器对 `*time.Time` + JSON `null` 的 `user-overruled` 记录 | 建议顺带确认，或授权编排器补 D / 删「等用户选择」残留句。不确认也不阻断 Root 关门。 |

**不要**把下列事项当成用户必须先裁的门禁：vision/roadmap 刷新、workspace.md 纲领段刷新、成功标准方框勾选（勾选应发生在用户确认**之后**）。

---

## 与 self（A-001）的异同

- **同意** A-001 对判据 1–5 的实质结论、`L-1`/`L-2`、`F-S-001`～`F-S-005`、判据 6 未越读、开放 required = 0。
- **不同意** A-001 把 GOAL-002 `D-002`/`D-003` 当作合同冻结证据路径（与矩阵同一误引）。
- **新增** 台账卫生 `F-I-001`～`F-I-006`（全部 recommended）。无新 required，**不**把 A-001 的 recommended 升格为必改。

## 建议给编排器 / 用户的下一步

用 **`/govern`** 响应本意见：把本条落盘为 `GOAL-008/03-audit/A-002-*.md` 并更新索引；闭合或登记 `F-I-001`～`F-I-005` 的处置（可与检查点 C 同批）；**然后把 `I-041-010` + `I-041-011` 交给用户书面确认**。在用户确认前不得把 Root 标 `done`。
```
