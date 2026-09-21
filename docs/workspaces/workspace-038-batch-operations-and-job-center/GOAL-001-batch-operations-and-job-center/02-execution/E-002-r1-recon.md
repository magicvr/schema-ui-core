---
id: E-002-r1-recon
doc: execution-entry
parent: GOAL-001-batch-operations-and-job-center
status: recorded
created: 2026-09-19
updated: 2026-09-19
version: 0.1.0
---

# E-002 · R1 只读侦察（`I-038-001`～`003` 证据收集）

## 事实（2026-09-19）

`I-038-001`～`003` 为 R1 冻结前的 required 信息门禁。本轮以**只读侦察**方式收集证据，不改动 `apps/**`，产出三份侦察报告：

| 报告 | 承接信息项 | 路径 |
|------|-----------|------|
| Job 种类与作用域 | `I-038-001` | `attachments/R1-recon-I-038-001-job-kinds-and-scopes.md` |
| 批量异步契约 | `I-038-002` | `attachments/R1-recon-I-038-002-batch-async-contract.md` |
| 批量操作与长操作清单 | `I-038-003` | `attachments/R1-recon-I-038-003-batch-operation-inventory.md` |

### 已核验事实（编排器独立复核，非仅采信报告）

| # | 事实 | 证据 |
|---|------|------|
| F-1 | Job 运行时是**通用**的，位于 `apps/api/internal/jobs`（六态 + lease + progress + cancel + retry + result 过期）；`apps/api/modules/jobs/` 只有迁移（`core.jobs`，`Provider.Register` 空实现，刻意不进运行期 Profile） | `internal/jobs/model.go:18-25,38-59`；`modules/jobs/migration/provider.go:20`；`migration.go:1-2` |
| F-2 | 全仓**唯一**已注册 Job kind = `wallet.reconcile`（`apps/api/modules/wallet/jobs.go:19`），由 `admin.wallet` 模块注册（`jobs.go:39`），只经 `/api/wallet/jobs/*` 访问 | `jobs.go:19,39,45,81,95,109`；`internal/handler/wallet.go:357-451` |
| F-3 | Job 现有读面仅 **actor 作用域**：`GetForActor(ctx, id, kind, actorID)` 三重限定；`RequestCancel`/`Retry` 亦限 `actor_id` | `internal/jobs/repository.go:66-74,124-149,220-238` |
| F-4 | Repository **没有任何列表查询**：无 `List`/`ListByActor`/`Count`；只有 `ListRunnable`（worker 内部取待执行任务，`limit` 必填）。缺 kind/status/actor/时间过滤、分页、总数、排序 | `repository.go:292-316`；索引 `idx_jobs_actor(actor_id, kind, updated_at DESC)`、`idx_jobs_expiry(status, expires_at)`、`idx_jobs_runnable(status, cancel_requested, lease_expires_at, created_at)`（`modules/jobs/migration/migration.go:45-47`） |
| F-5 | 同步 `batch-delete` 契约：`POST {path}/batch-delete`，body `{"ids":[...]}`（4 KiB 上限），权限 = 资源 `{id}.write`，成功 `200 {"deleted": n}`；仅 `users` / `roles` 实现原子 `DeleteBatch`，其余资源走**顺序回退**（非原子、首败即停、已删不回滚） | `internal/handler/resources.go:308,824-982`；`users.go:277`；`roles.go:201`；`resources.go:927-972` |
| F-6 | **生产页面 schema 零批量动作**：`batchMapping` / `requiresSelection` 的 schema 级声明只存在于 dev 范例 `dev/examples/schema/admin-list-batch.json` 与 pinned 协议工件/fixture。即「已交付的批量 UI 分母 = 0 个生产页面」 | 全量扫描 `apps/api/modules/**/schema/*.json`（27 个文件）仅 `admin-list-batch.json` 命中；`admin-list-batch.json:75-81` |
| F-7 | 前端批量执行链已实现且被测试钉住：`buildBatchRequest` → `runBatchRequest`（只判 `response.ok`）→ 成功 `reloadList()` + 清空选择；`$selection.keys` 仅允许出现在 body | `protocol/conformance/request-construction.ts:642,683-684`；`renderer/render.tsx:685,765,1244`；`representative-pages.integration.test.tsx:562-568` |
| F-8 | 异步先例形状：`POST /api/wallet/reconcile` → **202** + job 投影；`GET /jobs/{id}`、`POST /jobs/{id}/cancel`、`POST /jobs/{id}/retry`、`GET /jobs/{id}/result`（409 `JOB_RESULT_NOT_READY` / 410 `JOB_RESULT_EXPIRED` / succeeded 返回附件字节） | `internal/handler/wallet.go:357,384,387,400,413,426-451`；投影 `walletJobToMap` `:989-1006` |
| F-9 | **web 侧零 job 消费**：`apps/web/src` 对 `wallet/jobs` / `jobId` / `resultUrl` 零命中；wallet 页面只用同步 `request` action，202 响应体被忽略。真正可复用的轮询形态是 `monitoring-auto-refresh.tsx`（自定义组件 + `crud.refreshList`，5/10/30s） | `modules/wallet/schema/wallet.json:96-103`；`components/monitoring-auto-refresh.tsx:17-22,36-38` |
| F-10 | 上游 pin 规则：新增**本地**端点/页面**不需要**改任何 pinned 工件；改动 `docs/schemas/**` 或 `upstream/*.cases.json` 字节 = 上游协议变更，须重 pin 并更新 `stage3-fixtures.test.ts` 哈希 | `protocol/upstream/provenance-v2.9.json:2-4,6-127`；`docs/architecture/module-contribution-playbook.md:84,113`；`docs/vision/protocol-inventory-v2.7.0.md:213` |
| F-11 | 新 capability 字符串有守卫陷阱：写进 `meta.requiredCapabilities` 但未在 `MARKERS` 表登记标记 → `capability-declaration.guard.test.ts` 直接失败（`no usage marker defined for capability`） | `protocol/capability-declaration.guard.test.ts:38-58,110-114` |
| F-12 | 同步路径存在硬时间上界：HTTP `WriteTimeout = 10s`（默认，可配） | `internal/config/config.go:463-465` |
| F-13 | 长操作候选（非批量）：导出 CSV（`GET /api/export/{resource}`，同步、内存拼装非流式、上限 10000 行）、导入 CSV（`POST /api/import/{resource}`，同步、2 MiB、逐行 no-rollback 报告）、`recycle-bin` `purge-all`（同步、无界单条 `DELETE FROM recycle_items WHERE restored_at IS NULL`） | `internal/handler/export.go:25,42,191-199`；`import.go:30,43`；`recyclebin.go:128-146`；`modules/recyclebin/store/repository.go:227-242` |
| F-14 | `data-transfer` 模块只声明 3 条路由（`GET /api/export/{resource}`、`POST /api/import/{resource}`、`GET /api/import/{resource}/template`） | `modules/datatransfer/provider.go:44` |
| F-15 | 通用资源工厂给**每个非只读资源**挂 `batch-delete`；provider 契约面声明 6 条：users / roles / dict-types / dict-entries / scheduled-tasks | `resources.go:302-309`；`modules/users/provider.go:53`；`roles/provider.go:46`；`datadictionary/provider.go:53,56`；`scheduledtasks/provider.go:59` |
| F-16 | `admin.jobs` 目标模块尚**不存在**：`apps/api/modules/jobs/` 当前只有 `migration/`；现有 `ModuleID = "core.jobs"` 是迁移专用且不在任何 Profile 默认集 | `modules/jobs/migration/provider.go:12-20`；`modules/compiled/persistence.go:14,44`；`kernel/profile.go`（`mvp`/`admin`/`demo` 集合均无 jobs） |

### 本轮**未**做的事（边界）

- **未**改动 `apps/**`（`admin.jobs` 模块建立属 R2/R3 实现范围）。
- **未**关闭任何信息项：`I-038-001`～`003` 在本次侦察后仍为 `open`——证据已收集，但「冻结口径」需 R1 决策落盘并（对涉及方案选型的部分）经用户 P-004 裁决。
- **未**执行任何审计。
- **未**创建纲领阶段子目标：R1 子目标 `GOAL-002-r1-denominator-and-contract-freeze` 随后按 P-001 立项（见 `../../GOAL-002-r1-denominator-and-contract-freeze/02-execution/E-001-r1-establishment-and-recon.md`）。

## 待办（移交 R1 决策）

`I-038-002` 与 `I-038-003` 的**契约形态**与**首波分母**属方案选型，按 P-004 须经用户裁决；`I-038-001` 的作用域口径（管理作用域 vs actor 作用域）与读面清单同属 R1 冻结内容。
