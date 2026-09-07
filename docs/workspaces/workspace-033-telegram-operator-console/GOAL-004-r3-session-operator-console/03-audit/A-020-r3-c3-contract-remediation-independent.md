---
doc_type: goal-audit
id: A-020-r3-c3-contract-remediation-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-04
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: finding-closure
scope: workspace-033 R3 C3 合同修复后 independent re-audit（A-018 required F-001/F-002/F-003 合同侧闭合；D-002/D-008/D-009 v0.2.0/D-010；A-018 原件与 A-019 响应；当前 provider/runtime/auth/store 接缝；A-018 recommended F-004～F-007 是否纳入且无新矛盾；C3 是否可进入生产代码实施）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-020 · R3 C3 A-018 F-001/F-002/F-003 合同修复独立复审（2026-09-04）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C3 合同修复后复审（HEAD `3909d86b1045fc3dd1522251bd931157cd651e7d`；A-018 原件 F-001/F-002/F-003；用户 D-010「专用权限接管 lease」；D-009 v0.2.0；A-019 self 响应；当前 composition/auth/lease/runtime/store/PG 接缝；**三项 required 是否合同侧合法 fixed，C3 是否可进入生产代码实施**）
- **verdict**：pass
- **open_required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码。未读取或比较其他工作区正文。**A-001～A-019 原文及其 findings 全部保留、未改写。** 不把 A-017/A-019 self 或 D-010 选择本身当作合同已满足的成功依据；D-010 只证明用户选择了 A-018 F-002 方案 (a)，闭合证据来自独立核对 D-009 v0.2.0、D-010 正文与当前源码接缝。不接受 residual，不 overrule。不把尚未存在的 C3 代码、v69 或 operator 测试写成已完成。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD**：`3909d86b1045fc3dd1522251bd931157cd651e7d`（`docs(govern): close C3 contract findings`）。工作树干净。该提交仅 8 个 docs 文件；相对 A-018 原审 HEAD `6f935eba` 无 `.go`、无 `telegram_outbound`、无 `telegram.operator.*`、无 `/api/channel/telegram/operator/`
- **covered**：
  1. A-018 required F-001 / F-002 / F-003 是否经 D-009 v0.2.0（及 F-002 的 D-010 用户裁决）书面补全、按 P-003 `fixed` 合法闭合
  2. `Public: false` 不被当成认证；composition 必须 `a.Middleware` 包 operator handler；401→403→409；服务凭据缺 scope 语义
  3. 用户已选「专用权限接管 lease」是否明确为 `settings.read OR telegram.operator.read`；未绑定 polling 仍按心跳 `running` 门禁；settings API 不变
  4. PostgreSQL `request_id` / root-pending 冲突是否明确为方言无关 `INSERT ... ON CONFLICT DO NOTHING` + `RowsAffected` 后在未中止事务内读取；partial unique index 的 SQLite/PG 语义与 gated PG 证据
  5. A-018 recommended F-004～F-007 是否被合同纳入、是否制造与 D-002/D-008/D-010/VP-033 的新矛盾
  6. 当前源码接缝是否仍与合同一致（实施债，不是已修复事实）
- **excluded**：改写 A-001～A-019；采信 A-017/A-019 为交叉成功证据；替用户 residual / overrule；把未实施的 C3 代码写成完成；闭合 A-003 仍开放的 recommended 项；I-033-009/010 的 UI/`getChatMember` 实现（属 C4）

## 成果（有证据）

