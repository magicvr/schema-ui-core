---
id: r3d-root-exit-criteria-matrix-v0.1
doc_type: evidence-attachment
title: R3-D 退出判据证据矩阵 v0.1（Root 六条成功标准逐条落盘）
status: recorded
created: 2026-09-21
updated: 2026-09-21
parent: GOAL-008-r3-exit-matrix-and-root-closeout
version: 0.1.0
---

# R3-D 退出判据证据矩阵 v0.1

**读法**：本矩阵把 Root `GOAL-001` 的**六条成功标准**逐条落成「判据 → 主张 → 证据 → 结论」。每条证据都指向**已关门目标**的可核对产物（台账、附件、测试名或可复跑命令）。**禁止**把本矩阵读作「测试通过即判据满足」：结论列说明的是「证据是否覆盖该判据」，而**关门**由 R3-D 检查点 C 的用户确认（判据 6）决定。

**边界**：容器矩阵结果只作 **CI/reproducibility**（`D-017` §3 约束②、`GOAL-007/D-002` §4）；本矩阵**不**新增任何实证，只汇总已落盘证据。

## 汇总

| # | Root 成功标准（`GOAL-001/00-meta.md`） | 结论 | 主要证据 |
|--:|----------------------------------------|------|----------|
| 1 | PG 物理类型、SQLite 合同平等物理类型、UTC 语义、NULL/零值、编解码与公共面禁止泄漏**已书面冻结** | **满足** | Root `D-002`/`D-003`（承接 `GOAL-002/D-001`）；22 份冻结附件（含 v0.3 列清单与 wire inventory）；`internal/temporal` codec |
| 2 | R1 纳入分母的时间列**双方言迁移完成并通过 checksum**；非时间 `INTEGER` 显式排除 | **满足** | 15 个 v73–v87 descriptor；**15/15** checksum 冻结在 `internal/store/migrate_test.go`；分母 90 列 / 44 表 |
| 3 | 写入绝对时刻与读回一致；VP-020 展示/输入时区不漂移；双方言回归**含至少一条 PG 路径** | **满足** | R2 写入截断与读侧 fail-closed 适配器；R3-A/B wire formatter + 单位族矩阵 + VP-020 round-trip；**真实 PG 路径本轮实测非 skip**（见 §3） |
| 4 | R1 范围内 SQLite 快照 / PG dump 路径**可升级后恢复**，或用户书面 residual 点名 | **满足** | C3 SQLite/PG restore harness；**升级路径** `TestC3RecoveryAnchorsOnPostgresUpgrade`；R3-C 跨版本矩阵 18 个 supported 格（含 15→16/17） |
| 5 | **未**引入 ORM/第三库/Redis/MQ/多实例；**未**把 Admin 维护提示或业务域混入 | **满足** | 依赖清单在本工作区**零变更**；无同类依赖字符串；`kernel/` 仅 `backup.go` 被触碰（见 §5 反向核验） |
| 6 | 退出矩阵与必要独立意见落盘、**开放 required = 0**，并经**用户确认关门** | **部分满足**：矩阵与全部独立意见已落盘、跨目标开放 required = 0；**用户确认仍待**（`I-041-010`） | 本附件；六个已关门目标各自的 `03-audit` 索引；§6 汇总 |

## 1. 判据 1 · 合同已书面冻结

**主张**：方言物理类型、精度、UTC 表示、NULL/零值语义、编解码与公共面禁止泄漏均有冻结决策与逐列证据。

