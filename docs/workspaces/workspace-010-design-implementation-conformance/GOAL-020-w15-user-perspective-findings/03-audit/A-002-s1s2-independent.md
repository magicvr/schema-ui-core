---
id: GOAL-020-w15-user-perspective-findings
doc: audit-entry
record_id: A-002
source: independent
scope: S1 审视证据真实性 + S2 台账完整性（W15-F01～W15-F14 + E-001 证据链）
verdict: conditional
status: recorded
auditor: grok-build (grok-4.6 · reasoning high)
parent: GOAL-001-design-implementation-conformance
created: 2026-08-17
updated: 2026-08-17
version: 0.1.0
---

# A-002 · S1/S2 独立交叉审计（2026-08-17）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：ad-hoc / execution-facts · S1 审视证据真实性 + S2 台账完整性（W15-F01～W15-F14 ledger + E-001 证据链）。**不是**关门审计。
- **verdict**：**conditional**
- **工作区**：`workspace-010-design-implementation-conformance`（`root_goal` = `GOAL-001-design-implementation-conformance`；`canonical_scope` 已核对；`shared_materials_catalog: none`；`plan_refs` + `primary_plan` = `VP-010-design-implementation-conformance`）

## 范围与区间

- **covered**：GOAL-020 五件套（`00-meta` / `01-decision` + D-001 / `02-execution` + E-001 / `03-audit` + A-001）；`workspace.md` / `goal-tree.md` 绑定与 W15 行；D-001 全部 W15-F01～W15-F14 的 `file:line` 对照源码打开核验；E-001「729 键对齐」与「无代码修改」对照；I-001 是否被伪称为已裁决。
- **不 covered**：S4 响应、S5 用户裁决、整改实施、浏览器活栈点验、其他工作区上下文。
- **本意见不修改** `00-meta` status / 检查点 / `progress`、D-001 正文、`goal-tree.md` / `workspace.md` 状态列、或任何应用代码。

## 成果（有证据）

### 工作区绑定

| 字段 | 核对 |
|------|------|
| `workspace_id` | `workspace-010-design-implementation-conformance` |
| Root | `GOAL-001-design-implementation-conformance`（`parent: null`） |
| canonical | `docs/workspaces/workspace-010-design-implementation-conformance/` |
| `plan_refs` / `primary_plan` | `VP-010-design-implementation-conformance` |
| 共享资料 | `none`；无资料行被当作事实或关闭证据 |
| GOAL-020 归属 | `parent: GOAL-001-design-implementation-conformance`；树内 active · 2/5 |

### 「本波无代码修改」

- HEAD `59f454a`（`docs(workspace-010): 创建并登记 GOAL-020 W15…`）仅改本目标五件套 + `goal-tree.md` / `workspace.md`。
- 工作区干净；`apps/api` / `apps/web` 相对 HEAD 无本波 diff。
- 未见本波把 I-001 写成已裁决，也未见按台账实施修复。与 E-001 / `00-meta`「暂不启动代码修复」一致。

### I-001（P-005）

| 项 | 核对 |
|----|------|
| 级别 | `required` |
| 最晚需要阶段 | **S5**（用户裁决 / 整改启动） |
| 状态 | `open` |
| 是否伪装已决 | 否。D-001 §3 写「待用户裁决」；分批 A/B/C 仅「建议」 |
| 对本次 S1/S2 | **不构成 S1/S2 失败理由** |

### E-001 辅助主张

- `en-US.json` / `zh-CN.json` 均为 **729** 键，差集为空。E-001 该句成立。

### W15-F01～W15-F14 逐条打开核验

