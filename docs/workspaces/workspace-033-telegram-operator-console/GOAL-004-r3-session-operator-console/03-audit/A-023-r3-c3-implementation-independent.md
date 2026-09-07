---
doc_type: goal-audit
id: A-023-r3-c3-implementation-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: close-out
scope: workspace-033 R3 C3 实现关门（HEAD 与源码、测试、v69 outbound、operator routes/RBAC/runtime、会话/成绩单、幂等重试、D-010 lease、A-018 F-004～F-007；C3 检查点是否可关闭）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-023 · R3 C3 实现独立交叉审计（2026-09-05）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：close-out · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C3 实现关门（当前 HEAD `88d20ea112ebfe4cf02337b42190d88deea825bb`；实现提交 `7ddc97e1 feat(telegram): add C3 operator messaging surface`；对照 D-002/D-008/D-009 v0.2.0/D-010；独立核对当前源码、迁移、测试与 Git HEAD；用户书面要求把 A-018 F-004～F-007 纳入本实现审计范围；**C3 检查点是否可关闭**）
- **verdict**：pass
- **open_required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。**A-001～A-022 原文及其 findings 全部保留、未改写。** 不把 A-017/A-019/A-021/A-022 self、E-013 或任何更早审计当作成功依据。不接受 residual，不 overrule。不自行关闭 C3。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：本会话 `git rev-parse HEAD` = `88d20ea112ebfe4cf02337b42190d88deea825bb`（`docs(govern): record C3 implementation evidence`）。其父提交为 `7ddc97e1`（22 个 `.go` 文件，+2029/−79）。`88d20ea1` 只改 GOAL-004 的 E-013/A-022/索引。工作树干净。用户描述的两个提交与 HEAD 核对一致，本条以当前 HEAD 与工作树为准，不信任描述本身。
- **covered**：
  1. v69 `telegram_outbound_messages` SQLite/PostgreSQL DDL、迁移描述符、identity/fingerprint catalog、pending-root partial unique
  2. 四条 operator routes 的 composition `a.Middleware` 包装、`Public: false` 与认证的区别、401→403→409、descriptor/profile/permission 同步
  3. runtime gate、webhook/polling receiver、bot/chat 隔离、未绑定 `HasBusinessHandlers`、D-010 lease `settings.read OR telegram.operator.read`
  4. 会话列表/统一成绩单、64-bit decimal string ID、分页/排序/输入约束
  5. pending 先持久化再 sender、`(bot_id,request_id)` 幂等、同 payload/冲突/in-progress、发送失败、MarkSent 失败保持 pending、无重复外发
  6. `retry_of`/`retry_root`、failed-only、root pending/sent 禁止、同一 root 单 pending
  7. A-018 F-004 descriptor/profile sync、F-005 catalog/敏感诊断、F-006 mux-safe request id、F-007 post-send fail-closed
  8. 测试是否真覆盖上述边界；区分通过、skip/gated、环境波动和缺陷
- **excluded**：改写 A-001～A-022；采信 A-022 为交叉成功证据；替用户 residual/overrule；关闭 C3 检查点或改 `00-meta`/`goal-tree`；C4 UI/`getChatMember`/发言权缓存；把全仓 handler `-race` 的既有 wallet/SQLite 争用写成 C3 缺陷或通过

## 成果（有证据）

