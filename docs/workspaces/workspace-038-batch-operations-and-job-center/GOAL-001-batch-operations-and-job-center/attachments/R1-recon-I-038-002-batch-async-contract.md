---
title: R1 侦察 · I-038-002 批量异步契约（上游 pin / 同步 batch-delete / 异步先例 / A-B-C 选项）
status: draft
created: 2026-09-19
updated: 2026-09-19
parent: null
version: 0.1.0
---

# R1 侦察 · I-038-002 批量异步契约

> **性质**：只读侦察记录（read-only recon），为 `I-038-002`（批量异步契约 = 扩展 ADR-0022 vs 另立独立契约）提供 file:line 证据。
> **范围**：本仓 `apps/web/src/protocol/**`、`docs/schemas/**`、`apps/api/internal/handler/resources.go`、`apps/api/internal/handler/wallet.go`。
> **未改动任何代码**；本文件是本次唯一写入。
> **上游 pin 声明**：本仓协议固定点为 `schema-ui-docs` `v2.9.0` @ `81aa1d8954717f4ebdcc695eed6fafaeafcebe8d`（`apps/web/src/protocol/upstream/provenance-v2.9.json:3-4`）。

---

## 1. 上游协议面与「本地扩展 vs 上游变更」规则

### 1.1 钉住了什么（files + hashes/versions）

| 项 | 值 | 证据 |
|----|-----|------|
| 现行 pin 权威 | `schema-ui-docs` `v2.9.0` @ `81aa1d8954717f4ebdcc695eed6fafaeafcebe8d` | `apps/web/src/protocol/upstream/provenance-v2.9.json:2-4` |
| 工件数 | 10 件 `docs/schemas/*` + 1 件 `conformance/schemas/fixture-suite.schema.json` + 19 个 vendored fixture 套件 | `provenance-v2.9.json:6-127`；注 `:5` 声明「上游共 20 suites / 450 cases；scenarios 13 cases 未 vendor」 |
| 逐件 SHA-256 | 每个 `artifacts[].sha256`（如 `docs/schemas/component-registry.json` = `0714e133…d21a9`） | `provenance-v2.9.json:20-22` |
| 历史 pin（非现行） | v2.8 → `provenance-v2.8.json`；v2.7.0 @ `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b` → `provenance.json:1-4`（`note` 于 `:95` 明示「仅回归用；现行 pin 权威见 provenance-v2.9.json」） | `apps/web/src/protocol/README.md:13-15` |
| 机器校验点 | `stage3-fixtures.test.ts` 读取 `provenance-v2.9.json` 并逐件比对 SHA-256 | `apps/web/src/protocol/conformance/stage3-fixtures.test.ts:93`、`:115-149`；`SOURCE_COMMIT` 常量 `:95` |
| 哈希口径 | 以 LF 归一化后的字节计算（Git 可能签出 CRLF） | `stage3-fixtures.test.ts:75-78` |

**vendored 上游副本 vs 本仓本地文件**：

- **vendored 上游副本（不得手改）**：`docs/schemas/*.json`（10 件，含 `component-registry.json`）、`apps/web/src/protocol/upstream/*.cases.json`（19 套）、`docs/schemas/fixture-suite.schema.json`。它们的身份由 `provenance-v2.9.json` 的 `path` + `sha256` 定义。
- **本仓本地文件（可改，但不是上游）**：`apps/web/src/protocol/conformance/*.ts`（适配器，如 `request-construction.ts`）、`apps/web/src/protocol/*.guard.test.ts`、`apps/web/src/protocol/README.md`、`docs/vision/protocol-inventory-v2.7.0.md`（本地提取的清单，`status: active`，`source_tag: v2.7.0`，`:5-13`）、`apps/api/**`（后端实现）。
- **明确的「本地」判据**：`docs/vision/protocol-inventory-v2.7.0.md:26-27` 自述「**是**：**全量**能力/结构/fixture 清单（相对上游 pin）；**不是**：MVP 覆盖子集的冻结声明，也不是『已实现协议兼容』的证据」——即本仓清单是**消费侧映射**，不是上游权威。

### 1.2 规则：本地扩展 vs 上游变更（逐字引用）