| 证据 | 位置 | 内容 |
|------|------|------|
| 用户方案裁决 | `GOAL-001/01-decision/D-002-r1-contract-freeze-user-decisions.md`（R1 承接：`GOAL-002/01-decision/D-001-r1-contract-freeze.md`） | PG `timestamptz(6)`；SQLite fixed-6 UTC RFC3339 `TEXT`；sentinel `0 → NULL`；不提供 SQLite→PG 搬运器 |
| 公共 wire 输出裁决 | `GOAL-001/01-decision/D-003-r1-public-wire-contract-user-decision.md`（Root `D-003`） | 输出固定 `YYYY-MM-DDTHH:MM:SS.ffffffZ`；`VR-091` 同 |
| 输入兼容裁决 | Root `D-005`（`GOAL-001/01-decision/D-005-r1-wire-input-compat.md`；R1 承接 `GOAL-002/01-decision/D-003-wire-input-compat.md`） | 0/3/6/9 位小数与零等价 offset 接受；**非零 offset 与无时区拒绝** |
| 负瞬间与截断 | Root `D-015`（`GOAL-001/01-decision/D-015-negative-instant-truncation.md`） | 写入先在 Go codec 截断到微秒（向零）；毫秒族 PG 表达式为整数拆分式（原浮点式已弃用） |
| 非 DB 例外范围 | Root `D-009`（`GOAL-001/01-decision/D-009-wire-nondb-exceptions.md`） | 邮件正文 / audit `detail` / 任意 payload 不进固定 6 位合同 |
| 逐列表清单 | `GOAL-002/attachments/r1-time-column-inventory-v0.3.md` | 90 列、44 表、sec/ms 两族、D0 sentinel 与 D0-nullable 分类 |
| wire 清单 | `GOAL-002/attachments/r1-public-wire-inventory-v0.1.md` | Go/Web 输出面、parser、fixture 与例外 |
| 冻结列分母（机读） | `internal/temporalcontract/columns.go` | `Count = 90`、`Columns()`、`Tables()`；被 backup 与 store 共用 |
| codec | `internal/temporal/temporal.go` | `Parse`/`FormatWire`/`Truncate`/`Value`/`NullValue`；`Layout` 与 `CanonicalLen` |

**结论**：判据 1 满足；冻结物是**文本决策 + 机读分母 + codec** 三件套，可逐项对照。

## 2. 判据 2 · 双方言迁移完成并通过 checksum；非时间 INTEGER 显式排除

| 证据 | 位置 | 内容 |
|------|------|------|
| 转换 descriptor | `modules/*/migration/vp040_temporal.go`（**15** 个文件，版本 **73–87**） | 每个 descriptor 含 m0 哨兵普查 + m1–m3 重建 + m4 校验；PG 变体同文件（不进哈希，`D-017`） |
| checksum 台账（冻结） | `internal/store/migrate_test.go` | **v73–v87 的 15/15 checksum 全部被断言**（本轮实测：15 of 15 present） |
| checksum 算法 | `kernel.MigrationChecksum(stmts, transformID)` | `sha256(normalizeSQL(join(stmts,"\n")) + "\n" + transformID)` |
| 边界矩阵（真实迁移） | `internal/w040contracttest/migration_boundaries_test.go` | `TestRealMigrationsConvertBoundaryInstants`、`TestRealMigrationFailsClosedOnNegativeVoucherInstant`、`TestRealMigrationKeepsOrdinaryNegativeInstant`、`TestV73RefusesWhenRetiredRecordsTableIsPresent`、`TestCodecMatchesRealMigration` |
| 非时间 `INTEGER` 排除 | `r1-time-column-inventory-v0.3.md` §排除 + `contract_boundaries_test.go` | ID/duration/step/version/计数/金额/flag 不纳入分母；分母只含绝对时刻 |
| 分母规模 | `internal/temporalcontract/columns.go` | `Count = 90`，44 张表 |

**结论**：判据 2 满足；迁移台账可复跑（`go test -count=1 ./internal/store/ ./internal/w040contracttest/`）。

## 3. 判据 3 · 写入读回一致 + VP-020 不漂移 + 含真实 PG 路径

| 证据 | 位置 | 内容 |
|------|------|------|
| 写入截断（向零） | Root `D-008`/`D-015`；`internal/store/temporal_truncate_test.go` | 写入前 `temporal.Truncate`；PG 写路径亦截断（`GOAL-004/A-002` `F-I-001` fixed） |
| 双方言对拍 | `GOAL-004`（R2 仓储改造）台账 | 读写谓词改造 + 双方言对拍用例 |
| 读侧失败关闭 | `internal/store/scan.go`、`scan_test.go` | 只接受 canonical fixed-6；非规范值 fail closed（读适配器，不向 handler 泄漏驱动类型） |
| wire formatter / parser | `internal/temporal/temporal.go` `FormatWire`；`internal/handler/rfc3339.go` | 输出单一 fixed-6；输入按 `D-005`（拒绝非零 offset/无时区；拒绝逗号小数） |
| 公共 wire 面覆盖 | `GOAL-006/02-execution/E-003` + `internal/handler/w040_r3a_public_wire_fields_test.go` | 三个漏网字段（healthz/readyz、mail outbox、mail config GET+PUT）修复；两条常驻守卫防回归 |
| 单位族矩阵 | `internal/handler/w040_r3b_unit_family_matrix_test.go` | 秒/毫秒/可空/sentinel 四族四端点，族身份由冻结分母核对 |
| VP-020 展示/输入 round-trip | `internal/handler/w040_r3b_timezone_roundtrip_test.go`；`apps/web/src/i18n/utc-roundtrip.test.tsx` | Go：4 时区 × 6 瞬时（含 DST 边界、+05:45、负 epoch）wire 恒定、微秒精确；Web：L1/L2/L3 层解析驱动展示、秒粒度恒等 |
| **真实 PG 路径（本轮实测）** | 命令见下 | **四个目标测试全部 PASS、无一 skip** |

