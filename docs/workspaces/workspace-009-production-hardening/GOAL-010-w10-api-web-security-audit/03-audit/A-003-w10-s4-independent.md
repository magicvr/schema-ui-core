---
id: A-003-w10-s4-independent
doc: audit-entry
goal: GOAL-010-w10-api-web-security-audit
title: W10 S4 关门前 independent 复核（D-003 范围 3 条修复 + 4 条作废）
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-08-21
scope: D-003 调和后实施范围是否 genuine fixed / 作废是否有据，且修复未改错逻辑、未引入新缺陷（S4 / finding-closure / execution-facts）
verdict: pass
status: recorded
parent: GOAL-001-production-hardening
created: 2026-08-21
updated: 2026-08-21
version: 0.1.0
---

# A-003 · W10 S4 关门前 independent 复核（2026-08-21）

## 条目头

| 字段 | 值 |
|------|-----|
| **source** | independent |
| **auditor** | grok-build (grok-4.6 · reasoning high · `/audit`) |
| **类型** | close-out / finding-closure / execution-facts |
| **scope** | D-003 调和后的 3 条实施（F-001/F-002/F-007）是否 genuine fixed；4 条作废（F-003～F-006）是否有源码依据；修复是否改错既有逻辑或引入新缺陷。消费清单权威 = [D-002](../01-decision/D-002-w10-scope-and-go-hold.md) + [D-003](../01-decision/D-003-w10-scope-reconciliation.md)；实施 = [E-002](../02-execution/E-002-w10-s3-implementation.md)；self = [A-002](A-002-w10-self.md) |
| **verdict** | **pass** |
| **工作区** | `workspace-009-production-hardening`（Root `GOAL-001-production-hardening`；canonical `docs/workspaces/workspace-009-production-hardening/`；`plan_refs`/`primary_plan` = `VP-009-production-hardening`；`vision_ref` = `schema-ui-core-admin-foundation@0.2.0`；`shared_materials_catalog: none`） |

## 范围与区间

- **覆盖**：S3 实施后的代码与回归；D-003 作废判定的源码对质；超时包装是否改变认证/刷新/页面加载/资源请求的既有失败语义；选项源正则是否误伤合法 URL。
- **方法**：工作区绑定核对 → GOAL-010 五件套通读 → 点名路径源码抽验（含调用点）→ 本会话复跑回归（见下）→ Node 复现 F-007 WHATWG 逃逸形状。未做动态 exploit、未起 compose、未接 live Postgres、未改密。
- **不覆盖**：不改 `status` / `progress` / 检查点 / 方案正文 / goal-tree。不把 informational F-008～F-012 升格为本波 required。不自行恢复 VP-008 `go` 宣称，不把本意见当作已关门。
- **排除**：informational 不审闭合；I-003 仍为 non-blocking 用户裁决（本条即工作区惯例要求的 grok 复核腿，但不代 `/govern` 关闭该信息项）。

## P-005 / 工作区核对

