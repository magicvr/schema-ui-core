---
doc_type: goal-audit
id: A-018-r3-c3-contract-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: design-plan
scope: workspace-033 R3 C3 合同审视（D-002 专用权限与发送状态机、D-008 显式重试身份、D-009 人工台 API/权限/运行时合同；当前 Telegram provider/runtime/store/auth 接缝；v68 入站与 v69 outbound；C3/C4 与 VP-033 首波非目标；C3 是否可进入生产代码实施）
verdict: conditional
open_required: 3
version: 0.1.0
---

# A-018 · R3 C3 合同独立交叉审计（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：design-plan · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C3 合同审视（HEAD `6f935eba56a2326a9e9a70e1789086356147c963`；工作树含未提交 D-009/A-017；D-002 用户裁决；D-008 重试身份；D-009 人工台实施合同；A-017 self **不作为成功依据**；当前 provider/runtime/dispatcher/sender/store/auth；v68 入站；无 v69/operator 生产代码）
- **verdict**：conditional
- **open_required**：3（F-001、F-002、F-003）
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。A-001～A-017、D-002、D-008、D-009 原文均未改写。不把 A-017 self 或其他治理记录当作成功依据。不接受 residual，不 overrule。不把尚未存在的 C3 代码、v69 或 operator 测试写成已完成。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：`6f935eba56a2326a9e9a70e1789086356147c963`（`docs(govern): record C3 retry identity decision`）。工作树另有未跟踪/未提交的 D-009 与 A-017。源码中无 `telegram_outbound`、无 `telegram.operator.*`、无 `/api/channel/telegram/operator/`
- **covered**：
  1. D-002 专用 `telegram.operator.read/write` 与 D-008 新 `request_id` + `retry_of` 是否被 D-009 忠实保留
  2. 四条 operator 路由、认证顺序、未绑定且 `running`/`bot_id`/receiver 门禁、bot scope、64-bit JSON ID 与分页
  3. v69 SQLite/PostgreSQL outbound schema、`retry_root` pending 唯一并发边界、pending 先于 sender、成功/失败/持久化失败、同 request 幂等、显式重试与无自动重试
  4. C3/C4 边界与 VP-033 首波非目标
  5. D-009 遗漏、矛盾或不可实施之处；C3 是否可以进入生产代码实施
- **excluded**：C3～C4 生产实现（当前不存在，不得写成成功事实）；把 A-017 self `pass` 当交叉证据；改写 D-002/D-008/D-009/A-001～A-017；替用户 residual / overrule；闭合 A-003 仍开放的 recommended 项；I-033-009/010 的 UI/`getChatMember` 实现（属 C4）

## 成果（有证据）

| 主张 | 证据 |
|------|------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| Charter `active` 0.4.0；VP-033 `active`；`vision_ref` 对齐 | `docs/vision/charter.md` L5–6；VP-033 L5–7 |
| HEAD 为 D-008 提交；D-009 尚未进入 HEAD；无 operator/v69 业务 diff | `git rev-parse HEAD` = `6f935eba…`；`git status` 显示 D-009/A-017 未提交；源码检索无 `telegram_outbound` / `telegram.operator` / `operator/sessions` |
| D-002 以 `source: user` 冻结专用权限与发送状态机 | D-002 L6、L20、L23、L27–30 |
| D-008 以 `source: user` 冻结新 request + `retry_of`、无自动重试 | D-008 L6、L15–21 |
| D-009 承接四条鉴权路由、专用 read/write、pending→sent/failed、`(bot_id,request_id)` 幂等、retry 新行、无后台重试、C3 不建 page/nav/fragment | D-009 L21–107 |
| 现码占用位、连接状态、同步 sender、v68 入站与溢出安全分页均存在，可供 C3 接线 | `dispatcher.go` L27–36；`runtime.go` L20–41、L270–283；`kernel/telegram.go` L15–16、L101–105；`migration.go` L33–96、L153–185；`pagination.go` L36–42 |
| A-001～A-017 / D-002 / D-008 / D-009 原文仍在；本条未改写 | `03-audit/`、`01-decision/D-002-*.md`、`D-008-*.md`、`D-009-*.md` |

## 对照成功标准

### 1) D-002 专用权限与 D-008 重试身份是否被忠实保留