- **本地扩展必须显式标注、且不得改写上游工件字节**：
  `docs/architecture/module-contribution-playbook.md:84`：
  > 「…因此**不增加 `schema-ui-docs@v2.9.0` pinned schema 的字段**。本地消费者可省略该扩展，正式跨项目协议发行需另行升级/审视，**不得把本地扩展写成上游兼容声明**。」
- **禁止把业务页面文档塞进上游 schema 目录**：
  `docs/architecture/module-contribution-playbook.md:113`：
  > 「在 `docs/schemas/` 放业务页面文档 → 该目录是上游协议 JSON Schema，不是业务页面库」
- **上游 pin 变更 = 重提取清单 + 升 version**：
  `docs/vision/protocol-inventory-v2.7.0.md:213`：「上游 pin 变更 → 重新提取本清单并升 `version`；同步 Charter/VP 协议引用。」
- **扩大协议范围须另立决策**：
  `docs/architecture/module-architecture.md:127`：「协议范围以 VP-003 的继承协议基线为准；**扩大范围须先有新的决策、递增覆盖版本与验证**。」
- **渲染层本地扩展先例（不改 pin）**：`apps/web/src/protocol/conformance/runtime-schema-validate.ts:40` 以校验期 overlay 叠加本地 `component` 键，注释写明「(local-only addition)」——上游 schema 字节保持 byte-identical。

**结论（对本目标直接相关）**：**新增一个本地端点 / 本地能力字符串，不需要、也不允许改动任何 pinned 文件**。pinned 面的改动等价于「上游协议变更」，须走上游发布 + 重 pin（`provenance-v2.9.json` 重签 + `stage3-fixtures.test.ts` 哈希更新）；本地新增只要求「显式标注为本地扩展 + 不写入上游兼容声明」。

---

## 2. ADR-0022 批量契约（现状 = 同步 `batch-delete`）

### 2.1 后端

| 项 | 证据 |
|----|------|
| `BatchDeleter` 接口 | `apps/api/internal/handler/resources.go:81-83`：`type BatchDeleter interface { DeleteBatch(ids []string, user account.User) (int, error) }` |
| 接口语义注释 | `resources.go:77-80`：「optional atomic whole-batch delete capability (ADR-0022 D5d · D-001 P0)… commits or rolls back as one unit」 |
| 路由注册 | `resources.go:308`：`add("POST", res.Path+"/batch-delete", a.Middleware(h.batchDelete()))`（仅在 `!res.ReadOnly` 分支内，`:302-309`） |
| handler | `resources.go:824` `func (h *resourceHandler) batchDelete() http.Handler` |
| 请求体 | `resources.go:830-832` `struct { IDs []any \`json:"ids"\` }`；空 → 400 `EMPTY_SELECTION`（`:838-841`）；非标量 → 400 `INVALID_SELECTION_KEY`（`:850-863`） |
| 权限 | `resources.go:826` `requirePermission(w, r, h.writePerm)` |
| 成功响应 | `resources.go:893` / `:971` / `:980`：`writeJSON(w, http.StatusOK, map[string]any{"deleted": deleted})` |
| 实现 `DeleteBatch` 的实体 | 仅 2 个：`apps/api/internal/handler/users.go:277`（`func (e *usersEntity) DeleteBatch(...)`）、`apps/api/internal/handler/roles.go:201`（`func (e *rolesEntity) DeleteBatch(...)`）；其余实体走顺序删除回退分支 `resources.go:927-972` |

### 2.2 前端

