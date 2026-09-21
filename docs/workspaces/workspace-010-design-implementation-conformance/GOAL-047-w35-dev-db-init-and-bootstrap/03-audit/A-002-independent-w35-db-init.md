---
id: A-002-independent-w35-db-init
doc: audit-entry
status: active
parent: GOAL-047-w35-dev-db-init-and-bootstrap
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
source: independent
provider: grok-build
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-21
scope: GOAL-047 检查点 A/B（e2e-pgset 自举、cmd/dbsetup、dev.cmd init-db、3D000 hint、文档、reset→init→start 证据、相邻入口、治理台账）；不含 checkpoint C 关门
audit_type: execution-facts
verdict: conditional
---

# A-002 · independent 审计：W35 检查点 A/B

- **source**: `independent`
- **provider / auditor**: grok-build（model grok-4.6 · reasoning high）
- **日期**: 2026-09-21
- **scope**: GOAL-047 检查点 A/B
- **基线**: HEAD `a246dfb5`；工作树 clean；未改仓库文件
- **verdict**: `conditional`
- **开放 required**: **3**
- **A/B 是否无条件成立**: **否**
- **是否可进入 C**: **否**
- **是否需用户 P-004**: **否**（除非要把下面 required 降为 `accepted-residual` / `user-overruled`）

本意见不修改 `status` / `progress` / 方案正文 / goal-tree。响应归 `/govern`。

## 范围与区间

只审 workspace-010 的 `GOAL-047-w35-dev-db-init-and-bootstrap` 检查点 A/B。未读取其他工作区上下文。共享资料目录为 `none`。信息项 I-047-001 / I-047-002 均为 `non-blocking` 且已标 verified；本审不因它们放行 A/B。

独立复跑约束：未访问生产库；未删除共享实例上的 `schema_ui_dev` / `schema_ui_test`（避免污染用户现库）。一次性 scratch 库 `schema_ui_w35_audit_scratch` / `schema_ui_w35_fwd` / `schema_ui_w35_cwd` 已创建并清理；事后 `e2e-pgset list` = 0 条 `schema_ui_e2e_*`。

## 成果（有证据）

| # | 判据 | 独立结论 | 证据 |
|--:|------|----------|------|
| 1 | `DB_*` / `PG_TEST_*` namespace 分离 | **成立** | `ServerForRole` 按 role 读两套键；`TestServerForRoleUsesTheMatchingNamespace` 锁定。本轮 `go test -count=1 ./internal/pgsetup/` ok（含真实 PG `TestEnsureIsIdempotentAndDropRemoves`） |
| 2 | 维护连接目标库 → `postgres` → `template1`，且不用于业务连接 | **主路径成立；template1 第三候选仅有单元证据** | 候选表由 `TestMaintenanceCandidatesBootstrapOnEmptyServer` 锁定。本轮 `DB_NAME=schema_ui_w35_missing_zzz go run ./cmd/e2e-pgset list` **exit 0**（不再 3D000）。`openStore` / `openPostgres` 只用 `cfg.DBDSN` Ping，从不调用 `ConnectMaintenance`。未在共享实例上停用 `postgres` 来实打 `template1` 第三段 |
| 3 | `cmd/dbsetup` 默认一次建齐 dev+test、幂等、不删库 | **默认路径成立** | 本轮无额外参数 `go run ./cmd/dbsetup`：`exists schema_ui_dev (dev)` + `exists schema_ui_test (test)`，`0 created, 2 already present`。包注释与 usage 写明 never drops；`Drop` 只有 `e2e-pgset drop` 与测试调用 |
| 4 | 名称校验防注入 | **成立** | `identRe = ^[a-z_][a-z0-9_]*$`；`CREATE/DROP DATABASE` 仅在 `ValidName` 之后拼接。本轮 `e2e-pgset create 'bad;name'` → `invalid database name` / exit 2 |
| 5 | `e2e-pgset create` 改幂等、`drop` 缺失返回 `absent`、`list` 仍只看 `schema_ui_e2e_*` | **成立，不破坏既有调用方** | 本轮：`create` → `created`；再 `create` → `exists`（exit 0）；`drop` → `dropped`；再 `drop` → `absent`（exit 0）。`playwright.config.ts` 只判 `create.status !== 0`，幂等更安全。`list` 仍是 `datname LIKE 'schema_ui_e2e_%'` |
| 6 | 3D000 hint 只改提示、保留 error chain、不泄漏秘密、其它错误不加 init-db 文案 | **成立** | `missingDatabaseHint` 仅 `errors.As` + `Code=="3D000"`；`fmt.Errorf("%w\n...", err)`。`TestMissingDatabaseHintIsActionableAndPreservesClassification` PASS：`28P01` 原样返回。故意连 `127.0.0.1:1` 并注入 canary 密码：`PASSWORD_LEAK=NO`、`DSN_IN_ERROR=NO`；错误只含 `user=` / `database=` / host:port / 候选库名 |
| 7 | 文档主命令与实现一致（SQLite / PG / e2e 边界、CREATEDB、迁移由启动/测试应用） | **主命令一致；配置来源有缺口（见 F-002）** | README / QUICKSTART / `apps/api/README.md` / `dev.cmd help` 均写 `dev.cmd init-db` 与 `cd apps/api && go run ./cmd/dbsetup`。SQLite 不建库、e2e 自管、CREATEDB、dbsetup 不迁移，均有明文 |
| 8 | CWD 下能找到 `configs/.env` | **成立** | `LoadEnvFile` 用 `runtime.Caller` 上寻 `AGENTS.md`，不依赖 CWD。本轮从 `%TEMP%` 执行 `go run -C apps/api ./cmd/dbsetup --dev-only schema_ui_w35_cwd` 仍连到真实 dev 库并创建附加名，随后已 drop |
| 9 | 定向测试 | **成立** | `go test -count=1 ./internal/pgsetup/ ./internal/store/ ./internal/composition/` → 三包 ok（1.1s / 134s / 65s） |
| 10 | 默认全仓 Go suite | **未能一次复现 self 的 65/65** | 本轮 `go test -count=1 ./...` **FAIL**：`TestShutdownDrainHarnessPostgres`（EOF during drain）。同测试立即单跑 **PASS**（11.4s）。判定为并行压测 flake，不是 W35 回归；因此不能背书 B 的「全仓回归绿」 |