| 核对项 | 结论 |
|--------|------|
| 工作区绑定 | `workspace.md`：`id=workspace-009-production-hardening`；Root `GOAL-001-production-hardening`；canonical 与 `goal-tree.md` 一致；`vision_role: delivery`；`primary_plan` = `VP-009-production-hardening`。Charter `schema-ui-core-admin-foundation@0.2.0`；VP-009 `vision_ref` 精确匹配。共享资料目录 `none`，本 scope 未把资料引用当关闭证据。未读取其他工作区作为关闭依据。 |
| I-001（finding 清单） | **verified**（A-001 + D-003 调和）。消费实施 3 = F-001/F-002/F-007；作废 4 = F-003～F-006。 |
| I-002（范围 + VP-008 go） | **verified**（D-002 整单 7 条 + 暂挂 go；D-003 调和实施 3 + 作废 4）。`00-meta.md` 已标 verified。`01-decision.md` 与本文件索引表的 I-002 行仍写 open，属索引滞后（见 recommended F-001），不以滞后行否定 D-002。本条确认代码闭合条件已满足；**恢复对外 go 宣称仍须 `/govern` 书面决策**。 |
| I-003（provider 偏差 / 是否追加 grok 复核） | **open / non-blocking**；最晚阶段为关门前。本条满足 grok independent 腿，不替代用户对该信息项的书面关闭。 |
| 到期 required 信息项 | 无到期未关闭 required 信息项阻断本 S4 代码闭合 scope。 |
| 共享资料 | 无 |

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 3/3 实施项均有可核对代码改动，与 D-003 对齐 | 见逐条闭合判定；无偷换范围、无以作废项冒充 fixed |
| 4 条作废均有现行源码依据，不是假修复 | F-003 预览模板静态 + `noopener` 会丢窗口引用；F-004 防护式 UPDATE + `ErrAlreadyRevoked`；F-005 前导点已剥；F-006 去重后 1..64 |
| A-002 回归主张可重复核对 | 本会话：`apps/api` `go vet ./...` **exit 0**；`go test ./...` **全绿**（含 `config`/`authsession`/`auth` `-count=1`）；`apps/web` `npm test` **76 files / 1083 tests**；`npx tsc -b` **exit 0**。未再跑 `npm run build`（其 `prebuild` 会改写 conformance claim；类型检查已覆盖实现面） |
| 超时接线未改错既有认证语义 | `withTimeout` 仅包默认生产 `fetch`；注入 fetcher 原样；超时/中止走既有 network-throw 分支（login → `LOGIN_NETWORK`；refresh → 保 token，W15-F01）；auth-client **23/23** 仍绿（含 D-001 P2 generation、401 刷新、跨 tab 形状） |
| F-007 逃逸形状真实且已锁 | Node：`new URL("/\\evil.example/api/roles", "https://app.example")` → `https://evil.example/api/roles`（origin 逃逸成立）；修复后正则拒绝反斜杠；动态选项测试拒绝该组 URL 且不发 fetch |
| 未引入阻断性新缺陷 | 源码抽验未见把 fail-closed 改回 fail-open、未见刷新并发语义回退、未见合法 `/path?query` 选项源被误拒 |

## 对照成功标准

| 标准 | 本 scope | 状态 | 证据 |
|------|----------|------|------|
| S1 独立审计落盘 | 前序 | 达成 | A-001 + 附件全文 |
| S2 用户范围/go 裁决 | 前序 | 达成 | D-002 + D-003 |
| S3 按范围实施 + API/Web 回归 | 前序 + 本条复跑 | 达成 | E-002 + 本会话回归 |
| S4 self/independent 复核 required 合法闭合 | **本条判定代码闭合条件已满足** | 达成（意见层） | A-002 self pass + 本条 independent pass。本意见不改 status；合法闭合与关门由 `/govern` 响应 |

## 逐条闭合判定（D-003 消费）

### F-001 · env.example 硬编码真实凭据 → **fixed**

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 模板占位 | `apps/api/configs/env.example` | 真实 host/user/password/DSN 已换成占位符；加了「NEVER commit real credentials」注记；示例 `sslmode` 收紧为 `require`（仅注释模板，不改运行时默认） |
| 测试夹具 | `apps/api/internal/config/config_test.go` `TestDBPostgresExplodedParams` | 密码/host/name/user 改为 `test-only-password` / `db.example.internal` / `appdb` / `appuser`；断言同步。DSN 拼接逻辑未改，只换夹具字符串 |
| 跟踪文件残留 | `apps/` 下 `*.go`/`*.ts`/`*.tsx`/`*.yml`/`*.yaml`/`*.json`/`*.example` | 明文凭据不在跟踪的应用/配置文件中。本地 `apps/api/configs/.env` 仍有开发机值，但 **gitignored 且 `git ls-files` 未跟踪**（既有卫生，非本波回归） |
| 逻辑回归 | `go test -count=1 ./internal/config` | **ok**。占位符替换不改变 `Load`/`ValidateProd` 分支 |

原缺陷（版本控制模板携带真实库凭据）不再成立。Git 历史与本地 `.env` 仍视为已泄露——轮换是用户动作，A-002 已移交，本条不升格为开放 required。

