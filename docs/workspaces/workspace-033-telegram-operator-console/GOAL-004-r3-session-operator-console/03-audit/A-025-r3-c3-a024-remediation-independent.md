---
doc_type: goal-audit
id: A-025-r3-c3-a024-remediation-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: finding-closure
scope: workspace-033 R3 C3 A-024 修复后复审（当前 HEAD 与工作树、fa0caa70 源码/测试钉、A-023 F-001/F-002；不采信 A-024/E-014 为成功依据）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-025 · R3 C3 A-023 F-001/F-002 修复后独立复审（2026-09-05）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C3 A-024 修复后复审（当前 HEAD `279f0298f5d7ddb4a1fb43fa1cfad3f6fc861120`；修复提交 `fa0caa70`；对照 A-023 recommended F-001/F-002、D-009 v0.2.0 L126–139；独立核对当前源码、测试与本会话跑数；**A-023 F-001/F-002 是否可在响应侧记为 fixed；本条不关闭 C3**）
- **verdict**：pass
- **open_required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。**A-001～A-024 原文及其 findings 全部保留、未改写。** 不把 A-022/A-024 self、E-013/E-014 或任何更早审计当作成功依据。不把 recommended 升级为 required，不接受 residual，不 overrule。不自行关闭 C3。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：本会话 `git rev-parse HEAD` = `279f0298f5d7ddb4a1fb43fa1cfad3f6fc861120`（`docs(govern): record C3 recommended remediation`）。其父提交为 `fa0caa70`（`fix(telegram): close C3 recommended audit gaps`；8 files, +533/−8，含 `http_sender.go`、`telegram_operator.go` 与测试钉）。`279f0298` 只改 GOAL-004 的 E-014/A-024/索引。工作树干净。用户描述的两个提交与 HEAD 核对一致；本条以当前 HEAD 与工作树为准，不信任描述或 A-024 结论。
- **covered**：
  1. A-023 F-001 七项验证分母钉是否存在、失败即红、断言是否打到目标（composition mux 匿名 401、真实 service credential 缺 operator scope 403、page/pageSize、runtime 变体、未知 retry、gated PG 同 request/同 root 并发）
  2. A-023 F-002：HTTPSender 空 token 是否仍有无 runtime/无 fallback 的 `nil` 成功；CaptureSender fallback 与调用方契约；operator `Send()==nil` 后二次 runtime/token 检查是否在 send/retry 两条路径生效；pending 是否阻止重复外发；可证明残余窗口
  3. 本会话专项、隔离 `-race`、相关包、gated PostgreSQL 跑数（未把 skip 当通过）
- **excluded**：改写 A-001～A-024；采信 A-024；关闭 C3 检查点或改 `00-meta`/`goal-tree`；C4 UI/`getChatMember`/发言权；全仓 `go test ./...` 或全量 `./internal/handler -race`

## 成果（有证据）

| 主张 | 本条独立证据（不引用 A-024 结论） |
|------|----------------------------------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| Charter `active` 0.4.0 | `docs/vision/charter.md` L5–6 |
| HEAD 为治理记录；修复在 `fa0caa70`；无未提交 diff | `git rev-parse HEAD`；`git show --stat fa0caa70`；`git status` 空 |
| A-023 原文仍为 `pass` / `open_required: 0`；F-001/F-002 原件仍 recommended/open | `A-023-r3-c3-implementation-independent.md` L10–11、L202–229；本条不改写 |
| A-024 原文仍为 self `pass`，声称 F-001/F-002 `fixed` | `A-024-r3-c3-a023-response.md` L10–11、L24–50；**不是本条证据** |
| D-009 仍要求空 token 且无明确 CaptureSender 的 `nil` 不得当 sent | `D-009-r3-c3-operator-console-contract.md` L126–130 |
| `fa0caa70` 将无 runtime 的 `return nil` 改为 `ErrTelegramTokenMissing`；send/retry 在 `Send` 后增加 `senderReady()` | `git show fa0caa70 -- http_sender.go telegram_operator.go` |
| 本会话定向/race/包回归 **PASS**；gated PG **PASS（未 skip）** | 见下方「本会话验证」 |

### 本会话验证（独立执行，2026-09-05）

在 `apps/api`，HEAD `279f0298`：

