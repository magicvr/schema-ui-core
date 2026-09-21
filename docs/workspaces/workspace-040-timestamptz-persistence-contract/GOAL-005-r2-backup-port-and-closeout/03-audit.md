---
id: GOAL-005-r2-backup-port-and-closeout
doc: audit
status: active
parent: null
created: 2026-09-20
updated: 2026-09-21
version: 0.6.0
---

# 审计台账 · GOAL-005-r2-backup-port-and-closeout（R2 M4）

> 本文件是唯一正式审计台账索引：`self` 与 `independent` **共用** `A-NNN` 序列。
> 每条意见正文在 `03-audit/A-NNN-<slug>.md`；本文件登记条目头（`source` / 日期 / scope / `verdict`）。
> 独立审计默认只写意见，不修改 `status` / `progress` / 方案正文；响应归编排器。

## 信息就绪核对（本 scope）

| 核对项 | 状态 | 备注 |
|--------|------|------|
| I-041-006 | **verified**（`D-001` accepted，用户 2026-09-20） | A-002 F-I-007 指出 `00-meta` 未同步，已修正 |
| I-041-004 | deferred（R3 前，non-blocking） | 不阻断 R2；本轮实测 server 15.4 + 容器客户端 15.19（同主版本） |
| 共享资料引用 | none | `workspace.md` `shared_materials_catalog: none` |

## 意见索引

| A-ID | source | 日期 | scope | verdict | 摘要 | 文件 |
|------|--------|------|-------|---------|------|------|
| A-001 | self | 2026-09-20 | R2 关门自审：判据 M1–M4 逐项 + C3 §4.2/§4.3 + 偏差 8 项 | conditional | 交 independent 复审判定 | `03-audit/A-001-self-r2-closeout.md` |
| A-002 | independent | 2026-09-20 | R2 关门审计（grok-build grok-4.6 · high）：M1–M4 + C3 §3.1/§4.2/§4.3/§5/§5.1/§6 + 分母一致性 + 未声明回归 | **conditional** | M1–M3 与 `D-021` residual 可核对；开放 required = 3（PG 类 C 命名、composition 未注入、PG 样本别名）；另有 recommended F-I-004～F-I-007 | `03-audit/A-002-independent-r2-closeout.md` |
| A-003 | self（编排器响应） | 2026-09-20 | 响应 A-002 全部意见 | **pass** | 3 required 全部 `fixed`；F-I-004/005/006 fixed；F-I-007（元数据同步）fixed；待 independent 复审确认 | `03-audit/A-003-response-to-closeout-audit.md` |
| A-004 | independent | 2026-09-20 | 复审 A-002 三条 required 的闭合（grok-build grok-4.6 · high） | **pass** | **open required = 0**：F-I-001/002/003 均 `fixed`（真实 PG 复跑未 skip）；F-I-007 partial；新增 recommended F-I-008/009/010（不阻断检查点 C） | `03-audit/A-004-independent-a002-required-closure.md` |
| A-005 | self（编排器响应） | 2026-09-20 | 响应 A-004 + 检查点 C 关门 | **pass** | F-I-007/010 索引同步；F-I-008（PG 样本正向覆盖，实测 ms:2/sec:90）；F-I-009（marker 诊断并列保留 + PG 读失败不阻断）；**检查点 C 完成 → R2 阶段完成** | `03-audit/A-005-response-to-reaudit-and-checkpoint-c-closure.md` |
| A-006 | self | 2026-09-21 | PR #16 post-close PostgreSQL CI 缺口修复与发布候选复验 | **pass** | 首次 PG helper `localhost` 访问失败已修复；最终候选 CI 9/9 通过 | `03-audit/A-006-post-close-pg-client-network-self.md` |

## R2 关门审计范围

- 判据：Root `D-016` §4 M4（`D-021` residual 三项完成且经 independent 复审 → R2 self + independent 关门审计通过）。
- 必查：`GOAL-003`/`GOAL-004` 的 required 闭合状态、`D-021` residual 收口记录（`GOAL-002/03-audit/A-048`）、`I-041-006` 的用户裁决、Backup Port 的 `<recovery-artifact>` 校验与错误分类证据。

## A-006 · self · PR #16 post-close PostgreSQL CI 缺口修复（2026-09-21）

- **source**：self
- **auditor**：current-session
- **scope**：PR #16 中 `PgProvider` helper-container 网络选择、配置接线、PostgreSQL CI 验证；不复审 VP/Root 关门状态
- **verdict**：**pass**（本 scope open required = 0）
- **完整意见**：[`03-audit/A-006-post-close-pg-client-network-self.md`](A-006-post-close-pg-client-network-self.md)

### 结论摘要

