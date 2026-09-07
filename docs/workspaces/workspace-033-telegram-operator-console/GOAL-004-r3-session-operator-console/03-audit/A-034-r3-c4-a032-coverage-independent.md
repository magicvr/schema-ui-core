---
doc_type: goal-audit
id: A-034-r3-c4-a032-coverage-independent
parent: GOAL-004-r3-session-operator-console
date: 2026-09-05
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
audit_type: finding-closure
scope: workspace-033 R3 C4 A-032 新增 recommended F-001/F-002/F-003 修复响应复审（当前工作树源码/测试/双语 catalog/composition·settings；复核 A-030 F-001/F-002/F-003 与 A-032 原文；不采信 A-033/E-019；不选择 I-033-023；不关闭 C4）
verdict: pass
open_required: 0
version: 0.1.0
---

# A-034 · R3 C4 A-032 recommended 覆盖钉独立复审（2026-09-05）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型** / **scope**：finding-closure · `[workspace-033-telegram-operator-console]` `GOAL-004-r3-session-operator-console` 的 R3 C4 **A-032 新增 recommended F-001/F-002/F-003 修复响应**（当前 HEAD `8abe5aca5b174a4ff6aeb420dc2ce24844038a78`；C4 与 A-032 响应仍在**未提交工作树**；对照 A-030 原文 F-001/F-002/F-003、A-032 原文 F-001/F-002/F-003、当前源码/测试/双语 catalog/composition settings 接线与本会话跑数）。**不采信 A-033/E-019。****不选择 I-033-023**。**不关闭 C4**。
- **verdict**：pass
- **open_required**：0
- **完整意见**：本文件（未超 32 KiB，无附件）

本意见不修改 `status` / 检查点 / `progress` / 方案正文 / `goal-tree` / 生产代码 / 测试代码。未读取或比较其他工作区正文。**A-001～A-033 原文及其 findings 全部保留、未改写。** 不把 A-029/A-031/A-033 self、E-017/E-018/E-019 或 Root E-022/E-023 当作成功依据。不把 recommended 升级为 required，不接受 residual，不 overrule。不自行关闭 C4。不在三种 capability API 形状中替用户选择，不实现 `getChatMember` / 发送 / retry。

## 范围与区间

- **工作区**：`workspace-033-telegram-operator-console`；canonical `docs/workspaces/workspace-033-telegram-operator-console/`；Root `GOAL-001-telegram-operator-console`；`primary_plan = VP-033-telegram-operator-console`；`shared_materials_catalog: none`（本条未把任何共享资料当作关闭证据或跨区权限）
- **HEAD / 工作树**：本会话 `git rev-parse HEAD` = `8abe5aca5b174a4ff6aeb420dc2ce24844038a78`（`docs(govern): close R3 C3 checkpoint`）。C4 与 A-030/A-032 响应仍在工作树：`telegram-admin-tab.tsx` / `.test.tsx`、`catalog.test.ts`、`en-US.json` / `zh-CN.json`、`runtime.go` / `runtime_test.go` / `settings_handler.go`、`composition.go` / `composition_telegram_test.go`，以及未跟踪的 E-017～E-019、A-029～A-033。本条以**当前工作树源码与本会话跑数**为准，不信任描述或 A-033/E-019 结论。
- **covered**：
  1. A-032 F-001：en-US / zh-CN 的 `schema.telegram.operator.send` 是否有不会被 `t()` fallback 或 `"Send as bot"` 误满足的精确测试
  2. A-032 F-002：`timelineFlightsRef` 同 chat pending 请求是否有真实单飞测试
  3. A-032 F-003：composition mux 是否用同一 Dispatcher 注册后的 settings `business_occupied=true` 做真实 JSON 断言；UI 是否覆盖 `business_occupied` 缺省 fail-closed；lease effect 是否依赖 `business_occupied` 以处理热更新
  4. A-030 原文 F-001/F-002/F-003 与 A-032 原文是否仍可在当前工作树核对
  5. composer/retry 仍 fail-closed；未接通发送/`getChatMember`
  6. 本会话独立重跑 Web 定向/全量、API `telegram`/`composition`/`docscheck`；写集外 `tsc` 基线