未独立复做：删除共享 `schema_ui_dev`/`schema_ui_test` 后再 `dev.cmd start`。E-002 的 reset→start 叙述与当前代码路径一致，但本审用等价 scratch / `DB_NAME` 覆盖 / 幂等 exists 代替，避免动现库。

## 对照成功标准

| 检查点 | 00-meta 自称 | 本审 |
|--------|--------------|------|
| **A** | completed | **不能无条件成立**。自举、默认 `init-db`、hint、名称安全、e2e 幂等均有独立证据；但「同时接受显式库名/附加库名」在 `dev.cmd init-db` 上失败（F-001），且库名解析与 API/`config.yaml` 不一致（F-002） |
| **B** | completed | **不能无条件成立**。README/QUICKSTART 主命令正确；reset→start 未独立复做；goal-tree/Root/`02-execution.md` 仍写 `0/3` 而 `00-meta` 已 `2/3`（F-003）；全仓 `go test ./...` 本轮未一次绿 |
| **C** | pending | 正确保持 pending。本意见为 conditional，required 未闭合前不得进入关门 |

`00-meta` A 行引用 `TestMatrixClassify` 不是本目标测试（属 workspace-040 `TestMatrixClassifyOracle`）。E-002 把 `recovery_wiring_test.go` 算进本波次回归，但 `a246dfb5` diff **不含**该文件。B 行写 Go 64/64，A-001 写 65/65。这些是证据卫生问题，见 recommended。

## Findings

### F-001 · required · med · `dev.cmd init-db` 参数转发静默丢弃第一项

- **证据**：`dev.cmd` 在解析到 `CMD` 后立刻 `shift`，`:init_db` 却执行 `go run ./cmd/dbsetup %2 %3 %4 %5 %6 %7 %8 %9`。转发窗口从 `%2` 起，跳过 shift 后的 `%1`。
- **本轮复现**：`e2e-pgset drop schema_ui_w35_fwd` → `absent`；仓库根 `dev.cmd init-db schema_ui_w35_fwd` 只打印 `exists schema_ui_dev` / `exists schema_ui_test`（**未出现 extra 名**）；随后 `e2e-pgset create schema_ui_w35_fwd` 得到 **`created`**（证明 wrapper 没建这个库）；已 drop 清理。
- **影响**：`dev.cmd init-db --dev-only` / `--test-only` / 附加库名会被静默忽略，仍走默认 dev+test。不删库，但与 GOAL-047 范围第 2 条、D-001、I-047-001「同时接受显式库名」不符。`go run ./cmd/dbsetup --dev-only <extra>` 本身可用（本轮从 TEMP 验证）。
- **最小关闭动作**：`:init_db` 在 shift 之后改转发 `%*`（或 `%1` 起），并加一条真实 `cmd` 调用断言（`--dev-only` 与 extra name 至少各一条）。修完后独立复跑 `dev.cmd init-db <scratch>` 应变为 `created <scratch>`。

### F-002 · required · med · 快捷方式不读 `config.yaml` `db.name`，无 `DB_NAME` 时默认建 `postgres`