本轮实测（本环境存在常驻 PostgreSQL **15.4**，`PG_TEST_*` 已配置）：

```text
go test -count=1 -v -run "TestPGRestoreToNewDB|TestSQLiteRestoreToNewDB|TestC3RecoveryAnchorsOnPostgresUpgrade|TestCompositionPostgresStartup" \
  ./internal/backup/ ./internal/store/ ./internal/composition/

--- PASS: TestPGRestoreToNewDB (13.51s)            (真实 PG dump → restore → 校验)
--- PASS: TestSQLiteRestoreToNewDB (0.74s)          (SQLite VACUUM INTO → restore → 校验)
--- PASS: TestC3RecoveryAnchorsOnPostgresUpgrade (19.59s)  (升级路径)
--- PASS: TestCompositionPostgresStartup (14.13s)   (生产组装入口在 PG 上启动)
```

> 诚实边界：这些测试在**缺少 `PG_TEST_*` 的环境会 skip（并记录原因）**；本矩阵的「满足」绑定的是**本环境实测非 skip** 的事实，不声称任何环境都必然执行。
>
> **追加限定（2026-09-21，源自 `F-I-101`）**：这些测试使用 `t.TempDir()` 等**绝对**路径夹具，**不能代表「配置里 `db.path` 为相对路径」的真实入口**。用户报告的 `.\dev.cmd start` 失败正是该差异掩盖的启动级缺陷（`recoveryArtifactsDir` 相对路径 → `docker run -v` exit 125）。修复（绝对化 + 两条回归 + 一次性库端到端实测）见 `02-execution/E-003` 与 `03-audit/A-004`；据此判据 3/4 的结论仍成立，但证据面从「测试夹具」扩展到「真实入口实测」。

**结论**：判据 3 满足（含真实 PG 路径，非 skip）。

## 4. 判据 4 · 可升级后恢复（或用户书面 residual 点名）

| 路径 | 证据 | 结果 |
|------|------|------|
| SQLite 快照 → 新库恢复 | `TestSQLiteRestoreToNewDB`（`internal/backup`） | PASS（0.74s） |
| PG dump → 新库恢复 | `TestPGRestoreToNewDB`（`internal/backup`） | PASS（13.51s） |
| **升级后恢复**（C3 锚点） | `TestC3RecoveryAnchorsOnPostgresUpgrade`（`internal/store`） | PASS（19.59s） |
| 负例：legacy 类 A 产物必须失败 | `TestPGLegacyArtifactMustFail`、`TestPGMidBatchArtifactMustFail` | 通过（fail closed） |
| 跨版本升级恢复矩阵 | `GOAL-007/attachments/r3c-pg-cross-version-matrix-v0.1.md` | 15.19/16.15/17.11 上 9 dump + 54 restore 格；**18 个 supported 格形状校验全通过**（含 15→16、15→17、16→17）；两个失败模式逐格记录 |
| 容器证据定位 | `GOAL-007/D-002` §4；`D-017` §3 约束② | **明示**为 CI/reproducibility，不是生产就绪证据 |
| residual 台账 | R1 `D-021` 的 `F-I-005` → `GOAL-002/03-audit/A-048` **`fixed`**；R2 关门 `GOAL-005/A-004` **pass / open required = 0** | 无未闭合 residual；**不得**把 `F-I-005` 重记为 open |

**结论**：判据 4 满足；residual 已按三路径合法闭合，无需用户新裁。

## 5. 判据 5 · 反向核验（未引入 ORM/第三库/Redis/MQ/多实例；未混入 Admin 维护提示/业务域）

本轮**亲自执行**的扫描（命令与结果原样记录）：

