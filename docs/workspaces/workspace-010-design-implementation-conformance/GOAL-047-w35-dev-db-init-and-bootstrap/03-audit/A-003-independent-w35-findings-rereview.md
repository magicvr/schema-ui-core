---
id: A-003-independent-w35-findings-rereview
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
scope: A-002 F-001/F-002/F-003 定向复审（检查点 A/B 关闭证据；不含 C 关门）
audit_type: finding-closure
verdict: pass
---

# A-003 · independent 定向复审：A-002 三条 required（W35 DB init）

- **source**: `independent`
- **provider**: grok-build（model grok-4.6 · reasoning high）
- **日期**: 2026-09-21
- **类型**: finding-closure
- **scope**: `GOAL-047` 检查点 A/B；A-002 `F-001`/`F-002`/`F-003` 关闭复审；附带 03-audit 七条待复审事项与 A-001 recommended 复核。不含 checkpoint C 关门、不改 `status`/`progress`。
- **verdict**: `pass`（开放 required = 0；A-002 三条 required 均可 `fixed`）

## 范围与区间

- 工作区：`workspace-010-design-implementation-conformance`（Root `GOAL-001-design-implementation-conformance`；`primary_plan` = VP-010；`shared_materials_catalog: none`）。
- 被审目标：`GOAL-047-w35-dev-db-init-and-bootstrap`。
- 基线提交：`a246dfb5`（实现）→ HEAD `e3b90f58`（声称关闭 A-002）；工作树 **clean**。
- 上一轮意见身份：用户指定 `03-audit/A-002-independent-w35-db-init.md`（conditional，required=3）。**仓库内该文件与 `03-audit.md` 索引均不存在**（见 F-004）；本轮按用户指令编号 **A-003**，三条 required 以用户书面复述 + `e3b90f58` 提交说明为核对对象。
- 本轮只读验证：真实 `dev.cmd init-db`、仓库外临时 YAML、`ServerForRole`/`config.Load` 对照、定向 + 全仓 Go 测试。未访问生产库。一次性 scratch 已 drop；探针文件已删除。

## 成果（有证据）

| # | 主张 | 本轮核对 | 证据 |
|--:|------|----------|------|
| 1 | `dev.cmd init-db` 在 `shift` 后转发 `%1..%9` | **成立** | 仓库根 `dev.cmd init-db --dev-only schema_ui_w35_a003`：仅 `dev role`，`created schema_ui_w35_a003`，**无** `test role`（若仍从 `%2` 转发，`--dev-only` 会被丢掉并顺带处理 test）。`--test-only` 仅 `test role` / `exists schema_ui_test`。`--dev-only` 无 extra 只 `exists schema_ui_dev`。 |
| 2 | extra scratch 可创建并可 drop；未知/非法名不误建库 | **成立** | `e2e-pgset drop schema_ui_w35_a003` → `dropped`，再 drop → `absent`。`--no-such-flag` 报 `unknown flag` 且未建库。`--dev-only Bad-Name` 被 `ValidName` 拒绝。结束时 `yaml_only_db` / `schema_ui_w35_a003` 均 `absent`；`e2e-pgset list` = `0 schema_ui_e2e_*`。 |
| 3 | 无 `DB_NAME` 时 dev 读 YAML `db.name`，有则 DB_NAME 优先；test 始终 `PG_TEST_DB` | **成立** | 仓库外 `db.name: yaml_only_db`：`t.Setenv("DB_NAME","")` + `CONFIG_FILE` → `ServerForRole(dev)=yaml_only_db`；`DB_NAME=env_wins_db` 覆盖 YAML；test 在两种情况下都读 `PG_TEST_DB`。`LoadEnvFile()` 后空 `DB_NAME` 仍走 YAML（dbsetup 真实顺序）。`config.Load()` 同样：空 `DB_NAME` → `yaml_only_db`，有 env → `env_wins_db`。 |
| 4 | 本机 dbsetup 与 API 实际目标库一致 | **成立** | gitignored `.env`：`DB_NAME=schema_ui_dev`，`PG_TEST_DB=schema_ui_test`；committed `config.yaml` `db.name: schema_ui`（无 `DB_NAME` 时的 fallback）。真实 `init-db` 确保的是 `schema_ui_dev` / `schema_ui_test`，与 `config.Load` 的 `envOr("DB_NAME", yaml)` 一致，而不是维护库 `postgres`。 |
| 5 | 进度 2/3 与 A/B completed 已同步到治理主面 | **成立（主面）** | 见 F-003。`goal-tree.md` `updated: 2026-09-21` / `version: 0.60.0`。 |
| 6 | 回归 | **成立** | `go test -count=1 ./internal/pgsetup/ ./internal/store/ ./internal/composition/` 全 ok；`go build ./...` exit 0；全仓 `go test -count=1 ./...` **65 ok / 0 FAIL**。 |