| 项 | 证据 |
|----|------|
| 构造器 | `apps/web/src/protocol/conformance/request-construction.ts:642` `function buildBatchRequest(input: JsonObject): RequestConstructionResult`（分发点 `:765-766` `case "batchRequest"`） |
| 选择集归一 | `request-construction.ts:624` `export function normalizeSelection(keys: unknown[]): { keys: unknown[]; count: number }`（D3：仅标量、保序去重、`count = keys.length`，`:623` 注释） |
| 空选择拒绝 | `request-construction.ts:662-664` → `fail("EMPTY_SELECTION", "selection")` |
| `$selection.keys` 仅 body | `request-construction.ts:683-684` → `fail("SELECTION_KEYS_BODY_ONLY", …)`；`$selection.count` 可 query 可 body（`:734-736`） |
| 返回值 | `request-construction.ts:707-711`：`{ ok: true, request: { method, url, body }, selectionAfterSuccessReload: { keys: [], count: 0 } }` |
| 执行 | `apps/web/src/renderer/render.tsx:685` `async function runBatchRequest(...)`；构造入参 `:734-739`；`fetch` `:757-761`（`Content-Type: application/json`）；非 2xx → `readResourceApiError` 映射 `:765-777`；成功 `:778` `return { ok: true }` |
| 触发 | `render.tsx:1213` `invokeBatchAction`（D4：confirm → run）；成功后 `:1243` `batchSuccessMessageFor` + `:1244` `reloadList()` |
| 成功文案 | `render.tsx:329-335`：URL 以 `/batch-delete` 结尾 → `t("feedback.itemsDeleted")` |
| 选择清空 | `render.tsx:950-954`「any data reload success clears every table selection」→ `setSelections({})` |
| 表格侧 | `apps/web/src/renderer/schema-table.tsx:932-935` `selectionEnabled`（要求 `props.selection.mode === "multiple"`）；`:1329-1354` 批量触发判定（`batchMapping !== undefined \|\| requiresSelection === true`）与禁用逻辑 |

### 2.3 请求/响应形状（5 行摘要）

```
请求:  POST {action.url}                     # 例 /api/users/batch-delete，仅 protocol-relative path
       Content-Type: application/json
       body: { <batchMapping.body 键>: $selection.keys 数组 | $selection.count 标量 | 字面量 }
       query: 仅字面量 + $selection.count（$selection.keys 进 query → SELECTION_KEYS_BODY_ONLY）
       path:  仅字面量绑定（$ 前缀值 → INVALID_MAPPING_VALUE）
响应:  200 {"deleted": n}   ← 前端不读该字段，只判 response.ok；成功后 reloadList() 并清空选择
错误:  400 EMPTY_SELECTION / INVALID_SELECTION_KEY / INVALID_BODY；403 FORBIDDEN；500 INTERNAL
```

---

## 3. 异步先例：wallet reconcile（202 + jobId + 轮询）

### 3.1 后端路由（`apps/api/internal/handler/wallet.go`）

| 方法/路径 | 行 | 权限 | 返回 |
|-----------|----|------|------|
| `POST /api/wallet/reconcile` | `:357` | `wallet.write` | **`writeJSON(w, http.StatusAccepted, walletJobToMap(*job))` — `:384`（202 已确认）** |
| `GET /api/wallet/jobs/{id}` | `:387` | `wallet.read` | 200 job 投影（`:397`） |
| `POST /api/wallet/jobs/{id}/cancel` | `:400` | `wallet.write` | 200 job 投影（`:410`） |
| `POST /api/wallet/jobs/{id}/retry` | `:413` | `wallet.write` | 200 job 投影（`:423`） |
| `GET /api/wallet/jobs/{id}/result` | `:426` | `wallet.read` | queued/running → **409 `JOB_RESULT_NOT_READY`**（`:437-438`）；expired → **410 `JOB_RESULT_EXPIRED`**（`:439-440`）；failed/cancelled → 200 job 投影（`:441-442`）；succeeded → 200 + `Content-Disposition: attachment; filename="wallet-reconcile-{id}.json"` 原始结果字节（`:443-447`） |

### 3.2 响应 JSON（`walletJobToMap`，`wallet.go:989-1006`）

```json
{
  "id": "job-…", "kind": "wallet.reconcile", "status": "queued|running|succeeded|failed|cancelled|expired",
  "progress": 0, "attempt": 0, "maxAttempts": 3, "cancelRequested": false,
  "createdAt": "…Z", "updatedAt": "…Z",
  "error": { "code": "…", "message": "…" },      // 仅当 errorCode 非空
  "finishedAt": "…Z",                             // 仅当已终态
  "resultUrl": "/api/wallet/jobs/{id}/result"     // 仅当 status == succeeded (:1002-1004)
}
```

