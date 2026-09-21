---
title: R1 侦察 · 批量操作与长操作清单（I-038-003）
status: draft
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.3.1
---

# R1 侦察 · 批量操作与长操作清单（I-038-003）

> 只读侦察报告。范围：`apps/api`（Go）+ `apps/web`（React）。所有结论附 `file:line`。
> 未取到证据的项一律标注「未核实」，不臆测。生成于 2026-09-19。
>
> **版本沿革**
> - v0.1.0：初版（工厂批量面、模块浅扫、测试回归面、规模上界）。
> - v0.2.0：并入 `admin.data-transfer` / `admin.wallet` / Job 基础设施深挖。
> - **v0.3.0**：并入其余模块全量扫描（scheduledtasks / recyclebin / mfa / users / roles /
>   datadictionary / filelibrary / systemmonitoring / operationlog / notifications / settings），
>   新增 §2.9 后台周期任务全表；关闭 U-01、U-03～U-08；新增 3 个此前遗漏的长操作候选。
> - **v0.3.1**：响应 `GOAL-002` 独立审计 `A-002 F-001`——更正批量路由计数
>   （「4 处 / 6 条路由」→ **4 个模块 / 5 条路由 / 5 个非只读 Resource ID**）。结论方向不变。

## 0. 结论摘要

1. 全仓库实现 `DeleteBatch` 的实体**只有 2 个**：`usersEntity`、`rolesEntity`。
2. 通用工厂把 `POST {path}/batch-delete` 挂在**每一个非只读资源**上，注册面 = **4 个模块 / 5 条路由 /
   5 个非只读 Resource ID**（users、roles、dict-types、dict-entries、scheduled-tasks）；其余 `ResourceRoutes`
   调用点（files、task-runs、monitoring-errors、operations）均 `ReadOnly: true`，不挂批量路由。
   其中 dict-types、dict-entries、scheduled-tasks **没有** `DeleteBatch`，走**顺序删除回退路径**
   （非原子、首个失败即停、已删行不回滚）。

   > **v0.3.1 更正（2026-09-19，响应 `GOAL-002 A-002 F-001`）**：v0.1.0～v0.3.0 的「4 处 / 6 条路由」为**计数错误**
   > （多计 1 条，源自把 §1.2 的 catch-all 空行当作实际路由）。经独立审计复验，正确数字为
   > **4 个模块 / 5 条路由 / 5 个非只读 Resource ID**。见 `GOAL-002/03-audit/A-002-…` F-001 与 `A-003` 响应。
3. **生产页面 schema 没有任何批量动作**——`batchMapping` / `requiresSelection` 只出现在
   `dev.examples` 的 `admin-list-batch` 范例（且该模块**仅在 `demo` profile**）。
   ⇒ **批量 UI 的已交付分母为 0 个生产页面**。后端能力先于前端页面落地。
4. 唯一已交付的异步 Job 先例是 `admin.wallet` reconcile（202 + jobId + cancel/retry/result）。
   Job 运行时 `apps/api/internal/jobs` 是**通用**的（六态 + lease + progress + cancel + result TTL），
   且**前端目前完全不轮询它**（`apps/web/src` 对 `wallet/jobs`/`jobId`/`progress` 零命中）。
5. **全仓库只有 4 个周期后台循环**（§2.9），其中 1 个是 Job poll loop（10s）。
   **没有** recycle-bin 保留期 sweeper，**没有** cron 库，**没有** `time.AfterFunc`。
6. 导出/导入**同步**且导出**在内存拼整个 CSV**（非流式），上限 10000 行 / 2 MiB。

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
| id 校验 | `resources.go:838-875` | 空选 → 400 `EMPTY_SELECTION`；非标量 → 400 `INVALID_SELECTION_KEY`；去重保序；**无条数上限**（仅受 4 KiB 约束） |
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
| 1 | `users` | `/api/users` | `apps/api/internal/handler/users.go:277` | ✅ | **单事务原子**（`DeleteUsersBatch`，`apps/api/modules/authsession/users_repository.go:346-415`，`r.withTx` at `:352`） | 有（快照在整批提交后逐条记录，`resources.go:907-926`） | `users.write` |
| 2 | `roles` | `/api/roles` | `apps/api/internal/handler/roles.go:201` | ✅ | **单事务原子**（`DeleteRolesBatch`，`apps/api/modules/authsession/roles_repository.go:187-221`，`r.withTx` at `:193`） | 有（同上机制） | `roles.write` |
| 3 | `dict-types` | `/api/data-dictionary/types` | `apps/api/internal/handler/dictionary.go:284-303` | ❌ **未实现** | **顺序循环回退**（`resources.go:927-972`）：逐条删除、首个失败即 return（`:957-960`）、**已删行不回滚（可部分提交）** | 有（`Trash: recorder` at `dictionary.go:302`）；走 `TrashTxDeleter` 分支时**每行一个事务**（`resources.go:940-955`） | `dictionary.write`（`dictionary.go:299`） |
| 4 | `dict-entries` | `/api/data-dictionary/entries` | `apps/api/internal/handler/dictionary.go:304-321` | ❌ **未实现** | 同上顺序回退 | 有（`dictionary.go:320`） | `dictionary.write`（`dictionary.go:317`） |
| 5 | `scheduled-tasks` | `/api/scheduled-tasks` | `apps/api/internal/handler/scheduledtasks.go:344-361` | ❌ **未实现** | 同上顺序回退 | 有（`scheduledtasks.go:360`） | `tasks.write`（`scheduledtasks.go:357`） |
| 6 | 其他所有非只读资源 | 各自 path | — | ❌ | 同上顺序回退 | 视是否注入 `Trash` 而定 | 各自 `{id}.write` |

**明确不挂 batch-delete 的只读资源**（`resources.go:302`）：

- `system-monitoring` errors — `ReadOnly: true`（`apps/api/internal/handler/systemmonitoring.go:89`）
- `file-library` — `ReadOnly: true`（`apps/api/internal/handler/filelibrary.go:169`）⇒ **文件库没有批量删除端点**
- `operations` / `activity` 日志（`apps/api/internal/handler/operations.go:16`）

### 1.3 两条路径的行为差异（关键）

`apps/api/internal/handler/resources.go:898-981`：

- **BatchDeleter 路径**（`resources.go:899-926`）：先整批 `Get` 快照（`:907-913`），再 `batch.DeleteBatch(ids, user)`（`:915`），
  失败 → `writeEntityError` 直接返回（`:916`），**整批回滚**；成功后逐条 `Trash.Record`（`:919-926`）。
- **回退路径**（`resources.go:927-972`）：`for _, id := range ids`（`:929`）逐条 `Get`（`:934`）→
  `DeleteTrashTx`（`:943`）或 `Delete`（`:957`），任何一条失败即 `return`（`:957-960`），
  **前面已删的行已经提交、不回滚**。成功返回 `{"deleted": deleted}`（`:971`）。
- 自作用域（self scope）预过滤在两条路径之前：`resources.go:884-896`，逐 id `Get`（**N+1**），
  非本人行被**跳过**（不报 404），全被过滤时返回 `{"deleted":0}`（`:893`）。

> **风险提示（供异步化设计参考）**：回退路径的「部分提交 + 首个失败即停」意味着
> 客户端拿到 500/409 时**无法知道已删了哪些**——这在同步交互下已经是个已知缺口，
> 异步化（带逐项结果）反而能修好它。

### 1.4 users 的批级 last-admin 守卫（异步化时不可丢失的语义）

- 守卫实现：`apps/api/modules/authsession/users_repository.go:379-393`
  （`countAdminUsersExcludingBatch`，`:417` 起），批内含**全部** admin → `ErrLastAdmin`，整批回滚。
- 其他守卫：not-found `users_repository.go:359-361`、self `:362-364`。
- 背景：这是 W3 独立审计 F-001 修复
  （`docs/workspaces/workspace-009-production-hardening/GOAL-004-w3-security-audit-remediation/03-audit/A-002-w3-independent-cross.md:123-128`）。
- roles 侧守卫：`roles_repository.go:194-211` — 系统角色 `ErrRoleSystem`（`:201-203`）、被占用 `ErrRoleInUse`（`:204-210`）。

---

## 2. 其他批量类操作（OTHER BATCH-LIKE OPERATIONS）

### 2.1 `admin.data-transfer`（导出 / 导入）— **同步**

模块本身只有 `provider.go`（74 行，无 handler/service，无迁移、无页面）：

- 路由声明：`apps/api/modules/datatransfer/provider.go:44`
- 装配：`provider.go:55-64`（`ExportRoutes` / `ImportRoutes`）
- 权限策略：`provider.go:66`（`data.export` → `PolicyAdminEditor`）、`provider.go:67`（`data.import` → `PolicyAdmin`）
- `CompiledPersistence` 返回 nil（`provider.go:50-52`）⇒ **不写任何 job 行**

| 操作 | endpoint | 位置 | sync/async | 上限 | 权限 |
|------|----------|------|-----------|------|------|
| 导出 CSV | `GET /api/export/{resource}` | `apps/api/internal/handler/export.go:44` | **同步** | `maxExportRows = 10000`（`export.go:25`）；超限 → 400 `INVALID_EXPORT_LIMIT`（`export.go:130-134`） | `data.export`（`export.go:53`） |
| 导入 CSV | `POST /api/import/{resource}` | `apps/api/internal/handler/import.go:45` | **同步** | `maxImportBytes = 2 << 20`（2 MiB，`import.go:30`）；超限 → 413 `FILE_TOO_LARGE`（`import.go:179-182`）；空文件 → 400 `INVALID_FILE`（`:175-178`） | `data.import`（`import.go:60`） |
| 导入模板 | `GET /api/import/{resource}/template` | `apps/api/internal/handler/import.go:51` | 同步 | — | `data.import`（同一 gate，`import.go:52`） |

**导出请求/响应形状**：无 JSON body；query 参数 `q`/`sort`/`order`/`pageSize`（`export.go:129-142`）。
响应是**裸 CSV 附件，不是 envelope**：

```go
// export.go:210-214
w.Header().Set("Content-Type", "text/csv; charset=utf-8")
w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%q", resource+".csv"))
w.Header().Set("X-Content-Type-Options", "nosniff")
w.WriteHeader(http.StatusOK)
_, _ = w.Write([]byte(out.String()))
```