## 对照成功标准

| 检查点 | 判据 | 本轮 |
|--------|------|------|
| A | 自举 + 快捷方式 + 3D000 hint + 可执行测试 | **仍成立**（代码 + 本轮真实 init-db + pgsetup/store 测试） |
| B | 文档 + reset→init→start 证据 + 回归绿 | **仍成立**（文档命令一致；reset→start 本轮未再破坏常驻库，沿用 E-002；本轮全仓 65/65） |
| C | self + independent 落盘、required 闭合 → 关门 | **pending**（本意见通过后由编排器落盘/闭合；本轮不关门） |

`progress: 2/3` 不得作为 finding 闭合或放行依据；闭合依据是下面的逐条证据。

## A-002 required 逐条判定

### F-001 · `dev.cmd` 参数转发丢掉 `%1` → **fixed**

- **原缺陷**：`:init_db` 在顶层 `shift` 之后仍从 `%2` 转发，`--dev-only` / `--test-only` / extra name 被静默丢掉。
- **修复**：`dev.cmd` L316 现为 `go run ./cmd/dbsetup %1 %2 %3 %4 %5 %6 %7 %8 %9`。
- **独立实证（cmd.exe `shift` 语义）**：
  - `.\dev.cmd init-db --dev-only schema_ui_w35_a003` → 仅 dev；`created schema_ui_w35_a003 (dev)`；**没有** test 行。这是 `%1` vs `%2` 的判别实验。
  - `.\dev.cmd init-db --test-only` → 仅 test；`exists schema_ui_test (test)`。
  - `.\dev.cmd init-db --dev-only` → 仅确保 `schema_ui_dev`，无 extra。
  - 未知 flag / 非法名不建库（见上）。
- **残余（recommended，不阻断）**：`go run` 在 Windows 上把 dbsetup 的 `os.Exit(2)` 收成 wrapper `exit 1`；`dev.cmd` 对 usage 错误仍打印 “Check DB_HOST…”。不影响转发正确性。

### F-002 · dev 库名只读 `DB_NAME`、缺省 `postgres` → **fixed**

- **原缺陷**：`ServerForRole(dev)` 用 `envOr("DB_NAME", "postgres")`，无 `DB_NAME` 时可能 init 维护库，而 API `config.Load` 走 YAML `db.name`（本仓为 `schema_ui`）。
- **修复顺序（dev）**：`DB_NAME`（非空）→ `yamlDBName()`（`CONFIG_FILE` 或 `apps/api/configs/config.yaml`）→ `"schema_ui"`。test **仍**只读 `PG_TEST_DB`，缺省 `postgres`。
- **独立实证**：
  1. 仓库外临时 YAML `db.name: yaml_only_db` + 空 `DB_NAME` → `ServerForRole(dev) = yaml_only_db`。
  2. 同 YAML + `DB_NAME=env_wins_db` → `env_wins_db`（env 优先）。
  3. 上述两种情况下 `ServerForRole(test)` 都不吃 YAML，只吃 `PG_TEST_DB`。
  4. `LoadEnvFile()` 后空 `DB_NAME` 仍得 `yaml_only_db`（与 `cmd/dbsetup` 调用顺序相同）；`LoadEnvFile` 不覆盖已设置的 `DB_HOST` / `PG_TEST_HOST`；`.env` 同时含 `DB_*` 与 `PG_TEST_*`。
  5. `config.Load()` 在空 `DB_NAME` + 同一 YAML 时 `DBName=yaml_only_db`，有 env 时 `env_wins_db` —— **与 dbsetup 一致，不是只采信 pgsetup 单测**。
  6. 本机真实 CLI：`.env` 有 `DB_NAME=schema_ui_dev` 时，`init-db` 目标就是 `schema_ui_dev`，不是 `postgres`，也不是 YAML 字面量 `schema_ui`。这与 API `envOr("DB_NAME", cfg.DBName)` 一致。