| # | 已裁决选择 | D-009 | 判定 |
|---|------------|-------|------|
| 1 | 专用 `telegram.operator.read` / `telegram.operator.write`，独立于 `settings.read/write` | 读路由 `telegram.operator.read` + `system.admin-editor-viewer`；写/重试 `telegram.operator.write` + `system.admin`；由 `channel.telegram` 贡献并进 system-data reconcile | **路由键忠实**；polling 可用性仍绑在 `settings.read` lease 上，见 F-002 |
| 2 | `pending` → `sent`/`failed`；同一 `request_id` 不重复外发；失败显式重试、无自动重试 | 先提交 pending 再调 `kernel.TelegramSender`；`(bot_id,request_id)` 幂等；terminal 不重发；pending 返回 `TELEGRAM_REQUEST_IN_PROGRESS`；无 worker/定时器 | **方向忠实**；PG 同事务冲突读法未写死，见 F-003 |
| 3 | D-008：每次显式重试新建 `request_id` 与新记录；`retry_of` 关联原始请求；首发无 `retry_of` | 重试路由只接受 failed；新 pending 行；`retry_of` 指向链的 root；`retry_root` pending 唯一；不原地改写失败行 | **忠实**（把「原始请求」解释为 chain root，落在 D-008「合同须定义引用约束」授权内） |

未选方案（复用 `settings.read/write`、原地 `failed→pending`、无关联新 request、后台自动重试）未被写进 D-009。设置页/lease 继续走既有 `settings.*`（`provider.go` L48–66、L153–163），operator 不改 menu/page/fragment，符合 VP-033「不进默认集」与 A-003 F-004 的设置页红线。

### 2) 四条路由、认证顺序、门禁、bot scope、64-bit ID、分页

四条路径与方法与 C3 最小面匹配：列表、成绩单、新发送、显式重试。`chat_id` / Telegram 大整数以十进制字符串返回，可避免 `RuntimeStatus.BotID int64 \`json:"bot_id"\``（`runtime.go` L53）那种浏览器精度损失。列表排序 `last_message_at DESC, chat_id DESC` 与 v68 索引一致（`migration.go` L45–46）。`{items,total,page,pageSize}`、默认 20、最大 100、非法 400、使用 `pagination.Offset`，与现有 `INVALID_PAGE` / `INVALID_PAGE_SIZE` 及 W8 溢出辅助相容。

认证顺序的**意图**正确：先认证 401 `UNAUTHENTICATED`，再权限 403 `FORBIDDEN`，再运行时 409，避免用 409 向无权限者泄漏 bot 状态。现码 **不会**因 `Public: false` 自动套认证：

- `RouteContribution.Public` 只是声明字段（`kernel/contribution.go` L27–34）；`validateRoute` 不读它（L157–164）
- composition 挂载是直接 `mux.Handle(full, route.Handler)`（`composition.go` L673–678）
- settings/lease 能 401，是因为注入 Provider **之前**包了 `a.Middleware`（`composition.go` L629–631）；Middleware 失败写 catalog `UNAUTHENTICATED`（`auth.go` L590–602）
- 现 Telegram handler 在 Identity 后再查权限，用的是 `http.Error` 纯文本，不是 `FORBIDDEN`（`settings_handler.go` L30–48；`lease_handler.go` L33–42）

D-009 只写「非 Public」+「先认证」，未冻结 composition 包装。按字面只改 `provider.go` 会把四条人工发送/读取路由挂成无认证。见 **F-001**。

运行时门禁：`running` + `bot_id > 0` + receiver ∈ {webhook,polling} + `HasBusinessHandlers()==false`。占用位探测是现成包级方法，不扩 kernel（`dispatcher.go` L27–36）。`bot_id` 来自 `ConnectionStatus()`，客户端不得覆盖，与 C2「getMe 确认的稳定 bot_id」一致。

缺口是把 `running` 当成「连接成功」。VP-033 开发默认 polling；未绑定且无控制台心跳时，现码在 `getMe`+`deleteWebhook` 成功后进入 **`idle` / `receiver=none`，仍保留 BotID**（`connection_manager.go` L285–295、L303–317）。polling 的 `running` 只由 `hasPollingDemand()` 产生：已绑定 **或** 仍有效的 lease（L445–452）。lease 路由要求 `settings.read`（`lease_handler.go` L39–41）。因此仅持有 `telegram.operator.*`、没有 `settings.read` 的主体，在默认 polling 下会永远 409，发送也不会发生——HTTPSender 并不依赖 inbound loop（`http_sender.go` L69–80）。这与 D-002「专用权限独立于 settings」和 VP-033「未绑定且已连接可进入」冲突。见 **F-002**。

