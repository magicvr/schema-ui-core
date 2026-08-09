---
id: A-002
doc: audit
title: S6 · C1–C4 独立关门审计（recordSource + settings 表单页 + 证据）
status: open
source: independent
scope: GOAL-007-s6-settings-form-page · C1–C4 close-out readiness
verdict: pass
auditor: grok-4.5 (/audit independent)
created: 2026-08-09
updated: 2026-08-09
version: 0.1.0
parent: GOAL-007-s6-settings-form-page
---

# A-002 · S6 C1–C4 独立关门审计（independent · 2026-08-09）

## 结论

**pass**（C1–C3 实现与证据可核对；C4 流程门禁「independent 意见」本条已落盘；用户书面确认与 Root 恢复 `done` 仍属编排/用户裁决，不在本条代改）。

开放 **required** findings = **0**。开放 **recommended** = 2（初态 idle 竞态、e2e 浏览器补跑）。

## 范围与区间

| 项 | 值 |
|----|-----|
| 工作区 | `workspace-007-localization-and-system-settings`（`workspace.md` 校验：`root_goal` = `GOAL-001-localization-and-system-settings`，`canonical_scope` 匹配，`shared_materials_catalog: none`） |
| 被审目标 | `GOAL-007-s6-settings-form-page` |
| audit_type | `close-out`（D-001 §4：C4 关门 = independent） |
| scope | C1 renderer `recordSource`；C2 settings schema 重构；C3 测试/证据；C4 治理收口就绪度 |
| 对照台账 | `00-meta` C1–C4；`D-001`；`E-001`；`A-001`（self） |
| 信息门禁 | `I-001` non-blocking **verified**（无到期 required） |
| 共享资料 | 无 |
| 本审计未改 | `status` / `progress` / goal-tree / 方案正文 |

## 成果（有证据）

### C1 · renderer `recordSource` 预填（ADR-0021）

| 主张 | 证据 |
|------|------|
| 解析放行 `title`/`titleKey`/`recordSource` | `apps/web/src/renderer/render.ts` `parseRenderNode` form 分支 + `RenderFormNode.props` |
| `resolveResponsePath` 点号路径 | `render.ts` export；`render.test.ts` 3 例 |
| capability `form.record.load` 门禁 | `form-controls.ts` `FORM_RECORD_LOAD_CAPABILITY`；`render.tsx` `useRecordSourcePrefill` + `hasRequiredCapability` |
| GET 预填 + `responseMapping` | `useRecordSourcePrefill` → `constructRequest({kind:"recordSource"})` → 映射字段 |
| search 模式拒绝 | 同上；单测 “rejects recordSource on search-mode forms” |
| empty `responseMapping` 经 constructor fail-closed | `request-construction.ts` `EMPTY_RESPONSE_MAPPING`（非 renderer 手写旁路） |
| loading / error UI | `FormView` skeleton / `role="alert"` |
| reload 重预填 | `useEffect` deps 含 `reloadToken`；`FormInner` `key={reloadToken}` |
| form 标题 | `FormInner` `<h2>` + `titleKey`/`title` |
| `${formId}:submit` 只读门禁 | `canEdit` + `fieldDisabled` + submit `disabled`；权限树 `permissions.ts` formSubmit target |
| `actionId` 回退 + actionButton 权限 | `invokeAction` `actionRef \|\| actionId`；`ActionButtonView` 按 `props.key` |
| 单测 | `render.test.tsx` S6 块 + `render.test.ts` path 解析；本会话复跑相关套件 **60/60** |

### C2 · settings schema 重构

| 主张 | 证据 |
|------|------|
| 四类内联 form + `recordSource` identity mapping | `apps/api/internal/modules/settings/schema/settings.json`：`settings-general/branding/localization/appearance` |
| 5 个 request action 保留（4 PATCH + reset POST） | 同上 `actions`；无 `open*` modal / 无 table / 无 recordView |
| Restore defaults = actionButton | `settings-reset`：`actionId: resetSettings`、`permissionIntent: edit`、`key: reset`、`confirm`/`confirmKey` |
| body section 权限级联 | `permissionCascade.keys: ["edit"]` + `permissions.edit` → `settings.write` |
| meta 增 `form.record.load`，无 `table.sort` / `actions.row.request` | meta 列表核对 |
| catalog 键复用、无新增键 | `titleKey`/`labelKey`/`submitLabelKey` 均在 `zh-CN.json` / `en-US.json` 已有 `schema.settings.*` |
| 集成/结构测试 | `startup-config.test.tsx` 四类 + 预填 + PATCH + write 门禁；`schema-keys.structural`；`s5-denominator-render` 含 settings 双语面 |

### C3 · 测试与证据

| 主张 | 本会话独立复跑 |
|------|----------------|
| vitest **727/727**（40 files） | **确认** exit 0 |
| `npm run build`（tsc + vite） | **确认** exit 0 |
| `go test ./apps/api/...` | **确认** exit 0 |
| checkpoint `ebd0288` | 存在；stat 覆盖 E-001 所列路径（settings.json + renderer + tests + e2e） |
| e2e M3 已改写 | `apps/web/e2e/localization.spec.ts`：内联 General 表单 → 保存 → 投影；四类 heading + 恢复默认 button |
| e2e 环境降级 | `attachments/s6-e2e-env-block.log`；本会话 `netsh` 复核 **8011–8110** 仍排除（含 8080）；与 E-001 / A-001 诚实降级一致；C3 成功标准**已写明**单元覆盖为验收线 |

### 工作区 / 信息 / 历史意见