- **残余（recommended，不阻断）**：`yamlDBName` 不走 `config.Load` 的 `${VAR}` 插值；本仓 `config.yaml` 是字面量 `schema_ui`，插值失败会变成非法 ident 被 `ValidName` 拒绝（fail-closed），不会再静默落到 `postgres`。`Server` 结构体注释仍写 “defaulting to postgres”（见 F-006）。

### F-003 · 进度面 0/3 vs 2/3 → **fixed**（治理主面一致；一处正文残留见 F-004）

逐字核对：

| 位置 | A/B | progress |
|------|-----|----------|
| GOAL-047 `00-meta` 检查点表 | A **completed**；B **completed**；C pending | frontmatter `progress: 2/3` |
| `goal-tree.md` 树 | 「A/B 已完成，C 待 self + grok independent 审计」 | `(2/3)` |
| `goal-tree.md` 表 | （描述为实现已完成） | `active` / `2/3` / 2026-09-21 |
| Root `GOAL-001/00-meta` 波次台账 W35 | 「A/B 完成，待 self + grok independent 审计」 | **active（2/3 · 2026-09-21）** |
| GOAL-047 `02-execution` **当前事实** | 「检查点 A/B 完成」 | `progress: 2/3` |

- `goal-tree.md` frontmatter：`updated: 2026-09-21`，`version: 0.60.0`（该 version/日期在 `a246dfb5` 已存在；`e3b90f58` 未再 bump version，但 `updated` 已是当日，可接受）。
- Root `00-meta` frontmatter 仍 `updated: 2026-09-06` / `version: 0.9.0`（波次行已改、头未刷）—— recommended，不是 0/3 与 2/3 冲突。
- **未再构成 required 的残留**：`02-execution.md`「事实边界」仍写「尚未实施…`progress: 0/3`」。治理主面（00-meta / 树 / 表 / Root 台账 / 当前事实）已是 2/3。见 F-004。

## Findings（本轮）

| ID | 级别 | 严重度 | 状态 | 说明 |
|----|------|--------|------|------|
| A-002 F-001 | required | high | **fixed** | `shift` 后 `%1` 转发已在真实 cmd 路径验证。 |
| A-002 F-002 | required | high | **fixed** | dev：DB_NAME → YAML `db.name` → `schema_ui`；test：PG_TEST_DB；与 `config.Load` 一致。 |
| A-002 F-003 | required | med | **fixed** | 治理主面 A/B=completed、2/3 一致。 |
| F-004 | recommended | low | open | `02-execution.md`「事实边界」仍是立项快照（「尚未实施」+ `0/3`），与「当前事实」2/3 并存。建议改成历史快照或删除。 |
| F-005 | recommended | low | open | A-002 意见未落盘：`03-audit.md` 仅索引 A-001；`03-audit/` 无 A-002 文件。本意见按用户指令为 A-003。编排器应补登 A-002 或在索引中声明跳号。 |
| F-006 | recommended | low | open | `pgsetup.Server.Database` 注释仍写 default `postgres`，与 dev 的 `schema_ui` fallback 矛盾。GOAL-047 `00-meta` 检查点 B 仍写「Go 64/64」，E-002/本轮为 65/65。Root `00-meta` `updated` 未随 W35 行刷新。 |

无新 required。I-047-001 / I-047-002 均为 non-blocking 且已 verified；本 scope 无到期 required 信息门禁。共享资料目录为 `none`。

## 必改项汇总

- **开放 required：0**
- 关闭 A-002 F-001/F-002/F-003 为 `fixed` 不需要用户裁决。
- F-004/F-005/F-006 为 recommended，不阻断进入 C。

## 七条待复审事项（`03-audit.md` 索引）