### 3) v69 outbound、pending 先于 sender、幂等与显式重试

状态机顺序正确，且对齐 kernel「Send 同步、无队列、无后台重试」（`kernel/telegram.go` L15–16）：pending 提交成功才调用 sender；pending 写入失败不调用 sender；sender 错 → `failed` 且成绩单可见；状态更新失败不得伪装 `sent`、同一 request 不得再外发。`(bot_id,request_id)` 作幂等键；payload 不一致 `TELEGRAM_REQUEST_CONFLICT`；pending 中 `TELEGRAM_REQUEST_IN_PROGRESS`。重试只接受 failed、新 request、`retry_of`→root、同一 root 同时最多一个 pending。可实施的 DDL 形态在本仓已有先例：recycle-bin 双方言 partial unique index（`recyclebin/migration/migration.go` L30、L48）。

PostgreSQL 生产权威方言上，普通 `INSERT` 的唯一冲突会中止当前事务。本仓已用测试钉死：同事务继续执行必须走 `ON CONFLICT DO NOTHING`（`postgres_test.go` L1252–1270）。C2 入站最终合同把这一点写成 `RowsAffected` 分支（`store/repository.go` L61–89）。C3 的常路径同样要在**同一事务**里插入 pending 后读取已有行以比较 payload，或在 pending-root 冲突时返回 409 而不中止。D-009 只说「数据库约束与事务共同保证」，验证清单没有 PG 重复/并发运行时。按字面在 SQLite 可测绿、在 PG 上会把幂等重试和显式 retry 打成 5xx。`kernel.IsUniqueViolation` 只能映射整个 `Run()` 失败并回滚，不能在已中止 Tx 里 SELECT 比较 payload（`unique_violation.go` L22–36）。见 **F-003**。

### 4) C3/C4 边界与 VP-033 首波非目标

D-009 正确把 page/navigation/fragment、`getChatMember` TTL、composer、10 秒单飞/失焦留给 C4；不扩张 kernel Telegram port；不写历史回灌、媒体、FSM、群发、频道、多 bot、多实例 polling、独立进程、SSE/WebSocket。入站成绩单含 text/command、排除 callback/空文本/未建模媒体，与 v68 `message_kind` 及 VP-033「只文本 / 不把 callback 当成绩单成功条件」相容（命令作为已落盘收据展示，不是 C2 合同回退）。I-033-009/010 仍属 C4，本条不把它们标成 C3 可关闭。

未钉死的 C3/C4 交界只有 F-002 的 polling lease 权限：若不在 C3 合同写明，C4 要么继续用 `settings.read` 心跳（破坏专用权限），要么改 C3 门禁（破坏刚冻结的 409 语义）。

### 5) 信息门禁（P-005）

| 项 | 最晚阶段 | 状态 | 对本条 |
|----|----------|------|--------|
| I-033-021 | C1/C3 | verified (decision)；实现+测试未宣称完成 | 权限键选择忠实。C3 实施合同缺认证包装与 polling 可用性，不能把「decision verified」当成可开工 |
| I-033-022 | C1/C3 | verified (decision)；D-008 已钉死重试身份 | 状态机方向忠实。缺 PG 安全冲突读法，与 C2 A-008 F-001 同类 |
| I-033-009 | C3/C4 | verified (decision)；UI 测试属 C4 | D-009 未越界关闭 |
| I-033-010 | C1/C4 | verified (decision)；实现属 C4 | D-009 未实现预检；失败文本脱敏后须仍能让 C4 识别 403（recommended） |

无到期且未接受 residual 的 required 信息项被本条重新打开为「未知」。本条阻断的是 **C3 实施合同完整性**，不是把用户裁决打回未决。

## Findings

### F-001 · `Public: false` 不会认证；operator 路由必须在 composition 层包 Middleware

