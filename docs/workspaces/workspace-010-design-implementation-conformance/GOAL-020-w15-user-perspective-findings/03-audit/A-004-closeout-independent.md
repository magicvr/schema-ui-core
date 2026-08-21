---
id: GOAL-020-w15-user-perspective-findings
doc: audit-entry
record_id: A-004
source: independent
scope: close-out of GOAL-020 after I-001 D-002 and children GOAL-021/022/023 remediations (W15-F01～F14)
verdict: conditional
status: recorded
auditor: grok-build (grok-4.6 · reasoning high)
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-004 · W15 关门独立交叉审计（2026-08-17）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：close-out · I-001 / D-002 全部 in-scope 后，对照 GOAL-021/022/023 整改与 as-built（W15-F01～W15-F14）。**不是** /govern，不改 status。
- **verdict**：**conditional**
- **工作区**：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；`canonical_scope` 已核对；`shared_materials_catalog: none`；`plan_refs` + `primary_plan` = `VP-010-design-implementation-conformance`）

## 范围与区间

- **covered**：
  - 本区 `workspace.md` / `goal-tree.md` 绑定；GOAL-020 五件套（`00-meta` · D-001/D-002 · E-001～E-003 · A-001～A-003）
  - 下级 GOAL-021 / GOAL-022 / GOAL-023 的 `00-meta`、D-001 冻结、执行与关门自审
  - 对 D-002 in-scope F01～F14 打开 as-built（重点：F01 `authFetch`、F04 JSON 信封、F07 refresh 码、F11 GET 404、F03 RFC3339 不改名）
  - 本轮定向复跑：`apps/api` `TestJSONRouteErrors404And405` / `TestAuthRefreshInvalid` / `TestWalletByOwnerAutoCreate` / `TestWalletSelfAutoCreateAndIdempotency` / `TestFormatRFC3339Milli` / `TestWrapSecurityCORSAndNosniff` 绿；`apps/web` `auth-client.test.ts` + `data-table.test.tsx` + `schema-table.test.tsx` 65/65 绿
- **不 covered**：浏览器活栈点验；全量 `go test ./...` / 全量 vitest / e2e 复跑；其他工作区正文（含 workspace-011 钱包消费路径，只核对本区已记录的 go-impact 句）
- **本意见不修改** `00-meta` status / 检查点 / `progress`、`goal-tree.md` / `workspace.md`、方案正文、或任何应用代码

## 成果（有证据）

### 工作区与信息门禁

| 字段 | 核对 |
|------|------|
| `workspace_id` | `workspace-010-design-implementation-conformance` |
| Root | `GOAL-001-design-implementation-conformance`（`parent: null`） |
| canonical | `docs/workspaces/workspace-010-design-implementation-conformance/` |
| 共享资料 | `none`；无资料行被当作事实或关闭证据 |
| GOAL-020 归属 | `parent: GOAL-001-design-implementation-conformance`；树内 **active · 8/8**（待本关门） |
| I-001 | **closed**（D-002：F01～F14 全部 in-scope；A→B→C；父目标等子目标；F03 不改字段名；F05 留本区；F11 GET 404 + POST 创建） |
| 子目标 | GOAL-021 / 022 / 023 均为 `done` 4/4；D-002「子目标完成前不得标 done」对父目标仍成立，本条不改 status |
| 既有意见 | A-002 required F-001/F-002 已由 A-003 改写 D-001 闭合；GOAL-021 A-001 required 已由该目标 A-002 闭合。无未决冲突需 P-004 |

### 逐条 as-built（对照 D-002 + 各批冻结）