- **excluded**：改写 A-001～A-033；采信 A-033/E-019；选择独立 capability / 成绩单附带 / 会话列表附带；实现或修改 `getChatMember`、发送/retry、status/progress/goal-tree；全仓 `go test ./...`；C4 关门

## 成果（有证据）

| 主张 | 本条独立证据（不引用 A-033/E-019） |
|------|----------------------------------|
| 工作区绑定合格；共享资料目录为 `none` | `workspace.md` L1–16、L29–36、L47–51 |
| Charter `active` 0.4.0；VP-033 `active`；`vision_ref` 对齐 | `docs/vision/charter.md` L5–6；VP-033 L5–7 |
| HEAD 仍是 C3 关闭文档提交；C4/A-032 响应在工作树 | `git rev-parse HEAD`；`git status` 列出上述 web/api 文件为已修改；A-029～A-033 未跟踪 |
| A-030 原文仍为 independent `conditional` / `open_required: 1`；F-001 required/open，F-002/F-003 recommended/open | `A-030-r3-c4-ui-foundation-independent.md` L10–11、L86–119；本条不改写 |
| A-032 原文仍为 independent `pass` / `open_required: 0`；新增 recommended F-001/F-002/F-003 open | `A-032-r3-c4-a030-remediation-independent.md` L10–11、L121–157；本条不改写 |
| A-033 原文仍为 self `pass`，声称三项 `fixed` | `A-033-r3-c4-a032-response.md` L6–11、L26–41；**不是本条证据** |
| I-033-009 已用户裁决：10 秒单飞、失焦暂停、恢复立即刷新 | D-002；`00-meta.md` L52 |
| I-033-023 仍 `collecting`，本条不选择三种 API 形状 | `00-meta.md` L58 |
| `operatorCapability` 仍只被重置为 `"unknown"`，从不设 `"allowed"`；send 无 `onClick` | `telegram-admin-tab.tsx` L104、L376–378、L753 |
| 本会话定向/全量 Web 与 API 定向 **PASS**；`tsc -b` 被写集外错误阻断 | 见下方「本会话验证」 |

### 本会话验证（独立执行，2026-09-05）

对照当前工作树（非 HEAD blob）：

| 命令 | 结果 | 归类 |
|------|------|------|
| `apps/web` `npm test -- --run src/components/telegram-admin-tab.test.tsx src/i18n/catalog.test.ts src/i18n/schema-keys.structural.test.ts --reporter=verbose` | **PASS**（Test Files 3/3；Tests **28/28**：tab **12/12**、catalog **12/12** 含 send 精确钉、schema-keys 4/4） | 通过 |
| `apps/web` `npm test -- --run` | **PASS**：Test Files **92 passed (92)**；Tests **1208 passed (1208)**；Duration 10.07s | 通过 |
| `apps/api` `go test ./internal/channel/telegram ./internal/composition ./internal/docscheck -count=1` | **PASS**：telegram 5.176s；composition 23.329s；docscheck 0.226s | 通过 |
| `apps/web` `npx tsc -b --pretty false` | **FAIL**：`src/renderer/form-controls.tsx(946,11)` 与 `(947,11)` `TS2322`：`number \| undefined` 不能赋给 `string \| undefined` | **写集外**既有类型错误；`git status` 未列出该文件；不构成本切片 fail |
| `git diff --check`（C4 web/api 写集 10 文件） | 无 whitespace 报错（仅 CRLF 提示） | 通过 |

未把 skip 记为通过。未跑 e2e、未跑全仓 `go test ./...`、未跑 `npm run build` 的 vite 后半段（`tsc -b` 已失败）。`schema-keys.structural.test.ts` 只断言两套 catalog **键集合相同**，不能单独证明组件 `t()` 用到的每个 key 都存在。

## 对照成功标准