- Root 暂时回退承接 S6：GOAL-001 `D-003` 留痕；VP-007 保持 closed 边界与 D-001/D-003 一致。
- `I-001` verified；无到期 required 信息项。
- A-001（self，C1–C3，pass）与本条证据方向一致；本条独立复跑并补推荐项（见下）。

## 对照成功标准

| 检查点 | 判定 | 说明 |
|--------|------|------|
| **C1** | **pass** | 协议表面接线与单测/集成证据充分；见 F-001 推荐项（初态 idle） |
| **C2** | **pass** | schema 形态、action、capability、catalog 键与 D-001 一致 |
| **C3** | **pass** | 本会话复跑 vitest/build/go；e2e 环境 residual 已在 C3 成功标准内诚实承认 |
| **C4** | **partial（流程）** | 证据入库（E-001）+ self A-001 + 本 independent A-002 已齐；**用户书面确认**与 Root `done` 恢复**尚未**发生——属 `/govern` + 用户裁决，**不是**实现造假；本条不改 status |

## Findings

### F-001 · recordSource 首帧 `idle` 可短暂渲染空可编辑表单

- **严重度**：med
- **建议**：recommended（非 required）
- **状态**：open
- **描述**：`useRecordSourcePrefill` 初始 state 为 `{ status: "idle" }`。当 `node.props.recordSource` 已定义时，`FormView` 仅对 `loading`/`error` 短路；`idle` 会落入 `FormInner` 且 `prefillValues=null`，在 `useEffect` 将状态置为 `loading` 之前可出现**一帧（或极短窗口）空可编辑表单**。与 `FormView` 注释「never renders an editable blank form that could overwrite the record」及 fail-closed 意图不完全一致。reload 路径因先走 `loading` 较安全；首挂载路径是主要缺口。
- **证据**：`apps/web/src/renderer/render.tsx` `useRecordSourcePrefill` 初始 state + `FormView` 分支；对比 capability/search 错误路径会阻断 form。
- **建议修复**：`recordSource !== undefined` 时初始状态用 `loading`（或 `idle` 且有 recordSource 时与 loading 同等短路）；补单测「首帧不出现可提交空 form」。
- **为何非 required**：C1 验收条文未把「零帧 idle」写成硬门禁；实际用户提交需跨过极短窗口；capability/search/constructor 主 fail-closed 路径成立。关门可不阻塞，但建议在后续小修合入。

### F-002 · Playwright M3 浏览器运行仍受本机端口排除阻塞

- **严重度**：low
- **建议**：recommended
- **状态**：open（环境 residual；**非**代码缺陷主张）
- **描述**：`localization.spec.ts` M3 已改写，但本机 8080 落在 Windows 排除区间 8011–8110，Go API 无法绑定，浏览器 e2e 未在本会话复跑。C3 成功标准已接受「单元覆盖 + 降级留痕」；S5 同机历史曾跑通。
- **证据**：`attachments/s6-e2e-env-block.log`；本会话 `netsh` 仍显示 8011–8110；E-001 / A-001 同口径。
- **建议**：端口区间解除或改宿主/`PLAYWRIGHT` base 后补跑 admin M3 一次，附件补日志即可；**不必**阻塞 C4 用户确认（与已冻结 C3 口径一致）。

## 必改项汇总

| 级别 | 列表 |
|------|------|
| **required** | **无** |
| **recommended** | F-001（idle 首帧）；F-002（e2e 补跑） |

## 与既有意见的异同

| 条目 | 关系 |
|------|------|
| A-001 self · C1–C3 pass | **同意**实现完成与主证据链；本条独立复跑 vitest 727 / build / go test 再确认 |
| A-001 未列 idle 竞态 | 本条 **新增** F-001 recommended |
| A-001 e2e 降级 | **同意**并复核端口排除仍在；记为 F-002 recommended，不升级 required |

## 结论 + 建议给编排器 / 用户

1. **实现门禁（C1–C3）**：可视为达成；无未闭合 required finding。
2. **C4 剩余**：用户书面确认关门 → `/govern` 将 GOAL-007 `done`、C4 检查点勾选、Root `GOAL-001` 恢复 `done`（`7/7`）并同步 goal-tree / workspace.md（D-003 §4）；**勿**在未确认时静默改状态。
3. **可选跟进**：F-001 小修 + 单测；F-002 环境恢复后 e2e 补跑（recommended，不阻断书面确认）。
4. 响应入口：**`/govern`**（汇总 A-002，关闭/接受 recommended，驱动 C4 用户确认）。

## 声明

本意见 `source: independent`；**不修改**目标 `status` / 检查点 / 派生 `progress` / goal-tree。响应、finding 闭合与关门状态变更由 **`/govern`** 与用户书面裁决处理。

---

## 响应（对独立意见 A-002 · `/govern`）

| date | actor | summary |
|------|-------|---------|
| 2026-08-09 | `/govern` | 采纳 verdict `pass`。**F-001 → `fixed`**：`useRecordSourcePrefill` 首帧 `recordSource` 存在时进入 `loading`（skeleton），消除空可编辑表单 idle 竞态；新增回归单测；vitest **728/728** + `npm run build` exit 0；证据见 E-002 / commit `ac757c5`。**F-002 → `accepted-residual`**：环境 residual（本机 8080 落入端口排除区间 8011–8110），C3 成功标准已接受单元覆盖 + 降级留痕；范围 = 仅浏览器补跑证据，复审触发 = 区间解除/换宿主后补跑 admin M3 附日志；不阻塞 C4。**C4 剩余**：用户书面确认 → `/govern` 将 GOAL-007 `done`（C4 勾选，`progress 4/4`）并恢复 Root `GOAL-001` `done`（`7/7`），同步 goal-tree / workspace.md（D-003 §4）。本响应不改原 verdict 与 finding 原文。 |