| 主张 | 本条独立证据（不引用 A-022 结论） |
|------|----------------------------------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| Charter `active` 0.4.0；VP-033 `active`；`vision_ref` 对齐 | `docs/vision/charter.md` L5–6；VP-033 L5–7 |
| HEAD 为治理记录；实现在 `7ddc97e1`；无未提交生产 diff | `git log -5`；`git show --stat 7ddc97e1` 22 个 `.go`；`git show --stat 88d20ea1` 仅 4 个 docs；`git status` 空 |
| A-018 原文仍为 conditional / open_required=3；F-004～F-007 原件仍 recommended/open | A-018 L10–11、L139–185；本条不改写 |
| D-009 v0.2.0 与 D-010 用户裁决仍在 | D-009 L8、L23–40、L87–139；D-010 L6、L15–19 |
| v69 双方言表 + pending-root partial unique | `migration.go` L99–137、L193–251 |
| catalog head=69；fingerprint 含 `telegram_outbound_messages` | `identity.go` L93–116、L138；`identity_test.go` L41、L120；`migrate_test.go` L124–134、L192–193、L745 |
| composition 先 `a.Middleware` 再注入 Provider | `composition.go` L629–645 |
| operator `Public: false`；webhook 仍 `Public: true` | `provider.go` L140–155、L181–193 |
| 专用权限 + 默认策略；不进 mvp/admin/demo | `provider.go` L75、L159–175；`kernel/profile.go` L217；`provider_test.go` L137–148 |
| lease 接受 `settings.read` **或** `telegram.operator.read`；settings API 仍只认 `settings.*` | `lease_handler.go` L39–40；`settings_handler.go` L38–48 |
| 方言无关 `INSERT ... ON CONFLICT DO NOTHING` + `RowsAffected` | `store/repository.go` L608–660 |
| 本会话定向测试通过；gated PG **PASS（未 skip）** | 见下方「本会话验证」 |

### 本会话验证（独立执行，2026-09-05）

在 `apps/api`：

| 命令 | 结果 | 归类 |
|------|------|------|
| `go test ./modules/channel/telegram/store ./internal/handler -run 'TestRepositoryOperatorProjectionAndOutboundStateSQLite\|TestRepositoryOutboundRejectsUnsafeRequestIDSQLite\|TestRepositoryCreatePendingConcurrentRequestSQLite\|TestTelegramOperatorHandler' -count=1` | **PASS**（store 1.386s；handler 1.723s） | 通过 |
| `go test ./internal/store ./internal/handler ./internal/composition ./internal/channel/telegram ./modules/channel/telegram/... ./kernel -count=1` | **PASS**（store 47.502s；handler 30.889s；composition 25.835s；channel/telegram 6.382s；module 0.660s；store 包 1.903s；kernel 0.685s） | 通过 |
| 同上 C3 专项 `-race` | **PASS**（store 7.674s；handler 6.901s） | 通过 |
| `go test ./internal/store -run 'TestPostgresTelegramOutboundConflictAndRetryState\|TestPostgresTelegramIngressRepositoryIdempotency' -count=1 -v` | **PASS**（ingress 2.22s；outbound 2.06s） | **gated PostgreSQL 通过**，不是 skip |

未把 skip 记为通过。A-022 记载「无 PostgreSQL 环境时 skip」不作为本条证据；本机 `PG_TEST_*` 可用，outbound 冲突/重试路径实际执行。未重跑全仓 `go test ./...` 或全量 `./internal/handler -race`（A-022 记载的 wallet/SQLite `database locked` 不在 C3 定向面，也不构成本条 fail 或 pass）。

## 对照成功标准

C3 检查点（`00-meta`）：会话列表/成绩单/人工发送 API、权限与运行时接线。合同权威为 D-002、D-008、D-009 v0.2.0、D-010。GOAL 级成功标准中的 UI/发言权/`getChatMember` 属 C4，本条不写成已交付。

### 1) v69 outbound DDL、描述符、fingerprint、pending-root unique

| 项 | 判定 | 证据 |
|----|------|------|
| 表 `telegram_outbound_messages`；PK `(bot_id, request_id)` | **满足** | `migration.go` L100–112、L120–132 |
| 列：`retry_root`、可空 `retry_of`、`chat_id`、纯文本、`status` CHECK `pending\|sent\|failed`、脱敏 `error_message`、`created_at`/`updated_at` | **满足** | 同上；无 token/secret/raw JSON 列 |
| SQLite `INTEGER` / PostgreSQL `BIGINT`；其余语义一致 | **满足** | L99–117 vs L119–137 |
| `CREATE UNIQUE INDEX ... (bot_id, retry_root) WHERE status = 'pending'` 双方言同一形状 | **满足** | L115–116、L135–136 |
| 成绩单索引 `(bot_id, chat_id, created_at DESC, request_id DESC)` | **满足** | L113–114、L133–134 |
| 描述符 v69 `telegram_outbound`；checksum 冻结 SQLite DDL | **满足** | L218、L244–251；`migrate_test.go` L745 = `76f4fa39c39d796ec8f106ae08d152526216d3781d9db7689fee1273eb2c974d` |
| fingerprint head=69；restore 表清单含 outbound | **满足** | `identity.go` L93–116；`identity_test.go` L41、L120；fresh/reopen/upgrade 尾断言 `migrate_test.go` L124–125、L192–193；`migrate_telegram_upgrade_test.go` L48 |
| 插入使用无目标 `ON CONFLICT DO NOTHING`（占位符 `?`），覆盖 PK 与 partial unique | **满足** | `repository.go` L610–613。A-020 F-002 的实施陷阱未发生 |