- 严重度：high
- 建议：required
- 状态：open
- 关联：I-033-021
- 描述：D-009 L23–25、L40–49 把四条路由写成「非 Public」且「先经过认证和专用权限检查」。现码里 `Public` 不是安全边界。settings/lease 的 401 来自 `a.Middleware(...)` 预包装。D-009 未要求：operator handler 在 `telegrammodule.New` 之前由 composition 包 `Authenticator.Middleware`；权限检查在 `IdentityFrom` 之后、运行时门禁之前；403 使用 catalog `FORBIDDEN` 而不是 `http.Error` 纯文本；匿名与缺权限不得因运行时状态返回不同代码。`Provider.New` 目前用可变参数依次接收 settings、lease（`provider.go` L26–40）；再塞第四个 handler 且不写死顺序，容易把未包装 handler 挂上网。
- 证据：`kernel/contribution.go` L27–34、L157–164；`composition.go` L629–631、L673–678；`auth.go` L590–602；`settings_handler.go` L30–48；`lease_handler.go` L33–42；D-009 L23–49。
- 为何阻断 C3：C3 的核心交付就是这四条鉴权 API。按字面只改 Provider 贡献会把人工发送面做成匿名可调用。**不要** residual / overrule。
- 建议闭合：补进 D-009：composition 对 operator handler 使用与 settings/lease 相同的 `a.Middleware` 包装后再注入 Provider；`Public: false` 不得被当成认证实现；检查顺序固定为 Middleware 401 → 专用权限 403 `FORBIDDEN` → 运行时 409；验证必须覆盖匿名、有会话无权限、服务凭据缺 scope，且这三类不得返回 `TELEGRAM_OPERATOR_UNAVAILABLE` 或会话数据。

### F-002 · `running`+receiver 门禁把默认 polling 的人工台绑回 `settings.read` lease

- 严重度：high
- 建议：required
- 状态：open
- 关联：I-033-021；VP-033 判据 5；C3/C4 边界
- 描述：D-002 选择专用 operator 权限，正是为了不让 settings 与人工台互相蕴含。D-009 却要求 operator surface 仅在 `state==running` 且 receiver 为 webhook/polling 时可用。未绑定 polling 在无 lease 时是 `idle`/`none`，BotID 已有；lease 目前只认 `settings.read`。D-009 验证清单把「非 running fail-closed」写成 C3 必测项，会把「operator 用户在开发默认 polling 下 409」测成成功。发送路径并不需要 inbound receiver：`HTTPSender` 只读 token。
- 证据：D-002 L20；D-009 L30–34、L102–106；VP-033 L46、L58、L96–98；`connection_manager.go` L285–295、L303–317、L445–452；`lease_handler.go` L39–41；`http_sender.go` L69–80。
- 为何阻断 C3：这不是 C4 才出现的 UI 细节，而是 C3 可用性门禁与用户已选权限模型的矛盾。不写进合同，实现要么让专用权限在 polling 下不可用，要么在 C4 偷偷把 lease 留在 `settings.read`。**不要** residual / overrule。
- 建议闭合：D-009 必须**二选一写死**（可测）：**(a)** C3 将 polling lease 权限扩展为 `settings.read` **或** `telegram.operator.read`（写操作仍要 `telegram.operator.write` 才发送），使 operator 用户能把未绑定 polling 推到 `running`；**(b)** C3 门禁改为「已连接」：`bot_id>0` 且 state ∉ {unconfigured, error}，允许 `idle`/`receiver=none`，发送不依赖 inbound loop；polling 启停仍由既有 lease 负责，并明确 C4 必须让 `telegram.operator.read` 足以持有 lease。绑定占用位 fail-closed 两案都保留。

### F-003 · outbound 幂等与 pending-root 唯一缺少 PostgreSQL 安全的 `ON CONFLICT` 读法