| 检查 | 命令 | 结果 |
|------|------|------|
| 依赖清单是否变动 | `git diff --name-only <workspace-040 首个提交>..HEAD \| grep -E 'go\.(mod\|sum)\|package\.json\|pnpm-lock'` | **无变动**（工作区全程零依赖变更） |
| ORM/Redis/MQ/第三数据库 | `grep -E 'gorm\|entgo\|sqlx\|go-redis\|redigo\|sarama\|nats\|amqp\|mongo\|mysql\|lib/pq\|pgx/v4' apps/api/go.mod` | **无命中** |
| Web 同类依赖 | `grep -E 'redis\|amqp\|mqtt\|kafkajs\|mongodb\|mysql\|prisma\|typeorm\|sequelize' apps/web/package.json` | **无命中** |
| 驱动类型是否进入公共契约 | `grep -E '\bpgx\.\|\bpgtype\.\|modernc\.org/sqlite\|sqlite3\.\|\*sql\.Tx'`（非测试源） | `pgx.`/`sqlite` 仅出现在 `cmd/` 工具与适配器的空导入；`*sql.Tx` 仅 `internal/store` 适配器内部与 `internal/testsupport`；`kernel` 只出现于**注释**（`contribution.go:16` 说明公共面用 `kernel.Tx` 而非驱动 `*sql.Tx`） |
| `sql.NullTime` 的定位 | 同上扫描 | 属 `database/sql`（**非**驱动类型），出现在适配器与模块行类型上；`authsession.User.LockedUntil` 是 R2 **已记录**的 NULL 表示（`GOAL-004` 前置 `D-001` §2 #5），handler 只读其 `.Valid/.Time` 派生布尔，**从不**把它送上 wire |
| 业务域/维护提示是否混入 kernel | `git diff --name-only <range> \| grep 'apps/api/kernel/'` → 仅 `kernel/backup.go` | 本工作区只向 kernel 增加 C3 的 `RecoveryPointPort` 表面；未新增业务词汇（`kernel/telegram.go`、`kernel/profile.go` 的业务词均为**既有** VP-029/030 内容，未被本工作区改动） |
| 多实例/分布式协调 | 工作区生产代码变更范围（108 个非测试 `.go` 文件） | 全部落在 `internal/{temporal,temporalcontract,temporalmigrate,store,backup,jobs,auth,mail,channel,composition,handler}` 与 15 个 `modules/*/migration`；**未**引入锁服务、选举、pub/sub 或跨实例协调 |

**结论**：判据 5 满足。

## 6. 判据 6 · 独立意见、开放 required、用户确认

| 已关门目标 | 关门向 independent（裁决） | 开放 required |
|------------|----------------------------|---------------|
| `GOAL-002`（R1） | `A-046`（conditional → 闭合至 0）→ `A-047` self **pass** | **0** |
| `GOAL-003`（R2 M1/M2） | `A-003` **pass** | **0** |
| `GOAL-004`（R2 M3） | `A-002`（1 required：PG 写截断）→ `A-003` self **pass**（`F-I-002` 经用户 `user-overruled`） | **0** |
| `GOAL-005`（R2 M4） | `A-002`（3 required）→ `A-003` → **`A-004` pass** → `A-005` 关门 | **0** |
| `GOAL-006`（R3-A/B） | `A-002` **fail**（3 required）→ `A-003` → **`A-004` pass** → `A-005` 关门 | **0** |
| `GOAL-007`（R3-C） | `A-002` **conditional**（0 required；独立复跑矩阵）→ `A-003` 关门 | **0** |

- **本目标的审计**：`A-001` self 与 grok independent 关门审计属 R3-D 检查点 B（待办）。
- **用户确认关门**：`I-041-010`（required）仍 **open** —— Root `status: done` **必须**等待用户书面确认（判据 6）。本矩阵**不**代替该确认。

**结论**：判据 6 **部分满足** —— 证据矩阵与独立意见已落盘、跨目标开放 required = 0，**尚缺用户确认**。

## 7. 待用户裁决项（P-004）

| ID | 级别 | 问题 | 建议 |
|----|------|------|------|
| `I-041-010` | required | 是否接受本退出矩阵与残留清账后**关闭 Root**（判据 6） | 建议接受：判据 1–5 均有可核对证据，开放 required = 0 |
| `I-041-011` | non-blocking | 是否把 **PG 客户端/服务端支持组合**写入发布说明或退出矩阵（`GOAL-007/A-002` 移交） | 建议记入退出矩阵即可（矩阵已逐格记录）；若要求写入产品发布说明，请在裁决时明确范围与受众 |

## 8. 边界与不得越读

- 本矩阵是**证据汇总**，不新增实证、不改任何冻结决策、不关闭任何 finding。
- 容器矩阵（R3-C）是 CI/reproducibility 证据；生产就绪性不由本工作区声明。
- 「满足」不等于 Root 已关门：判据 6 需用户确认（`I-041-010`）。
