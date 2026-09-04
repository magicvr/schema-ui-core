---
doc_type: goal-audit
id: A-032-r3-c4-a030-remediation-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: finding-closure
scope: workspace-033 R3 C4 A-030 finding 修复后复审（当前工作树源码/测试/双语 catalog/composition·settings 占用接线；A-030 F-001/F-002/F-003；不采信 A-031/E-018；不选择 I-033-023；不关闭 C4）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-032 · R3 C4 A-030 修复后独立复审（2026-09-05）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C4 **A-030 修复响应复审**（当前 HEAD `8abe5aca5b174a4ff6aeb420dc2ce24844038a78`；C4 与 A-030 响应仍在**未提交工作树**；对照 A-030 原文 F-001/F-002/F-003、D-002 I-033-009、D-009 占用位 409、当前源码/测试/双语 catalog/composition settings 接线与本会话跑数）。**不选择 I-033-023**。**不关闭 C4**。
- **verdict**：pass
- **open_required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。**A-001～A-031 原文及其 findings 全部保留、未改写。** 不把 A-029/A-031 self、E-017/E-018 或 Root E-022/E-023 当作成功依据。不把 recommended 升级为 required，不接受 residual，不 overrule。不自行关闭 C4。不在三种 capability API 形状中替用户选择，不实现 `getChatMember` / 发送 / retry。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD / 工作树**：本会话 `git rev-parse HEAD` = `8abe5aca5b174a4ff6aeb420dc2ce24844038a78`（`docs(govern): close R3 C3 checkpoint`）。C4 与 A-030 响应仍在工作树：`telegram-admin-tab.tsx` / `.test.tsx`、`en-US.json` / `zh-CN.json`、`runtime.go` / `runtime_test.go` / `settings_handler.go`、`composition.go`，以及未跟踪的 E-017/E-018、A-029/A-030/A-031。本条以**当前工作树源码与本会话跑数**为准，不信任描述或 A-031 结论。
- **covered**：
  1. A-030 F-001：`schema.telegram.operator.send` 是否同时存在于 en-US / zh-CN，组件是否仍渲染该键
  2. A-030 F-002：可见态 10 秒边界（`9999ms` 不发、再 `1ms` 第二次 sessions）与 `operatorRefreshRef` 单飞合并是否有失败即红测试
  3. A-030 F-003：Dispatcher `HasBusinessHandlers` 是否经 settings GET/PATCH 以只读 `business_occupied` 暴露；占用或未知时 Admin 入口是否隐藏；polling lease 是否 fail-closed 不 acquire
  4. 双语 catalog 键集合、composition 与 settings 接线、composer/retry 仍 fail-closed、未接通发送/`getChatMember`
  5. 本会话独立重跑 Web 定向/全量、API `telegram`/`composition`/`docscheck`；写集外 `tsc` 基线
- **excluded**：改写 A-001～A-031；采信 A-031/E-018；选择独立 capability / 成绩单附带 / 会话列表附带；实现或修改 `getChatMember`、发送/retry、status/progress/goal-tree；全仓 `go test ./...`；C4 关门

## 成果（有证据）

| 主张 | 本条独立证据（不引用 A-031 结论） |
|------|----------------------------------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| Charter `active` 0.4.0；VP-033 `active`；`vision_ref` 对齐 | `docs/vision/charter.md` L5–6；VP-033 L5–7 |
| HEAD 仍是 C3 关闭文档提交；C4/A-030 响应在工作树 | `git rev-parse HEAD`；`git status` 列出上述 web/api 文件为已修改；A-029～A-031 未跟踪 |
| A-030 原文仍为 independent `conditional` / `open_required: 1`；F-001 required/open，F-002/F-003 recommended/open | `A-030-r3-c4-ui-foundation-independent.md` L10–11、L86–119；本条不改写 |
| A-031 原文仍为 self `pass`，声称三项 `fixed` | `A-031-r3-c4-a030-response.md` L6–11、L24–52；**不是本条证据** |
| I-033-009 已用户裁决：10 秒单飞、失焦暂停、恢复立即刷新 | D-002；`00-meta.md` L52 |
| I-033-023 仍 `collecting`，本条不选择三种 API 形状 | `00-meta.md` L58 |
| `operatorCapability` 仍只被重置为 `"unknown"`，从不设 `"allowed"`；send 无 `onClick` | `telegram-admin-tab.tsx` L104、L376–378、L753 |
| 本会话定向/全量 Web 与 API 定向 **PASS**；`tsc -b` 被写集外错误阻断 | 见下方「本会话验证」 |