checksum 仍只哈希 SQLite DDL（与 v66–v68 同一模式），不单独构成本条 required。本会话 gated PG outbound 测试 **PASS**，证明 `ApplyPostgres` 已落到可运行的 scratch DB。

### 2) 四条路由、Middleware、Public:false、401→403→409、descriptor/profile

| 项 | 判定 | 证据 |
|----|------|------|
| 四条路径/方法与 D-009 表一致 | **满足** | `provider.go` L140–155；`kernel/profile.go` L217 |
| `Public: false` 只是声明；composition 在注入前包 `a.Middleware` | **满足** | `composition.go` L631–645；`auth.go` L590–607 失败写 `401 UNAUTHENTICATED`；`mux.Handle` 仍直接挂 `route.Handler`，与 settings/lease 同形 |
| handler 在 Identity 后查专用权限，缺身份 `401 UNAUTHENTICATED`，缺权限 `403 FORBIDDEN`，再 runtime `409 TELEGRAM_OPERATOR_UNAVAILABLE` | **满足** | `telegram_operator.go` L87–103、L161–173；`resources.go` `requirePermission` L322–331 |
| 匿名不得因 runtime 差异泄漏 | **满足（源码+handler 测）** | 先权限后 runtime；`telegram_operator_test.go` L127–135 匿名 401、无权限 403、只读发送 403；随后才测 idle/business 409（L240–251） |
| Descriptor / `profile.go` / `reg.Authorization` 同步；`PolicyAdminEditorViewer` / `PolicyAdmin` | **满足** | `provider.go` L63–83、L159–175；`policy.go` L7–9；`provider_test.go` L60–90、L170–181 |
| 不进 mvp/admin/demo | **满足** | `provider_test.go` L137–148；`composition_telegram_test.go` L33–42 |
| webhook 仍 Public；operator 不是 Public | **满足** | `provider.go` L154 vs L190 |
| 无新 page/navigation/fragment | **满足** | `provider.go` L76–81、L195–230 仍只贡献 `telegram-settings` / `menu_telegram` |

A-018 F-001 作为**实施合同缺口**已在 A-020 响应侧 `fixed`。本条确认运行时包装与顺序已按该合同句落地。原件 F-001 仍 `open`，不改写。服务凭据缺 scope：合法凭据过 Middleware 后 `Permissions = scopes`，缺 `telegram.operator.*` 走 handler 403，不是 401/409。源码相容；缺专用 mux 集成钉，见本条 F-001 recommended。

### 3) runtime gate、receiver、占用位、D-010 lease、bot/chat 隔离

| 项 | 判定 | 证据 |
|----|------|------|
| 可用条件：`running` + `bot_id > 0` + receiver ∈ {webhook,polling} + `!HasBusinessHandlers()` | **满足** | `telegram_operator.go` L161–173；`dispatcher.go` L30–36 |
| `bot_id` 来自 `ConnectionStatus()`，客户端不可覆盖 | **满足** | `composition.go` L632–638 |
| 不满足不调用 sender | **满足** | L100–102 在 send 之前；idle 测试 L241–245 `callCount` 不变 |
| 占用位 true → 409 | **满足** | L247–251 |
| D-010：lease 授权 `settings.read OR telegram.operator.read`；operator API 不自启 lease | **满足** | `lease_handler.go` L39–40；operator handler 无 `AcquireLease`；`lease_handler_test.go` L48–53 operator reader 200 |
| settings API 权限不变 | **满足** | `settings_handler.go` L38–48；菜单仍骑 `settings.read`（`provider.go` L218–219） |
| webhook running 不需要 lease；polling 靠既有心跳 | **满足（源码）** | operator 只读 ConnectionStatus；D-010 L23–27 |
| 查询绑定当前 bot；跨 bot/跨 chat 不串读 | **满足** | `repository.go` L417–418、L641–646、L692–703；`telegram_outbound_test.go` L118–120、L153–158 |

