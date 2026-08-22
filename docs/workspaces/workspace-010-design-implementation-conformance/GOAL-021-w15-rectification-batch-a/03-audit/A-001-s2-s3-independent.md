---
id: GOAL-021-w15-rectification-batch-a
doc: audit-entry
record_id: A-001
source: independent
scope: S2/S3 implementation of W15-F01, F02, F04, F05, F07 against D-001 freeze + GOAL-020 D-002
verdict: conditional
status: recorded
auditor: grok-build (grok-4.6 · reasoning high)
parent: GOAL-020-w15-user-perspective-findings
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-001 · S2/S3 独立交叉审计（2026-08-17）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：execution-facts · S2/S3 实施与回归（W15-F01 / F02 / F04 / F05 / F07）对照本目标 D-001 冻结 + [GOAL-020 D-002](../../GOAL-020-w15-user-perspective-findings/01-decision/D-002-i001-user-adjudication.md)。**不是**关门审计，**不是** S4。
- **verdict**：**conditional**
- **工作区**：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；`canonical_scope` 已核对；`shared_materials_catalog: none`；`plan_refs` + `primary_plan` = `VP-010-design-implementation-conformance`）

## 范围与区间

- **covered**：GOAL-021 五件套（`00-meta` / `01-decision` + D-001 / `02-execution` + E-001/E-002 / `03-audit` 空索引）；`workspace.md` / `goal-tree.md` 绑定；GOAL-020 D-002 批 A 范围句；下列 shipped 代码与驱动测试：
  - F01：`apps/web/src/account/auth-client.ts` `doRefresh` + `auth-client.test.ts` `restoreSession` / `authFetch`
  - F02：`data-table.tsx` error+retry、`schema-table.tsx` `retryNonce`、对应测试
  - F04：`handler/route_envelope.go` `WithJSONRouteErrors` + `route_envelope_test.go`
  - F05：`server.WrapSecurity` + `server_test.go`；`config.HTTPCORSOrigins`
  - F07：`auth.go` refresh、`errorcatalog`、`error_contract_test` / `auth_test.go`
- **不 covered**：S4 自审与关门；批 B/C；浏览器活栈；全量 `go test ./...` / 全量 vitest 复跑（本轮只复跑驱动测试 + 一次 HEAD 探测）；其他工作区正文。
- **本意见不修改** `00-meta` status / 检查点 / `progress`、`goal-tree.md` / `workspace.md`、方案正文、或任何应用代码。

## 成果（有证据）

### 工作区与信息门禁

| 字段 | 核对 |
|------|------|
| `workspace_id` | `workspace-010-design-implementation-conformance` |
| Root | `GOAL-001-design-implementation-conformance`（`parent: null`） |
| canonical | `docs/workspaces/workspace-010-design-implementation-conformance/` |
| 共享资料 | `none`；无资料行被当作事实或关闭证据 |
| GOAL-021 归属 | `parent: GOAL-020-w15-user-perspective-findings`；树内 active · 3/4 |
| GOAL-020 D-002 | 批 A = F01/F02/F04/F05/F07；F05 留本区；全部 in-scope |
| I-001 | non-blocking · **closed**（D-001）；无到期且影响本 scope 的 required 信息项 |

### W15-F01 · `doRefresh` 层（部分达成）

- `doRefresh` 对 `postJSON` throw 只 `return false`，不再 `clearTokens()`：`apps/web/src/account/auth-client.ts:129-135`。
- `!response.ok` 仅 `401/403` 清会话；5xx 保留 Token 并 `return false`：同文件 `:141-153`。
- `restoreSession` 在 refresh 失败时返回 `{ kind: "reauth" }` 且 **不再二次清 Token**：`:379-386`。
- 测试经真实 `restoreSession`：网络 throw / 500 保留双 Token；401 清空：`auth-client.test.ts:238-266`。本轮 `npx vitest run` 该文件绿。

### W15-F02 · 表格重试（达成）

- `DataTable` 错误态 `onRetry` + `t("feedback.retry")`：`data-table.tsx:52-53,262-278`；i18n `feedback.retry` / `feedback.resourceFetchFailed` 中英均在。
- `schema-table` 以本地 `retryNonce` 接入既有 `useEffect` refetch（D-001 允许的「等价 refetch」，非必须 bump `crud.reloadToken`）：`schema-table.tsx:471,527,864-865`；错误主文案改为 `t("feedback.resourceFetchFailed")`，不再把 `err.message` 堆到表格。
- 测试：`data-table.test.tsx:144-159` 按钮回调；`schema-table.test.tsx:192-214` 点击后第二次 fetch 成功并渲染行。本轮 3 个 web 文件 63/63 绿。

### W15-F04 · JSON 信封主路径（核心达成，HEAD 见 F-002）

