---
id: A-002-independent-r3c-checkpoints
doc: audit-entry
status: active
parent: GOAL-007-r3-pg-cross-version-restore-matrix
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
---
# A-002 · independent 审计：R3-C 检查点 A/B（全量，非复审判定）

- **source**: `independent`
- **provider / auditor**: grok build · model grok-4.6 · reasoning high
- **日期**: 2026-09-21
- **scope**: `GOAL-007-r3-pg-cross-version-restore-matrix` 检查点 A（`D-001` 矩阵定义与判定口径）与 B（`E-003` 实测、`D-002` 规则更正、`I-041-004` 收口）。边界：Root `D-018` §2 第 4 项 / §4 R3-C / §5 `I-041-004`；`D-017` §3 约束②。不含 R3-D，不重开 R1/R2/R3-A/B。
- **类型**: execution-facts + finding-closure（对 `A-001` 待复审 12 条与 `F-S-003`～`F-S-006`、`N-001`～`N-005` 的独立取证）
- **verdict**: `conditional`（**开放 required = 0**；4 条 recommended。检查点 A/B 的技术主张经独立复跑成立，但冻结口径与实现有一处未入账偏差，且 `00-meta` / `goal-tree` 未与 `progress: 2/3` 对齐，不能无条件静默关门）

本意见不修改 `status` / `progress` / 方案正文 / `goal-tree`。仓库工作树在审计前后均为 clean（`HEAD = 358455f6`，其前 `9c6268c4`）。响应归 `/govern`。

---

## 成果（独立核实）

| # | 判据 | 独立结论 | 证据 |
|--:|------|----------|------|
| 1 | HEAD 与工作树 | **达成** | `git log --oneline -3` → `358455f6` / `9c6268c4` / `613afdc2`；`git status --short` 空 |
| 2 | 口径先于矩阵本体冻结 | **达成** | `D-001`（检查点 A）先于 `E-003`；`I-041-009` verified。探测附件用平凡 schema，不构成按结果反推口径 |
| 3 | 组合覆盖而非只跑同版本 | **达成** | 独立复跑：dump **9** 格、restore **54** 格（6 归档 × 3 client × 3 target），无跳过 |
| 4 | unsupported 逐格记录 | **达成** | 复跑表 36 个非 supported 格均有 exit + 首条诊断；`unexpected-failure` / `setup-failure` 计数为 0（出现即 `t.Fatalf`） |
| 5 | 升级后恢复有界核对 | **达成（抽样口径内）** | 18 个 supported 格（含 15→16/17、16→17）`shape = ok`。源库 **各自**建在 15.19 / 16.15 / 17.11 上，非 15 库再恢复 |
| 6 | 驱动 fail-closed（未知失败不记成 unsupported） | **达成（见 F-I-001）** | 仓库外 25 条文本反例：网络/磁盘/权限/认证均落入 `unexpected-failure`。`does not exist` 过宽 → `setup-failure`，仍 `Fatalf`，**不是** fail-open 进 toolgate/serverguc |
| 7 | 空库 / 无数据行不能骗过形状校验 | **达成** | 一次性容器：缺 `schema_migrations` → 查询失败；仅 ledger 无 `users` → 查询失败。对应 `matrixReadFacts` `t.Fatalf`（`pg_cross_version_matrix_test.go:255-272`） |
| 8 | 破坏性动作限于一次性库 | **达成** | 驱动只用容器内新建库 + `t.TempDir()`；不读 `PG_TEST_*`。复跑后 `docker ps -a` / `docker network ls` **无** `vp040mx*` |
| 9 | 未越读为生产就绪 | **达成** | `D-002` §4、附件「边界」节、`D-001` §7 均写 CI/reproducibility；C3 生产路径旗标与矩阵一致（`provider_pg.go:188`：`--exit-on-error --no-owner`） |
| 10 | `D-002` 三条规则覆盖 18+36 格 | **达成** | 逐格核对无漏格、无反例（见下） |
| 11 | 附带结论：迁移链在 16.15 / 17.11 至 v87 | **达成** | 复跑 `t.Logf`：三个 server **分别** `migrations at head 87, 87 ledger rows`，播种值读回 `.900000Z` / `.123456Z` / `.914000Z`。`store.Open` 对容器 published port 走 `Identify→Plan→Execute`（`internal/store/open.go:18-21`），catalog 来自 `compiledmodules.PersistenceCatalog()`（`restore_harness_test.go:76-82`） |
| 12 | 默认基线不依赖 docker | **达成** | 未设 `VP040_PG_MATRIX` → `--- SKIP: TestPGCrossVersionRestoreMatrix (0.00s)` |
| 13 | 独立复跑与附件一致 | **达成** | `PASS (104.81s)`；版本串、9/54、18/24/12、header 1.14/1.15/1.16 与 `attachments/r3c-pg-cross-version-matrix-v0.1.md` 一致 |

