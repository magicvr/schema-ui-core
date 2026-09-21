---
id: r3c-pg-cross-version-matrix-v0.1
doc_type: evidence-attachment
title: R3-C 实测：PG 15/16/17 跨版本 pg_dump/pg_restore 组合矩阵（真实 VP-040 迁移链）
status: recorded
created: 2026-09-21
updated: 2026-09-21
parent: GOAL-007-r3-pg-cross-version-restore-matrix
version: 0.1.0
---

# R3-C 实测：跨版本 `pg_dump`/`pg_restore` 组合矩阵

本附件是 `D-001` 冻结口径下由驱动**原样输出**的实测记录。与前置探测（`r3c-pg-tool-compatibility-probe-v0.1.md`）不同，本矩阵使用**本仓真实的 VP-040 迁移链**（`compiledmodules.PersistenceCatalog()` 全量应用到每个 server 容器的新库，至 head **87**），因此同一次运行也回答了「迁移链在 PG 16/17 上是否成立」。

## 复现入口

```text
cd apps/api
VP040_PG_MATRIX=1 VP040_PG_MATRIX_OUT=<out.md> \
  go test -count=1 -run TestPGCrossVersionRestoreMatrix -v -timeout 45m ./internal/backup/
```

驱动：`apps/api/internal/backup/pg_cross_version_matrix_test.go`（默认**跳过**；需 docker 与外部容器网络）。机制：`docker network` + `--network-alias`，server 容器同时发布一个随机宿主端口供宿主侧迁移/读取，client 容器走网内别名；归档经宿主临时目录 bind mount 交换；就绪门 `pg_isready`。资源在 `t.Cleanup` 中删除；常驻 15.4 实例**未**参与本矩阵（只作 `PG_TEST_*` 的 C3 库内测试）。

## 本次运行的形态

- 源库：3 个（server 15.19 / 16.15 / 17.11 各一个），均由真实迁移链建到 v87，并播种 3 行代表性数据：`users.created_at = 2026-09-20T12:57:15.900000Z`（微秒尾零）、`users.locked_until = 2026-01-02T03:04:05.123456Z`（D0 sentinel 列的非空方向）、`mx-u1.locked_until = NULL`（sentinel 的 NULL 方向）、`jobs.created_at = 2026-09-20T12:57:15.914000Z`（毫秒族）。
- dump 格：9（3 server × 3 client），其中 6 格产出归档。
- restore 格：**54**（6 归档 × 3 client × 3 目标 server）。
- 判定：`supported` 必须同时满足退出码 0 **且** 形状校验通过（`D-001` §4）；形状校验共执行 **18** 次，全部通过。
- `unexpected-failure` / `setup-failure`：**0**（驱动 fail closed，出现即整体失败）。

## 形状校验（每个 supported 格，D-001 §4）

| # | 校验 | 结果 |
|--:|------|------|
| 1 | `schema_migrations` head = 87，行数与源库一致 | 全部通过 |
| 2 | 抽样列 `data_type = timestamp with time zone` 且 `datetime_precision = 6`（`users.created_at`、`jobs.created_at`、`mail_config.updated_at`） | 全部通过 |
| 3 | 抽样瞬时逐字节等于源库规范形（`.900000` / `.123456` / `.914000`，微秒尾零不得被降为 `.9`） | 全部通过 |
| 4 | D0 sentinel 列 `users.locked_until` 的 NULL 保持 NULL（legacy `0` 不复活） | 全部通过 |

→ 由此得到一条**正向结论**：VP-040 迁移链在 **PG 16.15 / 17.11** 上可完整应用至 v87（三个源库分别建在三个 server 上，无失败），且恢复后的库保持转换后合同形状。

## 判定口径的修正（与 `D-001` §3 预期规则的差异）

`D-001` §3 曾把探测得到的预期规则写成 `dumper ≤ client ≤ target_server`。**本矩阵给出反例**：`src15_by_16.dump`（由 16 client 生成）用 **16 client 恢复到 15 server** 是 `supported` 且形状校验通过（同表 `src16_by_16` → 16 client → 15 server 亦同）。即「client ≤ target」并非一般规则。

实测支持的准确规则（**修正后**，详见 `01-decision/D-002`）：

1. **归档格式门**：恢复 client 的 major ≥ 生成归档的 client major（否则 `unsupported version (1.1x) in file header`）。
2. **server-GUC 门**：**PG 17 的 client** 在恢复会话中设置 PG 17 引入的 `transaction_timeout`，因此 17 client 只能恢复到 **17 server**；client ≤ 16 未观测到该门（16 client 可恢复到 15 server）。
3. 目标 server 的 major **不**受「必须 ≥ client」约束（反例见上）。

