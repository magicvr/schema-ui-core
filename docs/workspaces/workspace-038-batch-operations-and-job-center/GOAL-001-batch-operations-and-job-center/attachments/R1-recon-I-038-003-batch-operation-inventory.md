---
title: R1 侦察 · 批量操作与长操作清单（I-038-003）
status: draft
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.1.0
---

# R1 侦察 · 批量操作与长操作清单（I-038-003）

> 只读侦察报告。范围：`apps/api`（Go）+ `apps/web`（React）。所有结论附 `file:line`。
> 未取到证据的项一律标注「未核实」，不臆测。生成于 2026-09-19。

## 0. 结论摘要

1. 全仓库实现 `DeleteBatch` 的实体**只有 2 个**：`usersEntity`、`rolesEntity`。
2. 通用工厂把 `POST {path}/batch-delete` 挂在**每一个非只读资源**上，因此**注册面 4 处 / 6 条路由**，
   但其中 4 条（dict-types、dict-entries、scheduled-tasks 两处 + 其他）**没有** `DeleteBatch`，
   走的是**顺序删除回退路径**（非原子、首个失败即停、已删行不回滚）。
3. **生产页面 schema 目前没有任何一个声明批量动作**——`batchMapping` / `requiresSelection` 只出现在
   dev 范例 `admin-list-batch.json` 与协议 fixture 中。即：**批量 UI 的已交付分母为 0 个生产页面**，
   后端能力先于前端页面落地。这是首波筛选时最需要先对齐的事实。
4. 唯一已交付的异步 Job 先例是 `admin.wallet` 的 reconcile（202 + jobId + cancel/retry/result）。
   Job 运行时是 `apps/api/internal/jobs`（六态 + lease + progress + cancel），**通用**且可复用。
5. 导出/导入是**同步**的，且导出**在内存里拼整个 CSV**（非流式），上限 10000 行 / 2 MiB。

---

## 1. 批量面清单（BATCH SURFACE INVENTORY）

### 1.1 契约与路由注册

| 项 | 位置 | 说明 |
|----|------|------|
| `BatchDeleter` 接口 | `apps/api/internal/handler/resources.go:81-83` | `DeleteBatch(ids []string, user account.User) (int, error)` |
| 路由注册（唯一一处） | `apps/api/internal/handler/resources.go:308` | `add("POST", res.Path+"/batch-delete", a.Middleware(h.batchDelete()))` |
| 注册条件 | `apps/api/internal/handler/resources.go:302` | `if !res.ReadOnly` → 只读资源**不挂**批量路由 |
| handler 实现 | `apps/api/internal/handler/resources.go:824-982` | 见 §1.3 |
| 权限门 | `apps/api/internal/handler/resources.go:826` | `requirePermission(w, r, h.writePerm)`；`writePerm` 默认 `res.ID + ".write"`（`resources.go:279-282`） |
| 请求体 | `apps/api/internal/handler/resources.go:830-833` | `{"ids":[...]}`，`http.MaxBytesReader(w, r.Body, maxResourceBodyBytes)` = 4 KiB |
| 成功响应 | `apps/api/internal/handler/resources.go:980` | `200 {"deleted": n}` |

路由出现在 provider 声明（契约面）中的位置：

- `apps/api/modules/users/provider.go:53` — `POST /api/users/batch-delete`
- `apps/api/modules/roles/provider.go:46` — `POST /api/roles/batch-delete`
- `apps/api/modules/datadictionary/provider.go:53,56` — `POST /api/data-dictionary/types/batch-delete`、`.../entries/batch-delete`
- `apps/api/modules/scheduledtasks/provider.go:59` — `POST /api/scheduled-tasks/batch-delete`
- 冻结 profile 声明：`apps/api/kernel/profile.go:166`（users）、`:167`（roles）、`:189`（data-dictionary，两条）、`:193`（scheduled-tasks）

### 1.2 逐资源清单

| # | resource id | mount path | 实体文件 | 实现 `DeleteBatch`？ | 原子性 | 回收站快照 | 权限门 |
|---|-------------|-----------|---------|---------------------|--------|-----------|--------|
| 1 | `users` | `/api/users` | `apps/api/internal/handler/users.go:277` | ✅ | **单事务原子**（`DeleteUsersBatch`，`apps/api/modules/authsession/users_repository.go:346-415`，`r.withTx` at `:352`） | 有（`Trash` 由组合注入；快照在整批提交后逐条记录，`resources.go:907-926`） | `users.write`（`apps/api/internal/handler/users.go:56` 起的 Resource 定义默认派生） |
| 2 | `roles` | `/api/roles` | `apps/api/internal/handler/roles.go:201` | ✅ | **单事务原子**（`DeleteRolesBatch`，`apps/api/modules/authsession/roles_repository.go:187-221`，`r.withTx` at `:193`） | 有（同上机制） | `roles.write` |
| 3 | `dict-types` | `/api/data-dictionary/types` | `apps/api/internal/handler/dictionary.go:284-303` | ❌ **未实现** | **顺序循环回退**（`resources.go:927-972`）：逐条删除、首个失败即 return、**已删行不回滚** | 有（`Trash: recorder` at `dictionary.go:302`）；走 `TrashTxDeleter` 分支时**每行一个事务**（`resources.go:940-955`） | `dictionary.write`（`dictionary.go:299`） |
| 4 | `dict-entries` | `/api/data-dictionary/entries` | `apps/api/internal/handler/dictionary.go:304-321` | ❌ **未实现** | 同上顺序回退 | 有（`dictionary.go:320`） | `dictionary.write`（`dictionary.go:317`） |
| 5 | `scheduled-tasks` | `/api/scheduled-tasks` | `apps/api/internal/handler/scheduledtasks.go:344-361` | ❌ **未实现** | 同上顺序回退 | 有（`scheduledtasks.go:360`） | `tasks.write`（`scheduledtasks.go:357`） |
| 6 | 其他所有非只读资源（含 `file-library`、`settings`、`notifications` 等） | 各自 path | — | ❌ | 同上顺序回退 | 视是否注入 `Trash` 而定 | 各自 `{id}.write` |

