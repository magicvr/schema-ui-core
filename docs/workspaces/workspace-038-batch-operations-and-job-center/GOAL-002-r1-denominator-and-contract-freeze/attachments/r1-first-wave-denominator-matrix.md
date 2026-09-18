---
title: R1 冻结矩阵 · 首波批量操作分母（I-038-003 / GOAL-002 C3）
status: frozen
created: 2026-09-19
updated: 2026-09-19
parent: GOAL-002-r1-denominator-and-contract-freeze
version: 1.0.0
frozen_by: D-001-r1-contract-and-denominator-freeze
---

# R1 冻结矩阵 · 首波批量操作分母

> **性质**：GOAL-002 C3 的可机器核对交付物，冻结依据 `01-decision/D-001-r1-contract-and-denominator-freeze.md` §3。
> **证据基线**：`attachments/R1-recon-I-038-003-batch-operation-inventory.md`（v0.3.0）+ 编排器独立复核（Root `E-002` F-5、F-6、F-13、F-15）。
> **用户 P-004 裁决（2026-09-19）**：首波异步 = **仅「新建批量导出所选」**。

## 1. 分母口径修正（冻结）

| 口径 | 值 | 证据 |
|------|-----|------|
| ❌ 原口径「已有批量 UI 的生产页面」 | **0 个页面** | 全量扫描 `apps/api/modules/**/schema/*.json`（27 个文件），`batchMapping`/`requiresSelection` 仅命中 dev 范例 `dev/examples/schema/admin-list-batch.json:75-81` |
| ✅ 后端批量路由分母 | **4 个资源 / 6 条路由** | users、roles、dict-types、dict-entries、scheduled-tasks（`modules/*/provider.go`；工厂 `resources.go:302-309`） |
| ✅ **首波分母（冻结）** | **1 条新建操作** | §2 |

**关键澄清**：Root/VP-038 文本中「批量操作已在 users/roles/data-dictionary/scheduled-tasks 落地」指**后端路由 + 协议能力**，**不是页面 UI**。两者必须区分——首波要兑现的是「真实批量操作」的**端到端**路径。

## 2. 首波纳入（冻结 = 1 条）

| # | 操作 | 性质 | 形态 | 权限 | 结果形态 | 状态 |
|---|------|------|------|------|---------|------|
| W-1 | **批量导出所选** | **新建** | 列表页批量入口 → 异步 Job（202 + jobId）→ 进度 → 终态可下载 CSV | 待 R2/R3 冻结（读权限 + 导出权限口径） | Job `result` 承载导出产物 + `ResultExpiresAt` | **首波唯一纳入** |

**为什么是新建而不是改造既有**：生产页面零批量动作（§1）；既有后端 `batch-delete` 的页面未使用；且 `batch-delete` 改异步属 BREAKING（§3）。新建一条「所选 → 异步导出 → 可下载结果」是**唯一**同时满足「真实批量操作」「长耗时」「产出物可回看」三条件的路径，直接对应 VP-038 判据 3 与判据 4。

## 3. 保持同步（冻结）

| # | 操作 | 位置 | 冻结处置 | 理由 |
|---|------|------|---------|------|
| S-1 | `users` batch-delete | `resources.go:308` + `users.go:277` | **保持同步**（原子） | 改异步 = BREAKING；VP-038 显式非目标（`VP-038:51`）；`users_batch_test.go:119-184` 断言 409 `LAST_ADMIN` 且零删除 |
| S-2 | `roles` batch-delete | `resources.go:308` + `roles.go:201` | **保持同步**（原子） | 同上（系统角色/占用守卫） |
| S-3 | `dict-types` / `dict-entries` batch-delete | `dictionary.go:284/304` | **保持同步**（顺序回退，**观察项**） | 行数小、交互式；**非原子**是独立缺陷，应补 `DeleteBatch` 原子实现，而非异步化 |
| S-4 | `scheduled-tasks` batch-delete | `scheduledtasks.go:344` | **保持同步**（顺序回退，**观察项**） | 同上；任务定义表行数天然很小 |
| S-5 | 同步 `batch-delete` 契约本身 | ADR-0022 / `render.tsx:765` | **逐字冻结** | D-001 §1.1；改协议面即触碰 pinned 工件 |

> **诚实标注**：S-1/S-2 的「保持同步」按 VP-038 显式非目标记为**派生结论**（D-001 §3.3）——本轮 P-004 提问中该条为确认项，用户未勾选。

## 4. 不进首波（冻结排除清单）

保持**现状**；仍为**后续波次候选**，不构成本 VP 承诺。