| 主张 | 本条独立证据（不引用 A-019 结论） |
|------|----------------------------------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| Charter `active` 0.4.0；VP-033 `active`；`vision_ref` 对齐 | `docs/vision/charter.md` L5–6；VP-033 L5–7 |
| HEAD 为合同修复提交；无 operator/v69 业务 diff | `git rev-parse HEAD` = `3909d86b…`；`git show --stat` 仅 8 个 docs；全仓 `.go` 无 `telegram_outbound` / `telegram.operator` / `operator/sessions` |
| A-018 原文保留；F-001/F-002/F-003 在原件仍为 required/open | A-018 L10–11、L105–137、L176–186；本条不把它改写成已闭合 |
| D-002 专用权限与发送状态机原文未改 | D-002 L20、L23、L27–30 |
| D-008 新 request + `retry_of`、无自动重试原文未改 | D-008 L6、L15–21 |
| D-010 以 `source: user` 选择「专用权限接管 lease」 | D-010 L6、L15–19、L33–34 |
| D-009 为 v0.2.0，含 Middleware 包装、lease OR、`ON CONFLICT` | D-009 L8、L23–40、L87–92、L111–139 |
| 现码 `Public` 仍不是安全边界；settings/lease 401 来自 composition `a.Middleware` | `kernel/contribution.go` L27–34、L157–164；`composition.go` L629–631、L673–678；`auth.go` L590–607 |
| 未绑定 polling 无 lease 时仍是 `idle`/`none`；lease 现码仍只认 `settings.read` | `connection_manager.go` L285–295、L303–317、L445–452；`lease_handler.go` L39–41 |
| settings API 现码仍只受 `settings.read/write` | `settings_handler.go` L30–48；`provider.go` L48–66、L153–163 |
| C2 入站已用方言无关 `ON CONFLICT DO NOTHING` + `RowsAffected`；PG 同事务探针存在 | `modules/channel/telegram/store/repository.go` L61–89；`postgres_test.go` L1252–1270；`kernel/unique_violation.go` L22–36 |
| recycle-bin 双方言 partial unique index 先例仍在 | `recyclebin/migration/migration.go` L30、L48 |
| HTTPSender 无 token 且无 mock 时仍 `return nil` | `http_sender.go` L69–87 |
| A-019 声明合同响应不是 independent closure，且未把代码写成已完成 | A-019 L21–23、L41–45 |

## 对照成功标准

### 1) A-018 F-001 · `Public: false` 不是认证；composition 必须包 Middleware

A-018 F-001（required / high / open）要求在写生产代码前把下列句子补进 D-009：**不要** residual / overrule。原件仍保留；闭合写在本条响应侧。

| 闭合要件 | 状态 | 独立证据 |
|----------|------|----------|
| 未走 residual / overruled | **满足** | A-019 L44–45；D-009 书面补全。本条也不代用户接受 |
| `Public: false` 不得被当成认证实现 | **满足** | D-009 L23–25：「`Public: false` 只是路由声明，不是认证实现」 |
| composition 对 operator handler 使用与 settings/lease 相同的 `a.Middleware` 包装后再注入 Provider | **满足** | D-009 L23–25、L46–47。现码先例：`composition.go` L629–631 先包再 `telegrammodule.New(...)`；`mux.Handle` 直接挂 `route.Handler`（L673–678），`validateRoute` 不读 `Public`（`contribution.go` L157–164） |
| 顺序固定为 Middleware 401 → 专用权限 403 catalog `FORBIDDEN` → 运行时 409 | **满足** | D-009 L23–27、L32–36：先认证和专用权限，再运行时；匿名 `401 UNAUTHENTICATED`，缺权限 `403 FORBIDDEN`，不满足运行时 `409 TELEGRAM_OPERATOR_UNAVAILABLE` |
| 匿名、有会话无权限、服务凭据缺 scope 三类不得返回 `TELEGRAM_OPERATOR_UNAVAILABLE` 或会话数据 | **满足** | D-009 L25–27「不得通过错误差异泄漏运行时或会话信息」；L136–137 点名「匿名/无权限/服务凭据缺 scope 的 401/403 顺序」 |
| 服务凭据 scope 语义与现码相容 | **满足** | `auth.go` L605–607、L671–686：合法服务凭据先过 Middleware（失败仍 401），`Permissions = credential.Scopes`；缺 `telegram.operator.*` 是 handler 403，不是 401，也不是 409 |

**结论**：A-018 F-001 作为 **C3 实施合同缺口** 已按 `fixed` 合法闭合。这不是 operator 路由已实现的声明。原件 F-001 仍 `状态：open`；闭合状态写在本条。

### 2) A-018 F-002 · 专用权限接管 lease；保留心跳 `running`；settings API 不变

A-018 F-002 要求二选一写死。用户经 D-010 选择方案 **(a)**「专用权限接管 lease」，不是 (b) 放宽到 `idle`。