| 台账项 | 独立核对 | 结论 |
|--------|----------|------|
| **W15-F01** | `apps/web/src/account/auth-client.ts:129-136`：`doRefresh` 在 `postJSON` **throw**（网络失败）时 `clearTokens()`。同文件 `151` 对任意 `!response.ok`（含 5xx）也会清空。`382-385` 是 `restoreSession` 在 refresh 失败后返回 `{ kind: "reauth" }`，本身不再 `clearTokens`，但是清空后的启动后果。 | **成立**（382-385 为后果引用，可接受） |
| **W15-F02** | `schema-table.tsx:517-521` catch 只 `setError(err.message)`；`data-table.tsx:259-264` 红色 `<p>` 渲染 `{error}`，无重试。`resource.ts:295-297` 抛 `resource fetch failed: HTTP …`。 | **成立** |
| **W15-F03** | `handler/users.go:107` RFC3339 毫秒串；`scheduledtasks.go:329` / `dictionary.go:279` `Unix()` 秒整数；`filelibrary.go:177` 字段名 `created` + Unix 秒。格式/字段名割裂成立。`datetime.ts` 对非 ISO **返回 null**，表格走 `formatDisplayTime(value) ?? String(value)`，Unix 整数显示为数字，**不是** `Invalid Date`。 | **现象成立；「Invalid Date」过述** → F-003 |
| **W15-F04** | `composition.go:207` 为 `mux := http.NewServeMux()`，随后 `385` `mux.Handle(method+" "+pattern, …)`，无自定义 NotFound/405。默认 mux 对未注册路径返回 Go 原生 `text/plain` 404。但首方 `resource.ts:129-147` `readEnvelope` 对非 JSON `response.json()` **catch 后降级**，界面得到 `resource fetch failed: HTTP 404`，**不会崩溃**。 | **信封缺口成立；「崩溃」不成立** → **F-002** |
| **W15-F05** | `server.go:21-29` 只装配超时，无 CORS / 全局安全头。仓内无 `Access-Control` / API 级 CORS。`apps/web/nginx.conf` 才加 nosniff 等，且注释写明同域反代「needs no CORS」。纯 API 分域部署缺口成立。 | **成立** |
| **W15-F06** | `account_self.go:249-257` 注释与 `UpdateUser`：改密 **bump `token_version` + 撤销全部 refresh**，204。`authsession/users_repository.go:187-206` 写明已签发 access **立即**被中间件拒绝。`auth.go:463-471` `TokenVersion` 不匹配 → 401 `UNAUTHENTICATED`。`account.json:46-48` `onSuccess: reload`。默认 `AuthAccessTTL` 虽为 15m，但改密后 **不会**再活约 15 分钟。真实体验更接近：成功 → reload / 下一请求 → refresh 失败 → 立刻回登录。 | **行号在、机制叙述反了** → **F-001** |
| **W15-F07** | `auth.go:182-184` refresh 失败复用 `UNAUTHORIZED`；`errorcatalog.go:27` 文案「用户名或密码错误」。API 契约成立。首方 `doRefresh` 失败只 `clearTokens()`，**不把该文案铺到 Toast/登录页**。 | **契约成立；「用户看到误报文案」过述** → F-003 |
| **W15-F08** | `errorcatalog.go:221-224` 字段校验把 caller message 拼进 `message`；`resources.go:448` `Error()` = `field + " must " + reason`，`464` `"not be empty"` → `name must not be empty`。中文环境会出现「创建字段无效: name must not be empty」。 | **成立** |
| **W15-F09** | `render.tsx:1144-1153` + `1192`：`FEEDBACK_TOAST_MS = 4000`，error/success 一律定时消失。`1179-1186` **已有** `×` 关闭钮。核心问题是 4s 自动消失，不是「无关闭按钮」。 | **4s 消失成立；建议过时** → F-003 |
| **W15-F10** | `auth.go:104-106` 登录 429 `RATE_LIMITED`，全仓 **无** `Retry-After`。`upload.go:325-327` 配额用 429 + `UPLOAD_QUOTA_EXCEEDED`。目录文案是「上传被拒绝：超出每用户配额」/「登录失败次数过多，请稍后重试」，**不是**台账写的「请求过于频繁」。 | **状态码问题成立；文案张冠李戴** → F-003 |
| **W15-F11** | `wallet.go:82-103` `GET /api/wallet/by-owner/{ownerId}` → `GetOrCreateUserAccount` + 审计。`wallet_self.go:36-51` `GET /api/wallet/me` 同样懒创建。GET 写副作用成立。 | **成立** |
| **W15-F12** | `resources.go:393` 默认 10；`wallet.go:62` 默认 20；`scheduledtasks.go:305`（单任务 runs）默认 50。10/20/50 混用成立。 | **成立** |
| **W15-F13** | `account_self.go:312-328` 会话行仅 `id/createdAt/expiresAt/status/revokedAt`。authsession 无 UA/IP 字段；无 `current` 标记。 | **成立** |
| **W15-F14** | (1) `wallet.go:423-435` `accountToMap` 无 `decimals`；(2) `App.tsx:350-353` / `notification-bell.tsx:127-130` 仅 Escape，无方向键；(3) `schema-table.tsx:864` 空态固定 `feedback.noItemsMatch`；(4) `LoginPage.tsx:292-330` MFA 阶段仅校验钮、无取消；(5) `apps/web/nginx.conf;C/` 是**空目录**，不是空文件。 | **主体成立；第 5 点形状不准** → F-004 |