- **证据**：API `config.Load` 的库名是 yaml `db.name`（本仓 `schema_ui`），再用 `DB_NAME` 覆盖。`pgsetup.ServerForRole(RoleDev)` 只读 `DB_NAME`，缺省 **`postgres`**，从不读 yaml。`.env.example` 写明「database NAME belongs in configs/config.yaml」，并把 `# DB_NAME=schema_ui` 注释掉。
- **后果**：按文档复制 `.env.example`、只填 host/user/password、名称留在 yaml 时，`dev.cmd init-db` 会 `Ensure postgres`（已存在），然后 `dev.cmd start` 仍对 yaml 的 `schema_ui` 报 3D000。本机之所以绿，是因为 gitignored `.env` 已有 `DB_NAME=schema_ui_dev`（本轮 `dbsetup` 实际确保的也是这两个 env 名，不是 yaml 的 `schema_ui`）。
- **03-audit 待复审 #3** 要求「与 `.env`/`config.yaml` 的库名来源一致」。当前只对齐了 env，没有对齐 yaml 这条已文档化的权威来源。
- **最小关闭动作**（择一并测）：让 RoleDev 名称解析与 `config.Load` 相同（yaml `db.name` ← `DB_NAME`）；**或** `DB_NAME` 为空时 fail closed，禁止默认 `postgres`；并改 `.env.example` / README，使「init 建哪个库」与 API 实际连接的库逐字同一。用「仅 yaml 有 name、env 无 `DB_NAME`」的夹具复跑。

### F-003 · required · low · `progress` 台账不一致（AGENTS §7）

- **证据**：同一提交 `a246dfb5` 里，`GOAL-047/00-meta.md` 为 `progress: 2/3` 且 A/B = completed；`goal-tree.md` 树与表仍是 **`0/3`**；Root `GOAL-001/00-meta.md` W35 行仍是 **`0/3` · 立项**；`02-execution.md` 当前事实仍写 `progress: 0/3`。
- **AGENTS §7**：改 `progress` 必须同步当前工作区 `goal-tree.md`（树 + 表）。只改目标文件不改树 = 任务未完成。
- **最小关闭动作**：在 A/B 按本意见修正后重算检查点，把 `00-meta`、`goal-tree`、Root 波次表、`02-execution.md` 写成同一派生值。在 required 未闭合前，**不要**把树改成 `2/3` 来「对齐」一个尚未成立的 A/B。

### F-004 · recommended · roadmap 仍写「待立项 / 未实现」

- `docs/vision/roadmap.md` §一仍是「待立项 · dev/test PG 数据库初始化快捷方式」/「登记（未实现）」/「`/vision` 先定结构」。用户已裁决挂 workspace-010，GOAL-047 已实施。登记约定要求同步本节。
- **最小关闭**：把该行改为已由 `[workspace-010] GOAL-047` 承接（status 与 C 门禁一致：实现已落码、A/B 待修、C 未过），不要继续写「未实现」。

### F-005 · recommended · 证据卫生与全仓绿未独立确认

- `00-meta` A 引用 `TestMatrixClassify`（他区测试）；E-002 把未进入 `a246dfb5` 的 `recovery_wiring_test.go` 算进本波次；B 写 64/64、A-001 写 65/65。
- 本轮全仓 `go test ./...` 首跑 FAIL `TestShutdownDrainHarnessPostgres`；单跑 PASS。与 W35 无直接因果关系，但不能复述 self「65/65 · 0 fail」。
- **最小关闭**：删掉错误引用；全仓数字以一次可复跑日志为准；B 不要在 flake 未隔离时写「全仓绿」。

### F-006 · recommended · `Drop(name)` 在 `Server.Database == name` 时会连上目标库再 DROP

- `ConnectMaintenance` 先连配置库。`Drop` 在目标即配置库时，会对自己 `DROP DATABASE ... WITH (FORCE)`，PostgreSQL 不允许删除当前连接所在库。
- 集成测试 drop 的是 `schema_ui_w35_scratch`，不是 `PG_TEST_DB`，所以绿。`dbsetup` 不调用 `Drop`。`e2e-pgset drop schema_ui_e2e_*` 通常也不是 `DB_NAME`。
- **最小关闭**：维护连接跳过正要 DROP 的库名；或候选顺序对 `Drop` 固定为 `postgres`/`template1`。

### F-007 · recommended · 相邻文档/注释

- `.env.example` 的 PG 段没有 `init-db` / 3D000 指引（这是复制模板后的真实入口）。
- hint 写 `see README/QUICKSTART "initialize the databases"`，实际标题是「PostgreSQL 数据库初始化」。
- `e2e-pgset` 仍注释 `DROP DATABASE name (no active conns)`，实现已是 `WITH (FORCE)`。
- `workspace.md` 波次表停在 W34，W35 只出现在 Root/`goal-tree`。
- CI `r6-basic-matrix.yml` 无常驻 PG service，`pgsetup` 集成测试在 CI 会 skip；**不**需要为本目标改 CI 模板。`compose.yaml` 无 Postgres，SQLite 路径不需要 init-db。