Job 种类常量：`apps/api/modules/wallet/jobs.go:19` `const ReconcileJobKind = "wallet.reconcile"`；`SubmitReconcile` 于 `:45`，`Get` 带 `ExpireIfDue`（`:86-91`），`Cancel`/`Retry` 于 `:95`/`:109`。模块路由声明：`apps/api/modules/wallet/provider.go:273-274`；Profile 贡献键：`apps/api/kernel/profile.go:202`。

### 3.3 Web 侧驱动方式（关键发现）

- **没有前端轮询实现**。全仓 `apps/web/src` 搜 `wallet/jobs` / `jobId` / `resultUrl` / `progress` / `cancelRequested` **零命中**。
- 页面只声明一个**同步 `request` action**：`apps/api/modules/wallet/schema/wallet.json:96-103` `"runReconcile": { "type": "request", "method": "POST", "url": "/api/wallet/reconcile", "onSuccess": { "behavior": "reload" } }`；toolbar 入口 `:495-499`（`actionRef: "runReconcile"`）。
- 即：202 响应体被 `request` action 的通用成功路径忽略，用户看不到 jobId/进度；`/api/wallet/reconcile/runs` 与 `/api/wallet/jobs/*` 目前**只有 API 消费者**（测试：`apps/api/internal/handler/wallet_test.go:413,425,439-440,449,460`）。
- i18n 仅有 `schema.wallet.toolbar.reconcile`（`apps/web/src/i18n/messages/zh-CN.json:858`「对账」/ `en-US.json:858`「Reconcile」），无 job/进度相关文案。
- **真正的轮询先例在别处**：`apps/web/src/components/monitoring-auto-refresh.tsx:17-22`（5/10/30 秒选项）、`:36-38` `window.setInterval(() => crud?.refreshList("/api/system-monitoring/status"), intervalMs)`——自定义组件 + `refreshList` 的定向刷新，可作为作业中心轮询的形态先例。

---

## 4. A / B / C 决策选项（基于以上证据）

### 选项 A — 扩展 ADR-0022 批量语义，新增异步变体

- **做法**：让 `batchRequest` 支持「202 + jobId」响应分支（前端识别 202 并转入进度视图），或新增 `$selection` 之外的异步声明位。
- **触及文件**：`request-construction.ts`（batch 分支语义）、`render.tsx`（`runBatchRequest` 响应分支）、`schema-table.tsx`（触发态）；schema 侧需要新声明位 → `docs/schemas/action.schema.json` / `node.schema.json`。
- **协议 pin 影响**：**高**。上游 `action.schema.json` 为 `additionalProperties: false` 的严格校验（先例见 `docs/workspaces/.../GOAL-004-r2-f02-data-import-export/02-execution/E-003-s2-s3-implementation-verification.md:27`：`onSuccess.behavior` 枚举严格导致本地 `download` 行为无法过结构校验，被迫改用 `CustomAction` 扩展点）；新增字段 = 改 pinned 工件 = 上游协议变更，须重 pin（`provenance-v2.9.json` + `stage3-fixtures.test.ts:115-149` 哈希）。
- **回归风险**：**中高**——直接改写唯一已交付的批量语义路径（users/roles/data-dictionary/scheduled-tasks 四个页面共用）。
- **守卫测试**：`capability-declaration.guard.test.ts:57` 的 `actions.batch.request` 标记是 `/\/batch-delete\//`，语义扩展后需同步改标记；`stage3-fixtures.test.ts` 的 `request-construction` 套件（19 个 `batch-request-*` 用例，`request-construction.cases.json:1022-1569`）需保持全绿或显式排除。

### 选项 B — 保持 ADR-0022 同步语义不动，另立本地模块自有异步契约