- `WithJSONRouteErrors` 包装真实 `ServeMux`：未知路径 `404` + `NOT_FOUND`；已注册路径的错误方法 `405` + `METHOD_NOT_ALLOWED`：`route_envelope.go:10-24`。
- 生产装配：`composition.go:466-467` `newServer` → `handler.WithJSONRouteErrors(mux)` → `server.New`（Fx `:119`）。
- 目录 + 契约：`errorcatalog.go:29-30`；`error_contract_test.go:27` 冻结；`TestErrorCodeContractPinnedSet` / `TestErrorCatalogCoversFrozenCodesExceptInternal` 本轮绿。
- 测试打真实 mux：`route_envelope_test.go:11-61`（未知路径 / POST `/api/health` / 注册 GET 仍 200）。信封含 `messageKey`（catalog `error.notFound` / `error.methodNotAllowed`，中英 i18n 已齐）。

### W15-F05 · CORS + nosniff（达成）

- `WrapSecurity` 全局 `X-Content-Type-Options: nosniff`；白名单 Origin 才写 CORS，OPTIONS 204：`server.go:23,32-53`。
- 配置：`HTTPCORSOrigins`；YAML `http.cors_origins`（snake_case，与仓内其它 HTTP 键一致）；env `HTTP_CORS_ORIGINS` 覆盖；`splitCSV`；空 = 不发 CORS：`config.go:34-36,94,229-230,288-290`；`config.default.yaml:36-38`。
- 测试：`server_test.go:32-70` 允许 Origin 预检 204 + nosniff；外来 Origin 无 CORS 仍 nosniff。`go test ./internal/server ./internal/config` 本轮绿。
- 不改 Profile / 模块矩阵 / Manifest。与 D-002「F05 留本区」一致。

### W15-F07 · refresh 独立错误码（达成）

- `authHandler.refresh` 对 `ErrInvalidToken` / `ErrExpiredToken` / `ErrTokenRevoked` 写 `REFRESH_TOKEN_EXPIRED`：`auth.go:182-184`。
- 登录密码错误仍 `UNAUTHORIZED`：`auth.go:132-137`；`auth_test.go:41-42,52-53,271-272`。
- 目录文案与 D-001 字面一致：「登录已过期，请重新登录」/ `sign-in expired; please sign in again`：`errorcatalog.go:28`；前端 `error.refreshTokenExpired` 中英对齐。
- 基础设施失败仍走 `REFRESH_FAILED`（5xx），未误并入过期码：`auth.go:186-188`。
- `TestAuthRefreshInvalid` 断言 bogus token → `REFRESH_TOKEN_EXPIRED`：`auth_test.go:167-176`。

### 回归与 checkpoint（部分可核）

- git checkpoint **存在且为当前 HEAD**：`285c7e810881d80b3655ce2b07ba6edd614aa6f2`（`feat(workspace-010): W15 batch A F01/F02/F04/F05/F07 remediations`；owned 代码 + GOAL-020/021 文档；35 files）。与 E-002 声明一致。
- 本轮复跑（非全量）：handler 契约/信封/refresh、server/config、vitest 三文件 **全部 exit 0**。
- E-002 点名的 scratch `go-test-1.log` / `go-test-2.log` / `vitest-1.log` / `vitest-2.log` **在仓库与常见 Temp 路径未找到**（见 recommended F-004）。不据此否定上述定向绿，也不把 commit message 的「1046/1046 两遍」当作独立复跑证据。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| F01：`doRefresh` 仅 401/403 清 Token；网络/5xx 保留并允许重试；测试走 `restoreSession` / `authFetch` | **部分** | `doRefresh` + `restoreSession` 测试成立；`authFetch` 在 refresh `false` 时仍 `clearTokens()`（F-001） |
| F02：错误态可重试 + 友好 i18n 主文案 | **达成** | DataTable + schema-table + 两层测试 |
| F04：未匹配 404 JSON / 方法不允许 405 JSON + catalog | **部分** | GET 未知路径 / POST-on-GET 成立；HEAD-on-GET 被误判 405（F-002） |
| F05：nosniff + 可配置 CORS 白名单；空 = 不发 CORS | **达成** | `WrapSecurity` + YAML/env；键名与 D-001 字面不一致见 F-003 |
| F07：`REFRESH_TOKEN_EXPIRED`；登录仍 `UNAUTHORIZED`；契约冻结 | **达成** | handler + catalog + contract + auth_test |
| S3 全量两遍可独立复核 | **证据不足** | 定向测试绿；具名 scratch 日志缺失（F-004） |
| I-001 / 资料引用 | **不阻断** | I-001 closed non-blocking；catalog `none` |

## Findings