## 必改项汇总

| ID | 级别 | 一句话 | 阻断 |
|----|------|--------|------|
| F-001 | required | `dev.cmd init-db` 附加参数被 shift 后的 `%2` 丢掉 | A（附加库名/flag） |
| F-002 | required | 不读 yaml `db.name`，无 `DB_NAME` 时默认 `postgres` | A/B 文档与真实库名 |
| F-003 | required | `00-meta` 2/3 vs goal-tree/Root/execution 0/3 | 台账；C 前必须对齐 |

开放 required = **3**。无与 A-001 对同一 finding 一要一否的冲突：self 把转发列为 recommended `F-S-001`～`F-S-004`；本审把转发升级为 required，并新增库名来源与台账两项。编排器按 P-003 响应即可，不必先走 P-004。

## 待复审 7 条逐条结论

| # | 事项 | 结论 |
|--:|------|------|
| 1 | 空实例自举 + 不改变业务连接 | **通过（有限度）**。`DB_NAME` 不存在时 `e2e-pgset list` 可回退维护库。业务路径仍只 Ping `cfg.DBDSN`。`template1` 在 `postgres` 不可用时的活路径本轮未打。维护连接失败会遍历全部候选（不仅 3D000）；同凭据下认证失败会三次都失败，不会悄悄改业务 DSN |
| 2 | 快捷方式幂等且不破坏既有库 | **默认路径通过**。`Ensure` 不 drop；`dbsetup` 无删除入口。附加名经 `dev.cmd` **不执行**（F-001）。`e2e-pgset drop` 仍是显式独立动作，缺失为 `absent`/exit 0，合理 |
| 3 | 一次建齐 dev+test，库名来自配置不硬编码 | **部分失败**。有 `DB_NAME`/`PG_TEST_DB` 时一次建齐且不硬编码 `schema_ui_dev`。不读 `config.yaml` `db.name`，缺 `DB_NAME` 时落到 `postgres`（F-002） |
| 4 | 启动错误不泄漏秘密、不改分类 | **通过**。仅 3D000 包 hint；`%w` 保留 PgError；composition 仍是 `LIFECYCLE_START_FAILED [core.auth-session]`。canary 密码与 DSN 未进错误文本 |
| 5 | 文档与真实命令逐字一致，覆盖多方言 | **主命令通过；来源说明不够**。SQLite/PG/e2e、CREATEDB、迁移由启动/测试应用均有。缺「yaml name vs `DB_NAME`」和「`dev.cmd` 实际不转发 flag」 |
| 6 | 可复跑删库 → init → start | **不能视为独立已复现**。本审未删共享 `schema_ui_dev`/`schema_ui_test`。已独立复现：缺失 `DB_NAME` 自举、scratch Ensure/Drop 幂等、默认 `init-db` exists。E-002 的 start/readyz 叙述本轮未复跑 |
| 7 | 是否越界 | **未越界**。`a246dfb5` 20 files：pgsetup/dbsetup、e2e-pgset 瘦身、hint、dev.cmd、三份 README、GOAL-047 台账、goal-tree 插入行。无迁移链/checksum/codec/wire/备份 Port；`openStore` 仍不 `CREATE DATABASE` |

## 与 A-001 self 的异同

- **同意**：自举主路径、默认幂等、hint 边界、e2e create 幂等不破坏 Playwright、LoadEnvFile 双 namespace 仍由 `ServerForRole` 分流、不自动建库。
- **升级**：self `F-S-004`（recommended）在本轮被打成 **required F-001**（真实 cmd 复现，不是假设）。
- **新增 required**：F-002 库名来源、F-003 台账。self 写「开放 required = 0」在独立复跑后不成立。
- **self 高估**：A/B completed、`progress 2/3`、全仓 65/65、`TestMatrixClassify` 引用。

## 结论 + 建议给编排器/用户的下一步

检查点 A/B **不能放行**。核心缺陷（空实例无法自举、没有 init 入口、3D000 无下一步）在默认 `DB_NAME` 已写入 `.env` 的本机路径上已经修掉，但交付合同上的附加库名 wrapper、与 yaml 权威库名对齐、以及 progress 台账都还没过门禁。

**不要进入 C，不要把 GOAL-047 标 done。**

建议 `/govern` 下一句：响应 A-002，按 F-001 → F-002 → F-003 修代码与台账，再请独立复审关闭这三条 required。F-004～F-007 可同批处理。无需先问用户，除非要把其中 required 接受为残余。

## 声明

本意见 `source: independent`。不修改 status/progress/方案正文/goal-tree。响应由 `/govern` 处理。