### 本会话验证（独立执行，2026-09-05）

对照当前工作树（非 HEAD blob）：

| 命令 | 结果 | 归类 |
|------|------|------|
| `apps/web` `npm test -- --run src/components/telegram-admin-tab.test.tsx src/i18n/schema-keys.structural.test.ts --reporter=verbose` | **PASS**（Test Files 2/2；Tests **14/14**：tab **10/10**、catalog identical keys 4/4） | 通过 |
| `apps/web` `npm test -- --run` | **PASS**：Test Files **92 passed (92)**；Tests **1205 passed (1205)**；Duration 10.02s | 通过 |
| `apps/api` `go test ./internal/channel/telegram ./internal/composition ./internal/docscheck -count=1` | **PASS**：telegram 5.886s；composition 21.899s；docscheck 0.205s | 通过 |
| `apps/web` `npx tsc -b --pretty false` | **FAIL**：`src/renderer/form-controls.tsx(946,11)` 与 `(947,11)` `TS2322`：`number \| undefined` 不能赋给 `string \| undefined` | **写集外**既有类型错误；`git status` 未列出该文件；不构成本切片 fail |
| `git diff --check`（C4 web/api 写集 8 文件） | 无 whitespace 报错（仅 CRLF 提示） | 通过 |

未把 skip 记为通过。未跑 e2e、未跑全仓 `go test ./...`、未跑 `npm run build` 的 vite 后半段（`tsc -b` 已失败）。`schema-keys.structural.test.ts` 只断言两套 catalog **键集合相同**，不能单独证明组件 `t()` 用到的每个 key 都存在。

## 对照成功标准

本条审的是 **A-030 F-001/F-002/F-003 在当前工作树是否已有可核对修正**。C4 检查点关闭权在 `/govern`，不在本条。GOAL 级发言权/`getChatMember`/发送状态机仍属 C4 后续，不得写成已交付。

### 1) A-030 F-001 · 双语 `schema.telegram.operator.send`

A-030 F-001（required / med / 原件 **open**）：已渲染发送按钮缺少该键。

| 核对 | 当前工作树 | 本条判定 |
|------|------------|----------|
| 组件仍渲染该键 | `telegram-admin-tab.tsx` L753 `{t("schema.telegram.operator.send")}`；按钮仍 `disabled`、无 `onClick` | 键仍被使用 |
| en-US 有词条 | `en-US.json` L1036 `"schema.telegram.operator.send": "Send"` | **满足** |
| zh-CN 有词条 | `zh-CN.json` L1036 `"schema.telegram.operator.send": "发送"` | **满足** |
| 缺键时不再回退为生键 | `catalog.ts` L7–9、L87–99：locale → en-US → key itself。两边都有键，回退链不会落到 key 名 | **生产缺口闭合** |
| 失败即红断言 | tab 测试 L148 `toContain("Send")`。fieldset 图例是 `"Send as bot"`（`schema.telegram.operator.composer` L1031），缺 send 键时 `t()` 回退为 `schema.telegram.operator.send`，**子串 `"Send"` 仍会命中**。structural 测试只比键集合 | **生产已修**；断言不够失败即红，见本条 recommended F-001 |

**响应侧判定：fixed。** 原件 required 缺口（catalog 缺键导致可见生键）已不存在。不把弱断言升级为新的 required。