| 闭合要件 | 状态 | 独立证据 |
|----------|------|----------|
| 用户书面选择 (a)，非 residual / overruled | **满足** | D-010 L6 `source: user`；L33–34「专用权限接管 lease (Recommended)」 |
| lease 授权明确为 `settings.read` **或** `telegram.operator.read` | **满足** | D-010 L16–19；D-009 L37–40、L136。现码 `RouteContribution` 无 Permission 字段（`contribution.go` L27–34），OR 可在 `LeaseHandler` 于 `IdentityFrom` 之后检查，与 settings/lease 现状同形 |
| 未绑定 polling 仍按心跳把 receiver 推到 `running`；operator API 不隐式自启 | **满足** | D-010 L15–16、L23–24；D-009 L32–36、L40。VP-033 L41、L96 未绑定按心跳启停仍保留 |
| operator surface 继续要求 `running` + 有效 receiver + 已确认 `bot_id` + 未绑定 | **满足** | D-009 L32–36；D-010 L15–16。现码无 lease 时 polling 为 `idle`/`none` 且保留 BotID（`connection_manager.go` L285–295）——合同选择保持该门禁，改的是谁能持有 lease |
| settings API 权限与行为不变 | **满足** | D-010 L19；D-009 L39–40。现码 `settings_handler.go` L38–48 仍 `settings.read`/`settings.write`；菜单仍骑 `settings.read`（`provider.go` L153–163） |
| 写/发送仍要 `telegram.operator.write` | **满足** | D-010 L26–27；D-009 L28–31、L55–56 |
| 占用位 fail-closed 保留 | **满足** | D-009 L32–34；`dispatcher.go` L27–36 `HasBusinessHandlers()` 仍是包级探测 |

与 D-002「专用权限独立于 settings」相容：transcript/发送仍只用 `telegram.operator.*`；lease 是既有心跳引用计数，不是把人工台改回 `settings.*`。与 A-003 F-004「设置页继续 `settings.read` 红线」相容：设置页/设置 API 未改；lease OR 是用户对新运营面可用性的书面选择。

**结论**：A-018 F-002 合同侧 **fixed**。现码 lease 仍只认 `settings.read`（`lease_handler.go` L39–41）是 **C3 实施债**，不是合同缺口。

### 3) A-018 F-003 · 方言无关 `ON CONFLICT` + 未中止事务读取；partial unique；gated PG

| 闭合要件 | 状态 | 独立证据 |
|----------|------|----------|
| pending 插入使用方言无关 `INSERT ... ON CONFLICT DO NOTHING`（占位符 `?`） | **满足** | D-009 L87–89。与 C2 入站同形（`store/repository.go` L61–65） |
| `(bot_id,request_id)` 冲突：`RowsAffected()==0` 后同事务 SELECT 比较 chat/text/`retry_of` | **满足** | D-009 L88–90、L93–96 |
| 不得在 PostgreSQL 唯一失败后于同一已中止 Tx 继续无目标语句 | **满足** | D-009 L90–91。合读：禁止的是普通 `INSERT` 中止后再 SELECT 的反模式 |
| pending-root 冲突同样 `DO NOTHING` 后映射 `TELEGRAM_REQUEST_IN_PROGRESS` | **满足** | D-009 L89–90、L99–102 |
| 禁止把 `kernel.IsUniqueViolation` 当作「已成功接受」的同 Tx 控制流 | **满足** | D-009 L91–92。现码该谓词仍只适合整个 `Run()` 失败（`unique_violation.go` L22–36） |
| DDL：SQLite `INTEGER` / PG `BIGINT`；双方言 `CREATE UNIQUE INDEX ... WHERE status = 'pending'` | **满足** | D-009 L111–115。recycle-bin 先例双方言同一 `WHERE` 形状（`recyclebin/migration/migration.go` L30、L48） |
| 验证覆盖 SQLite **和** gated PostgreSQL；skip ≠ 通过 | **满足** | D-009 L132–139：双方言迁移/读写、同 request 并发与 payload 冲突、`retry_of`/root/并发 pending；「PG 未配置时只能标注 gated skip，不得把 skip 当作通过」 |
| PG 同事务 `ON CONFLICT` 探针仍存在，可供 C3 复用 | **满足** | `postgres_test.go` L1252–1270；钱包 `RowsAffected` 读法（`wallet/store/repository.go` L361–374） |

