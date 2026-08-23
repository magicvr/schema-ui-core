---
title: A-002 · W24 修复成果独立复审（independent）
source: independent
status: recorded
created: 2026-08-23
updated: 2026-08-23
parent: GOAL-035-w24-e2e-dual-dialect-matrix
version: 0.1.0
auditor: deepseek-harness independent audit（本地独立复跑复现）
scope: 复审 W24 修复成果 — S1 设计（D-001）→ S2 实施（E-002）→ S3 回归与 F-1 整改（E-003）→ 关门自审（A-001）；finding-closure / execution-facts
verdict: pass
---

# A-002 · W24 修复成果独立复审（2026-08-23，independent）

## 范围与区间

用户指令：`/audit 复审工作区10目标35的修复成果`。
复审对象：GOAL-035-w24-e2e-dual-dialect-matrix（W24）从设计冻结到关门的一整条修复成果链，重点是「强制 sqlite = 绕过」被修正为「双方言各测一次」的可信度：
方言契约（D-001/D1）、`cmd/e2e-pgset` provisioning（D2）、fail-fast 校验（D3）、CI 矩阵（D4）、npm 脚本与 README（D5）、回归证据与配置双载整改（E-003 F-1）、A-001 自审结论。

## 成果（有证据）

| 项 | 证据路径 | 核验方式 |
|----|----------|----------|
| 提交落地：工作树干净，W24 全部改动在 `018a4d7`（2026-08-23，dev 已推送 origin/dev） | `git status --short` 空；`git log`/`git show --stat 018a4d7`（25 文件 +709/−22） | 本复审直接核验 |
| 方言契约：默认 sqlite / pg 显式 opt-in，非法值抛错；webServer env **恒写** `DB_DIALECT` | `apps/web/playwright.config.ts:58-88,121-130` | 读码 |
| 「`.env` 无法再静默改道」的结构依据：config 装载 env-file 时跳过已存在的 process env | `apps/api/internal/config/config.go:868`（`if _, exists := os.LookupEnv(k); exists { continue }`） | 读码 |
| 本机威胁环境仍存在（`apps/api/configs/.env` 含 `DB_DIALECT` 键）→ 契约防御被实测 | 本机 `apps/api/configs/.env`（值已脱敏） | 检查键存在 |
| scratch-pg 工具：create/drop(含 WITH FORCE 回退)/verify/list；凭据 env → `configs/.env` 同源 | `apps/api/cmd/e2e-pgset/main.go`（212 行） | 读码 + `go build`/`go vet` 通过（BUILD-OK/VET-OK） |
| verify 语义真实：`to_regclass('public.schema_migrations')` 与 API postgres 存储的 ledger 表名一致 | `apps/api/internal/store/identity.go`（postgresLedgerDDL）、`internal/store/postgres.go:71`、`internal/store/postgres_test.go:160`（同款 `to_regclass` 检查） | 读码 |
| globalSetup 契约校验（sqlite=文件出现 / pg=verify，60s fail-fast 诊断）；globalTeardown drop 可见（stdio inherit） | `apps/web/e2e/global-setup.ts`、`global-teardown.ts` | 读码 |
| npm 脚本 + README 双方言说明 | `apps/web/package.json:13`、`apps/web/README.md:21-44`、`apps/web/scripts/run-e2e-postgres.mjs` | 读码 |
| CI `profile × dialect [sqlite, postgres]` 矩阵、service 容器、GITHUB_ENV 注入 | `.github/workflows/r6-basic-matrix.yml:118-174` | 读码 |
| I-001 先证实验：10 tests → 1 skipped → 9 passed（1.8m） | `apps/web/e2e-w23-pg-experiment.log`（UTF-16 原始日志；**被 `*.log` gitignore，未入库**） | 本复审读取（见 F-002） |
| **sqlite 腿独立复跑**：`npm run test:e2e` → 1 skipped + **9 passed**（2.1m），exit 0；且在本机 `.env` 陷阱仍存在时成立 | 本次复审现场运行（后台 job pwsh-62 输出） | 现场复跑 |
| **postgres 腿独立复跑**：`npm run test:e2e:postgres` → 1 skipped + **9 passed**（1.8m），exit 0；`dropped schema_ui_e2e_mt5a0ht1li16q8` 可见；事后 `e2e-pgset list` = **0 遗留**（单建单删，F-1 配置双载守卫生效） | 本次复审现场运行（后台 job pwsh-63 输出） | 现场复跑 |
| 台账终态：goal-tree 树/注解/表三处一致；workspace.md W24 节；conformance-claim buildId 随构建刷新 | `goal-tree.md:55,58,143`、`workspace.md` W24 节、commit 内 conformance-claim.json | 读文件 + 提交 diff |

## 对照成功标准

