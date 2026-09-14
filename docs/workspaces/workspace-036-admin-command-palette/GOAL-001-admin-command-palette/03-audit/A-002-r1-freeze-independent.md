---
doc_type: goal-audit
id: A-002-r1-freeze-independent
status: recorded
source: independent
auditor: grok-build (grok-4.6 · reasoning high)
date: 2026-09-14
scope: R1 范围与信息冻结 / SearchableItem 分母、四 Profile 算术、权限动作边界、provider v1、排序去重上限与可访问性口径
verdict: conditional
created: 2026-09-14
updated: 2026-09-14
parent: GOAL-001-admin-command-palette
version: 0.1.0
---

# A-002 · R1 范围与信息冻结独立审计（2026-09-14）

- **source**：independent
- **auditor**：grok-build (grok-4.6 · reasoning high)
- **类型**：design-plan
- **scope**：R1 范围与信息冻结（I-036-001～003；分母 / Profile 覆盖 / 权限与动作执行边界 / provider v1 / 排序·去重·上限 / 快捷键与 ARIA 口径）
- **verdict**：conditional

## 范围与区间

工作区：`[workspace-036-admin-command-palette]`，canonical `docs/workspaces/workspace-036-admin-command-palette/`，Root `GOAL-001-admin-command-palette`。共享资料目录 `none`，本轮无资料引用。

本意见只审已落盘的 R1 决策与矩阵是否可在进入 R2 前作为冻结基线；不把工作区未跟踪的 `searchable.ts` / `CommandPalette.tsx` 或 `App.tsx` 未完成接线写成 R2/R3 已完成。R1 不审 Palette 浏览器可用性。未读取其他工作区上下文。

## 成果（有证据）

- 工作区绑定成立：`workspace.md` 的 `id` / `root_goal` / `canonical_scope` / `plan_refs` / `primary_plan` 与 Root `00-meta.md`、VP-036 `active` v0.2.0、Charter `schema-ui-core-admin-foundation@0.4.0` 一致；`vision_role: delivery`；Vision open required = 0。
- D-002 与矩阵对实体搜索、Saved Views/最近项、`RT-X01`/`RT-X02`、不扩展 pinned AppManifest 的排除与用户书面确认一致。`AppManifest` 仅有 `protocolVersion` / `requiredCapabilities` / `app` / `pages` / `navigation`，未知字段走 `UNKNOWN_MANIFEST_FIELD`。
- 四 Profile 的模块数、页面注册数、持久化导航贡献、发布 Manifest 叶子、Shell 目标数、Schema 顶层 action 定义数均与仓库一致（见下表）。
- 页面/导航纳入规则与代码一致：全局目标来自 `projectNavigation` 可见 top/sidebar/user 叶子；`menu_notifications` 无公开 Manifest leaf，通知入口由 Shell `NotificationBell` 补 `/notifications`；inner/参数页（`users-invites`、`dictionary-entries`、`task-runs`、`wallet-entries`、`telegram-operator`）不进页面目标分母。
- 权限/执行边界的**政策**冻结可核对：收录与执行须复核 `visibleWhen` / `permissionIntent` / 局部 `permissions` / cascade；禁止直接打原始 URL 或 actionRef；后端鉴权仍是最终边界。现有 `invokeAction` 对 modal/navigate 无前置 gate、`runRequest` 对 custom 在 gate 之前返回，与 D-002 / A-001 F-001 描述相符。
- provider v1 / UX 口径本身可执行：稳定 id、kind、localized label、keywords、group、route/action context；冲突 id 全排除；查询不持久化；大小写/重音不敏感；精确 → 前缀 → 包含；上限 12；`Ctrl+K`/`Meta+K`；dialog + combobox/listbox；方向键 / Home / End / Enter / Escape / 外部点击 / Tab trap / 焦点恢复。现有代码中无已提交的冲突快捷键处理器。

## 对照成功标准 / 信息门禁

| 标准或信息项 | 状态 | 证据 |
|---|---|---|
| 工作区绑定与 VP-036 delivery | 达成 | `workspace.md`；VP-036 `lead_workspace`；Charter `@0.4.0` |
| I-036-001 精确分母 | 部分 | 模块/页面/导航/叶子/action 定义列正确；**可收录触发器列与 §1.3 规则差 1**（F-001） |
| I-036-002 权限/URL/动作语义 | 政策已冻结；实现缺口仍开放 | D-002；`render.tsx` invokeAction；`permissions.ts`；A-001 F-001 |
| I-036-003 字段/排序/快捷键/ARIA | 口径已冻结；浏览器证据仍待 R3/R4 | D-002；矩阵 §3 |
| I-036-004 / I-036-006 | verified，不影响本 scope | `00-meta.md`；VRev-092 |
| I-036-005 | deferred non-blocking，最晚 R3 | `00-meta.md`；D-002 |
| V-F123 recommended | 部分承接 | 有可机器核对计数表，但触发器列不是按可见页规则算出的 oracle |

