---
id: A-007-independent-pg-client-network-fix
goal_id: GOAL-005-r2-backup-port-and-closeout
doc: audit-entry
status: recorded
parent: GOAL-005-r2-backup-port-and-closeout
created: 2026-09-21
updated: 2026-09-21
version: 0.1.0
source: independent
verdict: pass
open_required: 0
auditor: Grok Build grok-4.6 (high)
---

# A-007 · independent · PR #16 PostgreSQL helper Docker 网络修复复审

- **source**：independent
- **auditor**：Grok Build grok-4.6 (high)
- **日期**：2026-09-21
- **类型** / **scope**：finding-closure + execution-facts · PR #16 中 `pg_dump`/`pg_restore` helper Docker 网络修复（commit `20026ffd`，当前 PR head `33b07b68`）。核对 `ClientDockerNetwork`、容器命令 `--network`、配置/环境传递、GitHub Hosted PostgreSQL CI 的 `DB_CLIENT_DOCKER_NETWORK=host`、测试与文档，以及 A-006 对 F-S-001 的 `fixed` 声明。
- **verdict**：**pass**
- **open required**：**0**
- **完整意见**：本文件

## 范围与区间

- 工作区：`workspace-040-timestamptz-persistence-contract`（`workspace.md`：`id` 匹配；`root_goal` = `GOAL-001-timestamptz-persistence-contract`；`canonical_scope` = `docs/workspaces/workspace-040-timestamptz-persistence-contract/`；`shared_materials_catalog: none`；`primary_plan` = `VP-040-timestamptz-persistence-contract`）。**未读其他工作区。**
- 被审目标：`GOAL-005-r2-backup-port-and-closeout`（已 `done · 3/3`）。本意见只审 post-close helper 网络修复与 A-006 闭合声明，不重开 R2 关门、不改目标状态。
- 对照：A-006（self · `pass` · F-S-001 `fixed`）；E-005；当前工作区树与五件套；AGENTS.md §6b / P-003～P-005；`docs/architecture/principles.md`；`docs/architecture/workspace-protocol.md`。
- 代码基线（本地工作树 = PR #16 head）：`33b07b68d419d5fc28b202aa1469428fb2a9283c`。`20026ffd` 是该 head 的祖先。
- P-005：`I-041-006` 仍为 `verified`（用户 2026-09-20 `D-001`）；`I-041-004` 在本目标仍登记为 R3 前 `non-blocking / deferred`，与 helper 网络缺口无关，不阻断本 scope。共享资料引用 none。
- **未改** `status` / `progress` / 方案正文 / goal-tree / 任何代码或测试。

### 证据分层

| 层 | 性质 | 本轮如何取得 |
|----|------|--------------|
| Hosted CI | 已验证事实 | `gh run view` / job logs：run `35567407333`、`35568956651`、`35569624987`、`35570762993`；PR #16 `headRefOid` = `33b07b68` |
| 仓库代码 | 本地代码判断 | 当前 HEAD 源文件、`20026ffd` diff、workflow YAML |
| 本地单测 | 本轮复跑 | `apps/api`：`TestPgProviderCommandConstruction`、`TestDBPostgresExplodedParams`、`TestRecoveryWiringHandsPgProviderAnAbsoluteWorkDir` 均 PASS |
| 本地真实 PG dump/restore | **未复跑** | 本机未执行 `TestPGRestoreToNewDB` 等四条集成用例；其通过与否以 Hosted `api + postgres` 为准 |

## 成果（有证据）

1. **根因可复核。** 首次 Hosted run `35567407333`（head `92cb6e41`）仅 `api + postgres` 失败，其余 8 个 job `success`。失败 job `106231880624` 日志中四条用例失败，命令均为 `docker run --rm -v … pg_dump`（无 `--network`），错误为 helper 容器内 `localhost:5432` Connection refused：
   - `TestPGRestoreToNewDB`
   - `TestPGLegacyArtifactMustFail`
   - `TestPGMidBatchArtifactMustFail`
   - `TestC3RecoveryAnchorsOnPostgresUpgrade`
2. **修复提交把网络选择收进共享参数构造。** `20026ffd` 为 `PgProvider` 增加 `ClientDockerNetwork`；`clientDockerArgs()` 在非空时追加 `"--network", network`；`Create`（`pg_dump`）与 `Restore`（`pg_restore`）共用该切片。空值不追加网络参数。`ClientVersion()` 仍是 `docker run --rm <image> pg_dump --version`，不连数据库，不需要网络参数。
3. **配置传递闭环。** `config.Config.DBClientDockerNetwork` ← YAML `db.client_docker_network` ← `envOr("DB_CLIENT_DOCKER_NETWORK")`；`composition.recoveryWiring` 写入 `PgProvider.ClientDockerNetwork`。`config.yaml` 为 `${DB_CLIENT_DOCKER_NETWORK:-}`；`.env.example` 给出 `host` 与 Compose 网络名示例。
4. **CI 显式设置 host。** `.github/workflows/r6-basic-matrix.yml`：`api-postgres` job env 与 `browser-e2e` job env 均有 `DB_CLIENT_DOCKER_NETWORK: host`。备份集成测试从 `os.Getenv("DB_CLIENT_DOCKER_NETWORK")` 注入 `PgProvider`。
5. **Hosted 复验绿。**
   | run | head | 结论 | `api + postgres` |
   |-----|------|------|------------------|
   | `35567407333` | `92cb6e41` | failure | failure（见上） |
   | `35568956651` | `20026ffd` | success | success；`internal/backup` `ok` 15.801s，`internal/store` `ok` 108.783s |
   | `35569624987` | `971ebbce` | success | success（A-006 所称最终候选） |
   | `35570762993` | `33b07b68`（当前 PR head） | success | success；PR #16 rollup 9/9 `SUCCESS` |
