---
doc_type: goal-audit
id: A-030-r3-c4-ui-foundation-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: execution-facts
scope: workspace-033 R3 C4 基础 UI 切片（未决前的 E-017/A-029、telegram-admin-tab.tsx 及其测试与双语文案；会话列表/成绩单接入、10 秒单飞、visibilitychange 失焦暂停与恢复即刷、请求去重、未得到 capability 时 composer/retry fail-closed；核对当前测试与工作树）。不含 I-033-023 三种 capability API 形状的选择，不含 C4 关门。
verdict: conditional
open_required: 1
version: 0.1.0
---

# A-030 · R3 C4 UI 基础切片独立交叉审计（2026-09-05）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：execution-facts · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C4 **基础 UI 切片**（当前 HEAD `8abe5aca5b174a4ff6aeb420dc2ce24844038a78`；C4 写集在**未提交工作树**；对照 D-002 I-033-009 / D-009 C3 JSON 合同；独立核对源码、测试、双语文案与本会话跑数）。**不选择 I-033-023**。**不关闭 C4**。
- **verdict**：conditional
- **open_required**：1
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。**A-001～A-029 原文及其 findings 全部保留、未改写。** 不把 A-029 self、E-017 或 Root E-022 当作成功依据。不接受 residual，不 overrule。不自行关闭 C4。不在三种 capability API 形状中替用户选择。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD / 工作树**：本会话 `git rev-parse HEAD` = `8abe5aca5b174a4ff6aeb420dc2ce24844038a78`（`docs(govern): close R3 C3 checkpoint`）。**HEAD 中的 `telegram-admin-tab.tsx` 不含 operator 会话/成绩单/capability 代码。** C4 写集仅在工作树：`apps/web/src/components/telegram-admin-tab.tsx`、`.test.tsx`、`en-US.json`、`zh-CN.json`，以及未跟踪的 E-017/A-029 与 Root E-022。本条以**当前工作树源码与本会话跑数**为准，不信任描述或 A-029 结论。
- **covered**：
  1. 会话列表与成绩单是否按 C3 `{items,total,page,pageSize}` / 十进制字符串 `chatId` 接入
  2. 10 秒单飞、`visibilitychange` 失焦暂停、恢复立即刷新
  3. operator refresh / 同一 chat 成绩单请求去重
  4. 未得到 capability 结果时 composer / retry fail-closed
  5. 双语文案是否覆盖本切片已渲染控件
  6. 当前测试是否覆盖上述行为，并与工作树一致
  7. I-033-009 / I-033-010 / I-033-023 对 C4 门禁的状态（不裁决 I-033-023）
- **excluded**：改写 A-001～A-029；采信 A-029/E-017；选择独立 capability / 成绩单附带 / 会话列表附带；实现或修改 `getChatMember`、发送/retry、status/progress/goal-tree；全仓 `go test`；C4 关门

## 成果（有证据）

| 主张 | 本条独立证据（不引用 A-029 结论） |
|------|----------------------------------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| Charter `active` 0.4.0；VP-033 `active`；`vision_ref` 对齐 | `docs/vision/charter.md` L5–6；VP-033 L5–7 |
| HEAD 仍是 C3 关闭文档提交；C4 UI 只在工作树 | `git rev-parse HEAD`；`git show HEAD:apps/web/src/components/telegram-admin-tab.tsx` 无 operator/sessions/capability；`git status` 列出上述 4 个 web 文件为已修改 |
| A-029 原文仍为 self `conditional` / open_required=0 | A-029 L6–11；**不是本条证据** |
| I-033-009 已用户裁决：10 秒单飞、失焦暂停、恢复立即刷新 | D-002 L21；`00-meta.md` L52 |
| I-033-023 仍 `collecting`，本条不选择三种 API 形状 | `00-meta.md` L58 |
| 本会话定向/全量 Web 测试 **PASS**；`tsc -b` 被写集外错误阻断 | 见下方「本会话验证」 |

### 本会话验证（独立执行，2026-09-05）

在 `apps/web`，对照当前工作树（非 HEAD blob）：

| 命令 | 结果 | 归类 |
|------|------|------|
| `npm test -- src/components/telegram-admin-tab.test.tsx src/i18n/schema-keys.structural.test.ts --reporter=verbose` | **PASS**（12 tests；其中 tab 8/8、catalog identical keys 4/4） | 通过 |
| `npm test`（全量 vitest） | **PASS**：Test Files **92 passed (92)**；Tests **1203 passed (1203)**；Duration 10.09s | 通过 |
| `npx tsc -b --pretty false` | **FAIL**：`src/renderer/form-controls.tsx(946,11)` 与 `(947,11)` `TS2322`：`number \| undefined` 不能赋给 `string \| undefined` | **写集外**既有类型错误；不构成本切片 fail |
| `git diff --check`（C4 四个 web 文件） | 无 whitespace 报错 | 通过 |