表头顺序冻结：users `export.go:164`、roles `export.go:185`。
单元格渲染 `exportRow` `export.go:69-87`；公式注入中和 `formulaSafe` `export.go:233-242`。
支持资源仅 users / roles（`export.go:125-128`），其他 → 404 `RESOURCE_NOT_FOUND`。

**导入请求/响应形状**：**非 multipart**——CSV 先经 `POST /api/upload`（multipart，`upload.go:264`）
上传并以 fileId 引用；body 仅 `{"fileId":"..."}`，上限 `maxResourceBodyBytes` 4 KiB（`import.go:144`）。
`/api/files/{id}` 前缀归一化 `import.go:157-159`；owner 校验 403（`import.go:171-174`）。
响应 envelope（200 OK，`import.go:200`）：

```go
// import.go:74-92
type importRowError struct {
    Row     int    `json:"row"`
    Field   string `json:"field,omitempty"`
    Message string `json:"message"`
}
type importFieldError struct {
    RowNumber int    `json:"rowNumber"`
    Field     string `json:"field"`
    Reason    string `json:"reason"`
}
type importResult struct {
    Applied     int                `json:"applied"`
    Failed      int                `json:"failed"`
    Total       int                `json:"total"`
    Errors      []importRowError   `json:"errors"`
    FieldErrors []importFieldError `json:"fieldErrors,omitempty"`
}
```

**导出不是流式的**：整个 CSV 在内存 `strings.Builder` 中拼装
（`export.go:191-199` → `writer.Flush()` 在 `:204` 刷到 **Builder 而非 ResponseWriter** →
`:214` 一次性 `w.Write`）。行还先累积进 `[][]string`（`export.go:145`、`:166`、`:187`），
**峰值内存 ≈ 2× CSV**。包内无 `io.Copy`、无 `http.Flusher`、无对 `w` 的 `Flush()`
（全 `internal/handler` 检索：`Flush()` 仅 `export.go:204` 与 `operations_export.go:106`，均为 csv writer flush）。
**文件头注释自述 "streams the filtered resource list as CSV"（`export.go:1-2`）与实现不符。**

**导入是逐行循环 + 逐行 INSERT，非批量插入**：

```go
// import.go:239-255（循环）
for {
    record, err := reader.Read()
    if err == io.EOF { break }
    ...
    rowNumber++
    result.Total++
```
```go
// import.go:301-310（写入：一次一行）
_, createErr := h.repository.CreateUserManagement(authsession.User{
    ID:       newUserIDValue(),
    Username: row["username"],
    ...
```
每行先 bcrypt 哈希（`import.go:281` `auth.HashPassword(row["password"], passwordHashCost)`）再单独 INSERT。
畸形 CSV 行 → 记录后 `break`（`import.go:244-253`）。列白名单 `import.go:234`；
逐行校验 `validateImportUser` `import.go:355-371`；逐行委派边界 `importRoleAssignmentError` `import.go:330-353`。
**导入行数无上限**（仅受 2 MiB 字节预算约束）。

**事务边界：导入逐行提交、非原子**。handler 从不接触 `kernel.Tx`
（`import.go` 内 `kernel.Tx|withTx|runner.Run` 零命中）；每行的 INSERT 各自开事务
（`apps/api/modules/authsession/users_repository.go:76` `CreateUserManagement` → `withTx`；
`apps/api/modules/authsession/repository.go:177-185` `withTx` 定义）。
⇒ N 行 = N 个独立已提交事务；后续失败**无法**回滚先前行（即刻意设计，见下）。

**无 dry-run / 预览**：`provider.go` 只声明 3 条路由（`:44`）；
`dry.?run|validateOnly` 在 `apps/api` 只命中无关 CLI（`apps/api/cmd/schema-ui/configpkg.go:10,451`、`main.go:52-56`）。
设计文档显式否决两段式：`docs/workspaces/workspace-011-admin-functional-modules/GOAL-004-r2-f02-data-import-export/01-decision/D-002-s1-freeze.md:85`
—「不做导入预览/两段式（importId + 错误表）… R2 采用单段提交 + 结构化错误报告」。
最接近的是模板端点（`import.go:120-131`，仅返回表头行）。

### 2.1.1 审计/操作日志导出 — **同步（此前遗漏的候选）**

| 项 | 位置 |
|----|------|
| endpoint | `GET /api/operations/export` — `apps/api/internal/handler/operations_export.go:19-24` |
| 上限 | **复用 `maxExportRows`**：`operations_export.go:29-33`（`pageSize > 10000` → 400） |
| 流式？ | 否，同为 builder-then-write（`operations_export.go:106` 仅 csv writer flush） |
| 前端 | `apps/web/src/components/activity-export.tsx:43` 调用 `/api/operations/export`；测试 `apps/web/src/components/activity-export.test.tsx:75` 断言 `?pageSize=10000` |
| 可见性 | **用户可见**（活动日志页导出按钮） |

> 这是**第二个导出面**，与 §2.1 的 `GET /api/export/{resource}` 共享同一 10000 行上限与同一
> buffer-then-write 实现。异步化若只做 data-transfer 而不做 operations，会留下不一致。

### 2.2 `admin.wallet` reconcile — **异步（唯一先例）**

| 项 | 位置 |
|----|------|
| 提交 reconcile | `POST /api/wallet/reconcile` → `apps/api/internal/handler/wallet.go:357-385`；`jobService.SubmitReconcile(...)` at `:379`；响应 **202 Accepted** at `:384` |
| 查询 Job | `GET /api/wallet/jobs/{id}` — `wallet.go:387-398`，权限 `wallet.read` |
| 取消 | `POST /api/wallet/jobs/{id}/cancel` — `wallet.go:400-411`，权限 `wallet.write` |
| 重试 | `POST /api/wallet/jobs/{id}/retry` — `wallet.go:413-424`，权限 `wallet.write` |
| 读结果 | `GET /api/wallet/jobs/{id}/result` — `wallet.go:426-451`；queued/running → 409 `JOB_RESULT_NOT_READY`（`:438`）；expired → 410 `JOB_RESULT_EXPIRED`（`:440`）；failed/cancelled → 200 元数据（`:441-442`）；succeeded → `Content-Disposition: attachment; filename="wallet-reconcile-<jobID>.json"`（`:445`） |
| 提交权限 | `wallet.write`（`wallet.go:358`），注释明确（`wallet.go:353-356`）；原冻结为 `wallet.read`，W11 F-006 收紧 |
| 运行历史列表 | `GET /api/wallet/reconcile/runs` — `wallet.go:454-478`，权限 `wallet.read` |
| Job 表示 | `walletJobToMap` `wallet.go:989-1006`；succeeded 时附 `"resultUrl"`（`:1002-1004`） |

**作用范围**：`body.AccountID` 为空 ⇒ **全账本 reconcile**（`wallet.go:366-370` 注释）。
body 上限 4 KiB（`wallet.go:365`）；非法 body → 400（`:371-374`，W11 F-006 修复：
此前垃圾 body 会静默退化为**全账本** reconcile）。

### 2.2.1 reconcile 每次运行实际做什么 / 能触及多少行

入口：`Service.ReconcileOnceTx`（`apps/api/modules/wallet/provider.go:187-189`）
→ `Repository.ReconcileOnceTx`（`apps/api/modules/wallet/store/repository.go:824-894`）。

1. 幂等预检（run id 已存在则返回既有 run）：`repository.go:828-832`。
2. 选取账户集合 —— **`accountID == ""` 时为全表**：

```go
// repository.go:840-865
if accountID == "" {
    rows, err := tx.Query(ctx, "SELECT id FROM wallet_accounts ORDER BY id")
    ...
    ids = append(ids, id)
} else {
    ... "SELECT COUNT(*) FROM wallet_accounts WHERE id = ?" ... if exists == 0 { return nil, ErrNotFound }
    ids = append(ids, accountID)
}
```

3. 核心循环 —— 逐账户重放账本链：

```go
// repository.go:866-870
for _, id := range ids {
    if reason, ok := checkAccountChain(tx, id); !ok {
        mismatches = append(mismatches, mismatch{AccountID: id, Reason: reason})
    }
}
```

4. 落一行 run（表 `wallet_reconciliation_runs`；DDL `apps/api/modules/wallet/migration/migration.go:56-64`）：