无目标 `ON CONFLICT DO NOTHING` 在 SQLite 与 PostgreSQL 上都覆盖普通唯一键与 partial unique index。若实现误写成 `ON CONFLICT (bot_id, retry_root) DO NOTHING` 且不带 `WHERE status = 'pending'`，PostgreSQL 无法匹配 partial unique——D-009 字面是无列清单形式，见本条 recommended F-002。不把 F-003 打回 required。

**结论**：A-018 F-003 合同侧 **fixed**。无 v69 表、无 outbound repository；这不是运行时已修复的声明。

### 4) A-018 recommended F-004～F-007 是否纳入、有无新矛盾

| A-018 | D-009 v0.2.0 | 判定 |
|-------|--------------|------|
| F-004 Descriptor + `profile.go` + `reg.Authorization`；`PolicyAdminEditorViewer`/`PolicyAdmin`；不进默认 Profile | L117–119 同步声明四条路由与两权限键；L28–31 默认 `system.admin-editor-viewer` / `system.admin`（即 `PolicyAdminEditorViewer` / `PolicyAdmin`）；不进 mvp/admin/demo；L137 验证点名 Descriptor/profile/catalog | **纳入**。与 VP-033 L62、现码 `provider_test.go` L121–131 红线一致 |
| F-005 四码入冻结 catalog + 未知 chat/request 稳定选择 | L119–122 点名四码及未知 chat/request「稳定选择必须登记」并要求中英条目 | **纳入四码**。未知 chat/request **尚未钉死** `NOT_FOUND` vs 409，见本条 F-001 recommended；不构成与 D-002/D-008 的矛盾 |
| F-006 mux-safe `request_id`；4096 **字节**上限写明有意收紧 | L122–124 `[A-Za-z0-9._-]{1,128}`；4096 字节是有意实施参数 | **纳入**。发送节 L80 仍留「安全可打印」旧句，后文收紧为准，见本条 F-003 recommended |
| F-007 persist-fail-after-send fail-closed；sender 前再确认 token/bot_id；无 token `nil` 不得当 `sent` | L126–130 | **纳入**。与保留 `running` 门禁无矛盾：卸载窗口仍要二次确认。现码 `http_sender.go` L82–87 的 `nil` 是实施时必须按合同当失败的接缝 |

未把未选方案（复用 `settings.*` 发送、原地 `failed→pending`、后台自动重试、operator 隐式自启 polling、把 `idle` 当可用）写回 D-009。C3 仍不建 page/nav/fragment，不扩张 kernel Telegram port。

### 5) 信息门禁（P-005）

| 项 | 最晚阶段 | 状态 | 对本条 |
|----|----------|------|--------|
| I-033-021 | C1/C3 | verified (decision + contract)；实现+测试未宣称完成 | 认证包装、lease OR、Descriptor/profile 已写入合同。不得把「decision verified」当成代码完成 |
| I-033-022 | C1/C3 | verified (decision + contract)；D-008 重试身份仍在 | PG 安全冲突读法已写入合同。实现证据仍待 C3 代码 |
| I-033-009 | C3/C4 | verified (decision)；UI 属 C4 | D-009 未越界关闭 |
| I-033-010 | C1/C4 | verified (decision)；实现属 C4 | D-009 未实现预检 |

无到期且未接受 residual 的 required 信息项阻断本条放行。本条放行的是 **C3 生产代码实施**，不是 C3 检查点关闭。

## Findings

本条无新的 required finding。下列 recommended 不阻断 C3 开工。

### F-001 · 未知 chat / 未知 request 的 catalog 码仍未钉死

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：A-018 F-005
- 描述：D-009 L119–122 要求未知 chat/request「稳定选择必须登记」，但未在 `NOT_FOUND` 与 409/`TELEGRAM_*` 之间点名。四码本身已写入合同。实现时必须选一个稳定码并加入 catalog/frozen 集合，否则 `error_contract_test.go` 仍会红。不阻断按已冻合同开工。
- 证据：D-009 L74–75、L119–122；`error_contract_test.go` L19–82；`errorcatalog/errorcatalog.go` L22–36（尚无 `TELEGRAM_*`，属实施债）。