本条审的是 **A-032 新增 recommended F-001/F-002/F-003 在当前工作树是否已有可核对修正**，并复核 A-030 原文三项。C4 检查点关闭权在 `/govern`，不在本条。GOAL 级发言权/`getChatMember`/发送状态机仍属 C4 后续，不得写成已交付。

### 0) A-030 原文三项（复核，不改写 A-030）

| A-030 finding | 当前工作树 | 本条判定 |
|---------------|------------|----------|
| F-001 required：缺 `schema.telegram.operator.send` | 组件 L753 仍渲染该键；`en-US.json` L1036 `"Send"`；`zh-CN.json` L1036 `"发送"` | **生产缺口仍闭合** |
| F-002 recommended：10 秒边界与 `operatorRefreshRef` 单飞 | 测试 L168–177（9999ms 仍 1 次、再 1ms 为 2 次）；L197–230 pending sessions 下 hidden→visible 合并。本会话 **PASS** | **仍满足** |
| F-003 recommended：占用信号 + 入口隐藏 + 占用时不 acquire | settings `json:"business_occupied"` 无 `omitempty`；组件 L638 仅 `=== false` 渲染；L159 `!== false` 不 acquire；测试 L289–302 occupied true。本会话 **PASS** | **仍满足** |

### 1) A-032 F-001 · 发送文案精确测试（不会被 fallback / `"Send as bot"` 误满足）

A-032 F-001（recommended / low / 原件 **open**）：tab 测试 L148 `toContain("Send")` 会与图例 `"Send as bot"` 及缺键回退键名碰撞。

| 核对 | 当前工作树 | 本条判定 |
|------|------------|----------|
| catalog 精确钉 | `catalog.test.ts` L40–43：`enUS["schema.telegram.operator.send"]` **全等** `"Send"`；`zhCN[...]` **全等** `"发送"`。读的是 JSON catalog 对象，**不经过** `t()` / fallback | **满足**。缺键时值为 `undefined`，不等于 `"Send"`/`"发送"`；值为 `"Send as bot"` 也不相等 |
| fallback 链 | `catalog.ts` L7–9、L87–103：locale → en-US → 键名本身。本钉不走该链 | **不会被 fallback 误绿** |
| `"Send as bot"` 碰撞 | 图例键是 `schema.telegram.operator.composer`（en-US L1031 `"Send as bot"`）。send 钉断言的是**另一键**的精确值 | **不会被图例误绿** |
| 组件仍渲染该键 | L753 `{t("schema.telegram.operator.send")}`；按钮仍 `disabled`、无 `onClick` | 键仍被使用 |
| tab UI 子串钉 | 测试 L148 仍是 composer `textContent` `toContain("Send")` | **仍弱**；见覆盖残余。不重开 A-032 F-001，因 catalog 精确钉已独立会红 |

**响应侧判定：fixed。** 用户本轮指定复核的「en-US/zh-CN 精确测试不会被 fallback 或 Send as bot 误满足」已存在于 `catalog.test.ts` L40–43。本会话该测 **PASS**。

### 2) A-032 F-002 · `timelineFlightsRef` 同 chat pending 单飞

A-032 F-002（recommended / low / 原件 **open**）：无测试让同一 `chatId` 的两次 messages 请求重叠。

| 核对 | 当前工作树 | 本条判定 |
|------|------------|----------|
| 源码去重 | `loadTimeline` L266–271：命中 in-flight 则 await 后 return；L293–298 完成后删除 | 仍按 chat 复用 |
| 真实 pending 单飞钉 | 测试 L238–287：messages 返回 unresolved promise；两次 `flushPromises` 后 `timelineCalls === 1`；再 `click` 已选中的 `[data-telegram-session='8001']`（onClick L670–674 **直接** `loadTimeline(chatId)`，不是走 `refreshOperatorSurface`）；`timelineCalls` 仍为 1；释放后仍为 1。本会话该测 **PASS** | **满足** 同 chat pending 单飞 |
| 与 sessions 单飞分离 | L197–230 仍只钉 `operatorRefreshRef`；L238–287 计数的是 `/sessions/8001/messages` | **不是** sessions 钉的重复 |

