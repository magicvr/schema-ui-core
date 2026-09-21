---
id: A-002-r1-freeze-independent
doc: audit-entry
parent: GOAL-002-r1-denominator-and-contract-freeze
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# A-002 · R1 冻结独立交叉审计（GOAL-002 C1～C3）

## A-002 · R1 分母与契约冻结独立复审（2026-09-19）

- **source**：independent
- **auditor**：grok-build（grok-4.6 · reasoning high · `/audit`）
- **类型** / **scope**：`stage` · GOAL-002 R1 冻结（C1～C3）全量复审——两项冻结矩阵 vs 代码；`I-038-001`～`003` 的 `verified` 关闭合法性；用户 P-004 裁决落盘与未选方案；派生结论 vs 用户裁决；未定项 O-1～O-3；边界（`apps/**` / pinned 协议工件）
- **verdict**：**conditional**（3 required · medium；4 recommended · low/medium）
- **完整意见**：本文件（未超长，无需附件）

### 范围与区间

- 被审目标：`docs/workspaces/workspace-038-batch-operations-and-job-center/GOAL-002-r1-denominator-and-contract-freeze/`
- 工作区校验：`workspace.md` 的 `id` = `workspace-038-batch-operations-and-job-center`、`root_goal` = `GOAL-001-batch-operations-and-job-center`、`canonical_scope` = 本区路径、`shared_materials_catalog: none`、`vision_role: delivery`、`plan_refs` / `primary_plan` = `VP-038-batch-operations-and-job-center`。绑定一致；无共享资料引用；未跨区读取。
- 审计区间：2026-09-19（R1 立项 → 只读侦察 → 用户 P-004 裁决 → C1～C3 冻结落盘 → A-001 self）。**不含** C4 投影、R2/R3 实施。
- 本条只写意见，不改 `status` / 检查点 / 派生 `progress` / 方案正文 / `goal-tree.md`。

### 范围核对（P-005 / 工作区）

| 项 | 核对结果 |
|----|---------|
| `I-038-001` | GOAL-002 台账 `verified`（2026-09-19）。证据：侦察报告 + `r1-job-kind-scope-matrix.md` + D-001 §2 + 用户 P-004（管理作用域 + `jobs.read`）。代码复验支持种类/作用域/读面缺口/索引/运行时门控主主张。 |
| `I-038-002` | GOAL-002 台账 `verified`（用户裁决方案 B）。证据：侦察报告 + D-001 §1 + 未选方案 §4。方案 B ⇒ 协议 pin 零改动是可由代码与 pin 规则推出的派生约束，不是把假设写成已验证。 |
| `I-038-003` | GOAL-002 台账 `verified`（用户裁决首波 = 仅新建批量导出所选）。首波 1 条与排除项主体成立；**现有** `GET /api/export/{resource}` 未进入 S/W/X 逐项口径（见 F-002）。 |
| Root / VP 台账 | Root `00-meta.md` 与 VP-038 的 `I-038-001`～`003` 仍为 `open`。属 C4 投影 / `/vision` 同步，**不**构成「GOAL-002 把假设写成 verified」。 |
| 共享资料 | `shared_materials_catalog: none`；本 scope 未把任何共享资料当关闭证据。 |

### 成果（有证据）