### F-002 · pending-root 部分唯一索引的 `ON CONFLICT` 目标不得丢掉 `WHERE`

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：A-018 F-003
- 描述：D-009 字面是无列清单 `INSERT ... ON CONFLICT DO NOTHING`，双方言均可覆盖 partial unique。若实现仿 C2 入站写成 `ON CONFLICT (bot_id, retry_root) DO NOTHING` 且不带 `WHERE status = 'pending'`，PostgreSQL 无法推断该 partial unique index。C3 必须保持无目标形式，或写带 `WHERE status = 'pending'` 的目标，并在 gated PG 上跑 root-pending 冲突路径。这是实施陷阱，不是合同缺口。
- 证据：D-009 L87–92、L111–115；`store/repository.go` L65（入站是**完整**唯一键，不是 partial）；`recyclebin/migration/migration.go` L30、L48。

### F-003 · 发送节「安全可打印」与验证分母 mux-safe 收紧并存

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：A-018 F-006
- 描述：D-009 L80 仍写「1～128 字节的安全可打印标识」；L122–123 收紧为 `[A-Za-z0-9._-]{1,128}`。后者为准。含 `/` 或空白的 id 仍会撞 Go 1.22 mux。实施与测试按 mux-safe 字符集，不要按「可打印」放行。
- 证据：D-009 L80–81、L122–123。

## 必改项汇总

| ID | 级别 | 阻断 |
|----|------|------|
| A-018 **F-001** | required / high | **否（响应侧 `fixed`）**：composition `a.Middleware`、禁止把 `Public: false` 当认证、401→403 `FORBIDDEN`→409、服务凭据缺 scope 已写入 D-009。原件保留。 |
| A-018 **F-002** | required / high | **否（响应侧 `fixed`）**：D-010 选择专用权限接管 lease；D-009 写死 `settings.read OR telegram.operator.read`、保留心跳 `running`、settings API 不变。原件保留。 |
| A-018 **F-003** | required / high | **否（响应侧 `fixed`）**：方言无关 `ON CONFLICT DO NOTHING` + `RowsAffected` 未中止事务读取、双方言 partial unique、gated PG skip≠通过已写入 D-009。原件保留。 |

开放 required = **0**。本条 recommended F-001～F-003 不单独阻断。

本条**不**把 A-018 F-001/F-002/F-003 标为 `accepted-residual` 或 `user-overruled`。闭合路径是合同侧 `fixed`。实现与测试事实仍待 C3 代码；本条不关闭 C3 检查点。

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-018 independent `conditional` / open_required=3 | **原件不改写**。本条确认其三项 required 在 D-009 v0.2.0 + D-010 后合同侧 `fixed`。A-018 原文 F-001～F-007 仍为 open |
| A-019 self `pass` / open_required=0 | **不作为证据**。本条同意其「三项已写入合同、未宣称代码完成、等待本 re-audit」的自我限制。独立核对后同意可开工 |
| A-017 self `pass` | **仍不作为证据**。A-017 原审的是补全前合同；本条不追溯改写它 |
| D-002 / D-008 | 原文保留。专用权限键与新 request/`retry_of` 仍被 D-009 忠实保留 |
| D-010 | 用户书面选择 (a)；本条核对其与 D-009/VP-033 心跳门禁一致，而不是用它替代合同句子 |
| A-003 F-004 | 设置页仍骑 `settings.read`。lease OR 是后续用户裁决，不废止设置页红线。A-003 原文不改 |

无 self/independent 对同一必改项的一要一否冲突需要当场 P-004。

## 结论 + 建议给编排器/用户的下一步

A-018 三项 required 已在合同侧合法 `fixed`。D-009 v0.2.0 现在是可实施的安全/方言合同：Middleware 包装与 401→403→409、lease `settings.read OR telegram.operator.read`、未绑定仍按心跳 `running`、settings API 不变、方言无关 `ON CONFLICT DO NOTHING` + 未中止事务读取、双方言 partial unique 与 gated PG 证据。F-004～F-007 已纳入，未制造与 D-002/D-008/D-010/VP-033 的新矛盾。HEAD 仍无 C3 生产代码。

**建议 `/govern`：** 响应本条 `pass`；进入 C3 生产代码实施。不要用本条关闭 C3 检查点。实现须遵守 D-009/D-010 全文，并处理本条 recommended：未知 chat/request 选稳定 catalog 码；pending-root `ON CONFLICT` 保持无目标或带 `WHERE status = 'pending'`；`request_id` 按 mux-safe 字符集，不要按「可打印」放行。C3 代码完成后再做 self + Grok independent 实现审计。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