| 命令 | 结果 | 归类 |
|------|------|------|
| `go test ./internal/channel/telegram -run 'TestHTTPSender_UnconfiguredToken_DowngradesToMock\|TestHTTPSender_UnconfiguredWithoutRuntimeFailsClosed' -count=1 -v` | **PASS**（1.609s；两测均 `--- PASS`） | 通过 |
| `go test ./internal/handler -run 'TestTelegramOperatorHandlerSessionsTimelineSendAndRetry\|TestTelegramOperatorHandlerKeepsPendingWhenSentFinalizationFails\|TestTelegramOperatorHandlerKeepsPendingWhenTokenDisappearsAfterSend\|TestTelegramOperatorHandlerPaginationRuntimeAndRetryGuards\|TestTelegramOperatorMiddlewareRejectsServiceCredentialWithoutOperatorScope' -count=1 -v` | **PASS**（2.135s；五测均 `--- PASS`） | 通过 |
| `go test ./internal/composition -run 'TestTelegramChannelComposition_RealWebhookMount' -count=1 -v` | **PASS**（2.313s；`--- PASS`） | 通过 |
| `go test ./modules/channel/telegram/store -run 'TestRepositoryCreatePendingConcurrentRequestSQLite' -count=1 -v` | **PASS**（1.234s；`--- PASS`） | 通过 |
| `go test ./internal/store -run 'TestPostgresTelegramIngressRepositoryIdempotency\|TestPostgresTelegramOutboundConflictAndRetryState\|TestPostgresTelegramOutboundConcurrentRequestAndRootIdempotency' -count=1 -v` | **PASS**（ingress 1.99s；outbound 1.97s；concurrent 2.02s） | **gated PostgreSQL 通过**，不是 skip |
| 上述 HTTPSender / operator / composition / SQLite concurrent 隔离 `-race -count=1` | **PASS**（telegram 3.759s；handler 11.829s；composition 5.041s；store 4.278s） | 通过 |
| 上述三条 PG 测试 `-race -count=1` | **PASS**（store 10.630s） | **gated PostgreSQL -race 通过**，不是 skip |
| `go test ./internal/store ./internal/handler ./internal/composition ./internal/channel/telegram ./modules/channel/telegram/... ./kernel -count=1` | **PASS**（store 47.238s；handler 30.606s；composition 22.629s；channel/telegram 6.299s；module 0.643s；telegram/store 1.850s；kernel 0.670s；manifest/migration/schema 为 `no test files`） | 包回归通过。`ok store` **不能**单独证明 PG 未 skip；PG 证据以上表 verbose 为准 |

未把 skip 记为通过。本机 `apps/api/configs/.env` 经 `internal/pgtest` 自动加载 `PG_TEST_*`，上述 PG 测试实际执行。未重跑全仓 `go test ./...` 或全量 `./internal/handler -race`（A-023 记载的 wallet/SQLite `database locked` 不在本定向面，也不构成本条 fail 或 pass）。

## 对照成功标准

本条审的是 **A-023 recommended F-001/F-002 在当前 HEAD 是否已有可核对修正**。C3 检查点关闭权在 `/govern`，不在本条。GOAL 级 UI/发言权属 C4，不写成已交付。

### 1) A-023 F-001 · 验证分母测试钉

A-023 F-001（recommended / low / 原件 **open**）列出七项「必须能单独失败」的回归锁。当前 HEAD 与本会话跑数：