| 台账项 | 批 | 独立核对 | 结论 |
|--------|----|----------|------|
| **W15-F01** | A | `doRefresh` 网络 throw / 非 401-403 不清 Token（`auth-client.ts:129-153`）。`authFetch` 仅在 refresh 后 `getRefreshToken() === null`（即 401/403 已清）时 `clearTokens`+`onAuthLost`（`:209-222`）。测试 `authFetch keeps tokens when refresh returns 500 after a 401` 本轮绿。 | **达成**（GOAL-021 A-001 F-001 闭合仍成立） |
| **W15-F02** | A | `DataTable` 错误态 `onRetry` + `feedback.retry`；`schema-table` `retryNonce` + `feedback.resourceFetchFailed`。本轮相关测试绿。 | **达成** |
| **W15-F03** | B | `formatRFC3339Milli`；scheduled-tasks / dictionary 改 RFC3339 串；filelibrary 字段名仍 `created`，值为 RFC3339（`filelibrary.go:92,177`）。handler JSON 输出无 Unix 秒时间字段。`TestFormatRFC3339Milli` 绿。 | **达成**（不改字段名） |
| **W15-F04** | A | `WithJSONRouteErrors`：未知路径 404 `NOT_FOUND`；错误方法 405 `METHOD_NOT_ALLOWED`；HEAD+已注册 GET 交给 mux（`route_envelope.go:12-15`）。生产 `composition.go:467`。本轮 `TestJSONRouteErrors404And405` 绿。 | **达成**（GOAL-021 A-001 F-002 闭合仍成立） |
| **W15-F05** | A | `WrapSecurity` 全局 nosniff；YAML `http.cors_origins` / `HTTP_CORS_ORIGINS`；空 = 不发 CORS。本轮 `TestWrapSecurityCORSAndNosniff` 绿。 | **核心达成**；白名单启用时 Allow-Headers 不覆盖首方实际头 → **F-003** |
| **W15-F06** | C | 改密仍立即 `token_version`+撤销 refresh（`account_self.go:249-257`）。首方 `authFetch` 在 `/api/account/password` 成功后写 `password.changedNotice`；`LoginPage` 读出 `schema.account.passwordChangedReauth`。schema `onSuccess` 仍仅为 `reload`（与冻结字面「onSuccess messageKey」不同，但登录页提示路径可核）。 | **用户意图达成** |
| **W15-F07** | A | `authHandler.refresh` 对 invalid/expired/revoked 写 `REFRESH_TOKEN_EXPIRED`（`auth.go:186-188`）；登录仍 `UNAUTHORIZED`。catalog 中英与冻结一致。本轮 `TestAuthRefreshInvalid` 绿。 | **达成** |
| **W15-F08** | C | create 的 `fieldErrors.reason` 为 `required` / `string`（`resources.go:466-474`）。PATCH 仍 `must not be empty`；前端把 reason **原文**写到输入框（`render.tsx:1548-1549`），catalog 仍拼接 `创建字段无效: required`。 | **部分** → **F-004** |
| **W15-F09** | C | `FeedbackRegion` 对 `kind === "error"` 不设 4s timer（`render.tsx:1151-1153`）；关闭钮仍在。 | **达成** |
| **W15-F10** | B | 登录 429 写 `Retry-After`（`auth.go:106-108`，`auth_test.go:309-310`）。配额 `http.StatusRequestEntityTooLarge` 413 + `UPLOAD_QUOTA_EXCEEDED`（`upload.go:326`）。 | **达成** |
| **W15-F11** | B | GET by-owner / me / me/entries 走 `GetUserAccountByOwner`，缺失 `writeWalletError` → 404 `WALLET_NOT_FOUND`。创建走 POST。本轮 `TestWalletByOwnerAutoCreate` / `TestWalletSelfAutoCreateAndIdempotency` 绿（断言 404，未锁 error 码）。 | **API 达成；首方面未接 POST** → **F-001** |
| **W15-F12** | B | `DefaultPageSize = 20`；resources / wallet / scheduled-tasks runs 已用该常量。会话列表仍默认 **10**（`account_self.go:272`）。 | **原引用点达成**；残余 → **F-006** |
| **W15-F13** | C | 会话行有 `current`；当前行带本次请求 `userAgent` / `ip`（`account_self.go:311-332`）。`account.json` 会话表列仍只有 createdAt / expiresAt / status；无 i18n 列。`account_self_test.go` 无 `current` / `X-Refresh-Token` 锁。 | **API 达成；用户列表仍无法辨认当前设备** → **F-002** |
| **W15-F14** | C | (1) `accountToMap` 有 `decimals: 2`；(2) 顶栏 menu 仍仅 Escape，无方向键（D-023 明示本批不落地）；(3) 空态区分 `listEmpty` / `noItemsMatch`；(4) MFA 取消钮在 `LoginPage.tsx:328-337`；(5) `apps/web/nginx.conf;C/` 已不存在。 | **4/5 达成**；第 2 点无父级 residual → **F-005** |

