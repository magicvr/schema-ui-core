---
id: GOAL-007-r5-examples-contract-verification
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-08-01
version: 0.9.0
---

# 决策记录 · GOAL-007

## 信息需求与阶段门禁

P-005 信息台账维护在 [00-meta.md](00-meta.md)。本目标 `I-007-001` 在 R5 **验收前**须验证；父目标 `I-PROTO-003` 在 R5 **验收/关门前**须合法闭合。本目标不修改 Root 门禁状态。

## D-001 · 立项 R5 纳入域范例与契约验证子目标

**日期**：2026-07-31
**状态**：accepted

**决定**：

在 `GOAL-001-mvp-admin-foundation` 下创建 `GOAL-007-r5-examples-contract-verification`，将 R5 范围限定为：为 R2 冻结基线 v0.1.3 的 11 个纳入域交付可观察范例页/场景 + 结构/行为验证路径，并在 R5 验收前闭合父目标 `I-PROTO-003`。Root 保持 `active`，纲领进度维持 `4/6`。

**为什么**：

- Root 路线图把 R5 定义为「纳入域范例与契约验证」，且 `I-PROTO-003` 要求「每条纳入能力的范例页路径与自动化/手工验证入口」在验收前闭合。
- R3/R4（GOAL-005/006）已完成 Admin 外壳与账号权限链路；R5 承接这两条既有路径，把其余纳入域补齐为可观察范例与验证入口。
- [I-PROTO-001 v0.1.3] §3 已给出 P0/P1/P2 范例候选，协议清单 §2.5 提供可复用的信息性场景，契约发现有明确来源。

**未选方案**：

- 在 R4 目标中补入范例页：会越过 R4 已关门边界，并提前触碰 `I-PROTO-003` 验收门禁。
- 一次性实现「完整协议 Renderer / 全 registry」：超出 v0.1.3 明确 include/include-partial 边界，违反 MVP 非目标。
- 以「构建可跑」作为 R5 验收证据：无法覆盖各纳入域的结构/行为契约，`I-PROTO-003` 要求的是逐域范例 + 验证入口。

**影响**：

本目标进入 `active` 规划阶段；R5 验收前验证 `I-007-001` 并闭合 `I-PROTO-003`。当前不修改 `apps/*` 实现，不改变父目标 `I-PROTO-003` / `I-PROTO-004` 状态。

## D-002 · 采用「契约发现登记 → 范例实现 → 结构/行为验证 → 验收关门」的 R5 路线

**日期**：2026-07-31
**状态**：accepted

**决定**：

先建立 `I-007-001` 登记表（每纳入域范例路径 + 自动化/手工验证入口），再据此实现范例页与验证路径；在此之前不把待确认取舍写成实现事实。实现完成后必须补结构/行为证据与阶段自审，才可讨论 `done`。

**为什么**：

- 沿袭 R3（GOAL-005 D-002/D-005）与 R4（GOAL-006 D-002）已建立的信息门禁纪律：开放 required 信息项是验收/关门的门禁。
- `I-PROTO-003` 的闭合要求「每纳入能力都有范例路径与验证入口」可逐项核对；先登记再实现可避免「实现了但没有验证入口」或「登记了但不可执行」的验收落差。

**未选方案**：

- 直接实现页面再补登记：会把规划缺失伪装成已决定行为，验收时难逐域核对。
- 仅以 dev server 可打开作为 R5 关门证据：无法覆盖各纳入域的结构/行为契约。

**影响**：

`I-007-001` 与父目标 `I-PROTO-003` 是 R5 验收/关门门禁；必要时须按 P-004 由用户裁决 fixed、accepted-residual 或 user-overruled，不能静默放行。

## D-003 · 落地 R5 阶段 1：契约发现登记表（2026-07-31）

**日期**：2026-07-31
**状态**：accepted

**决定**：

执行 R5 路线图第 1 阶段，落盘 `I-007-001` 登记表为 [attachments/I-007-001-registry.md](attachments/I-007-001-registry.md) v0.1.0：对照 [I-PROTO-001 v0.1.3 §3] 候选与协议清单 §2.5，为 11 个纳入域登记范例页/场景路径与结构/行为验证入口；D-APP / D-PERM 复用产物（R3/R4）与其可执行验证命令（`npm test`、`npm run build`、`go test ./...`、`go build ./...`）一并核验入账。`I-007-001` → `verified`（登记层面），执行验证仍属阶段 3，不据此放行验收/关门。