> **注**：`system-monitoring` 为 `ReadOnly: true`（`apps/api/internal/handler/systemmonitoring.go:85`），
> 因此**不挂** `batch-delete` 路由（`resources.go:302`）。

### 1.3 两条路径的行为差异（关键）

`apps/api/internal/handler/resources.go:898-981`：

- **BatchDeleter 路径**（`resources.go:899-926`）：先整批 `Get` 快照（`:907-913`），再 `batch.DeleteBatch(ids, user)`（`:915`），
  失败 → `writeEntityError` 直接返回（`:916`），**整批回滚**；成功后逐条 `Trash.Record`（`:919-926`）。
- **回退路径**（`resources.go:927-972`）：`for _, id := range ids`（`:929`）逐条 `Delete`，
  任何一条失败即 `return`（`:957-960`），**前面已删的行已经提交、不回滚**。
  成功返回 `{"deleted": deleted}`（`:971`）。

自作用域（self scope）预过滤在两条路径之前：`resources.go:884-896`，非本人行被**跳过**（不报 404），全被过滤时返回 `{"deleted":0}`（`:893`）。

### 1.4 users 的批级 last-admin 守卫（异步化时不可丢失的语义）

- 守卫实现：`apps/api/modules/authsession/users_repository.go:379-393`
  （`countAdminUsersExcludingBatch`，`:417` 起），批内含**全部** admin → `ErrLastAdmin`，整批回滚。
- 背景：这是 W3 独立审计 F-001 修复（`docs/workspaces/workspace-009-production-hardening/GOAL-004-w3-security-audit-remediation/03-audit/A-002-w3-independent-cross.md:123-128`）。
- roles 侧守卫：系统角色 `ErrRoleSystem`、被占用 `ErrRoleInUse`（`roles_repository.go:201-210`）。

---

## 2. 其他批量类操作（OTHER BATCH-LIKE OPERATIONS）

### 2.1 `admin.data-transfer`（导出 / 导入）— **同步**

| 操作 | endpoint | 位置 | sync/async | 上限 | 权限 |
|------|----------|------|-----------|------|------|
| 导出 CSV | `GET /api/export/{resource}` | `apps/api/internal/handler/export.go:42-45` | **同步** | `maxExportRows = 10000`（`export.go:25`）；`pageSize > maxExportRows` → 400 `INVALID_EXPORT_LIMIT`（`export.go:130-134`） | `data.export`（`export.go:53`） |
| 导入 CSV | `POST /api/import/{resource}` | `apps/api/internal/handler/import.go:43-47` | **同步** | `maxImportBytes = 2 << 20`（2 MiB，`import.go:30`）；超限 → 413（测试 `data_transfer_test.go:326-338`） | `data.import`（`import.go:60`） |
| 导入模板 | `GET /api/import/{resource}/template` | `apps/api/internal/handler/import.go:49-52` | 同步 | — | `data.import`（`import.go:60`） |

**导出不是流式的**：整个 CSV 在内存 `strings.Builder` 中拼装
（`export.go:191-199`：`var out strings.Builder` → `csv.NewWriter(&out)` → 逐行 `writer.Write`），
先 `ListUsers`/`ListRoles` 全量取回（`export.go:156-159`、`:177-180`），最后一次性写出。
文件头注释自述为「streams the filtered resource list as CSV」（`export.go:1-5`），
**与实现不符**——实现是 buffer-then-write。这是导出异步化的直接动机之一。

**导入为显式 no-rollback 语义**：`import.go:1-6` 文件头自述
「Explicit no-rollback semantics: valid rows commit even when other rows fail；
the response carries the full per-row error report ({applied, failed, total, errors})」。
即导入**不是**原子的，天然适合异步 + 进度 + 逐行错误报告。

导出内容仅支持 `users` / `roles`（`export.go:125-128`，其他 → 404 `RESOURCE_NOT_FOUND`）。

### 2.2 `admin.wallet` reconcile — **异步（唯一先例）**

| 项 | 位置 |
|----|------|
| 提交 reconcile | `POST /api/wallet/reconcile` → `apps/api/internal/handler/wallet.go:357-385`；`jobService.SubmitReconcile(...)` at `:379`；响应 **202 Accepted** at `:384` |
| 查询 Job | `GET /api/wallet/jobs/{id}` — `wallet.go:387-398`，权限 `wallet.read` |
| 取消 | `POST /api/wallet/jobs/{id}/cancel` — `wallet.go:400-411`，权限 `wallet.write` |
| 重试 | `POST /api/wallet/jobs/{id}/retry` — `wallet.go:413-424`，权限 `wallet.write` |
| 读结果 | `GET /api/wallet/jobs/{id}/result` — `wallet.go:426-451`；queued/running → 409 `JOB_RESULT_NOT_READY`（`:438`）；expired → 410（`:440`）；succeeded → `Content-Disposition: attachment; filename="wallet-reconcile-<jobID>.json"`（`:445`） |
| 提交权限 | `wallet.write`（`wallet.go:358`），注释明确：提交/取消/重试都是写操作，只读角色不得触发（`wallet.go:353-356`） |
| 运行历史列表 | `GET /api/wallet/reconcile/runs` — `wallet.go:454-478`，权限 `wallet.read` |

**作用范围**：`body.AccountID` 为空 ⇒ **全账本 reconcile**（`wallet.go:366-370` 注释：
「an empty body (io.EOF) stays valid and means the documented all-accounts reconcile」）。
body 上限 4 KiB（`wallet.go:365`）。空 accountId 走**全表**，这正是它必须异步的原因。

### 2.3 Job 运行时（异步模板）

`apps/api/internal/jobs`（**不是** `apps/api/modules/jobs`——后者只有迁移，见
`apps/api/modules/jobs/migration/migration.go`）：

- 六态：`apps/api/internal/jobs/model.go:18-25` — `queued` / `running` / `succeeded` / `failed` / `cancelled` / `expired`
- `Job` 结构：`model.go:38-59` — 含 `Kind`、`Payload`、**`Progress int`**（`:43`）、`CancelRequested`（`:44`）、
  `Attempt`/`MaxAttempts`、`LeaseOwner`/`LeaseVersion`/`LeaseExpiresAt`（租约）、`Result`、`ErrorCode`/`ErrorMessage`、
  `ActorID`、`CorrelationID`、`ResultExpiresAt`