`HTTPSender` 仍只读 token（`http_sender.go` L76–87），发送不依赖 inbound loop。未绑定 polling 无 lease 时仍是 `idle`/`none` → operator 409，符合 D-010「不把 idle 当可用」。

### 4) 会话列表、成绩单、64-bit ID、分页/排序/输入

| 项 | 判定 | 证据 |
|----|------|------|
| 仅 inbound text/command 且非空文本的 session 可见；callback-only 不可见 | **满足** | `repository.go` L190–199、L232–285；SQLite 测试只返回 chat 8001；handler L137–144 |
| 成绩单 UNION inbound text/command + 全部 outbound；排除 callback/空文本 | **满足** | L318–339；handler 初始 timeline 一条 inbound（L146–154） |
| JSON 大整数用十进制字符串；时间为 UTC RFC3339 | **满足** | `telegram_operator.go` L190–196、L211–217、L259–266；测试 `chatId=="8001"`、`updateId=="9201"` |
| `{items,total,page,pageSize}`；默认 20、最大 100；非法 400 | **满足（源码）** | L432–444；`DefaultPageSize=20`、`maxPageSize=100`（`resources.go` L38–40）；`pagination.Offset` L254、L365 |
| 会话排序 `last_message_at DESC, chat_id DESC`；成绩单 `occurred_at DESC` + 来源 ident | **满足** | L259、L367–368 |
| `request_id` mux-safe `[A-Za-z0-9._-]{1,128}`；文本非空且 ≤4096 UTF-8 字节 | **满足** | handler L28、L486–501；store L59、L582–605；测试拒绝 `bad+id` / `../escape` / 空白 / `plus+sign` |
| 未知 chat 发送 404 `TELEGRAM_CHAT_NOT_FOUND` | **满足** | handler L222–225、L505–509 |

HTTP 层未单独钉非法 `page`/`pageSize` 与 JS 不安全大整数；实现使用 `strconv.FormatInt` 与既有溢出安全 offset。见 F-001。

### 5) pending 先于 sender、幂等、失败、MarkSent 失败、无重复外发

| 项 | 判定 | 证据 |
|----|------|------|
| 先事务写入 pending 并提交，成功才 `sender.Send` | **满足** | `sendMessage` L319–339；测试 `sender.before` 读到 pending 后再外发（L156–177） |
| pending 写入失败不调用 sender | **满足** | L320–326 在 Send 之前 return |
| 同 `(bot_id,request_id)` + 同 payload：pending → 409 in-progress；terminal → 200 不重发 | **满足** | `createPendingTx` L627–636；handler L179–185、L211–217；store L75–88 |
| 同 request 不同 payload → 409 conflict，不外发 | **满足** | L186–189；PG L119–121 |
| sender 错 → `failed` + `TELEGRAM_SEND_FAILED`，诊断为固定类别 | **满足** | handler L339–348、L527–544；store `safeErrorMessage` L731–741；测试 L191–199、L125–134（token/raw downstream 被收敛为 `telegram send failed`） |
| MarkSent 失败：5xx、行保持 pending、重放 409、不再外发 | **满足** | handler L350–354；`TestTelegramOperatorHandlerKeepsPendingWhenSentFinalizationFails` L287–304 |
| SQLite 8 路同 request 并发恰好 1 条 created | **满足** | `TestRepositoryCreatePendingConcurrentRequestSQLite` |
| PG 首次写入/payload 冲突/retry pending/sent 后禁止 | **满足（本会话 PASS）** | `TestPostgresTelegramOutboundConflictAndRetryState` |

调用 sender 前 `senderReady()` 再次要求 `runtimeAvailable()` 且 token 非空（L175–181、L331–337）。token 在探测后、Send 前被清空时，既有 `HTTPSender` 仍可能 `return nil` 被当成成功。这是 A-018 F-007 的剩余接缝，见本条 F-002；不升 required。