**为什么**：

- A-001（independent，pass）建议 `/govern` 响应并推进 `I-007-001` 登记工作；用户裁决「不需要自审，直接推进」。
- 登记表以冻结基线 [I-PROTO-001 v0.1.3] 与协议清单 §2.5 为权威输入，避免把待确认取舍写成实现事实；D-APP/D-PERM 已有产物只登记复核，不重复实现。

**未选方案**：

- 直接进入阶段 2 实现范例页：会越过「先登记后实现」的既定路线（D-002），且 `I-007-001` 登记为阶段 2 的输入。
- 将 `I-007-001` 直接标 `verified` 且把待执行的逐域验证写成已完成：违反 P-005，本决策只确认登记完成。

**影响**：

阶段 1 完成；阶段 2（范例页/场景实现）、阶段 3（结构/行为验证）、阶段 4（验收与关门）未开始。`I-PROTO-003` 仍 open；R5 验收/关门前须以阶段 3 可执行证据闭合。

## D-004 · 落地 R5 阶段 2：范例页/场景实现子方案（2026-07-31）

**日期**：2026-07-31
**状态**：accepted

**决定**：

采用「三个批次 + I-PROTO-004 决策时点」的阶段 2 实施子方案，按登记表为未覆盖域落地可观察范例页：

- **批次 2a · D-DATA / D-TABLE 范例与 Go 数据支撑**：Go 新增列表/详情数据 API（对照 `I-PROTO-002` 最小 API 纪律，落在 `apps/api/internal/handler/`）；Web 新增数据表格组件（排序声明 + 基础列表交互）；范例页 `search-form-table` / `data-table`。
- **批次 2b · D-FORM 控件与 D-ACT 动作**：Web 实现 §5 白名单的 2.6/2.7 表单控件子集（`form-controls-*`）；复用 R4 `executeAction` 时序引擎，落地非批量行动作；范例页 `admin-list-edit-lifecycle`。
- **批次 2c · D-EXPR 与 Renderer 接线**：落地 `form-with-reactions` 范例页（复用 `evaluateExpression`）；在 `apps/web/src/renderer/` 完成最小 Renderer 接线（D-COMP，resolve R4 推荐跟踪项 F-002），仅覆盖已纳入白名单 type，不做完整 registry。
- **`I-PROTO-004` 决策时点**：**阶段 3 结构校验实现前**决策 vendor vs pin（影响 `node`/`page`/`action`/`reaction` schema 校验命令落地方式）；阶段 2 范例页不依赖该决策，可按已 vendor 的 `app-manifest.schema.json` 与本地样例先行实现。`I-PROTO-004` 保持 `non-blocking` / open。
- 每个批次完成即更新 `I-007-001-registry.md` 对应行「现状」为已发生事实，并在 `02-execution` 记录实际命令/证据；阶段 2 全部落地后再进入阶段 3。

**为什么**：

- 登记表 §3 将 P1/P2 候选归组为 `search-form-table` / `data-table`（D-DATA+D-TABLE）、`admin-list-edit-lifecycle`（D-ACT+D-DATA+D-FORM）、`form-with-reactions`（D-EXPR+D-FORM）——按依赖排序：先数据/表格（Go 支撑 + 前端组件），再控件/动作，最后表达式与 Renderer 接线。
- Go 现仅 `accounts/me` + `health`，无列表/详情 API（登记表 D-DATA 行「需补」），故批次 2a 先补数据支撑，避免前端范例无真实数据源。
- `apps/web/src/renderer/` 现仅 `permissions.ts`（空分层，登记表 D-COMP 行已标注），Renderer 接线是 D-COMP 范例页的前置，故放在批次 2c 且限定白名单 type，不越 R5 非目标。
- 用户确认「不需要自审，直接推进」，且选择「先落实施子方案再实现」——本决策是实施前的最小子方案固化，不替代逐批次的 03 执行记录。

**未选方案**：

- 一次性实现全部 6 个未覆盖域范例页：范围与证据分散，验收时难逐域核对，且 Go 数据支撑是共同前置，按批落地更可验证。
- 先决策 `I-PROTO-004`（vendor vs pin）再实现：该决策只影响阶段 3 的 schema 校验命令形态，不阻断阶段 2 范例页；提前决策会把 non-blocking 项升级为阶段 2 前置，无必要。
- 阶段 2 直接复用 dev server 可打开作为关门证据：不满足 `I-PROTO-003`「每纳入域范例路径 + 验证入口可核对」，验证执行仍归阶段 3。