### 2) A-030 F-002 · 10 秒边界与 `operatorRefreshRef` 单飞

A-030 F-002（recommended / low / 原件 **open**）：测试只锁「10 秒前不发」，未钉 10000ms 第二次 sessions，也未钉 refresh 并发合并。

| 核对 | 当前工作树 | 本条判定 |
|------|------------|----------|
| 源码周期 | L67 `telegramLeaseIntervalMs = 10_000`；L360–366 完成后 `setTimeout(..., telegramLeaseIntervalMs)` | 仍是完成后再 schedule |
| 10 秒边界钉 | 测试 L168–177：推进 `9_999ms` 仍 1 次 sessions，再推进 `1ms` + `flushPromises` 后为 2 次。本会话该测 **PASS** | **满足** A-030 要求的可见态第二次请求 |
| `operatorRefreshRef` 单飞 | 源码 L335–352：已有 promise 则 await 后 return。测试 L197–230：pending sessions 下 hidden→visible，`sessionCalls` 保持 1；释放后仍为 1。本会话该测 **PASS** | **满足** |
| 失焦暂停/恢复即刷 | 测试 L179–189：hidden 后 20s 仍 2 次，unhide 后 3 次。本会话 **PASS** | 仍成立 |
| `timelineFlightsRef` 同 chat 重叠 | 源码 L266–271 仍按 chat 复用 in-flight；**无**测试让两次 timeline 重叠 | 原 F-002 第三点仍未钉，见本条 recommended F-002 |

**响应侧判定：fixed**（用户本轮指定复核的 10 秒边界与 `operatorRefreshRef` 单飞）。不把 `timelineFlightsRef` 未钉升为 required。

### 3) A-030 F-003 · Dispatcher 占用信号、入口隐藏、lease fail-closed

A-030 F-003（recommended / low / 原件 **open**）：settings 无 occupancy 字段；入口在 configured+`bot_id>0` 即渲染；已绑定表现为 sessions 409/`loadFailed`。

| 核对 | 当前工作树 | 本条判定 |
|------|------------|----------|
| Dispatcher 探针 | `dispatcher.go` L27–36：`len(commands)>0 \|\| len(callbacks)>0` | 业务占用信号存在 |
| settings 只读字段 | `runtime.go` L53 `BusinessOccupied bool \`json:"business_occupied"\``（无 `omitempty`）；`RuntimeManager.Status()` 仍不填该字段 | 字段已进入 JSON 合同 |
| settings 接线 | `settings_handler.go` L19–24 接受 probe；L108–113 GET/PATCH 经 `status()` 覆盖。probe 为 nil 时保持 Go 零值 `false` | composition 路径有 probe |
| composition 与 operator 同一探针 | `composition.go` L629–633 settings 与 L644–646 operator 均为 `tr.DispatcherState != nil && tr.DispatcherState.HasBusinessHandlers()` | **满足**「同一进程级探针」 |
| operator 仍 409 | `telegram_operator.go` L172 `return !h.businessProbe()` | 占用时 API 仍不可用；UI 不再依赖 409 文案隐藏 |
| Admin 入口隐藏 | 组件 L638：仅 `configured && business_occupied === false && bot_id>0` 渲染 `[data-telegram-operator]`。未知/`true` 不渲染 | **满足** 占用或未知隐藏 |
| polling lease fail-closed | L159：`business_occupied !== false` 则 `leaseState=inactive` 并 return，不 `acquire`。测试 L238–250：polling + `business_occupied: true` 时 operator 为 null 且 calls 不含 `/lease/acquire`。本会话 **PASS** | **满足** 占用时不 acquire |
| API GET/PATCH 占用信号 | `runtime_test.go` L67 probe `return true`；L125–127 GET、L143–145 PATCH 断言 `BusinessOccupied`。本会话 `go test ./internal/channel/telegram` **PASS** | **满足** handler 层 |
| composition mux 是否断言 settings 占用 | `composition_telegram_test.go` 无 `business_occupied` 断言 | 覆盖缺口，见本条 F-003 |
| lease effect deps | L252 `[callLease, loadState, status?.mode]`，**不含** `status.business_occupied`。首次 `loadState→ready` 时能读到占用快照；页内占用翻转不重跑 effect | 首次加载路径满足原 finding；热更新缺口见本条 F-003 |