### 6) retry_of / retry_root、failed-only、单 pending

| 项 | 判定 | 证据 |
|----|------|------|
| 首发 `retry_root = request_id`、`retry_of` 空 | **满足** | `CreatePending` L467；store 测试 L71–73 |
| 重试新 request + 新 pending 行；`retry_of` 指向 chain root | **满足** | `CreateRetry` L517；handler 测试 `retryOf=="failed-1"`（L202–205）；store L97–99 |
| 只接受 failed；pending/sent 源 → `TELEGRAM_RETRY_NOT_ALLOWED` | **满足** | L500–502；root 已 sent 再查 L503–509 |
| 同 root 已有 pending → `TELEGRAM_REQUEST_IN_PROGRESS` | **满足** | L510–516；SQLite L100–102；PG L129–131；handler 对 pending 重放 L211–217 |
| 同 root 任一 sent → 禁止再 retry | **满足** | handler L206–209；PG L135–137 |
| 无后台 worker/定时器/自动重试 | **满足** | telegram 包无 outbound ticker；每次尝试都是显式 HTTP |
| 正文/chat 从失败行派生，客户端不能改 | **满足** | `CreateRetry` 用 `source.Text` / `source.RetryRoot`；retry body 只有新 `requestId` |

D-008「retry_of 关联原始请求」按 D-009 解释为 chain root，与 A-018 当时接受的合同解释一致。实现未退回原地 `failed→pending`。

### 7) A-018 F-004～F-007（用户要求纳入本实现范围）

| A-018 | 本条独立判定 | 证据 |
|-------|--------------|------|
| F-004 Descriptor + `profile.go` + 权限贡献；不进默认 Profile | **运行时已落地** | 见 §2；`provider_test.go` / `composition_telegram_test.go` |
| F-005 稳定 catalog + 未知 chat/request + 中英条目 | **运行时已落地** | `errorcatalog.go` L217–224 七码（含 `TELEGRAM_CHAT_NOT_FOUND` / `TELEGRAM_REQUEST_NOT_FOUND` / `TELEGRAM_SEND_FAILED`）；`error_contract_test.go` L82–84；handler 映射 L505–524 |
| F-006 mux-safe `request_id`；4096 **字节**上限 | **运行时已落地** | 正则 L28 / store L59；非法 id 测试 handler L226–229、store L161–172 |
| F-007 persist-fail-after-send fail-closed；sender 前再确认 token；HTTPSender 空 token `nil` 不得当 sent | **主路径已落地；空 token `nil` 仍是接缝** | MarkSent 失败测试 L263–304；`senderReady` L175–181；`http_sender.go` L82–87 **未改**，见本条 F-002 |

A-018 原件 F-004～F-007 仍为 recommended/open。本条不在原件上改状态。响应侧合同纳入仍以 A-020 为准。本条确认其中 F-004/F-005/F-006 与 F-007 的 pending 卡住语义已在代码中可核对。

### 8) 信息门禁（P-005）

| 项 | 最晚阶段 | 对本条 |
|----|----------|--------|
| I-033-021 | C1/C3 | 决策+合同已 verified；本条独立看到权限键、Middleware、lease OR、Descriptor/profile 的代码与测试。不改 `00-meta` 信息表 |
| I-033-022 | C1/C3 | 决策+合同已 verified；本条独立看到 pending→sent/failed、ON CONFLICT、retry_of/root、PG PASS。不改 `00-meta` |
| I-033-009 | C3/C4 | UI 属 C4；D-009 未越界关闭 |
| I-033-010 | C1/C4 | `getChatMember` 属 C4；本条未实现预检 |

无到期且影响本 scope 的 required 信息项被静默当作已实现。无共享资料引用。

### 9) 治理投影