**影响**：

阶段 2 进入实施准备；本决策不修改成功标准、不改变 `I-PROTO-003` / `I-PROTO-004` 状态、不抬升 Root `progress`（仍 4/6）。批次 2a 完成后应补阶段自审（self）再进入批次 2b/2c。

## D-005 · 落地批次 2b：D-FORM 控件与 D-ACT 动作（2026-07-31）

**日期**：2026-07-31
**状态**：accepted

**决定**：

按 D-004 批次划分实施并记录批次 2b：实现 D-FORM §5 白名单控件表面（`apps/web/src/renderer/form-controls.ts` / `.tsx`，含 2.6/2.7 版本/capability 门禁与 `defaultValue` wire 类型校验）与 D-ACT 非批量行动作（`row-action.ts` 复用 R4 `executeAction` 时序引擎），补 Go 编辑生命周期支撑（`PATCH`/`DELETE /api/records/{id}`，mutex 保护 + `validatePatch` fail-closed），落地 `form-controls` 与 `list-edit-lifecycle` 范例页及接线；批次完成后更新 `I-007-001-registry.md` 与 `02-execution` 为已发生事实，并补阶段自审（A-003）。

**为什么**：

- D-004 将 P1 `admin-list-edit-lifecycle`（D-ACT+D-DATA+D-FORM）归入批次 2b；批次 2a（数据/表格）完成后按依赖推进控件/动作，Go 仅需扩展编辑/删除，无新数据源前置。
- D-FORM 依赖冻结的 §5 白名单与 2.6/2.7 版本/capability 规则，实现时须 fail-closed（非白名单 type、wire 类型失配、能力缺失均拒绝），不越 R5 非目标（完整 registry / 多选批量）。
- D-ACT 复用 R4 已审计的 `executeAction` 引擎，仅做 UI 参数包装，不重新发明权限时序；Q1=否排除批量。

**未选方案**：

- 批次 2b 与 2c 合并实现：证据分散，Renderer 接线（D-COMP）是独立前置，仍按 D-004 分批落地更可验证。
- 在 Go 侧引入 DB/持久化：MVP 无 DB，`sync.RWMutex` 进程内数据集足够支持范例生命周期，越界即偏离非目标。

**影响**：

批次 2b 完成；`I-PROTO-003`（父目标）与 `I-PROTO-004` 仍 open，验收/关门须以阶段 3 可执行证据闭合，`I-PROTO-004` 在阶段 3 结构校验实现前决策。Root `progress` 仍 4/6。批次 2c（D-EXPR + Renderer 接线）待推进。

## D-006 · 落地批次 2c：D-EXPR 反应引擎与 D-COMP 最小 Renderer 接线（2026-08-01）

**日期**：2026-08-01
**状态**：accepted

**决定**：

按 D-004 批次划分实施并记录批次 2c：实现 D-EXPR 反应引擎（`apps/web/src/renderer/reactions.ts`，复用 `evaluateExpression`，frozen $context 表达式子集）与 D-COMP 最小 Renderer 接线（`render.ts`/`render.tsx`，whitelist form/section/table fail-closed，resolve R4 推荐跟踪项 F-002），落地 `form-with-reactions` 范例页（Admin/Viewer 角色 + audit feature 切换演示字段显隐/禁用）及注册接线；批次完成后更新 `I-007-001-registry.md`（v0.4.0）与 `02-execution` 为已发生事实，并补阶段自审（A-004）。

**为什么**：

- D-004 将 P2 `form-with-reactions`（D-EXPR+D-FORM）与 Renderer 接线（D-COMP，resolve R4 F-002/F-003「Renderer 集成层尚未消费引擎」跟踪项）归入批次 2c；批次 2b（控件/动作）完成后按依赖推进表达式与 Renderer。
- D-EXPR 复用 R3 已审计的 `evaluateExpression`（导航 visibleWhen 同源），以 `isValidExpression` 导出同一语法校验，不重复发明表达式语义；reactions 仅操作 frozen $context 快照，Q=否排除 field-value 触发。
- D-COMP 以「最小 Renderer 接线」落地：只 dispatch whitelist 的 form/section/table 节点，未知 type fail-closed 出 alert；不做完整 component registry，不越 R5 非目标。
- `reaction.schema.json` 结构校验仍随 `I-PROTO-004`（阶段 3 前决策），批次 2c 不假装已完成 schema 级验证。