`D-001` 的冻结口径（判定分类、组合枚举、形状校验）**未被推翻**——被修正的只是 §3 中对探测规则的**外推陈述**；矩阵本体按冻结口径正常分类，未出现 `unexpected-failure`。

## 驱动原样输出

### Measured versions

| server | image | server version | client tools |
|---|---|---|---|
| 15 | `postgres:15-alpine` | 15.19 | 15.19 / 15.19 |
| 16 | `postgres:16` | 16.15 (Debian 16.15-1.pgdg13+2) | 16.15 (Debian 16.15-1.pgdg13+2) / 16.15 (Debian 16.15-1.pgdg13+2) |
| 17 | `postgres:17-alpine` | 17.11 | 17.11 / 17.11 |

### pg_dump: server x client (9 cells)

| server | client | exit | class | diagnostic |
|---|---|---|---|---|
| 15 | 15 | 0 | supported |  |
| 15 | 16 | 0 | supported |  |
| 15 | 17 | 0 | supported |  |
| 16 | 15 | 1 | unsupported-toolgate | pg_dump: error: aborting because of server version mismatch |
| 16 | 16 | 0 | supported |  |
| 16 | 17 | 0 | supported |  |
| 17 | 15 | 1 | unsupported-toolgate | pg_dump: error: aborting because of server version mismatch |
| 17 | 16 | 1 | unsupported-toolgate | pg_dump: error: aborting because of server version mismatch |
| 17 | 17 | 0 | supported |  |

### Archive format version (PGDMP header, decimal)

| archive | dumper | header |
|---|---|---|
| src15_by_15.dump | 15 | 80 71 68 77 80 1 14 (1.14) |
| src15_by_16.dump | 15 | 80 71 68 77 80 1 15 (1.15) |
| src15_by_17.dump | 15 | 80 71 68 77 80 1 16 (1.16) |
| src16_by_16.dump | 16 | 80 71 68 77 80 1 15 (1.15) |
| src16_by_17.dump | 16 | 80 71 68 77 80 1 16 (1.16) |
| src17_by_17.dump | 17 | 80 71 68 77 80 1 16 (1.16) |

### pg_restore: archive x client x target server (54 cells)

| archive | dumper | client | target | exit | class | shape | diagnostic |
|---|---|---|---|---|---|---|---|
| src15_by_15.dump | 15 | 15 | 15 | 0 | supported | ok |  |
| src15_by_15.dump | 15 | 15 | 16 | 0 | supported | ok |  |
| src15_by_15.dump | 15 | 15 | 17 | 0 | supported | ok |  |
| src15_by_15.dump | 15 | 16 | 15 | 0 | supported | ok |  |
| src15_by_15.dump | 15 | 16 | 16 | 0 | supported | ok |  |
| src15_by_15.dump | 15 | 16 | 17 | 0 | supported | ok |  |
| src15_by_15.dump | 15 | 17 | 15 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src15_by_15.dump | 15 | 17 | 16 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src15_by_15.dump | 15 | 17 | 17 | 0 | supported | ok |  |
| src15_by_16.dump | 15 | 15 | 15 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.15) in file header |
| src15_by_16.dump | 15 | 15 | 16 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.15) in file header |
| src15_by_16.dump | 15 | 15 | 17 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.15) in file header |
| src15_by_16.dump | 15 | 16 | 15 | 0 | supported | ok |  |
| src15_by_16.dump | 15 | 16 | 16 | 0 | supported | ok |  |
| src15_by_16.dump | 15 | 16 | 17 | 0 | supported | ok |  |
| src15_by_16.dump | 15 | 17 | 15 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src15_by_16.dump | 15 | 17 | 16 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src15_by_16.dump | 15 | 17 | 17 | 0 | supported | ok |  |
| src15_by_17.dump | 15 | 15 | 15 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src15_by_17.dump | 15 | 15 | 16 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src15_by_17.dump | 15 | 15 | 17 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src15_by_17.dump | 15 | 16 | 15 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src15_by_17.dump | 15 | 16 | 16 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src15_by_17.dump | 15 | 16 | 17 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src15_by_17.dump | 15 | 17 | 15 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src15_by_17.dump | 15 | 17 | 16 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src15_by_17.dump | 15 | 17 | 17 | 0 | supported | ok |  |
| src16_by_16.dump | 16 | 15 | 15 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.15) in file header |
| src16_by_16.dump | 16 | 15 | 16 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.15) in file header |
| src16_by_16.dump | 16 | 15 | 17 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.15) in file header |
| src16_by_16.dump | 16 | 16 | 15 | 0 | supported | ok |  |
| src16_by_16.dump | 16 | 16 | 16 | 0 | supported | ok |  |
| src16_by_16.dump | 16 | 16 | 17 | 0 | supported | ok |  |
| src16_by_16.dump | 16 | 17 | 15 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src16_by_16.dump | 16 | 17 | 16 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src16_by_16.dump | 16 | 17 | 17 | 0 | supported | ok |  |
| src16_by_17.dump | 16 | 15 | 15 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src16_by_17.dump | 16 | 15 | 16 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src16_by_17.dump | 16 | 15 | 17 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src16_by_17.dump | 16 | 16 | 15 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src16_by_17.dump | 16 | 16 | 16 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src16_by_17.dump | 16 | 16 | 17 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src16_by_17.dump | 16 | 17 | 15 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src16_by_17.dump | 16 | 17 | 16 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src16_by_17.dump | 16 | 17 | 17 | 0 | supported | ok |  |
| src17_by_17.dump | 17 | 15 | 15 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src17_by_17.dump | 17 | 15 | 16 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src17_by_17.dump | 17 | 15 | 17 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src17_by_17.dump | 17 | 16 | 15 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src17_by_17.dump | 17 | 16 | 16 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src17_by_17.dump | 17 | 16 | 17 | 1 | unsupported-toolgate | - | pg_restore: error: unsupported version (1.16) in file header |
| src17_by_17.dump | 17 | 17 | 15 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src17_by_17.dump | 17 | 17 | 16 | 1 | unsupported-serverguc | - | pg_restore: error: could not execute query: ERROR:  unrecognized configuration parameter "transaction_timeout" |
| src17_by_17.dump | 17 | 17 | 17 | 0 | supported | ok |  |