```go
// repository.go:886-889
if _, err := tx.Exec(ctx,
    `INSERT INTO wallet_reconciliation_runs (id, account_id, result, mismatch_count, details, actor_id, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
    runID, acctID, result, len(mismatches), details, actorID, now.Unix(),
```

**单账户重放读取全量流水，无 LIMIT**：

```go
// repository.go:921-923
rows, err := tx.Query(context.Background(),
    `SELECT entry_type, amount_delta, balance_after_total, balance_after_available, balance_after_frozen FROM wallet_ledger_entries WHERE account_id = ? ORDER BY created_at ASC, id ASC`,
    accountID,
```

校验项：快照不变量 `afterTotal != afterAvail+afterFrozen`（`:936-938`）；重放等价性（`Apply` `:944-950`）；
账户不变量 `:963-965`；最后快照 == 当前余额 `:966-968`（读取 `:957-960`）。
结果值 `ResultConsistent`/`ResultInconsistent`（`repository.go:142-145`）。

> **行数结论：无界**。账户扫描无 LIMIT（`repository.go:841`），
> 单账户流水重放无 LIMIT（`repository.go:921-923`）——**从创世块全量捞入内存重放**。
> 无分页、无游标、无 max-accounts 常量。
> 独立审计已就此定性：`docs/workspaces/workspace-011-admin-functional-modules/GOAL-019-r3-s14-wallet-ledger/attachments/audit-A-008-independent.md:32`
> —「**【性能隐患】对账 O(N) 全量内存扫描**：checkAccountChain 从创世块全量捞取重放；
> 热点账户 … 数十万流水 → 长事务锁定、API 超时、OOM」。
> 每次 run 触及：**读** `wallet_accounts`、`wallet_ledger_entries`；**写** `wallet_reconciliation_runs`
> （job 路径另写 `operation_log`、`jobs`），全部在一个事务内。

### 2.2.2 reconcile Job 的 payload / result / 入队 / 轮询

Kind 常量：`apps/api/modules/wallet/jobs.go:19` — `const ReconcileJobKind = "wallet.reconcile"`。

**输入 payload**：

```go
// jobs.go:21-23
type reconcileJobPayload struct {
    AccountID string `json:"accountId,omitempty"`
}
```
```go
// jobs.go:65-72
payload, err := json.Marshal(reconcileJobPayload{AccountID: accountID})
...
job, err := s.runner.Submit(ctx, jobs.CreateInput{
    ID: id, Kind: ReconcileJobKind, Payload: payload, ActorID: actor.ID,
    CorrelationID: correlationID, MaxAttempts: jobs.DefaultMaxAttempts, Now: now,
})
```

**输出 result —— 没有结构体，是 ad-hoc `map[string]any`**：

```go
// jobs.go:214-221
func reconciliationResult(run walletstore.ReconciliationRun) map[string]any {
    return map[string]any{
        "id": run.ID, "accountId": run.AccountID, "result": run.Result,
        "mismatchCount": run.MismatchCount, "details": run.Details,
        "actorId":   run.ActorID,
        "createdAt": run.CreatedAt.UTC().Format("2006-01-02T15:04:05.000Z07:00"),
    }
}
```

返回 `return json.Marshal(reconciliationResult(*run))`（`jobs.go:158`）。

**入队路径**：`SubmitReconcile` `jobs.go:45-79` → `runner.Submit` → `repo.Create` INSERT `status='queued'`
（`apps/api/internal/jobs/repository.go:40-48`）；`Submit` 内部唤醒 scan loop
（`apps/api/internal/jobs/runner.go:152-158`）；queued 审计 id 预分配以保证审计顺序
（`jobs.go:51-64`，事件 `wallet.reconcile.queued` `jobs.go:76`）。

**worker 绑定**（run id == job id）：

```go
// jobs.go:127-128
return func(tx kernel.Tx) (json.RawMessage, error) {
    run, err := s.service.ReconcileOnceTx(context.Background(), tx, payload.AccountID, job.ID, job.ActorID, s.now().UTC())
```
注册：`runner.RegisterWithTerminalHook(ReconcileJobKind, s.runReconcile, s.recordTerminal)` `jobs.go:39`；
组合装配 `apps/api/internal/composition/composition.go:584-588`。

**原子性**：run INSERT + 成功审计（`jobs.go:149-155`，`TransactionalRecorder`）+ `jobs` 状态 UPDATE
在**同一事务**内完成——`CompleteWithCommit` 在 `repository.go:195` 调 `commit(tx)`，
jobs UPDATE 在 `:202`，同属一次 `r.runner.Run`。

**Actor 隔离**：所有读取经 `GetForActor(ctx, id, ReconcileJobKind, actorID)`
（`jobs.go:82`；`repository.go:66-74`）→ 跨 actor / 不存在 ⇒ 404 `JOB_NOT_FOUND`。
Result TTL 默认 24h（`runner.go:51`），读取时惰性过期 `ExpireIfDue`（`jobs.go:86-91`）。

**可取消 / 进度**：

```go
// apps/api/modules/wallet/jobs.go:95-107
func (s *JobService) Cancel(ctx context.Context, id, actorID string) (*jobs.Job, error) {
    if _, err := s.repository.GetForActor(ctx, id, ReconcileJobKind, actorID); err != nil { return nil, err }
    job, err := s.runner.Cancel(ctx, id, actorID)
    ...
    if job.Status == jobs.StatusCancelled { s.recordTerminal(*job) }
```
Runner：`Cancel` `runner.go:160-175`（running 时取消在飞 handler context）；
仓储：queued → 立即 `cancelled`；running → `cancel_requested=1`（`repository.go:124-149`）；
finalize `:151-160`；租约过期回收 `RecoverCancelledDueJobs` `:264-273`。
Retry 仅从 `failed` 且 `attempt < max_attempts`（`repository.go:220-238`）。

**进度：极粗——硬编码 10，成功时 100**：

```go
// apps/api/modules/wallet/jobs.go:121-126
if err := reporter.Progress(10); err != nil { return nil, err }
if reporter.Cancelled() { return nil, ctx.Err() }
```
`CompleteWithCommit` 置 `progress=100`（`repository.go:202`）。
**循环内不报进度**（循环在 store 层，无 `Reporter`）。`UpdateProgress` 限 0..99 且单调
（`progress <= ?`，`repository.go:115-122`）；runner 每 `HeartbeatInterval`（默认 10s，`runner.go:50`）
心跳 + 轮询取消（`runner.go:296-313`）。
⇒ UI 只能显示 queued/running/终态 + 粗糙 10→100，**拿不到真实百分比**。

### 2.3 Job 运行时（异步模板）

`apps/api/internal/jobs`（**不是** `apps/api/modules/jobs`——后者**只有迁移**：
`apps/api/modules/jobs/migration/{migration.go,provider.go,migration_test.go}`）：

- 六态：`apps/api/internal/jobs/model.go:18-25` — `queued` / `running` / `succeeded` / `failed` / `cancelled` / `expired`
- `Job` 结构：`model.go:38-59` — 含 `Kind`、`Payload`、**`Progress int`**（`:43`）、`CancelRequested`（`:44`）、
  `Attempt`/`MaxAttempts`、`LeaseOwner`/`LeaseVersion`/`LeaseExpiresAt`（租约）、`Result`、`ErrorCode`/`ErrorMessage`、
  `ActorID`、`CorrelationID`、`ResultExpiresAt`
- `CreateInput`：`model.go:61-69`；`Lease`：`model.go:71-75`
- `DefaultMaxAttempts = 3`：`model.go:27`
- 错误：`ErrNotFound`/`ErrInvalid`/`ErrTransition`/`ErrLeaseLost`/`ErrNotCancellable`/`ErrNotRetryable` `model.go:29-36`
- `NewID`（可排序 id）`model.go:78-84`；列清单 `jobColumns` `model.go:90-93`；`scanJob` `model.go:95-127`

**Handler / Reporter 契约**（新消费者要实现的东西）：

```go
// apps/api/internal/jobs/runner.go:15-27
type CommitFunc func(kernel.Tx) (json.RawMessage, error)
type Reporter interface { Progress(int) error; Cancelled() bool }
type Handler func(context.Context, Job, Reporter) (CommitFunc, error)
type TerminalHook func(Job)
```

### 2.3.1 Job 表 DDL 与 runner 默认参数

`apps/api/modules/jobs/migration/migration.go` — `const ModuleID = "core.jobs"` `:12`；
SQLite DDL `jobsDDL` `:14-48`；Postgres 变体 `jobsPGDDL` `:53-87`（时间列 BIGINT）；
descriptor `Descriptors()` `:89-98`（`Version: 42`，`Name: "async_jobs"`）。
provider `migration/provider.go:12-20`（`Register` 是 no-op `:20` —
「deliberately absent from runtime profiles」，`migration.go:1-3`）。

```sql
-- migration.go:15-35（节选）
CREATE TABLE jobs (
  id TEXT PRIMARY KEY, kind TEXT NOT NULL, status TEXT NOT NULL CHECK (status IN ('queued','running','succeeded','failed','cancelled','expired')),
  payload TEXT NOT NULL DEFAULT '{}', progress INTEGER NOT NULL DEFAULT 0 CHECK (progress BETWEEN 0 AND 100),
  cancel_requested INTEGER NOT NULL DEFAULT 0 CHECK (cancel_requested IN (0,1)),
  attempt INTEGER NOT NULL DEFAULT 0, max_attempts INTEGER NOT NULL DEFAULT 3 CHECK (max_attempts > 0 AND attempt <= max_attempts),
  lease_owner TEXT, lease_version INTEGER NOT NULL DEFAULT 0, lease_expires_at INTEGER,
  result TEXT, error_code TEXT, error_message TEXT,
  actor_id TEXT NOT NULL, correlation_id TEXT NOT NULL,
  created_at INTEGER NOT NULL, updated_at INTEGER NOT NULL, finished_at INTEGER, expires_at INTEGER,
```

外加**按状态的状态机 CHECK**（`:36-43`），使非法形态不可表示
（如 `succeeded` 要求 `result IS NOT NULL AND progress = 100 AND finished_at IS NOT NULL AND expires_at IS NOT NULL`）；
三个索引（`:45-47`）：`idx_jobs_runnable(status, cancel_requested, lease_expires_at, created_at)`、
`idx_jobs_actor(actor_id, kind, updated_at DESC)`、`idx_jobs_expiry(status, expires_at)`。
`payload` 与 `result` 为 **TEXT JSON 列**。

**Runner 默认参数**（`apps/api/internal/jobs/runner.go:48-54`）：

| 参数 | 默认值 | 行 |
|------|--------|----|
| `LeaseDuration` | 30s | `runner.go:48` |
| `HeartbeatInterval` | 10s | `runner.go:50` |
| `ScanInterval` | **10s** | `runner.go:51` |
| `ResultTTL` | **24h** | `runner.go:52` |
| `BatchSize` | 32 | `runner.go:53` |

注册必须在 `Start` 之前：`Register` / `RegisterWithTerminalHook` `runner.go:89-109`
（重复 kind 拒绝 `:104-106`；start 之后注册拒绝 `:101-103`）。
`ScanOnce` `:185-212` = 回收已取消（`:187`）+ 耗尽超次（`:194`）+ `ExpireDue`（`:201`）
+ `ListRunnable(now, BatchSize)`（`:204`）→ 逐 job 起 goroutine（`:247`）。
认领带租约围栏：`repository.go:76-104`。执行 + panic 兜底 + 终态通知：
`execute` `runner.go:258-319`（panic → 持久 `JOB_HANDLER_FAILED` `:282-286`）、
`finish` `:335-381`、`notifyTerminal` `:383-398`、`runnerReporter` `:421-439`。
组合装配：`composition.go:170-183`（`jobRuntime` + `newJobRuntime`）、
启用 `:588`、生命周期 Start `:1101`、Stop `:1123`/`:1160`。

### 2.3.2 新增异步操作的复用配方（由代码归纳）

1. 加 `const XxxJobKind = "module.op"` + payload struct（范式：`apps/api/modules/wallet/jobs.go:19-23`）。
2. 建 `JobService`，在 runner **启动前** `runner.RegisterWithTerminalHook(kind, handler, terminalHook)`
   （`jobs.go:34-43`）；`Submit` 把 payload marshal 进
   `jobs.CreateInput{ID: jobs.NewID(now), Kind, Payload, ActorID, CorrelationID, MaxAttempts: jobs.DefaultMaxAttempts, Now}`（`jobs.go:45-79`）。
3. handler 返回 `jobs.CommitFunc`，在给定的 `kernel.Tx` 内做领域写入并返回 JSON 结果
   （`jobs.go:116-159`）；可选 `reporter.Progress(n)` / `reporter.Cancelled()`（`jobs.go:121-126`）。
4. HTTP：提交 → `writeJSON(w, http.StatusAccepted, …)`；读取 → actor 隔离的
   `GetForActor(ctx, id, kind, actorID)`；结果 → 409 not-ready / 410 expired / 200 attachment
   （`apps/api/internal/handler/wallet.go:357-451`）。
5. **无需**改 `apps/api/modules/jobs`——`jobs` 表（迁移 42）已存在。

### 2.4 `admin.scheduled-tasks` — 手动触发是**同步**的

- `POST /api/scheduled-tasks/{id}/run` — 路由 `apps/api/internal/handler/scheduledtasks.go:445-465`
  （`Method: "POST"`，`Pattern` at `:448`）；handler `:449-464`：
  **在请求内联执行** `runner.Execute(*task, time.Now().UTC())`（`:459`），
  失败 → 500 `INTERNAL`（`:460`），成功 → **204 No Content**（`:463`）。
  ⇒ **同步**，**无 job 行、无 goroutine、无返回值、无 runId、无进度**。
- 权限 `tasks.write`（`scheduledtasks.go:450`）。
- `Execute` 实现：`apps/api/modules/scheduledtasks/scheduler.go:133-189`；
  panic 兜底 `:161-168`；记录 `task_runs` 行 `:184`。
  handler 通过 `context.Background()` 调用（`scheduler.go:167`）——**不继承请求 context**，
  客户端断开不会取消，**无超时**。
- 定时循环：`scheduler.go:65-85`，`tickInterval = 30 * time.Second`（`scheduler.go:24`）；
  单实例 best-effort，错过窗口不回补（`scheduler.go:3-5`）；按分钟槽去重（`:123-126`）。
- v1 只注册 `system.noop` handler（`scheduler.go:46-50`）⇒ 当前耗时 ≈ 0，
  但 `Execute` 是**通用同步派发点**，未来真实 handler 将在请求 goroutine 里跑。
- **批量 enable/disable：不存在**。只有通用 `batch-delete`（`scheduledtasks/provider.go:59`），
  启用/停用是逐行 `PATCH /api/scheduled-tasks/{id}` → `repository.UpdateTask(id, cron, name, enabled, ...)`
  （`scheduledtasks.go:249,257`；SQL `apps/api/modules/scheduledtasks/store/repository.go:190-193`）。
  列表有 `enabled` 过滤（`scheduledtasks.go:351`）但**无批量切换**。
- 调度扫描**无界**：`apps/api/modules/scheduledtasks/store/repository.go:233-239`
  — `SELECT ... FROM scheduled_tasks WHERE enabled = 1`，**无 LIMIT**。

### 2.5 `admin.recycle-bin` — 单行 restore/purge + **有 purge-all**

路由声明：`apps/api/modules/recyclebin/provider.go:44-48`。

| 操作 | endpoint | 位置 | 单/批 | sync/async | 权限 |
|------|----------|------|-------|-----------|------|
| 列表 | `GET /api/recycle-bin` | `apps/api/internal/handler/recyclebin.go:45-47` | — | 同步 | `recycle.read` |
| 详情 | `GET /api/recycle-bin/{id}` | `recyclebin.go:92-94` | 单行 | 同步 | `recycle.read` |
| 恢复 | `POST /api/recycle-bin/{id}/restore` | `recyclebin.go:109-111`（handler `:108-126`） | **单行** | 同步 | `recycle.write` |
| 物理删除 | `DELETE /api/recycle-bin/{id}` | `recyclebin.go:148-165`（route `:151`） | **单行** | 同步 | `recycle.write`（`:153`） |
| **清空回收站** | `POST /api/recycle-bin/purge-all` | `recyclebin.go:128-146`（route `:131`） | **全量批量** | **同步** | `recycle.write`（`:133`） |

- 列表**分页单查询**：`apps/api/modules/recyclebin/store/repository.go:96-172`；
  pageSize 在 `repository.go:135-137` 被**硬夹到 100**（`if pageSize > 100 { pageSize = 100 }`）。
- **恢复是单行且原子的，不是循环**：`apps/api/modules/recyclebin/service.go:95-124`；
  行 INSERT + 快照标记在**同一事务**（`service.go:103-108` `s.runner.Run(...)` + `MarkRestoredTx`）。
  分派是 `switch`，**只支持 3 种资源**（`service.go:185-196`：`dict-types`、`dict-entries`、
  `scheduled-tasks`；default → 错误）。**没有 restore-all / 批量恢复**。
- 单行 purge：`apps/api/modules/recyclebin/store/repository.go:245-256`
  — `DELETE FROM recycle_items WHERE id = ?`（0 行 → `ErrItemNotFound`）。
- **purge-all 是单条无界 SQL，同步、一个事务、inline 在 handler 里**：

```go
// apps/api/modules/recyclebin/store/repository.go:227-242
func (r *Repository) PurgeAllUnrestored() (int, error) {
    purged := 0
    err := r.runner.Run(context.Background(), func(tx kernel.Tx) error {
        res, err := tx.Exec(context.Background(), `DELETE FROM recycle_items WHERE restored_at IS NULL`)
        ...
```
  服务层 `Service.PurgeAll()` `apps/api/modules/recyclebin/service.go:132-134`；
  返回 `{"purged": n}`（`recyclebin.go:144`）。**无 LIMIT、无分批、无进度。**
- **保留期自动清理 sweeper：不存在**（已全量核实）。检索 `recycle.*retention`、
  `recycle_items.*deleted_at <`、`expireRecycle` → 零命中；`apps/api` 内**没有**任何
  goroutine / ticker / cron 触碰 `recycle_items`。快照**永久保留**直到人工 purge。

### 2.6 其他模块的批量/长操作

| 模块 | 操作 | 位置 | endpoint | sync/async | 说明 |
|------|------|------|----------|-----------|------|
| `notifications` | **全部标记已读** | handler `apps/api/internal/handler/notifications.go:162-177`；SQL `apps/api/modules/authsession/notifications_repository.go:216-233` | `POST /api/notifications/read-all`（`notifications.go:52`） | **同步** | 单条 `UPDATE notifications SET read_at=? WHERE user_id=? AND read_at IS NULL`（`notifications_repository.go:220`），返回 `{"updated": n}`（`notifications.go:175`）。**无界**（该用户全部未读）。有天然上限 `maxNotificationsPerUser = 500`（`notifications_repository.go:21-22`）⇒ 实际有界，**keep-sync** |
| `notifications` | 保留期剪枝 | `notifications_repository.go:88-136` | 无端点（写入时内联） | 同步 | `excess = count - maxNotificationsPerUser`（`:104`），**逐行删**（`:128-132`）；有界 |
| `mail` | 发件箱保留上限 | `apps/api/internal/mail/outbox.go:39` `DefaultOutboxCap = 500`；强制 `outbox.go:95-109` | 无端点（写入时内联） | 同步 | `DELETE FROM mail_outbox WHERE id NOT IN (SELECT id ... LIMIT ?)`；有界 |
| `wallet` | **代金券批量生成** | handler `apps/api/internal/handler/wallet.go:484-`；校验 `apps/api/modules/wallet/voucher/service.go:41` | `POST /api/wallet/vouchers/batches`（`wallet.go:484`） | **同步** | `count <= 0 \|\| count > 1000` 拒绝（`voucher/service.go:41`；handler 同判 `wallet.go:506`）⇒ **上限 1000**。权限 `wallet.voucher.issue`（`wallet.go:485`）；body 上限 16 KiB（`wallet.go:496`） |
| `settings` | **重置设置（删除全部品牌资源）** | handler `apps/api/internal/handler/settings.go:266-292`；`:286` `_ = assets.DeleteAll()` → `apps/api/internal/handler/raster_assets.go:379-390` | `POST /api/settings/{id}/reset`（`settings.go:41`） | **同步** | `DeleteAll` = `List` + 逐 id `Delete` 循环（`raster_assets.go:384-388`），**无界**（全部品牌资源） |
| `users` | 邀请（创建/撤销/重发） | `apps/api/internal/handler/invites.go:98-101` | `POST /api/users/invites`、`DELETE /api/users/invites/{id}`、`POST /api/users/invites/{id}/resend` | 同步 | **单目标**：创建 body 仅一个 `Email` + `Roles`（`invites.go:123-127`），一次 `CreateInvite`（`:134`）；重发 `r.PathValue("id")`（`:258`）。**无批量邀请、无批量重发** |
| `users` | 启用/停用/解锁 | `apps/api/internal/handler/users_state.go:43-45` | `POST /api/users/{id}/enable\|disable\|unlock` | 同步 | **逐用户单条**（`users_state.go:66,75,106,113`），无批量变体 |
| `users` | 角色分配 | `apps/api/internal/handler/users.go:~240-243`（`patch.Roles = &roles`） | `PATCH /api/users/{id}` | 同步 | **单目标**；`roles.assign` 门在 `users.go:355-375`。**无批量角色分配** |
| `mfa` | 重置他人 MFA | `apps/api/internal/handler/mfa.go:358-396`（route `:358`，`targetID := r.PathValue("id")` `:363`，`AdminReset(targetID)` `:380`） | `POST /api/users/{id}/mfa/reset` | 同步 | **单用户**；权限 `users.mfa-reset`（`apps/api/kernel/profile.go:172`）。**无批量重置 / 批量禁用 / 强制重注册**（`mfa.*[Bb]ulk\|AdminResetAll\|DisableAll` 零命中）；service（`modules/mfa/service.go:98-321`）与 store（`modules/mfa/store/repository.go:60-333`）**全部单 `userID` 签名** |
| `mfa` | 自服务 enroll/confirm/disable/recovery rotate | `mfa.go:204/256/279/321` | `POST /api/mfa/*` | 同步 | 全部单用户、body 上限 4 KiB |
| `roles` | 权限/菜单授权 | `apps/api/modules/authsession/roles_repository.go:357-380`（`replaceRolePermissions`）、`:382-402`（`replaceRoleMenuItems`） | `PATCH /api/roles/{id}`（`roles.go:187` `UpdateRoleWithGrants`） | 同步 | 循环的是**单个角色的权限 key**（`:360` 校验、`:371` `DELETE FROM role_permissions WHERE role_id = ?`、`:374-378` 插入），**不是多个角色**。**无批量授权、无批量用户指派**（`modules/roles/provider.go:43-49` 无此路由） |
| `filelibrary` | 上传/下载/删除 | 上传 `apps/api/internal/handler/filelibrary.go:255-298`；下载 `:218-253`；删除 `:181-216` | `POST /api/library/files/upload`、`GET .../download`、`DELETE /api/library/files/{id}` | 同步 | **无批量上传**（`file-library.json:17-21` `"multiple": false`）、**无批量删除**（资源 `ReadOnly: true` `filelibrary.go:169` ⇒ 不挂 batch-delete）、**无 zip 打包下载**（全 `apps/api` 无 `archive/zip`）、**无移动/重命名路由**。下载整包读入内存（`filelibrary.go:228` `store.load(id)`，`Content-Length` 由 `len(body)` `:249`） |
| `filelibrary` | **列表无界扫描** | `apps/api/internal/handler/filelibrary.go:106-125` | `GET /api/library/files` | 同步 | `objects.List(...)` 后**逐 id `Stat`**（`:112-123`），无 LIMIT；再内存过滤排序（`:52-68`）。同样模式见配额检查 `apps/api/internal/handler/upload.go:207-242`（`List` + `Stat` 循环 `:219-234`；注释 `:205-206`「Scanning is O(files) per upload」） |
| `systemmonitoring` | 状态/错误查询 | `apps/api/internal/handler/systemmonitoring.go:98-137` | `GET /api/system-monitoring/*` | 同步 | **无批量、无后台扫描、无 goroutine/ticker**。`status` 单行 envelope，store ping 有 **1s 超时**（`systemmonitoring.go:110`）；`errors` 资源 `ReadOnly: true`（`:89`）⇒ 不挂 batch-delete |
| `datadictionary` | **删除类型级联删条目** | FK `apps/api/modules/datadictionary/migration/migration.go:32`（`dict_key TEXT NOT NULL REFERENCES dict_types(key) ON DELETE CASCADE`；PG 变体 `:60`）；显式层 `apps/api/modules/datadictionary/store/repository.go:212-252` `DeleteTypeTx` | 经 `DELETE /api/data-dictionary/types/{id}` | 同步 | 先收集条目 id（`:215-232` `SELECT e.id FROM dict_entries e JOIN dict_types t ON e.dict_key = t.key WHERE t.id = ?` — **无界，该类型全部条目**），再 `:233` `DELETE FROM dict_types WHERE id = ?`。handler 包装 `dictionary.go:121-140` `DeleteTrashTx`（删除 + 回收站快照同一事务；级联条目 id 写入审计 detail `:134-138`） |
| `datapermission` | 策略/作用域读写 | `apps/api/internal/handler/datapermission.go:53/95/135/156` | `GET/PATCH /api/data-permission/*` | 同步 | 配置面，非批量 |

### 2.7 循环 / 长操作检索结论

- 面向选择的循环：`apps/api/internal/handler/resources.go:929`（`for _, id := range ids`，顺序删除回退路径）、
  `resources.go:886`（self-scope 逐 id `Get`，N+1）、`resources.go:908`（快照预读逐 id `Get`）、
  `resources.go:921`、`:975`（逐条 `OnWrite`）。
- 无界批量 SQL：
  - `DELETE FROM recycle_items WHERE restored_at IS NULL`
    （`apps/api/modules/recyclebin/store/repository.go:230`）——唯一一条无界 `DELETE` 端点。
  - `UPDATE notifications SET read_at=? WHERE user_id=? AND read_at IS NULL`
    （`apps/api/modules/authsession/notifications_repository.go:220`）——有 500 行上限兜底。
  - `DELETE FROM operation_log WHERE created_at < ?`（`apps/api/modules/operationlog/retention.go:66`）——**后台**，见 §2.9。
- 无界扫描：`SELECT id FROM wallet_accounts ORDER BY id`（`wallet/store/repository.go:841`）、
  `wallet_ledger_entries WHERE account_id = ?`（`:921-923`）、
  `SELECT ... FROM scheduled_tasks WHERE enabled = 1`（`scheduledtasks/store/repository.go:233-239`）、
  filelibrary `List`+`Stat`（`filelibrary.go:106-125`）。
- `... WHERE ... IN (...)` 动态占位符：`apps/api/modules/authsession/roles_repository.go:231-237`
  （`PermissionsForRoles`，角色数有界）。
- 事务内逐 id 循环：`users_repository.go:354-410`、`roles_repository.go:194-216`（均**单事务**，见 §1.2）。
- 逐行独立事务循环：`import.go:239-324`（N 行 = N 个事务，刻意 no-rollback）。

### 2.8 用户可见性矩阵（前端已接线 vs 仅 API）

| 操作 | 前端接线 | 证据 |
|------|---------|------|
| batch-delete | **仅 demo 页** | `apps/api/modules/dev/examples/schema/admin-list-batch.json:48,72-84`；该模块**仅在 `demo` profile**（`apps/api/kernel/profile.go:99-113`）。真实 `users.json`/`roles.json`/`scheduled-tasks.json`/`data-dictionary.json` **均无** `selection`/`batchMapping` |
| users 导出/导入 | **是** | `apps/api/modules/users/schema/users.json:446-463`（toolbar `exportUsers`/`openImport`）；导入 modal `users.json:119-129`（`submitImport` → `/api/import/users`）；导出 URL 映射 `apps/web/src/renderer/render.tsx:353-354`（`"export.users": "/api/export/users"`） |
| 活动日志导出 | **是** | `apps/web/src/components/activity-export.tsx:43`（调 `/api/operations/export`） |
| recycle-bin 清空 | **是** | `apps/api/modules/recyclebin/schema/recycle-bin.json:31-34`（`purgeAll` action）、toolbar `:147-154`（`"key": "purgeAll", "label": "Purge all"`）；行操作 restore/purge `:119-146` |
| scheduled-tasks 立即运行 | **是** | `apps/api/modules/scheduledtasks/schema/scheduled-tasks.json:289-303`（row action `run`，label "Run now"，`confirm: "Run this task now?"`）；action 定义 `:57-64` |
| notifications 全部已读 | **是** | `apps/web/src/components/notification-center.tsx:123`（`fetcher("/api/notifications/read-all", { method: "POST" })`）；schema `apps/api/modules/notifications/schema/notifications.json:16` |
| wallet reconcile | **否（不轮询）** | `apps/api/modules/wallet/schema/wallet.json:96-103` — `runReconcile` 为 `{"type":"request","method":"POST","url":"/api/wallet/reconcile","onSuccess":{"behavior":"reload"}}`；toolbar 项 `:494-499`。**202 body 被通用成功路径丢弃**；`apps/web/src` 对 `wallet/jobs`/`jobId`/`progress` **零命中** |
| filelibrary 下载/预览/复制链接 | **是** | `apps/web/src/renderer/render.tsx:352-363` |

> **关键缺口**：**没有任何前端轮询 `/api/wallet/jobs/*`**。
> 已交付的 Job 能力目前**只有测试消费者**在用（`docs/workspaces/workspace-038-.../R1-recon-I-038-002-batch-async-contract.md:131` 同结论）。
> ⇒ VP-038 的「结果中心」是**从零建 UI**，不是接现成前端。

### 2.9 后台周期任务全表（`apps/api`）

全仓库**只有 4 个周期循环**（`time.NewTicker(` 5 处命中，其中 1 处是 per-job 心跳）：

| # | file:line | 做什么 | 间隔 |
|---|-----------|--------|------|
| 1 | `apps/api/modules/scheduledtasks/scheduler.go:66`（goroutine `:59-63`，由 `apps/api/modules/scheduledtasks/provider.go:38` `scheduler.Start()` 启动） | Cron 调度器。`tick()` `:95-129`：`EnabledTasks()`（无界，`store/repository.go:233-239`）→ `ParseCron` → `fields.Matches(slot)` → 逐到期任务 `Execute`。按分钟槽去重（`:123-126`）；每日不可调度诊断（`:114-120`）。best-effort 单实例，错过不回补（`:3-5`） | **30s** — `scheduler.go:24` |
| 2 | `apps/api/modules/operationlog/retention.go:103`（goroutine `:82`，由 `apps/api/internal/composition/composition.go:1108-1117` 启动） | 审计日志保留期清理。每 tick 读站点设置策略（`retention.go:84`），`ApplyRetention`（`:23-74`）：可选归档 `INSERT INTO operation_log_archive ... SELECT ... WHERE created_at < ?`（`:36-64`），再 `DELETE FROM operation_log WHERE created_at < ?`（`:66`）——**无界、按龄全表**。天数/动作由管理员配置（`settings/migration/migration.go:103-106`：1..3650，`archive`\|`delete`），**从不硬编码**。启动后立即跑一次再按 tick（`:102,109-111`） | **1h** — `composition.go:1117` |
| 3 | `apps/api/internal/jobs/runner.go:216`（goroutine `:121` `go r.loop()`，由 `:111` `Start()`，`composition.go:1101` 调用） | Job 轮询循环。`ScanOnce`（`:185-212`）：回收已取消（`:187`）、耗尽超次（`:194`）、`ExpireDue`（`:201`）、`ListRunnable`（`:204`，LIMIT = `BatchSize`）、逐 job 派发（`:208-210`）→ 每 job 一 goroutine（`:247`）。`Submit`/`Cancel`/`Retry` 会提前唤醒（`:155,173,180`） | **10s** — `runner.go:51` |
| 4 | `apps/api/internal/channel/telegram/connection_manager.go:390` | Telegram 连接需求调和。每 tick `m.reconcileDemand(...)`（`:408`）按活跃租约启停长轮询接收器；`pruneExpiredLeasesLocked`（`:455`）也在租约获取（`:219`）/释放（`:451`）时跑 | **1s** — `connection_manager.go:18` |
| (4b) | `apps/api/internal/jobs/runner.go:291` | **非全局循环**——单个 running job 的租约心跳 + 取消轮询（`:293-318`） | 10s — `runner.go:50` |

**非周期 goroutine**（非调度器）：`composition.go:1139`、`apps/api/server/serve.go:237`（HTTP `srv.Serve`）、
`apps/api/internal/obs/server.go:80`（metrics listener）、`apps/api/internal/eventbus/memory.go:194`（关停 drain）。
**一次性启动任务**：品牌资源 GC `composition.go:471-474` `_ = brandAssets.GC([...])` → `raster_assets.go:395-415`。

**`apps/api` 中不存在**：`cron.New` / `robfig` / `gocron` / `time.AfterFunc`（零命中）；
**recycle-bin 保留期 sweeper**；邮件发件箱 sweeper goroutine（邮件保留在插入时内联强制，
`apps/api/internal/mail/outbox.go:95-109`，`DefaultOutboxCap = 500` `:39`）。

---

## 3. Web 侧批量 UI（WEB-SIDE BATCH UI）

### 3.1 关键结论：**生产页面 schema 没有任何批量动作**

全仓库 `batchMapping` / `requiresSelection` / `$selection.keys` 的 schema 级声明只有两处：

| 文件 | 性质 | 内容 |
|------|------|------|
| `apps/api/modules/dev/examples/schema/admin-list-batch.json:75-81` | **dev 范例**（`demo` profile only） | `requiresSelection: true`、`confirm`、`batchMapping.body.ids = "$selection.keys"`；action url `POST /api/users/batch-delete`（`:20`）；`onSuccess.behavior = "reload"`（`:22`） |
| `docs/schemas/component-registry.json:701,707,711` + `apps/web/src/protocol/upstream/*.cases.json` | **协议 schema / fixture** | 声明 `requiresSelection`、`batchMapping`（body 值可为 `$selection.keys`；query 可为 `$selection.count`） |

**生产页面 schema 清单（均无 `selection` / `batchMapping` / 批量动作）**：
`apps/api/modules/users/schema/users.json`（toolbar 仅 `openCreate`/`openInvites`/`exportUsers`/`openImport`，`:429-463`）、
`apps/api/modules/roles/schema/roles.json`、`apps/api/modules/datadictionary/schema/data-dictionary.json`、
`apps/api/modules/scheduledtasks/schema/scheduled-tasks.json`、`apps/api/modules/recyclebin/schema/recycle-bin.json`、
`apps/api/modules/wallet/schema/wallet.json`、`apps/api/modules/filelibrary/schema/file-library.json`。

> **推论（对 I-038-003 重要）**：`batch-delete` 的**后端 + 协议能力已交付**，
> 但**生产页面尚未有任何一处使用它**。因此「批量 UI 的既有分母 = 0 个生产页面」，
> 首波筛选若以「已有批量 UI 的页面」为分母会得到空集；分母应改为
> **「已挂载 batch-delete 路由的资源」= 4 个模块 / 5 条路由 / 5 个非只读 Resource ID**（§1.2；v0.3.1 更正），
> 或「具备 `table.selection` + `actions.batch.request` 能力声明的页面」。

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

> 确认弹窗、成功 reload、清选、失败保留选择的**逐行实现未核实**（U-02）。
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

### 3.4 异步 Job UI：**不存在**

- **没有任何前端轮询 `/api/wallet/jobs/*`**：`apps/web/src` 对 `wallet/jobs`/`jobId`/`progress` 零命中；
  `wallet.json:96-103` 的 `runReconcile` 是普通 request + `onSuccess.behavior: "reload"`，
  **202 响应体被丢弃**。
- **没有通用 Job 管理页 / 作业中心 / 进度条 / 轮询循环**。
  VP-038 的立项描述本身即把「通用 Job 管理页」列为 VP-012 的显式非目标
  （`docs/vision/plans/VP-038-batch-operations-and-job-center.md:39` 附近），与「当前不存在」一致。
- ⇒ 首波若要「可观察进度 + 结果中心」，**UI 是从零建**；后端 Job 面已完全就绪。

---

## 4. 候选筛选（CANDIDATE SCREENING）

| # | 操作 | 位置 | 当前 | 筛选 | 一句话理由 |
|---|------|------|------|------|-----------|
| 1 | `users` batch-delete | `resources.go:308` + `users.go:277` | 同步 | **keep-sync** | 单事务原子、4 KiB body 上限、批级 last-admin 守卫已交付且有回归测试；异步化即破坏既有合同 |
| 2 | `roles` batch-delete | `resources.go:308` + `roles.go:201` | 同步 | **keep-sync** | 同上（系统角色/占用守卫），原子且行数小 |
| 3 | `dict-types` / `dict-entries` batch-delete | `dictionary.go:284/304` | 同步（顺序回退） | **keep-sync（观察项）** | 行数小、交互式；但**非原子、可部分提交**——若要改进应补 `DeleteBatch` 原子实现，而非异步化 |
| 4 | `scheduled-tasks` batch-delete | `scheduledtasks.go:344` | 同步（顺序回退） | **keep-sync（观察项）** | 同上；任务定义表行数天然很小 |
| 5 | **导出 CSV（data-transfer）** | `export.go:44` | 同步 | **async-candidate（最强）** | 已在内存拼整个 CSV（`export.go:191-214`，**非流式、峰值 ≈2×**）、上限 10000 行、天然「产出物 + 可下载结果」形态；与 Job 的 `Result` + `ResultExpiresAt`（`model.go:50,58`）高度契合 |
| 6 | **导入 CSV** | `import.go:45` | 同步 | **async-candidate（最强）** | 已是「逐行处理 + 逐行错误报告 + no-rollback + 每行独立事务」形态（`import.go:1-6`、`:239-324`），异步化只需把 per-row 报告搬进 Job result；2 MiB 但**行数无上限**，且逐行 bcrypt（`:281`）使耗时随行数线性增长 |
| 7 | **recycle-bin `purge-all`** | `recyclebin.go:131` | 同步 | **async-candidate** | 唯一一条无界 `DELETE`（`repository.go:230`，无分页/无上限），不可逆，全表扫描；结果只需 `{"purged": n}` |
| 8 | **操作日志导出** | `operations_export.go:19-24` | 同步 | **async-candidate** | 与 #5 同一实现形态与同一 10000 上限；只做 #5 会造成两个导出面不一致 |
| 9 | **settings 重置（删全部品牌资源）** | `settings.go:41` / `:286` | 同步 | **async-candidate（弱）** | `DeleteAll` 是无界 `List`+逐 id `Delete` 循环（`raster_assets.go:384-388`）；但品牌资源实际数量很小，收益有限 |
| 10 | wallet reconcile | `wallet.go:357` | **已异步** | **out-of-scope（已是先例/模板）** | 202 + jobId + cancel/retry/result 已交付；是本波要复用的模板而非候选。**注意其进度只有硬编码 10→100**（`jobs.go:121-126`） |
| 11 | scheduled-tasks `POST /{id}/run` | `scheduledtasks.go:446` | 同步 | **async-candidate（弱）** | 用 `context.Background()` 执行、不随请求取消（`scheduler.go:167`），无 runId 返回；但 v1 只有 `system.noop`（`scheduler.go:46-50`），当前耗时 ≈ 0 |
| 12 | notifications `read-all` | `notifications.go:52` | 同步 | **keep-sync** | 单条 UPDATE + 每用户 500 行硬上限（`notifications_repository.go:21-22`） |
| 13 | wallet 代金券批量生成 | `wallet.go:484` | 同步 | **keep-sync** | `count > 1000` 直接拒绝（`voucher/service.go:41`），行数有界 |
| 14 | recycle-bin 单行 restore/purge | `recyclebin.go:109/149` | 同步 | **out-of-scope** | 单行、交互式、原子，非批量 |
| 15 | users enable/disable/unlock、MFA reset、invites、角色分配 | `users_state.go:43-45`、`mfa.go:358`、`invites.go:98-101`、`users.go:~240` | 同步 | **out-of-scope** | 全部单目标；**批量变体不存在**（属新功能，非本波承接对象） |
| 16 | roles 权限/菜单授权 | `roles_repository.go:357-402` | 同步 | **out-of-scope** | 循环的是单角色的权限 key，不是多角色 |
| 17 | filelibrary 上传/下载/删除/列表 | `filelibrary.go:181/218/255/106` | 同步 | **out-of-scope（列表扫描可另立观察项）** | 单文件、无批量、无 zip、无移动；但 `List`+`Stat` 无界扫描（`:106-125`）与 `quotaReached` O(files)（`upload.go:207-242`）是独立性能隐患 |
| 18 | operationlog 保留期清理 | `retention.go:103` | **后台周期（1h）** | **out-of-scope（但需进 Job 矩阵）** | 已是后台异步，非用户触发；但「Job 种类×作用域矩阵」应说明它与 `jobs` 表的关系（当前**不走** `jobs`，无进度/可见性） |
| 19 | scheduledtasks 调度器 tick | `scheduler.go:66` | 后台周期（30s） | **out-of-scope** | 同上 |
| 20 | Telegram 租约调和 | `connection_manager.go:390` | 后台周期（1s） | **out-of-scope** | 同上 |

### 4.1 Breaking-change 标记

| 操作 | 是否 breaking | 说明 |
|------|--------------|------|
| users / roles `batch-delete` | **改异步即 BREAKING** | 已交付合同：`200 {"deleted": n}` + 同步原子回滚（`resources.go:980`）。改成 `202 + jobId` 会破坏：① 前端 `onSuccess.behavior = "reload"` 依赖响应即终态；② 原子回滚语义（`users_batch_test.go:119-184` 断言 409 `LAST_ADMIN` 且**零删除**）在异步下不再可同步返回；③ 协议 fixture `request-construction.cases.json` 的 batch 用例。VP-038 已把它写成**显式非目标**（`docs/vision/plans/VP-038-batch-operations-and-job-center.md:51`：「把已交付的同步 `batch-delete` 改成 breaking 异步语义」）。**结论：保持同步，另立异步变体（新增端点），不改既有路由。** |
| 导出 / 导入 | **非 breaking** | 当前为同步 `GET /api/export/{resource}` / `POST /api/import/{resource}`。若**新增**异步端点并保留旧端点，属纯增量；若直接改旧端点返回 202 则 breaking。注意响应形态：导出是**裸 CSV 附件**（`export.go:210-214`），导入是 `{applied,failed,total,errors,fieldErrors}` envelope（`import.go:74-92`）——异步版的 result 必须能承载这两种形态。 |
| 操作日志导出 | **非 breaking（同上）** | `GET /api/operations/export`；前端 `activity-export.tsx:43` 直接触发下载，改成 202 会破坏该组件。 |
| `purge-all` | **非 breaking（需谨慎）** | 返回 `{"purged": n}`（`recyclebin.go:144`）；测试 `TestRecycleBinPurgeAll`（`recyclebin_test.go:365-405`）断言同步返回 `purged=2`。前端 `recycle-bin.json:31-34,147-154` 已接线。 |
| scheduled-tasks `POST /{id}/run` | **改异步即 BREAKING** | 当前返回 **204 No Content**（`scheduledtasks.go:463`），测试 `scheduledtasks_test.go:54`、`:133` 依赖 204；前端 `scheduled-tasks.json:289-303` 已接线。 |
| settings reset | 非 breaking（不建议改） | 前端已接线（`settings.go:41`）。 |
| notifications `read-all` | 非 breaking（不建议改） | 前端 `notification-center.tsx:123` 已接线。 |

---

## 5. 规模 / 时间上界证据（SIZE/TIME EVIDENCE）

### 5.1 请求体上限

| 常量 / 调用 | 值 | 位置 |
|-------------|-----|------|
| `maxResourceBodyBytes` | `4 << 10` = **4 KiB** | `apps/api/internal/handler/resources.go:37` |
| ↳ create / update / **batch-delete** | — | `resources.go:611` / `:696` / `:833` |
| ↳ import body（仅 `{"fileId":...}`） | — | `import.go:144` |
| ↳ account_self / email_identity / password_policy / service_credentials / settings / notifications / filelibrary | — | `account_self.go:142,240`；`email_identity.go:72,94`；`password_policy_settings.go:96`；`service_credentials.go:139`；`settings.go:181`；`notifications.go:220`；`filelibrary.go:267` |
| `maxImportBytes` | `2 << 20` = **2 MiB** | `apps/api/internal/handler/import.go:30`；检查 `import.go:179-182`（→ 413） |
| `maxUploadBytes` | `8 << 20` = **8 MiB**（单次 multipart） | `apps/api/internal/handler/upload.go:34`；应用 `upload.go:294`、`:305` |
| `UploadMaxFilesPerUser` 默认 | **1000 文件** | `apps/api/internal/config/config.go:485` |
| `UploadMaxBytesPerUser` 默认 | `256 << 20` = **256 MiB** | `apps/api/internal/config/config.go:486`（env `UPLOAD_MAX_BYTES_PER_USER` `:742`；装配 `composition.go:456`） |
| `BrandingMaxBytes` 默认 | `4 << 20` = **4 MiB** | `apps/api/internal/config/config.go:490`（env `:743`）；`apps/api/internal/handler/branding_assets.go:29-30,44`；强制 `raster_assets.go:186-193` |
| wallet reconcile body | 4 KiB | `apps/api/internal/handler/wallet.go:365` |
| wallet 代金券批量生成 body | 16 KiB | `apps/api/internal/handler/wallet.go:496` |
| auth 登录/改密 body | 4 KiB | `apps/api/internal/handler/auth.go:94,223,261` |
| Telegram webhook body | 1 MiB | `apps/api/internal/channel/telegram/webhook.go:142` |
| Telegram callback_data | 64 B | `apps/api/kernel/telegram.go:38-39` |

> **batch-delete 的 4 KiB 上限是硬约束**：`{"ids":[...]}` 在 4 KiB 内大约只能装
> **数百个 UUID 级 id**（未实测；按每 id 约 40 字节计 ≈ 100 个）。这是「同步批量天然有界」的直接证据，
> 也是异步批量必须换端点的原因之一。

### 5.2 行数 / 分页上限

| 常量 | 值 | 位置 |
|------|-----|------|
| `maxPageSize` | **100** | `apps/api/internal/handler/resources.go:38`；校验 `:434-437` |
| recycle-bin 列表页大小硬夹 | **100** | `apps/api/modules/recyclebin/store/repository.go:135-137` |
| `DefaultPageSize` | 20 | `apps/api/internal/handler/resources.go:40` |
| `maxExportRows` | **10000** | `apps/api/internal/handler/export.go:25`；校验 `:130-134`；注释：「full exports beyond this require additional filtering. Kept in sync with the frozen limit.」 |
| ↳ 操作日志导出复用 | — | `apps/api/internal/handler/operations_export.go:29-33` |
| `maxNotificationsPerUser` | **500** | `apps/api/modules/authsession/notifications_repository.go:21-22`；剪枝 `:103-104` |
| `DefaultOutboxCap`（邮件发件箱） | **500** | `apps/api/internal/mail/outbox.go:39`；强制 `:95-109` |
| `UploadMaxFilesPerUser` | **1000** | `apps/api/internal/config/config.go:485` |
| 代金券批量生成上限 | **1000** | `apps/api/modules/wallet/voucher/service.go:41`；handler `wallet.go:506` |
| `resourceIDRetries` | 3 | `apps/api/internal/handler/resources.go:42` |
| `DefaultMaxAttempts`（Job） | 3 | `apps/api/internal/jobs/model.go:27` |
| Job `BatchSize` | 32 | `apps/api/internal/jobs/runner.go:53` |
| Job `ResultTTL` | 24h | `apps/api/internal/jobs/runner.go:52` |
| `recycle_items` purge-all | **无上限** | `apps/api/modules/recyclebin/store/repository.go:230` |
| 导入行数 | **无上限**（仅 2 MiB 字节预算） | `apps/api/internal/handler/import.go:29-30` |
| wallet reconcile 账户数 / 流水数 | **无上限** | `apps/api/modules/wallet/store/repository.go:841`、`:921-923` |
| scheduled tasks 扫描 | **无上限** | `apps/api/modules/scheduledtasks/store/repository.go:233-239` |
| operation_log 保留期删除 | **无上限**（按龄全表） | `apps/api/modules/operationlog/retention.go:66` |

### 5.3 超时 / 限流

| 项 | 值 | 位置 |
|----|-----|------|
| HTTP `ReadTimeout` 默认 | **5s** | `apps/api/internal/config/config.go:463`；env `HTTP_READ_TIMEOUT` `:694` |
| HTTP `WriteTimeout` 默认 | **10s** | `apps/api/internal/config/config.go:464`；env `:695` |
| HTTP `IdleTimeout` 默认 | 60s | `apps/api/internal/config/config.go:465` |
| `HTTPShutdownTimeout` 默认 | 10s（drain 预算，≤0 fail-closed） | `apps/api/internal/config/config.go:468`、校验 `:709-711` |
| ↳ 服务器装配 | — | `apps/api/server/serve.go:203-206`；`apps/api/server/config.go:125-128` |
| system-monitoring store ping | **1s** | `apps/api/internal/handler/systemmonitoring.go:110` |
| Telegram 出站 HTTP | **10s** 严格预算 | `apps/api/internal/channel/telegram/http_sender.go:21` |
| SQLite `busy_timeout` | 5000 ms | `apps/api/internal/store/store.go:57`（DSN `_busy_timeout=5000`） |
| 调度器 tick | 30s | `apps/api/modules/scheduledtasks/scheduler.go:24` |
| 审计保留 sweeper tick | 1h | `apps/api/internal/composition/composition.go:1117` |
| Job scan tick | 10s | `apps/api/internal/jobs/runner.go:51` |
| Job 心跳 | 10s | `apps/api/internal/jobs/runner.go:50` |
| Job 租约 | 30s | `apps/api/internal/jobs/runner.go:48` |

> **重要**：`WriteTimeout = 10s`（`config.go:464`）意味着**任何同步请求超过 10s 都会被服务器写超时掐断**。
> 这为「同步批量/导出必须异步化」提供了硬性时间上界证据：
> 全量 reconcile、10000 行导出、无界 purge-all、无界 operation_log 保留删除都可能逼近或超过该预算。
> **注意**：scheduled-tasks 的手动触发用 `context.Background()`（`scheduler.go:167`），
> **不受**该超时约束（但客户端会先超时）——这是当前最不安全的同步长操作形态。

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
| `apps/api/internal/handler/data_transfer_test.go:326` `TestImportSizeLimit` | 3 MiB 上传（上传允许 8 MiB）→ 导入 **413**（钉住 2 MiB；注释 `:328` 同时钉住两个上限） |
| `apps/api/internal/handler/data_transfer_test.go:340` `TestImportMissingHeaderInvalidCsv` | 缺表头 → 400 `INVALID_CSV` |
| `apps/api/internal/handler/data_transfer_test.go:361` `TestImportRoleAssignmentBoundary` | `roles=admin` 逐行失败（非整请求 403） |
| `apps/api/internal/handler/data_transfer_test.go:398` `TestImportUnknownResource404` | 未知资源 → 404 |
| `apps/api/internal/handler/w16_batch_b_test.go:13,31` | 导入模板 + 字段错误 |
| `apps/api/internal/handler/operations_test.go:195` `TestOperationLogStructuredFiltersAndExport` | 操作日志导出复用同一 `maxExportRows` |
| `apps/web/src/components/activity-export.test.tsx:75` | 断言 `?pageSize=10000` |

### 6.3 recycle-bin 同步语义

| 文件 | 断言 |
|------|------|
| `apps/api/internal/handler/recyclebin_test.go:81` `TestRecycleBinRestoreAndPurge` | 恢复 200 / 物理删除 204 / 删除后详情 404（`:115-121`） |
| `apps/api/internal/handler/recyclebin_test.go:365` `TestRecycleBinPurgeAll` | 非 admin → **403**；admin → **200 `{"purged":2}`**（`:374-391`），之后列表 total=0（`:395-405`） |
| `apps/api/modules/recyclebin/store/repository_test.go:116` `TestPurge` | 单行物理删除 |
| `apps/api/modules/recyclebin/store/repository_test.go:132` `TestPurgeAllUnrestored` | 全量删除仅针对 `restored_at IS NULL` |
| `apps/api/modules/recyclebin/store/repository_test.go:163` `TestMarkRestoredTwiceFails` | 重复恢复失败 |
| `apps/api/modules/recyclebin/service_test.go:84` `TestRestoreDictTypeAndConflict` | 恢复 + 冲突 |
| `apps/api/modules/recyclebin/service_test.go:125` `TestRestoreTaskRoundTrip` | 任务恢复往返 |
| `apps/api/modules/recyclebin/service_test.go:173` `TestRestoreDictEntryRoundTrip` | 字典条目恢复往返 |
| `apps/api/modules/recyclebin/service_test.go:208` `TestRestoreOrphanDictEntryReturnsDomainError` | 孤儿条目恢复 → 领域错误 |
| `apps/api/modules/recyclebin/service_test.go:253` `TestRestoreAtomicityRollsBackOnFailedMark` | 恢复失败回滚 |
| `apps/api/modules/recyclebin/provider_test.go:99` `TestRecycleRealServiceRestoreConflictHTTP` | 恢复冲突 HTTP 面 |

### 6.4 wallet 异步 Job（作为模板的回归面）

| 文件 | 断言 |
|------|------|
| `apps/api/modules/wallet/jobs_test.go:124` `TestWalletJobServiceCompletesAtomicallyAndAudits` | job 达 `succeeded` 且 `Progress == 100`、`Result` 非空（`:138-141`）；result JSON `id == job.ID`、`result == "consistent"`（`:146`）；恰好 1 行 `wallet_reconciliation_runs` 且 `ID == job.ID`（`:149-152`）；跨 actor `Get` → `jobs.ErrNotFound`（`:153-155`）；审计链 `[wallet.reconcile.queued, wallet.reconcile]`（`:156-159`）；+25h 后 `expired` 且 result 空（`:160-164`） |
| `apps/api/modules/wallet/jobs_test.go:167` `TestWalletJobServiceCancelsQueuedWithoutBusinessRun` | 取消 queued → `StatusCancelled` 且**零** reconcile run（`:173-180`）；对其 `Retry` → `jobs.ErrNotRetryable`（`:181`）；第二事件 `wallet.reconcile.cancelled`（`:184-187`） |
| `apps/api/modules/wallet/jobs_test.go:190` `TestWalletJobServiceRollsBackConsumerFailureAndRetries` | 账户不存在 ⇒ `failed`，`ErrorCode == "JOB_HANDLER_FAILED"`，`Attempt == 1`，**零** run 行（消费者回滚）（`:196-203`）；`Retry` → `queued` attempt 1（`:204-207`）；二次失败 attempt 2（`:208-211`）；事件 `[queued, failed, failed]`（`:212-215`） |
| `apps/api/internal/handler/wallet_test.go:319` `TestWalletLifecycleAndAdjustFlow` | `POST /api/wallet/reconcile` → **202**（`:401-404`）；轮询 `GET /api/wallet/jobs/{id}` 至 `succeeded`（`:410-423`）；`GET …/result` → 200 且带 `Content-Disposition`、`result == "consistent"`（`:424-434`）；`cancel` → 409 `JOB_NOT_CANCELLABLE` / `retry` → 409 `JOB_NOT_RETRYABLE`（`:435-447`）；缺 job → 404 `JOB_NOT_FOUND`（`:448-452`）；强制过期 → 410 `JOB_RESULT_EXPIRED`（`:453-463`）；operationlog 含全部六个 wallet 事件（`:465-490`） |
| `apps/api/internal/handler/wallet_test.go:257` `TestWalletRoutesGates` | 17 条 admin 路由（含 4 条 job 路由）匿名 401 / editor 403（`:272-277`、`:301-306`） |
| `apps/api/internal/handler/wallet_test.go:680` `TestWalletReconcileBadBodyAndWriteGate` | 垃圾 body → 400（`:684-688`）；editor 在 reconcile/cancel/retry → 403（`:694-707`） |
| `apps/api/modules/wallet/store/repository_test.go:212` `TestReconcileConsistentAndInconsistent` | 全局 run（`accountID == ""`）、单账户 run、run 列表倒序 |
| `apps/api/modules/wallet/store/repository_test.go:263` `TestReconcileDetectsMismatch` | 篡改后 → `inconsistent` / `mismatchCount == 1` / details 含 `mismatches` |
| `apps/api/internal/jobs/runner_test.go:75,131,156,180,207,222,267`；`repository_test.go:40,94,147,187,212`；`shutdown_reclaim_test.go:19,61`；`runner_panic_test.go:14`；`runner_failure_test.go:12` | 通用 Job 六态 / 租约 / 取消 / 关停回收 / panic / 失败重试 |

---

## 7. 对首波筛选的直接影响（供 I-038-003 决策）

1. **分母必须先修正**：生产页面**零**批量动作（§3.1）。可选分母口径：
   （a）已挂载 `batch-delete` 路由的 **4 个模块 / 5 条路由 / 5 个非只读 Resource ID**（v0.3.1 更正）；
   （b）已具备 `table.selection` + `actions.batch.request` 能力声明的页面；
   （c）§2 中的**非批量类长操作**（导出/导入/purge-all/settings reset）。
2. **最强的 3 个异步候选**（§4 #5/#6/#7）：**导出**、**导入**、**recycle-bin purge-all**。
   三者都满足「长/行数无上界或上界很大 + 产出物可回看 + 与 Job 的 result/progress 天然契合」。
   **第 4 个应一并考虑**：#8 操作日志导出——它与 #5 共享同一实现形态与上限，
   只做一半会留下两个不一致的导出面。
3. **batch-delete 必须保持同步**，且 VP-038 已把它写成显式非目标
   （`docs/vision/plans/VP-038-batch-operations-and-job-center.md:51`）。若需异步删除，
   应**新增**端点而非改既有路由。
4. **同步路径的硬性时间上界是 `WriteTimeout = 10s`**（`config.go:464`）——
   这是论证「哪些操作必须异步」最有力的既有约束。
   例外：scheduled-tasks 手动触发用 `context.Background()`（`scheduler.go:167`），**不受**该约束。
5. **Job 运行时已通用**（`apps/api/internal/jobs`，六态 + progress + cancel + result TTL + 租约围栏，
   `jobs` 表迁移 42 已存在），首波无需新建 Job 基础设施，只需新增 `Kind` + handler（配方见 §2.3.2）。
6. **进度语义是既有缺口**：唯一先例 reconcile 只报硬编码 10 → 100（`jobs.go:121-126`），
   循环内不报进度。「进度」判据若要真正可观察，需在**新** handler 里实现细粒度 `reporter.Progress`。
7. **UI 是从零建**：前端**没有任何** Job 轮询/进度/结果中心（§3.4），
   `wallet.json:96-103` 甚至丢弃 202 响应体。「结果中心」是本波净新增工作量。
8. **Job 矩阵需覆盖 3 个既有后台周期任务**（§2.9 #1/#2/#4）：它们**不走** `jobs` 表，
   无进度、无可见性、无取消。I-038-003 的「Job 种类×作用域矩阵」应明确它们是 out-of-scope
   还是未来迁移对象。

---

## 待确认 / 未知

| # | 未确认项 | 影响 | 建议核实动作 | 状态 |
|---|---------|------|-------------|------|
| U-01 | 前端是否有任何 job 轮询 / 进度 UI | 首波「结果中心」UI 工作量 | 全量 grep `apps/web/src` 的 `job`/`poll`/`progress` | ✅ **已关闭**：`wallet/jobs`/`jobId`/`progress` 零命中；`wallet.json:96-103` 丢弃 202 body。**不存在**（§3.4） |
| U-02 | `render.tsx` 批量确认弹窗、成功 reload、清选、失败保留选择的逐行实现 | 异步化后 UX 兼容口径 | 通读 `apps/web/src/renderer/render.tsx` 批量相关段落 | ⬜ 未关闭（仅确认判定分支 `:331` 与 selection 传递 `:738`） |
| U-03 | wallet reconcile 每次运行实际触及的行数与耗时预期 | 判断它是否真是「大 N」先例 | 读 reconcile 服务实现 | ✅ **已关闭**：**无界**——全表账户扫描（`wallet/store/repository.go:841`）+ 单账户创世块全量重放（`:921-923`），无 LIMIT（§2.2.1） |
| U-04 | reconcile Job 的 payload / result 结构体与注册方式 | 复用模板的具体形状 | 读 `SubmitReconcile` 与 worker 注册点 | ✅ **已关闭**：payload `reconcileJobPayload{AccountID}`（`wallet/jobs.go:21-23`）；**result 无结构体**，是 ad-hoc `map[string]any`（`:214-221`）；注册 `RegisterWithTerminalHook`（`:39`）（§2.2.2） |
| U-05 | `datadictionary` 删除类型是否级联删除其条目 | batch-delete 规模估计 | 读 `dictTypeEntity.Delete` / store 层 | ✅ **已关闭**：**两层级联**——FK `ON DELETE CASCADE`（`datadictionary/migration/migration.go:32,60`）+ 显式 `DeleteTypeTx`（`store/repository.go:212-252`），级联无界但限于单类型条目（§2.6） |
| U-06 | 是否存在保留期自动清理 sweeper | purge-all 异步优先级 | 检索 `recycle_items` 定时清理 | ✅ **已关闭**：**recycle-bin 无 sweeper**（快照永久保留）。但存在**审计日志**保留 sweeper（`operationlog/retention.go:103`，1h）（§2.9） |
| U-07 | 各生产页面前端是否已暴露对应长操作按钮 | 「用户可见性」列 | 读 schema + `apps/web/src` | ✅ **已关闭**：见 §2.8 可见性矩阵（purge-all / run-now / read-all / 两个导出 / 导入均已接线；**wallet reconcile 不轮询**） |
| U-08 | `data-transfer` provider 路由注册行 | 证据完整性 | 读 `datatransfer/provider.go` | ✅ **已关闭**：`provider.go:44` 声明，`:55-64` 装配，`:66-67` 权限策略（§2.1） |
| U-09 | 是否存在 `DELETE ... WHERE id IN` 形式的多行删除 | 长操作完备性 | 全量 grep `DELETE FROM` | ⬜ 未关闭（已找到 3 条无界批量 DML，见 §2.7；未做全量 `DELETE FROM` 穷举） |
| U-10 | 导出「文件大小/耗时」的实测或文档预期 | 异步化收益量化 | 查 `docs/` data-transfer 决策文档 | ⬜ 未关闭（仅 10000 行上限 + 10s WriteTimeout；D-002 冻结文档给出上限口径但无实测耗时） |
| U-11 | `apps/api/internal/jobs` 是否已接入 Profile/模块可见性（`core.jobs` 迁移 no-op、`admin.jobs` 是否新建） | VP-038 的 Profile 边界（I-038-004 已裁决新建 `admin.jobs`） | 读 `kernel/profile.go` 与 composition 装配 | ⬜ 未关闭（本报告未核实；VP-038 `I-038-004` 用户 2026-09-19 已裁决新建 `admin.jobs` 进 admin 默认集） |
| U-12 | 4 个后台周期任务（scheduler / retention / jobs / telegram）是否应在 Job 矩阵中登记为可见作业 | I-038-003「Job 种类×作用域矩阵」的完备性 | 与 I-038-001 对齐口径 | ⬜ 未关闭（§2.9 已给出全表事实，口径待定） |