## 对照成功标准

| 标准 | 状态 | 证据 |
|------|------|------|
| S1 审视有可打开的代码证据，非空泛编造 | **部分达成** | 14 条均能指到真实文件；F06 影响机制与源码相反；F04「崩溃」与首方错误读取不符 |
| S2 台账 W15-F01～W15-F14 完整（描述 / 证据 / 影响 / 优先级）且为候选而非已改 | **部分达成** | 14 条齐；语气是改进候选；F06/F04 影响句需在 S5 前改正 |
| 本波无代码修复、与「先落盘再裁决」一致 | **达成** | E-001 + git `59f454a` 仅文档 |
| I-001 未伪称已决；S1/S2 不因 I-001 开放而失败 | **达成** | `00-meta` I-001 `open` · 最晚 S5；D-001 §3 待裁决 |
| A-001 声称「行号与上下文准确、无过度推论」 | **未完全成立** | 见 F-001 / F-002；A-001 过宽 |

## Findings

### F-001 · W15-F06 把「改密后 access 仍活约 15 分钟」写成事实，与 `token_version` as-built 相反

- 严重度：**med**
- 建议：**required**
- 状态：open
- 描述：D-001 W15-F06 写「服务端立即撤销全部 Refresh Token，但当前的 Access Token 依然有效约 15 分钟……10 多分钟后刷新 Token 失败被突然踢回登录页」。打开的代码是：改密走 `UpdateUser` + `PasswordHash`，**同一事务** `token_version + 1` 并撤销 refresh（`account_self.go:249-257`，`users_repository.go:187-206`）；中间件按 JWT `tv` 立刻拒绝旧 access（`auth.go:463-471`）。账户页 `changePassword.onSuccess.behavior = reload`（`account.json:46-48`）。15m 只是默认 `AuthAccessTTL`，**不是**改密后的残存窗口。真实用户路径是成功 204 后 reload / 下一请求立即 401 → refresh 失败 → 回登录。条目作为「改密后会话 UX 不完整」仍可保留，但必须改写机制与影响，否则 S5 会按错误威胁模型裁决（例如再做一遍已经存在的 access 作废）。
- 证据：`apps/api/internal/handler/account_self.go:249-257`；`apps/api/internal/modules/authsession/users_repository.go:187-206`；`apps/api/internal/auth/auth.go:463-471`；`apps/api/internal/modules/account/schema/account.json:46-48`；`01-decision/D-001-w15-findings-ledger.md` W15-F06 行。

### F-002 · W15-F04 「前端 `response.json()` 崩溃」对首方 Web 不成立

- 严重度：**med**
- 建议：**required**
- 状态：open
- 描述：未注册路径走默认 `ServeMux`、返回 Go 原生 `text/plain` 404/405、缺少统一 JSON 信封——这句成立，引用 `composition.go:207`（mux 创建、无自定义 NotFound）可接受。但「导致前端 `response.json()` 崩溃」与首方实现不符：`readEnvelope` 将 `response.json()` 包在 try/catch，非 JSON 时返回 `{ code: "UNKNOWN", message: "" }`，调用方得到 `resource fetch failed: HTTP 404`，不抛未捕获异常。High 分组若主要靠「崩溃」支撑，S5 前须改写影响（契约/第三方客户端 / 辨识度），不得再写成首方运行时崩溃。
- 证据：`apps/api/internal/composition/composition.go:207,385`；`apps/web/src/renderer/resource.ts:129-162`；`01-decision/D-001-w15-findings-ledger.md` W15-F04 行。

