---
id: GOAL-007-r5-examples-contract-verification
doc: decision
status: active
parent: GOAL-001-mvp-admin-foundation
created: 2026-07-31
updated: 2026-07-31
version: 0.4.0
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