| 标准 | 核验 | 结论 |
|------|------|------|
| C1 设计冻结（D-001） | 契约显式化 + 未选方案记录 + 用户书面复审承接链完整（GOAL-034 → GOAL-035） | pass |
| C2 实施 | 六项产物全部存在且结构与声明一致；`go build`/`go vet` 通过 | pass |
| C3 回归 | sqlite 与 postgres 两腿**独立复跑均 9/9（1 跳过）**，provisioning 生命周期建→验→跑→删闭环、遗留 0；vitest/go/tsc 声明未重复执行（W24 无产品代码改动，提交足迹佐证） | pass |
| C4 关门 | A-001 自审落盘；台账同步；I-001 closed 证据（实验日志）存在 | pass |

## Findings

| F-ID | 级别 | 内容 | 证据 | 处置建议 |
|------|------|------|------|----------|
| F-001 | recommended（低） | `apps/web/playwright.config.ts` L42/L51 注释中「→」被写成乱码「閳?」（U+95B3+U+003F），由 W24 引入（对比 `1a81bac` 版无此字符，W24 diff 新增）；纯注释、不影响执行 | `git diff 1a81bac 018a4d7 -- apps/web/playwright.config.ts` | 顺手修正为 UTF-8「→」 |
| F-002 | recommended（低） | I-001 先证实验日志 `apps/web/e2e-w23-pg-experiment.log` 被 `apps/web/.gitignore:16 (*.log)` 忽略，实验证据仅存本地；同样，E-003 最终回归也无持久化日志。CI 矩阵合入 main 后将成为持久证据，但当前 I-001/回归的原始证据不在仓库内 | `git check-ignore -v` 输出；本机日志 10→9+1 内容 | 将实验摘要（库名/结果/日期）落盘进 00-meta 或 attachments；CI 首跑后回填运行记录 |
| F-003 | recommended（中低） | CI `browser-e2e` 矩阵**尚无运行记录**：W24 提交在 `origin/dev`（已推送），而工作流仅 `push/PR to main` 触发（r6-basic-matrix.yml:3-7），gh 查询无 2026-08-23 之后运行。矩阵语法/行为仅经静态阅读 + 本地双腿复跑推定 | `gh run list --workflow=r6-basic-matrix.yml`（最新运行 2026-08-21）；`git branch -a --contains 018a4d7` | dev 合入 main 或 PR 后核验矩阵双腿全绿并回填证据；如首次运行失败属遗留 pre-main 状态，不影响本次修复成果判定 |

## 必改项汇总

无 required/必改 findings。三项均为 recommended，不阻断关门，也不需要 reopen。

## 与既有意见的异同

- 与 A-001（self，pass）**结论一致**：修复成果成立。差异在于证据强度——A-001 依据自报回归；本复审对 C3 关键主张做了**现场独立复跑**（sqlite + postgres 双腿 9/9、遗留 0），A-001 F-1「配置双载」修复亦被复跑行为间接证实（单建单删）。
- 新增 A-001 未记录的三点：注释乱码（F-001）、实验证据未入库（F-002）、CI 矩阵首跑证据缺失（F-003）。

## 结论 + 建议给编排器/用户的下一步

**verdict: pass**。
用户复审要求「收尾层 e2e 双方言各测一次」的修复成果**成立且有可重复证据**：「强制 sqlite = 绕过」已从结构上消除（process-env 契约 + fail-fast 校验双保险，本机 `.env` 陷阱仍在时两腿复跑均绿），postgres 成为一等公民（provisioning 生命周期闭环），CI 矩阵设计就位等待 main 首跑。

建议 `/govern` 响应：
1. 将三项 recommended（F-001 乱码修正、F-002 证据固化）列入顺手整改或记录接受；
2. F-003 在 dev 合入 main 的 PR/推送后回填 CI 矩阵运行记录（双方言腿全绿截图/URL）；
3. 无需 reopen GOAL-035；无 required 残留，关门维持。

## 声明

本意见不修改 status/progress/方案正文；响应与任何状态变更由 `/govern` 处理。

---

## 编排器响应（2026-08-23，source: orchestrator · 非独立意见）

用户指令：`/govern 响应 GOAL-035 A-002，把 recommended 顺手改了。然后提交。`

| F-ID | 级别 | 响应处置 | 证据/位置 |
|------|------|----------|-----------|
| F-001 | recommended | **fixed** | `apps/web/playwright.config.ts` L42/L51 乱码已改回 UTF-8「—」/「→」；非 ASCII 残留扫描仅剩预期字符；`--list` 正常 |
| F-002 | recommended | **fixed** | 证据摘要落盘 `attachments/I-001-evidence.md`（先证实验 / W24 双腿终验 / 本复审现场复跑 / CI 回填节）；00-meta I-001 备注已链 |
| F-003 | recommended | **fixed（2026-08-23 合入后回填）** | PR #5 合并 `cdb2308`；`browser-e2e` 矩阵首跑 https://github.com/magicvr/schema-ui-core/actions/runs/32617287887 9/9 SUCCESS（含 mvp/admin × sqlite/postgres 四腿）；证据见 `attachments/I-001-evidence.md` §4 |

响应事实全量见 `02-execution/E-004-a002-response.md`。A-002 三项 recommended 全部闭合；GOAL-035 维持 `done`，无 reopen。