### F-003 · 若干「真实影响」句过述（不否定条目本身）

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：S5 前建议收紧措辞，避免用户按字面误判严重度。  
  1. **W15-F03**：`formatDisplayTime` 对 Unix 秒返回 null，单元格显示原始数字，不是 `Invalid Date`（`datetime.ts:23-30`，`data-table.tsx:70`）。格式不统一仍在。  
  2. **W15-F07**：首方 refresh 失败不清渲染 catalog 文案，只清会话；「用户看到用户名或密码错误」对内置 Web 过强。API 复用 `UNAUTHORIZED` 仍在。  
  3. **W15-F09**：error Toast 已有显式关闭钮（`render.tsx:1179-1186`）；缺口是 4s 自动消失，不是「无留存且无关闭」。  
  4. **W15-F10**：配额 429 的目录文案是 `UPLOAD_QUOTA_EXCEEDED`，登录 429 是「登录失败次数过多…」，台账「请求过于频繁」不是这两处文案。缺 `Retry-After`、配额不宜 429 仍成立。
- 证据：上表对应源码行；D-001 各「真实影响」句。

### F-004 · 台账引用卫生（缩写路径 / 遗留物形状 / S3 编号）

- 严重度：low
- 建议：**recommended**
- 状态：open
- 描述：W15-F03/F08/F10/F12 部分证据写成 `users.go:107`、`resources.go:464` 等无目录短名；本仓库能唯一落到 `apps/api/internal/handler/`，但与同表全路径不一致，复审成本高。W15-F14 第 5 点写「空文件 `apps/web/nginx.conf;C`」，实际是空目录 `apps/web/nginx.conf;C/`。`00-meta` 路线图 S3 写成「independent 审计（A-001）」——A-001 已是 self，本条才是 S3 independent（A-002）。
- 证据：`01-decision/D-001-w15-findings-ledger.md` 证据列；`apps/web/nginx.conf;C/`（list_dir 空）；`00-meta.md` S3 检查点行。

## 必改项汇总

| ID | 建议 | 要做什么 |
|----|------|----------|
| **F-001** | required | 改写 W15-F06 影响：改密已立即作废 access（`token_version`）+ 撤销 refresh；UX 是成功后立刻掉登录，不是 15 分钟后突然过期 |
| **F-002** | required | 改写 W15-F04 影响：保留 404/405 非 JSON 信封缺口；删除或限定「首方 `response.json()` 崩溃」 |

I-001 仍为 S5 required 信息项，**不是**本条 S1/S2 必改 finding。

## 与既有意见的异同（vs A-001）

- **A-001 self · pass**：同意工作区归属、本波无代码改动、I-001 留待裁决、14 条大多能指到真实行号。
- **不同意** A-001「行号与上下文准确无误；杜绝空泛假设与过度推论」的全称。W15-F06 机制写反；W15-F04 崩溃句对首方不成立。因此不能原样接受 A-001 对 S1/S2 的无条件 pass。
- 无 P-004 意见冲突需用户在两条审计之间二选一：本条收紧证据精度，不否定「先落盘、不修代码」边界。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** S1/S2 作为「只读审视 + 14 条候选台账 + 无本波代码修复」总体真实；E-001 双语 729 键与 git 边界成立。但有 **2 条 med required**（F-001 / F-002）必须在把台账交给 S5 用户裁决前由 `/govern` 改写 D-001 对应行。不得用本意见推进 S5 或关门。

建议 `/govern`：

1. 响应并闭合 F-001 / F-002（改 D-001 影响句；不要改应用代码——本波仍禁止修复）。
2. 顺手处理 F-003 / F-004（措辞与引用卫生）。
3. I-001 保持 open，直到用户书面裁决 in-scope / 分批。
4. 不要把 A-001 的 pass 当作 S1/S2 证据已全部过关。

## 声明

本意见不修改 status / progress / 检查点 / 方案正文 / goal-tree / workspace。响应由 `/govern` 处理。