- `CreateInput`：`model.go:61-69`
- `DefaultMaxAttempts = 3`：`model.go:27`
- 取消语义错误：`ErrNotCancellable`（`model.go:34`）、`ErrLeaseLost`（`:33`）
- 仓储与租约：`apps/api/internal/jobs/repository.go:136`（状态分支）、`:192`、`:364`（租约校验）
- 运行器：`apps/api/internal/jobs/runner.go:165`（running 分支）
- 该运行时**已具备** progress + cancel + retry + result 过期，是首波异步批量操作可直接复用的模板。

### 2.4 `admin.scheduled-tasks` — 手动触发是**同步**的

- `POST /api/scheduled-tasks/{id}/run` — 路由 `apps/api/internal/handler/scheduledtasks.go:446-448`；
  handler `:449-464`：**在请求内联执行** `runner.Execute(*task, time.Now().UTC())`（`:459`），
  失败 → 500 `INTERNAL`（`:460`），成功 → **204 No Content**（`:463`）。
  ⇒ **同步**，且**无返回值、无 runId、无进度**。
- 权限 `tasks.write`（`scheduledtasks.go:450`）。
- `Execute` 实现：`apps/api/modules/scheduledtasks/scheduler.go:133-189`；
  handler 通过 `context.Background()` 调用（`scheduler.go:167`）——**不继承请求 context**，
  因此客户端断开不会取消任务。
- 定时循环：`scheduler.go:65-85`，`tickInterval = 30 * time.Second`（`scheduler.go:24`）；
  单实例 best-effort，错过窗口不回补（`scheduler.go:1-5`）。
- **批量 enable/disable：未发现**。只有通用 `batch-delete`（`scheduledtasks/provider.go:59`），
  以及创建/编辑时的 `enabled` 字段（`scheduledtasks.go:355` JSONFields）。
  ⇒ 「批量启用/停用」目前**不存在**。

### 2.5 `admin.recycle-bin` — 单行 restore/purge + **有 purge-all**

| 操作 | endpoint | 位置 | 单/批 | sync/async | 权限 |
|------|----------|------|-------|-----------|------|
| 列表 | `GET /api/recycle-bin` | `apps/api/internal/handler/recyclebin.go:45-47` | — | 同步 | `recycle.read` |
| 详情 | `GET /api/recycle-bin/{id}` | `recyclebin.go:92-94` | 单行 | 同步 | `recycle.read` |
| 恢复 | `POST /api/recycle-bin/{id}/restore` | `recyclebin.go:109-111` | **单行** | 同步 | `recycle.write` |
| 物理删除 | `DELETE /api/recycle-bin/{id}` | `recyclebin.go:149-151` | **单行** | 同步 | `recycle.write`（`recyclebin.go:153`） |
| **清空回收站** | `POST /api/recycle-bin/purge-all` | `recyclebin.go:128-146` | **全量批量** | **同步** | `recycle.write`（`recyclebin.go:133`） |

- `purge-all` 是**单条 SQL 全表删除**，在一个事务内：
  `apps/api/modules/recyclebin/store/repository.go:227-242` —
  `DELETE FROM recycle_items WHERE restored_at IS NULL`（`:230`），返回 `{"purged": n}`（`recyclebin.go:144`）。
  **无分页、无上限、无进度**。
- 服务层：`Service.PurgeAll()` `apps/api/modules/recyclebin/service.go:132-134`；
  单行 `Purge` `service.go:127-129`。
- **恢复是单行的，没有批量恢复**；恢复内部走事务（`service.go:185` `restoreRowTx`）。
- **未发现**基于保留期的自动清理 sweeper（未核实：仅确认无 `PurgeAll` 的定时调用点；
  全仓库周期任务只找到 scheduler 的 30s tick）。

### 2.6 其他模块的批量/长操作

| 模块 | 操作 | 位置 | endpoint | sync/async | 说明 |
|------|------|------|----------|-----------|------|
| `notifications` | **全部标记已读** | handler `apps/api/internal/handler/notifications.go:162-177`；仓储 `apps/api/modules/authsession/notifications_repository.go:216-232` | `POST /api/notifications/read-all`（`notifications.go:52`） | **同步** | 单条 `UPDATE notifications SET read_at=? WHERE user_id=? AND read_at IS NULL`（`notifications_repository.go:220`），返回 `{"updated": n}`（`notifications.go:175`）。有天然上限：`maxNotificationsPerUser = 500`（`notifications_repository.go:22`）⇒ **行数有界，keep-sync** |
| `wallet` | **代金券批量生成** | handler `apps/api/internal/handler/wallet.go:484-`；服务校验 `apps/api/modules/wallet/voucher/service.go:41` | `POST /api/wallet/vouchers/batches`（`wallet.go:484`） | **同步** | `count <= 0 \|\| count > 1000` 拒绝（`voucher/service.go:41`；handler 侧同判 `wallet.go:506`）⇒ **上限 1000，行数有界**。权限 `wallet.voucher.issue`（`wallet.go:485`）。body 上限 16 KiB（`wallet.go:496`） |
| `users` | 邀请（单条创建/撤销/重发） | `apps/api/internal/handler/invites.go:98-101` | `POST /api/users/invites`、`DELETE /api/users/invites/{id}`、`POST /api/users/invites/{id}/resend` | 同步 | **无批量邀请、无批量重发**（未发现循环） |
| `users` | 启用/停用/解锁 | `apps/api/internal/handler/users_state.go:43-45` | `POST /api/users/{id}/enable\|disable\|unlock` | 同步 | **逐用户单条**，无批量变体 |
| `mfa` | 重置他人 MFA | `apps/api/internal/handler/mfa.go:358` | `POST /api/users/{id}/mfa/reset` | 同步 | **单用户**；权限 `users.mfa-reset`（`apps/api/kernel/profile.go:172`）。**无批量重置** |
| `mfa` | 自服务 enroll/confirm/disable/recovery rotate | `mfa.go:204/256/279/321` | `POST /api/mfa/*` | 同步 | 全部单用户、body 上限 4 KiB |
| `roles` | 权限/菜单授权 | `apps/api/internal/handler/roles.go`（PATCH 单角色） | `PATCH /api/roles/{id}` | 同步 | **无批量授权**（未发现循环） |
| `filelibrary` | 上传登记 | `apps/api/internal/handler/filelibrary.go:255-259` | `POST /api/library/files/upload` | 同步 | body 上限 `maxResourceBodyBytes` 4 KiB（`filelibrary.go:267`）——只登记**已上传**的 fileId，不是上传本身 |
| `filelibrary` | 下载 | `filelibrary.go`（`GET /api/library/files/{id}/download`，见 `kernel/profile.go:187`） | 单文件 | 同步 | **无批量下载 / 无 zip 打包**（未发现） |
| `filelibrary` | 删除 | `DELETE /api/library/files/{id}`（`kernel/profile.go:187`） | 单文件 | 同步 | 无批量变体（除通用 batch-delete，但该资源未声明） |
| `systemmonitoring` | 状态/错误查询 | `apps/api/internal/handler/systemmonitoring.go:85` 起 | `GET /api/system-monitoring/*` | 同步 | `ReadOnly: true` ⇒ **不挂** batch-delete（`resources.go:302`） |
| `datadictionary` | 类型删除是否级联删条目 | **未核实** | — | — | 未在本次侦察中确认；`dictTypeEntity`/`dictEntryEntity` 未实现 `DeleteBatch` |
| `datapermission` | 策略/作用域读写 | `apps/api/internal/handler/datapermission.go:53/95/135/156` | `GET/PATCH /api/data-permission/*` | 同步 | 配置面，非批量 |