未把 skip 记为通过。未跑 e2e、未跑 `npm run build` 的 vite 后半段（`tsc -b` 已失败）。`schema-keys.structural.test.ts` 只断言两套 catalog **键集合相同**，不能证明组件 `t()` 用到的每个 key 都存在。

## 对照成功标准

本条只审 C4 **基础切片**，不是 C4 检查点关门。GOAL 级成功标准中「列出私聊/群并展示成绩单」的只读部分属于本切片；「已绑定时入口隐藏」「发言权 composer 按 can_send/403 禁用」「端到端验证」仍属后续 C4，不得写成已交付。

| 标准 | 状态 | 证据 |
|------|------|------|
| 会话列表接入 C3 operator sessions，`chatId` 为字符串 | **本切片满足** | `telegram-admin-tab.tsx` L73–74、L301–331、L659–679；C3 `telegram_operator.go` L183–219 `json:"chatId"`；测试 L118–141 使用 `"8001"` |
| 成绩单接入 C3 messages，按选中 chat 加载 | **本切片满足** | L264–298、L682–728；路径 `` `${sessions}/${encodeURIComponent(chatId)}/messages?page=1&pageSize=100` ``；测试渲染 `hello` / `reply` |
| `{items,total,page,pageSize}`；pageSize=100 ≤ D-009 上限 | **本切片满足** | L59–64、L74；D-009 L59–60 默认 20、最大 100 |
| 10 秒单飞 | **源码满足；测试只锁「10 秒前不发」** | L66–67 `telegramLeaseIntervalMs = 10_000`；L333–351 `operatorRefreshRef`；L353–372 完成后 `setTimeout(..., 10000)`。测试 L166–169 推进 9999ms 仍 1 次调用，**未**推进 10000ms 断言第二次调用 |
| 失焦暂停、恢复立即刷新 | **本切片满足** | L113–119 `visibilitychange` → `pageVisible`；L353–354 `!pageVisible` 直接 return 并 cleanup timer；恢复后 effect 重跑并 `void refresh()`。测试 L171–181：hidden 后 20s 仍 1 次，unhide 后 2 次。本会话该测 PASS |
| 请求去重 | **源码满足；测试未钉并发合并** | `operatorRefreshRef` 命中已有 promise 则 await 后 return；`timelineFlightsRef` 按 chat 去重。无测试让两次 refresh/timeline 重叠 |
| 未得到 capability 时 composer/retry fail-closed | **本切片满足** | `operatorCapability` 初值 `"unknown"`，唯一写入是 L376 重置为 `"unknown"`，从不设 `"allowed"`；fieldset `disabled={operatorCapability !== "allowed" \|\| !operatorReady}`；retry `disabled={operatorCapability !== "allowed"}`；send 按钮硬编码 `disabled`；retry/send **无 onClick、无 POST**。测试 L144–147 fieldset 与 retry 均为 disabled |
| 双语文案覆盖已渲染控件 | **部分不满足** | 见 F-001 |
| I-033-023 capability API 形状 | **未决 / 本条不选择** | `00-meta.md` L58 `collecting` |
| C4 检查点可关闭 | **否** | 见结论 |

## Findings

### 必改（required）

#### F-001 · 已渲染的发送按钮缺少 `schema.telegram.operator.send` 双语词条

- 严重度：med
- 建议：required
- 状态：open
- 关联：C4 UI 文案；不关联 I-033-023 形状选择
- **是否阻断 C4 关门**：**是**（文案完备）。**不阻断**本切片后续 capability 实施，也**不构成**对 I-033-023 的裁决。
- 描述：工作树 `telegram-admin-tab.tsx` L751 渲染 `{t("schema.telegram.operator.send")}`。`en-US.json` / `zh-CN.json` 本切片新增了 22 个 `schema.telegram.operator.*` 键，**两边都没有** `schema.telegram.operator.send`。catalog 缺失时 `t()` 回退为键名本身（`catalog.ts` L7–9、L82）。发送按钮虽保持 disabled（fail-closed 仍成立），但人工台可见控件会显示生键。`schema-keys.structural.test.ts` 因两套 catalog 键集合相同而通过，定向测试也不断言 send 文案，故 8/8 与 92/1203 **不能**覆盖本缺口。
- 证据：组件 L751；`git diff` 两套 catalog 新增键止于 `capabilityUnavailable`；本会话 `rg schema.telegram.operator.send` 仅命中组件一处；structural test 本会话 PASS。
- 建议：在 en-US / zh-CN 同时补 `schema.telegram.operator.send`，并加一条失败即红断言。不要顺手接通发送 API。

### 建议（recommended）