独立复跑命令（2026-09-21，本会话）：

```text
cd apps/api
$env:VP040_PG_MATRIX="1"
$env:VP040_PG_MATRIX_OUT="$env:TEMP\a002-matrix.md"
go test -count=1 -run TestPGCrossVersionRestoreMatrix -v -timeout 45m ./internal/backup/
```

复跑源库日志（一手，非采信叙述）：

```text
server 15 (15.19): migrations at head 87, 87 ledger rows, samples 2026-09-20T12:57:15.900000Z / 2026-01-02T03:04:05.123456Z / 2026-09-20T12:57:15.914000Z
server 16 (16.15 (Debian 16.15-1.pgdg13+2)): migrations at head 87, 87 ledger rows, samples …
server 17 (17.11): migrations at head 87, 87 ledger rows, samples …
```

`D-002` 规则 vs 54 格（无漏）：

```text
dump:     client_major >= server_major          → dump 6 supported + 3 toolgate
restore:  client_major >= dumper_client_major   → 24 toolgate
          client_major == 17 ⇒ target == 17     → 12 serverguc
其余                                      → 18 supported（均 shape ok）
```

18 = `src15_by_15`(7) + `src15_by_16`(4) + `src15_by_17`(1) + `src16_by_16`(4) + `src16_by_17`(1) + `src17_by_17`(1)。含 `src16_by_16` → client 16 → target **15**（`D-001` 外推句的反例格，复跑仍 supported）。

---

## Findings

### F-I-001 · recommended · `matrixClassify` 的 `does not exist` 过宽

- **问题**：未知失败**不会**被记成 `unsupported-toolgate` / `unsupported-serverguc`（不 fail-open）。但 `strings.Contains(out, "does not exist")`（`pg_cross_version_matrix_test.go:379-381`）会把 `role/schema/extension/relation/function … does not exist` 标成 `setup-failure`。当前仍 `t.Fatalf`，矩阵不会因此放行；若诊断被读成「环境没搭好」会掩盖真实恢复不兼容（本 schema 含 `CREATE EXTENSION IF NOT EXISTS citext`，`modules/authsession/migration/migration.go:443`）。
- **证据**：仓库外 25 条 oracle（拷贝 `matrixClassify`，未改仓库）。Connection refused / timeout / disk full / permission denied / password authentication failed → `unexpected-failure`。`ERROR: extension "citext" does not exist`、`ERROR: role "app" does not exist`、以及「连接断开 + database does not exist」混排 → `setup-failure`。`"unsupported version"` 缺 `in file header` 或反之 → `unexpected-failure`（子串门未误触发）。
- **关闭要求**：把 `does not exist` 收到明确的 database/连接 setup 句式；对象/扩展/角色不存在走 `unexpected-failure`。不阻塞本目标用「本轮 0 unexpected」作为分类鲁棒性的完整证明。

### F-I-002 · recommended · `D-001` §4 第 4 项未按落盘列执行

- **问题**：冻结口径要求 sentinel 核 `mail_config.updated_at` / `telegram_config.updated_at` 保持 NULL（`D-001` 第 59–64 行）。驱动实际核的是播种行 `users.locked_until`（`pg_cross_version_matrix_test.go:317-323`）；`mail_config.updated_at` 只做 `timestamptz(6)` 类型抽样；`telegram_config` 完全未触及。`D-002` 未记录这次替换。C3 生产校验 `verify.go:337-341` 三者都查。
- **证据**：矩阵源库由 **v87 catalog 直接 Apply**，不是 C3 harness 那种「先插入 legacy `0` 再转换」。`mail_config` 在 corepersistence 里是建表（`migration.go:94-106`），新鲜库上这两列的「保持 NULL」可能空转。`users.locked_until` 的 NULL 播种反而更硬。偏差在**口径追溯**，不是「NULL 复活漏检」。
- **关闭要求**：在新决策条目（或给 `D-002` 追加修订段）写明 §4.4 的实际探针是 `users.locked_until`（mx-u1 NULL / mx-u2 非空），并说明为何不采用空的 singleton 表。不要回改 `D-001` 正文。