### 回归与子目标关门（部分可核）

- GOAL-021：有 independent A-001 + 响应 A-002 + 关门 A-003；E-002 记录 checkpoint `285c7e8` 与定向修复。本轮复跑批 A 驱动测试绿。
- GOAL-022 / GOAL-023：仅 self 关门一段话；E-001 写「回归见 E-002」但 **E-002 不存在**。本轮对 F03/F11/F04/F07 定向测试绿，**不**把子目标「Go/vitest 两遍绿」当作本轮独立全量证据。
- GOAL-020 父级执行台账停在 E-003（开启批 A、当时 5/8）；**没有**记录 R1～R3 完成的 E 条。`00-meta` 8/8 来自子目标 `done`，不是父级时间线事实。

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| I-001 已书面裁决；无到期未决 required 信息项 | **达成** | D-002 closed |
| F01～F14 按 D-002 全部 in-scope 实施，无静默 defer | **部分** | API/核心路径大多落地；F11 首方面、F13 会话表、F14-2 方向键未闭环 |
| F01 `authFetch` 网络/5xx 不清会话 | **达成** | 代码 + 本轮测试 |
| F04 404/405 JSON 信封；HEAD 不误伤 | **达成** | `route_envelope.go` + 本轮测试 |
| F07 `REFRESH_TOKEN_EXPIRED`；登录仍 `UNAUTHORIZED` | **达成** | `auth.go` + 本轮测试 |
| F11 GET 只读 404 `WALLET_NOT_FOUND`；创建走 POST | **部分** | API 是；首方 `my-wallet.json` 仍只 GET（F-001） |
| F03 RFC3339 毫秒串且不改 `created` 字段名 | **达成** | scheduled-tasks / dictionary / filelibrary |
| D-002 闭门等待：子目标均 `done` 后才可把父目标标 done | **结构达成** | 三子目标 done；父目标仍 active。本意见不放行关门 |
| 无未闭合 high required；无到期 required 信息项 | **未达无条件关门** | 本条 2 条 med required 未闭合 |

## Findings

### F-001 · 首方「我的钱包」仍只 GET，缺失账户对用户是死 404

- 严重度：**med**
- 建议：**required**
- 状态：open
- 关联：W15-F11；D-002「GET 只读 / 404 `WALLET_NOT_FOUND` / 创建仍由 POST 触发」；GOAL-022 D-001
- 描述：API 层 GET `/api/wallet/me`、`/api/wallet/me/entries`、`/api/wallet/by-owner/{id}` 已只读，缺失走 `GetUserAccountByOwner` → `WALLET_NOT_FOUND`。这与 D-002 的 HTTP 语义一致，本轮测试也锁了 GET 404 + POST 创建。但同仓首方面 `apps/api/internal/modules/wallet/schema/my-wallet.json` 三个 statCard `dataSource` 与流水表仍只 GET，**没有任何 POST `/api/wallet/me` 开户动作**。整改前 GET 会懒创建，首访即有零余额账户；整改后新用户打开「我的钱包」得到 404，F02 重试仍是 GET 404。D-002 把 go-impact 写成 workspace-011 自动开户读路径，**漏记本仓首方 schema**。W15 是真实用户视角关门：API 正确但首方面退化，不能把 F11 标为对终端用户已闭环。
- 证据：`apps/api/internal/modules/wallet/schema/my-wallet.json:36,47,58,101`；`apps/api/internal/handler/wallet_self.go:34-46,66-76`；`apps/api/internal/handler/wallet.go:81-97`；GOAL-022 `01-decision/D-001-s1-freeze.md`；本轮 `TestWalletSelfAutoCreateAndIdempotency`（GET 404 / POST 200）。
- 建议修复：首方面在读之前（或 404 后）显式 POST `/api/wallet/me`，或提供可发现的开户操作；补 schema/渲染测试。测试补锁 404 的 `error == WALLET_NOT_FOUND`。