### Tally

- dump cells measured: 9 (supported 6)
- restore cells measured: 54
- restore supported (shape-verified): 18
- restore unsupported-toolgate: 24
- restore unsupported-serverguc: 12
- shape checks executed: 18

## 边界与证据定位（重申）

- 本矩阵的证据定位是 **CI/reproducibility**，不是生产就绪证据（`D-017` §3 约束②，用户裁决）。
- 结论**仅**由本次运行支持，且绑定上表记录的版本串（镜像 tag 可变；`D-001` §1 要求每次重新记录）。
- `unsupported` 组合是**工具/服务端行为事实**，本附件不对其作产品级「支持承诺」；如 R3-D 退出矩阵需要该层面的书面表述，按 P-004 另行裁决。
- 破坏性动作只作用于容器内新建的一次性 database 与宿主临时目录；常驻实例未被本矩阵修改。

## 原始运行日志

- 驱动退出码 0；`--- PASS: TestPGCrossVersionRestoreMatrix (105.30s)`。
- 源库建立日志（驱动 `t.Logf`）：三个 server 均报 `migrations at head 87`，播种的规范形分别读回 `.900000Z` / `.123456Z` / `.914000Z`。

## 修正后复跑（响应 A-002 F-I-001～F-I-003，2026-09-21）

在按 D-003 字面化 sentinel 探针、增加 44 张分母表存在性检查、收窄 matrixClassify（并新增 oracle 测试）之后，**重新完整运行**门控矩阵：

- 结果与本节上方表格**一致**：dump 9 格（6 supported）/ restore 54 格（18 supported 且形状校验全通过 / 24 toolgate / 12 serverguc / **0 unexpected**）。
- 三个源库仍分别报 migrations at head 87, 87 ledger rows，规范形样本逐字节一致。
- --- PASS: TestPGCrossVersionRestoreMatrix (105.07s)。

实测版本（复跑，与首次运行相同）：

| server | image | server version | client tools |
|---|---|---|---|
| 15 | `postgres:15-alpine` | 15.19 | 15.19 / 15.19 |
| 16 | `postgres:16` | 16.15 (Debian 16.15-1.pgdg13+2) | 16.15 (Debian 16.15-1.pgdg13+2) / 16.15 (Debian 16.15-1.pgdg13+2) |
| 17 | `postgres:17-alpine` | 17.11 | 17.11 / 17.11 |

复跑摘要：

- dump cells measured: 9 (supported 6)
- restore cells measured: 54
- restore supported (shape-verified): 18
- restore unsupported-toolgate: 24
- restore unsupported-serverguc: 12
- shape checks executed: 18