- **做法**：`admin.jobs` 模块（`I-038-004` 已裁决方案 A）提供本地端点，例如 `POST /api/jobs/{kind}` → **202 + jobId**，`GET /api/jobs/{id}` 轮询，`GET /api/jobs/{id}/result` 取结果；页面用**既有 `request` action** 提交 + 一个结果中心页面（`admin.jobs` 自有 schema）。
- **触及文件**：`apps/api/modules/jobs/**`（+ `kernel/profile.go` 模块矩阵内容扩展）、新模块 schema JSON；web 侧**零协议代码改动**（或仅新增一个自定义组件用于轮询，形态照 `monitoring-auto-refresh.tsx`）。
- **协议 pin 影响**：**零**——不新增 capability 字符串、不改 `docs/schemas/**`、不改 `upstream/*.cases.json`。
- **回归风险**：**低**——`batchDelete()` 路径完全不碰。
- **守卫测试**：无需改动；若新页面声明 `actions.batch.request` 或 `table.selection`，则由 `capability-declaration.guard.test.ts:57/:51` 的既有标记自然覆盖（文本标记 `/batch-delete/`、`"requiresSelection"`）。**风险点**：若新 capability 字符串（如 `actions.batch.async`）被写进 `meta.requiredCapabilities`，`capability-declaration.guard.test.ts:110-114` 会以 `no usage marker defined for capability "…"` **直接失败**——必须同时在该 guard 的 `MARKERS` 表登记标记，或改用既有 capability。

### 选项 C — 混合：ADR-0022 语义不变，异步能力以本地扩展位承载

- **做法**：同步 `batch-delete` 完全冻结；长操作走本地模块端点（同 B），但**由批量工具栏以本地扩展属性触发**（如 `batchMapping` 之外的本地键，或 `actionRef` 指向本地 action），从而复用选择集与确认弹窗。
- **触及文件**：`schema-table.tsx`（本地扩展键放行）、`render.tsx`（本地分支）、模块端点；`docs/architecture/module-contribution-playbook.md` 需补本地扩展标注。
- **协议 pin 影响**：**零到低**——只要本地键走 `props`（`props` 无 `additionalProperties: false`，见 `docs/workspaces/.../GOAL-018-mfa-manager-ui/03-audit/A-003-s5-independent.md:80` 的裁决），不改 pinned schema。
- **回归风险**：**低到中**——需证明新增分支不影响既有 batch 判定（`schema-table.tsx:1329-1354` 的 `isBatch` 判定顺序）。
- **守卫测试**：若沿用 `table.selection` + `requiresSelection`，既有标记已覆盖；若引入新 capability 字符串，同 B 的 `MARKERS` 表缺口问题。

### 一句话取舍

| 选项 | pin 影响 | 前端改动 | 同步 batch-delete 回归风险 | 守卫改动 |
|------|----------|----------|---------------------------|----------|
| A | 高（改 pinned schema，须上游发布/重 pin） | 大（改批量核心路径） | 中高 | 需改 `MARKERS` + fixture 口径 |
| B | 零 | 小（新页面 + 可选自定义轮询组件） | 低 | 无（除非新增 capability 字符串） |
| C | 零到低 | 中（复用选择集，本地扩展键） | 低到中 | 无（除非新增 capability 字符串） |

---

## 待确认 / 未知

1. **上游 ADR-0022 是否已有异步语义预留**：本仓未见 `actions.batch.request` 的 202/异步定义；上游 `schema-ui-docs@v2.9.0` 的 ADR 全文未在本仓 vendored，**无法从本仓断言上游无此规划**——需在 R1 冻结前查上游 `docs/decisions/0022-*` 原文。
2. **`admin.jobs` 的作业作用域**：`wallet.reconcile` 走 `GetForActor`（actor 作用域，`jobs.go:82`）；管理作用域的 Job 列表是否复用同一 repository 尚未确认（对应 `I-038-001`，本轮未扫描 `internal/jobs` 全量）。
3. **202 响应在前端 `request` action 路径下的行为**：`runBatchRequest` 只看 `response.ok`（`render.tsx:765`），202 属 ok；但通用 `request` action（非 batch）是否读 body 未在本轮核对——影响选项 B 是否能「零前端改动」拿到 jobId。
4. **`capability-registry.json` 是否允许本地 capability**：`docs/schemas/capability-registry.json` 为 pinned 工件（`provenance-v2.9.json:16-18`），本仓是否已有「本地 capability 登记表」未确认；若无，选项 B/C 应避免引入新 capability 字符串。
5. **首波纳入异步的批量操作分母**（`I-038-003`）本轮未盘点（`data-transfer` 导出/导入、`scheduled-tasks` 批量面）。
6. **轮询间隔与终态语义**：现有先例只有 `monitoring-auto-refresh.tsx` 的 5/10/30 秒；作业中心应采用何种间隔/退避、以及 409/410 在 UI 上的呈现未定。