### F-002 · F13 的 `current` / UA / IP 未进入个人中心会话表

- 严重度：**med**
- 建议：**required**
- 状态：open
- 关联：W15-F13；GOAL-023 D-001
- 描述：原题是用户在「在线会话列表」无法识别当前设备、易误撤销自己。as-built 只在 JSON 增加 `current`，且仅当前行附带**本次请求** UA/IP（`account_self.go:311-332`）。`account.json` 会话表列仍是 createdAt / expiresAt / status，无 `current` / `userAgent` / `ip` 列，也无对应 i18n。首方表格只渲染声明列，用户界面与整改前一样。`account_self_test.go` 对 `current` / `X-Refresh-Token` 无断言。
- 证据：`apps/api/internal/handler/account_self.go:311-332`；`apps/api/internal/modules/account/schema/account.json:161-176`；`apps/web/src/i18n/messages/zh-CN.json` 仅有 created/expires/status 列键。
- 建议修复：会话表增加当前标记（及可选 UA/IP）列并补测试：带 `X-Refresh-Token` 的会话 `current: true`，其余 `false`。

### F-003 · F05 CORS Allow-Headers 覆盖不了本波首方实际请求头

- 严重度：**med**
- 建议：**recommended**
- 状态：open
- 描述：默认 `cors_origins: ""` 不发 CORS，同域 Nginx 不受影响。一旦按 F05 用途打开白名单，`WrapSecurity` 只允许 `Authorization, Content-Type`（`server.go:45`）。同波 `authFetch`/`withAuth` 固定带 `X-Refresh-Token`（F13）与 `Accept-Language`，预检会失败，认证请求与当前会话标记一并不可用。`TestWrapSecurityCORSAndNosniff` 未探测这两头。
- 证据：`apps/api/internal/server/server.go:32-53`；`apps/web/src/account/auth-client.ts:171-183`。

### F-004 · F08 规则码未在前端本地化；PATCH 仍是英文 reason

- 严重度：**med**
- 建议：**recommended**
- 状态：open
- 描述：create 已下发 `required` / `string`。`errorcatalog.Body` 对字段码仍把 caller message 拼进 `message`（现为「创建字段无效: required」）；`render.tsx` 把 `fe.reason` 原文写入输入框。PATCH 仍 `must not be empty`（`resources.go:510,660`）。中英混杂从长句换成码，用户可见问题未闭环。冻结字面只要求规则码，故不升 required。
- 证据：`apps/api/internal/handler/resources.go:466-474,510,660`；`apps/api/internal/errorcatalog/errorcatalog.go:221-227`；`apps/web/src/renderer/render.tsx:1548-1549`。

### F-005 · F14 方向键导航被批 C 冻结裁掉，父级无 residual

- 严重度：**low**
- 建议：**recommended**
- 状态：open
- 描述：D-002 写 F01～F14 全部 in-scope、无 defer。GOAL-023 D-001 把「菜单方向键」留到 `App.tsx` / `notification-bell.tsx` 且本批不落地。as-built 仍仅 Escape。F14 其余 4 点已做。父级未把第 2 点写成 `accepted-residual` / 显式 defer。
- 证据：GOAL-023 `01-decision/D-001-s1-freeze.md`；`apps/web/src/app/App.tsx:350-353`；`apps/web/src/app/notification-bell.tsx:127-130`。