| # | 操作 | 位置 | 当前 | 为何不进首波 |
|---|------|------|------|-------------|
| X-1 | 回收站 `purge-all` | `recyclebin.go:131` | 同步 | 唯一无界 `DELETE`（`repository.go:230`）、不可逆；且是「全量」而非「所选批量」，不契合首波形态 |
| X-2 | CSV 导入 | `import.go:45` | 同步 | 单文件上传，非选择集批量；逐行 no-rollback 语义需专门设计 |
| X-3 | 操作日志导出 | `operations_export.go:19-24` | 同步 | 与 data-transfer 导出同形态同上限；首波只做「所选批量」这一条形态 |
| X-4 | `settings` 重置 | `settings.go:41` | 同步 | 品牌资源实际数量小，收益有限 |
| X-5 | `scheduled-tasks` 手动触发 | `scheduledtasks.go:446` | 同步（204） | 改异步 = BREAKING（测试与前端依赖 204） |
| X-6 | notifications `read-all` | `notifications.go:52` | 同步 | 单条 UPDATE + 每用户 500 行硬上限 |
| X-7 | wallet 代金券批量生成 | `wallet.go:484` | 同步 | `count > 1000` 直接拒绝，行数有界 |
| X-8 | wallet reconcile | `wallet.go:357` | **已异步** | 已是先例/模板，非候选 |
| X-9 | 单目标操作（enable/disable/unlock、MFA reset、invites、角色分配） | `users_state.go:43-45` 等 | 同步 | 批量变体**不存在**；属新功能 |
| X-10 | filelibrary 上传/下载/删除/列表 | `filelibrary.go:106/181/218/255` | 同步 | 单文件；但列表无界扫描（`:106-125`）与 `quotaReached` O(files) 是独立性能观察项 |
| X-11 | 4 个后台周期任务（cron 30s / 日志保留 1h / Job 轮询 10s / Telegram 租约 1s） | 见 `r1-job-kind-scope-matrix.md` §5 | 后台 | 非用户触发；首波不迁移进 `jobs` |

## 5. Breaking-change 冻结标记

| 操作 | 标记 | 依据 |
|------|------|------|
| users / roles `batch-delete` → 异步 | **BREAKING** | 前端 `onSuccess.behavior="reload"` 依赖响应即终态；原子回滚语义无法同步返回；协议 fixture 钉住批量构造 |
| `scheduled-tasks POST /{id}/run` → 异步 | **BREAKING** | 现 204；`scheduledtasks_test.go:54,133` 与前端 `scheduled-tasks.json:289-303` 依赖 |
| 导出 / 导入 / 操作日志导出 / `purge-all` → 异步 | **改旧端点即 BREAKING；新增端点则纯增量** | 导出为裸 CSV 附件（`export.go:210-214`）、导入为 `{applied,failed,total,errors,fieldErrors}` envelope（`import.go:74-92`）、`purge-all` 返回 `{"purged": n}` |
| 首波 W-1（新建批量导出所选） | **纯增量** | 新增页面批量入口 + 新端点，不改任何既有路由 |

## 6. 首波回归面（冻结，不得破坏）

| 组 | 文件 | 钉住内容 |
|----|------|---------|
| 同步 batch-delete 原子语义 | `users_batch_test.go:14,72,119,193`；`users_repository_test.go:14,145,178`；`roles_repository_test.go:13` | 200 + `{"deleted":n}`；400 空选/非法 key；409 `LAST_ADMIN` 零删除；批级守卫 |
| 自作用域批量 | `resources_test.go:365,434-444` | self scope 仅删本人行 |
| 顺序回退快照 | `recyclebin_test.go:319` | 回退路径须为**每个** id 记录快照 |
| provider 路由清单（**逐字**） | `users/provider_test.go:65`；`roles/provider_test.go:66`；`datadictionary/provider_test.go:71,74`；`scheduledtasks/provider_test.go:71` | 路由集合含 `POST .../batch-delete`——改路由即失败 |
| 前端批量端到端 | `representative-pages.integration.test.tsx:562-568` | 单次 POST + 归一化 keys body + 成功反馈 + 清选 |
| 协议守卫 | `capability-declaration.guard.test.ts:57`；`request-construction.cases.json:1029-1192` | `actions.batch.request` marker = `/batch-delete/`；11 个 batch fixture |
| wallet Job 面（模板回归） | `wallet_test.go:257-315,435-452,662-670,680-708`；`jobs_test.go:153-155` | 4 条路由权限门；终态/结果码；18 个错误码入 catalog；actor 隔离不可放宽 |
| 导出/导入同步语义 | `data_transfer_test.go:50,86,100,110,122,166,254,266,290,305,326,340,351,361,398` | 同步 CSV/envelope、权限门、10000/2MiB 上限、逐行 no-rollback |
| recycle-bin 同步语义 | `recyclebin_test.go:81,365` | 单行 restore/purge；`purge-all` 同步 `{"purged":2}` |

## 7. 承接 `V-F126`

`V-F126`（`open · recommended`，由 `I-038-003` 承接）的意图是「首波批量操作分母必须明确」。本矩阵以「首波 1 条 + 明确排除清单 X-1～X-11 + 保持同步 S-1～S-5」给出明确分母，**承接动作已完成**。`V-F126` 的闭合登记属愿景层（`/vision`），不在本目标台账内自行改判。