### 2.7 循环 / 长操作检索结论

- 面向选择的循环：`apps/api/internal/handler/resources.go:929`（`for _, id := range ids`，顺序删除回退路径）、
  `resources.go:886`（self-scope 逐 id `Get`）、`resources.go:908`（快照预读逐 id `Get`）、
  `resources.go:921`、`:975`（逐条 `OnWrite`）。
- 批量 SQL：`DELETE FROM recycle_items WHERE restored_at IS NULL`
  （`apps/api/modules/recyclebin/store/repository.go:230`）——唯一一条无界 `DELETE`。
- `... WHERE ... IN (...)` 动态占位符：`apps/api/modules/authsession/roles_repository.go:231-237`
  （`PermissionsForRoles`，角色数有界）。
- 事务内逐 id 循环：`users_repository.go:354-410`、`roles_repository.go:194-216`（均**单事务**，见 §1.2）。

---

## 3. Web 侧批量 UI（WEB-SIDE BATCH UI）

### 3.1 关键结论：**生产页面 schema 没有任何批量动作**

全仓库 `batchMapping` / `requiresSelection` / `$selection.keys` 的 schema 级声明只有两处：

| 文件 | 性质 | 内容 |
|------|------|------|
| `apps/api/modules/dev/examples/schema/admin-list-batch.json:75-81` | **dev 范例**（非生产页面） | `requiresSelection: true`、`confirm`、`batchMapping.body.ids = "$selection.keys"`，action url `POST /api/users/batch-delete`（`:20`），`onSuccess.behavior = "reload"`（`:22`） |
| `docs/schemas/component-registry.json:701,707,711` + `apps/web/src/protocol/upstream/*.cases.json` | **协议 schema / fixture** | 声明 `requiresSelection`、`batchMapping`（body 值可为 `$selection.keys`；query 可为 `$selection.count`） |

**生产页面 schema 清单（均无批量动作）**：
`apps/api/modules/users/schema/users.json`（toolbar 仅 `openCreate`/`openInvites`/`exportUsers`/`openImport`，`:429-463`）、
`apps/api/modules/roles/schema/roles.json`、`apps/api/modules/datadictionary/schema/data-dictionary.json`、
`apps/api/modules/scheduledtasks/schema/scheduled-tasks.json`、`apps/api/modules/recyclebin/schema/recycle-bin.json`、
`apps/api/modules/wallet/schema/wallet.json`、`apps/api/modules/filelibrary/schema/file-library.json`。

> **推论（对 I-038-003 重要）**：`batch-delete` 的**后端 + 协议能力已交付**，
> 但**生产页面尚未有任何一处使用它**。因此「批量 UI 的既有分母 = 0 个生产页面」，
> 首波筛选若以「已有批量 UI 的页面」为分母会得到空集；分母应改为
> **「已挂载 batch-delete 路由的资源」= 4 个资源 / 6 条路由**（§1.2），
> 或「具备 table.selection + actions.batch.request 能力的页面」。

### 3.2 渲染器运行时行为（协议能力已实现）

`apps/web/src/renderer/render.tsx`：
- 批量 action 的判定：`render.tsx:331` — `if (url.endsWith("/batch-delete"))`
- `batchMapping` 透传：`render.tsx:180-181`（类型）、`:729-737`（收集）、`:1233`（携带）、
  `:1297` — 确认通过后重建 item：`const item = { actionRef, key: actionKey, batchMapping: pendingConfirm.batchMapping };`
- selection 传入构造器：`render.tsx:738` — `selection: { keys: selection.keys, count: selection.count }`
- 表格侧批量触发判定：`apps/web/src/renderer/schema-table.tsx:1333`
  （`(trigger as Record<string, unknown>).batchMapping !== undefined || ...`）
- 构造器：`apps/web/src/protocol/conformance/request-construction.ts:654-700`（`buildBatchRequest`）
  - `:655-657` 缺 `batchMapping` → `EMPTY_BATCH_MAPPING`
  - `:683` `$selection.keys` **仅允许出现在 body** → 否则 `SELECTION_KEYS_BODY_ONLY`（`:684`）
  - `:734-735` `$selection.count` 为标量（query 可用）
  - `:669-677` `batchMapping.path` 绑定校验（`MISSING/EXTRA_PATH_BINDING`）
- 可搜索性：`apps/web/src/app/searchable.ts:470` — `requiresSelection === true || batchMapping !== undefined` 视为批量候选

> 确认弹窗、成功 reload、清选的具体实现行号**未逐一核实**（`render.tsx` 对应段落未逐行读完）。
> 已确认的事实：`onSuccess.behavior = "reload"`（范例 `admin-list-batch.json:22`）、
> 后端注释「Success returns {"deleted": n} so the client can reload (which clears selection)」
> （`apps/api/internal/handler/resources.go:822-823`）。