### F-006 · 会话列表默认 pageSize 仍为 10

- 严重度：**low**
- 建议：**recommended**
- 状态：open
- 描述：D-001 原证据点（resources 10 / wallet 20 / task-runs 50）已归一到 `DefaultPageSize = 20`。`account_self` 会话列表仍 `intParam(..., 10)`，未引用该常量。原台账未点名此文件，故不升 required。
- 证据：`apps/api/internal/handler/account_self.go:272`；`apps/api/internal/handler/resources.go:37-38,395`。

### F-007 · 父级关门证据链偏薄（执行台账 / 树摘要 / 批 B·C 回归）

- 严重度：**low**
- 建议：**recommended**
- 状态：open
- 描述：GOAL-020 `02-execution` 无 R1～R3 完成事实（E-003 停在开批 A、5/8）。`goal-tree.md` 正文仍写「W15 进行中 5/8」，与树/表 8/8 不一致（本条不改 goal-tree）。GOAL-022/023 E-001 指向不存在的 E-002；关门自审各一段话。`wallet_self_test.go` 注释仍写「首次 /me 自动开户」，与 GET 404 测试体相反。不据此否定本轮定向绿。
- 证据：GOAL-020 `02-execution.md`；`goal-tree.md:46` vs `:40,99`；GOAL-022/023 `02-execution/` 仅 E-001；`wallet_self_test.go:54-56`。

## 必改项汇总

| ID | 建议 | 要做什么 |
|----|------|----------|
| **F-001** | required · med | 首方「我的钱包」在缺失账户时显式 POST 开户（或等价可发现创建），不要只把 404 丢给 statCard/表格 |
| **F-002** | required · med | 个人中心会话表展示 `current`（及冻结已有的 UA/IP），并补测试 |

Recommended（不阻断「带条件继续」，但关门前应修或书面 residual）：F-003 CORS 头；F-004 F08 本地化/PATCH；F-005 F14-2 residual；F-006 会话 pageSize；F-007 父级 E 条与树摘要。

## 与既有意见的异同

- **vs A-002（S1/S2 independent · conditional）**：彼时审台账真实性；required 已由 A-003 改写闭合。本条审 **整改后 as-built**，不重开 F06 机制写反 / F04 崩溃句。
- **vs GOAL-021 A-001**：同意 F01 `authFetch` 与 F04 HEAD 修复仍在；本轮复跑未推翻。
- **vs GOAL-022/023 A-001 self · pass**：不同意把批 B/C 写成用户视角已全部闭环。API 合同大体成立；F11 首方面与 F13 会话表对终端用户仍是原题。
- 无 P-004 意见冲突需用户在两条审计之间二选一。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** D-002 范围、I-001 关闭、三子目标 `done`、以及 F01/F02/F03/F04/F07/F09/F10 与 F05/F06/F12/F14 的主体，均有可打开代码与（本轮）定向测试。不能无条件把 GOAL-020 标 `done`：F11 把懒创建从 GET 挪到 POST 后，首方「我的钱包」没有创建路径（**F-001**）；F13 的当前会话标记没有进个人中心表（**F-002**）。

建议 `/govern`：

1. 响应本 A-004。F-001 / F-002 走 `fixed`（改首方 schema/UI + 测试）或用户书面 `accepted-residual` / `user-overruled`（写清范围与复审触发）。未合法闭合前 **不得** 将 GOAL-020 标 `done`。
2. F-003～F-007 可一并修或 residual。F-005 若维持「方向键另波」，须在父级书面 defer/residual。
3. 补父级 E 条记录 R1～R3 事实；不要用 `progress: 8/8` 代替 finding 闭合。
4. F11 的 go-impact 应同时覆盖本仓 `my-wallet.json`，不能只写 workspace-011。

## 声明

本意见 `source: independent`，不修改 status / progress / 检查点 / goal-tree / 应用代码。响应由 **/govern** 处理。