### F-I-003 · recommended · 四项抽样不能检测「未抽样对象缺失」

- **问题**：空库、只恢复 ledger、缺 `users`/`jobs` 行、精度降为 3、微秒截断，四项校验（或 `matrixReadFacts` 的 Fatalf）**都能抓住**。但只要抽样对象在、值对、类型为 `timestamptz(6)`，其它 VP-040 表可以全部不存在，四项仍全绿。这是对**校验函数**的攻击，不是已证明 `pg_restore --exit-on-error` 对完整 TOC 会静默丢表。
- **证据**：一次性容器 `postgres:15-alpine`（名 `a002-shape`，已删）手建库：`schema_migrations` 1–87、`users`/`jobs`/`mail_config` 列与三行播种值，**无** `telegram_config` / `login_failures`。四项查询全部通过；`to_regclass('telegram_config')` 为空。同容器：空库 → `schema_migrations` 不存在；仅 ledger → `users` 不存在；`timestamp(3) without time zone` 与 `.123000Z` 会被类型/值项抓住。
- **关闭要求**：接受为 `D-001`「抽样」的已知边界即可（`--exit-on-error` 下完整 dump 丢 TOC 项通常非 0）。若要把「部分对象」从理论缺口收掉，至少加表计数或 `telegram_config` 存在性。不要求为关门重跑 54 格。

### F-I-004 · recommended · 台账未与检查点 B 对齐

- **问题**：`00-meta.md` `progress: 2/3` 且检查点 B 文案写 `I-041-004 → verified`，但同文件信息表仍写 `I-041-004` **open**（第 65 行）。`goal-tree.md` ASCII 树为 `0/3`，纲领叙述与说明为 `0/3` 或 `1/3`，状态表为 `1/3`。AGENTS §7：改 progress 必须同步树+表。
- **证据**：`goal-tree.md:27,35,38,50,54`；`GOAL-007/00-meta.md:9,56,65`；`01-decision.md` 信息表为 verified。Root `D-018` §5 仍为 open——子目标未关门前可以暂留，但子目标自己的 `00-meta` 不应自相矛盾。
- **关闭要求**：编排器响应本意见时同步：GOAL-007 信息表 `I-041-004 = verified`；`goal-tree` 树+表+叙述 → `2/3`（C 未完成前不要写成 3/3）。Root `I-041-004` 可在子目标关门时一并改。

---

## 逐条结论表

### `03-audit.md` 待复审 12 条

| # | 事项 | 独立结论 |
|--:|------|----------|
| 1 | 组合是否覆盖关键 server×client，而非只跑同版本 | **成立**。复跑 9+54，含全部交叉 major；无跳过 |
| 2 | supported/unsupported 是否事先冻结、可复现 | **成立**。`D-001` 先于 `E-003`；`D-002` 只改外推句。命令/版本串/退出码可复跑复现 |
| 3 | unsupported 是否逐条记录 | **成立**。36 格均有 exit + 诊断；不是「全部支持」概括 |
| 4 | 升级后恢复是否真跨版本且有形状证据 | **成立（抽样范围内）**。18 格含 15→16/17、16→17；四项在抽样对象上复跑通过。见 F-I-002/F-I-003 |
| 5 | 破坏性动作是否限于一次性/专用库 | **成立**。容器新建库 + `t.TempDir()`；常驻 15.4 未入轴、驱动不读 `PG_TEST_*` |
| 6 | 是否把容器结果当成生产就绪 | **未越读**。`D-002` §4 与附件边界节遵守 `D-017` §3 约束② |
| 7 | `I-041-004` 关闭/residual 是否有证据 | **同意 verified**（非 residual）。信息需求是逐格矩阵，已收集。`00-meta` 表未改，见 F-I-004 |
| 8 | `matrixClassify` 是否把未知失败判成 unsupported | **不会 fail-open 进前两类**。网络/磁盘/权限/认证 → `unexpected-failure` → Fatalf。`does not exist` → `setup-failure`（仍 Fatalf）。见 F-I-001 |
| 9 | 四项校验能否被空库骗过 | **不能**。空库/`schema_migrations` 缺失 → Fatalf。见下「反例」 |
| 10 | 结论是否绑在本次镜像与旗标 | **成立**。附件与 `D-002` §4 绑定 15.19/16.15/17.11 与 `-F c --no-owner` / `--exit-on-error --no-owner` |
| 11 | 源库数据能否区分成功 vs 空表 | **能**。3 个固定 id；缺行 → `ErrNoRows` Fatalf。两库都空不能通过 |
| 12 | 16/17 迁移链附带结论是否有独立证据 | **有**。三个源库分别在对应 server 上 `store.Open` 到 head 87，不是复用 15 的库再 restore |