- 首次 Hosted CI 的 PG 失败有明确根因：默认 bridge 下 helper 容器中的 `localhost` 不指向 runner 的 PostgreSQL service。
- `DB_CLIENT_DOCKER_NETWORK` 可选配置向 `pg_dump` 与 `pg_restore` 注入 Docker 网络；本地配置与 composition 测试通过。
- Hosted `api + postgres` job 实际运行成功；最终候选 `971ebbce` 的 `r6-basic-matrix` 9/9 `success`。
- 此为已关闭目标的后续证据补录，不改目标状态或 progress。

## A-002 · independent · R2 关门（2026-09-20）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **verdict**：conditional（open required = 3）
- **完整意见**：[`03-audit/A-002-independent-r2-closeout.md`](A-002-independent-r2-closeout.md)

### 结论摘要

- M1–M3 与 `D-021` residual 三项 **可视为已满足**（`GOAL-002/03-audit/A-048` `fixed`；`GOAL-003`/`GOAL-004` required 已闭合）。
- Port 表面、SQLite harness、PG 快乐路径、§5.1 的 A 类反向断言 **真实 PG 复跑绿**；90 列分母与 inventory v0.3 **逐行一致**（90 行 / 44 表）。
- 开放 required = **3**：F-I-001（PG 类 C 固定文件名挡住「可续跑」）、F-I-002（默认 composition 不注入 Port，§4.3 在生产入口是空操作）、F-I-003（PG `SampleVerified` 是类型检查别名）。
- recommended：F-I-004（marker 只看存在不解析）、F-I-005（`verifyIntegrityPG` 写死 87）、F-I-006（PG 缺类 C 反向断言）、F-I-007（`00-meta` 未同步 `I-041-006`）。

## A-003 · self（编排器响应）· 2026-09-20

- **verdict**：**pass**（A-002 的 3 条 required 与 4 条 recommended 全部处置）
- **完整响应**：[`03-audit/A-003-response-to-closeout-audit.md`](A-003-response-to-closeout-audit.md)

| finding | 处置 |
|---------|------|
| `F-I-001`（required，PG 类 C 命名） | **fixed**：`<db>.pre-v%04d-<ts(ms)>-<hex4>.dump`；真实 PG 升级链（A=1/C=1/B=1）与 composition 全绿 |
| `F-I-002`（required，composition 未注入） | **fixed**：`composition.recoveryWiring` 按方言构造 `backup.Service` + `PgProvider`，并传 `ArtifactDir`（sqlite：DB 同级 `recovery/`；PG：用户缓存目录按 DSN 派生） |
| `F-I-003`（required，PG 样本） | **fixed**：新增 `verifyPostgresSamples`（sentinel 列全 NULL + 抽样微秒精度），service 的 PG 分支调用它 |
| `F-I-004`（recommended） | **fixed**：`probeRecoveryState` 解析 marker；不可解析 = 未满足并把原因写入 `Detail` |
| `F-I-005`（recommended） | **fixed（数据驱动）**：`conversionCompletionVersion(catalog)` 以冻结台账最后一个转换描述符为门禁版本，裁短历史（v1–v72 / v1–v86）不触发全量要求 |
| `F-I-006`（recommended） | **fixed**：`TestPGMidBatchArtifactMustFail`（真实 PG 混合形状 → `TimeContractMismatch`） |
| `F-I-007`（recommended） | **fixed**：`00-meta` 的 `I-041-006` 同步为 verified（本文件信息就绪表同步） |

## A-004 · independent · A-002 required 闭合复审（2026-09-20）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **verdict**：**pass**（open required = 0）
- **完整意见**：[`03-audit/A-004-independent-a002-required-closure.md`](A-004-independent-a002-required-closure.md)

| A-002 finding | 本审判定 |
|---------------|----------|
| `F-I-001`（required/high） | **fixed** |
| `F-I-002`（required/med） | **fixed** |
| `F-I-003`（required/med） | **fixed** |
| `F-I-004` / `F-I-005` / `F-I-006` | 主体 **fixed**（F-I-004 Detail/PG 阻断残余 → 新 F-I-009） |
| `F-I-007` | **partial**（`00-meta` 已同步；`01-decision.md` 索引仍 open → 新 F-I-010） |

本轮新 recommended（不阻断 C）：`F-I-008`（PG 样本未正向锁定「各至少一例」）、`F-I-009`、`F-I-010`。

## 关门审计状态（2026-09-20）

`self` A-001（`conditional`）→ `independent` A-002（`conditional`，required = 3）→ 编排器 A-003（**pass**）→ `independent` A-004（**pass**，open required = **0**）。检查点 C 的 independent 门禁在 finding-closure scope 内可成立；**本意见不修改** GOAL-005 `status` / `progress` / Root R2 状态，响应归 `/govern`。