| # | 成果 | 证据 |
|---|------|------|
| 1 | 用户 P-004 三点裁决如实落盘，含未选方案 | `01-decision.md` P-004 表；`01-decision/D-001-…` §1.1 / §2.1 / §3.1 / §4 |
| 2 | 派生结论与用户裁决分离 | D-001 §3.3；`02-execution/E-002-…` §5；对照 VP-038 显式非目标（`docs/vision/plans/VP-038-batch-operations-and-job-center.md` 首波范围表「把已交付的同步 `batch-delete` 改成 breaking 异步语义」） |
| 3 | 未定项 O-1～O-3 未冒充已裁决 | D-001 §1.3；`01-decision.md` 未定项段 |
| 4 | 生产 Job 种类分母 = 1；HTTP 读面无跨 actor 路径 | 生产 `RegisterWithTerminalHook` 仅 `modules/wallet/jobs.go:39`；`GetForActor` SQL `id+kind+actor_id`（`internal/jobs/repository.go:66-74`）；wallet HTTP 一律经 `JobService.Get/Cancel/Retry` → `GetForActor` |
| 5 | 三索引 sqlite/postgres 两侧一致；跨 actor `ORDER BY updated_at DESC` 不被覆盖 | `modules/jobs/migration/migration.go:45-47` 与 `:84-86`；`idx_jobs_actor` 以 `actor_id` 打头 |
| 6 | `jobRuntime.enabled` 仅 `admin.wallet` 置 true；`Start`/`Stop` 有 no-op 守卫 | `internal/composition/composition.go:185-197`、`:582-588`、`:1101` |
| 7 | `core.jobs` 两处冻结断言如文档 | `internal/store/migrate_test.go:693` checksum `55e1d3f8…`；`modules/jobs/migration/migration_test.go:15` `Version == 42` / `Name == "async_jobs"` |
| 8 | 生产页面 schema 批量 UI 分母 = 0 | 全量 `apps/api/modules/**/schema/*.json`：`batchMapping`/`requiresSelection` 仅 `dev/examples/schema/admin-list-batch.json:75-81` |
| 9 | 同步 `batch-delete` 形状与原子实现者 | 工厂 `resources.go:302-309,824-980`；`maxResourceBodyBytes = 4<<10`（`:37`）；权限 = 资源写权限；成功 `200 {"deleted": n}`；`DeleteBatch` 仅 `users.go:277`、`roles.go:201` |
| 10 | Breaking 标记可核对 | `users_batch_test.go:119-184`（409 `LAST_ADMIN` 零删除）；`internal/handler/scheduledtasks_test.go:54,133`（204）；`scheduled-tasks.json:289-303`（`runTask` 行动作，无 202 处理） |
| 11 | 两次 checkpoint 仅本区文档，`apps/**` 与 pinned 工件零改动 | `e6258dc6`（立项+侦察）、`d09a206f`（C1～C3 冻结）；`git diff-tree` 无 `apps/**`、`docs/schemas/**`、`apps/web/src/protocol/upstream/**` |
| 12 | 方案 B 下「零改动 pinned 工件」成立 | D-001 §1.1；ADR-0022 成功路径 `render.tsx:765` 只判 `response.ok`（202 属 ok）随后 `:1244` `reloadList()`，异步不得走 `batchMapping` |

### 对照检查点（C1～C3）

| 检查点 | 状态 | 独立复核 |
|--------|------|---------|
| C1 分母与作用域矩阵 | **大体达成，有计数/台账缺口** | 种类/作用域/缺口/索引/门控与代码一致；矩阵可 grep。goal-tree 未随 `progress: 3/4` 同步（F-003）。 |
| C2 契约形态冻结 | **达成** | 方案 B + 未选 A/C + K-1～K-7 均可指回代码；O-1～O-3 未写成已裁决。 |
| C3 首波分母冻结 | **大体达成，有分母计数误差与漏项** | 首波 = 1 条新建「批量导出所选」与用户裁决一致；「4 个资源 / 6 条路由」不成立（F-001）；现有 data-transfer 导出未逐项落口径（F-002）。 |
| C4 审计与投影 | 进行中 | 本条为 independent 腿；Root 投影不在本条范围。 |

### 代码复验摘要（不采信文档）

| 主张 | 复验 |
|------|------|
| 生产 Job kind = 1 | **成立**。生产注册点仅 `wallet.jobs.go:39`；测试 `panic.kind` 不计入。 |
| `GetForActor` 三重限定；无跨 actor HTTP 读路径 | **成立**。SQL `WHERE id=? AND kind=? AND actor_id=?`。`Repository.Get(id)` 存在但仅 runner/测试使用（`runner.go:384`），无 HTTP 暴露。`RequestCancel`/`Retry` 经 `getForActorTx`（`id+actor_id`，**不校验 kind**；kind 在 `wallet/jobs.go:96`）——矩阵已诚实标注。 |
| 三索引与「跨 actor ORDER BY updated_at DESC 不被覆盖」 | **成立**。 |
| R-1 运行时门控 | **成立**。不含 `admin.wallet` 时 runner 不启动；R2 须为 `admin.jobs` 显式 `enabled`。 |
| 生产页面零批量动作 | **成立**。`schema/*.json` 共 35 个（其中 `dev/examples` 8 个）；生产 27 个零命中。矩阵写「扫描 27 个」与 glob 全量 35 混用，结论仍对。 |
| 后端批量路由「4 / 6」 | **不成立**。见 F-001。实际 **4 个模块 / 5 条 `POST …/batch-delete` / 5 个非只读 Resource ID**。 |
| 仅 users/roles 原子 `DeleteBatch` | **成立**。 |
| X-1～X-11 主体 | **成立**（purge-all 无界 DELETE `:230`；导入 `:45`；操作日志导出 `:19-24`；settings reset `:41`；tasks run 204 `:446`；notifications `read-all` `:52` + 每用户 500 上限；代金券 `count > 1000` 拒绝；reconcile 已异步；单目标 enable/disable/unlock；filelibrary 列表无界扫描 `:106-125`；4 个 `NewTicker` 周期循环）。遗漏见 F-002。 |
| users/roles `batch-delete` 与 tasks `run` 改异步 = BREAKING | **成立**。 |

### Findings