### 3.3 i18n 键

| key | zh-CN | en-US | 定义位置 |
|-----|-------|-------|----------|
| `schema.admin-list-batch.toolbar.batchDelete` | 批量删除 | Batch delete | `apps/web/src/i18n/messages/zh-CN.json:477` / `en-US.json:477` |
| `schema.admin-list-batch.confirm.batchDelete` | 确认删除所选用户？ | Delete the selected users? | `zh-CN.json:478` / `en-US.json:478` |

使用位置：`apps/api/modules/dev/examples/schema/admin-list-batch.json:82-83`（`labelKey` / `confirmKey`）。
**仅此 2 个键**，且仅服务于 dev 范例；生产页面无批量 i18n 键。

### 3.4 其他长操作 UI

- 导出/导入 UI：`apps/api/modules/users/schema/users.json:449-462`（`exportUsers` / `openImport` toolbar 项）。
  导出是 `GET` 触发下载；导入是 `actions.upload`（`users.json:14`）上传 CSV 后同步调用。
  **未发现** job/progress 轮询。
- wallet reconcile UI：`apps/api/modules/wallet/schema/wallet.json`。后端已提供 202 + job 轮询面，
  **前端是否有 job 状态轮询未核实**。
- scheduled-tasks「立即运行」按钮：后端 `POST /{id}/run` 同步 204；前端 UI **未核实**。
- recycle-bin 清空按钮：后端 `POST /purge-all` 同步；前端 UI **未核实**。
- **通用异步 Job 中心 UI：本次未发现**（未核实：`apps/web/src` 未做完整 job/poll 检索）。
  VP-038 的立项描述本身即把「通用 Job 管理页」列为 VP-012 的显式非目标
  （`docs/vision/plans/VP-038-batch-operations-and-job-center.md:39` 附近），
  与「当前不存在」一致。

---

## 4. 候选筛选（CANDIDATE SCREENING）

| # | 操作 | 位置 | 当前 | 筛选 | 一句话理由 |
|---|------|------|------|------|-----------|
| 1 | `users` batch-delete | `resources.go:308` + `users.go:277` | 同步 | **keep-sync** | 单事务原子、4 KiB body 上限、批级 last-admin 守卫已交付且有回归测试；异步化即破坏既有合同 |
| 2 | `roles` batch-delete | `resources.go:308` + `roles.go:201` | 同步 | **keep-sync** | 同上（系统角色/占用守卫），原子且行数小 |
| 3 | `dict-types` / `dict-entries` batch-delete | `dictionary.go:284/304` | 同步（顺序回退） | **keep-sync（观察项）** | 行数小、交互式；但**非原子**——若要改进应补 `DeleteBatch` 原子实现，而非异步化 |
| 4 | `scheduled-tasks` batch-delete | `scheduledtasks.go:344` | 同步（顺序回退） | **keep-sync（观察项）** | 同上；任务定义表行数天然很小 |
| 5 | **导出 CSV** | `export.go:42` | 同步 | **async-candidate（最强）** | 已在内存拼整个 CSV（`export.go:191-199`，**非流式**）、上限 10000 行、天然「产出物 + 可下载结果」形态；与 Job 的 `Result` + `ResultExpiresAt`（`model.go:50,58`）高度契合 |
| 6 | **导入 CSV** | `import.go:43` | 同步 | **async-candidate（最强）** | 已是「逐行处理 + 逐行错误报告 + no-rollback」形态（`import.go:1-6`），异步化只需把 per-row 报告搬进 Job result；2 MiB / 行数无上限 |
| 7 | **recycle-bin `purge-all`** | `recyclebin.go:128` | 同步 | **async-candidate** | 唯一一条无界 `DELETE`（`repository.go:230`，无分页/无上限），不可逆，全表扫描；结果只需 `{"purged": n}` |
| 8 | wallet reconcile | `wallet.go:357` | **已异步** | **out-of-scope（已是先例/模板）** | 202 + jobId + cancel/retry/result 已交付，是本波要复用的模板而非候选 |
| 9 | scheduled-tasks `POST /{id}/run` | `scheduledtasks.go:446` | 同步 | **async-candidate（弱）** | 用 `context.Background()` 执行、不随请求取消（`scheduler.go:167`），且无 runId 返回；但 v1 只有 `system.noop` handler（`scheduler.go:46-50`），当前耗时≈0 |
| 10 | notifications `read-all` | `notifications.go:52` | 同步 | **keep-sync** | 单条 UPDATE + 每用户 500 行硬上限（`notifications_repository.go:22`） |
| 11 | wallet 代金券批量生成 | `wallet.go:484` | 同步 | **keep-sync** | `count > 1000` 直接拒绝（`voucher/service.go:41`），行数有界 |
| 12 | recycle-bin 单行 restore/purge | `recyclebin.go:109/149` | 同步 | **out-of-scope** | 单行、交互式，非批量 |
| 13 | users enable/disable/unlock、MFA reset、invites | `users_state.go:43-45`、`mfa.go:358`、`invites.go:98-101` | 同步 | **out-of-scope** | 全部单用户；**批量变体不存在**（属新功能，非本波承接对象） |
| 14 | filelibrary 上传/下载/删除 | `filelibrary.go:255` 等 | 同步 | **out-of-scope** | 单文件；无批量下载/zip（未发现） |

### 4.1 Breaking-change 标记