- 严重度：high
- 建议：required
- 状态：open
- 关联：I-033-022
- 描述：D-009 L80–99 要求 `(bot_id,request_id)` 唯一，相同 payload 返回已有记录（terminal 不重发，pending 409），不同 payload 409 conflict；并对 `(bot_id,retry_root)` 的 `pending` 建数据库级唯一。比较 payload 必须在插入冲突后 **同一事务** 读出行。PostgreSQL 上无 `ON CONFLICT` 的唯一冲突会中止 Tx。C2 已因此补过合同；C3 未继承该写法。pending 部分唯一索引在 SQLite/PG 都可行（recycle-bin 先例），但 PG 推断冲突目标必须带 `WHERE status = 'pending'`。验证清单只说「SQLite 与 gated PostgreSQL 的迁移/读写」和「同 request 并发」，没有要求冲突路径不得 5xx、不得回滚已读状态。
- 证据：D-009 L80–106；`postgres_test.go` L1252–1270；`store/repository.go` L61–89；`kernel/unique_violation.go` L22–36；`recyclebin/migration/migration.go` L30、L48。
- 为何阻断 C3：C3 的核心就是这条幂等/重试边界。客户端网络重试与显式 retry 是常路径。SQLite 测绿、PG 5xx 会把「同一 request 不重复外发」落实成「重复提交变成持久化失败」。**不要** residual / overrule。
- 建议闭合：补进 D-009 并列入必测项：pending 插入使用方言无关 `INSERT ... ON CONFLICT DO NOTHING`（占位符 `?`）；`(bot_id,request_id)` 冲突用 `RowsAffected()==0` 后同事务 SELECT 比较 chat/text/`retry_of`，不得在唯一失败后于同一 Tx 继续无目标语句；pending-root 冲突（含 partial unique）同样 `DO NOTHING` 后映射 `TELEGRAM_REQUEST_IN_PROGRESS`，禁止 `kernel.IsUniqueViolation` 当作「已成功接受」的同 Tx 控制流。DDL 写明 SQLite INTEGER / PG BIGINT、双方言 `CREATE UNIQUE INDEX ... WHERE status = 'pending'`。验证必须覆盖 SQLite **和** gated PostgreSQL 的首次写入、同 request 重放、payload 冲突、同 root 并发 pending、显式 retry。

### F-004 · Permission/Route 贡献必须同步 Descriptor 与 `kernel/profile.go`，且不进默认 Profile

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-021；A-003 F-004
- 描述：`descriptorsMatch` 要求 Provider `Descriptor()` 与 `BuiltinModules()` 的贡献键完全一致（`provider.go` L107–135）。当前 `channel.telegram` 声明 6 条路由、**零** Permissions（`profile.go` L216；`provider.go` L55–70）。C3 新增四条路由和两个权限键却不改这两处，`RegisterContributions` 会 `MODULE_API_MISMATCH` / undeclared key。`mvp`/`admin` 默认集必须继续不含本模块（`provider_test.go` L121–131）。这是实施清单，不是新的产品方案。
- 证据：上引；`provider.go` L92–110。
- 建议：D-009 验证清单点名：`Descriptor` + `profile.go` 同步 Routes/Permissions；`reg.Authorization` 使用 `PolicyAdminEditorViewer` / `PolicyAdmin`；不把模块写入默认 Profile。

### F-005 · 新的 `TELEGRAM_*` 错误码会撞上冻结 catalog，合同未列入交付物

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：D-009 引入 `TELEGRAM_OPERATOR_UNAVAILABLE`、`TELEGRAM_REQUEST_IN_PROGRESS`、`TELEGRAM_REQUEST_CONFLICT`、`TELEGRAM_RETRY_NOT_ALLOWED`。`error_contract_test.go` 禁止未登记字面量，catalog 也必须中英条目（L169–225）。未知 chat / 未知 request / 非法 `requestId` 的代码未钉（`NOT_FOUND` vs 409）。不阻断按上表开工，但 `go test` 会红。
- 证据：D-009 L33、L81–88；`error_contract_test.go` L19–82、L169–225；`errorcatalog/errorcatalog.go` L22–36。
- 建议：把上述四码（及未知 chat/request 的选择）写入 C3 合同与 catalog/frozen 集合。

### F-006 · `request_id` 字符集未钉死，重试路径会撞 ServeMux

- 严重度：low
- 建议：recommended
- 状态：open
- 描述：重试路由是 `.../messages/{request_id}/retry`。D-009 L73 只说「1～128 字节的安全可打印标识」。若包含 `/` 或空白，Go 1.22 mux 无法把该段当作单一 `{request_id}`。正文上限 4096 UTF-8 **字节**比 Bot API 的 4096 **字符**更严，属 D-002 允许的实施参数，可保留，但应写明是有意收紧。
- 证据：D-009 L48–49、L73–74；D-002 L29。
- 建议：把 `request_id` 收成 mux 安全字母数字（例如 `[A-Za-z0-9._-]{1,128}`），并在合同注明 4096 字节上限。