**未选方案**：

- 批次 2c 实现完整 component registry / 全量 node 树渲染：超出 R5 非目标（完整 registry），且 `node`/`page` schema 校验须先决策 `I-PROTO-004`。
- 在 reactions 中引入 field-value 触发：`evaluateExpression` frozen 语法仅支持 `$context.user`/`$context.features`，扩展即偏离冻结子集。
- 跳过 Renderer 接线只做范例页：F-002 跟踪项要求 renderer 集成层消费引擎，否则「页面运行时已应用 D-PERM/D-EXPR」主张不成立。

**影响**：

批次 2c 完成；`I-PROTO-003`（父目标）与 `I-PROTO-004`（non-blocking）仍 open，验收/关门须以阶段 3 可执行证据闭合，`I-PROTO-004` 在阶段 3 结构校验实现前决策。Root `progress` 仍 4/6。阶段 2 全部落地，阶段 3（结构/行为验证）待推进。

## D-007 · 响应 A-005：以 fixed 路径闭合 F-001～F-004（2026-08-01）

**日期**：2026-08-01
**状态**：accepted

**决定**：

用户裁决「不需要自审，直接推进」（P-004 §3.1），并按 **`fixed`** 路径响应 A-005（independent, conditional）的 4 项 required findings，F-005（recommended）同步：

- **F-001（fixed）**：`render.ts` 以 `resolveActionGate` 区分「属性缺省（→ absent default）」与「显式非法表达式（→ fail-closed + `ACTION_GATE_EXPRESSION_INVALID` 可核对错误）」；`tableActionGate` 改为返回 `{ visible, disabled, errors }`，显式非法 `visibleWhen` 拒绝/隐藏、非法 `disabledWhen` fail-closed 禁用；补非法 visible/disabled 回归测试。
- **F-002（fixed）**：`render.tsx` `FormView` 在渲染前经 `gateRenderFormFields` 解析字段并执行 D-FORM §5 whitelist、版本与 capability 门禁；非法 type、低版本或缺 capability 拒绝受影响字段并在页内出 `role="alert"` 确定错误；补 Renderer 级测试（缺 capability 拒绝、未知 type 拒绝）。
- **F-003（fixed）**：`checkFormCapabilities` 为任一出现 `defaultValue` 的字段增加 `protocolVersion >= 2.7` + `form.controls.advanced` 双门禁（沿用冻结 §5「2.7 属性能力」行）；补 2.6、2.7 缺 capability、2.7 完整 capability 三组回归测试。
- **F-004（fixed）**：Renderer node whitelist 从 form/section/table **扩展至冻结 §5 全部 node type**（layout：grid/section/tabs；data/action：text/table/recordView/actionButton；form），逐 type 落 `parseRenderNode` dispatch 与 `RenderPage` 渲染 + 测试；**不做范围缩小**，故不需要 Root 覆盖表修订或新版 v0.1.3 冻结。
- **F-005（fixed）**：Root `00-meta.md` R5 行同步「批次 2c 完成 / 阶段 2 全部落地」。

**为什么**：

- A-005 为独立审计，其对同一冻结范围的 F-001～F-004 主张经 `/govern` 直接读代码核验成立（`render.ts` `gateAction` 对非法表达式返回 `defaultValue`、`render.tsx` `FormView` 未调 `checkFormCapabilities`、`checkFormCapabilities` 对 `defaultValue` 无 2.7+advanced 门禁、renderer whitelist 静默窄于 §5 初表），与 A-002～A-004（self, pass）构成同 scope 分歧，须按 P-004 §3.2 用户裁决。
- 四项均为真实 fail-open/契约缺口，`fixed` 是最低成本、最可核对的闭合路径；不改变 v0.1.3 冻结基线，不抬升 `progress`，不触碰 `I-PROTO-003` / `I-PROTO-004` 状态。
- 用户本轮裁决「不需要自审」，延续此前节奏（D-003/D-004）。

**未选方案**：

- `accepted-residual`：接受残余风险需书面范围与复审触发，且四项缺口会在阶段 3 schema/fixture 对照中复发，成本高于直接修正。
- `user-overruled`：代码证据支持 A-005 主张，驳回会使 R5 验收证据链失真。
- 先做 self 审计再响应：批次自审 A-002～A-004 已覆盖各批实现事实，本轮以代码核验 + 回归测试替代新增自审，证据充分。