| 操作 | 是否 breaking | 说明 |
|------|--------------|------|
| users / roles `batch-delete` | **改异步即 BREAKING** | 已交付合同：`200 {"deleted": n}` + 同步原子回滚（`resources.go:980`）。改成 `202 + jobId` 会破坏：① 前端 `onSuccess.behavior = "reload"` 依赖响应即终态；② 原子回滚语义（`users_batch_test.go:119-184` 断言 409 `LAST_ADMIN` 且**零删除**）在异步下不再可同步返回；③ 协议 fixture `request-construction.cases.json` 的 batch 用例。VP-038 已把它写成**显式非目标**（`docs/vision/plans/VP-038-batch-operations-and-job-center.md:51`：「把已交付的同步 `batch-delete` 改成 breaking 异步语义」）。**结论：保持同步，另立异步变体（如新增 `POST {path}/export-jobs` 或 `?async=true`），不改既有路由。** |
| 导出 / 导入 | **非 breaking** | 当前为同步 `GET /api/export/{resource}` / `POST /api/import/{resource}`。若**新增**异步端点（如 `POST /api/export/{resource}/jobs`）并保留旧端点，属纯增量；若直接改旧端点返回 202 则 breaking。 |
| `purge-all` | **非 breaking（需谨慎）** | 返回 `{"purged": n}`（`recyclebin.go:144`）；若改异步需保留同步端点或提供兼容。测试 `TestRecycleBinPurgeAll`（`recyclebin_test.go:365-405`）断言同步返回 `purged=2`。 |
| scheduled-tasks `POST /{id}/run` | **改异步即 BREAKING** | 当前返回 **204 No Content**（`scheduledtasks.go:463`），测试 `scheduledtasks_test.go:54`、`:133` 依赖 204。 |
| notifications `read-all` | 非 breaking（不建议改） | — |

---

## 5. 规模 / 时间上界证据（SIZE/TIME EVIDENCE）

### 5.1 请求体上限

| 常量 / 调用 | 值 | 位置 |
|-------------|-----|------|
| `maxResourceBodyBytes` | `4 << 10` = **4 KiB** | `apps/api/internal/handler/resources.go:37` |
| ↳ create 应用 | — | `resources.go:611` |
| ↳ update 应用 | — | `resources.go:696` |
| ↳ **batch-delete 应用** | — | `resources.go:833` |
| ↳ account_self 应用 | — | `account_self.go:142,240` |
| ↳ email_identity / password_policy / service_credentials / settings / notifications / filelibrary 应用 | — | `email_identity.go:72,94`；`password_policy_settings.go:96`；`service_credentials.go:139`；`settings.go:181`；`notifications.go:220`；`filelibrary.go:267` |
| `maxImportBytes` | `2 << 20` = **2 MiB** | `apps/api/internal/handler/import.go:30`；检查点 `import.go:179`（超限 → 413） |
| `maxUploadBytes` | `8 << 20` = **8 MiB**（单次 multipart 上传） | `apps/api/internal/handler/upload.go:34`；应用 `upload.go:294`、`upload.go:305` |
| `UploadMaxBytesPerUser` 默认 | `256 << 20` = **256 MiB**（每用户配额） | `apps/api/internal/config/config.go:486`（env `UPLOAD_MAX_BYTES_PER_USER`，`config.go:742`） |
| `BrandingMaxBytes` 默认 | `4 << 20` = **4 MiB** | `apps/api/internal/config/config.go:490`（env `BRANDING_MAX_BYTES`，`config.go:743`） |
| wallet reconcile body | `4 << 10` = 4 KiB | `apps/api/internal/handler/wallet.go:365` |
| wallet 代金券批量生成 body | `16 << 10` = 16 KiB | `apps/api/internal/handler/wallet.go:496` |
| auth 登录/改密 body | `4 << 10` | `apps/api/internal/handler/auth.go:94,223,261` |
| Telegram webhook body | `1 << 20` | `apps/api/internal/channel/telegram/webhook.go:142` |

> **batch-delete 的 4 KiB 上限是硬约束**：`{"ids":[...]}` 在 4 KiB 内大约只能装
> **数百个 UUID 级 id**（未实测；按每 id 约 40 字节计 ≈ 100 个）。这是「同步批量天然有界」的直接证据，
> 也是异步批量必须换端点的原因之一。

### 5.2 行数 / 分页上限

| 常量 | 值 | 位置 |
|------|-----|------|
| `maxPageSize` | **100** | `apps/api/internal/handler/resources.go:38`；校验 `resources.go:434-437` |
| `DefaultPageSize` | 20 | `apps/api/internal/handler/resources.go:40` |
| `maxExportRows` | **10000** | `apps/api/internal/handler/export.go:25`；校验 `export.go:130-134`；注释：「full exports beyond this require additional filtering. Kept in sync with the frozen limit.」 |
| ↳ operation-log 导出复用 | — | `apps/api/internal/handler/operations_export.go:29-30` |
| `maxNotificationsPerUser` | **500** | `apps/api/modules/authsession/notifications_repository.go:22`；剪枝 `:103-104` |
| 代金券批量生成上限 | **1000** | `apps/api/modules/wallet/voucher/service.go:41`；handler 侧 `wallet.go:506` |
| `resourceIDRetries` | 3 | `apps/api/internal/handler/resources.go:42` |
| `DefaultMaxAttempts`（Job） | 3 | `apps/api/internal/jobs/model.go:27` |
| `recycle_items` purge-all | **无上限** | `apps/api/modules/recyclebin/store/repository.go:230` |

### 5.3 超时 / 限流

| 项 | 值 | 位置 |
|----|-----|------|
| HTTP `ReadTimeout` 默认 | **5s** | `apps/api/internal/config/config.go:463`；env `HTTP_READ_TIMEOUT` `config.go:694` |
| HTTP `WriteTimeout` 默认 | **10s** | `apps/api/internal/config/config.go:464`；env `config.go:695` |
| HTTP `IdleTimeout` 默认 | 60s | `apps/api/internal/config/config.go:465` |
| `HTTPShutdownTimeout` 默认 | 10s（drain 预算，≤0 fail-closed） | `apps/api/internal/config/config.go:468`、校验 `:709-711` |
| ↳ 服务器装配 | — | `apps/api/server/serve.go:203-206`；`apps/api/server/config.go:125-128` |
| Telegram 出站 HTTP | **10s** 严格预算 | `apps/api/internal/channel/telegram/http_sender.go:21` |
| SQLite `busy_timeout` | 5000 ms | `apps/api/internal/store/store.go:57`（DSN `_busy_timeout=5000`） |
| 调度器 tick | 30s | `apps/api/modules/scheduledtasks/scheduler.go:24` |

> **重要**：`WriteTimeout = 10s`（`config.go:464`）意味着**任何同步请求超过 10s 都会被服务器写超时掐断**。
> 这为「同步批量/导出必须异步化」提供了硬性时间上界证据：
> 全量 reconcile、10000 行导出、无界 purge-all 都可能逼近或超过该预算。