**响应侧判定：fixed**（原缺陷：无占用信号 + 入口仍显示 + 用 409 当隐藏）。首次加载占用隐藏与不 acquire 有测试钉。

## Findings

本条 **无新增 required**。下列为覆盖紧密度，不重开 A-030 F-001，不阻断本 finding-closure 的 pass，也**不构成**对 I-033-023 的裁决。

### 建议（recommended）

#### F-001 · 发送文案测试钉不能在缺键时单独变红

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：A-030 F-001 的回归锁；非 I-033-023
- **是否阻断 C4 关门**：**否**（catalog 两边已有键；不升 required）
- 描述：`telegram-admin-tab.test.tsx` L148 对 composer `textContent` 断言 `toContain("Send")`。同 fieldset 图例 `schema.telegram.operator.composer` = `"Send as bot"`（`en-US.json` L1031），缺 `schema.telegram.operator.send` 时 `t()` 回退键名也含 `"Send"`（`catalog.ts` L7–9）。该钉在 catalog 回退到生键时仍绿。zh-CN `"发送"` 只被 identical-keys 间接锁住，没有渲染断言。
- 证据：测试 L115–148；`en-US.json` L1031 与 L1036；`catalog.ts` L87–99。本会话 10/10 PASS 只证明已写断言为绿。
- 建议：断言按钮自身文本等于 catalog 值（或 `not.toContain("schema.telegram.operator.send")`），不要用会与 `"Send as bot"` 碰撞的子串。不要顺手接通发送 API。

#### F-002 · `timelineFlightsRef` 同 chat 重叠请求仍无失败即红钉

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：A-030 F-002 原第三点；I-033-009
- **是否阻断 C4 关门**：**否**（源码仍去重；10 秒与 `operatorRefreshRef` 已钉）
- 描述：`loadTimeline` L266–271 在 `timelineFlightsRef` 命中已有 promise 时 await 后 return。当前 10 个 tab 测试没有让同一 `chatId` 的两次 messages 请求重叠。
- 证据：组件 L266–300；测试文件无 timeline 并发合并用例。本会话 10/10 不能覆盖该分支。
- 建议：补一条 pending messages + 二次 `loadTimeline` 的单飞钉。不升 required。

#### F-003 · 占用信号的 composition mux、缺字段 UI 与 lease 热更新未钉

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：A-030 F-003 覆盖紧密度；D-009 占用位
- **是否阻断 C4 关门**：**否**（首次加载占用隐藏与不 acquire 已钉；composition 源码已接同一探针）
- 描述：
  1. `go test ./internal/composition` 本会话 PASS，但没有任何用例断言真实 mux 的 GET/PATCH settings JSON 含 `business_occupied`，或 Dispatcher 注册 command 后该字段翻转为 true。handler 单测用的是字面 `return true` 探针（`runtime_test.go` L67）。
  2. UI 测试只覆盖 `business_occupied: true`（L243），没有「字段缺省 / `undefined`」时入口隐藏且不 acquire 的钉。源码 L159/L260/L638 对未知是 fail-closed，但无会红测试。
  3. lease effect 依赖 L252 省略 `status.business_occupied`。首次 `loadState === "ready"` 能看到占用快照；若同一挂载周期内 settings 快照从 `false` 变为 `true` 且 `mode` 不变，effect 不会清理已持有的 polling lease。占用字段目前只在 mount/`save`/`clearSecrets` 刷新。