| # | A-023 要求 | 当前钉 | 本条判定 |
|---|------------|--------|----------|
| 1 | 经 composition mux 打四条 operator 路由的匿名 401 | `TestTelegramChannelComposition_RealWebhookMount`（`composition_telegram_test.go` L612–619）对真实 `newMuxWithExtraProviders` 发匿名 `GET /api/channel/telegram/operator/sessions`，断言 HTTP 401。四条路由共用已包装的 `tgOperator`（`composition.go` L631–645；`provider.go` L140–157 `Public: false` + 同一 `p.operatorHandler`）。**只打了 GET sessions**，只断言状态码、不断言 `UNAUTHENTICATED`。handler 直挂匿名同样是 401（`resources.go` L324–326），故该钉不能单独证明 Middleware；服务凭据 403 钉可以 | **存在且失败即红**（非 401 会 Fail）；覆盖偏窄，见本条 recommended F-001 |
| 2 | 服务凭据缺 `telegram.operator.*` → 403 且非 409 | `TestTelegramOperatorMiddlewareRejectsServiceCredentialWithoutOperatorScope`（`telegram_operator_test.go` L399–429）：真实 `auth.Authenticator.Middleware` + `POST /api/service-credentials` 创建仅 `users.read` 的凭据，再用 secret 打 GET sessions，断言 `403 FORBIDDEN`。无 Middleware 时 Identity 不会被写入，handler 会走 401；runtime 先于权限则会 409。本会话 **PASS** | **满足** |
| 3 | HTTP `page`/`pageSize` 非法与 `pageSize>100` → `INVALID_PAGE` / `INVALID_PAGE_SIZE` | `TestTelegramOperatorHandlerPaginationRuntimeAndRetryGuards` L355–367：`page=0` → `INVALID_PAGE`；`pageSize=101` 与 `pageSize=not-a-number` → `INVALID_PAGE_SIZE`；源码 `parseTelegramOperatorPagination` L445–468。`page=not-a-number` 与 `pageSize` 共用 `parsePositiveOperatorParam` | **满足** |
| 4 | runtime `unconfigured` / `receiver=none` / `bot_id<=0` / `receiver=webhook` | 同测：`webhook` → 200 可用（L369–373）；`none` → 409（L374–378）；`botID=0` → 409（L379–384）。既有测已钉 `idle` 与 `HasBusinessHandlers`（L240–251）。`state=unconfigured` 未具名；`runtimeAvailable` 对任何非 `running` 走同一分支（L166–167） | **满足**（unconfigured 与 idle 同分支） |
| 5 | pending 写入后 token 变空：`senderReady` 失败 → failed + 409，且不调用 sender | 同测 L386–391：token 在请求前清空，断言 409 `TELEGRAM_OPERATOR_UNAVAILABLE` 且 `callCount` 不变。`runtimeAvailable` 不读 token，故会进入 `CreatePending` 后再 `senderReady`（L331–337 调 `MarkFailed`）。**测试未读 durable 行是否 `failed`** | **HTTP/sender 钉满足**；durable `failed` 未锁，见本条 F-001 |
| 6 | 重试未知 request → `TELEGRAM_REQUEST_NOT_FOUND` | 同测 L393–396：对 `missing-request` retry 断言 404 且无 sender。本会话 **PASS** | **满足** |
| 7 | gated PG 同 request 并发与同 root 并发 pending | `TestPostgresTelegramOutboundConcurrentRequestAndRootIdempotency`（`postgres_telegram_test.go` L140–216）：8 路同 `request_id` `CreatePending` 恰好 1 个 created；`MarkFailed` 后 8 路不同 retry id 的 `CreateRetry` 恰好 1 个 pending root。本会话 verbose **PASS 2.02s，未 skip**；`-race` 亦 PASS。既有 SQLite `TestRepositoryCreatePendingConcurrentRequestSQLite` 本会话 PASS | **满足** |

**F-001 响应侧结论：fixed。** 原件 A-023 F-001 仍 open，不改写。七项均有失败即红钉且本会话绿；残余紧密度见本条 recommended F-001，不阻断、不升 required。

### 2) A-023 F-002 · HTTPSender 空 token 与 Send nil 接缝

A-023 F-002（recommended / low / 原件 **open**）描述：无 mock 时 `HTTPSender` 空 token 仍 `return nil`；operator 把 `Send()==nil` 当成功，可在 `senderReady` 与 `Send` 之间的卸载窗口 `MarkSent`。