**响应侧判定：fixed。**

### 3) A-032 F-003 · composition mux JSON、缺字段 UI、lease 热更新依赖

A-032 F-003（recommended / low / 原件 **open**）：composition mux 无占用 JSON 断言；UI 无缺字段钉；lease effect 省略 `business_occupied`。

| 核对 | 当前工作树 | 本条判定 |
|------|------------|----------|
| 同一 Dispatcher | `composition.go` L907–936：`disp := NewDispatcher()`，`Dispatcher` 与 `DispatcherState` 同指针。测试 L527 `tr.Dispatcher.RegisterCommand("status", ...)` 后，probe `tr.DispatcherState.HasBusinessHandlers()` 读的是同一实例 | **满足** |
| 真实 mux GET JSON | `composition_telegram_test.go` L525–533 先注册 command；L636–654 用**同一** `tr` 建 `muxLease`；L658–672 GET `/api/channel/telegram/settings`，`json.Decode` 后 `BusinessOccupied` 必须为 true。本会话 `go test ./internal/composition` **PASS** | **满足**「注册后 settings `business_occupied=true` 真实 JSON 断言」 |
| settings 与 operator 同一探针 | `composition.go` L629–633 与 L644–646 均为 `tr.DispatcherState != nil && tr.DispatcherState.HasBusinessHandlers()` | 仍成立 |
| handler 覆盖 | `settings_handler.go` L108–113 GET/PATCH 经 `status()` 覆盖；`RuntimeManager.Status()` L438–451 **仍不填**该字段 | 字段只从 probe 进入 JSON |
| UI 缺省 fail-closed | 测试 L304–317：`business_occupied: undefined` 经 `JSON.stringify` **省略该键**；断言 `[data-telegram-operator]` 为 null 且 calls 不含 `/lease/acquire`。源码 L159 / L260 / L638 均 `=== false` 才放行。本会话 **PASS** | **满足** |
| lease effect 热更新依赖 | L252 `[callLease, loadState, status?.business_occupied, status?.mode]`。占用从 `false`→非 `false` 时 effect 重跑：cleanup 在 `leaseHeld` 时 `release`，新 effect 走 L159 不 acquire | **源码满足** A-032「把 `business_occupied` 纳入依赖」 |

**响应侧判定：fixed**（A-032 点名的三条：mux 注册后 JSON true、缺字段 UI、lease deps）。

## 覆盖残余（不升 required，不新开 recommended）

下列为紧密度，**不**构成对本 finding-closure 的阻断，也**不**重开 A-030/A-032 findings：

1. **tab UI 仍用弱子串**：`telegram-admin-tab.test.tsx` L148 `toContain("Send")` 仍会与 `"Send as bot"` 碰撞。A-032 F-001 要的精确锁已在 `catalog.test.ts` L40–43；本条不因 leftover UI 子串重开 F-001。没有独立断言发送**按钮**自身文本。
2. **composition mux 未钉 PATCH JSON**：真实 mux 只断言 **GET** 占用字段（`composition_telegram_test.go` L658–672）。PATCH 占用仍只在 handler 单测、且探针是字面 `return true`（`runtime_test.go` L67、L143–145）。A-032 原文写的是 GET/PATCH **或** 注册后翻转；后者已由 GET + `RegisterCommand` 满足。
3. **占用翻转无行为测试**：lease deps 已含 `status?.business_occupied`，但没有任何用例在同一挂载周期把占用从 `false` 改成 `true` 再断言 `release`。源码路径可核对；测试只锁首次加载的 occupied/missing。
4. **zh-CN 发送文案无渲染断言**：精确锁在 catalog 对象上，不经过 `I18nProvider` 默认 en-US 的 tab 渲染。缺键会红；渲染路径仍只覆盖 en-US。

以上不阻断本条 pass，也不构成对 I-033-023 的裁决。

## Findings