#### F-001 · 冻结矩阵「4 个资源 / 6 条路由」与代码不一致

- 严重度：med
- 建议：**required**
- 描述：`r1-first-wave-denominator-matrix.md` §1 把后端批量路由分母冻成 **「4 个资源 / 6 条路由」**，证据列却写出 5 个名字（users、roles、dict-types、dict-entries、scheduled-tasks）。独立复验：
  - 非只读 `ResourceRoutes` 实挂 `POST {path}/batch-delete` 的 Resource ID = **5**（users / roles / dict-types / dict-entries / scheduled-tasks）。
  - provider / `kernel/profile.go` 声明的 `POST …/batch-delete` = **5 条**（users、roles、types、entries、scheduled-tasks）。
  - 模块数 = **4**（users / roles / data-dictionary / scheduled-tasks）时路由仍是 5，不是 6。
  - 侦察报告 §1.2 第 6 行「其他所有非只读资源」当前为空集（其余 `ResourceRoutes` 均为 `ReadOnly: true`：operations、monitoring-errors、files、task-runs）。
  该数字从侦察「4 处 / 6 条」原样写入 **frozen** 矩阵，破坏「可机器核对」承诺。列出的名字是对的，冻结的计数是错的。
- 证据：`attachments/r1-first-wave-denominator-matrix.md` §1；`resources.go:302-309`；`modules/users/provider.go:53`；`roles/provider.go:46`；`datadictionary/provider.go:53,56`；`scheduledtasks/provider.go:59`；`kernel/profile.go:166,167,189,193`；Root `attachments/R1-recon-I-038-003-…md` §0.2 / §1.2。
- 关联：`I-038-003`
- 状态：open

#### F-002 · 现有 `GET /api/export/{resource}` 未进入 S/W/X 逐项口径

- 严重度：med
- 建议：**required**
- 描述：`I-038-003` 要求对现有批量/长操作逐项给出「保持同步 / 改异步 / 不进首波」。侦察 §2.1 已把 data-transfer **同步 CSV 导出**（`GET /api/export/{resource}`，内存拼装、10000 行上限、`WriteTimeout=10s`）列为与导入并列的长操作。冻结矩阵：
  - 首波 W-1 是**新建**「批量导出所选」（不同操作）；
  - X-2 覆盖导入、X-3 覆盖操作日志导出；
  - Breaking 表笼统写「导出 / 导入」；
  - **没有** S-* 或 X-* 行给既有 `GET /api/export/{resource}`。
  未列项虽可被「保持现状」兜底，但该端点是与 W-1 最近的既有导出面，漏项会使「逐项冻结」不可机器核对。不改变用户「首波仅 1 条新建」的裁决，但须补一行明确处置（建议：**保持同步 / 不进首波**，与 X-2/X-3 同级）。
- 证据：`attachments/r1-first-wave-denominator-matrix.md` §2～§4；Root `attachments/R1-recon-I-038-003-…md` §2.1（`export.go:25,44,191-214`）；`internal/config/config.go:464`。
- 关联：`I-038-003`
- 状态：open

#### F-003 · `00-meta` 已标 `progress: 3/4`，`goal-tree.md` 仍为 `0/4`；A-001 把「已同步」写成成果

- 严重度：med
- 建议：**required**
- 描述：C1～C3 勾选后 GOAL-002 `00-meta.md` frontmatter `progress: 3/4`，但工作区 `goal-tree.md` 树与状态表仍写 `active · 0/4`。冻结提交 `d09a206f` 未改 `goal-tree.md`。AGENTS §7：改 progress 必须同步 goal-tree，否则视为任务未完成。A-001 成果表第 8 行「goal-tree 树/表/纲领路线图与事实同步」与磁盘事实不符。另：`00-meta.md` 正文仍写「构成 `progress: 0/4` 的派生来源」，与 frontmatter 自相矛盾。本条不改 tree；由 `/govern` 响应时修正。
- 证据：`GOAL-002/00-meta.md` frontmatter `progress: 3/4` 与正文「0/4」；`goal-tree.md` 树/表 `0/4`；`03-audit/A-001-r1-freeze-self.md` 成果 #8；`git show --stat d09a206f`。
- 状态：open

#### F-004 · 侦察报告仍为 `status: draft`，与「已冻结」权威易混淆

- 严重度：low
- 建议：**recommended**
- 描述：同意 A-001 F-001。三份侦察报告 frontmatter 均为 `draft`，决策与矩阵才是冻结权威。不阻断放行；建议在报告头或交叉链接上标明「只读证据，非决策」。
- 证据：Root `attachments/R1-recon-I-038-00{1,2,3}-*.md` `status: draft`；本目标两份矩阵 `status: frozen`。
- 状态：open（recommended，不阻断）

#### F-005 · `V-F126` 闭合登记未在本目标台账写成交接项