独立复核计数（对照 `apps/api/kernel/profile.go`、`apps/api/configs/config.yaml`、各模块 `manifest/fragment.json` 与 `schema/*.json`）：

| profile | 模块 | 页面 | 持久化导航 | Manifest 叶子 | Shell 目标 | action 定义 | 矩阵「可收录触发器」 | 按 §1.3 可见/可解析页触发器 |
|---|---:|---:|---:|---:|---:|---:|---:|---:|
| mvp | 11 | 6 | 5 | 4 | 5 | 32 | 6 | **6** |
| admin | 23 | 22 | 18 | 17 | 18 | 93 | 17 | **16** |
| demo | 12 | 14 | 5 | 12 | 13 | 34 | 6 | **6** |
| custom（admin + telegram + digital-offer） | 25 | 27 | 22 | 21 | 22 | 101 | 19 | **18** |

admin/custom 多出的 1 条均为 `dictionary-entries.json` 的 toolbar `openCreate`（「New entry」）。该页路由为 `/dictionary-entries/{dictKey}`，不在发布导航叶子中，且登记于 `NAVIGATION_PAGE_PARENTS`。按矩阵 §1.1/§1.3，它不是可见、可解析的全局目标页，其 toolbar 也不应进入全局命令分母。行级 `table.props.actions[]`、`admin-list-batch` 的 `requiresSelection` toolbar 均已正确排除。

按 §1.3，admin 可见页直接触发器为 16 条：users×4、roles×2、settings `resetSettings`、file-library `openUpload`、data-dictionary `openCreate`、scheduled-tasks `openCreate`、recycle-bin `purgeAll`、data-permission `openRegister`、my-wallet `openRedeem`、wallet-vouchers `openGenerate`、wallet `openCreate`/`runReconcile`。custom 另加 telegram-settings `openTelegramOperator` 与 digitaloffer-offers `openCreate`，共 18 条。

## Findings

### F-001 · 可收录触发器分母与 §1.3 规则不一致（required · med · open）

- **描述**：矩阵把 admin=17、custom=19 写成可收录页面级触发器，但该计数包含 inner/参数页 `dictionary-entries` 的 toolbar。§1.3 明确只从「上述可见、可解析页面」的 `toolbar[]` / `actionButton` 收录；`openCreate` 依赖 `dictKey` 路由上下文，不能作为无行/无参数的全局命令。I-036-001 已被 `00-meta.md` 标为 `verified`，但精确分母在触发器列上仍自相矛盾。R2 若把 17/19 当作回归 oracle，会把不可全局执行的命令编入分母。
- **证据**：`attachments/r1-searchable-item-matrix.md` §2；`apps/api/modules/datadictionary/schema/dictionary-entries.json` toolbar `openCreate`；`apps/api/modules/datadictionary/manifest/fragment.json` route `/dictionary-entries/{dictKey}`；`apps/web/src/app/navigation.ts` `NAVIGATION_PAGE_PARENTS`；矩阵 §1.1/§1.3。
- **关联**：I-036-001（最晚阶段 R1，影响 R1 冻结 / R2 聚合 / R4 回归）。
- **建议闭合**：修正矩阵触发器列为 admin **16** / custom **18**（或改为显式枚举可见页触发器并单独列出 excluded inner toolbar）；同步 I-036-001 证据。不得在未修正时把 17/19 写入 R2 测试期望。
- **状态**：open

### F-002 · 分母附件是 Profile 计数而非 item→module→profile 清单（recommended · med · open）

- **描述**：I-036-001 的验证动作是「建立 item→module→profile 矩阵」；V-F123 要求可机器核对的分母。当前附件只有 Profile 合计。计数在模块/页面/导航列上正确，但缺少条目清单使 F-001 的 1 条误差无法被测试直接钉死。
- **证据**：`00-meta.md` I-036-001；`attachments/r1-searchable-item-matrix.md` §2；VRev-090/VRev-092 V-F123。
- **建议闭合**：在修正 F-001 时追加可见页/动作的稳定 id 清单（至少四 Profile 的 pageRef/action 键），作为 R2/R4 oracle。
- **状态**：open