### `A-001` 提交项与非问题项

| ID | 独立结论 |
|----|----------|
| `F-S-003` | **部分成立，降为已知限制**。未知失败不会被记成「已知不支持」。`does not exist` 过宽见 F-I-001。setup-failure 与 unexpected 目前都整体失败，本轮计数 0 不能证明分类标签永远正确 |
| `F-S-004` | **针对「空库 / 只有 ledger / 没有数据行」：驱动不被骗**。针对「只恢复了抽样对象」：校验函数能被手建库骗过（F-I-003）。未证明真实 `pg_restore --exit-on-error` 会产出这种库 |
| `F-S-005` | **同意为 note**。`D-001` §2 的 27 格是 3×3×3（每源一档同 client dump）；实测 6 归档 ×3×3=54，覆盖 strictly 更大。`D-002` §2 已交叉引用。可选把 §2 改写成「归档数 × 3 × 3」，非门禁 |
| `F-S-006` | **同意非缺口**。常驻 15.4 与容器 15.19 同 major；`D-001` §7 明确未选。判据 4 的有界核对按 `D-018` 是固定版本临时容器。C3 同 major 路径由 `provider_pg.go` + harness 承担 |
| `N-001` | **同意非问题**。复跑 `src16_by_16.dump` → 16 → 15 仍 supported/ok。归档来源 server 不是额外目标约束 |
| `N-002` | **同意**。矩阵以实际 restore 为准，不以 `pg_restore -l` 当证据 |
| `N-003` | **同意**。本次复跑重新量了版本串，与附件相同 |
| `N-004` | **同意**。`provider_pg.go:19,188` 与矩阵同为 `--exit-on-error --no-owner` |
| `N-005` | **同意**。复跑后无 `vp040mx` 容器/网络；基线仅 `gf-pg`（Exited，3 weeks，与本矩阵无关） |

`F-S-001` / `F-S-002`：本轮只核 git 中重写后的驱动与 `D-002` 更正链；不重开。`D-002` 不重写 `D-001`、分类未被推翻，合规。

---

## 反例尝试记录

### 1. `matrixClassify`（仓库外副本，25 条）

| 输入（exit≠0 除非注明） | 实际分类 | 审计含义 |
|--------------------------|----------|----------|
| `aborting because of server version mismatch` | toolgate | 真阳性 |
| `unsupported version (1.16) in file header` | toolgate | 真阳性 |
| `unrecognized configuration parameter "transaction_timeout"` | serverguc | 真阳性 |
| Connection refused / timed out / server closed unexpectedly | unexpected | **未** fail-open |
| No space left on device | unexpected | 未 fail-open |
| Permission denied | unexpected | 未 fail-open |
| password / Ident authentication failed | unexpected | 未 fail-open |
| `pg_restore: warning: errors ignored on restore: 3`（exit 1） | unexpected | 未 fail-open |
| warning、exit 0 | supported | 交给形状校验 |
| role/schema/extension/relation/function does not exist | **setup-failure** | 标签过宽（F-I-001），仍 Fatalf |
| 连接断开与 `database does not exist` 混排 | setup-failure | 可能掩盖网络故障 |
| 只有 `unsupported version` 或只有 `in file header` | unexpected | 子串门足够窄 |
| citext `.control` 找不到（无 “does not exist”） | unexpected | 正确 |

`--exit-on-error` 之外的形态：若去掉该旗标，GUC 失败可能变成 warning + 非 0/或继续；口径已冻结为 C3 实际旗标（N-004），不把未测旗标算进 supported。

### 2. `matrixVerifyShape`（一次性 `postgres:15-alpine`，已 `docker rm -f`）

| 反例库 | 四项结果 | 能否骗过 |
|--------|----------|----------|
| 空库（无 `schema_migrations`） | 查询 ERROR | **否**（Fatalf） |
| 仅 ledger 87/87，无 `users` | ledger 过，users 查询 ERROR | **否** |
| 抽样对象齐全、值与类型正确，缺 `telegram_config`/`login_failures` | 四项全过 | **是（校验函数）** |
| `timestamp(3) without time zone` | 类型项失败 | 否 |
| 值为 `.123000Z` 对 `.123456Z` | 值项失败 | 否 |