### F-002 · Web fetch 统一超时 → **fixed**（未改错认证/资源语义）

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 包装 | `apps/web/src/lib/fetch-timeout.ts` | 默认 30s；每调用解析 `globalThis.fetch`（兼容测试后置 stub）；已中止的 caller signal 快速失败且不发请求；与 caller signal 组合；`finally` 清 timer |
| auth 接线 | `apps/web/src/account/auth-client.ts` | `postJSON` 与 `authFetch` 首次/401 重试均走 `timeoutFetch`。超时表现为 throw：login/mfa 外层 `catch` → `LOGIN_NETWORK`；`doRefresh` `catch` → `false` 且不清 token（W15-F01 保持） |
| 页面 schema | `apps/web/src/protocol/load-page.ts` | 默认 `withTimeout()`；注入 fetcher 不变。生产 `App.tsx` 未传 `schemaFetcher`，因此默认路径**确实**走到超时包装 |
| 资源请求 | `main.tsx` → `createConfigAwareFetcher(authFetch)` → `runRequest`/`library.preview` | A-001 点名的 `render.tsx` `runRequest` 本身未改文件，但生产 fetcher 即 `authFetch`，超时覆盖资源 CRUD、预览/下载 blob。`runRequest` 已有 `REQUEST_FAILED` catch，超时不会变成未处理拒绝 |
| 动态选项 | `form-controls.tsx` | 无注入时走 `defaultTimeoutFetch`；生产 `fetcher={crud?.fetcher}` 为已超时的 `authFetch`，不双重包装 |
| 测试隔离 | load-page / 动态选项 / conformance | 注入 fetcher 原样，测试确定性保持。新增 `fetch-timeout.test.ts` 5 例（透传/超时/组合/已中止/30s 默认） |
| 行为未回退 | `auth-client.test.ts` 23/23 | 含 generation 守卫、401 刷新、refresh 5xx 保 token、登出后丢弃 in-flight 旋转 |

`fetch()` 在响应头到达时即 settle，随后 `clearTimeout`，因此大文件 **body** 读取不受 30s 墙限制（预览/下载 blob 不因本包装在传完头后被切断）。连接挂起（头不到）才会 abort——与 F-002 原问题对齐。

### F-007 · 选项源 URL 反斜杠缺口 → **fixed**（核实后的真实缺口，不是误修）