- 证据：`composition.go` L629–646；`composition_telegram_test.go` 无 occupancy 断言；组件 L159、L252、L638；测试 L238–250。
- 建议：composition 定向钉 + 缺字段 UI 钉；若要锁热更新，把 `business_occupied` 纳入 lease effect 依赖。不在本条设计新信号，不升 required。

## 必改项汇总

| ID | 级别 | 阻断 C4 关门 |
|----|------|----------------|
| A-030 F-001 | required / med · **响应侧 fixed**（原件仍 open） | 原文「是」；**本条复核后生产缺口已闭合** |
| A-030 F-002 | recommended / low · **响应侧 fixed**（10 秒 + `operatorRefreshRef`） | 否 |
| A-030 F-003 | recommended / low · **响应侧 fixed**（占用信号 + 入口隐藏 + 占用时不 acquire） | 否（本切片）；C4 关门仍须保持占用隐藏 |
| 本条 F-001 | recommended / low | **否** |
| 本条 F-002 | recommended / low | **否** |
| 本条 F-003 | recommended / low | **否** |
| I-033-023 | required 信息项 / collecting | **是**（C4 API/UI/关门门禁；本条不选择形状） |

开放 required finding = **0**。开放 required 信息项 I-033-023 仍为 `collecting`。本条不把任何 finding 标为 `accepted-residual` 或 `user-overruled`。

## 与既有意见的异同

- A-030 independent `conditional` / open_required=1：原件 F-001/F-002/F-003 全文保留。本条独立复核后，三项在**响应侧**记为 `fixed`，不改写 A-030。
- A-031 self `pass` / open_required=0：方向与本条一致（三项已处理；不选 I-033-023；不关 C4）。本条**不采信**其测试结论，但本会话独立重跑后 10/10、14/14、92/1205、telegram/composition/docscheck 数字可重复。
- A-031 未记录发送文案断言与 `"Send as bot"` 碰撞、也未记录 `timelineFlightsRef` / composition mux / 缺字段 UI / lease deps。本条列为 recommended，不重开 A-030 F-001。
- A-029 self 不是本条证据。A-001～A-028 与本切片无关；原文保留。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** A-030 的 required F-001 与 recommended F-002/F-003 在当前工作树有可重复核对的修正：双语 `schema.telegram.operator.send` 已存在；10 秒边界与 `operatorRefreshRef` 单飞有会红测试；settings 经与 operator 相同的 `HasBusinessHandlers` 探针暴露 `business_occupied`，占用或未知时 Admin 入口不渲染且不 acquire polling lease。本条 **open required finding = 0**。

即使如此，**C4 仍不能关门**，原因不是本条替用户做了 I-033-023 选择，而是：

1. **I-033-023 仍为 `collecting`**：C4 API/UI/关门的 required 信息项未闭合。三种互斥形状必须由用户书面裁决；本条明确不选，也不实现 `getChatMember` / 发送 / retry。
2. D-002 已冻结的混合发言权仍未落地：`getChatMember` 预检、60 秒 bot/chat 缓存、403 失效、显式重探、composer 按真实 `can_send` 启用。`operatorCapability` 唯一写入仍是重置为 `"unknown"`。
3. 发送/失败/retry 状态机尚未接到 UI（send 仍硬编码 `disabled`、无 `onClick`）。
4. 本条 recommended F-001～F-003 为覆盖紧密度，不阻断本 finding-closure。
5. 写集外 `form-controls.tsx` L946–947 `tsc` 失败是既有基线，不是本切片回归。

建议 `/govern`：响应本条；可将 A-030 F-001/F-002/F-003 在响应侧记为 `fixed`（原件不改写）；**询问用户 I-033-023**（三种形状 + 建议）；在用户裁决并实现发言权/发送路径之前，保持 C4 `进行中`。不要把本条当作 C4 或 R3 关门证据。recommended 覆盖钉可在后续 C4 工作中顺手补，不升 required。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