`ledger rows = 87` 且 `max = 87`、version 唯一时，鸽笼原理下集合即为 1..87，ledger 项对「少版本」较紧。

### 3. 矩阵复跑与清理

| 项 | 结果 |
|----|------|
| 是否复跑 | **已复跑**（不是「未复跑」） |
| 退出 | `PASS` 104.81s / 104.940s |
| 与附件 | 版本串、9/54、18 supported shape ok、24 toolgate、12 serverguc、0 unexpected **一致** |
| 源库 | 15.19 / 16.15 / 17.11 **各自** head 87、87 行 |
| 清理 | `docker ps -a` / `docker network ls` 无 `vp040mx`；形状探针容器 `a002-shape` 已删 |
| 常驻 15.4 | 驱动不引用 `PG_TEST_*`；本会话未连该实例 |
| 仓库 | `git status` clean；临时探针（classify / shape sql / 报告副本）已删 |

---

## 明确结论

| 问题 | 答案 |
|------|------|
| 开放 required 数量 | **0** |
| 是否同意检查点 A 成立 | **同意**。口径与组合边界先于本体冻结；`I-041-009` verified 成立 |
| 是否同意检查点 B 成立 | **同意技术判据**。9+54 逐格、18 形状 ok、规则更正正确且覆盖全部格、迁移链附带结论有一手证据。**不同意**在 `goal-tree`/`00-meta` 信息表未对齐前把「B 的治理落盘」视为完成（F-I-004） |
| 是否同意 `I-041-004` 维持 verified | **同意**。non-blocking 信息需求是逐组合记录，不是生产支持承诺。不存在「应记 residual 却静默关闭」。容器/旗标/未测 minor 是已写明的边界，不是未收集的必需信息 |
| 是否同意可进入检查点 C 关门 | **有条件同意**：本意见落盘后，编排器同步 F-I-004 的台账，并以新决策段记录 F-I-002 的探针替换，即可静默关门。F-I-001 / F-I-003 不阻断 C。关门 ≠ Root 关门，≠ 产品支持矩阵 |
| 需用户 P-004 的点 | **本目标没有必须现在裁的 P-004**。产品级「支持哪些组合」（尤其「17 client 不可恢复到 15/16」「dump 要求 client ≥ server」是否写入发布说明/退出矩阵）应作为 **R3-D 的 P-004**，`D-001` §7 / `D-002` §4 已预留。独立审计不把该裁决提前塞进 GOAL-007 |

`progress: 2/3` 由 A～C 三个显式检查点等权派生，公式正确；C 仍 pending。`D-002` 用新条目更正 `D-001` 外推句、不重写已落盘决策、分类未推翻，合规。`03-audit.md` 索引在本意见落盘前只有 A-001，完整。不存在「未合法闭合 required 却推进 A/B 技术工作」；存在 progress 已写 2/3 而 `goal-tree` 未跟（F-I-004）。

---

## 必改项汇总与最小补齐

**Required：无。**

关门前最小补齐（编排器，不需要重跑矩阵）：

1. 同步 `00-meta` 信息表 `I-041-004 → verified`，同步 `goal-tree` 树+表到 `2/3`（C 完成前不要标 `done`/`3/3`）。
2. 用新决策段记下 `D-001` §4.4 的实际探针（`users.locked_until`），不要改 `D-001` 原文。

可选（非门禁）：收窄 `does not exist`；给形状校验加表存在性。

---

## 与 `A-001` self 的异同

- **同意** self 的 12 项技术成果、fail-closed 方向、`I-041-004` verified、未越读生产就绪、`D-002` 更正方式。
- **加强**：独立复跑 + 分类/形状反例，而不是只读附件。
- **不同意 self 把形状校验写成「足以区分部分恢复」而不加限定**：空库/空表可以挡住；「只恢复抽样对象」挡不住。
- **新增** F-I-002（口径列 vs 实现）、F-I-004（台账漂移）。self 未提。

---

## 声明

本意见 `source: independent`，不修改 status/progress/方案/goal-tree。请用 **`/govern`** 响应本意见、闭合 recommended、同步台账，再决定检查点 C 是否静默关门。