### F-007 · 发送成功后状态更新失败会留下不可重试的 pending；HTTPSender 在无 token 时返回 nil

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-022
- 描述：D-009 L76–78 禁止把持久化失败伪装成 `sent`，也禁止同一 request 再外发。记录将停在 pending；重试路由只接受 failed；pending-root 唯一又挡住新的 retry。成绩单能看见 pending，但没有合法恢复动作。另：`HTTPSender` 在 token 为空且无 mock 时 `return nil`（`http_sender.go` L82–87），若门禁与发送之间发生卸载，C3 会把「没打出去」写成成功——F-002 若改成允许 idle，这个缝更大。
- 证据：D-009 L76–88；`http_sender.go` L69–87。
- 建议：合同写明 persist-fail-after-send 的 API 状态（5xx、行保持 pending、不得 retry/不得重发）为有意 fail-closed；C3 在调用 sender 前再次确认 token/bot_id，把 HTTPSender 的无 token `nil` 视为失败而不是 `sent`。

## 必改项汇总

| ID | 级别 | 阻断 |
|----|------|------|
| **F-001** | **required / high** | **是：进入 C3 生产代码实施**。冻结 composition `Middleware` 包装与 401→403→409 顺序，禁止把 `Public: false` 当成认证。 |
| **F-002** | **required / high** | **是：进入 C3 生产代码实施**。解开 `running` 门禁与 `settings.read` lease 对专用 operator 权限的回绑，或把 C3/C4 的 lease 权限改写写死。 |
| **F-003** | **required / high** | **是：进入 C3 生产代码实施**。冻结 PG 安全的 `ON CONFLICT DO NOTHING` + `RowsAffected` 幂等/pending-root 冲突，并要求 gated PG 运行时证据。 |

开放 required = **3**。F-004～F-007 为 recommended/open，不单独阻断，但不得假装已冻全。

本条**不**把 F-001 / F-002 / F-003 标为 `accepted-residual` 或 `user-overruled`。闭合路径只有 `fixed`（书面补全 D-009 或等价 C3 合同，再实施）。

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-017 self `pass` / `open_required: 0` | **不作为证据**。本条同意 D-002/D-008 主方向被 D-009 保留、四条路由与 v69 字段形态可识别、C3 未越界做 C4 UI、尚未实施代码的声明。**不同意**「D-009 足以作为 C3 代码实施合同」。A-017 原文不改。 |
| A-003 F-002 | 重试身份已由 D-008+D-009 在决策层钉死（新 request + `retry_of`→root）。本条不闭合 A-003 F-002；实现证据仍待 C3 代码。 |
| A-003 F-004 | D-009 开始贡献专用权限键，设置页仍骑 `settings.read`。认证包装与 profile 同步仍开放（本条 F-001/F-004）。A-003 原文不改。 |
| A-008 F-001 | 同类 PG 同事务唯一冲突问题。C2 已 `fixed`；C3 outbound 未继承写法（本条 F-003）。A-008 原文不改。 |
| D-002 / D-008 / D-009 | 原文保留。F-001/F-002/F-003 要求 `/govern` 补全，而不是本条改写方案。 |

无 self/independent 对同一必改项的一要一否冲突需要当场 P-004 裁 residual。A-017 认为可开工、本条认为不可：这是合同是否足够的意见差，由编排器响应 F-001/F-002/F-003，而不是把 A-017 改写成 fail。

## 结论 + 建议给编排器/用户的下一步

D-009 作为 C3 **方向**合同是可读的：专用权限键、四条 API、字符串化 64-bit ID、分页、pending 先于同步 sender、D-008 的新 request/`retry_of`、无自动重试、C4 UI 与 VP 非目标均未混入。它还不是可实施的安全/方言合同。

**建议 `/govern`：** 按用户选择 `fixed` 补全 D-009（或等价 C3 合同）以闭合 F-001/F-002/F-003；不要进入 C3 生产代码；不要用 A-017 降级本条。补全后再做一次 Grok independent 合同复审。F-004～F-007 可与补全一并写上，以免实施时再撞 catalog/mux/Descriptor。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