### F-001 · `authFetch` 在 refresh 非凭证失败时仍永久清会话

| 字段 | 值 |
|------|-----|
| 严重度 | **high** |
| 建议 | **required** |
| 状态 | open |
| 关联 | W15-F01；D-001「网络/5xx 保留 Token 并允许重试」；D-001「测试走真实 restoreSession / **authFetch** 入口」 |

- **描述**：D-001 把容灾写在 `doRefresh`（throw / 非 401-403 不清 Token、`return false` 以便重试）。`restoreSession` 遵守该约定。但中会话几乎所有 API 都走 `authFetch`（`main.tsx` resourceFetcher、`fetchMe` 等）。`authFetch` 在业务 401 后调用 `refreshAccess()`，**只要 refresh 返回 false 就无条件 `clearTokens()` + `onAuthLost`**，不论失败原因是 401 凭证还是网络 throw / 5xx：

```196:211:apps/web/src/account/auth-client.ts
export async function authFetch(input: RequestInfo | URL, init?: RequestInit): Promise<Response> {
  let response = await fetch(input, withAuth(init));
  if (response.status === 401 && !isAuthEndpoint(input)) {
    const refreshed = await refreshAccess();
    if (refreshed) {
      // ...
    } else {
      clearTokens();
      onAuthLost?.();
    }
  }
  return response;
}
```

  这会抵消 `doRefresh` 刚保住的 Token：典型路径是 access 过期 → 业务 401 → refresh 遇弱网/5xx/`doRefresh` 返回 false → 本地会话仍被抹掉 → `AuthContext` `onAuthLost` 进 `reauth-required`。W15-F01 原题「网络偶发抖动时刷新 Token 会永久抹除本地会话」在**开页**路径已缓解，在**使用中**路径仍成立。
- **测试缺口**：D-001 要求 `authFetch` 入口。现有 `authFetch notifies auth loss when the refresh fails`（`auth-client.test.ts:157-170`）只打 refresh **401**，不断言 Token，也不覆盖 throw/500 应保留。没有与 `restoreSession keeps tokens when refresh throws/500` 对等的 `authFetch` 用例。
- **相关（不另开 required）**：`restoreSession` 在网络/5xx 时仍返回 `{ kind: "reauth" }`（测试也如此断言），boot 会进 reauth 屏；Token 还在，用户 reload 或再进 `/login` 后若网络恢复可再旋转。这是 UX 残余，不是二次清空。见 F-005。
- **建议修复**：`authFetch` 仅在 refresh 因 **401/403**（或明确凭证/撤销）失败时清会话并 `onAuthLost`；网络/5xx 应保留 Token、把原 401（或 refresh 失败）返回给调用方，允许重试。补 `authFetch` 网络 throw / 500 保 Token、401 仍清 Token 的测试。

### F-002 · `WithJSONRouteErrors` 把 GET 路由上的 HEAD 误判为 405

| 字段 | 值 |
|------|-----|
| 严重度 | **med** |
| 建议 | **required** |
| 状态 | open |
| 关联 | W15-F04；D-001「方法不允许 → 405」 |

- **描述**：Go 1.22+ `ServeMux` 对已注册 `GET` 自动服务 `HEAD`。包装器在 `mux.Handler(r)` 得到 `GET /api/health` 后，因 pattern 方法 `GET != HEAD` 且 `pathHasOtherMethod` 为真，写成 `405 METHOD_NOT_ALLOWED`。
- **证据（本轮独立探测，探测文件未留在树内）**：同一 mux 注册 `GET /api/health` 后：
  - 裸 mux `HEAD /api/health` → **200**
  - `WithJSONRouteErrors` `HEAD /api/health` → **405** JSON `METHOD_NOT_ALLOWED`（含 `messageKey`）
- **影响**：F04 目标是把「未匹配 / 真·方法不允许」收成 JSON 信封，不是收窄合法 HEAD。负载均衡、健康检查、部分客户端的 HEAD 探测会从 200 变成 405。现有 `route_envelope_test.go` 只覆盖未知路径与 POST-on-GET，未锁 HEAD。
- **建议修复**：HEAD 在存在 GET 注册时交给 mux（或显式等价处理）；补回归：`HEAD` 已注册 GET 路径 = 200（或与裸 mux 一致），真正未注册方法仍 405。

### F-003 · D-001 CORS YAML 键名与 as-built / 操作员文件不一致

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| 状态 | open |

- D-001 写 YAML `http.corsOrigins`；实现与 E-002、`config.default.yaml` 为 `http.cors_origins`（`yaml:"cors_origins"`）。仓内 YAML 使用 KnownFields：若按冻结字面写 `corsOrigins` 会启动失败。
- 操作员副本 `apps/api/configs/config.yaml` 的 `http:` 段**没有**该键或注释；空缺行为正确（不发 CORS），但可发现性差。
- 功能本身（env `HTTP_CORS_ORIGINS`、空=关闭、精确 Origin、不反射 `*`）成立。建议改 D-001 字面或同时接受两键，并在操作员 YAML 补一行注释。