| 项 | 判定 | 证据 |
|----|------|------|
| 无 runtime、无 fallback：不再 `nil` 成功 | **满足** | `http_sender.go` L90–94 返回 `ErrTelegramTokenMissing`。`fa0caa70` 将旧 `return nil` 删掉。`TestHTTPSender_UnconfiguredWithoutRuntimeFailsClosed` 断言 `errors.Is(err, ErrTelegramTokenMissing)`；本会话 PASS。该钉在恢复 `return nil` 时会红 |
| 显式 CaptureSender fallback 保留，且与调用方一致 | **满足** | 空 token 且 `runtime.Mock() != nil` 时走 mock（L91–92）。`TestHTTPSender_UnconfiguredToken_DowngradesToMock` PASS。生产 `composition.go` L903–908 **显式** `NewCaptureSender()` 再交给 `NewRuntimeManagerWithSettings` / `NewHTTPSender`。`RuntimeManager` 构造器在 `mock==nil` 时也会 `NewCaptureSender()`（`runtime.go` L100–102），故「有 runtime、无 mock」分支在现有构造器下不可达；这与「保留显式 fallback、避免扩大非 operator 调用方行为」一致，不是无 runtime 的 `nil` 成功 |
| operator send 路径：`Send()==nil` 后再确认 runtime/token | **满足** | `telegram_operator.go` L350–356：`Send` 成功后若 `!senderReady()` 则 409 且 **不** `MarkSent`。`TestTelegramOperatorHandlerKeepsPendingWhenTokenDisappearsAfterSend`：sender `before` 清空 token，断言 409、一行 pending、`callCount==1`；恢复 token 后重放 409 `TELEGRAM_REQUEST_IN_PROGRESS` 且无第二次 Send。本会话 PASS |
| operator retry 路径：同样的二次检查 | **源码满足；无独立钉** | `retryMessage` L413–418 与 send 对称。没有「retry 发送窗口 token 消失 → pending」测试。见本条 F-001 |
| pending 阻止重复外发 | **满足** | 上述 token-race 重放；以及既有 `TestTelegramOperatorHandlerKeepsPendingWhenSentFinalizationFails` L299–304。`CreatePending` 对已有 pending 返回 `ErrRequestInProgress`（`repository.go` L633–634） |
| 空 token 进入发送前即失败、零 sender | **满足** | `PaginationRuntimeAndRetryGuards` L386–391；`senderReady` L175–181 |

**可证明残余（不升 required）：**

1. **retry 的 Send-nil 窗口无测试钉**（源码有对称检查）。
2. **HTTPSender 在 RuntimeManager 上对空 token 仍走 CaptureSender 并返回 nil**。D-009 允许「明确的测试 CaptureSender」；生产 composition 始终注入该 fallback。operator 在 `Send` 前/后都查 token，空 token 不会 `MarkSent`。这不是 A-023 所写的「无 runtime/无 mock 的 nil 成功」。
3. **理论窗口**：`Send` 期间 token 被清空（HTTPSender 落入 mock、未打 Bot API），若在 `senderReady()` 二次检查前 token 又被写回，operator 仍会 `MarkSent`。本会话没有测试能证明该窗口；运行时 `running` 门禁使其很窄。不把它写成 required 缺陷。

**F-002 响应侧结论：fixed。** 原件 A-023 F-002 仍 open，不改写。无 runtime 的 `nil` 成功已删除；operator send 路径的卸载窗口已 fail-closed 且有失败即红钉；retry 路径源码对称。残余钉见本条 F-001。

### 3) 信息门禁（P-005）

| 项 | 最晚阶段 | 对本条 |
|----|----------|--------|
| I-033-021 | C1/C3 | 决策+合同已 verified。本条独立看到 composition mux 401、真实 service credential 403、Middleware→权限顺序。不改 `00-meta` |
| I-033-022 | C1/C3 | 决策+合同已 verified。本条独立看到 HTTPSender fail-closed、pending 卡住、PG 并发 PASS。不改 `00-meta` |
| I-033-009 / I-033-010 | C3/C4、C1/C4 | UI / `getChatMember` 属 C4；本条未写成已交付 |

无到期且影响本 scope 的 required 信息项被静默当作已实现。无共享资料引用。

### 4) 治理投影

| 投影 | 本条 |
|------|------|
| A-024 self「F-001/F-002 fixed」与 `fa0caa70` 描述 | **不是证据**。独立核对后同意响应侧 fixed；不同意用 self 代替本条 |
| E-014「专项测试通过」 | 本会话重跑后确认；不以 E-014 文本代替跑数 |
| 用户描述 HEAD=`279f0298` / 修复=`fa0caa70` | 核对后一致。本条以 `rev-parse` 为准 |
| 「PG skip 即通过」 | **未发生在本条**；三条 gated 测试 verbose 均为 `--- PASS` |
| A-023 原文 F-001/F-002 open | **保留**。本条只在响应侧记录 fixed |

## Findings

### 必改（required）

无。开放 required = **0**。

### 建议（recommended）