---

## 6. 钉住当前同步行为的测试（TESTS THAT PIN CURRENT SYNC BEHAVIOR）

这些测试构成异步化的**回归面**——任何改动都必须让它们保持绿。

### 6.1 batch-delete 同步 / 原子语义

| 文件 | 断言 |
|------|------|
| `apps/api/internal/handler/users_batch_test.go:14` `TestUsersBatchDelete` | `POST /api/users/batch-delete` 返回 **200**，body `{"deleted": n}`（`:48-53`） |
| `apps/api/internal/handler/users_batch_test.go:72` `TestUsersBatchDeleteFailClosed` | 空选/非法 body → 400 `EMPTY_SELECTION` / `INVALID_SELECTION_KEY`（`:88-106`） |
| `apps/api/internal/handler/users_batch_test.go:119` `TestUsersBatchDeleteAtomicRollbackHTTP` | 批含 last admin → **409 `LAST_ADMIN`**（`:164-172`），且**选择内其他 id 全部未删除**（`:176-184`） |
| `apps/api/internal/handler/users_batch_test.go:193` `TestUsersBatchDeleteRejectsRemovingAllAdminsHTTP` | 同批含全部 admin → 拒绝，零删除（`:224`、`:253`） |
| `apps/api/modules/authsession/users_repository_test.go:14` `TestDeleteUsersBatchAtomicRollback` | 仓库层：中途失败（not-found/self/last-admin）整批回滚（`:70`、`:188`） |
| `apps/api/modules/authsession/users_repository_test.go:145` `TestDeleteUsersBatchCleansRoleAndMfaLinks` | 批删同时清理 `user_roles` / `user_mfa`（`users_repository.go:401-406`） |
| `apps/api/modules/authsession/users_repository_test.go:178` `TestDeleteUsersBatchRejectsRemovingAllAdmins` | 批级 last-admin 守卫 |
| `apps/api/modules/authsession/roles_repository_test.go:13` `TestRolesRepositoryBatchDeleteAtomicRollback` | 仓库层：系统角色/占用角色导致整批回滚（`roles_repository.go:201-210`） |
| `apps/api/internal/handler/resources_test.go:365,434-444` | 自作用域下 batch-delete **仅删本人行**（`deleted:1`，`o-4` 仍在） |
| `apps/api/internal/handler/recyclebin_test.go:319` `TestRecycleFactoryHookBatchDeleteSnapshots` | **顺序回退路径**必须为**每个** id 记录快照（N≥2，断言 `len(trash.calls) == 2`，`:352-360`）——这条直接钉住 dict-types 走的是回退路径 |
| `apps/api/modules/recyclebin/service_test.go:58` 附近 | 同 `now` 下批量删除记录多条快照不撞 PK |
| `apps/web/src/app/representative-pages.integration.test.tsx:401,500,564` | 前端：`POST /api/users/batch-delete` **只发一次请求**（`:562-564` 注释「One logical POST with the normalized `$selection.keys` body」），body 为归一化 keys |
| `apps/web/src/protocol/capability-declaration.guard.test.ts:57` | `actions.batch.request` 能力声明必须以 `/batch-delete` 为证据（`/\/batch-delete/.test(text)`） |
| `apps/web/src/protocol/upstream/request-construction.cases.json:1029-1192` | 11 个 batch request fixture，钉住 `batchMapping` 构造（body `$selection.keys`、query `$selection.count`） |
| `apps/api/modules/users/provider_test.go:65`、`roles/provider_test.go:66`、`datadictionary/provider_test.go:71,74`、`scheduledtasks/provider_test.go:71` | provider 路由清单**逐字**包含 `POST .../batch-delete`——改动路由即失败 |

### 6.2 导出 / 导入同步语义

| 文件 | 断言 |
|------|------|
| `apps/api/internal/handler/data_transfer_test.go:50` `TestExportUsersCSV` | `GET /api/export/users` 同步返回 CSV（UTF-8 BOM + 表头 + 行） |
| `apps/api/internal/handler/data_transfer_test.go:86` `TestExportRolesCSV` | 同上，roles |
| `apps/api/internal/handler/data_transfer_test.go:100` `TestExportUnknownResource404` | 未知资源 → 404 |
| `apps/api/internal/handler/data_transfer_test.go:110` `TestExportPermissionGated` | 无 `data.export` → 403 |
| `apps/api/internal/handler/data_transfer_test.go:254` `TestEditorCanExport` | editor 可导出（权限面） |
| `apps/api/internal/handler/data_transfer_test.go:266` `TestExportAuditLogged` | 导出写审计事件 |
| `apps/api/internal/handler/data_transfer_test.go:290` `TestExportCSVEscaping` | RFC 4180 转义 |
| `apps/api/internal/handler/data_transfer_test.go:305` `TestExportFormulaInjectionGuarded` | 公式注入中和（`'=HYPERLINK`） |
| `apps/api/internal/handler/data_transfer_test.go:351` `TestExportLimitExceeded` | `pageSize=20000` → **400**（钉住 10000 上限） |
| `apps/api/internal/handler/data_transfer_test.go:122` `TestImportUsersPartialApply` | **部分成功**语义：合法行提交、非法行报错（no-rollback） |
| `apps/api/internal/handler/data_transfer_test.go:166` `TestImportUsersValidationErrors` | 逐行错误报告 `{applied, failed, total, errors}` |
| `apps/api/internal/handler/data_transfer_test.go:197` `TestImportForeignFileForbidden` | 他人文件 → 403 |
| `apps/api/internal/handler/data_transfer_test.go:213` `TestImportMissingFile404` | 文件不存在 → 404 |
| `apps/api/internal/handler/data_transfer_test.go:223` `TestImportPermissionGated` | 无 `data.import` → 403 |
| `apps/api/internal/handler/data_transfer_test.go:326` `TestImportSizeLimit` | 3 MiB 上传（上传允许 8 MiB）→ 导入 **413**（钉住 2 MiB） |
| `apps/api/internal/handler/data_transfer_test.go:340` `TestImportMissingHeaderInvalidCsv` | 缺表头 → 400 `INVALID_CSV` |
| `apps/api/internal/handler/data_transfer_test.go:361` `TestImportRoleAssignmentBoundary` | `roles=admin` 逐行失败（非整请求 403） |
| `apps/api/internal/handler/data_transfer_test.go:398` `TestImportUnknownResource404` | 未知资源 → 404 |
| `apps/api/internal/handler/w16_batch_b_test.go:13,31` | 导入模板 + 字段错误 |
| `apps/api/internal/handler/operations_test.go:195` | 操作日志导出复用同一 `maxExportRows` |