### F-004 · E-002 具名 S3 scratch 日志无法独立复核

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| 状态 | open |

- E-002 写 scratch `go-test-1.log` / `go-test-2.log` / `vitest-1.log` / `vitest-2.log`。仓库、`apps/api`、`apps/web`、本机 Temp 中**无同名文件**。仓内 `go-test*.log` / `vitest*.log` 时间戳为 2026-08-14；Temp `go-full-test*.log` / `web-full-test*.log` 为 2026-08-16 且 vitest 为 1027/1027，对不上「1046/1046 · 2026-08-17」。
- 不把缺失日志升级为「S3 未做」。checkpoint 与本轮定向测试支持「这些驱动测试当前为绿」。关门前若要坚持「全量两遍」，请把日志挂到本目标 `attachments/` 或给出可打开路径。

### F-005 · 开页网络失败仍走 `reauth-required` 终端（Token 已保留）

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| 状态 | open |

- `restoreSession` 在 refresh throw/500 时保留 Token，但返回 `{ kind: "reauth" }`；`AuthContext.tsx:82-86` 将其映射为 `reauth-required`（失败屏 + reauth 去 `/login`）。D-001 未要求新的 restore 结果种别，故不升 required。
- 与 F-001 不同：此处**不清** Token。若 F-001 修好而此处不变，弱网开页仍会看到「请重新登录」，只是 reload 后有机会自动恢复。可选：区分「凭证失效」与「暂时不可达」并提供 retry/reload。

### F-006 · 若干冻结/安全路径缺测试锁

| 字段 | 值 |
|------|-----|
| 严重度 | low |
| 建议 | **recommended** |
| 状态 | open |

- F01：无 `403` 清 Token 用例；无 `authFetch` 保 Token 用例（主缺口已在 F-001）。
- F05：无「空白名单不发 CORS 头」；无 YAML `cors_origins` / env `HTTP_CORS_ORIGINS` 解析测试。
- F07：`TestAuthRefreshRotates` 只断言旧 refresh 401，不断言撤销后的 `error == REFRESH_TOKEN_EXPIRED`（handler 三错误共一支，实现上应已覆盖）。

## 必改项汇总

| ID | 级别 | 要做什么 |
|----|------|----------|
| **F-001** | required · high | `authFetch` 不得在 refresh 网络/5xx（`doRefresh` 已 `false` 且未清 Token）时 `clearTokens`/`onAuthLost`；补 `authFetch` 入口测试 |
| **F-002** | required · med | HEAD + 已注册 GET 不得被包装器打成 405；补与裸 mux 一致的测试 |

Recommended（不阻断无条件放行，但建议在 S4 前处理或书面 residual）：F-003 键名/操作员 YAML；F-004 挂 S3 日志；F-005 reauth UX；F-006 补测。

## 与既有意见的异同

- 本目标此前 **无** self / independent 条目（索引为「尚未到达审计节点」）。本条为 **A-001**，共用序列首条。
- 与 [GOAL-020 A-002](../../GOAL-020-w15-user-perspective-findings/03-audit/A-002-s1s2-independent.md) 不冲突：彼时审的是台账真实性；W15-F04「首方崩溃」已改写；本次审的是批 A **实施**。
- 项目级决策要求交叉审默认先 self 再 independent。本次由用户直接 `/audit` 指定本 scope；本意见不补写 self，也不因此把 S4 标完成。

## 结论 + 建议给编排器/用户的下一步

F02、F05、F07 与 F01 的 `doRefresh`/`restoreSession` 层、F04 的未知路径/错误方法 JSON 信封均有可打开的代码与驱动测试，E-002 对上述主张大体属实。不能无条件放行 S4 / 把批 A 当 done：中会话 `authFetch` 仍会在 refresh 抖动时抹会话（**F-001**，W15-F01 未闭环），且 F04 包装器回归性拒绝合法 HEAD（**F-002**）。

建议 `/govern`：

1. 响应本 A-001；F-001 / F-002 走 `fixed`（改代码+测试）或用户书面 `accepted-residual` / `user-overruled`（写清范围与复审触发）。未合法闭合前 **不得** 将 S4 / 目标标 `done`。
2. F-003～F-006 可一并修或 residual。
3. 不要用 `progress: 3/4` 或 checkpoint 哈希代替 finding 闭合。

## 声明

本意见 `source: independent`，不修改 status / progress / 检查点 / goal-tree / 应用代码。响应由 **/govern** 处理。