- 严重度：low
- 建议：**recommended**
- 描述：同意 A-001 F-002。D-001 §3.5 与首波矩阵 §7 声明承接动作已完成、闭合属 `/vision`，但未登记待办交接。R5 关门前可能无人回看。
- 证据：D-001 §3.5；`r1-first-wave-denominator-matrix.md` §7。
- 状态：open（recommended，不阻断）

#### F-006 · 冻结执行台账索引与证据路径不完整

- 严重度：low
- 建议：**recommended**
- 描述：冻结事实文件 `02-execution/E-002-c1-c3-freeze.md` 存在，但 `02-execution.md` 索引仍只有 E-001。E-002 §7 指向不存在的 `E-003`；实际 checkpoint 为 `d09a206f`（已核路径）。GOAL-002 `00-meta.md` / D-001 把侦察报告写成 `attachments/R1-recon-*.md`（相对本目标则不存在），真实位置是 Root `GOAL-001-…/attachments/`。E-002 §6 写「本轮写入全部为 GOAL-002/」，但 `d09a206f` 同时改了 Root 侦察报告 v0.3.0——仍在本区、未碰 `apps/**`。
- 证据：`02-execution.md` 索引；`02-execution/E-002-c1-c3-freeze.md` §6～§7；D-001 顶部证据路径；`git show --name-only d09a206f`。
- 状态：open（recommended，不阻断）

### 必改项汇总（required）

1. **F-001**：把 C3 矩阵后端批量路由分母改成与代码一致的可核对数字（建议：**4 个模块 / 5 条 `POST …/batch-delete` / 5 个非只读 Resource ID**），并同步侦察报告若仍被引用为权威计数。
2. **F-002**：为既有 `GET /api/export/{resource}` 补 S-* 或 X-* 一行（建议：保持同步 / 不进首波），使 `I-038-003` 的逐项口径闭合。
3. **F-003**：同步 `goal-tree.md` 树与表到 GOAL-002 `3/4`，并修正 `00-meta.md` 正文仍写 `0/4` 的句子。A-001 成果 #8 在响应中更正。

未发现 high 级 required。`I-038-001`/`002` 的 verified 关闭合法。`I-038-003` 的用户首波裁决合法，但矩阵计数与导出漏项使「可机器核对」不完整，故 **不可无条件放行 C4 投影**。

### 与既有意见的异同（A-001 self）

| 点 | A-001 self | 本条 independent |
|----|------------|------------------|
| verdict | pass（0 required，2 recommended） | **conditional**（3 required · medium，4 recommended） |
| P-004 落盘 / 未选方案 / O-1～O-3 / 派生 vs 裁决 | 通过 | **同意** |
| 矩阵 vs 代码主主张 | 抽查通过 | **同意**主主张；**不同意**「4/6 路由」与「goal-tree 已同步」 |
| 生产零批量 UI、原子 DeleteBatch、无跨 actor HTTP 读面、R-1 门控、冻结断言 | 通过 | **同意**（独立复验） |
| `apps/**` / pin 零改动 | 通过 | **同意**（`e6258dc6`、`d09a206f`） |
| A-001 F-001 draft 侦察报告 | recommended | **同意** → 本条 F-004 |
| A-001 F-002 V-F126 交接 | recommended | **同意** → 本条 F-005 |
| 后端路由计数 | 未发现 | **新增 F-001 required** |
| data-transfer 导出逐项口径 | 未发现 | **新增 F-002 required** |
| goal-tree / progress 同步 | 误列为成果 | **新增 F-003 required** |
| 执行索引 / 证据路径 | 未写 | **新增 F-006 recommended** |

无 verdict 相反的冲突需要 P-004 裁两审；本条把 A-001 的 pass 降为独立侧 conditional，由编排器响应 required 项。

### 结论 + 建议给编排器/用户的下一步

C1～C3 的**决策内容**（方案 B、首波仅新建批量导出所选、管理作用域 + `jobs.read`、batch-delete 保持同步为派生、O-1～O-3 移交 R2）**如实、可核对、边界干净**。信息项关闭不是把假设写成 verified：有侦察证据 + 书面 P-004。

独立侧不能给 pass：冻结矩阵作为「可机器核对」交付物有一处错误计数、一处长操作漏行，且 progress 未按 AGENTS §7 同步 goal-tree。**verdict = conditional**。

**建议 `/govern` 下一步**：响应本条 F-001～F-003（修正矩阵计数、补导出口径、同步 goal-tree 与 00-meta 正文），并顺手处理 F-004～F-006；闭合 required 后再投影 Root R1。不要在本条 required 开放时把 Root R1 标完成。

### 声明

本意见不修改 status/progress/检查点/goal-tree/方案正文；响应由 `/govern` 处理。