### 6.3 recycle-bin 同步语义

| 文件 | 断言 |
|------|------|
| `apps/api/internal/handler/recyclebin_test.go:81` `TestRecycleBinRestoreAndPurge` | 恢复 200 / 物理删除 204 / 删除后详情 404（`:115-121`） |
| `apps/api/internal/handler/recyclebin_test.go:365` `TestRecycleBinPurgeAll` | 非 admin → **403**；admin → **200 `{"purged":2}`**（`:374-391`），之后列表 total=0（`:395-405`） |
| `apps/api/modules/recyclebin/store/repository_test.go:116` `TestPurge` | 单行物理删除 |
| `apps/api/modules/recyclebin/store/repository_test.go:132` `TestPurgeAllUnrestored` | 全量删除仅针对 `restored_at IS NULL` |
| `apps/api/modules/recyclebin/store/repository_test.go:163` `TestMarkRestoredTwiceFails` | 重复恢复失败 |
| `apps/api/modules/recyclebin/service_test.go:253` `TestRestoreAtomicityRollsBackOnFailedMark` | 恢复失败回滚 |
| `apps/api/modules/recyclebin/provider_test.go:99` | 恢复冲突 HTTP 面 |

### 6.4 wallet 异步 Job（作为模板的回归面）

| 文件 | 断言 |
|------|------|
| `apps/api/internal/handler/wallet_test.go:680` `TestWalletReconcileBadBodyAndWriteGate` | 非法 body → 400；无 `wallet.write` → 403 |
| `apps/api/internal/handler/wallet_test.go:273,302` | 路由清单含 `GET /api/wallet/reconcile/runs` |
| `apps/api/internal/jobs/runner_test.go:299`、`repository_test.go:51,111`、`shutdown_reclaim_test.go:61` | Job 六态/租约/取消/关停回收 |
| `apps/api/modules/wallet/store/repository_test.go:212,263` | reconcile 一致/不一致检测 |

---

## 7. 对首波筛选的直接影响（供 I-038-003 决策）

1. **分母必须先修正**：生产页面**零**批量动作（§3.1）。可选分母口径：
   （a）已挂载 `batch-delete` 路由的 **4 个资源**；（b）已具备
   `table.selection` + `actions.batch.request` 能力声明的页面；（c）§2 中的**非批量类长操作**。
2. **最强的 3 个异步候选**（§4 #5/#6/#7）：**导出**、**导入**、**recycle-bin purge-all**。
   三者都满足「长/行数无上界或上界很大 + 产出物可回看 + 与 Job 的 result/progress 天然契合」。
3. **batch-delete 必须保持同步**，且 VP-038 已把它写成显式非目标
   （`docs/vision/plans/VP-038-batch-operations-and-job-center.md:51`）。若需异步删除，
   应**新增**端点而非改既有路由。
4. **同步路径的硬性时间上界是 `WriteTimeout = 10s`**（`config.go:464`）——
   这是论证「哪些操作必须异步」最有力的既有约束。
5. **Job 运行时已通用**（`apps/api/internal/jobs`，六态 + progress + cancel + result 过期），
   首波无需新建 Job 基础设施，只需新增 `Kind` 与 handler。

---

## 待确认 / 未知

| # | 未确认项 | 影响 | 建议核实动作 |
|---|---------|------|-------------|
| U-01 | **前端是否有任何 job 轮询 / 进度 UI**（`apps/web/src` 未做完整 `job`/`poll`/`setInterval` 检索） | 决定首波是否需要「结果中心」UI 工作量的前置判断 | 全量 grep `apps/web/src` 的 `job`/`poll`/`progress` |
| U-02 | `render.tsx` 批量确认弹窗、成功 reload、清选、失败保留选择的**逐行实现**（仅确认了判定分支 `:331` 与 selection 传递 `:738`） | 异步化后 UX 兼容口径 | 通读 `apps/web/src/renderer/render.tsx` 批量相关段落 |
| U-03 | wallet reconcile **每次运行实际触及的行数**、是否有分批/上限、per-run 耗时预期 | 判断它是否真是「大 N」先例 | 读 `apps/api/modules/wallet` 的 reconcile 服务实现 |
| U-04 | reconcile Job 的 **payload / result 结构体定义**与 handler 注册方式（`SubmitReconcile` 内部） | 复用模板的具体形状 | 读 `SubmitReconcile` 与其 worker 注册点 |
| U-05 | `datadictionary` 删除类型是否**级联删除**其条目 | 影响 batch-delete 的规模估计 | 读 `dictTypeEntity.Delete` / store 层 |
| U-06 | **未发现保留期自动清理 sweeper**（仅确认无 `PurgeAll` 定时调用点） | 若有则 purge-all 的异步优先级更高 | 全量检索 `recycle_items` 的定时/过期清理 |
| U-07 | 各生产页面（wallet / scheduled-tasks / recycle-bin）**前端 UI 是否已暴露**对应长操作按钮 | 影响「用户可见性」列 | 读对应 `apps/api/modules/*/schema/*.json` 与 `apps/web/src` |
| U-08 | `data-transfer` 的 `provider.go` 具体路由注册行（本次引用 `:44` 来自上级同步，未自行复核） | 证据完整性 | 读 `apps/api/modules/datatransfer/provider.go` |
| U-09 | 是否存在 `DELETE ... WHERE id IN` 形式的多行删除（本次只找到 `recycle_items` 的无界 DELETE 与动态 `IN` 的 SELECT） | 长操作完备性 | 全量 grep `DELETE FROM` |
| U-10 | 导出「文件大小/耗时」的实测或文档预期（仅找到 10000 行上限与 10s WriteTimeout） | 异步化收益的量化依据 | 查 `docs/` 中 data-transfer 相关决策文档 |