| 投影 | 本条 |
|------|------|
| A-022 self `pass` | **不是证据**。独立核对后，核心实现主张与本条一致；A-022 把本机当时的 PG skip 写成验证事实，本会话 PG 实际 PASS，二者不得混写 |
| E-013「专项测试通过」 | 本会话重跑后确认；不以 E-013 文本代替跑数 |
| 「PG skip 即通过」 | **未发生在本条**；gated outbound 测试实际 PASS |
| 用户描述 HEAD=`7ddc97e1` / 治理=`88d20ea1` | 核对后：当前 HEAD 是 `88d20ea1`，实现父提交是 `7ddc97e1`。本条以 `rev-parse` 为准 |
| 全量 handler `-race` wallet/SQLite locked | 不在 C3 定向面；C3 隔离 `-race` PASS |

## Findings

### 必改（required）

无。开放 required = **0**。

### 建议（recommended）

#### F-001 · D-009 若干验证分母仍缺失败即红的测试钉

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-021；I-033-022；D-009 L132–139
- 描述：源码已实现下列语义，但缺少可单独失败的测试钉，不能把「读源码」或相邻 settings/lease mux 测试当成与现有 handler/PG 测试同级的回归锁。不阻断 C3：匿名/无权限 401/403、idle/占用位 409、pending 先于 send、terminal 不重发、conflict、failed 目录码、retry root、MarkSent 失败卡住、SQLite 并发、本会话 gated PG 冲突/重试，均已有测试且本会话绿。
  1. 经 composition mux 打四条 operator 路由的匿名 401（现只测 handler 直挂；settings 401 在同一 `a.Middleware` 包装函数里已测）
  2. 服务凭据缺 `telegram.operator.*` scope → 403 且非 409
  3. HTTP `page`/`pageSize` 非法与 `pageSize>100` → `INVALID_PAGE` / `INVALID_PAGE_SIZE`
  4. runtime `unconfigured` / `receiver=none` / `bot_id<=0` / `receiver=webhook` 变体（现钉 idle 与 HasBusinessHandlers）
  5. pending 写入后 token 变空：`senderReady` 失败 → failed + 409，且 HTTPSender 不被调用
  6. 重试未知 request → `TELEGRAM_REQUEST_NOT_FOUND`（store sentinel 已有，HTTP 未钉）
  7. gated PG 同 request **并发**与同 root 并发 pending（本会话 PG 为串行冲突/重试 PASS，不是 skip；SQLite 8 路同 request 已绿）
- 证据：`telegram_operator_test.go`；`composition_telegram_test.go` L604–649（settings/lease mux，无 operator 路径）；`postgres_telegram_test.go` L103–138；`telegram_operator.go` L175–181、L432–444、L510–511。
- 建议：C4 前把 1–4、6 补成失败即红；5 与 7 可与 F-002 一并处理。不升 required。

#### F-002 · `HTTPSender` 在空 token 且无 mock 时仍 `return nil`；operator 把 `Send()==nil` 当成功

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：A-018 F-007；I-033-022；D-009 L126–130
- 描述：D-009 要求调用 sender 前再确认 token/bot_id，且 HTTPSender 空 token 的 `nil` 不得被当成外部发送成功。operator 已做 `senderReady()`（token 空则 `MarkFailed` + 409，不调用 Send），post-send 持久化失败也已 fail-closed。但 `http_sender.go` L82–87 未改：若 `senderReady` 看到 token 之后、`Send` 读到空 token（卸载窗口），`nil` 仍会被 `MarkSent`。运行时 `running` 门禁使该窗口很窄。不阻断 C3。
- 证据：`http_sender.go` L82–87；`telegram_operator.go` L175–181、L331–358；`7ddc97e1` `--stat` 不含 `http_sender.go`。
- 建议：operator 将 `Send()==nil` 且无明确 CaptureSender 视为失败并 `MarkFailed`，或让 HTTPSender 在无 token/无 mock 时返回 error。

A-018 / A-020 原件 findings **全部保留**。本条不把它们改写成从未提出，也不在原件上改状态。

## 必改项汇总

| ID | 级别 | 阻断 C3 关门 |
|----|------|----------------|
| （无） | — | — |
| 本条 F-001 | recommended / low | **否** |
| 本条 F-002 | recommended / low | **否** |
| A-018 F-001 | required / high | **否（响应侧合同 `fixed` · A-020；本条确认 composition `a.Middleware` + 401→403→409 已落地）** |
| A-018 F-002 | required / high | **否（响应侧合同 `fixed` · A-020；本条确认 lease OR 与 running 门禁已落地）** |
| A-018 F-003 | required / high | **否（响应侧合同 `fixed` · A-020；本条确认无目标 `ON CONFLICT DO NOTHING` + 本会话 gated PG PASS）** |