#### F-001 · retry token 窗口与部分分母钉仍偏松

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：A-023 F-001 项 1/5；A-023 F-002；I-033-022
- 描述：生产接缝与 A-023 所列主钉已落地且本会话绿，不阻断 C3。仍有三处紧密度缺口，不能把「读源码」或相邻钉当成同级回归锁：
  1. `retryMessage` 在 `Send()==nil` 后的 `senderReady()` 检查（`telegram_operator.go` L413–418）没有独立失败即红钉；send 路径已有 `TestTelegramOperatorHandlerKeepsPendingWhenTokenDisappearsAfterSend`
  2. 空 token 进入发送前的 409 钉未断言 durable 行为 `failed`（源码 L331–337 调 `MarkFailed`）
  3. composition mux 匿名 401 只覆盖 GET sessions，只断言 HTTP 401；因 handler 直挂匿名同样 401，该钉不能单独证明 Middleware（服务凭据 403 钉已证明认证后缺 scope）
- 证据：`telegram_operator.go` L331–337、L413–418；`telegram_operator_test.go` L307–338、L386–391；`composition_telegram_test.go` L612–619；`resources.go` L324–326
- 建议：C4 前按需补 retry token-race 钉，以及空 token 路径的 durable `failed` 断言。不升 required。

A-023 / A-024 原件 findings **全部保留**。本条不把它们改写成从未提出，也不在原件上改状态。

## 必改项汇总

| ID | 级别 | 阻断 C3 关门 |
|----|------|----------------|
| （无） | — | — |
| 本条 F-001 | recommended / low | **否** |
| A-023 F-001 | recommended / low | **否（响应侧 `fixed`；原件 open）** |
| A-023 F-002 | recommended / low | **否（响应侧 `fixed`；原件 open）** |

开放 required = **0**。本条不把任何 finding 标为 `accepted-residual` 或 `user-overruled`。

## 仍存在的 findings（台账）

| 条目 | 级别 | 状态 | 说明 |
|------|------|------|------|
| A-023 F-001 | recommended / low | 原件 **open**；响应侧 **fixed** | 七项钉已存在且本会话绿；原意见保留 |
| A-023 F-002 | recommended / low | 原件 **open**；响应侧 **fixed** | 无 runtime `nil` 成功已删除；send 路径 pending 卡住已测；原意见保留 |
| A-025 F-001 | recommended / low | open | 本条新增紧密度；不阻断 C3 |
| A-018 F-001～F-007 | 见 A-018/A-020/A-023 | 原件状态保留 | 本条不重审合同原件 |
| A-020 F-001～F-003 | recommended / low | 原件 open | 本条不重审 |

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| **A-023** independent `pass` / F-001/F-002 recommended | **原文不改、findings 原文保留。** 同意当时 HEAD `88d20ea1` 上两项为 recommended。本条审的是 `279f0298`（修复 `fa0caa70`）是否补上那些钉与接缝 |
| **A-024** self `pass` / 声称 fixed | **不作为证据。** 独立核对后同意 F-001/F-002 响应侧 fixed；不同意用 self 或 E-014 代替本条；本条另记 recommended 紧密度 |
| D-009 L126–130 | 对照合同，不改写。无 runtime/无明确 CaptureSender 的 `nil` 成功已消除；显式 CaptureSender fallback 仍在 |

无 self/independent 对同一必改项的一要一否冲突。无需 P-004 裁 residual。

## 结论 + 建议给编排器/用户的下一步

**verdict：pass。A-023 F-001/F-002 在响应侧 fixed。本 scope 无未关闭 high/med required。本条不关闭 C3 检查点。**

当前 HEAD `279f0298` 的生产修复（`fa0caa70`）使 HTTPSender 在无 runtime 时返回 `ErrTelegramTokenMissing`；operator 在 send/retry 两条路径于 `Send()==nil` 后再次确认 runtime/token，不确定时保持 durable pending 并阻止重放。A-023 所列 composition mux 401、真实 service credential 403、分页、runtime 变体、未知 retry、PG 同 request/同 root 并发钉均存在、失败即红，且本会话（含 gated PostgreSQL 与隔离 `-race`）为 PASS 而非 skip。残余为 retry token 窗口缺独立钉、空 token 路径未锁 durable `failed`、composition 401 覆盖偏窄——recommended，不阻断。

**建议 `/govern`：** 响应本条 `pass`；将 A-023 F-001/F-002 在响应侧记为 `fixed`（不要改写 A-023 原文）；关闭 C3 检查点并放行 C4。不要用 A-024 替代本条。不要把本条 recommended F-001 当成阻断。不要改写 A-001～A-024。本条不修改 status/progress。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