#### F-002 · C4 刷新/去重行为缺少失败即红的测试钉

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：I-033-009
- **是否阻断 C4 关门**：**否**（源码已实现；不升 required）
- 描述：tab 文件共 8 个测试，其中 6 个是既有 settings/lease，**仅 2 个**覆盖本切片（会话/成绩单/composer fail-closed；失焦暂停/恢复即刷）。源码中的 10 秒周期（完成后再 schedule）与 `operatorRefreshRef` / `timelineFlightsRef` 去重没有单独会红的测试：可见态推进到 10000ms 应出现第二次 sessions 请求；两次并发 refresh 应合并为一次网络；同一 chat 的重叠 timeline 请求应复用 in-flight promise。
- 证据：`telegram-admin-tab.test.tsx` L114–187 vs L94–317 其余 6 测；组件 L264–372。本会话 8/8 PASS 只证明已写断言为绿。
- 建议：补 10 秒触发与单飞合并钉。不升 required。

#### F-003 · 已绑定占用位时人工台入口仍显示，409 表现为加载失败

- 严重度：low
- 建议：recommended
- 状态：open
- 关联：GOAL 成功标准「已绑定时人工台入口隐藏」；D-009 占用位 409；非 I-033-023
- **是否阻断 C4 关门**：**是，作为 C4 剩余工作，不作为本切片实施缺陷的高优先级 required。** 本条不把它升为 required，以免把未声称交付的占用位隐藏写成切片作假。
- 描述：GOAL 成功标准要求已绑定时隐藏入口。`operatorReady`（L257–262）检查 configured / running / `bot_id>0` / polling lease，**不含** `HasBusinessHandlers()`。operator 区块在 configured 且 `bot_id>0` 时即渲染（L636）。settings `RuntimeStatus`（`runtime.go` L45–58）无 occupancy 字段。已绑定时 C3 API 仍会 409，UI 会走 `sessionsLoadState === "error"` 显示 `operator.loadFailed`，而不是隐藏入口。
- 证据：组件 L257–262、L636–654；`runtime.go` L45–58；D-009 L32–36。
- 建议：C4 后续用服务端占用位信号隐藏入口；不要用 409 文案代替「入口隐藏」。不在本条设计该信号。

## 必改项汇总

| ID | 级别 | 阻断 C4 关门 |
|----|------|----------------|
| 本条 F-001 | required / med | **是**（双语发送键；不阻断 I-033-023 裁决） |
| 本条 F-002 | recommended / low | **否** |
| 本条 F-003 | recommended / low | **否（本切片）**；C4 关门仍须交付占用位隐藏 |
| I-033-023 | required 信息项 / collecting | **是**（C4 API/UI/关门门禁；本条不选择形状） |

开放 required finding = **1**（F-001）。开放 required 信息项 I-033-023 仍为 `collecting`。本条不把任何 finding 标为 `accepted-residual` 或 `user-overruled`。

## 与既有意见的异同

- A-029 self `conditional` / open_required=0：方向与本条一致（切片不是关门；capability 形状不在本条选择；composer fail-closed）。本条**不采信**其测试结论，但本会话独立重跑后 8/8 与 92/1203 数字可重复。
- A-029 **未记录**缺失的 `schema.telegram.operator.send`。本条将其列为 required F-001。
- A-029 将 I-033-023 写为未决边界而非 finding。本条同样**不把它写成实现作假**，但按 P-005 将其登记为阻断 C4 关门的 required 信息项。
- A-001～A-028 与 C3 无关本切片；原文保留。

## 结论 + 建议给编排器/用户的下一步

**verdict = conditional。** 当前工作树已把 C3 会话列表/成绩单接到既有 Telegram Admin 页，并实现 I-033-009 的 10 秒单飞、失焦暂停、恢复即刷与请求去重；未得到 capability 结果时 composer/retry 保持 fail-closed，且没有任何代码把 `operatorCapability` 设为 `allowed`。本切片**没有**接通发送/retry，也**没有**实现 `getChatMember`。

即使 F-001 修复后，**C4 仍不能关门**，原因不是本条替用户做了 I-033-023 选择，而是：

1. **I-033-023 仍为 `collecting`**：C4 API/UI/关门的 required 信息项未闭合。三种互斥形状（独立 capability、成绩单附带、会话列表附带）必须由用户书面裁决；本条明确不选。
2. D-002 已冻结的混合发言权仍未落地：`getChatMember` 预检、60 秒 bot/chat 缓存、403 失效、显式重探、composer 按真实 `can_send` 启用。
3. 发送/失败/retry 状态机尚未接到 UI（本切片有意 fail-closed）。
4. 已绑定占用位时入口隐藏尚未交付（F-003）。
5. C4 端到端验证与关门向 independent 尚未发生。
6. F-001 发送键缺失在修复前也阻断关门文案完备。

建议 `/govern`：响应本条；补 F-001 双语键；**询问用户 I-033-023**（三种形状 + 建议）；在用户裁决并实现发言权/发送路径之前，保持 C4 `进行中`、R3 `active · 3/4`。不要把本条当作 C4 或 R3 关门证据。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