6. **本轮本地单测绿。** 命令构造、YAML/env 覆盖、composition wiring 三条均 PASS（不依赖真实 PG）。

## 对照成功标准（本 scope）

本目标检查点 A～C 已在 2026-09-20 关门；本 scope 只对照 A-006 的修复主张与 CI 接线。

| 主张 | 状态 | 证据 |
|------|------|------|
| helper 仅在配置非空时追加 `--network` | **成立**（代码） | `provider_pg.go` `clientDockerArgs()`；Create/Restore 共用 |
| CI Hosted PostgreSQL 使用 `host` | **成立**（代码 + CI） | workflow YAML；修复后 `api + postgres` 由红转绿 |
| 默认空配置不改变既有命令 | **成立**（代码） | TrimSpace 为空则不追加；未在 Hosted 上空配置路径做负向集成 |
| A-006 F-S-001 `fixed` | **成立**（可重复核对） | 失败日志无 `--network` + 修复后 Hosted PG job 成功 |
| 未以 skip 替代修复 | **成立**（CI） | 失败 run 记录四条 FAIL 而非 skip；修复后 backup/store 包 `ok` 且耗时与失败 run 同量级 |

## Findings

### F-S-001 · PG helper 容器默认 bridge 无法访问 host-published PostgreSQL（A-006）

- 严重度：high
- 建议：required（原 A-006）
- 状态：**closed · fixed**（本独立复审确认）
- 描述：A-006 的闭合声明与本轮独立核对一致。失败命令不含 `--network`；`20026ffd` 之后 Create/Restore 在配置非空时注入网络；Hosted `api-postgres` 设置 `host`；run `35568956651` / `35569624987` / `35570762993` 的 `api + postgres` 均为 `success`。
- 证据：job `106231880624` 日志；`provider_pg.go` L64–69 / L152–159 / L207–214；`.github/workflows/r6-basic-matrix.yml` L115、L147；PR #16 head `33b07b68` rollup。

### F-I-001 · 命令构造单测未钉住 Restore 的 `--network`，也未钉住空配置省略

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`TestPgProviderCommandConstruction` 把 `ClientDockerNetwork: "host"` 后的 Create 拼接串断言含 `"--network host"`，不调用 Restore，也没有空网络不得出现 `--network` 的负向用例。生产代码 Create/Restore 共用 `clientDockerArgs()`，Hosted `TestPGRestoreToNewDB` 覆盖真实 restore，故不是现行缺陷；回归时单测层会漏掉 Restore 参数漂移或默认行为被改成始终带 `--network`。
- 证据：`apps/api/internal/backup/provider_pg_test.go` L324–387；对照 `provider_pg.go` `clientDockerArgs()`。

### F-I-002 · 嵌入默认 YAML 与 Compose 未登记该配置键

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：`apps/api/configs/config.yaml` 与 `.env.example` 已说明 `DB_CLIENT_DOCKER_NETWORK`。嵌入 `internal/config/config.default.yaml` 的 `db:` 段没有 `client_docker_network`；仓库根 `compose.yaml` 也没有该环境变量。裸二进制走嵌入默认时仍可通过进程环境覆盖（`envOr`），CI 与测试也走环境变量。缺口是文档/默认文件对称性，不是 Hosted PG 门禁缺口。
- 证据：`config.default.yaml` L73–95；`compose.yaml` 无 `DB_CLIENT_DOCKER_NETWORK`；`config.go` L742 `envOr("DB_CLIENT_DOCKER_NETWORK", …)`。

## 必改项汇总

无。本 scope 开放 required = 0。

## 与既有意见的异同

- **与 A-006**：同意 F-S-001 `fixed` 与 verdict `pass`。A-006 把两组 PostgreSQL browser E2E 成功当作 helper 网络修复的旁证；独立核对显示 **同一失败 run `35567407333` 上两组 PG E2E 已经是 `success`**，因此 E2E 不能证明 dump/restore 网络已通。真正转绿的门禁是 `api + postgres`。A-006 所称最终候选 `971ebbce` / run `35569624987` 属实；当前 PR head 已前移到 `33b07b68`，其 run `35570762993` 仍为 9/9 `success`。
- **与 A-002 / A-004**：原 R2 关门审计在本机真实 PG 上跑过 dump/restore，未覆盖 GitHub service-container + helper 默认 bridge 拓扑。本缺口是关门后 Hosted 矩阵才暴露的环境接线问题，不回溯否定 A-004 对当时 scope 的 `pass`。
- **信息项**：`I-041-006` 仍 verified；`I-041-004` 不在本 scope 门禁内。

## 结论 + 建议给编排器/用户的下一步

本 scope 判定 **pass**。A-006 的 F-S-001 闭合证据充分、可重复核对；当前 PR head `33b07b68` 仍包含 `20026ffd`，Hosted `api + postgres` 为 `success`。两条 recommended 不阻断本修复的独立确认。

建议 `/govern` 响应本意见：将 F-S-001 的独立确认记入响应节；F-I-001 / F-I-002 可排入后续文档/单测卫生，无需重开 GOAL-005 状态。

## 声明

本意见 `source: independent`，不修改目标 `status` / `progress` / goal-tree / 决策或执行正文 / 代码或测试。响应由 `/govern` 处理。