### F-003 · 程序化动作 gate 缺口独立复核成立（recommended · med · open）

- **描述**：与 A-001 F-001 同范围，独立复核确认：`invokeAction` 对 `modal`/`navigate` 直接执行、无 `effectivePermission` 前置；`runRequest` 对 `custom` 在 permission gate 之前返回；`ActionButtonView` 的 UI 目标 id 回退 `node.id`，默认 `onAction` 只把 `node.props` 交给 executor。当前 settings 的 `resetSettings` 有 `props.key=reset`，telegram 的 `openTelegramOperator` 无 key 且为 navigate。Palette 若直接调用 `invokeAction`，会绕过 `schema-table.tsx` 在 UI 层做的 toolbar gate。这是 R2/R3 实现约束，不单独阻断 R1 口径冻结，但在 R3 动作调用前必须 `fixed`。
- **证据**：`apps/web/src/renderer/render.tsx` `invokeAction` / `runRequest` / `ActionButtonView` / default `onAction`；`apps/web/src/renderer/schema-table.tsx` toolbar `effectivePermission`；`apps/web/src/renderer/permissions.ts` actionButton targetId。
- **状态**：open（与 A-001 F-001 并行；本条不升为 R1 required）

### F-004 · 工作区投影与未审查实现草稿不同步（recommended · low · open）

- **描述**：`00-meta.md` 已将 I-036-001～003 标为 verified，且 R1 检查点仍未勾选（等待独立合并）——这一点符合门禁。但 `workspace.md` 纲领表与 `goal-tree.md` 备注仍写「R1 not started / I-036-001～003 collecting」。同时工作区存在未跟踪的 `apps/web/src/app/searchable.ts`、`CommandPalette.tsx` 以及 `App.tsx` 中未使用的 provider/palette import。这些草稿**不能**当作 R2/R3 完成证据（无 i18n key、App 未渲染 Palette、无本轮测试），但 E-002「R2 尚未发生」作为工作区快照已不精确。编排器在响应本意见后应同步投影，并把草稿当作未验收 WIP，而不是空白起步或已完成实现。
- **证据**：`workspace.md` 纲领表；`goal-tree.md` 状态表备注；`git status` 中上述路径；`App.tsx` 仅 import / 未使用的 `searchableProviders` prop。
- **状态**：open

## 必改项汇总

1. **F-001（required）**：按矩阵 §1.3 修正 admin/custom 可收录触发器分母（16 / 18），或给出与规则一致的枚举；在此之前不得把 17/19 用作 R2 oracle，也不得把 I-036-001 视为已精确关闭。

## 与既有意见的异同

- A-001 self `pass`，open required = 0；其 F-001（程序化 gate）独立复核同意，保持 recommended。
- 本意见新增 **F-001 required**：self 未核对「可收录触发器」列与 inner/参数页排除规则的差 1。因此独立 verdict 为 **conditional**，与 self `pass` 在 R1 放行结论上不一致。冲突点是 I-036-001 是否已经精确关闭，不是动作 gate 政策本身。
- 无用户书面 residual / overruled。按 P-004，编排器须展示该冲突并给出建议，等用户确认修正矩阵或接受残余后再勾选 R1 / 进入以该列为 oracle 的 R2。

## 结论 + 建议给编排器/用户的下一步

**verdict: conditional。** 工作区绑定、四 Profile 的模块/页面/导航/叶子/Shell/action 定义算术、页面纳入排除、权限政策、provider v1 与 UX 口径总体可核对；不能无条件放行 R1，因为触发器分母与自己的纳入规则矛盾。

建议 `/govern`：

1. 先响应本条 F-001（修正矩阵 ± 枚举清单），再决定是否保持 I-036-001 `verified`。
2. 不要把工作区未跟踪的 `searchable.ts` / `CommandPalette.tsx` 当成已验收 R2/R3。
3. R2 测试的触发器期望值用 16/18（外加 mvp/demo=6），不要用 17/19。
4. A-001/A-002 的程序化 gate 在 R3 动作调用前必须留下 `fixed` 证据。

## 声明

本意见 `source: independent`，不修改 `status` / 检查点 / 派生 `progress` / 方案正文 / goal-tree。响应、finding 闭合与是否勾选 R1 由 `/govern` 处理。