开放 required = **0**。本条不把任何 finding 标为 `accepted-residual` 或 `user-overruled`。

## 仍存在的 findings（台账）

| 条目 | 级别 | 状态 | 说明 |
|------|------|------|------|
| A-018 F-001 | required / high | 原件 **open**；响应侧合同 **fixed**（A-020）；**本条确认运行时已按该合同句实现** | 原意见保留，不改写 A-018 |
| A-018 F-002 | required / high | 原件 **open**；响应侧合同 **fixed**（A-020）；**本条确认 lease OR 已实现** | 原意见保留 |
| A-018 F-003 | required / high | 原件 **open**；响应侧合同 **fixed**（A-020）；**本条确认 ON CONFLICT + 本会话 PG PASS** | 原意见保留 |
| A-018 F-004 | recommended / low | 原件 open；本条确认 Descriptor/profile/权限同步已落地 | 原意见保留 |
| A-018 F-005 | recommended / low | 原件 open；本条确认七码入 catalog，未知 chat/request 为 404 | 原意见保留 |
| A-018 F-006 | recommended / low | 原件 open；本条确认 mux-safe 正则与非法 id 测试 | 原意见保留 |
| A-018 F-007 | recommended / low | 原件 open；pending 卡住已测；HTTPSender `nil` 见本条 F-002 | 原意见保留 |
| A-020 F-001～F-003 | recommended / low | 原件 open；未知码已选 404；ON CONFLICT 无目标；mux-safe 已收紧 | 原意见保留 |
| A-023 F-001～F-002 | recommended / low | open | 本条新增；不阻断 C3 |
| A-003 仍开放 recommended | recommended | open | 本条不重审、不闭合 |

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| **A-018** independent `conditional` / `open_required: 3` | **原文不改、findings 原文保留。** 同意当时 HEAD `6f935eba` 上 F-001/F-002/F-003 阻断开工。本条审的是 `88d20ea1` 的**实现**是否满足那些闭合句，不是重开合同审计。 |
| **A-020** independent `pass` | **原文不改。** 同意当时合同 `fixed`、当时尚无 C3 代码。本条补的是 A-020 明确排除的实现审计。 |
| A-022 self `pass` | **不作为证据。** 独立核对后同意 v69/operator/RBAC/runtime/幂等/F-004～F-006 主路径已落盘；不同意把 self 或「PG skip」写成独立通过；本条另记 F-001～F-002，且本会话 PG 为 PASS。 |
| A-017 / A-019 / A-021 | 历史合同/响应记录；不当实现通过证据。 |
| D-002 / D-008 / D-009 / D-010 | 对照合同，不改写。专用权限、新 request+`retry_of`、Middleware、lease OR、ON CONFLICT 在代码中仍忠实。 |

无 self/independent 对同一必改项的一要一否冲突。无需 P-004 裁 residual。

## 结论 + 建议给编排器/用户的下一步

**verdict：pass。C3 实现关门范围内无未关闭 high/med required。C3 检查点可以关闭。**

当前 HEAD `88d20ea1` 的生产实现（`7ddc97e1`）忠实于 D-009/D-010：v69 双方言 outbound 与 pending-root 部分唯一、composition `a.Middleware`、401→403→409、专用权限与 profile 同步、lease OR、占用位 fail-closed、pending 先于同步 sender、`(bot_id,request_id)` 幂等、显式 retry 新行、MarkSent 失败卡住、本会话 SQLite/race **以及 gated PostgreSQL PASS**。A-018 F-004～F-007 中除 HTTPSender 空 token `nil` 接缝外，均已在代码中可核对。C4 UI 与发言权不在本条范围。

**建议 `/govern`：** 响应本条 `pass`；关闭 C3 检查点并放行 C4。不要用 A-022 替代本条。不要把本条 recommended F-001/F-002 当成阻断。不要改写 A-001～A-022。本条不修改 status/progress。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