**影响**：

F-001～F-004 以 `fixed` 合法闭合，F-005 同步；`npm test` 14 文件 **173 项** / `npm run build` / `go test ./...` / `go build ./...` / Edge headless 实测全绿。`I-PROTO-003`（父目标，required，R5 验收/关门）与 `I-PROTO-004`（non-blocking，阶段 3 前决策）仍 open；阶段 3（结构/行为验证）与阶段 4（验收/关门）未开始。Root `progress` 仍 `4/6`。

## D-008 · 响应 A-006 + I-PROTO-004=vendor + 进入阶段 3（2026-08-01）

**日期**：2026-08-01
**状态**：accepted

**决定**：

1. **P-004 §3.1**：对 A-006（independent, pass）用户裁决「**不需要自审，直接推进**」。
2. **响应 A-006**：采纳 `verdict: pass`；开放 required=0；A-005 F-001～F-005 闭合成立。A-006 recommended F-001/F-002（UI/范例页一致性）跟踪至后续产品化或验收补强，**不阻断**阶段 3；F-003（登记表 §4 计数）随本轮登记表同步闭合。
3. **`I-PROTO-004` = vendor**（父目标信息项）：将上游 `schema-ui-docs@2.7.0`（commit `ca9e5fe207c169d6957bdd4f9a968deaf3bd2d7b`）的结构 schema（`node`/`page`/`action`/`reaction`/`component-registry`，叠加既有 `app-manifest`）与纳入 fixture suites **vendor 进仓**并 SHA pin（`apps/web/src/protocol/upstream/provenance.json`），校验离线可复现；**不**采用远程 pin-only 校验路径。
4. **进入阶段 3**：落地可执行结构校验（Ajv draft-07）与纳入行为 fixture 对照（`apps/web/src/protocol/conformance/*` + `stage3-fixtures.test.ts`）。

**为什么**：

- A-006 复核确认阶段 2 修复可重复，建议主路径进入阶段 3 且先决 `I-PROTO-004`。
- vendor 与既有 R3 `app-manifest.schema.json` + fixture pin 模式一致，CI/离线可复核，维护成本明确（升级需新 commit pin + SHA）。
- 阶段 3 是 `I-PROTO-003` / R5 验收前的结构与行为证据层；不抬升 Root `progress`，不关闭 `I-PROTO-003`。

**未选方案**：

- **pin 远程校验**：依赖网络与上游可用性，与「可 fork 基架 / 离线证据」方向不符。
- 先做 A-006 self 审计再推进：独立 pass + 既有批次 self（A-002～A-004）已覆盖；用户明确跳过。
- 将 upstream `reactions` 全套强制映射到 MVP host：上游为 multi-round `$deps` 字段写引擎，超出冻结 MVP D-EXPR（仅 `$context` visible/disabled）；选择 **vendor+account 排除** 而非静默伪通过。

**影响**：

- 父目标 `I-PROTO-004` → **verified**（vendor 策略与 pin 证据落盘）。
- 阶段 3 实施完成（结构校验可跑 + 纳入 suites 执行/覆盖记账）；阶段 4（验收/关门）与 `I-PROTO-003` 仍 open。
- Root `progress` 仍 `4/6`。

## D-009 · 响应 A-007：F-002 fixed + F-001 矩阵/入口升格 + request-construction residual（2026-08-01）

**日期**：2026-08-01
**状态**：accepted

**决定**：