| 核对项 | 路径 | 结论 |
|--------|------|------|
| 原正则 | 修复前 `/^\/(?!\/)[^\s\#]*$/` | 允许 `/\host`。本会话 Node 复现：`new URL("/\\evil.example/api/roles", "https://app.example")` → **`https://evil.example/api/roles`**（origin 逃逸）。`/%5C…` 百分号编码**不会**逃逸，保持同 origin |
| 现行正则 | `form-controls.tsx:85` `[^\s\\\#]` | 拒绝空白、反斜杠、`#`；仍允许 query（`?`），与「MAY carry a query string」注释及既有 `/api/roles?pageSize=100` 用例一致 |
| 同型 11 处 | `app-manifest` / `branding` / `resource` / `request-construction` / `upload-orchestration` / `render.tsx` / `form-controls.ts` / `form-controls.tsx` `isPreviewableUrl` | 其余 10 处修复前已含 `\\`；仅 optionsSource 这一处是缺口。本修复未改那 10 处 |
| 回归锁 | `form-controls.dynamic-options.test.tsx` | `/\evil.example`、`/api\evil`、`/api/roles\` 拒绝且 `fetcher` 不被调用；原合法加载用例仍绿 |

A-001 原文建议禁 `..`；D-003 升级为反斜杠 origin 逃逸是正确再诊断。同 origin 的 `..` 不是本缺口的利用面。

### 作废判定（未做假修复）

| F-ID | D-003 结论 | 本条核实 |
|------|------------|----------|
| F-003 `window.open` 无 noopener | 作废（不成立） | **同意**。`render.tsx:350-384`：预览窗写入固定模板 + `blob:` iframe `sandbox=""`（无脚本、无 same-origin）。第三个参数 `"noopener"` 会使 `window.open` 返回 `null`，功能依赖该引用 `document.write`。不是「该加却没加」的可实施缺口 |
| F-004 刷新旋转 PG 原子性 | 作废（误报） | **同意**。`accounts.go:337-359`：`UPDATE … WHERE id=? AND revoked_at IS NULL` + `RowsAffected`；0 行且记录存在 → `ErrAlreadyRevoked`。`auth.go:259-264` 将该错误映射为 `ErrTokenRevoked`，失败方不签发第二对。A-001 只读了调用层 |
| F-005 文件名点前缀 | 作废（误报） | **同意**。`sanitizeClientFilename`：`/^[._-]+|[._-]+$/g` 已剥前导/尾部分隔符；空结果回退 `download` |
| F-006 凭据作用域无上限 | 作废（误报） | **同意**。`service_credentials.go:151-155`：`normalizedCredentialScopes` 去重后强制 1..64；无 update 通路可绕过（仅 list/detail/create/revoke） |

## Findings

### F-001 · 审计索引与工作区波次表滞后于 S3 事实

- 严重度：**low** ｜ 建议：**recommended** ｜ 状态：open
- **描述**：`00-meta.md` / D-002 / D-003 / E-002 / A-002 已把 I-002 标为 verified、S3 完成；但 `01-decision.md` 信息表 I-002 仍为 `open`，本文件 `03-audit.md` 信息核对表 I-002 仍为 `open`，`workspace.md` 波次表 W10 行仍停在「S1 完成 · S2 待确认」。不影响代码闭合，但会误导下一轮编排扫描。
- **证据**：三处索引 vs `00-meta.md` 检查点 3/4 与 D-002 §4。
- **建议**：`/govern` 响应时同步这三处索引（不改 status/progress 则仅文档卫生）。

### F-002 · `withTimeout` 未在 settle 后移除 caller signal 监听

- 严重度：**low** ｜ 建议：**recommended** ｜ 状态：open
- **描述**：`outer.addEventListener("abort", …, { once: true })` 在 `finally` 中只 `clearTimeout`，不 `removeEventListener`。生产默认路径几乎不传长期 signal（auth/load-page/form-controls 均未传），因此不是当前利用面。若未来大量请求共享一个长寿命 `AbortSignal`，监听器会累积到该 signal abort/GC。
- **证据**：`apps/web/src/lib/fetch-timeout.ts:33-40`。
- **建议**：`finally` 中移除监听；可选补一条「共享 signal 多次调用后 listener 不增长」的测试。不阻断本波闭合。

### F-003 · 预览窗可在写入后将 `opener` 置空（纵深，非推翻作废）

- 严重度：**low** ｜ 建议：**recommended** ｜ 状态：open
- **描述**：F-003 作废成立（不能用 `window.open` 第三参数 `noopener`）。写入静态文档后执行 `previewWindow.opener = null` 可在不丢失引用的前提下切断 opener，作为后续若模板改入不可信内容的纵深。当前模板无脚本，不是 required。
- **证据**：`render.tsx:353-384`。
- **建议**：非本波范围；若做，单独小改 + 预览回归。

## 必改项汇总

无。开放 required = **0**。

## 与既有意见的异同

| 来源 | 异同 |
|------|------|
| A-001（independent · DSH） | 原文 1 HIGH + 6 MED + 5 info。本条不重开审计面，只复核 D-003 消费结果。同意 F-001/F-002/F-007 为真实缺口且已修；同意 F-003～F-006 作废。F-007 从「可禁 `..`」升级为反斜杠 origin 逃逸，本条用 Node 复现确认升级正确 |
| A-002（self） | 同意 3 fixed + 4 作废 + 开放 required = 0 + 回归全绿。本条复跑 1083/1083 与 go 全绿一致。补充：生产 schema 加载**确实**走默认 `withTimeout`（`App` 不传 `schemaFetcher`）；`runRequest` 经 `authFetch` 间接受益。本条新增 recommended F-001（索引滞后）、F-002（listener 清理）、F-003（`opener=null` 纵深），均不升格 required |
| 密码轮换残余 | 与 A-002 一致：用户动作，本条不因 git 历史或 gitignored `.env` 把 F-001 重开为 required |

## 结论 + 建议给编排器/用户的下一步

- **verdict = pass**：D-003 范围内 3 条 genuine fixed、4 条作废有据、required 开放 0、修复未改错认证/刷新/选项源/DSN 逻辑，本会话回归与 A-002 主张一致。
- **不是关门、不是恢复 go**：I-003 仍待用户书面关闭；VP-008 go 宣称恢复须用户书面（D-002/D-003；D-002 原文「7 条全部 fixed」应以 D-003 调和后的 3 fixed + 4 作废为闭合条件）。
- 建议 `/govern`：① 将 A-001 F-001/F-002/F-007 标 fixed、F-003～F-006 按 D-003 作废闭合；② 同步 `01-decision.md` / `03-audit.md` I-002 与 `workspace.md` W10 行；③ 询问用户是否关门 + 是否书面恢复 go + 是否接受本条 3 条 recommended（修 / residual / 留待下波）；④ 密码轮换仍为用户侧残余。

## 声明

本意见不修改 status/progress；响应由 `/govern` 处理。