本条 **无新增 required**。本条 **无新增 recommended**。开放 required finding = **0**。

## 必改项汇总

| ID | 级别 | 阻断 C4 关门 |
|----|------|----------------|
| A-030 F-001 | required / med · 响应侧仍 **fixed**（原件仍 open） | 原文「是」；生产缺口仍闭合 |
| A-030 F-002 | recommended / low · 响应侧仍 **fixed** | 否 |
| A-030 F-003 | recommended / low · 响应侧仍 **fixed** | 否（本切片）；C4 关门仍须保持占用隐藏 |
| A-032 F-001 | recommended / low · **本条响应侧 fixed**（原件仍 open） | 否 |
| A-032 F-002 | recommended / low · **本条响应侧 fixed**（原件仍 open） | 否 |
| A-032 F-003 | recommended / low · **本条响应侧 fixed**（原件仍 open） | 否 |
| I-033-023 | required 信息项 / collecting | **是**（C4 API/UI/关门门禁；本条不选择形状） |

开放 required finding = **0**。开放 required 信息项 I-033-023 仍为 `collecting`。本条不把任何 finding 标为 `accepted-residual` 或 `user-overruled`。

## 与既有意见的异同

- A-030 independent `conditional` / open_required=1：原件 F-001/F-002/F-003 全文保留。本条复核生产缺口仍闭合，不改写 A-030。
- A-032 independent `pass` / open_required=0：原件新增 recommended F-001/F-002/F-003 全文保留。本条独立核对当前工作树后，三项在**响应侧**记为 `fixed`，不改写 A-032。
- A-033 self `pass`：方向与本条一致（三项已处理；不选 I-033-023；不关 C4）。本条**不采信**其测试结论或「fixed」声明；上表与跑数为独立核对。
- A-031/A-029 self 不是本条证据。A-001～A-028 与本切片无关；原文保留。

## 结论 + 建议给编排器/用户的下一步

**verdict = pass。** A-032 的三条 recommended 覆盖钉在当前工作树有可重复核对的修正：

1. `catalog.test.ts` L40–43 对 en-US `"Send"` / zh-CN `"发送"` 做 catalog 全等断言，缺键或 `"Send as bot"` 都会红，不经过 `t()` fallback。
2. tab 测试 L238–287 在 messages 请求 pending 时再次 click 同一 chat，`timelineCalls` 保持 1，钉的是 `loadTimeline` / `timelineFlightsRef`，不是 sessions 单飞。
3. composition 真实 mux 在同一 Dispatcher `RegisterCommand` 之后 GET settings，JSON `business_occupied` 为 true；UI 缺字段隐藏入口且不 acquire；lease effect 依赖含 `status?.business_occupied`。

本条 **open required finding = 0**。覆盖残余见上节，不升 required。

即使如此，**C4 仍不能关门**，原因不是本条替用户做了 I-033-023 选择，而是：

1. **I-033-023 仍为 `collecting`**：C4 API/UI/关门的 required 信息项未闭合。三种互斥形状必须由用户书面裁决；本条明确不选，也不实现 `getChatMember` / 发送 / retry。
2. D-002 已冻结的混合发言权仍未落地：`getChatMember` 预检、60 秒 bot/chat 缓存、403 失效、显式重探、composer 按真实 `can_send` 启用。`operatorCapability` 唯一写入仍是重置为 `"unknown"`。
3. 发送/失败/retry 状态机尚未接到 UI（send 仍硬编码 `disabled`、无 `onClick`）。
4. 写集外 `form-controls.tsx` L946–947 `tsc` 失败是既有基线，不是本切片回归。

建议 `/govern`：响应本条；可将 A-032 F-001/F-002/F-003 在响应侧记为 `fixed`（原件不改写）；**询问用户 I-033-023**（三种形状 + 建议）；在用户裁决并实现发言权/发送路径之前，保持 C4 `进行中`。不要把本条当作 C4 或 R3 关门证据。

### 声明

本意见不修改 status/progress；响应由 /govern 处理。