| # | 事项 | 结论 |
|--:|------|------|
| 1 | 空实例自举、维护连接不污染业务路径 | **仍成立**。`MaintenanceCandidates` = 目标库 → `postgres` → `template1`（单测钉死）。业务 `openPostgres` 仍只 Ping 配置库，不 `CREATE DATABASE`。本轮未再 wipe 常驻实例；空实例路径沿用 E-002 + 单测。 |
| 2 | 快捷方式幂等、不破坏既有库；drop 独立 | **仍成立**。对已存在的 `schema_ui_dev`/`schema_ui_test` 报 `exists`。`dbsetup` 无 drop。`e2e-pgset drop` 对存在/缺失分别为 `dropped`/`absent`。 |
| 3 | 一次建齐 dev/test，库名来自配置而非写死 | **仍成立（且 F-002 已补 YAML 路径）**。默认命令仍是 dev+test；本轮用 `--dev-only`/`--test-only` 证明角色开关有效。dev 名：DB_NAME → YAML → `schema_ui`（与 committed `db.name` 相同）；test 名：PG_TEST_DB。不是写死 `schema_ui_dev`。 |
| 4 | 启动 hint 不泄密、只对 3D000、不改分类 | **仍成立**。`missingDatabaseHint` 仅 `pgconn.PgError.Code=="3D000"`；`%w` 保留原错误；文案无 DSN/口令。`28P01` 原样返回（`postgres_open_test.go` 本轮随 store 包绿）。 |
| 5 | 文档与命令一致；SQLite 无需初始化；失败排查 | **仍成立**。根 README / QUICKSTART / `apps/api/README.md` / `dev.cmd` help 均可复制 `dev.cmd init-db` 与 `go run ./cmd/dbsetup`。SQLite 走 `sql.Open("sqlite", dsn)` 建文件，无 `CREATE DATABASE`。e2e `schema_ui_e2e_*` 仍独立。文档未写 YAML fallback（非 required）。 |
| 6 | 可复跑「删库 → init → `dev.cmd start`」 | **证据仍在，本轮未复做破坏性 reset**。E-002 记录 reset 两库 → init-db → start `/readyz 200` + Web 200 → stop。本轮只对 scratch 做 create/drop，未 drop `schema_ui_dev`/`schema_ui_test`。 |
| 7 | 是否越界（迁移/checksum/codec/wire/备份 Port/自动建库） | **未越界**。`e3b90f58` 只改 `pgsetup.go`/`pgsetup_test.go`/`dev.cmd` 与三处进度文档。启动路径仍不自动建库。备份 Port 的 `CREATE DATABASE` 是既有 restore 路径，不在本 diff。 |

## A-001 recommended 复核

| ID | 结论 |
|----|------|
| F-S-001 | **可接受**。`e2e-pgset create` 现为幂等 `Ensure`；drop 仍显式。本轮 `drop`/`absent` 语义符合文档。e2e 使用一次性名字，不依赖「已存在必须失败」。 |
| F-S-002 | **已核对**。`LoadEnvFile` 读 `DB_` + `PG_TEST_`，不覆盖 process env；与 `internal/pgtest` 的 LookupEnv 守卫同构。角色 namespace 分离单测 + 真实 `--dev-only`/`--test-only` 均通过。 |
| F-S-003 | **仍成立**。见待复审 #4。 |
| F-S-004 | **已升级为 A-002 F-001 并在本轮 fixed**。 |

## 与既有意见的异同

- A-001 self：conditional，开放 required=0，四条 recommended 交独立审。本轮把 F-S-004 对应的真实 cmd 缺陷视为已修，其余 recommended 未升格。
- A-002 independent（用户书面 / 提交说明，**未落盘**）：conditional，required=3。本轮关闭这三条。
- 本意见不改 A-001 正文；不把 2/3 当作 C 已完成。

## 结论 + 建议给编排器/用户的下一步

- **开放 required = 0**。
- **A/B 成立**（实现 + 本轮复验证据）。
- **可进入 C**（self A-001 已落盘且 required=0；独立意见以本 A-003 为准，落地后即可做关门审计）。不要用 progress 推导 `done`。
- **无需 P-004**：无意见冲突、无 required residual/overrule、无信息冲突。
- 建议 `/govern`：
  1. 落盘本意见为 `03-audit/A-003-independent-w35-findings-rereview.md`，索引增加 A-003（`source: independent`，`verdict: pass`）。
  2. 将 A-002 F-001/F-002/F-003 标 `fixed`（若 A-002 仍不在仓内，先补登或在 A-003 索引中声明空洞）。
  3. 可选修 F-004/F-006 正文卫生；F-005 为台账手续。
  4. 再跑 C（关门审计）。本轮**不**授权把 GOAL-047 标 `done`。

## 声明

本意见不修改 `status` / `progress` / 方案正文 / goal-tree 状态列。响应、finding 闭合与是否进入 C 由 `/govern` 处理。审计过程写入的临时测试文件与 scratch 库已清理；结束时 `git status` clean，HEAD `e3b90f58`。