1. **P-004 §3.1**：对 A-007（independent, conditional）用户裁决「**不需要自审，直接推进**」。
2. **F-002 → `fixed`**：为 `reaction` 与 `action` 结构校验各补可核对断言——合法样本 `expect(ok).toBe(true)`、非法样本 `expect(ok).toBe(false)`（reaction 合法样本须含 `dependencies`）；落在 `apps/web/src/protocol/conformance/stage3-fixtures.test.ts`。
3. **F-001 → 组合闭合**：
   - 在登记表 [I-007-001-registry.md](attachments/I-007-001-registry.md) **§2b** 落盘 **include suite 执行矩阵**（suite → executed / excluded / deferred|residual），禁止把「阶段 3 完成 + 测试全绿」读成全部冻结 include suite 已 host 对照通过。
   - **`reactions`**：将 MVP `$context` 验证入口（`reactions.test.ts` + `/form-with-reactions`）升格为 **正式 D-EXPR 行为验证入口**；显式声明上游 multi-round `$deps` suite **不在** MVP 语义子集（对齐 D-008，**不**伪通过）。
   - **`request-construction` non-batch** → **`accepted-residual`**（用户书面「… request-construction … residual」；本决策将 R1/R2 标签歧义按 **residual** 语义解释，**不**采用本轮落地统一 host adapter 的 fixed 路径，也**不**修订 v0.1.3 冻结表）：
     - **范围**：冻结 §2 include 的 non-batch case **不**作为 stage3 全量 host 对照绿项；batch 仍 Q1 排除。
     - **缓解 / 监控**：suite vendor+SHA pin 保留；partial 路径由 `records.ts` / `row-action.ts` 单元与范例覆盖；矩阵与登记表诚实记账 executed=0。
     - **复审触发**：落地统一 request-construction host/adapter 时；或进入 R6 集成验收前；或用户要求撤回 residual。
4. **A-007 recommended** F-003～F-005 继续跟踪，**不阻断**进入阶段 4 材料准备。
5. **不**关闭父目标 `I-PROTO-003`；**不**放行 R5 验收/关门；**不**抬升 Root `progress`（仍 4/6）。

**为什么**：

- A-007 核验阶段 3 vendor/Ajv/多数 suite 真实可复核，同时指出 include 面 0 执行与 action/reaction 断言空洞——F-002 低成本可 `fixed`；F-001 需显式矩阵 + 入口/residual 裁决，避免验收叙述静默覆盖缺口。
- `reactions` 排除理由已在 D-008 成立；升格 MVP 入口满足「有可执行行为验证」而不假装上游 suite 已跑。
- non-batch request-construction 统一引擎超出本轮 R5 最短路径；用户选 residual 并书面留痕，符合 P-003/P-004 三路径之一。

**未选方案**：

- **request-construction R1 fixed（统一 adapter）**：工作量大，可在 residual 复审触发时再做；本轮用户明确 residual。
- **修订冻结 v0.1.3 为 include-partial**：触及 Root 覆盖表新版本与重估，范围大于必要。
- 先做阶段 3 self 再响应：用户裁决跳过。

**影响**：

- A-007 开放 required → **0**（F-002 fixed；F-001 矩阵 + reactions 升格 + request-construction accepted-residual）。
- 登记表 → **v0.7.0**；`npm test` **331** 项。
- 阶段 4 / `I-PROTO-003` 仍 open；可进入验收材料准备，但闭合 `I-PROTO-003` 仍须逐域证据评审与用户确认。

> **部分 supersede**：`request-construction` non-batch 的 **accepted-residual** 路径由 **D-010** 更正为 **`fixed`**（用户澄清）。矩阵 / reactions / F-002 部分仍有效。

## D-010 · 更正 A-007 F-001：request-construction non-batch → fixed（2026-08-01）

**日期**：2026-08-01
**状态**：accepted

**决定**：

用户澄清：`request-construction` 要求 **`fixed`**（可执行 host/adapter 对照），**不是** accepted-residual。

1. **supersede** D-009 中关于 `request-construction` non-batch 的 residual 范围/缓解/复审触发。
2. 落地 `apps/web/src/protocol/conformance/request-construction.ts`（`constructRequest`），在 stage3 对 **全部 non-batch** case 执行 `expect(constructRequest(input)).toEqual(expected)`。
3. **batchRequest** 仍按冻结 Q1=否 **排除**（不改 v0.1.3）。
4. 登记表 §2b 矩阵更新为 non-batch **executed**；审计响应与执行记录更正 residual 表述。

**为什么**：

- 用户原文「R1 residual」存在 R1=fixed 与 residual 的标签冲突；澄清后以 **fixed** 为准。
- fixed 满足 A-007 F-001 必改项「落地可执行 host/adapter 对照」，且不静默伪通过 batch。

**未选方案**：

- 维持 residual：与用户澄清冲突。
- 同时执行 batch：违反 Q1 / include-partial D-ACT。

**影响**：

- A-007 F-001 中 request-construction 闭合路径 = **`fixed`**（+ 既有矩阵与 reactions 升格）。
- `npm test` **395** 项（+64 non-batch）；登记表 **v0.8.0**。
- 仍 **不**关闭 `I-PROTO-003`；阶段 4 未开始；Root `progress` 仍 4/6。
